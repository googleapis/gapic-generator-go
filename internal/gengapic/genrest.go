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
	"sort"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/gapic-generator-go/internal/errors"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"google.golang.org/genproto/googleapis/api/annotations"
)

func (g *generator) restClientInit(serv *descriptor.ServiceDescriptorProto, servName string, imp pbinfo.ImportSpec, hasRPCForLRO bool) {
	p := g.printf

	lowcaseServName := lowerFirst(servName)

	p("// %sRestClient is a client for interacting with %s over REST transport.", lowcaseServName, g.apiName)
	p("// It satisfies the internal%sClient interface.", servName)
	p("//")
	p("// Methods, except Close, may be called concurrently. However, fields must not be modified concurrently with method calls.")
	p("type %sRestClient struct {", lowcaseServName)
	p("  host string")
	p("}")
	p("")
	g.restClientUtilities(serv, servName, imp, hasRPCForLRO)
	g.genRestMethods(serv, servName)
}

func (g *generator) genRestMethods(serv *descriptor.ServiceDescriptorProto, servName string) error {
	// Clear LROs between services
	// TODO: handle LRO types shared between rest and grpc clients with a map[bool]
	g.addMetadataServiceForTransport(serv.GetName(), "rest", servName)

	methods := append(serv.GetMethod(), g.getMixinMethods()...)

	for _, m := range methods {
		g.methodDoc(m)
		if err := g.genRestMethod(servName, serv, m); err != nil {
			return errors.E(err, "method: %s", m.GetName())
		}
		g.addMetadataMethod(serv.GetName(), "rest", m.GetName())
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

func (g *generator) restClientUtilities(serv *descriptor.ServiceDescriptorProto, servName string, imp pbinfo.ImportSpec, hasRPCForLRO bool) {
	// p := g.printf

	// clientName := camelToSnake(serv.GetName())
	// clientName = strings.Replace(clientName, "_", " ", -1)

	// NOTE(dovs): enable this function only when we can verify that it works.
	// If we leave it commented out, we can change it, and we don't have to
	// guarantee correctness or stability.
	//
	// p("// New%sRestClient creates a new %s rest client.", servName, clientName)
	// p("//")
	// g.comment(g.comments[serv])
	// p("func New%[1]sRestClient(ctx context.Context, opts ...option.ClientOption) (*%[1]sRestClient, error) {", servName)
	// p("    c := &%sRestClient{", servName)
	// p("    }")
	// p("")
	// p("    return &%sClient{internal%[1]sClient: c, CallOptions: default%[1]sCallOptions()}, nil", servName)
	// p("}")

}

// genRestMethod generates a single method from a client. m must be a method declared in serv.
// If the generated method requires an auxiliary type, it is added to aux.
func (g *generator) genRestMethod(servName string, serv *descriptor.ServiceDescriptorProto, m *descriptor.MethodDescriptorProto) error {
	if g.isLRO(m) {
		// TODO(dovs)
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
		return g.unaryRestCall(servName, m)
	}
}

func doAThing(m *descriptor.MethodDescriptorProto) (error, error) {
	eHTTP, err := proto.GetExtension(m.GetOptions(), annotations.E_Http)
	if m == nil || m.GetOptions() == nil || err == proto.ErrMissingExtension {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	http := eHTTP.(*annotations.HttpRule)
	rules := []*annotations.HttpRule{http}
	rules = append(rules, http.GetAdditionalBindings()...)
	_ = rules

	return nil, nil
}

func (g *generator) unaryRestCall(servName string, m *descriptor.MethodDescriptorProto) error {
	inType := g.descInfo.Type[*m.InputType]
	outType := g.descInfo.Type[*m.OutputType]

	doAThing(m)

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
	p("func (c *%sRestClient) %s(ctx context.Context, req *%s.%s, opts ...gax.CallOption) (*%s.%s, error) {",
		lowcaseServName, m.GetName(), inSpec.Name, inType.GetName(), outSpec.Name, outType.GetName())

	// TODO(dovs): handle cancellation, metadata, osv.
	// TODO(dovs): handle request from http body opt
	p("m := jsonpb.Marshaler{}") // TODO(dovs): add appropriate options
	p("if jsonReq, err := m.MarshalToString(req); err != nil {")
	p("    return nil, err")
	p("}")

	// _ := http_opt(m)

	// * JSONify the request
	// * determine the url
	// ** Handle query parameters
	// * send the request, receive the response, maybe return an error
	// * deserialize the result and return
	p("}")

	return nil
}
