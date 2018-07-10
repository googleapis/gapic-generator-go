#!/bin/bash

# Usage: COMMON_PROTO=path/to/api-common-protos GOOGLEAPIS=path/to/googleapis [OUT=out/dir] ./test.sh
#
# If OUT is not set, files are written to testdata/out, which is gitignore'd.
# To integration test, set OUT=$GOPATH/src. The script will overwrite old files,
# and you can see changes by git-diff-ing the cloud.google.com/go repo.
#
# We need proto annotations that's not stable yet. Use the input-contract branch of both repos.

set -e

if [ -z $GOOGLEAPIS ]; then
	echo >&2 "need GOOGLEAPIS variable"
	exit 1
fi
if [ -z $COMMON_PROTO ]; then
	echo >&2 "need COMMON_PROTO variable"
	exit 1
fi

OUT=${OUT:-testdata/out}

mkdir -p "$OUT"

generate() {
	protoc --gogapic_out "$OUT" -I "$COMMON_PROTO" -I "$GOOGLEAPIS" $*
}

generate --gogapic_opt 'cloud.google.com/go/vision/apiv1;vision' $GOOGLEAPIS/google/cloud/vision/v1/*.proto
generate --gogapic_opt 'cloud.google.com/go/speech/apiv1;speech' $GOOGLEAPIS/google/cloud/speech/v1/*.proto
generate --gogapic_opt 'cloud.google.com/go/pubsub/apiv1;pubsub' $GOOGLEAPIS/google/pubsub/v1/*.proto
generate --gogapic_opt 'cloud.google.com/go/logging/apiv2;logging' $GOOGLEAPIS/google/logging/v2/*.proto
