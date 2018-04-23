VERSION=0.8.0
BINARY_NAME=zbctl
ZBC_PATH=$(GOPATH)/src/github.com/zeebe-io/zbc-go
PREFIX=/usr/local

build:
	@mkdir -p target/bin
	@go build -o target/bin/$(BINARY_NAME) ./cmd/*.go
	@cp cmd/config.toml target/bin/

.PHONY: install
install:
	mkdir -p $(PREFIX)/zeebe
	install -m 644 target/bin/config.toml $(PREFIX)/zeebe/
	install -m 755 target/bin/zbctl $(PREFIX)/zeebe/
	ln -sf $(PREFIX)/zeebe/zbctl $(PREFIX)/bin/zbctl

.PHONY: uninstall
uninstall:
	rm -rf $(PREFIX)/bin/zbctl
	rm -rf $(PREFIX)/zeebe

run:
	@go run cmd/main.go

cov:
	cat .coverage/*.txt > coverage.txt
	rm .coverage/*.txt

.PHONY: cov-setup
cov-setup:
	mkdir -p .coverage

test-protocol: cov-setup
	go test -race -coverprofile=.coverage/protocol.txt -covermode=atomic zbc/models/zbprotocol/*.go -v

test-sbe: cov-setup
	go test -race -coverprofile=.coverage/sbe.txt -covermode=atomic zbc/models/zbsbe/*.go -v

test-msgpack: cov-setup
	go test -race -coverprofile=.coverage/msgpack.txt -covermode=atomic zbc/models/zbmsgpack/*.go -v

test-hexdump: cov-setup
	go test -race -coverprofile=.coverage/hexdump.txt -covermode=atomic tests/test-zbdump/*.go -v

test-exchange: cov-setup
	go test -race -coverprofile=.coverage/exchange.txt -covermode=atomic tests/test-exchange/*.go -v

test-tasksub: cov-setup
	go test -race -coverprofile=.coverage/tasksub.txt -covermode=atomic tests/test-task-subscriptions/*.go -v

test-topicsub: cov-setup
	go test -race -coverprofile=.coverage/topicsub.txt -covermode=atomic tests/test-topic-subscriptions/*.go -v

test-socket: cov-setup
	go test -race -coverprofile=.coverage/socket.txt -covermode=atomic tests/test-socket/*.go -v

test-topology: cov-setup
	go test -race -coverprofile=.coverage/topology.txt -covermode=atomic tests/test-topology/*.go -v

test-setup: cov-setup
	go test -race -coverprofile=.coverage/setup.txt -covermode=atomic tests/*.go -v

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
	@rm -rf ./target *.tar.gz $(BINARY_NAME)
