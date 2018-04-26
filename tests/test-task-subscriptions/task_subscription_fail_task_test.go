package test_task_subscriptions

import (
	"testing"

	"github.com/zeebe-io/zbc-go/zbc"
	"github.com/zeebe-io/zbc-go/zbc/common"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsubscriptions"
	"github.com/zeebe-io/zbc-go/zbc/services/zbsubscribe"

	"fmt"
	. "github.com/zeebe-io/zbc-go/tests/test-helpers"
	"os"
	"sync/atomic"
	"time"
)

func TestTaskSubscriptionFailTask(t *testing.T) {
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

	t.Log("Creating 1 workflow instance")
	createdInstance, err := zbClient.CreateWorkflowInstance(hash, zbc.NewWorkflowInstance("demoProcess", -1, payload))
	Assert(t, nil, err, true)
	Assert(t, nil, createdInstance, false)
	Assert(t, zbcommon.WorkflowInstanceCreated, createdInstance.State, true)
	t.Log("Instance created")

	t.Log("Create task subscription")
	var ops uint64
	subStart := time.Now()
	subscription, err := zbClient.TaskSubscription(hash, "task_subscription_test", "foo", 30000, 30,
		func(client zbsubscribe.ZeebeAPI, event *zbsubscriptions.SubscriptionEvent) {
			atomic.AddUint64(&ops, 1)
			_, err := zbClient.FailTask(event)
			if err != nil {
				fmt.Printf("ERROR: Failed to FailTask: %+v\n", err)
				os.Exit(1)
			}
		})

	Assert(t, nil, err, true)
	Assert(t, nil, *subscription, false)
	t.Logf("Subscription creation took %v", time.Since(subStart))

	t.Log("Starting to consume subscription")
	go subscription.Start()

	t.Log("Starting to monitor subscription")
	for {
		op := atomic.LoadUint64(&ops)
		t.Log("Subscription processed tasks ", op)

		if op >= 1 {
			errs := subscription.Close()
			Assert(t, 0, len(errs), true)
			break
		}
		time.Sleep(time.Duration(1 * time.Second))
	}
}
