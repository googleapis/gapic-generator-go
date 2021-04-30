// Copyright (C) 2021  Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package gengapic

import (
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/gapic-generator-go/internal/errors"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
)

func (g *generator) genGRPCMethods(serv *descriptor.ServiceDescriptorProto, servName string) error {
	g.addMetadataServiceForTransport(serv.GetName(), "grpc", servName)

	methods := append(serv.GetMethod(), g.getMixinMethods()...)
	for _, m := range methods {
		g.methodDoc(m)
		if err := g.genMethod(servName, serv, m); err != nil {
			return errors.E(err, "method: %s", m.GetName())
		}
		g.addMetadataMethod(serv.GetName(), "grpc", m.GetName())
	}
	sort.Slice(g.aux.lros, func(i, j int) bool {
		return g.aux.lros[i].GetName() < g.aux.lros[j].GetName()
	})
	for _, m := range g.aux.lros {
		if err := g.lroType(servName, serv, m); err != nil {
			return err
		}
	}

	return nil
}
func (g *generator) unaryCall(servName string, m *descriptor.MethodDescriptorProto) error {
	s := g.descInfo.ParentElement[m]
	sFQN := fmt.Sprintf("%s.%s", g.descInfo.ParentFile[s].GetPackage(), s.GetName())
	inType := g.descInfo.Type[*m.InputType]
	outType := g.descInfo.Type[*m.OutputType]

	inSpec, err := g.descInfo.ImportSpec(inType)
	if err != nil {
		return err
	}
	outSpec, err := g.descInfo.ImportSpec(outType)
	if err != nil {
		return err
	}

	p := g.printf

	lowcaseServName := lowerFirst(servName)

	p("func (c *%sGRPCClient) %s(ctx context.Context, req *%s.%s, opts ...gax.CallOption) (*%s.%s, error) {",
		lowcaseServName, m.GetName(), inSpec.Name, inType.GetName(), outSpec.Name, outType.GetName())

	g.deadline(sFQN, m.GetName())

	err = g.insertMetadata(m)
	if err != nil {
		return err
	}
	g.appendCallOpts(m)

	p("var resp *%s.%s", outSpec.Name, outType.GetName())
	p("err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {")
	p("  var err error")
	p("  resp, err = %s", g.grpcStubCall(m))
	p("  return err")
	p("}, opts...)")
	p("if err != nil {")
	p("  return nil, err")
	p("}")
	p("return resp, nil")

	p("}")
	p("")

	g.imports[inSpec] = true
	g.imports[outSpec] = true

	return nil
}

func (g *generator) emptyUnaryCall(servName string, m *descriptor.MethodDescriptorProto) error {
	s := g.descInfo.ParentElement[m]
	sFQN := fmt.Sprintf("%s.%s", g.descInfo.ParentFile[s].GetPackage(), s.GetName())
	inType := g.descInfo.Type[*m.InputType]

	inSpec, err := g.descInfo.ImportSpec(inType)
	if err != nil {
		return err
	}

	p := g.printf

	lowcaseServName := lowerFirst(servName)

	p("func (c *%sGRPCClient) %s(ctx context.Context, req *%s.%s, opts ...gax.CallOption) error {",
		lowcaseServName, m.GetName(), inSpec.Name, inType.GetName())

	g.deadline(sFQN, m.GetName())

	err = g.insertMetadata(m)
	if err != nil {
		return err
	}
	g.appendCallOpts(m)
	p("err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {")
	p("  var err error")
	p("  _, err = %s", g.grpcStubCall(m))
	p("  return err")
	p("}, opts...)")
	p("return err")

	p("}")
	p("")

	g.imports[inSpec] = true
	return nil
}
func (g *generator) insertMetadata(m *descriptor.MethodDescriptorProto) error {
	headers, err := parseRequestHeaders(m)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		seen := map[string]bool{}
		var formats, values strings.Builder
		for _, h := range headers {
			field := h[1]
			// skip fields that have multiple patterns, they use the same accessor
			if _, dupe := seen[field]; dupe {
				continue
			}
			seen[field] = true

			accessor := fmt.Sprintf("req%s", buildAccessor(field))
			typ := g.lookupFieldType(m.GetInputType(), field)

			// TODO(noahdietz): need to handle []byte for TYPE_BYTES.
			if typ == descriptor.FieldDescriptorProto_TYPE_STRING {
				accessor = fmt.Sprintf("url.QueryEscape(%s)", accessor)
			} else if typ == descriptor.FieldDescriptorProto_TYPE_DOUBLE || typ == descriptor.FieldDescriptorProto_TYPE_FLOAT {
				// Format the floating point value with mode 'g' to allow for
				// exponent formatting when necessary, and decimal when adequate.
				// QueryEscape the resulting string in case there is a '+' in the
				// exponent.
				// See golang.org/pkg/fmt for more information on formatting.
				accessor = fmt.Sprintf(`url.QueryEscape(fmt.Sprintf("%%g", %s))`, accessor)
			}

			// URL encode key & values separately per aip.dev/4222.
			// Encode the key ahead of time to reduce clutter
			// and because it will likely never be necessary
			fmt.Fprintf(&values, " %q, %s,", url.QueryEscape(field), accessor)
			formats.WriteString("%s=%v&")
		}
		f := formats.String()[:formats.Len()-1]
		v := values.String()[:values.Len()-1]

		g.printf("md := metadata.Pairs(\"x-goog-request-params\", fmt.Sprintf(%q,%s))", f, v)
		g.printf("ctx = insertMetadata(ctx, c.xGoogMetadata, md)")

		g.imports[pbinfo.ImportSpec{Path: "fmt"}] = true
		g.imports[pbinfo.ImportSpec{Path: "net/url"}] = true

		return nil
	}

	g.printf("ctx = insertMetadata(ctx, c.xGoogMetadata)")

	return nil
}

// genMethod generates a single method from a client. m must be a method declared in serv.
// If the generated method requires an auxillary type, it is added to aux.
func (g *generator) genMethod(servName string, serv *descriptor.ServiceDescriptorProto, m *descriptor.MethodDescriptorProto) error {
	// Check if the RPC returns google.longrunning.Operation.
	if g.isLRO(m) {
		g.aux.lros = append(g.aux.lros, m)
		return g.lroCall(servName, m)
	}

	if m.GetOutputType() == emptyType {
		return g.emptyUnaryCall(servName, m)
	}

	if pf, err := g.pagingField(m); err != nil {
		return err
	} else if pf != nil {
		iter, err := g.iterTypeOf(pf)
		if err != nil {
			return err
		}

		return g.pagingCall(servName, m, pf, iter)
	}

	switch {
	case m.GetClientStreaming():
		return g.noRequestStreamCall(servName, serv, m)
	case m.GetServerStreaming():
		return g.serverStreamCall(servName, serv, m)
	default:
		return g.unaryCall(servName, m)
	}
}

func (g *generator) grpcStubCall(method *descriptor.MethodDescriptorProto) string {
	service := g.descInfo.ParentElement[method]
	stub := pbinfo.ReduceServName(service.GetName(), g.opts.pkgName)
	return fmt.Sprintf("c.%s.%s(ctx, req, settings.GRPC...)", grpcClientField(stub), method.GetName())
}
func (g *generator) grpcClientInit(serv *descriptor.ServiceDescriptorProto, servName string, imp pbinfo.ImportSpec, hasRPCForLRO bool) {
	p := g.printf

	// We DON'T want to export the transport layers.
	lowcaseServName := lowerFirst(servName)

	p("// %sGRPCClient is a client for interacting with %s over gRPC transport.", lowcaseServName, g.apiName)
	p("//")
	p("// Methods, except Close, may be called concurrently. However, fields must not be modified concurrently with method calls.")
	p("type %sGRPCClient struct {", lowcaseServName)
	p("// Connection pool of gRPC connections to the service.")
	p("connPool gtransport.ConnPool")
	p("")

	p("// flag to opt out of default deadlines via %s", disableDeadlinesVar)
	p("disableDeadlines bool")
	p("")

	p("// Points back to the CallOptions field of the containing %sClient", servName)
	p("CallOptions **%sCallOptions", servName)
	p("")

	p("// The gRPC API client.")
	p("%s %s.%sClient", grpcClientField(servName), imp.Name, serv.GetName())
	p("")

	if hasRPCForLRO {
		p("// LROClient is used internally to handle long-running operations.")
		p("// It is exposed so that its CallOptions can be modified if required.")
		p("// Users should not Close this client.")
		p("LROClient **lroauto.OperationsClient")
		p("")
		g.imports[pbinfo.ImportSpec{Name: "lroauto", Path: "cloud.google.com/go/longrunning/autogen"}] = true
	}

	if g.hasLROMixin() {
		p("operationsClient longrunningpb.OperationsClient")
		p("")

		g.imports[pbinfo.ImportSpec{Name: "longrunningpb", Path: "google.golang.org/genproto/googleapis/longrunning"}] = true
	}

	if g.hasIAMPolicyMixin() {

		p("iamPolicyClient iampb.IAMPolicyClient")
		p("")

		g.imports[pbinfo.ImportSpec{Name: "iampb", Path: "google.golang.org/genproto/googleapis/iam/v1"}] = true
	}

	if g.hasLocationMixin() {

		p("locationsClient locationpb.LocationsClient")
		p("")

		g.imports[pbinfo.ImportSpec{Name: "locationpb", Path: "google.golang.org/genproto/googleapis/cloud/location"}] = true
	}

	p("// The x-goog-* metadata to be sent with each request.")
	p("xGoogMetadata metadata.MD")

	p("}")
	p("")

	g.imports[pbinfo.ImportSpec{Path: "google.golang.org/grpc"}] = true
	g.imports[pbinfo.ImportSpec{Path: "google.golang.org/grpc/metadata"}] = true

	g.grpcClientUtilities(serv, servName, imp, hasRPCForLRO)
}

func (g *generator) grpcClientUtilities(serv *descriptor.ServiceDescriptorProto, servName string, imp pbinfo.ImportSpec, hasRPCForLRO bool) {
	p := g.printf

	clientName := camelToSnake(serv.GetName())
	clientName = strings.Replace(clientName, "_", " ", -1)
	lowcaseServName := lowerFirst(servName)

	// Factory function
	p("// New%sClient creates a new %s client based on gRPC.", servName, clientName)
	p("//")
	g.comment(g.comments[serv])
	p("func New%[1]sClient(ctx context.Context, opts ...option.ClientOption) (*%[1]sClient, error) {", servName)
	p("  clientOpts := default%[1]sGRPCClientOptions()", servName)

	p("  if new%sClientHook != nil {", servName)
	p("    hookOpts, err := new%sClientHook(ctx, clientHookParams{})", servName)
	p("    if err != nil {")
	p("      return nil, err")
	p("    }")
	p("    clientOpts = append(clientOpts, hookOpts...)")
	p("  }")
	p("")
	p("  disableDeadlines, err := checkDisableDeadlines()")
	p("  if err != nil {")
	p("    return nil, err")
	p("  }")
	p("")
	p("  connPool, err := gtransport.DialPool(ctx, append(clientOpts, opts...)...)")
	p("  if err != nil {")
	p("    return nil, err")
	p("  }")
	p("  client := %[1]sClient{CallOptions: default%[1]sCallOptions()}", servName)
	p("")
	p("  c := &%sGRPCClient{", lowcaseServName)
	p("    connPool:    connPool,")
	p("    disableDeadlines: disableDeadlines,")
	p("    %s: %s.New%sClient(connPool),", grpcClientField(servName), imp.Name, serv.GetName())
	p("    CallOptions: &client.CallOptions,")
	if g.hasLROMixin() {
		p("    operationsClient: longrunningpb.NewOperationsClient(connPool),")
	}
	if g.hasIAMPolicyMixin() {
		p("    iamPolicyClient: iampb.NewIAMPolicyClient(connPool),")
	}
	if g.hasLocationMixin() {
		p("    locationsClient: locationpb.NewLocationsClient(connPool),")
	}
	p("")
	p("  }")
	p("  c.setGoogleClientInfo()")
	p("")
	p("  client.internal%sClient = c", servName)
	p("")

	if hasRPCForLRO {
		p("  client.LROClient, err = lroauto.NewOperationsClient(ctx, gtransport.WithConnPool(connPool))")
		p("  if err != nil {")
		p("    // This error \"should not happen\", since we are just reusing old connection pool")
		p("    // and never actually need to dial.")
		p("    // If this does happen, we could leak connp. However, we cannot close conn:")
		p("    // If the user invoked the constructor with option.WithGRPCConn,")
		p("    // we would close a connection that's still in use.")
		p("    // TODO: investigate error conditions.")
		p("    return nil, err")
		p("  }")
		p("  c.LROClient = &client.LROClient")
	}

	p("  return &client, nil")
	p("}")
	p("")

	g.imports[pbinfo.ImportSpec{Name: "gtransport", Path: "google.golang.org/api/transport/grpc"}] = true
	g.imports[pbinfo.ImportSpec{Path: "context"}] = true

	// Connection method
	p("// Connection returns a connection to the API service.")
	p("//")
	p("// Deprecated.")
	p("func (c *%sGRPCClient) Connection() *grpc.ClientConn {", lowcaseServName)
	p("  return c.connPool.Conn()")
	p("}")
	p("")

	// setGoogleClientInfo method
	p("// setGoogleClientInfo sets the name and version of the application in")
	p("// the `x-goog-api-client` header passed on each request. Intended for")
	p("// use by Google-written clients.")
	p("func (c *%sGRPCClient) setGoogleClientInfo(keyval ...string) {", lowcaseServName)
	p(`  kv := append([]string{"gl-go", versionGo()}, keyval...)`)
	p(`  kv = append(kv, "gapic", versionClient, "gax", gax.Version, "grpc", grpc.Version)`)
	p(`  c.xGoogMetadata = metadata.Pairs("x-goog-api-client", gax.XGoogHeader(kv...))`)
	p("}")
	p("")

	// Close method
	p("// Close closes the connection to the API service. The user should invoke this when")
	p("// the client is no longer required.")
	p("func (c *%sGRPCClient) Close() error {", lowcaseServName)
	p("  return c.connPool.Close()")
	p("}")
	p("")
}
