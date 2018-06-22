all:
	@echo "no default"

test/server:
	go test -v core/server/*.go

run/example:
	go run example/go/main.go

docker/build/example:
	docker build --no-cache -f example/go/Dockerfile -t goexample .

docker/run/example:
	docker run -p 3333 -t goexample
