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
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/gapic-generator-go/internal/errors"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"google.golang.org/genproto/googleapis/api/annotations"
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
	p("  endpoint string")
	p("  httpClient http.Client")
	p("}")
	p("")
	g.restClientUtilities(serv, servName, imp, hasRPCForLRO)

	g.imports[pbinfo.ImportSpec{Path: "google.golang.org/protobuf/encoding/protojson"}] = true
	g.imports[pbinfo.ImportSpec{Path: "net/http"}] = true
	g.imports[pbinfo.ImportSpec{Path: "io/ioutil"}] = true
	g.imports[pbinfo.ImportSpec{Path: "fmt"}] = true
	g.imports[pbinfo.ImportSpec{Path: "bytes"}] = true
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
	// g.serviceDoc(serv)
	// p("func New%[1]sRESTClient(ctx context.Context, opts ...option.ClientOption) (*%[1]sClient, error) {", servName)
	// p("    c := &%s{", servName)
	// p("    }")
	// p("")
	// p("    return &%sClient{internal%[1]sClient: c, CallOptions: default%[1]sCallOptions()}, nil", servName)
	// p("}")

}

type httpInfo struct {
	verb, url, body string
}

func getHTTPInfo(m *descriptor.MethodDescriptorProto) (*httpInfo, error) {
	if m == nil || m.GetOptions() == nil {
		return nil, nil
	}

	eHTTP, err := proto.GetExtension(m.GetOptions(), annotations.E_Http)
	if err == proto.ErrMissingExtension {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	httpRule := eHTTP.(*annotations.HttpRule)
	info := httpInfo{body: httpRule.Body}

	switch httpRule.GetPattern().(type) {
	case *annotations.HttpRule_Get:
		info.verb = "get"
		info.url = httpRule.GetGet()
	case *annotations.HttpRule_Post:
		info.verb = "post"
		info.url = httpRule.GetPost()
	case *annotations.HttpRule_Patch:
		info.verb = "patch"
		info.url = httpRule.GetPatch()
	case *annotations.HttpRule_Put:
		info.verb = "put"
		info.url = httpRule.GetPatch()
	case *annotations.HttpRule_Delete:
		info.verb = "delete"
		info.url = httpRule.GetDelete()
	}

	return &info, nil
}

// genRESTMethod generates a single method from a client. m must be a method declared in serv.
// If the generated method requires an auxiliary type, it is added to aux.
func (g *generator) genRESTMethod(servName string, serv *descriptor.ServiceDescriptorProto, m *descriptor.MethodDescriptorProto) error {
	if g.isLRO(m) {
		g.aux.lros[m] = true
		return g.lroRESTCall(servName, m)
	}

	if m.GetOutputType() == emptyType {
		return g.emptyUnaryRESTCall(servName, m)
	}

	if pf, err := g.pagingField(m); err != nil {
		return err
	} else if pf != nil {
		iter, err := g.iterTypeOf(pf)
		if err != nil {
			return err
		}

		return g.pagingRESTCall(servName, m, pf, iter)
	}

	switch {
	case m.GetClientStreaming():
		return g.noRequestStreamRESTCall(servName, serv, m)
	case m.GetServerStreaming():
		return g.serverStreamRESTCall(servName, serv, m)
	default:
		return g.unaryRESTCall(servName, m)
	}
}

func (g *generator) serverStreamRESTCall(servName string, s *descriptor.ServiceDescriptorProto, m *descriptor.MethodDescriptorProto) error {
	// Streaming calls are not currently supported for REST clients,
	// but the interface signature must be preserved.
	// Unimplemented REST methods will always error.

	inType := g.descInfo.Type[m.GetInputType()]

	inSpec, err := g.descInfo.ImportSpec(inType)
	if err != nil {
		return err
	}
	g.imports[inSpec] = true

	servSpec, err := g.descInfo.ImportSpec(s)
	if err != nil {
		return err
	}
	g.imports[servSpec] = true

	p := g.printf
	lowcaseServName := lowcaseRestClientName(servName)
	p("func (c *%s) %s(ctx context.Context, req *%s.%s, opts ...gax.CallOption) (%s.%s_%sClient, error) {",
		lowcaseServName, m.GetName(), inSpec.Name, inType.GetName(), servSpec.Name, s.GetName(), m.GetName())
	p(`  return nil, fmt.Errorf("%s not yet supported for REST clients")`, m.GetName())
	p("}")
	p("")

	return nil
}

func (g *generator) noRequestStreamRESTCall(servName string, s *descriptor.ServiceDescriptorProto, m *descriptor.MethodDescriptorProto) error {
	// Streaming calls are not currently supported for REST clients,
	// but the interface signature must be preserved.
	// Unimplemented REST methods will always error.

	p := g.printf

	servSpec, err := g.descInfo.ImportSpec(s)
	if err != nil {
		return err
	}
	g.imports[servSpec] = true

	lowcaseServName := lowcaseRestClientName(servName)

	p("func (c *%s) %s(ctx context.Context, opts ...gax.CallOption) (%s.%s_%sClient, error) {",
		lowcaseServName, m.GetName(), servSpec.Name, s.GetName(), m.GetName())
	p(`  return nil, fmt.Errorf("%s not yet supported for REST clients")`, m.GetName())
	p("}")
	p("")

	return nil
}

func (g *generator) pagingRESTCall(servName string, m *descriptor.MethodDescriptorProto, elemField *descriptor.FieldDescriptorProto, pt *iterType) error {
	lowcaseServName := lowcaseRestClientName(servName)
	p := g.printf

	inType := g.descInfo.Type[m.GetInputType()].(*descriptor.DescriptorProto)
	// outType := g.descInfo.Type[m.GetOutputType()].(*descriptor.DescriptorProto)

	inSpec, err := g.descInfo.ImportSpec(inType)
	if err != nil {
		return err
	}

	// outSpec, err := g.descInfo.ImportSpec(outType)
	// if err != nil {
	// 	return err
	// }

	p("func (c *%s) %s(ctx context.Context, req *%s.%s, opts ...gax.CallOption) *%s {",
		lowcaseServName, m.GetName(), inSpec.Name, inType.GetName(), pt.iterTypeName)
	p("it := &%s{}", pt.iterTypeName)
	p("req = proto.Clone(req).(*%s.%s)", inSpec.Name, inType.GetName())
	p("it.InternalFetch = func(pageSize int, pageToken string) ([]%s, string, error) {", pt.elemTypeName)
	p(`return nil, "", fmt.Errorf("%s not yet supported for REST clients")`, m.GetName())
	p("}")
	p("")
	p("fetch := func(pageSize int, pageToken string) (string, error) {")
	p("  items, nextPageToken, err := it.InternalFetch(pageSize, pageToken)")
	p("  if err != nil {")
	p(`    return "", err`)
	p("  }")
	p("  it.items = append(it.items, items...)")
	p("  return nextPageToken, nil")
	p("}")
	p("")
	p("it.pageInfo, it.nextFunc = iterator.NewPageInfo(fetch, it.bufLen, it.takeBuf)")
	p("it.pageInfo.MaxSize = int(req.GetPageSize())")
	p("it.pageInfo.Token = req.GetPageToken()")
	p("")
	p("return it")
	p("}")

	return nil
}

func (g *generator) lroRESTCall(servName string, m *descriptor.MethodDescriptorProto) error {
	lowcaseServName := lowcaseRestClientName(servName)
	p := g.printf
	inType := g.descInfo.Type[m.GetInputType()].(*descriptor.DescriptorProto)
	// outType := g.descInfo.Type[m.GetOutputType()].(*descriptor.DescriptorProto)

	inSpec, err := g.descInfo.ImportSpec(inType)
	if err != nil {
		return err
	}

	// outSpec, err := g.descInfo.ImportSpec(outType)
	// if err != nil {
	// 	return err
	// }

	lroType := lroTypeName(m.GetName())
	p("func (c *%s) %s(ctx context.Context, req *%s.%s, opts ...gax.CallOption) (*%s, error) {",
		lowcaseServName, m.GetName(), inSpec.Name, inType.GetName(), lroType)
	p(`    return nil, fmt.Errorf("%s not yet supported for REST clients")`, m.GetName())
	p("}")
	p("")

	return nil
}

func (g *generator) emptyUnaryRESTCall(servName string, m *descriptor.MethodDescriptorProto) error {
	info, err := getHTTPInfo(m)
	if err != nil {
		return err
	}
	if info == nil {
		return errors.E(nil, "method has no http info: %s", m.GetName())
	}

	inType := g.descInfo.Type[m.GetInputType()]
	inSpec, err := g.descInfo.ImportSpec(inType)
	if err != nil {
		return err
	}

	p := g.printf
	lowcaseServName := lowcaseRestClientName(servName)
	p("func (c *%s) %s(ctx context.Context, req *%s.%s, opts ...gax.CallOption) error {",
		lowcaseServName, m.GetName(), inSpec.Name, inType.GetName())

	// TODO(dovs): handle cancellation, metadata, osv.
	// TODO(dovs): handle http headers
	// TODO(dovs): handle deadlines
	// TODO(dovs): handle call options
	p("// The default (false) for the other options are fine.")
	p("// Field names should be lowerCamel, not snake.")
	p("m := protojson.MarshalOptions{AllowPartial: true, EmitUnpopulated: true, UseProtoNames: false}")
	p("jsonReq, err := m.Marshal(req)")
	p("if err != nil {")
	p("  return err")
	p("}")
	p("")
	// TODO(dovs): handle path parameters
	p(`url := fmt.Sprintf("%%s%s", c.endpoint)`, info.url)
	p(`httpReq, err := http.NewRequestWithContext(ctx, "%s", url, bytes.NewReader(jsonReq))`, strings.ToUpper(info.verb))
	p("if err != nil {")
	p("    return err")
	p("}")
	p("")
	p("httpRsp, err := c.httpClient.Do(httpReq)")
	p("if err != nil{")
	p(" return err")
	p("}")
	p("defer httpRsp.Body.Close()")
	p("if httpRsp.StatusCode != http.StatusOK {")
	// TODO(dovs): handle this error more
	p("  return fmt.Errorf(httpRsp.Status)")
	p("}")
	p("")
	p("return nil")
	p("}")

	return nil
}

func (g *generator) unaryRESTCall(servName string, m *descriptor.MethodDescriptorProto) error {
	info, err := getHTTPInfo(m)
	if err != nil {
		return err
	}
	if info == nil {
		return errors.E(nil, "method has no http info: %s", m.GetName())
	}

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
	// TODO(dovs): handle http headers
	// TODO(dovs): handle deadlines?
	// TODO(dovs): handle call options
	p("// The default (false) for the other options are fine.")
	p("// TODO(dovs): handle path parameters")
	p("marshaler := protojson.MarshalOptions{AllowPartial: true, EmitUnpopulated: true}")
	p("jsonReq, err := marshaler.Marshal(req)")
	p("if err != nil {")
	p("  return nil, err")
	p("}")
	p("")
	p(`url := fmt.Sprintf("%%s%s", c.endpoint)`, info.url)
	p(`httpReq, err := http.NewRequestWithContext(ctx, "%s", url, bytes.NewReader(jsonReq))`, strings.ToUpper(info.verb))
	p("if err != nil {")
	p("    return nil, err")
	p("}")
	p("")
	p("httpRsp, err := c.httpClient.Do(httpReq)")
	p("if err != nil{")
	p(" return nil, err")
	p("}")
	p("defer httpRsp.Body.Close()")
	p("if httpRsp.StatusCode >= http.StatusOK {")
	// TODO(dovs): handle this error more
	p("  return nil, fmt.Errorf(httpRsp.Status)")
	p("}")
	p("")
	p("buf, err := ioutil.ReadAll(httpRsp.Body)")
	p("if err != nil {")
	p("  return nil, err")
	p("}")
	p("")
	p("unmarshaler := protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}")
	p("rsp := &%s.%s{}", outSpec.Name, outType.GetName())
	p("")
	p("return rsp, unmarshaler.Unmarshal(buf, rsp)")
	p("}")
	return nil
}
