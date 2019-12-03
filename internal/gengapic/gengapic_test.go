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
	"path/filepath"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"github.com/googleapis/gapic-generator-go/internal/txtdiff"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/genproto/googleapis/longrunning"
)

func TestComment(t *testing.T) {
	var g generator

	for _, tst := range []struct {
		in, want string
	}{
		{
			in:   "",
			want: "",
		},
		{
			in:   "abc\ndef\n",
			want: "// abc\n// def\n",
		},
	} {
		g.pt.Reset()
		g.comment(tst.in)
		if got := g.pt.String(); got != tst.want {
			t.Errorf("comment(%q) = %q, want %q", tst.in, got, tst.want)
		}
	}
}

func TestMethodDoc(t *testing.T) {
	m := &descriptor.MethodDescriptorProto{
		Name: proto.String("MyMethod"),
	}

	var g generator
	g.comments = make(map[proto.Message]string)

	for _, tst := range []struct {
		in, want string
	}{
		{
			in:   "",
			want: "",
		},
		{
			in:   "Does stuff.\n It also does other stuffs.",
			want: "// MyMethod does stuff.\n// It also does other stuffs.\n",
		},
	} {
		g.comments[m] = tst.in
		g.pt.Reset()
		g.methodDoc(m)
		if got := g.pt.String(); got != tst.want {
			t.Errorf("comment(%q) = %q, want %q", tst.in, got, tst.want)
		}
	}
}

func TestReduceServName(t *testing.T) {
	for _, tst := range []struct {
		in, pkg, want string
	}{
		{"Foo", "", "Foo"},
		{"Foo", "foo", ""},

		{"FooV2", "", "Foo"},
		{"FooV2", "foo", ""},

		{"FooService", "", "Foo"},
		{"FooService", "foo", ""},

		{"FooServiceV2", "", "Foo"},
		{"FooServiceV2", "foo", ""},

		{"FooV2Bar", "", "FooV2Bar"},
	} {
		if got := pbinfo.ReduceServName(tst.in, tst.pkg); got != tst.want {
			t.Errorf("pbinfo.ReduceServName(%q, %q) = %q, want %q", tst.in, tst.pkg, got, tst.want)
		}
	}
}

func TestGRPCClientField(t *testing.T) {
	for _, tst := range []struct {
		in, pkg, want string
	}{
		{"Foo", "foo", "client"},
		{"FooV2", "foo", "client"},
		{"FooService", "foo", "client"},
		{"FooServiceV2", "foo", "client"},
		{"FooV2Bar", "", "fooV2BarClient"},
	} {
		if got := grpcClientField(pbinfo.ReduceServName(tst.in, tst.pkg)); got != tst.want {
			t.Errorf("grpcClientField(pbinfo.ReduceServName(%q, %q)) = %q, want %q", tst.in, tst.pkg, got, tst.want)
		}
	}
}

func TestGenMethod(t *testing.T) {
	inputType := &descriptor.DescriptorProto{
		Name: proto.String("InputType"),
	}
	outputType := &descriptor.DescriptorProto{
		Name: proto.String("OutputType"),
	}

	typep := func(t descriptor.FieldDescriptorProto_Type) *descriptor.FieldDescriptorProto_Type {
		return &t
	}
	labelp := func(l descriptor.FieldDescriptorProto_Label) *descriptor.FieldDescriptorProto_Label {
		return &l
	}

	pageInputType := &descriptor.DescriptorProto{
		Name: proto.String("PageInputType"),
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
	}
	paginatedField := &descriptor.FieldDescriptorProto{
		Name:  proto.String("items"),
		Type:  typep(descriptor.FieldDescriptorProto_TYPE_STRING),
		Label: labelp(descriptor.FieldDescriptorProto_LABEL_REPEATED),
	}
	pageOutputType := &descriptor.DescriptorProto{
		Name: proto.String("PageOutputType"),
		Field: []*descriptor.FieldDescriptorProto{
			{
				Name:  proto.String("next_page_token"),
				Type:  typep(descriptor.FieldDescriptorProto_TYPE_STRING),
				Label: labelp(descriptor.FieldDescriptorProto_LABEL_OPTIONAL),
			},
			paginatedField,
		},
	}

	opts := &descriptor.MethodOptions{}
	ext := &annotations.HttpRule{
		Pattern: &annotations.HttpRule_Get{
			Get: "/v1/{field_name.nested=projects/*/foo/*}/bars/{other=bar/*/baz/*}/buz",
		},
		AdditionalBindings: []*annotations.HttpRule{
			&annotations.HttpRule{
				Pattern: &annotations.HttpRule_Post{
					Post: "/v1/{field_name.nested=projects/*/foo/*}/bars/{another=bar/*/baz/*}/buz",
				},
			},
		},
	}
	proto.SetExtension(opts, annotations.E_Http, ext)

	file := &descriptor.FileDescriptorProto{
		Package: proto.String("my.pkg"),
		Options: &descriptor.FileOptions{
			GoPackage: proto.String("mypackage"),
		},
	}
	serv := &descriptor.ServiceDescriptorProto{}

	var g generator
	g.imports = map[pbinfo.ImportSpec]bool{}

	commonTypes(&g)
	for _, typ := range []*descriptor.DescriptorProto{
		inputType, outputType, pageInputType, pageOutputType,
	} {
		g.descInfo.Type[".my.pkg."+*typ.Name] = typ
		g.descInfo.ParentFile[typ] = file
	}
	g.descInfo.ParentFile[serv] = file
	g.descInfo.ParentElement = map[pbinfo.ProtoType]pbinfo.ProtoType{
		paginatedField: pageOutputType,
	}

	meths := []*descriptor.MethodDescriptorProto{
		{
			Name:       proto.String("GetEmptyThing"),
			InputType:  proto.String(".my.pkg.InputType"),
			OutputType: proto.String(emptyType),
			Options:    opts,
		},
		{
			Name:       proto.String("GetOneThing"),
			InputType:  proto.String(".my.pkg.InputType"),
			OutputType: proto.String(".my.pkg.OutputType"),
			Options:    opts,
		},
		{
			Name:       proto.String("GetManyThings"),
			InputType:  proto.String(".my.pkg.PageInputType"),
			OutputType: proto.String(".my.pkg.PageOutputType"),
			Options:    opts,
		},
		{
			Name:            proto.String("ServerThings"),
			InputType:       proto.String(".my.pkg.InputType"),
			OutputType:      proto.String(".my.pkg.OutputType"),
			ServerStreaming: proto.Bool(true),
			Options:         opts,
		},
		{
			Name:            proto.String("ClientThings"),
			InputType:       proto.String(".my.pkg.InputType"),
			OutputType:      proto.String(".my.pkg.OutputType"),
			ClientStreaming: proto.Bool(true),
			Options:         opts,
		},
		{
			Name:            proto.String("BidiThings"),
			InputType:       proto.String(".my.pkg.InputType"),
			OutputType:      proto.String(".my.pkg.OutputType"),
			ServerStreaming: proto.Bool(true),
			ClientStreaming: proto.Bool(true),
			Options:         opts,
		},
	}

methods:
	for _, m := range meths {
		g.pt.Reset()

		g.aux = &auxTypes{
			iters: map[string]*iterType{},
		}
		if err := g.genMethod("Foo", serv, m); err != nil {
			t.Error(err)
			continue
		}

		for _, m := range g.aux.lros {
			if err := g.lroType("MyService", serv, m); err != nil {
				t.Error(err)
				continue methods
			}
		}

		for _, iter := range g.aux.iters {
			g.pagingIter(iter)
		}

		txtdiff.Diff(t, m.GetName(), g.pt.String(), filepath.Join("testdata", "method_"+m.GetName()+".want"))
	}
}

func TestGenLRO(t *testing.T) {
	inputType := &descriptor.DescriptorProto{
		Name: proto.String("InputType"),
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
	serv := &descriptor.ServiceDescriptorProto{}

	var g generator
	g.imports = map[pbinfo.ImportSpec]bool{}

	commonTypes(&g)
	for _, typ := range []*descriptor.DescriptorProto{
		inputType, outputType,
	} {
		g.descInfo.Type[".my.pkg."+*typ.Name] = typ
		g.descInfo.ParentFile[typ] = file
	}
	g.descInfo.ParentFile[serv] = file

	emptyLRO := &longrunning.OperationInfo{
		ResponseType: emptyValue,
	}
	emptyLROOpts := &descriptor.MethodOptions{}
	proto.SetExtension(emptyLROOpts, longrunning.E_OperationInfo, emptyLRO)

	respLRO := &longrunning.OperationInfo{
		ResponseType: "OutputType",
	}
	respLROOpts := &descriptor.MethodOptions{}
	proto.SetExtension(respLROOpts, longrunning.E_OperationInfo, respLRO)

	lros := []*descriptor.MethodDescriptorProto{
		{
			Name:       proto.String("EmptyLRO"),
			InputType:  proto.String(".my.pkg.InputType"),
			OutputType: proto.String(".google.longrunning.Operation"),
			Options:    emptyLROOpts,
		},
		{
			Name:       proto.String("RespLRO"),
			InputType:  proto.String(".my.pkg.InputType"),
			OutputType: proto.String(".google.longrunning.Operation"),
			Options:    respLROOpts,
		},
	}

lros:
	for _, m := range lros {
		g.pt.Reset()

		g.aux = &auxTypes{}

		if err := g.genMethod("Foo", serv, m); err != nil {
			t.Error(err)
			continue
		}

		for _, m := range g.aux.lros {
			if err := g.lroType("MyService", serv, m); err != nil {
				t.Error(err)
				continue lros
			}
		}

		txtdiff.Diff(t, m.GetName(), g.pt.String(), filepath.Join("testdata", "method_"+m.GetName()+".want"))
	}
}

func Test_buildAccessor(t *testing.T) {
	tests := []struct {
		name  string
		field string
		want  string
	}{
		{name: "simple", field: "foo_foo", want: ".GetFooFoo()"},
		{name: "nested", field: "foo_foo.bar_bar", want: ".GetFooFoo().GetBarBar()"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildAccessor(tt.field); got != tt.want {
				t.Errorf("buildAccessor() = %v, want %v", got, tt.want)
			}
		})
	}
}
