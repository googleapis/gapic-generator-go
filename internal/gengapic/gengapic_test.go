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
	"bytes"
	"path/filepath"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/google/go-cmp/cmp"
	conf "github.com/googleapis/gapic-generator-go/internal/grpc_service_config"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"github.com/googleapis/gapic-generator-go/internal/txtdiff"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/pluginpb"
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

		// IAM should be replaced with Iam
		{"IAMCredentials", "credentials", "IamCredentials"},
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
	typep := func(t descriptor.FieldDescriptorProto_Type) *descriptor.FieldDescriptorProto_Type {
		return &t
	}
	labelp := func(l descriptor.FieldDescriptorProto_Label) *descriptor.FieldDescriptorProto_Label {
		return &l
	}

	extra := &descriptor.DescriptorProto{
		Name: proto.String("ExtraMessage"),
		Field: []*descriptor.FieldDescriptorProto{
			{
				Name:  proto.String("nested"),
				Type:  typep(descriptor.FieldDescriptorProto_TYPE_STRING),
				Label: labelp(descriptor.FieldDescriptorProto_LABEL_OPTIONAL),
			},
		},
	}

	inputType := &descriptor.DescriptorProto{
		Name: proto.String("InputType"),
		Field: []*descriptor.FieldDescriptorProto{
			{
				Name:  proto.String("other"),
				Type:  typep(descriptor.FieldDescriptorProto_TYPE_STRING),
				Label: labelp(descriptor.FieldDescriptorProto_LABEL_OPTIONAL),
			},
			{
				Name:  proto.String("another"),
				Type:  typep(descriptor.FieldDescriptorProto_TYPE_STRING),
				Label: labelp(descriptor.FieldDescriptorProto_LABEL_OPTIONAL),
			},
			{
				Name:  proto.String("biz"),
				Type:  typep(descriptor.FieldDescriptorProto_TYPE_DOUBLE),
				Label: labelp(descriptor.FieldDescriptorProto_LABEL_OPTIONAL),
			},
			{
				Name:     proto.String("field_name"),
				Type:     typep(descriptor.FieldDescriptorProto_TYPE_MESSAGE),
				Label:    labelp(descriptor.FieldDescriptorProto_LABEL_OPTIONAL),
				TypeName: proto.String(".my.pkg.ExtraMessage"),
			},
		},
	}
	outputType := &descriptor.DescriptorProto{
		Name: proto.String("OutputType"),
	}

	pageInputType := &descriptor.DescriptorProto{
		Name: proto.String("PageInputType"),
		Field: append(inputType.GetField(), &descriptor.FieldDescriptorProto{
			Name:  proto.String("page_size"),
			Type:  typep(descriptor.FieldDescriptorProto_TYPE_INT32),
			Label: labelp(descriptor.FieldDescriptorProto_LABEL_OPTIONAL),
		}, &descriptor.FieldDescriptorProto{
			Name:  proto.String("page_token"),
			Type:  typep(descriptor.FieldDescriptorProto_TYPE_STRING),
			Label: labelp(descriptor.FieldDescriptorProto_LABEL_OPTIONAL),
		}),
	}
	pageInputTypeOptional := &descriptor.DescriptorProto{
		Name: proto.String("PageInputTypeOptional"),
		Field: append(inputType.GetField(), &descriptor.FieldDescriptorProto{
			Name:           proto.String("page_size"),
			Type:           typep(descriptor.FieldDescriptorProto_TYPE_INT32),
			Label:          labelp(descriptor.FieldDescriptorProto_LABEL_OPTIONAL),
			Proto3Optional: proto.Bool(true),
		}, &descriptor.FieldDescriptorProto{
			Name:           proto.String("page_token"),
			Type:           typep(descriptor.FieldDescriptorProto_TYPE_STRING),
			Label:          labelp(descriptor.FieldDescriptorProto_LABEL_OPTIONAL),
			Proto3Optional: proto.Bool(true),
		}),
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
			{
				Pattern: &annotations.HttpRule_Post{
					Post: "/v1/{field_name.nested=projects/*/foo/*}/bars/{another=bar/*/baz/*}/buz",
				},
			},
			{
				Pattern: &annotations.HttpRule_Post{
					Post: "/v1/{field_name.nested=projects/*/foo/*}/bars/{another=bar/*/baz/*}/buz/{biz}/booz",
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
	serv := &descriptor.ServiceDescriptorProto{
		Name: proto.String("Foo"),
	}

	var g generator
	g.opts = &options{
		pkgName: "pkg",
	}
	g.imports = map[pbinfo.ImportSpec]bool{}
	cpb := &conf.ServiceConfig{
		MethodConfig: []*conf.MethodConfig{
			{
				Name: []*conf.MethodConfig_Name{
					{
						Service: "my.pkg.Foo",
					},
				},
				Timeout: &duration.Duration{Seconds: 10},
			},
		},
	}
	data, err := protojson.Marshal(cpb)
	if err != nil {
		t.Error(err)
	}
	in := bytes.NewReader(data)
	g.grpcConf, err = conf.New(in)
	if err != nil {
		t.Error(err)
	}

	commonTypes(&g)
	for _, typ := range []*descriptor.DescriptorProto{
		inputType, outputType, pageInputType, pageInputTypeOptional, pageOutputType, extra,
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
			Name:       proto.String("GetManyThingsOptional"),
			InputType:  proto.String(".my.pkg.PageInputTypeOptional"),
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
		g.descInfo.ParentElement[m] = serv

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
	serv := &descriptor.ServiceDescriptorProto{
		Name: proto.String("Foo"),
	}

	var g generator
	g.imports = map[pbinfo.ImportSpec]bool{}
	g.opts = &options{
		pkgName: "pkg",
	}
	cpb := &conf.ServiceConfig{
		MethodConfig: []*conf.MethodConfig{
			{
				Name: []*conf.MethodConfig_Name{
					{
						Service: "my.pkg.Foo",
					},
				},
				Timeout: &duration.Duration{Seconds: 10},
			},
		},
	}
	data, err := protojson.Marshal(cpb)
	if err != nil {
		t.Error(err)
	}
	in := bytes.NewReader(data)
	g.grpcConf, err = conf.New(in)
	if err != nil {
		t.Error(err)
	}

	commonTypes(&g)
	for _, typ := range []*descriptor.DescriptorProto{
		inputType, outputType,
	} {
		g.descInfo.Type[".my.pkg."+*typ.Name] = typ
		g.descInfo.ParentFile[typ] = file
	}
	g.descInfo.ParentFile[serv] = file
	g.descInfo.ParentElement = map[pbinfo.ProtoType]pbinfo.ProtoType{}

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
		g.descInfo.ParentElement[m] = serv

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
		{name: "numbers", field: "foo_foo64", want: ".GetFooFoo64()"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildAccessor(tt.field); got != tt.want {
				t.Errorf("buildAccessor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsLRO(t *testing.T) {
	lroGetOp := &descriptor.MethodDescriptorProto{
		Name:       proto.String("GetOperation"),
		OutputType: proto.String(".google.longrunning.Operation"),
	}

	actualLRO := &descriptor.MethodDescriptorProto{
		Name:       proto.String("SuperLongRPC"),
		OutputType: proto.String(".google.longrunning.Operation"),
	}

	var g generator
	g.descInfo.ParentFile = map[proto.Message]*descriptor.FileDescriptorProto{
		lroGetOp: {
			Package: proto.String("google.longrunning"),
		},
		actualLRO: {
			Package: proto.String("my.pkg"),
		},
	}

	for _, tst := range []struct {
		in   *descriptor.MethodDescriptorProto
		want bool
	}{
		{lroGetOp, false},
		{actualLRO, true},
	} {
		if got := g.isLRO(tst.in); got != tst.want {
			t.Errorf("isLRO(%v) = %v, want %v", tst.in, got, tst.want)
		}
	}
}

func Test_parseRequestHeaders(t *testing.T) {
	for _, tst := range []struct {
		name, pattern string
		want          [][]string
	}{
		{"not annotated", "", nil},
		{"no params", "/no/params", nil},
		{"not key-value", "/{foo}/{bar}", [][]string{
			{"{foo", "foo"},
			{"{bar", "bar"},
		}},
		{"key-value pair", "/{foo=blah}/{bar=blah}", [][]string{
			{"{foo", "foo"},
			{"{bar", "bar"},
		}},
		{"with underscore", "/{foo_foo}/{bar_bar}", [][]string{
			{"{foo_foo", "foo_foo"},
			{"{bar_bar", "bar_bar"},
		}},
		{"with dot", "/{foo.foo}/{bar.bar}", [][]string{
			{"{foo.foo", "foo.foo"},
			{"{bar.bar", "bar.bar"},
		}},
		{"numbers in field names", "/{foo123}/{bar456}", [][]string{
			{"{foo123", "foo123"},
			{"{bar456", "bar456"},
		}},
		{"everything", "/{foo.foo_foo123=blah}/{bar.bar_bar456=blah}", [][]string{
			{"{foo.foo_foo123", "foo.foo_foo123"},
			{"{bar.bar_bar456", "bar.bar_bar456"},
		}},
	} {
		m := &descriptor.MethodDescriptorProto{}
		if tst.pattern != "" {
			m.Options = &descriptor.MethodOptions{}

			err := setHTTPOption(m.Options, tst.pattern)
			if err != nil {
				t.Errorf("parseRequestHeaders(%s): failed to set http annotation: %v", tst.name, err)
			}
		}

		got, err := parseRequestHeaders(m)
		if err != nil {
			t.Error(err)
			continue
		}

		if diff := cmp.Diff(got, tst.want); diff != "" {
			t.Errorf("parseRequestHeaders(%s) = got(-), want(+):\n%s", tst.name, diff)
		}
	}
}

func Test_lookupFieldType(t *testing.T) {
	typep := func(t descriptor.FieldDescriptorProto_Type) *descriptor.FieldDescriptorProto_Type {
		return &t
	}

	extra := &descriptor.DescriptorProto{
		Name: proto.String("ExtraMessage"),
		Field: []*descriptor.FieldDescriptorProto{
			{
				Name: proto.String("leaf"),
				Type: typep(descriptor.FieldDescriptorProto_TYPE_STRING),
			},
		},
	}

	inputType := &descriptor.DescriptorProto{
		Name: proto.String("InputType"),
		Field: []*descriptor.FieldDescriptorProto{
			{
				Name: proto.String("str"),
				Type: typep(descriptor.FieldDescriptorProto_TYPE_STRING),
			},
			{
				Name: proto.String("bool"),
				Type: typep(descriptor.FieldDescriptorProto_TYPE_BOOL),
			},
			{
				Name: proto.String("int"),
				Type: typep(descriptor.FieldDescriptorProto_TYPE_INT32),
			},
			{
				Name: proto.String("double"),
				Type: typep(descriptor.FieldDescriptorProto_TYPE_DOUBLE),
			},
			{
				Name:     proto.String("msg"),
				Type:     typep(descriptor.FieldDescriptorProto_TYPE_MESSAGE),
				TypeName: proto.String("ExtraMessage"),
			},
		},
	}

	var g generator
	g.descInfo.Type = map[string]pbinfo.ProtoType{
		inputType.GetName(): inputType,
		extra.GetName():     extra,
	}

	for _, tst := range []struct {
		name, msg, field string
		want             descriptor.FieldDescriptorProto_Type
	}{
		{
			name:  "string",
			msg:   "InputType",
			field: "str",
			want:  descriptor.FieldDescriptorProto_TYPE_STRING,
		},
		{
			name:  "boolean",
			msg:   "InputType",
			field: "bool",
			want:  descriptor.FieldDescriptorProto_TYPE_BOOL,
		},
		{
			name:  "integer",
			msg:   "InputType",
			field: "int",
			want:  descriptor.FieldDescriptorProto_TYPE_INT32,
		},
		{
			name:  "double",
			msg:   "InputType",
			field: "double",
			want:  descriptor.FieldDescriptorProto_TYPE_DOUBLE,
		},
		{
			name:  "nested field",
			msg:   "InputType",
			field: "msg.leaf",
			want:  descriptor.FieldDescriptorProto_TYPE_STRING,
		},
	} {
		got := g.lookupFieldType(tst.msg, tst.field)

		if got != tst.want {
			t.Errorf("Test_lookupFieldType(%s): got %v want %v", tst.name, got, tst.want)
		}
	}
}

func TestGRPCStubCall(t *testing.T) {
	getFoo := &descriptor.MethodDescriptorProto{
		Name:       proto.String("GetFoo"),
		InputType:  proto.String("google.protobuf.Empty"),
		OutputType: proto.String("google.protobuf.Empty"),
	}
	getBar := &descriptor.MethodDescriptorProto{
		Name:       proto.String("GetBar"),
		InputType:  proto.String("google.protobuf.Empty"),
		OutputType: proto.String("google.protobuf.Empty"),
	}
	foo := &descriptor.FileDescriptorProto{
		Package: proto.String("google.foo.v1"),
		Options: &descriptor.FileOptions{
			GoPackage: proto.String("google.golang.org/genproto/googleapis/foo/v1"),
		},
		Dependency: []string{"google.protobuf.Empty"},
		Service: []*descriptor.ServiceDescriptorProto{
			{
				Name:   proto.String("FooService"),
				Method: []*descriptor.MethodDescriptorProto{getFoo},
			},
			{
				Name:   proto.String("BarService"),
				Method: []*descriptor.MethodDescriptorProto{getBar},
			},
		},
	}
	var g generator
	err := g.init(&pluginpb.CodeGeneratorRequest{
		ProtoFile: []*descriptor.FileDescriptorProto{foo},
		Parameter: proto.String("go-gapic-package=cloud.google.com/go/foo/apiv1;foo"),
	})
	if err != nil {
		t.Error(err)
	}

	for _, tst := range []struct {
		name, want string
		in         *descriptor.MethodDescriptorProto
	}{
		{
			name: "foo.FooService.GetFoo",
			want: "c.client.GetFoo(ctx, req, settings.GRPC...)",
			in:   getFoo,
		},
		{
			name: "foo.BarService.GetBar",
			want: "c.barClient.GetBar(ctx, req, settings.GRPC...)",
			in:   getBar,
		},
	} {
		got := g.grpcStubCall(tst.in)
		if diff := cmp.Diff(got, tst.want); diff != "" {
			t.Errorf("TestGRPCStubCall(%s) got(-),want(+):\n%s", tst.name, diff)
		}
	}
}

func setHTTPOption(o *descriptor.MethodOptions, pattern string) error {
	return proto.SetExtension(o, annotations.E_Http, &annotations.HttpRule{
		Pattern: &annotations.HttpRule_Get{
			Get: pattern,
		},
	})
}
