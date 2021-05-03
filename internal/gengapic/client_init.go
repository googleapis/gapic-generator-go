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

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
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

	for _, v := range g.opts.transports {
		switch v {
		case grpc:
			if err := g.grpcClientOptions(serv, servName); err != nil {
				return err
			}
			g.grpcCallOptions(serv, servName)
		case rest:
			// TODO(dovs)
			continue
		default:
			return fmt.Errorf("unexpected transport variant (supported variants are %q, %q): %d",
				v, grpc, rest)
		}
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
