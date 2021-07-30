// Copyright 2021 Google LLC
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

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/google/go-cmp/cmp"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"github.com/googleapis/gapic-generator-go/internal/txtdiff"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoiface"
)

func TestCustomOpProtoName(t *testing.T) {
	pkg := "google.cloud.foo.v1"
	op := &descriptor.DescriptorProto{
		Name: proto.String("Operation"),
	}
	g := &generator{
		aux: &auxTypes{
			customOp: &customOp{
				message: op,
			},
		},
		descInfo: pbinfo.Info{
			ParentFile: map[protoiface.MessageV1]*descriptor.FileDescriptorProto{
				op: {
					Package: proto.String(pkg),
				},
			},
		},
	}
	got := g.customOpProtoName()
	want := fmt.Sprintf(".%s.%s", pkg, op.GetName())
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("got(-),want(+):\n%s", diff)
	}
}

func TestCustomPointerTyp(t *testing.T) {
	op := &descriptor.DescriptorProto{
		Name: proto.String("Operation"),
	}
	g := &generator{
		aux: &auxTypes{
			customOp: &customOp{
				message: op,
			},
		},
		descInfo: pbinfo.Info{
			ParentFile: map[protoiface.MessageV1]*descriptor.FileDescriptorProto{
				op: {
					Package: proto.String("google.cloud.foo.v1"),
					Options: &descriptor.FileOptions{
						GoPackage: proto.String("google.golang.org/genproto/cloud/foo/v1;foo"),
					},
				},
			},
		},
	}
	got, err := g.customOpPointerType()
	if err != nil {
		t.Fatal(err)
	}
	want := "*foopb.Operation"
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("got(-),want(+):\n%s", diff)
	}
}

func TestCustomOpInit(t *testing.T) {
	op := &descriptor.DescriptorProto{
		Name: proto.String("Operation"),
	}
	g := &generator{
		aux: &auxTypes{
			customOp: &customOp{
				message: op,
			},
		},
	}
	got := g.customOpInit("foo")
	want := "&Operation{proto: foo}"
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("got(-),want(+):\n%s", diff)
	}
}

func TestCustomOperationType(t *testing.T) {
	op := &descriptor.DescriptorProto{
		Name: proto.String("Operation"),
	}
	g := &generator{
		aux: &auxTypes{
			customOp: &customOp{
				message: op,
			},
		},
		descInfo: pbinfo.Info{
			ParentFile: map[protoiface.MessageV1]*descriptor.FileDescriptorProto{
				op: {
					Package: proto.String("google.cloud.foo.v1"),
					Options: &descriptor.FileOptions{
						GoPackage: proto.String("google.golang.org/genproto/cloud/foo/v1;foo"),
					},
				},
			},
		},
	}
	err := g.customOperationType()
	if err != nil {
		t.Fatal(err)
	}
	txtdiff.Diff(t, "custom_op_type", g.pt.String(), filepath.Join("testdata", "custom_op_type.want"))
}
