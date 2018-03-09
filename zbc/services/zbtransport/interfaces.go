package zbtransport

import (
	"github.com/zeebe-io/zbc-go/zbc/services/zbsocket"
)

type LikeTransportSvc interface {
	ExecTransport(request *zbsocket.RequestWrapper)
}
