FROM debian:stable-slim@sha256:1710bde34461551a19a47c787885ec9ad7058d9a5bead2affb8d088fa2f8502b

# Add protoc and our common protos.
COPY --from=gcr.io/gapic-images/api-common-protos:latest@sha256:bff39e8eb3f95c117aaeb7fa36e7f118612a27bad041b2cb87627915cd7498cd /usr/local/bin/protoc /usr/local/bin/protoc
COPY --from=gcr.io/gapic-images/api-common-protos:latest@sha256:bff39e8eb3f95c117aaeb7fa36e7f118612a27bad041b2cb87627915cd7498cd /protos/ /protos/

# Add protoc-gen-go_gapic binary
COPY protoc-gen-go_gapic /usr/local/bin

# Add plugin invocation script for entrypoint
COPY docker-entrypoint.sh /usr/local/bin

# Set entry point script
ENTRYPOINT [ "docker-entrypoint.sh" ]
