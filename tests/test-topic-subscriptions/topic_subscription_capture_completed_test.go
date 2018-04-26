package test_topic_subscriptions

import (
	"github.com/zeebe-io/zbc-go/zbc"
	"github.com/zeebe-io/zbc-go/zbc/common"
	"testing"

	. "github.com/zeebe-io/zbc-go/tests/test-helpers"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsubscriptions"
	"github.com/zeebe-io/zbc-go/zbc/services/zbsubscribe"
	"sync/atomic"
	"time"
)

var completeTester *testing.T

type SenderType struct {
	ID   int    `msgpack:"a"`
	Name string `msgpack:"name"`
}

type ReceivingType struct {
	ID   int    `msgpack:"foo"`
	Name string `msgpack:"name"`
}

func FooHandler(client zbsubscribe.ZeebeAPI, event *zbsubscriptions.SubscriptionEvent) {
	Assert(completeTester, nil, event, false)
	Assert(completeTester, nil, client, false)

	// var receiverUserType ReceivingType
	// event.LoadTask(&receiverUserType)
	//
	// receiverUserType.ID++
	// Assert(completeTester, 11, receiverUserType.ID, true)
	// Assert(completeTester, "", receiverUserType.Name, true)
	//
	// event.UpdateTask(&receiverUserType)

	task, err := client.CompleteTask(event)
	Assert(completeTester, nil, task, false)
	Assert(completeTester, nil, err, true)
}

func BarHandler(client zbsubscribe.ZeebeAPI, event *zbsubscriptions.SubscriptionEvent) {
	Assert(completeTester, nil, event, false)
	Assert(completeTester, nil, client, false)

	task, err := client.CompleteTask(event)
	Assert(completeTester, nil, task, false)
	Assert(completeTester, nil, err, true)
}

func FooBarHandler(client zbsubscribe.ZeebeAPI, event *zbsubscriptions.SubscriptionEvent) {
	Assert(completeTester, nil, event, false)
	Assert(completeTester, nil, client, false)

	task, err := client.CompleteTask(event)
	Assert(completeTester, nil, task, false)
	Assert(completeTester, nil, err, true)
}

func TestTopicSubscriptionCaptureCompleted(t *testing.T) {
	completeTester = t

	t.Log("Creating client")
	zbClient, err := zbc.NewClient(BrokerAddr)
	Assert(t, nil, err, true)
	Assert(t, nil, zbClient, false)
	t.Log("Client created")

	NumberOfPartitions = 5
	hash := CreateRandomTopicWithTimeout(t, zbClient)

	t.Log("Creating workflow")
	workflow, err := zbClient.CreateWorkflowFromFile(hash, zbcommon.BpmnXml, "../../examples/demoProcess.bpmn")
	Assert(t, nil, err, true)
	Assert(t, nil, workflow, false)
	Assert(t, nil, workflow.State, false)
	Assert(t, zbcommon.DeploymentCreated, workflow.State, true)
	t.Logf("Workflow created")

	user := SenderType{ID: 10, Name: "zeebe-io"}
	instance := zbc.NewWorkflowInstance("demoProcess", -1, user)

	wfStart := time.Now()
	t.Log("Creating 50 workflow instances")
	for i := 0; i < 50; i++ {
		createdInstance, err := zbClient.CreateWorkflowInstance(hash, instance)
		Assert(t, nil, err, true)
		Assert(t, nil, createdInstance, false)
		Assert(t, zbcommon.WorkflowInstanceCreated, createdInstance.State, true)
	}
	t.Logf("Workflow instances created in %v", time.Since(wfStart))

	t.Log("Create task subscription on type 'foo'")
	fooSub, err := zbClient.TaskSubscription(hash, "task_subscription_test", "foo", 30000, 30, FooHandler)
	Assert(t, nil, err, true)
	Assert(t, nil, *fooSub, false)

	t.Log("Create task subscription on type 'bar'")
	barSub, err := zbClient.TaskSubscription(hash, "task_subscription_test", "bar", 30000, 30, BarHandler)
	Assert(t, nil, err, true)
	Assert(t, nil, *barSub, false)

	t.Log("Create task subscription on type 'foobar'")
	foobarSub, err := zbClient.TaskSubscription(hash, "task_subscription_test", "foobar", 30000, 30, FooBarHandler)
	Assert(t, nil, err, true)
	Assert(t, nil, *foobarSub, false)

	go foobarSub.Start()
	go barSub.Start()
	go fooSub.Start()

	var ops uint64
	subscription, err := zbClient.TopicSubscription(hash, "default-name", 0, 0, false,
		func(client zbsubscribe.ZeebeAPI, event *zbsubscriptions.SubscriptionEvent) error {
			Assert(t, nil, event, false)
			Assert(t, nil, client, false)

			eventPayload, err := event.GetEvent()
			if eventPayload["state"] == "WORKFLOW_INSTANCE_COMPLETED" {
				atomic.AddUint64(&ops, 1)
			}

			Assert(t, nil, err, true)
			Assert(t, nil, eventPayload, false)

			return nil
		})
	Assert(t, nil, subscription, false)
	Assert(t, nil, err, true)

	go subscription.Start()

	for {
		op := atomic.LoadUint64(&ops)
		t.Log("Subscription processed events ", op)
		if op >= 50 {
			errs := subscription.Close()
			Assert(t, 0, len(errs), true)
			break
		}

		time.Sleep(time.Duration(time.Second * 1))
	}
}
