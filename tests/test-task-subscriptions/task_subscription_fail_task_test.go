package test_task_subscriptions

import (
	"testing"

	"github.com/zeebe-io/zbc-go/zbc"
	"github.com/zeebe-io/zbc-go/zbc/common"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsubscriptions"
	"github.com/zeebe-io/zbc-go/zbc/services/zbsubscribe"

	. "github.com/zeebe-io/zbc-go/tests/test-helpers"
	"sync/atomic"
	"time"
)

func TestTaskSubscriptionFailTask(t *testing.T) {
	t.Log("creating new client")
	zbClient, err := zbc.NewClient(BrokerAddr)
	Assert(t, nil, err, true)
	Assert(t, nil, zbClient, false)

	workflow, err := zbClient.CreateWorkflowFromFile(TopicName, zbcommon.BpmnXml, "../../examples/demoProcess.bpmn")
	Assert(t, nil, err, true)
	Assert(t, nil, workflow, false)
	Assert(t, nil, workflow.State, false)
	Assert(t, zbcommon.DeploymentCreated, workflow.State, true)

	payload := make(map[string]interface{})
	payload["a"] = "b"

	instance := zbc.NewWorkflowInstance("demoProcess", -1, payload)
	t.Log("Creating 1 workflow instances")
	createdInstance, err := zbClient.CreateWorkflowInstance(TopicName, instance)
	Assert(t, nil, err, true)
	Assert(t, nil, createdInstance, false)
	Assert(t, zbcommon.WorkflowInstanceCreated, createdInstance.State, true)

	var ops uint64
	subscription, err := zbClient.TaskSubscription(TopicName, "task_subscription_test", "foo",
		func(client zbsubscribe.ZeebeAPI, event *zbsubscriptions.SubscriptionEvent) {
			Assert(t, nil, event, false)
			atomic.AddUint64(&ops, 1)

			task, err := zbClient.FailTask(event)
			Assert(t, nil, task, false)
			Assert(t, nil, err, true)
		})

	for {
		op := atomic.LoadUint64(&ops)
		t.Log("Subscription processed tasks ", op)

		if op >= 1 {
			errs := zbClient.CloseTaskSubscription(subscription)
			Assert(t, 0, len(errs), true)
			break
		}
		time.Sleep(time.Duration(1 * time.Second))
	}
}
