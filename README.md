API Client Generator for Go
===========================

![latest release](https://img.shields.io/github/v/release/googleapis/gapic-generator-go)
![Go version](https://img.shields.io/github/go-mod/go-version/googleapis/gapic-generator-go)

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
`go install github.com/googleapis/gapic-generator-go/cmd/protoc-gen-go_gapic@latest`.
If you are using Go 1.11 and see error `cannot find main module`, see this [FAQ page](https://github.com/golang/go/wiki/Modules#why-does-installing-a-tool-via-go-get-fail-with-error-cannot-find-main-module).

Or to install from source:
```
git pull https://github.com/googleapis/gapic-generator-go.git
cd gapic-generator-go
go install ./cmd/protoc-gen-go_gapic
```

The generator works as a `protoc` plugin, get `protoc` from [google/protobuf](https://github.com/protocolbuffers/protobuf).

Configuration
-------------
The generator is configured via protobuf annotations found at [googleapis/api-common-protos](https://github.com/googleapis/api-common-protos).
The generator follows the guidance defined in [AIP-4210](https://aip.dev/4210).

The only *required* annotation to generate a client is the service annotation `google.api.default_host` ([here](https://github.com/googleapis/googleapis/blob/82528cf321ed0d09b7d93b7cee9122ccea422ad2/google/api/client.proto#L67-L76)).

The value of `google.api.default_host` must be just a host name, excluding a scheme. For example,
```
import "google/api/client.proto";
...

service Foo {
    option (google.api.default_host) = "api.foo.com";
    ...
}  
```

If a RPC returns a `google.longrunning.Operation`, the RPC must be annotated with `google.longrunning.operation_info` in accordance with [AIP-151](https://aip.dev/151).

The [Cloud Natural Language API](https://github.com/googleapis/googleapis/blob/master/google/cloud/language/v1/language_service.proto) is an example of a fully configured API that has a [Go client](https://github.com/googleapis/google-cloud-go/tree/main/language/apiv1) generated by gapic-generator-go.

The supported configuration annotations include:
* Service Options
  * `google.api.default_host`: host name used in the default service client initialization
  * `google.api.oauth_scopes`: OAuth scopes needed by the client to auth'n/z
* Method Options
  * `google.longrunning.operation_info`: used to determine response & metadata types of LRO methods

Invocation
----------
`protoc -I $GOOGLEAPIS --go_gapic_out [OUTPUT_DIR] --go_gapic_opt 'go-gapic-package=package/path/url;name' a.proto b.proto`

**Note:** The `$GOOGLEAPIS` variable represents a path to the [googleapis/googleapis](https://github.com/googleapis/googleapis) directory to import the configuration annotations.

The `go_gapic_opt` protoc plugin option flag is necessary to convey configuration information not present in the protos. 
The plugin option's value is a key-value pair delimited by an equal sign `=`.
The configuration supported by the plugin option includes:
  
  * `go-gapic-package`: the Go package of the generated client library.
    *  The substring preceding the semicolon is the import path of the package, e.g. `github.com/username/awesomeness`.
    *  The substring after the semicolon is the name of the package used in the `package` statement.
    
    **Note:** Idiomatically the name is last element of the path but it need not be.
    For instance, the last element of the path might be the package's version, and the package would benefit
    from a more descriptive name.

  * `module`: prefix to be stripped from the `go-gapic-package` used in the generated filenames.
     * _Note: This option is not supported from the Bazel interface._

  * `metadata`: enable generation of [GapicMetadata](https://github.com/googleapis/googleapis/blob/master/gapic/metadata/gapic_metadata.proto) in JSON form. The default is `false`.
  
  * `grpc-service-config`: the path to a gRPC ServiceConfig JSON file.
    * This is used for client-side retry configuration in accordance with [AIP-4221](http://aip.dev/4221)

  * `release-level`: the client library release level.
    * Defaults to empty, which is essentially the GA release level.
    * Acceptable values are `alpha` and `beta`.

  * `api-service-config`: the path the service YAML file.
    * This is used for service-level client documentation.

  * `transport`: the desired transport(s) to generate, delimited by `+` e.g. `grpc+rest`.
    * Acceptable values are `grpc` and `rest`.
    * Defaults to `grpc`.

  * `rest-numeric-enums`: enables requesting response enums be encoded as numbers.
    * Not enabled by default.
    * Only effective when `rest` is included as a `transport` to be generated.

  * `omit-snippets`: disable generation of code snippets to the `internal/generated/snippets` path. The default is `false`.

Bazel
-----

The generator can be executed via a Bazel BUILD file using the macro in this repo.

Add the following to your WORKSPACE to import this project.

```
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "com_googleapis_gapic_generator_go",
    strip_prefix = "gapic-generator-go-main",
    urls = ["https://github.com/googleapis/gapic-generator-go/archive/main.zip"],
)

load("@com_googleapis_gapic_generator_go//:repositories.bzl", "com_googleapis_gapic_generator_go_repositories")

com_googleapis_gapic_generator_go_repositories()
```

_Note: do not use `main`, use a commit hash or a release tag._

And invoke it in a BUILD file like so, using an example based on the [googleapis repo](https://github.com/googleapis/googleapis/tree/92bebf78345af8b2d3585220527115bda8bdedf8/google/cloud/language/v1).

```
load("@com_googleapis_gapic_generator_go//rules_go_gapic:go_gapic.bzl", "go_gapic_library")

go_gapic_library(
  name = "language_go_gapic",
  srcs = [
    # BUILD target for proto_library
    "//google/cloud/language/v1:language_proto",
  ],
  deps = [
    # BUILD target for go_library_proto
    "//google/cloud/language/v1:language_go_proto",
  ],
  # go-gapic-package parameter value
  importpath = "cloud.google.com/go/language/apiv1;language",
)
```

The generator options defined in [Invocation](#Invocation) are supported as the
following attributes:

  * `grpc_service_config`: a label for a gRPC ServiceConfig JSON file.

  * `release_level`: the client library release level.

  * `service_yaml`: a label for a service YAML file.
    * _Note: This option will eventually be deprecated._

  * `metadata`: if `True`, [GapicMetadata](https://github.com/googleapis/googleapis/blob/master/gapic/metadata/gapic_metadata.proto) will be generated in JSON form. The default is `False`.

  * `transport`: the desired transport(s) to generate, delimited by `+` e.g. `grpc+rest`.
    * Acceptable values are `grpc` and `rest`.
    * Defaults to `grpc`.

  * `rest_numeric_enums`: if `True`, enables generation of system parameter requesting
  response enums be encoded as numbers.
    * Default is `False`.
    * Only effective when `rest` is included as a `transport` to be generated.

  * `omit_snippets`: if `True`, code snippets will be generated to the `internal/generated/snippets` path. The default is `True`.

Docker Wrapper
--------------
The generator can also be executed via a Docker container. The image containes `protoc`, the microgenerator
binary, and the standard API protos.

```bash
$ docker run \
  --rm \
  --user $UID \
  --mount type=bind,source=</abs/path/to/protos>,destination=/in,readonly \
  --mount type=bind,source=</abs/path/to/configs>,destination=/conf,readonly \
  --mount type=bind,source=$GOPATH/src,destination=/out/ \
  gcr.io/gapic-images/gapic-generator-go \
  --go-gapic-package "github.com/package/import/path;name"
```

Replace `/abs/path/to/protos` with the absolute path to the input protos and `github.com/package/import/path;name`
with the desired import path & name for the `gapic`, as described in [Invocation](#Invocation).

For convenience, the [gapic.sh](./gapic.sh) script wraps the above `docker` invocation.
An equivalent invocation using `gapic.sh` is:

```bash
$ gapic.sh \
  --image gcr.io/gapic-images/gapic-generator-go \
  --in /abs/path/to/protos \
  --out $GOPATH/src \
  --go-gapic-package "<github.com/package/import/path;name>"
```

Use `gapic.sh --help` to print the usage documentation.

Code Generation
---------------

This is an explanation of the Go GAPIC generator for those interested in how it works and possibly those using it as a reference.

### Plugin interface

`gapic-generator-go` is a `protoc` [plugin](https://developers.google.com/protocol-buffers/docs/reference/other). It consumes a serialzed `CodeGeneratorRequest` on `stdin` and produces a serialized `CodeGeneratorResponse` on `stdout`. The `CodeGeneratorResponse` contains all of the generated Go code and/or any error(s) that might of occured during generation. All logs are emitted on `stderr`.

The plugin implementation can be found in [cmd/protoc-gen-go_gapic](/cmd/protoc-gen-go_gapic).

### Generated Artifacts

A single invocation of the code generator creates a `doc.go` file package level documentation according to [godoc](https://blog.golang.org/godoc-documenting-go-code).  This documentation is (currently) pulled from a given service config.

Each service found in the input protos gets two generated artifacts:

* `{service}_client.go`: contains the GAPIC implementation
* `{service}_client_example_test.go`: contains example code for each service method, consumed by [godoc](https://blog.golang.org/examples)

There is no directory structure in the generated output. All files are placed directly in the designated output directory by `protoc`.

### Generation Process

The generator implementation can be found in [internal/gengapic](/internal/gengapic).

The service client type, initialization code and any standard helpers are generated first. Then each method is generated. Any relevant helper types (i.e. pagination [Iterator](https://github.com/googleapis/google-cloud-go/wiki/Iterator-Guidelines) types, LRO helpers, etc.) for the service methods are generated following the methods.

Following the client implementation, the client example file is generated, and after all services have been generated the single `doc.go` file is created.

Go Version Supported
--------------------
The generator itself supports the latest version.

The generated code is compatible with Go 1.6.

Contributing
------------

If you are looking to contribute to the project, please see CONTRIBUTING.md for guidelines.