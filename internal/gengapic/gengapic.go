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
		genServs = append(genServs, f.GetService()...)
	}

	if g.hasIAMPolicyMixin() {
		g.hasIAMPolicyOverrides = hasIAMPolicyOverrides(genServs)
	}

	if g.serviceConfig != nil {
		g.apiName = g.serviceConfig.GetTitle()
	}

	// Use the proto package from the parent file of the first Service seen.
	if len(genServs) > 0 {
		g.metadata.ProtoPackage = g.descInfo.ParentFile[genServs[0]].GetPackage()
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
		if err := g.genExampleFile(s, g.opts.pkgName); err != nil {
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
	g.genDocFile(time.Now().Year(), scopes)
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

	// clear LRO types between services
	g.aux.lros = []*descriptor.MethodDescriptorProto{}

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

	return nil
}

// auxTypes gathers details of types we need to generate along with the client
type auxTypes struct {
	// List of LRO methods. For each method "Foo", we use this to create the "FooOperation" type.
	lros []*descriptor.MethodDescriptorProto

	// "List" of iterator types. We use these to generate FooIterator returned by paging methods.
	// Since multiple methods can page over the same type, we dedupe by the name of the iterator,
	// which is in turn determined by the element type name.
	iters map[string]*iterType
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

func buildAccessor(field string) string {
	var ax strings.Builder
	split := strings.Split(field, ".")
	for _, s := range split {
		fmt.Fprintf(&ax, ".Get%s()", snakeToCamel(s))
	}
	return ax.String()
}

func (g *generator) lookupFieldType(msgName, field string) descriptor.FieldDescriptorProto_Type {
	var typ descriptor.FieldDescriptorProto_Type
	msg := g.descInfo.Type[msgName]
	msgProto := msg.(*descriptor.DescriptorProto)
	msgFields := msgProto.GetField()

	// Split the key name for nested fields, and traverse the message chain.
	for _, seg := range strings.Split(field, ".") {
		// Look up the desired field by name, stopping if the leaf field is
		// found, continuing if the field is a nested message.
		for _, f := range msgFields {
			if f.GetName() == seg {
				typ = f.GetType()

				// Search the nested message for the next segment of the
				// nested field chain.
				if typ == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
					msg = g.descInfo.Type[f.GetTypeName()]
					msgProto = msg.(*descriptor.DescriptorProto)
					msgFields = msgProto.GetField()
				}
				break
			}
		}
	}
	return typ
}

func (g *generator) appendCallOpts(m *descriptor.MethodDescriptorProto) {
	g.printf("opts = append(%[1]s[0:len(%[1]s):len(%[1]s)], opts...)", "(*c.CallOptions)."+*m.Name)
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
