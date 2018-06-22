all:
	@echo "no default"

test/server:
	go test -v core/server/*.go

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
	docker run -dt snapshot_test:1

test/run/snapshot:
	node snapshot_test/index.js

