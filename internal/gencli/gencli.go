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
	Flags            []*Flag
	ShortDesc        string
	LongDesc         string
}

// Gen is the main entry point for code generation of a command line utility
func Gen(genReq *plugin.CodeGeneratorRequest) (*plugin.CodeGeneratorResponse, error) {
	var g gcli

	root, _, _, err := parseParameters(genReq.Parameter)
	if err != nil {
		errStr := fmt.Sprintf("Error in parsing params: %s", err.Error())
		g.response.Error = &errStr
		return &g.response, err
	}

	g.init(genReq.ProtoFile)

	// gather services, message types and field comments for generation
	for _, f := range genReq.ProtoFile {
		if strContains(genReq.FileToGenerate, f.GetName()) {
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
}

func (g *gcli) init(f []*descriptor.FileDescriptorProto) {
	g.comments = make(map[string]string)
	g.descInfo = pbinfo.Of(f)
}

func (g *gcli) buildCommands() {
	// TODO(ndietz) weird result for names containing acronyms
	// i.e. SearchByID -> [Search, By, I, D]
	camelCaseRegex := regexp.MustCompile("[A-Z]+[a-z]*")

	// build commands for desird services
	for _, srv := range g.services {
		for _, mthd := range srv.GetMethod() {
			cmd := Command{}

			// parse Service name for base subcommand
			cmd.Service = pbinfo.ReduceServName(srv.GetName(), "")

			// build input fields into flags if not Empty
			if mthd.GetInputType() != EmptyProtoType {
				cmd.InputMessageType = mthd.GetInputType()

				msg := g.descInfo.Type[cmd.InputMessageType].(*descriptor.DescriptorProto)
				cmd.Flags = append(cmd.Flags, g.buildFieldFlags(mthd.GetName(), msg, "")...)
			}

			cmd.Method = mthd.GetName()

			// format Method into command line form
			methodSplit := camelCaseRegex.FindAllString(mthd.GetName(), -1)
			cmd.MethodCmd = strings.ToLower(strings.Join(methodSplit, "-"))

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

func (g *gcli) buildFieldFlags(method string, msg *descriptor.DescriptorProto, prefix string) []*Flag {
	var flags []*Flag

	for _, field := range msg.GetField() {
		fieldType := field.GetType()
		required := false

		// evaluate field comments for API behavior
		if cmt, ok := g.descInfo.Comments[pbinfo.BuildFieldCommentKey(msg, field)]; ok {
			// output-only fields are not added as input flags
			if strings.Contains(cmt, OutputOnlyStr) {
				continue
			}

			if strings.Contains(cmt, RequiredStr) {
				required = true
			}
		}

		repeated := field.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED

		// expand singular nested message fields into dot-notation input flags
		if fieldType == descriptor.FieldDescriptorProto_TYPE_MESSAGE && !repeated {
			// recursively add nested message fields
			nested := g.descInfo.Type[field.GetTypeName()].(*descriptor.DescriptorProto)
			flags = append(flags, g.buildFieldFlags(method, nested, prefix+field.GetName()+".")...)
			continue
		}

		name := field.GetName()
		if prefix != "" {
			name = prefix + name
		}

		flag := Flag{
			Name:     name,
			Type:     fieldType,
			Repeated: repeated,
			Required: required,
		}

		// compose Method-scoped variable name for template use
		flag.ComposeFlagVarName(method)

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
