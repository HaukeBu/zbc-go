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
	Task  *zbmsgpack.Task
	Event *zbsbe.SubscribedEvent
}

func (event *SubscriptionEvent) GetEvent() (map[string]interface{}, error) {
	if event.Event == nil {
		return nil, zbcommon.ErrEventIsEmpty
	}
	var m map[string]interface{}
	err := msgpack.Unmarshal(event.Event.Event, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (event *SubscriptionEvent) UpdateTask(data interface{}) error {
	if event.Task == nil {
		return zbcommon.ErrEventNotTask
	}
	b, err := msgpack.Marshal(data)
	if err != nil {
		return err
	}

	event.Task.Payload = b
	return nil
}

func (event *SubscriptionEvent) LoadTask(userType interface{}) error {
	if event.Task == nil {
		return zbcommon.ErrEventNotTask
	}
	err := msgpack.Unmarshal(event.Task.Payload, userType)
	return err
}

func (se *SubscriptionEvent) String() string {
	b, _ := json.MarshalIndent(se, "", "  ")
	return fmt.Sprintf("%+v", string(b))
}

func NewSubscriptionEvent(task *zbmsgpack.Task, event *zbsbe.SubscribedEvent) *SubscriptionEvent {
	return &SubscriptionEvent{
		Task:  task,
		Event: event,
	}
}
