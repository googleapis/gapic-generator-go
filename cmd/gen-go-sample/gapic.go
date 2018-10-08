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

package main

type GAPICConfig struct {
	Interfaces []GAPICInterface
}

type GAPICInterface struct {
	Name    string
	Methods []GAPICMethod
}

type GAPICMethod struct {
	Name            string
	SampleValueSets []SampleValueSet `yaml:"sample_value_sets"`
	Samples         struct {
		Standalone []GAPICSample
	}
}

type SampleValueSet struct {
	ID         string
	Parameters SampleParameter
	OnSuccess  []OutSpec `yaml:"on_success"`
}

type SampleParameter struct {
	Defaults   []string
	Attributes []SampleAttribute
}

type SampleAttribute struct {
	Parameter      string
	SampleArgument bool `yaml:"sample_argument"`
}

type GAPICSample struct {
	ValueSets []string `yaml:"value_sets"`

	// TODO(pongad): Does this mean multiple samples have the same tag?
	RegionTag string `yaml:"region_tag"`
}

type OutSpec struct {
	Define string
	Print  []string
}
