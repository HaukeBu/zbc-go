package zbexchange

import (
	"io/ioutil"
	"path/filepath"

	"github.com/zeebe-io/zbc-go/zbc/models/zbdispatch"
	"github.com/zeebe-io/zbc-go/zbc/models/zbmsgpack"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsubscriptions"

	"github.com/zeebe-io/zbc-go/zbc/services/zbtopology"

	"github.com/zeebe-io/zbc-go/zbc/common"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsbe"
	"github.com/zeebe-io/zbc-go/zbc/services/zbsocket"
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
	zbcommon.ZBL.Debug().Str("component", "ExchangeSvc").Str("method", "CreateWorkflowInstance").Msg("creating workflow instance")
	pid, err := rm.NextPartitionID(topic)
	if err != nil {
		zbcommon.ZBL.Error().Str("component", "ExchangeSvc").Str("method", "CreateWorkflowInstance").Msgf("cannot find next partitionID: %+v\n", err)
		return nil, err
	}
	message := rm.CreateWorkflowInstanceRequest(*pid, 0, topic, wfi)
	request := zbsocket.NewRequestWrapper(message)
	resp, err := rm.ExecuteRequest(request)
	if err != nil {
		zbcommon.ZBL.Error().Str("component", "ExchangeSvc").Str("method", "CreateWorkflowInstance").Msgf("cannot execute request: %+v\n", err)
		return nil, err
	}
	zbcommon.ZBL.Debug().Str("component", "ExchangeSvc").Str("method", "CreateWorkflowInstance").Msg("workflow instance created")
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
func (rm *ExchangeSvc) OpenTaskPartition(
	ch chan *zbsubscriptions.SubscriptionEvent,
	partitionID uint16,
	lockOwner,
	taskType string,
	lockDuration uint64,
	credits int32) *zbmsgpack.TaskSubscriptionInfo {

	message := rm.OpenTaskSubscriptionRequest(partitionID, lockOwner, taskType, lockDuration, credits)
	request := zbsocket.NewRequestWrapper(message)
	resp, err := rm.ExecuteRequest(request)

	if err != nil {
		return nil
	}

	if taskSubInfo := rm.UnmarshalTaskSubscriptionInfo(resp); taskSubInfo != nil {
		request.Sock.AddTaskSubscription(taskSubInfo.SubscriberKey, ch)
		return taskSubInfo
	}
	return nil
}

func (rm *ExchangeSvc) OpenTopicPartition(
	ch chan *zbsubscriptions.SubscriptionEvent,
	partitionID uint16,
	topic, subscriptionName string,
	startPosition int64,
	forceStart bool,
	prefetchCapacity int32) *zbmsgpack.TopicSubscriptionInfo {

	message := rm.OpenTopicSubscriptionRequest(partitionID, topic, subscriptionName, startPosition, forceStart, prefetchCapacity)
	request := zbsocket.NewRequestWrapper(message)
	resp, err := rm.ExecuteRequest(request)
	if err != nil {
		return nil
	}

	subscriberKey := (resp.SbeMessage).(*zbsbe.ExecuteCommandResponse).Key
	if topicSubInfo := rm.UnmarshalTopicSubscriptionInfo(resp); topicSubInfo != nil {
		request.Sock.AddTopicSubscription(subscriberKey, ch)
		topicSubInfo.SubscriberKey = subscriberKey
		return topicSubInfo
	}
	return nil
}

func (rm *ExchangeSvc) IncreaseTaskSubscriptionCredits(partitionID uint16, task *zbmsgpack.TaskSubscriptionInfo) (*zbmsgpack.TaskSubscriptionInfo, error) {
	zbcommon.ZBL.Debug().Msgf("increasing task credits :: %+v", task)
	message := rm.IncreaseTaskSubscriptionCreditsRequest(partitionID, task)

	request := zbsocket.NewRequestWrapper(message)
	resp, err := rm.ExecuteRequest(request)
	if err != nil {
		zbcommon.ZBL.Debug().Msgf("credits refresh failed. What do we do?")
		return nil, err
	}

	zbcommon.ZBL.Debug().Msgf("credits increased")
	return rm.UnmarshalTaskSubscriptionInfo(resp), nil
}

func (rm *ExchangeSvc) CloseTaskSubscriptionPartition(partitionID uint16, task *zbmsgpack.TaskSubscriptionInfo) (*zbdispatch.Message, error) {
	message := rm.CloseTaskSubscriptionRequest(partitionID, task)
	request := zbsocket.NewRequestWrapper(message)
	resp, err := rm.ExecuteRequest(request)

	request.Sock.RemoveTaskSubscription(task.SubscriberKey)
	return resp, err
}

func (rm *ExchangeSvc) CloseTopicSubscriptionPartition(topicPartition *zbmsgpack.TopicSubscriptionCloseRequest) (*zbdispatch.Message, error) {
	message := rm.CloseTopicSubscriptionRequest(topicPartition)
	request := zbsocket.NewRequestWrapper(message)
	resp, err := rm.ExecuteRequest(request)

	request.Sock.RemoveTopicSubscription(topicPartition.SubscriberKey)
	return resp, err
}

func (rm *ExchangeSvc) TopicSubscriptionAck(subName string, s *zbsubscriptions.SubscriptionEvent) (*zbmsgpack.TopicSubscriptionAckRequest, error) {
	message := rm.TopicSubscriptionAckRequest(subName, s.Event.Position, s.Event.PartitionId)
	request := zbsocket.NewRequestWrapper(message)
	resp, err := rm.ExecuteRequest(request)
	return rm.UnmarshalTopicSubAck(resp), err
}

func (rm *ExchangeSvc) CreateTopic(name string, partitionNum int) (*zbmsgpack.CreateTopic, error) {
	topic := zbmsgpack.NewTopic(name, zbcommon.TopicCreate, partitionNum)
	message := rm.CreateTopicRequest(topic)
	request := zbsocket.NewRequestWrapper(message)

	resp, err := rm.ExecuteRequest(request)
	if err != nil {
		return nil, err
	}
	return rm.UnmarshalTopic(resp), nil
}

func NewExchangeSvc(bootstrapAddr string) (*ExchangeSvc, error) {
	zbcommon.ZBL.Debug().Msg("creating new request manager")

	topologySvc, err := zbtopology.NewTopologySvc(bootstrapAddr)
	if err != nil {
		return nil, err
	}

	return &ExchangeSvc{
		LikeTopologySvc: topologySvc,
		RequestFactory:  zbdispatch.NewRequestFactory(),
		ResponseHandler: zbdispatch.NewResponseHandler(),
	}, nil
}
