package zbsubscribe


import (
	"github.com/zeebe-io/zbc-go/zbc/models/zbmsgpack"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsubscriptions"
)

// TopicSubscriptionCtrl is controller object for subscription behaviour.
type TopicSubscriptionCallbackCtrl struct {
	callback TopicSubscriptionCallback
}

// ExecuteCallback will execute attached callback.
func (ts *TopicSubscriptionCallbackCtrl) ExecuteCallback(event *zbsubscriptions.SubscriptionEvent) error {
	return ts.callback(clientInstance, event)
}

// NewTopicSubscriptionCtrl is a constructor for TopicSubscriptionCtrl object.
func NewTopicSubscriptionCallbackCtrl(cb TopicSubscriptionCallback) *TopicSubscriptionCallbackCtrl {
	return &TopicSubscriptionCallbackCtrl{
		callback: cb,
	}
}

type TopicSubscriptionAckCtrl struct {
	messageCount  uint64
	Subscriptions map[uint16]*zbmsgpack.TopicSubscription
}

func (tc *TopicSubscriptionAckCtrl) PeekMessageCount() uint64 {
	return tc.messageCount
}

func (tc *TopicSubscriptionAckCtrl) AddSubscription(partition uint16, sub *zbmsgpack.TopicSubscription) {
	tc.Subscriptions[partition] = sub
}

func (tc *TopicSubscriptionAckCtrl) GetSubscription(partition uint16) (*zbmsgpack.TopicSubscription, bool) {
	value, ok := tc.Subscriptions[partition]
	return value, ok
}

func (tc *TopicSubscriptionAckCtrl) IncrementMessageCount() {
	tc.messageCount++
}

func (tc *TopicSubscriptionAckCtrl) GetMessageCount() uint64 {
	return tc.messageCount
}

func (tc *TopicSubscriptionAckCtrl) NormalizeMessageCount() {
	tc.messageCount /= 3
}

func NewTopicSubscriptionAckCtrl() *TopicSubscriptionAckCtrl {
	return &TopicSubscriptionAckCtrl{
		messageCount:  0,
		Subscriptions: make(map[uint16]*zbmsgpack.TopicSubscription),
	}
}