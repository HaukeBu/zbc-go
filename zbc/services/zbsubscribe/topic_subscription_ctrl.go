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

type TopicSubscriptionCtrl struct {
	CloseRequests     map[uint16]*zbmsgpack.TopicSubscriptionCloseRequest
	SubscriptionsInfo map[uint16]*zbmsgpack.TopicSubscriptionInfo
}

func (ts *TopicSubscriptionCtrl) AddSubscriptionInfo(partition uint16, request *zbmsgpack.TopicSubscriptionInfo) {
	ts.SubscriptionsInfo[partition] = request
}

func (tc *TopicSubscriptionCtrl) AddCloseRequest(partition uint16, sub *zbmsgpack.TopicSubscriptionCloseRequest) {
	tc.CloseRequests[partition] = sub
}

func (tc *TopicSubscriptionCtrl) GetSubscriptionCloseRequest(partition uint16) (*zbmsgpack.TopicSubscriptionCloseRequest, bool) {
	value, ok := tc.CloseRequests[partition]
	return value, ok
}

func NewTopicSubscriptionAckCtrl() *TopicSubscriptionCtrl {
	return &TopicSubscriptionCtrl{
		CloseRequests:     make(map[uint16]*zbmsgpack.TopicSubscriptionCloseRequest),
		SubscriptionsInfo: make(map[uint16]*zbmsgpack.TopicSubscriptionInfo),
	}
}
