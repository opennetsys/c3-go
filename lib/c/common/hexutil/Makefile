all: build

.PHONY: build
build: clean
	@go build -o hexutil.so -buildmode=c-shared hexutil.go

.PHONY: clean
clean:
	@-find . -type f -name *.so -delete
