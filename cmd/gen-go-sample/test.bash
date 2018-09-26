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

if [ -z $GOOGLEAPIS ]; then
	echo >&2 "need GOOGLEAPIS variable"
	exit 1
fi
if [ -z $COMMON_PROTO ]; then
	echo >&2 "need COMMON_PROTO variable"
	exit 1
fi

./gen-go-sample $* \
  -clientpkg 'cloud.google.com/go/language/apiv1;language' \
  -gapic "$GOOGLEAPIS/google/cloud/language/v1/language_gapic.yaml" \
  -desc <(protoc -o /dev/stdout --include_imports -I "$COMMON_PROTO" -I "$GOOGLEAPIS" "$GOOGLEAPIS"/google/cloud/language/v1/*.proto)
