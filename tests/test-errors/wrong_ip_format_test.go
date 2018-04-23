package test_errors

import (
	"testing"
	"github.com/zeebe-io/zbc-go/zbc"
	. "github.com/zeebe-io/zbc-go/tests/test-helpers"
	"github.com/zeebe-io/zbc-go/zbc/common"
)

func TestWrongIPFormat(t *testing.T) {
	t.Log("Creating client")
	zbClient, err := zbc.NewClient("not.a.broker.address:51015")
	Assert(t, zbcommon.ErrWrongIPFormat, err, true)
	Assert(t, nil, zbClient, false)
}

func TestWrongIPv6Format(t *testing.T) {
	t.Log("Creating client")
	zbClient, err := zbc.NewClient("not.a.ipv6.address:51015")
	Assert(t, zbcommon.ErrWrongIPFormat, err, true)
	Assert(t, nil, zbClient, false)
}
