package zbsubscribe

import (
	"github.com/zeebe-io/zbc-go/zbc/common"
	"github.com/zeebe-io/zbc-go/zbc/models/zbmsgpack"
	"github.com/zeebe-io/zbc-go/zbc/services/zbexchange"
)

type TopicSubscriptionSvc struct {
	zbexchange.LikeExchangeSvc
}

func (ts *TopicSubscriptionSvc) topicConsumer(topic, subName string, startPosition int64, forceStart bool, prefetchCapacity int32) (*TopicSubscription, error) {
	partitions, err := ts.TopicPartitionsAddrs(topic)
	if err != nil {
		return nil, err
	}

	if partitions == nil || len(*partitions) == 0 {
		zbcommon.ZBL.Error().Msgf("topic %s topology not found", topic)
		return nil, zbcommon.BrokerNotFound
	}

	topicSubscription := NewTopicSubscription()
	zbcommon.ZBL.Debug().Msg("new topic subscription created")

	for partitionID := range *partitions {
		sub := ts.OpenTopicPartition(topicSubscription.OutCh, partitionID, topic, subName, startPosition, forceStart, prefetchCapacity)
		if sub != nil {
			closeRequest := &zbmsgpack.TopicSubscriptionCloseRequest{
				TopicName:        topic,
				PartitionID:      partitionID,
				SubscriberKey:    sub.SubscriberKey,
				SubscriptionName: subName,
			}

			topicSubscription.AddCloseRequest(partitionID, closeRequest)
			topicSubscription.AddSubscriptionInfo(partitionID, sub)
			topicSubscription.AddPartition(partitionID)
		} else {
			zbcommon.ZBL.Debug().Msg("there was an error during opening topic subscription, closing all partitions")
			topicSubscription.Close()
			return nil, zbcommon.ErrSubscriptionPipelineFailed
		}
	}

	return topicSubscription, nil
}

func (ts *TopicSubscriptionSvc) TopicSubscription(topic, subName string, prefetchCapacity int32, startPosition int64, forceStart bool, cb TopicSubscriptionCallback) (*TopicSubscription, error) {
	zbcommon.ZBL.Debug().Msg("Opening topic subscription")
	topicConsumer, err := ts.topicConsumer(topic, subName, startPosition, forceStart, prefetchCapacity)
	zbcommon.ZBL.Debug().Msg("CreateTopic subscription opened.")

	if err != nil {
		return nil, err
	}
	topicConsumer = topicConsumer.WithCallback(cb).WithTopicSubscriptionSvc(ts)

	return topicConsumer, err
}

func (ts *TopicSubscriptionSvc) CloseTopicSubscription(sub *TopicSubscription) []error {
	var errs []error
	for _, request := range sub.CloseRequests {
		_, err := ts.CloseTopicSubscriptionPartition(request)
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
