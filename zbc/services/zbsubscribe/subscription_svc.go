package zbsubscribe

import (
	"github.com/zeebe-io/zbc-go/zbc/services/zbexchange"
)

type SubscriptionSvc struct {
	zbexchange.LikeExchangeSvc

	LikeTaskSubscriptionSvc
	LikeTopicSubscriptionSvc
}

func NewSubscriptionSvc(exchange zbexchange.LikeExchangeSvc) *SubscriptionSvc {
	return &SubscriptionSvc{
		LikeExchangeSvc:exchange,
		LikeTaskSubscriptionSvc: NewTaskSubscriptionSvc(exchange),
		LikeTopicSubscriptionSvc: NewTopicSubscribtionSvc(exchange),
	}
}
