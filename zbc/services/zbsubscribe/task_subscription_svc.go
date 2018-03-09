package zbsubscribe

import (
	"github.com/zeebe-io/zbc-go/zbc/services/zbexchange"
	"github.com/zeebe-io/zbc-go/zbc/common"
)

// TaskSubscriptionSvc is responsible for test-task-subscriptions subscription management.
type TaskSubscriptionSvc struct {
	zbexchange.LikeExchangeSvc
}

func (ts *TaskSubscriptionSvc) taskConsumer(topic, lockOwner, taskType string, credits int32) (*TaskSubscription, error) {
	partitions, err := ts.TopicPartitionsAddrs(topic)
	if err != nil {
		return nil, err
	}

	taskSubscription := NewTaskSubscription()

	for partitionID := range *partitions {
		sub, ch := ts.OpenTaskPartition(partitionID, lockOwner, taskType, credits)
		if sub != nil && ch != nil {
			taskSubscription.AddSubscription(partitionID, sub)
			taskSubscription.AddPartitionChannel(partitionID, ch)
		} else {
			// TODO: send close to all previous subs
			return nil, zbcommon.ErrSubscriptionPipelineFailed
		}
	}

	taskSubscription.StartPipeline()
	return taskSubscription, nil
}

// TaskSubscription will open a test-task-subscriptions subscription with specified handler/callback function.
func (ts *TaskSubscriptionSvc) TaskSubscription(topic, lockOwner, taskType string, cb TaskSubscriptionCallback) (*TaskSubscription, error) {
	taskConsumer, err := ts.taskConsumer(topic, lockOwner, taskType, 32)
	if err != nil {
		return nil, err
	}

	taskConsumer = taskConsumer.WithCallback(cb)
	taskConsumer.StartCallbackHandler(ts)

	return taskConsumer, err
}

func (tp *TaskSubscriptionSvc) CloseTaskSubscription(sub *TaskSubscription) []error {
	var errs []error
	for _, taskSub := range sub.Subscriptions {
		_, err := tp.CloseTaskSubscriptionPartition(taskSub)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func NewTaskSubscriptionSvc(exchange zbexchange.LikeExchangeSvc) *TaskSubscriptionSvc {
	return &TaskSubscriptionSvc{
		exchange,
	}
}
