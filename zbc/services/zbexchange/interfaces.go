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
	CreateWorkflow(topic string, resources ...*zbmsgpack.Resource) (*zbmsgpack.DeployWorkflow, error)
	CreateWorkflowFromFile(topic, resourceType, path string) (*zbmsgpack.DeployWorkflow, error)
	CreateWorkflowInstance(topic string, workflowInstance *zbmsgpack.CreateWorkflowInstance) (*zbmsgpack.CreateWorkflowInstance, error)

	OpenTopicPartition(ch chan *zbsubscriptions.SubscriptionEvent, partitionID uint16, topic, subscriptionName string, startPosition int64, forceStart bool, prefetchCapacity int32) *zbmsgpack.TopicSubscriptionInfo
	OpenTaskPartition(ch chan *zbsubscriptions.SubscriptionEvent, partitionID uint16, lockOwner, taskType string, lockDuration uint64, credits int32) *zbmsgpack.TaskSubscriptionInfo

	IncreaseTaskSubscriptionCredits(partitionID uint16, task *zbmsgpack.TaskSubscriptionInfo) (*zbmsgpack.TaskSubscriptionInfo, error)
	TopicSubscriptionAck(subName string, s *zbsubscriptions.SubscriptionEvent) (*zbmsgpack.TopicSubscriptionAckRequest, error)

	CloseTopicSubscriptionPartition(topicPartition *zbmsgpack.TopicSubscriptionCloseRequest) (*zbdispatch.Message, error)
	CloseTaskSubscriptionPartition(partitionID uint16, task *zbmsgpack.TaskSubscriptionInfo) (*zbdispatch.Message, error)

	CompleteTask(task *zbsubscriptions.SubscriptionEvent) (*zbmsgpack.Task, error)
	FailTask(task *zbsubscriptions.SubscriptionEvent) (*zbmsgpack.Task, error)

	UpdatePayload(event *zbsubscriptions.SubscriptionEvent) (*zbmsgpack.CreateWorkflowInstance, error)
}
