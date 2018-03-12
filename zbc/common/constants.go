package zbcommon

import "time"

// RequestTimeout specifies default timeout for responder.
const RequestTimeout = 5

// TopologyRefreshInterval defines time to live of cluster topology object.
const TopologyRefreshInterval = 60

// Retry constants
const (
	BackoffMin      = 20 * time.Millisecond
	BackoffMax      = 1000 * time.Millisecond
	BackoffDeadline = 30 * time.Second
)

// Sbe template ID constants
const (
	TemplateIDExecuteCommandRequest  = 20
	TemplateIDExecuteCommandResponse = 21
	TemplateIDControlMessageResponse = 11
	TemplateIDSubscriptionEvent      = 30
)

// Zeebe protocol constants
const (
	FrameHeaderSize           = 12
	TransportHeaderSize       = 2
	RequestResponseHeaderSize = 8
	SBEMessageHeaderSize      = 8

	TotalHeaderSizeNoFrame = 18
	TotalHeaderSize        = 30
	LengthFieldSize        = 2
)

// Message pack states for Task, Deployment and WorkflowInstance
const (
	TaskCreate  = "CREATE"
	TaskCreated = "CREATED"

	TaskComplete  = "COMPLETE"
	TaskCompleted = "COMPLETED"

	TaskFail           = "FAIL"
	TaskFailed         = "FAILED"
	TaskFailedRejected = "FAIL_REJECTED"

	CreateDeployment   = "CREATE"
	DeploymentCreated  = "CREATED"
	DeploymentRejected = "REJECTED"

	CreateWorkflowInstance   = "CREATE_WORKFLOW_INSTANCE"
	WorkflowInstanceCreated  = "WORKFLOW_INSTANCE_CREATED"
	WorkflowInstanceRejected = "WORKFLOW_INSTANCE_REJECTED"
)

//
const (
	TopicCreate   = "CREATE"
	TopicCreated  = "CREATED"
	TopicRejected = "CREATE_REJECTED"
)

// TopicSubscription states
const (
	TopicSubscriptionSubscribeState  = "SUBSCRIBE"
	TopicSubscriptionSubscribedState = "SUBSCRIBED"
)

// TopicSubscriptionAck states
const (
	TopicSubscriptionAckState          = "ACKNOWLEDGE"
	TopicSubscriptionAcknowledgedState = "ACKNOWLEDGED"
)

const (
	SystemTopic = "internal-system"
)

// Workflow resource types
const (
	BpmnXml      = "BPMN_XML"
	YamlWorkflow = "YAML_WORKFLOW"
)

const (
	SocketChunkSize = 4096
)

const RequestQueueSize uint64 = 4096
const SubscriptionPipelineQueueSize = ^uint16(0)

const StateLeader = "LEADER"


const TaskSubscriptionRefreshCreditsThreshold float32 = 0.3
const TopicSubscriptionAckThreshold float32 = 0.7
