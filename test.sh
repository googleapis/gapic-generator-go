#!/bin/bash

# Usage: GOOGLEAPIS=path/to/googleapis [OUT=out/dir] ./test.sh
#
# If OUT is not set, files are written to testdata/out, which is gitignore'd.
# To integration test, set OUT=$GOPATH/src. The script will overwrite old files,
# and you can see changes by git-diff-ing the cloud.google.com/go repo.

set -e

if [ -z $GOOGLEAPIS ]; then
	echo >&2 "need GOOGLEAPIS variable"
	exit 1
fi

OUT=${OUT:-testdata/out}

mkdir -p "$OUT"

protoc --gogapic_out "$OUT" -I $GOOGLEAPIS --gogapic_opt cloud.google.com/go/vision/apiv1 $GOOGLEAPIS/google/cloud/vision/v1/*.proto
protoc --gogapic_out "$OUT" -I $GOOGLEAPIS --gogapic_opt cloud.google.com/go/speech/apiv1 $GOOGLEAPIS/google/cloud/speech/v1/*.proto
protoc --gogapic_out "$OUT" -I $GOOGLEAPIS --gogapic_opt cloud.google.com/go/pubsub/apiv1 $GOOGLEAPIS/google/pubsub/v1/*.proto
protoc --gogapic_out "$OUT" -I $GOOGLEAPIS --gogapic_opt cloud.google.com/go/logging/apiv2 $GOOGLEAPIS/google/logging/v2/*.proto
