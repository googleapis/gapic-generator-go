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

package schema_v1p2

const (
	// These values are used in Sample.SampleType to denote taht
	// the config applies to standalone samples or to in-code
	// (language doc) samples.
	sampleTypeStandalone = "standalone"
	sampleTypeDoc        = "incode"
)

type SampleConfig struct {
	Type    string
	Version string `yaml:"schema_version"`
	Samples []*Sample
}

type Sample struct {
	ID              string `yaml:"id"`
	Title           string
	RegionTag       string `yaml:"region_tag"`
	Description     string
	Service         string
	Rpc             string
	CallingPatterns []string `yaml:"calling_patterns"`
	Request         []RequestConfig
	Response        []ResponseConfig
	SampleType      []string `yaml:"sample_type"`
}

type RequestConfig struct {
	Field          string
	Comment        string
	Value          string
	ValueIsFile    bool   `yaml:"value_is_file"`
	InputParameter string `yaml:"input_parameter"`
}

type ResponseConfig struct {
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
	Body       []ResponseConfig
}

type WriteFileSpec struct {
	Contents string
	FileName []string `yaml: file_name`
}

// IsStandaloneSample returns true iff s.SampleType specifies s should
// generate standalone samples.
func (s *Sample) IsStandaloneSample() bool {
	if len(s.SampleType) == 0 {
		return true
	}
	for _, t := range s.SampleType {
		if t == sampleTypeStandalone {
			return true
		}
	}
	return false
}

// IsDocSample returns true iff s.SampleType specifies s should
// generate in-code (language doc) samples.
func (s *Sample) IsDocSample() bool {
	for _, t := range s.SampleType {
		if t == sampleTypeDoc {
			return true
		}
	}
	return false
}
