package zbsubscriptions

import (
	"encoding/json"
	"fmt"

	"github.com/zeebe-io/zbc-go/zbc/models/zbmsgpack"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsbe"
	"github.com/zeebe-io/zbc-go/zbc/common"
	"github.com/vmihailenco/msgpack"
)

// SubscriptionEvent is used on test-task-subscriptions and topic subscription.
type SubscriptionEvent struct {
	Task  *zbmsgpack.Task
	Event *zbsbe.SubscribedEvent
}

func (event *SubscriptionEvent) Update( data interface{}) error {
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

func (event *SubscriptionEvent) Load(userType interface{}) error {
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
		Task: task,
		Event: event,
	}
}