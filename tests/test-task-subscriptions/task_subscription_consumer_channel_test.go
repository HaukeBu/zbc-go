package test_task_subscriptions

import (
	"github.com/zeebe-io/zbc-go/zbc/models/zbsubscriptions"
	"github.com/zeebe-io/zbc-go/zbc/services/zbsubscribe"
	"testing"

	. "github.com/zeebe-io/zbc-go/tests/test-helpers"
	"github.com/zeebe-io/zbc-go/zbc"
	"github.com/zeebe-io/zbc-go/zbc/common"
	"sync/atomic"
	"time"
)

func TestTaskSubscriptionConsumerChannelTest(t *testing.T) {
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

	t.Log("Create task subscription")
	var ops uint64
	subStart := time.Now()
	subscription, err := zbClient.TaskSubscription(hash, "task_subscription_test", "foo", 30,
		func(client zbsubscribe.ZeebeAPI, event *zbsubscriptions.SubscriptionEvent) {
			atomic.AddUint64(&ops, 1)
		})

	Assert(t, nil, err, true)
	Assert(t, nil, *subscription, false)
	t.Logf("Subscription creation took %v", time.Since(subStart))

	t.Log("Starting to consume subscription")
	go subscription.Start()

	payload := make(map[string]interface{})
	payload["a"] = "b"

	instance := zbc.NewWorkflowInstance("demoProcess", -1, payload)

	wfStart := time.Now()
	t.Log("Creating 50 workflow instances")
	for i := 0; i < 50; i++ {
		createdInstance, err := zbClient.CreateWorkflowInstance(hash, instance)
		Assert(t, nil, err, true)
		Assert(t, nil, createdInstance, false)
		Assert(t, zbcommon.WorkflowInstanceCreated, createdInstance.State, true)
	}
	t.Logf("Workflow instances created in %v", time.Since(wfStart))

	for {
		op := atomic.LoadUint64(&ops)
		t.Log("Subscription processed instances ", op)

		if op >= 25 {
			errs := zbClient.CloseTaskSubscription(subscription)
			Assert(t, 0, len(errs), true)
			break
		}
		time.Sleep(time.Duration(1 * time.Second))
	}
}
