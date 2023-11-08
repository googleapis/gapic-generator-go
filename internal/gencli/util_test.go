// Copyright 2023 Google LLC
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

package gencli

import "testing"

func TestTitle(t *testing.T) {
	testCases := []struct {
		name, input, want string
	}{
		{
			name:  "simple",
			input: "title",
			want:  "Title",
		},
		{
			name:  "underscores",
			input: "title_with_underscores",
			want:  "TitleWithUnderscores",
		},
		{
			name:  "dots",
			input: "title.with.dots",
			want:  "Title.With.Dots",
		},
		{
			name:  "dots and underscores",
			input: "title.with_dots.and_underscores",
			want:  "Title.WithDots.AndUnderscores",
		},
	}
	for _, tst := range testCases {
		t.Run(tst.name, func(t *testing.T) {
			if got := title(tst.input); got != tst.want {
				t.Errorf("title(%s): got %s, want %s", tst.input, got, tst.want)
			}
		})
	}
}
