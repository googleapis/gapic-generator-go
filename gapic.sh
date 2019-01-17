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
set -e 

CMD="$0"

# Set variables used by this script.
# All of these are set in options below, and all but $PATH are required.
IMAGE=
IN=
OUT=
PLUGIN_OPTIONS=
PROTO_PATH=`pwd`

# Print help and exit.
function show_help {
  cat << EOF
Usage: $CMD --image IMAGE --in IN_DIR --out OUT_DIR [--path PATH_DIR]

Required arguments:
      --image     The Docker image to use. The script will attempt to pull
                    it if it is not present.
  -i, --in        A directory containing the protos describing the API
                    to be generated.
  -o, --out       Destination directory for the completed client library.

Optional arguments:
  -p, --path      The base import path for the protos. Assumed to be the
                    current working directory if unspecified.
  -h, --help      This help information.
EOF
 
  exit 0
}

# Parse out options.
while true; do
  case "$1" in
    -h | --help ) show_help ;;
    --image ) IMAGE="$2"; shift 2 ;;
    -i | --in ) IN="$2"; shift 2 ;;
    -o | --out ) OUT="$2"; shift 2 ;;
    -p | --path ) PROTO_PATH=$2; shift 2 ;;
    --* ) PLUGIN_OPTIONS="$PLUGIN_OPTIONS $1 $2"; shift 2 ;;
    -- ) shift; break; ;;
    * ) break ;;
  esac
done

# Ensure that all required options are set.
if [ -z "$IMAGE" ] || [ -z "$IN" ] || [ -z "$OUT" ]; then
  cat << EOF
Required argument missing.
The --image, --in, and --out arguments are all required.
Run $CMD --help for more information.
EOF

  exit 64
fi

# Ensure that the input directory exists (and is a directory).
if ! [ -d $IN ]; then
  cat << EOF
Directory does not exist: $IN
EOF
  exit 2
fi

# Ensure Docker is running and seems healthy.
# This is mostly a check to bubble useful errors quickly.
docker ps > /dev/null

# If the output directory does not exist, create it.
mkdir -p $OUT

# If the output directory is not empty, warn (but continue).
if [ "$(ls -A $OUT )" ]; then
  cat << EOF
Warning: Output directory is not empty.
EOF
fi

# Generate the client library.
docker run \
  --mount type=bind,source=${PROTO_PATH},destination=/conf,readonly \
  --mount type=bind,source=${PROTO_PATH}/${IN},destination=/in/${IN},readonly \
  --mount type=bind,source=$OUT,destination=/out \
  --rm \
  --user $UID \
  $IMAGE \
  $PLUGIN_OPTIONS
