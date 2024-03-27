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

	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"google.golang.org/protobuf/types/descriptorpb"
)

const (
	googleDefaultUniverse     = "googleapis.com"
	universeDomainPlaceholder = "UNIVERSE_DOMAIN"
)

func (g *generator) clientHook(servName string) {
	p := g.printf

	p("var new%sClientHook clientHook", servName)
	p("")
}

func (g *generator) clientOptions(serv *descriptorpb.ServiceDescriptorProto, servName string) error {
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

	for _, v := range g.opts.transports {
		switch v {
		case grpc:
			if err := g.grpcClientOptions(serv, servName); err != nil {
				return err
			}
			g.grpcCallOptions(serv, servName)
		case rest:
			g.restCallOptions(serv, servName)
		default:
			return fmt.Errorf("unexpected transport variant (supported variants are %q, %q): %d",
				v, grpc, rest)
		}
	}

	return nil
}

func (g *generator) internalClientIntfInit(serv *descriptorpb.ServiceDescriptorProto, servName string) error {
	p := g.printf

	p("// internal%sClient is an interface that defines the methods available from %s.", servName, g.apiName)
	p("type internal%sClient interface {", servName)
	p("Close() error")
	p("setGoogleClientInfo(...string)")
	p("Connection() *grpc.ClientConn")
	g.imports[pbinfo.ImportSpec{Path: "google.golang.org/grpc"}] = true

	// The mixin methods are for manipulating LROs, IAM, and Location.
	methods := append(serv.GetMethod(), g.getMixinMethods()...)

	if len(methods) > 0 {
		g.imports[pbinfo.ImportSpec{Path: "context"}] = true
		g.imports[pbinfo.ImportSpec{Path: "google.golang.org/api/option"}] = true
	}
	for _, m := range methods {

		inType := g.descInfo.Type[m.GetInputType()]
		inSpec, err := g.descInfo.ImportSpec(inType)
		if err != nil {
			return err
		}
		g.imports[inSpec] = true
		if m.GetOutputType() == emptyType {
			p("%s(context.Context, *%s.%s, ...gax.CallOption) error",
				m.GetName(),
				inSpec.Name,
				inType.GetName())
			continue
		}

		if pf, _, err := g.getPagingFields(m); err != nil {
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
			lroType := lroTypeName(m)
			p("%s(context.Context, *%s.%s, ...gax.CallOption) (*%s, error)",
				m.GetName(), inSpec.Name, inType.GetName(), lroType)
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
			retTyp, err := g.returnType(m)
			if err != nil {
				return err
			}

			p("%s(context.Context, *%s.%s, ...gax.CallOption) (%s, error)",
				m.GetName(), inSpec.Name, inType.GetName(), retTyp)
		}
	}
	p("}")
	p("")

	return nil
}

// serviceDoc is a helper function similar to methodDoc that includes a deprecation notice for deprecated services.
func (g *generator) serviceDoc(serv *descriptorpb.ServiceDescriptorProto) {
	com := g.comments[serv]

	// If there's no comment and the service is not deprecated, return.
	if com == "" && !serv.GetOptions().GetDeprecated() {
		return
	}

	// If the service is marked as deprecated and there is no comment, then add default deprecation comment.
	// If the service has a comment but it does not include a deprecation notice, then append a default deprecation notice.
	// If the service includes a deprecation notice at the beginning of the comment, prepend a comment stating the service is deprecated and use the included deprecation notice.
	if serv.GetOptions().GetDeprecated() {
		if com == "" {
			com = fmt.Sprintf("\n%s is deprecated.\n\nDeprecated: %[1]s may be removed in a future version.", serv.GetName())
		} else if strings.HasPrefix(com, "Deprecated:") {
			com = fmt.Sprintf("\n%s is deprecated.\n\n%s", serv.GetName(), com)
		} else if !containsDeprecated(com) {
			com = fmt.Sprintf("%s\n\nDeprecated: %s may be removed in a future version.", com, serv.GetName())
		}
	}
	com = strings.TrimSpace(com)

	// Prepend new line break before existing service comments.
	g.printf("//")
	g.comment(com)
}

func (g *generator) clientInit(serv *descriptorpb.ServiceDescriptorProto, servName string, hasRPCForLRO bool) {
	p := g.printf

	// client struct
	p("// %sClient is a client for interacting with %s.", servName, g.apiName)
	p("// Methods, except Close, may be called concurrently. However, fields must not be modified concurrently with method calls.")
	g.serviceDoc(serv)
	p("type %sClient struct {", servName)

	// Fields
	p("// The internal transport-dependent client.")
	p("internalClient internal%sClient", servName)
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
	p("// Wrapper methods routed to the internal client.")
	p("")
	p("// Close closes the connection to the API service. The user should invoke this when")
	p("// the client is no longer required.")
	p("func (c *%sClient) Close() error {", servName)
	p("  return c.internalClient.Close()")
	p("}")
	p("")
	p("// setGoogleClientInfo sets the name and version of the application in")
	p("// the `x-goog-api-client` header passed on each request. Intended for")
	p("// use by Google-written clients.")
	p("func (c *%sClient) setGoogleClientInfo(keyval ...string) {", servName)
	p("  c.internalClient.setGoogleClientInfo(keyval...)")
	p("}")
	p("")
	p("// Connection returns a connection to the API service.")
	p("//")
	p("// Deprecated: Connections are now pooled so this method does not always")
	p("// return the same resource.")
	p("func (c *%sClient) Connection() *grpc.ClientConn {", servName)
	p("  return c.internalClient.Connection()")
	p("}")
	p("")
	methods := append(serv.GetMethod(), g.getMixinMethods()...)
	for _, m := range methods {
		g.genClientWrapperMethod(m, serv, servName)
	}
}

func (g *generator) genClientWrapperMethod(m *descriptorpb.MethodDescriptorProto, serv *descriptorpb.ServiceDescriptorProto, servName string) error {
	p := g.printf

	clientTypeName := fmt.Sprintf("%sClient", servName)
	inType := g.descInfo.Type[m.GetInputType()]
	inSpec, err := g.descInfo.ImportSpec(inType)
	if err != nil {
		return err
	}

	// Generate method documentation just before any method is generated.
	g.methodDoc(m, serv)

	if m.GetOutputType() == emptyType {
		reqTyp := fmt.Sprintf("%s.%s", inSpec.Name, inType.GetName())
		p("func (c *%s) %s(ctx context.Context, req *%s, opts ...gax.CallOption) error {",
			clientTypeName, m.GetName(), reqTyp)
		p("    return c.internalClient.%s(ctx, req, opts...)", m.GetName())
		p("}")
		p("")

		g.addSnippetsMetadataParams(m, serv.GetName(), reqTyp)
		return nil
	}

	if g.isLRO(m) {
		reqTyp := fmt.Sprintf("%s.%s", inSpec.Name, inType.GetName())
		lroType := lroTypeName(m)
		p("func (c *%s) %s(ctx context.Context, req *%s, opts ...gax.CallOption) (*%s, error) {",
			clientTypeName, m.GetName(), reqTyp, lroType)
		p("    return c.internalClient.%s(ctx, req, opts...)", m.GetName())
		p("}")
		p("")
		p("// %s returns a new %[1]s from a given name.", lroType)
		p("// The name must be that of a previously created %s, possibly from a different process.", lroType)
		p("func (c *%s) %s(name string) *%[2]s {", clientTypeName, lroType)
		p("  return c.internalClient.%s(name)", lroType)
		p("}")
		p("")

		g.addSnippetsMetadataParams(m, serv.GetName(), reqTyp)
		g.addSnippetsMetadataResult(m, serv.GetName(), lroType)
		return nil
	}

	if pf, _, err := g.getPagingFields(m); err != nil {
		return err
	} else if pf != nil {
		reqTyp := fmt.Sprintf("%s.%s", inSpec.Name, inType.GetName())
		iter, err := g.iterTypeOf(pf)
		if err != nil {
			return err
		}
		p("func (c *%s) %s(ctx context.Context, req *%s, opts ...gax.CallOption) *%s {",
			clientTypeName, m.GetName(), reqTyp, iter.iterTypeName)
		p("    return c.internalClient.%s(ctx, req, opts...)", m.GetName())
		p("}")
		p("")

		g.addSnippetsMetadataParams(m, serv.GetName(), reqTyp)
		g.addSnippetsMetadataResult(m, serv.GetName(), iter.iterTypeName)
		return nil
	}

	switch {
	case m.GetClientStreaming():
		servSpec, err := g.descInfo.ImportSpec(serv)
		if err != nil {
			return err
		}

		retTyp := fmt.Sprintf("%s.%s_%sClient", servSpec.Name, serv.GetName(), m.GetName())
		p("func (c *%s) %s(ctx context.Context, opts ...gax.CallOption) (%s, error) {",
			clientTypeName, m.GetName(), retTyp)
		p("    return c.internalClient.%s(ctx, opts...)", m.GetName())
		p("}")
		p("")

		g.addSnippetsMetadataParams(m, serv.GetName(), "")
		g.addSnippetsMetadataResult(m, serv.GetName(), retTyp)
		return nil
	case m.GetServerStreaming():
		servSpec, err := g.descInfo.ImportSpec(serv)
		if err != nil {
			return err
		}

		reqTyp := fmt.Sprintf("%s.%s", inSpec.Name, inType.GetName())
		retTyp := fmt.Sprintf("%s.%s_%sClient", servSpec.Name, serv.GetName(), m.GetName())
		p("func (c *%s) %s(ctx context.Context, req *%s, opts ...gax.CallOption) (%s, error) {",
			clientTypeName, m.GetName(), reqTyp, retTyp)
		p("    return c.internalClient.%s(ctx, req, opts...)", m.GetName())
		p("}")
		p("")

		g.addSnippetsMetadataParams(m, serv.GetName(), reqTyp)
		g.addSnippetsMetadataResult(m, serv.GetName(), retTyp)
		return nil
	default:
		reqTyp := fmt.Sprintf("%s.%s", inSpec.Name, inType.GetName())
		retTyp, err := g.returnType(m)
		if err != nil {
			return err
		}

		p("func (c *%s) %s(ctx context.Context, req *%s, opts ...gax.CallOption) (%s, error) {",
			clientTypeName, m.GetName(), reqTyp, retTyp)
		p("    return c.internalClient.%s(ctx, req, opts...)", m.GetName())
		p("}")
		p("")

		g.addSnippetsMetadataParams(m, serv.GetName(), reqTyp)
		g.addSnippetsMetadataResult(m, serv.GetName(), retTyp)
		return nil
	}

}

func (g *generator) makeClients(serv *descriptorpb.ServiceDescriptorProto, servName string) error {
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
	g.clientInit(serv, servName, hasLRO)

	for _, v := range g.opts.transports {
		switch v {
		case grpc:
			g.grpcClientInit(serv, servName, imp, hasLRO)
		case rest:
			g.restClientInit(serv, servName, hasLRO)
		default:
			return fmt.Errorf("unexpected transport variant (supported variants are %q, %q): %d",
				v, grpc, rest)
		}
	}

	return nil
}

// generateDefaultEndpointTemplate returns the default endpoint with the
// Google Default Universe (googleapis.com) replaced with the placeholder
// UNIVERSE_DOMAIN for universe domain substitution.
//
// We need to apply the following type of transformation:
// 1. pubsub.googleapis.com to pubsub.UNIVERSE_DOMAIN
// 2. pubsub.sandbox.googleapis.com to pubsub.sandbox.UNIVERSE_DOMAIN
//
// This function is needed because the default endpoint template is currently
// not part of the service proto. In the future, we should update the
// service proto to include a new "google.api.default_endpoint_template" option.
func generateDefaultEndpointTemplate(defaultEndpoint string) string {
	return strings.Replace(defaultEndpoint, googleDefaultUniverse, universeDomainPlaceholder, 1)
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
