package zbsocket

import (
	"github.com/zeebe-io/zbc-go/zbc/models/zbmsgpack"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsbe"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsubscriptions"

	"github.com/zeebe-io/zbc-go/zbc/common"
	"github.com/zeebe-io/zbc-go/zbc/models/zbdispatch"
	"sync"
)

type safeTaskSubscriptionDispatcher struct {
	*sync.Mutex
	taskSubscriptions  map[uint64]chan *zbsubscriptions.SubscriptionEvent
}

func (sd *safeTaskSubscriptionDispatcher) AddTaskSubscription(key uint64, ch chan *zbsubscriptions.SubscriptionEvent) {
	sd.Lock()
	defer sd.Unlock()

	sd.taskSubscriptions[key] = ch
}


func (sd *safeTaskSubscriptionDispatcher) RemoveTaskSubscription(key uint64) {
	sd.Lock()
	defer sd.Unlock()

	if ch, ok := sd.taskSubscriptions[key]; ok {
		close(ch)
	}
	delete(sd.taskSubscriptions, key)
}


func (sd *safeTaskSubscriptionDispatcher) GetTaskChannel(key uint64) chan *zbsubscriptions.SubscriptionEvent {
	sd.Lock()
	defer sd.Unlock()

	if ch, ok := sd.taskSubscriptions[key]; ok {
		return ch
	}

	return nil
}


func (sd *safeTaskSubscriptionDispatcher) DispatchTaskEvent(key uint64, message *zbsbe.SubscribedEvent, task *zbmsgpack.Task) error {
	sd.Lock()
	defer sd.Unlock()

	if ch, ok := sd.taskSubscriptions[key]; ok {
		task := &zbsubscriptions.SubscriptionEvent{
			Task: task, Event: message,
		}
		ch <- task
		return nil
	}
	return zbcommon.ErrWrongSubscriptionKey
}



type safeTopicSubscriptionDispatcher struct {
	*sync.Mutex
	topicSubscriptions map[uint64]chan *zbsubscriptions.SubscriptionEvent
}

func (sd *safeTopicSubscriptionDispatcher) AddTopicSubscription(key uint64, ch chan *zbsubscriptions.SubscriptionEvent) {
	sd.Lock()
	defer sd.Unlock()

	sd.topicSubscriptions[key] = ch
}


func (sd *safeTopicSubscriptionDispatcher) RemoveTopicSubscription(key uint64) {
	sd.Lock()
	defer sd.Unlock()

	delete(sd.topicSubscriptions, key)
}


func (sd *safeTopicSubscriptionDispatcher) GetTopicChannel(key uint64) chan *zbsubscriptions.SubscriptionEvent {
	sd.Lock()
	defer sd.Unlock()

	if ch, ok := sd.topicSubscriptions[key]; ok {
		return ch
	}

	return nil
}


func (sd *safeTopicSubscriptionDispatcher) DispatchTopicEvent(key uint64, message *zbsbe.SubscribedEvent) error {
	sd.Lock()
	defer sd.Unlock()

	if ch, ok := sd.topicSubscriptions[key]; ok {
		event := &zbsubscriptions.SubscriptionEvent{
			Event: message,
		}
		ch <- event
		return nil
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
			Mutex:&sync.Mutex{},
			taskSubscriptions: make(map[uint64]chan *zbsubscriptions.SubscriptionEvent),
		},
		safeTopicSubscriptionDispatcher: &safeTopicSubscriptionDispatcher{
			Mutex: &sync.Mutex{},
			topicSubscriptions: make(map[uint64]chan *zbsubscriptions.SubscriptionEvent),
		},
	}
}

type transactionDispatcher struct {
	*sync.RWMutex
	lastTransactionSeed uint64
	activeTransactions  []*RequestWrapper
}

func (d *transactionDispatcher) AddTransaction(request *RequestWrapper) {
	d.Lock()
	defer d.Unlock()

	d.lastTransactionSeed += zbcommon.RequestQueueSize + 1
	request.WithRequestID(d.lastTransactionSeed)

	index := d.lastTransactionSeed & (zbcommon.RequestQueueSize - 1)
	d.activeTransactions[index] = request
}

func (d *transactionDispatcher) GetTransaction(requestID uint64) *RequestWrapper {
	d.Lock()
	defer d.Unlock()

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
		RWMutex:             &sync.RWMutex{},
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
