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

	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"github.com/googleapis/gapic-generator-go/internal/txtdiff"
)

func TestClientUtils(t *testing.T) {
	var imports = map[pbinfo.ImportSpec]bool{
		pbinfo.ImportSpec{Path: "runtime"}: true,
		pbinfo.ImportSpec{Path: "strings"}: true,
		pbinfo.ImportSpec{Path: "unicode"}: true,
	}
	var g generator
	g.imports = map[pbinfo.ImportSpec]bool{}

	g.clientUtils()
	txtdiff.Diff(t, "util_file", g.pt.String(), filepath.Join("testdata", "util_file.want"))

	for k, v := range imports {
		if g.imports[k] != v {
			t.Errorf("clientUtils() - g.imports[%v] = %v, want %v", k, g.imports[k], v)
		}
	}
}
