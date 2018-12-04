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

package gencli

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"github.com/googleapis/gapic-generator-go/internal/printer"
	goimport "golang.org/x/tools/imports"
)

const (
	// EmptyProtoType is the type name for the Empty message type
	EmptyProtoType = ".google.protobuf.Empty"
	// LROProtoType is the type name for the LRO message type
	LROProtoType = ".google.longrunning.Operation"
	// OutputOnlyStr represents the comment-level string that
	// indicates a field is only present as output, never input
	OutputOnlyStr = "Output only"
	// RequiredStr represents the comment-level string that
	// indicates a field is required
	RequiredStr = "Required"
)

// Command intermediate representation of a RPC/Method as a CLI command
type Command struct {
	Service           string
	Method            string
	MethodCmd         string
	InputMessageType  string
	InputMessage      string
	ShortDesc         string
	LongDesc          string
	Imports           map[string]*pbinfo.ImportSpec
	Flags             []*Flag
	OneOfSelectors    map[string]*Flag
	NestedMessages    []*NestedMessage
	EnvPrefix         string
	OutputMessageType string
	ServerStreaming   bool
	ClientStreaming   bool
	Paged             bool
	IsLRO             bool
	HasEnums          bool
	SubCommands       []*Command
}

// NestedMessage represents a nested message that will need to be initialized
// in the generated code
type NestedMessage struct {
	FieldName string
	FieldType string
}

// Gen is the main entry point for code generation of a command line utility
func Gen(genReq *plugin.CodeGeneratorRequest) (*plugin.CodeGeneratorResponse, error) {
	var g gcli

	err := g.init(genReq)
	if err != nil {
		return &g.response, err
	}

	// build commands from proto
	g.buildCommands()

	// write root.go
	g.genRootCmdFile()

	// write completion.go
	g.genCompletionCmdFile()

	// generate Service-level subcommands
	g.genServiceCmdFiles()

	// generate Method-level subcommand files
	g.genCommands()

	return &g.response, nil
}

type gcli struct {
	comments    map[proto.Message]string
	commands    []*Command
	descInfo    pbinfo.Info
	imports     map[string]*pbinfo.ImportSpec
	pt          printer.P
	response    plugin.CodeGeneratorResponse
	root        string
	services    []*descriptor.ServiceDescriptorProto
	subcommands map[string][]*Command
}

func (g *gcli) init(req *plugin.CodeGeneratorRequest) error {
	g.comments = make(map[proto.Message]string)
	g.descInfo = pbinfo.Of(req.ProtoFile)
	g.imports = make(map[string]*pbinfo.ImportSpec)
	g.subcommands = make(map[string][]*Command)

	root, gapic, err := parseParameters(req.Parameter)
	if err != nil {
		errStr := fmt.Sprintf("Error in parsing params: %s", err.Error())
		g.response.Error = &errStr
		return err
	}
	g.root = root

	putImport(g.imports, &pbinfo.ImportSpec{
		Name: "gapic",
		Path: gapic,
	})

	// gather services & imports for generation
	for _, f := range req.ProtoFile {
		if strContains(req.FileToGenerate, f.GetName()) {
			pkg := f.GetOptions().GetGoPackage()
			if pkg == "" {
				errStr := fmt.Sprintf("Error missing Go package option for: %s", f.GetName())
				g.response.Error = &errStr
				return err
			}

			spec := pbinfo.ImportSpec{
				Name: pkg[strings.LastIndexByte(pkg, '/')+1:] + "pb",
				Path: pkg,
			}

			if sep := strings.LastIndex(pkg, ";"); sep != -1 {
				spec.Name = pkg[sep+1:] + "pb"
				spec.Path = pkg[:sep]
			}

			putImport(g.imports, &spec)

			// Comments
			for _, loc := range f.GetSourceCodeInfo().GetLocation() {
				if loc.LeadingComments == nil {
					continue
				}

				g.addComments(f, loc)
			}

			g.services = append(g.services, f.Service...)
		}
	}

	return nil
}

func (g *gcli) buildCommands() {
	// TODO(ndietz) weird result for names containing acronyms
	// i.e. SearchByID -> [Search, By, I, D]
	camelCaseRegex := regexp.MustCompile("[A-Z]+[a-z]*")

	// build commands for desird services
	for _, srv := range g.services {
		for _, mthd := range srv.GetMethod() {
			cmd := Command{
				Imports:          make(map[string]*pbinfo.ImportSpec),
				Service:          pbinfo.ReduceServName(srv.GetName(), ""),
				Method:           mthd.GetName(),
				InputMessageType: mthd.GetInputType(),
				ServerStreaming:  mthd.GetServerStreaming(),
				ClientStreaming:  mthd.GetClientStreaming(),
				OneOfSelectors:   make(map[string]*Flag),
				MethodCmd: strings.ToLower(strings.Join(
					camelCaseRegex.FindAllString(mthd.GetName(), -1), "-")),
			}

			// copy top level imports
			copyImports(g.imports, cmd.Imports)

			// add any available comment as usage
			if cmt, ok := g.comments[mthd]; ok {
				cmt = sanitizeComment(cmt)

				cmd.LongDesc = cmt
				cmd.ShortDesc = toShortUsage(cmt)
			}

			// add input message import
			msg := g.descInfo.Type[cmd.InputMessageType].(*descriptor.DescriptorProto)
			pkg, err := g.addImport(&cmd, msg)
			if err != nil {
				continue
			}

			cmd.InputMessage = fmt.Sprintf("%s.%s",
				pkg.Name,
				cmd.InputMessageType[strings.LastIndex(cmd.InputMessageType, ".")+1:])

			// build input fields into flags
			if cmd.InputMessageType != EmptyProtoType && !cmd.ClientStreaming {
				cmd.Flags = append(cmd.Flags, g.buildFieldFlags(&cmd, msg, "", false)...)
			}

			// capture output type for template formatting reasons
			if out := mthd.GetOutputType(); out == LROProtoType {
				cmd.IsLRO = true
				cmd.OutputMessageType = out
			} else if out != EmptyProtoType {
				msg := g.descInfo.Type[out].(*descriptor.DescriptorProto)

				// buildFieldFlags identifies if a Method is paged
				// while iterating over the fields
				if cmd.Paged {
					var f *descriptor.FieldDescriptorProto

					// find repeated field in paged response
					for _, f = range msg.GetField() {

						if f.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
							break
						}
					}

					// primitive repeated type
					if fType := f.GetType(); fType != descriptor.FieldDescriptorProto_TYPE_MESSAGE {
						cmd.OutputMessageType = pbinfo.GoTypeForPrim[fType]
						g.commands = append(g.commands, &cmd)
						g.subcommands[srv.GetName()] = append(g.subcommands[srv.GetName()], &cmd)
						continue
					}

					// set message being evaluated to the repeated type
					out = f.GetTypeName()
					msg = g.descInfo.Type[f.GetTypeName()].(*descriptor.DescriptorProto)
				}

				pkg, err := g.addImport(&cmd, msg)
				if err != nil {
					continue
				}

				cmd.OutputMessageType = pkg.Name + "." + out[strings.LastIndex(out, ".")+1:]
			}

			g.commands = append(g.commands, &cmd)
			g.subcommands[srv.GetName()] = append(g.subcommands[srv.GetName()], &cmd)
		}
	}
}

func (g *gcli) buildOneOfSelectors(cmd *Command, msg *descriptor.DescriptorProto, prefix string) {
	for _, field := range msg.GetOneofDecl() {
		flag := Flag{
			Name:     prefix + field.GetName(),
			Type:     descriptor.FieldDescriptorProto_TYPE_STRING,
			OneOfs:   make(map[string]*Flag),
			Required: true,
		}

		if _, ok := cmd.OneOfSelectors[field.GetName()]; !ok {
			cmd.OneOfSelectors[field.GetName()] = &flag
		}
	}
}

func (g *gcli) buildOneOfFlag(cmd *Command, msg *descriptor.DescriptorProto, field *descriptor.FieldDescriptorProto, prefix string, isNested bool) (flags []*Flag) {
	var output bool

	oneOfField := msg.GetOneofDecl()[field.GetOneofIndex()].GetName()
	oneOfPrefix := prefix + oneOfField + "."

	flag := Flag{
		Name:         oneOfPrefix + field.GetName(),
		Type:         field.GetType(),
		Repeated:     field.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED,
		IsOneOfField: true,
		IsNested:     isNested,
	}

	// evaluate field comments for API behavior
	output, flag.Required, flag.Usage = g.getFieldBehavior(msg, field)
	if output {
		return
	}

	cmd.HasEnums = cmd.HasEnums || flag.IsEnum()

	// construct oneof gRPC struct type info
	m := msg
	if !isNested {
		m = g.descInfo.Type[cmd.InputMessageType].(*descriptor.DescriptorProto)
	}

	parent, err := g.getImport(m)
	if err != nil {
		return
	}
	flag.Message = m.GetName()[strings.LastIndex(m.GetName(), ".")+1:]
	flag.MessageImport = *parent

	// handle oneof message or enum fields
	if flag.IsMessage() || flag.IsEnum() {
		flag.Message = parseMessageName(field, msg)

		nested := g.descInfo.Type[field.GetTypeName()]

		// add nested type import
		pkg, err := g.addImport(cmd, nested)
		if err != nil {
			return
		}
		flag.MessageImport = *pkg

		if flag.IsMessage() {
			cmd.NestedMessages = append(cmd.NestedMessages, &NestedMessage{
				FieldName: flag.GenOneOfVarName("") + "." + flag.OneOfInputFieldName(),
				FieldType: fmt.Sprintf("%s.%s", pkg.Name, flag.Message),
			})

			p := oneOfPrefix + field.GetName() + "."
			m := nested.(*descriptor.DescriptorProto)

			// recursively add singular, nested message fields
			flags = append(flags, g.buildFieldFlags(cmd, m, p, true)...)

			cmd.OneOfSelectors[oneOfField].OneOfs[field.GetName()] = &flag
			return
		}
	}

	cmd.OneOfSelectors[oneOfField].OneOfs[field.GetName()] = &flag
	flags = append(flags, &flag)

	return flags
}

func (g *gcli) buildFieldFlags(cmd *Command, msg *descriptor.DescriptorProto, prefix string, isOneOf bool) []*Flag {
	var flags []*Flag
	var output bool

	for _, field := range msg.GetField() {
		if field.OneofIndex != nil {
			g.buildOneOfSelectors(cmd, msg, prefix)

			in := cmd.InputMessage[strings.LastIndex(cmd.InputMessage, ".")+1:]
			flags = append(flags, g.buildOneOfFlag(cmd, msg, field, prefix, in != msg.GetName())...)

			continue
		}

		flag := Flag{
			Name:         prefix + field.GetName(),
			Type:         field.GetType(),
			Repeated:     field.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED,
			IsOneOfField: isOneOf,
		}

		// evaluate field comments for API behavior
		output, flag.Required, flag.Usage = g.getFieldBehavior(msg, field)
		if output {
			continue
		}

		cmd.HasEnums = cmd.HasEnums || flag.IsEnum()

		if flag.IsMessage() || flag.IsEnum() {
			flag.Message = parseMessageName(field, msg)

			nested := g.descInfo.Type[field.GetTypeName()]

			// add nested message import
			pkg, err := g.addImport(cmd, nested)
			if err != nil {
				continue
			}
			flag.MessageImport = *pkg

			// recursively add singular, nested message fields
			if !flag.Repeated && !flag.IsEnum() {
				n := &NestedMessage{
					FieldName: "." + flag.InputFieldName(),
					FieldType: fmt.Sprintf("%s.%s", pkg.Name, flag.Message),
				}

				if isOneOf {
					n.FieldName = flag.GenOneOfVarName("") + "." + flag.OneOfInputFieldName()
				}

				cmd.NestedMessages = append(cmd.NestedMessages, n)

				p := prefix + field.GetName() + "."
				m := nested.(*descriptor.DescriptorProto)
				flags = append(flags, g.buildFieldFlags(cmd, m, p, isOneOf)...)
				continue
			}
		}

		if name := field.GetName(); name == "page_token" || name == "page_size" {
			cmd.Paged = true
			putImport(cmd.Imports, &pbinfo.ImportSpec{
				Name: "iterator",
				Path: "google.golang.org/api/iterator",
			})
		}

		flags = append(flags, &flag)
	}

	return flags
}

func (g *gcli) getFieldBehavior(msg *descriptor.DescriptorProto, field *descriptor.FieldDescriptorProto) (output bool, required bool, cmt string) {
	if cmt, ok := g.comments[field]; ok {
		cmt = sanitizeComment(cmt)
		output = strings.Contains(cmt, OutputOnlyStr)
		required = strings.Contains(cmt, RequiredStr)
	}

	return
}

func (g *gcli) getImport(t pbinfo.ProtoType) (*pbinfo.ImportSpec, error) {
	_, pkg, err := g.descInfo.NameSpec(t)
	if err != nil {
		errStr := fmt.Sprintf("Error retrieving import for message: %s", err.Error())
		g.response.Error = &errStr
		return nil, err
	}

	return &pkg, nil
}

func (g *gcli) addImport(cmd *Command, t pbinfo.ProtoType) (*pbinfo.ImportSpec, error) {
	pkg, err := g.getImport(t)
	if err != nil {
		return nil, err
	}
	putImport(cmd.Imports, pkg)

	return pkg, nil
}

func (g *gcli) addComments(f *descriptor.FileDescriptorProto, loc *descriptor.SourceCodeInfo_Location) {
	var key proto.Message
	p := loc.Path

	switch {
	// Service comment
	case len(p) == 2 && p[0] == 6:
		key = f.Service[p[1]]

	// Method comment
	case len(p) == 4 && p[0] == 6 && p[2] == 2:
		key = f.Service[p[1]].Method[p[3]]

	// Message comment
	case len(p) == 2 && p[0] == 4:
		key = f.MessageType[p[1]]

	// Field comment
	case len(p) == 4 && p[0] == 4 && p[2] == 2:
		key = f.MessageType[p[1]].Field[p[3]]

	// Extension comment
	case len(p) == 2 && p[0] == 7:
		key = f.Extension[p[1]]

	// Enum comment
	case len(p) == 2 && p[0] == 5:
		key = f.EnumType[p[1]]

	// Enum Value comment
	case len(p) == 4 && p[0] == 5 && p[2] == 2:
		key = f.EnumType[p[1]].Value[p[3]]

	// TODO(ndietz): this is an incomplete mapping of comments
	default:
		return
	}

	g.comments[key] = *loc.LeadingComments
}

func (g *gcli) addGoFile(name string) {
	file := &plugin.CodeGeneratorResponse_File{
		Name: proto.String(name),
	}

	// format and prune unused imports in generatd code
	data, err := goimport.Process(*file.Name, g.pt.Bytes(), nil)
	if err != nil {
		errStr := fmt.Sprintf("Error in formatting output: %s", err.Error())
		g.response.Error = &errStr
		return
	}

	file.Content = proto.String(string(data))

	g.response.File = append(g.response.File, file)
}
