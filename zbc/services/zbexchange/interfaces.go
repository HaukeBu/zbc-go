package zbexchange

import (
	"github.com/zeebe-io/zbc-go/zbc/models/zbdispatch"
	"github.com/zeebe-io/zbc-go/zbc/models/zbmsgpack"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsubscriptions"
	"github.com/zeebe-io/zbc-go/zbc/services/zbtopology"
)

type LikeExchangeSvc interface {
	zbtopology.LikeTopologySvc

	CreateTopic(name string, partitionNum int) (*zbmsgpack.CreateTopic, error)
	CreateTask(topic string, task *zbmsgpack.Task) (*zbmsgpack.Task, error)
	CreateWorkflow(topic string, resources ...*zbmsgpack.Resource) (*zbmsgpack.Workflow, error)
	CreateWorkflowFromFile(topic, resourceType, path string) (*zbmsgpack.Workflow, error)
	CreateWorkflowInstance(topic string, workflowInstance *zbmsgpack.WorkflowInstance) (*zbmsgpack.WorkflowInstance, error)

	OpenTopicPartition(ch chan *zbsubscriptions.SubscriptionEvent, partitionID uint16, topic, subscriptionName string, startPosition int64) *zbmsgpack.TopicSubscriptionInfo
	OpenTaskPartition(ch chan *zbsubscriptions.SubscriptionEvent, partitionID uint16, lockOwner, taskType string, credits int32) *zbmsgpack.TaskSubscriptionInfo

	IncreaseTaskSubscriptionCredits(task *zbmsgpack.TaskSubscriptionInfo) (*zbmsgpack.TaskSubscriptionInfo, error)
	TopicSubscriptionAck(subName string, s *zbsubscriptions.SubscriptionEvent) (*zbmsgpack.TopicSubscriptionAckRequest, error)

	CloseTopicSubscriptionPartition(topicPartition *zbmsgpack.TopicSubscriptionCloseRequest) (*zbdispatch.Message, error)
	CloseTaskSubscriptionPartition(task *zbmsgpack.TaskSubscriptionInfo) (*zbdispatch.Message, error)

	CompleteTask(task *zbsubscriptions.SubscriptionEvent) (*zbmsgpack.Task, error)
	FailTask(task *zbsubscriptions.SubscriptionEvent) (*zbmsgpack.Task, error)
}