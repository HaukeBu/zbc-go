package testbroker

import (
	"testing"

	"github.com/zeebe-io/zbc-go/zbc"

	. "github.com/zeebe-io/zbc-go/tests/test-helpers"
)

func TestTopology(t *testing.T) {
	zbClient, err := zbc.NewClient(BrokerAddr)
	Assert(t, nil, err, true)
	Assert(t, nil, zbClient, false)
}
