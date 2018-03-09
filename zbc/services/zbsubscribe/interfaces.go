package zbsubscribe

import (
	"github.com/zeebe-io/zbc-go/zbc/services/zbexchange"
)

// LikeTaskSubscriptionSvc defines behaviour of TaskSubscriptionSvc.
type LikeTaskSubscriptionSvc interface {
	TaskSubscription(topic, lockOwner, taskType string, cb TaskSubscriptionCallback) (*TaskSubscription, error)
	CloseTaskSubscription(task *TaskSubscription) []error

	// TODO: CompleteTask(test-task-subscriptions *SubscriptionEvent) (*zbmsgpack.Task, error)
	// TODO: FailTask(test-task-subscriptions *zbcommon.SubscriptionEvent) (*zbmsgpack.Task, error)
}


// LikeTopicSubscriptionSvc defines behaviour of TopicSubscriptionSvc.
type LikeTopicSubscriptionSvc interface {
	TopicSubscription(topic, subName string, startPosition int64, cb TopicSubscriptionCallback) (*TopicSubscription, error)
	CloseTopicSubscription(sub *TopicSubscription) []error

	//TopicSubscriptionAck(ts *zbmsgpack.TopicSubscription, s *SubscriptionEvent) (*zbmsgpack.TopicSubscriptionAck, error)
}

// ZeebeAPI defines behaviour of client API behaviour.
type ZeebeAPI interface {
	zbexchange.LikeExchangeSvc
	LikeTaskSubscriptionSvc
	LikeTopicSubscriptionSvc
}
