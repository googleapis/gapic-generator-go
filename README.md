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

Disclaimer
----------
This generator is currently experimental. Please don't use it for anything mission-critical.

Go Version Supported
--------------------
The generator itself supports the latest version.

The generated code is compatible with Go 1.6.
