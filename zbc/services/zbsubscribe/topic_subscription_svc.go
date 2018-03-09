package zbsubscribe


import (
	"github.com/zeebe-io/zbc-go/zbc/services/zbexchange"
	"github.com/zeebe-io/zbc-go/zbc/models/zbmsgpack"
	"github.com/zeebe-io/zbc-go/zbc/common"
)

type TopicSubscriptionSvc struct {
	zbexchange.LikeExchangeSvc
}

func (ts *TopicSubscriptionSvc) topicConsumer(topic, subName string, startPosition int64) (*TopicSubscription, error) {
	partitions, err := ts.TopicPartitionsAddrs(topic)
	if err != nil {
		return nil, err
	}
	topicSubscription := NewTopicSubscription()

	for partitionID := range *partitions {
		sub, ch := ts.OpenTopicPartition(partitionID, topic, subName, startPosition)
		if sub != nil && ch != nil {
			closeableSub := &zbmsgpack.TopicSubscription{
				TopicName:        topic,
				PartitionID:      partitionID,
				SubscriberKey:    sub.SubscriberKey,
				SubscriptionName: subName,
			}

			topicSubscription.AddSubscription(partitionID, closeableSub)
			topicSubscription.AddPartitionChannel(partitionID, ch)
		} else {
			zbcommon.ZBL.Debug().Msg("there was an error during opening topic subscription")
			// TODO: send close to all previous subs
			return nil, zbcommon.ErrSubscriptionPipelineFailed
		}
	}

	topicSubscription.StartPipeline()
	return topicSubscription, nil
}

func (ts *TopicSubscriptionSvc) TopicSubscription(topic, subName string, startPosition int64, cb TopicSubscriptionCallback) (*TopicSubscription, error) {
	zbcommon.ZBL.Debug().Msg("Opening topic subscription")
	topicConsumer, err := ts.topicConsumer(topic, subName, startPosition)
	zbcommon.ZBL.Debug().Msg("Topic subscription opened.")

	if err != nil {
		return nil, err
	}
	topicConsumer = topicConsumer.WithCallback(cb)
	topicConsumer.StartCallbackHandler(ts)

	return topicConsumer, err
}

func (ts *TopicSubscriptionSvc) CloseTopicSubscription(sub *TopicSubscription) []error {
	var errs []error
	for _, taskSub := range sub.Subscriptions {
		_, err := ts.CloseTopicSubscriptionPartition(taskSub)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func NewTopicSubscribtionSvc(exchange zbexchange.LikeExchangeSvc) *TopicSubscriptionSvc {
	return &TopicSubscriptionSvc{
		exchange,
	}
}
