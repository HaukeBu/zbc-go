package zbtopology

import (
	"github.com/zeebe-io/zbc-go/zbc/common"
	"github.com/zeebe-io/zbc-go/zbc/models/zbmsgpack"
	"sync"
	"time"
)

// RoundRobinCtrl is used for implementation of round robin algorithm on TopologySvc.
type RoundRobinCtrl struct {
	cluster          *zbmsgpack.ClusterTopology
	lastIndexByTopic *sync.Map
}

func (rr *RoundRobinCtrl) UpdateClusterTopology(c *zbmsgpack.ClusterTopology) {
	rr.cluster = c
	rr.cluster.UpdatedAt = time.Now()
}

// PeekPartitionIndex is used to check current status of round robin.
func (rr *RoundRobinCtrl) PeekPartitionIndex(topic string) *uint16 {
	lastPartitionUsedIndex, ok := rr.lastIndexByTopic.Load(topic)
	if !ok {
		if index := rr.cluster.GetRandomPartitionIndex(topic); index != nil {
			rr.lastIndexByTopic.Store(topic, index)
			return index
		}

	}
	return lastPartitionUsedIndex.(*uint16)
}

func (rr *RoundRobinCtrl) incrementLastPartitionIndex(topic string) {
	if index, ok := rr.lastIndexByTopic.Load(topic); ok {
		*index.(*uint16)++
		rr.lastIndexByTopic.Store(topic, index)
	}
}

func (rr *RoundRobinCtrl) setLastPartitionIndex(topic string, val uint16) {
	rr.lastIndexByTopic.Store(topic, &val)
}

func (rr *RoundRobinCtrl) Partitions(topic string) ([]uint16, error) {
	brokers, ok := rr.cluster.PartitionIDByTopicName[topic]

	if !ok {
		return nil, zbcommon.ErrTopicPartitionNotFound
	}
	return brokers, nil
}

func (rr *RoundRobinCtrl) nextPartitionID(topic string) (*uint16, error) {
	partitions, err := rr.Partitions(topic)
	if err != nil {
		return nil, err
	}
	partitionsLength := uint16(len(partitions) - 1)

	lastIndex := rr.PeekPartitionIndex(topic)
	if lastIndex != nil {
		partitionID := partitions[*lastIndex]

		if *lastIndex == partitionsLength {
			rr.setLastPartitionIndex(topic, 0)
		} else {
			rr.incrementLastPartitionIndex(topic)
		}

		// TODO: zbc-go/issues#40 + zbc-go/issues#48
		return &partitionID, nil
	}
	return nil, zbcommon.ErrClusterPartialInformation
}

func NewRoundRobinCtl(cluster *zbmsgpack.ClusterTopology) *RoundRobinCtrl {
	return &RoundRobinCtrl{cluster, &sync.Map{}}
}
