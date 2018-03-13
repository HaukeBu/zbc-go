package zbsubscribe

import (
	"math"
	"sync"

	"github.com/zeebe-io/zbc-go/zbc/common"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsubscriptions"
)

// TopicSubscriptionCloseRequest is responsible for handling of topic subscription mechanics.
type TopicSubscription struct {
	*zbsubscriptions.SubscriptionPipelineCtrl

	*TopicSubscriptionCtrl
	*TopicSubscriptionCallbackCtrl

	svc                *TopicSubscriptionSvc
	lastSucessfulEvent *zbsubscriptions.SubscriptionEvent
	*sync.Mutex
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

			zbcommon.ZBL.Debug().Str("component", "TopicSubscription").Msg("received message from socket")
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
			ts.lastSucessfulEvent = lastProcessed
			ts.Unlock()

			zbcommon.ZBL.Debug().Str("component", "TopicSubscription").Msgf("%d/%d", i, n)
		}
	}

	return lastProcessed, i, nil
}

func (ts *TopicSubscription) ProcessNext(n uint64) (bool, error) {
	bsize := float32(zbcommon.RequestQueueSize) * float32(0.01)
	var toProcess, processed, blockSize uint64 = 0, 0, uint64(bsize)
	var threshold int = int(math.Ceil(float64(float64(zbcommon.RequestQueueSize) * float64(zbcommon.TopicSubscriptionAckThreshold))))

	for {
		toProcess = ts.EventsToProcess(blockSize, n, processed)
		zbcommon.ZBL.Info().Msgf("trying to process %d events", toProcess)
		lastProcessed, processCount, err := ts.processNext(toProcess)
		zbcommon.ZBL.Info().Msgf("processed %d events", processCount)

		if err != nil {
			zbcommon.ZBL.Error().Msg("handler cannot process event")
			zbcommon.ZBL.Error().Msg("closing subscription")
			ts.Close()
			return true, zbcommon.ErrSubscriptionClosed
		}

		if lastProcessed == nil {
			return true, nil
		}

		processed += processCount

		zbcommon.ZBL.Info().Msgf("processed %d events in total", processed)
		if n <= processed {
			break
		}

		if len(ts.OutCh) < threshold { // or n > processed
			zbcommon.ZBL.Info().Msgf("acknowledging events")
			_, err := ts.svc.TopicSubscriptionAck(ts.CloseRequests[lastProcessed.Event.PartitionId].SubscriptionName, lastProcessed)
			if err != nil {
				return false, err
			}
		}
	}

	return false, nil
}

func (ts *TopicSubscription) Start() error {
	for {
		done, err := ts.ProcessNext(zbcommon.RequestQueueSize)
		if err != nil {
			return err
		}

		if done {
			return nil
		}
	}
}

func (ts *TopicSubscription) Close() []error {
	zbcommon.ZBL.Debug().Msg("Closing Subscription")
	var errors []error
	ts.Lock()
	if ts.lastSucessfulEvent != nil {
		zbcommon.ZBL.Debug().Msgf("acking last successful event %v\n", ts.lastSucessfulEvent.String())
		_, err := ts.svc.TopicSubscriptionAck(ts.CloseRequests[ts.lastSucessfulEvent.Event.PartitionId].SubscriptionName, ts.lastSucessfulEvent)
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
func NewTopicSubscription() *TopicSubscription {
	return &TopicSubscription{
		SubscriptionPipelineCtrl:      zbsubscriptions.NewSubscriptionPipelineCtrl(),
		TopicSubscriptionCallbackCtrl: nil,
		TopicSubscriptionCtrl:         NewTopicSubscriptionAckCtrl(),
		Mutex: &sync.Mutex{},
	}
}
