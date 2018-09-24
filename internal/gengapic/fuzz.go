// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build gofuzz

package gengapic

import plugin "github.com/golang/protobuf/protoc-gen-go/plugin"

func Fuzz(data []byte) int {
	var genReq plugin.CodeGeneratorRequest
	if err := genReq.Unmarshal(data); err != nil {
		// -1: don't include these in corpus even if they create new coverage.
		// We're not interested in bytes that don't deserialize properly.
		return -1
	}

	// It's OK if we error. Just test that we don't crash.
	Gen(&genReq)
	return 1
}
