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

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

var tabsCache = strings.Repeat("\t", 20)
var spacesCache = strings.Repeat(" ", 100)

func main() {
	reqBytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	var genReq plugin.CodeGeneratorRequest
	if err := genReq.Unmarshal(reqBytes); err != nil {
		log.Fatal(err)
	}

	outDir := ""
	if p := genReq.Parameter; p != nil {
		outDir = *p
	}

	var g generator
	g.init(genReq.ProtoFile)
	for _, f := range genReq.ProtoFile {
		if strContains(genReq.FileToGenerate, *f.Name) {
			for _, s := range f.Service {
				g.gen(s)
				g.commit(filepath.Join(outDir, camelToSnake(*s.Name)+"_client.go"))
			}
		}
	}

	outBytes, err := proto.Marshal(&g.resp)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stdout.Write(outBytes); err != nil {
		log.Fatal(err)
	}
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
	sb strings.Builder

	// current indentation level
	in int

	resp plugin.CodeGeneratorResponse

	// Maps services and messages to the file containing them,
	// so we can figure out the import.
	parentFile map[proto.Message]*descriptor.FileDescriptorProto

	// Maps type names to their messages
	types map[string]*descriptor.DescriptorProto

	// Maps proto elements to their comments
	comments map[proto.Message]string
}

func (g *generator) init(files []*descriptor.FileDescriptorProto) {
	g.parentFile = map[proto.Message]*descriptor.FileDescriptorProto{}
	g.types = map[string]*descriptor.DescriptorProto{}
	g.comments = map[proto.Message]string{}

	for _, f := range files {
		// parentFile
		for _, m := range f.MessageType {
			g.parentFile[m] = f
		}
		for _, s := range f.Service {
			g.parentFile[s] = f
		}

		// types
		for _, m := range f.MessageType {
			fullyQualifiedName := fmt.Sprintf(".%s.%s", *f.Package, *m.Name)
			g.types[fullyQualifiedName] = m
		}

		// comment
		for _, loc := range f.SourceCodeInfo.Location {
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

// pkgName reports the package name of protobuf element e.
func (g *generator) pkgName(e proto.Message) string {
	fdesc := g.parentFile[e]
	pkg := *fdesc.Options.GoPackage
	if p := strings.IndexByte(pkg, ';'); p >= 0 {
		return pkg[p+1:] + "pb"
	}

	for {
		p := strings.LastIndexByte(pkg, '/')
		if p < 0 {
			return pkg + "pb"
		}
		elem := pkg[p+1:]
		if len(elem) >= 2 && elem[0] == 'v' && elem[1] >= '0' && elem[1] <= '9' {
			// It's a version number; skip so we get a more meaningful name
			pkg = pkg[:p]
			continue
		}
		return elem + "pb"
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
	s = strings.TrimSpace(s)
	if s == "" {
		g.sb.WriteByte('\n')
		return
	}

	for i := 0; i < len(s) && s[i] == '}'; i++ {
		g.in--
	}

	in := g.in
	for in > len(tabsCache) {
		g.sb.WriteString(tabsCache)
		in -= len(tabsCache)
	}
	g.sb.WriteString(tabsCache[:in])

	fmt.Fprintf(&g.sb, s, a...)
	g.sb.WriteByte('\n')

	for i := len(s) - 1; i >= 0 && s[i] == '{'; i-- {
		g.in++
	}
}

func (g *generator) commit(fileName string) {
	const license = `// Copyright %d Google LLC
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

// AUTO-GENERATED CODE. DO NOT EDIT.

`

	var header strings.Builder
	fmt.Fprintf(&header, license, time.Now().Year())
	// TODO(pongad): imports
	// TODO(pongad): read package name from somewhere
	header.WriteString("package foo\n\n")

	g.resp.File = append(g.resp.File, &plugin.CodeGeneratorResponse_File{
		Name:    &fileName,
		Content: proto.String(header.String()),
	})
	g.resp.File = append(g.resp.File, &plugin.CodeGeneratorResponse_File{
		Content: proto.String(g.sb.String()),
	})
}

func (g *generator) gen(serv *descriptor.ServiceDescriptorProto) {
	g.sb.Reset()
	g.in = 0

	p := g.printf

	// TODO(pongad): maybe individual sections should be their own methods.

	// CallOptions struct
	{
		var maxNameLen int
		for _, m := range serv.Method {
			if l := len(*m.Name); maxNameLen < l {
				maxNameLen = l
			}
		}

		p("// %[1]sCallOptions contains the retry settings for each method of %[1]sClient.", *serv.Name)
		p("type %sCallOptions struct {", *serv.Name)
		for _, m := range serv.Method {
			p("%s%s[]gax.CallOption", *m.Name, spaces(maxNameLen-len(*m.Name)+1))
		}
		p("}")
		p("")
	}

	// defaultClientOptions
	{
		// TODO(pongad): read URL from somewhere
		p("func default%sClientOptions() []option.ClientOption {", *serv.Name)
		p("  return []option.ClientOption{")
		p(`    option.WithEndpoint("foo.googleapis.com:443"),`)
		p("    option.WithScopes(DefaultAuthScopes()...),")
		p("  }")
		p("}")
		p("")
	}

	// defaultCallOptions
	{
		// TODO(pongad): read retry params from somewhere
		p("func default%[1]sCallOptions() *%[1]sCallOptions {", *serv.Name)
		p("  return &%sCallOptions{", *serv.Name)
		p("  }")
		p("}")
		p("")
	}

	// client struct
	{
		// TODO(pongad): read "human" API name from somewhere
		p("// %sClient is a client for interacting with Foo API.", *serv.Name)
		p("//")
		p("// Methods, except Close, may be called concurrently. However, fields must not be modified concurrently with method calls.")
		p("type %sClient struct {", *serv.Name)

		p("// The connection to the service.")
		p("conn *grpc.ClientConn")
		p("")

		p("// The gRPC API client.")
		p("%sClient %s.%sClient", lowerFirst(*serv.Name), g.pkgName(serv), *serv.Name)
		p("")

		p("// LROClient is used internally to handle longrunning operations.")
		p("// It is exposed so that its CallOptions can be modified if required.")
		p("// Users should not Close this client.")
		p("LROClient *lroauto.OperationsClient")
		p("")

		p("// The call options for this service.")
		p("CallOptions *%sCallOptions", *serv.Name)
		p("")

		p("// The x-goog-* metadata to be sent with each request.")
		p("xGoogMetadata metadata.MD")
		p("}")
		p("")
	}

	// Client constructor
	{
		// TODO(pongad): client name
		p("// New%sClient creates a new foo client.", *serv.Name)
		p("//")
		g.comment(g.comments[serv])
		p("func New%[1]sClient(ctx context.Context, opts ...option.ClientOption) (*%[1]sClient, error) {", *serv.Name)
		p("  conn, err := transport.DialGRPC(ctx, append(default%sClientOptions(), opts...)...)", *serv.Name)
		p("  if err != nil {")
		p("    return nil, err")
		p("  }")
		p("  c := &%sClient{", *serv.Name)
		p("    conn:        conn,")
		p("    CallOptions: default%sCallOptions(),", *serv.Name)
		p("")
		p("    %sClient: %s.New%sClient(conn),", lowerFirst(*serv.Name), g.pkgName(serv), *serv.Name)
		p("  }")
		p("  c.setGoogleClientInfo()")
		p("")
		p("  c.LROClient, err = lroauto.NewOperationsClient(ctx, option.WithGRPCConn(conn))")
		p("  if err != nil {")
		p("    // This error \"should not happen\", since we are just reusing old connection")
		p("    // and never actually need to dial.")
		p("    // If this does happen, we could leak conn. However, we cannot close conn:")
		p("    // If the user invoked the function with option.WithGRPCConn,")
		p("    // we would close a connection that's still in use.")
		p("    // TODO(pongad): investigate error conditions.")
		p("    return nil, err")
		p("  }")
		p("  return c, nil")
		p("}")
		p("")
	}

	// Connection()
	{
		p("// Connection returns the client's connection to the API service.")
		p("func (c *%sClient) Connection() *grpc.ClientConn {", *serv.Name)
		p("  return c.conn")
		p("}")
		p("")
	}

	// Close()
	{
		p("// Close closes the connection to the API service. The user should invoke this when")
		p("// the client is no longer required.")
		p("func (c *%sClient) Close() error {", *serv.Name)
		p("  return c.conn.Close()")
		p("}")
		p("")
	}

	// setGoogleClientInfo
	{
		p("// setGoogleClientInfo sets the name and version of the application in")
		p("// the `x-goog-api-client` header passed on each request. Intended for")
		p("// use by Google-written clients.")
		p("func (c *%sClient) setGoogleClientInfo(keyval ...string) {", *serv.Name)
		p(`  kv := append([]string{"gl-go", version.Go()}, keyval...)`)
		p(`  kv = append(kv, "gapic", version.Repo, "gax", gax.Version, "grpc", grpc.Version)`)
		p(`  c.xGoogMetadata = metadata.Pairs("x-goog-api-client", gax.XGoogHeader(kv...))`)
		p("}")
		p("")
	}

	for _, m := range serv.Method {
		g.unaryCall(*serv.Name, m)
	}
}

func (g *generator) unaryCall(servName string, m *descriptor.MethodDescriptorProto) {
	inType := g.types[*m.InputType]
	outType := g.types[*m.OutputType]
	p := g.printf

	g.methodDoc(m)
	p("func (c *%sClient) %s(ctx context.Context, req *%s.%s, opts ...gax.CallOption) (*%s.%s, error) {",
		servName, *m.Name, g.pkgName(inType), *inType.Name, g.pkgName(outType), *outType.Name)

	p("ctx = insertMetadata(ctx, c.xGoogMetadata)")
	p("opts = append(%[1]s[0:len(%[1]s):len(%[1]s)], opts...)", "c.CallOptions."+*m.Name)
	p("var resp *%s.%s", g.pkgName(outType), *outType.Name)
	p("err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {")
	p("  var err error")
	p("  resp, err = c.%sClient.%s(ctx, req, settings.GRPC...)", lowerFirst(servName), *m.Name)
	p("  return err")
	p("}, opts...)")
	p("if err != nil {")
	p("  return nil, err")
	p("}")
	p("return resp, nil")

	p("}")
	p("")
}

// TODO(pongad): escape markdown
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
	lines := strings.Split(s, "\n")
	for _, l := range lines {
		g.printf("// %s", strings.TrimSpace(l))
	}
}

func spaces(n int) string {
	if n > len(spacesCache) {
		return strings.Repeat(" ", n)
	}
	return spacesCache[:n]
}

func lowerFirst(s string) string {
	if s == "" {
		return ""
	}
	r, w := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[w:]
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
