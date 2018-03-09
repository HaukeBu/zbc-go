package zbsubscribe


import (
	"github.com/zeebe-io/zbc-go/zbc/models/zbsubscriptions"
)

// TopicSubscription is responsible for handling of topic subscription mechanics.
type TopicSubscription struct {
	*zbsubscriptions.SubscriptionPipelineCtrl

	*TopicSubscriptionCallbackCtrl
	*TopicSubscriptionAckCtrl
}

func (ts *TopicSubscription) StartCallbackHandler(svc *TopicSubscriptionSvc) {
	go func(ts *TopicSubscription) {
		for {
			msg := ts.GetNextMessage()
			for errCount := 0; errCount < 3; errCount++ {
				err := ts.ExecuteCallback(msg)
				if err != nil && errCount >= 3 {
					ts.StopPipeline()
					return
				}
				if err != nil {
					continue
				}

				break
			}

			ts.IncrementMessageCount()
			if ts.GetMessageCount()%3 == 0 {
				subInfo, ok := ts.GetSubscription(msg.Event.PartitionId)
				if ok {
					svc.TopicSubscriptionAck(subInfo, msg)
				}
			}
			ts.NormalizeMessageCount()

		}
	}(ts)
}

// WithCallback is builder to set callback handler to topic subscription.
func (ts *TopicSubscription) WithCallback(cb TopicSubscriptionCallback) *TopicSubscription {
	ts.TopicSubscriptionCallbackCtrl = NewTopicSubscriptionCallbackCtrl(cb)
	return ts
}

// NewTopicSubscription is a constructor for TopicSubscription object.
func NewTopicSubscription() *TopicSubscription {
	return &TopicSubscription{
		SubscriptionPipelineCtrl:      zbsubscriptions.NewSubscriptionPipelineCtrl(),
		TopicSubscriptionCallbackCtrl: nil,
		TopicSubscriptionAckCtrl:      NewTopicSubscriptionAckCtrl(),
	}
}
