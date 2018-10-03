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

# Setup a test workspace. We gitignore this.
rm -rf showcase-testdir
mkdir showcase-testdir/

cd showcase-testdir
go mod init showcase-test

protoc \
	--go_gapic_out . \
	--go_gapic_opt 'cloud.google.com/go/showcase/apiv1alpha2;showcase' \
	--descriptor_set_in=<(curl -sSL https://github.com/googleapis/gapic-showcase/releases/download/v0.0.7/gapic-showcase-0.0.7.desc) \
	google/showcase/v1alpha2/echo.proto

# With Go modules, we're transitioning off version package.
# For now this hack will keep things passing.
mkdir -p cloud.google.com/go/internal/version
cat > cloud.google.com/go/internal/version/version.go <<EOF
package version

const Repo = "UNKNOWN"
func Go() string {return "UNKNOWN"}
EOF

# TODO(pongad): Move this file into this repository once we deprecate the old generator.
curl -sSL https://raw.githubusercontent.com/googleapis/gapic-generator/753ff9d8a04a59962f3b8c2c06cb79be7df344c8/showcase/go/showcase_integration_test.go \
	> showcase_integration_test.go

# If the package name of the generated package is "a/b/c", we write files under directory "$OUT/a/b/c",
# where OUT is the directory passed to protoc. Prior to Go 1.11, we set OUT=$GOPATH/src.
# This works well because that's where we expect to find the package a/b/c.
#
# This plays badly with modules however. In `go mod init` above, we set our module name to "showcase-test",
# and set module root to showcase-testdir (because that's where the file is). If we set OUT=showcase-testdir,
# we generate the files to showcase-testdir/a/b/c, but the import path of this package would actually be showcase-testdir/a/b/c,
# not a/b/c as expected from the previous paragraph. We have to rewrite imports to keep this working.
#
# We copy this behavior from protoc-gen-go, so they're probably having the same problem with modules. Therefore...
# TODO(pongad): figure out what protoc-gen-go is doing to solve this then do the same thing.
find -name '*.go' | xargs sed -i 's,cloud.google.com/go/internal/version,showcase-test/&,g'
find -name '*.go' | xargs sed -i 's,cloud.google.com/go/showcase/apiv1alpha2,showcase-test/&,g'

curl -sSL https://github.com/googleapis/gapic-showcase/releases/download/v0.0.7/gapic-showcase-0.0.7-linux-amd64.tar.gz | tar xz
./gapic-showcase start &
showcase_pid=$!

stop_showcase() {
	kill $showcase_pid
}
trap stop_showcase EXIT

go test -count=1 ./...
