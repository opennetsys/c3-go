all: build

# INSTALL

.PHONY: install
install:
	@go get -u github.com/c3systems/c3-go

# /END INSTALL

# DEPS

.PHONY: deps
deps:
	@rm -rf ./vendor && \
		echo "running dep ensure..." && \
		dep ensure && \
		$(MAKE) gxundo && \
		(cd vendor/github.com/libp2p/go-libp2p-pubsub/pb \
		&& rm rpc.pb.go && rm rpc.proto \
		&& wget https://github.com/c3systems/go-libp2p-pubsub/raw/master/pb/rpc.pb.go \
		&& wget https://github.com/c3systems/go-libp2p-pubsub/raw/master/pb/rpc.proto) && \
		git clone https://github.com/gxed/pubsub.git vendor/github.com/gxed/pubsub && \
		rm -rf vendor/github.com/gxed/pubsub/.git && \
		$(MAKE) deps/copy/ethereum/crypto

.PHONY: gxundo
gxundo:
	@bash scripts/gxundo.sh vendor/

.PHONY: install/gxundo
install/gxundo:
	@wget https://raw.githubusercontent.com/c3systems/gxundo/master/gxundo.sh \
	-O scripts/gxundo.sh && \
	chmod +x scripts/gxundo.sh

.PHONY: deps/copy/ethereum/crypto
deps/copy/ethereum/crypto:
	@mkdir -p vendor/github.com/ethereum/go-ethereum/crypto/secp256k1 && \
		go get github.com/ethereum/go-ethereum/crypto/secp256k1/... && \
		cp -r "${GOPATH}/src/github.com/ethereum/go-ethereum/crypto/secp256k1/libsecp256k1" "vendor/github.com/ethereum/go-ethereum/crypto/secp256k1/"

# /END DEPS

# BUILD

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

# /END BUILD

# TEST ALL

.PHONY: test/check
test/check:
	# TODO: kill script if required commands and daemons not found
	@command -v ipfs &>/dev/null || echo "IPFS is required"
	@command -v docker &>/dev/null || echo "Docker daemon is required"
	@pgrep -f ipfs > /dev/null || echo "IPFS daemon is not running"
	@pgrep -f docker > /dev/null || echo "Docker daemon is not running"

.PHONY: test
test: test/check test/common test/cleanup
	# test/unit test/integration test/e2e
	# test/core
	# test/registry

.PHONY: test/cleanup
test/cleanup:
	@chmod +x scripts/test_cleanup.sh
	@. scripts/test_cleanup.sh

# /END TEST ALL

# TEST TYPES

.PHONY: test/unit
test/unit:
	@go test ./... -tags=unit

.PHONY: test/integration
test/integration:
	@go test ./... -tags=integration

.PHONY: test/e2e
test/e2e:
	@go test ./... -tags=e2e

# /END TEST TYPES

# COMMON

.PHONY: test/common
test/common: test/common/netutil test/common/stringutil test/common/hexutil test/common/hashing test/common/c3crypto

.PHONY: test/common/netutil
test/common/netutil:
	@go test -v common/netutil/*.go $(ARGS)

.PHONY: test/common/stringutil
test/common/stringutil:
	@go test -v common/stringutil/*.go $(ARGS)

.PHONY: test/common/hexutil
test/common/hexutil:
	@go test -v common/hexutil/*.go $(ARGS)

.PHONY: test/common/hashing
test/common/hashing:
	@go test -v common/hashing/*.go $(ARGS)

.PHONY: test/common/command
test/common/command:
	@go test -v common/command/*.go $(ARGS)

.PHONY: test/common/c3crypto
test/common/c3crypto:
	@go test -v -parallel 1 common/c3crypto/*.go $(ARGS)

# /END COMMON

# CORE

.PHONY: test/core
test/core: test/core/server test/core/ipfs
	# test/core/sandbox
	# test/core/docker
	# test/core/chain/mainchain/miner
	# test/core/diffing

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
test/core/sandbox:
	@go test -v -parallel 1 core/sandbox/*.go $(ARGS)

.PHONY: test/core/sandbox/with-build
test/core/sandbox/with-build: docker/build/example
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

.PHONY: test/core/chain/mainchain/miner
test/core/chain/mainchain/miner:
	@go test -v core/chain/mainchain/miner/*.go $(ARGS)

.PHONY: test/core/diffing
test/core/diffing:
	@go test -v core/diffing/*.go $(ARGS)

# /END CORE


# REGISTRY

.PHONY: test/registry
test/registry:
	@docker pull hello-world && \
	go test -v -parallel 1 registry/*.go $(ARGS)

.PHONY: test/registry/server
test/registry/server:
	@go test -v registry/server/*.go $(ARGS)

# /END REGISTRY

.PHONY: test/node
test/node:
	@IMAGEID="$(IMAGEID)" PEERID="$(PEERID)" METHOD="$(METHOD)" go test -v node/*.go $(ARGS)

# /END REGISTRY

# LOGGER

.PHONY: test/logger/color
test/logger/color:
	@go test -v logger/color/*.go $(ARGS)

# /END LOGGER

# NODE

.PHONY: run/node
run/node:
	@go run main.go node start --pem=node/test_data/priv1.pem --uri /ip4/0.0.0.0/tcp/9005 --data-dir ~/.c3-1 --difficulty 5

.PHONY: run/node/2
run/node/2:
	@go run main.go node start --pem=node/test_data/priv2.pem --uri /ip4/0.0.0.0/tcp/9006 --data-dir ~/.c3-2 --difficulty 5

.PHONY: node/save/testimage
node/save/testimage:
	@docker save goexample -o node/test_data/go_example_image.tar

# /END NODE

# DOCKER

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
	@docker run -p 5555:3333 --mount type=bind,src=/tmp,target=/tmp -t goexample

.PHONY: docker/example/cat
docker/example/cat:
	@docker exec -it e795688997e9 bash -c cat /tmp/state.json

.PHONY: docker/example/send
docker/example/send:
	@echo '["0x1e51aea686ccbea473b94a662b980644601831cf1a390a4fb08b1793bc6c6463","0x666f6f","0x626172"]' | nc localhost 3333

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

# /END DOCKER

# IPFS

.PHONY: ipfs/daemon
ipfs/daemon:
	@ipfs daemon

# /END IPFS

# MISC

.PHONY: localhostproxy
localhostproxy:
	# proxy localhost to 123.123.123.123 required so that docker container can communicate with host machine
	@sudo ifconfig lo0 alias 123.123.123.123/24

# /END MISC
