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

// TODO(pongad): maybe move markdown stuff into its own package
// once we have a proper import path.

package main

import "testing"

func TestMDPlain(t *testing.T) {
	for _, tst := range []struct {
		in, want string
	}{
		{
			in:   "plain text",
			want: "plain text",
		},
		{
			in:   "code `and` plain text",
			want: "code and plain text",
		},
		{
			in:   "link to [a search engine](https://www.google.com)",
			want: "link to a search engine (at https://www.google.com)",
		},
		{
			in:   "paragraph\n\nanother paragraph",
			want: "paragraph\n\nanother paragraph",
		},
	} {
		got := MDPlain(tst.in)
		if got != tst.want {
			t.Errorf("MDPlain(%q)=%q, want %q", tst.in, got, tst.want)
		}
	}
}
