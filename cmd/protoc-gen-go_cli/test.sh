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

set -e

# setup variables
GCLI="github.com/googleapis/gapic-generator-go/cmd/protoc-gen-go_cli"
GCLI_SRC=${GCLI_SRC:-cmd/protoc-gen-go_cli}

OUT=${OUT:-$GCLI_SRC/testdata}
TESTPROTOS=${TESTPROTOS:-$GCLI_SRC/testprotos}
SHOW_PROTOS=${SHOW_PROTOS:-$GCLI_SRC/testprotos/showcase}

# make test output directories
mkdir -p "$OUT/showcase"

# make test proto directories
mkdir -p $SHOW_PROTOS

SHOWCASE_VERSION=0.19.0

# download gapic-showcase proto descriptor set
curl -L -O https://github.com/googleapis/gapic-showcase/releases/download/v$SHOWCASE_VERSION/gapic-showcase-$SHOWCASE_VERSION.desc
mv gapic-showcase-$SHOWCASE_VERSION.desc $SHOW_PROTOS/

# install CLI generator plugin
go install $GCLI

# generate gapic-showcase gapic & go_cli
protoc --descriptor_set_in=$SHOW_PROTOS/gapic-showcase-$SHOWCASE_VERSION.desc \
  --go_cli_out $OUT/showcase \
  --go_cli_opt "gapic=github.com/googleapis/gapic-showcase/client;client" \
  --go_cli_opt "root=testshowctl" \
  --go_cli_opt "fmt=false" \
  google/showcase/v1beta1/echo.proto \
  google/showcase/v1beta1/identity.proto \
  google/showcase/v1beta1/messaging.proto \
  google/showcase/v1beta1/testing.proto

# build showcase go_cli for sanity check
d=$(pwd)
cd $OUT/showcase
go mod init github.com/googleapis/gapic-generator-go/testshowctl
go build
cd $d

# clean up
rm -r $OUT
rm -r $TESTPROTOS