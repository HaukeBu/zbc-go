package zbsubscribe

import (
	"github.com/zeebe-io/zbc-go/zbc/services/zbexchange"
)

// LikeTaskSubscriptionSvc defines behaviour of TaskSubscriptionSvc.
type LikeTaskSubscriptionSvc interface {
	TaskSubscription(topic, lockOwner, taskType string, lockDuration uint64, credits int32, cb TaskSubscriptionCallback) (*TaskSubscription, error)
	CloseTaskSubscription(task *TaskSubscription) []error
}

// LikeTopicSubscriptionSvc defines behaviour of TopicSubscriptionSvc.
type LikeTopicSubscriptionSvc interface {
	TopicSubscription(topic, subName string, prefetchCapacity uint32, startPosition int64, forceStart bool, cb TopicSubscriptionCallback) (*TopicSubscription, error)
	CloseTopicSubscription(sub *TopicSubscription) []error
}

// ZeebeAPI defines behaviour of client API behaviour.
type ZeebeAPI interface {
	zbexchange.LikeExchangeSvc
	LikeTaskSubscriptionSvc
	LikeTopicSubscriptionSvc
}
