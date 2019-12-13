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
	"strings"
	"testing"

	yaml "gopkg.in/yaml.v2"

	"github.com/googleapis/gapic-generator-go/internal/errors"
	"github.com/googleapis/gapic-generator-go/internal/gensample/schema_v1p2"
)

func TestDisambiguateUniqueSampleIDs(t *testing.T) {
	gen := &generator{
		sampleConfig: schema_v1p2.SampleConfig{
			Samples: []*schema_v1p2.Sample{
				&schema_v1p2.Sample{ID: "sample0", Service: "foo.FooService", Rpc: "CreateFoo"},
				&schema_v1p2.Sample{ID: "sample1", Service: "foo.FooService", Rpc: "UpdateFoo"},
				&schema_v1p2.Sample{ID: "sample2", Service: "foo.FooService", Rpc: "DeleteFoo"},
			},
		},
	}

	gen.disambiguateSampleIDs()
	if gen.sampleConfig.Samples[0].ID != "sample0" {
		t.Fatal(errors.E(nil, `expected "sample0", got %q`, gen.sampleConfig.Samples[0].ID))
	}
	if gen.sampleConfig.Samples[1].ID != "sample1" {
		t.Fatal(errors.E(nil, `expected "sample1", got %q`, gen.sampleConfig.Samples[1].ID))
	}
	if gen.sampleConfig.Samples[2].ID != "sample2" {
		t.Fatal(errors.E(nil, `expected "sample2", got %q`, gen.sampleConfig.Samples[2].ID))
	}
}

func TestDisambiguateRepeatedSampleIDs(t *testing.T) {
	gen := &generator{
		sampleConfig: schema_v1p2.SampleConfig{
			Samples: []*schema_v1p2.Sample{
				&schema_v1p2.Sample{ID: "sample0", Service: "foo.FooService", Rpc: "CreateFoo"},
				&schema_v1p2.Sample{ID: "sample0", Service: "foo.FooService", Rpc: "UpdateFoo"},
				&schema_v1p2.Sample{ID: "sample1", Service: "foo.FooService", Rpc: "DeleteFoo"},
			},
		},
	}

	gen.disambiguateSampleIDs()
	if expected := "sample0C57ISC5Q"; gen.sampleConfig.Samples[0].ID != expected {
		t.Fatal(errors.E(nil, `expected %q, got %q`, expected, gen.sampleConfig.Samples[0].ID))
	}
	if expected := "sample05L2GRN22"; gen.sampleConfig.Samples[1].ID != expected {
		t.Fatal(errors.E(nil, `expected %q, got %q`, expected, gen.sampleConfig.Samples[1].ID))
	}
	if expected := "sample1"; gen.sampleConfig.Samples[2].ID != expected {
		t.Fatal(errors.E(nil, `expected %q, got %q`, expected, gen.sampleConfig.Samples[2].ID))
	}
}

func TestDisambiguateRepeatedSampleIDsFromRegionTags(t *testing.T) {
	gen := &generator{
		sampleConfig: schema_v1p2.SampleConfig{
			Samples: []*schema_v1p2.Sample{
				&schema_v1p2.Sample{ID: "sample0", Service: "foo.FooService", Rpc: "CreateFoo"},
				&schema_v1p2.Sample{RegionTag: "sample0", Service: "foo.FooService", Rpc: "UpdateFoo"},
				&schema_v1p2.Sample{ID: "sample1", Service: "foo.FooService", Rpc: "DeleteFoo"},
			},
		},
	}

	gen.disambiguateSampleIDs()
	if expected := "sample0C57ISC5Q"; gen.sampleConfig.Samples[0].ID != expected {
		t.Fatal(errors.E(nil, `expected %q, got %q`, expected, gen.sampleConfig.Samples[0].ID))
	}
	if expected := "sample0MHV4NHU4"; gen.sampleConfig.Samples[1].ID != expected {
		t.Fatal(errors.E(nil, `expected %q, got %q`, expected, gen.sampleConfig.Samples[1].ID))
	}
	if expected := "sample1"; gen.sampleConfig.Samples[2].ID != expected {
		t.Fatal(errors.E(nil, `expected %q, got %q`, expected, gen.sampleConfig.Samples[2].ID))
	}
}

func TestReadSampleConfigSingleDoc(t *testing.T) {
	yamlStr := `
type: com.google.api.codegen.samplegen.v1p2.SampleConfigProto
schema_version: 1.2.0
samples:
- service: FooService
  rpc: CreateFoo
  region_tag: awesome_region
`
	decoder := yaml.NewDecoder(strings.NewReader(yamlStr))
	config, err := readSampleConfig(decoder, "some.yaml")
	if err != nil {
		t.Fatal(errors.E(nil, "unexpected error: %q", err.Error()))
	}
	if len(config.Samples) != 1 {
		t.Fatal(errors.E(nil, "expected one valid sample, got %d", len(config.Samples)))
	}
}

func TestReadSampleConfigSingleDoc_BadTypeError(t *testing.T) {
	yamlStr := `
type: com.google.api.codegen.samplegen.v1p2.NotSampleConfigProto
schema_version: 1.2.0
samples:
- service: FooService
  rpc: CreateFoo
  region_tag: awesome_region
`
	decoder := yaml.NewDecoder(strings.NewReader(yamlStr))
	_, err := readSampleConfig(decoder, "some.yaml")
	expectedErrMsg := `Found no valid sample config in "some.yaml"`
	if err == nil || err.Error() != expectedErrMsg {
		t.Fatal(errors.E(nil, "expected error with message %q, got %q", expectedErrMsg, err.Error()))
	}
}

func TestReadSampleConfigSingleDoc_BadSchemaVersionError(t *testing.T) {
	yamlStr := `
type: com.google.api.codegen.samplegen.v1p2.SampleConfigProto
schema_version: 1.3.0
samples:
- service: FooService
  rpc: CreateFoo
  region_tag: awesome_region
`
	decoder := yaml.NewDecoder(strings.NewReader(yamlStr))
	_, err := readSampleConfig(decoder, "some.yaml")
	expectedErrMsg := `Found no valid sample config in "some.yaml"`
	if err == nil || err.Error() != expectedErrMsg {
		t.Fatal(errors.E(nil, "expected error with message %q, got %q", expectedErrMsg, err.Error()))
	}
}

func TestReadSampleConfig_EmptyYamlError(t *testing.T) {
	yamlStr := ""
	decoder := yaml.NewDecoder(strings.NewReader(yamlStr))
	_, err := readSampleConfig(decoder, "some.yaml")
	expectedErrMsg := `Found no valid sample config in "some.yaml"`
	if err == nil || err.Error() != expectedErrMsg {
		t.Fatal(errors.E(nil, "expected error with message %q, got %q", expectedErrMsg, err.Error()))
	}
}

func TestReadSampleConfigMultiDoc(t *testing.T) {
	yamlStr := `
type: com.google.api.codegen.samplegen.v1p2.SampleConfigProto
schema_version: 1.2.0
samples:
- service: FooService
  rpc: CreateFoo
  region_tag: awesome_region
---
type: com.google.api.codegen.samplegen.v1p2.SampleConfigProto
schema_version: 1.2.0
samples:
- service: FooService
  rpc: GetFoo
  region_tag: lame_region
`
	decoder := yaml.NewDecoder(strings.NewReader(yamlStr))
	config, err := readSampleConfig(decoder, "some.yaml")
	if err != nil {
		t.Fatal(errors.E(nil, "unexpected error: %q", err.Error()))
	}
	if config.Samples[0].Rpc != "CreateFoo" {
		t.Fatal(errors.E(nil, `expected "CreateFoo", got %q`, config.Samples[0].Rpc))
	}
	if config.Samples[1].Rpc != "GetFoo" {
		t.Fatal(errors.E(nil, `expected "GetFoo", got %q`, config.Samples[1].Rpc))
	}
}

func TestReadSampleConfigMultiDocSkipNonSampleConfig(t *testing.T) {
	yamlStr := `
type: com.google.api.codegen.samplegen.v1p2.SampleConfigProto
schema_version: 1.2.0
samples:
- service: FooService
  rpc: CreateFoo
  region_tag: awesome_region
---
type: another_type
`
	decoder := yaml.NewDecoder(strings.NewReader(yamlStr))
	config, err := readSampleConfig(decoder, "some.yaml")
	if err != nil {
		t.Fatal(errors.E(nil, "unexpected error: %q", err.Error()))
	}
	if len(config.Samples) != 1 {
		t.Fatal(errors.E(nil, "expected one valid sample, got %d", len(config.Samples)))
	}
}

func TestReadSampleConfigMultiDoc_BadFormatError(t *testing.T) {
	yamlStr := `
type: com.google.api.codegen.samplegen.v1p2.SampleConfigProto
schema_version: 1.2.0
samples:
- service: FooService
  rpc: CreateFoo
  region_tag: awesome_region
---
type: another_type
  bad_field: derp
`
	decoder := yaml.NewDecoder(strings.NewReader(yamlStr))
	_, err := readSampleConfig(decoder, "some.yaml")
	if err == nil {
		t.Fatal(errors.E(nil, "expected error"))
	}
}
