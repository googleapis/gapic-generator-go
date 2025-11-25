# GAPIC Generator Go (`gapic-generator-go`) Context

## Repository Overview
This is the **code generator**. It reads Protocol Buffers and generates the Go client code found in `google-cloud-go`.

## Key Directories
*   `cmd/protoc-gen-go_gapic/`: The entry point for the `protoc` plugin. Start here to trace execution.
*   `internal/gengapic/`: This directory contains the **core logic** for generating the Go code.
    *   If you need to change *how* code is generated (e.g., changing the format of a client method), this is where you look.
*   `internal/snippets/`: Code for generating documentation snippets.

## Architecture & Wiring (Code Generation)
*   **Client Constructors (`internal/gengapic/client_init.go`):** Generates the `NewClient` and `NewRESTClient` functions.
    *   **Edit here:** If you need to change how clients are initialized (e.g., adding default options or transport configuration).
*   **Method Generation:**
    *   `internal/gengapic/gengrpc.go`: Generates gRPC methods.
    *   `internal/gengapic/genrest.go`: Generates REST (HTTP) methods.
    *   **Edit here:** To inject logic at the call site (e.g., adding attributes to the context before the RPC).

## Workflow
*   To test changes, you typically need to:
    1.  Build this generator.
    2.  Run it against a sample proto (often found in `google-cloud-go` or a test workspace).
    3.  Inspect the generated output.
