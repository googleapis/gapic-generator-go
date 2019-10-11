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
)

const (
	latestSchemaVersion = "3"
)


func (gen *generator) genManifest() {
	// Do not generate sample manifest when there are no sample configs
	if len(gen.sampleConfig.Samples) == 0 {
		return
	}
	
	var b bytes.Buffer
	p := func(s string, args ...interface{}) {
		fmt.Fprintf(&b, s, args...)
		b.WriteByte('\n')
	}

	p("type: manifest/samples")
	p("schema_version: %s", latestSchemaVersion)
	p("go: &go")
	p("  environment: go")
	p("  bin: go run")
	p("  chdir: {@manifest_dir}/")
	p("samples:")

	for _, sp := range gen.sampleConfig.Samples {
		p("- <<: *go")
		p("  sample: %s", sp.ID)
		p("  path: %s.go", sp.ID)
	}

	gen.Outputs["samples.manifest.yaml"] = b.Bytes()
}
