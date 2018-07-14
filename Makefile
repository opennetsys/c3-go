all: build

.PHONY: install
install:
	@go get -u github.com/c3systems/c3

.PHONY: deps
deps:
	@rm -rf ./vendor && \
	dep ensure && \
	gxundo ./vendor && \
	(cd vendor/github.com/libp2p/go-libp2p-pubsub/pb \
	&& rm rpc.pb.go && rm rpc.proto \
	&& wget https://github.com/c3systems/go-libp2p-pubsub/raw/master/pb/rpc.pb.go \
	&& wget https://github.com/c3systems/go-libp2p-pubsub/raw/master/pb/rpc.proto)


.PHONY: build
build:
	@go build -v -ldflags "-s -w" -o bin/c3 .

.PHONY: build/mac
build/mac: clean/mac
	@env GOARCH=amd64 go build -ldflags "-s -w" -o build/macos/c3 && upx build/macos/c3

.PHONY: build/linux
build/linux: clean/linux
	@env GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o build/linux/c3 && upx build/linux/c3

.PHONY: clean/mac
clean/mac:
	@go clean && \
	rm -rf build/mac

.PHONY: clean/linux
clean/linux:
	@go clean && \
	rm -rf build/linux

.PHONY: clean
clean:
	@go clean && \
	rm -rf bin/

.PHONY: ipfs/daemon
ipfs/daemon:
	@ipfs daemon

.PHONY: localhostproxy
localhostproxy:
	# proxy localhost to 123.123.123.123 required so that docker container can communicate with host machine
	@sudo ifconfig lo0 alias 123.123.123.123/24

.PHONY: test/check
test/check:
	# TODO: kill script if required commands and daemons not found
	@command -v ipfs &>/dev/null || echo "IPFS is required"
	@command -v docker &>/dev/null || echo "Docker daemon is required"
	@pgrep -f ipfs > /dev/null || echo "IPFS daemon is not running"
	@pgrep -f docker > /dev/null || echo "Docker daemon is not running"

.PHONY: test/cleanup
test/cleanup:
	@chmod +x scripts/test_cleanup.sh
	@. scripts/test_cleanup.sh

.PHONY: test
test: test/check test/c3 test/common test/common test/registry test/core test/node test/cleanup

.PHONY: test/c3
test/c3:
	@go test -v c3/*.go $(ARGS)

.PHONY: test/common
test/common: test/common/network test/common/stringutil

.PHONY: test/common/network
test/common/network:
	@go test -v common/network/*.go $(ARGS)

.PHONY: test/common/stringutil
test/common/stringutil:
	@go test -v common/stringutil/*.go $(ARGS)

.PHONY: test/common/command
test/common/command:
	@go test -v common/command/*.go $(ARGS)

.PHONY: test/core
test/core: test/core/server test/core/docker test/core/ipfs test/core/sandbox

.PHONY: test/core/server
test/core/server:
	@go test -v core/server/*.go $(ARGS)

.PHONY: test/core/docker
test/core/docker:
	@go test -v -parallel 1 core/docker/*.go $(ARGS)

.PHONY: test/core/ipfs
test/core/ipfs:
	@go test -v -parallel 1 core/ipfs/*.go $(ARGS)

.PHONY: test/core/sandbox
test/core/sandbox: docker/build/example
	@IMAGEID=$$(docker images -q | grep -m1 "") go test -v -parallel 1 core/sandbox/*.go $(ARGS)

.PHONY: test/core/chain
test/core/chain: test/core/chain/mainchain test/core/chain/statechain
	@echo "done"

.PHONY: test/core/chain/mainchain
test/core/chain/mainchain:
	@go test -v core/chain/mainchain/*.go $(ARGS)

.PHONY: test/core/chain/statechain
test/core/chain/statechain:
	@go test -v core/chain/statechain/*.go $(ARGS)

.PHONY: test/registry
test/registry:
	@docker pull hello-world && \
	go test -v -parallel 1 registry/*.go $(ARGS)

.PHONY: test/node
test/node:
	@go test -v node/*.go $(ARGS)

.PHONY: run/node
run/node:
	go run main.go node start --pem=node/test_data/key.pem

.PHONY: test/docker/build/snapshot
test/docker/build/snapshot:
	@docker build --no-cache -f snapshot_test/Dockerfile -t snapshot_test:1 ./snapshot_test

.PHONY: test/docker/run/snapshot
demo/docker/run/snapshot:
	@docker run -t snapshot_test:1

.PHONY: test/docker/run/snapshot/daemon
demo/docker/run/snapshot/daemon:
	@docker run -d snapshot_test:1

.PHONY: test/run/snapshot
demo/run/snapshot:
	@node snapshot_test/index.js

.PHONY: docker/run/localregistry
docker/run/localregistry:
	@docker run -d -p 5000:5000 --restart=always --name registry registry:2

.PHONY: docker/push/localregistry
docker/push/localregistry:
	@docker push localhost:5000/$(IMAGE)

.PHONY: docker/list/localregistry
docker/list/localregistry:
	@curl -X GET -k "http://$(docker-machine ip):5000/v2/_catalog"

.PHONY: docker/gcr/images/list
docker/gcr/images/list:
	@curl http://gcr.c3labs.io:5000/v2/_catalog

.PHONY: run/example
run/example:
	@go run example/go/main.go

.PHONY: docker/build/example
docker/build/example:
	@docker build --no-cache -f example/go/Dockerfile -t goexample ./example

.PHONY: docker/run/example
docker/run/example:
	@docker run -p 3333 --mount type=bind,src=/tmp/633029102,target=/tmp -t goexample

.PHONY: docker/example/cat
docker/example/cat:
	@docker exec -it e795688997e9 bash -c cat /tmp/state.json

.PHONY: docker/example/send
docker/example/send:
	@echo '["setItem", "foo", "bar"]' | nc localhost 32776

.PHONY: docker/build/example/bash
docker/build/example/bash:
	@$(MAKE) -C example/bash build

.PHONY: docker/run/example/bash
docker/run/example/bash:
	@$(MAKE) -C example/bash run

.PHONY: docker/deploy/example
docker/deploy/example: docker/build/example
	# build image -> push to ipfs -> pull from ipfs
	@go run main.go pull $$(go run main.go push $$(docker images -q | grep -m1 "") | grep "uploaded to" | sed -E 's/.*uploaded to \/ipfs\///g' | tr -d ' ' | tr -d '\n')
