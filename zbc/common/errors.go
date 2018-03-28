package zbcommon

import "errors"

// TODO: Organize errors by services

// Reader Errors
var (
	ErrFrameHeaderRead    = errors.New("cannot read bytes for frame header")
	ErrFrameHeaderDecode  = errors.New("cannot decode bytes into frame header")
	ErrProtocolIDNotFound = errors.New("ProtocolID not found")
	ErrSocketRead         = errors.New("failed to read requested number of bytes")
)

// Client Errors
var (
	ErrTimeout                = errors.New("request timeout")
	ErrTopicLeaderNotFound    = errors.New("topic leader not found")
	ErrTopicPartitionNotFound = errors.New("topic partition not found")
	ErrResourceNotFound       = errors.New("resource not found")
)

// Socket Errors
var (
	ErrSocketWrite    = errors.New("tried to write more bytes to socket")
	ErrConnectionDead = errors.New("connection is dead")
)

// Topology Service
var (
	ErrPartitionsNotFound = errors.New("partitions for topic not found")
	ErrNoBrokersFound     = errors.New("broker not found")
)

// RetryDeadlineReached is error which occurs when requestWrapper failed multiple times unsuccessfully.
var RetryDeadlineReached = errors.New("message retry deadline reached")

var (
	// BrokerNotFound is used in transport manager to denote that broker cannot be contacted
	BrokerNotFound             = errors.New("cannot contact the broker")
	ErrNoPartitionsFound       = errors.New("partitions not found")
	ErrRoundRobinCtrlNotFound  = errors.New("round robin controller not initialized")
	ErrDestinationAddrNotFound = errors.New("destination address could not be determined")
	ErrClusterInfoNotFound     = errors.New("cluster not initialized")
)

var ErrFailedToOpenTaskSubscription = errors.New("failed to open task subscription")

// TaskSubscriptionCtrl errors
var (
	ErrCallbackNotAttached = errors.New("callback not attached to the controler")
)

// Socket dispatcher errors
var (
	ErrWrongTransactionIndex = errors.New("wrong transaction index, possible problem with RequestID")
	ErrWrongSubscriptionKey  = errors.New("cannot dispatch subscriptions event, possible problem with SubscriberKey")
)

// Round Robin Controller
var (
	ErrClusterPartialInformation = errors.New("round robin cluster information is not complete. try refreshing topology")
)

var (
	ErrSubscriptionPipelineFailed = errors.New("opening on all partitions failed")
	ErrSubscriptionClosed         = errors.New("error processing event in handler, subscription closed")
)

var (
	ErrEventNotTask = errors.New("provided events is not a task")
	ErrEventIsEmpty = errors.New("event is empty")
)

var (
	SubscriptionIsClosed = errors.New("subscription is closed")
)
