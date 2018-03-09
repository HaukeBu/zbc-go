package zbsocket

import (
	"github.com/zeebe-io/zbc-go/zbc/models/zbdispatch"
	"github.com/zeebe-io/zbc-go/zbc/models/zbmsgpack"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsbe"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsubscriptions"
)

type LikeMessageDispatcher interface {
	AddTaskSubscription(key uint64, ch chan *zbsubscriptions.SubscriptionEvent)
	AddTopicSubscription(key uint64, ch chan *zbsubscriptions.SubscriptionEvent)

	RemoveTaskSubscription(key uint64)
	RemoveTopicSubscription(key uint64)

	GetTaskChannel(key uint64) chan *zbsubscriptions.SubscriptionEvent
	GetTopicChannel(key uint64) chan *zbsubscriptions.SubscriptionEvent

	DispatchTaskEvent(key uint64, message *zbsbe.SubscribedEvent, task *zbmsgpack.Task) error
	DispatchTopicEvent(key uint64, message *zbsbe.SubscribedEvent) error

	AddTransaction(request *RequestWrapper)
	GetTransaction(requestID uint64) *RequestWrapper
	DispatchTransaction(requestID uint64, response *zbdispatch.Message) error
}
type LikeSocket interface {
	LikeMessageDispatcher

	Dial(addr string) error
	Send(message *RequestWrapper) error
}
