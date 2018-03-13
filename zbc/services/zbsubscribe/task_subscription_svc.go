package zbsubscribe

import (
	"github.com/zeebe-io/zbc-go/zbc/common"
	"github.com/zeebe-io/zbc-go/zbc/services/zbexchange"
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

	if partitions == nil || len(*partitions) == 0 {
		zbcommon.ZBL.Error().Msgf("topic %s topology not found", topic)
		return nil, zbcommon.BrokerNotFound
	}

	taskSubscription := NewTaskSubscription()
	zbcommon.ZBL.Debug().Msg("new task subscription created")

	for partitionID := range *partitions {
		sub := ts.OpenTaskPartition(taskSubscription.OutCh, partitionID, lockOwner, taskType, credits)

		if sub != nil {
			taskSubscription.AddSubscription(partitionID, sub)
			taskSubscription.AddPartition(partitionID)
		} else {
			zbcommon.ZBL.Error().Msg("opening subscription on all partitions failed")
			taskSubscription.Close()
			return nil, zbcommon.ErrSubscriptionPipelineFailed
		}
	}

	return taskSubscription, nil
}

// TaskSubscription will open a test-task-subscriptions subscription with specified handler/callback function.
func (ts *TaskSubscriptionSvc) TaskSubscription(topic, lockOwner, taskType string, credits int32, cb TaskSubscriptionCallback) (*TaskSubscription, error) {
	var err error
	var taskSubscription *TaskSubscription
	var count int = 0

	for {
		taskSubscription, err = ts.taskConsumer(topic, lockOwner, taskType, credits)
		if taskSubscription != nil && err == nil {
			break
		}
		ts.RefreshTopology()

		if count < 3 {
			count++
			continue
		} else {
			zbcommon.ZBL.Error().Msg("Failed to create TaskSubscription")
			return nil, err
		}
	}

	zbcommon.ZBL.Debug().Msg("TaskSubscription created")

	taskSubscription = taskSubscription.WithCallback(cb).WithTaskSubscriptionSvc(ts)
	if taskSubscription == nil {
		return nil, zbcommon.BrokerNotFound
	}
	return taskSubscription, err
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
