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

	// SourceSubItemPathLength the length of a source code path
	// referencing a sub-item of a Message or Service
	SourceSubItemPathLength = 4
	// SourceMessageTypePathValue the value that indicates a
	// source code path belongs to a Message
	SourceMessageTypePathValue = 4
	// SourceMessagePathIndex the index of a source code path
	// referncing a specific Message
	SourceMessagePathIndex = 1
	// SourceFieldPathIndex the index of a source code path
	// referencing a specific Field
	SourceFieldPathIndex = 3
	// SourceFieldTypePathValue represents the value associated
	// with the Field type in the source code path
	SourceFieldTypePathValue = 2
	// SourceItemTypePathIndex the index of the source code
	// item path type
	SourceItemTypePathIndex = 0
	// SourceSubItemTypePathIndex the index of the sub-item type
	// in a source code path
	SourceSubItemTypePathIndex = 2

	// OutputOnlyStr represents the comment-level string that
	// indicates a field is only present as output, never input
	OutputOnlyStr = "Output only"
	// RequiredStr represents the comment-level string that
	// indicates a field is required
	RequiredStr = "Required"
)

// Command intermediate representation of a RPC/Method as a CLI command
type Command struct {
	Service          string
	Method           string
	MethodCmd        string
	InputMessageType string
	InputMessage     string
	Flags            []*Flag
	ShortDesc        string
	LongDesc         string
	Imports          map[string]*pbinfo.ImportSpec
	NestedMessages   []*NestedMessage
	EnvPrefix        string
	OutputType       string
	ServerStreaming  bool
	ClientStreaming  bool
	Paged            bool
	OneOfTypes       map[string]*Flag
}

// NestedMessage represents a nested message that will need to be initialized
// in the generated code
type NestedMessage struct {
	FieldName    string
	FieldType    string
	IsOneOfField bool
}

// Gen is the main entry point for code generation of a command line utility
func Gen(genReq *plugin.CodeGeneratorRequest) (*plugin.CodeGeneratorResponse, error) {
	var g gcli

	root, gapicPkg, err := parseParameters(genReq.Parameter)
	if err != nil {
		errStr := fmt.Sprintf("Error in parsing params: %s", err.Error())
		g.response.Error = &errStr
		return &g.response, err
	}
	g.Root = root

	g.init(genReq.ProtoFile)
	putImport(g.imports,
		&pbinfo.ImportSpec{Name: "gapic", Path: gapicPkg})

	// gather services for generation
	for _, f := range genReq.ProtoFile {
		if strContains(genReq.FileToGenerate, f.GetName()) {
			// gather imports for target proto gRPC libs
			if goPkg := f.GetOptions().GetGoPackage(); goPkg != "" {
				var pkgSpec pbinfo.ImportSpec

				if sep := strings.LastIndex(goPkg, ";"); sep != -1 {
					pkgSpec.Name = goPkg[sep+1:] + "pb"
					pkgSpec.Path = goPkg[:sep]
				} else {
					pkgSpec.Name = goPkg[strings.LastIndexByte(goPkg, '/')+1:] + "pb"
					pkgSpec.Path = goPkg
				}

				putImport(g.imports, &pkgSpec)
			}

			g.services = append(g.services, f.Service...)
		}
	}

	// write root.go
	g.genRootCmdFile()

	// generate Service-level subcommands
	g.genServiceCmdFiles()

	// build commands from proto
	g.buildCommands()

	// generate Method-level subcommand files
	g.genCommands()

	return &g.response, nil
}

type gcli struct {
	comments map[string]string
	services []*descriptor.ServiceDescriptorProto
	commands []*Command
	response plugin.CodeGeneratorResponse
	pt       printer.P
	descInfo pbinfo.Info
	imports  map[string]*pbinfo.ImportSpec
	Root     string
}

func (g *gcli) init(f []*descriptor.FileDescriptorProto) {
	g.comments = make(map[string]string)
	g.descInfo = pbinfo.Of(f)
	g.imports = make(map[string]*pbinfo.ImportSpec)
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
				OneOfTypes:       make(map[string]*Flag),
				MethodCmd: strings.ToLower(strings.Join(
					camelCaseRegex.FindAllString(mthd.GetName(), -1), "-")),
			}

			// TODO(ndietz) take this out and support
			if cmd.ServerStreaming && cmd.ClientStreaming {
				continue
			}

			// copy top level imports
			for _, val := range g.imports {
				putImport(cmd.Imports, val)
			}

			// add any available comment as usage
			key := pbinfo.BuildElementCommentKey(g.descInfo.ParentFile[srv], mthd)
			if cmt, ok := g.descInfo.Comments[key]; ok {
				cmt = strings.TrimSpace(strings.Replace(cmt, "\n", " ", -1))

				cmd.LongDesc = cmt
				cmd.ShortDesc = toShortUsage(cmt)
			}

			// gather necessary imports for input message
			msg := g.descInfo.Type[cmd.InputMessageType].(*descriptor.DescriptorProto)
			_, pkg, err := g.descInfo.NameSpec(msg)
			if err != nil {
				errStr := fmt.Sprintf("Error retrieving import for message: %s", err.Error())
				g.response.Error = &errStr
				continue
			}
			putImport(cmd.Imports, &pkg)

			cmd.InputMessage = fmt.Sprintf("%s.%s",
				pkg.Name,
				cmd.InputMessageType[strings.LastIndex(cmd.InputMessageType, ".")+1:])

			// build input fields into flags if not Empty
			if cmd.InputMessageType != EmptyProtoType && !cmd.ClientStreaming {
				cmd.Flags = append(cmd.Flags, g.buildFieldFlags(&cmd, msg, "", false)...)
				g.buidOneOfTypeFlags(&cmd, msg)
				cmd.Flags = append(cmd.Flags, g.buildOneOfFlags(&cmd, msg, "")...)
			}

			// capture output type for template formatting reasons
			if out := mthd.GetOutputType(); out != EmptyProtoType {
				// gather necessary imports for output message
				msg := g.descInfo.Type[out].(*descriptor.DescriptorProto)

				// buildFieldFlags identifies if a Method is paged
				if cmd.Paged {
					var f *descriptor.FieldDescriptorProto

					// find repeated field in paged response
					for _, f = range msg.GetField() {

						if f.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
							break
						}
					}

					// primitive type
					if fType := f.GetType(); fType != descriptor.FieldDescriptorProto_TYPE_MESSAGE {
						cmd.OutputType = pbinfo.GoTypeForPrim[fType]
						g.commands = append(g.commands, &cmd)
						continue
					}

					// set message being evaluated to the repeated type
					out = f.GetTypeName()
					msg = g.descInfo.Type[f.GetTypeName()].(*descriptor.DescriptorProto)
				}

				_, pkg, err := g.descInfo.NameSpec(msg)
				if err != nil {
					errStr := fmt.Sprintf("Error retrieving import for message: %s", err.Error())
					g.response.Error = &errStr
					continue
				}
				putImport(cmd.Imports, &pkg)

				cmd.OutputType = pkg.Name + "." + out[strings.LastIndex(out, ".")+1:]
			}

			g.commands = append(g.commands, &cmd)
		}
	}
}

func (g *gcli) buidOneOfTypeFlags(cmd *Command, msg *descriptor.DescriptorProto) {
	for _, field := range msg.GetOneofDecl() {
		flag := Flag{
			Name:   field.GetName(),
			Type:   descriptor.FieldDescriptorProto_TYPE_STRING,
			OneOfs: make(map[string]*Flag),
		}

		cmd.OneOfTypes[field.GetName()] = &flag
	}
}

func (g *gcli) buildOneOfFlags(cmd *Command, msg *descriptor.DescriptorProto, prefix string) []*Flag {
	var flags []*Flag

	for _, field := range msg.GetField() {
		if field.OneofIndex == nil {
			continue
		}

		oneOfField := msg.GetOneofDecl()[field.GetOneofIndex()].GetName()
		oneOfPrefix := prefix + oneOfField + "."

		flag := Flag{
			Name:         oneOfPrefix + field.GetName(),
			Type:         field.GetType(),
			Repeated:     field.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED,
			IsOneOfField: true,
		}

		// evaluate field comments for API behavior
		if cmt, ok := g.descInfo.Comments[pbinfo.BuildFieldCommentKey(msg, field)]; ok {
			// output-only fields are not added as input flags
			if strings.Contains(cmt, OutputOnlyStr) {
				continue
			}

			flag.Usage = strings.TrimSpace(strings.Replace(cmt, "\n", " ", -1))
			flag.Required = strings.Contains(flag.Usage, RequiredStr)
		}

		// expand singular nested message fields into dot-notation input flags
		if flag.Type == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
			nested := g.descInfo.Type[field.GetTypeName()].(*descriptor.DescriptorProto)

			// gather necessary imports for nested messages
			_, pkg, err := g.descInfo.NameSpec(nested)
			if err != nil {
				errStr := fmt.Sprintf("Error getting import for message: %s", err.Error())
				g.response.Error = &errStr
				continue
			}
			putImport(cmd.Imports, &pkg)

			flag.MessageImport = pkg
			flag.Message = field.GetTypeName()[strings.LastIndex(field.GetTypeName(), ".")+1:]

			// recursively add singular, nested message fields
			cmd.NestedMessages = append(cmd.NestedMessages, &NestedMessage{
				FieldName:    flag.GenOneOfVarName("") + "." + flag.OneOfInputFieldName(),
				FieldType:    fmt.Sprintf("%s.%s", pkg.Name, flag.Message),
				IsOneOfField: true,
			})

			flags = append(flags, g.buildFieldFlags(cmd, nested, oneOfPrefix+field.GetName()+".", true)...)

			cmd.OneOfTypes[oneOfField].OneOfs[field.GetName()] = &flag
			continue
		}

		flags = append(flags, &flag)
		cmd.OneOfTypes[oneOfField].OneOfs[field.GetName()] = &flag
	}

	return flags
}

func (g *gcli) buildFieldFlags(cmd *Command, msg *descriptor.DescriptorProto, prefix string, isOneOf bool) []*Flag {
	var flags []*Flag

	for _, field := range msg.GetField() {
		// TODO(ndietz) remove and add support for oneof
		if field.OneofIndex != nil {
			continue
		}

		flag := Flag{
			Name:         prefix + field.GetName(),
			Type:         field.GetType(),
			Repeated:     field.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED,
			IsOneOfField: isOneOf,
		}

		// evaluate field comments for API behavior
		if cmt, ok := g.descInfo.Comments[pbinfo.BuildFieldCommentKey(msg, field)]; ok {
			// output-only fields are not added as input flags
			if strings.Contains(cmt, OutputOnlyStr) {
				continue
			}

			flag.Usage = strings.TrimSpace(strings.Replace(cmt, "\n", " ", -1))
			flag.Required = strings.Contains(flag.Usage, RequiredStr)
		}

		// expand singular nested message fields into dot-notation input flags
		if flag.Type == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
			nested := g.descInfo.Type[field.GetTypeName()].(*descriptor.DescriptorProto)

			// gather necessary imports for nested messages
			_, pkg, err := g.descInfo.NameSpec(nested)
			if err != nil {
				errStr := fmt.Sprintf("Error getting import for message: %s", err.Error())
				g.response.Error = &errStr
				continue
			}
			putImport(cmd.Imports, &pkg)

			flag.MessageImport = pkg
			flag.Message = field.GetTypeName()[strings.LastIndex(field.GetTypeName(), ".")+1:]

			// recursively add singular, nested message fields
			if !flag.Repeated {
				n := &NestedMessage{
					FieldName:    "." + flag.InputFieldName(),
					FieldType:    fmt.Sprintf("%s.%s", pkg.Name, flag.Message),
					IsOneOfField: isOneOf,
				}

				if isOneOf {
					n.FieldName = flag.GenOneOfVarName("") + "." + flag.OneOfInputFieldName()
				}

				cmd.NestedMessages = append(cmd.NestedMessages, n)

				flags = append(flags, g.buildFieldFlags(cmd, nested, prefix+field.GetName()+".", isOneOf)...)
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
