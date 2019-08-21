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

package main

import (
	"testing"

	"github.com/googleapis/gapic-generator-go/cmd/gen-go-sample/schema_v1p2"
	"github.com/googleapis/gapic-generator-go/internal/errors"
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
	if gen.sampleConfig.Samples[0].ID != "sample0869f447dbdb23e5b6b08dec51059198ce28c361a" {
		t.Fatal(errors.E(nil, `expected "sample0869f447dbdb23e5b6b08dec51059198ce28c361a", got %q`, gen.sampleConfig.Samples[0].ID))
	}
	if gen.sampleConfig.Samples[1].ID != "sample019c7d6dbf4468cafdf72661f33c5b415108b807e" {
		t.Fatal(errors.E(nil, `expected "sample019c7d6dbf4468cafdf72661f33c5b415108b807e", got %q`, gen.sampleConfig.Samples[1].ID))
	}
	if gen.sampleConfig.Samples[2].ID != "sample1" {
		t.Fatal(errors.E(nil, `expected "sample1", got %q`, gen.sampleConfig.Samples[2].ID))
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
	if gen.sampleConfig.Samples[0].ID != "sample0869f447dbdb23e5b6b08dec51059198ce28c361a" {
		t.Fatal(errors.E(nil, `expected "sample0869f447dbdb23e5b6b08dec51059198ce28c361a", got %q`, gen.sampleConfig.Samples[0].ID))
	}
	if gen.sampleConfig.Samples[1].ID != "sample01cd8282ddcf344fc77617c2ed0cbae7c7a0bd637" {
		t.Fatal(errors.E(nil, `expected "sample01cd8282ddcf344fc77617c2ed0cbae7c7a0bd637", got %q`, gen.sampleConfig.Samples[1].ID))
	}
	if gen.sampleConfig.Samples[2].ID != "sample1" {
		t.Fatal(errors.E(nil, `expected "sample2", got %q`, gen.sampleConfig.Samples[2].ID))
	}
}
