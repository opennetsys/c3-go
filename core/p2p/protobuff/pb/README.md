## Info
See:
1. [golang/protobuf](https://github.com/golang/protobuf)
2. [gogo/protobutf](https://github.com/gogo/protobuf)
3. [libp2p/go-libp2p/examples/multipro](https://github.com/libp2p/go-libp2p/tree/91fec896549430b7d93a82368b3bcd1ab71320a3/examples/multipro)

## Installation
First, [install protobuffs](https://github.com/golang/protobuf). Then:

```bash
$ go get github.com/gogo/protobuf/protoc-gen-gogofast
```

## Usage
```bash
$ protoc -I=. -I=$GOPATH/src -I=$GOPATH/src/github.com/gogo/protobuf/protobuf --gogofast_out=. p2p.proto
```
