#!/bin/bash

# Copyright 2018 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


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
	protoc --go_gapic_out "$OUT" -I "$COMMON_PROTO" -I "$GOOGLEAPIS" $*
}

generate --go_gapic_opt 'cloud.google.com/go/vision/apiv1;vision' $GOOGLEAPIS/google/cloud/vision/v1/*.proto
generate --go_gapic_opt 'cloud.google.com/go/speech/apiv1;speech' $GOOGLEAPIS/google/cloud/speech/v1/*.proto
# generate --go_gapic_opt 'cloud.google.com/go/language/apiv1;language' $GOOGLEAPIS/google/cloud/language/v1/*.proto
# generate --go_gapic_opt 'cloud.google.com/go/pubsub/apiv1;pubsub' $GOOGLEAPIS/google/pubsub/v1/*.proto
# generate --go_gapic_opt 'cloud.google.com/go/logging/apiv2;logging' $GOOGLEAPIS/google/logging/v2/*.proto
