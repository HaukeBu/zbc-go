package zbexchange

import (
	"github.com/zeebe-io/zbc-go/zbc/models/zbdispatch"
	"github.com/zeebe-io/zbc-go/zbc/models/zbmsgpack"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsubscriptions"
	"github.com/zeebe-io/zbc-go/zbc/services/zbtopology"
)

type LikeExchangeSvc interface {
	zbtopology.LikeTopologySvc

	CreateTopic(name string, partitionNum int) (*zbmsgpack.Topic, error)
	CreateTask(topic string, task *zbmsgpack.Task) (*zbmsgpack.Task, error)
	CreateWorkflow(topic string, resources ...*zbmsgpack.Resource) (*zbmsgpack.Workflow, error)
	CreateWorkflowFromFile(topic, resourceType, path string) (*zbmsgpack.Workflow, error)
	CreateWorkflowInstance(topic string, workflowInstance *zbmsgpack.WorkflowInstance) (*zbmsgpack.WorkflowInstance, error)

	OpenTopicPartition(partitionID uint16, topic, subscriptionName string, startPosition int64) (*zbmsgpack.TopicSubscriptionInfo, chan *zbsubscriptions.SubscriptionEvent)
	TopicSubscriptionAck(ts *zbmsgpack.TopicSubscription, s *zbsubscriptions.SubscriptionEvent) (*zbmsgpack.TopicSubscriptionAck, error)
	CloseTopicSubscriptionPartition(topicPartition *zbmsgpack.TopicSubscription) (*zbdispatch.Message, error)

	OpenTaskPartition(partitionID uint16, lockOwner, taskType string, credits int32) (*zbmsgpack.TaskSubscriptionInfo, chan *zbsubscriptions.SubscriptionEvent)
	IncreaseTaskSubscriptionCredits(task *zbmsgpack.TaskSubscriptionInfo) (*zbmsgpack.TaskSubscriptionInfo, error)
	CloseTaskSubscriptionPartition(task *zbmsgpack.TaskSubscriptionInfo) (*zbdispatch.Message, error)
	CompleteTask(task *zbsubscriptions.SubscriptionEvent) (*zbmsgpack.Task, error)
	FailTask(task *zbsubscriptions.SubscriptionEvent) (*zbmsgpack.Task, error)
}
