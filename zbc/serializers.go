package zbc

import (
	"github.com/zeebe-io/zbc-go/zbc/models/zbdispatch"
	"github.com/zeebe-io/zbc-go/zbc/services/zbsocket"
)

type ProtocolSerializer struct{}

// UnmarshalFromFile will read binary message from disk.
func (ps *ProtocolSerializer) UnmarshalFromFile(path string) (*zbdispatch.Message, error) {
	return zbsocket.NewMessageReader(nil).UnmarshalFromFile(path)
}

func NewProtocolSerializer() *ProtocolSerializer {
	return &ProtocolSerializer{}
}

