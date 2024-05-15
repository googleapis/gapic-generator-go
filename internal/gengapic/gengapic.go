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
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"github.com/googleapis/gapic-generator-go/internal/printer"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

const (
	emptyValue = "google.protobuf.Empty"
	// protoc puts a dot in front of name, signaling that the name is fully qualified.
	emptyType               = "." + emptyValue
	operationType           = ".google.longrunning.Operation"
	httpBodyType            = ".google.api.HttpBody"
	alpha                   = "alpha"
	beta                    = "beta"
	deprecated              = "deprecated"
	disableDeadlinesVar     = "GOOGLE_API_GO_EXPERIMENTAL_DISABLE_DEFAULT_DEADLINE"
	fieldTypeBool           = descriptorpb.FieldDescriptorProto_TYPE_BOOL
	fieldTypeString         = descriptorpb.FieldDescriptorProto_TYPE_STRING
	fieldTypeBytes          = descriptorpb.FieldDescriptorProto_TYPE_BYTES
	fieldTypeMessage        = descriptorpb.FieldDescriptorProto_TYPE_MESSAGE
	fieldLabelRepeated      = descriptorpb.FieldDescriptorProto_LABEL_REPEATED
	defaultPollInitialDelay = "time.Second" // 1 second
	defaultPollMaxDelay     = "time.Minute" // 1 minute
)

var headerParamRegexp = regexp.MustCompile(`{([_.a-z0-9]+)`)

// Gen is the entry point for GAPIC generation via the protoc pluginpb.
func Gen(genReq *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error) {
	g, err := newGenerator(genReq)
	if err != nil {
		return &g.resp, err
	}

	genServs := g.collectServices(genReq)

	if len(genServs) == 0 {
		return &g.resp, nil
	}

	g.checkIAMPolicyOverrides(genServs)

	if g.serviceConfig != nil {
		g.apiName = g.serviceConfig.GetTitle()
	}

	protoPkg := g.descInfo.ParentFile[genServs[0]].GetPackage()

	if op, ok := g.descInfo.Type[fmt.Sprintf(".%s.Operation", protoPkg)]; g.opts.diregapic && ok {
		g.aux.customOp = &customOp{
			message:       op.(*descriptorpb.DescriptorProto),
			handles:       []*descriptorpb.ServiceDescriptorProto{},
			pollingParams: map[*descriptorpb.ServiceDescriptorProto][]string{}}
		g.loadCustomOpServices(genServs)
	}

	// Use the proto package from the parent file of the first Service seen.
	if len(genServs) > 0 {
		g.metadata.ProtoPackage = protoPkg
	}
	g.metadata.LibraryPackage = g.opts.pkgPath
	// Initialize the model that will collect snippet metadata.
	g.snippetMetadata = g.newSnippetsMetadata(protoPkg)

	for _, s := range genServs {
		// TODO(pongad): gapic-generator does not remove the package name here,
		// so even though the client for LoggingServiceV2 is just "Client"
		// the file name is "logging_client.go".
		// Keep the current behavior for now, but we could revisit this later.
		servName := pbinfo.ReduceServName(s.GetName(), "")
		outFile := camelToSnake(servName)
		outFile = filepath.Join(g.opts.outDir, outFile)

		if err := g.genAndCommitSnippets(s); err != nil {
			return &g.resp, fmt.Errorf("error generating snippets for %s: %v ", s.GetName(), err)
		}

		g.reset()
		// If the service has no REST-able RPCs, then a REGAPIC should not be
		// generated for it, even if REST is an enabled transport.
		transports := g.opts.transports
		hasREST := hasRESTMethod(s)
		if !hasREST {
			g.opts.transports = []transport{grpc}
		}
		if err := g.gen(s); err != nil {
			return &g.resp, err
		}
		g.commit(outFile+"_client.go", g.opts.pkgName)

		g.reset()
		if err := g.genExampleFile(s); err != nil {
			return &g.resp, fmt.Errorf("error generating example for %q; %v", s.GetName(), err)
		}
		g.imports[pbinfo.ImportSpec{Name: g.opts.pkgName, Path: g.opts.pkgPath}] = true
		g.commit(outFile+"_client_example_test.go", g.opts.pkgName+"_test")

		// Replace original set of transports for the next service that may have
		// REST-able RPCs.
		if !hasREST {
			g.opts.transports = transports
		}
	}
	if err := g.genAndCommitSnippetMetadata(protoPkg); err != nil {
		return &g.resp, err
	}
	g.reset()
	scopes := collectScopes(genServs)
	serv := genServs[0]

	g.genDocFile(time.Now().Year(), scopes, serv)
	g.resp.File = append(g.resp.File, &pluginpb.CodeGeneratorResponse_File{
		Name:    proto.String(filepath.Join(g.opts.outDir, "doc.go")),
		Content: proto.String(g.pt.String()),
	})

	if g.opts.metadata {
		g.reset()
		g.genGapicMetadataFile()
		g.resp.File = append(g.resp.File, &pluginpb.CodeGeneratorResponse_File{
			Name:    proto.String(filepath.Join(g.opts.outDir, "gapic_metadata.json")),
			Content: proto.String(g.pt.String()),
		})
	}

	if g.aux.customOp != nil {
		g.reset()
		if err := g.customOperationType(); err != nil {
			return &g.resp, err
		}
		g.commit(filepath.Join(g.opts.outDir, "operations.go"), g.opts.pkgName)
	}

	g.reset()
	if err := g.genAuxFile(); err != nil {
		return &g.resp, err
	}

	return &g.resp, nil
}

// Collects the proto services to generate GAPICs for from the CodeGeneratorRequest.
func (g *generator) collectServices(genReq *pluginpb.CodeGeneratorRequest) (genServs []*descriptorpb.ServiceDescriptorProto) {
	for _, f := range genReq.GetProtoFile() {
		if !strContains(genReq.GetFileToGenerate(), f.GetName()) {
			continue
		}
		if !g.includeMixinInputFile(f.GetName()) {
			continue
		}
		genServs = append(genServs, f.GetService()...)
	}
	return
}

// gen generates client for the given service.
func (g *generator) gen(serv *descriptorpb.ServiceDescriptorProto) error {
	servName := pbinfo.ReduceServName(serv.GetName(), g.opts.pkgName)

	g.clientHook(servName)
	if err := g.clientOptions(serv, servName); err != nil {
		return err
	}
	if err := g.makeClients(serv, servName); err != nil {
		return err
	}

	for _, v := range g.opts.transports {
		switch v {
		case grpc:
			if err := g.genGRPCMethods(serv, servName); err != nil {
				return err
			}
		case rest:
			if err := g.genRESTMethods(serv, servName); err != nil {
				return err
			}
		}
	}

	return g.genOperationBuilders(serv, servName)
}

func (g *generator) getFormattedValue(m *descriptorpb.MethodDescriptorProto, field string, accessor string) (string, error) {
	f := g.lookupField(m.GetInputType(), field)
	value := ""
	// TODO(noahdietz): need to handle []byte for TYPE_BYTES.
	switch f.GetType() {
	case descriptorpb.FieldDescriptorProto_TYPE_STRING:
		value = fmt.Sprintf("url.QueryEscape(%s)", accessor)
	case descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:
		// Double and float are handled the same way.
		fallthrough
	case descriptorpb.FieldDescriptorProto_TYPE_FLOAT:
		// Format the floating point value with mode 'g' to allow for
		// exponent formatting when necessary, and decimal when adequate.
		// QueryEscape the resulting string in case there is a '+' in the
		// exponent.
		// See golang.org/pkg/fmt for more information on formatting.
		value = fmt.Sprintf(`url.QueryEscape(fmt.Sprintf("%%g", %s))`, accessor)
	case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
		en := g.descInfo.Type[f.GetTypeName()]

		n, imp, err := g.descInfo.NameSpec(en)
		if err != nil {
			return "", err
		}
		g.imports[imp] = true

		// protobuf Go generates a mapping from number to string
		// representation of an enum, in UPPER_SNAKE_CASE form. The map
		// is named with the enum name and the _name suffix. If it is a
		// nested enum, the name is prefixed with the parent message name.
		// For example, Severity_name or Error_Severity_name.
		value = fmt.Sprintf("%s.%s_name[int32(%s)]", imp.Name, n, accessor)
	default:
		value = accessor
	}
	return value, nil
}

func (g *generator) insertImplicitRequestHeaders(m *descriptorpb.MethodDescriptorProto, headers [][]string) error {
	if len(headers) == 0 {
		return nil
	}

	seen := map[string]bool{}
	var formats, values strings.Builder
	for _, h := range headers {
		field := h[1]
		// skip fields that have multiple patterns, they use the same accessor
		if _, dupe := seen[field]; dupe {
			continue
		}
		seen[field] = true
		accessor := fmt.Sprintf("req%s", fieldGetter(field))
		accessor, err := g.getFormattedValue(m, field, accessor)
		if err != nil {
			return err
		}
		fmt.Fprintf(&values, " %q, %s,", url.QueryEscape(field), accessor)
		formats.WriteString("%s=%v&")
	}
	// Trim the trailing comma and ampersand symbols.
	f := formats.String()[:formats.Len()-1]
	v := values.String()[:values.Len()-1]
	g.printf("hds := []string{\"x-goog-request-params\", fmt.Sprintf(%q,%s)}", f, v)
	return nil
}

func (g *generator) insertDynamicRequestHeaders(m *descriptorpb.MethodDescriptorProto, headers [][]string) error {
	if len(headers) == 0 {
		return nil
	}

	g.printf(`routingHeaders := ""`)
	g.printf("routingHeadersMap := make(map[string]string)")
	for i := range headers {
		namedCaptureRegex := headers[i][0]
		field := headers[i][1]
		headerName := headers[i][2]
		accessor := fmt.Sprintf("req%s", fieldGetter(field))
		regexHelper := fmt.Sprintf("reg.FindStringSubmatch(%s)[1]", accessor)
		regexHelper, err := g.getFormattedValue(m, field, regexHelper)
		if err != nil {
			return err
		}
		// There could be an edge case where the request field is empty and the path template is a wildcard. In that case, we still don't want to send an empty header name.
		g.printf("if reg := regexp.MustCompile(%q); reg.MatchString(%s) && len(%s) > 0 {", namedCaptureRegex, accessor, regexHelper)
		g.printf("  routingHeadersMap[%q] = %s", headerName, regexHelper)
		g.printf("}")
	}
	g.printf("for headerName, headerValue := range routingHeadersMap {")
	g.printf(`  routingHeaders = fmt.Sprintf("%%s%%s=%%s&", routingHeaders, headerName, headerValue)`)
	g.printf("}")
	g.printf(`routingHeaders = strings.TrimSuffix(routingHeaders, "&")`)
	g.imports[pbinfo.ImportSpec{Path: "strings"}] = true
	g.printf(`hds := []string{"x-goog-request-params", routingHeaders}`)
	g.imports[pbinfo.ImportSpec{Path: "regexp"}] = true
	return nil
}

func (g *generator) insertRequestHeaders(m *descriptorpb.MethodDescriptorProto, t transport) {
	p := g.printf
	// Implicit headers are parsed from the google.api.http annotations and are a default
	// behavior for generators per https://google.aip.dev/client-libraries/4222.
	headers := parseImplicitRequestHeaders(m)
	insertHeaders := g.insertImplicitRequestHeaders

	// If dynamic routing annotation exists, it supercedes injection of implicit headers
	// as described in https://google.aip.dev/client-libraries/4222.
	if dynamicRequestHeadersExist(m) {
		headers = parseDynamicRequestHeaders(m)
		insertHeaders = g.insertDynamicRequestHeaders
	}

	// Has request-based header parameters
	if len(headers) > 0 {
		insertHeaders(m, headers)
		p("")
		switch t {
		case grpc:
			p("hds = append(c.xGoogHeaders, hds...)")
			p("ctx = gax.InsertMetadataIntoOutgoingContext(ctx, hds...)")
		case rest:
			p(`hds = append(c.xGoogHeaders, hds...)`)
			p(`hds = append(hds, "Content-Type", "application/json")`)
			p(`headers := gax.BuildHeaders(ctx, hds...)`)
		}
		g.imports[pbinfo.ImportSpec{Path: "fmt"}] = true
		g.imports[pbinfo.ImportSpec{Path: "net/url"}] = true
		return
	}
	// No request-based header parameters.
	switch t {
	case grpc:
		p("ctx = gax.InsertMetadataIntoOutgoingContext(ctx, c.xGoogHeaders...)")
	case rest:
		p(`hds := append(c.xGoogHeaders, "Content-Type", "application/json")`)
		p(`headers := gax.BuildHeaders(ctx, hds...)`)
	}
}

// insertAutoPopulatedFields generates the conditional initialization of any
// default-value request fields (for the given method) that are specified for
// auto-population by autoPopulatedFields.
//
// If the field value is not equal to default value at the time of sending the
// request, implying it was set by the user, or if the field has explicit
// presence and is set by the user, the field must not be auto-populated by
// the client. Values automatically populated this way must be reused for
// retries of the same request.
func (g *generator) initializeAutoPopulatedFields(servName string, m *descriptorpb.MethodDescriptorProto) {
	apfs := g.autoPopulatedFields(servName, m)
	if len(apfs) == 0 {
		return
	}
	g.imports[pbinfo.ImportSpec{Path: "github.com/google/uuid"}] = true
	p := g.printf
	for _, apf := range apfs {
		f := buildAccessor(apf.GetName(), true)
		if apf.GetProto3Optional() {
			// Type will be *string if field has explicit presence.
			p("if req != nil && req%s == nil {", f)
			p("  req%s = proto.String(uuid.NewString())", f)
		} else {
			// Type will be string if field does not have explicit presence.
			p(`if req != nil && req%s == "" {`, buildAccessor(apf.GetName(), false))
			p("  req%s = uuid.NewString()", f)
		}
		p("}")
	}
}

func buildAccessor(field string, rawFinal bool) string {
	// Corner case if passed the result of strings.Join on an empty slice.
	if field == "" {
		return ""
	}

	var ax strings.Builder
	split := strings.Split(field, ".")
	idx := len(split)
	if rawFinal {
		idx--
	}
	for _, s := range split[:idx] {
		fmt.Fprintf(&ax, ".Get%s()", snakeToCamel(s))
	}
	if rawFinal {
		fmt.Fprintf(&ax, ".%s", snakeToCamel(split[len(split)-1]))
	}
	return ax.String()
}

// Given a chained description for a field in a proto message,
// e.g. squid.mantle.mass_kg
// return the string description of the go expression
// describing idiomatic access to the terminal field
// i.e. .GetSquid().GetMantle().GetMassKg()
//
// This is the normal way to retrieve values.
func fieldGetter(field string) string {
	return buildAccessor(field, false)
}

// Given a chained description for a field in a proto message,
// e.g. squid.mantle.mass_kg
// return the string description of the go expression
// describing direct access to the terminal field
// i.e. .GetSquid().GetMantle().MassKg
//
// This is used for determining field presence for terminal optional fields.
func directAccess(field string) string {
	return buildAccessor(field, true)
}

func (g *generator) lookupField(msgName, field string) *descriptorpb.FieldDescriptorProto {
	var desc *descriptorpb.FieldDescriptorProto
	msg := g.descInfo.Type[msgName]

	// If the message doesn't exist, fail cleanly.
	if msg == nil {
		return desc
	}

	msgProto := msg.(*descriptorpb.DescriptorProto)
	msgFields := msgProto.GetField()

	// Split the key name for nested fields, and traverse the message chain.
	for _, seg := range strings.Split(field, ".") {
		// Look up the desired field by name, stopping if the leaf field is
		// found, continuing if the field is a nested message.
		for _, f := range msgFields {
			if f.GetName() == seg {
				desc = f

				// Search the nested message for the next segment of the
				// nested field chain.
				if f.GetType() == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE {
					msg = g.descInfo.Type[f.GetTypeName()]
					msgProto = msg.(*descriptorpb.DescriptorProto)
					msgFields = msgProto.GetField()
				}
				break
			}
		}
	}
	return desc
}

func (g *generator) appendCallOpts(m *descriptorpb.MethodDescriptorProto) {
	g.printf("opts = append(%[1]s[0:len(%[1]s):len(%[1]s)], opts...)", "(*c.CallOptions)."+*m.Name)
}

// This is a helper function that checks whether a description contains a Deprecated header.
func containsDeprecated(com string) bool {
	for _, s := range strings.Split(com, "\n") {
		if strings.HasPrefix(s, "Deprecated:") {
			return true
		}
	}
	return false
}

func (g *generator) methodDoc(m *descriptorpb.MethodDescriptorProto, serv *descriptorpb.ServiceDescriptorProto) {
	com := g.comments[m]

	// If there's no comment and the method is not deprecated, adding method name is just confusing.
	if !m.GetOptions().GetDeprecated() && com == "" {
		return
	}

	if containsTransport(g.opts.transports, rest) && m.GetClientStreaming() {
		com = fmt.Sprintf("%s\n\nThis method is not supported for the REST transport.", com)
	}
	// If the method is marked as deprecated and there is no comment, then add default deprecation comment.
	// If the method has a comment but it does not include a deprecation notice, then append a default deprecation notice.
	// If the method includes a deprecation notice at the beginning of the comment, prepend a comment stating the method is deprecated and use the included deprecation notice.
	if m.GetOptions().GetDeprecated() {
		if com == "" {
			com = fmt.Sprintf("\n is deprecated.\n\nDeprecated: %s may be removed in a future version.", m.GetName())
		} else if strings.HasPrefix(com, "Deprecated:") {
			com = fmt.Sprintf("\n is deprecated.\n\n%s", com)
		} else if !containsDeprecated(com) {
			com = fmt.Sprintf("%s\n\nDeprecated: %s may be removed in a future version.", com, m.GetName())
		}
	}
	com = strings.TrimSpace(com)

	// Prepend the method name to all non-empty comments.
	com = m.GetName() + " " + lowerFirst(com)
	g.addSnippetsMetadataDoc(m, serv.GetName(), com)
	g.comment(com)
}

// printComment accepts a string to print in code comments and a printFunc
// function to format-print the comment slashes plus the string.
func (g *generator) printComment(s string, p *printer.P) {
	s = strings.TrimSpace(s)
	if s == "" {
		return
	}

	s = mdPlain(s)

	lines := strings.Split(s, "\n")
	for _, l := range lines {
		if strings.TrimSpace(l) == "" {
			p.Printf("//")
		} else {
			p.Printf("// %s", l)
		}
	}
}

// comment format-prints a string in code comments to pt using printf.
func (g *generator) comment(s string) {
	g.printComment(s, &g.pt)
}

// comment format-prints a string in code comments to headerComments.
func (g *generator) headerComment(s string) {
	g.printComment(s, &g.headerComments)
}

// Similar functionality to 'comment', except specifically used for generating formatted code snippets.
func (g *generator) codesnippet(s string) {
	if s == "" {
		return
	}

	lines := strings.Split(s, "\n")
	for _, l := range lines {
		if strings.TrimSpace(l) == "" {
			g.printf("//")
		} else {
			// At least 2 spaces after the // are necessary for GoDoc code snippet formatting.
			g.printf("//  %s", l)
		}
	}
}

// isLRO determines if a given Method is a longrunning operation, ignoring
// those defined by the longrunning proto package.
func (g *generator) isLRO(m *descriptorpb.MethodDescriptorProto) bool {
	return m.GetOutputType() == operationType && g.descInfo.ParentFile[m].GetPackage() != "google.longrunning"
}

func (g *generator) isPaginated(m *descriptorpb.MethodDescriptorProto) bool {
	pf, _, err := g.getPagingFields(m)
	if err != nil {
		return false
	}

	return pf != nil
}

func (g *generator) returnType(m *descriptorpb.MethodDescriptorProto) (string, error) {
	outType := g.descInfo.Type[m.GetOutputType()]
	outSpec, err := g.descInfo.ImportSpec(outType)
	if err != nil {
		return "", err
	}
	info := getHTTPInfo(m)

	// Regular return type.
	retTyp := fmt.Sprintf("*%s.%s", outSpec.Name, outType.GetName())

	// Wrap the raw operation with a custom operation type.
	if g.isCustomOp(m, info) {
		// This will only be *Operation to start.
		retTyp = fmt.Sprintf("*%s", g.aux.customOp.message.GetName())
	}

	return retTyp, nil
}

func parseImplicitRequestHeaders(m *descriptorpb.MethodDescriptorProto) [][]string {
	var matches [][]string

	eHTTP := proto.GetExtension(m.GetOptions(), annotations.E_Http)
	if m == nil || m.GetOptions() == nil {
		return nil
	}

	http := eHTTP.(*annotations.HttpRule)
	rules := []*annotations.HttpRule{http}
	rules = append(rules, http.GetAdditionalBindings()...)

	for _, rule := range rules {
		pattern := ""

		switch rule.GetPattern().(type) {
		case *annotations.HttpRule_Get:
			pattern = rule.GetGet()
		case *annotations.HttpRule_Post:
			pattern = rule.GetPost()
		case *annotations.HttpRule_Patch:
			pattern = rule.GetPatch()
		case *annotations.HttpRule_Put:
			pattern = rule.GetPut()
		case *annotations.HttpRule_Delete:
			pattern = rule.GetDelete()
		}

		matches = append(matches, headerParamRegexp.FindAllStringSubmatch(pattern, -1)...)
	}

	return matches
}

// Determine whether routing annotation exists
func dynamicRequestHeadersExist(m *descriptorpb.MethodDescriptorProto) bool {
	return proto.HasExtension(m.GetOptions(), annotations.E_Routing)
}

// Parse routing annotations to be used as request headers
func parseDynamicRequestHeaders(m *descriptorpb.MethodDescriptorProto) [][]string {
	var matches [][]string

	reqHeaders := proto.GetExtension(m.GetOptions(), annotations.E_Routing)
	routingRule := reqHeaders.(*annotations.RoutingRule)
	routingParameter := routingRule.GetRoutingParameters()

	for _, param := range routingParameter {
		pathTemplateRegex := convertPathTemplateToRegex(param.GetPathTemplate())
		fieldReq := param.Field
		headerName := getHeaderName(param.GetPathTemplate())
		if len(headerName) < 1 {
			headerName = fieldReq
		}
		paramSlice := []string{pathTemplateRegex, fieldReq, headerName}
		matches = append(matches, paramSlice)
	}

	return matches
}
