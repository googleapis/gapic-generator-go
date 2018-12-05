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
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
)

func (g *generator) clientUtils() {
	p := g.printf

	// versionGo
	{
		p(`const versionClient = "UNKNOWN"`)
		p("")
		p("// versionGo returns the Go runtime version. The returned string")
		p("// has no whitespace, suitable for reporting in header.")
		p("func versionGo() string {")
		p(`  const develPrefix = "devel +"`)
		p("")
		p("  s := runtime.Version()")
		p("  if strings.HasPrefix(s, develPrefix) {")
		p("    s = s[len(develPrefix):]")
		p("    if p := strings.IndexFunc(s, unicode.IsSpace); p >= 0 {")
		p("      s = s[:p]")
		p("    }")
		p("    return s")
		p("  }")
		p("")
		p("  notSemverRune := func(r rune) bool {")
		p(`    return strings.IndexRune("0123456789.", r) < 0`)
		p("  }")
		p("")
		p(`  if strings.HasPrefix(s, "go1") {`)
		p("    s = s[2:]")
		p("    var prerelease string")
		p("    if p := strings.IndexFunc(s, notSemverRune); p >= 0 {")
		p("      s, prerelease = s[:p], s[p:]")
		p("    }")
		p(`    if strings.HasSuffix(s, ".") {`)
		p(`      s += "0"`)
		p(`    } else if strings.Count(s, ".") < 2 {`)
		p(`      s += ".0"`)
		p("    }")
		p(`    if prerelease != "" {`)
		p(`      s += "-" + prerelease`)
		p("    }")
		p("    return s")
		p("  }")
		p(`  return "UNKNOWN"`)
		p("}")
		p("")

		g.imports[pbinfo.ImportSpec{Path: "runtime"}] = true
		g.imports[pbinfo.ImportSpec{Path: "strings"}] = true
		g.imports[pbinfo.ImportSpec{Path: "unicode"}] = true
	}
}
