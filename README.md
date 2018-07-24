<p align="center">
	<img src="https://user-images.githubusercontent.com/168240/42129996-3bd8e646-7c8a-11e8-940d-89cea5ef87b8.png" width="250" alt="sandbox" />
	<br>
	<br>
</p>

# C3 Go

> Go implementation of the C3 protocol

[![License](http://img.shields.io/badge/license-Apache-blue.svg)](https://raw.githubusercontent.com/c3systems/c3/master/LICENSE.md) [![CircleCI](https://circleci.com/gh/c3systems/c3.svg?style=svg)](https://circleci.com/gh/c3systems/c3) [![Coverage Status](https://coveralls.io/repos/github/c3systems/c3/badge.svg?branch=master)](https://coveralls.io/github/c3systems/c3?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/c3systems/c3?)](https://goreportcard.com/report/github.com/c3systems/c3) [![GoDoc](https://godoc.org/github.com/c3systems/c3?status.svg)](https://godoc.org/github.com/c3systems/c3) [![Automated Release Notes by gren](https://img.shields.io/badge/%F0%9F%A4%96-release%20notes-00B2EE.svg)](https://github-tools.github.io/github-release-notes/) [![stability-experimental](https://img.shields.io/badge/stability-experimental-orange.svg)](https://github.com/emersion/stability-badges#experimental)

## Install

### Requirements

- Docker
  - [Install instructions](https://docs.docker.com/install/)
- IPFS
  - [Install instructions](https://ipfs.io/docs/install/)
- Patchutils
  - MacOS
    - `brew install patchutils`

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

## Docker config

Configure `daemon.json` to include the private registry as insecured (momentarily).

```json
{
  "insecure-registries" : [
    "{YOUR_LOCAL_IP}:5000"
  ]
}
```

- Linux
  - `/etc/docker/daemon.json`
- macOS
  - `~/.docker/daemon.json`

Restart the docker daemon after configuring `daemon.json`

## Test

```bash
make test
```

Tests require docker daemon and IPFS daemon to be running

## License

[Apache 2.0](LICENSE)
