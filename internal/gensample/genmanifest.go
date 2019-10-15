// Copyright 2019 Google LLC
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

package gensample

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/googleapis/gapic-generator-go/internal/errors"
)

const (
	latestSchemaVersion = "3"
)

// genManifest generates manifest file to be consumed by sample-tester.
// See https://sample-tester.readthedocs.io/en/stable/defining-tests/manifest-reference.html
// for the format of the manifest file.
func (gen *generator) genManifest() error {
	// Do not generate sample manifest when there are no sample configs
	if len(gen.sampleConfig.Samples) == 0 {
		return nil
	}

	var b bytes.Buffer
	p := func(s string, args ...interface{}) {
		fmt.Fprintf(&b, s, args...)
		b.WriteByte('\n')
	}

	p("---")
	p("type: manifest/samples")
	p("schema_version: %s", latestSchemaVersion)
	p("go: &go")
	p("  environment: go")
	p("  invocation: go run {path}")
	p("  chdir: {@manifest_dir}/")
	p("samples:")

	for _, sp := range gen.sampleConfig.Samples {
		p("- <<: *go")
		p("  sample: %s", sp.ID)
		p("  path: %s.go", sp.ID)
	}

	api := gen.clientPkg.Name
	pos := strings.LastIndexByte(gen.clientPkg.Path, '/')
	if pos < 0 {
		return errors.E(nil, "expecting clientPkg path in 'url/to/client/apiversion/' format, got %q", gen.clientPkg.Path)
	}
	v := strings.Replace(gen.clientPkg.Path[pos+1:], "api", "", -1)
	gen.Outputs[fmt.Sprintf("%s.%s.go.manifest.yaml", api, v)] = b.Bytes()
	return nil
}
