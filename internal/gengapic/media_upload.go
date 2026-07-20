// Copyright 2026 Google LLC
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

package gengapic

import "google.golang.org/protobuf/types/descriptorpb"

// isMediaUpload evaluates if a given RPC method relies on media upload.
func (g *generator) isMediaUpload(m *descriptorpb.MethodDescriptorProto) bool {
	if !g.featureEnabled(MediaUploadFeature) {
		// Disallow detection if the feature isn't enabled.
		return false
	}

	// TODO: implement detection and config flow.
	return false
}
