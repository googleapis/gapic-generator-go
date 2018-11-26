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

# silence pushd & popd
pushd () {
    command pushd "$@" > /dev/null
}

popd () {
    command popd "$@" > /dev/null
}

set -e

if [ -z $GOOGLEAPIS ]; then
	echo >&2 "need GOOGLEAPIS variable"
	exit 1
fi
if [ -z $COMMON_PROTO ]; then
	echo >&2 "need COMMON_PROTO variable"
	exit 1
fi

OUT=${OUT:-$(pwd)/testdata}
KIOSK_PROTOS=${KIOSK_PROTOS:-$GOOGLEAPIS/kiosk/protos}
SHOW_PROTOS=${SHOW_PROTOS:-$GOOGLEAPIS/gapic-showcase/schema}
LANG_PROTOS=${LANG_PROTOS:-$GOOGLEAPIS/googleapis/google/cloud/language/v1}

mkdir -p "$OUT/kiosk"
mkdir -p "$OUT/showcase"
mkdir -p "$OUT/language"

generate() {
	protoc -I "$COMMON_PROTO" -I "$GOOGLEAPIS" $*
}

generate -I $KIOSK_PROTOS \
 --go_gcli_out $OUT/kiosk \
 --go_gcli_opt 'gapic:github.com/googleapis/kiosk/kioskgapic' \
 --go_gcli_opt 'root:testkctl' \
 $KIOSK_PROTOS/kiosk.proto

generate -I $SHOW_PROTOS \
  --go_gcli_out $OUT/showcase \
  --go_gcli_opt 'gapic:github.com/googleapis/gapic-showcase/showgapic' \
  --go_gcli_opt 'root:testshowctl' \
  $SHOW_PROTOS/*.proto

generate -I $LANG_PROTOS \
  --go_gcli_out $OUT/language \
  --go_gcli_opt 'gapic:github.com/googleapis/googleapis/google/cloud/language/v1/gapic' \
  --go_gcli_opt 'root:testlang' \
  $LANG_PROTOS/*.proto

d=$(pwd)
cd $OUT/kiosk
go build
cd ../showcase
go build
cd ../language
go build
cd $d

rm -r $OUT