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
	"path/filepath"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"github.com/googleapis/gapic-generator-go/internal/txtdiff"
)

func TestUnary(t *testing.T) {
	g := initTestGenerator()

	vs := SampleValueSet{
		ID: "my_value_set",
		Parameters: SampleParameter{
			Defaults: []string{
				`a.x = 42`,
				`a.y = 3.14159`,
				`b = "foobar"`,
				`e = FOO`,
				`f = "in a oneof"`,

				`a_array[0].x = 0`,
				`a_array[0].y = 1`,
				`a_array[1].x = 2`,
				`a_array[1].y = 3`,
			},
			Attributes: []SampleAttribute{
				{"a.x", "the_x"},
				{"b", "the_b"},
			},
		},
		OnSuccess: []OutputSpec{
			{Define: "out_a = $resp.a"},
			{Print: []string{"x = %s", "$resp.a.x"}},
			{Print: []string{"y = %s", "out_a.y"}},
			{
				Loop: &LoopSpec{
					Variable:   "r",
					Collection: "$resp.r",
					Body: []OutputSpec{
						{Print: []string{"resp.r contains %s", "r"}},
					},
				},
			},
		},
	}
	if err := g.genSample("foo.FooService", "UnaryMethod", "awesome_region", vs); err != nil {
		t.Fatal(err)
	}

	// Don't format. Format can change with Go version.
	gofmt := false
	year := 2018
	content, err := g.commit(gofmt, year)
	if err != nil {
		t.Fatal(err)
	}
	txtdiff.Diff(t, "TestUnary", string(content), filepath.Join("testdata", "sample_unary.want"))
}

func TestSample_InitError(t *testing.T) {
	g := initTestGenerator()
	for _, tst := range []string{
		// bad type
		`a = "string"`,
		`a.x = "string"`,
		`a_array[0].x = "string"`,

		// try to access array
		`a_array.x = 0`,
		// try to index singular
		`a[0].x = 0`,
		// missing ']'
		`a_array[0.x = 0`,
		// index can only be ints
		`a_array[0.0].x = 0`,
		`a_array["0"].x = 0`,

		// bad enum variant
		`e = DERP`,
	} {
		g.reset()

		vs := SampleValueSet{
			ID: "my_value_set",
			Parameters: SampleParameter{
				Defaults: []string{tst},
			},
		}

		if err := g.genSample("foo.FooService", "UnaryMethod", "awesome_region", vs); err == nil {
			t.Errorf("expected error from init config: %s", tst)
		}
	}
}

func TestPaging(t *testing.T) {
	g := initTestGenerator()

	vs := SampleValueSet{
		ID: "my_value_set",
	}
	if err := g.genSample("foo.FooService", "PagingMethod", "awesome_region", vs); err != nil {
		t.Fatal(err)
	}

	// Don't format. Format can change with Go version.
	gofmt := false
	year := 2018
	content, err := g.commit(gofmt, year)
	if err != nil {
		t.Fatal(err)
	}
	txtdiff.Diff(t, "TestUnary", string(content), filepath.Join("testdata", "sample_paging.want"))
}

func initTestGenerator() *generator {
	labelp := func(l descriptor.FieldDescriptorProto_Label) *descriptor.FieldDescriptorProto_Label {
		return &l
	}

	inType := &descriptor.DescriptorProto{
		Name: proto.String("InputType"),
		OneofDecl: []*descriptor.OneofDescriptorProto{
			{Name: proto.String("Group")},
		},
		Field: []*descriptor.FieldDescriptorProto{
			{Name: proto.String("a"), TypeName: proto.String(".foo.AType")},
			{Name: proto.String("a_array"), TypeName: proto.String(".foo.AType"), Label: labelp(descriptor.FieldDescriptorProto_LABEL_REPEATED)},
			{Name: proto.String("b"), Type: typep(descriptor.FieldDescriptorProto_TYPE_STRING)},
			{Name: proto.String("e"), TypeName: proto.String(".foo.AType.EType")},
			{Name: proto.String("f"), Type: typep(descriptor.FieldDescriptorProto_TYPE_STRING), OneofIndex: proto.Int32(0)},
			{Name: proto.String("r"), Type: typep(descriptor.FieldDescriptorProto_TYPE_STRING), Label: labelp(descriptor.FieldDescriptorProto_LABEL_REPEATED)},
		},
	}
	eType := &descriptor.EnumDescriptorProto{
		Name: proto.String("EType"),
		Value: []*descriptor.EnumValueDescriptorProto{
			{Name: proto.String("FOO")},
		},
	}
	aType := &descriptor.DescriptorProto{
		Name: proto.String("AType"),
		Field: []*descriptor.FieldDescriptorProto{
			{Name: proto.String("x"), Type: typep(descriptor.FieldDescriptorProto_TYPE_INT64)},
			{Name: proto.String("y"), Type: typep(descriptor.FieldDescriptorProto_TYPE_FLOAT)},
		},
		EnumType: []*descriptor.EnumDescriptorProto{eType},
	}

	pageInType := &descriptor.DescriptorProto{
		Name: proto.String("PageInType"),
		Field: []*descriptor.FieldDescriptorProto{
			{Name: proto.String("a"), TypeName: proto.String("AType")},
			{Name: proto.String("page_size"), Type: typep(descriptor.FieldDescriptorProto_TYPE_INT32)},
			{Name: proto.String("page_token"), Type: typep(descriptor.FieldDescriptorProto_TYPE_STRING)},
		},
	}
	pageOutType := &descriptor.DescriptorProto{
		Name: proto.String("PageOutType"),
		Field: []*descriptor.FieldDescriptorProto{
			{Name: proto.String("a"), TypeName: proto.String("AType"), Label: labelp(descriptor.FieldDescriptorProto_LABEL_REPEATED)},
			{Name: proto.String("next_page_token"), Type: typep(descriptor.FieldDescriptorProto_TYPE_STRING)},
		},
	}

	serv := &descriptor.ServiceDescriptorProto{
		Name: proto.String("FooService"),
		Method: []*descriptor.MethodDescriptorProto{
			{
				Name:       proto.String("UnaryMethod"),
				InputType:  proto.String(".foo.InputType"),
				OutputType: proto.String(".foo.InputType"),
			},
			{
				Name:       proto.String("PagingMethod"),
				InputType:  proto.String(".foo.PageInType"),
				OutputType: proto.String(".foo.PageOutType"),
			},
		},
	}
	file := &descriptor.FileDescriptorProto{
		Options: &descriptor.FileOptions{
			GoPackage: proto.String("path.to/pb/foo;foo"),
		},
		Package:     proto.String("foo"),
		Service:     []*descriptor.ServiceDescriptorProto{serv},
		MessageType: []*descriptor.DescriptorProto{inType, aType, pageInType, pageOutType},
	}

	return &generator{
		clientPkg: pbinfo.ImportSpec{Path: "path.to/client/foo", Name: "foo"},
		imports:   map[pbinfo.ImportSpec]bool{},
		descInfo:  pbinfo.Of([]*descriptor.FileDescriptorProto{file}),
	}
}
