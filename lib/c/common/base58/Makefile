all: build

.PHONY: build
build: clean
	@go build -o base58.so -buildmode=c-shared base58.go

.PHONY: clean
clean:
	@-find . -type f -name *.so -delete
