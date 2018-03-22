package test_topic_subscriptions

import (
"testing"

"github.com/zeebe-io/zbc-go/zbc/services/zbsubscribe"

"github.com/zeebe-io/zbc-go/zbc"
"github.com/zeebe-io/zbc-go/zbc/common"
"github.com/zeebe-io/zbc-go/zbc/models/zbsubscriptions"

"sync/atomic"
"time"

. "github.com/zeebe-io/zbc-go/tests/test-helpers"
	"fmt"
)

func TestTopicSubscriptionCloseOnError(t *testing.T) {
	t.Log("Creating client")
	zbClient, err := zbc.NewClient(BrokerAddr)
	Assert(t, nil, err, true)
	Assert(t, nil, zbClient, false)
	t.Log("Client created")

	partitionCount := 5
	t.Log("Creating topic")
	hash := RandStringBytes(25)
	topic, err := zbClient.CreateTopic(hash, partitionCount)
	Assert(t, nil, err, true)
	Assert(t, nil, topic, false)
	t.Logf("Topic %s created with %d partitions", hash, partitionCount)

	t.Log("Creating workflow")
	workflow, err := zbClient.CreateWorkflowFromFile(hash, zbcommon.BpmnXml, "../../examples/demoProcess.bpmn")
	Assert(t, nil, err, true)
	Assert(t, nil, workflow, false)
	Assert(t, nil, workflow.State, false)
	Assert(t, zbcommon.DeploymentCreated, workflow.State, true)
	t.Logf("Workflow created")

	payload := make(map[string]interface{})
	payload["a"] = "b"

	var i, wiCount uint64 = 0, 5
	instance := zbc.NewWorkflowInstance("demoProcess", -1, payload)

	wfStart := time.Now()
	t.Log("Creating 5 workflow instances")
	for ; i < wiCount; i++ {
		createdInstance, err := zbClient.CreateWorkflowInstance(hash, instance)
		Assert(t, nil, err, true)
		Assert(t, nil, createdInstance, false)
		Assert(t, zbcommon.WorkflowInstanceCreated, createdInstance.State, true)
	}
	t.Logf("Workflow instances created in %v", time.Since(wfStart))

	var ops uint64
	now := time.Now()
	subscription, err := zbClient.TopicSubscription(hash, "callback-test", 0, 0, false,
		func(client zbsubscribe.ZeebeAPI, event *zbsubscriptions.SubscriptionEvent) error {
			Assert(t, nil, event, false)
			Assert(t, nil, client, false)
			atomic.AddUint64(&ops, 1)
			return fmt.Errorf("processing failed")
		})
	Assert(t, nil, subscription, false)
	Assert(t, nil, err, true)

	go subscription.Start()

	for {
		op := atomic.LoadUint64(&ops)
		t.Log("Subscription processed events ", op)
		if op == 3 {
			time.Sleep(time.Duration(time.Millisecond * 500))
			if op == 3 {
				break
			}
		}

		time.Sleep(time.Duration(time.Millisecond * 900))
	}
	t.Logf("executed in %v", time.Since(now))
}
