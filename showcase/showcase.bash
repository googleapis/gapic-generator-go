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

SHOWCASE_SEMVER=0.7.0

rm -rf gen
mkdir gen

curl -L -O https://github.com/googleapis/gapic-showcase/releases/download/v$SHOWCASE_SEMVER/showcase_grpc_service_config.json

protoc \
	--go_out=plugins=grpc:./gen \
	--go_gapic_out ./gen \
	--go_gapic_opt 'go-gapic-package=github.com/googleapis/gapic-showcase/client;client' \
	--go_gapic_opt 'grpc-service-config=showcase_grpc_service_config.json' \
	--descriptor_set_in=<(curl -sSL https://github.com/googleapis/gapic-showcase/releases/download/v$SHOWCASE_SEMVER/gapic-showcase-$SHOWCASE_SEMVER.desc) \
	google/showcase/v1beta1/echo.proto google/showcase/v1beta1/identity.proto

pushd gen/github.com/googleapis/gapic-showcase
go mod init github.com/googleapis/gapic-showcase
popd

go mod edit -replace=github.com/googleapis/gapic-showcase=./gen/github.com/googleapis/gapic-showcase

hostos=$(go env GOHOSTOS)
hostarch=$(go env GOHOSTARCH)

curl -sSL https://github.com/googleapis/gapic-showcase/releases/download/v$SHOWCASE_SEMVER/gapic-showcase-$SHOWCASE_SEMVER-$hostos-$hostarch.tar.gz | tar xz
./gapic-showcase run &
showcase_pid=$!

go test -count=1 ./...
test_exit=$?

cleanup() {
	kill $showcase_pid
	# Wait for the process to die, but don't report error from the kill.
	wait $showcase_pid || exit $test_exit
}
trap cleanup EXIT
