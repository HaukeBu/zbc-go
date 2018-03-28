package zbsocket

import (
	"bytes"
	"net"
	"time"

	"github.com/zeebe-io/zbc-go/zbc/models/zbdispatch"
	"github.com/zeebe-io/zbc-go/zbc/models/zbsbe"

	"github.com/zeebe-io/zbc-go/zbc/common"
)

type socket struct {
	*messageDispatcher

	connection net.Conn
	stream     []byte
	closeCh    chan bool
}

func (s *socket) Send(request *RequestWrapper) error {
	if s.connection == nil {
		return zbcommon.ErrConnectionDead
	}

	s.AddTransaction(request) // INFO: AddTransaction will write to payload for RequestID
	request.Sock = s

	writer := NewMessageWriter(request.Payload)
	byteBuff := &bytes.Buffer{}
	writer.Write(byteBuff)
	n, err := s.connection.Write(byteBuff.Bytes())
	if err != nil {
		return err
	}

	if n != len(byteBuff.Bytes()) {
		return zbcommon.ErrSocketWrite
	}

	return nil
}

func (s *socket) receiver() {
	reader := NewMessageReader(s)
	responseHandler := zbdispatch.ResponseHandler{}

	for {
		select {
		case <-s.closeCh:
			s.connection.Close()
			return

		default:
			headers, tail, err := reader.readHeaders()
			if err != nil {
				continue
			}
			message, err := reader.parseMessage(headers, tail)

			if err != nil && !headers.IsSingleMessage() {
				continue
			}

			if !headers.IsSingleMessage() && message != nil && len(message.Data) > 0 {
				s.DispatchTransaction(headers.RequestResponseHeader.RequestID, message)
				continue
			}

			if err != nil && headers.IsSingleMessage() {
				continue
			}

			if headers.IsSingleMessage() && message != nil && len(message.Data) > 0 {
				event := message.SbeMessage.(*zbsbe.SubscribedEvent)
				partitionID := message.SbeMessage.(*zbsbe.SubscribedEvent).PartitionId
				zbcommon.ZBL.Debug().Str("component", "socket").Msgf("subscriptions event from partitionID %d, subscriptions type %d", partitionID, event.SubscriptionType)
				switch event.SubscriptionType {
				case zbsbe.SubscriptionType.TASK_SUBSCRIPTION:
					task := responseHandler.UnmarshalTask(message)
					if task == nil {
						zbcommon.ZBL.Error().Str("component", "socket").Msg("task is nil")
						continue
					}
					zbcommon.ZBL.Debug().Str("component", "socket").Msg("task received -> dispatching task")
					err := s.DispatchTaskEvent(event.SubscriberKey, event, task)
					if err != nil {
						zbcommon.ZBL.Debug().Str("component", "socket").Msg("dispatching failed")
					} else {
						zbcommon.ZBL.Debug().Str("component", "socket").Msg("task dispatched")
					}

					break
				case zbsbe.SubscriptionType.TOPIC_SUBSCRIPTION:
					zbcommon.ZBL.Debug().Str("component", "socket").Msg("event received -> dispatching event")
					err := s.DispatchTopicEvent(event.SubscriberKey, event)
					if err != nil {
						zbcommon.ZBL.Error().Str("component", "socket").Msgf("Error while dispatching topic event: %v", err)
					}
					zbcommon.ZBL.Debug().Str("component", "socket").Msg("event dispatched")
					break
				default:
					zbcommon.ZBL.Debug().Str("component", "socket").Msg("undispatchable message received")
				}

			}
		}
	}
}

func (s *socket) teardown() {
	close(s.closeCh)
}

func (s *socket) readChunk() {
	for {
		s.connection.SetReadDeadline(time.Now().Add(time.Millisecond * 100))

		total := make([]byte, zbcommon.SocketChunkSize)
		returnedSize, err := s.connection.Read(total)
		if err != nil {
			continue
		}

		s.stream = append(s.stream, total[:returnedSize]...)
		break
	}
}

func (s *socket) getBytes(start, end int) []byte {
	for {
		if end-start > len(s.stream) || end > len(s.stream) {
			s.readChunk()
		} else {
			break
		}
	}

	frame := s.stream[start:end]
	return frame
}

func (s *socket) popBytes(pos int) {
	s.stream = s.stream[pos:]
}

func (s *socket) Dial(addr string) error {
	tcpAddr, wrongAddr := net.ResolveTCPAddr("tcp4", addr)
	if wrongAddr != nil {
		return wrongAddr
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}

	s.connection = conn

	go s.receiver()

	return nil
}

func NewSocket(addr string) LikeSocket {
	ss := &socket{
		messageDispatcher: newMessageDispatcher(),
		connection:        nil,
		stream:            make([]byte, 0),
		closeCh:           make(chan bool),
	}

	return ss
}
