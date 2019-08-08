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
	Interfaces  []GAPICInterface
	Collections []ResourceName
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

	// map[fieldName]ResourceName.EntityName
	FieldNamePatterns map[string]string `yaml:"field_name_patterns"`

	LongRunning LongRunningConfig `yaml: "long_running"`
}

type SampleValueSet struct {
	ID         string
	Parameters SampleParameter
	OnSuccess  []OutputSpec `yaml:"on_success"`
}

type SampleParameter struct {
	Defaults   []string
	Attributes []SampleAttribute
}

type SampleAttribute struct {
	Parameter          string
	SampleArgumentName string `yaml:"sample_argument_name"`
	ReadFile           bool   `yaml:"read_file"`
}

type GAPICSample struct {
	ValueSets []string `yaml:"value_sets"`

	// TODO(pongad): Does this mean multiple samples have the same tag?
	RegionTag string `yaml:"region_tag"`
}

type OutputSpec struct {
	Comment   []string
	Define    string
	Print     []string
	Loop      *LoopSpec
	WriteFile *WriteFileSpec
}

type LoopSpec struct {
	Collection string
	Map        string
	Variable   string
	Key        string
	Value      string
	Body       []OutputSpec
}

type WriteFileSpec struct {
	Contents string
	FileName []string `yaml: file_name`
}

type ResourceName struct {
	EntityName  string `yaml:"entity_name"`
	NamePattern string `yaml:"name_pattern"`
}

// All other fields are left out because samples do not need to know polling config, and we are moving to annotations anyway
type LongRunningConfig struct {
	ReturnType   string `yaml: "return_type"`
	MetadataType string `yaml: "metadata_type"`
}
