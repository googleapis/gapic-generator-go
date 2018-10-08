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

// Package txtdiff provides text-related test helpers.
package txtdiff

import (
	"flag"
	"io/ioutil"
	"strings"
	"testing"
)

var updateGolden = flag.Bool("update_golden", false, "update golden files")

func Diff(t *testing.T, name, got, goldenFile string) {
	t.Helper()

	if *updateGolden {
		if err := ioutil.WriteFile(goldenFile, []byte(got), 0644); err != nil {
			t.Fatal(err)
		}
		return
	}

	wantBytes, err := ioutil.ReadFile(goldenFile)
	if err != nil {
		t.Fatal(err)
	}
	want := string(wantBytes)

	if got == want {
		return
	}

	gotLines := strings.Split(got, "\n")
	wantLines := strings.Split(want, "\n")
	t.Errorf("%s: (+got,-want)\n%s", name, lcsDiff(wantLines, gotLines))
}

func lcsDiff(aLines, bLines []string) string {
	// Algorithm is described by https://en.wikipedia.org/wiki/Longest_common_subsequence_problem.

	// We require O(n^2) space to memoize LCS. This is not great, however
	// imagine we have 10,000-line baseline; 1e4^2 = 1e8 ints ~= 1e9 bytes = 1GB.
	// Most development computers have more memory than this and
	// our baselines are orders of magnitude smaller; we should we fine.

	// The article uses 1-based index and use index 0 to refer to the conceptual empty element.
	// Instead of dancing around the index, we just create the empty element.
	aLines = append([]string{""}, aLines...)
	bLines = append([]string{""}, bLines...)

	c := make([][]int, len(aLines))
	for i := range c {
		c[i] = make([]int, len(bLines))
	}
	for i := 1; i < len(aLines); i++ {
		for j := 1; j < len(bLines); j++ {
			if aLines[i] == bLines[j] {
				c[i][j] = c[i-1][j-1] + 1
			} else if c[i][j-1] < c[i-1][j] {
				c[i][j] = c[i-1][j]
			} else {
				c[i][j] = c[i][j-1]
			}
		}
	}

	// The article uses recursion. I think iteration is more clear.
	var diff []string
	var sign []byte

	i := len(aLines) - 1
	j := len(bLines) - 1
	for i > 0 || j > 0 {
		if i > 0 && j > 0 && aLines[i] == bLines[j] {
			diff = append(diff, aLines[i])
			sign = append(sign, ' ')
			i--
			j--
		} else if j > 0 && (i == 0 || c[i][j-1] >= c[i-1][j]) {
			diff = append(diff, bLines[j])
			sign = append(sign, '+')
			j--
		} else if i > 0 && (j == 0 || c[i][j-1] < c[i-1][j]) {
			diff = append(diff, aLines[i])
			sign = append(sign, '-')
			i--
		}
	}

	var sb strings.Builder
	for i := len(diff) - 1; i >= 0; i-- {
		sb.WriteByte(sign[i])
		sb.WriteByte(' ')
		sb.WriteString(diff[i])
		sb.WriteByte('\n')
	}
	return sb.String()
}
