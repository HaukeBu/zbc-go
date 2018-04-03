package zbsubscribe

import (
	"github.com/zeebe-io/zbc-go/zbc/common"
	"github.com/zeebe-io/zbc-go/zbc/models/zbmsgpack"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsubscriptions"
	"sync"
)

// TaskSubscription is object for handling overall test-task-subscriptions subscription.
type TaskSubscription struct {
	*zbsubscriptions.SubscriptionPipelineCtrl

	*TaskSubscriptionCtrl
	*TaskSubscriptionCallbackCtrl

	svc *TaskSubscriptionSvc

	credits *sync.Map
}

// totalCredits returns len(partitions) * subscription.Credits
// in case of no partitions present or partitionID is not found in subscription information it will return 0
func (ts *TaskSubscription) totalCredits() uint64 {
	if len(ts.Partitions) == 0 {
		return 0
	}
	if sub, ok := ts.GetSubscription(ts.Partitions[0]); ok {
		return uint64(sub.Credits) * uint64(len(ts.Partitions))
	}
	return 0
}

func (ts *TaskSubscription) increaseCredits() {
	zbcommon.ZBL.Info().
		Str("component", "TaskSubscription").
		Msg("increasing credits on all partitions")

	ts.Subscriptions.Range(func(key, value interface{}) bool {
		s := value.(*zbmsgpack.TaskSubscriptionInfo)

		for i := 0; i < 3; i++ {
			_, creditsErr := ts.svc.IncreaseTaskSubscriptionCredits(s)
			if creditsErr != nil {
				zbcommon.ZBL.Error().Msg("increasing credits failed on partition")
				continue
			}
			return true
		}

		ts.Close()
		return false
	})

	zbcommon.ZBL.Info().
		Str("component", "TaskSubscription").
		Msg("credits successfully increased")

}

func (ts *TaskSubscription) processNext(n uint64) uint64 {
	zbcommon.ZBL.Info().
		Str("component", "TaskSubscription").
		Str("method", "processNext").
		Msgf("About to process %d events", n)

	var i uint64 = 0
	for ; i < n; i++ {
		select {
		case msg, ok := <-ts.OutCh:
			if !ok {
				zbcommon.ZBL.Debug().
					Str("component", "TaskSubscription").
					Msgf("cannot read from OutCh(%d/%d)", len(ts.OutCh), cap(ts.OutCh))

				return i
			}

			ts.consumeCredit(msg.Event.PartitionId)
			ts.ExecuteCallback(msg)
		}
	}

	zbcommon.ZBL.Debug().
		Str("component", "TaskSubscription").
		Msgf("finished processing %d events", n)

	return n
}

func (ts *TaskSubscription) initCredits(size int32) {
	for _, partitionID := range ts.Partitions {
		ts.credits.Store(partitionID, size)
	}
}

func (ts *TaskSubscription) consumeCredit(partitionID uint16) {
	credits, ok := ts.credits.Load(partitionID)
	if ok {
		c := credits.(int32)
		c--
		ts.credits.Store(partitionID, c)
	}
}

func (ts *TaskSubscription) checkCredits(creditsSize int32) {
	ts.credits.Range(func(key, value interface{}) bool {
		partitionID, credits := key.(uint16), value.(int32)
		zbcommon.ZBL.Debug().Str("component", "TaskSubscription").Msgf("partition %d = %d credits", partitionID, credits)
		if float32(credits) <= float32(creditsSize)*zbcommon.TaskSubscriptionRefreshCreditsThreshold {
			s, ok := ts.Subscriptions.Load(partitionID)
			if ok {
				sub := s.(*zbmsgpack.TaskSubscriptionInfo)
				for i := 0; i < 3; i++ {
					_, creditsErr := ts.svc.IncreaseTaskSubscriptionCredits(sub)
					if creditsErr != nil {
						zbcommon.ZBL.Error().Msg("increasing credits failed on partition")
						continue
					}
					break
				}
				ts.credits.Store(partitionID, sub.Credits)
			}
		}

		// iterate over all partitions in credits map
		return true
	})
}

func (ts *TaskSubscription) ProcessNext(n uint64) error {
	var toProcess, processed uint64 = 0, 0
	sub, ok := ts.GetSubscription(ts.Partitions[0])
	if !ok {
		return zbcommon.SubscriptionIsClosed
	}

	bsize := uint64(sub.Credits)
	for {
		toProcess = ts.EventsToProcess(bsize, n, processed)
		p := ts.processNext(toProcess)
		processed += p

		ts.checkCredits(sub.Credits)
		zbcommon.ZBL.Info().Msgf("processed %d tasks", processed)
		if n <= processed || p < toProcess {
			break
		}
	}
	return nil
}

func (ts *TaskSubscription) Start() {
	for {
		ts.ProcessNext(uint64(cap(ts.OutCh)))
	}
}

func (ts *TaskSubscription) Close() []error {
	return ts.svc.CloseTaskSubscription(ts)
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
func NewTaskSubscription(channelSize uint64) *TaskSubscription {
	subscription := &TaskSubscription{
		SubscriptionPipelineCtrl:     zbsubscriptions.NewSizedSubscriptionPipelineCtrl(channelSize),
		TaskSubscriptionCallbackCtrl: nil,
		TaskSubscriptionCtrl:         NewTaskSubscriptionCreditsCtrl(),
		credits:                      &sync.Map{},
	}
	return subscription
}
