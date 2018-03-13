package test_socket

import (
	"github.com/zeebe-io/zbc-go/zbc/models/zbdispatch"
	"github.com/zeebe-io/zbc-go/zbc/services/zbsocket"
	"testing"

	. "github.com/zeebe-io/zbc-go/tests/test-helpers"
	"os"
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

	brokerAddr := os.Getenv("ZB_BROKER_ADDR")
	if len(brokerAddr) == 0 {
		Assert(t, 1, len(topology.Brokers), true)
	} else {
		Assert(t, 3, len(topology.Brokers), true)
	}

	Assert(t, nil, response, false)
	Assert(t, nil, topology, false)
}
