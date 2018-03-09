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

func checkRoundRobinSequence(order uint16, sequence []uint16) bool {
	sequenced := true
	last := sequence[0]
	for i := 1; i < len(sequence); i++ {
		if last == order && sequence[i] != 0 {
			sequenced = false
			break
		}
		if last != order && last+1 != sequence[i] {
			sequenced = false
			break
		}
		last = sequence[i]
	}
	return sequenced
}

func TestRoundRobin(t *testing.T) {
	client := &zbexchange.ExchangeSvc{
		LikeTopologySvc:zbtopology.NewTopologySvc(BrokerAddr),
		RequestFactory: zbdispatch.NewRequestFactory(),
		ResponseHandler: zbdispatch.NewResponseHandler(),
	}

	t.Log("RefreshTopology starting ")
	topology, err := client.RefreshTopology()
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

	Assert(t, true, checkRoundRobinSequence(partitionsLength, sequence), true)
}
