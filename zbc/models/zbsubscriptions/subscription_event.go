package zbsubscriptions

import (
	"encoding/json"
	"fmt"

	"github.com/zeebe-io/zbc-go/zbc/models/zbmsgpack"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsbe"
)

// SubscriptionEvent is used on test-task-subscriptions and topic subscription.
type SubscriptionEvent struct {
	Task  *zbmsgpack.Task
	Event *zbsbe.SubscribedEvent
}

func (se *SubscriptionEvent) String() string {
	b, _ := json.MarshalIndent(se, "", "  ")
	return fmt.Sprintf("%+v", string(b))
}
