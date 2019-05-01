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
	OneOfSelector string
	VarName       string
	FieldName     string
	SliceAccessor string
	IsOneOfField  bool
	IsNested      bool
}

// GenFlag generates the pflag API call for this flag
func (f *Flag) GenFlag() string {
	var str, def string

	tStr := pbinfo.GoTypeForPrim[f.Type]
	if f.IsEnum() {
		tStr = "string"
	}
	fType := strings.Title(tStr)

	if f.Repeated {
		if f.Type == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
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

	str = fmt.Sprintf(`%sVar(&%s, "%s", %s, "%s")`, fType, name, f.Name, def, f.Usage)

	return str
}

// IsMessage is a template helper that reports if the flag is a message type
func (f *Flag) IsMessage() bool {
	return f.Type == descriptor.FieldDescriptorProto_TYPE_MESSAGE
}

// IsEnum is a template helper that reports if the flag is of an enum type
func (f *Flag) IsEnum() bool {
	return f.Type == descriptor.FieldDescriptorProto_TYPE_ENUM
}

// IsBytes is a helper that reports if the flag is of a type bytes
func (f *Flag) IsBytes() bool {
	return f.Type == descriptor.FieldDescriptorProto_TYPE_BYTES
}
