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

package pbinfo

import (
	"fmt"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

const (
	// SourcePathTopItemLength the length of a top-level type path
	SourcePathTopItemLength = 2
	// SourcePathTopItemTypeIndex the index of the source code
	// item path type
	SourcePathTopItemTypeIndex = 0
	// SourcePathTopItemIndex the index of a source code path
	// referncing a specific Message
	SourcePathTopItemIndex = 1
	// SourcePathSubItemLength the length of a source code path
	// referencing a sub-item of a Message or Service
	SourcePathSubItemLength = 4
	// SourcePathSubItemTypeIndex the index of the sub-item type
	// in a source code path
	SourcePathSubItemTypeIndex = 2
	// SourcePathSubItemIndex the index of a source code path
	// referencing a specific Field
	SourcePathSubItemIndex = 3
	// SourcePathServiceTypeValue the value that indicates a
	// a source code path points to a Service
	SourcePathServiceTypeValue = 6
	// SourcePathMessageTypeValue the value that indicates a
	// source code path points to a Message
	SourcePathMessageTypeValue = 4
	// SourcePathMethodTypeValue the value the indicates a
	// source code path points to a Method
	SourcePathMethodTypeValue = 2
	// SourcePathFieldTypeValue represents the value associated
	// with the Field type in a source code path
	SourcePathFieldTypeValue = 2
	// SourcePathExtensionTypeValue represents the value associated
	// with the Extension type in a source code path
	SourcePathExtensionTypeValue = 7
	// SourcePathEnumTypeValue represents the value associated
	// with the Enum type in a source code path
	SourcePathEnumTypeValue = 5
	// SourcePathEnumValueTypeValue represents the value associated
	// with the Enum Value type in a source code path
	SourcePathEnumValueTypeValue = 2
)

func addComments(comments map[string]string, f *descriptor.FileDescriptorProto, loc *descriptor.SourceCodeInfo_Location) {
	var key string
	p := loc.Path

	switch {
	// Service comment
	case len(p) == SourcePathTopItemLength &&
		p[SourcePathTopItemTypeIndex] == SourcePathServiceTypeValue:

		key = BuildElementCommentKey(f, f.Service[p[SourcePathTopItemIndex]])
	// Method comment
	case len(p) == SourcePathSubItemLength &&
		p[SourcePathTopItemTypeIndex] == SourcePathServiceTypeValue &&
		p[SourcePathSubItemTypeIndex] == SourcePathMethodTypeValue:

		key = BuildElementCommentKey(f, f.Service[p[SourcePathTopItemIndex]].Method[p[SourcePathSubItemIndex]])
	// Message comment
	case len(p) == SourcePathTopItemLength &&
		p[SourcePathTopItemTypeIndex] == SourcePathMessageTypeValue:

		key = BuildElementCommentKey(f, f.MessageType[p[SourcePathTopItemIndex]])
	// Field comment
	case len(p) == SourcePathSubItemLength &&
		p[SourcePathTopItemTypeIndex] == SourcePathMessageTypeValue &&
		p[SourcePathSubItemTypeIndex] == SourcePathFieldTypeValue:

		key = BuildFieldCommentKey(f.MessageType[p[SourcePathTopItemIndex]], f.MessageType[p[SourcePathTopItemIndex]].Field[p[SourcePathSubItemIndex]])
	// Extension comment
	case len(p) == SourcePathTopItemLength &&
		p[SourcePathTopItemTypeIndex] == SourcePathExtensionTypeValue:

		key = BuildElementCommentKey(f, f.Extension[p[SourcePathTopItemIndex]])
	// Enum comment
	case len(p) == SourcePathTopItemLength &&
		p[SourcePathTopItemTypeIndex] == SourcePathEnumTypeValue:

		key = BuildElementCommentKey(f, f.EnumType[p[SourcePathTopItemIndex]])
	// Enum Value comment
	case len(p) == SourcePathSubItemLength &&
		p[SourcePathTopItemTypeIndex] == SourcePathEnumTypeValue &&
		p[SourcePathSubItemTypeIndex] == SourcePathFieldTypeValue:

		key = BuildEnumValueCommentKey(f.EnumType[p[SourcePathTopItemIndex]], f.EnumType[p[SourcePathTopItemIndex]].Value[p[SourcePathSubItemIndex]])
	// TODO(ndietz): this is an incomplete mapping of comments
	default:
		return
	}

	if _, ok := comments[key]; !ok {
		comments[key] = *loc.LeadingComments
	}
}

// BuildFieldCommentKey builds the map key for a given the parent message and child field
func BuildFieldCommentKey(msg *descriptor.DescriptorProto, field *descriptor.FieldDescriptorProto) string {
	return fmt.Sprintf("%s.%s", msg.GetName(), field.GetName())
}

// BuildEnumValueCommentKey builds the map key for a given the parent enum and child enum value
func BuildEnumValueCommentKey(enum *descriptor.EnumDescriptorProto, value *descriptor.EnumValueDescriptorProto) string {
	return fmt.Sprintf("%s.%s", enum.GetName(), value.GetName())
}

// BuildElementCommentKey builds the map key for an arbitrary, top-level element within a FileDescriptor
func BuildElementCommentKey(file *descriptor.FileDescriptorProto, item ProtoType) string {
	return fmt.Sprintf("%s.%s", file.GetPackage(), item.GetName())
}
