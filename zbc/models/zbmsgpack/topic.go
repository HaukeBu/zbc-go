package zbmsgpack

import (
	"encoding/json"
	"fmt"
)

type CreateTopic struct {
	Name       string `msgpack:"name"`
	State      string `msgpack:"state"`
	Partitions int    `msgpack:"partitions"`
	ReplicationFactor int `msgpack:"replicationFactor"`
}

func (t *CreateTopic) String() string {
	b, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return fmt.Sprintf("json marshaling failed\n")
	}
	return fmt.Sprintf("%+v", string(b))
}

func NewTopic(name, state string, partitionsNum, replicationFactor int) *CreateTopic {
	return &CreateTopic{
		Name: name,
		State: state,
		Partitions: partitionsNum,
		ReplicationFactor: replicationFactor,
	}
}
