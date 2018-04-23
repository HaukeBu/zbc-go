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

	t.Log("Creating client")
	zbClient, err := zbc.NewClient(BrokerAddr)
	Assert(t, nil, err, true)
	Assert(t, nil, zbClient, false)
	t.Log("Client created")

	hash := CreateRandomTopicWithTimeout(t, zbClient)

	t.Log("Creating workflow")
	workflow, err := zbClient.CreateWorkflowFromFile(hash, zbcommon.BpmnXml, "../../examples/demoProcess.bpmn")
	Assert(t, nil, err, true)
	Assert(t, nil, workflow, false)
	Assert(t, nil, workflow.State, false)
	Assert(t, zbcommon.DeploymentCreated, workflow.State, true)
	t.Log("Workflow created")

	payload := make(map[string]interface{})
	payload["a"] = "b"

	instance := zbc.NewWorkflowInstance("demoProcess", -1, payload)

	wfStart := time.Now()
	t.Log("Creating 10 workflow instances")
	for i := 0; i < 10; i++ {
		createdInstance, err := zbClient.CreateWorkflowInstance(hash, instance)
		Assert(t, nil, err, true)
		Assert(t, nil, createdInstance, false)
		Assert(t, zbcommon.WorkflowInstanceCreated, createdInstance.State, true)
	}
	t.Logf("Workflow instances created in %v", time.Since(wfStart))

	subscription, err := zbClient.TaskSubscription(hash, "task_subscription_test", "foo", 30, TaskHandler)

	t.Log("Starting to consume subscription")
	go subscription.Start()

	for {
		op := atomic.LoadUint64(&ops)
		t.Log("Subscription processed tasks ", op)
		if op >= 10 {
			errs := subscription.Close()
			Assert(t, 0, len(errs), true)
			break
		}
		time.Sleep(time.Duration(1 * time.Second))
	}

}
