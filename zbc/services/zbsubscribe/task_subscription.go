package zbsubscribe

import (
	"github.com/zeebe-io/zbc-go/zbc/models/zbsubscriptions"
	"github.com/zeebe-io/zbc-go/zbc/common"
)

// TaskSubscription is object for handling overall test-task-subscriptions subscription.
type TaskSubscription struct {
	*zbsubscriptions.SubscriptionPipelineCtrl

	*TaskSubscriptionCtrl
	*TaskSubscriptionCallbackCtrl

	svc *TaskSubscriptionSvc
}

func (ts *TaskSubscription) totalCredits() uint64 {
	var total uint64 = 0
	for _, p := range ts.Partitions {
		total += uint64(ts.Subscriptions[p].Credits)
	}
	return total
}

func (ts *TaskSubscription) increaseCredits() {
	for _, sub := range ts.Subscriptions {
		s := sub
		_, creditsErr := ts.svc.IncreaseTaskSubscriptionCredits(s)
		if creditsErr != nil {
			zbcommon.ZBL.Error().Msg("increasing credits failed on partition")
		}
	}
}

func (ts *TaskSubscription) processNext(n uint64) {
	var i uint64 = 0
	for ; i < n; i++ {
		select {
		case msg := <-ts.OutCh:
			ts.ExecuteCallback(msg)
		}
	}
}

func (ts *TaskSubscription) ProcessNext(n uint64) {
	var toProcess, processed uint64 = 0, 0
	totalCredits := uint64(ts.Subscriptions[ts.Partitions[0]].Credits)

	threshold := float32(zbcommon.RequestQueueSize) * zbcommon.TaskSubscriptionRefreshCreditsThreshold
	for {
		toProcess = ts.EventsToProcess(totalCredits, n, processed)
		ts.processNext(toProcess)
		processed += toProcess

		zbcommon.ZBL.Info().Msgf("processed %d tasks", processed)

		if n <= processed {
			break
		}
		if len(ts.OutCh) < int(threshold)  { // or n > processed
			zbcommon.ZBL.Info().Msg("increasing credits on all partitions")
			ts.increaseCredits()
		}
	}
}

func (ts *TaskSubscription) Start() {
	for {
		ts.ProcessNext(zbcommon.RequestQueueSize)
	}
}

func (ts *TaskSubscription) Close() []error {
	return ts.svc.CloseTaskSubscription(ts) //clientInstance.CloseTaskSubscription(ts)
}

func (ts *TaskSubscription) WithTaskSubscriptionSvc(svc *TaskSubscriptionSvc) *TaskSubscription {
	ts.svc = svc
	return ts
}

func (ts *TaskSubscription) WithCallback(cb TaskSubscriptionCallback) *TaskSubscription {
	ts.TaskSubscriptionCallbackCtrl = NewTaskSubscriptionCallbackCtrl(cb)
	return ts
}

// NewTaskSubscription is constructor for TaskSubscription object.
func NewTaskSubscription() *TaskSubscription {
	return &TaskSubscription{
		SubscriptionPipelineCtrl:     zbsubscriptions.NewSubscriptionPipelineCtrl(),
		TaskSubscriptionCallbackCtrl: nil,
		TaskSubscriptionCtrl:         NewTaskSubscriptionCreditsCtrl(),
	}
}
