package zbsubscriptions

import (
	"sync/atomic"
	"github.com/zeebe-io/zbc-go/zbc/common"
)

type SubscriptionPipelineCtrl struct {
	MessagesProcessed uint64
	OutCh             chan *SubscriptionEvent
	Partitions []uint16
}

func (sc *SubscriptionPipelineCtrl) CloseOutChannel() {
	zbcommon.ZBL.Debug().Msg("Stoping pipeline")
	close(sc.OutCh)
}

func (sc *SubscriptionPipelineCtrl) AddPartition(partitionID uint16) {
	zbcommon.ZBL.Debug().Msg("Adding partition to subscription")
	sc.Partitions = append(sc.Partitions, partitionID)
}

// GetNextMessage retrieves next message from the channel and gives us the guarantee that the returned message
// will not be null, however this routine will block the calling thread.
func (sc *SubscriptionPipelineCtrl) GetNextMessage() *SubscriptionEvent {
	atomic.AddUint64(&sc.MessagesProcessed, 1)
	zbcommon.ZBL.Debug().Str("component", "SubscriptionPipelineCtrl").Str("method", "GetNextMessage").Msgf("fetching new message")
	count := atomic.LoadUint64(&sc.MessagesProcessed)
	zbcommon.ZBL.Debug().Str("component", "SubscriptionPipelineCtrl").Str("method", "GetNextMessage").Msgf("message count :: %+v", count)
	for {
		msg, ok := <-sc.OutCh // MARK: this should not be nil
		if !ok {
			zbcommon.ZBL.Error().Msg("received nil message")
			continue
		}
		return msg
	}
}

func (ts *SubscriptionPipelineCtrl) EventsToProcess(blockSize, wantToProcess, processed uint64) uint64 {
	remaining := wantToProcess - processed
	if remaining > blockSize {
		return blockSize
	} else {
		return remaining
	}
}


func NewSubscriptionPipelineCtrl() *SubscriptionPipelineCtrl {
	return &SubscriptionPipelineCtrl{
		MessagesProcessed: 0,
		OutCh:             make(chan *SubscriptionEvent, zbcommon.RequestQueueSize),
	}
}
