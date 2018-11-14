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
	comments    map[string]string
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
	g.comments = make(map[string]string)
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

			// TODO(ndietz) take this out and support
			if cmd.ServerStreaming && cmd.ClientStreaming {
				continue
			}

			// copy top level imports
			copyImports(g.imports, cmd.Imports)

			// add any available comment as usage
			key := pbinfo.BuildElementCommentKey(g.descInfo.ParentFile[srv], mthd)
			if cmt, ok := g.descInfo.Comments[key]; ok {
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
			g.buildFlags(&cmd, msg)

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

func (g *gcli) buildFlags(cmd *Command, msg *descriptor.DescriptorProto) {
	if cmd.InputMessageType == EmptyProtoType || cmd.ClientStreaming {
		return
	}

	// build oneof type selector flags, stored separately from field flags
	g.buildOneOfSelectors(cmd, msg)

	// build standard field flags
	cmd.Flags = append(cmd.Flags, g.buildFieldFlags(cmd, msg, "", false)...)

	// build oneof field flags
	cmd.Flags = append(cmd.Flags, g.buildOneOfFlags(cmd, msg, "")...)
}

func (g *gcli) buildOneOfSelectors(cmd *Command, msg *descriptor.DescriptorProto) {
	for _, field := range msg.GetOneofDecl() {
		flag := Flag{
			Name:     field.GetName(),
			Type:     descriptor.FieldDescriptorProto_TYPE_STRING,
			OneOfs:   make(map[string]*Flag),
			Required: true,
		}

		cmd.OneOfSelectors[field.GetName()] = &flag
	}
}

func (g *gcli) buildOneOfFlags(cmd *Command, msg *descriptor.DescriptorProto, prefix string) []*Flag {
	var flags []*Flag

	for _, field := range msg.GetField() {
		// standard fields handled by buildFieldFlags
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

			flag.Usage = sanitizeComment(cmt)
			flag.Required = strings.Contains(flag.Usage, RequiredStr)
		}

		// expand singular nested message fields into dot-notation input flags
		if flag.Type == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
			flag.Message = field.GetTypeName()[strings.LastIndex(field.GetTypeName(), ".")+1:]
			nested := g.descInfo.Type[field.GetTypeName()].(*descriptor.DescriptorProto)

			// add nested message import
			pkg, err := g.addImport(cmd, nested)
			if err != nil {
				continue
			}
			flag.MessageImport = *pkg

			cmd.NestedMessages = append(cmd.NestedMessages, &NestedMessage{
				FieldName: flag.GenOneOfVarName("") + "." + flag.OneOfInputFieldName(),
				FieldType: fmt.Sprintf("%s.%s", pkg.Name, flag.Message),
			})

			// recursively add singular, nested message fields
			flags = append(flags, g.buildFieldFlags(cmd, nested, oneOfPrefix+field.GetName()+".", true)...)

			cmd.OneOfSelectors[oneOfField].OneOfs[field.GetName()] = &flag
			continue
		}

		flags = append(flags, &flag)
		cmd.OneOfSelectors[oneOfField].OneOfs[field.GetName()] = &flag
	}

	return flags
}

func (g *gcli) buildFieldFlags(cmd *Command, msg *descriptor.DescriptorProto, prefix string, isOneOf bool) []*Flag {
	var flags []*Flag

	for _, field := range msg.GetField() {
		// oneof fields handled by buildOneOfFlags
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

			flag.Usage = sanitizeComment(cmt)
			flag.Required = strings.Contains(flag.Usage, RequiredStr)
		}

		// expand singular nested message fields into dot-notation input flags
		if flag.Type == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
			flag.Message = field.GetTypeName()[strings.LastIndex(field.GetTypeName(), ".")+1:]
			nested := g.descInfo.Type[field.GetTypeName()].(*descriptor.DescriptorProto)

			// add nested message import
			pkg, err := g.addImport(cmd, nested)
			if err != nil {
				continue
			}
			flag.MessageImport = *pkg

			// recursively add singular, nested message fields
			if !flag.Repeated {
				n := &NestedMessage{
					FieldName: "." + flag.InputFieldName(),
					FieldType: fmt.Sprintf("%s.%s", pkg.Name, flag.Message),
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

func (g *gcli) addImport(cmd *Command, msg *descriptor.DescriptorProto) (*pbinfo.ImportSpec, error) {
	_, pkg, err := g.descInfo.NameSpec(msg)
	if err != nil {
		errStr := fmt.Sprintf("Error retrieving import for message: %s", err.Error())
		g.response.Error = &errStr
		return nil, err
	}
	putImport(cmd.Imports, &pkg)

	return &pkg, nil
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
