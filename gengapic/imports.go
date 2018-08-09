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
)

type importSpec struct {
	name, path string
}

// sortImports sorts the import specs,
// and returns the index of the first non-standard import.
func sortImports(a []importSpec) int {
	sort.Slice(a, func(i, j int) bool {
		iDot := strings.IndexByte(a[i].path, '.') >= 0
		jDot := strings.IndexByte(a[j].path, '.') >= 0

		// standard import (without dots) comes first
		if iDot != jDot {
			return jDot
		}

		if a[i].path != a[j].path {
			return a[i].path < a[j].path
		}
		return a[i].name < a[j].name
	})
	return sort.Search(len(a), func(i int) bool {
		return strings.IndexByte(a[i].path, '.') >= 0
	})
}
