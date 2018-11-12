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
func (f *Flag) GenRepeatedMessageVarName(inputVar string) string {
	return inputVar + strings.Replace(f.InputFieldName(), ".", "", -1)
}

func (f *Flag) GenOneOfVarName(inputVar string) string {
	// EchoInput & response.error.message --> EchoInputResponseError
	name := f.Name
	if strings.Count(f.Name, ".") > 1 {
		name = f.Name[:strings.LastIndex(f.Name, ".")]
	}

	for _, tkn := range strings.Split(name, ".") {
		inputVar += strings.Title(tkn)
	}

	return inputVar
}

// IsMessage is a template helper that reports if the flag is a message type
func (f *Flag) IsMessage() bool {
	return f.Type == descriptor.FieldDescriptorProto_TYPE_MESSAGE
}

// GenFlag generates the pflag API call for this flag
func (f *Flag) GenFlag(inputVar string) string {
	var str, defaultVal string
	typeStr := pbinfo.GoTypeForPrim[f.Type]
	flagType := strings.Title(typeStr)

	if f.Repeated {
		if f.Type == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
			field := f.GenRepeatedMessageVarName(inputVar)
			// repeated Messages are entered as JSON strings and unmarshaled into the Message type later
			return fmt.Sprintf(`StringArrayVar(&%s, "%s", []string{}, "%s")`, field, f.Name, f.Usage)
		}

		flagType += "Slice"
		defaultVal = "[]" + typeStr + "{}"
	} else {
		switch typeStr {
		case "bool":
			defaultVal = "false"
		case "string":
			defaultVal = `""`
		case "int32", "int64", "int", "uint32", "uint64":
			defaultVal = "0"
		case "float32", "float64":
			defaultVal = "0.0"
		case "[]byte":
			defaultVal = "[]byte{}"
			flagType = "BytesHex"
		default:
			return ""
		}
	}

	var name string
	if f.IsOneOfField {
		name = f.GenOneOfVarName(inputVar) + "." + f.OneOfInputFieldName()
	} else if len(f.OneOfs) > 1 {
		name = inputVar + f.InputFieldName()
	} else {
		name = inputVar + "." + f.InputFieldName()
	}

	str = fmt.Sprintf(`%sVar(&%s, "%s", %s, "%s")`, flagType, name, f.Name, defaultVal, f.Usage)

	return str
}

// GenRequired generates the code to mark the flag as required
func (f *Flag) GenRequired() string {
	return fmt.Sprintf(`cmd.MarkFlagRequired("%s")`, f.Name)
}

func (f *Flag) OneOfInputFieldName() string {
	split := strings.Split(f.Name, "_")
	for ndx, tkn := range split {
		split[ndx] = strings.Title(tkn)
	}
	name := strings.Join(split, "")

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
