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
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	longrunning "cloud.google.com/go/longrunning/autogen/longrunningpb"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/google/go-cmp/cmp"
	conf "github.com/googleapis/gapic-generator-go/internal/grpc_service_config"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"github.com/googleapis/gapic-generator-go/internal/snippets"
	"github.com/googleapis/gapic-generator-go/internal/txtdiff"
	"google.golang.org/genproto/googleapis/api/annotations"
	metadatapb "google.golang.org/genproto/googleapis/gapic/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoiface"
	duration "google.golang.org/protobuf/types/known/durationpb"
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

func TestCodeSnippet(t *testing.T) {
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
			want: "//  abc\n//  def\n//\n",
		},
	} {
		g.pt.Reset()
		g.codesnippet(tst.in)
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

// TODO(chrisdsmith): Expand this test to invoke Gen (or at least gen) in gengapic.go,
// otherwise move it to gengrpc_test.go. (genGRPCMethods was extracted to gengrpc.go on 2021-05-04.)
func TestGenGRPCMethods(t *testing.T) {
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

	topLevelEnum := &descriptor.EnumDescriptorProto{
		Name: proto.String("TopLevelEnum"),
	}
	nestedEnum := &descriptor.EnumDescriptorProto{
		Name: proto.String("NestedEnum"),
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
			{
				Name:     proto.String("top_level_enum"),
				Type:     typep(descriptor.FieldDescriptorProto_TYPE_ENUM),
				TypeName: proto.String(".my.pkg.TopLevelEnum"),
			},
			{
				Name:     proto.String("nested_enum"),
				Type:     typep(descriptor.FieldDescriptorProto_TYPE_ENUM),
				TypeName: proto.String(".my.pkg.InputType.NestedEnum"),
			},
		},
		EnumType: []*descriptor.EnumDescriptorProto{
			nestedEnum,
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
			{
				Pattern: &annotations.HttpRule_Post{
					Post: "/v1/{top_level_enum}/enums/{nested_enum}",
				},
			},
		},
	}
	proto.SetExtension(opts, annotations.E_Http, ext)

	optsGetAnotherThing := &descriptor.MethodOptions{}
	extGetAnotherThing := &annotations.RoutingRule{
		RoutingParameters: []*annotations.RoutingParameter{
			{
				Field: "other",
			},
			{
				Field:        "other",
				PathTemplate: "{name=projects/*}/foos",
			},
			{
				Field:        "another",
				PathTemplate: "{foo_name=projects/*}/bars/*/**",
			},
			{
				Field:        "another",
				PathTemplate: "{foo_name=projects/*/foos/*}/bars/*/**",
			},
			{
				Field:        "another",
				PathTemplate: "{foo_name=**}",
			},
			{
				Field:        "field_name.nested",
				PathTemplate: "{nested_name=**}",
			},
			{
				Field:        "field_name.nested",
				PathTemplate: "{part_of_nested=projects/*}/bars",
			},
		},
	}
	proto.SetExtension(optsGetAnotherThing, annotations.E_Routing, extGetAnotherThing)

	optsGetManyOtherThings := &descriptor.MethodOptions{}
	extGetManyOtherThings := &annotations.RoutingRule{}
	proto.SetExtension(optsGetManyOtherThings, annotations.E_Routing, extGetManyOtherThings)

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
	g.metadata = &metadatapb.GapicMetadata{
		Services: make(map[string]*metadatapb.GapicMetadata_ServiceForTransport),
	}
	g.mixins = mixins{
		"google.longrunning.Operations":   operationsMethods(),
		"google.cloud.location.Locations": locationMethods(),
		"google.iam.v1.IAMPolicy":         iamPolicyMethods(),
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
	for _, typ := range []pbinfo.ProtoType{
		inputType, outputType, pageInputType, pageInputTypeOptional, pageOutputType, extra, topLevelEnum,
	} {
		g.descInfo.Type[".my.pkg."+typ.GetName()] = typ
		g.descInfo.ParentFile[typ] = file
	}
	g.descInfo.ParentFile[serv] = file
	g.descInfo.ParentElement = map[pbinfo.ProtoType]pbinfo.ProtoType{
		paginatedField: pageOutputType,
		nestedEnum:     inputType,
	}
	n := fmt.Sprintf(".my.pkg.%s.%s", inputType.GetName(), nestedEnum.GetName())
	g.descInfo.Type[n] = nestedEnum

methods:
	for _, tst := range []struct {
		m       *descriptor.MethodDescriptorProto
		imports map[pbinfo.ImportSpec]bool
	}{
		{
			m: &descriptor.MethodDescriptorProto{
				Name:       proto.String("GetEmptyThing"),
				InputType:  proto.String(".my.pkg.InputType"),
				OutputType: proto.String(emptyType),
				Options:    opts,
			},
			imports: map[pbinfo.ImportSpec]bool{
				{Path: "fmt"}:                            true,
				{Path: "net/url"}:                        true,
				{Path: "time"}:                           true,
				{Name: "mypackagepb", Path: "mypackage"}: true,
			},
		},
		{
			m: &descriptor.MethodDescriptorProto{
				Name:       proto.String("GetOneThing"),
				InputType:  proto.String(".my.pkg.InputType"),
				OutputType: proto.String(".my.pkg.OutputType"),
				Options:    opts,
			},
			imports: map[pbinfo.ImportSpec]bool{
				{Path: "fmt"}:                            true,
				{Path: "net/url"}:                        true,
				{Path: "time"}:                           true,
				{Name: "mypackagepb", Path: "mypackage"}: true,
			},
		},
		{
			m: &descriptor.MethodDescriptorProto{
				Name:       proto.String("GetManyThings"),
				InputType:  proto.String(".my.pkg.PageInputType"),
				OutputType: proto.String(".my.pkg.PageOutputType"),
				Options:    opts,
			},
			imports: map[pbinfo.ImportSpec]bool{
				{Path: "fmt"}:                              true,
				{Path: "google.golang.org/api/iterator"}:   true,
				{Path: "google.golang.org/protobuf/proto"}: true,
				{Path: "net/url"}:                          true,
				{Name: "mypackagepb", Path: "mypackage"}:   true,
			},
		},
		{
			m: &descriptor.MethodDescriptorProto{
				Name:       proto.String("GetManyThingsOptional"),
				InputType:  proto.String(".my.pkg.PageInputTypeOptional"),
				OutputType: proto.String(".my.pkg.PageOutputType"),
				Options:    opts,
			},
			imports: map[pbinfo.ImportSpec]bool{
				{Path: "fmt"}:                              true,
				{Path: "google.golang.org/api/iterator"}:   true,
				{Path: "google.golang.org/protobuf/proto"}: true,
				{Path: "net/url"}:                          true,
				{Name: "mypackagepb", Path: "mypackage"}:   true,
			},
		},
		{
			m: &descriptor.MethodDescriptorProto{
				Name:            proto.String("ServerThings"),
				InputType:       proto.String(".my.pkg.InputType"),
				OutputType:      proto.String(".my.pkg.OutputType"),
				ServerStreaming: proto.Bool(true),
				Options:         opts,
			},
			imports: map[pbinfo.ImportSpec]bool{
				{Path: "fmt"}:                            true,
				{Path: "net/url"}:                        true,
				{Name: "mypackagepb", Path: "mypackage"}: true,
			},
		},
		{
			m: &descriptor.MethodDescriptorProto{
				Name:            proto.String("ClientThings"),
				InputType:       proto.String(".my.pkg.InputType"),
				OutputType:      proto.String(".my.pkg.OutputType"),
				ClientStreaming: proto.Bool(true),
				Options:         opts,
			},
			imports: map[pbinfo.ImportSpec]bool{
				{Name: "mypackagepb", Path: "mypackage"}: true,
			},
		},
		{
			m: &descriptor.MethodDescriptorProto{
				Name:            proto.String("BidiThings"),
				InputType:       proto.String(".my.pkg.InputType"),
				OutputType:      proto.String(".my.pkg.OutputType"),
				ServerStreaming: proto.Bool(true),
				ClientStreaming: proto.Bool(true),
				Options:         opts,
			},
			imports: map[pbinfo.ImportSpec]bool{
				{Name: "mypackagepb", Path: "mypackage"}: true,
			},
		},
		// Test for dynamic routing header annotations.
		{
			m: &descriptor.MethodDescriptorProto{
				Name:       proto.String("GetAnotherThing"),
				InputType:  proto.String(".my.pkg.InputType"),
				OutputType: proto.String(".my.pkg.OutputType"),
				Options:    optsGetAnotherThing,
			},
			imports: map[pbinfo.ImportSpec]bool{
				{Path: "fmt"}:                            true,
				{Path: "net/url"}:                        true,
				{Path: "regexp"}:                         true,
				{Path: "strings"}:                        true,
				{Path: "time"}:                           true,
				{Name: "mypackagepb", Path: "mypackage"}: true,
			},
		},
		// Test for empty dynamic routing annotation, so no headers should be sent.
		{
			m: &descriptor.MethodDescriptorProto{
				Name:       proto.String("GetManyOtherThings"),
				InputType:  proto.String(".my.pkg.PageInputType"),
				OutputType: proto.String(".my.pkg.PageOutputType"),
				Options:    optsGetManyOtherThings,
			},
			imports: map[pbinfo.ImportSpec]bool{
				{Path: "google.golang.org/api/iterator"}:   true,
				{Path: "google.golang.org/protobuf/proto"}: true,
				{Name: "mypackagepb", Path: "mypackage"}:   true,
			},
		},
	} {
		g.reset()
		g.descInfo.ParentElement[tst.m] = serv
		serv.Method = []*descriptor.MethodDescriptorProto{
			tst.m,
		}

		g.aux = &auxTypes{
			iters: map[string]*iterType{},
		}
		if err := g.genGRPCMethods(serv, "Foo"); err != nil {
			t.Error(err)
			continue
		}

		var lros []*descriptor.MethodDescriptorProto
		for m := range g.aux.lros {
			lros = append(lros, m)
		}
		sort.Slice(lros, func(i, j int) bool {
			return lros[i].GetName() < lros[j].GetName()
		})
		for _, m := range lros {
			if err := g.lroType("MyService", serv, m); err != nil {
				t.Error(err)
				continue methods
			}
		}

		for _, iter := range g.aux.iters {
			g.pagingIter(iter)
		}

		if diff := cmp.Diff(g.imports, tst.imports); diff != "" {
			t.Errorf("TestGenMethod(%s): imports got(-),want(+):\n%s", tst.m.GetName(), diff)
		}

		txtdiff.Diff(t, tst.m.GetName(), g.pt.String(), filepath.Join("testdata", "method_"+tst.m.GetName()+".want"))
	}
}

func TestContainsDeprecated(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{
			in:   "",
			want: false,
		},
		{
			in:   "This comment does not contain the prefix Deprecated.\nNot even on this line.",
			want: false,
		},
		{
			in:   "Deprecated: this is a proper deprecation notice.",
			want: true,
		},
		{
			in:   "This is a comment that includes a deprecation notice.\nDeprecated: This is a proper deprecation notice.",
			want: true,
		},
		{
			in:   "This is not a properly formatted Deprecated: notice.\nNeither is this one - Deprecated:",
			want: false,
		},
	}

	for _, tst := range tests {
		if diff := cmp.Diff(containsDeprecated(tst.in), tst.want); diff != "" {
			t.Errorf("comment() got(-),want(+):\n%s", diff)
		}
	}
}

func TestMethodDoc(t *testing.T) {
	serv := &descriptor.ServiceDescriptorProto{
		Name: proto.String("Foo"),
	}
	servName := "Foo"
	methodName := "MyMethod"
	m := &descriptor.MethodDescriptorProto{
		Name: proto.String(methodName),
	}

	g := generator{
		comments: make(map[protoiface.MessageV1]string),
	}

	for _, tst := range []struct {
		in, want                    string
		clientStreaming, deprecated bool
		opts                        options
	}{
		{
			in:   "",
			want: "",
		},
		{
			in:   "Does stuff.\nIt also does other stuffs.",
			want: "// MyMethod does stuff.\n// It also does other stuffs.\n",
		},
		{
			in:         "This is deprecated.\nIt does not have a proper comment.",
			want:       "// MyMethod this is deprecated.\n// It does not have a proper comment.\n//\n// Deprecated: MyMethod may be removed in a future version.\n",
			deprecated: true,
		},
		{
			in:         "Deprecated: this is a proper deprecation notice.",
			want:       "// MyMethod is deprecated.\n//\n// Deprecated: this is a proper deprecation notice.\n",
			deprecated: true,
		},
		{
			in:         "Does my thing.\nDeprecated: this is a proper deprecation notice.",
			want:       "// MyMethod does my thing.\n// Deprecated: this is a proper deprecation notice.\n",
			deprecated: true,
		},
		{
			in:         "",
			want:       "// MyMethod is deprecated.\n//\n// Deprecated: MyMethod may be removed in a future version.\n",
			deprecated: true,
		},
		{
			in:              "Does client streaming stuff.\nIt also does other stuffs.",
			want:            "// MyMethod does client streaming stuff.\n// It also does other stuffs.\n//\n// This method is not supported for the REST transport.\n",
			clientStreaming: true,
			opts:            options{transports: []transport{rest}},
		},
	} {
		g.opts = &tst.opts
		sm, err := snippets.NewMetadata("mypackage", "github.com/googleapis/mypackage", "mypackage.googleapis.com")
		if err != nil {
			t.Fatal(err)
		}
		sm.AddService(servName)
		sm.AddMethod(servName, methodName, 50)
		g.snippetMetadata = sm
		g.comments[m] = tst.in
		m.Options = &descriptor.MethodOptions{
			Deprecated: proto.Bool(tst.deprecated),
		}
		m.ClientStreaming = proto.Bool(tst.clientStreaming)
		g.pt.Reset()
		g.methodDoc(serv, servName, m)
		if diff := cmp.Diff(g.pt.String(), tst.want); diff != "" {
			t.Errorf("comment() got(-),want(+):\n%s", diff)
		}
		mi := g.snippetMetadata.ToMetadataIndex()
		if got := len(mi.Snippets); got != 1 {
			t.Errorf("%s: wanted len 1 Snippets, got %d", t.Name(), got)
		}
		snp := mi.Snippets[0]
		// remove slashes to compare with snippet description.
		want := strings.Replace(tst.want, "// ", "", -1)
		want = strings.Replace(want, "//", "", -1)
		want = strings.Trim(want, "\n")
		if got := snp.Description; !tst.clientStreaming && got != want {
			t.Errorf("%s: wanted %s, got %s", t.Name(), want, got)
		}
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
	g.mixins = mixins{
		"google.longrunning.Operations":   operationsMethods(),
		"google.cloud.location.Locations": locationMethods(),
		"google.iam.v1.IAMPolicy":         iamPolicyMethods(),
	}
	g.opts = &options{
		pkgName:    "pkg",
		transports: []transport{grpc},
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

		g.aux = &auxTypes{
			lros: map[*descriptor.MethodDescriptorProto]bool{},
		}

		if err := g.genGRPCMethod("Foo", serv, m); err != nil {
			t.Error(err)
			continue
		}

		var genLros []*descriptor.MethodDescriptorProto
		for m := range g.aux.lros {
			genLros = append(genLros, m)
		}
		sort.Slice(genLros, func(i, j int) bool {
			return genLros[i].GetName() < genLros[j].GetName()
		})
		for _, m := range genLros {
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
		name    string
		field   string
		want    string
		variant func(string) string
	}{
		{name: "simple", field: "foo_foo", want: ".GetFooFoo()", variant: fieldGetter},
		{name: "nested", field: "foo_foo.bar_bar", want: ".GetFooFoo().GetBarBar()", variant: fieldGetter},
		{name: "numbers", field: "foo_foo64", want: ".GetFooFoo64()", variant: fieldGetter},
		{name: "raw_final", field: "foo.bar.baz.bif", want: ".GetFoo().GetBar().GetBaz().Bif", variant: directAccess},
		{name: "independent_number", field: "display_video_360_advertiser_links", want: ".GetDisplayVideo_360AdvertiserLinks()", variant: fieldGetter},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.variant(tt.field); got != tt.want {
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
	g.descInfo.ParentFile = map[protoiface.MessageV1]*descriptor.FileDescriptorProto{
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

func Test_parseImplicitRequestHeaders(t *testing.T) {
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

			setHTTPOption(m.Options, tst.pattern)
		}

		got := parseImplicitRequestHeaders(m)

		if diff := cmp.Diff(got, tst.want); diff != "" {
			t.Errorf("parseImplicitRequestHeaders(%s) = got(-), want(+):\n%s", tst.name, diff)
		}
	}
}

func Test_ContainsDynamicRequestHeaders(t *testing.T) {
	methodNoRule := &descriptor.MethodDescriptorProto{
		Name:       proto.String("MethodNoRule"),
		InputType:  proto.String(".my.pkg.InputType"),
		OutputType: proto.String(".google.longrunning.Operation"),
		Options:    &descriptor.MethodOptions{},
	}
	proto.SetExtension(methodNoRule.GetOptions(), annotations.E_Http, &annotations.HttpRule{
		Pattern: &annotations.HttpRule_Post{
			Post: "/v1/operations",
		},
	})
	methodOneRule := &descriptor.MethodDescriptorProto{
		Name:       proto.String("MethodOneRule"),
		InputType:  proto.String(".my.pkg.InputType"),
		OutputType: proto.String(".google.longrunning.Operation"),
		Options:    &descriptor.MethodOptions{},
	}
	extRoutingOneRule := &annotations.RoutingRule{
		RoutingParameters: []*annotations.RoutingParameter{
			{
				Field:        "other",
				PathTemplate: "{routing_id=**}",
			},
		},
	}
	proto.SetExtension(methodOneRule.GetOptions(), annotations.E_Routing, extRoutingOneRule)
	proto.SetExtension(methodOneRule.GetOptions(), annotations.E_Http, &annotations.HttpRule{
		Pattern: &annotations.HttpRule_Post{
			Post: "/v1/operations",
		},
	})
	methodEmptyRule := &descriptor.MethodDescriptorProto{
		Options: &descriptor.MethodOptions{},
	}
	extRoutingEmptyRule := &annotations.RoutingRule{
		RoutingParameters: []*annotations.RoutingParameter{
			{},
		},
	}
	proto.SetExtension(methodEmptyRule.GetOptions(), annotations.E_Routing, extRoutingEmptyRule)
	methodMultipleRules := &descriptor.MethodDescriptorProto{
		Options: &descriptor.MethodOptions{},
	}
	extRoutingMultipleRules := &annotations.RoutingRule{
		RoutingParameters: []*annotations.RoutingParameter{
			{
				Field: "other",
			},
			{
				Field:        "other",
				PathTemplate: "{name=projects/*}/foos",
			},
			{
				Field:        "another",
				PathTemplate: "{foo_name=projects/*}/bars/*/**",
			},
		},
	}
	proto.SetExtension(methodMultipleRules.GetOptions(), annotations.E_Routing, extRoutingMultipleRules)
	for _, tst := range []struct {
		want   bool
		method *descriptor.MethodDescriptorProto
	}{
		{
			method: methodOneRule,
			want:   true,
		},
		{
			method: methodEmptyRule,
			want:   true,
		},
		{
			method: methodMultipleRules,
			want:   true,
		},
		{
			method: methodNoRule,
			want:   false,
		},
	} {
		got := dynamicRequestHeadersExist(tst.method)
		if diff := cmp.Diff(got, tst.want); diff != "" {
			t.Errorf("%s: got(-),want(+):\n%s", tst.method, diff)
		}
	}
}

func Test_lookupField(t *testing.T) {
	topLevelEnum := &descriptor.EnumDescriptorProto{
		Name: proto.String("TopLevelEnum"),
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
			{
				Name:     proto.String("top_level_enum"),
				Type:     typep(descriptor.FieldDescriptorProto_TYPE_ENUM),
				TypeName: proto.String("TopLevelEnum"),
			},
			{
				Name:     proto.String("nested_enum"),
				Type:     typep(descriptor.FieldDescriptorProto_TYPE_ENUM),
				TypeName: proto.String("NestedEnum"),
			},
		},
	}

	var g generator
	g.descInfo.Type = map[string]pbinfo.ProtoType{
		inputType.GetName():    inputType,
		extra.GetName():        extra,
		topLevelEnum.GetName(): topLevelEnum,
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
		{
			name:  "top level enum",
			msg:   "InputType",
			field: "top_level_enum",
			want:  descriptor.FieldDescriptorProto_TYPE_ENUM,
		},
	} {
		got := g.lookupField(tst.msg, tst.field)

		if got.GetType() != tst.want {
			t.Errorf("Test_lookupField(%s): got %v want %v", tst.name, got.GetType(), tst.want)
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

func TestReturnType(t *testing.T) {
	op := &descriptor.DescriptorProto{
		Name: proto.String("Operation"),
	}
	foo := &descriptor.DescriptorProto{
		Name: proto.String("Foo"),
	}
	com := &descriptor.MethodDescriptorProto{
		OutputType: proto.String(".google.cloud.foo.v1.Operation"),
		Options:    &descriptor.MethodOptions{},
	}
	proto.SetExtension(com.GetOptions(), annotations.E_Http, &annotations.HttpRule{
		Pattern: &annotations.HttpRule_Post{
			Post: "/v1/operations",
		},
	})
	reg := &descriptor.MethodDescriptorProto{
		OutputType: proto.String(".google.cloud.foo.v1.Foo"),
		Options:    &descriptor.MethodOptions{},
	}
	proto.SetExtension(reg.GetOptions(), annotations.E_Http, &annotations.HttpRule{
		Pattern: &annotations.HttpRule_Post{
			Post: "/v1/operations",
		},
	})
	get := &descriptor.MethodDescriptorProto{
		OutputType: proto.String(".google.cloud.foo.v1.Operation"),
		Options:    &descriptor.MethodOptions{},
	}
	proto.SetExtension(get.GetOptions(), annotations.E_Http, &annotations.HttpRule{
		Pattern: &annotations.HttpRule_Get{
			Get: "/v1/operations",
		},
	})
	wait := &descriptor.MethodDescriptorProto{
		Name:       proto.String("Wait"),
		OutputType: proto.String(".google.cloud.foo.v1.Operation"),
		Options:    &descriptor.MethodOptions{},
	}
	proto.SetExtension(wait.GetOptions(), annotations.E_Http, &annotations.HttpRule{
		Pattern: &annotations.HttpRule_Post{
			Post: "/v1/operations",
		},
	})
	f := &descriptor.FileDescriptorProto{
		Package: proto.String("google.cloud.foo.v1"),
		Options: &descriptor.FileOptions{
			GoPackage: proto.String("google.golang.org/genproto/cloud/foo/v1;foo"),
		},
	}
	g := &generator{
		aux: &auxTypes{
			customOp: &customOp{
				message: op,
			},
		},
		opts: &options{
			diregapic: true,
		},
		descInfo: pbinfo.Info{
			ParentFile: map[protoiface.MessageV1]*descriptor.FileDescriptorProto{
				op:  f,
				foo: f,
			},
			Type: map[string]pbinfo.ProtoType{
				com.GetOutputType(): op,
				reg.GetOutputType(): foo,
			},
		},
	}

	for _, tst := range []struct {
		name, want string
		method     *descriptor.MethodDescriptorProto
	}{
		{
			name:   "custom_operation",
			method: com,
			want:   "*Operation",
		},
		{
			name:   "regular_unary",
			method: reg,
			want:   "*foopb.Foo",
		},
		{
			name:   "get_custom_op",
			method: get,
			want:   "*foopb.Operation",
		},
		{
			name:   "wait_custom_op",
			method: wait,
			want:   "*foopb.Operation",
		},
	} {
		got, err := g.returnType(tst.method)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(got, tst.want); diff != "" {
			t.Errorf("%s: got(-),want(+):\n%s", tst.name, diff)
		}
	}
}

func setHTTPOption(o *descriptor.MethodOptions, pattern string) {
	proto.SetExtension(o, annotations.E_Http, &annotations.HttpRule{
		Pattern: &annotations.HttpRule_Get{
			Get: pattern,
		},
	})
}
