// Copyright 2023 Google LLC
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

package gengapic

import (
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/typepb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var wellKnownTypeNames = []string{
	".google.protobuf.FieldMask",
	".google.protobuf.Timestamp",
	".google.protobuf.Duration",
	".google.protobuf.DoubleValue",
	".google.protobuf.FloatValue",
	".google.protobuf.Int64Value",
	".google.protobuf.UInt64Value",
	".google.protobuf.Int32Value",
	".google.protobuf.UInt32Value",
	".google.protobuf.BoolValue",
	".google.protobuf.StringValue",
	".google.protobuf.BytesValue",
	".google.protobuf.Value",
	".google.protobuf.ListValue",
}

// TODO: Look into if we need to support ListValue/Value fields with
// StringValue or BytesValue values.
var wellKnownStringTypes = []string{
	".google.protobuf.FieldMask",
	".google.protobuf.Timestamp",
	".google.protobuf.Duration",
	".google.protobuf.StringValue",
	".google.protobuf.BytesValue",
}

var wellKnownTypeFiles = []*descriptorpb.FileDescriptorProto{
	protodesc.ToFileDescriptorProto(emptypb.File_google_protobuf_empty_proto),
	protodesc.ToFileDescriptorProto(anypb.File_google_protobuf_any_proto),
	protodesc.ToFileDescriptorProto(structpb.File_google_protobuf_struct_proto),
	protodesc.ToFileDescriptorProto(durationpb.File_google_protobuf_duration_proto),
	protodesc.ToFileDescriptorProto(timestamppb.File_google_protobuf_timestamp_proto),
	protodesc.ToFileDescriptorProto(fieldmaskpb.File_google_protobuf_field_mask_proto),
	protodesc.ToFileDescriptorProto(wrapperspb.File_google_protobuf_wrappers_proto),
	protodesc.ToFileDescriptorProto(typepb.File_google_protobuf_type_proto),
}
