package zbsubscribe

import (
	"github.com/zeebe-io/zbc-go/zbc/models/zbsubscriptions"
)

// TaskSubscription is object for handling overall test-task-subscriptions subscription.
type TaskSubscription struct {
	*zbsubscriptions.SubscriptionPipelineCtrl

	*TaskSubscriptionCallbackCtrl
	*TaskSubscriptionCreditsCtrl

}

func (ts *TaskSubscription) StartCallbackHandler(svc *TaskSubscriptionSvc) {
	go func(ts *TaskSubscription) {
		for {
			msg := ts.GetNextMessage()
			if msg == nil {
				continue
			}

			PartitionID := msg.Event.PartitionId
			sub, err := ts.ReduceCredit(PartitionID)
			if sub != nil && err == nil {
				svc.IncreaseTaskSubscriptionCredits(sub)
				ts.ResetCredits(PartitionID)
			}
			ts.ExecuteCallback(msg)
		}
	}(ts)
}

func (ts *TaskSubscription) WithCallback(cb TaskSubscriptionCallback) *TaskSubscription {
	ts.TaskSubscriptionCallbackCtrl = NewTaskSubscriptionCallbackCtrl(cb)
	return ts
}

// NewTaskSubscription is constructor for TaskSubscription object.
func NewTaskSubscription() *TaskSubscription {
	return &TaskSubscription{
		SubscriptionPipelineCtrl:     zbsubscriptions.NewSubscriptionPipelineCtrl(),
		TaskSubscriptionCallbackCtrl: nil,
		TaskSubscriptionCreditsCtrl:  NewTaskSubscriptionCreditsCtrl(),
	}
}
