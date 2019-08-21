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

package main

import (
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

// validPrims contains functions to check value-validity for primitive types.
// Functions only perform cursory checking; since we use scanner for tokenizing,
// they must already be a valid token of some type.
var validPrims = [...]func(string) bool{
	descriptor.FieldDescriptorProto_TYPE_BOOL: func(s string) bool { return s == "true" || s == "false" || s == "True" || s == "False" },

	descriptor.FieldDescriptorProto_TYPE_BYTES:  validStrings,
	descriptor.FieldDescriptorProto_TYPE_STRING: validStrings,

	descriptor.FieldDescriptorProto_TYPE_DOUBLE: validFloat,
	descriptor.FieldDescriptorProto_TYPE_FLOAT:  validFloat,

	descriptor.FieldDescriptorProto_TYPE_INT64:    validInt,
	descriptor.FieldDescriptorProto_TYPE_UINT64:   validInt,
	descriptor.FieldDescriptorProto_TYPE_INT32:    validInt,
	descriptor.FieldDescriptorProto_TYPE_FIXED64:  validInt,
	descriptor.FieldDescriptorProto_TYPE_FIXED32:  validInt,
	descriptor.FieldDescriptorProto_TYPE_UINT32:   validInt,
	descriptor.FieldDescriptorProto_TYPE_SFIXED32: validInt,
	descriptor.FieldDescriptorProto_TYPE_SFIXED64: validInt,
	descriptor.FieldDescriptorProto_TYPE_SINT32:   validInt,
	descriptor.FieldDescriptorProto_TYPE_SINT64:   validInt,
}

func validStrings(s string) bool {
	return true
}

func validFloat(s string) bool {
	return strings.Trim(s, "+-.0123456789") == ""
}

func validInt(s string) bool {
	return strings.Trim(s, "+-0123456789") == ""
}
