package zbsubscriptions

import (
	"sync"
)

type SubscriptionPipelineCtrl struct {
	outCh      chan *SubscriptionEvent
	consumerCh map[uint16]chan *SubscriptionEvent
}


func (sc *SubscriptionPipelineCtrl) StopPipeline() {
	for _, ch := range sc.consumerCh {
		close(ch)
	}
	close(sc.outCh)
}

func (sc *SubscriptionPipelineCtrl) StartPipeline() {
	var wg sync.WaitGroup
	out := make(chan *SubscriptionEvent)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan *SubscriptionEvent) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}

	wg.Add(len(sc.consumerCh))
	for _, c := range sc.consumerCh {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()

	sc.outCh = out
}

func (sc *SubscriptionPipelineCtrl) AddPartitionChannel(partitionID uint16, ch chan *SubscriptionEvent) {
	sc.consumerCh[partitionID] = ch
}


func (sc *SubscriptionPipelineCtrl) GetNextMessage() *SubscriptionEvent {
	return <-sc.outCh
}


func NewSubscriptionPipelineCtrl() *SubscriptionPipelineCtrl {
	return &SubscriptionPipelineCtrl{
		outCh:      make(chan *SubscriptionEvent),
		consumerCh: make(map[uint16]chan *SubscriptionEvent),
	}
}
