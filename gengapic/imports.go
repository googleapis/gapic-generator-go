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

package gengapic

import (
	"sort"
	"strings"

	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
)

// sortImports sorts the import specs,
// and returns the index of the first non-standard import.
func sortImports(a []pbinfo.ImportSpec) int {
	sort.Slice(a, func(i, j int) bool {
		iDot := strings.IndexByte(a[i].Path, '.') >= 0
		jDot := strings.IndexByte(a[j].Path, '.') >= 0

		// standard import (without dots) comes first
		if iDot != jDot {
			return jDot
		}

		if a[i].Path != a[j].Path {
			return a[i].Path < a[j].Path
		}
		return a[i].Name < a[j].Name
	})
	return sort.Search(len(a), func(i int) bool {
		return strings.IndexByte(a[i].Path, '.') >= 0
	})
}
