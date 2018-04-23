package test_topology

import (
	"testing"

	"github.com/zeebe-io/zbc-go/zbc"

	. "github.com/zeebe-io/zbc-go/tests/test-helpers"
)

func TestGetPartitions(t *testing.T) {
	zbClient, err := zbc.NewClient(BrokerAddr)
	Assert(t, nil, err, true)
	Assert(t, nil, zbClient, false)

	topology, err := zbClient.GetTopology()
	Assert(t, nil, err, true)
	Assert(t, nil, topology, false)

	partitions := topology.PartitionIDByTopicName
	Assert(t, nil, partitions, false)
	t.Logf("%+v", partitions)
}
