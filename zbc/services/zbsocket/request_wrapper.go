package zbsocket

import (
	"github.com/zeebe-io/zbc-go/zbc/models/zbdispatch"
)

type RequestWrapper struct {
	Sock LikeSocket

	RequestID  uint64
	Addr       string
	ResponseCh chan *zbdispatch.Message
	ErrorCh    chan error
	Payload    *zbdispatch.Message
}

func (rw *RequestWrapper) ExecutedWithSocket(socket LikeSocket) {
	rw.Sock = socket
}

func (rw *RequestWrapper) WithRequestID(id uint64) *RequestWrapper {
	rw.RequestID = id
	rw.Payload.Headers.RequestResponseHeader.RequestID = id
	return rw
}

func NewRequestWrapper(payload *zbdispatch.Message) *RequestWrapper {
	return &RequestWrapper{
		Addr:       "",
		ResponseCh: make(chan *zbdispatch.Message),
		ErrorCh:    make(chan error),
		Payload:    payload,
	}
}
