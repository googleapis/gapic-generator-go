// Copyright 2026 Google LLC
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

	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func TestResourceNameField(t *testing.T) {
	optsResRef := &descriptorpb.FieldOptions{}
	proto.SetExtension(optsResRef, annotations.E_ResourceReference, &annotations.ResourceReference{
		Type: "foo.googleapis.com/Bar",
	})

	extraMsg := &descriptorpb.DescriptorProto{
		Name: proto.String("ExtraMessage"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:    proto.String("name"),
				Type:    typep(descriptorpb.FieldDescriptorProto_TYPE_STRING),
				Options: optsResRef,
			},
		},
	}

	inputType := &descriptorpb.DescriptorProto{
		Name: proto.String("InputType"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name: proto.String("top_level"),
				Type: typep(descriptorpb.FieldDescriptorProto_TYPE_STRING),
			},
			{
				Name:     proto.String("nested"),
				Type:     typep(descriptorpb.FieldDescriptorProto_TYPE_MESSAGE),
				TypeName: proto.String(".my.pkg.ExtraMessage"),
			},
		},
	}

	inputTypeWithTopLevel := &descriptorpb.DescriptorProto{
		Name: proto.String("InputTypeWithTopLevel"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:    proto.String("name"),
				Type:    typep(descriptorpb.FieldDescriptorProto_TYPE_STRING),
				Options: optsResRef,
			},
			{
				Name:     proto.String("nested"),
				Type:     typep(descriptorpb.FieldDescriptorProto_TYPE_MESSAGE),
				TypeName: proto.String(".my.pkg.ExtraMessage"),
			},
		},
	}

	for _, tst := range []struct {
		name       string
		inputType  *descriptorpb.DescriptorProto
		httpRule   *annotations.HttpRule
		wantField  string
		extraTypes map[string]pbinfo.ProtoType
	}{
		{
			name:      "nested_resource_field",
			inputType: inputType,
			wantField: "nested.name",
			extraTypes: map[string]pbinfo.ProtoType{
				".my.pkg.ExtraMessage": extraMsg,
			},
		},
		{
			name:      "top_level_resource_field",
			inputType: inputTypeWithTopLevel,
			wantField: "name",
			extraTypes: map[string]pbinfo.ProtoType{
				".my.pkg.ExtraMessage": extraMsg,
			},
		},
		{
			name:      "tie_breaking_with_http_path",
			inputType: inputTypeWithTopLevel,
			httpRule: &annotations.HttpRule{
				Pattern: &annotations.HttpRule_Get{
					Get: "/v1/{nested.name=*}",
				},
			},
			wantField: "nested.name",
			extraTypes: map[string]pbinfo.ProtoType{
				".my.pkg.ExtraMessage": extraMsg,
			},
		},
	} {
		t.Run(tst.name, func(t *testing.T) {
			m := &descriptorpb.MethodDescriptorProto{
				Name:      proto.String("MyMethod"),
				InputType: proto.String(".my.pkg." + tst.inputType.GetName()),
				Options:   &descriptorpb.MethodOptions{},
			}
			if tst.httpRule != nil {
				proto.SetExtension(m.Options, annotations.E_Http, tst.httpRule)
			}

			g := &generator{
				descInfo: pbinfo.Info{
					Type: map[string]pbinfo.ProtoType{
						".my.pkg." + tst.inputType.GetName(): tst.inputType,
					},
				},
			}
			for k, v := range tst.extraTypes {
				g.descInfo.Type[k] = v
			}

			if got := g.resourceNameField(m); got != tst.wantField {
				t.Errorf("resourceNameField(%s) = %q, want %q", tst.name, got, tst.wantField)
			}
		})
	}
}
