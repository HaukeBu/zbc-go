package zbsubscribe

import (
	"errors"

	"github.com/zeebe-io/zbc-go/zbc/common"
	"github.com/zeebe-io/zbc-go/zbc/models/zbmsgpack"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsubscriptions"
	"sync/atomic"
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

// TaskSubscriptionCreditsCtrl is controller structure for credits management on test-task-subscriptions subscription.
type TaskSubscriptionCreditsCtrl struct {
	credits       map[uint16]int32
	Subscriptions map[uint16]*zbmsgpack.TaskSubscriptionInfo
}

// PeekCredits will return current state of credits per partition.
func (ts *TaskSubscriptionCreditsCtrl) PeekCredits() map[uint16]int32 {
	cp := ts.credits
	return cp
}

// ReduceCredit will decrement credits of given partition.
func (ts *TaskSubscriptionCreditsCtrl) ReduceCredit(partitionID uint16) (*zbmsgpack.TaskSubscriptionInfo, error) {
	var taskSub *zbmsgpack.TaskSubscriptionInfo

	if value, ok := ts.credits[partitionID]; ok {
		atomic.AddInt32(&value, -1)
		ts.credits[partitionID] = value
		if value <= 0 {
			taskSub = ts.Subscriptions[partitionID]
		}
	} else {
		return nil, errors.New("cannot find subscription")
	}

	return taskSub, nil
}

// ResetCredits will set credits to initial value.
func (ts *TaskSubscriptionCreditsCtrl) ResetCredits(partitionID uint16) {
	if sub, ok := ts.Subscriptions[partitionID]; ok {
		ts.credits[partitionID] = sub.Credits
	}
}

// AddSubscription will add test-task-subscriptions subscription information to its belonging partitionID.
func (ts *TaskSubscriptionCreditsCtrl) AddSubscription(partition uint16, sub *zbmsgpack.TaskSubscriptionInfo) {
	ts.Subscriptions[partition] = sub
	ts.credits[partition] = sub.Credits
}

// GetSubscription will return test-task-subscriptions subscription information.
func (ts *TaskSubscriptionCreditsCtrl) GetSubscription(partition uint16) (*zbmsgpack.TaskSubscriptionInfo, bool) {
	value, ok := ts.Subscriptions[partition]
	return value, ok
}

// NewTaskSubscriptionCreditsCtrl  is a constructor for TaskSubscriptionCreditsCtrl.
func NewTaskSubscriptionCreditsCtrl() *TaskSubscriptionCreditsCtrl {
	return &TaskSubscriptionCreditsCtrl{
		credits:       make(map[uint16]int32),
		Subscriptions: make(map[uint16]*zbmsgpack.TaskSubscriptionInfo),
	}
}

// taskSubscriptionCloseCtrl is responsible for controlling of closing of the channel.
type taskSubscriptionCloseCtrl struct {
	closeCh chan interface{}
}

func (ts *taskSubscriptionCloseCtrl) getCloseMessage() chan interface{} {
	return ts.closeCh
}

func (ts *taskSubscriptionCloseCtrl) sendCloseMessage() {
	ts.closeCh <- struct{}{}
}

// NewTaskSubscriptionCloseCtrl is constructor for taskSubscriptionCloseCtrl object.
func newTaskSubscriptionCloseCtrl() *taskSubscriptionCloseCtrl {
	return &taskSubscriptionCloseCtrl{
		closeCh: make(chan interface{}),
	}
}
