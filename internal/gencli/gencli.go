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

	root, gapicPkg, err := parseParameters(genReq.Parameter)
	if err != nil {
		errStr := fmt.Sprintf("Error in parsing params: %s", err.Error())
		g.response.Error = &errStr
		return &g.response, err
	}

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
	g.genRootCmdFile(root)

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
			cmd := Command{}
			cmd.Imports = make(map[string]*pbinfo.ImportSpec)

			// copy top level imports
			for _, val := range g.imports {
				putImport(cmd.Imports, val)
			}

			// parse Service name for base subcommand
			cmd.Service = pbinfo.ReduceServName(srv.GetName(), "")

			cmd.Method = mthd.GetName()

			// format Method into command line form
			methodSplit := camelCaseRegex.FindAllString(mthd.GetName(), -1)
			cmd.MethodCmd = strings.ToLower(strings.Join(methodSplit, "-"))

			// build input fields into flags if not Empty
			if mthd.GetInputType() != EmptyProtoType {
				cmd.InputMessageType = mthd.GetInputType()

				msg := g.descInfo.Type[cmd.InputMessageType].(*descriptor.DescriptorProto)

				// gather necessary imports for input message
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

				cmd.Flags = append(cmd.Flags, g.buildFieldFlags(&cmd, msg, "")...)
			}

			// add any available comment as usage
			key := pbinfo.BuildElementCommentKey(g.descInfo.ParentFile[srv], mthd)
			if cmt, ok := g.descInfo.Comments[key]; ok {
				cmt = strings.TrimSpace(strings.Replace(cmt, "\n", " ", -1))

				cmd.LongDesc = cmt
				cmd.ShortDesc = toShortUsage(cmt)
			}

			g.commands = append(g.commands, &cmd)
		}
	}
}

func (g *gcli) buildFieldFlags(cmd *Command, msg *descriptor.DescriptorProto, prefix string) []*Flag {
	var flags []*Flag

	for _, field := range msg.GetField() {
		flag := Flag{
			Name:     prefix + field.GetName(),
			Type:     field.GetType(),
			Repeated: field.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED,
		}

		// evaluate field comments for API behavior
		if cmt, ok := g.descInfo.Comments[pbinfo.BuildFieldCommentKey(msg, field)]; ok {
			// output-only fields are not added as input flags
			if strings.Contains(cmt, OutputOnlyStr) {
				continue
			}

			cmt = strings.TrimSpace(strings.Replace(cmt, "\n", " ", -1))
			flag.Required = strings.Contains(cmt, RequiredStr)
			if flag.Required {
				cmt = cmt[strings.Index(cmt, RequiredStr)+1:]
			}

			flag.Usage = cmt
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
				cmd.NestedMessages = append(cmd.NestedMessages, &NestedMessage{
					FieldName: flag.InputFieldName(),
					FieldType: fmt.Sprintf("%s.%s", pkg.Name, flag.Message),
				})

				flags = append(flags, g.buildFieldFlags(cmd, nested, prefix+field.GetName()+".")...)
				continue
			}
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
