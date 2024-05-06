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

	"github.com/googleapis/gapic-generator-go/internal/license"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"github.com/googleapis/gapic-generator-go/internal/printer"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// genDocFile generates doc.go
//
// Since it's the only file that needs to write package documentation and canonical import,
// it does not use g.commit().
func (g *generator) genDocFile(year int, scopes []string, serv *descriptorpb.ServiceDescriptorProto) {
	p := g.printf

	p(license.Apache, year)
	p("")

	if g.apiName != "" {
		p("// Package %s is an auto-generated package for the ", g.opts.pkgName)
		p("// %s.", g.apiName)
	}

	if g.serviceConfig != nil && g.serviceConfig.GetDocumentation() != nil {
		summary := g.serviceConfig.GetDocumentation().GetSummary()
		summary = mdPlain(summary)
		wrapped := wrapString(summary, 75)

		if len(wrapped) > 0 && g.apiName != "" {
			p("//")
		}

		for _, line := range wrapped {
			p("// %s", strings.TrimSpace(line))
		}
	}

	switch g.opts.relLvl {
	case alpha:
		p("//")
		p("//   NOTE: This package is in alpha. It is not stable, and is likely to change.")
	case beta:
		p("//")
		p("//   NOTE: This package is in beta. It is not stable, and may be subject to changes.")
	case deprecated:
		p("//")
		p("// Deprecated: Find the newer version of this package in the module.")
	}

	p("//")
	p("// General documentation")
	p("//")
	p("// For information that is relevant for all client libraries please reference")
	p("// https://pkg.go.dev/cloud.google.com/go#pkg-overview. Some information on this")
	p("// page includes:")
	p("//")
	p("//  - [Authentication and Authorization]")
	p("//  - [Timeouts and Cancellation]")
	p("//  - [Testing against Client Libraries]")
	p("//  - [Debugging Client Libraries]")
	p("//  - [Inspecting errors]")
	p("//")
	p("// Example usage")
	p("//")
	p("// To get started with this package, create a client.")
	// Code block for client creation
	servName := pbinfo.ReduceServName(serv.GetName(), g.opts.pkgName)
	tmpClient := g.pt
	g.pt = printer.P{}
	g.exampleInitClient(g.opts.pkgName, servName)
	snipClient := g.pt.String()
	g.pt = tmpClient
	g.codesnippet(snipClient)
	p("// The client will use your default application credentials. Clients should be reused instead of created as needed.")
	p("// The methods of Client are safe for concurrent use by multiple goroutines.")
	p("// The returned client must be Closed when it is done being used.")
	p("//")
	// If the service does not have any methods, then do not generate sample method snippet.
	if len(serv.GetMethod()) > 0 {
		p("// Using the Client")
		p("//")
		p("// The following is an example of making an API call with the newly created client.")
		p("//")
		// Code block for client using the first method of the service
		tmpMethod := g.pt
		g.pt = printer.P{}
		g.exampleMethodBody(g.opts.pkgName, servName, serv.GetMethod()[0])
		snipMethod := g.pt.String()
		g.pt = tmpMethod
		g.codesnippet(snipMethod)
	}

	p("// Use of Context")
	p("//")
	p("// The ctx passed to New%sClient is used for authentication requests and", servName)
	p("// for creating the underlying connection, but is not used for subsequent calls.")
	p("// Individual methods on the client use the ctx given to them.")
	p("//")
	p("// To close the open connection, use the Close() method.")
	p("//")
	p("// [Authentication and Authorization]: https://pkg.go.dev/cloud.google.com/go#hdr-Authentication_and_Authorization")
	p("// [Timeouts and Cancellation]: https://pkg.go.dev/cloud.google.com/go#hdr-Timeouts_and_Cancellation")
	p("// [Testing against Client Libraries]: https://pkg.go.dev/cloud.google.com/go#hdr-Testing")
	p("// [Debugging Client Libraries]: https://pkg.go.dev/cloud.google.com/go#hdr-Debugging")
	p("// [Inspecting errors]: https://pkg.go.dev/cloud.google.com/go#hdr-Inspecting_errors")
	p("package %s // import %q", g.opts.pkgName, g.opts.pkgPath)
	p("")

	p("import (")
	p("%s%q", "\t", "context")
	p("")
	p("%s%q", "\t", "google.golang.org/api/option")
	p(")")
	p("")

	p("// For more information on implementing a client constructor hook, see")
	p("// https://github.com/googleapis/google-cloud-go/wiki/Customizing-constructors.")
	p("type clientHookParams struct{}")
	p("type clientHook func(context.Context, clientHookParams) ([]option.ClientOption, error)")
	p("")

	p("var versionClient string")
	p("")
	p("func getVersionClient() string {")
	p(`  if versionClient == "" {`)
	p(`    return "UNKNOWN"`)
	p("  }")
	p("  return versionClient")
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
}

func collectScopes(servs []*descriptorpb.ServiceDescriptorProto) []string {
	scopeSet := map[string]bool{}
	for _, s := range servs {
		eOauthScopes := proto.GetExtension(s.Options, annotations.E_OauthScopes)
		scopes := strings.Split(eOauthScopes.(string), ",")
		for _, sc := range scopes {
			scopeSet[sc] = true
		}
	}

	var scopes []string
	for sc := range scopeSet {
		scopes = append(scopes, sc)
	}
	sort.Strings(scopes)
	return scopes
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
