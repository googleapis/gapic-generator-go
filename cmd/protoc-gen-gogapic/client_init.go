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

import "github.com/golang/protobuf/protoc-gen-go/descriptor"

func (g *generator) clientInit(serv *descriptor.ServiceDescriptorProto, servName string) {
	p := g.printf

	var hasLRO bool
	for _, m := range serv.Method {
		if isLRO(m) {
			hasLRO = true
			break
		}
	}

	// CallOptions struct
	{
		var maxNameLen int
		for _, m := range serv.Method {
			if l := len(*m.Name); maxNameLen < l {
				maxNameLen = l
			}
		}

		p("// %[1]sCallOptions contains the retry settings for each method of %[1]sClient.", servName)
		p("type %sCallOptions struct {", servName)
		for _, m := range serv.Method {
			p("%s%s[]gax.CallOption", *m.Name, spaces(maxNameLen-len(*m.Name)+1))
		}
		p("}")
		p("")

		g.imports[importSpec{"gax", "github.com/googleapis/gax-go"}] = true
	}

	// defaultClientOptions
	{
		// TODO(pongad): read URL from somewhere
		p("func default%sClientOptions() []option.ClientOption {", servName)
		p("  return []option.ClientOption{")
		p(`    option.WithEndpoint("foo.googleapis.com:443"),`)
		p("    option.WithScopes(DefaultAuthScopes()...),")
		p("  }")
		p("}")
		p("")

		g.imports[importSpec{path: "google.golang.org/api/option"}] = true
	}

	// defaultCallOptions
	{
		// TODO(pongad): read retry params from somewhere
		p("func default%[1]sCallOptions() *%[1]sCallOptions {", servName)
		p("  return &%sCallOptions{", servName)
		p("  }")
		p("}")
		p("")
	}

	// client struct
	{
		// TODO(pongad): read "human" API name from somewhere
		p("// %sClient is a client for interacting with Foo API.", servName)
		p("//")
		p("// Methods, except Close, may be called concurrently. However, fields must not be modified concurrently with method calls.")
		p("type %sClient struct {", servName)

		p("// The connection to the service.")
		p("conn *grpc.ClientConn")
		p("")

		p("// The gRPC API client.")
		p("%sClient %s.%sClient", lowerFirst(servName), g.pkgName(serv), servName)
		p("")

		if hasLRO {
			p("// LROClient is used internally to handle longrunning operations.")
			p("// It is exposed so that its CallOptions can be modified if required.")
			p("// Users should not Close this client.")
			p("LROClient *lroauto.OperationsClient")
			p("")

			g.imports[importSpec{name: "lroauto", path: "cloud.google.com/go/longrunning/autogen"}] = true
		}

		p("// The call options for this service.")
		p("CallOptions *%sCallOptions", servName)
		p("")

		p("// The x-goog-* metadata to be sent with each request.")
		p("xGoogMetadata metadata.MD")
		p("}")
		p("")

		g.imports[importSpec{path: "google.golang.org/grpc"}] = true
		g.imports[importSpec{path: "google.golang.org/grpc/metadata"}] = true
	}

	// Client constructor
	{
		// TODO(pongad): client name
		p("// New%sClient creates a new foo client.", servName)
		p("//")
		g.comment(g.comments[serv])
		p("func New%[1]sClient(ctx context.Context, opts ...option.ClientOption) (*%[1]sClient, error) {", servName)
		p("  conn, err := transport.DialGRPC(ctx, append(default%sClientOptions(), opts...)...)", servName)
		p("  if err != nil {")
		p("    return nil, err")
		p("  }")
		p("  c := &%sClient{", servName)
		p("    conn:        conn,")
		p("    CallOptions: default%sCallOptions(),", servName)
		p("")
		p("    %sClient: %s.New%sClient(conn),", lowerFirst(servName), g.pkgName(serv), servName)
		p("  }")
		p("  c.setGoogleClientInfo()")
		p("")

		if hasLRO {
			p("  c.LROClient, err = lroauto.NewOperationsClient(ctx, option.WithGRPCConn(conn))")
			p("  if err != nil {")
			p("    // This error \"should not happen\", since we are just reusing old connection")
			p("    // and never actually need to dial.")
			p("    // If this does happen, we could leak conn. However, we cannot close conn:")
			p("    // If the user invoked the function with option.WithGRPCConn,")
			p("    // we would close a connection that's still in use.")
			p("    // TODO(pongad): investigate error conditions.")
			p("    return nil, err")
			p("  }")
		}

		p("  return c, nil")
		p("}")
		p("")

		g.imports[importSpec{path: "google.golang.org/api/transport"}] = true
		g.imports[importSpec{path: "golang.org/x/net/context"}] = true
	}

	// Connection()
	{
		p("// Connection returns the client's connection to the API service.")
		p("func (c *%sClient) Connection() *grpc.ClientConn {", servName)
		p("  return c.conn")
		p("}")
		p("")
	}

	// Close()
	{
		p("// Close closes the connection to the API service. The user should invoke this when")
		p("// the client is no longer required.")
		p("func (c *%sClient) Close() error {", servName)
		p("  return c.conn.Close()")
		p("}")
		p("")
	}

	// setGoogleClientInfo
	{
		p("// setGoogleClientInfo sets the name and version of the application in")
		p("// the `x-goog-api-client` header passed on each request. Intended for")
		p("// use by Google-written clients.")
		p("func (c *%sClient) setGoogleClientInfo(keyval ...string) {", servName)
		p(`  kv := append([]string{"gl-go", version.Go()}, keyval...)`)
		p(`  kv = append(kv, "gapic", version.Repo, "gax", gax.Version, "grpc", grpc.Version)`)
		p(`  c.xGoogMetadata = metadata.Pairs("x-goog-api-client", gax.XGoogHeader(kv...))`)
		p("}")
		p("")

		g.imports[importSpec{path: "cloud.google.com/go/internal/version"}] = true
	}
}
