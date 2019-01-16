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

package gengapic

import (
	"path/filepath"
	"testing"

	"github.com/googleapis/gapic-generator-go/internal/txtdiff"
)

func TestDocFile(t *testing.T) {
	var g generator
	g.apiName = "Awesome Foo API"
	g.serviceConfig = &serviceConfig{
		Documentation: &configDocumentation{
			Summary: "The Awesome Foo API is really really awesome. It enables the use of Foo with Buz and Baz to acclerate bar.",
		},
	}

	g.genDocFile("path/to/awesome", "awesome", 42, []string{"https://foo.bar.com/auth", "https://zip.zap.com/auth"})
	txtdiff.Diff(t, "doc_file", g.pt.String(), filepath.Join("testdata", "doc_file.want"))
}
