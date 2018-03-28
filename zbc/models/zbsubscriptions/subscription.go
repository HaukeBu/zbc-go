package zbsubscriptions

import (
	"github.com/zeebe-io/zbc-go/zbc/common"
)

type SubscriptionPipelineCtrl struct {
	MessagesProcessed uint64
	OutCh             chan *SubscriptionEvent
	Partitions        []uint16
}

func (sc *SubscriptionPipelineCtrl) CloseOutChannel() {
	defer func() {
		if r := recover(); r != nil {
			zbcommon.ZBL.Error().Msg("Recovered in close out channel")
		}
	}()

	zbcommon.ZBL.Debug().Msg("Stoping pipeline")
	close(sc.OutCh)
}

func (sc *SubscriptionPipelineCtrl) AddPartition(partitionID uint16) {
	zbcommon.ZBL.Debug().Msg("Adding partition to subscription")
	sc.Partitions = append(sc.Partitions, partitionID)
}

func (ts *SubscriptionPipelineCtrl) EventsToProcess(blockSize, wantToProcess, processed uint64) uint64 {
	remaining := wantToProcess - processed
	if remaining > blockSize {
		return blockSize
	} else {
		return remaining
	}
}

func NewSizedSubscriptionPipelineCtrl(channelSize uint64) *SubscriptionPipelineCtrl {
	return &SubscriptionPipelineCtrl{
		MessagesProcessed: 0,
		OutCh:             make(chan *SubscriptionEvent, channelSize),
	}
}
