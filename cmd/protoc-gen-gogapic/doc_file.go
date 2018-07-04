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

import "time"

// genDocFile generates doc.go
//
// Since it's the only file that needs to write package documentation and canonical import,
// it does not use g.commit().
func (g *generator) genDocFile(pkgPath, pkgName string) {
	p := g.printf

	p(apacheLicense, time.Now().Year())
	p("")

	p("// Package %s is an auto-generated package for the", pkgName)
	p("// Foo API.") // TODO(pongad): read service name and description
	p("package %s // import %q", pkgName, pkgPath)
	p("")

	p("import (")
	p("%s%q", "\t", "golang.org/x/net/context")
	p("%s%q", "\t", "google.golang.org/grpc/metadata")
	p(")")
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
	// TODO(pongad): read scopes
	p("  }")
	p("}")
}
