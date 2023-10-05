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
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var updateGolden = flag.Bool("update_golden", false, "update golden files")

// Diff is a test helper, testing got against contents of goldenFile. If they do not match,
// Diff fails the test, reporting the diff.
func Diff(t *testing.T, got, goldenFile string) {
	t.Helper()

	if *updateGolden {
		if err := os.WriteFile(goldenFile, []byte(got), 0644); err != nil {
			t.Fatal(err)
		}
		return
	}

	wantBytes, err := os.ReadFile(goldenFile)
	if err != nil {
		t.Fatal(err)
	}
	want := string(wantBytes)

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch(-want, +got): %s", diff)
	}
}
