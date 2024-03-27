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
	"net/http"
	"regexp"
	"sort"
	"strings"

	conf "github.com/googleapis/gapic-generator-go/internal/grpc_service_config"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

var httpPatternVarRegex = regexp.MustCompile(`{([a-zA-Z0-9_.]+?)(=[^{}]+)?}`)

// Derived from internal source: go/http-canonical-mapping.
var gRPCToHTTP map[code.Code]string = map[code.Code]string{
	code.Code_OK:                  "http.StatusOK",
	code.Code_CANCELLED:           "499", // There isn't a Go constant ClientClosedConnection
	code.Code_UNKNOWN:             "http.StatusInternalServerError",
	code.Code_INVALID_ARGUMENT:    "http.StatusBadRequest",
	code.Code_DEADLINE_EXCEEDED:   "http.StatusGatewayTimeout",
	code.Code_NOT_FOUND:           "http.StatusNotFound",
	code.Code_ALREADY_EXISTS:      "http.StatusConflict",
	code.Code_PERMISSION_DENIED:   "http.StatusForbidden",
	code.Code_UNAUTHENTICATED:     "http.StatusUnauthorized",
	code.Code_RESOURCE_EXHAUSTED:  "http.StatusTooManyRequests",
	code.Code_FAILED_PRECONDITION: "http.StatusBadRequest",
	code.Code_ABORTED:             "http.StatusConflict",
	code.Code_OUT_OF_RANGE:        "http.StatusBadRequest",
	code.Code_UNIMPLEMENTED:       "http.StatusNotImplemented",
	code.Code_INTERNAL:            "http.StatusInternalServerError",
	code.Code_UNAVAILABLE:         "http.StatusServiceUnavailable",
	code.Code_DATA_LOSS:           "http.StatusInternalServerError",
}

func lowcaseRestClientName(servName string) string {
	if servName == "" {
		return "restClient"
	}

	return lowerFirst(servName + "RESTClient")
}

func (g *generator) restClientInit(serv *descriptorpb.ServiceDescriptorProto, servName string, hasRPCForLRO bool) {
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
	if hasRPCForLRO {
		p("// LROClient is used internally to handle long-running operations.")
		p("// It is exposed so that its CallOptions can be modified if required.")
		p("// Users should not Close this client.")
		p("LROClient **lroauto.OperationsClient")
		p("")
		g.imports[pbinfo.ImportSpec{Name: "lroauto", Path: "cloud.google.com/go/longrunning/autogen"}] = true
	}
	if opServ, ok := g.customOpServices[serv]; ok {
		opServName := pbinfo.ReduceServName(opServ.GetName(), g.opts.pkgName)
		p("// operationClient is used to call the operation-specific management service.")
		p("operationClient *%sClient", opServName)
		p("")
	}
	p("	 // The x-goog-* headers to be sent with each request.")
	p("	 xGoogHeaders []string")
	p("")
	p("  // Points back to the CallOptions field of the containing %sClient", servName)
	p("  CallOptions **%sCallOptions", servName)
	p("}")
	p("")
	g.restClientUtilities(serv, servName, hasRPCForLRO)

	g.imports[pbinfo.ImportSpec{Path: "net/http"}] = true
	g.imports[pbinfo.ImportSpec{Name: "httptransport", Path: "google.golang.org/api/transport/http"}] = true
	g.imports[pbinfo.ImportSpec{Path: "google.golang.org/api/option/internaloption"}] = true
}

func (g *generator) genRESTMethods(serv *descriptorpb.ServiceDescriptorProto, servName string) error {
	g.addMetadataServiceForTransport(serv.GetName(), "rest", servName)

	methods := append(serv.GetMethod(), g.getMixinMethods()...)

	for _, m := range methods {
		g.methodDoc(m, serv)
		if err := g.genRESTMethod(servName, serv, m); err != nil {
			return fmt.Errorf("error generating method %q: %v", m.GetName(), err)
		}
		g.addMetadataMethod(serv.GetName(), "rest", m.GetName())
	}

	return nil
}

func (g *generator) restClientOptions(serv *descriptorpb.ServiceDescriptorProto, servName string) {
	if !proto.HasExtension(serv.GetOptions(), annotations.E_DefaultHost) {
		// Not an error, just doesn't apply to us.
		return
	}

	p := g.printf

	eHost := proto.GetExtension(serv.GetOptions(), annotations.E_DefaultHost)

	// Default to https, just as gRPC defaults to a secure connection.
	host := fmt.Sprintf("https://%s", eHost.(string))

	p("func default%sRESTClientOptions() []option.ClientOption {", servName)
	p("  return []option.ClientOption{")
	p("    internaloption.WithDefaultEndpoint(%q),", host)
	p("    internaloption.WithDefaultEndpointTemplate(%q),", generateDefaultEndpointTemplate(host))
	p("    internaloption.WithDefaultMTLSEndpoint(%q),", generateDefaultMTLSEndpoint(host))
	p("    internaloption.WithDefaultUniverseDomain(%q),", googleDefaultUniverse)
	p("    internaloption.WithDefaultAudience(%q),", generateDefaultAudience(host))
	p("    internaloption.WithDefaultScopes(DefaultAuthScopes()...),")
	p("  }")
	p("}")
}

func (g *generator) restClientUtilities(serv *descriptorpb.ServiceDescriptorProto, servName string, hasRPCForLRO bool) {
	p := g.printf
	lowcaseServName := lowcaseRestClientName(servName)
	clientName := camelToSnake(serv.GetName())
	clientName = strings.Replace(clientName, "_", " ", -1)
	opServ, hasCustomOp := g.customOpServices[serv]

	p("// New%sRESTClient creates a new %s rest client.", servName, clientName)
	g.serviceDoc(serv)
	p("func New%[1]sRESTClient(ctx context.Context, opts ...option.ClientOption) (*%[1]sClient, error) {", servName)
	p("    clientOpts := append(default%sRESTClientOptions(), opts...)", servName)
	p("    httpClient, endpoint, err := httptransport.NewClient(ctx, clientOpts...)")
	p("    if err != nil {")
	p("        return nil, err")
	p("    }")
	p("")
	p("    callOpts := default%sRESTCallOptions()", servName)
	p("    c := &%s{", lowcaseServName)
	p("        endpoint: endpoint,")
	p("        httpClient: httpClient,")
	p("        CallOptions: &callOpts,")
	p("    }")
	p("    c.setGoogleClientInfo()")
	p("")
	if hasRPCForLRO {
		p("lroOpts := []option.ClientOption{")
		p("  option.WithHTTPClient(httpClient),")
		p("  option.WithEndpoint(endpoint),")
		p("}")
		p("opClient, err := lroauto.NewOperationsRESTClient(ctx, lroOpts...)")
		p("if err != nil {")
		p("    return nil, err")
		p("}")
		p("c.LROClient = &opClient")
		p("")
	}
	if hasCustomOp {
		opServName := pbinfo.ReduceServName(opServ.GetName(), g.opts.pkgName)
		p("o := []option.ClientOption{")
		p("  option.WithHTTPClient(httpClient),")
		p("  option.WithEndpoint(endpoint),")
		p("}")
		p("opC, err := New%sRESTClient(ctx, o...)", opServName)
		p("if err != nil {")
		p("  return nil, err")
		p("}")
		p("c.operationClient = opC")
		p("")
		g.imports[pbinfo.ImportSpec{Path: "google.golang.org/api/option"}] = true
	}

	// TODO(dovs): set the LRO client
	p("    return &%[1]sClient{internalClient: c, CallOptions: callOpts}, nil", servName)
	p("}")
	p("")

	g.restClientOptions(serv, servName)

	// setGoogleClientInfo method
	p("// setGoogleClientInfo sets the name and version of the application in")
	p("// the `x-goog-api-client` header passed on each request. Intended for")
	p("// use by Google-written clients.")
	p("func (c *%s) setGoogleClientInfo(keyval ...string) {", lowcaseServName)
	p(`  kv := append([]string{"gl-go", gax.GoVersion}, keyval...)`)
	p(`  kv = append(kv, "gapic", getVersionClient(), "gax", gax.Version, "rest", "UNKNOWN")`)
	p(`  c.xGoogHeaders = []string{"x-goog-api-client", gax.XGoogHeader(kv...)}`)
	p("}")
	p("")

	// Close method
	p("// Close closes the connection to the API service. The user should invoke this when")
	p("// the client is no longer required.")
	p("func (c *%s) Close() error {", lowcaseServName)
	p("    // Replace httpClient with nil to force cleanup.")
	p("    c.httpClient = nil")
	if hasCustomOp {
		p("if err := c.operationClient.Close(); err != nil {")
		p("  return err")
		p("}")
	}
	p("    return nil")
	p("}")
	p("")

	p("// Connection returns a connection to the API service.")
	p("//")
	p("// Deprecated: This method always returns nil.")
	p("func (c *%s) Connection() *grpc.ClientConn {", lowcaseServName)
	p("    return nil")
	p("}")
}

type httpInfo struct {
	verb, url, body string
}

func (g *generator) pathParams(m *descriptorpb.MethodDescriptorProto) map[string]*descriptorpb.FieldDescriptorProto {
	pathParams := map[string]*descriptorpb.FieldDescriptorProto{}
	info := getHTTPInfo(m)
	if info == nil {
		return pathParams
	}

	// Match using the curly braces but don't include them in the grouping.
	re := regexp.MustCompile("{([^}]+)}")
	for _, p := range re.FindAllStringSubmatch(info.url, -1) {
		// In the returned slice, the zeroth element is the full regex match,
		// and the subsequent elements are the sub group matches.
		// See the docs for FindStringSubmatch for further details.
		param := strings.Split(p[1], "=")[0]
		field := g.lookupField(m.GetInputType(), param)
		if field == nil {
			continue
		}
		pathParams[param] = field
	}

	return pathParams
}

func (g *generator) queryParams(m *descriptorpb.MethodDescriptorProto) map[string]*descriptorpb.FieldDescriptorProto {
	queryParams := map[string]*descriptorpb.FieldDescriptorProto{}
	info := getHTTPInfo(m)
	if info == nil {
		return queryParams
	}
	if info.body == "*" {
		// The entire request is the REST body.
		return queryParams
	}

	pathParams := g.pathParams(m)
	// Minor hack: we want to make sure that the body parameter is NOT a query parameter.
	pathParams[info.body] = &descriptorpb.FieldDescriptorProto{}

	request := g.descInfo.Type[m.GetInputType()].(*descriptorpb.DescriptorProto)
	// Body parameters are fields present in the request body.
	// This may be the request message itself or a subfield.
	// Body parameters are not valid query parameters,
	// because that means the same param would be sent more than once.
	bodyField := g.lookupField(m.GetInputType(), info.body)

	// Possible query parameters are all leaf fields in the request or body.
	pathToLeaf := g.getLeafs(request, bodyField)
	// Iterate in sorted order to
	for path, leaf := range pathToLeaf {
		// If, and only if, a leaf field is not a path parameter or a body parameter,
		// it is a query parameter.
		if _, ok := pathParams[path]; !ok && g.lookupField(request.GetName(), leaf.GetName()) == nil {
			queryParams[path] = leaf
		}
	}

	return queryParams
}

// Returns a map from fully qualified path to field descriptor for all the leaf fields of a message 'm',
// where a "leaf" field is a non-message whose top message ancestor is 'm'.
// e.g. for a message like the following
//
//	message Mollusc {
//	    message Squid {
//	        message Mantle {
//	            int32 mass_kg = 1;
//	        }
//	        Mantle mantle = 1;
//	    }
//	    Squid squid = 1;
//	}
//
// The one entry would be
// "squid.mantle.mass_kg": *descriptorpb.FieldDescriptorProto...
func (g *generator) getLeafs(msg *descriptorpb.DescriptorProto, excludedFields ...*descriptorpb.FieldDescriptorProto) map[string]*descriptorpb.FieldDescriptorProto {
	pathsToLeafs := map[string]*descriptorpb.FieldDescriptorProto{}

	contains := func(fields []*descriptorpb.FieldDescriptorProto, field *descriptorpb.FieldDescriptorProto) bool {
		for _, f := range fields {
			if field == f {
				return true
			}
		}
		return false
	}

	// We need to declare and define this function in two steps
	// so that we can use it recursively.
	var recurse func([]*descriptorpb.FieldDescriptorProto, *descriptorpb.DescriptorProto)

	handleLeaf := func(field *descriptorpb.FieldDescriptorProto, stack []*descriptorpb.FieldDescriptorProto) {
		elts := []string{}
		for _, f := range stack {
			elts = append(elts, f.GetName())
		}
		elts = append(elts, field.GetName())
		key := strings.Join(elts, ".")
		pathsToLeafs[key] = field
	}

	handleMsg := func(field *descriptorpb.FieldDescriptorProto, stack []*descriptorpb.FieldDescriptorProto) {
		if field.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED {
			// Repeated message fields must not be mapped because no
			// client library can support such complicated mappings.
			// https://cloud.google.com/endpoints/docs/grpc-service-config/reference/rpc/google.api#grpc-transcoding
			return
		}
		if contains(excludedFields, field) {
			return
		}
		// Short circuit on infinite recursion
		if contains(stack, field) {
			return
		}

		subMsg := g.descInfo.Type[field.GetTypeName()].(*descriptorpb.DescriptorProto)
		recurse(append(stack, field), subMsg)
	}

	recurse = func(
		stack []*descriptorpb.FieldDescriptorProto,
		m *descriptorpb.DescriptorProto,
	) {
		for _, field := range m.GetField() {
			if field.GetType() == fieldTypeMessage && !strContains(wellKnownTypeNames, field.GetTypeName()) {
				handleMsg(field, stack)
			} else {
				handleLeaf(field, stack)
			}
		}
	}

	recurse([]*descriptorpb.FieldDescriptorProto{}, msg)
	return pathsToLeafs
}

func (g *generator) generateQueryString(m *descriptorpb.MethodDescriptorProto) {
	p := g.printf
	queryParams := g.queryParams(m)

	// We want to iterate over fields in a deterministic order
	// to prevent spurious deltas when regenerating gapics.
	fields := make([]string, 0, len(queryParams))
	for p := range queryParams {
		fields = append(fields, p)
	}
	sort.Strings(fields)

	if g.opts.restNumericEnum || len(fields) > 0 {
		g.imports[pbinfo.ImportSpec{Path: "net/url"}] = true
		p("params := url.Values{}")
	}
	if g.opts.restNumericEnum {
		p(`params.Add("$alt", "json;enum-encoding=int")`)
	}
	for _, path := range fields {
		field := queryParams[path]
		required := isRequired(field)
		accessor := fieldGetter(path)
		singularPrimitive := field.GetType() != fieldTypeMessage &&
			field.GetType() != fieldTypeBytes &&
			field.GetLabel() != fieldLabelRepeated
		key := lowerFirst(snakeToCamel(path))

		var paramAdd string
		// Handle well known protobuf types with special JSON encodings.
		if strContains(wellKnownTypeNames, field.GetTypeName()) {
			b := strings.Builder{}
			b.WriteString(fmt.Sprintf("%s, err := protojson.Marshal(req%s)\n", field.GetJsonName(), accessor))
			b.WriteString("if err != nil {\n")
			if m.GetOutputType() == emptyType {
				b.WriteString("  return err\n")
			} else if g.isPaginated(m) {
				b.WriteString("  return nil, \"\", err\n")
			} else {
				b.WriteString("  return nil, err\n")
			}
			b.WriteString("}\n")
			// Only some of the well known types will be encoded as strings, remove the wrapping quotations for those.
			if strContains(wellKnownStringTypes, field.GetTypeName()) {
				b.WriteString(fmt.Sprintf("params.Add(%q, string(%s[1:len(%[2]s)-1]))", key, field.GetJsonName()))
			} else {
				b.WriteString(fmt.Sprintf("params.Add(%q, string(%s))", key, field.GetJsonName()))
			}
			paramAdd = b.String()
		} else {
			paramAdd = fmt.Sprintf("params.Add(%q, fmt.Sprintf(%q, req%s))", key, "%v", accessor)
			g.imports[pbinfo.ImportSpec{Path: "fmt"}] = true
		}

		// Only required, singular, primitive field types should be added regardless.
		if required && singularPrimitive {
			// Use string format specifier here in order to allow %v to be a raw string.
			p("%s", paramAdd)
			continue
		}

		if field.GetLabel() == fieldLabelRepeated {
			// It's a slice, so check for len > 0, nil slice returns 0.
			p("if items := req%s; len(items) > 0 {", accessor)
			b := strings.Builder{}
			b.WriteString("for _, item := range items {\n")
			b.WriteString(fmt.Sprintf("  params.Add(%q, fmt.Sprintf(%q, item))\n", key, "%v"))
			b.WriteString("}")
			paramAdd = b.String()

		} else if field.GetProto3Optional() {
			// Split right before the raw access
			toks := strings.Split(path, ".")
			toks = toks[:len(toks)-1]
			parentField := fieldGetter(strings.Join(toks, "."))
			directLeafField := directAccess(path)
			p("if req%s != nil && req%s != nil {", parentField, directLeafField)
		} else {
			// Default values are type specific
			switch field.GetType() {
			// Degenerate case, field should never be a message because that implies it's not a leaf.
			case fieldTypeMessage, fieldTypeBytes:
				p("if req%s != nil {", accessor)
			case fieldTypeString:
				p(`if req%s != "" {`, accessor)
			case fieldTypeBool:
				p(`if req%s {`, accessor)
			default: // Handles all numeric types including enums
				p(`if req%s != 0 {`, accessor)
			}
		}

		// Split on newline so that multi-line param adders will be formatted properly.
		for _, s := range strings.Split(paramAdd, "\n") {
			p("    %s", s)
		}
		p("}")
	}

	if g.opts.restNumericEnum || len(fields) > 0 {
		p("")
		p("baseUrl.RawQuery = params.Encode()")
		p("")
	}
}

func (g *generator) generateBaseURL(info *httpInfo, ret string) {
	p := g.printf

	fmtStr := info.url
	// TODO(noahdietz): handle more complex path urls involving = and *,
	// e.g. v1beta1/repeat/{info.f_string=first/*}/{info.f_child.f_string=second/**}:pathtrailingresource
	fmtStr = httpPatternVarRegex.ReplaceAllStringFunc(fmtStr, func(s string) string { return "%v" })

	g.imports[pbinfo.ImportSpec{Path: "net/url"}] = true
	p("baseUrl, err := url.Parse(c.endpoint)")
	p("if err != nil {")
	p("  %s", ret)
	p("}")

	tokens := []string{fmt.Sprintf("%q", fmtStr)}
	// Can't just reuse pathParams because the order matters
	for _, path := range httpPatternVarRegex.FindAllStringSubmatch(info.url, -1) {
		// In the returned slice, the zeroth element is the full regex match,
		// and the subsequent elements are the sub group matches.
		// See the docs for FindStringSubmatch for further details.
		tokens = append(tokens, fmt.Sprintf("req%s", fieldGetter(path[1])))
	}
	g.imports[pbinfo.ImportSpec{Path: "fmt"}] = true
	p("baseUrl.Path += fmt.Sprintf(%s)", strings.Join(tokens, ", "))
	p("")
}

func getHTTPInfo(m *descriptorpb.MethodDescriptorProto) *httpInfo {
	if m == nil || m.GetOptions() == nil {
		return nil
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
		info.url = httpRule.GetPut()
	case *annotations.HttpRule_Delete:
		info.verb = "delete"
		info.url = httpRule.GetDelete()
	}

	return &info
}

// genRESTMethod generates a single method from a client. m must be a method declared in serv.
// If the generated method requires an auxiliary type, it is added to aux.
func (g *generator) genRESTMethod(servName string, serv *descriptorpb.ServiceDescriptorProto, m *descriptorpb.MethodDescriptorProto) error {
	if g.isLRO(m) {
		if err := g.maybeAddOperationWrapper(m); err != nil {
			return err
		}
		return g.lroRESTCall(servName, m)
	}

	if m.GetOutputType() == emptyType {
		return g.emptyUnaryRESTCall(servName, m)
	}

	if pf, ps, err := g.getPagingFields(m); err != nil {
		return err
	} else if pf != nil {
		iter, err := g.iterTypeOf(pf)
		if err != nil {
			return err
		}

		return g.pagingRESTCall(servName, m, pf, ps, iter)
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

func (g *generator) serverStreamRESTCall(servName string, s *descriptorpb.ServiceDescriptorProto, m *descriptorpb.MethodDescriptorProto) error {
	info := getHTTPInfo(m)
	if info == nil {
		return fmt.Errorf("method has no http info: %s", m.GetName())
	}

	inType := g.descInfo.Type[m.GetInputType()]

	inSpec, err := g.descInfo.ImportSpec(inType)
	if err != nil {
		return err
	}
	g.imports[inSpec] = true

	outType := g.descInfo.Type[m.GetOutputType()]

	outSpec, err := g.descInfo.ImportSpec(outType)
	if err != nil {
		return err
	}
	g.imports[outSpec] = true

	servSpec, err := g.descInfo.ImportSpec(s)
	if err != nil {
		return err
	}
	g.imports[servSpec] = true

	p := g.printf
	lowcaseServName := lowcaseRestClientName(servName)
	streamClient := fmt.Sprintf("%sRESTClient", lowerFirst(m.GetName()))

	// rest-client method
	p("func (c *%s) %s(ctx context.Context, req *%s.%s, opts ...gax.CallOption) (%s.%s_%sClient, error) {",
		lowcaseServName, m.GetName(), inSpec.Name, inType.GetName(), servSpec.Name, s.GetName(), m.GetName())
	body := "nil"
	verb := strings.ToUpper(info.verb)

	// Marshal body for HTTP methods that take a body.
	if info.body != "" {
		if verb == http.MethodGet || verb == http.MethodDelete {
			return fmt.Errorf("invalid use of body parameter for a get/delete method %q", m.GetName())
		}
		g.protoJSONMarshaler()
		requestObject := "req"
		if info.body != "*" {
			requestObject = "body"
			p("body := req%s", fieldGetter(info.body))
		}
		p("jsonReq, err := m.Marshal(%s)", requestObject)
		p("if err != nil {")
		p("  return nil, err")
		p("}")
		p("")

		body = "bytes.NewReader(jsonReq)"
		g.imports[pbinfo.ImportSpec{Path: "bytes"}] = true
		g.imports[pbinfo.ImportSpec{Path: "google.golang.org/protobuf/encoding/protojson"}] = true
	}

	g.generateBaseURL(info, "return nil, err")
	g.generateQueryString(m)
	p("// Build HTTP headers from client and context metadata.")
	g.insertRequestHeaders(m, rest)
	p("var streamClient *%s", streamClient)
	p("e := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {")
	p(`  if settings.Path != "" {`)
	p("    baseUrl.Path = settings.Path")
	p("  }")
	p(`  httpReq, err := http.NewRequest("%s", baseUrl.String(), %s)`, verb, body)
	p("  if err != nil {")
	p("      return err")
	p("  }")
	p("  httpReq = httpReq.WithContext(ctx)")
	p("  httpReq.Header = headers")
	p("")
	p("  httpRsp, err := c.httpClient.Do(httpReq)")
	p("  if err != nil{")
	p("   return err")
	p("  }")
	p("")
	p("  if err = googleapi.CheckResponse(httpRsp); err != nil {")
	p("    return err")
	p("  }")
	p("")
	p("  streamClient = &%s{", streamClient)
	p("    ctx: ctx,")
	p("    md: metadata.MD(httpRsp.Header),")
	p("    stream: gax.NewProtoJSONStreamReader(httpRsp.Body, (&%s.%s{}).ProtoReflect().Type()),", outSpec.Name, outType.GetName())
	p("  }")
	p("  return nil")
	p("}, opts...)")
	p("")
	p("return streamClient, e")
	p("}")
	p("")

	g.imports[pbinfo.ImportSpec{Path: "google.golang.org/api/googleapi"}] = true

	// server-stream wrapper client
	p("// %s is the stream client used to consume the server stream created by", streamClient)
	p("// the REST implementation of %s.", m.GetName())
	p("type %s struct {", streamClient)
	p("  ctx context.Context")
	p("  md metadata.MD")
	p("  stream *gax.ProtoJSONStream")
	p("}")
	p("")
	p("func (c *%s) Recv() (*%s.%s, error) {", streamClient, outSpec.Name, outType.GetName())
	p("  if err := c.ctx.Err(); err != nil {")
	p("    defer c.stream.Close()")
	p("    return nil, err")
	p("  }")
	p("  msg, err := c.stream.Recv()")
	p("  if err != nil {")
	p("    defer c.stream.Close()")
	p("    return nil, err")
	p("  }")
	p("  res := msg.(*%s.%s)", outSpec.Name, outType.GetName())
	p("  return res, nil")
	p("}")
	p("")
	p("func (c *%s) Header() (metadata.MD, error) {", streamClient)
	p("  return c.md, nil")
	p("}")
	p("")
	p("func (c *%s) Trailer() metadata.MD {", streamClient)
	p("  return c.md")
	p("}")
	p("")
	p("func (c *%s) CloseSend() error {", streamClient)
	p("  // This is a no-op to fulfill the interface.")
	p(`  return fmt.Errorf("this method is not implemented for a server-stream")`)
	p("}")
	p("")
	p("func (c *%s) Context() context.Context {", streamClient)
	p("  return c.ctx")
	p("}")
	p("")
	p("func (c *%s) SendMsg(m interface{}) error {", streamClient)
	p("  // This is a no-op to fulfill the interface.")
	p(`  return fmt.Errorf("this method is not implemented for a server-stream")`)
	p("}")
	p("")
	p("func (c *%s) RecvMsg(m interface{}) error {", streamClient)
	p("  // This is a no-op to fulfill the interface.")
	p(`  return fmt.Errorf("this method is not implemented, use Recv")`)
	p("}")
	p("")

	g.imports[pbinfo.ImportSpec{Path: "context"}] = true
	g.imports[pbinfo.ImportSpec{Path: "fmt"}] = true
	g.imports[pbinfo.ImportSpec{Path: "google.golang.org/grpc/metadata"}] = true

	return nil
}

func (g *generator) noRequestStreamRESTCall(servName string, s *descriptorpb.ServiceDescriptorProto, m *descriptorpb.MethodDescriptorProto) error {
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

	g.imports[pbinfo.ImportSpec{Path: "context"}] = true
	g.imports[pbinfo.ImportSpec{Path: "fmt"}] = true

	return nil
}

func (g *generator) pagingRESTCall(servName string, m *descriptorpb.MethodDescriptorProto, elemField, pageSize *descriptorpb.FieldDescriptorProto, pt *iterType) error {
	lowcaseServName := lowcaseRestClientName(servName)
	p := g.printf

	inType := g.descInfo.Type[m.GetInputType()].(*descriptorpb.DescriptorProto)
	outType := g.descInfo.Type[m.GetOutputType()].(*descriptorpb.DescriptorProto)

	inSpec, err := g.descInfo.ImportSpec(inType)
	if err != nil {
		return err
	}

	outSpec, err := g.descInfo.ImportSpec(outType)
	if err != nil {
		return err
	}
	info := getHTTPInfo(m)
	if err != nil {
		return err
	}
	if info == nil {
		return fmt.Errorf("method has no http info: %q", m.GetName())
	}

	verb := strings.ToUpper(info.verb)

	max := "math.MaxInt32"
	g.imports[pbinfo.ImportSpec{Path: "math"}] = true
	psTyp := pbinfo.GoTypeForPrim[pageSize.GetType()]
	ps := fmt.Sprintf("%s(pageSize)", psTyp)
	if isOptional(inType, pageSize.GetName()) {
		max = fmt.Sprintf("proto.%s(%s)", upperFirst(psTyp), max)
		ps = fmt.Sprintf("proto.%s(%s)", upperFirst(psTyp), ps)
	}
	tok := "pageToken"
	if isOptional(inType, "page_token") {
		tok = fmt.Sprintf("proto.String(%s)", tok)
	}

	pageSizeFieldName := snakeToCamel(pageSize.GetName())
	p("func (c *%s) %s(ctx context.Context, req *%s.%s, opts ...gax.CallOption) *%s {",
		lowcaseServName, m.GetName(), inSpec.Name, inType.GetName(), pt.iterTypeName)
	p("it := &%s{}", pt.iterTypeName)
	p("req = proto.Clone(req).(*%s.%s)", inSpec.Name, inType.GetName())

	maybeReqBytes := "nil"
	if info.body != "" {
		g.protoJSONMarshaler()
		maybeReqBytes = "bytes.NewReader(jsonReq)"
		g.imports[pbinfo.ImportSpec{Path: "bytes"}] = true
	}

	p("unm := protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}")
	p("it.InternalFetch = func(pageSize int, pageToken string) ([]%s, string, error) {", pt.elemTypeName)
	g.internalFetchSetup(outType, outSpec, tok, pageSizeFieldName, max, ps)

	if info.body != "" {
		p("  jsonReq, err := m.Marshal(req)")
		p("  if err != nil {")
		p(`    return nil, "", err`)
		p("  }")
		p("")
	}

	g.generateBaseURL(info, `return nil, "", err`)
	g.generateQueryString(m)
	p("  // Build HTTP headers from client and context metadata.")
	p(`  hds := append(c.xGoogHeaders, "Content-Type", "application/json")`)
	p(`  headers := gax.BuildHeaders(ctx, hds...)`)
	p("  e := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {")
	p(`    if settings.Path != "" {`)
	p("      baseUrl.Path = settings.Path")
	p("    }")
	p(`    httpReq, err := http.NewRequest("%s", baseUrl.String(), %s)`, verb, maybeReqBytes)
	p("    if err != nil {")
	p(`      return err`)
	p("    }")
	// TODO: Should this http.Request use WithContext?
	p("    httpReq.Header = headers")
	p("")
	p("    httpRsp, err := c.httpClient.Do(httpReq)")
	p("    if err != nil{")
	p(`     return err`)
	p("    }")
	p("    defer httpRsp.Body.Close()")
	p("")
	p("    if err = googleapi.CheckResponse(httpRsp); err != nil {")
	p(`      return err`)
	p("    }")
	p("")
	p("    buf, err := io.ReadAll(httpRsp.Body)")
	p("    if err != nil {")
	p(`      return err`)
	p("    }")
	p("")
	p("    if err := unm.Unmarshal(buf, resp); err != nil {")
	p("      return err")
	p("    }")
	p("")
	p("    return nil")
	p("  }, opts...)")
	p("  if e != nil {")
	p(`    return nil, "", e`)
	p("  }")
	p("  it.Response = resp")
	elems := g.maybeSortMapPage(elemField, pt)
	p("  return %s, resp.GetNextPageToken(), nil", elems)
	p("}")
	p("")
	g.makeFetchAndIterUpdate(pageSizeFieldName)
	p("}")

	g.imports[pbinfo.ImportSpec{Path: "google.golang.org/api/iterator"}] = true
	g.imports[pbinfo.ImportSpec{Path: "io"}] = true
	g.imports[pbinfo.ImportSpec{Path: "google.golang.org/protobuf/proto"}] = true
	g.imports[pbinfo.ImportSpec{Path: "google.golang.org/protobuf/encoding/protojson"}] = true
	g.imports[pbinfo.ImportSpec{Path: "google.golang.org/api/googleapi"}] = true
	g.imports[inSpec] = true
	g.imports[outSpec] = true

	return nil
}

func (g *generator) lroRESTCall(servName string, m *descriptorpb.MethodDescriptorProto) error {
	info := getHTTPInfo(m)
	if info == nil {
		return fmt.Errorf("method has no http info: %s", m.GetName())
	}

	lowcaseServName := lowcaseRestClientName(servName)
	p := g.printf
	inType := g.descInfo.Type[m.GetInputType()].(*descriptorpb.DescriptorProto)
	outType := g.descInfo.Type[m.GetOutputType()]

	inSpec, err := g.descInfo.ImportSpec(inType)
	if err != nil {
		return err
	}
	// This is always the longrunningp.Operation/google.longrunning.Operation.
	outSpec, err := g.descInfo.ImportSpec(outType)
	if err != nil {
		return err
	}

	// TODO(codyoss): This if can be removed once the public protos
	// have been migrated to their new package. This should be soon after this
	// code is merged.
	if outSpec.Path == "google.golang.org/genproto/googleapis/longrunning" {
		outSpec.Path = "cloud.google.com/go/longrunning/autogen/longrunningpb"
	}
	g.imports[outSpec] = true

	opWrapperType := lroTypeName(m)
	p("func (c *%s) %s(ctx context.Context, req *%s.%s, opts ...gax.CallOption) (*%s, error) {",
		lowcaseServName, m.GetName(), inSpec.Name, inType.GetName(), opWrapperType)

	// TODO(noahdietz): handle cancellation, metadata, osv.
	// TODO(noahdietz): handle http headers
	// TODO(noahdietz): handle deadlines?
	// TODO(noahdietz): handle calloptions

	body := "nil"
	verb := strings.ToUpper(info.verb)

	// Marshal body for HTTP methods that take a body.
	if info.body != "" {
		if verb == http.MethodGet || verb == http.MethodDelete {
			return fmt.Errorf("invalid use of body parameter for a get/delete method %q", m.GetName())
		}
		g.protoJSONMarshaler()
		requestObject := "req"
		if info.body != "*" {
			requestObject = "body"
			p("body := req%s", fieldGetter(info.body))
		}
		p("jsonReq, err := m.Marshal(%s)", requestObject)
		p("if err != nil {")
		p("  return nil, err")
		p("}")
		p("")

		body = "bytes.NewReader(jsonReq)"
		g.imports[pbinfo.ImportSpec{Path: "bytes"}] = true
	}

	g.generateBaseURL(info, "return nil, err")
	g.generateQueryString(m)
	p("// Build HTTP headers from client and context metadata.")
	g.insertRequestHeaders(m, rest)
	p("unm := protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}")
	p("resp := &%s.%s{}", outSpec.Name, outType.GetName())
	p("e := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {")
	p(`  if settings.Path != "" {`)
	p("    baseUrl.Path = settings.Path")
	p("  }")
	p(`  httpReq, err := http.NewRequest("%s", baseUrl.String(), %s)`, verb, body)
	p("  if err != nil {")
	p("      return err")
	p("  }")
	p("  httpReq = httpReq.WithContext(ctx)")
	p("  httpReq.Header = headers")
	p("")
	p("  httpRsp, err := c.httpClient.Do(httpReq)")
	p("  if err != nil{")
	p("   return err")
	p("  }")
	p("  defer httpRsp.Body.Close()")
	p("")
	p("  if err = googleapi.CheckResponse(httpRsp); err != nil {")
	p("    return err")
	p("  }")
	p("")
	p("  buf, err := io.ReadAll(httpRsp.Body)")
	p("  if err != nil {")
	p("    return err")
	p("  }")
	p("")
	p("  if err := unm.Unmarshal(buf, resp); err != nil {")
	p("    return err")
	p("  }")
	p("")
	p("  return nil")
	p("}, opts...)")
	p("if e != nil {")
	p("  return nil, e")
	p("}")
	p("")
	override := g.getOperationPathOverride(g.descInfo.ParentFile[m].GetPackage())
	p("override := fmt.Sprintf(%q, resp.GetName())", override)
	p("return &%s{", opWrapperType)
	p("  lro: longrunning.InternalNewOperation(*c.LROClient, resp),")
	p("  pollPath: override,")
	p("}, nil")
	p("}")
	p("")

	g.imports[pbinfo.ImportSpec{Path: "fmt"}] = true
	g.imports[pbinfo.ImportSpec{Path: "io"}] = true
	g.imports[pbinfo.ImportSpec{Path: "cloud.google.com/go/longrunning"}] = true
	g.imports[pbinfo.ImportSpec{Path: "google.golang.org/api/googleapi"}] = true
	g.imports[pbinfo.ImportSpec{Path: "google.golang.org/protobuf/encoding/protojson"}] = true

	return nil
}

func (g *generator) emptyUnaryRESTCall(servName string, m *descriptorpb.MethodDescriptorProto) error {
	info := getHTTPInfo(m)
	if info == nil {
		return fmt.Errorf("method has no http info: %s", m.GetName())
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

	g.initializeAutoPopulatedFields(servName, m)
	// TODO(dovs): handle cancellation, metadata, osv.
	// TODO(dovs): handle http headers
	// TODO(dovs): handle deadlines
	// TODO(dovs): handle call options

	body := "nil"
	verb := strings.ToUpper(info.verb)

	// Marshal body for HTTP methods that take a body.
	// TODO(dovs): add tests generating methods with(out) a request body.
	if info.body != "" {
		if verb == http.MethodGet || verb == http.MethodDelete {
			return fmt.Errorf("invalid use of body parameter for a get/delete method %q", m.GetName())
		}
		g.protoJSONMarshaler()
		requestObject := "req"
		if info.body != "*" {
			requestObject = "body"
			p("body := req%s", fieldGetter(info.body))
		}
		p("jsonReq, err := m.Marshal(%s)", requestObject)
		p("if err != nil {")
		p("  return err")
		p("}")
		p("")
		body = "bytes.NewReader(jsonReq)"
		g.imports[pbinfo.ImportSpec{Path: "bytes"}] = true
		g.imports[pbinfo.ImportSpec{Path: "google.golang.org/protobuf/encoding/protojson"}] = true
	}

	g.generateBaseURL(info, "return err")
	g.generateQueryString(m)
	p("// Build HTTP headers from client and context metadata.")
	g.insertRequestHeaders(m, rest)
	p("return gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {")
	p(`  if settings.Path != "" {`)
	p("    baseUrl.Path = settings.Path")
	p("  }")
	p(`  httpReq, err := http.NewRequest("%s", baseUrl.String(), %s)`, verb, body)
	p("  if err != nil {")
	p("      return err")
	p("  }")
	p("  httpReq = httpReq.WithContext(ctx)")
	p("  httpReq.Header = headers")
	p("")
	p("  httpRsp, err := c.httpClient.Do(httpReq)")
	p("  if err != nil{")
	p("   return err")
	p("  }")
	p("  defer httpRsp.Body.Close()")
	p("")
	p("  // Returns nil if there is no error, otherwise wraps")
	p("  // the response code and body into a non-nil error")
	p("  return googleapi.CheckResponse(httpRsp)")
	p("  }, opts...)")
	p("}")

	g.imports[pbinfo.ImportSpec{Path: "google.golang.org/api/googleapi"}] = true
	g.imports[inSpec] = true
	return nil
}

func (g *generator) unaryRESTCall(servName string, m *descriptorpb.MethodDescriptorProto) error {
	info := getHTTPInfo(m)
	if info == nil {
		return fmt.Errorf("method has no http info: %s", m.GetName())
	}

	inType := g.descInfo.Type[m.GetInputType()]
	outType := g.descInfo.Type[m.GetOutputType()]

	inSpec, err := g.descInfo.ImportSpec(inType)
	if err != nil {
		return err
	}
	outSpec, err := g.descInfo.ImportSpec(outType)
	if err != nil {
		return err
	}
	outFqn := fmt.Sprintf(".%s.%s", g.descInfo.ParentFile[outType].GetPackage(), outType.GetName())
	isHTTPBodyMessage := outFqn == httpBodyType

	// Ignore error because the only possible error would be from looking up
	// the ImportSpec for the OutputType, which has already happened above.
	retTyp, _ := g.returnType(m)

	isCustomOp := g.isCustomOp(m, info)

	p := g.printf
	lowcaseServName := lowcaseRestClientName(servName)
	p("func (c *%s) %s(ctx context.Context, req *%s.%s, opts ...gax.CallOption) (%s, error) {",
		lowcaseServName, m.GetName(), inSpec.Name, inType.GetName(), retTyp)

	g.initializeAutoPopulatedFields(servName, m)
	// TODO(dovs): handle cancellation, metadata, osv.
	// TODO(dovs): handle http headers
	// TODO(dovs): handle deadlines?
	// TODO(dovs): handle calloptions

	body := "nil"
	verb := strings.ToUpper(info.verb)

	// Marshal body for HTTP methods that take a body.
	// TODO(dovs): add tests generating methods with(out) a request body.
	if info.body != "" {
		if verb == http.MethodGet || verb == http.MethodDelete {
			return fmt.Errorf("invalid use of body parameter for a get/delete method %q", m.GetName())
		}
		g.protoJSONMarshaler()
		requestObject := "req"
		if info.body != "*" {
			requestObject = "body"
			p("body := req%s", fieldGetter(info.body))
		}
		p("jsonReq, err := m.Marshal(%s)", requestObject)
		p("if err != nil {")
		p("  return nil, err")
		p("}")
		p("")

		body = "bytes.NewReader(jsonReq)"
		g.imports[pbinfo.ImportSpec{Path: "bytes"}] = true
		g.imports[pbinfo.ImportSpec{Path: "google.golang.org/protobuf/encoding/protojson"}] = true

	}

	g.generateBaseURL(info, "return nil, err")
	g.generateQueryString(m)
	p("// Build HTTP headers from client and context metadata.")
	g.insertRequestHeaders(m, rest)
	g.appendCallOpts(m)
	if !isHTTPBodyMessage {
		p("unm := protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}")
		g.imports[pbinfo.ImportSpec{Path: "google.golang.org/protobuf/encoding/protojson"}] = true

	}
	p("resp := &%s.%s{}", outSpec.Name, outType.GetName())
	p("e := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {")
	p(`  if settings.Path != "" {`)
	p("    baseUrl.Path = settings.Path")
	p("  }")
	p(`  httpReq, err := http.NewRequest("%s", baseUrl.String(), %s)`, verb, body)
	p("  if err != nil {")
	p("      return err")
	p("  }")
	p("  httpReq = httpReq.WithContext(ctx)")
	p("  httpReq.Header = headers")
	p("")
	p("  httpRsp, err := c.httpClient.Do(httpReq)")
	p("  if err != nil{")
	p("   return err")
	p("  }")
	p("  defer httpRsp.Body.Close()")
	p("")
	p("  if err = googleapi.CheckResponse(httpRsp); err != nil {")
	p("    return err")
	p("  }")
	p("")
	p("  buf, err := io.ReadAll(httpRsp.Body)")
	p("  if err != nil {")
	p("    return err")
	p("  }")
	p("")
	if isHTTPBodyMessage {
		p("resp.Data = buf")
		p(`if headers := httpRsp.Header; len(headers["Content-Type"]) > 0 {`)
		p(`  resp.ContentType = headers["Content-Type"][0]`)
		p("}")
	} else {
		p("if err := unm.Unmarshal(buf, resp); err != nil {")
		p("  return err")
		p("}")
	}

	p("")
	p("  return nil")
	p("}, opts...)")
	p("if e != nil {")
	p("  return nil, e")
	p("}")
	ret := "return resp, nil"
	if isCustomOp {
		opVar := "op"
		g.customOpInit("resp", "req", opVar, inType.(*descriptorpb.DescriptorProto), g.customOpService(m))
		ret = fmt.Sprintf("return %s, nil", opVar)
	}
	p(ret)
	p("}")

	g.imports[pbinfo.ImportSpec{Path: "io"}] = true
	g.imports[pbinfo.ImportSpec{Path: "google.golang.org/api/googleapi"}] = true
	g.imports[inSpec] = true
	g.imports[outSpec] = true
	return nil
}

func (g *generator) protoJSONMarshaler() {
	marshalOpts := "AllowPartial: true, UseEnumNumbers: true"
	// Today, DIREGAPIC protos are generated to represent enum fields as strings
	// except for the (very static) Operation.state field, which is kept as an
	// enum. However, since they are synthetic protos, the enum numbers likely
	// do not match the internal proto enum value, so we do not want to send
	// any enums as numbers with DIREGAPIC clients.
	if g.opts.diregapic {
		marshalOpts = "AllowPartial: true"
	}
	g.pt.Printf("m := protojson.MarshalOptions{%s}", marshalOpts)
}

func (g *generator) restCallOptions(serv *descriptorpb.ServiceDescriptorProto, servName string) {
	p := g.printf

	// defaultCallOptions
	c := g.grpcConf

	methods := append(serv.GetMethod(), g.getMixinMethods()...)

	// read retry params from gRPC ServiceConfig
	p("func default%[1]sRESTCallOptions() *%[1]sCallOptions {", servName)
	p("  return &%sCallOptions{", servName)
	for _, m := range methods {
		sFQN := g.fqn(g.descInfo.ParentElement[m])
		mn := m.GetName()
		p("%s: []gax.CallOption{", mn)

		if timeout, ok := c.Timeout(sFQN, mn); ok {
			p("gax.WithTimeout(%d * time.Millisecond),", timeout)
			g.imports[pbinfo.ImportSpec{Path: "time"}] = true
		}

		if rp, ok := c.RetryPolicy(sFQN, mn); ok && rp != nil && len(rp.GetRetryableStatusCodes()) > 0 {
			p("gax.WithRetry(func() gax.Retryer {")
			p("  return gax.OnHTTPCodes(gax.Backoff{")
			// this ignores max_attempts
			p("    Initial:    %d * time.Millisecond,", conf.ToMillis(rp.GetInitialBackoff()))
			p("    Max:        %d * time.Millisecond,", conf.ToMillis(rp.GetMaxBackoff()))
			p("    Multiplier: %.2f,", rp.GetBackoffMultiplier())
			p("	 },")

			rc := rp.GetRetryableStatusCodes()
			for ndx, c := range rc {
				s := fmt.Sprintf("%s,", gRPCToHTTP[c])
				if ndx == len(rc)-1 {
					s = strings.ReplaceAll(s, ",", ")")
				}

				p(s)
			}
			p("}),")

			// include imports necessary for retry configuration
			g.imports[pbinfo.ImportSpec{Path: "time"}] = true
		}
		p("},")
	}
	p("  }")
	p("}")
	p("")
}
