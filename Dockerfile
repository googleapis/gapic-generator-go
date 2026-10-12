FROM debian:stable-slim@sha256:eb593cf2c358cacef45ca0a424bbc7d30cfa3466265fc2662b9466a0ca6ba1c5

# Add protoc and our common protos.
COPY --from=gcr.io/gapic-images/api-common-protos:latest@sha256:bff39e8eb3f95c117aaeb7fa36e7f118612a27bad041b2cb87627915cd7498cd /usr/local/bin/protoc /usr/local/bin/protoc
COPY --from=gcr.io/gapic-images/api-common-protos:latest@sha256:bff39e8eb3f95c117aaeb7fa36e7f118612a27bad041b2cb87627915cd7498cd /protos/ /protos/

# Add protoc-gen-go_gapic binary
COPY protoc-gen-go_gapic /usr/local/bin

# Add plugin invocation script for entrypoint
COPY docker-entrypoint.sh /usr/local/bin

# Set entry point script
ENTRYPOINT [ "docker-entrypoint.sh" ]
