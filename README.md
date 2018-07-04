# C3

> Implementation of the C3 protocol

## Install

### CLI

Install using `go get` (must have [Go](https://golang.org/doc/install) installed).

```bash
go get -u github.com/c3systems/c3
```

Show help for C3

```bash
$ c3 help
```

### Push image to IPFS

```bash
$ c3 push {imageID}
```

### Pull image from IPFS

```bash
$ c3 pull {ipfsHash}
```

## Test

```bash
make test
```

Tests require docker daemon and IPFS daemon to be running

## License

-
