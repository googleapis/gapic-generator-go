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
	"regexp"
	"sort"
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/gapic-generator-go/internal/errors"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
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
	p("  // The http endpoint to connect to.")
	p("  endpoint string")
	p("")
	p("  // The http client.")
	p("  httpClient *http.Client")
	p("")
	p("	 // The x-goog-* metadata to be sent with each request.")
	p("	 xGoogMetadata metadata.MD")
	p("}")
	p("")
	g.restClientUtilities(serv, servName, imp, hasRPCForLRO)

	g.imports[pbinfo.ImportSpec{Path: "google.golang.org/protobuf/encoding/protojson"}] = true
	g.imports[pbinfo.ImportSpec{Path: "net/http"}] = true
	g.imports[pbinfo.ImportSpec{Path: "io/ioutil"}] = true
	g.imports[pbinfo.ImportSpec{Path: "fmt"}] = true
	g.imports[pbinfo.ImportSpec{Path: "bytes"}] = true
	g.imports[pbinfo.ImportSpec{Path: "strings"}] = true
	g.imports[pbinfo.ImportSpec{Path: "google.golang.org/grpc/metadata"}] = true
	g.imports[pbinfo.ImportSpec{Name: "httptransport", Path: "google.golang.org/api/transport/http"}] = true
	g.imports[pbinfo.ImportSpec{Path: "google.golang.org/api/option/internaloption"}] = true
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

func (g *generator) restClientOptions(serv *descriptor.ServiceDescriptorProto, servName string) error {
	if !proto.HasExtension(serv.GetOptions(), annotations.E_DefaultHost) {
		// Not an error, just doesn't apply to us.
		return nil
	}

	p := g.printf

	var host string
	eHost := proto.GetExtension(serv.GetOptions(), annotations.E_DefaultHost)

	host = eHost.(string)

	p("func default%sRESTClientOptions() []option.ClientOption {", servName)
	p("  return []option.ClientOption{")
	p("    internaloption.WithDefaultEndpoint(%q),", host)
	p("    internaloption.WithDefaultMTLSEndpoint(%q),", generateDefaultMTLSEndpoint(host))
	p("    internaloption.WithDefaultAudience(%q),", generateDefaultAudience(host))
	p("    internaloption.WithDefaultScopes(DefaultAuthScopes()...),")
	p("  }")
	p("}")

	return nil
}

func (g *generator) restClientUtilities(serv *descriptor.ServiceDescriptorProto, servName string, imp pbinfo.ImportSpec, hasRPCForLRO bool) {
	p := g.printf
	lowcaseServName := lowcaseRestClientName(servName)
	clientName := camelToSnake(serv.GetName())
	clientName = strings.Replace(clientName, "_", " ", -1)

	p("// New%sRESTClient creates a new %s rest client.", servName, clientName)
	g.serviceDoc(serv)
	p("func New%[1]sRESTClient(ctx context.Context, opts ...option.ClientOption) (*%[1]sClient, error) {", servName)
	p("    clientOpts := append(default%sRESTClientOptions(), opts...)", servName)
	p("    httpClient, endpoint, err := httptransport.NewClient(ctx, clientOpts...)")
	p("    if err != nil {")
	p("        return nil, err")
	p("    }")
	p("")
	p("    c := &%s{", lowcaseServName)
	p("        endpoint: endpoint,")
	p("        httpClient: httpClient,")
	p("    }")
	p("    c.setGoogleClientInfo()")
	p("")
	// TODO(dovs): make rest default call options
	// TODO(dovs): set the LRO client
	p("    return &%[1]sClient{internalClient: c, CallOptions: default%[1]sCallOptions()}, nil", servName)
	p("}")
	p("")

	g.restClientOptions(serv, servName)

	// setGoogleClientInfo method
	p("// setGoogleClientInfo sets the name and version of the application in")
	p("// the `x-goog-api-client` header passed on each request. Intended for")
	p("// use by Google-written clients.")
	p("func (c *%s) setGoogleClientInfo(keyval ...string) {", lowcaseServName)
	p(`  kv := append([]string{"gl-go", versionGo()}, keyval...)`)
	p(`  kv = append(kv, "gapic", versionClient, "gax", gax.Version, "rest", "UNKNOWN")`)
	p(`  c.xGoogMetadata = metadata.Pairs("x-goog-api-client", gax.XGoogHeader(kv...))`)
	p("}")
	p("")

	// Close method
	p("// Close closes the connection to the API service. The user should invoke this when")
	p("// the client is no longer required.")
	p("func (c *%s) Close() error {", lowcaseServName)
	p("    c.httpClient.CloseIdleConnections()")
	p("    return nil")
	p("}")
	p("")

	p("// Connection returns a connection to the API service.")
	p("//")
	p("// Deprecated.")
	p("func (c *%s) Connection() *grpc.ClientConn {", lowcaseServName)
	p("    return nil")
	p("}")
}

type httpInfo struct {
	verb, url, body string
}

func (g *generator) pathParams(m *descriptor.MethodDescriptorProto) map[string]*descriptor.FieldDescriptorProto {
	pathParams := map[string]*descriptor.FieldDescriptorProto{}
	info, err := getHTTPInfo(m)
	if info == nil || err != nil {
		return pathParams
	}

	// Match using the curly braces but don't include them in the grouping.
	re := regexp.MustCompile("{([^}]+)}")
	for _, p := range re.FindAllStringSubmatch(info.url, -1) {
		// In the returned slice, the zeroth element is the full regex match,
		// and the subsequent elements are the sub group matches.
		// See the docs for FindStringSubmatch for further details.
		param := p[1]
		field := g.lookupField(m.GetInputType(), param)
		if field == nil {
			continue
		}
		pathParams[param] = field
	}

	return pathParams
}

func (g *generator) queryParams(m *descriptor.MethodDescriptorProto) map[string]*descriptor.FieldDescriptorProto {
	queryParams := map[string]*descriptor.FieldDescriptorProto{}
	info, err := getHTTPInfo(m)
	if info == nil || err != nil {
		return queryParams
	}
	pathParams := g.pathParams(m)
	// Minor hack: we want to make sure that the body parameter, if one exists,
	// is NOT a query parameter.
	dummyField := &descriptor.FieldDescriptorProto{}
	pathParams[info.body] = dummyField
	msg := g.descInfo.Type[m.GetInputType()].(*descriptor.DescriptorProto)
	if info.body != "" && info.body != "*" {
		field := g.lookupField(m.GetInputType(), info.body)
		msg = g.descInfo.Type[field.GetTypeName()].(*descriptor.DescriptorProto)
	}

	// If, and only if, a field is a path parameter, it's not a query parameter.
	for _, f := range msg.GetField() {
		if _, ok := pathParams[f.GetName()]; !ok {
			queryParams[f.GetName()] = f
		}
	}

	return queryParams
}

func (g *generator) generateQueryString(m *descriptor.MethodDescriptorProto) {
	p := g.printf
	queryParams := g.queryParams(m)
	if queryParams == nil {
		return
	}

	// We want to iterate over fields in a deterministic order
	// to prevent spurious deltas when regenerating gapics.
	fields := make([]string, 0, len(queryParams))
	for p, _ := range queryParams {
		fields = append(fields, p)
	}
	sort.Strings(fields)
	info, _ := getHTTPInfo(m)

	requestObject := "req"
	if info.body != "" && info.body != "*" {
		requestObject = "body"
		p("body := req.Get%s()", snakeToCamel(info.body))
	}
	p("queryParamStrs := []string{}")
	// TODO(dovs): make this work for nested fields
	for _, f := range fields {
		field := queryParams[f]
		f = snakeToCamel(f)
		if field.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			// It's a slice, so check for nil
			p("if %s.Get%s() != nil {", requestObject, f)
		} else {
			// Default values are type specific
			switch field.GetType() {
			case descriptor.FieldDescriptorProto_TYPE_MESSAGE, descriptor.FieldDescriptorProto_TYPE_BYTES:
				p("if %s.Get%s() != nil {", requestObject, f)
			case descriptor.FieldDescriptorProto_TYPE_STRING:
				p(`if %s.Get%s() != "" {`, requestObject, f)
			case descriptor.FieldDescriptorProto_TYPE_BOOL:
				p(`if %s.Get%s() {`, requestObject, f)
			default: // Handles all numeric types
				p(`if %s.Get%s() != 0 {`, requestObject, f)
			}
		}
		p("    queryParamStrs = append(queryParamStrs, fmt.Sprintf(\"%s=%%v\", %s.Get%s()))", lowerFirst(f), requestObject, f)
		p("}")
	}

	p(`query := strings.ReplaceAll(strings.Join(queryParamStrs, "&"), " ", "+")`)
	p("")
}

func (g *generator) generateUrlString(m *descriptor.MethodDescriptorProto) error {
	info, err := getHTTPInfo(m)
	if err != nil {
		return err
	}
	if info == nil {
		return errors.E(nil, "method has no http info: %s", m.GetName())
	}

	p := g.printf

	requestObject := "req"
	if info.body != "" && info.body != "*" {
		// Setting the local variable 'body' is already handled by generateQueryString
		requestObject = "body"
	}

	fmtStr := fmt.Sprintf("%%s%s", info.url)
	// TODO(dovs): handle more complex path urls involving = and *,
	// e.g. v1beta1/repeat/{info.f_string=first/*}/{info.f_child.f_string=second/**}:pathtrailingresource
	re := regexp.MustCompile(`{([a-zA-Z0-9_.]+?)}`)
	fmtStr = re.ReplaceAllStringFunc(fmtStr, func(s string) string { return "%v" })

	p(`url := fmt.Sprintf("%s",`, fmtStr)
	p("c.endpoint,")
	// Can't just reuse pathParams because the order matters
	for _, path := range re.FindAllStringSubmatch(info.url, -1) {
		// In the returned slice, the zeroth element is the full regex match,
		// and the subsequent elements are the sub group matches.
		// See the docs for FindStringSubmatch for further details.
		//
		// buildAccessor handles nested fields for us.
		p("%s%s,", requestObject, buildAccessor(path[1]))
	}
	p(")")
	p("")
	return nil
}

func getHTTPInfo(m *descriptor.MethodDescriptorProto) (*httpInfo, error) {
	if m == nil || m.GetOptions() == nil {
		return nil, nil
	}

	eHTTP := proto.GetExtension(m.GetOptions(), annotations.E_Http)

	httpRule := eHTTP.(*annotations.HttpRule)
	info := httpInfo{body: httpRule.GetBody()}

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

	p("// %[1]s returns a new %[1]s from a given name.", lroType)
	p("// The name must be that of a previously created %s, possibly from a different process.", lroType)
	p("func (c *%s) %[2]s(name string) *%[2]s {", lowcaseServName, lroType)
	p("  return &%s{", lroType)
	// TODO(dovs): return a non-empty and useful object
	p("  }")
	p("}")
	p("")

	g.imports[pbinfo.ImportSpec{Name: "longrunningpb", Path: "google.golang.org/genproto/googleapis/longrunning"}] = true

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
	g.generateQueryString(m)
	g.generateUrlString(m)
	p(`if query != "" {`)
	p(`    url += "?" + query`)
	p("}")
	p("")
	p(`httpReq, err := http.NewRequestWithContext(ctx, "%s", url, bytes.NewReader(jsonReq))`, strings.ToUpper(info.verb))
	p("if err != nil {")
	p("    return err")
	p("}")
	p("")
	p("// Set the headers")
	p("for k, v := range c.xGoogMetadata {")
	p("  httpReq.Header[k] = v")
	p("}")
	// The content type is separate from x-goog metadata.
	// It's a little more verbose to set here explicitly, but
	// it's conceptually cleaner.
	// TODO(dovs) add field headers.
	p(`httpReq.Header["Content-Type"] = []string{"application/json"}`)
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
	// TODO(dovs): handle calloptions
	p("m := protojson.MarshalOptions{AllowPartial: true, EmitUnpopulated: true}")
	p("jsonReq, err := m.Marshal(req)")
	p("if err != nil {")
	p("  return nil, err")
	p("}")
	p("")
	// TOOD(dovs) reenable
	g.generateQueryString(m)
	g.generateUrlString(m)
	p(`if query != "" {`)
	p(`    url += "?" + query`)
	p("}")
	p("")
	p(`httpReq, err := http.NewRequestWithContext(ctx, "%s", url, bytes.NewReader(jsonReq))`, strings.ToUpper(info.verb))
	p("if err != nil {")
	p("    return nil, err")
	p("}")
	p("// Set the headers")
	p("for k, v := range c.xGoogMetadata {")
	p("  httpReq.Header[k] = v")
	p("}")
	// The content type is separate from x-goog metadata.
	// It's a little more verbose to set here explicitly, but
	// it's conceptually cleaner.
	// TODO(dovs) add field headers.
	p(`httpReq.Header["Content-Type"] = []string{"application/json",}`)
	p("")
	p("httpRsp, err := c.httpClient.Do(httpReq)")
	p("if err != nil{")
	p(" return nil, err")
	p("}")
	p("defer httpRsp.Body.Close()")
	p("")
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
	p("unm := protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}")
	p("rsp := &%s.%s{}", outSpec.Name, outType.GetName())
	p("")
	p("return rsp, unm.Unmarshal(buf, rsp)")
	p("}")

	return nil
}
