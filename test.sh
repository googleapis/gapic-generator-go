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

generate() {
	protoc --go_gapic_out "$OUT" -I "$GOOGLEAPIS" $*
}

echo "Generating Cloud KMS v1"
generate --go_gapic_opt 'go-gapic-package=cloud.google.com/go/kms/apiv1;kms' $GOOGLEAPIS/google/cloud/kms/v1/*.proto

echo "Generating Cloud Data Catalog v1beta1"
generate --go_gapic_opt 'go-gapic-package=cloud.google.com/go/datacatalog/apiv1beta1;datacatalog' $GOOGLEAPIS/google/cloud/datacatalog/v1beta1/*.proto

echo "Generating Cloud Text-to-Speech v1 w/gRPC ServiceConfig"
generate --go_gapic_opt 'go-gapic-package=cloud.google.com/go/texttospeech/apiv1;texttospeech' --go_gapic_opt "grpc-service-config=$GOOGLEAPIS/google/cloud/texttospeech/v1/texttospeech_grpc_service_config.json" $GOOGLEAPIS/google/cloud/texttospeech/v1/*.proto

echo "Generation complete"

echo "Running gofmt to check for syntax errors"
gofmt -w -e $OUT

echo "No syntax errors"
