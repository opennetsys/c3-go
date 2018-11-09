all: build

# INSTALL

.PHONY: install
install:
	@go get -u github.com/c3systems/c3-go

# /END INSTALL

# DEPS

.PHONY: deps
deps:
	@echo "running dep ensure..." && \
		dep ensure -update -v && \
		$(MAKE) gxundo && \
		#git clone https://github.com/gxed/pubsub.git vendor/github.com/gxed/pubsub && \
		#rm -rf vendor/github.com/gxed/pubsub/.git && \
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
test: test/check test/common test/trie test/cleanup
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

# CONFIG

.PHONY: test/config
test/config:
	@go test -v -parallel 1 config/*.go $(ARGS)

# /END CONFIG

# COMMON

.PHONY: test/common
test/common: test/common/netutil test/common/stringutil test/common/hexutil test/common/hashutil test/common/c3crypto

.PHONY: test/common/netutil
test/common/netutil:
	@go test -v common/netutil/*.go $(ARGS)

.PHONY: test/common/stringutil
test/common/stringutil:
	@go test -v common/stringutil/*.go $(ARGS)

.PHONY: test/common/hexutil
test/common/hexutil:
	@go test -v common/hexutil/*.go $(ARGS)

.PHONY: test/common/hashutil
test/common/hashutil:
	@go test -v common/hashutil/*.go $(ARGS)

.PHONY: test/common/command
test/common/command:
	@go test -v common/command/*.go $(ARGS)

.PHONY: test/common/c3crypto
test/common/c3crypto:
	@go test -v -parallel 1 common/c3crypto/*.go $(ARGS)

.PHONY: test/common/ipns
test/common/ipns:
	@go test -v -parallel 1 common/ipns/*.go $(ARGS)

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
test/core/docker: test/core/docker/util
	@go test -v -parallel 1 core/docker/*.go $(ARGS)

.PHONY: test/core/docker/util
test/core/docker/util:
	@go test -v -parallel 1 core/docker/util*.go $(ARGS)

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

.PHONY: test/core/eosclient
test/core/eosclient:
	@go test -v core/eosclient/*.go $(ARGS)

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

# TRIE

.PHONY: test/trie
test/trie:
	go test -v -parallel 1 trie/*.go $(ARGS)

# /END TRIE

# LOGGER

.PHONY: test/logger/color
test/logger/color:
	@go test -v logger/color/*.go $(ARGS)

# /END LOGGER

# NODE

.PHONY: run/node
run/node:
	@go run main.go node start --pem=node/test_data/priv1.pem --uri /ip4/0.0.0.0/tcp/3330 --data-dir .tmp --mempool-type memory --rpc ":5005"

.PHONY: run/node/2
run/node/2:
	@go run main.go node start --pem=node/test_data/priv2.pem --uri /ip4/0.0.0.0/tcp/3001 --data-dir ~/.c3-2 --peer "$(PEER)"

.PHONY: run/node/eos
run/node/eos:
	@go run main.go node start --pem=node/test_data/priv1.pem --uri /ip4/0.0.0.0/tcp/3330 --data-dir .tmp --mempool-type memory --rpc ":5005" \
		--checkpoint-eos-url "http://api.kylin.alohaeos.com" \
		--checkpoint-eos-account-name "helloworld54" \
		--checkpoint-eos-action-name "chkpointroot" \
		--checkpoint-eos-action-permissions "helloworld54@active" \
		--checkpoint-eos-wif-private-key "5Jh9tD4Fp1EpVn3EzEW6ura5NV3NddY8NNBcfpCZTvPDsKd9i5c"

.PHONY: run/node/ethereum
run/node/ethereum:
	@go run main.go node start --pem=node/test_data/priv1.pem --uri /ip4/0.0.0.0/tcp/3330 --data-dir .tmp --mempool-type memory --rpc ":5005" \
		--checkpoint-ethereum-url "https://rinkeby.infura.io" \
		--checkpoint-ethereum-contract-address "0x1e8cb4885e3aae139fea9aab55925f8e5ace2840" \
		--checkpoint-ethereum-method-name "checkpoint" \
		--checkpoint-ethereum-private-key ""

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

# proxy localhost to 123.123.123.123 required so that docker container can communicate with host machine
.PHONY: localhostproxy
localhostproxy:
	@sudo ifconfig $$(ifconfig | grep LOOPBACK | awk '{print $1}' | sed -E 's/[^a-zA-Z0-9]+//g') 123.123.123.123/24
	#@sudo ifconfig lo0 123.123.123.123/24

# /END MISC

# COVERAGE

.PHONY: coverage
coverage: coverage/install test/coverage

.PHONY: coverage/install
coverage/install:
	@go get golang.org/x/tools/cmd/cover
	@go get github.com/mattn/goveralls

.PHONY: test/coverage
test/coverage:
	@go test -v ./common/... -covermode=count -coverprofile=coverage.out
	@goveralls -coverprofile=coverage.out -service=travis-ci -repotoken="$$COVERALLS_TOKEN"

# /END COVERAGE

# NOTE: Temp fix till PR is merged:
# https://github.com/libp2p/go-libp2p-crypto/pull/35
.PHONY: fix/libp2pcrypto
fix/libp2pcrypto:
	@rm -rf vendor/github.com/libp2p/go-libp2p-crypto/
	@git clone -b forPR https://github.com/c3systems/go-libp2p-crypto.git
	@mv go-libp2p-crypto vendor/github.com/libp2p/go-libp2p-crypto
	@find "./vendor" -name "*.go" -print0 | xargs -0 perl -pi -e "s/c3systems\/go-libp2p-crypto/libp2p\/go-libp2p-crypto/g"
	@git clone git@github.com:gogo/protobuf.git
	@rm -rf vendor/github.com/gogo/protobuf
	@mv protobuf vendor/github.com/gogo/
	@git clone git@github.com:libp2p/go-libp2p-netutil.git
	@rm -rf vendor/github.com/libp2p/go-libp2p-netutil
	@mv go-libp2p-netutil vendor/github.com/libp2p/
	#@sed -iE 's/k1, k2 :=/k1, k2, _ :=/g' vendor/github.com/libp2p/go-libp2p-secio/protocol.go
	#@sed -iE 's/s.local.keys = k1/\/\/s.local.keys = k1/g' vendor/github.com/libp2p/go-libp2p-secio/protocol.go
	#@sed -iE 's/s.remote.keys = k2/\/\/s.remote.keys = k2/g' vendor/github.com/libp2p/go-libp2p-secio/protocol.go

.PHONY: fix/libp2ppubsub
fix/libp2pubsub:
	@rm vendor/github.com/libp2p/go-libp2p-pubsub/pb/rpc.pb.go
	@(cd vendor/github.com/libp2p/go-libp2p-pubsub/pb && wget https://raw.githubusercontent.com/libp2p/go-libp2p-pubsub/master/pb/rpc.pb.go)

# RPC

.PHONY: build/rpc/proto
build/rpc/proto:
	@protoc --go_out=plugins=grpc:. rpc/pb/c3.proto

# build protobufs for web client
# example:
# make build/rpc/proto/web OUT_DIR=./tmp
.PHONY: build/rpc/proto/web
build/rpc/proto/web:
	@protoc -I=rpc/pb/ --js_out=import_style=commonjs:$(OUT_DIR) --grpc-web_out=import_style=commonjs,mode=grpcwebtext:$(OUT_DIR) c3.proto

.PHONY: test/rpc
test/rpc:
	@go test -v rpc/*.go $(ARGS)

RPC_HOST := "localhost:5005"

.PHONY: run/rpc/ping
run/rpc/ping:
	@grpcurl -v -plaintext -d '{"jsonrpc":"2.0","id":"1","method":"c3_ping"}' $(RPC_HOST) protos.C3Service/Send

.PHONY: run/rpc/latestBlock
run/rpc/latestBlock:
	@grpcurl -v -plaintext -d '{"jsonrpc":"2.0","id":"1","method":"c3_latestBlock"}' $(RPC_HOST) protos.C3Service/Send

.PHONY: run/rpc/getblock
run/rpc/getblock:
	@grpcurl -v -plaintext -d '{"jsonrpc":"2.0","id":"1","method":"c3_getBlock","params":["0x3"]}' $(RPC_HOST) protos.C3Service/Send

.PHONY: run/rpc/getstateblock
run/rpc/getstateblock:
	@grpcurl -v -plaintext -d '{"jsonrpc":"2.0","id":"1","method":"c3_getStateBlock","params":["65cb6a153dd5", "0x1"]}' $(RPC_HOST) protos.C3Service/Send

.PHONY: run/rpc/pushImage
run/rpc/pushImage:
	@grpcurl -v -plaintext -d '{"jsonrpc":"2.0","id":"1","method":"c3_pushImage","params":[]}' $(RPC_HOST) protos.C3Service/Send

.PHONY: run/rpc/getstateblock
run/rpc/getstateblock:
	@grpcurl -v -plaintext -d '{"jsonrpc":"2.0","id":"1","method":"c3_getStateBlock","params":["65cb6a153dd5", "0x1"]}' $(RPC_HOST) protos.C3Service/Send

.PHONY: install/grpcwebproxy
install/grpcwebproxy:
	@go get github.com/improbable-eng/grpc-web/go/grpcwebproxy

.PHONY: run/grpcwebproxy
run/grpcwebproxy:
	@grpcwebproxy --backend_addr=localhost:5005 --run_tls_server=false

# /END RPC

# LOC

# get total lines of code
.PHONY: loc
loc:
	@find ./ -name '*.go' ! -path ".//vendor/*" ! -path ".//.git/*" | xargs wc -l

# /END LOC

# CLI
# Example
# $ make snapshot IMAGE=d50ada614c01 STATEBLOCK=2
.PHONY: snapshot
snapshot:
	@go run main.go snapshot --priv priv.pem --image $(IMAGE) --stateblock $(STATEBLOCK)

.PHONY: install/grpcwebproxy
install/grpcwebproxy:
	@go get github.com/improbable-eng/grpc-web/go/grpcwebproxy

.PHONY: run/grpcwebproxy
run/grpcwebproxy:
	@grpcwebproxy --backend_addr=localhost:5005 --run_tls_server=false

# /END RPC

# LOC

# get total lines of code
.PHONY: loc
loc:
	@find ./ -name '*.go' ! -path ".//vendor/*" ! -path ".//.git/*" | xargs wc -l

# /END LOC

# CLI
# Example
# $ make snapshot IMAGE=d50ada614c01 STATEBLOCK=2
.PHONY: snapshot
snapshot:
	@go run main.go snapshot --priv priv.pem --image $(IMAGE) --stateblock $(STATEBLOCK)
