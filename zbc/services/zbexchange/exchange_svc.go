package zbexchange

import (
	"io/ioutil"
	"path/filepath"

	"github.com/zeebe-io/zbc-go/zbc/models/zbdispatch"
	"github.com/zeebe-io/zbc-go/zbc/models/zbmsgpack"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsubscriptions"

	"github.com/zeebe-io/zbc-go/zbc/services/zbtopology"

	"github.com/zeebe-io/zbc-go/zbc/common"
	"github.com/zeebe-io/zbc-go/zbc/services/zbsocket"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsbe"
)

type ExchangeSvc struct {
	zbtopology.LikeTopologySvc

	*zbdispatch.RequestFactory
	*zbdispatch.ResponseHandler
}

func (rm *ExchangeSvc) CreateTask(topic string, task *zbmsgpack.Task) (*zbmsgpack.Task, error) {
	pid, err := rm.NextPartitionID(topic)
	if err != nil {
		return nil, err
	}

	PartitionID := *pid

	message := rm.CreateTaskRequest(PartitionID, 0, task)
	request := zbsocket.NewRequestWrapper(message)
	resp, err := rm.ExecuteRequest(request)
	if err != nil {
		return nil, err
	}
	return rm.UnmarshalTask(resp), nil
}

func (rm *ExchangeSvc) CreateWorkflow(topic string, resource ...*zbmsgpack.Resource) (*zbmsgpack.Workflow, error) {
	message := rm.DeployWorkflowRequest(topic, resource)
	request := zbsocket.NewRequestWrapper(message)
	resp, err := rm.ExecuteRequest(request)
	if err != nil {
		return nil, err
	}
	return rm.UnmarshalWorkflow(resp), nil
}

// CreateWorkflowFromFile will read workflow file and return message pack workflow object.
func (rm *ExchangeSvc) CreateWorkflowFromFile(topic, resourceType, path string) (*zbmsgpack.Workflow, error) {
	if len(path) == 0 {
		return nil, zbcommon.ErrResourceNotFound
	}

	filename, _ := filepath.Abs(path)
	definition, err := ioutil.ReadFile(filename)
	resource := &zbmsgpack.Resource{
		ResourceName: path,
		ResourceType: resourceType,
		Resource:     definition,
	}

	if err != nil {
		return nil, zbcommon.ErrResourceNotFound
	}
	return rm.CreateWorkflow(topic, resource)
}

func (rm *ExchangeSvc) CreateWorkflowInstance(topic string, wfi *zbmsgpack.WorkflowInstance) (*zbmsgpack.WorkflowInstance, error) {
	zbcommon.ZBL.Debug().Msg("creating workflow instance")
	pid, err := rm.NextPartitionID(topic)
	if err != nil {
		return nil, err
	}

	PartitionID := *pid
	message := rm.CreateWorkflowInstanceRequest(PartitionID, 0, topic, wfi)
	request := zbsocket.NewRequestWrapper(message)
	resp, err := rm.ExecuteRequest(request)
	if err != nil {
		return nil, err
	}
	return rm.UnmarshalWorkflowInstance(resp), nil
}

func (rm *ExchangeSvc) CompleteTask(task *zbsubscriptions.SubscriptionEvent) (*zbmsgpack.Task, error) {
	message := rm.CompleteTaskRequest(task)
	request := zbsocket.NewRequestWrapper(message)
	resp, err := rm.ExecuteRequest(request)
	if err != nil {
		return nil, err
	}
	return rm.UnmarshalTask(resp), nil
}

func (rm *ExchangeSvc) FailTask(task *zbsubscriptions.SubscriptionEvent) (*zbmsgpack.Task, error) {
	message := rm.FailTaskRequest(task)
	request := zbsocket.NewRequestWrapper(message)
	resp, err := rm.ExecuteRequest(request)
	if err != nil {
		return nil, err
	}
	return rm.UnmarshalTask(resp), nil
}

// OpenTaskPartition will open a test-task-subscriptions subscription on one partition.
func (rm *ExchangeSvc) OpenTaskPartition(partitionID uint16, lockOwner, taskType string, credits int32) (*zbmsgpack.TaskSubscriptionInfo, chan *zbsubscriptions.SubscriptionEvent) {
	subscriptionCh := make(chan *zbsubscriptions.SubscriptionEvent, zbcommon.SubscriptionPipelineQueueSize)

	message := rm.OpenTaskSubscriptionRequest(partitionID, lockOwner, taskType, credits)
	request := zbsocket.NewRequestWrapper(message)
	resp, err := rm.ExecuteRequest(request)

	if err != nil {
		return nil, nil
	}

	if taskSubInfo := rm.UnmarshalTaskSubscriptionInfo(resp); taskSubInfo != nil {
		// MARK: the caller is responsible for attaching subscriptionCh to the right socket
		request.Sock.AddTaskSubscription(taskSubInfo.SubscriberKey, subscriptionCh)
		return taskSubInfo, subscriptionCh
	}
	return nil, nil
}

func (rm *ExchangeSvc) OpenTopicPartition(partitionID uint16, topic, subscriptionName string, startPosition int64) (*zbmsgpack.TopicSubscriptionInfo, chan *zbsubscriptions.SubscriptionEvent) {
	subscriptionCh := make(chan *zbsubscriptions.SubscriptionEvent, zbcommon.SubscriptionPipelineQueueSize)

	message := rm.OpenTopicSubscriptionRequest(partitionID, topic, subscriptionName, startPosition)
	request := zbsocket.NewRequestWrapper(message)
	resp, err := rm.ExecuteRequest(request)
	if err != nil {
		return nil, nil
	}

	subscriberKey := (resp.SbeMessage).(*zbsbe.ExecuteCommandResponse).Key
	if topicSubInfo := rm.UnmarshalTopicSubscriptionInfo(resp); topicSubInfo != nil {
		request.Sock.AddTopicSubscription(subscriberKey, subscriptionCh)
		topicSubInfo.SubscriberKey = subscriberKey
		return topicSubInfo, subscriptionCh
	}
	return nil, nil
}

func (rm *ExchangeSvc) IncreaseTaskSubscriptionCredits(task *zbmsgpack.TaskSubscriptionInfo) (*zbmsgpack.TaskSubscriptionInfo, error) {
	zbcommon.ZBL.Debug().Msgf("increasing task credits ", task.Credits)
	message := rm.IncreaseTaskSubscriptionCreditsRequest(task)

	request := zbsocket.NewRequestWrapper(message)

	resp, err := rm.ExecuteRequest(request)
	if err != nil {
		zbcommon.ZBL.Debug().Msgf("task increasage failed. What do we do?")
		return nil, err
	}

	return rm.UnmarshalTaskSubscriptionInfo(resp), nil
}

func (rm *ExchangeSvc) CloseTaskSubscriptionPartition(task *zbmsgpack.TaskSubscriptionInfo) (*zbdispatch.Message, error) {
	message := rm.CloseTaskSubscriptionRequest(task)
	request := zbsocket.NewRequestWrapper(message)
	resp, err := rm.ExecuteRequest(request)

	request.Sock.RemoveTaskSubscription(task.SubscriberKey)
	return resp, err
}

func (rm *ExchangeSvc) CloseTopicSubscriptionPartition(topicPartition *zbmsgpack.TopicSubscription) (*zbdispatch.Message, error) {
	message := rm.CloseTopicSubscriptionRequest(topicPartition)
	request := zbsocket.NewRequestWrapper(message)
	resp, err := rm.ExecuteRequest(request)

	request.Sock.RemoveTopicSubscription(topicPartition.SubscriberKey)
	return resp, err
}

func (rm *ExchangeSvc) TopicSubscriptionAck(ts *zbmsgpack.TopicSubscription, s *zbsubscriptions.SubscriptionEvent) (*zbmsgpack.TopicSubscriptionAck, error) {
	message := rm.TopicSubscriptionAckRequest(ts, s)
	request := zbsocket.NewRequestWrapper(message)
	resp, err := rm.ExecuteRequest(request)
	return rm.UnmarshalTopicSubAck(resp), err
}

func (rm *ExchangeSvc) CreateTopic(name string, partitionNum int) (*zbmsgpack.Topic, error) {
	topic := zbmsgpack.NewTopic(name, zbcommon.TopicCreate, partitionNum)
	message := rm.CreateTopicRequest(topic)
	request := zbsocket.NewRequestWrapper(message)

	resp, err := rm.ExecuteRequest(request)
	if err != nil {
		return nil, err
	}
	return rm.UnmarshalTopic(resp), nil
}

func NewExchangeSvc(bootstrapAddr string) *ExchangeSvc {

	zbcommon.ZBL.Debug().Msg("creating new request manager")

	return &ExchangeSvc{
		LikeTopologySvc: zbtopology.NewTopologySvc(bootstrapAddr),
		RequestFactory:  zbdispatch.NewRequestFactory(),
		ResponseHandler: zbdispatch.NewResponseHandler(),
	}
}
