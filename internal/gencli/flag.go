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
	"unicode"
	"unicode/utf8"

	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/protobuf/types/descriptorpb"
)

// Flag is used to represent fields as flags
type Flag struct {
	Name          string
	Type          descriptorpb.FieldDescriptorProto_Type
	Message       string
	Repeated      bool
	Required      bool
	Usage         string
	MessageImport pbinfo.ImportSpec
	OneOfs        map[string]*Flag
	OneOfSelector string
	OneOfDesc     *desc.OneOfDescriptor
	VarName       string
	FieldName     string
	SliceAccessor string
	IsOneOfField  bool
	IsNested      bool
	IsMap         bool
	Optional      bool

	// Accessor is only set after calling GenFlag
	Accessor string
	MsgDesc  *desc.MessageDescriptor
}

// GenFlag generates the pflag API call for this flag
func (f *Flag) GenFlag() string {
	var str, def string

	tStr := pbinfo.GoTypeForPrim[f.Type]
	if f.IsEnum() {
		tStr = "string"
	}
	fType := toTitle(tStr)

	if f.Repeated {
		if f.Type == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE {
			// repeated Messages are entered as JSON strings and unmarshaled into the Message type later
			return fmt.Sprintf(`StringArrayVar(&%s, "%s", []string{}, "%s")`, f.VarName, f.Name, f.Usage)
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
			// set default page_size to 10 because 0 seems to short circuit the RPC
			if f.FieldName == "PageSize" {
				def = "10"
				f.Usage = fmt.Sprintf("Default is %s. %s", def, f.Usage)
				break
			}

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

	name := f.VarName + "." + f.FieldName
	if len(f.OneOfs) > 0 || f.IsEnum() {
		name = f.VarName
	}

	f.Accessor = name

	if f.Optional && !f.IsEnum() {
		name = f.OptionalVarName()
	}

	str = fmt.Sprintf(`%sVar(&%s, "%s", %s, "%s")`, fType, name, f.Name, def, f.Usage)

	return str
}

// IsMessage is a template helper that reports if the flag is a message type
func (f *Flag) IsMessage() bool {
	return f.Type == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE
}

// IsEnum is a template helper that reports if the flag is of an enum type
func (f *Flag) IsEnum() bool {
	return f.Type == descriptorpb.FieldDescriptorProto_TYPE_ENUM
}

// IsBytes is a helper that reports if the flag is of a type bytes
func (f *Flag) IsBytes() bool {
	return f.Type == descriptorpb.FieldDescriptorProto_TYPE_BYTES
}

// EnumFieldAccess constructs the input message field accessor for an enum
// assignment.
func (f *Flag) EnumFieldAccess(inputVar string) string {
	if f.IsOneOfField {
		seg := strings.LastIndex(f.FieldName, ".")
		inputVar = strings.TrimSuffix(f.VarName, f.FieldName[seg+1:])
	}
	return fmt.Sprintf("%s.%s", inputVar, f.FieldName)
}

// OptionalVarName constructs the place holder variable name for a proto3
// optional field.
func (f *Flag) OptionalVarName() string {
	s := f.VarName + "." + f.FieldName
	s = dotToCamel(s)
	r, w := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[w:]
}

// GoTypeForPrim returns the name of the Go type for a primitive proto type.
func (f *Flag) GoTypeForPrim() string {
	return pbinfo.GoTypeForPrim[f.Type]
}
