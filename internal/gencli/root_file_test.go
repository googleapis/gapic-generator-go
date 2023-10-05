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

package gencli

import (
	"path/filepath"
	"testing"

	"github.com/googleapis/gapic-generator-go/internal/txtdiff"
)

func TestRootFile(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Generating the root_file panicked: %v", r)
		}
	}()

	g := &gcli{
		root:   "Root",
		format: true,
	}

	g.genRootCmdFile()
	if g.response.GetError() != "" {
		t.Errorf("Error generating the root_file: %s", g.response.GetError())
		return
	}

	file := g.response.File[0]

	if file.GetName() != "root.go" {
		t.Errorf("(%s).genRootCmdFile() = %s, want %s", g.root, file.GetName(), "root.go")
	}
	txtdiff.Diff(t, file.GetContent(), filepath.Join("testdata", "root_file.want"))
}
