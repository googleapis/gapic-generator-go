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
GCLI="github.com/googleapis/gapic-generator-go/cmd/protoc-gen-go_cli"
GCLI_SRC=${GCLI_SRC:-$GOPATH/src/$GCLI}

OUT=${OUT:-$GCLI_SRC/testdata}
TESTPROTOS=${TESTPROTOS:-$GCLI_SRC/testprotos}
KIOSK_PROTOS=${KIOSK_PROTOS:-$GCLI_SRC/testprotos/kiosk}
SHOW_PROTOS=${SHOW_PROTOS:-$GCLI_SRC/testprotos/showcase}
COMMON_PROTOS=${COMMON_PROTOS:-$OUT/common}

KIOSK_GAPIC=${KIOSK_GAPIC:-$GCLI/testdata/kiosk/gapic}
SHOWCASE_GAPIC=${SHOWCASE_GAPIC:-$GCLI/testdata/showcase/gapic}

# make test output directories
mkdir -p "$OUT/kiosk"
mkdir -p "$OUT/showcase"

# make test proto directories
mkdir -p $KIOSK_PROTOS
mkdir -p $SHOW_PROTOS

# download api-common-protos:input-contract branch
curl -L -O https://github.com/googleapis/api-common-protos/archive/input-contract.zip
unzip -q input-contract.zip
rm -f input-contract.zip
mv -f api-common-protos-input-contract $COMMON_PROTOS

# download kiosk proto
curl -L -O https://raw.githubusercontent.com/googleapis/kiosk/master/protos/kiosk.proto
mv kiosk.proto $KIOSK_PROTOS/

# download gapic-showcase proto descriptor set
curl -L -O https://github.com/googleapis/gapic-showcase/releases/download/v0.0.16/gapic-showcase-0.0.16.desc
mv gapic-showcase-0.0.16.desc $SHOW_PROTOS/

# install gapic microgenerator plugin
go install "github.com/googleapis/gapic-generator-go/cmd/protoc-gen-go_gapic"

# install CLI generator plugin
go install $GCLI

# generate kiosk gapic & go_cli
protoc -I $KIOSK_PROTOS \
  -I $COMMON_PROTOS \
  --go_out=plugins=grpc:$GOPATH/src \
  --go_gapic_out $GOPATH/src \
  --go_gapic_opt "go-gapic-package=$KIOSK_GAPIC"';gapic' \
  --go_cli_out $OUT/kiosk \
  --go_cli_opt "gapic=$KIOSK_GAPIC" \
  --go_cli_opt "root=testkctl" \
  $KIOSK_PROTOS/kiosk.proto

# generate gapic-showcase gapic & go_cli
protoc --descriptor_set_in=$SHOW_PROTOS/gapic-showcase-0.0.16.desc \
  --go_out=plugins=grpc:$GOPATH/src \
  --go_gapic_out $GOPATH/src \
  --go_gapic_opt "go-gapic-package=$SHOWCASE_GAPIC"';gapic' \
  --go_cli_out $OUT/showcase \
  --go_cli_opt "gapic=$SHOWCASE_GAPIC" \
  --go_cli_opt "root=testshowctl" \
  --go_cli_opt "fmt=false" \
  google/showcase/v1alpha3/echo.proto \
  google/showcase/v1alpha3/identity.proto \
  google/showcase/v1alpha3/messaging.proto \
  google/showcase/v1alpha3/testing.proto

# build each go_cli for sanity check
d=$(pwd)
cd $OUT/kiosk
go build
cd ../showcase
go build
cd $d

# clean up
rm -r $OUT
rm -r $TESTPROTOS