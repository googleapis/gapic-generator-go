API Client Generator for Go
===========================

[![CircleCI](https://circleci.com/gh/googleapis/gapic-generator-go.svg?style=svg)](https://circleci.com/gh/googleapis/gapic-generator-go) 
![release level](https://img.shields.io/badge/release%20level-%20beta-blue.svg)

A generator for protocol buffer described APIs for and in Go.

This is a generator for API client libraries for APIs specified by protocol buffers, such as those inside Google.
It takes a protocol buffer (with particular annotations) and uses it to generate a client library.

Purpose
-------
We aim for this generator to replace the [older monolithic generator](https://github.com/googleapis/gapic-generator).
Some areas we hope to improve over the old generator are:
- using explicit normalized format for specifying APIs,
- simpler, faster implementation, and
- better error reporting.

Installation
------------
`go get github.com/googleapis/gapic-generator-go/cmd/protoc-gen-go_gapic`.
If you are using Go 1.11 and see error `cannot find main module`, see this [FAQ page](https://github.com/golang/go/wiki/Modules#why-does-installing-a-tool-via-go-get-fail-with-error-cannot-find-main-module).

Or to install from source:
```
git pull https://github.com/googleapis/gapic-generator-go.git
cd gapic-generator-go
go install ./cmd/protoc-gen-go_gapic
```

The generator works as a `protoc` plugin, get `protoc` from [google/protobuf](https://github.com/protocolbuffers/protobuf).

Invocation
----------
`protoc --go_gapic_out [OUTPUT_DIR] --go_gapic_opt 'package/path/url;name' a.proto b.proto`

The `go_gapic_opt` flag is necessary because we need to know where to generated file will live.
The substring before the semicolon is the import path of the package, e.g. `github.com/username/awesomeness`.
The substring after the semicolon is the name of the package used in the `package` statement.
Idiomatically the name is last element of the path but it need not be.
For instance, the last element of the path might be the package's version, and the package would benefit
from a more descriptive name.

Docker Wrapper
--------------
The generator can also be executed via a Docker container. The image containes `protoc`, the microgenerator
binary, and the standard API protos.

```bash
$ docker pull gcr.io/gapic-images/gapic-generator-go
$ docker run \
  --rm \
  --user $UID \
  --mount type=bind,source=</abs/path/to/protos>,destination=/in,readonly \
  --mount type=bind,source=$GOPATH/src,destination=/out/ \
  --env "GO_GAPIC_OPT=<github.com/package/import/path;name>" \
  gcr.io/gapic-images/gapic-generator-go
```

Replace `/abs/path/to/protos` with the absolute path to the input protos and `github.com/package/import/path;name`
with the desired import path & name for the `gapic`, as described in [Invocation](#Invocation).

For convenience, the [gapic.sh](./gapic.sh) script wraps the above `docker` invocation.
An equivalent invocation using `gapic.sh` is:

```bash
$ gapic.sh \
  --go_gapic_opt 'github.com/package/import/path;name' \
  --image gcr.io/gapic-images/gapic-generator-go
  --in /abs/path/to/protos \
  --out $GOPATH/src
```

Use `gapic.sh --help` to print the usage documentation.

Disclaimer
----------
This generator is currently experimental. Please don't use it for anything mission-critical.

Go Version Supported
--------------------
The generator itself supports the latest version.

The generated code is compatible with Go 1.6.
