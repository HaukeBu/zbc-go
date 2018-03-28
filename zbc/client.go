package zbc

import (
	"github.com/zeebe-io/zbc-go/zbc/models/zbmsgpack"
	"github.com/zeebe-io/zbc-go/zbc/services/zbexchange"
	"github.com/zeebe-io/zbc-go/zbc/services/zbsubscribe"

	"github.com/vmihailenco/msgpack"
	"github.com/zeebe-io/zbc-go/zbc/common"
)

// Client for Zeebe broker with support for clustered deployment.
type Client struct {
	zbsubscribe.ZeebeAPI
}

// NewClient is constructor for Client structure. It will resolve IP address and dial the provided tcp address.
func NewClient(bootstrapAddr string) (*Client, error) {
	if zbcommon.ZBL == nil {
		zbcommon.InitLogger()
	}

	exchangeSvc := zbexchange.NewExchangeSvc(bootstrapAddr)
	topology, err := exchangeSvc.RefreshTopology()
	if err != nil {
		return nil, err
	}
	if topology == nil {
		return nil, zbcommon.ErrNoBrokersFound
	}

	c := &Client{
		ZeebeAPI: zbsubscribe.NewSubscriptionSvc(exchangeSvc),
	}

	// MARK: Initialize subscriptions client for callbacks
	zbsubscribe.SetClientInstance(c)
	return c, nil
}

// NewTask is constructor for Task object. Function signature denotes mandatory fields.
func NewTask(typeName string) *zbmsgpack.Task {
	return &zbmsgpack.Task{
		State:        zbcommon.TaskCreate,
		Headers:      make(map[string]interface{}),
		CustomHeader: make(map[string]interface{}),

		Type:    typeName,
		Retries: 3,
	}
}

// NewWorkflowInstance will create new workflow instance.
func NewWorkflowInstance(bpmnProcessID string, version int, payload interface{}) *zbmsgpack.WorkflowInstance {
	b, err := msgpack.Marshal(payload)
	if err != nil {
		return nil
	}

	return &zbmsgpack.WorkflowInstance{
		State:         zbcommon.CreateWorkflowInstance,
		BPMNProcessID: bpmnProcessID,
		Version:       version,
		Payload:       b,
	}
}

// NewResource will create new message pack resource.
func NewResource(resourceName, resourceType string, resource []byte) *zbmsgpack.Resource {
	return &zbmsgpack.Resource{
		ResourceName: resourceName,
		ResourceType: resourceType,
		Resource:     resource,
	}
}
