# Showcase integration testing

This submodule contains the gapic-showcase integration tests for
gapic-generator-go, as well as the script used to setup and execute them.

## How the tests work

The tests can be run out-of-the-box in the normal `go test` fashion, so long as
there is a `gapic-showcase` server already running locally. This does not test
the a client generated with the locally checked out generator - it would use the
released version of the GAPIC. Running `make test` will execute the tests
against the locally install generator.

To test the local generator, the `showcase.bash` script first downloads the
Showcase artifacts associated with the targeted version. These are the compiled
proto descriptor set, the retry configuration, and a pre-compiled server binary.

Then using protoc and the retrieved artifacts as input, the Go protobuf/gRPC
bindings, and the Go GAPIC are generated, the latter with the
**locally installed** gapic-generator-go. Make sure to `make install` so that
any new changes are utilized during generation.

This submodule's `go.mod` is temporarily edited to replace the remote dependency
on `gapic-showcase` with the locally generated artifacts.

The server binary is run in the background and the tests in this directory are
executed against it.

### Adding tests

Adding tests for an existing service is as easy as adding a new Go test to the
respective file: `echo_test.go`, `identity_test.go`, or `sequence_test.go`.

Adding a new service for testiing requires creating an appropriate
`{service}_test.go` file and adding client creation to the `main_test.go`.

Any newly added tests should make sure to add a `leakcheck` statement to the
test function: `defer check(t)`.

### Updating the Showcase version

To update the version of Showcase referenced for artifact retrieval, change the
`VERSION` variable in `showcase.bash`. It is also good practice to update the
`go.mod` to this version.

### Cleaning up after the tests

One can use `make clean` to revert any changes and delete any of the test
artifacts created by `showcase.bash`.
