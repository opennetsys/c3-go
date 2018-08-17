all: build

.PHONY: build
build: clean
	@go build -o config.so -buildmode=c-shared config.go

.PHONY: clean
clean:
	@-find . -type f -name *.so -delete
