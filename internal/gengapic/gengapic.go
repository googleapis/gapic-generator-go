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
	"sort"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/googleapis/gapic-generator-go/internal/errors"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"google.golang.org/genproto/googleapis/api/annotations"
)

const (
	emptyValue = "google.protobuf.Empty"
	// protoc puts a dot in front of name, signaling that the name is fully qualified.
	emptyType           = "." + emptyValue
	lroType             = ".google.longrunning.Operation"
	alpha               = "alpha"
	beta                = "beta"
	disableDeadlinesVar = "GOOGLE_API_GO_EXPERIMENTAL_DISABLE_DEFAULT_DEADLINE"
	fieldTypeBool       = descriptor.FieldDescriptorProto_TYPE_BOOL
	fieldTypeString     = descriptor.FieldDescriptorProto_TYPE_STRING
	fieldTypeBytes      = descriptor.FieldDescriptorProto_TYPE_BYTES
	fieldTypeMessage    = descriptor.FieldDescriptorProto_TYPE_MESSAGE
	fieldLabelRepeated  = descriptor.FieldDescriptorProto_LABEL_REPEATED
)

var headerParamRegexp = regexp.MustCompile(`{([_.a-z0-9]+)`)

// Gen is the entry point for GAPIC generation via the protoc plugin.
func Gen(genReq *plugin.CodeGeneratorRequest) (*plugin.CodeGeneratorResponse, error) {
	var g generator
	if err := g.init(genReq); err != nil {
		return &g.resp, err
	}

	var genServs []*descriptor.ServiceDescriptorProto
	for _, f := range genReq.GetProtoFile() {
		if !strContains(genReq.GetFileToGenerate(), f.GetName()) {
			continue
		}
		if !g.includeMixinInputFile(f.GetName()) {
			continue
		}
		genServs = append(genServs, f.GetService()...)
	}

	g.checkIAMPolicyOverrides(genServs)

	if g.serviceConfig != nil {
		g.apiName = g.serviceConfig.GetTitle()
	}

	protoPkg := g.descInfo.ParentFile[genServs[0]].GetPackage()

	if op, ok := g.descInfo.Type[fmt.Sprintf(".%s.Operation", protoPkg)]; g.opts.diregapic && ok {
		g.aux.customOp = &customOp{op.(*descriptor.DescriptorProto), false}
	}

	// Use the proto package from the parent file of the first Service seen.
	if len(genServs) > 0 {
		g.metadata.ProtoPackage = protoPkg
	}
	g.metadata.LibraryPackage = g.opts.pkgPath

	for _, s := range genServs {
		// TODO(pongad): gapic-generator does not remove the package name here,
		// so even though the client for LoggingServiceV2 is just "Client"
		// the file name is "logging_client.go".
		// Keep the current behavior for now, but we could revisit this later.
		outFile := pbinfo.ReduceServName(s.GetName(), "")
		outFile = camelToSnake(outFile)
		outFile = filepath.Join(g.opts.outDir, outFile)

		g.reset()
		if err := g.gen(s); err != nil {
			return &g.resp, err
		}
		g.commit(outFile+"_client.go", g.opts.pkgName)

		g.reset()
		if err := g.genExampleFile(s); err != nil {
			return &g.resp, errors.E(err, "example: %s", s.GetName())
		}
		g.imports[pbinfo.ImportSpec{Name: g.opts.pkgName, Path: g.opts.pkgPath}] = true
		g.commit(outFile+"_client_example_test.go", g.opts.pkgName+"_test")
	}

	g.reset()
	scopes, err := collectScopes(genServs)
	if err != nil {
		return &g.resp, err
	}

	serv := genServs[0]

	g.genDocFile(time.Now().Year(), scopes, serv)
	g.resp.File = append(g.resp.File, &plugin.CodeGeneratorResponse_File{
		Name:    proto.String(filepath.Join(g.opts.outDir, "doc.go")),
		Content: proto.String(g.pt.String()),
	})

	if g.opts.metadata {
		g.reset()
		g.genGapicMetadataFile()
		g.resp.File = append(g.resp.File, &plugin.CodeGeneratorResponse_File{
			Name:    proto.String(filepath.Join(g.opts.outDir, "gapic_metadata.json")),
			Content: proto.String(g.pt.String()),
		})
	}

	return &g.resp, nil
}

// gen generates client for the given service.
func (g *generator) gen(serv *descriptor.ServiceDescriptorProto) error {
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

	// g.aux.lros is a map (set)
	// so that we generate types at most once,
	// but we want a deterministic order to prevent
	// spurious regenerations.
	var lros []*descriptor.MethodDescriptorProto
	for m := range g.aux.lros {
		lros = append(lros, m)
	}
	sort.Slice(lros, func(i, j int) bool {
		return lros[i].GetName() < lros[j].GetName()
	})
	for _, m := range lros {
		if err := g.lroType(servName, serv, m); err != nil {
			return err
		}
	}

	// clear LRO types between services
	g.aux.lros = map[*descriptor.MethodDescriptorProto]bool{}

	var iters []*iterType
	for _, iter := range g.aux.iters {
		// skip iterators that have already been generated in this package
		//
		// TODO(ndietz): investigate generating auxiliary types in a
		// separate file in the same package to avoid keeping this state
		if iter.generated {
			continue
		}

		iter.generated = true
		iters = append(iters, iter)
	}
	sort.Slice(iters, func(i, j int) bool {
		return iters[i].iterTypeName < iters[j].iterTypeName
	})
	for _, iter := range iters {
		g.pagingIter(iter)
	}
	if g.aux.customOp != nil && !g.aux.customOp.generated {
		if err := g.customOperationType(); err != nil {
			return err
		}
		g.aux.customOp.generated = true
	}

	return nil
}

// auxTypes gathers details of types we need to generate along with the client
type auxTypes struct {
	// List of LRO methods. For each method "Foo", we use this to create the "FooOperation" type.
	lros map[*descriptor.MethodDescriptorProto]bool

	// "List" of iterator types. We use these to generate FooIterator returned by paging methods.
	// Since multiple methods can page over the same type, we dedupe by the name of the iterator,
	// which is in turn determined by the element type name.
	iters map[string]*iterType

	customOp *customOp
}

func (g *generator) deadline(s, m string) {
	t, ok := g.grpcConf.Timeout(s, m)
	if !ok {
		return
	}

	g.printf("if _, ok := ctx.Deadline(); !ok && !c.disableDeadlines {")
	g.printf("  cctx, cancel := context.WithTimeout(ctx, %d * time.Millisecond)", t)
	g.printf("  defer cancel()")
	g.printf("  ctx = cctx")
	g.printf("}")

	g.imports[pbinfo.ImportSpec{Path: "time"}] = true
}

func (g *generator) insertMetadata(m *descriptor.MethodDescriptorProto) error {
	headers, err := parseRequestHeaders(m)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
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
			f := g.lookupField(m.GetInputType(), field)

			// TODO(noahdietz): need to handle []byte for TYPE_BYTES.
			switch f.GetType() {
			case descriptor.FieldDescriptorProto_TYPE_STRING:
				accessor = fmt.Sprintf("url.QueryEscape(%s)", accessor)
			case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
				// Double and float are handled the same way.
				fallthrough
			case descriptor.FieldDescriptorProto_TYPE_FLOAT:
				// Format the floating point value with mode 'g' to allow for
				// exponent formatting when necessary, and decimal when adequate.
				// QueryEscape the resulting string in case there is a '+' in the
				// exponent.
				// See golang.org/pkg/fmt for more information on formatting.
				accessor = fmt.Sprintf(`url.QueryEscape(fmt.Sprintf("%%g", %s))`, accessor)
			case descriptor.FieldDescriptorProto_TYPE_ENUM:
				en := g.descInfo.Type[f.GetTypeName()]

				n, imp, err := g.descInfo.NameSpec(en)
				if err != nil {
					return err
				}
				g.imports[imp] = true

				// protobuf Go generates a mapping from number to string
				// representation of an enum, in UPPER_SNAKE_CASE form. The map
				// is named with the enum name and the _name suffix. If it is a
				// nested enum, the name is prefixed with the parent message name.
				// For example, Severity_name or Error_Severity_name.
				accessor = fmt.Sprintf("%s.%s_name[int32(%s)]", imp.Name, n, accessor)
			}

			// URL encode key & values separately per aip.dev/4222.
			// Encode the key ahead of time to reduce clutter
			// and because it will likely never be necessary
			fmt.Fprintf(&values, " %q, %s,", url.QueryEscape(field), accessor)
			formats.WriteString("%s=%v&")
		}
		// Trim the trailing comma and ampersand symbols.
		f := formats.String()[:formats.Len()-1]
		v := values.String()[:values.Len()-1]

		g.printf("md := metadata.Pairs(\"x-goog-request-params\", fmt.Sprintf(%q,%s))", f, v)
		g.printf("ctx = insertMetadata(ctx, c.xGoogMetadata, md)")

		g.imports[pbinfo.ImportSpec{Path: "fmt"}] = true
		g.imports[pbinfo.ImportSpec{Path: "net/url"}] = true

		return nil
	}

	g.printf("ctx = insertMetadata(ctx, c.xGoogMetadata)")

	return nil
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

func (g *generator) lookupField(msgName, field string) *descriptor.FieldDescriptorProto {
	var desc *descriptor.FieldDescriptorProto
	msg := g.descInfo.Type[msgName]

	// If the message doesn't exist, fail cleanly.
	if msg == nil {
		return desc
	}

	msgProto := msg.(*descriptor.DescriptorProto)
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
				if f.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
					msg = g.descInfo.Type[f.GetTypeName()]
					msgProto = msg.(*descriptor.DescriptorProto)
					msgFields = msgProto.GetField()
				}
				break
			}
		}
	}
	return desc
}

func (g *generator) appendCallOpts(m *descriptor.MethodDescriptorProto) {
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

func (g *generator) methodDoc(m *descriptor.MethodDescriptorProto) {
	com := g.comments[m]

	// If there's no comment and the method is not deprecated, adding method name is just confusing.
	if !m.GetOptions().GetDeprecated() && com == "" {
		return
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

	g.comment(com)
}

func (g *generator) comment(s string) {
	s = strings.TrimSpace(s)
	if s == "" {
		return
	}

	s = mdPlain(s)

	lines := strings.Split(s, "\n")
	for _, l := range lines {
		if strings.TrimSpace(l) == "" {
			g.printf("//")
		} else {
			g.printf("// %s", l)
		}
	}
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
func (g *generator) isLRO(m *descriptor.MethodDescriptorProto) bool {
	return m.GetOutputType() == lroType && g.descInfo.ParentFile[m].GetPackage() != "google.longrunning"
}

func (g *generator) getServiceName(m *descriptor.MethodDescriptorProto) string {
	f := g.descInfo.ParentFile[m].GetPackage()
	s := g.descInfo.ParentElement[m].GetName()
	return fmt.Sprintf("%s.%s", f, s)
}

func (g *generator) returnType(m *descriptor.MethodDescriptorProto) (string, error) {
	outType := g.descInfo.Type[m.GetOutputType()]
	outSpec, err := g.descInfo.ImportSpec(outType)
	if err != nil {
		return "", err
	}
	info, err := getHTTPInfo(m)
	if err != nil {
		return "", err
	}

	// Regular return type.
	retTyp := fmt.Sprintf("*%s.%s", outSpec.Name, outType.GetName())

	// Wrap the raw operation with a custom operation type.
	if g.isCustomOp(m, info) {
		// This will only be *Operation to start.
		retTyp = fmt.Sprintf("*%s", g.aux.customOp.message.GetName())
	}

	return retTyp, nil
}

func parseRequestHeaders(m *descriptor.MethodDescriptorProto) ([][]string, error) {
	var matches [][]string

	eHTTP, err := proto.GetExtension(m.GetOptions(), annotations.E_Http)
	if m == nil || m.GetOptions() == nil || err == proto.ErrMissingExtension {
		return nil, nil
	} else if err != nil {
		return nil, err
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

	return matches, nil
}
