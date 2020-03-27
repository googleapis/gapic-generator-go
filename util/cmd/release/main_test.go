// Copyright 2020 Google LLC
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
)

func Test_format(t *testing.T) {
	in := []string{
		"gapic: foo bar baz",
		"bazel: foo bar baz",
		"gencli: foo bar baz",
		"samples: foo bar baz",
		"chore(deps): foo bar baz",
		"chore: foo bar baz",
		"unknown: foo bar baz",
		"malformed message",
	}

	want := `# gapic

* foo bar baz

# bazel

* foo bar baz

# gencli

* foo bar baz

# samples

* foo bar baz

# chores

* foo bar baz
* update dependencies (see history)

# other

* unknown: foo bar baz
* malformed message`

	got := format(in)
	if got != want {
		t.Errorf("format(%q)=%q, want %q", in, got, want)
	}
}
