package zbmsgpack

import (
	"fmt"
	"encoding/json"
)

// IncidentEvent
type IncidentEvent struct {
	ActivityID           string  `msgpack:"activityID"`
	ActivityInstanceKey  uint64  `msgpack:"activityInstanceKey"`
	BPMNProcessID        string  `msgpack:"bpmnProcessId"`
	ErrorMessage         string  `msgpack:"errorMessage"`
	ErrorType            string  `msgpack:"errorType"`
	FailureEventPosition uint64  `msgpack:"failureEventPosition"`
	Payload              []uint8 `msgpack:"payload"`
	State                string  `msgpack:"state"`
	TaskKey              int64   `msgpack:"taskKey"`
	WorkflowInstanceKey  uint64  `msgpack:"workflowInstanceKey"`
}


func (e IncidentEvent) String() string {
	b, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return fmt.Sprintf("json marshaling failed\n")
	}
	return fmt.Sprintf("%+v", string(b))
}
