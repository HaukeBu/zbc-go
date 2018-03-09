package testbroker

import (
	"testing"

	"github.com/zeebe-io/zbc-go/zbc"
	"github.com/zeebe-io/zbc-go/zbc/common"

	. "github.com/zeebe-io/zbc-go/tests/test-helpers"
)

func TestCreateWorkflow(t *testing.T) {
	zbClient, err := zbc.NewClient(BrokerAddr)
	Assert(t, nil, err, true)
	Assert(t, nil, zbClient, false)

	workflow, err := zbClient.CreateWorkflowFromFile(TopicName, zbcommon.BpmnXml, "../../examples/demoProcess.bpmn")
	Assert(t, nil, err, true)
	Assert(t, nil, workflow, false)
	Assert(t, zbcommon.DeploymentCreated, workflow.State, true)
}
