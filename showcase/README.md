# Showcase integration testing

This submodule contains the gapic-showcase integration tests for
gapic-generator-go, as well as the script used to setup and execute them.

## How the tests work

The tests can be run out-of-the-box in the normal `go test` fashion, so long as
there is a `gapic-showcase` server already running locally. This does not test
the Showcase client generated with the locally installed generator - it would
use the released version of the GAPIC. Running `make test` (from the repository
root) will execute the tests against the locally installed generator.

To test the local generator, use the `showcase.bash` script. It does the
following:

1. Downloads the Showcase artifacts associated with the targeted version.
These are the compiled proto descriptor set, the retry configuration, and a
pre-compiled server binary.

1. Using protoc and the retrieved artifacts as input, the Go protobuf/gRPC
bindings, and the Go GAPIC are generated, the latter with the **locally
installed** gapic-generator-go. Make sure to `make install` (from the repository
root) so that any new changes are utilized during generation.

1. The submodule's `go.mod` is temporarily edited to replace the remote
dependency on `github.com/googleapis/gapic-showcase` with the locally generated
artifacts.

1. The server binary is started in the background.

1. The tests in this directory are executed against the locally running server.

1. The server process is stopped.

### Adding tests

Adding tests for an existing service is as easy as adding a new Go test to the
respective file: `echo_test.go`, `identity_test.go`, or `sequence_test.go`.

Adding a new service for testing requires creating an appropriate
`{service}_test.go` file and adding client creation to the `main_test.go`.

Any newly added tests should make sure to add a `leakcheck` statement to the
test function: `defer check(t)`.

### Updating the Showcase version

To update the version of Showcase referenced for artifact retrieval, update the
version in the `go.mod`. This version is extracted for use in artifact retrieval
and ensures all the necessary dependencies are up to date.

### Cleaning up after the tests

One can use `make clean` to revert any changes and delete any of the test
artifacts created by `showcase.bash`.
