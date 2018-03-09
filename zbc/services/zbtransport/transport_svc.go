package zbtransport

import (
	"github.com/zeebe-io/zbc-go/zbc/common"
	"github.com/zeebe-io/zbc-go/zbc/models/zbdispatch"
	"github.com/zeebe-io/zbc-go/zbc/services/zbsocket"
)

// TransportSvc is responsible for executing requests on right socket.
type TransportSvc struct {
	workQueue   chan *zbsocket.RequestWrapper
	connections map[string]zbsocket.LikeSocket
}

// ExecTransport will enqueue request to transport workload queue.
func (tm *TransportSvc) ExecTransport(request *zbsocket.RequestWrapper) {
	tm.workQueue <- request
}

// TODO: check if socket in map is unreachable/unhealthy than delete that record from the map
// TODO: keep-alive messages
func (tm *TransportSvc) getSocket(addr string) (zbsocket.LikeSocket, error) {
	if conn, ok := tm.connections[addr]; ok && conn != nil {
		return conn, nil
	}

	sock := zbsocket.NewSocket(addr)
	err := sock.Dial(addr)

	if err != nil {
		return nil, err
	}

	tm.connections[addr] = sock
	return sock, nil
}

func (tm *TransportSvc) transportWorker() {
	for {
		select {
		case request := <-tm.workQueue:
			sock, err := tm.getSocket(request.Addr)
			if err != nil {
				request.ErrorCh <- err
				continue
			}

			MessageRetry(func() (*zbdispatch.Message, error) {
				sock.AddTransaction(request)
				if err := sock.Send(request); err != nil {
					return nil, err
				}
				return nil, nil
			})

		}
	}
}

// NewTransportSvc will construct new TransportSvc object.
func NewTransportSvc() *TransportSvc {
	tm := &TransportSvc{
		workQueue:   make(chan *zbsocket.RequestWrapper, zbcommon.RequestQueueSize),
		connections: make(map[string]zbsocket.LikeSocket),
	}

	go tm.transportWorker()
	return tm
}
