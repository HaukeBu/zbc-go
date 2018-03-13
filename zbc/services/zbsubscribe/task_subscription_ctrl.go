package zbsubscribe

import (
	"github.com/zeebe-io/zbc-go/zbc/common"
	"github.com/zeebe-io/zbc-go/zbc/models/zbmsgpack"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsubscriptions"
)

// TaskSubscriptionCallbackCtrl is controller structure for test-task-subscriptions subscription callback handling.
type TaskSubscriptionCallbackCtrl struct {
	callback TaskSubscriptionCallback
}

// ExecuteCallback will execute attached callback.
func (ts *TaskSubscriptionCallbackCtrl) ExecuteCallback(event *zbsubscriptions.SubscriptionEvent) error {
	if ts.callback == nil {
		return zbcommon.ErrCallbackNotAttached
	}
	ts.callback(clientInstance.GetClientInstance(), event)
	return nil
}

// NewTaskSubscriptionCallbackCtrl is constructor for TaskSubscriptionCtrl object.
func NewTaskSubscriptionCallbackCtrl(cb TaskSubscriptionCallback) *TaskSubscriptionCallbackCtrl {
	return &TaskSubscriptionCallbackCtrl{
		callback: cb,
	}
}

// TaskSubscriptionCtrl is controller structure for credits management on test-task-subscriptions subscription.
type TaskSubscriptionCtrl struct {
	Subscriptions map[uint16]*zbmsgpack.TaskSubscriptionInfo
}

// AddSubscription will add test-task-subscriptions subscription information to its belonging partitionID.
func (ts *TaskSubscriptionCtrl) AddSubscription(partition uint16, sub *zbmsgpack.TaskSubscriptionInfo) {
	ts.Subscriptions[partition] = sub
}

// GetSubscription will return test-task-subscriptions subscription information.
func (ts *TaskSubscriptionCtrl) GetSubscription(partition uint16) (*zbmsgpack.TaskSubscriptionInfo, bool) {
	value, ok := ts.Subscriptions[partition]
	return value, ok
}

// NewTaskSubscriptionCreditsCtrl  is a constructor for TaskSubscriptionCtrl.
func NewTaskSubscriptionCreditsCtrl() *TaskSubscriptionCtrl {
	return &TaskSubscriptionCtrl{
		Subscriptions: make(map[uint16]*zbmsgpack.TaskSubscriptionInfo),
	}
}
