all:
	@echo "no default"

test/core/server:
	go test -v core/server/*.go

test/core/dockerclient:
	go test -v core/dockerclient/*.go

test/core/registry:
	go test -v core/registry/*.go

test/ditto:
	go test -v ditto/*.go

run/example:
	go run example/go/main.go

docker/build/example:
	docker build --no-cache -f example/go/Dockerfile -t goexample ./example

docker/run/example:
	docker run -p 3333 -t goexample

test/docker/build/snapshot:
	docker build --no-cache -f snapshot_test/Dockerfile -t snapshot_test:1 ./snapshot_test

test/docker/run/snapshot:
	docker run -t snapshot_test:1

test/docker/run/snapshot/daemon:
	docker run -d snapshot_test:1

test/run/snapshot:
	node snapshot_test/index.js

docker/run/localregistry:
	docker run -d -p 5000:5000 --restart=always --name registry registry:2

docker/push/localregistry:
	docker push localhost:5000/$(IMAGE)

docker/list/localregistry:
	curl -X GET -k "http://$(docker-machine ip):5000/v2/_catalog"

docker/gcr/images/list:
	curl http://gcr.c3labs.io:5000/v2/_catalog
