package test_helpers

import (
	"github.com/zeebe-io/zbc-go/zbc"
	"testing"
	"time"
)

func CreateRandomTopicWithTimeout(t *testing.T, client *zbc.Client) string {
	t.Log("Creating topic")
	hash := RandStringBytes(25)
	topic, err := client.CreateTopic(hash, NumberOfPartitions, 1)
	Assert(t, nil, err, true)
	Assert(t, nil, topic, false)
	t.Logf("Topic %s created with %d partitions", hash, NumberOfPartitions)

	time.Sleep(time.Duration(time.Second * 10))
	return hash
}
