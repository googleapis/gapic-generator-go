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

package gensample

import (
	"testing"

	"github.com/googleapis/gapic-generator-go/internal/errors"
	"github.com/googleapis/gapic-generator-go/internal/gensample/schema_v1p2"
	"github.com/googleapis/gapic-generator-go/internal/txtdiff"
)

func TestGenManifest(t *testing.T) {
	t.Parallel()

	g := &generator{
		sampleConfig: schema_v1p2.SampleConfig {
			Samples: []*schema_v1p2.Sample{
				&schema_v1p2.Sample {
					ID: "first_sample",
				},
				&schema_v1p2.Sample {
					ID: "second_sample",
				},
			},
		},
		Outputs: make(map[string][]byte),
	}

	g.genManifest()
	got, ok := g.Outputs["samples.manifest.yaml"]
	if !ok {
		t.Fatal(errors.E(nil, "manifest file not generated"))
	}

	compareManifest(t, got, "testdata/manifest.want")
}

func TestGenManifest_NoSamples(t *testing.T) {
	t.Parallel()

	g := &generator{
		Outputs: make(map[string][]byte),
	}

	g.genManifest()
	_, ok := g.Outputs["samples.manifest.yaml"]
	if ok {
		t.Fatal(errors.E(nil, "manifest file should not be generated when there are no sample configs."))
	}
}

func compareManifest(t *testing.T, got []byte, goldenPath string) {
	t.Helper()
	txtdiff.Diff(t, t.Name(), string(got), goldenPath)
}
