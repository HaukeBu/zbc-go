package zbdispatch

import (
	"github.com/vmihailenco/msgpack"

	"github.com/zeebe-io/zbc-go/zbc/models/zbmsgpack"
	"github.com/zeebe-io/zbc-go/zbc/models/zbprotocol"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsbe"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsubscriptions"

	"github.com/zeebe-io/zbc-go/zbc/common"
)

type RequestFactory struct{}

func (rf *RequestFactory) headers(t interface{}) *Headers {
	switch v := t.(type) {

	case *zbsbe.ExecuteCommandRequest:
		length := uint32(zbcommon.LengthFieldSize + len(v.Command))
		length += uint32(v.SbeBlockLength()) + zbcommon.TotalHeaderSize

		var headers Headers
		headers.SetSbeMessageHeader(&zbsbe.MessageHeader{
			BlockLength: v.SbeBlockLength(),
			TemplateId:  v.SbeTemplateId(),
			SchemaId:    v.SbeSchemaId(),
			Version:     v.SbeSchemaVersion(),
		})

		headers.SetRequestResponseHeader(zbprotocol.NewRequestResponseHeader())
		headers.SetTransportHeader(zbprotocol.NewTransportHeader(zbprotocol.RequestResponse))

		headers.SetFrameHeader(zbprotocol.NewFrameHeader(uint32(length), 0, 0, 0, 0))
		return &headers

	case *zbsbe.ControlMessageRequest:
		length := uint32(zbcommon.LengthFieldSize + len(v.Data))
		length += uint32(v.SbeBlockLength()) + zbcommon.TotalHeaderSize

		var headers Headers
		headers.SetSbeMessageHeader(&zbsbe.MessageHeader{
			BlockLength: v.SbeBlockLength(),
			TemplateId:  v.SbeTemplateId(),
			SchemaId:    v.SbeSchemaId(),
			Version:     v.SbeSchemaVersion(),
		})
		headers.SetRequestResponseHeader(zbprotocol.NewRequestResponseHeader())
		headers.SetTransportHeader(zbprotocol.NewTransportHeader(zbprotocol.RequestResponse))

		// Writer will set FrameHeader after serialization to byte array.
		headers.SetFrameHeader(zbprotocol.NewFrameHeader(uint32(length), 0, 0, 0, 0))
		return &headers
	}

	return nil
}

func (rf *RequestFactory) newCommandMessage(commandRequest *zbsbe.ExecuteCommandRequest, command interface{}) *Message {
	var msg Message

	b, err := msgpack.Marshal(command)
	if err != nil {
		return nil
	}
	commandRequest.Command = b
	msg.SetSbeMessage(commandRequest)
	msg.SetHeaders(rf.headers(commandRequest))

	return &msg
}

func (rf *RequestFactory) newControlMessage(req *zbsbe.ControlMessageRequest, payload interface{}) *Message {
	var msg Message

	b, err := msgpack.Marshal(payload)
	if err != nil {
		return nil
	}
	req.Data = b
	msg.SetSbeMessage(req)
	msg.SetHeaders(rf.headers(req))

	return &msg
}

func (rf *RequestFactory) CreatePartitionRequest() *Message {
	controlMessage := &zbsbe.ControlMessageRequest{
		MessageType: zbsbe.ControlMessageType.REQUEST_PARTITIONS,
		PartitionId: 0,
		Data:        nil,
	}
	return rf.newControlMessage(controlMessage, nil)

}

func (rf *RequestFactory) CreateTaskRequest(partition uint16, position uint64, task *zbmsgpack.Task) *Message {
	commandRequest := &zbsbe.ExecuteCommandRequest{
		PartitionId: partition,
		Position:    position,
		Command:     []uint8{},
	}

	commandRequest.Key = commandRequest.KeyNullValue()
	commandRequest.EventType = zbsbe.EventTypeEnum(0)

	return rf.newCommandMessage(commandRequest, task)
}

func (rf *RequestFactory) CompleteTaskRequest(subEvent *zbsubscriptions.SubscriptionEvent) *Message {
	subEvent.Event.(*zbmsgpack.Task).State = zbcommon.TaskComplete

	cmdReq := &zbsbe.ExecuteCommandRequest{
		PartitionId: subEvent.Metadata.PartitionId,
		Position:    subEvent.Metadata.Position,
		Key:         subEvent.Metadata.Key,
	}

	return rf.newCommandMessage(cmdReq, subEvent.Event.(*zbmsgpack.Task))
}

func (rf *RequestFactory) FailTaskRequest(subEvent *zbsubscriptions.SubscriptionEvent) *Message {
	subEvent.Event.(*zbmsgpack.Task).State = zbcommon.TaskFail

	cmdReq := &zbsbe.ExecuteCommandRequest{
		PartitionId: subEvent.Metadata.PartitionId,
		Position:    subEvent.Metadata.Position,
		Key:         subEvent.Metadata.Key,
	}

	return rf.newCommandMessage(cmdReq, subEvent.Event.(*zbmsgpack.Task))
}

func (rf *RequestFactory) UpdatePayloadRequest(event *zbsubscriptions.SubscriptionEvent) *Message {
	wfi := event.Event.(*zbmsgpack.WorkflowInstanceEvent)
	wfi.State = zbcommon.UpdatePayload

	cmdReq := &zbsbe.ExecuteCommandRequest{
		PartitionId: event.Metadata.PartitionId,
		Position:    event.Metadata.Position,
		Key:         event.Metadata.Key,
		EventType:   zbsbe.EventType.WORKFLOW_INSTANCE_EVENT,
	}

	return rf.newCommandMessage(cmdReq, wfi)
}

func (rf *RequestFactory) CreateWorkflowInstanceRequest(partition uint16, position uint64, topic string, wf *zbmsgpack.CreateWorkflowInstance) *Message {
	commandRequest := &zbsbe.ExecuteCommandRequest{
		PartitionId: partition,
		Position:    position,
		Command:     []uint8{},
	}

	commandRequest.Key = commandRequest.KeyNullValue()
	commandRequest.EventType = zbsbe.EventTypeEnum(5)

	return rf.newCommandMessage(commandRequest, wf)

}

func (rf *RequestFactory) TopologyRequest() *Message {
	t := &zbmsgpack.TopologyRequest{}

	cmr := &zbsbe.ControlMessageRequest{
		MessageType: zbsbe.ControlMessageType.REQUEST_TOPOLOGY,
		// MARK: PartitionID is by default 0, which is excatly what we want since partitionID 0 is system-partition
		Data: nil,
	}

	return rf.newControlMessage(cmr, t)
}

func (rf *RequestFactory) DeployWorkflowRequest(topic string, resources []*zbmsgpack.Resource) *Message {
	deployment := zbmsgpack.DeployWorkflow{
		State:     zbcommon.CreateDeployment,
		Resources: resources,
		TopicName: topic,
	}

	commandRequest := &zbsbe.ExecuteCommandRequest{
		PartitionId: 0,
		Position:    0,
		Command:     []uint8{},
	}

	commandRequest.Key = commandRequest.KeyNullValue()
	commandRequest.EventType = zbsbe.EventTypeEnum(4)

	return rf.newCommandMessage(commandRequest, deployment)
}

func (rf *RequestFactory) OpenTaskSubscriptionRequest(partitionID uint16, lockOwner, taskType string, lockDuration uint64, credits int32) *Message {
	taskSub := &zbmsgpack.TaskSubscriptionInfo{
		Credits:       credits,
		LockDuration:  lockDuration,
		LockOwner:     lockOwner,
		SubscriberKey: 0,
		TaskType:      taskType,
	}

	var msg Message
	b, err := msgpack.Marshal(taskSub)
	if err != nil {
		return nil
	}

	controlRequest := &zbsbe.ControlMessageRequest{
		MessageType: zbsbe.ControlMessageType.ADD_TASK_SUBSCRIPTION,
		PartitionId: partitionID,
		Data:        b,
	}

	msg.SetSbeMessage(controlRequest)
	msg.SetHeaders(rf.headers(controlRequest))

	return &msg
}

func (rf *RequestFactory) IncreaseTaskSubscriptionCreditsRequest(partitionID uint16, ts *zbmsgpack.TaskSubscriptionInfo) *Message {
	var msg Message

	b, err := msgpack.Marshal(ts)
	if err != nil {
		return nil
	}

	controlRequest := &zbsbe.ControlMessageRequest{
		MessageType: zbsbe.ControlMessageType.INCREASE_TASK_SUBSCRIPTION_CREDITS,
		PartitionId: partitionID,
		Data:        b,
	}

	msg.SetSbeMessage(controlRequest)
	msg.SetHeaders(rf.headers(controlRequest))

	return &msg
}

func (rf *RequestFactory) CloseTaskSubscriptionRequest(partitionID uint16, ts *zbmsgpack.TaskSubscriptionInfo) *Message {
	var msg Message

	b, err := msgpack.Marshal(ts)
	if err != nil {
		return nil
	}

	controlRequest := &zbsbe.ControlMessageRequest{
		MessageType: zbsbe.ControlMessageType.REMOVE_TASK_SUBSCRIPTION,
		PartitionId: partitionID,
		Data:        b,
	}

	msg.SetSbeMessage(controlRequest)
	msg.SetHeaders(rf.headers(controlRequest))

	return &msg
}

func (rf *RequestFactory) CloseTopicSubscriptionRequest(ts *zbmsgpack.TopicSubscriptionCloseRequest) *Message {
	var msg Message

	b, err := msgpack.Marshal(ts)
	if err != nil {
		return nil
	}

	controlRequest := &zbsbe.ControlMessageRequest{
		MessageType: zbsbe.ControlMessageType.REMOVE_TOPIC_SUBSCRIPTION,
		PartitionId: ts.PartitionID,
		Data:        b,
	}

	msg.SetSbeMessage(controlRequest)
	msg.SetHeaders(rf.headers(controlRequest))

	return &msg
}

func (rf *RequestFactory) OpenTopicSubscriptionRequest(partitionID uint16, topic, subName string, startPosition int64, forceStart bool, prefetchCapacity int32) *Message {
	ts := &zbmsgpack.TopicSubscriptionInfo{
		StartPosition:    startPosition,
		Name:             subName,
		PrefetchCapacity: prefetchCapacity,
		ForceStart:       forceStart,
		State:            zbcommon.TopicSubscriptionSubscribeState,
	}

	execCommandRequest := &zbsbe.ExecuteCommandRequest{
		PartitionId: partitionID,
		Position:    0,
		EventType:   zbsbe.EventType.SUBSCRIBER_EVENT,
	}

	execCommandRequest.Key = execCommandRequest.KeyNullValue()

	var msg Message
	b, err := msgpack.Marshal(ts)
	if err != nil {
		return nil
	}

	execCommandRequest.Command = b
	msg.SetSbeMessage(execCommandRequest)
	msg.SetHeaders(rf.headers(execCommandRequest))

	return &msg
}

func (rf *RequestFactory) TopicSubscriptionAckRequest(subscriptionName string, position uint64, partitionID uint16) *Message {
	tsa := &zbmsgpack.TopicSubscriptionAckRequest{
		Name:        subscriptionName,
		AckPosition: position,
		State:       zbcommon.TopicSubscriptionAckState,
	}

	execCommandRequest := &zbsbe.ExecuteCommandRequest{
		PartitionId: partitionID,
		Position:    0,
		EventType:   zbsbe.EventType.SUBSCRIPTION_EVENT,
	}

	execCommandRequest.Key = execCommandRequest.KeyNullValue()

	var msg Message
	b, err := msgpack.Marshal(tsa)
	if err != nil {
		return nil
	}

	execCommandRequest.Command = b
	msg.SetSbeMessage(execCommandRequest)
	msg.SetHeaders(rf.headers(execCommandRequest))

	return &msg
}

func (rf *RequestFactory) CreateTopicRequest(topic *zbmsgpack.CreateTopic) *Message {
	execCommandRequest := &zbsbe.ExecuteCommandRequest{
		PartitionId: 0,
		Position:    0,
		EventType:   zbsbe.EventType.TOPIC_EVENT,
	}

	execCommandRequest.Key = execCommandRequest.KeyNullValue()

	var msg Message
	b, err := msgpack.Marshal(topic)
	if err != nil {
		return nil
	}

	execCommandRequest.Command = b
	msg.SetSbeMessage(execCommandRequest)
	msg.SetHeaders(rf.headers(execCommandRequest))

	return &msg
}

func NewRequestFactory() *RequestFactory {
	return &RequestFactory{}
}
