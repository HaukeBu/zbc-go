package zbsocket

import (
	"bytes"
	"encoding/binary"
	"io"
	"io/ioutil"

	"github.com/zeebe-io/zbc-go/zbc/models/zbdispatch"
	"github.com/zeebe-io/zbc-go/zbc/models/zbprotocol"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsbe"

	"github.com/zeebe-io/zbc-go/zbc/common"
)

// MessageReader is builder which will read byte array and construct zbdispatch.Message with all their parts.
type MessageReader struct {
	*socket
}

func (mr *MessageReader) readFrameHeader(data io.Reader) (*zbprotocol.FrameHeader, error) {
	var frameHeader zbprotocol.FrameHeader
	if frameHeader.Decode(data, binary.LittleEndian, 0) != nil {
		return nil, zbcommon.ErrFrameHeaderDecode
	}
	return &frameHeader, nil
}

func (mr *MessageReader) readTransportHeader(data io.Reader) (*zbprotocol.TransportHeader, error) {
	var transport zbprotocol.TransportHeader
	err := transport.Decode(data, binary.LittleEndian, 0)
	if err != nil {
		return nil, err
	}
	if transport.ProtocolID == zbprotocol.RequestResponse || transport.ProtocolID == zbprotocol.FullDuplexSingleMessage {
		return &transport, nil
	}

	return nil, zbcommon.ErrProtocolIDNotFound
}

func (mr *MessageReader) readRequestResponseHeader(data io.Reader) (*zbprotocol.RequestResponseHeader, error) {
	var requestResponse zbprotocol.RequestResponseHeader
	err := requestResponse.Decode(data, binary.LittleEndian, 0)
	if err != nil {
		return nil, err
	}

	return &requestResponse, nil
}

func (mr *MessageReader) readSbeMessageHeader(data io.Reader) (*zbsbe.MessageHeader, error) {
	var sbeMessageHeader zbsbe.MessageHeader
	err := sbeMessageHeader.Decode(data, binary.LittleEndian, 0)
	if err != nil {
		return nil, err
	}
	return &sbeMessageHeader, nil
}

func (mr *MessageReader) readMessage(data []byte) (*zbdispatch.Message, error) {
	var header zbdispatch.Headers

	headerByte := data[0:zbcommon.FrameHeaderSize]
	frameHeader, err := mr.readFrameHeader(bytes.NewReader(headerByte))

	if err != nil {
		return nil, err
	}
	if frameHeader.Length < zbcommon.TotalHeaderSizeNoFrame {
		return nil, zbcommon.ErrFrameHeaderRead
	}
	header.SetFrameHeader(frameHeader)
	message := data[zbcommon.FrameHeaderSize : int(frameHeader.Length)+zbcommon.FrameHeaderSize]

	if int(frameHeader.Length) != len(message) || len(message) == 0 {
		return nil, zbcommon.ErrFrameHeaderRead
	}

	transportReader := bytes.NewReader(message[:zbcommon.TransportHeaderSize])
	transport, err := mr.readTransportHeader(transportReader)

	if err != nil {
		return nil, err
	}
	header.SetTransportHeader(transport)

	sbeIndex := zbcommon.TransportHeaderSize
	switch transport.ProtocolID {
	case zbprotocol.RequestResponse:
		reqRespReader := bytes.NewReader(message[zbcommon.TransportHeaderSize : zbcommon.TransportHeaderSize+zbcommon.RequestResponseHeaderSize])
		requestResponse, errHeader := mr.readRequestResponseHeader(reqRespReader)
		if errHeader != nil {
			return nil, err
		}
		header.SetRequestResponseHeader(requestResponse)
		sbeIndex = zbcommon.TransportHeaderSize + zbcommon.RequestResponseHeaderSize
		break

	case zbprotocol.FullDuplexSingleMessage:
		header.SetRequestResponseHeader(nil)
		break
	}

	sbeHeaderReader := bytes.NewReader(message[sbeIndex : sbeIndex+zbcommon.SBEMessageHeaderSize])
	sbeMessageHeader, err := mr.readSbeMessageHeader(sbeHeaderReader)
	if err != nil {
		return nil, err
	}
	header.SetSbeMessageHeader(sbeMessageHeader)

	body := message[sbeIndex+8:]

	currentSize := int64(frameHeader.Length) + zbcommon.FrameHeaderSize
	expectedSize := (currentSize + 7) & ^7
	data = data[expectedSize:]

	return mr.parseMessage(&header, &body)
}

// readHeaders will read entire message and interpret all headers. It will return pointer to headers object and tail of the message as byte array.
func (mr *MessageReader) readHeaders() (*zbdispatch.Headers, *[]byte, error) {
	var header zbdispatch.Headers

	headerByte := mr.getBytes(0, zbcommon.FrameHeaderSize)
	frameHeader, err := mr.readFrameHeader(bytes.NewReader(headerByte))
	frameHeader.Length = frameHeader.Length - zbcommon.FrameHeaderSize

	if err != nil {
		return nil, nil, err
	}
	if frameHeader.Length < zbcommon.TotalHeaderSizeNoFrame {
		return nil, nil, zbcommon.ErrFrameHeaderRead
	}
	header.SetFrameHeader(frameHeader)
	message := mr.getBytes(zbcommon.FrameHeaderSize, int(frameHeader.Length)+zbcommon.FrameHeaderSize)

	if int(frameHeader.Length) != len(message) || len(message) == 0 {
		return nil, nil, zbcommon.ErrFrameHeaderRead
	}

	transportReader := bytes.NewReader(message[:zbcommon.TransportHeaderSize])
	transport, err := mr.readTransportHeader(transportReader)
	if err != nil {
		return nil, nil, err
	}
	header.SetTransportHeader(transport)

	sbeIndex := zbcommon.TransportHeaderSize
	switch transport.ProtocolID {
	case zbprotocol.RequestResponse:
		reqRespReader := bytes.NewReader(message[zbcommon.TransportHeaderSize : zbcommon.TransportHeaderSize+zbcommon.RequestResponseHeaderSize])
		requestResponse, errHeader := mr.readRequestResponseHeader(reqRespReader)
		if errHeader != nil {
			return nil, nil, err
		}
		header.SetRequestResponseHeader(requestResponse)
		sbeIndex = zbcommon.TransportHeaderSize + zbcommon.RequestResponseHeaderSize
		break

	case zbprotocol.FullDuplexSingleMessage:
		header.SetRequestResponseHeader(nil)
		break
	}

	sbeHeaderReader := bytes.NewReader(message[sbeIndex : sbeIndex+zbcommon.SBEMessageHeaderSize])
	sbeMessageHeader, err := mr.readSbeMessageHeader(sbeHeaderReader)
	if err != nil {
		return nil, nil, err
	}
	header.SetSbeMessageHeader(sbeMessageHeader)

	body := message[sbeIndex+8:]

	currentSize := int64(frameHeader.Length) + zbcommon.FrameHeaderSize
	expectedSize := (currentSize + 7) & ^7
	mr.popBytes(int(expectedSize))

	return &header, &body, nil
}

func (mr *MessageReader) decodeCmdRequest(reader *bytes.Reader, header *zbsbe.MessageHeader) (*zbsbe.ExecuteCommandRequest, error) {
	var commandRequest zbsbe.ExecuteCommandRequest

	err := commandRequest.Decode(reader,
		binary.LittleEndian,
		header.Version,
		header.BlockLength,
		true)
	if err != nil {
		return nil, err
	}
	return &commandRequest, nil
}

func (mr *MessageReader) decodeCmdResponse(reader *bytes.Reader, header *zbsbe.MessageHeader) (*zbsbe.ExecuteCommandResponse, error) {
	var commandResponse zbsbe.ExecuteCommandResponse
	err := commandResponse.Decode(reader,
		binary.LittleEndian,
		header.Version,
		header.BlockLength,
		true)
	if err != nil {
		return nil, err
	}
	return &commandResponse, nil
}

func (mr *MessageReader) decodeCtlResponse(reader *bytes.Reader, header *zbsbe.MessageHeader) (*zbsbe.ControlMessageResponse, error) {
	var controlResponse zbsbe.ControlMessageResponse
	err := controlResponse.Decode(reader, binary.LittleEndian, header.Version, header.BlockLength, true)
	if err != nil {
		return nil, err
	}
	return &controlResponse, nil
}

func (mr *MessageReader) decodeSubEvent(reader *bytes.Reader, header *zbsbe.MessageHeader) (*zbsbe.SubscribedEvent, error) {
	var subEvent zbsbe.SubscribedEvent
	err := subEvent.Decode(reader, binary.LittleEndian, header.Version, header.BlockLength, true)
	if err != nil {
		return nil, err
	}
	return &subEvent, nil
}

// parseMessage will take the headers and tail and construct zbdispatch.Message.
func (mr *MessageReader) parseMessage(headers *zbdispatch.Headers, message *[]byte) (*zbdispatch.Message, error) {
	var msg zbdispatch.Message
	msg.SetHeaders(headers)
	reader := bytes.NewReader(*message)

	switch headers.SbeMessageHeader.TemplateId {

	case zbcommon.TemplateIDExecuteCommandRequest: // Testing purposes.
		commandRequest, err := mr.decodeCmdRequest(reader, headers.SbeMessageHeader)
		if err != nil {
			return nil, err
		}
		msg.SetSbeMessage(commandRequest)
		msg.SetData([]byte(commandRequest.Command))

		break

	case zbcommon.TemplateIDExecuteCommandResponse: // Read response from the socket.
		commandResponse, err := mr.decodeCmdResponse(reader, headers.SbeMessageHeader)
		if err != nil {
			return nil, err
		}
		msg.SetSbeMessage(commandResponse)
		msg.SetData([]byte(commandResponse.Event))

		break

	case zbcommon.TemplateIDControlMessageResponse:
		ctlResponse, err := mr.decodeCtlResponse(reader, headers.SbeMessageHeader)
		if err != nil {
			return nil, err
		}
		msg.SetSbeMessage(ctlResponse)
		msg.SetData([]byte(ctlResponse.Data))

		break

	case zbcommon.TemplateIDSubscriptionEvent:
		subscribedEvent, err := mr.decodeSubEvent(reader, headers.SbeMessageHeader)
		if err != nil {
			return nil, err
		}
		msg.SetSbeMessage(subscribedEvent)
		msg.SetData([]byte(subscribedEvent.Event))

		break
	}
	return &msg, nil
}

// UnmarshalFromFile will read binary message from disk.
func (mr *MessageReader) UnmarshalFromFile(path string) (*zbdispatch.Message, error) {
	data, fsErr := ioutil.ReadFile(path)
	if fsErr != nil {
		return nil, fsErr
	}
	return mr.readMessage(data)
}

// NewMessageReader is constructor for MessageReader builder.
func NewMessageReader(socket *socket) *MessageReader {
	return &MessageReader{
		socket,
	}
}
