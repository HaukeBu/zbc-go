package zbsubscribe

import (
	"github.com/zeebe-io/zbc-go/zbc/services/zbexchange"
)

// LikeTaskSubscriptionSvc defines behaviour of TaskSubscriptionSvc.
type LikeTaskSubscriptionSvc interface {
	TaskSubscription(topic, lockOwner, taskType string, credits int32, cb TaskSubscriptionCallback) (*TaskSubscription, error)
	CloseTaskSubscription(task *TaskSubscription) []error

	// TODO: Move CompleteTask behind LikeTaskSubscriptionSvc interface
	// CompleteTask(test-task-subscriptions *SubscriptionEvent) (*zbmsgpack.Task, error)
	// TODO: Move FailTask behind LikeTaskSubscriptionSvc interface
	// FailTask(test-task-subscriptions *zbcommon.SubscriptionEvent) (*zbmsgpack.Task, error)
}


// LikeTopicSubscriptionSvc defines behaviour of TopicSubscriptionSvc.
type LikeTopicSubscriptionSvc interface {
	TopicSubscription(topic, subName string, startPosition int64, cb TopicSubscriptionCallback) (*TopicSubscription, error)
	CloseTopicSubscription(sub *TopicSubscription) []error

	// TODO: Move TopicSubscriptionAckRequest behind this interface
	// TopicSubscriptionAckRequest(ts *zbmsgpack.TopicSubscriptionCloseRequest, s *SubscriptionEvent) (*zbmsgpack.TopicSubscriptionAckRequest, error)
}

// ZeebeAPI defines behaviour of client API behaviour.
type ZeebeAPI interface {
	zbexchange.LikeExchangeSvc
	LikeTaskSubscriptionSvc
	LikeTopicSubscriptionSvc
}
