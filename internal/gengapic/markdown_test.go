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

package gengapic

import (
	"testing"
)

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
		{
			in:   "<b>html</b> <a href=\"/link/to/some#thing\">value</a> <br> test",
			want: "html value (at /link/to/some#thing) \n test",
		},
		{
			in:   "not <actually: html, just some> docs",
			want: "not <actually: html, just some> docs",
		},
		{
			// basic list
			in:   "List:\n- item1\n- item2",
			want: "List:\n\n  item1\n\n  item2",
		},
		{
			// list with nested list
			in:   "List:\n* item1\n  * item2",
			want: "List:\n\n  item1\n\n    item2",
		},
		{
			// list with nested list, inline code and following text
			in:   "List:\n* item1\nabc\n  * item2\n`def`, ghi\n\ndone",
			want: "List:\n\n  item1\n  abc\n\n    item2\n    def, ghi\n\ndone",
		},
		{
			// ignore headings
			in:   "## Heading",
			want: "Heading",
		},
		{
			in:   "html <a href=\"/link/to/some#thing\">\n with a softbreak</a> <br> test",
			want: "html with a softbreak (at /link/to/some#thing) \n test",
		},
		{
			in:   "html <a \n  href=\"/link/to/some#thing\">\n with a linebreak in tag</a> <br> test",
			want: "html with a linebreak in tag (at /link/to/some#thing) \n test",
		},
		{
			in:   "link to [a search engine](https://www.google.com) with request type [Search][foo.bar_baz.v1.Search], [biz][][buz][baz].",
			want: "link to a search engine (at https://www.google.com) with request type Search, bizbuz.",
		},
	} {
		got := mdPlain(tst.in)
		if got != tst.want {
			t.Errorf("MDPlain(%q)=%q, want %q", tst.in, got, tst.want)
		}
	}
}
