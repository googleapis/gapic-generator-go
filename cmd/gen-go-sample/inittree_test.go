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

import (
	"strings"
	"testing"
)

func TestTree(t *testing.T) {
	specs := []string{
		`a.b = 1`,
		`a.c = "xyz"`,
		`a.d = 2.718281828`,
		`x = 123`,
	}

	var root initTree
	for _, s := range specs {
		if err := root.Parse(s); err != nil {
			t.Fatal(err)
		}
	}

	for _, tst := range []struct {
		path []string
		val  string
	}{
		{[]string{"a", "b"}, "1"},
		{[]string{"a", "c"}, `"xyz"`},
		{[]string{"a", "d"}, "2.718281828"},
		{[]string{"x"}, "123"},
	} {
		node := &root
		for _, p := range tst.path {
			node = node.get(p)
		}
		if node.leafVal != tst.val {
			t.Errorf("%s = %q, want %q", strings.Join(tst.path, "->"), node.leafVal, tst.val)
		}
	}

	if t.Failed() {
		t.SkipNow()
	}

	var buf strings.Builder
	buf.WriteByte('\n')
	if err := root.Print(&buf); err != nil {
		t.Error(err)
	}
	t.Log(buf.String())
}
