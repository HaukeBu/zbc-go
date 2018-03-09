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

var taskTopic *testing.T
var ops uint64

func TaskHandler(client zbsubscribe.ZeebeAPI, event *zbsubscriptions.SubscriptionEvent) {
	Assert(taskTopic, client, nil, false)
	Assert(taskTopic, event, nil, false)
	atomic.AddUint64(&ops, 1)
}

func TestTaskSubscriptionUnscopedCallback(t *testing.T) {
	taskTopic = t

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

	wfStart := time.Now()
	t.Log("Creating 10 workflow instances")
	for i := 0; i < 10; i++ {
		createdInstance, err := zbClient.CreateWorkflowInstance(TopicName, instance)
		Assert(t, nil, err, true)
		Assert(t, nil, createdInstance, false)
		Assert(t, zbcommon.WorkflowInstanceCreated, createdInstance.State, true)
	}
	t.Logf("Workflow instances created in %v", time.Since(wfStart))

	subscription, err := zbClient.TaskSubscription(TopicName, "task_subscription_test", "foo", TaskHandler)

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
