package test_helpers

import (
	"github.com/zeebe-io/zbc-go/zbc/common"
	"testing"
)

func TestRandStringBytes(t *testing.T) {
}

func TestAssert(t *testing.T) {
}

func TestIsIPv4(t *testing.T) {
	result := zbcommon.IsIPv4("127.0.0.1:51015")
	Assert(t, true, result, true)
}
