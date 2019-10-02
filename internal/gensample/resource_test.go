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

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseNamePattern(t *testing.T) {
	t.Run("Error", func(t *testing.T) {
		for _, s := range []string{
			"foos/{foo",
			"foos/{foo}/{foo}",
		} {
			if _, err := parseNamePattern(s); err == nil {
				t.Errorf("parseNamePattern(%q) OK, want error", s)
			}
		}
	})

	for _, tst := range []struct {
		s      string
		pieces []string
	}{
		{"foos/{foo}", []string{"foos/", "foo"}},
		{"foos/{foo}/bars/{bar}", []string{"foos/", "foo", "/bars/", "bar"}},
	} {
		np, err := parseNamePattern(tst.s)
		if err != nil {
			t.Errorf("parseNamePattern(%q) errors %q, want %q", tst.s, err, tst.pieces)
			continue
		}
		if diff := cmp.Diff(tst.pieces, np.pieces); diff != "" {
			t.Errorf("parseNamePattern(%q).pieces: (got=-, want=+):\n%s", tst.s, diff)
			continue
		}
		for i := 1; i < len(np.pieces); i += 2 {
			p := np.pieces[i]
			if got := np.pos[p]; got != i {
				t.Errorf("parseNamePattern(%q).pos[%q] = %d, want %d", tst.s, p, got, i)
			}
		}
	}
}
