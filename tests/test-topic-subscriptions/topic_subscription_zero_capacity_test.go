package test_topic_subscriptions

import (
	"testing"

	"github.com/zeebe-io/zbc-go/zbc/services/zbsubscribe"

	"github.com/zeebe-io/zbc-go/zbc"
	"github.com/zeebe-io/zbc-go/zbc/common"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsubscriptions"

	. "github.com/zeebe-io/zbc-go/tests/test-helpers"
)

func TestTopicSubscriptionZeroCapacity(t *testing.T) {
	t.Log("Creating client")
	zbClient, err := zbc.NewClient(BrokerAddr)
	Assert(t, nil, err, true)
	Assert(t, nil, zbClient, false)
	t.Log("Client created")

	NumberOfPartitions = 3
	hash := CreateRandomTopicWithTimeout(t, zbClient)

	t.Log("Creating workflow")
	workflow, err := zbClient.CreateWorkflowFromFile(hash, zbcommon.BpmnXml, "../../examples/demoProcess.bpmn")
	Assert(t, nil, err, true)
	Assert(t, nil, workflow, false)
	Assert(t, nil, workflow.State, false)
	Assert(t, zbcommon.DeploymentCreated, workflow.State, true)
	t.Logf("Workflow created")

	payload := make(map[string]interface{})
	payload["a"] = "b"

	instance := zbc.NewWorkflowInstance("demoProcess", -1, payload)
	createdInstance, err := zbClient.CreateWorkflowInstance(hash, instance)
	Assert(t, nil, err, true)
	Assert(t, nil, createdInstance, false)
	Assert(t, zbcommon.WorkflowInstanceCreated, createdInstance.State, true)

	subscription, err := zbClient.TopicSubscription(hash, "callback-test", 0, 0, false,
		func(client zbsubscribe.ZeebeAPI, event *zbsubscriptions.SubscriptionEvent) error { return nil })

	Assert(t, nil, subscription, false)
	Assert(t, zbcommon.ErrZeroCapacity, err, true)

}
