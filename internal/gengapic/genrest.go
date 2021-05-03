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
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/gapic-generator-go/internal/errors"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
)

func lowcaseRestClientName(servName string) string {
	if servName == "" {
		return "restClient"
	}
	return lowerFirst(servName + "RESTClient")
}

func (g *generator) restClientInit(serv *descriptor.ServiceDescriptorProto, servName string, imp pbinfo.ImportSpec, hasRPCForLRO bool) {
	p := g.printf
	lowcaseServName := lowcaseRestClientName(servName)

	p("// Methods, except Close, may be called concurrently. However, fields must not be modified concurrently with method calls.")
	p("type %s struct {", lowcaseServName)
	p("  host string")
	p("}")
	p("")
	g.restClientUtilities(serv, servName, imp, hasRPCForLRO)
	g.genRESTMethods(serv, servName)
}

func (g *generator) genRESTMethods(serv *descriptor.ServiceDescriptorProto, servName string) error {
	g.addMetadataServiceForTransport(serv.GetName(), "rest", servName)

	methods := append(serv.GetMethod(), g.getMixinMethods()...)

	for _, m := range methods {
		g.methodDoc(m)
		if err := g.genRESTMethod(servName, serv, m); err != nil {
			return errors.E(err, "method: %s", m.GetName())
		}
		g.addMetadataMethod(serv.GetName(), "rest", m.GetName())
	}

	return nil
}

func (g *generator) restClientUtilities(serv *descriptor.ServiceDescriptorProto, servName string, imp pbinfo.ImportSpec, hasRPCForLRO bool) {
	// p := g.printf

	// clientName := camelToSnake(serv.GetName())
	// clientName = strings.Replace(clientName, "_", " ", -1)

	// NOTE(dovs): enable this function only when we can verify that it works.
	// If we leave it commented out, we can change it, and we don't have to
	// guarantee correctness or stability.
	//
	// p("// New%sRESTClient creates a new %s rest client.", servName, clientName)
	// p("//")
	// g.comment(g.comments[serv])
	// p("func New%[1]sRESTClient(ctx context.Context, opts ...option.ClientOption) (*%[1]sClient, error) {", servName)
	// p("    c := &%s{", servName)
	// p("    }")
	// p("")
	// p("    return &%sClient{internal%[1]sClient: c, CallOptions: default%[1]sCallOptions()}, nil", servName)
	// p("}")

}

// genRESTMethod generates a single method from a client. m must be a method declared in serv.
// If the generated method requires an auxiliary type, it is added to aux.
func (g *generator) genRESTMethod(servName string, serv *descriptor.ServiceDescriptorProto, m *descriptor.MethodDescriptorProto) error {
	if g.isLRO(m) {
		// TODO(dovs
		g.aux.lros[m] = true
		return nil
	}

	if m.GetOutputType() == emptyType {
		// TODO(dovs)
		return nil
	}

	if pf, err := g.pagingField(m); err != nil {
		return err
	} else if pf != nil {
		// TODO(dovs)
		return nil
	}

	switch {
	case m.GetClientStreaming():
		// TODO(dovs)
		return nil
	case m.GetServerStreaming():
		// TODO(dovs)
		return nil
	default:
		// TODO(dovs)
		return g.unaryRESTCall(servName, m)
	}
}

func (g *generator) unaryRESTCall(servName string, m *descriptor.MethodDescriptorProto) error {
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
	lowcaseServName := lowcaseRestClientName(servName)
	p("func (c *%s) %s(ctx context.Context, req *%s.%s, opts ...gax.CallOption) (*%s.%s, error) {",
		lowcaseServName, m.GetName(), inSpec.Name, inType.GetName(), outSpec.Name, outType.GetName())

	// TODO(dovs): handle cancellation, metadata, osv.
	// TODO(dovs): handle request from http body opt
	p("m := jsonpb.Marshaler{}") // TODO(dovs): add appropriate options
	p("if jsonReq, err := m.MarshalToString(req); err != nil {")
	p("    return nil, err")
	p("")
	p("}")
	// TODO(dovs): add real method logic.
	// _ := http_opt(m)

	// * JSONify the request
	// * determine the url
	// ** Handle query parameters
	// * send the request, receive the response, maybe return an error
	// * deserialize the result and return

	p("return nil, nil")
	p("")
	p("}")

	return nil
}
