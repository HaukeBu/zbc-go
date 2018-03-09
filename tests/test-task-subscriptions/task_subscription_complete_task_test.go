package test_task_subscriptions

import (
	"testing"

	. "github.com/zeebe-io/zbc-go/tests/test-helpers"
	"github.com/zeebe-io/zbc-go/zbc"
	"github.com/zeebe-io/zbc-go/zbc/common"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsubscriptions"
	"github.com/zeebe-io/zbc-go/zbc/services/zbsubscribe"
	"sync/atomic"
	"time"
)

// type UserType struct {
// 	ID int
// 	Name string
// }

func TestTaskSubscriptionCompleteTask(t *testing.T) {
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

	t.Log("Creating 1 workflow instances")
	createdInstance, err := zbClient.CreateWorkflowInstance(TopicName, zbc.NewWorkflowInstance("demoProcess", -1, payload))
	Assert(t, nil, err, true)
	Assert(t, nil, createdInstance, false)
	Assert(t, zbcommon.WorkflowInstanceCreated, createdInstance.State, true)

	var ops uint64
	subscription, err := zbClient.TaskSubscription(TopicName, "task_subscription_test", "foo",
		func(client zbsubscribe.ZeebeAPI, event *zbsubscriptions.SubscriptionEvent) {
			atomic.AddUint64(&ops, 1)
			Assert(t, nil, event, false)
			Assert(t, nil, client, false)

			// TODO: Implement unmarshaling of custom user struct and updating event with the struct change
			// ut := &UserType{}
			// err := zbClient.UnmarshalCustomStruct(event, ut)

			task, err := zbClient.CompleteTask(event)
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
