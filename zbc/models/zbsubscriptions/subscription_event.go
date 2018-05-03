package zbsubscriptions

import (
	"encoding/json"
	"fmt"

	"github.com/vmihailenco/msgpack"
	"github.com/zeebe-io/zbc-go/zbc/common"
	"github.com/zeebe-io/zbc-go/zbc/models/zbmsgpack"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsbe"
)

// SubscriptionEvent is used on test-task-subscriptions and topic subscription.
type SubscriptionEvent struct {
	Event    interface{}
	Metadata *zbsbe.SubscribedEvent
}

func (event *SubscriptionEvent) GetTask() (*zbmsgpack.Task, error) {
	switch e := event.Event.(type) {
	case *zbmsgpack.Task:
		return e, nil
	}
	return nil, zbcommon.ErrWrongTypeAssertion
}

func (event *SubscriptionEvent) GetWorkflowInstanceEvent() (*zbmsgpack.WorkflowInstanceEvent, error) {
	switch e := event.Event.(type) {
	case *zbmsgpack.WorkflowInstanceEvent:
		return e, nil
	}
	return nil, zbcommon.ErrWrongTypeAssertion
}

func (event *SubscriptionEvent) GetWorkflowEvent() (*zbmsgpack.WorkflowEvent, error) {
	switch e := event.Event.(type) {
	case *zbmsgpack.WorkflowEvent:
		return e, nil
	}
	return nil, zbcommon.ErrWrongTypeAssertion
}

func (event *SubscriptionEvent) GetIncident() (*zbmsgpack.IncidentEvent, error) {
	switch e := event.Event.(type) {
	case *zbmsgpack.IncidentEvent:
		return e, nil
	}
	return nil, zbcommon.ErrWrongTypeAssertion
}

func (event *SubscriptionEvent) GetEvent() (interface{}, error) {
	if event.Metadata == nil {
		return nil, zbcommon.ErrEventIsEmpty
	}

	switch event.Metadata.EventType {

	case zbsbe.EventType.TASK_EVENT:
		var task zbmsgpack.Task
		err := msgpack.Unmarshal(event.Metadata.Event, &task)
		if err != nil {
			return nil, err
		}
		return &task, nil
		break

	case zbsbe.EventType.RAFT_EVENT:
		// TODO
		break

	case zbsbe.EventType.WORKFLOW_INSTANCE_EVENT:
		var instance zbmsgpack.WorkflowInstanceEvent
		err := msgpack.Unmarshal(event.Metadata.Event, &instance)
		if err != nil {
			return nil, err
		}
		return &instance, nil
		break

	case zbsbe.EventType.INCIDENT_EVENT:
		var incident zbmsgpack.IncidentEvent
		err := msgpack.Unmarshal(event.Metadata.Event, &incident)
		if err != nil {
			return nil, err
		}
		return &incident, nil
		break

	case zbsbe.EventType.WORKFLOW_EVENT:
		var workflow zbmsgpack.WorkflowEvent
		err := msgpack.Unmarshal(event.Metadata.Event, &workflow)
		if err != nil {
			return nil, err
		}
		return &workflow, nil
		break

	}

	var m map[string]interface{}
	err := msgpack.Unmarshal(event.Metadata.Event, &m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (event *SubscriptionEvent) UpdatePayload(payload interface{}) error {
	if event.Event == nil {
		return zbcommon.ErrEventIsEmpty
	}

	b, err := msgpack.Marshal(payload)
	if err != nil {
		return err
	}

	switch  e := event.Event.(type) {
	case *zbmsgpack.Task:
		e.Payload = b
	case *zbmsgpack.CreateWorkflowInstance:
		e.Payload = b
	case *zbmsgpack.WorkflowEvent:
		return zbcommon.ErrNoPayloadOnWorkflowEvent
	case *zbmsgpack.IncidentEvent:
		e.Payload = b
	}

	return nil
}

func (event *SubscriptionEvent) LoadTask(userType interface{}) error {
	if event.Event == nil {
		return zbcommon.ErrEventNotTask
	}
	err := msgpack.Unmarshal(event.Event.(*zbmsgpack.Task).Payload, userType)
	return err
}

func (se *SubscriptionEvent) String() string {
	b, _ := json.MarshalIndent(se, "", "  ")
	return fmt.Sprintf("%+v", string(b))
}

func NewSubscriptionEvent(task *zbmsgpack.Task, metadata *zbsbe.SubscribedEvent) *SubscriptionEvent {
	return &SubscriptionEvent{
		Event:    task,
		Metadata: metadata,
	}
}
