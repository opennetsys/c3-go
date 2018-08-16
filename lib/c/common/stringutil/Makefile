all: build

.PHONY: build
build: clean
	@go build -o stringutil.so -buildmode=c-shared stringutil.go

.PHONY: clean
clean:
	@-find . -type f -name *.so -delete
