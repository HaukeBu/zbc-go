package test_task_subscriptions

import (
	"testing"
	"time"

	"github.com/zeebe-io/zbc-go/zbc"

	"github.com/zeebe-io/zbc-go/zbc/common"

	. "github.com/zeebe-io/zbc-go/tests/test-helpers"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsubscriptions"
	"github.com/zeebe-io/zbc-go/zbc/services/zbsubscribe"
	"sync/atomic"
)

func TestTaskSubscriptionIncreaseCreditsLoop(t *testing.T) {
	t.Log("Creating client")
	zbClient, err := zbc.NewClient(BrokerAddr)
	Assert(t, nil, err, true)
	Assert(t, nil, zbClient, false)
	t.Log("Client created")

	t.Log("Creating topic")
	hash := RandStringBytes(25)
	topic, err := zbClient.CreateTopic(hash, NumberOfPartitions)
	Assert(t, nil, err, true)
	Assert(t, nil, topic, false)
	t.Logf("Topic %s created with %d partitions", hash, NumberOfPartitions)

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
	var i, wiCount uint64 = 0, 5000
	t.Log("Creating 5000 workflow instances")
	for ; i < wiCount; i++ {
		createdInstance, err := zbClient.CreateWorkflowInstance(hash, instance)
		Assert(t, nil, err, true)
		Assert(t, nil, createdInstance, false)
		Assert(t, zbcommon.WorkflowInstanceCreated, createdInstance.State, true)
	}
	t.Logf("Workflow instances created in %v", time.Since(wfStart))

	// MARK: we make payload sufficiently large (wiCount) and credits small.
	// If everything passes test should finish in time.
	// Otherwise go test runtime will teardown the test after 10 minutes
	var ops uint64
	subStart := time.Now()
	subscription, err := zbClient.TaskSubscription(hash, "task_subscription_test", "foo", 30,
		func(client zbsubscribe.ZeebeAPI, event *zbsubscriptions.SubscriptionEvent) {
			atomic.AddUint64(&ops, 1)
		})

	Assert(t, nil, err, true)
	Assert(t, nil, *subscription, false)
	go subscription.Start()
	t.Logf("Subscriptio§n opening took %v", time.Since(subStart))

	pStart := time.Now()
	for {
		op := atomic.LoadUint64(&ops)
		if op == wiCount {
			break
		}
		time.Sleep(time.Duration(50 * time.Millisecond))
	}
	t.Logf("Subscription processing took %v", time.Since(pStart))

	op := atomic.LoadUint64(&ops)
	Assert(t, uint64(wiCount), op, true)

	errs := subscription.Close()
	Assert(t, 0, len(errs), true)
}
