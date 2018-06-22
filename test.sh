#!/bin/bash

# Usage: GOOGLEAPIS=path/to/googleapis ./testh.sh

set -e

if [ -z $GOOGLEAPIS ]; then
	echo >&2 "need GOOGLEAPIS variable"
	exit 1
fi

mkdir -p testdata/out

#vision
protoc --gogapic_out testdata/out -I $GOOGLEAPIS $GOOGLEAPIS/google/cloud/vision/v1/*.proto

#speech
#pubsub
#logging
