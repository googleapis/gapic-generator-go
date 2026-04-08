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

	"github.com/google/go-cmp/cmp"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/genproto/googleapis/api/serviceconfig"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/descriptorpb"
)

func TestAutoPopulatedFields(t *testing.T) {
	optsUUID4 := &descriptorpb.FieldOptions{}
	proto.SetExtension(optsUUID4, annotations.E_FieldInfo, &annotations.FieldInfo{Format: annotations.FieldInfo_UUID4})

	optsIPV4 := &descriptorpb.FieldOptions{}
	proto.SetExtension(optsIPV4, annotations.E_FieldInfo, &annotations.FieldInfo{Format: annotations.FieldInfo_IPV4})

	optsRequiredAndUUID4 := &descriptorpb.FieldOptions{}
	proto.SetExtension(optsRequiredAndUUID4, annotations.E_FieldBehavior, []annotations.FieldBehavior{annotations.FieldBehavior_REQUIRED})
	proto.SetExtension(optsRequiredAndUUID4, annotations.E_FieldInfo, &annotations.FieldInfo{Format: annotations.FieldInfo_UUID4})

	requestIDField := &descriptorpb.FieldDescriptorProto{
		Name:           proto.String("request_id"),
		Type:           typep(descriptorpb.FieldDescriptorProto_TYPE_STRING),
		Label:          labelp(descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL),
		Proto3Optional: proto.Bool(true),
		Options:        optsUUID4,
	}
	inputType := &descriptorpb.DescriptorProto{
		Name: proto.String("InputType"),
		Field: []*descriptorpb.FieldDescriptorProto{
			requestIDField,
			{
				Name:  proto.String("invalid_auto_populated_not_in_serviceconfig"),
				Type:  typep(descriptorpb.FieldDescriptorProto_TYPE_STRING),
				Label: labelp(descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL),
			},
			{
				Name:  proto.String("invalid_auto_populated_no_annotation"),
				Type:  typep(descriptorpb.FieldDescriptorProto_TYPE_STRING),
				Label: labelp(descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL),
			},
			{
				Name:    proto.String("invalid_auto_populated_required"),
				Type:    typep(descriptorpb.FieldDescriptorProto_TYPE_STRING),
				Label:   labelp(descriptorpb.FieldDescriptorProto_LABEL_REQUIRED),
				Options: optsRequiredAndUUID4,
			},
			{
				Name:    proto.String("invalid_auto_populated_int"),
				Type:    typep(descriptorpb.FieldDescriptorProto_TYPE_INT64),
				Label:   labelp(descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL),
				Options: optsUUID4,
			},
			{
				Name:    proto.String("invalid_auto_populated_ipv4"),
				Type:    typep(descriptorpb.FieldDescriptorProto_TYPE_STRING),
				Label:   labelp(descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL),
				Options: optsIPV4,
			},
		},
	}
	outputType := &descriptorpb.DescriptorProto{
		Name: proto.String("OutputType"),
	}
	file := &descriptorpb.FileDescriptorProto{
		Package: proto.String("my.pkg"),
		Options: &descriptorpb.FileOptions{
			GoPackage: proto.String("mypackage"),
		},
	}
	serv := &descriptorpb.ServiceDescriptorProto{
		Name: proto.String("Foo"),
	}

	var g generator
	g.cfg = &generatorConfig{
		pkgName: "pkg",
	}
	g.imports = map[pbinfo.ImportSpec]bool{}

	g.cfg.APIServiceConfig = &serviceconfig.Service{
		Name: "foo.googleapis.com",
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

	m := &descriptorpb.MethodDescriptorProto{
		Name:       proto.String("GetOneThing"),
		InputType:  proto.String(".my.pkg.InputType"),
		OutputType: proto.String(".my.pkg.OutputType"),
		Options:    &descriptorpb.MethodOptions{},
	}
	g.descInfo.ParentElement[m] = serv
	serv.Method = []*descriptorpb.MethodDescriptorProto{m}

	got := g.autoPopulatedFields(serv.GetName(), m)

	want := []*descriptorpb.FieldDescriptorProto{requestIDField}
	if diff := cmp.Diff(got, want, protocmp.Transform()); diff != "" {
		t.Errorf("got(-),want(+):\n%s", diff)
	}
}

func TestSelectiveGapicGeneration(t *testing.T) {
	file := &descriptorpb.FileDescriptorProto{
		Package: proto.String("my.pkg"),
		Options: &descriptorpb.FileOptions{
			GoPackage: proto.String("mypackage"),
		},
	}
	m := &descriptorpb.MethodDescriptorProto{
		Name: proto.String("InternalMethod"),
	}
	m2 := &descriptorpb.MethodDescriptorProto{
		Name: proto.String("PublicMethod"),
	}
	serv := &descriptorpb.ServiceDescriptorProto{
		Name:   proto.String("Foo"),
		Method: []*descriptorpb.MethodDescriptorProto{m, m2},
	}

	var g generator
	g.cfg = &generatorConfig{
		pkgName: "pkg",
		featureEnablement: map[featureID]struct{}{
			SelectiveGapicGenerationFeature: {},
		},
	}
	g.descInfo.ParentElement = map[pbinfo.ProtoType]pbinfo.ProtoType{
		m:  serv,
		m2: serv,
	}
	g.descInfo.ParentFile = map[proto.Message]*descriptorpb.FileDescriptorProto{
		serv: file,
		m:    file,
		m2:   file,
	}

	if g.isInternalService(serv) {
		t.Errorf("expected isInternalService to be false when config is missing")
	}

	g.cfg.APIServiceConfig = &serviceconfig.Service{
		Publishing: &annotations.Publishing{
			LibrarySettings: []*annotations.ClientLibrarySettings{
				{
					GoSettings: &annotations.GoSettings{
						Common: &annotations.CommonLanguageSettings{
							SelectiveGapicGeneration: &annotations.SelectiveGapicGeneration{
								GenerateOmittedAsInternal: true,
								Methods:                   []string{"my.pkg.Foo.PublicMethod"},
							},
						},
					},
				},
			},
		},
	}

	if !g.isInternalService(serv) {
		t.Errorf("expected isInternalService to be true")
	}

	if g.methodName(m) != "internalMethod" {
		t.Errorf("expected methodName(m) to be 'internalMethod', got %q", g.methodName(m))
	}

	if g.methodName(m2) != "PublicMethod" {
		t.Errorf("expected methodName(m2) to be 'PublicMethod', got %q", g.methodName(m2))
	}

	if got := g.clientName(serv, "pkg"); got != "BaseFoo" {
		t.Errorf("expected clientName to be 'BaseFoo', got %q", got)
	}

	if got := g.clientName(serv, "Foo"); got != "Base" {
		t.Errorf("expected clientName to be 'Base', got %q", got)
	}
}
