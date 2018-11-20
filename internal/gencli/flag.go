package gencli

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
)

// Flag is used to represent fields as flags
type Flag struct {
	Name          string
	Type          descriptor.FieldDescriptorProto_Type
	Message       string
	Repeated      bool
	Required      bool
	Usage         string
	MessageImport pbinfo.ImportSpec
	OneOfs        map[string]*Flag
	IsOneOfField  bool
}

// GenRepeatedMessageVarName generates the Go variable to store repeated Message string values
func (f *Flag) GenRepeatedMessageVarName(in string) string {
	return in + strings.Replace(f.InputFieldName(), ".", "", -1)
}

// GenOneOfVarName generates the variable name for a oneof entry type
func (f *Flag) GenOneOfVarName(in string) string {
	name := f.Name
	if strings.Count(f.Name, ".") > 1 {
		name = f.Name[:strings.LastIndex(f.Name, ".")]
	}

	for _, tkn := range strings.Split(name, ".") {
		in += strings.Title(tkn)
	}

	return in
}

// GenEnumVarName generates the variable name for a enum property
func (f *Flag) GenEnumVarName(in string) string {
	return in + strings.Replace(f.InputFieldName(), ".", "", -1)
}

// GenFlag generates the pflag API call for this flag
func (f *Flag) GenFlag(in string) string {
	var str, def string

	tStr := pbinfo.GoTypeForPrim[f.Type]
	if f.IsEnum() {
		tStr = "string"
	}
	fType := strings.Title(tStr)

	if f.Repeated {
		if f.Type == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
			field := f.GenRepeatedMessageVarName(in)
			// repeated Messages are entered as JSON strings and unmarshaled into the Message type later
			return fmt.Sprintf(`StringArrayVar(&%s, "%s", []string{}, "%s")`, field, f.Name, f.Usage)
		}

		fType += "Slice"
		def = "[]" + tStr + "{}"
	} else {
		switch tStr {
		case "bool":
			def = "false"
		case "string":
			def = `""`
		case "int32", "int64", "int", "uint32", "uint64":
			def = "0"
		case "float32", "float64":
			def = "0.0"
		case "[]byte":
			def = "[]byte{}"
			fType = "BytesHex"
		default:
			return ""
		}
	}

	name := in + "." + f.InputFieldName()
	if f.IsOneOfField {
		name = f.GenOneOfVarName(in) + "." + f.OneOfInputFieldName()
	} else if len(f.OneOfs) > 0 {
		name = in + f.InputFieldName()
	} else if f.IsEnum() {
		name = f.GenEnumVarName(in)
	}

	str = fmt.Sprintf(`%sVar(&%s, "%s", %s, "%s")`, fType, name, f.Name, def, f.Usage)

	return str
}

// GenRequired generates the code to mark the flag as required
func (f *Flag) GenRequired() string {
	return fmt.Sprintf(`cmd.MarkFlagRequired("%s")`, f.Name)
}

// IsMessage is a template helper that reports if the flag is a message type
func (f *Flag) IsMessage() bool {
	return f.Type == descriptor.FieldDescriptorProto_TYPE_MESSAGE
}

// IsEnum is a template helper that reports if the flag is of an enum type
func (f *Flag) IsEnum() bool {
	return f.Type == descriptor.FieldDescriptorProto_TYPE_ENUM
}

// OneOfInputFieldName converts the field name into the Go struct property
// name for a oneof field, which excludes the oneof definition name
func (f *Flag) OneOfInputFieldName() string {
	name := f.InputFieldName()

	return name[strings.Index(name, ".")+1:]
}

// InputFieldName converts the field name into the Go struct property name
func (f *Flag) InputFieldName() string {
	split := strings.Split(f.Name, "_")
	for ndx, tkn := range split {
		split[ndx] = strings.Title(tkn)
	}

	return strings.Join(split, "")
}
