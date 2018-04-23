package test_topology

import (
	. "github.com/zeebe-io/zbc-go/tests/test-helpers"
	"github.com/zeebe-io/zbc-go/zbc/services/zbtopology"
	"testing"
)

func TestTopologyInit(t *testing.T) {
	svc, err := zbtopology.NewTopologySvc(BrokerAddr)
	Assert(t, nil, svc, false)
	Assert(t, nil, err, true)

	topology, err := svc.GetTopology()
	Assert(t, nil, topology, false)
	Assert(t, nil, err, true)
}
