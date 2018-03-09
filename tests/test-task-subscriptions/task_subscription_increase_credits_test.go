package test_task_subscriptions

import (
	"sync/atomic"
	"testing"

	"github.com/zeebe-io/zbc-go/zbc"

	"github.com/zeebe-io/zbc-go/zbc/common"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsubscriptions"
	"github.com/zeebe-io/zbc-go/zbc/services/zbsubscribe"

	. "github.com/zeebe-io/zbc-go/tests/test-helpers"
	"time"
)

func TestTaskSubscriptionIncreaseCredits(t *testing.T) {
	zbClient, _ := zbc.NewClient(BrokerAddr)
	zbClient.CreateWorkflowFromFile(TopicName, zbcommon.BpmnXml, "../../examples/demoProcess.bpmn")

	payload := make(map[string]interface{})
	payload["a"] = "b"

	instance := zbc.NewWorkflowInstance("demoProcess", -1, payload)

	wfStart := time.Now()
	t.Log("Creating 50 workflow instances")
	for i := 0; i < 50; i++ {
		createdInstance, err := zbClient.CreateWorkflowInstance(TopicName, instance)
		Assert(t, nil, err, true)
		Assert(t, nil, createdInstance, false)
		Assert(t, zbcommon.WorkflowInstanceCreated, createdInstance.State, true)
	}
	t.Logf("Workflow instances created in %v", time.Since(wfStart))

	var ops uint64
	subscription, _ := zbClient.TaskSubscription(TopicName, "task_subscription_test", "foo",
		func(client zbsubscribe.ZeebeAPI, event *zbsubscriptions.SubscriptionEvent) {
			atomic.AddUint64(&ops, 1)
			// TODO: check credits
		})

	for {
		op := atomic.LoadUint64(&ops)
		t.Log("Subscription processed tasks ", op)
		//credits := subscription.PeekCredits()
		//t.Logf("Credits :: %+v", credits)

		if op >= 25 {
			errs := zbClient.CloseTaskSubscription(subscription)
			Assert(t, 0, len(errs), true)
			break
		}
		time.Sleep(time.Duration(1 * time.Second))
	}
}
