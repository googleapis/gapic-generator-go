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

package gencli

import (
	"strings"

	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
)

const (
	// ShortDescMax is the maximum length accepted for
	// the Short usage docs
	ShortDescMax = 50

	// LongDescMax is the maximum length accepted for
	// the Long usage docs
	LongDescMax = 150
)

func toLongUsage(cmt string) string {
	return shorten(cmt, LongDescMax)
}

func toShortUsage(cmt string) string {
	return shorten(cmt, ShortDescMax)
}

func shorten(cmt string, limit int) string {
	if len(cmt) > limit {
		sep := strings.LastIndex(cmt[:limit], " ")
		if sep == -1 {
			sep = limit
		}
		cmt = cmt[:sep] + "..."
	}

	return cmt
}

func sanitizeComment(cmt string) string {
	cmt = strings.Replace(cmt, "\\", `\\`, -1)
	cmt = strings.Replace(cmt, "\n", " ", -1)
	cmt = strings.Replace(cmt, "\"", "'", -1)
	cmt = strings.TrimSpace(cmt)
	return cmt
}

func strContains(a []string, s string) bool {
	for _, as := range a {
		if as == s {
			return true
		}
	}
	return false
}

func putImport(imports map[string]*pbinfo.ImportSpec, pkg *pbinfo.ImportSpec) {
	key := pkg.Name
	if key == "" {
		key = pkg.Path
	}

	imports[key] = pkg
}

func title(name string) string {
	split := strings.Split(name, "_")
	for i, s := range split {
		split[i] = strings.Title(s)
	}

	return strings.Join(split, "")
}
