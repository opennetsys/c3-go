<p align="center">
	<img src="https://user-images.githubusercontent.com/168240/44509676-396bab80-a667-11e8-89d1-69e8599af9af.png" width="250" alt="sandbox" /> <br>
	<br>
</p>

# C3 Go

> Go implementation of the C3 protocol

[![License](http://img.shields.io/badge/license-GNU%20AGPL%203.0-blue.svg)](https://raw.githubusercontent.com/c3systems/c3/master/LICENSE.md) [![CircleCI](https://circleci.com/gh/c3systems/c3-go.svg?style=svg)](https://circleci.com/gh/c3systems/c3-go) [![Coverage Status](https://coveralls.io/repos/github/c3systems/c3-go/badge.svg?branch=master)](https://coveralls.io/github/c3systems/c3-go?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/c3systems/c3-go?)](https://goreportcard.com/report/github.com/c3systems/c3-go) [![GoDoc](https://godoc.org/github.com/c3systems/c3-go?status.svg)](https://godoc.org/github.com/c3systems/c3-go) [![Automated Release Notes by gren](https://img.shields.io/badge/%F0%9F%A4%96-release%20notes-00B2EE.svg)](https://github-tools.github.io/github-release-notes/) [![stability-experimental](https://img.shields.io/badge/stability-experimental-orange.svg)](https://github.com/emersion/stability-badges#experimental)

## Install

### Requirements

- Docker
  - [Install instructions](https://docs.docker.com/install/)
- IPFS
  - [Install instructions](https://ipfs.io/docs/install/)
- Patchutils
  - MacOS
    - `brew install patchutils`
  - Linux
    - `[sudo] apt-get install patchutils`

#### Docker config

Set up a localhost proxy route to `123.123.123.123`

```bash
$ sudo ifconfig $(ifconfig | grep LOOPBACK | awk '{print $1}' | sed -E 's/[^a-zA-Z0-9]+//g') 123.123.123.123/24
```

Then configure `daemon.json` to include the private registry as insecured (because it's localhost).

```json
{
  "insecure-registries" : [
    "123.123.123.123:5000"
  ]
}
```

- Linux
  - `/etc/docker/daemon.json`
- macOS
  - `~/.docker/daemon.json`

Restart the docker daemon after configuring `daemon.json`

### Install c3-go

Install using `go get` (must have [Go](https://golang.org/doc/install) installed).

```bash
$ go get -u github.com/c3systems/c3
```

## Hello world tutorial

Warning: Tutorial is not the most intuitive. Please bear with me as we're still in the proof-of-concept stage.

1. Generate new private keys

```bash
$ c3-go generate key -o priv.pem
$ c3-go generate key -o priv2.pem
```

2. Run the C3 node

```bash
$ c3-go node start --pem=priv.pem --uri /ip4/0.0.0.0/tcp/3330 --data-dir ~/.c3

INFO[0002] [node] 0: /ip4/127.0.0.1/tcp/3330/ipfs/QmNRR7uLZ2bZXjjQqEY5fcm5BXubBEne3bkq6pYwg1QR18
  source="start.go:152:node.Start"
```

3. In another terminal, clone and build the hello-world dapp

```bash
$ git clone https://github.com/c3systems/c3-sdk-go-example-hello-world.git hello-world

$ (cd hello-world && docker build .)
```

4. Push the image to IPFS

```bash
$ c3-go push $(docker images -q | grep -m1 "")

[registry] uploaded to /ipfs/QmWJF5MYtnjb76P1CXQsn8MHpT26tjdBcs6CzKfR7zjRBm
  source="registry.go:101:registry.(*Registry).PushImage"
```

5. Deploy the image to the local C3 testnet (update `image` and `peer` to your own values)

```bash
$ c3-go deploy --priv priv2.pem --genesis '' --image QmWJF5MYtnjb76P1CXQsn8MHpT26tjdBcs6CzKfR7zjRBm --peer "/ip4/127.0.0.1/tcp/3330/ipfs/QmZQ3cJMMjA7HUyEvsMXmN73LZ7fKsrQUmyKwsxrpecb7Z"
```

- The private key is the second key derived from step 1.
- The peer multihash was derived from step 2.
- The image hash was derived from step 4.

6. Invoke a method on the dApp (update `image` and `peer` to your own values from the previous step)

```bash
$ go run main.go invokeMethod --priv priv2.pem --payload '["setItem", "foo", "bar"]' --image QmWJF5MYtnjb76P1CXQsn8MHpT26tjdBcs6CzKfR7zjRBm --peer "/ip4/127.0.0.1/tcp/3330/ipfs/QmZQ3cJMMjA7HUyEvsMXmN73LZ7fKsrQUmyKwsxrpecb7Z"
```

- In this example we're invoking the `setItem` method which accepts two arguments; the values are `foo` and `bar`. The example dApp code is found [here](https://github.com/c3systems/c3-sdk-go-example-hello-world/blob/master/main.go).

- The private key is the second key derived from step 1.
- The peer multihash was derived from step 2.
- The image hash was derived from step 4.

7. That's it; you've just deployed a dApp and invoked a method by broadcasting a transaction.

## CLI commands

Show help for C3

```bash
$ c3-go help
```

#### Push image to IPFS

```bash
$ c3-go push {imageID}
```

#### Pull image from IPFS

```bash
$ c3-go pull {ipfsHash}
```

#### Run a node

```bash
$ c3-go node start [options]
```

#### Generate a private key

```bash
$ c3-go generate key
```

#### Encode data

```go
$ c3-go encode [options]
```

#### Deploy dApp

```go
$ c3-go deploy [options]
```

#### Invoke dApp method

```go
$ c3-go invokeMethod [options]
```

## Test

```bash
$ make test
```

Tests require docker daemon and IPFS daemon to be running.

## Development

Install new dependencies to `vendor/` (will take a couple of minutes)

```bash
$ make deps
```

## License

[GNU AGPL 3.0](LICENSE)
