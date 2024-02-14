// Copyright 2024 Google LLC
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
	"testing"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/genproto/googleapis/api/serviceconfig"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func TestAutoPopulatedFields(t *testing.T) {
	optsUUID4 := &descriptorpb.FieldOptions{}
	proto.SetExtension(optsUUID4, annotations.E_FieldInfo, &annotations.FieldInfo{Format: annotations.FieldInfo_UUID4})
	optsIPV4 := &descriptorpb.FieldOptions{}
	proto.SetExtension(optsIPV4, annotations.E_FieldInfo, &annotations.FieldInfo{Format: annotations.FieldInfo_IPV4})
	inputType := &descriptor.DescriptorProto{
		Name: proto.String("InputType"),
		Field: []*descriptor.FieldDescriptorProto{
			{
				Name:           proto.String("request_id"),
				Type:           typep(descriptor.FieldDescriptorProto_TYPE_STRING),
				Label:          labelp(descriptor.FieldDescriptorProto_LABEL_OPTIONAL),
				Proto3Optional: proto.Bool(true),
				Options:        optsUUID4,
			},
			{
				Name:  proto.String("invalid_auto_populated_not_in_serviceconfig"),
				Type:  typep(descriptor.FieldDescriptorProto_TYPE_STRING),
				Label: labelp(descriptor.FieldDescriptorProto_LABEL_OPTIONAL),
			},
			{
				Name:  proto.String("invalid_auto_populated_no_annotation"),
				Type:  typep(descriptor.FieldDescriptorProto_TYPE_STRING),
				Label: labelp(descriptor.FieldDescriptorProto_LABEL_OPTIONAL),
			},
			{
				Name:    proto.String("invalid_auto_populated_required"),
				Type:    typep(descriptor.FieldDescriptorProto_TYPE_STRING),
				Label:   labelp(descriptor.FieldDescriptorProto_LABEL_REQUIRED),
				Options: optsUUID4,
			},
			{
				Name:    proto.String("invalid_auto_populated_int"),
				Type:    typep(descriptor.FieldDescriptorProto_TYPE_INT64),
				Label:   labelp(descriptor.FieldDescriptorProto_LABEL_OPTIONAL),
				Options: optsUUID4,
			},
			{
				Name:    proto.String("invalid_auto_populated_ipv4"),
				Type:    typep(descriptor.FieldDescriptorProto_TYPE_STRING),
				Label:   labelp(descriptor.FieldDescriptorProto_LABEL_OPTIONAL),
				Options: optsIPV4,
			},
		},
	}
	outputType := &descriptor.DescriptorProto{
		Name: proto.String("OutputType"),
	}
	file := &descriptor.FileDescriptorProto{
		Package: proto.String("my.pkg"),
		Options: &descriptor.FileOptions{
			GoPackage: proto.String("mypackage"),
		},
	}
	serv := &descriptor.ServiceDescriptorProto{
		Name: proto.String("Foo"),
	}

	var g generator
	g.opts = &options{
		pkgName: "pkg",
	}
	g.imports = map[pbinfo.ImportSpec]bool{}

	g.serviceConfig = &serviceconfig.Service{
		Publishing: &annotations.Publishing{
			MethodSettings: []*annotations.MethodSettings{
				{
					Selector: "my.pkg.Foo.GetOneThing",
					AutoPopulatedFields: []string{
						"request_id",
						"invalid_auto_populated_no_annotation",
						"invalid_auto_populated_required",
						"invalid_auto_populated_int",
						"invalid_auto_populated_ipv4",
					},
				},
			},
		},
	}

	commonTypes(&g)
	g.descInfo.Type[".my.pkg."+inputType.GetName()] = inputType
	g.descInfo.ParentFile[inputType] = file
	g.descInfo.Type[".my.pkg."+outputType.GetName()] = outputType
	g.descInfo.ParentFile[outputType] = file
	g.descInfo.ParentFile[serv] = file

	m := &descriptor.MethodDescriptorProto{
		Name:       proto.String("GetOneThing"),
		InputType:  proto.String(".my.pkg.InputType"),
		OutputType: proto.String(".my.pkg.OutputType"),
		Options:    &descriptor.MethodOptions{},
	}
	g.descInfo.ParentElement[m] = serv
	serv.Method = []*descriptor.MethodDescriptorProto{m}

	got := g.autoPopulatedFields(serv.GetName(), m)
	if want := 1; len(got) != want {
		t.Errorf("len(got) = %d, want: %d, got: %v", len(got), want, got)
	}
	if want := "request_id"; got[0].GetName() != want {
		t.Errorf("got[0].GetName() = %s, want: %s", got[0].GetName(), want)
	}
	if !got[0].GetProto3Optional() {
		t.Error("got[0].GetProto3Optional() = false, want: true")
	}
}
