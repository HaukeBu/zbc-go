package zbdispatch

import (
	"encoding/json"
	"fmt"

	"github.com/zeebe-io/zbc-go/zbc/models/zbsbe"
)

// Message is Zeebe message which will contain pointers to all parts of the message. Data is Message Pack layer.
type Message struct {
	Headers    *Headers
	SbeMessage SimpleBinaryEncoding
	Data       []byte
}

// SetHeaders is a setter for Headers attribute.
func (m *Message) SetHeaders(headers *Headers) {
	m.Headers = headers
}

// SetSbeMessage is a setter for SBE attribute.
func (m *Message) SetSbeMessage(data SimpleBinaryEncoding) {
	m.SbeMessage = data
}

// SetData is a setter for unmarshaled message pack data.
func (m *Message) SetData(data []byte) {
	m.Data = data
}

func (m *Message) ForPartitionId() *uint16 {
	switch sbe := (m.SbeMessage).(type) {
	case *zbsbe.ControlMessageRequest:
		return &sbe.PartitionId
	case *zbsbe.ExecuteCommandRequest:
		return &sbe.PartitionId
	default:
		return nil
	}
}

func (m *Message) IsTopologyMessage() bool {
	switch sbe := (m.SbeMessage).(type) {
	case *zbsbe.ControlMessageRequest:
		if sbe.MessageType == zbsbe.ControlMessageType.REQUEST_TOPOLOGY {
			return true
		}
	}
	return false
}

func (m *Message) JsonString(data interface{}) string {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Sprintf("json marshaling failed\n")
	}
	return fmt.Sprintf("%+v", string(b))
}
