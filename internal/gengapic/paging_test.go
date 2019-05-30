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

package gengapic

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/google/go-cmp/cmp"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
)

func TestPagingField(t *testing.T) {
	typep := func(t descriptor.FieldDescriptorProto_Type) *descriptor.FieldDescriptorProto_Type {
		return &t
	}
	labelp := func(l descriptor.FieldDescriptorProto_Label) *descriptor.FieldDescriptorProto_Label {
		return &l
	}

	resField := &descriptor.FieldDescriptorProto{
		Name:  proto.String("resource"),
		Type:  typep(descriptor.FieldDescriptorProto_TYPE_STRING),
		Label: labelp(descriptor.FieldDescriptorProto_LABEL_REPEATED),
	}

	g := &generator{}
	g.descInfo.Type = map[string]pbinfo.ProtoType{
		"Foo": &descriptor.DescriptorProto{
			Name: proto.String("Foo"),
		},
		"PageIn": &descriptor.DescriptorProto{
			Name: proto.String("PageIn"),
			Field: []*descriptor.FieldDescriptorProto{
				{
					Name:  proto.String("page_size"),
					Type:  typep(descriptor.FieldDescriptorProto_TYPE_INT32),
					Label: labelp(descriptor.FieldDescriptorProto_LABEL_OPTIONAL),
				},
				{
					Name:  proto.String("page_token"),
					Type:  typep(descriptor.FieldDescriptorProto_TYPE_STRING),
					Label: labelp(descriptor.FieldDescriptorProto_LABEL_OPTIONAL),
				},
			},
		},
		"PageOut": &descriptor.DescriptorProto{
			Name: proto.String("PageOut"),
			Field: []*descriptor.FieldDescriptorProto{
				{
					Name:  proto.String("next_page_token"),
					Type:  typep(descriptor.FieldDescriptorProto_TYPE_STRING),
					Label: labelp(descriptor.FieldDescriptorProto_LABEL_OPTIONAL),
				},
				resField,
			},
		},
		"BadPageOut1": &descriptor.DescriptorProto{
			Name: proto.String("BadPageOut1"),
			Field: []*descriptor.FieldDescriptorProto{
				{
					Name:  proto.String("next_page_token"),
					Type:  typep(descriptor.FieldDescriptorProto_TYPE_STRING),
					Label: labelp(descriptor.FieldDescriptorProto_LABEL_OPTIONAL),
				},
				// No repeated field
			},
		},
		"BadPageOut2": &descriptor.DescriptorProto{
			Name: proto.String("BadPageOut2"),
			Field: []*descriptor.FieldDescriptorProto{
				{
					Name:  proto.String("next_page_token"),
					Type:  typep(descriptor.FieldDescriptorProto_TYPE_STRING),
					Label: labelp(descriptor.FieldDescriptorProto_LABEL_OPTIONAL),
				},
				// Too many repeated field
				resField,
				resField,
			},
		},
	}

	for _, tst := range []struct {
		in, out string
		field   *descriptor.FieldDescriptorProto
		err     bool
	}{
		{
			in:  "Foo",
			out: "Foo",
		},
		{
			in:  "PageIn",
			out: "Foo",
		},
		{
			in:  "Foo",
			out: "PageOut",
		},
		{
			in:    "PageIn",
			out:   "PageOut",
			field: resField,
		},
		{
			in:  "PageIn",
			out: "BadPageOut1",
			err: true,
		},
		{
			in:  "PageIn",
			out: "BadPageOut2",
			err: true,
		},
	} {
		meth := &descriptor.MethodDescriptorProto{
			Name:       proto.String("TestPagingField"),
			InputType:  proto.String(tst.in),
			OutputType: proto.String(tst.out),
		}
		f, err := g.pagingField(meth)
		if tst.err && err == nil {
			t.Errorf("pagingField(%v)=%v, expected error", meth, f)
		} else if !tst.err && err != nil {
			t.Errorf("pagingField(%v) errors %q, expected %v", meth, err, tst.field)
		} else if f != tst.field {
			t.Errorf("pagingField(%v)=%v, want %v", meth, f, tst.field)
		}
	}
}

func TestIterTypeOf(t *testing.T) {
	typep := func(t descriptor.FieldDescriptorProto_Type) *descriptor.FieldDescriptorProto_Type {
		return &t
	}
	msgType := &descriptor.DescriptorProto{
		Name: proto.String("Foo"),
	}
	g := &generator{
		descInfo: pbinfo.Info{
			Type: map[string]pbinfo.ProtoType{
				msgType.GetName(): msgType,
			},
			ParentElement: map[pbinfo.ProtoType]pbinfo.ProtoType{},
			ParentFile: map[proto.Message]*descriptor.FileDescriptorProto{
				msgType: &descriptor.FileDescriptorProto{
					Options: &descriptor.FileOptions{
						GoPackage: proto.String("path/to/foo;foo"),
					},
				},
			},
		},
	}

	for i, tst := range []struct {
		field *descriptor.FieldDescriptorProto
		want  iterType
	}{
		{
			field: &descriptor.FieldDescriptorProto{
				Type: typep(descriptor.FieldDescriptorProto_TYPE_STRING),
			},
			want: iterType{
				iterTypeName: "StringIterator",
				elemTypeName: "string",
			},
		},
		{
			field: &descriptor.FieldDescriptorProto{
				Type: typep(descriptor.FieldDescriptorProto_TYPE_BYTES),
			},
			want: iterType{
				iterTypeName: "BytesIterator",
				elemTypeName: "[]byte",
			},
		},
		{
			field: &descriptor.FieldDescriptorProto{
				Type:     typep(descriptor.FieldDescriptorProto_TYPE_MESSAGE),
				TypeName: proto.String(msgType.GetName()),
			},
			want: iterType{
				iterTypeName: "FooIterator",
				elemTypeName: "*foopb.Foo",
				elemImports:  []pbinfo.ImportSpec{{Name: "foopb", Path: "path/to/foo"}},
			},
		},
	} {
		g.descInfo.ParentElement[tst.field] = msgType
		got, err := g.iterTypeOf(tst.field)
		if err != nil {
			t.Error(err)
		} else if diff := cmp.Diff(tst.want, got, cmp.AllowUnexported(got)); diff != "" {
			t.Errorf("%d: (got=-, want=+):\n%s", i, diff)
		}
	}
}
