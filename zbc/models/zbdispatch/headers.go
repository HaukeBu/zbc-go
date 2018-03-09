package zbdispatch

import (
	"github.com/zeebe-io/zbc-go/zbc/models/zbprotocol"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsbe"
)

// Headers is aggregator for all headers. It holds pointer to every layer. If RequestResponseHeader is nil, then IsSingleMessage will always return true.
type Headers struct {
	FrameHeader           *zbprotocol.FrameHeader
	TransportHeader       *zbprotocol.TransportHeader
	RequestResponseHeader *zbprotocol.RequestResponseHeader // If this is nil then struct is equal to SingleMessage
	SbeMessageHeader      *zbsbe.MessageHeader
}

// SetFrameHeader is a setter for FrameHeader.
func (h *Headers) SetFrameHeader(header *zbprotocol.FrameHeader) {
	h.FrameHeader = header
}

// SetTransportHeader is a setter for TransportHeader.
func (h *Headers) SetTransportHeader(header *zbprotocol.TransportHeader) {
	h.TransportHeader = header
}

// SetRequestResponseHeader is a setting for RequestResponseHeader.
func (h *Headers) SetRequestResponseHeader(header *zbprotocol.RequestResponseHeader) {
	h.RequestResponseHeader = header
}

// IsSingleMessage is helper to determine which model of communication we use.
func (h *Headers) IsSingleMessage() bool {
	return h.RequestResponseHeader == nil
}

// SetSbeMessageHeader is a setter for SBEMessageHeader.
func (h *Headers) SetSbeMessageHeader(header *zbsbe.MessageHeader) {
	h.SbeMessageHeader = header
}
