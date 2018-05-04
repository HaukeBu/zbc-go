package zbtopology

import (
	"github.com/zeebe-io/zbc-go/zbc/models/zbdispatch"
	"github.com/zeebe-io/zbc-go/zbc/models/zbmsgpack"
	"github.com/zeebe-io/zbc-go/zbc/services/zbsocket"
)

type LikeTopologySvc interface {
	GetTopology() (*zbmsgpack.ClusterTopology, error)
	TopicPartitionsAddrs(topic string) (*map[uint16]string, error)

	GetRoundRobinCtl() *RoundRobinCtrl
	PeekPartitionIndex(topic string) *uint16
	NextPartitionID(topic string) (*uint16, error)

	ExecuteRequest(request *zbsocket.RequestWrapper) (*zbdispatch.Message, error)
}
