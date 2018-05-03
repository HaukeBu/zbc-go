package zbmsgpack

import (
	"encoding/json"
	"fmt"
)

// CreateWorkflowInstance
type CreateWorkflowInstance struct {
	BPMNProcessID string  `msgpack:"bpmnProcessId"`
	Payload       []uint8 `msgpack:"payload"`
	State         string  `msgpack:"state"`
	Version       int     `msgpack:"version"`

	WorkflowInstanceKey uint64 `msgpack:"-"`
}

// String
func (t *CreateWorkflowInstance) String() string {
	b, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return fmt.Sprintf("json marshaling failed\n")
	}
	return fmt.Sprintf("%+v", string(b))
}

// WorkflowInstanceEvent
type WorkflowInstanceEvent struct {
	ActivityID          string  `msgpack:"activityID"`
	BPMNProcessID       string  `msgpack:"bpmnProcessId"`
	Payload             []uint8 `msgpack:"payload"`
	State               string  `msgpack:"state"`
	Version             int     `msgpack:"version"`
	WorkflowInstanceKey uint64  `msgpack:"workflowInstanceKey"`
	workflowKey         uint64  `msgpack:"workflowKey"`
}

// String
func (t *WorkflowInstanceEvent) String() string {
	b, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return fmt.Sprintf("json marshaling failed\n")
	}
	return fmt.Sprintf("%+v", string(b))
}
