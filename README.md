API Client Generator for Go
===========================

[![CircleCI](https://circleci.com/gh/googleapis/gapic-generator-go.svg?style=svg)](https://circleci.com/gh/googleapis/gapic-generator-go)

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

Disclaimer
----------
This generator is currently experimental. Please don't use it for anything mission-critical.

Go Version Supported
--------------------
The generator itself supports the latest version.

The generated code is compatible with Go 1.6.
