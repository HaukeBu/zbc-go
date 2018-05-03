package zbmsgpack

import (
	"encoding/json"
	"fmt"
)

// Task structure is used when creating or read a test-task-subscriptions.
type Task struct {
	State        string                 `msgpack:"state"`
	LockTime     uint64                 `msgpack:"lockTime"`
	LockOwner    string                 `msgpack:"lockOwner"`
	Headers      map[string]interface{} `msgpack:"headers"`
	CustomHeader map[string]interface{} `msgpack:"customHeaders"`
	Retries      int                    `msgpack:"retries"`
	Type         string                 `msgpack:"type"`

	Payload      []uint8                `msgpack:"payload"`
	//PayloadJSON  map[string]interface{} `yaml:"payload" msgpack:"-" json:"-"`
}

func (t *Task) String() string {
	b, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return fmt.Sprintf("json marshaling failed\n")
	}
	return fmt.Sprintf("%+v", string(b))
}
