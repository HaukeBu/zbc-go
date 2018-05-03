package zbmsgpack

import (
	"encoding/json"
	"fmt"
)

// Resource is message pack structure used to represent a workflow definition.
type Resource struct {
	Resource     []byte `msgpack:"resource"`
	ResourceType string `msgpack:"resourceType"`
	ResourceName string `msgpack:"resourceName"`
}

// DeployWorkflow is message pack structure used when creating a workflow.
type DeployWorkflow struct {
	State string `msgpack:"state"`

	TopicName string      `msgpack:"topicName"`
	Resources []*Resource `msgpack:"resources"`
}

func (t *DeployWorkflow) String() string {
	b, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return fmt.Sprintf("json marshaling failed\n")
	}
	return fmt.Sprintf("%+v", string(b))
}

// WorkflowEvent
type WorkflowEvent struct {
	State         string  `msgpack:"state"`
	BPMNProcessID string  `msgpack:"bpmnProcessId"`
	Version       int     `msgpack:"version"`
	BPMNXml       []uint8 `msgpack:"bpmnXml"`
	DeploymentKey int64   `msgpack:"deploymentKey"`
}

func (e WorkflowEvent) String() string {
	b, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return fmt.Sprintf("json marshaling failed\n")
	}
	return fmt.Sprintf("%+v", string(b))
}
