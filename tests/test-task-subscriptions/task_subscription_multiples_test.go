package test_task_subscriptions

import (
	"fmt"
	. "github.com/zeebe-io/zbc-go/tests/test-helpers"
	"github.com/zeebe-io/zbc-go/zbc"
	"github.com/zeebe-io/zbc-go/zbc/common"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsubscriptions"
	"github.com/zeebe-io/zbc-go/zbc/services/zbsubscribe"
	"sync/atomic"
	"testing"
	"time"
)

func TestTaskSubscriptionMultipleSubscriptions(t *testing.T) {
	var counter uint64

	type ProcessSenderType struct {
		ID   int    `msgpack:"a"`
		Name string `msgpack:"name"`
	}

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

	user := ProcessSenderType{ID: 10, Name: "zeebe-io"}
	var instances, i uint64 = 1000, 0

	t.Logf("Creating %v workflow instances", instances)
	instance := zbc.NewWorkflowInstance("demoProcess", -1, user)
	for ; i < instances; i++ {
		createdInstance, err := zbClient.CreateWorkflowInstance(hash, instance)
		Assert(t, nil, err, true)
		Assert(t, nil, createdInstance, false)
		Assert(t, zbcommon.WorkflowInstanceCreated, createdInstance.State, true)
	}
	t.Logf("Instances created")

	t.Log("Create task subscription on type 'foo'")
	fooSub, err := zbClient.TaskSubscription(hash, "task_subscription_test", "foo", 30000,  30,
		func(client zbsubscribe.ZeebeAPI, event *zbsubscriptions.SubscriptionEvent) {
			_, err := zbClient.CompleteTask(event)
			if err != nil {
				fmt.Println(err)
			}
		})
	Assert(t, nil, err, true)
	go fooSub.Start()

	t.Log("Create task subscription on type 'bar'")
	barSub, err := zbClient.TaskSubscription(hash, "task_subscription_task_b", "bar", 30000, 30,
		func(client zbsubscribe.ZeebeAPI, event *zbsubscriptions.SubscriptionEvent) {
			_, err := zbClient.CompleteTask(event)
			if err != nil {
				fmt.Println(err)
			}
		})
	Assert(t, nil, err, true)
	go barSub.Start()

	t.Log("Create task subscription on type 'foobar'")
	foobarSub, err := zbClient.TaskSubscription(hash, "task_subscription_test_c", "foobar", 30000, 30,
		func(client zbsubscribe.ZeebeAPI, event *zbsubscriptions.SubscriptionEvent) {
			_, err := zbClient.CompleteTask(event)
			if err != nil {
				fmt.Println(err)
			}

			atomic.AddUint64(&counter, 1)
		})

	Assert(t, nil, err, true)
	go foobarSub.Start()

	subStart := time.Now()
	for {
		val := atomic.LoadUint64(&counter)
		fmt.Println(val)
		if val == instances {
			t.Logf("Completing took %v", time.Since(subStart))
			break
		}
		time.Sleep(time.Duration(500 * time.Millisecond))
	}
}
