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
	"strings"
	"testing"

	"cloud.google.com/go/longrunning/autogen/longrunningpb"
	"github.com/google/go-cmp/cmp"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"github.com/googleapis/gapic-generator-go/internal/txtdiff"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestWrapperExists(t *testing.T) {
	testCases := []struct {
		desc, wantErrStr string
		existing         map[string]operationWrapper
		target           operationWrapper
		want             bool
	}{
		{
			desc: "existing match",
			existing: map[string]operationWrapper{
				"FooOperation": {
					name:         "FooOperation",
					responseName: protoreflect.FullName("google.example.v1.Foo"),
					metadataName: protoreflect.FullName("google.example.v1.FooMetadata"),
				},
			},
			target: operationWrapper{
				name:         "FooOperation",
				responseName: protoreflect.FullName("google.example.v1.Foo"),
				metadataName: protoreflect.FullName("google.example.v1.FooMetadata"),
			},
			want: true,
		},
		{
			desc: "existing response mismatch",
			existing: map[string]operationWrapper{
				"FooOperation": {
					name:         "FooOperation",
					responseName: protoreflect.FullName("google.example.v1.Bar"),
					metadataName: protoreflect.FullName("google.example.v1.FooMetadata"),
				},
			},
			target: operationWrapper{
				name:         "FooOperation",
				responseName: protoreflect.FullName("google.example.v1.Foo"),
				metadataName: protoreflect.FullName("google.example.v1.FooMetadata"),
			},
			want:       true,
			wantErrStr: "mismatched response_types",
		},
		{
			desc: "existing metadata mismatch",
			existing: map[string]operationWrapper{
				"FooOperation": {
					name:         "FooOperation",
					responseName: protoreflect.FullName("google.example.v1.Foo"),
					metadataName: protoreflect.FullName("google.example.v1.BarMetadata"),
				},
			},
			target: operationWrapper{
				name:         "FooOperation",
				responseName: protoreflect.FullName("google.example.v1.Foo"),
				metadataName: protoreflect.FullName("google.example.v1.FooMetadata"),
			},
			want:       true,
			wantErrStr: "mismatched metadata_types",
		},
		{
			desc:     "doesn't exist",
			existing: map[string]operationWrapper{},
			target: operationWrapper{
				name:         "FooOperation",
				responseName: protoreflect.FullName("google.example.v1.Foo"),
				metadataName: protoreflect.FullName("google.example.v1.FooMetadata"),
			},
			want: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			a := auxTypes{
				opWrappers: tC.existing,
			}
			got, gotErr := a.wrapperExists(tC.target)
			if tC.want != got {
				t.Errorf("wrapperExists(%+v): got %v, want %v, with existing: %+v", tC.target, got, tC.want, tC.existing)
			}
			if tC.wantErrStr == "" && gotErr != nil {
				t.Errorf("wrapperExists(%+v): got error %v, want no error, with existing: %+v", tC.target, gotErr, tC.existing)
			}
			if tC.wantErrStr != "" && gotErr == nil {
				t.Errorf("wrapperExists(%+v): got no error, want an error with %q, with existing: %+v", tC.target, tC.wantErrStr, tC.existing)
			}
			if tC.wantErrStr != "" && gotErr != nil && !strings.Contains(gotErr.Error(), tC.wantErrStr) {
				t.Errorf("wrapperExists(%+v): got %q, want %q, with existing: %+v", tC.target, gotErr.Error(), tC.wantErrStr, tC.existing)
			}
		})
	}
}

func TestMaybeAddOperationWrapper(t *testing.T) {
	testCases := []struct {
		desc, wantErrStr                           string
		existing, want                             map[string]operationWrapper
		res, meta, oiRes, oiMeta, m, pkg, otherPkg string
	}{
		{
			desc:     "add new",
			existing: map[string]operationWrapper{},
			want: map[string]operationWrapper{
				"CreateFooOperation": {
					name: "CreateFooOperation",
					response: &descriptorpb.DescriptorProto{
						Name: proto.String("Foo"),
					},
					metadata: &descriptorpb.DescriptorProto{
						Name: proto.String("FooMetadata"),
					},
					responseName: protoreflect.FullName("google.example.v1.Foo"),
					metadataName: protoreflect.FullName("google.example.v1.FooMetadata"),
				},
			},
			res:    "Foo",
			meta:   "FooMetadata",
			oiRes:  "Foo",
			oiMeta: "FooMetadata",
			m:      "CreateFoo",
			pkg:    "google.example.v1",
		},
		{
			desc:     "add new fully qualified",
			existing: map[string]operationWrapper{},
			want: map[string]operationWrapper{
				"CreateFooOperation": {
					name: "CreateFooOperation",
					response: &descriptorpb.DescriptorProto{
						Name: proto.String("Foo"),
					},
					metadata: &descriptorpb.DescriptorProto{
						Name: proto.String("FooMetadata"),
					},
					responseName: protoreflect.FullName("google.example.v1.Foo"),
					metadataName: protoreflect.FullName("google.example.v1.FooMetadata"),
				},
			},
			res:    "Foo",
			meta:   "FooMetadata",
			oiRes:  "google.example.v1.Foo",
			oiMeta: "google.example.v1.FooMetadata",
			m:      "CreateFoo",
			pkg:    "google.example.v1",
		},
		{
			desc:     "add new different package",
			existing: map[string]operationWrapper{},
			want: map[string]operationWrapper{
				"CreateFooOperation": {
					name: "CreateFooOperation",
					response: &descriptorpb.DescriptorProto{
						Name: proto.String("Foo"),
					},
					metadata: &descriptorpb.DescriptorProto{
						Name: proto.String("FooMetadata"),
					},
					responseName: protoreflect.FullName("google.other.v1.Foo"),
					metadataName: protoreflect.FullName("google.other.v1.FooMetadata"),
				},
			},
			res:      "Foo",
			meta:     "FooMetadata",
			oiRes:    "google.other.v1.Foo",
			oiMeta:   "google.other.v1.FooMetadata",
			m:        "CreateFoo",
			pkg:      "google.example.v1",
			otherPkg: "google.other.v1",
		},
		{
			desc:       "err missing response_type",
			existing:   map[string]operationWrapper{},
			want:       map[string]operationWrapper{},
			wantErrStr: "missing option google.longrunning.operation_info.response_type",
			res:        "Foo",
			meta:       "FooMetadata",
			oiRes:      "",
			oiMeta:     "FooMetadata",
			m:          "CreateFoo",
			pkg:        "google.example.v1",
		},
		{
			desc:       "err missing metadata_type",
			existing:   map[string]operationWrapper{},
			want:       map[string]operationWrapper{},
			wantErrStr: "missing option google.longrunning.operation_info.metadata_type",
			res:        "Foo",
			meta:       "FooMetadata",
			oiRes:      "Foo",
			oiMeta:     "",
			m:          "CreateFoo",
			pkg:        "google.example.v1",
		},
		{
			desc:       "err unresolvable response_type",
			existing:   map[string]operationWrapper{},
			want:       map[string]operationWrapper{},
			wantErrStr: "unable to resolve google.longrunning.operation_info.response_type",
			res:        "Foo",
			meta:       "FooMetadata",
			oiRes:      "DoesNotExist",
			oiMeta:     "FooMetadata",
			m:          "CreateFoo",
			pkg:        "google.example.v1",
		},
		{
			desc:       "err unresolvable metadata_type",
			existing:   map[string]operationWrapper{},
			want:       map[string]operationWrapper{},
			wantErrStr: "unable to resolve google.longrunning.operation_info.metadata_type",
			res:        "Foo",
			meta:       "FooMetadata",
			oiRes:      "Foo",
			oiMeta:     "DoesNotExist",
			m:          "CreateFoo",
			pkg:        "google.example.v1",
		},
		{
			desc: "err mismatch collision",
			existing: map[string]operationWrapper{
				"CreateFooOperation": {
					name: "CreateFooOperation",
					response: &descriptorpb.DescriptorProto{
						Name: proto.String("Bar"),
					},
					metadata: &descriptorpb.DescriptorProto{
						Name: proto.String("FooMetadata"),
					},
					responseName: protoreflect.FullName("google.other.v1.Bar"),
					metadataName: protoreflect.FullName("google.other.v1.FooMetadata"),
				},
			},
			want: map[string]operationWrapper{
				"CreateFooOperation": {
					name: "CreateFooOperation",
					response: &descriptorpb.DescriptorProto{
						Name: proto.String("Bar"),
					},
					metadata: &descriptorpb.DescriptorProto{
						Name: proto.String("FooMetadata"),
					},
					responseName: protoreflect.FullName("google.other.v1.Bar"),
					metadataName: protoreflect.FullName("google.other.v1.FooMetadata"),
				},
			},
			wantErrStr: "duplicate operation wrapper types",
			res:        "Foo",
			meta:       "FooMetadata",
			oiRes:      "Foo",
			oiMeta:     "FooMetadata",
			m:          "CreateFoo",
			pkg:        "google.example.v1",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			r := &descriptorpb.DescriptorProto{
				Name: proto.String(tC.res),
			}
			md := &descriptorpb.DescriptorProto{
				Name: proto.String(tC.meta),
			}
			oi := &longrunningpb.OperationInfo{
				ResponseType: tC.oiRes,
				MetadataType: tC.oiMeta,
			}
			m := &descriptorpb.MethodDescriptorProto{
				Name:    proto.String(tC.m),
				Options: &descriptorpb.MethodOptions{},
			}
			proto.SetExtension(m.GetOptions(), longrunningpb.E_OperationInfo, oi)

			f := &descriptorpb.FileDescriptorProto{
				Package: proto.String(tC.pkg),
			}
			messageParent := f
			if tC.otherPkg != "" {
				of := &descriptorpb.FileDescriptorProto{
					Package: proto.String(tC.otherPkg),
				}
				messageParent = of
			}

			g := &generator{
				descInfo: pbinfo.Info{
					Type: map[string]pbinfo.ProtoType{
						fmt.Sprintf(".%s.%s", messageParent.GetPackage(), r.GetName()):  r,
						fmt.Sprintf(".%s.%s", messageParent.GetPackage(), md.GetName()): md,
					},
					ParentFile: map[protoreflect.ProtoMessage]*descriptorpb.FileDescriptorProto{
						m:  f,
						r:  messageParent,
						md: messageParent,
					},
				},
				aux: &auxTypes{
					opWrappers:      tC.existing,
					methodToWrapper: map[*descriptorpb.MethodDescriptorProto]operationWrapper{},
				},
			}

			err := g.maybeAddOperationWrapper(m)
			if tC.wantErrStr == "" && err != nil {
				t.Errorf("maybeAddOperationWrapper(%+v): unexpected error %v", m, err)
			} else if tC.wantErrStr != "" && err == nil {
				t.Errorf("maybeAddOperationWrapper(%+v): got no error, want %q, with existing %+v", m, tC.wantErrStr, tC.existing)
			} else if tC.wantErrStr != "" && err != nil && !strings.Contains(err.Error(), tC.wantErrStr) {
				t.Errorf("maybeAddOperationWrapper(%+v): got %q, want %q, with existing %+v", m, err.Error(), tC.wantErrStr, tC.existing)
			} else if diff := cmp.Diff(g.aux.opWrappers, tC.want, cmp.Comparer(proto.Equal), cmp.AllowUnexported(operationWrapper{})); diff != "" {
				t.Errorf("maybeAddOperationWrapper(%+v): got(-),want(+)\n%s", m, diff)
			}
		})
	}
}

func TestGenOperations(t *testing.T) {
	wrappers := map[string]operationWrapper{
		"CreateFooOperation": {
			name: "CreateFooOperation",
			response: &descriptorpb.DescriptorProto{
				Name: proto.String("Foo"),
			},
			metadata: &descriptorpb.DescriptorProto{
				Name: proto.String("FooMetadata"),
			},
			responseName: protoreflect.FullName("google.example.v1.Foo"),
			metadataName: protoreflect.FullName("google.example.v1.FooMetadata"),
		},
		"DeleteFooOperation": {
			name: "DeleteFooOperation",
			response: &descriptorpb.DescriptorProto{
				Name: proto.String("Empty"),
			},
			metadata: &descriptorpb.DescriptorProto{
				Name: proto.String("FooMetadata"),
			},
			responseName: protoreflect.FullName("google.protobuf.Empty"),
			metadataName: protoreflect.FullName("google.example.v1.FooMetadata"),
		},
	}
	file := &descriptorpb.FileDescriptorProto{
		Package: proto.String("google.example.v1"),
		Options: &descriptorpb.FileOptions{
			GoPackage: proto.String("cloud.google.com/go/example/apiv1/examplepb"),
		},
	}
	parentFiles := make(map[protoreflect.ProtoMessage]*descriptorpb.FileDescriptorProto)
	for _, ow := range wrappers {
		if ow.responseName == emptyValue {
			parentFiles[ow.response] = protodesc.ToFileDescriptorProto(emptypb.File_google_protobuf_empty_proto)
		} else {
			parentFiles[ow.response] = file
		}
		parentFiles[ow.metadata] = file
	}
	g := &generator{
		aux: &auxTypes{
			opWrappers: wrappers,
		},
		descInfo: pbinfo.Info{
			ParentFile:    parentFiles,
			ParentElement: make(map[pbinfo.ProtoType]pbinfo.ProtoType),
		},
		imports: make(map[pbinfo.ImportSpec]bool),
		opts:    &options{transports: []transport{grpc, rest}},
	}

	wantImports := map[pbinfo.ImportSpec]bool{
		{Path: "context"}:                         true,
		{Path: "cloud.google.com/go/longrunning"}: true,
		{Name: "examplepb", Path: "cloud.google.com/go/example/apiv1/examplepb"}: true,
		{Name: "gax", Path: "github.com/googleapis/gax-go/v2"}:                   true,
		{Path: "time"}: true,
	}

	if err := g.genOperations(); err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(g.imports, wantImports); diff != "" {
		t.Errorf("imports got(-),want(+):\n%s", diff)
	}

	txtdiff.Diff(t, g.pt.String(), filepath.Join("testdata", "gen_operations.want"))
}

func TestGenIterators(t *testing.T) {
	g := &generator{
		aux: &auxTypes{
			iters: map[string]*iterType{
				"FooIterator": {
					iterTypeName: "FooIterator",
					elemTypeName: "*examplepb.Foo",
					elemImports: []pbinfo.ImportSpec{
						{Name: "examplepb", Path: "cloud.google.com/go/example/apiv1/examplepb"},
					},
				},
			},
		},
		imports: make(map[pbinfo.ImportSpec]bool),
	}

	wantImports := map[pbinfo.ImportSpec]bool{
		{Name: "examplepb", Path: "cloud.google.com/go/example/apiv1/examplepb"}: true,
		{Path: "google.golang.org/api/iterator"}:                                 true,
	}

	if err := g.genIterators(); err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(g.imports, wantImports); diff != "" {
		t.Errorf("imports got(-),want(+):\n%s", diff)
	}

	txtdiff.Diff(t, g.pt.String(), filepath.Join("testdata", "gen_iterators.want"))
}

func TestSortOperationWrapperMap(t *testing.T) {
	in := map[string]operationWrapper{
		"FooOperation":   {name: "FooOperation"},
		"ZzzzzOperation": {name: "ZzzzzOperation"},
		"BarOperation":   {name: "BarOperation"},
		"AaaaaOperation": {name: "AaaaaOperation"},
	}
	want := []operationWrapper{
		{name: "AaaaaOperation"},
		{name: "BarOperation"},
		{name: "FooOperation"},
		{name: "ZzzzzOperation"},
	}
	got := sortOperationWrapperMap(in)
	if diff := cmp.Diff(got, want, cmp.AllowUnexported(operationWrapper{})); diff != "" {
		t.Errorf("sortOperationWrapperMap: got(-),want(+):\n%s", diff)
	}
}

func TestSortIteratorMap(t *testing.T) {
	in := map[string]*iterType{
		"FooIterator":   {iterTypeName: "FooIterator"},
		"ZzzzzIterator": {iterTypeName: "ZzzzzIterator"},
		"BarIterator":   {iterTypeName: "BarIterator"},
		"AaaaaIterator": {iterTypeName: "AaaaaIterator"},
	}
	want := []*iterType{
		{iterTypeName: "AaaaaIterator"},
		{iterTypeName: "BarIterator"},
		{iterTypeName: "FooIterator"},
		{iterTypeName: "ZzzzzIterator"},
	}
	got := sortIteratorMap(in)
	if diff := cmp.Diff(got, want, cmp.AllowUnexported(iterType{})); diff != "" {
		t.Errorf("sortIteratorMap: got(-),want(+):\n%s", diff)
	}
}
