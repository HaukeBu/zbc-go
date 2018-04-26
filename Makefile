VERSION=0.8.0
ZBC_PATH=$(GOPATH)/src/github.com/zeebe-io/zbc-go

cov:
	@cat .coverage/*.txt > coverage.txt
	@rm .coverage/*.txt

.PHONY: dir-setup
dir-setup:
	@mkdir -p .coverage .profiles .traces .logs

test-protocol: dir-setup
	go test -cpuprofile .profiles/test-protocol-cpu.prof -memprofile .profiles/test-protocol-mem.prof -race -trace .traces/test-protocol-trace.out -coverprofile=.coverage/protocol.txt -covermode=atomic zbc/models/zbprotocol/*.go -v

test-sbe: dir-setup
	go test -cpuprofile .profiles/test-sbe-cpu.prof -memprofile .profiles/test-sbe-mem.prof -race -trace .traces/test-sbe-trace.out -coverprofile=.coverage/sbe.txt -covermode=atomic zbc/models/zbsbe/*.go -v

test-msgpack: dir-setup
	go test -cpuprofile .profiles/test-msgpack-cpu.prof -memprofile .profiles/test-msgpack-mem.prof -race -trace .traces/test-msgpack-trace.out -coverprofile=.coverage/msgpack.txt -covermode=atomic zbc/models/zbmsgpack/*.go -v

test-hexdump: dir-setup
	go test -cpuprofile .profiles/test-hexdump-cpu.prof -memprofile .profiles/test-hexdump-mem.prof -race -trace .traces/test-hexdump-trace.out -coverprofile=.coverage/hexdump.txt -covermode=atomic tests/test-zbdump/*.go -v

test-exchange: dir-setup
	go test -cpuprofile .profiles/test-exchange-cpu.prof -memprofile .profiles/test-exchange-mem.prof -race -trace .traces/test-exchange-trace.out -coverprofile=.coverage/exchange.txt -covermode=atomic tests/test-exchange/*.go -v

test-tasksub: dir-setup
	go test -cpuprofile .profiles/test-tasksub-cpu.prof -memprofile .profiles/test-tasksub-mem.prof -race -trace .traces/test-tasksub-trace.out -coverprofile=.coverage/tasksub.txt -covermode=atomic tests/test-task-subscriptions/*.go -v

test-topicsub: dir-setup
	go test -cpuprofile .profiles/test-topicsub-cpu.prof -memprofile .profiles/test-topicsub-mem.prof -race -trace .traces/test-topicsub-trace.out -coverprofile=.coverage/topicsub.txt -covermode=atomic tests/test-topic-subscriptions/*.go -v

test-socket: dir-setup
	go test -cpuprofile .profiles/test-socket-cpu.prof -memprofile .profiles/test-socket-mem.prof -race -trace .traces/test-socket-trace.out -coverprofile=.coverage/socket.txt -covermode=atomic tests/test-socket/*.go -v

test-topology: dir-setup
	go test -cpuprofile .profiles/test-topology-cpu.prof -memprofile .profiles/test-topology-mem.prof -race -trace .traces/test-topology-trace.out -coverprofile=.coverage/topology.txt -covermode=atomic tests/test-topology/*.go -v

test-setup: dir-setup
	go test -cpuprofile .profiles/test-cpu.prof -memprofile .profiles/test-mem.prof -race -trace .traces/test-trace.out -coverprofile=.coverage/setup.txt -covermode=atomic tests/*.go -v

test:
	make test-setup
	make test-socket
	make test-topology
	make test-exchange
	make test-tasksub
	make test-topicsub
	make test-protocol
	make test-sbe
	make test-msgpack
	make test-hexdump


clean:
	@rm -rf .profiles .coverage .traces .logs
