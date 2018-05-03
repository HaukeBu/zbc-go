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
)

func TestTopicSubscriptionRecord(t *testing.T) {
	t.Log("Creating client")
	zbClient, err := zbc.NewClient(BrokerAddr)
	Assert(t, nil, err, true)
	Assert(t, nil, zbClient, false)
	t.Log("Client created")

	NumberOfPartitions = 1
	hash := CreateRandomTopicWithTimeout(t, zbClient)

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

	t.Log("Creating 5 workflow instances")

	wfStart := time.Now()
	for ; i < wiCount; i++ {
		createdInstance, err := zbClient.CreateWorkflowInstance(hash, instance)
		Assert(t, nil, err, true)
		Assert(t, nil, createdInstance, false)
		Assert(t, zbcommon.WorkflowInstanceCreated, createdInstance.State, true)
	}
	t.Logf("Workflow instances created in %v", time.Since(wfStart))

	var recorder zbcommon.EventRecorder
	var ops uint64
	now := time.Now()
	subscription, err := zbClient.TopicSubscription(hash, "callback-test", 30, 0, false,
		func(client zbsubscribe.ZeebeAPI, event *zbsubscriptions.SubscriptionEvent) error {
			Assert(t, nil, event, false)
			Assert(t, nil, client, false)
			recorder.Add(event)
			atomic.AddUint64(&ops, 1)
			return nil
		})

	Assert(t, nil, subscription, false)
	Assert(t, nil, err, true)

	go subscription.Start()

	for {
		op := atomic.LoadUint64(&ops)
		t.Log("Subscription processed events ", op)
		if op == (wiCount*8)+(2*uint64(NumberOfPartitions)) {
			errs := subscription.Close()
			Assert(t, 0, len(errs), true)
			break
		}

		time.Sleep(time.Duration(time.Millisecond * 900))
	}

	t.Logf("executed in %v", time.Since(now))
	recorder.Dump("../../.logs/topic.json")
}
