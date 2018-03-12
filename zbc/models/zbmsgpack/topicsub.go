package zbmsgpack

import (
	"encoding/json"
	"fmt"
)

// TopicSubscriptionInfo is used to open a topic subscription.
type TopicSubscriptionInfo struct {
	StartPosition    int64  `msgpack:"startPosition"`
	PrefetchCapacity int32  `msgpack:"prefetchCapacity"`
	Name             string `msgpack:"name"`

	ForceStart bool   `msgpack:"forceStart"`
	State      string `msgpack:"state"`

	SubscriberKey uint64 `msgpack:"-"`
}

func (t *TopicSubscriptionInfo) String() string {
	b, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return fmt.Sprintf("json marshaling failed\n")
	}
	return fmt.Sprintf("%+v", string(b))
}

// TopicSubscriptionAckRequest is used to acknowledge receiving of an event.
type TopicSubscriptionAckRequest struct {
	Name        string `msgpack:"name"`
	AckPosition uint64 `msgpack:"ackPosition"`
	State       string `msgpack:"state"`
}

func (t *TopicSubscriptionAckRequest) String() string {
	b, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return fmt.Sprintf("json marshaling failed\n")
	}
	return fmt.Sprintf("%+v", string(b))
}

type TopicSubscriptionCloseRequest struct {
	TopicName     string `msgpack:"topicName"`
	PartitionID   uint16 `msgpack:"partitionId"`
	SubscriberKey uint64 `msgpack:"subscriberKey"`

	SubscriptionName string `msgpack:"-"`
}

func (t *TopicSubscriptionCloseRequest) String() string {
	b, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return fmt.Sprintf("json marshaling failed\n")
	}
	return fmt.Sprintf("%+v", string(b))
}
