package test_topic_subscriptions

import (
	"testing"
	"github.com/zeebe-io/zbc-go/zbc"
	"github.com/zeebe-io/zbc-go/zbc/common"
	"time"
	"github.com/zeebe-io/zbc-go/zbc/services/zbsubscribe"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsubscriptions"
	"sync/atomic"

	. "github.com/zeebe-io/zbc-go/tests/test-helpers"
	"fmt"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsbe"
)

func TestTopicSubscriptionCaptureIncident(t *testing.T) {
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
	payload["b"] = "a"

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

	//subscription.SetFilter(zbsbe.EventType.INCIDENT_EVENT)
	go subscription.Start()

	for {
		op := atomic.LoadUint64(&ops)
		fmt.Println(op)
		t.Log("Subscription processed events ", op)
		if op == 37 {
			errs := subscription.Close()
			Assert(t, 0, len(errs), true)
			break
		}

		time.Sleep(time.Duration(time.Millisecond * 900))
	}
	t.Logf("executed in %v", time.Since(now))
	recorder.Dump("../../.logs/incidents.json")

	for index, e := range recorder.Records {
		event := e.(*zbsubscriptions.SubscriptionEvent)
		if event.Metadata.EventType == zbsbe.EventType.INCIDENT_EVENT  {

			for {
				index--
				e := recorder.Records[index].(*zbsubscriptions.SubscriptionEvent)
				if e.Metadata.EventType == zbsbe.EventType.WORKFLOW_INSTANCE_EVENT {
					if evt, _ := e.GetWorkflowInstanceEvent(); evt.State == "WORKFLOW_INSTANCE_CREATED" {
						updated := make(map[string]interface{})
						updated["a"] = "b"
						e.UpdatePayload(updated)

						wi, err := zbClient.UpdatePayload(e)
						Assert(t, nil, wi, false)
						Assert(t, nil, err,true)
						Assert(t, "PAYLOAD_UPDATED", wi.State, true)
					}

					break
				} else {
					continue
				}
			}


		}
	}
}
