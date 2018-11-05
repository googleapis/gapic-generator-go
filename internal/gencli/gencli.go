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
	"golang.org/x/tools/imports"
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
		g.parseFieldComments(f)

		if strContains(genReq.FileToGenerate, f.GetName()) {
			g.services = append(g.services, f.Service...)
		}
	}

	// write main.go
	g.genMainFile()

	// write root.go
	g.genRootCmdFile(root)

	// build commands from proto
	g.buildCommands()

	// generate command files
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
	camelCaseRegex := regexp.MustCompile("[A-Z]+[a-z]*")

	for _, srv := range g.services {
		for _, rpc := range srv.GetMethod() {
			cmd := Command{}

			// parse Service name for base subcommand
			cmd.Service = pbinfo.ReduceServName(srv.GetName(), "")

			// build input fields into flags if not Empty
			if rpc.GetInputType() != EmptyProtoType {
				cmd.InputMessage = rpc.GetInputType()
				msg := g.descInfo.Type[cmd.InputMessage].(*descriptor.DescriptorProto)

				cmd.Flags = append(cmd.Flags, g.buildFieldFlags(rpc.GetName(), msg, "")...)
			}

			cmd.Method = rpc.GetName()

			// format Method into command line form
			rpcSplit := camelCaseRegex.FindAllString(rpc.GetName(), -1)
			cmd.MethodCmd = strings.ToLower(strings.Join(rpcSplit, "-"))

			g.commands = append(g.commands, &cmd)
		}
	}
}

func (g *gcli) buildFieldFlags(method string, msg *descriptor.DescriptorProto, prefix string) []*Flag {
	var flags []*Flag

	for _, field := range msg.GetField() {
		fieldType := field.GetType()
		required := false

		if cmt, ok := g.comments[msg.GetName()+"."+field.GetName()]; ok {
			if strings.Contains(cmt, OutputOnlyStr) {
				continue
			}

			if strings.Contains(cmt, RequiredStr) {
				required = true
			}
		}

		repeated := field.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED

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

		flag.ComposeFlagVarName(method)

		flags = append(flags, &flag)
	}

	return flags
}

func (g *gcli) parseFieldComments(f *descriptor.FileDescriptorProto) {
	for _, loc := range f.GetSourceCodeInfo().GetLocation() {
		if loc.LeadingComments == nil {
			continue
		}

		// loc points to a Field's SourceCode
		if len(loc.Path) == SourceSubItemPathLength &&
			loc.Path[SourceItemTypePathIndex] == SourceMessageTypePathValue &&
			loc.Path[SourceSubItemTypePathIndex] == SourceFieldTypePathValue {
			// make map key
			key := f.MessageType[loc.Path[SourceMessagePathIndex]].GetName() +
				"." +
				f.MessageType[loc.Path[SourceMessagePathIndex]].Field[loc.Path[SourceFieldPathIndex]].GetName()

			// add new Field comments
			if _, ok := g.comments[key]; !ok {
				g.comments[key] = *loc.LeadingComments
			}
		}
	}
}

func (g *gcli) addGoFile(name string) {
	file := &plugin.CodeGeneratorResponse_File{
		Name: proto.String(name),
	}

	data, err := imports.Process(*file.Name, g.pt.Bytes(), nil)
	if err != nil {
		errStr := fmt.Sprintf("Error in formatting output: %s", err.Error())
		g.response.Error = &errStr
		return
	}

	file.Content = proto.String(string(data))

	g.response.File = append(g.response.File, file)
}
