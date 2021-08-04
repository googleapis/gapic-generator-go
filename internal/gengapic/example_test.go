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

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/google/go-cmp/cmp"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"github.com/googleapis/gapic-generator-go/internal/txtdiff"
	"google.golang.org/genproto/googleapis/api/serviceconfig"
	"google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/apipb"
)

func TestExample(t *testing.T) {
	var g generator
	g.opts = &options{}
	g.imports = map[pbinfo.ImportSpec]bool{}
	g.mixins = mixins{
		"google.longrunning.Operations":   operationsMethods(),
		"google.cloud.location.Locations": locationMethods(),
		"google.iam.v1.IAMPolicy":         iamPolicyMethods(),
	}
	g.serviceConfig = &serviceconfig.Service{
		Apis: []*apipb.Api{
			{Name: "foo.bar.Baz"},
			{Name: "google.iam.v1.IAMPolicy"},
			{Name: "google.cloud.location.Locations"},
			{Name: "google.longrunning.Operations"},
		},
	}

	inputType := &descriptor.DescriptorProto{
		Name: proto.String("InputType"),
	}
	outputType := &descriptor.DescriptorProto{
		Name: proto.String("OutputType"),
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
	pageOutputType := &descriptor.DescriptorProto{
		Name: proto.String("PageOutputType"),
		Field: []*descriptor.FieldDescriptorProto{
			{
				Name:  proto.String("next_page_token"),
				Type:  typep(descriptor.FieldDescriptorProto_TYPE_STRING),
				Label: labelp(descriptor.FieldDescriptorProto_LABEL_OPTIONAL),
			},
			{
				Name:  proto.String("items"),
				Type:  typep(descriptor.FieldDescriptorProto_TYPE_STRING),
				Label: labelp(descriptor.FieldDescriptorProto_LABEL_REPEATED),
			},
		},
	}

	file := &descriptor.FileDescriptorProto{
		Options: &descriptor.FileOptions{
			GoPackage: proto.String("mypackage"),
		},
	}

	emptyLRO := &longrunning.OperationInfo{
		ResponseType: emptyValue,
	}
	emptyLROOpts := &descriptor.MethodOptions{}
	proto.SetExtension(emptyLROOpts, longrunning.E_OperationInfo, emptyLRO)

	respLRO := &longrunning.OperationInfo{
		ResponseType: "my.pkg.OutputType",
	}
	respLROOpts := &descriptor.MethodOptions{}
	proto.SetExtension(respLROOpts, longrunning.E_OperationInfo, respLRO)

	commonTypes(&g)
	for _, typ := range []*descriptor.DescriptorProto{
		inputType, outputType, pageInputType, pageOutputType,
	} {
		g.descInfo.Type[".my.pkg."+*typ.Name] = typ
		g.descInfo.ParentFile[typ] = file
	}

	serv := &descriptor.ServiceDescriptorProto{
		Name: proto.String("Foo"),
		Method: []*descriptor.MethodDescriptorProto{
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
		},
	}
	for _, tst := range []struct {
		tstName, pkgName string
		transports       []transport
		imports          map[pbinfo.ImportSpec]bool
	}{
		{
			tstName:    "empty_example",
			pkgName:    "Foo",
			transports: []transport{grpc, rest},
			imports: map[pbinfo.ImportSpec]bool{
				{Path: "context"}:                        true,
				{Path: "google.golang.org/api/iterator"}: true,
				{Path: "io"}:                             true,
				{Name: "iampb", Path: "google.golang.org/genproto/googleapis/iam/v1"}:              true,
				{Name: "locationpb", Path: "google.golang.org/genproto/googleapis/cloud/location"}: true,
				{Name: "longrunningpb", Path: "google.golang.org/genproto/googleapis/longrunning"}: true,
				{Name: "mypackagepb", Path: "mypackage"}:                                           true,
			},
		},
		{
			tstName:    "empty_example_grpc",
			pkgName:    "Foo",
			transports: []transport{grpc},
			imports: map[pbinfo.ImportSpec]bool{
				{Path: "context"}:                        true,
				{Path: "google.golang.org/api/iterator"}: true,
				{Path: "io"}:                             true,
				{Name: "iampb", Path: "google.golang.org/genproto/googleapis/iam/v1"}:              true,
				{Name: "locationpb", Path: "google.golang.org/genproto/googleapis/cloud/location"}: true,
				{Name: "longrunningpb", Path: "google.golang.org/genproto/googleapis/longrunning"}: true,
				{Name: "mypackagepb", Path: "mypackage"}:                                           true,
			},
		},
		{
			tstName:    "foo_example",
			pkgName:    "Bar",
			transports: []transport{grpc, rest},
			imports: map[pbinfo.ImportSpec]bool{
				{Path: "context"}:                        true,
				{Path: "google.golang.org/api/iterator"}: true,
				{Path: "io"}:                             true,
				{Name: "iampb", Path: "google.golang.org/genproto/googleapis/iam/v1"}:              true,
				{Name: "locationpb", Path: "google.golang.org/genproto/googleapis/cloud/location"}: true,
				{Name: "longrunningpb", Path: "google.golang.org/genproto/googleapis/longrunning"}: true,
				{Name: "mypackagepb", Path: "mypackage"}:                                           true,
			},
		},
		{
			tstName:    "foo_example_rest",
			pkgName:    "Bar",
			transports: []transport{rest},
			imports: map[pbinfo.ImportSpec]bool{
				{Path: "context"}:                        true,
				{Path: "google.golang.org/api/iterator"}: true,
				{Path: "io"}:                             true,
				{Name: "iampb", Path: "google.golang.org/genproto/googleapis/iam/v1"}:              true,
				{Name: "locationpb", Path: "google.golang.org/genproto/googleapis/cloud/location"}: true,
				{Name: "longrunningpb", Path: "google.golang.org/genproto/googleapis/longrunning"}: true,
				{Name: "mypackagepb", Path: "mypackage"}:                                           true,
			},
		},
	} {
		g.reset()
		g.opts.pkgName = tst.pkgName
		g.opts.transports = tst.transports
		g.genExampleFile(serv)
		if diff := cmp.Diff(g.imports, tst.imports); diff != "" {
			t.Errorf("TestExample(%s): imports got(-),want(+):\n%s", tst.tstName, diff)
		}
		txtdiff.Diff(t, tst.tstName, g.pt.String(), filepath.Join("testdata", tst.tstName+".want"))
	}
}

func commonTypes(g *generator) {
	empty := &descriptor.DescriptorProto{
		Name: proto.String("Empty"),
	}
	emptyFile := &descriptor.FileDescriptorProto{
		Package: proto.String("google.protobuf"),
		Options: &descriptor.FileOptions{
			GoPackage: proto.String("google.golang.org/protobuf/types/known/emptypb"),
		},
		MessageType: []*descriptor.DescriptorProto{empty},
	}

	files := append(g.getMixinFiles(), emptyFile)

	g.descInfo = pbinfo.Of(files)
}
