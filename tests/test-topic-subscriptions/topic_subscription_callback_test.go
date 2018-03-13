package test_topic_subscriptions

import (
	"testing"

	"github.com/zeebe-io/zbc-go/zbc/services/zbsubscribe"

	"github.com/zeebe-io/zbc-go/zbc"
	"github.com/zeebe-io/zbc-go/zbc/common"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsubscriptions"

	. "github.com/zeebe-io/zbc-go/tests/test-helpers"
	"time"
	"sync/atomic"
)

func TestTopicSubscriptionCallback(t *testing.T) {
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
	t.Logf("Workflow created")

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

	var ops uint64
	subscription, err := zbClient.TopicSubscription(hash, "default-name", 0,0,
		func(client zbsubscribe.ZeebeAPI, event *zbsubscriptions.SubscriptionEvent) error {
			Assert(t, nil, event, false)
			Assert(t, nil, client, false)
			atomic.AddUint64(&ops, 1)
			return nil
		})
	Assert(t, nil, subscription, false)
	Assert(t, nil, err, true)

	go subscription.Start()

	for {
		op := atomic.LoadUint64(&ops)
		t.Log("Subscription processed events ", op)
		if op >= 25 {
			errs := zbClient.CloseTopicSubscription(subscription)
			Assert(t, 0, len(errs), true)
			break
		}

		time.Sleep(time.Duration(time.Second * 1))
	}
}
