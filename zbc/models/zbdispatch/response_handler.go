package zbdispatch

import (
	"time"

	"github.com/vmihailenco/msgpack"
	"github.com/zeebe-io/zbc-go/zbc/common"
	"github.com/zeebe-io/zbc-go/zbc/models/zbmsgpack"
)

type ResponseHandler struct{}

func (rf *ResponseHandler) UnmarshalPartition(msg *Message) *zbmsgpack.PartitionCollection {
	var resp zbmsgpack.PartitionCollection
	msgpack.Unmarshal(msg.Data, &resp)
	return &resp
}

func (rf *ResponseHandler) UnmarshalTopology(msg *Message) *zbmsgpack.ClusterTopology {
	var resp zbmsgpack.ClusterTopologyResponse
	msgpack.Unmarshal(msg.Data, &resp)

	ct := zbmsgpack.ClusterTopology{
		AddrByPartitionID:      make(map[uint16]string),   // contains partitionID: brokerAddr
		PartitionIDByTopicName: make(map[string][]uint16), // contains topicName: [PartitionID]
		Brokers:                resp.Brokers,
		UpdatedAt:              time.Now(),
	}

	for _, broker := range resp.Brokers {
		for _, partition := range broker.Partitions {

			ct.PartitionIDByTopicName[partition.TopicName] = append(ct.PartitionIDByTopicName[partition.TopicName], partition.PartitionID)

			if partition.State == zbcommon.StateLeader {
				ct.AddrByPartitionID[partition.PartitionID] = broker.Addr()
			}
		}
	}

	ct.Brokers = resp.Brokers
	return &ct
}

func (rf *ResponseHandler) UnmarshalTopic(msg *Message) *zbmsgpack.CreateTopic {
	var topic zbmsgpack.CreateTopic
	msgpack.Unmarshal(msg.Data, &topic)

	return &topic
}

func (rf *ResponseHandler) UnmarshalTask(m *Message) *zbmsgpack.Task {
	var d zbmsgpack.Task
	err := msgpack.Unmarshal(m.Data, &d)
	if err != nil {
		return nil
	}
	if len(d.Type) > 0 {
		return &d
	}
	return nil
}

func (rf *ResponseHandler) UnmarshalTaskSubscriptionInfo(m *Message) *zbmsgpack.TaskSubscriptionInfo {
	var d zbmsgpack.TaskSubscriptionInfo
	err := msgpack.Unmarshal(m.Data, &d)
	if err != nil {
		return nil
	}
	return &d
}

func (rf *ResponseHandler) UnmarshalTopicSubscriptionInfo(m *Message) *zbmsgpack.TopicSubscriptionInfo {
	var d zbmsgpack.TopicSubscriptionInfo
	err := msgpack.Unmarshal(m.Data, &d)
	if err != nil {
		return nil
	}
	return &d
}

func (rf *ResponseHandler) UnmarshalWorkflow(m *Message) *zbmsgpack.Workflow {
	var d zbmsgpack.Workflow
	err := msgpack.Unmarshal(m.Data, &d)
	if err != nil {
		return nil
	}
	if len(d.State) > 0 {
		return &d
	}
	return nil
}

func (rf *ResponseHandler) UnmarshalWorkflowInstance(m *Message) *zbmsgpack.WorkflowInstance {
	var d zbmsgpack.WorkflowInstance
	err := msgpack.Unmarshal(m.Data, &d)
	if err != nil {
		return nil
	}
	if len(d.BPMNProcessID) > 0 {
		return &d
	}
	return nil
}

func (rf *ResponseHandler) UnmarshalTopicSubAck(m *Message) *zbmsgpack.TopicSubscriptionAckRequest {
	var d zbmsgpack.TopicSubscriptionAckRequest
	err := msgpack.Unmarshal(m.Data, &d)
	if err != nil {
		return nil
	}
	if d.State == zbcommon.TopicSubscriptionAcknowledgedState {
		return &d
	}
	return nil
}

func NewResponseHandler() *ResponseHandler {
	zbcommon.ZBL.Debug().Msg("creating new response handler")
	return &ResponseHandler{}
}
