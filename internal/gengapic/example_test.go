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

	longrunning "cloud.google.com/go/longrunning/autogen/longrunningpb"
	"github.com/google/go-cmp/cmp"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"github.com/googleapis/gapic-generator-go/internal/snippets"
	"github.com/googleapis/gapic-generator-go/internal/testing/sample"
	"github.com/googleapis/gapic-generator-go/internal/txtdiff"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/genproto/googleapis/api/serviceconfig"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/apipb"
)

func TestExample(t *testing.T) {
	var g generator
	g.imports = map[pbinfo.ImportSpec]bool{}
	g.aux = &auxTypes{}
	mix := mixins{
		"google.longrunning.Operations":   operationsMethods(),
		"google.cloud.location.Locations": locationMethods(),
		"google.iam.v1.IAMPolicy":         iamPolicyMethods(),
	}
	g.mixins = mix
	g.serviceConfig = &serviceconfig.Service{
		Apis: []*apipb.Api{
			{Name: "foo.bar.Baz"},
			{Name: "google.iam.v1.IAMPolicy"},
			{Name: "google.cloud.location.Locations"},
			{Name: "google.longrunning.Operations"},
		},
	}

	inputType := &descriptorpb.DescriptorProto{
		Name: proto.String("InputType"),
	}
	outputType := &descriptorpb.DescriptorProto{
		Name: proto.String("OutputType"),
	}

	pageInputType := &descriptorpb.DescriptorProto{
		Name: proto.String("PageInputType"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:  proto.String("page_size"),
				Type:  typep(descriptorpb.FieldDescriptorProto_TYPE_INT32),
				Label: labelp(descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL),
			},
			{
				Name:  proto.String("page_token"),
				Type:  typep(descriptorpb.FieldDescriptorProto_TYPE_STRING),
				Label: labelp(descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL),
			},
		},
	}
	pageOutputType := &descriptorpb.DescriptorProto{
		Name: proto.String("PageOutputType"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:  proto.String("next_page_token"),
				Type:  typep(descriptorpb.FieldDescriptorProto_TYPE_STRING),
				Label: labelp(descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL),
			},
			{
				Name:  proto.String("items"),
				Type:  typep(descriptorpb.FieldDescriptorProto_TYPE_STRING),
				Label: labelp(descriptorpb.FieldDescriptorProto_LABEL_REPEATED),
			},
		},
	}

	cop := &descriptorpb.DescriptorProto{
		Name: proto.String("Operation"),
	}

	file := &descriptorpb.FileDescriptorProto{
		Options: &descriptorpb.FileOptions{
			GoPackage: proto.String("mypackage"),
		},
		Package: proto.String("my.pkg"),
	}

	emptyLRO := &longrunning.OperationInfo{
		ResponseType: emptyValue,
	}
	emptyLROOpts := &descriptorpb.MethodOptions{}
	proto.SetExtension(emptyLROOpts, longrunning.E_OperationInfo, emptyLRO)

	respLRO := &longrunning.OperationInfo{
		ResponseType: "my.pkg.OutputType",
	}
	respLROOpts := &descriptorpb.MethodOptions{}
	proto.SetExtension(respLROOpts, longrunning.E_OperationInfo, respLRO)

	customOpOpts := &descriptorpb.MethodOptions{}
	proto.SetExtension(customOpOpts, annotations.E_Http, &annotations.HttpRule{
		Pattern: &annotations.HttpRule_Post{
			Post: "/customOp",
		},
	})

	commonTypes(&g)
	for _, typ := range []*descriptorpb.DescriptorProto{
		inputType, outputType, pageInputType, pageOutputType, cop,
	} {
		g.descInfo.Type[".my.pkg."+typ.GetName()] = typ
		g.descInfo.ParentFile[typ] = file
	}

	serv := &descriptorpb.ServiceDescriptorProto{
		Name: proto.String("Foo"),
		Method: []*descriptorpb.MethodDescriptorProto{
			{
				Name:       proto.String("GetEmptyThing"),
				InputType:  proto.String(".my.pkg.InputType"),
				OutputType: proto.String(emptyType),
			},
			{
				Name:       proto.String("GetOneThing"),
				InputType:  proto.String(".my.pkg.InputType"),
				OutputType: proto.String(".my.pkg.OutputType"),
			},
			{
				Name:       proto.String("GetBigThing"),
				InputType:  proto.String(".my.pkg.InputType"),
				OutputType: proto.String(".google.longrunning.Operation"),
				Options:    respLROOpts,
			},
			{
				Name:       proto.String("GetManyThings"),
				InputType:  proto.String(".my.pkg.PageInputType"),
				OutputType: proto.String(".my.pkg.PageOutputType"),
			},
			{
				Name:            proto.String("BidiThings"),
				InputType:       proto.String(".my.pkg.InputType"),
				OutputType:      proto.String(".my.pkg.OutputType"),
				ServerStreaming: proto.Bool(true),
				ClientStreaming: proto.Bool(true),
			},
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
			{
				Name:       proto.String("CustomOp"),
				InputType:  proto.String(".my.pkg.InputType"),
				OutputType: proto.String(".my.pkg.Operation"),
				Options:    customOpOpts,
			},
		},
	}
	for _, tst := range []struct {
		tstName string
		imports map[pbinfo.ImportSpec]bool
		options options
		op      *customOp
	}{
		{
			tstName: "empty_example",
			options: options{
				pkgName:    "Foo",
				transports: []transport{grpc, rest},
			},
			imports: map[pbinfo.ImportSpec]bool{
				{Path: "context"}:                        true,
				{Path: "google.golang.org/api/iterator"}: true,
				{Path: "io"}:                             true,
				{Name: "iampb", Path: "cloud.google.com/go/iam/apiv1/iampb"}:                           true,
				{Name: "locationpb", Path: "google.golang.org/genproto/googleapis/cloud/location"}:     true,
				{Name: "longrunningpb", Path: "cloud.google.com/go/longrunning/autogen/longrunningpb"}: true,
				{Name: "mypackagepb", Path: "mypackage"}:                                               true,
			},
		},
		{
			tstName: "empty_example_grpc",
			options: options{
				pkgName:    "Foo",
				transports: []transport{grpc},
			},
			imports: map[pbinfo.ImportSpec]bool{
				{Path: "context"}:                        true,
				{Path: "google.golang.org/api/iterator"}: true,
				{Path: "io"}:                             true,
				{Name: "iampb", Path: "cloud.google.com/go/iam/apiv1/iampb"}:                           true,
				{Name: "locationpb", Path: "google.golang.org/genproto/googleapis/cloud/location"}:     true,
				{Name: "longrunningpb", Path: "cloud.google.com/go/longrunning/autogen/longrunningpb"}: true,
				{Name: "mypackagepb", Path: "mypackage"}:                                               true,
			},
		},
		{
			tstName: "foo_example",
			options: options{
				pkgName:    "Bar",
				transports: []transport{grpc, rest},
			},
			imports: map[pbinfo.ImportSpec]bool{
				{Path: "context"}:                        true,
				{Path: "google.golang.org/api/iterator"}: true,
				{Path: "io"}:                             true,
				{Name: "iampb", Path: "cloud.google.com/go/iam/apiv1/iampb"}:                           true,
				{Name: "locationpb", Path: "google.golang.org/genproto/googleapis/cloud/location"}:     true,
				{Name: "longrunningpb", Path: "cloud.google.com/go/longrunning/autogen/longrunningpb"}: true,
				{Name: "mypackagepb", Path: "mypackage"}:                                               true,
			},
		},
		{
			tstName: "foo_example_rest",
			options: options{
				pkgName:    "Bar",
				transports: []transport{rest},
			},
			imports: map[pbinfo.ImportSpec]bool{
				{Path: "context"}:                        true,
				{Path: "google.golang.org/api/iterator"}: true,
				{Path: "io"}:                             true,
				{Name: "iampb", Path: "cloud.google.com/go/iam/apiv1/iampb"}:                           true,
				{Name: "locationpb", Path: "google.golang.org/genproto/googleapis/cloud/location"}:     true,
				{Name: "longrunningpb", Path: "cloud.google.com/go/longrunning/autogen/longrunningpb"}: true,
				{Name: "mypackagepb", Path: "mypackage"}:                                               true,
			},
		},
		{
			tstName: "custom_op_example",
			options: options{
				pkgName:    "Bar",
				transports: []transport{rest},
				diregapic:  true,
			},
			imports: map[pbinfo.ImportSpec]bool{
				{Path: "context"}:                        true,
				{Path: "google.golang.org/api/iterator"}: true,
				{Path: "io"}:                             true,
				{Name: "mypackagepb", Path: "mypackage"}: true,
			},
			op: &customOp{message: cop},
		},
	} {
		t.Run(tst.tstName, func(t *testing.T) {
			g.reset()
			g.opts = &tst.options
			g.mixins = mix
			if tst.options.diregapic {
				g.mixins = nil
			}
			g.aux.customOp = tst.op
			g.genExampleFile(serv)
			if diff := cmp.Diff(g.imports, tst.imports); diff != "" {
				t.Errorf("TestExample(%s): imports got(-),want(+):\n%s", tst.tstName, diff)
			}
			txtdiff.Diff(t, g.pt.String(), filepath.Join("testdata", tst.tstName+".want"))
		})
	}
}

func TestGenSnippetFile(t *testing.T) {
	var g generator
	g.imports = map[pbinfo.ImportSpec]bool{}
	g.serviceConfig = sample.ServiceConfig()

	serv := sample.Service()
	g.snippetMetadata = snippets.NewMetadata(sample.ProtoPackagePath, sample.GoPackagePath, sample.GoPackageName)
	g.snippetMetadata.AddService(serv.GetName(), sample.ServiceURL)

	inputType := sample.InputType(sample.CreateRequest)
	outputType := sample.OutputType(sample.Resource)

	files := []*descriptorpb.FileDescriptorProto{}
	g.descInfo = pbinfo.Of(files)

	file := sample.File()
	for _, typ := range []*descriptorpb.DescriptorProto{
		inputType, outputType,
	} {
		g.descInfo.Type[sample.DescriptorInfoTypeName(typ.GetName())] = typ
		g.descInfo.ParentFile[typ] = file
	}

	for _, test := range []struct {
		name    string
		options options
		imports map[pbinfo.ImportSpec]bool
	}{
		{
			name: "snippet",
			options: options{
				pkgName:    "migration",
				transports: []transport{grpc, rest},
			},
			imports: map[pbinfo.ImportSpec]bool{
				{Path: "context"}: true,
				{Name: sample.GoProtoPackageName, Path: sample.GoProtoPackagePath}: true,
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			g.reset()
			g.opts = &test.options
			err := g.genSnippetFile(serv, serv.Method[0])
			if err != nil {
				t.Fatal(err)
			}
			g.commit(filepath.Join("cloud.google.com/go", "internal", "generated", "snippets", "bigquery", "main.go"), "main")
			if diff := cmp.Diff(test.imports, g.imports); diff != "" {
				t.Errorf("TestExample(%s) imports mismatch: (-want +got):\n%s", test.name, diff)
			}
			got := *g.resp.File[0].Content + *g.resp.File[1].Content
			txtdiff.Diff(t, got, filepath.Join("testdata", test.name+".want"))
		})
	}
}

func commonTypes(g *generator) {
	empty := &descriptorpb.DescriptorProto{
		Name: proto.String("Empty"),
	}
	emptyFile := &descriptorpb.FileDescriptorProto{
		Package: proto.String("google.protobuf"),
		Options: &descriptorpb.FileOptions{
			GoPackage: proto.String("google.golang.org/protobuf/types/known/emptypb"),
		},
		MessageType: []*descriptorpb.DescriptorProto{empty},
	}

	files := append(g.getMixinFiles(), emptyFile)

	g.descInfo = pbinfo.Of(files)
}
