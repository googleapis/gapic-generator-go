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

type GAPICConfig struct {
	Interfaces  []GAPICInterface
	Collections []ResourceName
}

type GAPICInterface struct {
	Name    string
	Methods []GAPICMethod
}

type GAPICMethod struct {
	Name string

	// map[fieldName]ResourceName.EntityName
	FieldNamePatterns map[string]string `yaml:"field_name_patterns"`

	LongRunning LongRunningConfig `yaml:"long_running"`
}

type ResourceName struct {
	EntityName  string `yaml:"entity_name"`
	NamePattern string `yaml:"name_pattern"`
}

// All other fields are left out because samples do not need to know polling config, and we are moving to annotations anyway
type LongRunningConfig struct {
	ReturnType   string `yaml:"return_type"`
	MetadataType string `yaml:"metadata_type"`
}
