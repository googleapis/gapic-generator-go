package gencli

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
)

// Flag is used to represent fields as flags
type Flag struct {
	Name     string
	VarName  string
	Type     descriptor.FieldDescriptorProto_Type
	Repeated bool
	Required bool
	Usage    string
}

// ComposeFlagVarName composes the variable name for the generated flag
func (f *Flag) ComposeFlagVarName(method string) {
	f.VarName = strings.Replace(method+strings.Title(f.Name), ".", "", -1)
}

// GenFlagVar generates the Go variable to store the flag value
func (f *Flag) GenFlagVar() string {
	str := "var " + f.VarName
	typeStr := pbinfo.GoTypeForPrim[f.Type]

	if f.Repeated {
		if f.Type == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
			return str + " string // JSON to be parsed"
		} else if strings.Contains(typeStr, "int") {
			typeStr = "int"
		}

		str += " []" + typeStr
	} else {
		str += " " + typeStr
	}

	return str
}

// GenFlag generates the pflag API call for this flag
func (f *Flag) GenFlag() string {
	var str, defaultVal string
	typeStr := pbinfo.GoTypeForPrim[f.Type]
	flagType := strings.Title(typeStr)

	if f.Repeated {
		if strings.Contains(typeStr, "int") {
			flagType = "Int"
			typeStr = "int"
		}

		flagType += "Slice"
		defaultVal = "[]" + typeStr + "{}"
	} else {
		switch typeStr {
		case "bool":
			defaultVal = "false"
		case "string":
			defaultVal = `""`
		case "int32":
			defaultVal = "0"
		case "int":
			defaultVal = "0"
		case "float64":
			defaultVal = "0.0"
		case "[]byte":
			defaultVal = "[]byte{}"
			flagType = "BytesHex"
		default:
			return ""
		}
	}

	str = fmt.Sprintf(`%sVar(&%s, "%s", %s, "%s")`, flagType, f.VarName, f.Name, defaultVal, f.Usage)

	return str
}

// GenRequired generates the code to mark the flag as required
func (f *Flag) GenRequired() string {
	return fmt.Sprintf(`cmd.MarkFlagRequired("%s")`, f.Name)
}
