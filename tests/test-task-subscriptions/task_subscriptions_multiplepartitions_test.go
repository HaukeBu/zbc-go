package test_task_subscriptions

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/zeebe-io/zbc-go/zbc"

	"github.com/zeebe-io/zbc-go/zbc/common"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsubscriptions"
	"github.com/zeebe-io/zbc-go/zbc/services/zbsubscribe"

	. "github.com/zeebe-io/zbc-go/tests/test-helpers"
)

func TestTaskSubscriptionMultiplePartitions(t *testing.T) {
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
	t.Log("Creating 1000 workflow instances")
	for i := 0; i < 1000; i++ {
		createdInstance, err := zbClient.CreateWorkflowInstance(hash, instance)
		Assert(t, nil, err, true)
		Assert(t, nil, createdInstance, false)
		Assert(t, zbcommon.WorkflowInstanceCreated, createdInstance.State, true)
	}
	t.Logf("Workflow instances created in %v", time.Since(wfStart))

	mut := &sync.Mutex{}
	var partitionIDs []uint16
	var ops uint64
	subscription, _ := zbClient.TaskSubscription(hash, "task_subscription_test", "foo", 30000, 30,
		func(client zbsubscribe.ZeebeAPI, event *zbsubscriptions.SubscriptionEvent) {
			atomic.AddUint64(&ops, 1)

			// MARK: User is responsible for this kind of behaviour.
			mut.Lock()
			partitionIDs = append(partitionIDs, event.Event.PartitionId)
			mut.Unlock()

		})

	t.Log("Starting to consume subscription")
	go subscription.Start()

	for {
		op := atomic.LoadUint64(&ops)
		t.Log("Subscription processed tasks ", op)

		if op >= 1000 {
			errs := zbClient.CloseTaskSubscription(subscription)
			Assert(t, 0, len(errs), true)
			break
		}
		time.Sleep(time.Duration(1 * time.Second))
	}

	mut.Lock()
	Assert(t, true, len(partitionIDs) >= 1000, true)
	mut.Unlock()
}
