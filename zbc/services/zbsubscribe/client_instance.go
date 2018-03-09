package zbsubscribe

import (
	"github.com/zeebe-io/zbc-go/zbc/models/zbsubscriptions"
	"sync"
)

type safeClient struct {
	*sync.Mutex
	ZeebeAPI
}

func (client *safeClient) GetClientInstance() ZeebeAPI {
	client.Lock()
	defer client.Unlock()

	return client.ZeebeAPI
}

var clientInstance *safeClient

func SetClientInstance(client ZeebeAPI) {
	if clientInstance == nil {
		clientInstance = &safeClient{
			Mutex:    &sync.Mutex{},
			ZeebeAPI: client,
		}
	}
	return
}

type TaskSubscriptionCallback func(ZeebeAPI, *zbsubscriptions.SubscriptionEvent)
type TopicSubscriptionCallback func(ZeebeAPI, *zbsubscriptions.SubscriptionEvent) error
