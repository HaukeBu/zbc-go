package tests

import (
	"testing"

	. "github.com/zeebe-io/zbc-go/tests/test-helpers"
	"github.com/zeebe-io/zbc-go/zbc"
	"github.com/zeebe-io/zbc-go/zbc/common"
)

func TestCreateTopic(t *testing.T) {
	zbClient, err := zbc.NewClient(BrokerAddr)
	Assert(t, nil, err, true)
	Assert(t, nil, zbClient, false)

	hash := RandStringBytes(25)
	topic, err := zbClient.CreateTopic(hash, NumberOfPartitions, 1)
	Assert(t, nil, err, true)
	Assert(t, nil, topic, false)

	Assert(t, zbcommon.TopicCreating, topic.State, true)

	topic, _ = zbClient.CreateTopic("default-topic", NumberOfPartitions, 1)
	Assert(t, nil, topic, false)
}
