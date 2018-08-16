all: build

.PHONY: build
build: clean
	@-find . -mindepth 2 -type f -name Makefile -execdir make \;

.PHONY: clean
clean:
	@-find . -type f -name *.so -delete
