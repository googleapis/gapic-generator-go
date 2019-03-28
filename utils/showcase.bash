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

go install ./cmd/protoc-gen-go_gapic

SHOWCASE_SEMVER=0.0.16

pushd showcase
rm -rf gen
mkdir gen

protoc \
	--go_gapic_out ./gen \
	--go_gapic_opt 'go-gapic-package=cloud.google.com/go/showcase/apiv1alpha3;showcase' \
	--descriptor_set_in=<(curl -sSL https://github.com/googleapis/gapic-showcase/releases/download/v$SHOWCASE_SEMVER/gapic-showcase-$SHOWCASE_SEMVER.desc) \
	google/showcase/v1alpha3/echo.proto

pushd gen/cloud.google.com/go/showcase
go mod init cloud.google.com/go/showcase
popd

hostos=$(go env GOHOSTOS)
hostarch=$(go env GOHOSTARCH)

curl -sSL https://github.com/googleapis/gapic-showcase/releases/download/v$SHOWCASE_SEMVER/gapic-showcase-$SHOWCASE_SEMVER-$hostos-$hostarch.tar.gz | tar xz
./gapic-showcase run &
showcase_pid=$!

stop_showcase() {
	kill $showcase_pid
	# Wait for the process to die, but don't report error from the kill.
	wait $showcase_pid || true
}
trap stop_showcase EXIT

go test -count=1 ./...
popd
