package test_socket

import (
	"github.com/zeebe-io/zbc-go/zbc/models/zbdispatch"
	"github.com/zeebe-io/zbc-go/zbc/services/zbsocket"
	"testing"

	. "github.com/zeebe-io/zbc-go/tests/test-helpers"
)

func TestNewSocket(t *testing.T) {
	socket := zbsocket.NewSocket(BrokerAddr)

	requestFactory := zbdispatch.RequestFactory{}
	responseHandler := zbdispatch.ResponseHandler{}
	request := zbsocket.NewRequestWrapper(requestFactory.TopologyRequest())

	err := socket.Dial(BrokerAddr)
	Assert(t, nil, err, true)

	err = socket.Send(request)
	Assert(t, nil, err, true)

	response := <-request.ResponseCh
	topology := responseHandler.UnmarshalTopology(response)

	Assert(t, 1, len(topology.Brokers), true)
	Assert(t, nil, response, false)
	Assert(t, nil, topology, false)
}
