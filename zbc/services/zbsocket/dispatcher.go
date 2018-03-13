package zbsocket

import (
	"time"

	"github.com/zeebe-io/zbc-go/zbc/models/zbmsgpack"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsbe"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsubscriptions"

	"strconv"
	"sync"

	"github.com/zeebe-io/zbc-go/zbc/common"
	"github.com/zeebe-io/zbc-go/zbc/models/zbdispatch"
)

type safeTaskSubscriptionDispatcher struct {
	subscriptions *sync.Map
}

func (sd *safeTaskSubscriptionDispatcher) AddTaskSubscription(key uint64, ch chan *zbsubscriptions.SubscriptionEvent) {
	sd.subscriptions.Store(key, ch)
}

func (sd *safeTaskSubscriptionDispatcher) RemoveTaskSubscription(key uint64) {
	sd.subscriptions.Delete(key)
}

func (sd *safeTaskSubscriptionDispatcher) GetTaskChannel(key uint64) chan *zbsubscriptions.SubscriptionEvent {
	if ch, ok := sd.subscriptions.Load(key); ok {
		return ch.(chan *zbsubscriptions.SubscriptionEvent)
	}
	return nil
}

func (sd *safeTaskSubscriptionDispatcher) DispatchTaskEvent(key uint64, message *zbsbe.SubscribedEvent, task *zbmsgpack.Task) error {
	for i := 0; i < 3; i++ {
		if ch, ok := sd.subscriptions.Load(key); ch != nil && ok {
			c := ch.(chan *zbsubscriptions.SubscriptionEvent)
			task := &zbsubscriptions.SubscriptionEvent{
				Task: task, Event: message,
			}
			zbcommon.ZBL.Debug().Str("component", "safeTaskSubscriptionDispatcher").Str("len(ch)", strconv.Itoa(len(c))).Msgf("Dispatching task event")
			c <- task
			zbcommon.ZBL.Debug().Str("component", "safeTaskSubscriptionDispatcher").Str("len(ch)", strconv.Itoa(len(c))).Msgf("Dispatched")
			return nil
		}

		sleep := time.Duration(100 * i)
		time.Sleep(time.Duration(time.Millisecond * sleep))
	}

	return zbcommon.ErrWrongSubscriptionKey
}

type safeTopicSubscriptionDispatcher struct {
	subscriptions *sync.Map
}

func (sd *safeTopicSubscriptionDispatcher) AddTopicSubscription(key uint64, ch chan *zbsubscriptions.SubscriptionEvent) {
	zbcommon.ZBL.Debug().Msgf("Adding subscription %d to map", key)
	sd.subscriptions.Store(key, ch)
}

func (sd *safeTopicSubscriptionDispatcher) RemoveTopicSubscription(key uint64) {
	sd.subscriptions.Delete(key)
}

func (sd *safeTopicSubscriptionDispatcher) GetTopicChannel(key uint64) chan *zbsubscriptions.SubscriptionEvent {
	if ch, ok := sd.subscriptions.Load(key); ch != nil && ok {
		return ch.(chan *zbsubscriptions.SubscriptionEvent)
	}
	return nil
}

func (sd *safeTopicSubscriptionDispatcher) DispatchTopicEvent(key uint64, message *zbsbe.SubscribedEvent) error {
	for i := 0; i < 3; i++ {
		if ch, ok := sd.subscriptions.Load(key); ch != nil && ok {
			event := &zbsubscriptions.SubscriptionEvent{
				Event: message,
			}
			zbcommon.ZBL.Debug().Msgf("dispatch topic event")
			ch.(chan *zbsubscriptions.SubscriptionEvent) <- event
			zbcommon.ZBL.Debug().Msgf("topic event dispatched")
			return nil
		}

		sleep := time.Duration(100 * i)
		time.Sleep(time.Duration(time.Millisecond * sleep))
	}

	return zbcommon.ErrWrongSubscriptionKey
}

type subscriptionsDispatcher struct {
	*safeTaskSubscriptionDispatcher
	*safeTopicSubscriptionDispatcher
}

func newSubscriptionsDispatcher() *subscriptionsDispatcher {
	return &subscriptionsDispatcher{
		safeTaskSubscriptionDispatcher: &safeTaskSubscriptionDispatcher{
			subscriptions: &sync.Map{},
		},
		safeTopicSubscriptionDispatcher: &safeTopicSubscriptionDispatcher{
			subscriptions: &sync.Map{},
		},
	}
}

type transactionDispatcher struct {
	lastTransactionSeed uint64
	activeTransactions  []*RequestWrapper
}

func (d *transactionDispatcher) AddTransaction(request *RequestWrapper) {
	d.lastTransactionSeed += zbcommon.RequestQueueSize + 1
	request.WithRequestID(d.lastTransactionSeed)

	index := d.lastTransactionSeed & (zbcommon.RequestQueueSize - 1)
	d.activeTransactions[index] = request
}

func (d *transactionDispatcher) GetTransaction(requestID uint64) *RequestWrapper {
	index := requestID & (zbcommon.RequestQueueSize - 1)
	return d.activeTransactions[index]
}

func (d *transactionDispatcher) DispatchTransaction(requestID uint64, response *zbdispatch.Message) error {
	if transaction := d.GetTransaction(requestID); transaction != nil {
		transaction.ResponseCh <- response
		return nil
	}
	return zbcommon.ErrWrongTransactionIndex

}

func newTransactionDispatcher() *transactionDispatcher {
	return &transactionDispatcher{
		lastTransactionSeed: 0,
		activeTransactions:  make([]*RequestWrapper, zbcommon.RequestQueueSize),
	}
}

type messageDispatcher struct {
	*subscriptionsDispatcher
	*transactionDispatcher
}

func newMessageDispatcher() *messageDispatcher {
	return &messageDispatcher{
		subscriptionsDispatcher: newSubscriptionsDispatcher(),
		transactionDispatcher:   newTransactionDispatcher(),
	}
}
