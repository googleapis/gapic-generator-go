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


# Usage: COMMON_PROTO=path/to/api-common-protos [OUT=out/dir] ./test.sh
#
# If OUT is not set, files are written to testdata/out, which is gitignore'd.
# To integration test, set OUT=$GOPATH/src. The script will overwrite old files,
# and you can see changes by git-diff-ing the cloud.google.com/go repo.
#
# We need proto annotations that's not stable yet. Use the input-contract branch of both repos.

set -e

if [ -z $COMMON_PROTO ]; then
	echo >&2 "need COMMON_PROTO variable"
	exit 1
fi

OUT=${OUT:-$(pwd)/testdata}
KIOSK_PROTOS=${KIOSK_PROTOS:-$GOPATH/src/github.com/googleapis/kiosk/protos}
SHOW_PROTOS=${SHOW_PROTOS:-$GOPATH/src/github.com/googleapis/gapic-showcase/schema}

KIOSK_GAPIC=${KIOSK_GAPIC:-github.com/googleapis/kiosk/kioskgapic}
SHOWCASE_GAPIC=${SHOWCASE_GAPIC:-github.com/googleapis/gapic-showcase/showgapic}

mkdir -p "$OUT/kiosk"
mkdir -p "$OUT/showcase"

generate() {
	protoc -I "$COMMON_PROTO" $*
}

generate -I $KIOSK_PROTOS \
  --gcli_out $OUT/kiosk \
  --gcli_opt "gapic:$KIOSK_GAPIC" \
  --gcli_opt 'root:testkctl' \
  $KIOSK_PROTOS/kiosk.proto

generate -I $SHOW_PROTOS \
  --gcli_out $OUT/showcase \
  --gcli_opt "gapic:$SHOWCASE_GAPIC" \
  --gcli_opt 'root:testshowctl' \
  --gcli_opt 'fmt:false' \
  $SHOW_PROTOS/*.proto

d=$(pwd)
cd $OUT/kiosk
go build
cd ../showcase
go build
cd $d

rm -r $OUT