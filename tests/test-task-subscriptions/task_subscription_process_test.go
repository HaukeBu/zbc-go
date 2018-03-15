package test_task_subscriptions

import (
	"testing"
)

var processTester *testing.T
//
//type SenderType struct {
//	ID   int `msgpack:"a"`
//	Name string `msgpack:"name"`
//}
//
//type ReceivingType struct {
//	ID   int `msgpack:"foo"`
//	Name string `msgpack:"name"`
//
//}
//
//func FooHandler(client zbsubscribe.ZeebeAPI, event *zbsubscriptions.SubscriptionEvent) {
//	Assert(marshallTester, nil, event, false)
//	Assert(marshallTester, nil, client, false)
//
//	var receiverUserType ReceivingType
//	event.Load(&receiverUserType) // move it to event
//
//	receiverUserType.ID++
//	Assert(marshallTester, 11, receiverUserType.ID, true)
//	Assert(marshallTester, "", receiverUserType.Name, true)
//
//	event.Update(&receiverUserType)
//
//	task, err := client.CompleteTask(event)
//	Assert(marshallTester, nil, task, false)
//	Assert(marshallTester, nil, err, true)
//}
//
//func BarHandler(client zbsubscribe.ZeebeAPI, event *zbsubscriptions.SubscriptionEvent) {
//	Assert(marshallTester, nil, event, false)
//	Assert(marshallTester, nil, client, false)
//
//	var receiverUserType SenderType
//	event.Load(&receiverUserType)
//
//	receiverUserType.ID++
//	event.Update(&receiverUserType)
//
//	task, err := client.CompleteTask(event)
//	Assert(marshallTester, nil, task, false)
//	Assert(marshallTester, nil, err, true)
//}
//
//func FooBarHandler(client zbsubscribe.ZeebeAPI, event *zbsubscriptions.SubscriptionEvent) {
//	Assert(marshallTester, nil, event, false)
//	Assert(marshallTester, nil, client, false)
//
//	var receiverUserType SenderType
//	event.Load(&receiverUserType)
//
//	receiverUserType.ID++
//	event.Update(&receiverUserType)
//
//	task, err := client.CompleteTask(event)
//	Assert(marshallTester, nil, task, false)
//	Assert(marshallTester, nil, err, true)
//}

func TestTaskSubscriptionProcessMarshaller(t *testing.T) {
	processTester = t

	//t.Log("Creating client")
	//zbClient, err := zbc.NewClient(BrokerAddr)
	//Assert(t, nil, err, true)
	//Assert(t, nil, zbClient, false)
	//t.Log("Client created")
	//
	//t.Log("Creating topic")
	//hash := RandStringBytes(25)
	//topic, err := zbClient.CreateTopic(hash, NumberOfPartitions)
	//Assert(t, nil, err, true)
	//Assert(t, nil, topic, false)
	//t.Logf("Topic %s created with %d partitions", hash, NumberOfPartitions)
	//
	//t.Log("Creating workflow")
	//workflow, err := zbClient.CreateWorkflowFromFile(hash, zbcommon.BpmnXml, "../../examples/demoProcess.bpmn")
	//Assert(t, nil, err, true)
	//Assert(t, nil, workflow, false)
	//Assert(t, nil, workflow.State, false)
	//Assert(t, zbcommon.DeploymentCreated, workflow.State, true)
	//t.Log("Workflow created")
	//
	//user := SenderType{ID: 10, Name: "zeebe-io"}
	////payload := make(map[string]interface{})
	////payload["a"] = "b"
	////payload["user"] = "sam"
	//
	//instance := zbc.NewWorkflowInstance("demoProcess", -1, user)
	//
	//t.Log("Creating 1 workflow instance")
	//createdInstance, err := zbClient.CreateWorkflowInstance(hash, instance)
	//Assert(t, nil, err, true)
	//Assert(t, nil, createdInstance, false)
	//Assert(t, zbcommon.WorkflowInstanceCreated, createdInstance.State, true)
	//t.Logf("Instance created with key %d", createdInstance.WorkflowInstanceKey)
	//
	//t.Log("Create task subscription on type 'foo'")
	//subStart := time.Now()
	//fooSub, err := zbClient.TaskSubscription(hash, "task_subscription_test", "foo", 30, FooHandler)
	//Assert(t, nil, err, true)
	//Assert(t, nil, *fooSub, false)
	//t.Logf("Subscription creation took %v", time.Since(subStart))
	//
	//t.Log("Starting to consume subscription")
	//pStart := time.Now()
	//fooSub.ProcessNext(1)
	//t.Logf("Subscription processing took %v", time.Since(pStart))
	//fooSub.Close()

	//t.Log("Create task subscription on type 'bar'")
	//subStart = time.Now()
	//barSub, err := zbClient.TaskSubscription(hash, "task_subscription_task_b", "bar", 30, BarHandler)
	//Assert(t, nil, err, true)
	//Assert(t, nil, *barSub, false)
	//t.Logf("Subscription creation took %v", time.Since(subStart))
	//
	//t.Log("Starting to consume subscription")
	//barSub.ProcessNext(1)
	//barSub.Close()

	//t.Log("Create task subscription on type 'foobar'")
	//subStart = time.Now()
	//foobarSub, err := zbClient.TaskSubscription(hash, "task_subscription_test", "foobar", 30, FooBarHandler)
	//Assert(t, nil, err, true)
	//Assert(t, nil, *foobarSub, false)
	//t.Logf("Subscription creation took %v", time.Since(subStart))
	//t.Log("Starting to consume subscription")
	//go foobarSub.Start()

}
