package zbmsgpack

import (
	"encoding/json"
	"fmt"
)

type CreateTopic struct {
	Name       string `msgpack:"name"`
	State      string `msgpack:"state"`
	Partitions int    `msgpack:"partitions"`
}

func (t *CreateTopic) String() string {
	b, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return fmt.Sprintf("json marshaling failed\n")
	}
	return fmt.Sprintf("%+v", string(b))
}

func NewTopic(name, state string, partitionsNum int) *CreateTopic {
	return &CreateTopic{
		name,
		state,
		partitionsNum,
	}
}
