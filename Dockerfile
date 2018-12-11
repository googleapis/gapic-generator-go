FROM debian:stable-slim

# Add protoc and our common protos.
COPY --from=gcr.io/gapic-images/api-common-protos:latest /usr/local/bin/protoc /usr/local/bin/protoc
COPY --from=gcr.io/gapic-images/api-common-protos:latest /protos/ /protos/

# Add protoc-gen-go_gapic binary
COPY protoc-gen-go_gapic /usr/local/bin

# Define the generator as an entry point.
ENTRYPOINT protoc --proto_path=/protos/ --proto_path=/in/ \
                  --go_gapic_out=/out/ \
                  --go_gapic_opt=$GO_GAPIC_OPT \
                  `find /in/ -name *.proto`
