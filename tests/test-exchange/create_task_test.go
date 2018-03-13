package testbroker

import (
	"testing"

	"github.com/zeebe-io/zbc-go/zbc"
	"github.com/zeebe-io/zbc-go/zbc/common"

	. "github.com/zeebe-io/zbc-go/tests/test-helpers"
)

func TestCreateTask(t *testing.T) {
	zbClient, err := zbc.NewClient(BrokerAddr)
	Assert(t, nil, err, true)
	Assert(t, nil, zbClient, false)

	task := zbc.NewTask("testType")
	responseTask, err := zbClient.CreateTask(TopicName, task)
	Assert(t, nil, err, true)
	Assert(t, nil, responseTask, false)

	Assert(t, nil, responseTask, false)
	Assert(t, nil, responseTask.State, false)
	Assert(t, responseTask.State, zbcommon.TaskCreated, true)
}
