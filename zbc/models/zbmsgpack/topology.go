package zbmsgpack

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

// BrokerPartition contains information about partition contained in topology request.
type BrokerPartition struct {
	State             string `msgpack:"state"`
	TopicName         string `msgpack:"topicName"`
	PartitionID       uint16 `msgpack:"partitionId"`
	ReplicationFactor uint16 `msgpack:"replicationFactor"`
}

// Broker is used to hold broker contact information.
type Broker struct {
	Host       string            `msgpack:"host"`
	Port       uint64            `msgpack:"port"`
	Partitions []BrokerPartition `msgpack:"partitions"`
}

func (t *Broker) String() string {
	b, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return fmt.Sprintf("json marshaling failed\n")
	}
	return fmt.Sprintf("%+v", string(b))
}

func (t *Broker) Addr() string {
	return fmt.Sprintf("%s:%d", t.Host, t.Port)
}

// TopologyRequest is used to make a topology request.
type TopologyRequest struct{}

func (t *TopologyRequest) String() string {
	b, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return fmt.Sprintf("json marshaling failed\n")
	}
	return fmt.Sprintf("%+v", string(b))
}

// ClusterTopologyResponse is used to parse the topology response from the broker.
type ClusterTopologyResponse struct {
	Brokers []Broker `msgpack:"brokers"`
}

func (t *ClusterTopologyResponse) String() string {
	b, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return fmt.Sprintf("json marshaling failed\n")
	}
	return fmt.Sprintf("%+v", string(b))
}

// ClusterTopology is structure used by the client object to hold information about cluster.
type ClusterTopology struct {
	AddrByPartitionID      map[uint16]string
	PartitionIDByTopicName map[string][]uint16

	Brokers   []Broker
	UpdatedAt time.Time
}

func (t *ClusterTopology) FromPartitionCollection(collection *PartitionCollection) {
	mm := make(map[string][]uint16)
	for _, partition := range collection.Partitions {
		mm[partition.Topic] = append(mm[partition.Topic], uint16(partition.ID))
	}
	t.PartitionIDByTopicName = mm
}

func (t *ClusterTopology) GetRandomPartitionIndex(topic string) *uint16 {
	partitions, ok := t.PartitionIDByTopicName[topic]
	if !ok || len(partitions) <= 0 {
		return nil
	}

	rand.Seed(time.Now().Unix())
	index := uint16(rand.Intn(len(partitions)))
	return &index

}

// GetRandomBroker will select one broker on random using time.Now().UTC().UnixNano() as the seed for randomness.
func (t *ClusterTopology) GetRandomBroker() *Broker {
	rand.Seed(time.Now().UTC().UnixNano())
	if len(t.Brokers) <= 0 {
		return nil
	}
	index := rand.Intn(len(t.Brokers))
	broker := t.Brokers[index]
	return &broker
}

// String will marshal the data structure as series of characters.
func (t *ClusterTopology) String() string {
	b, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return fmt.Sprintf("json marshaling failed\n")
	}
	return fmt.Sprintf("%+v", string(b))
}
