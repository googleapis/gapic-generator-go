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

package main

import (
	"flag"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/google/go-cmp/cmp"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
)

var updateGolden = flag.Bool("update_golden", false, "update golden files")

func diff(t *testing.T, name, got, goldenFile string) {
	t.Helper()

	if *updateGolden {
		if err := ioutil.WriteFile(goldenFile, []byte(got), 0644); err != nil {
			t.Fatal(err)
		}
	}
	want, err := ioutil.ReadFile(goldenFile)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(got, string(want)); diff != "" {
		t.Errorf("%s: (-got,+want)\n%s", name, diff)
	}
}

// TODO(pongad): maybe we should baseline test the whole file.
func TestSample(t *testing.T) {
	inType := &descriptor.DescriptorProto{
		Name: proto.String("InputType"),
		Field: []*descriptor.FieldDescriptorProto{
			{Name: proto.String("a"), TypeName: proto.String("AType")},
			{Name: proto.String("b"), Type: typep(descriptor.FieldDescriptorProto_TYPE_STRING)},
		},
	}
	aType := &descriptor.DescriptorProto{
		Name: proto.String("AType"),
		Field: []*descriptor.FieldDescriptorProto{
			{Name: proto.String("x"), Type: typep(descriptor.FieldDescriptorProto_TYPE_INT64)},
			{Name: proto.String("y"), Type: typep(descriptor.FieldDescriptorProto_TYPE_FLOAT)},
		},
	}

	serv := &descriptor.ServiceDescriptorProto{
		Name: proto.String("FooService"),
		Method: []*descriptor.MethodDescriptorProto{
			{
				Name:      proto.String("MyMethod"),
				InputType: inType.Name,
			},
		},
	}
	file := &descriptor.FileDescriptorProto{
		Options: &descriptor.FileOptions{
			GoPackage: proto.String("path/to/foo;foo"),
		},
	}

	g := generator{
		clientPkg: pbinfo.ImportSpec{Path: "path/to/foo", Name: "foo"},
		imports:   map[pbinfo.ImportSpec]bool{},
		descInfo: pbinfo.Info{
			Serv: map[string]*descriptor.ServiceDescriptorProto{
				".MyService": serv,
			},
			ParentFile: map[proto.Message]*descriptor.FileDescriptorProto{
				serv:           file,
				serv.Method[0]: file,
				inType:         file,
				aType:          file,
			},
			Type: map[string]pbinfo.ProtoType{
				"InputType": inType,
				"AType":     aType,
			},
		},
	}

	vs := SampleValueSet{
		ID: "my_value_set",
		Parameters: SampleParameter{
			Defaults: []string{
				`a.x = 42`,
				`a.y = 3.14159`,
				`b = "foobar"`,
			},
		},
	}
	if err := g.genSample("MyService", "MyMethod", "awesome_region", vs); err != nil {
		t.Fatal(err)
	}
	diff(t, "TestSample", string(g.pt.Bytes()), filepath.Join("testdata", "sample.want"))
}
