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
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/gapic-generator-go/internal/errors"
	"google.golang.org/genproto/googleapis/api/annotations"
)

func (g *generator) clientOptions(serv *descriptor.ServiceDescriptorProto, servName string) error {
	p := g.printf

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
		eHost, err := proto.GetExtension(serv.Options, annotations.E_DefaultHost)
		if err != nil {
			return errors.E(err, "cannot read default host")
		}

		p("func default%sClientOptions() []option.ClientOption {", servName)
		p("  return []option.ClientOption{")
		p(`    option.WithEndpoint("%s:443"),`, *eHost.(*string))
		p("    option.WithScopes(DefaultAuthScopes()...),")
		p("  }")
		p("}")
		p("")

		g.imports[importSpec{path: "google.golang.org/api/option"}] = true
	}

	// defaultCallOptions
	{
		var retryables []string
		for _, m := range serv.GetMethod() {
			eHttp, err := proto.GetExtension(m.GetOptions(), annotations.E_Http)
			if err != nil {
				// Some methods are not annotated, this is not an error.
				continue
			}
			// Generator spec mandates we should only retry on GET, unless there is an override.
			// TODO(pongad): implement the override.
			if _, ok := eHttp.(*annotations.HttpRule).Pattern.(*annotations.HttpRule_Get); ok {
				retryables = append(retryables, m.GetName())
			}
		}

		// TODO(pongad): read retry params from somewhere
		p("func default%[1]sCallOptions() *%[1]sCallOptions {", servName)

		if len(retryables) > 0 {
			p("retry := []gax.CallOption{")
			p("  gax.WithRetry(func() gax.Retryer {")
			p("    return gax.OnCodes([]codes.Code{")
			p("      codes.Internal,")
			p("      codes.Unavailable,")
			p("    }, gax.Backoff{")
			p("      Initial: 100 * time.Millisecond,")
			p("      Max: time.Minute,")
			p("      Multiplier: 1.3,")
			p("    })")
			p("  }),")
			p("}")
			p("")

			g.imports[importSpec{path: "time"}] = true
			g.imports[importSpec{path: "google.golang.org/grpc/codes"}] = true
		}

		p("  return &%sCallOptions{", servName)
		for _, m := range retryables {
			p("%s: retry,", m)
		}
		p("  }")
		p("}")
		p("")
	}

	return nil
}

func (g *generator) clientInit(serv *descriptor.ServiceDescriptorProto, servName string) error {
	p := g.printf

	var hasLRO bool
	for _, m := range serv.Method {
		if *m.OutputType == lroType {
			hasLRO = true
			break
		}
	}

	imp, err := g.importSpec(serv)
	if err != nil {
		return err
	}

	// client struct
	{
		p("// %sClient is a client for interacting with %s API.", servName, g.apiName)
		p("//")
		p("// Methods, except Close, may be called concurrently. However, fields must not be modified concurrently with method calls.")
		p("type %sClient struct {", servName)

		p("// The connection to the service.")
		p("conn *grpc.ClientConn")
		p("")

		p("// The gRPC API client.")
		p("%s %s.%sClient", grpcClientField(servName), imp.name, serv.GetName())
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
		clientName := camelToSnake(serv.GetName())
		clientName = strings.Replace(clientName, "_", " ", -1)

		p("// New%sClient creates a new %s client.", servName, clientName)
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
		p("    %s: %s.New%sClient(conn),", grpcClientField(servName), imp.name, serv.GetName())
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
	return nil
}
