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

if [[ "${BASH_SOURCE[0]}" != "${0}" ]]; then
	echo >&2 "don't source; execute the script instead"
	return
fi

# Capture the version of gapic-showcase dependency listed in the go.mod to
# download release assets.
SHOWCASE_SEMVER=
regex="showcase v([0-9]+\.[0-9]+\.[0-9]+)"
mod=`cat go.mod`

if [[ $mod =~ $regex ]]
then
	SHOWCASE_SEMVER="${BASH_REMATCH[1]}"
	echo "Using Showcase version ${SHOWCASE_SEMVER}"
else
	echo "Could not find a proper version for showcase in go.mod: ${mod}" >&2
	exit 1
fi

rm -rf gen
mkdir gen

curl -L -O https://github.com/googleapis/gapic-showcase/releases/download/v$SHOWCASE_SEMVER/showcase_grpc_service_config.json

curl -L -O https://github.com/googleapis/gapic-showcase/releases/download/v$SHOWCASE_SEMVER/showcase_v1beta1.yaml

protoc \
	--experimental_allow_proto3_optional \
	--go_out ./gen \
	--go-grpc_out ./gen \
	--go_gapic_out ./gen \
	--go_gapic_opt 'transport=rest+grpc' \
	--go_gapic_opt 'rest-numeric-enums' \
	--go_gapic_opt 'go-gapic-package=github.com/googleapis/gapic-showcase/client;client' \
	--go_gapic_opt 'grpc-service-config=showcase_grpc_service_config.json' \
	--go_gapic_opt 'api-service-config=showcase_v1beta1.yaml' \
	--descriptor_set_in=<(curl -sSL https://github.com/googleapis/gapic-showcase/releases/download/v$SHOWCASE_SEMVER/gapic-showcase-$SHOWCASE_SEMVER.desc) \
	google/showcase/v1beta1/echo.proto \
	google/showcase/v1beta1/identity.proto \
	google/showcase/v1beta1/sequence.proto \
	google/showcase/v1beta1/testing.proto \
	google/showcase/v1beta1/messaging.proto \
	google/showcase/v1beta1/compliance.proto

hostos=$(go env GOHOSTOS)
hostarch=$(go env GOHOSTARCH)

pushd gen/github.com/googleapis/gapic-showcase
go mod init github.com/googleapis/gapic-showcase
# Fixes a name collision with the operation helper WaitOperation by renaming the mixin method.
if [[ "$hostos" == "darwin" ]]; then
    SEDARGS="-i ''"
else
    SEDARGS="-i"
fi
sed $SEDARGS '1,/WaitOperation(ctx/{s/WaitOperation(ctx/WaitOperationMixin(ctx/;}' client/echo_client*

popd

go mod edit -replace=github.com/googleapis/gapic-showcase=./gen/github.com/googleapis/gapic-showcase

curl -sSL https://github.com/googleapis/gapic-showcase/releases/download/v$SHOWCASE_SEMVER/gapic-showcase-$SHOWCASE_SEMVER-$hostos-$hostarch.tar.gz | tar xz
curl -sSL -O https://github.com/googleapis/gapic-showcase/releases/download/v$SHOWCASE_SEMVER/compliance_suite.json

./gapic-showcase run &
showcase_pid=$!

cleanup() {
	kill $showcase_pid
	# Wait for the process to die, but don't report error from the kill.
	wait $showcase_pid || exit $exit_code
}
trap cleanup EXIT

go test -mod=mod -count=1 ./...
exit_code=$?
