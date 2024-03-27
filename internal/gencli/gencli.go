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

	"github.com/jhump/protoreflect/desc"

	longrunning "cloud.google.com/go/longrunning/autogen/longrunningpb"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"github.com/googleapis/gapic-generator-go/internal/printer"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

const (
	emptyValue = "google.protobuf.Empty"
	// EmptyProtoType is the type name for the Empty message type
	EmptyProtoType = "google.protobuf.Empty"
	// LROProtoType is the type name for the LRO message type
	LROProtoType = "google.longrunning.Operation"
)

// Command intermediate representation of a RPC/Method as a CLI command
type Command struct {
	Service           string
	ServiceClientType string
	Method            string
	MethodCmd         string
	InputMessageType  string
	InputMessage      string
	InputMessageVar   string
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
	HasPageSize       bool
	HasPageToken      bool
	IsLRO             bool
	IsLRORespEmpty    bool
	HasEnums          bool
	HasOptional       bool
	SubCommands       []*Command
}

// NestedMessage represents a nested message that will need to be initialized
// in the generated code
type NestedMessage struct {
	FieldName string
	FieldType string
}

// Gen is the main entry point for code generation of a command line utility
func Gen(genReq *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error) {
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
	protos      map[string]*desc.FileDescriptor
	imports     map[string]*pbinfo.ImportSpec
	pt          printer.P
	response    pluginpb.CodeGeneratorResponse
	root        string
	services    []*desc.ServiceDescriptor
	subcommands map[string][]*Command
	format      bool
	gapicName   string
}

func (g *gcli) init(req *pluginpb.CodeGeneratorRequest) error {
	var err error

	g.comments = make(map[proto.Message]string)
	g.imports = make(map[string]*pbinfo.ImportSpec)
	g.subcommands = make(map[string][]*Command)

	g.protos, err = desc.CreateFileDescriptors(req.GetProtoFile())
	if err != nil {
		return err
	}

	err = g.parseParameters(req.Parameter)
	if err != nil {
		errStr := fmt.Sprintf("error in parsing params: %s", err.Error())
		g.response.Error = &errStr
		return err
	}

	// gather services & imports for generation
	for _, f := range req.FileToGenerate {
		file, ok := g.protos[f]
		if !ok {
			errStr := fmt.Sprintf("Target file %q did not have a parsed descriptor", f)
			g.response.Error = &errStr
			return err
		}

		g.services = append(g.services, file.GetServices()...)
	}

	return nil
}

func (g *gcli) genCommands() {
	// TODO(ndietz) weird result for names containing acronyms
	// i.e. SearchByID -> [Search, By, I, D]
	camelCaseRegex := regexp.MustCompile("[A-Z]+[a-z]*")

	// build commands for desird services
	for _, srv := range g.services {
		for _, mthd := range srv.GetMethods() {
			cmd := Command{
				Imports:          make(map[string]*pbinfo.ImportSpec),
				Service:          pbinfo.ReduceServName(srv.GetName(), ""),
				Method:           mthd.GetName(),
				InputMessageType: mthd.GetInputType().GetFullyQualifiedName(),
				InputMessageVar:  mthd.GetName() + "Input",
				ServerStreaming:  mthd.IsServerStreaming(),
				ClientStreaming:  mthd.IsClientStreaming(),
				OneOfSelectors:   make(map[string]*Flag),
				MethodCmd: strings.ToLower(strings.Join(
					camelCaseRegex.FindAllString(mthd.GetName(), -1), "-")),
			}

			// add any available comment as usage
			if cmt := mthd.GetSourceInfo().GetLeadingComments(); cmt != "" {
				cmt = sanitizeComment(cmt)

				cmd.LongDesc = toLongUsage(cmt)
				cmd.ShortDesc = toShortUsage(cmt)
			}

			// add input message import
			msg := mthd.GetInputType()
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
					cmd.Flags = append(cmd.Flags, g.buildFieldFlags(&cmd, msg, nil, "", false)...)

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
			if out := mthd.GetOutputType().GetFullyQualifiedName(); out == LROProtoType {
				cmd.IsLRO = true
				cmd.OutputMessageType = out

				operationInfo := proto.GetExtension(mthd.GetMethodOptions(), longrunning.E_OperationInfo)
				opInfo := operationInfo.(*longrunning.OperationInfo)
				cmd.IsLRORespEmpty = opInfo == nil || opInfo.GetResponseType() == emptyValue

				// add fmt for verbose printing
				putImport(cmd.Imports, &pbinfo.ImportSpec{
					Path: "fmt",
				})
			} else if out != EmptyProtoType {
				msg := mthd.GetOutputType()

				// buildFieldFlags identifies if a Method is paged
				// while iterating over the fields
				if cmd.Paged {
					var f *desc.FieldDescriptor

					// find repeated field in paged response
					for _, f = range msg.GetFields() {
						if f.IsRepeated() {
							break
						}
					}

					// primitive repeated type
					if fType := f.GetType(); fType != descriptorpb.FieldDescriptorProto_TYPE_MESSAGE {
						cmd.OutputMessageType = pbinfo.GoTypeForPrim[fType]
						g.subcommands[srv.GetName()] = append(g.subcommands[srv.GetName()], &cmd)
						putImport(cmd.Imports, &pbinfo.ImportSpec{
							Path: "google.golang.org/api/iterator",
						})
						// add fmt for verbose printing
						putImport(cmd.Imports, &pbinfo.ImportSpec{
							Path: "fmt",
						})
						g.genCommandFile(&cmd)
						continue
					}

					// set message being evaluated to the repeated type
					msg = f.GetMessageType()
					out = msg.GetName()

					putImport(cmd.Imports, &pbinfo.ImportSpec{
						Path: "google.golang.org/api/iterator",
					})
				}

				pkg, err := getImport(msg)
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
		if cmt := srv.GetSourceInfo().GetLeadingComments(); cmt != "" {
			cmt = sanitizeComment(cmt)

			cmd.LongDesc = toLongUsage(cmt)
			cmd.ShortDesc = toShortUsage(cmt)
		}

		g.genServiceCmdFile(&cmd)
	}
}

func (g *gcli) buildOneOfSelectors(cmd *Command, msg *desc.MessageDescriptor, prefix string) {
	for _, field := range msg.GetOneOfs() {
		// proto3_optional fields are represented as a oneof
		// with the same field name preceded by a "_"
		if strings.HasPrefix(field.GetName(), "_") {
			continue
		}

		flag := Flag{
			Name:     prefix + field.GetName(),
			Type:     descriptorpb.FieldDescriptorProto_TYPE_STRING,
			OneOfs:   make(map[string]*Flag),
			Required: true,
			Usage:    buildOneOfUsage(field),
		}

		flag.FieldName = title(flag.Name)

		n := title(flag.Name)
		n = dotToCamel(n)
		flag.VarName = cmd.InputMessageVar + n

		if _, ok := cmd.OneOfSelectors[field.GetName()]; !ok {
			cmd.OneOfSelectors[field.GetName()] = &flag
		}
	}
}

func (g *gcli) buildOneOfFlag(cmd *Command, msg *desc.MessageDescriptor, field *desc.FieldDescriptor, prefix string, isNested bool) (flags []*Flag) {
	var outputOnly bool

	oneOfField := field.GetOneOf().GetName()
	oneOfPrefix := prefix + oneOfField + "."

	flag := Flag{
		Name:          oneOfPrefix + field.GetName(),
		Type:          field.GetType(),
		Repeated:      field.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED,
		IsOneOfField:  true,
		IsNested:      isNested,
		OneOfSelector: prefix + oneOfField,
		OneOfDesc:     field.GetOneOf(),
		Usage:         toShortUsage(sanitizeComment(field.GetSourceInfo().GetLeadingComments())),
		VarName:       cmd.InputMessageVar + dotToCamel(title(oneOfPrefix+field.GetName())),
	}

	// evaluate field behavior
	outputOnly, flag.Required = g.getFieldBehavior(field)
	if outputOnly {
		return
	}

	cmd.HasEnums = cmd.HasEnums || flag.IsEnum()

	if flag.Required && !strings.HasPrefix(flag.Usage, "Required. ") {
		flag.Usage = "Required. " + flag.Usage
	}

	// build input FieldName
	fieldName := title(strings.TrimPrefix(flag.Name, oneOfPrefix))
	if !isNested {
		ndx := strings.Index(fieldName, ".")
		fieldName = fieldName[ndx+1:]
	}
	flag.FieldName = fieldName

	// construct oneof gRPC struct type info
	parent, err := getImport(msg)
	if err != nil {
		return
	}
	flag.Message = msg.GetName()[strings.LastIndex(msg.GetName(), ".")+1:]
	flag.MessageImport = *parent

	// handle oneof message or enum fields
	if flag.IsMessage() {
		nested := field.GetMessageType()

		// add nested type import
		pkg, err := g.addImport(cmd, nested)
		if err != nil {
			return
		}
		flag.MessageImport = *pkg
		flag.MsgDesc = nested

		cmd.NestedMessages = append(cmd.NestedMessages, &NestedMessage{
			FieldName: flag.VarName + "." + flag.FieldName,
			FieldType: fmt.Sprintf("%s.%s", pkg.Name, g.prepareName(nested)),
		})

		p := oneOfPrefix + field.GetName() + "."

		// recursively add singular, nested message fields
		flags = append(flags, g.buildFieldFlags(cmd, nested, field, p, true)...)

		cmd.OneOfSelectors[oneOfField].OneOfs[field.GetName()] = &flag

		return
	} else if flag.IsEnum() {
		e := field.GetEnumType()
		flag.Message = g.prepareName(e)

		// add enum type import
		pkg, err := g.addImport(cmd, e)
		if err != nil {
			return
		}
		flag.MessageImport = *pkg
	}

	cmd.OneOfSelectors[oneOfField].OneOfs[field.GetName()] = &flag
	flags = append(flags, &flag)

	return flags
}

func (g *gcli) buildFieldFlags(cmd *Command, msg *desc.MessageDescriptor, parent *desc.FieldDescriptor, prefix string, isOneOf bool) []*Flag {
	var flags []*Flag
	var outputOnly bool

	// check if we've recursed into a nested message
	isInNested := msg.GetFullyQualifiedName() != cmd.InputMessageType

	for _, field := range msg.GetFields() {
		proto3optional := field.AsFieldDescriptorProto().GetProto3Optional()

		if oneof := field.GetOneOf(); oneof != nil && !proto3optional {
			// add fmt for oneof choice error formatting
			putImport(cmd.Imports, &pbinfo.ImportSpec{
				Path: "fmt",
			})

			// build oneof option selector flags
			g.buildOneOfSelectors(cmd, msg, prefix)

			// build flags for oneof option fields
			oneofs := g.buildOneOfFlag(cmd, msg, field, prefix, isInNested)

			// post-process oneof field flags in context of parent oneof
			for _, o := range oneofs {
				o.OneOfSelector = prefix + oneof.GetName()
				o.FieldName = title(strings.TrimPrefix(o.Name, o.OneOfSelector+"."))

				// oneof fields containing repeated message fields
				// need to update the corresponding slice accessor
				if o.Repeated && o.IsMessage() {
					n := title(o.Name)
					n = n[:strings.LastIndex(n, ".")]
					n = dotToCamel(n)

					o.SliceAccessor = fmt.Sprintf("%s%s.%s", cmd.InputMessageVar, n, o.FieldName)
				}
			}

			flags = append(flags, oneofs...)

			continue
		}

		flag := Flag{
			Name:         prefix + field.GetName(),
			FieldName:    title(prefix + field.GetName()),
			Type:         field.GetType(),
			IsMap:        field.IsMap(),
			Repeated:     field.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED,
			IsOneOfField: isOneOf,
			IsNested:     isInNested,
			Usage:        toShortUsage(sanitizeComment(field.GetSourceInfo().GetLeadingComments())),
		}

		// Optional messages aren't handled differently because they are already
		// pointers to a struct.
		flag.Optional = proto3optional && !flag.IsMessage()

		// evaluate field behavior
		outputOnly, flag.Required = g.getFieldBehavior(field)
		if flag.Required && !strings.HasPrefix(flag.Usage, "Required. ") {
			flag.Usage = "Required. " + flag.Usage
		}

		// skip repeated bytes, they end up being [][]byte which
		// isn't a supported pFlag flag, and skip output only fields
		if (flag.IsBytes() && flag.Repeated) || outputOnly {
			continue
		}

		cmd.HasOptional = cmd.HasOptional || flag.Optional
		cmd.HasEnums = cmd.HasEnums || flag.IsEnum()

		if flag.IsMap {
			// add "strings" import for key=value string split
			putImport(cmd.Imports, &pbinfo.ImportSpec{
				Path: "strings",
			})
			flag.Usage = "key=value pairs. " + flag.Usage
		}

		// build the variable name this field belongs to
		n := title(flag.Name)

		// oneof option fields exclude the actual field name
		// from the var name
		if flag.IsOneOfField && !flag.IsMessage() && !flag.IsEnum() {
			n = n[:strings.LastIndex(n, ".")]

			// A primitive field of a nested message belonging to a nested
			// message oneof field shouldn't use its parent's name in VarName.
			if parent != nil && parent.GetOneOf() == nil {
				n = strings.TrimSuffix(n, title(parent.GetName()))
			}
		}

		flag.VarName = cmd.InputMessageVar + dotToCamel(n)

		// top-level, primitive type fields reference their parent directly
		if !flag.IsOneOfField && !flag.IsMessage() && !flag.IsEnum() {
			flag.VarName = cmd.InputMessageVar
		}

		// handle a field of another Message type
		if flag.IsMessage() {
			// only actually used when repeated
			flag.SliceAccessor = fmt.Sprintf("%s.%s", cmd.InputMessageVar, flag.FieldName)

			// handle nested message information
			nested := field.GetMessageType()
			flag.Message = g.prepareName(nested)
			pkg, err := g.addImport(cmd, nested)
			if err != nil {
				continue
			}
			flag.MessageImport = *pkg

			// recursively add non-repeated, nested message fields
			if !flag.Repeated {
				flag.VarName = cmd.InputMessageVar

				n := &NestedMessage{
					FieldName: flag.VarName + "." + flag.FieldName,
					FieldType: fmt.Sprintf("%s.%s", pkg.Name, flag.Message),
				}

				// fields belonging to a oneof option need
				// the selector prefix to be trimmed
				if isOneOf {
					fieldName := title(strings.TrimPrefix(flag.Name, flag.OneOfSelector+"."))

					// Nested message fields that belong to a nested message
					// oneof but aren't a oneof themselves don't have a
					// OneOfSelector to TrimPrefix.
					//
					// NOTE(ndietz): this is not a solid fix for further depth.
					// A rewrite might be needed.
					if isInNested && field.GetOneOf() == nil {
						flag.VarName += title(dotToCamel(prefix))
						split := strings.Split(prefix, ".")
						fieldName = title(split[len(split)-2]) + "." + title(field.GetName())
					}

					n.FieldName = flag.VarName + "." + fieldName
				}

				cmd.NestedMessages = append(cmd.NestedMessages, n)

				p := prefix + field.GetName() + "."

				flags = append(flags, g.buildFieldFlags(cmd, nested, field, p, isOneOf)...)
				continue
			}
		} else if flag.IsEnum() {
			e := field.GetEnumType()
			flag.Message = g.prepareName(e)

			// add enum type import
			pkg, err := g.addImport(cmd, e)
			if err != nil {
				continue
			}
			flag.MessageImport = *pkg
		}

		if name := field.GetName(); name == "page_token" {
			cmd.HasPageToken = true
		} else if name == "page_size" {
			cmd.HasPageSize = true
		}

		cmd.Paged = cmd.HasPageSize && cmd.HasPageToken

		flags = append(flags, &flag)
	}

	return flags
}

func (g *gcli) getFieldBehavior(field *desc.FieldDescriptor) (output bool, required bool) {
	if field.GetFieldOptions() == nil {
		return
	}

	eBehav := proto.GetExtension(field.GetFieldOptions(), annotations.E_FieldBehavior)
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

func getImport(m desc.Descriptor) (*pbinfo.ImportSpec, error) {
	pkg := m.GetFile().GetFileOptions().GetGoPackage()
	appendpb := func(n string) string {
		if !strings.HasSuffix(n, "pb") {
			n += "pb"
		}
		return n
	}

	// the below logic is copied from pbinfo.NameSpec()
	if pkg == "" {
		return &pbinfo.ImportSpec{}, fmt.Errorf("can't determine import path for %q, file %q missing `option go_package`", m.GetName(), m.GetFile().GetName())
	}

	if p := strings.IndexByte(pkg, ';'); p >= 0 {
		return &pbinfo.ImportSpec{Path: pkg[:p], Name: appendpb(pkg[p+1:])}, nil
	}

	for {
		p := strings.LastIndexByte(pkg, '/')
		if p < 0 {
			return &pbinfo.ImportSpec{Path: pkg, Name: appendpb(pkg)}, nil
		}
		elem := pkg[p+1:]
		if len(elem) >= 2 && elem[0] == 'v' && elem[1] >= '0' && elem[1] <= '9' {
			// It's a version number; skip so we get a more meaningful name
			pkg = pkg[:p]
			continue
		}
		return &pbinfo.ImportSpec{Path: pkg, Name: appendpb(elem)}, nil
	}
}

func (g *gcli) addImport(cmd *Command, m desc.Descriptor) (*pbinfo.ImportSpec, error) {
	pkg, err := getImport(m)
	if err != nil {
		return nil, err
	}
	putImport(cmd.Imports, pkg)

	return pkg, nil
}

func (g *gcli) addGoFile(name string) {
	file := &pluginpb.CodeGeneratorResponse_File{
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

func (g *gcli) prepareName(m desc.Descriptor) string {
	name := m.GetName()

	// prepend parent name for nested types
	for par := m.GetParent(); par != nil; par = par.GetParent() {
		if _, ok := par.(*desc.MessageDescriptor); !ok {
			break
		}

		name = par.GetName() + "_" + name
	}

	return name
}

func (g *gcli) parseParameters(params *string) (err error) {
	// by default formatting is enabled
	g.format = true

	if params == nil {
		return fmt.Errorf("parameters should not be nil")
	}

	for _, str := range strings.Split(*params, ",") {
		argSep := strings.Index(str, "=")
		if argSep == -1 {
			return fmt.Errorf("unknown parameter: %s", str)
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
			return fmt.Errorf("unknown parameter: %s", str)
		}
	}

	if _, ok := g.imports["gapic"]; !ok {
		return fmt.Errorf("missing option \"gapic=[import path]\". Got %q", *params)
	}

	if g.root == "" {
		return fmt.Errorf("missing option \"root=[root cmd]\". Got %q", *params)
	}

	return
}

func buildOneOfUsage(oneof *desc.OneOfDescriptor) string {
	var usage strings.Builder
	fmt.Fprint(&usage, "Choices:")

	for _, choice := range oneof.GetChoices() {
		fmt.Fprintf(&usage, " %s,", choice.GetName())
	}

	// remove trailing comma
	u := usage.String()
	u = u[:len(u)-1]

	return u
}
