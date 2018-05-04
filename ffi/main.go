package main

/*
struct CreateTopic {
	char* Error;
	char* Name;
	char* State;
	int Partitions;
	int ReplicationFactor;
};
*/
import "C"

import (
	"github.com/zeebe-io/zbc-go/zbc"
	"fmt"
)

//export Add
func Add(a, b int) int { return a + b }

var client *zbc.Client

//export InitClient
func InitClient(bootstrapAddr string) string {
	fmt.Println(bootstrapAddr, len(bootstrapAddr))
	var err error
	client, err = zbc.NewClient(bootstrapAddr)
	if err != nil {
		return err.Error()
	}
	return "CLIENT_RUNNING"
}

//export CreateTopic
func CreateTopic(name string, partitionNum, replicationFactor int) C.struct_CreateTopic {
	fmt.Println(name, partitionNum, replicationFactor)
	var createTopic C.struct_CreateTopic
	if client == nil {
		createTopic.Error = C.CString("client is not initialized")
		return createTopic
	}

	topic, err := client.CreateTopic(name, partitionNum, replicationFactor)
	if err != nil {
		createTopic.Error = C.CString(err.Error())
		return createTopic
	}

	createTopic.Name = C.CString(topic.Name)
	createTopic.State = C.CString(topic.State)
	createTopic.Partitions = C.int(topic.Partitions)
	createTopic.ReplicationFactor = C.int(topic.ReplicationFactor)

	return createTopic
}

//export TopicInfo
func TopicInfo() C.struct_CreateTopic {
	return C.struct_CreateTopic{
		C.CString(""),
		C.CString("default-topic"),
		C.CString("CREATING"),
		1,
		15,
	}
}

func main() {}
