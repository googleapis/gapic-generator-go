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

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/gapic-generator-go/internal/errors"
	"github.com/googleapis/gapic-generator-go/internal/license"
	"google.golang.org/genproto/googleapis/api/annotations"
)

// genDocFile generates doc.go
//
// Since it's the only file that needs to write package documentation and canonical import,
// it does not use g.commit().
func (g *generator) genDocFile(pkgPath, pkgName string, year int, scopes []string) {
	p := g.printf

	p(license.Apache, year)
	p("")

	if g.apiName != "" {
		p("// Package %s is an auto-generated package for the", pkgName)
		p("// %s.", g.apiName)
	}

	// TODO(ndietz) figure out how to include this without the service config
	if g.serviceConfig != nil && g.serviceConfig.Documentation != nil {
		wrapped := wrapString(g.serviceConfig.Documentation.Summary, 75)

		if len(wrapped) > 0 && g.apiName != "" {
			p("//")
		}

		for _, line := range wrapped {
			p("// %s", strings.TrimSpace(line))
		}
	}

	switch g.relLvl {
	case alpha:
		p("//")
		p("//   NOTE: This package is in alpha. It is not stable, and is likely to change.")
	case beta:
		p("//")
		p("//   NOTE: This package is in beta. It is not stable, and may be subject to changes.")
	}

	p("//")
	p("// Use of Context")
	p("//")
	p("// The ctx passed to NewClient is used for authentication requests and")
	p("// for creating the underlying connection, but is not used for subsequent calls.")
	p("// Individual methods on the client use the ctx given to them.")
	p("//")
	p("// To close the open connection, use the Close() method.")
	p("//")
	p("// For information about setting deadlines, reusing contexts, and more")
	p("// please visit pkg.go.dev/cloud.google.com/go.")
	p("package %s // import %q", pkgName, pkgPath)
	p("")

	p("import (")
	p("%s%q", "\t", "context")
	p("%s%q", "\t", "runtime")
	p("%s%q", "\t", "strings")
	p("%s%q", "\t", "unicode")
	p("")
	p("%s%q", "\t", "google.golang.org/api/option")
	p("%s%q", "\t", "google.golang.org/grpc/metadata")
	p(")")
	p("")

	p("// For more information on implementing a client constructor hook, see")
	p("// https://github.com/googleapis/google-cloud-go/wiki/Customizing-constructors.")
	p("type clientHookParams struct{}")
	p("type clientHook func(context.Context, clientHookParams) ([]option.ClientOption, error)")
	p("")

	p("const versionClient = %q", "UNKNOWN")
	p("")

	p("func insertMetadata(ctx context.Context, mds ...metadata.MD) context.Context {")
	p("  out, _ := metadata.FromOutgoingContext(ctx)")
	p("  out = out.Copy()")
	p("  for _, md := range mds {")
	p("    for k, v := range md {")
	p("      out[k] = append(out[k], v...)")
	p("    }")
	p("  }")
	p("  return metadata.NewOutgoingContext(ctx, out)")
	p("}")
	p("")

	p("// DefaultAuthScopes reports the default set of authentication scopes to use with this package.")
	p("func DefaultAuthScopes() []string {")
	p("  return []string{")
	for _, sc := range scopes {
		p("%q,", sc)
	}
	p("  }")
	p("}")

	// versionGo
	{
		p("")
		p("// versionGo returns the Go runtime version. The returned string")
		p("// has no whitespace, suitable for reporting in header.")
		p("func versionGo() string {")
		p("  const develPrefix = %q", "devel +")
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
		p("    return !strings.ContainsRune(%q, r)", "0123456789.")
		p("  }")
		p("")
		p("  if strings.HasPrefix(s, %q) {", "go1")
		p("    s = s[2:]")
		p("    var prerelease string")
		p("    if p := strings.IndexFunc(s, notSemverRune); p >= 0 {")
		p("      s, prerelease = s[:p], s[p:]")
		p("    }")
		p("    if strings.HasSuffix(s, %q) {", ".")
		p("      s += %q", "0")
		p("    } else if strings.Count(s, %q) < 2 {", ".")
		p("      s += %q", ".0")
		p("    }")
		p("    if prerelease != %q {", "")
		p("      s += %q + prerelease", "-")
		p("    }")
		p("    return s")
		p("  }")
		p("  return %q", "UNKNOWN")
		p("}")
		p("")
	}
}

func collectScopes(servs []*descriptor.ServiceDescriptorProto, config *serviceConfig) ([]string, error) {
	scopeSet := map[string]bool{}
	for _, s := range servs {
		eOauthScopes, err := proto.GetExtension(s.Options, annotations.E_OauthScopes)
		if err == proto.ErrMissingExtension {
			continue
		}
		if err != nil {
			return nil, errors.E(err, "cannot find scopes for service: %q", s.GetName())
		}
		scopes := strings.Split(*eOauthScopes.(*string), ",")
		for _, sc := range scopes {
			scopeSet[sc] = true
		}
	}

	var scopes []string
	for sc := range scopeSet {
		scopes = append(scopes, sc)
	}
	sort.Strings(scopes)
	return scopes, nil
}

func wrapString(str string, max int) []string {
	var lines []string
	var line string

	if str == "" {
		return lines
	}

	split := strings.Fields(str)
	for _, w := range split {
		if len(line)+len(w)+1 > max {
			lines = append(lines, line)
			line = ""
		}

		line += " " + w
	}
	lines = append(lines, line)

	return lines
}
