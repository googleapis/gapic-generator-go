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
	"reflect"
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
	pageInputTypeOptional := &descriptor.DescriptorProto{
		Name: proto.String("PageInputTypeOptional"),
		Field: []*descriptor.FieldDescriptorProto{
			{
				Name:           proto.String("page_size"),
				Type:           typep(descriptor.FieldDescriptorProto_TYPE_INT32),
				Label:          labelp(descriptor.FieldDescriptorProto_LABEL_OPTIONAL),
				Proto3Optional: proto.Bool(true),
			},
			{
				Name:           proto.String("page_token"),
				Type:           typep(descriptor.FieldDescriptorProto_TYPE_STRING),
				Label:          labelp(descriptor.FieldDescriptorProto_LABEL_OPTIONAL),
				Proto3Optional: proto.Bool(true),
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
		inputType, outputType, pageInputType, pageInputTypeOptional, pageOutputType,
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

func Test_camelToSnake(t *testing.T) {
	for _, tst := range []struct {
		in, want string
	}{
		{"IAMCredentials", "iam_credentials"},
		{"DLP", "dlp"},
		{"OsConfig", "os_config"},
	} {
		if got := camelToSnake(tst.in); got != tst.want {
			t.Errorf("camelToSnake(%q) = %q, want %q", tst.in, got, tst.want)
		}
	}
}

func Test_isLRO(t *testing.T) {
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
		lroGetOp: &descriptor.FileDescriptorProto{
			Package: proto.String("google.longrunning"),
		},
		actualLRO: &descriptor.FileDescriptorProto{
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

func Test_optionsParse(t *testing.T) {
	for _, tst := range []struct {
		param        string
		expectedOpts *options
		expectErr    bool
	}{
		{
			param: "transport=grpc,go-gapic-package=path;pkg",
			expectedOpts: &options{
				transports: []Transport{grpc},
				pkgPath:    "path",
				pkgName:    "pkg",
				outDir:     "path",
			},
			expectErr: false,
		},
		{
			param: "transport=rest+grpc,go-gapic-package=path;pkg",
			expectedOpts: &options{
				transports: []Transport{rest, grpc},
				pkgPath:    "path",
				pkgName:    "pkg",
				outDir:     "path",
			},
			expectErr: false,
		},
		{
			param: "go-gapic-package=path;pkg",
			expectedOpts: &options{
				transports: []Transport{grpc},
				pkgPath:    "path",
				pkgName:    "pkg",
				outDir:     "path",
			},
		},
		{
			param: "module=path,go-gapic-package=path/to/out;pkg",
			expectedOpts: &options{
				transports:   []Transport{grpc},
				pkgPath:      "path/to/out",
				pkgName:      "pkg",
				outDir:       "to/out",
				modulePrefix: "path",
			},
			expectErr: false,
		},
		{
			param:     "transport=tcp,go-gapic-package=path;pkg",
			expectErr: true,
		},
		{
			param:     "go-gapic-package=pkg;",
			expectErr: true,
		},
		{
			param:     "go-gapic-package=;path",
			expectErr: true,
		},
		{
			param:     "go-gapic-package=bogus",
			expectErr: true,
		},
		{
			param:     "module=different_path,go-gapic-package=path;pkg",
			expectErr: true,
		},
	} {
		opts, err := ParseOptions(&tst.param)
		if tst.expectErr && err == nil {
			t.Errorf("ParseOptions(%s) expected error", tst.param)
			continue
		}

		if !tst.expectErr && err != nil {
			t.Errorf("ParseOptions(%s) got unexpected error: %v", tst.param, err)
			continue
		}

		if !reflect.DeepEqual(opts, tst.expectedOpts) {
			t.Errorf("ParseOptions(%s) = %v, expected %v", tst.param, opts, tst.expectedOpts)
			continue
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

		if got, err := parseRequestHeaders(m); err != nil {
			t.Error(err)
			continue
		}
		if diff := cmp.Diff(got, tst.want); diff != "" {
			t.Errorf("parseRequestHeaders(%s) = %v, want %v, diff %s", tst.name, got, tst.want, diff)
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
