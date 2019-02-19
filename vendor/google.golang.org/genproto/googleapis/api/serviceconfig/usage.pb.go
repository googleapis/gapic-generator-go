// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/api/usage.proto

package serviceconfig

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Configuration controlling usage of a service.
type Usage struct {
	// Requirements that must be satisfied before a consumer project can use the
	// service. Each requirement is of the form <service.name>/<requirement-id>;
	// for example 'serviceusage.googleapis.com/billing-enabled'.
	Requirements []string `protobuf:"bytes,1,rep,name=requirements,proto3" json:"requirements,omitempty"`
	// A list of usage rules that apply to individual API methods.
	//
	// **NOTE:** All service configuration rules follow "last one wins" order.
	Rules []*UsageRule `protobuf:"bytes,6,rep,name=rules,proto3" json:"rules,omitempty"`
	// The full resource name of a channel used for sending notifications to the
	// service producer.
	//
	// Google Service Management currently only supports
	// [Google Cloud Pub/Sub](https://cloud.google.com/pubsub) as a notification
	// channel. To use Google Cloud Pub/Sub as the channel, this must be the name
	// of a Cloud Pub/Sub topic that uses the Cloud Pub/Sub topic name format
	// documented in https://cloud.google.com/pubsub/docs/overview.
	ProducerNotificationChannel string   `protobuf:"bytes,7,opt,name=producer_notification_channel,json=producerNotificationChannel,proto3" json:"producer_notification_channel,omitempty"`
	XXX_NoUnkeyedLiteral        struct{} `json:"-"`
	XXX_unrecognized            []byte   `json:"-"`
	XXX_sizecache               int32    `json:"-"`
}

func (m *Usage) Reset()         { *m = Usage{} }
func (m *Usage) String() string { return proto.CompactTextString(m) }
func (*Usage) ProtoMessage()    {}
func (*Usage) Descriptor() ([]byte, []int) {
	return fileDescriptor_701aa74a03c68f0a, []int{0}
}

func (m *Usage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Usage.Unmarshal(m, b)
}
func (m *Usage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Usage.Marshal(b, m, deterministic)
}
func (m *Usage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Usage.Merge(m, src)
}
func (m *Usage) XXX_Size() int {
	return xxx_messageInfo_Usage.Size(m)
}
func (m *Usage) XXX_DiscardUnknown() {
	xxx_messageInfo_Usage.DiscardUnknown(m)
}

var xxx_messageInfo_Usage proto.InternalMessageInfo

func (m *Usage) GetRequirements() []string {
	if m != nil {
		return m.Requirements
	}
	return nil
}

func (m *Usage) GetRules() []*UsageRule {
	if m != nil {
		return m.Rules
	}
	return nil
}

func (m *Usage) GetProducerNotificationChannel() string {
	if m != nil {
		return m.ProducerNotificationChannel
	}
	return ""
}

// Usage configuration rules for the service.
//
// NOTE: Under development.
//
//
// Use this rule to configure unregistered calls for the service. Unregistered
// calls are calls that do not contain consumer project identity.
// (Example: calls that do not contain an API key).
// By default, API methods do not allow unregistered calls, and each method call
// must be identified by a consumer project identity. Use this rule to
// allow/disallow unregistered calls.
//
// Example of an API that wants to allow unregistered calls for entire service.
//
//     usage:
//       rules:
//       - selector: "*"
//         allow_unregistered_calls: true
//
// Example of a method that wants to allow unregistered calls.
//
//     usage:
//       rules:
//       - selector: "google.example.library.v1.LibraryService.CreateBook"
//         allow_unregistered_calls: true
type UsageRule struct {
	// Selects the methods to which this rule applies. Use '*' to indicate all
	// methods in all APIs.
	//
	// Refer to [selector][google.api.DocumentationRule.selector] for syntax details.
	Selector string `protobuf:"bytes,1,opt,name=selector,proto3" json:"selector,omitempty"`
	// If true, the selected method allows unregistered calls, e.g. calls
	// that don't identify any user or application.
	AllowUnregisteredCalls bool `protobuf:"varint,2,opt,name=allow_unregistered_calls,json=allowUnregisteredCalls,proto3" json:"allow_unregistered_calls,omitempty"`
	// If true, the selected method should skip service control and the control
	// plane features, such as quota and billing, will not be available.
	// This flag is used by Google Cloud Endpoints to bypass checks for internal
	// methods, such as service health check methods.
	SkipServiceControl   bool     `protobuf:"varint,3,opt,name=skip_service_control,json=skipServiceControl,proto3" json:"skip_service_control,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UsageRule) Reset()         { *m = UsageRule{} }
func (m *UsageRule) String() string { return proto.CompactTextString(m) }
func (*UsageRule) ProtoMessage()    {}
func (*UsageRule) Descriptor() ([]byte, []int) {
	return fileDescriptor_701aa74a03c68f0a, []int{1}
}

func (m *UsageRule) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UsageRule.Unmarshal(m, b)
}
func (m *UsageRule) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UsageRule.Marshal(b, m, deterministic)
}
func (m *UsageRule) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UsageRule.Merge(m, src)
}
func (m *UsageRule) XXX_Size() int {
	return xxx_messageInfo_UsageRule.Size(m)
}
func (m *UsageRule) XXX_DiscardUnknown() {
	xxx_messageInfo_UsageRule.DiscardUnknown(m)
}

var xxx_messageInfo_UsageRule proto.InternalMessageInfo

func (m *UsageRule) GetSelector() string {
	if m != nil {
		return m.Selector
	}
	return ""
}

func (m *UsageRule) GetAllowUnregisteredCalls() bool {
	if m != nil {
		return m.AllowUnregisteredCalls
	}
	return false
}

func (m *UsageRule) GetSkipServiceControl() bool {
	if m != nil {
		return m.SkipServiceControl
	}
	return false
}

func init() {
	proto.RegisterType((*Usage)(nil), "google.api.Usage")
	proto.RegisterType((*UsageRule)(nil), "google.api.UsageRule")
}

func init() { proto.RegisterFile("google/api/usage.proto", fileDescriptor_701aa74a03c68f0a) }

var fileDescriptor_701aa74a03c68f0a = []byte{
	// 331 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x91, 0xc1, 0x4b, 0xfb, 0x30,
	0x14, 0xc7, 0xe9, 0xf6, 0xdb, 0x7e, 0x5b, 0x14, 0x0f, 0x41, 0x47, 0x99, 0x0a, 0x65, 0xa7, 0x82,
	0xd0, 0x8a, 0x5e, 0x04, 0x4f, 0x6e, 0x88, 0x78, 0x91, 0x51, 0xd9, 0xc5, 0x4b, 0x89, 0xd9, 0x5b,
	0x0c, 0x66, 0x79, 0x35, 0x49, 0xf5, 0x0f, 0xf1, 0xea, 0xc9, 0xbf, 0x54, 0x9a, 0xcc, 0xd9, 0x1d,
	0xdf, 0xfb, 0x7c, 0xbe, 0xef, 0xb5, 0x2f, 0x64, 0x24, 0x10, 0x85, 0x82, 0x9c, 0x55, 0x32, 0xaf,
	0x2d, 0x13, 0x90, 0x55, 0x06, 0x1d, 0x52, 0x12, 0xfa, 0x19, 0xab, 0xe4, 0xf8, 0xa4, 0xe5, 0x30,
	0xad, 0xd1, 0x31, 0x27, 0x51, 0xdb, 0x60, 0x4e, 0xbe, 0x22, 0xd2, 0x5b, 0x34, 0x49, 0x3a, 0x21,
	0xfb, 0x06, 0xde, 0x6a, 0x69, 0x60, 0x0d, 0xda, 0xd9, 0x38, 0x4a, 0xba, 0xe9, 0xb0, 0xd8, 0xe9,
	0xd1, 0x33, 0xd2, 0x33, 0xb5, 0x02, 0x1b, 0xf7, 0x93, 0x6e, 0xba, 0x77, 0x71, 0x94, 0xfd, 0xed,
	0xc9, 0xfc, 0x94, 0xa2, 0x56, 0x50, 0x04, 0x87, 0x4e, 0xc9, 0x69, 0x65, 0x70, 0x59, 0x73, 0x30,
	0xa5, 0x46, 0x27, 0x57, 0x92, 0xfb, 0xd5, 0x25, 0x7f, 0x61, 0x5a, 0x83, 0x8a, 0xff, 0x27, 0x51,
	0x3a, 0x2c, 0x8e, 0x7f, 0xa5, 0x87, 0x96, 0x33, 0x0b, 0xca, 0xe4, 0x33, 0x22, 0xc3, 0xed, 0x60,
	0x3a, 0x26, 0x03, 0x0b, 0x0a, 0xb8, 0x43, 0x13, 0x47, 0x3e, 0xbc, 0xad, 0xe9, 0x15, 0x89, 0x99,
	0x52, 0xf8, 0x51, 0xd6, 0xda, 0x80, 0x90, 0xd6, 0x81, 0x81, 0x65, 0xc9, 0x99, 0x52, 0x36, 0xee,
	0x24, 0x51, 0x3a, 0x28, 0x46, 0x9e, 0x2f, 0x5a, 0x78, 0xd6, 0x50, 0x7a, 0x4e, 0x0e, 0xed, 0xab,
	0xac, 0x4a, 0x0b, 0xe6, 0x5d, 0x72, 0x28, 0x39, 0x6a, 0x67, 0x50, 0xc5, 0x5d, 0x9f, 0xa2, 0x0d,
	0x7b, 0x0c, 0x68, 0x16, 0xc8, 0x54, 0x91, 0x03, 0x8e, 0xeb, 0xd6, 0xcf, 0x4f, 0x89, 0xff, 0xc8,
	0x79, 0x73, 0xd2, 0x79, 0xf4, 0x74, 0xbb, 0x21, 0x02, 0x15, 0xd3, 0x22, 0x43, 0x23, 0x72, 0x01,
	0xda, 0x1f, 0x3c, 0x0f, 0x88, 0x55, 0xd2, 0xfa, 0x17, 0xd9, 0x2c, 0xe5, 0xa8, 0x57, 0x52, 0x5c,
	0xef, 0x54, 0xdf, 0x9d, 0x7f, 0x77, 0x37, 0xf3, 0xfb, 0xe7, 0xbe, 0x0f, 0x5e, 0xfe, 0x04, 0x00,
	0x00, 0xff, 0xff, 0x9c, 0x4b, 0x8c, 0x57, 0xed, 0x01, 0x00, 0x00,
}
