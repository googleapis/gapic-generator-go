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
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/googleapis/gapic-generator-go/internal/errors"
	"github.com/googleapis/gapic-generator-go/internal/license"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"github.com/googleapis/gapic-generator-go/internal/printer"
	"google.golang.org/genproto/googleapis/api/annotations"
)

const (
	// protoc puts a dot in front of name, signaling that the name is fully qualified.
	emptyType  = ".google.protobuf.Empty"
	lroType    = ".google.longrunning.Operation"
	paramError = "need parameter in format: go-gapic-package=client/import/path;packageName"
)

func Gen(genReq *plugin.CodeGeneratorRequest) (*plugin.CodeGeneratorResponse, error) {
	var pkgPath, pkgName, outDir string
	if genReq.Parameter == nil {
		return nil, errors.E(nil, paramError)
	}

	// parse plugin params, ignoring unknown values
	for _, s := range strings.Split(*genReq.Parameter, ",") {
		if e := strings.IndexByte(*genReq.Parameter, '='); e > 0 && s[:e] == "go-gapic-package" {
			p := strings.IndexByte(*genReq.Parameter, ';')

			if p < 0 {
				return nil, errors.E(nil, paramError)
			}

			pkgPath = (*genReq.Parameter)[e+1 : p]
			pkgName = (*genReq.Parameter)[p+1:]
			outDir = filepath.FromSlash(pkgPath)
		}
	}

	if pkgPath == "" || pkgName == "" || outDir == "" {
		return nil, errors.E(nil, paramError)
	}

	var g generator
	g.init(genReq.ProtoFile)

	var genServs []*descriptor.ServiceDescriptorProto
	var eMeta *annotations.Metadata
	for _, f := range genReq.ProtoFile {
		if !strContains(genReq.FileToGenerate, f.GetName()) {
			continue
		}
		genServs = append(genServs, f.Service...)

		// TODO(pongad): check if first-one-wins is the right strategy here.
		if eMeta == nil {
			if em, err := proto.GetExtension(f.GetOptions(), annotations.E_Metadata); err == nil {
				eMeta = em.(*annotations.Metadata)
			}
		}
	}
	if eMeta != nil {
		// Without this, the doc is going to be a little bad but this is not an error.
		nameParts := append([]string(nil), eMeta.PackageNamespace...)
		nameParts = append(nameParts, eMeta.ProductName)
		g.apiName = strings.Join(nameParts, " ")
	}

	for _, s := range genServs {
		// TODO(pongad): gapic-generator does not remove the package name here,
		// so even though the client for LoggingServiceV2 is just "Client"
		// the file name is "logging_client.go".
		// Keep the current behavior for now, but we could revisit this later.
		outFile := pbinfo.ReduceServName(s.GetName(), "")
		outFile = camelToSnake(outFile)
		outFile = filepath.Join(outDir, outFile)

		g.reset()
		if err := g.gen(s, pkgName); err != nil {
			return nil, errors.E(err, "service: %s", s.GetName())
		}
		g.commit(outFile+"_client.go", pkgName)

		g.reset()
		if err := g.genExampleFile(s, pkgName); err != nil {
			return nil, errors.E(err, "example: %s", s.GetName())
		}
		g.imports[pbinfo.ImportSpec{Path: pkgPath}] = true
		g.commit(outFile+"_client_example_test.go", pkgName+"_test")
	}

	g.reset()
	scopes, err := collectScopes(genServs)
	if err != nil {
		return nil, err
	}
	g.genDocFile(pkgPath, pkgName, time.Now().Year(), scopes)
	g.resp.File = append(g.resp.File, &plugin.CodeGeneratorResponse_File{
		Name:    proto.String(filepath.Join(outDir, "doc.go")),
		Content: proto.String(g.pt.String()),
	})

	return &g.resp, nil
}

func strContains(a []string, s string) bool {
	for _, as := range a {
		if as == s {
			return true
		}
	}
	return false
}

type generator struct {
	pt printer.P

	descInfo pbinfo.Info

	// Maps proto elements to their comments
	comments map[proto.Message]string

	resp plugin.CodeGeneratorResponse

	imports map[pbinfo.ImportSpec]bool

	// Human-readable name of the API used in docs
	apiName string
}

func (g *generator) init(files []*descriptor.FileDescriptorProto) {
	g.descInfo = pbinfo.Of(files)

	g.comments = map[proto.Message]string{}
	g.imports = map[pbinfo.ImportSpec]bool{}

	for _, f := range files {
		for _, loc := range f.GetSourceCodeInfo().GetLocation() {
			if loc.LeadingComments == nil {
				continue
			}

			// p is an array with format [f1, i1, f2, i2, ...]
			// - f1 refers to the protobuf field tag
			// - if field refer to by f1 is a slice, i1 refers to an element in that slice
			// - f2 and i2 works recursively.
			// So, [6, x] refers to the xth service defined in the file,
			// since the field tag of Service is 6.
			// [6, x, 2, y] refers to the yth method in that service,
			// since the field tag of Method is 2.
			p := loc.Path
			switch {
			case len(p) == 2 && p[0] == 6:
				g.comments[f.Service[p[1]]] = *loc.LeadingComments
			case len(p) == 4 && p[0] == 6 && p[2] == 2:
				g.comments[f.Service[p[1]].Method[p[3]]] = *loc.LeadingComments
			}
		}
	}
}

// printf formatted-prints to sb, using the print syntax from fmt package.
//
// It automatically keeps track of indentation caused by curly-braces.
// To make nested blocks easier to write elsewhere in the code,
// leading and trailing whitespaces in s are ignored.
// These spaces are for humans reading the code, not machines.
//
// Currently it's not terribly difficult to confuse the auto-indenter.
// To fix-up, manipulate g.in or write to g.sb directly.
func (g *generator) printf(s string, a ...interface{}) {
	g.pt.Printf(s, a...)
}

func (g *generator) commit(fileName, pkgName string) {
	var header strings.Builder
	fmt.Fprintf(&header, license.Apache, time.Now().Year())
	fmt.Fprintf(&header, "package %s\n\n", pkgName)

	var imps []pbinfo.ImportSpec
	for imp := range g.imports {
		imps = append(imps, imp)
	}
	impDiv := sortImports(imps)

	writeImp := func(is pbinfo.ImportSpec) {
		s := "\t%[2]q\n"
		if is.Name != "" {
			s = "\t%s %q\n"
		}
		fmt.Fprintf(&header, s, is.Name, is.Path)
	}

	header.WriteString("import (\n")
	for _, imp := range imps[:impDiv] {
		writeImp(imp)
	}
	if impDiv != 0 && impDiv != len(imps) {
		header.WriteByte('\n')
	}
	for _, imp := range imps[impDiv:] {
		writeImp(imp)
	}
	header.WriteString(")\n\n")

	g.resp.File = append(g.resp.File, &plugin.CodeGeneratorResponse_File{
		Name:    &fileName,
		Content: proto.String(header.String()),
	})

	// Trim trailing newlines so we have only one.
	// NOTE(pongad): This might be an overkill since we have gofmt,
	// but the rest of the file already conforms to gofmt, so we might as well?
	body := g.pt.String()
	if !strings.HasSuffix(body, "\n") {
		body += "\n"
	}
	for i := len(body) - 1; i >= 0; i-- {
		if body[i] != '\n' {
			body = body[:i+2]
			break
		}
	}

	g.resp.File = append(g.resp.File, &plugin.CodeGeneratorResponse_File{
		Content: proto.String(body),
	})
}

func (g *generator) reset() {
	g.pt.Reset()
	for k := range g.imports {
		delete(g.imports, k)
	}
}

// gen generates client for the given service.
func (g *generator) gen(serv *descriptor.ServiceDescriptorProto, pkgName string) error {
	servName := pbinfo.ReduceServName(*serv.Name, pkgName)
	if err := g.clientOptions(serv, servName); err != nil {
		return err
	}
	if err := g.clientInit(serv, servName); err != nil {
		return err
	}

	aux := auxTypes{
		iters: map[string]iterType{},
	}
	for _, m := range serv.Method {
		g.methodDoc(m)
		if err := g.genMethod(servName, serv, m, &aux); err != nil {
			return errors.E(err, "method: %s", m.GetName())
		}
	}

	sort.Slice(aux.lros, func(i, j int) bool {
		return aux.lros[i].GetName() < aux.lros[j].GetName()
	})
	for _, m := range aux.lros {
		if err := g.lroType(servName, serv, m); err != nil {
			return errors.E(err, "while generating LRO type for %q", m.GetName())
		}
	}

	var iters []iterType
	for _, iter := range aux.iters {
		iters = append(iters, iter)
	}
	sort.Slice(iters, func(i, j int) bool {
		return iters[i].iterTypeName < iters[j].iterTypeName
	})
	for _, iter := range iters {
		g.pagingIter(iter)
	}

	return nil
}

// auxTypes gathers details of types we need to generate along with the client
type auxTypes struct {
	// List of LRO methods. For each method "Foo", we use this to create the "FooOperation" type.
	lros []*descriptor.MethodDescriptorProto

	// "List" of iterator types. We use these to generate FooIterator returned by paging methods.
	// Since multiple methods can page over the same type, we dedupe by the name of the iterator,
	// which is in turn determined by the element type name.
	iters map[string]iterType
}

// genMethod generates a single method from a client. m must be a method declared in serv.
// If the generated method requires an auxillary type, it is added to aux.
func (g *generator) genMethod(servName string, serv *descriptor.ServiceDescriptorProto, m *descriptor.MethodDescriptorProto, aux *auxTypes) error {
	if m.GetOutputType() == lroType {
		aux.lros = append(aux.lros, m)
		return g.lroCall(servName, m)
	}

	if m.GetOutputType() == emptyType {
		return g.emptyUnaryCall(servName, m)
	}

	if pf, err := g.pagingField(m); err != nil {
		return err
	} else if pf != nil {
		iter, err := g.iterTypeOf(pf)
		if err != nil {
			return err
		}
		aux.iters[iter.iterTypeName] = iter
		return g.pagingCall(servName, m, pf, iter)
	}

	switch {
	case m.GetClientStreaming():
		return g.noRequestStreamCall(servName, serv, m)
	case m.GetServerStreaming():
		return g.serverStreamCall(servName, serv, m)
	default:
		return g.unaryCall(servName, m)
	}
}

func (g *generator) unaryCall(servName string, m *descriptor.MethodDescriptorProto) error {
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

	p("func (c *%sClient) %s(ctx context.Context, req *%s.%s, opts ...gax.CallOption) (*%s.%s, error) {",
		servName, *m.Name, inSpec.Name, inType.GetName(), outSpec.Name, outType.GetName())

	g.insertMetadata()
	g.appendCallOpts(m)
	p("var resp *%s.%s", outSpec.Name, outType.GetName())
	p("err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {")
	p("  var err error")
	p("  resp, err = %s", grpcClientCall(servName, *m.Name))
	p("  return err")
	p("}, opts...)")
	p("if err != nil {")
	p("  return nil, err")
	p("}")
	p("return resp, nil")

	p("}")
	p("")

	g.imports[inSpec] = true
	g.imports[outSpec] = true

	return nil
}

func (g *generator) emptyUnaryCall(servName string, m *descriptor.MethodDescriptorProto) error {
	inType := g.descInfo.Type[*m.InputType]

	inSpec, err := g.descInfo.ImportSpec(inType)
	if err != nil {
		return err
	}

	p := g.printf

	p("func (c *%sClient) %s(ctx context.Context, req *%s.%s, opts ...gax.CallOption) error {",
		servName, m.GetName(), inSpec.Name, inType.GetName())

	g.insertMetadata()
	g.appendCallOpts(m)
	p("err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {")
	p("  var err error")
	p("  _, err = %s", grpcClientCall(servName, m.GetName()))
	p("  return err")
	p("}, opts...)")
	p("return err")

	p("}")
	p("")

	g.imports[inSpec] = true
	return nil
}

func (g *generator) insertMetadata() {
	g.printf("ctx = insertMetadata(ctx, c.xGoogMetadata)")
}

func (g *generator) appendCallOpts(m *descriptor.MethodDescriptorProto) {
	g.printf("opts = append(%[1]s[0:len(%[1]s):len(%[1]s)], opts...)", "c.CallOptions."+*m.Name)
}

func (g *generator) methodDoc(m *descriptor.MethodDescriptorProto) {
	com := g.comments[m]
	com = strings.TrimSpace(com)

	// If there's no comment, adding method name is just confusing.
	if com == "" {
		return
	}

	g.comment(*m.Name + " " + lowerFirst(com))
}

func (g *generator) comment(s string) {
	s = strings.TrimSpace(s)
	if s == "" {
		return
	}

	s = MDPlain(s)

	lines := strings.Split(s, "\n")
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" {
			g.printf("//")
		} else {
			g.printf("// %s", l)
		}
	}
}

// grpcClientField reports the field name to store gRPC client.
func grpcClientField(reducedServName string) string {
	// Not the same as pbinfo.ReduceServName(*serv.Name, pkg)+"Client".
	// If the service name is reduced to empty string, we should
	// lower-case "client" so that the field is not exported.
	return lowerFirst(reducedServName + "Client")
}

func grpcClientCall(reducedServName, methName string) string {
	return fmt.Sprintf("c.%s.%s(ctx, req, settings.GRPC...)", grpcClientField(reducedServName), methName)
}

func lowerFirst(s string) string {
	if s == "" {
		return ""
	}
	r, w := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[w:]
}

func upperFirst(s string) string {
	if s == "" {
		return ""
	}
	r, w := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[w:]
}

func camelToSnake(s string) string {
	var sb strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) && i != 0 {
			sb.WriteByte('_')
		}
		sb.WriteRune(unicode.ToLower(r))
	}
	return sb.String()
}

// snakeToCamel converts snake_case and SNAKE_CASE to CamelCase.
func snakeToCamel(s string) string {
	var sb strings.Builder
	up := true
	for _, r := range s {
		if r == '_' {
			up = true
		} else if up {
			sb.WriteRune(unicode.ToUpper(r))
			up = false
		} else {
			sb.WriteRune(unicode.ToLower(r))
		}
	}
	return sb.String()
}
