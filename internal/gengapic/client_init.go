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
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	conf "github.com/googleapis/gapic-generator-go/internal/grpc_service_config"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/genproto/googleapis/rpc/code"
)

func (g *generator) clientHook(servName string) {
	p := g.printf

	p("var new%sClientHook clientHook", servName)
	p("")
}

func (g *generator) clientOptions(serv *descriptor.ServiceDescriptorProto, servName string) error {
	p := g.printf

	// CallOptions struct
	{
		methods := append(serv.GetMethod(), g.getMixinMethods()...)

		p("// %[1]sCallOptions contains the retry settings for each method of %[1]sClient.", servName)
		p("type %sCallOptions struct {", servName)
		for _, m := range methods {
			p("%s []gax.CallOption", m.GetName())
		}
		p("}")
		p("")

		g.imports[pbinfo.ImportSpec{Name: "gax", Path: "github.com/googleapis/gax-go/v2"}] = true
	}

	// defaultClientOptions
	{
		var host string
		if eHost, err := proto.GetExtension(serv.Options, annotations.E_DefaultHost); err == nil {
			host = *eHost.(*string)
		} else {
			fqn := g.descInfo.ParentFile[serv].GetPackage() + "." + serv.GetName()
			return fmt.Errorf("service %q is missing option google.api.default_host", fqn)
		}

		if !strings.Contains(host, ":") {
			host += ":443"
		}

		p("func default%sGRPCClientOptions() []option.ClientOption {", servName)
		p("  return []option.ClientOption{")
		p("    internaloption.WithDefaultEndpoint(%q),", host)
		p("    internaloption.WithDefaultMTLSEndpoint(%q),", generateDefaultMTLSEndpoint(host))
		p("    internaloption.WithDefaultAudience(%q),", generateDefaultAudience(host))
		p("    internaloption.WithDefaultScopes(DefaultAuthScopes()...),")
		p("    option.WithGRPCDialOption(grpc.WithDisableServiceConfig()),")
		p("    option.WithGRPCDialOption(grpc.WithDefaultCallOptions(")
		p("      grpc.MaxCallRecvMsgSize(math.MaxInt32))),")
		p("  }")
		p("}")
		p("")

		g.imports[pbinfo.ImportSpec{Path: "math"}] = true
		g.imports[pbinfo.ImportSpec{Path: "google.golang.org/api/option"}] = true
		g.imports[pbinfo.ImportSpec{Path: "google.golang.org/api/option/internaloption"}] = true
	}

	// defaultCallOptions
	{
		c := g.grpcConf

		methods := append(serv.GetMethod(), g.getMixinMethods()...)

		// read retry params from gRPC ServiceConfig
		p("func default%[1]sCallOptions() *%[1]sCallOptions {", servName)
		p("  return &%sCallOptions{", servName)
		for _, m := range methods {
			sFQN := g.getServiceName(m)
			mn := m.GetName()
			p("%s: []gax.CallOption{", mn)
			if maxReq, ok := c.RequestLimit(sFQN, mn); ok {
				p("gax.WithGRPCOptions(grpc.MaxCallSendMsgSize(%d)),", maxReq)
			}

			if maxRes, ok := c.ResponseLimit(sFQN, mn); ok {
				p("gax.WithGRPCOptions(grpc.MaxCallRecvMsgSize(%d)),", maxRes)
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

	return nil
}

func (g *generator) internalClientIntfInit(serv *descriptor.ServiceDescriptorProto, servName string) error {
	p := g.printf

	p("// internal%sClient is an interface that defines the methods availaible from %s.", servName, g.apiName)
	p("type internal%sClient interface {", servName)
	p("Close() error")
	p("setGoogleClientInfo(...string)")
	p("Connection() *grpc.ClientConn")

	// The mixin methods are for manipulating LROs, IAM, and Location.
	methods := append(serv.GetMethod(), g.getMixinMethods()...)
	// methods := serv.GetMethod()

	for _, m := range methods {

		inType := g.descInfo.Type[m.GetInputType()]
		inSpec, err := g.descInfo.ImportSpec(inType)
		if err != nil {
			return err
		}
		if m.GetOutputType() == emptyType {
			p("%s(context.Context, *%s.%s, ...gax.CallOption) error",
				m.GetName(),
				inSpec.Name,
				inType.GetName())
			continue
		}

		outType := g.descInfo.Type[m.GetOutputType()]
		outSpec, err := g.descInfo.ImportSpec(outType)
		if err != nil {
			return err
		}

		if pf, err := g.pagingField(m); err != nil {
			return err
		} else if pf != nil {
			iter, err := g.iterTypeOf(pf)
			if err != nil {
				return err
			}
			p("%s(context.Context, *%s.%s, ...gax.CallOption) *%s",
				m.GetName(),
				inSpec.Name,
				inType.GetName(),
				iter.iterTypeName)
			continue
		}

		switch {
		case g.isLRO(m):
			// Unary call where the return type is a wrapper of
			// longrunning.Operation and more precise types
			lroType := lroTypeName(m.GetName())
			p("%s(context.Context, *%s.%s, ...gax.CallOption) (*%s, error)",
				m.GetName(), inSpec.Name, inType.GetName(), lroTypeName(m.GetName()))
			p("%[1]s(name string) *%[1]s", lroType)

		case m.GetClientStreaming():
			// Handles both client-streaming and bidi-streaming
			p("%s(context.Context, ...gax.CallOption) (%s.%s_%sClient, error)",
				m.GetName(), inSpec.Name, serv.GetName(), m.GetName())
		case m.GetServerStreaming():
			// Handles _just_ server streaming
			p("%s(context.Context, *%s.%s, ...gax.CallOption) (%s.%s_%sClient, error)",
				m.GetName(), inSpec.Name, inType.GetName(), inSpec.Name, serv.GetName(), m.GetName())
		default:
			// Regular, unary call
			p("%s(context.Context, *%s.%s, ...gax.CallOption) (*%s.%s, error)",
				m.GetName(), inSpec.Name, inType.GetName(), outSpec.Name, outType.GetName())
		}
	}
	p("}")
	p("")

	return nil
}

func (g *generator) clientInit(serv *descriptor.ServiceDescriptorProto, servName string, imp pbinfo.ImportSpec, hasRPCForLRO bool) {
	p := g.printf

	// client struct
	p("// %sClient is a client for interacting with %s.", servName, g.apiName)
	p("// Methods, except Close, may be called concurrently. However, fields must not be modified concurrently with method calls.")
	p("type %sClient struct {", servName)

	// Fields
	p("// The internal transport-dependent client.")
	p("internal%sClient", servName)
	p("")
	p("// The call options for this service.")
	p("CallOptions *%sCallOptions", servName)
	p("")

	// Need to keep for back compat
	if hasRPCForLRO {
		p("// LROClient is used internally to handle long-running operations.")
		p("// It is exposed so that its CallOptions can be modified if required.")
		p("// Users should not Close this client.")
		p("LROClient *lroauto.OperationsClient")
		p("")
		g.imports[pbinfo.ImportSpec{Name: "lroauto", Path: "cloud.google.com/go/longrunning/autogen"}] = true
	}

	p("}")
	p("")
}

func (g *generator) grpcClientInit(serv *descriptor.ServiceDescriptorProto, servName string, imp pbinfo.ImportSpec, hasRPCForLRO bool) {
	p := g.printf

	// We DON'T want to export the transport layers.
	lowcaseServName := lowerFirst(servName + "GRPCClient")

	p("// %s is a client for interacting with %s over gRPC transport.", lowcaseServName, g.apiName)
	p("//")
	p("// Methods, except Close, may be called concurrently. However, fields must not be modified concurrently with method calls.")
	p("type %s struct {", lowcaseServName)
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

	g.mixinStubs()

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
	lowcaseServName := lowerFirst(servName + "GRPCClient")

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
	p("  c := &%s{", lowcaseServName)
	p("    connPool:    connPool,")
	p("    disableDeadlines: disableDeadlines,")
	p("    %s: %s.New%sClient(connPool),", grpcClientField(servName), imp.Name, serv.GetName())
	p("    CallOptions: &client.CallOptions,")
	g.mixinStubsInit()
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
	p("func (c *%s) Connection() *grpc.ClientConn {", lowcaseServName)
	p("  return c.connPool.Conn()")
	p("}")
	p("")

	// setGoogleClientInfo method
	p("// setGoogleClientInfo sets the name and version of the application in")
	p("// the `x-goog-api-client` header passed on each request. Intended for")
	p("// use by Google-written clients.")
	p("func (c *%s) setGoogleClientInfo(keyval ...string) {", lowcaseServName)
	p(`  kv := append([]string{"gl-go", versionGo()}, keyval...)`)
	p(`  kv = append(kv, "gapic", versionClient, "gax", gax.Version, "grpc", grpc.Version)`)
	p(`  c.xGoogMetadata = metadata.Pairs("x-goog-api-client", gax.XGoogHeader(kv...))`)
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

func (g *generator) makeClients(serv *descriptor.ServiceDescriptorProto, servName string) error {
	var hasLRO bool
	for _, m := range serv.GetMethod() {
		if g.isLRO(m) {
			hasLRO = true
			break
		}
	}

	imp, err := g.descInfo.ImportSpec(serv)
	if err != nil {
		return err
	}

	err = g.internalClientIntfInit(serv, servName)
	if err != nil {
		return err
	}
	g.clientInit(serv, servName, imp, hasLRO)

	for _, v := range g.opts.transports {
		switch v {
		case grpc:
			g.grpcClientInit(serv, servName, imp, hasLRO)
		case rest:
			g.restClientInit(serv, servName, imp, hasLRO)
		default:
			return fmt.Errorf("unexpected transport variant (supported variants are %q, %q): %d",
				v, grpc, rest)
		}
	}

	return nil
}

// generateDefaultMTLSEndpoint attempts to derive the mTLS version of the
// defaultEndpoint via regex, and returns defaultEndpoint if unsuccessful.
//
// We need to applying the following 2 transformations:
// 1. pubsub.googleapis.com to pubsub.mtls.googleapis.com
// 2. pubsub.sandbox.googleapis.com to pubsub.mtls.sandbox.googleapis.com
//
// This function is needed because the default mTLS endpoint is currently
// not part of the service proto. In the future, we should update the
// service proto to include a new "google.api.default_mtls_host" option.
func generateDefaultMTLSEndpoint(defaultEndpoint string) string {
	var domains = []string{
		".sandbox.googleapis.com", // must come first because .googleapis.com is a substring
		".googleapis.com",
	}
	for _, domain := range domains {
		if strings.Contains(defaultEndpoint, domain) {
			return strings.Replace(defaultEndpoint, domain, ".mtls"+domain, -1)
		}
	}
	return defaultEndpoint
}

// generateDefaultAudience transforms a host into a an audience that can be used
// as the `aud` claim in a JWT.
func generateDefaultAudience(host string) string {
	aud := host
	// Add a scheme if not present.
	if !strings.Contains(aud, "://") {
		aud = "https://" + aud
	}
	// Remove port, and everything after, if present.
	if strings.Count(aud, ":") > 1 {
		firstIndex := strings.Index(aud, ":")
		secondIndex := strings.Index(aud[firstIndex+1:], ":") + firstIndex + 1
		aud = aud[:secondIndex]
	}
	// Add trailing slash if not present.
	if !strings.HasSuffix(aud, "/") {
		aud = aud + "/"
	}
	return aud
}
