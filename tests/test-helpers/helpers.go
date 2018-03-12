package test_helpers

import (
	"errors"
	"math/rand"
	"net"
	"os"
	"reflect"
	"runtime/debug"
	"strconv"
	"testing"
	"time"
)

var errClientStartFailed = errors.New("cannot connect to the broker")

const TopicName = "default-topic"

// To change package variables export specified environment variables:
// export ZB_BROKER_ADDR="192.168.0.11:51015"
// export ZB_DEFAULT_NUM_PARTITIONS="5"

var (
	BrokerAddr         = "0.0.0.0:51015"
	NumberOfPartitions = 1
)

func init() {
	brokerAddr := os.Getenv("ZB_BROKER_ADDR")
	if len(brokerAddr) > 0 {
		// validate if it's valid URI
		ip := net.ParseIP(brokerAddr)
		if ip != nil {
			// TODO: validate if it's valid IP
		}
		BrokerAddr = brokerAddr
	}

	numberOfPartitions := os.Getenv("ZB_DEFAULT_NUM_PARTITIONS")
	if len(numberOfPartitions) > 0 {
		num, err := strconv.Atoi(numberOfPartitions)
		if err != nil {
			panic("wrong ZB_DEFAULT_NUM_PARTITIONS value")
		}
		NumberOfPartitions = num

	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringBytes(n int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func Assert(t *testing.T, exp, got interface{}, equal bool) {
	if reflect.DeepEqual(exp, got) != equal {
		debug.PrintStack()
		t.Fatalf("\x1b[31;1mExpecting '%v' got '%v'\x1b[0m\n", exp, got)
	}
}

func CheckRoundRobinSequence(order uint16, sequence []uint16) bool {
	sequenced := true
	last := sequence[0]
	for i := 1; i < len(sequence); i++ {
		if last == order && sequence[i] != 0 {
			sequenced = false
			break
		}
		if last != order && last+1 != sequence[i] {
			sequenced = false
			break
		}
		last = sequence[i]
	}
	return sequenced
}
