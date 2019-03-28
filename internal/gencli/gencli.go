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
	"go/format"
	"regexp"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"github.com/googleapis/gapic-generator-go/internal/printer"
	"google.golang.org/genproto/googleapis/api/annotations"
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
	ServiceClientType string
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
	g.genCommands()

	// write root.go
	g.genRootCmdFile()

	// write completion.go
	g.genCompletionCmdFile()

	return &g.response, nil
}

type gcli struct {
	comments    map[proto.Message]string
	descInfo    pbinfo.Info
	imports     map[string]*pbinfo.ImportSpec
	pt          printer.P
	response    plugin.CodeGeneratorResponse
	root        string
	services    []*descriptor.ServiceDescriptorProto
	subcommands map[string][]*Command
	format      bool
	gapicName   string
}

func (g *gcli) init(req *plugin.CodeGeneratorRequest) error {
	g.comments = make(map[proto.Message]string)
	g.descInfo = pbinfo.Of(req.ProtoFile)
	g.imports = make(map[string]*pbinfo.ImportSpec)
	g.subcommands = make(map[string][]*Command)

	err := g.parseParameters(req.Parameter)
	if err != nil {
		errStr := fmt.Sprintf("Error in parsing params: %s", err.Error())
		g.response.Error = &errStr
		return err
	}

	// gather services & imports for generation
	for _, f := range req.ProtoFile {
		if strContains(req.FileToGenerate, f.GetName()) {
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

func (g *gcli) genCommands() {
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

			if cmd.InputMessageType != EmptyProtoType {
				// imports handling input
				putImport(cmd.Imports, &pbinfo.ImportSpec{
					Path: "os",
				})

				putImport(cmd.Imports, &pbinfo.ImportSpec{
					Path: "github.com/golang/protobuf/jsonpb",
				})

				if !cmd.ClientStreaming {
					// build input fields into flags
					cmd.Flags = append(cmd.Flags, g.buildFieldFlags(&cmd, msg, "", false)...)

					if cmd.HasEnums {
						putImport(cmd.Imports, &pbinfo.ImportSpec{
							Path: "strings",
						})
					}
				} else {
					putImport(cmd.Imports, &pbinfo.ImportSpec{
						Path: "bufio",
					})
				}
			}

			// capture output type for template formatting reasons
			if out := mthd.GetOutputType(); out == LROProtoType {
				cmd.IsLRO = true
				cmd.OutputMessageType = out

				// add fmt for verbose printing
				putImport(cmd.Imports, &pbinfo.ImportSpec{
					Path: "fmt",
				})
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
						g.subcommands[srv.GetName()] = append(g.subcommands[srv.GetName()], &cmd)
						continue
					}

					// set message being evaluated to the repeated type
					out = f.GetTypeName()
					msg = g.descInfo.Type[f.GetTypeName()].(*descriptor.DescriptorProto)

					putImport(cmd.Imports, &pbinfo.ImportSpec{
						Path: "google.golang.org/api/iterator",
					})
				}

				pkg, err := g.getImport(msg)
				if err != nil {
					continue
				}

				// only need the actual import for server stream unmarshaling
				if cmd.ServerStreaming && !cmd.ClientStreaming {
					putImport(cmd.Imports, pkg)
					putImport(cmd.Imports, &pbinfo.ImportSpec{
						Path: "io",
					})
				}

				cmd.OutputMessageType = pkg.Name + "." + out[strings.LastIndex(out, ".")+1:]

				// add fmt for verbose printing
				putImport(cmd.Imports, &pbinfo.ImportSpec{
					Path: "fmt",
				})
			}

			g.subcommands[srv.GetName()] = append(g.subcommands[srv.GetName()], &cmd)
			g.genCommandFile(&cmd)
		}

		name := pbinfo.ReduceServName(srv.GetName(), "")
		cmd := Command{
			Service:           name,
			ServiceClientType: pbinfo.ReduceServName(srv.GetName(), g.gapicName) + "Client",
			MethodCmd:         strings.ToLower(name),
			ShortDesc:         "Sub-command for Service: " + name,
			Imports: map[string]*pbinfo.ImportSpec{
				"gapic": g.imports["gapic"],
			},
			EnvPrefix:   strings.ToUpper(g.root + "_" + name),
			SubCommands: g.subcommands[srv.GetName()],
		}

		// add any available comment as usage
		if cmt, ok := g.comments[srv]; ok {
			cmt = sanitizeComment(cmt)

			cmd.LongDesc = cmt
			cmd.ShortDesc = toShortUsage(cmt)
		}
		g.genServiceCmdFile(&cmd)
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
		Usage:        sanitizeComment(g.comments[field]),
	}

	// evaluate field behavior
	output, flag.Required = g.getFieldBehavior(field)
	if output {
		return
	}

	if flag.Required {
		flag.Usage = "Required. " + flag.Usage
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
		flag.Message = g.prepareMessageName(field)

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

			// check if we've recursed into a nested message's oneof
			isInNested := g.descInfo.Type[cmd.InputMessageType] != msg

			flags = append(flags, g.buildOneOfFlag(cmd, msg, field, prefix, isInNested)...)

			continue
		}

		flag := Flag{
			Name:         prefix + field.GetName(),
			Type:         field.GetType(),
			Repeated:     field.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED,
			IsOneOfField: isOneOf,
			Usage:        sanitizeComment(g.comments[field]),
		}

		// skip repeated bytes, they end up being [][]byte which isn't a supported pFlag flag
		if flag.IsBytes() && flag.Repeated {
			continue
		}

		// evaluate field behavior
		output, flag.Required = g.getFieldBehavior(field)
		if output {
			continue
		}

		if flag.Required {
			flag.Usage = "Required. " + flag.Usage
		}

		cmd.HasEnums = cmd.HasEnums || flag.IsEnum()

		if flag.IsMessage() || flag.IsEnum() {
			flag.Message = g.prepareMessageName(field)

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
		}

		flags = append(flags, &flag)
	}

	return flags
}

func (g *gcli) getFieldBehavior(field *descriptor.FieldDescriptorProto) (output bool, required bool) {
	eBehav, err := proto.GetExtension(field.GetOptions(), annotations.E_FieldBehavior)
	if err == proto.ErrMissingExtension || field.GetOptions() == nil {
		return
	} else if err != nil {
		errStr := fmt.Sprintf("Error parsing the %s field_behavior: %v", field.GetName(), err)
		g.response.Error = &errStr
		return
	}

	behavior := eBehav.([]annotations.FieldBehavior)

	for _, b := range behavior {
		if b == annotations.FieldBehavior_REQUIRED {
			required = true
		} else if b == annotations.FieldBehavior_OUTPUT_ONLY {
			output = true
		}
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
		Name:    proto.String(name),
		Content: proto.String(g.pt.String()),
	}

	if g.format {
		// format generated code
		data, err := format.Source(g.pt.Bytes())
		if err != nil {
			errStr := fmt.Sprintf("Error in formatting output: %s", err.Error())
			g.response.Error = &errStr
			return
		}

		file.Content = proto.String(string(data))
	}

	g.response.File = append(g.response.File, file)
}

func (g *gcli) prepareMessageName(field *descriptor.FieldDescriptorProto) string {
	f := g.descInfo.Type[field.GetTypeName()]
	name := f.GetName()

	// prepend parent name for nested message types
	for p, ok := g.descInfo.ParentElement[f]; ok; p, ok = g.descInfo.ParentElement[p] {
		name = p.GetName() + "_" + name
	}

	return name
}

func (g *gcli) parseParameters(params *string) (err error) {
	// by default formatting is enabled
	g.format = true

	if params == nil {
		return fmt.Errorf("Missing required parameters. See usage")
	}

	for _, str := range strings.Split(*params, ",") {
		argSep := strings.Index(str, "=")
		if argSep == -1 {
			return fmt.Errorf("Unknown parameter: %s", str)
		}

		switch str[:argSep] {
		case "gapic":
			pkg := str[argSep+1:]
			pkgSep := strings.Index(pkg, ";")
			if pkgSep >= 0 {
				// save the package name for Service name reduction later
				g.gapicName = pkg[pkgSep+1:]
				pkg = pkg[:pkgSep]
			}

			putImport(g.imports, &pbinfo.ImportSpec{
				Name: "gapic",
				Path: pkg,
			})
		case "root":
			g.root = str[argSep+1:]
		case "fmt":
			g.format, err = strconv.ParseBool(str[argSep+1:])
			if err != nil {
				return
			}
		default:
			return fmt.Errorf("Unknown parameter: %s", str)
		}
	}

	if _, ok := g.imports["gapic"]; !ok {
		return fmt.Errorf("Missing option \"gapic=[import path]\". Got %q", *params)
	}

	if g.root == "" {
		return fmt.Errorf("Missing option \"root=[root cmd]\". Got %q", *params)
	}

	return
}
