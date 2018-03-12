package zbtopology

import (
	"net"

	"github.com/zeebe-io/zbc-go/zbc/common"
	"github.com/zeebe-io/zbc-go/zbc/models/zbdispatch"
	"github.com/zeebe-io/zbc-go/zbc/models/zbmsgpack"
	"github.com/zeebe-io/zbc-go/zbc/services/zbsocket"
	"github.com/zeebe-io/zbc-go/zbc/services/zbtransport"
	"sync"
)

type TopologySvc struct {
	*sync.RWMutex
	*zbdispatch.RequestFactory
	*zbdispatch.ResponseHandler

	transportSvc     zbtransport.LikeTransportSvc
	topologyWorkload chan *zbsocket.RequestWrapper

	lastIndexes   map[string]uint16
	bootstrapAddr string
	Cluster       *zbmsgpack.ClusterTopology
	RoundRobin    *RoundRobinCtrl
}

func (svc *TopologySvc) SetPartitionCollection(collection *zbmsgpack.PartitionCollection) {
	svc.Cluster.FromPartitionCollection(collection)
}

func (svc *TopologySvc) GetRoundRobinCtl() *RoundRobinCtrl {
	return svc.RoundRobin
}

func (svc *TopologySvc) TopicPartitionsAddrs(topic string) (*map[uint16]string, error) {
	if svc.Cluster == nil {
		return nil, zbcommon.BrokerNotFound
	}

	addrs := make(map[uint16]string)
	if partitions, ok := svc.Cluster.PartitionIDByTopicName[topic]; ok {
		for _, partitionID := range partitions {
			if addr, ok := svc.Cluster.AddrByPartitionID[partitionID]; ok {
				addrs[partitionID] = addr
			} else {
				return nil, zbcommon.ErrPartitionsNotFound
			}
		}
	}

	return &addrs, nil
}

func (svc *TopologySvc) PeekPartitionIndex(topic string) *uint16 {
	return svc.RoundRobin.PeekPartitionIndex(topic)
}

func (svc *TopologySvc) NextPartitionID(topic string) (*uint16, error) {
	if svc.RoundRobin == nil {
		return nil, zbcommon.BrokerNotFound
	}
	partitionID, err := svc.RoundRobin.nextPartitionID(topic)
	if err != nil {
		zbcommon.ZBL.Debug().Msg("fetching nextPartitionID failed")
		zbcommon.ZBL.Debug().Msg(" NextPartitionID :: refreshing topology")
		_, err := svc.RefreshTopology()
		newPartitionID, err := svc.RoundRobin.nextPartitionID(topic)
		if err != nil {
			return nil, err
		}
		return newPartitionID, nil
	}

	return partitionID, nil
}

func (svc *TopologySvc) getPartitions() (*zbmsgpack.PartitionCollection, error) {
	message := svc.CreatePartitionRequest()
	request := zbsocket.NewRequestWrapper(message)
	if svc.Cluster == nil || len(svc.Cluster.Brokers) == 0 {
		request.Addr = svc.bootstrapAddr
	} else {
		broker := svc.Cluster.GetRandomBroker()
		request.Addr = broker.Addr()
	}
	resp, err := svc.ExecuteRequest(request)
	if err != nil {
		return nil, err
	}

	partitionCollection := svc.UnmarshalPartition(resp)
	return partitionCollection, nil
}

func (svc *TopologySvc) getTopology() (*zbmsgpack.ClusterTopology, error) {
	request := zbsocket.NewRequestWrapper(svc.TopologyRequest())
	if svc.Cluster == nil || len(svc.Cluster.Brokers) == 0 {
		request.Addr = svc.bootstrapAddr
	} else {
		broker := svc.Cluster.GetRandomBroker()
		request.Addr = broker.Addr()
	}

	resp, err := svc.ExecuteRequest(request)
	if err != nil {
		return nil, err
	}

	topology := svc.UnmarshalTopology(resp)

	partitionCollection, err := svc.getPartitions()
	if err != nil {
		return nil, err
	}
	if len(partitionCollection.Partitions) == 0 {
		zbcommon.ZBL.Warn().Msg("PartitionCollection is empty. No topics are created.")
	}

	svc.Lock()
	defer svc.Unlock()

	svc.Cluster = topology
	svc.Cluster.FromPartitionCollection(partitionCollection)
	return topology, err
}

func (svc *TopologySvc) RefreshTopology() (*zbmsgpack.ClusterTopology, error) {
	zbcommon.ZBL.Debug().Msg("refreshing topology")

	topology, err := svc.getTopology()
	if err != nil {
		return nil, err
	}

	if svc.RoundRobin == nil {
		zbcommon.ZBL.Debug().Msg("round robing controller is nil ... initializing")
		svc.Lock()
		svc.RoundRobin = NewRoundRobinCtl(svc.Cluster)
		svc.Unlock()

		//go func() {
		//	zbcommon.ZBL.Debug().Msg("starting topology ticker")
		//	for {
		//		select {
		//		case <-time.After(zbcommon.TopologyRefreshInterval * time.Second * 30):
		//			svc.Lock()
		//			lastUpdate := svc.Cluster.UpdatedAt
		//			svc.Unlock()
		//
		//			if time.Since(lastUpdate) > zbcommon.TopologyRefreshInterval*time.Second {
		//				zbcommon.ZBL.Debug().Msg("topology ticker :: refreshing topology")
		//				_, err := svc.RefreshTopology()
		//				if err != nil {
		//					// TODO: do something with error here
		//				}
		//			}
		//
		//		}
		//	}
		//}()
	} else {
		svc.Lock()
		svc.RoundRobin.UpdateClusterTopology(topology)
		svc.Unlock()
	}
	zbcommon.ZBL.Debug().Msg("topology refreshed")
	return svc.Cluster, nil
}

func (svc *TopologySvc) ExecuteRequest(request *zbsocket.RequestWrapper) (*zbdispatch.Message, error) {
	if len(request.Addr) == 0 {
		addr, err := svc.getDestinationAddr(request.Payload)

		if err == zbcommon.BrokerNotFound {
			return nil, zbcommon.BrokerNotFound
		}
		request.Addr = addr
	}

	svc.transportSvc.ExecTransport(request)

	select {
	case resp := <-request.ResponseCh:
		return resp, nil

	case err := <-request.ErrorCh:
		netErr, ok := err.(*net.OpError)
		if ok {
			// INFO: exception is cause of network failure
			return nil, netErr
		} else if svc.Cluster != nil {
			// INFO: error is caused cause of broker
			zbcommon.ZBL.Debug().Msg("topology_svc::error: cluster is not initialized")
			topology, err := svc.RefreshTopology()
			if err != nil {
				return nil, err
			}
			if len(topology.PartitionIDByTopicName) == 0 {
				return nil, zbcommon.ErrTopicPartitionNotFound
			}
		}
		return nil, err

	}
}

func (svc *TopologySvc) getDestinationAddr(msg *zbdispatch.Message) (string, error) {
	svc.Lock()
	defer svc.Unlock()

	if svc.Cluster == nil {
		return "", zbcommon.BrokerNotFound
	}


	if msg.IsTopologyMessage() && svc.Cluster == nil {
		return svc.bootstrapAddr, nil
	}

	partitionID := msg.ForPartitionId()

	if partitionID != nil {
		if addr, ok := svc.Cluster.AddrByPartitionID[*partitionID]; ok {
			return addr, nil
		}
	}

	return "", zbcommon.BrokerNotFound
}

func NewTopologySvc(bootstrapAddr string) *TopologySvc {
	// TODO: validate bootstrapAddr as valid URI/IP

	return &TopologySvc{
		RWMutex: &sync.RWMutex{},
		RequestFactory:   zbdispatch.NewRequestFactory(),
		ResponseHandler:  zbdispatch.NewResponseHandler(),
		transportSvc:     zbtransport.NewTransportSvc(),
		topologyWorkload: make(chan *zbsocket.RequestWrapper, zbcommon.RequestQueueSize),
		lastIndexes:      make(map[string]uint16),
		bootstrapAddr:    bootstrapAddr,
		Cluster:          nil,
		RoundRobin:       nil,
	}
}
