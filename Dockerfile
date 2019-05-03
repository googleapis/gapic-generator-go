FROM debian:stable-slim

# Add protoc and our common protos.
COPY --from=gcr.io/gapic-images/api-common-protos:0.1.0 /usr/local/bin/protoc /usr/local/bin/protoc
COPY --from=gcr.io/gapic-images/api-common-protos:0.1.0 /protos/ /protos/

# Add gapic-config-validator plugin
COPY --from=gcr.io/gapic-images/gapic-config-validator /usr/local/bin/protoc-gen-gapic-validator /usr/local/bin/protoc-gen-gapic-validator

# Add protoc-gen-go_gapic binary
COPY protoc-gen-go_gapic /usr/local/bin

# Add plugin invocation script for entrypoint
COPY docker-entrypoint.sh /usr/local/bin

# Set entry point script
ENTRYPOINT [ "docker-entrypoint.sh" ]
