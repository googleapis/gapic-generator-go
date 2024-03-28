// Copyright 2023 Google LLC
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
	"fmt"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"github.com/googleapis/gapic-generator-go/internal/snippets"
	"github.com/googleapis/gapic-generator-go/internal/txtdiff"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func TestSnippetsOutDir(t *testing.T) {
	for _, tst := range []struct {
		opts options
		want string
	}{
		{
			opts: options{
				outDir:  "cloud.google.com/go/video/stitcher/apiv1",
				pkgPath: "cloud.google.com/go/video/stitcher/apiv1",
			},
			want: "cloud.google.com/go/internal/generated/snippets/video/stitcher/apiv1",
		},
		{
			opts: options{
				outDir:  "example.com/my/package",
				pkgPath: "example.com/my/package",
			},
			want: "example.com/my/package/internal/snippets",
		},
	} {
		var g generator
		g.opts = &tst.opts
		if s := g.snippetsOutDir(); s != tst.want {
			t.Errorf("TestGenAndCommitSnippets(g.opts.pkgPath = %s): got %s, want %s", g.opts.pkgPath, s, tst.want)
		}
	}
}

func TestGenAndCommitSnippets(t *testing.T) {
	inputType := &descriptorpb.DescriptorProto{
		Name: proto.String("InputType"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:  proto.String("biz"),
				Type:  typep(descriptorpb.FieldDescriptorProto_TYPE_DOUBLE),
				Label: labelp(descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL),
			},
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
		pkgName:    "pkg",
		transports: []transport{grpc},
	}
	// TODO(chrisdsmith): Add a test case below for ListLocations or similar.
	g.mixins = mixins{
		"google.longrunning.Operations":   operationsMethods(),
		"google.cloud.location.Locations": locationMethods(),
		"google.iam.v1.IAMPolicy":         iamPolicyMethods(),
	}
	g.imports = map[pbinfo.ImportSpec]bool{}
	commonTypes(&g)
	for _, typ := range []pbinfo.ProtoType{
		inputType, outputType, pageInputType, pageOutputType,
	} {
		g.descInfo.Type[".my.pkg."+typ.GetName()] = typ
		g.descInfo.ParentFile[typ] = file
	}
	g.descInfo.ParentFile[serv] = file
	g.descInfo.ParentElement = map[pbinfo.ProtoType]pbinfo.ProtoType{
		paginatedField: pageOutputType,
	}
	g.snippetMetadata = snippets.NewMetadata("my.pkg", "github.com/googleapis/mypackage", "pkg")

	for _, tst := range []struct {
		wantNil bool
		m       *descriptorpb.MethodDescriptorProto
		imports map[pbinfo.ImportSpec]bool
	}{
		{
			m: &descriptorpb.MethodDescriptorProto{
				Name:       proto.String("GetEmptyThing"),
				InputType:  proto.String(".my.pkg.InputType"),
				OutputType: proto.String(emptyType),
			},
			imports: map[pbinfo.ImportSpec]bool{
				{Path: "context"}:                        true,
				{Name: "mypackagepb", Path: "mypackage"}: true,
				{Name: "pkg"}:                            true,
			},
		},
		{
			m: &descriptorpb.MethodDescriptorProto{
				Name:       proto.String("GetOneThing"),
				InputType:  proto.String(".my.pkg.InputType"),
				OutputType: proto.String(".my.pkg.OutputType"),
			},
			imports: map[pbinfo.ImportSpec]bool{
				{Path: "context"}:                        true,
				{Name: "mypackagepb", Path: "mypackage"}: true,
				{Name: "pkg"}:                            true,
			},
		},
		{
			m: &descriptorpb.MethodDescriptorProto{
				Name:       proto.String("GetManyThings"),
				InputType:  proto.String(".my.pkg.PageInputType"),
				OutputType: proto.String(".my.pkg.PageOutputType"),
			},
			imports: map[pbinfo.ImportSpec]bool{
				{Path: "context"}:                        true,
				{Path: "google.golang.org/api/iterator"}: true,
				{Name: "mypackagepb", Path: "mypackage"}: true,
				{Name: "pkg"}:                            true,
			},
		},
		{
			// TODO(chrisdsmith): implement streaming examples correctly, see example.go TODOs.
			wantNil: true,
			m: &descriptorpb.MethodDescriptorProto{
				Name:            proto.String("ServerThings"),
				InputType:       proto.String(".my.pkg.InputType"),
				OutputType:      proto.String(".my.pkg.OutputType"),
				ServerStreaming: proto.Bool(true),
			},
			imports: map[pbinfo.ImportSpec]bool{},
		},
		{
			// TODO(chrisdsmith): implement streaming examples correctly, see example.go TODOs.
			wantNil: true,
			m: &descriptorpb.MethodDescriptorProto{
				Name:            proto.String("ClientThings"),
				InputType:       proto.String(".my.pkg.InputType"),
				OutputType:      proto.String(".my.pkg.OutputType"),
				ClientStreaming: proto.Bool(true),
			},
			imports: map[pbinfo.ImportSpec]bool{},
		},
		{
			m: &descriptorpb.MethodDescriptorProto{
				Name:            proto.String("BidiThings"),
				InputType:       proto.String(".my.pkg.InputType"),
				OutputType:      proto.String(".my.pkg.OutputType"),
				ServerStreaming: proto.Bool(true),
				ClientStreaming: proto.Bool(true),
			},
			imports: map[pbinfo.ImportSpec]bool{
				{Path: "context"}:                        true,
				{Path: "io"}:                             true,
				{Name: "mypackagepb", Path: "mypackage"}: true,
				{Name: "pkg"}:                            true,
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

			if err := g.genAndCommitSnippets(serv); err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(g.imports, tst.imports); diff != "" {
				t.Errorf("TestGenAndCommitSnippets(%s): imports got(-),want(+):\n%s", tst.m.GetName(), diff)
			}
			if !tst.wantNil {
				wantRegionTag := fmt.Sprintf("// [START _pkg_generated_Foo_%s_sync]\n", tst.m.GetName())
				if g.headerComments.String() != wantRegionTag {
					t.Errorf("TestGenAndCommitSnippets(%s): got %s, want %s", tst.m.GetName(), g.headerComments.String(), wantRegionTag)
				}
				txtdiff.Diff(t, g.pt.String(), filepath.Join("testdata", "snippet_"+tst.m.GetName()+".want"))
			}
		})
	}
}
