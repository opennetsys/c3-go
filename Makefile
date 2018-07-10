all: build

.PHONY: install
install:
	@go get -u github.com/c3systems/c3

.PHONY: deps
deps:
	@rm -rf ./vendor && \
	dep ensure && \
	gxundo ./vendor

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

.PHONY: test/before
test/before:
	# proxy localhost to 123.123.123.123 required so that docker container can communicate with host machine
	@sudo ifconfig lo0 alias 123.123.123.123/24

.PHONY: test/cleanup
test/cleanup:
	@chmod +x scripts/test_cleanup.sh
	@. scripts/test_cleanup.sh

.PHONY: test
test: test/c3 test/core/server test/core/docker test/core/registry test/core/sandbox test/ditto test/cleanup

.PHONY: test/c3
test/c3:
	@go test -v c3/*.go

.PHONY: test/core/server
test/core/server:
	@go test -v core/server/*.go

.PHONY: test/core/docker
test/core/docker:
	@go test -v core/docker/*.go $(ARGS)

.PHONY: test/core/sandbox
test/core/sandbox:
	@go test -v core/sandbox/*.go $(ARGS)

.PHONY: test/core/registry
test/core/registry:
	@go test -v core/registry/*.go

.PHONY: test/ditto
test/ditto:
	@docker pull hello-world && \
	go test -v ditto/*.go $(ARGS)

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
	@docker run -p 3333 --mount type=bind,src=/tmp,target=/tmp -t goexample

.PHONY: docker/example/send
docker/example/send:
	@echo '["setItem", "foo", "bar"]' | nc localhost 32776

.PHONY: docker/build/example/bash
docker/build/example/bash:
	@$(MAKE) -C example/bash build

.PHONY: docker/run/example/bash
docker/run/example/bash:
	@$(MAKE) -C example/bash run
