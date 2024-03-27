// Copyright 2021 Google LLC
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
	"fmt"
	"strings"

	conf "github.com/googleapis/gapic-generator-go/internal/grpc_service_config"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func lowcaseGRPCClientName(servName string) string {
	if servName == "" {
		return "gRPCClient"
	}

	return lowerFirst(servName + "GRPCClient")
}

func (g *generator) genGRPCMethods(serv *descriptorpb.ServiceDescriptorProto, servName string) error {
	g.addMetadataServiceForTransport(serv.GetName(), "grpc", servName)

	methods := append(serv.GetMethod(), g.getMixinMethods()...)
	for _, m := range methods {
		if err := g.genGRPCMethod(servName, serv, m); err != nil {
			return fmt.Errorf("error generating method %q: %v", m.GetName(), err)
		}
		g.addMetadataMethod(serv.GetName(), "grpc", m.GetName())
	}
	return nil
}

// genGRPCMethod generates a single method from a client. m must be a method declared in serv.
// If the generated method requires an auxillary type, it is added to aux.
func (g *generator) genGRPCMethod(servName string, serv *descriptorpb.ServiceDescriptorProto, m *descriptorpb.MethodDescriptorProto) error {
	// Check if the RPC returns google.longrunning.Operation.
	if g.isLRO(m) {
		if err := g.maybeAddOperationWrapper(m); err != nil {
			return err
		}
		return g.lroCall(servName, m)
	}

	if m.GetOutputType() == emptyType {
		return g.emptyUnaryGRPCCall(servName, m)
	}

	if pf, ps, err := g.getPagingFields(m); err != nil {
		return err
	} else if pf != nil {
		iter, err := g.iterTypeOf(pf)
		if err != nil {
			return err
		}

		return g.pagingCall(servName, m, pf, ps, iter)
	}

	switch {
	case m.GetClientStreaming():
		return g.noRequestStreamCall(servName, serv, m)
	case m.GetServerStreaming():
		return g.serverStreamCall(servName, serv, m)
	default:
		return g.unaryGRPCCall(servName, m)
	}
}

func (g *generator) unaryGRPCCall(servName string, m *descriptorpb.MethodDescriptorProto) error {
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

	lowcaseServName := lowcaseGRPCClientName(servName)
	retTyp := fmt.Sprintf("%s.%s", outSpec.Name, outType.GetName())
	p("func (c *%s) %s(ctx context.Context, req *%s.%s, opts ...gax.CallOption) (*%s, error) {",
		lowcaseServName, m.GetName(), inSpec.Name, inType.GetName(), retTyp)

	g.insertRequestHeaders(m, grpc)
	g.initializeAutoPopulatedFields(servName, m)
	g.appendCallOpts(m)

	p("var resp *%s", retTyp)
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

func (g *generator) emptyUnaryGRPCCall(servName string, m *descriptorpb.MethodDescriptorProto) error {
	inType := g.descInfo.Type[*m.InputType]

	inSpec, err := g.descInfo.ImportSpec(inType)
	if err != nil {
		return err
	}

	p := g.printf

	lowcaseServName := lowcaseGRPCClientName(servName)

	p("func (c *%s) %s(ctx context.Context, req *%s.%s, opts ...gax.CallOption) error {",
		lowcaseServName, m.GetName(), inSpec.Name, inType.GetName())

	g.insertRequestHeaders(m, grpc)
	g.initializeAutoPopulatedFields(servName, m)
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

func (g *generator) grpcStubCall(method *descriptorpb.MethodDescriptorProto) string {
	service := g.descInfo.ParentElement[method]
	stub := pbinfo.ReduceServName(service.GetName(), g.opts.pkgName)
	return fmt.Sprintf("c.%s.%s(ctx, req, settings.GRPC...)", grpcClientField(stub), method.GetName())
}

func (g *generator) grpcClientOptions(serv *descriptorpb.ServiceDescriptorProto, servName string) error {
	p := g.printf

	// defaultClientOptions
	eHost := proto.GetExtension(serv.Options, annotations.E_DefaultHost)
	host := eHost.(string)

	if !strings.Contains(host, ":") {
		host += ":443"
	}

	p("func default%sGRPCClientOptions() []option.ClientOption {", servName)
	p("  return []option.ClientOption{")
	p("    internaloption.WithDefaultEndpoint(%q),", host)
	p("    internaloption.WithDefaultEndpointTemplate(%q),", generateDefaultEndpointTemplate(host))
	p("    internaloption.WithDefaultMTLSEndpoint(%q),", generateDefaultMTLSEndpoint(host))
	p("    internaloption.WithDefaultUniverseDomain(%q),", googleDefaultUniverse)
	p("    internaloption.WithDefaultAudience(%q),", generateDefaultAudience(host))
	p("    internaloption.WithDefaultScopes(DefaultAuthScopes()...),")
	p("    internaloption.EnableJwtWithScope(),")
	p("    option.WithGRPCDialOption(grpc.WithDefaultCallOptions(")
	p("      grpc.MaxCallRecvMsgSize(math.MaxInt32))),")
	p("  }")
	p("}")
	p("")

	g.imports[pbinfo.ImportSpec{Path: "math"}] = true
	g.imports[pbinfo.ImportSpec{Path: "google.golang.org/api/option"}] = true
	g.imports[pbinfo.ImportSpec{Path: "google.golang.org/api/option/internaloption"}] = true

	return nil
}

func (g *generator) grpcCallOptions(serv *descriptorpb.ServiceDescriptorProto, servName string) {
	p := g.printf

	// defaultCallOptions
	c := g.grpcConf

	methods := append(serv.GetMethod(), g.getMixinMethods()...)

	// read retry params from gRPC ServiceConfig
	p("func default%[1]sCallOptions() *%[1]sCallOptions {", servName)
	p("  return &%sCallOptions{", servName)
	for _, m := range methods {
		sFQN := g.fqn(g.descInfo.ParentElement[m])
		mn := m.GetName()
		p("%s: []gax.CallOption{", mn)
		if maxReq, ok := c.RequestLimit(sFQN, mn); ok {
			p("gax.WithGRPCOptions(grpc.MaxCallSendMsgSize(%d)),", maxReq)
		}

		if maxRes, ok := c.ResponseLimit(sFQN, mn); ok {
			p("gax.WithGRPCOptions(grpc.MaxCallRecvMsgSize(%d)),", maxRes)
		}

		streaming := m.GetClientStreaming() || m.GetServerStreaming()
		if timeout, ok := c.Timeout(sFQN, mn); !streaming && ok {
			p("gax.WithTimeout(%d * time.Millisecond),", timeout)
			g.imports[pbinfo.ImportSpec{Path: "time"}] = true
		}

		if rp, ok := c.RetryPolicy(sFQN, mn); ok && rp != nil {
			p("gax.WithRetry(func() gax.Retryer {")
			p("  return gax.OnCodes([]codes.Code{")
			for _, c := range rp.GetRetryableStatusCodes() {
				cstr := c.String()

				// Go uses the American-English spelling with a single "L"
				if c == code.Code_CANCELLED {
					cstr = "Canceled"
				}

				p("    codes.%s,", snakeToCamel(cstr))
			}
			p("	 }, gax.Backoff{")
			// this ignores max_attempts
			p("		Initial:    %d * time.Millisecond,", conf.ToMillis(rp.GetInitialBackoff()))
			p("		Max:        %d * time.Millisecond,", conf.ToMillis(rp.GetMaxBackoff()))
			p("		Multiplier: %.2f,", rp.GetBackoffMultiplier())
			p("	 })")
			p("}),")

			// include imports necessary for retry configuration
			g.imports[pbinfo.ImportSpec{Path: "time"}] = true
			g.imports[pbinfo.ImportSpec{Path: "google.golang.org/grpc/codes"}] = true
		}
		p("},")
	}
	p("  }")
	p("}")
	p("")
}

func (g *generator) grpcClientInit(serv *descriptorpb.ServiceDescriptorProto, servName string, imp pbinfo.ImportSpec, hasRPCForLRO bool) {
	p := g.printf

	// We DON'T want to export the transport layers.
	lowcaseServName := lowcaseGRPCClientName(servName)

	p("// %s is a client for interacting with %s over gRPC transport.", lowcaseServName, g.apiName)
	p("//")
	p("// Methods, except Close, may be called concurrently. However, fields must not be modified concurrently with method calls.")
	p("type %s struct {", lowcaseServName)
	p("// Connection pool of gRPC connections to the service.")
	p("connPool gtransport.ConnPool")
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

	g.mixinStubs()

	p("// The x-goog-* metadata to be sent with each request.")
	p("xGoogHeaders []string")

	p("}")
	p("")

	g.imports[pbinfo.ImportSpec{Path: "google.golang.org/grpc"}] = true
	g.imports[imp] = true

	g.grpcClientUtilities(serv, servName, imp, hasRPCForLRO)
}

func (g *generator) grpcClientUtilities(serv *descriptorpb.ServiceDescriptorProto, servName string, imp pbinfo.ImportSpec, hasRPCForLRO bool) {
	p := g.printf

	clientName := camelToSnake(serv.GetName())
	clientName = strings.Replace(clientName, "_", " ", -1)
	lowcaseServName := lowcaseGRPCClientName(servName)

	// Factory function
	p("// New%sClient creates a new %s client based on gRPC.", servName, clientName)
	p("// The returned client must be Closed when it is done being used to clean up its underlying connections.")
	g.serviceDoc(serv)
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
	p("  connPool, err := gtransport.DialPool(ctx, append(clientOpts, opts...)...)")
	p("  if err != nil {")
	p("    return nil, err")
	p("  }")
	p("  client := %[1]sClient{CallOptions: default%[1]sCallOptions()}", servName)
	p("")
	p("  c := &%s{", lowcaseServName)
	p("    connPool:    connPool,")
	p("    %s: %s.New%sClient(connPool),", grpcClientField(servName), imp.Name, serv.GetName())
	p("    CallOptions: &client.CallOptions,")
	g.mixinStubsInit()
	p("")
	p("  }")
	p("  c.setGoogleClientInfo()")
	p("")
	p("  client.internalClient = c")
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
	p("// Deprecated: Connections are now pooled so this method does not always")
	p("// return the same resource.")
	p("func (c *%s) Connection() *grpc.ClientConn {", lowcaseServName)
	p("  return c.connPool.Conn()")
	p("}")
	p("")

	// setGoogleClientInfo method
	p("// setGoogleClientInfo sets the name and version of the application in")
	p("// the `x-goog-api-client` header passed on each request. Intended for")
	p("// use by Google-written clients.")
	p("func (c *%s) setGoogleClientInfo(keyval ...string) {", lowcaseServName)
	p(`  kv := append([]string{"gl-go", gax.GoVersion}, keyval...)`)
	p(`  kv = append(kv, "gapic", getVersionClient(), "gax", gax.Version, "grpc", grpc.Version)`)
	p(`  c.xGoogHeaders = []string{"x-goog-api-client", gax.XGoogHeader(kv...)}`)
	p("}")
	p("")

	// Close method
	p("// Close closes the connection to the API service. The user should invoke this when")
	p("// the client is no longer required.")
	p("func (c *%s) Close() error {", lowcaseServName)
	p("  return c.connPool.Close()")
	p("}")
	p("")
}
