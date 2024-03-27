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

import "google.golang.org/protobuf/types/descriptorpb"

// GoTypeForPrim maps protobuf primitive types to Go primitive types.
var GoTypeForPrim = map[descriptorpb.FieldDescriptorProto_Type]string{
	descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:   "float64",
	descriptorpb.FieldDescriptorProto_TYPE_FLOAT:    "float32",
	descriptorpb.FieldDescriptorProto_TYPE_INT64:    "int64",
	descriptorpb.FieldDescriptorProto_TYPE_UINT64:   "uint64",
	descriptorpb.FieldDescriptorProto_TYPE_INT32:    "int32",
	descriptorpb.FieldDescriptorProto_TYPE_FIXED64:  "uint64",
	descriptorpb.FieldDescriptorProto_TYPE_FIXED32:  "uint32",
	descriptorpb.FieldDescriptorProto_TYPE_BOOL:     "bool",
	descriptorpb.FieldDescriptorProto_TYPE_STRING:   "string",
	descriptorpb.FieldDescriptorProto_TYPE_BYTES:    "[]byte",
	descriptorpb.FieldDescriptorProto_TYPE_UINT32:   "uint32",
	descriptorpb.FieldDescriptorProto_TYPE_SFIXED32: "int32",
	descriptorpb.FieldDescriptorProto_TYPE_SFIXED64: "int64",
	descriptorpb.FieldDescriptorProto_TYPE_SINT32:   "int32",
	descriptorpb.FieldDescriptorProto_TYPE_SINT64:   "int64",
}
