package zbsubscribe

import (
	"sync"

	"github.com/zeebe-io/zbc-go/zbc/common"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsubscriptions"
)

// TopicSubscriptionCloseRequest is responsible for handling of topic subscription mechanics.
type TopicSubscription struct {
	*sync.Mutex
	*zbsubscriptions.SubscriptionPipelineCtrl

	*TopicSubscriptionCtrl
	*TopicSubscriptionCallbackCtrl

	svc *TopicSubscriptionSvc

	lastSuccessful map[uint16]*zbsubscriptions.SubscriptionEvent
}

func (ts *TopicSubscription) processNext(n uint64) (*zbsubscriptions.SubscriptionEvent, uint64, error) {
	var i, errCount uint64 = 0, 0
	var lastProcessed *zbsubscriptions.SubscriptionEvent

	for ; i < n; i++ {
		select {
		case msg, ok := <-ts.OutCh:
			if !ok {
				return nil, i, nil
			}

			zbcommon.ZBL.Debug().
				Str("component", "TopicSubscription").
				Msgf("received message from OutCh(%d/%d)", len(ts.OutCh), cap(ts.OutCh))

			errCount = 0
			for ; errCount < 3; errCount++ {
				err := ts.ExecuteCallback(msg)
				if err != nil && errCount == 2 {
					return nil, i, err
				}

				if err == nil {
					break
				}
			}
			lastProcessed = msg

			ts.Lock()
			ts.lastSuccessful[lastProcessed.Event.PartitionId] = lastProcessed
			ts.Unlock()
		}
	}

	return lastProcessed, i, nil
}

func (ts *TopicSubscription) ProcessNext(n uint64) (bool, error) {
	if n == 0 {
		return true, zbcommon.ErrCannotProcessZeroEvents
	}

	bsize := ts.SubscriptionsInfo[ts.Partitions[0]].PrefetchCapacity - 1
	var toProcess, processed, blockSize uint64 = 0, 0, uint64(bsize)

	for {
		zbcommon.ZBL.Info().Str("component", "TopicSubscription").Msgf("processed (%d/%d)", processed, n)
		if n <= processed {
			zbcommon.ZBL.Debug().Str("component", "TopicSubscription").Msgf("n is less than processed (%d/%d) exiting", n, processed)
			break
		}

		toProcess = ts.EventsToProcess(blockSize, n, processed)
		zbcommon.ZBL.Info().Str("component", "TopicSubscription").Msgf("trying to process %d events", toProcess)
		lastProcessed, processCount, err := ts.processNext(toProcess)

		if err != nil {
			zbcommon.ZBL.Error().Str("component", "TopicSubscription").Msg("handler cannot process event")
			zbcommon.ZBL.Error().Str("component", "TopicSubscription").Msg("closing subscription")
			ts.Close()
			return true, zbcommon.ErrSubscriptionClosed
		}

		if lastProcessed == nil {
			return true, nil
		}

		processed += processCount

		for partitionID, lastSuccessful := range ts.lastSuccessful {
			zbcommon.ZBL.Info().Str("component", "TopicSubscription").Msgf("acknowledging events for partition %d", partitionID)
			_, err = ts.svc.TopicSubscriptionAck(ts.CloseRequests[lastSuccessful.Event.PartitionId].SubscriptionName, lastSuccessful)
			if err != nil {
				return false, err
			}
		}

	}

	return false, nil
}

func (ts *TopicSubscription) Start() []error {
	for {
		done, err := ts.ProcessNext(uint64(cap(ts.OutCh)))
		if err != nil {
			errs := ts.Close()
			errs = append(errs, err)
			return errs
		}

		if done {
			return nil
		}
	}
}

func (ts *TopicSubscription) Close() []error {
	zbcommon.ZBL.Debug().Str("component", "TopicSubscription").Msg("Closing Subscription")
	var errors []error
	ts.Lock()
	for partitionID, lastSuccessful := range ts.lastSuccessful {
		zbcommon.ZBL.Debug().Str("component", "TopicSubscription").Msgf("acking last successful event %v for partitionID %d\n", lastSuccessful.String(), partitionID)
		_, err := ts.svc.TopicSubscriptionAck(ts.CloseRequests[lastSuccessful.Event.PartitionId].SubscriptionName, lastSuccessful)
		if err != nil {
			errors = append(errors, err)
		}
	}
	ts.Unlock()

	err := ts.svc.CloseTopicSubscription(ts)
	errors = append(errors, err...)

	ts.CloseOutChannel()
	return errors
}

func (ts *TopicSubscription) WithTopicSubscriptionSvc(svc *TopicSubscriptionSvc) *TopicSubscription {
	ts.svc = svc
	return ts
}

// WithCallback is builder to set callback handler to topic subscription.
func (ts *TopicSubscription) WithCallback(cb TopicSubscriptionCallback) *TopicSubscription {
	ts.TopicSubscriptionCallbackCtrl = NewTopicSubscriptionCallbackCtrl(cb)
	return ts
}

// NewTopicSubscription is a constructor for TopicSubscriptionCloseRequest object.
func NewTopicSubscription(channelSize uint64) *TopicSubscription {
	return &TopicSubscription{
		Mutex: &sync.Mutex{},
		SubscriptionPipelineCtrl:      zbsubscriptions.NewSizedSubscriptionPipelineCtrl(channelSize),
		TopicSubscriptionCallbackCtrl: nil,
		TopicSubscriptionCtrl:         NewTopicSubscriptionAckCtrl(),
		lastSuccessful:                make(map[uint16]*zbsubscriptions.SubscriptionEvent),
	}
}
