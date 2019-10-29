#!/bin/bash
# Copyright 2018 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# THIS SCRIPT IS MEANT ONLY TO BE USED IN THE GAPIC-GENERATOR-GO DOCKER IMAGE

GO_GAPIC_PACKAGE=
GAPIC_SERVICE_CONFIG=
GRPC_SERVICE_CONFIG=
GAPIC_CONFIG=
SAMPLES=
SAMPLE_ONLY=
RELEASE_LEVEL=

# enable extended globbing for flag pattern matching
shopt -s extglob

# Parse out options.
while true; do
  case "$1" in
    --go-gapic-package ) GO_GAPIC_PACKAGE="go-gapic-package=$2"; shift 2 ;;
    --gapic-service-config ) GAPIC_SERVICE_CONFIG="gapic-service-config=/conf/$2"; shift 2;;
    --grpc-service-config ) GRPC_SERVICE_CONFIG="grpc-service-config=/conf/$2"; shift 2;;
    --gapic-config ) GAPIC_CONFIG="gapic-config=/conf/$2"; shift 2;;
    --sample ) SAMPLES="${SAMPLES}sample=/conf/$2,"; shift 2;;
    --sample-only ) SAMPLE_ONLY="sample-only"; shift 1;;
    --release-level ) RELEASE_LEVEL="release-level=$2"; shift 2 ;; 
    --go-gapic* ) echo "Skipping unrecognized go-gapic flag: $1" >&2; shift ;;
    --* | +([[:word:][:punct:]]) ) shift ;;
    * ) break ;;
  esac
done

if [ -z "$GO_GAPIC_PACKAGE" ]; then
  echo "Required argument --go-gapic-package missing."
  exit 64
fi

# Trim the last comma from $SAMPLES
if [ ! -z "$SAMPLES" ]; then
  SAMPLES=${SAMPLES::-1}
fi

protoc --proto_path=/protos/ --proto_path=/in/ \
                  --go_gapic_out=/out/ \
                  --go_gapic_opt="$GO_GAPIC_PACKAGE" \
                  --go_gapic_opt="$RELEASE_LEVEL" \
                  --go_gapic_opt="$GAPIC_SERVICE_CONFIG" \
                  --go_gapic_opt="$GRPC_SERVICE_CONFIG" \
                  --go_gapic_opt="$GAPIC_CONFIG" \
                  --go_gapic_opt="$SAMPLES" \
                  --go_gapic_opt="$SAMPLE_ONLY" \
                  `find /in/ -name *.proto`
