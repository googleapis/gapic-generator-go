// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/api/field_behavior.proto

package annotations

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// An indicator of the behavior of a given field (for example, that a field
// is required in requests, or given as output but ignored as input).
// This DOES NOT change the behavior in protocol buffers itself; it only
// denotes the behavior and may affect how API tooling handles the field.
type FieldBehavior int32

const (
	FieldBehavior_FIELD_BEHAVIOR_UNSPECIFIED FieldBehavior = 0
	// Specifically denotes a field as optional.
	// Because fields are optional by default, this annotation is unnecessary;
	// however, it may be provided if doing so is useful for humans reading
	// the file.
	FieldBehavior_OPTIONAL FieldBehavior = 1
	// Denotes a field as required.
	// This indicates that the field *must* be provided as part of the request,
	// and failure to do so will cause an error (usually `INVALID_ARGUMENT`).
	FieldBehavior_REQUIRED FieldBehavior = 2
	// Denotes a field as output only.
	// This indicates that the field is provided in responses, but including the
	// field in a request does nothing (the server *must* ignore it and
	// *must not* throw an error as a result of the field's presence).
	FieldBehavior_OUTPUT_ONLY FieldBehavior = 3
	// Denotes a field as input only.
	// This indicates that the field is provided in requests, and the
	// corresponding field is not included in output.
	FieldBehavior_INPUT_ONLY FieldBehavior = 4
	// Denotes a field as immutable.
	// This indicates that the field may be set once in a request to create a
	// resource, but may not be changed thereafter.
	FieldBehavior_IMMUTABLE FieldBehavior = 5
)

var FieldBehavior_name = map[int32]string{
	0: "FIELD_BEHAVIOR_UNSPECIFIED",
	1: "OPTIONAL",
	2: "REQUIRED",
	3: "OUTPUT_ONLY",
	4: "INPUT_ONLY",
	5: "IMMUTABLE",
}

var FieldBehavior_value = map[string]int32{
	"FIELD_BEHAVIOR_UNSPECIFIED": 0,
	"OPTIONAL":                   1,
	"REQUIRED":                   2,
	"OUTPUT_ONLY":                3,
	"INPUT_ONLY":                 4,
	"IMMUTABLE":                  5,
}

func (x FieldBehavior) String() string {
	return proto.EnumName(FieldBehavior_name, int32(x))
}

func (FieldBehavior) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_4648f18fd5079967, []int{0}
}

func init() {
	proto.RegisterEnum("google.api.FieldBehavior", FieldBehavior_name, FieldBehavior_value)
}

func init() { proto.RegisterFile("google/api/field_behavior.proto", fileDescriptor_4648f18fd5079967) }

var fileDescriptor_4648f18fd5079967 = []byte{
	// 244 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x8f, 0xbd, 0x4e, 0xc3, 0x30,
	0x14, 0x85, 0x69, 0x29, 0x08, 0x2e, 0xb4, 0x44, 0x1e, 0x19, 0x60, 0x67, 0x48, 0x06, 0x46, 0x26,
	0x9b, 0x38, 0x60, 0x29, 0x8d, 0x4d, 0x88, 0x91, 0x60, 0x89, 0x5c, 0x08, 0xc6, 0x52, 0xf0, 0xb5,
	0xd2, 0x88, 0x85, 0xb7, 0xe1, 0x49, 0x51, 0xd2, 0x8a, 0x9f, 0xed, 0x1e, 0xdd, 0xef, 0x48, 0xdf,
	0x81, 0x73, 0x8b, 0x68, 0xdb, 0x26, 0x31, 0xc1, 0x25, 0xaf, 0xae, 0x69, 0x5f, 0xea, 0x55, 0xf3,
	0x66, 0x3e, 0x1c, 0x76, 0x71, 0xe8, 0xb0, 0x47, 0x02, 0x1b, 0x20, 0x36, 0xc1, 0x5d, 0x7c, 0xc2,
	0x3c, 0x1b, 0x18, 0xb6, 0x45, 0xc8, 0x19, 0x9c, 0x66, 0x82, 0xe7, 0x69, 0xcd, 0xf8, 0x2d, 0x7d,
	0x10, 0xb2, 0xac, 0x75, 0x71, 0xaf, 0xf8, 0xb5, 0xc8, 0x04, 0x4f, 0xa3, 0x1d, 0x72, 0x0c, 0x07,
	0x52, 0x55, 0x42, 0x16, 0x34, 0x8f, 0x26, 0x43, 0x2a, 0xf9, 0x9d, 0x16, 0x25, 0x4f, 0xa3, 0x29,
	0x39, 0x81, 0x23, 0xa9, 0x2b, 0xa5, 0xab, 0x5a, 0x16, 0xf9, 0x63, 0xb4, 0x4b, 0x16, 0x00, 0xa2,
	0xf8, 0xc9, 0x33, 0x32, 0x87, 0x43, 0xb1, 0x5c, 0xea, 0x8a, 0xb2, 0x9c, 0x47, 0x7b, 0x2c, 0xc0,
	0xe2, 0x19, 0xdf, 0xe3, 0x5f, 0x1d, 0x46, 0xfe, 0xc9, 0xa8, 0x41, 0x57, 0x4d, 0x9e, 0xe8, 0x96,
	0xb0, 0xd8, 0x1a, 0x6f, 0x63, 0xec, 0x6c, 0x62, 0x1b, 0x3f, 0x8e, 0x49, 0x36, 0x2f, 0x13, 0xdc,
	0x7a, 0x1c, 0x6c, 0xbc, 0xc7, 0xde, 0xf4, 0x0e, 0xfd, 0xfa, 0xea, 0xcf, 0xfd, 0x35, 0x9d, 0xdd,
	0x50, 0x25, 0x56, 0xfb, 0x63, 0xe9, 0xf2, 0x3b, 0x00, 0x00, 0xff, 0xff, 0xd3, 0x0e, 0xf7, 0xb6,
	0x24, 0x01, 0x00, 0x00,
}
