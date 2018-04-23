package testbroker

import (
	"testing"

	"github.com/zeebe-io/zbc-go/zbc/services/zbtopology"

	"github.com/zeebe-io/zbc-go/zbc"
	"github.com/zeebe-io/zbc-go/zbc/common"
	"github.com/zeebe-io/zbc-go/zbc/models/zbdispatch"
	"github.com/zeebe-io/zbc-go/zbc/services/zbexchange"

	. "github.com/zeebe-io/zbc-go/tests/test-helpers"
)

func TestRoundRobin(t *testing.T) {
	topologySvc, err := zbtopology.NewTopologySvc(BrokerAddr)
	Assert(t, nil, err, true)

	client := &zbexchange.ExchangeSvc{
		LikeTopologySvc: topologySvc,
		RequestFactory:  zbdispatch.NewRequestFactory(),
		ResponseHandler: zbdispatch.NewResponseHandler(),
	}

	t.Log("GetTopology starting ")
	topology, err := client.GetTopology()
	Assert(t, nil, topology, false)
	Assert(t, nil, err, true)
	t.Logf("Topology :: %+v", topology)

	fatRobin := client.GetRoundRobinCtl()
	partitions, _ := fatRobin.Partitions(TopicName)
	partitionsLength := uint16(len(partitions) - 1)

	workflow, err := client.CreateWorkflowFromFile(TopicName, zbcommon.BpmnXml, "../../examples/demoProcess.bpmn")
	Assert(t, nil, err, true)
	Assert(t, nil, workflow, false)
	Assert(t, nil, workflow.State, false)
	Assert(t, zbcommon.DeploymentCreated, workflow.State, true)

	payload := make(map[string]interface{})
	payload["a"] = "b"

	var sequence []uint16
	instance := zbc.NewWorkflowInstance("demoProcess", -1, payload)
	for i := 0; i < NumberOfPartitions+NumberOfPartitions; i++ {
		createdInstance, err := client.CreateWorkflowInstance(TopicName, instance)
		Assert(t, nil, err, true)
		Assert(t, nil, createdInstance, false)
		Assert(t, zbcommon.WorkflowInstanceCreated, createdInstance.State, true)

		val := fatRobin.PeekPartitionIndex(TopicName)
		if val != nil {
			sequence = append(sequence, *val)
		}
	}

	Assert(t, true, CheckRoundRobinSequence(partitionsLength, sequence), true)
}
