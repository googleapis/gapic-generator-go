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
	"strings"
	"testing"

	longrunning "cloud.google.com/go/longrunning/autogen/longrunningpb"
	"github.com/google/go-cmp/cmp"
	conf "github.com/googleapis/gapic-generator-go/internal/grpc_service_config"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"github.com/googleapis/gapic-generator-go/internal/snippets"
	"github.com/googleapis/gapic-generator-go/internal/txtdiff"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/genproto/googleapis/api/serviceconfig"
	metadatapb "google.golang.org/genproto/googleapis/gapic/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoiface"
	"google.golang.org/protobuf/types/descriptorpb"
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
	extra := &descriptorpb.DescriptorProto{
		Name: proto.String("ExtraMessage"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:  proto.String("nested"),
				Type:  typep(descriptorpb.FieldDescriptorProto_TYPE_STRING),
				Label: labelp(descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL),
			},
		},
	}

	topLevelEnum := &descriptorpb.EnumDescriptorProto{
		Name: proto.String("TopLevelEnum"),
	}
	nestedEnum := &descriptorpb.EnumDescriptorProto{
		Name: proto.String("NestedEnum"),
	}

	optsUUID4 := &descriptorpb.FieldOptions{}
	proto.SetExtension(optsUUID4, annotations.E_FieldInfo, &annotations.FieldInfo{Format: annotations.FieldInfo_UUID4})

	inputType := &descriptorpb.DescriptorProto{
		Name: proto.String("InputType"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:  proto.String("other"),
				Type:  typep(descriptorpb.FieldDescriptorProto_TYPE_STRING),
				Label: labelp(descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL),
			},
			{
				Name:  proto.String("another"),
				Type:  typep(descriptorpb.FieldDescriptorProto_TYPE_STRING),
				Label: labelp(descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL),
			},
			{
				Name:  proto.String("biz"),
				Type:  typep(descriptorpb.FieldDescriptorProto_TYPE_DOUBLE),
				Label: labelp(descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL),
			},
			{
				Name:     proto.String("field_name"),
				Type:     typep(descriptorpb.FieldDescriptorProto_TYPE_MESSAGE),
				Label:    labelp(descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL),
				TypeName: proto.String(".my.pkg.ExtraMessage"),
			},
			{
				Name:     proto.String("top_level_enum"),
				Type:     typep(descriptorpb.FieldDescriptorProto_TYPE_ENUM),
				TypeName: proto.String(".my.pkg.TopLevelEnum"),
			},
			{
				Name:     proto.String("nested_enum"),
				Type:     typep(descriptorpb.FieldDescriptorProto_TYPE_ENUM),
				TypeName: proto.String(".my.pkg.InputType.NestedEnum"),
			},
			{
				Name:           proto.String("request_id"),
				Type:           typep(descriptorpb.FieldDescriptorProto_TYPE_STRING),
				Label:          labelp(descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL),
				Proto3Optional: proto.Bool(true),
				Options:        optsUUID4,
			},
			{
				Name:    proto.String("non_proto3optional_request_id"),
				Type:    typep(descriptorpb.FieldDescriptorProto_TYPE_STRING),
				Label:   labelp(descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL),
				Options: optsUUID4,
			},
		},
		EnumType: []*descriptorpb.EnumDescriptorProto{
			nestedEnum,
		},
	}
	outputType := &descriptorpb.DescriptorProto{
		Name: proto.String("OutputType"),
	}

	pageInputType := &descriptorpb.DescriptorProto{
		Name: proto.String("PageInputType"),
		Field: append(inputType.GetField(), &descriptorpb.FieldDescriptorProto{
			Name:  proto.String("page_size"),
			Type:  typep(descriptorpb.FieldDescriptorProto_TYPE_INT32),
			Label: labelp(descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL),
		}, &descriptorpb.FieldDescriptorProto{
			Name:  proto.String("page_token"),
			Type:  typep(descriptorpb.FieldDescriptorProto_TYPE_STRING),
			Label: labelp(descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL),
		}),
	}
	pageInputTypeOptional := &descriptorpb.DescriptorProto{
		Name: proto.String("PageInputTypeOptional"),
		Field: append(inputType.GetField(), &descriptorpb.FieldDescriptorProto{
			Name:           proto.String("page_size"),
			Type:           typep(descriptorpb.FieldDescriptorProto_TYPE_INT32),
			Label:          labelp(descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL),
			Proto3Optional: proto.Bool(true),
		}, &descriptorpb.FieldDescriptorProto{
			Name:           proto.String("page_token"),
			Type:           typep(descriptorpb.FieldDescriptorProto_TYPE_STRING),
			Label:          labelp(descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL),
			Proto3Optional: proto.Bool(true),
		}),
	}
	paginatedField := &descriptorpb.FieldDescriptorProto{
		Name:  proto.String("items"),
		Type:  typep(descriptorpb.FieldDescriptorProto_TYPE_STRING),
		Label: labelp(descriptorpb.FieldDescriptorProto_LABEL_REPEATED),
	}
	pageOutputType := &descriptorpb.DescriptorProto{
		Name: proto.String("PageOutputType"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:  proto.String("next_page_token"),
				Type:  typep(descriptorpb.FieldDescriptorProto_TYPE_STRING),
				Label: labelp(descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL),
			},
			paginatedField,
		},
	}

	opts := &descriptorpb.MethodOptions{}
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

	optsGetAnotherThing := &descriptorpb.MethodOptions{}
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

	optsGetManyOtherThings := &descriptorpb.MethodOptions{}
	extGetManyOtherThings := &annotations.RoutingRule{}
	proto.SetExtension(optsGetManyOtherThings, annotations.E_Routing, extGetManyOtherThings)

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

	g.serviceConfig = &serviceconfig.Service{
		Publishing: &annotations.Publishing{
			MethodSettings: []*annotations.MethodSettings{
				{
					Selector: "my.pkg.Foo.GetEmptyThing",
					AutoPopulatedFields: []string{
						"request_id",
					},
				},
				{
					Selector: "my.pkg.Foo.GetOneThing",
					AutoPopulatedFields: []string{
						"request_id",
						"non_proto3optional_request_id",
					},
				},
			},
		},
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

	for _, tst := range []struct {
		m       *descriptorpb.MethodDescriptorProto
		imports map[pbinfo.ImportSpec]bool
	}{
		{
			m: &descriptorpb.MethodDescriptorProto{
				Name:       proto.String("GetEmptyThing"),
				InputType:  proto.String(".my.pkg.InputType"),
				OutputType: proto.String(emptyType),
				Options:    opts,
			},
			imports: map[pbinfo.ImportSpec]bool{
				{Path: "fmt"}:                            true,
				{Path: "github.com/google/uuid"}:         true,
				{Path: "net/url"}:                        true,
				{Name: "mypackagepb", Path: "mypackage"}: true,
			},
		},
		{
			m: &descriptorpb.MethodDescriptorProto{
				Name:       proto.String("GetOneThing"),
				InputType:  proto.String(".my.pkg.InputType"),
				OutputType: proto.String(".my.pkg.OutputType"),
				Options:    opts,
			},
			imports: map[pbinfo.ImportSpec]bool{
				{Path: "fmt"}:                            true,
				{Path: "github.com/google/uuid"}:         true,
				{Path: "net/url"}:                        true,
				{Name: "mypackagepb", Path: "mypackage"}: true,
			},
		},
		{
			m: &descriptorpb.MethodDescriptorProto{
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
			m: &descriptorpb.MethodDescriptorProto{
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
			m: &descriptorpb.MethodDescriptorProto{
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
			m: &descriptorpb.MethodDescriptorProto{
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
			m: &descriptorpb.MethodDescriptorProto{
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
			m: &descriptorpb.MethodDescriptorProto{
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
				{Name: "mypackagepb", Path: "mypackage"}: true,
			},
		},
		// Test for empty dynamic routing annotation, so no headers should be sent.
		{
			m: &descriptorpb.MethodDescriptorProto{
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
		t.Run(tst.m.GetName(), func(t *testing.T) {
			g.reset()
			g.descInfo.ParentElement[tst.m] = serv
			serv.Method = []*descriptorpb.MethodDescriptorProto{
				tst.m,
			}

			g.aux = &auxTypes{
				iters: map[string]*iterType{},
			}
			if err := g.genGRPCMethods(serv, "Foo"); err != nil {
				t.Fatal(err)
			}
			if err := g.genOperationBuilders(serv, "MyService"); err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(g.imports, tst.imports); diff != "" {
				t.Errorf("TestGenMethod(%s): imports got(-),want(+):\n%s", tst.m.GetName(), diff)
			}
			txtdiff.Diff(t, g.pt.String(), filepath.Join("testdata", "method_"+tst.m.GetName()+".want"))
		})
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
	servName := "Foo"
	methodName := "MyMethod"
	m := &descriptorpb.MethodDescriptorProto{
		Name: proto.String(methodName),
	}
	serv := &descriptorpb.ServiceDescriptorProto{
		Name: proto.String(servName),
	}

	g := generator{
		comments: make(map[protoiface.MessageV1]string),
	}
	g.descInfo.ParentElement = map[pbinfo.ProtoType]pbinfo.ProtoType{}
	g.descInfo.ParentElement[m] = serv

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
		sm := snippets.NewMetadata("mypackage", "github.com/googleapis/mypackage", "mypackagego")
		sm.AddService(servName, "mypackage.googleapis.com")
		sm.AddMethod(servName, methodName, "mypackage", servName, 50)
		g.snippetMetadata = sm
		g.comments[m] = tst.in
		m.Options = &descriptorpb.MethodOptions{
			Deprecated: proto.Bool(tst.deprecated),
		}
		m.ClientStreaming = proto.Bool(tst.clientStreaming)
		g.pt.Reset()
		g.methodDoc(m, serv)
		if diff := cmp.Diff(g.pt.String(), tst.want); diff != "" {
			t.Errorf("comment() got(-),want(+):\n%s", diff)
		}
		mi := g.snippetMetadata.ToMetadataIndex()
		if got := len(mi.Snippets); got != 1 {
			t.Errorf("%s: got %d want 1,", t.Name(), got)
		}
		snp := mi.Snippets[0]
		// remove slashes to compare with snippet description.
		want := strings.Replace(tst.want, "// ", "", -1)
		want = strings.Replace(want, "//", "", -1)
		want = strings.Trim(want, "\n")
		if got := snp.Description; !tst.clientStreaming && got != want {
			t.Errorf("%s: got %s want %s", t.Name(), got, want)
		}
	}
}

func TestGenOperationBuilders(t *testing.T) {
	inputType := &descriptorpb.DescriptorProto{
		Name: proto.String("InputType"),
	}
	outputType := &descriptorpb.DescriptorProto{
		Name: proto.String("OutputType"),
	}
	metadataType := &descriptorpb.DescriptorProto{
		Name: proto.String("MetadataType"),
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
	for _, typ := range []*descriptorpb.DescriptorProto{
		inputType, outputType, metadataType,
	} {
		g.descInfo.Type[".my.pkg."+*typ.Name] = typ
		g.descInfo.ParentFile[typ] = file
	}
	g.descInfo.ParentFile[serv] = file
	g.descInfo.ParentElement = map[pbinfo.ProtoType]pbinfo.ProtoType{}

	emptyLRO := &longrunning.OperationInfo{
		ResponseType: emptyValue,
		MetadataType: "MetadataType",
	}
	emptyLROOpts := &descriptorpb.MethodOptions{}
	proto.SetExtension(emptyLROOpts, longrunning.E_OperationInfo, emptyLRO)

	respLRO := &longrunning.OperationInfo{
		ResponseType: "OutputType",
		MetadataType: "MetadataType",
	}
	respLROOpts := &descriptorpb.MethodOptions{}
	proto.SetExtension(respLROOpts, longrunning.E_OperationInfo, respLRO)

	for _, m := range []*descriptorpb.MethodDescriptorProto{
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
	} {
		t.Run(m.GetName(), func(t *testing.T) {
			g.pt.Reset()
			g.descInfo.ParentElement[m] = serv
			g.descInfo.ParentFile[m] = file
			g.aux = &auxTypes{
				methodToWrapper: map[*descriptorpb.MethodDescriptorProto]operationWrapper{},
				opWrappers:      map[string]operationWrapper{},
			}

			if err := g.genGRPCMethod("Foo", serv, m); err != nil {
				t.Fatal(err)
			}

			if err := g.genOperationBuilders(serv, "MyService"); err != nil {
				t.Fatal(err)
			}

			txtdiff.Diff(t, g.pt.String(), filepath.Join("testdata", "method_"+m.GetName()+".want"))
		})
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
	lroGetOp := &descriptorpb.MethodDescriptorProto{
		Name:       proto.String("GetOperation"),
		OutputType: proto.String(".google.longrunning.Operation"),
	}

	actualLRO := &descriptorpb.MethodDescriptorProto{
		Name:       proto.String("SuperLongRPC"),
		OutputType: proto.String(".google.longrunning.Operation"),
	}

	var g generator
	g.descInfo.ParentFile = map[protoreflect.ProtoMessage]*descriptorpb.FileDescriptorProto{
		lroGetOp: {
			Package: proto.String("google.longrunning"),
		},
		actualLRO: {
			Package: proto.String("my.pkg"),
		},
	}

	for _, tst := range []struct {
		in   *descriptorpb.MethodDescriptorProto
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
		m := &descriptorpb.MethodDescriptorProto{}
		if tst.pattern != "" {
			m.Options = &descriptorpb.MethodOptions{}

			setHTTPOption(m.Options, tst.pattern)
		}

		got := parseImplicitRequestHeaders(m)

		if diff := cmp.Diff(got, tst.want); diff != "" {
			t.Errorf("parseImplicitRequestHeaders(%s) = got(-), want(+):\n%s", tst.name, diff)
		}
	}
}

func Test_ContainsDynamicRequestHeaders(t *testing.T) {
	methodNoRule := &descriptorpb.MethodDescriptorProto{
		Name:       proto.String("MethodNoRule"),
		InputType:  proto.String(".my.pkg.InputType"),
		OutputType: proto.String(".google.longrunning.Operation"),
		Options:    &descriptorpb.MethodOptions{},
	}
	proto.SetExtension(methodNoRule.GetOptions(), annotations.E_Http, &annotations.HttpRule{
		Pattern: &annotations.HttpRule_Post{
			Post: "/v1/operations",
		},
	})
	methodOneRule := &descriptorpb.MethodDescriptorProto{
		Name:       proto.String("MethodOneRule"),
		InputType:  proto.String(".my.pkg.InputType"),
		OutputType: proto.String(".google.longrunning.Operation"),
		Options:    &descriptorpb.MethodOptions{},
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
	methodEmptyRule := &descriptorpb.MethodDescriptorProto{
		Options: &descriptorpb.MethodOptions{},
	}
	extRoutingEmptyRule := &annotations.RoutingRule{
		RoutingParameters: []*annotations.RoutingParameter{
			{},
		},
	}
	proto.SetExtension(methodEmptyRule.GetOptions(), annotations.E_Routing, extRoutingEmptyRule)
	methodMultipleRules := &descriptorpb.MethodDescriptorProto{
		Options: &descriptorpb.MethodOptions{},
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
		method *descriptorpb.MethodDescriptorProto
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
	topLevelEnum := &descriptorpb.EnumDescriptorProto{
		Name: proto.String("TopLevelEnum"),
	}

	extra := &descriptorpb.DescriptorProto{
		Name: proto.String("ExtraMessage"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name: proto.String("leaf"),
				Type: typep(descriptorpb.FieldDescriptorProto_TYPE_STRING),
			},
		},
	}

	inputType := &descriptorpb.DescriptorProto{
		Name: proto.String("InputType"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name: proto.String("str"),
				Type: typep(descriptorpb.FieldDescriptorProto_TYPE_STRING),
			},
			{
				Name: proto.String("bool"),
				Type: typep(descriptorpb.FieldDescriptorProto_TYPE_BOOL),
			},
			{
				Name: proto.String("int"),
				Type: typep(descriptorpb.FieldDescriptorProto_TYPE_INT32),
			},
			{
				Name: proto.String("double"),
				Type: typep(descriptorpb.FieldDescriptorProto_TYPE_DOUBLE),
			},
			{
				Name:     proto.String("msg"),
				Type:     typep(descriptorpb.FieldDescriptorProto_TYPE_MESSAGE),
				TypeName: proto.String("ExtraMessage"),
			},
			{
				Name:     proto.String("top_level_enum"),
				Type:     typep(descriptorpb.FieldDescriptorProto_TYPE_ENUM),
				TypeName: proto.String("TopLevelEnum"),
			},
			{
				Name:     proto.String("nested_enum"),
				Type:     typep(descriptorpb.FieldDescriptorProto_TYPE_ENUM),
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
		want             descriptorpb.FieldDescriptorProto_Type
	}{
		{
			name:  "string",
			msg:   "InputType",
			field: "str",
			want:  descriptorpb.FieldDescriptorProto_TYPE_STRING,
		},
		{
			name:  "boolean",
			msg:   "InputType",
			field: "bool",
			want:  descriptorpb.FieldDescriptorProto_TYPE_BOOL,
		},
		{
			name:  "integer",
			msg:   "InputType",
			field: "int",
			want:  descriptorpb.FieldDescriptorProto_TYPE_INT32,
		},
		{
			name:  "double",
			msg:   "InputType",
			field: "double",
			want:  descriptorpb.FieldDescriptorProto_TYPE_DOUBLE,
		},
		{
			name:  "nested field",
			msg:   "InputType",
			field: "msg.leaf",
			want:  descriptorpb.FieldDescriptorProto_TYPE_STRING,
		},
		{
			name:  "top level enum",
			msg:   "InputType",
			field: "top_level_enum",
			want:  descriptorpb.FieldDescriptorProto_TYPE_ENUM,
		},
	} {
		got := g.lookupField(tst.msg, tst.field)

		if got.GetType() != tst.want {
			t.Errorf("Test_lookupField(%s): got %v want %v", tst.name, got.GetType(), tst.want)
		}
	}
}

func TestGRPCStubCall(t *testing.T) {
	getFoo := &descriptorpb.MethodDescriptorProto{
		Name:       proto.String("GetFoo"),
		InputType:  proto.String("google.protobuf.Empty"),
		OutputType: proto.String("google.protobuf.Empty"),
	}
	getBar := &descriptorpb.MethodDescriptorProto{
		Name:       proto.String("GetBar"),
		InputType:  proto.String("google.protobuf.Empty"),
		OutputType: proto.String("google.protobuf.Empty"),
	}
	foo := &descriptorpb.FileDescriptorProto{
		Package: proto.String("google.foo.v1"),
		Options: &descriptorpb.FileOptions{
			GoPackage: proto.String("google.golang.org/genproto/googleapis/foo/v1"),
		},
		Dependency: []string{"google.protobuf.Empty"},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name:   proto.String("FooService"),
				Method: []*descriptorpb.MethodDescriptorProto{getFoo},
			},
			{
				Name:   proto.String("BarService"),
				Method: []*descriptorpb.MethodDescriptorProto{getBar},
			},
		},
	}
	var g generator
	err := g.init(&pluginpb.CodeGeneratorRequest{
		ProtoFile: []*descriptorpb.FileDescriptorProto{foo},
		Parameter: proto.String("go-gapic-package=cloud.google.com/go/foo/apiv1;foo"),
	})
	if err != nil {
		t.Error(err)
	}

	for _, tst := range []struct {
		name, want string
		in         *descriptorpb.MethodDescriptorProto
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
	op := &descriptorpb.DescriptorProto{
		Name: proto.String("Operation"),
	}
	foo := &descriptorpb.DescriptorProto{
		Name: proto.String("Foo"),
	}
	com := &descriptorpb.MethodDescriptorProto{
		OutputType: proto.String(".google.cloud.foo.v1.Operation"),
		Options:    &descriptorpb.MethodOptions{},
	}
	proto.SetExtension(com.GetOptions(), annotations.E_Http, &annotations.HttpRule{
		Pattern: &annotations.HttpRule_Post{
			Post: "/v1/operations",
		},
	})
	reg := &descriptorpb.MethodDescriptorProto{
		OutputType: proto.String(".google.cloud.foo.v1.Foo"),
		Options:    &descriptorpb.MethodOptions{},
	}
	proto.SetExtension(reg.GetOptions(), annotations.E_Http, &annotations.HttpRule{
		Pattern: &annotations.HttpRule_Post{
			Post: "/v1/operations",
		},
	})
	get := &descriptorpb.MethodDescriptorProto{
		OutputType: proto.String(".google.cloud.foo.v1.Operation"),
		Options:    &descriptorpb.MethodOptions{},
	}
	proto.SetExtension(get.GetOptions(), annotations.E_Http, &annotations.HttpRule{
		Pattern: &annotations.HttpRule_Get{
			Get: "/v1/operations",
		},
	})
	wait := &descriptorpb.MethodDescriptorProto{
		Name:       proto.String("Wait"),
		OutputType: proto.String(".google.cloud.foo.v1.Operation"),
		Options:    &descriptorpb.MethodOptions{},
	}
	proto.SetExtension(wait.GetOptions(), annotations.E_Http, &annotations.HttpRule{
		Pattern: &annotations.HttpRule_Post{
			Post: "/v1/operations",
		},
	})
	f := &descriptorpb.FileDescriptorProto{
		Package: proto.String("google.cloud.foo.v1"),
		Options: &descriptorpb.FileOptions{
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
			ParentFile: map[protoreflect.ProtoMessage]*descriptorpb.FileDescriptorProto{
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
		method     *descriptorpb.MethodDescriptorProto
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

func TestCollectServices(t *testing.T) {
	libraryServ := &descriptorpb.ServiceDescriptorProto{
		Name: proto.String("Library"),
	}
	library := &descriptorpb.FileDescriptorProto{
		Name: proto.String("google/cloud/library/v1/library.proto"),
		Options: &descriptorpb.FileOptions{
			GoPackage: proto.String("cloud.google.com/go/library/apiv1/librarypb;librarypb"),
		},
		Service: []*descriptorpb.ServiceDescriptorProto{libraryServ},
	}

	for _, tst := range []struct {
		name, goPkgPath string
		toGen           []string
		fileSet         []*descriptorpb.FileDescriptorProto
		want            []*descriptorpb.ServiceDescriptorProto
	}{
		{
			name:      "simple",
			goPkgPath: "cloud.google.com/go/library/apiv1",
			toGen:     []string{library.GetName()},
			fileSet:   []*descriptorpb.FileDescriptorProto{library},
			want:      []*descriptorpb.ServiceDescriptorProto{libraryServ},
		},
		{
			name:      "ignore-mixins",
			goPkgPath: "cloud.google.com/go/library/apiv1",
			toGen:     []string{library.GetName()},
			fileSet: []*descriptorpb.FileDescriptorProto{
				library,
				mixinFiles["google.longrunning.Operations"][0],
				mixinFiles["google.iam.v1.IAMPolicy"][0],
				mixinFiles["google.cloud.location.Locations"][0],
			},
			want: []*descriptorpb.ServiceDescriptorProto{libraryServ},
		},
		{
			name:      "include-iam-mixin",
			goPkgPath: "cloud.google.com/go/iam/apiv1",
			toGen:     []string{mixinFiles["google.iam.v1.IAMPolicy"][0].GetName()},
			fileSet: []*descriptorpb.FileDescriptorProto{
				mixinFiles["google.longrunning.Operations"][0],
				mixinFiles["google.iam.v1.IAMPolicy"][0],
				mixinFiles["google.cloud.location.Locations"][0],
			},
			want: []*descriptorpb.ServiceDescriptorProto{mixinFiles["google.iam.v1.IAMPolicy"][0].GetService()[0]},
		},
	} {
		g := &generator{opts: &options{pkgPath: tst.goPkgPath}}
		got := g.collectServices(&pluginpb.CodeGeneratorRequest{
			FileToGenerate: tst.toGen,
			ProtoFile:      tst.fileSet,
		})
		if diff := cmp.Diff(got, tst.want, cmp.Comparer(proto.Equal)); diff != "" {
			t.Errorf("%s: got(-),want(+):\n%s", tst.name, diff)
		}
	}
}

func setHTTPOption(o *descriptorpb.MethodOptions, pattern string) {
	proto.SetExtension(o, annotations.E_Http, &annotations.HttpRule{
		Pattern: &annotations.HttpRule_Get{
			Get: pattern,
		},
	})
}
