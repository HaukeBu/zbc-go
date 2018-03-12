package zbsubscribe

import (
	"github.com/zeebe-io/zbc-go/zbc/common"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsubscriptions"
	"math"
)

// TopicSubscriptionCloseRequest is responsible for handling of topic subscription mechanics.
type TopicSubscription struct {
	*zbsubscriptions.SubscriptionPipelineCtrl

	*TopicSubscriptionCtrl
	*TopicSubscriptionCallbackCtrl

	svc *TopicSubscriptionSvc
}

func (ts *TopicSubscription) processNext(n uint64) *zbsubscriptions.SubscriptionEvent {
	var i, errCount uint64 = 0, 0
	var lastProcessed *zbsubscriptions.SubscriptionEvent
	for ; i < n; i++ {
		select {
		case msg := <-ts.OutCh:
			zbcommon.ZBL.Debug().Str("component", "TopicSubscription").Msg("received message from socket")
			for ; errCount < 3; errCount++ {
				err := ts.ExecuteCallback(msg)
				if err != nil && errCount == 2 {
					return nil
				}
			}
			lastProcessed = msg
			zbcommon.ZBL.Debug().Str("component", "TopicSubscription").Msgf("%d/%d", i, n)
		}
	}
	return lastProcessed
}

func (ts *TopicSubscription) ProcessNext(n uint64) error {
	bsize := float32(zbcommon.RequestQueueSize) * float32(0.01)
	var toProcess, processed, blockSize uint64 = 0, 0, uint64(bsize)
	var threshold int = int(math.Ceil(float64(float64(zbcommon.RequestQueueSize) * float64(zbcommon.TopicSubscriptionAckThreshold))))
	var lastProcessed, lastSuccessfulProcessed *zbsubscriptions.SubscriptionEvent

	for {
		toProcess = ts.EventsToProcess(blockSize, n, processed)
		zbcommon.ZBL.Info().Msgf("trying to process %d events", toProcess)
		lastProcessed = ts.processNext(toProcess)
		zbcommon.ZBL.Info().Msgf("processed %d events", toProcess)

		if lastProcessed == nil {
			zbcommon.ZBL.Error().Msg("handler cannot process event")
			if lastSuccessfulProcessed != nil {
				zbcommon.ZBL.Error().Msg("acking last successful event")
				_, err := ts.svc.TopicSubscriptionAck(ts.CloseRequests[lastProcessed.Event.PartitionId].SubscriptionName, lastSuccessfulProcessed)
				if err != nil {
					return err
				}
			}
			zbcommon.ZBL.Error().Msg("closing subscription")
			ts.Close()
			return zbcommon.ErrSubscriptionClosed
		}
		processed += toProcess
		lastSuccessfulProcessed = lastProcessed

		zbcommon.ZBL.Info().Msgf("processed %d tasks", processed)
		if n <= processed {
			break
		}

		if len(ts.OutCh) < threshold { // or n > processed
			zbcommon.ZBL.Info().Msgf("acknowledging events")
			_, err := ts.svc.TopicSubscriptionAck(ts.CloseRequests[lastProcessed.Event.PartitionId].SubscriptionName, lastProcessed)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (ts *TopicSubscription) Start() error {
	for {
		err := ts.ProcessNext(zbcommon.RequestQueueSize)
		if err != nil {
			return err
		}
	}
}

func (ts *TopicSubscription) Close() []error {
	err := ts.svc.CloseTopicSubscription(ts)
	ts.CloseOutChannel()
	return err
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
	}
}
