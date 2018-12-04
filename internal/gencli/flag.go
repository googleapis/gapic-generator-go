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
	IsNested      bool
}

// GenOtherVarName generates the Go variable to store repeated Message & Enum string values
func (f *Flag) GenOtherVarName(in string) string {
	return in + strings.Replace(f.InputFieldName(), ".", "", -1)
}

// GenOneOfVarName generates the variable name for a oneof entry type
func (f *Flag) GenOneOfVarName(in string) string {
	name := f.InputFieldName()
	if !f.IsNested && strings.Count(f.Name, ".") > 1 {
		name = f.Name[:strings.LastIndex(f.Name, ".")]
	}

	for _, tkn := range strings.Split(name, ".") {
		in += strings.Title(tkn)
	}

	return in
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
			field := f.GenOtherVarName(in)
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
		name = f.GenOneOfVarName(in)
	} else if f.IsEnum() {
		name = f.GenOtherVarName(in)
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
	ndx := strings.Index(name, ".")

	if f.IsNested {
		ndx = strings.LastIndex(name, ".")
	}

	return name[ndx+1:]
}

// InputFieldName converts the field name into the Go struct property name
func (f *Flag) InputFieldName() string {
	split := strings.Split(f.Name, "_")
	for ndx, tkn := range split {
		split[ndx] = strings.Title(tkn)
	}

	return strings.Join(split, "")
}
