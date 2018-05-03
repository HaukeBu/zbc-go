package zbsubscribe

import (
	"github.com/zeebe-io/zbc-go/zbc/common"
	"github.com/zeebe-io/zbc-go/zbc/models/zbmsgpack"
	"github.com/zeebe-io/zbc-go/zbc/services/zbexchange"
)

// TaskSubscriptionSvc is responsible for test-task-subscriptions subscription management.
type TaskSubscriptionSvc struct {
	zbexchange.LikeExchangeSvc
}

func (ts *TaskSubscriptionSvc) taskConsumer(topic, lockOwner, taskType string, lockDuration uint64, credits int32) (*TaskSubscription, error) {
	partitions, err := ts.TopicPartitionsAddrs(topic)
	if err != nil {
		return nil, err
	}

	if partitions == nil || len(*partitions) == 0 {
		zbcommon.ZBL.Error().Str("component", "TaskSubscriptionSvc").Str("method", "taskConsumer").Msgf("topic %s topology not found", topic)
		return nil, zbcommon.ErrNoPartitionsFound
	}

	sz := uint64(len(*partitions)) * uint64(credits)
	taskSubscription := NewTaskSubscription(sz * sz)
	zbcommon.ZBL.Debug().Str("component", "TaskSubscriptionSvc").Str("method", "taskConsumer").Msg("new task subscription created")

	for partitionID := range *partitions {
		sub := ts.OpenTaskPartition(taskSubscription.OutCh, partitionID, lockOwner, taskType, lockDuration, credits)

		if sub != nil {
			taskSubscription.AddSubscription(partitionID, sub)
			taskSubscription.AddPartition(partitionID)
		} else {
			zbcommon.ZBL.Error().Str("component", "TaskSubscriptionSvc").Str("method", "taskConsumer").Msg("opening subscription on all partitions failed")
			taskSubscription.Close()
			return nil, zbcommon.ErrSubscriptionPipelineFailed
		}
	}

	return taskSubscription, nil
}

// TaskSubscription will open a test-task-subscriptions subscription with specified handler/callback function.
func (ts *TaskSubscriptionSvc) TaskSubscription(topic, lockOwner, taskType string, lockDuration uint64, credits int32, cb TaskSubscriptionCallback) (*TaskSubscription, error) {
	var err error
	var taskSubscription *TaskSubscription
	var count int = 0

	for {
		taskSubscription, err = ts.taskConsumer(topic, lockOwner, taskType, lockDuration, credits)
		if taskSubscription != nil && err == nil {
			break
		}

		zbcommon.ZBL.Error().Str("component", "TaskSubscriptionSvc").Str("method", "TaskSubscription").Msgf("failed to create task consumer: %+v", err)
		ts.GetTopology()

		if count < 3 {
			count++
			continue
		} else {
			zbcommon.ZBL.Error().Str("component", "TaskSubscriptionSvc").Str("method", "TaskSubscription").Msg("failed to create task consumer after 3 tries. bailing")
			return nil, err
		}
	}

	zbcommon.ZBL.Debug().Str("component", "TaskSubscriptionSvc").Msgf("TaskSubscription created %s", taskType)
	taskSubscription = taskSubscription.WithCallback(cb).WithTaskSubscriptionSvc(ts)
	taskSubscription.initCredits(credits)

	if taskSubscription == nil {
		return nil, zbcommon.ErrFailedToOpenTaskSubscription
	}
	return taskSubscription, err
}

func (tp *TaskSubscriptionSvc) CloseTaskSubscription(sub *TaskSubscription) []error {
	var errs []error

	sub.Subscriptions.Range(func(key, value interface{}) bool {
		taskSub := value.(*zbmsgpack.TaskSubscriptionInfo)
		_, err := tp.CloseTaskSubscriptionPartition(key.(uint16), taskSub)
		if err != nil {
			errs = append(errs, err)
		}
		return true
	})

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
