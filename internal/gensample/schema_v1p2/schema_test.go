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

import "testing"

func TestDetermineSampleType(t *testing.T) {
	for idx, tcase := range []struct {
		sampleType           []string
		expectTypeStandalone bool
		expectTypeDoc        bool
	}{
		{
			sampleType:           nil,
			expectTypeStandalone: true,
		},
		{
			sampleType:           []string{},
			expectTypeStandalone: true,
		},
		{
			sampleType:           []string{sampleTypeStandalone},
			expectTypeStandalone: true,
		},
		{
			sampleType:           []string{"foo", sampleTypeStandalone},
			expectTypeStandalone: true,
		},
		{
			sampleType:    []string{sampleTypeDoc},
			expectTypeDoc: true,
		},
		{
			sampleType:    []string{"bar", sampleTypeDoc},
			expectTypeDoc: true,
		},
		{
			sampleType:           []string{sampleTypeDoc, sampleTypeStandalone},
			expectTypeStandalone: true,
			expectTypeDoc:        true,
		},
		{
			sampleType:           []string{"foo", sampleTypeStandalone, "bar", sampleTypeDoc, "baz"},
			expectTypeStandalone: true,
			expectTypeDoc:        true,
		},
		{
			sampleType: []string{"foo"},
		},
	} {
		conf := Sample{SampleType: tcase.sampleType}
		if actual := conf.IsStandaloneSample(); actual != tcase.expectTypeStandalone {
			t.Errorf("case %d standalone sample type: expected %v, got %v", idx, tcase.expectTypeStandalone, actual)
		}
		if actual := conf.IsDocSample(); actual != tcase.expectTypeDoc {
			t.Errorf("case %d doc sample type: expected %v, got %v", idx, tcase.expectTypeDoc, actual)
		}
	}
}
