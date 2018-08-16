all: build

.PHONY: build
build: clean
	@go build -o hashing.so -buildmode=c-shared hashing.go

.PHONY: clean
clean:
	@-find . -type f -name *.so -delete
