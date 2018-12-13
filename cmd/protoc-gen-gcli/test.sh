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

# setup variables
GCLI="github.com/googleapis/gapic-generator-go/cmd/protoc-gen-gcli"
GCLI_SRC=${GCLI_SRC:-$GOPATH/src/$GCLI}

OUT=${OUT:-$GCLI_SRC/testdata}
KIOSK_PROTOS=${KIOSK_PROTOS:-$GCLI_SRC/testprotos/kiosk}
SHOW_PROTOS=${SHOW_PROTOS:-$GCLI_SRC/testprotos/showcase}
COMMON_PROTOS=${COMMON_PROTOS:-$OUT/common}

KIOSK_GAPIC=${KIOSK_GAPIC:-$GCLI/testdata/kiosk/gapic}
SHOWCASE_GAPIC=${SHOWCASE_GAPIC:-$GCLI/testdata/showcase/gapic}

# make test output directories
mkdir -p "$OUT/kiosk"
mkdir -p "$OUT/showcase"

# download api-common-protos:input-contract branch
curl -L -O https://github.com/googleapis/api-common-protos/archive/input-contract.zip
unzip -q input-contract.zip
rm -f input-contract.zip
mv -f api-common-protos-input-contract $COMMON_PROTOS

# install gapic microgenerator plugin
go install "github.com/googleapis/gapic-generator-go/cmd/protoc-gen-go_gapic"

# install CLI generator plugin
go install $GCLI

generate() {
	protoc -I "$COMMON_PROTOS" $*
}

# generate kiosk gapic & gcli
generate -I $KIOSK_PROTOS \
  --go_out=plugins=grpc:$GOPATH/src \
  --go_gapic_out $GOPATH/src \
  --go_gapic_opt $KIOSK_GAPIC';gapic' \
  --gcli_out $OUT/kiosk \
  --gcli_opt "gapic=$KIOSK_GAPIC" \
  --gcli_opt "root=testkctl" \
  $KIOSK_PROTOS/kiosk.proto

# generate gapic-showcase gapic & gcli
generate -I $SHOW_PROTOS \
  --go_out=plugins=grpc:$GOPATH/src \
  --go_gapic_out $GOPATH/src \
  --go_gapic_opt $SHOWCASE_GAPIC';gapic' \
  --gcli_out $OUT/showcase \
  --gcli_opt "gapic=$SHOWCASE_GAPIC" \
  --gcli_opt "root=testshowctl" \
  --gcli_opt "fmt=false" \
  $SHOW_PROTOS/*.proto

# build each gcli for sanity check
d=$(pwd)
cd $OUT/kiosk
go build
cd ../showcase
go build
cd $d

# clean up
rm -r $OUT