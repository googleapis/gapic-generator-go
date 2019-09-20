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
	"strings"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
)

func typep(t descriptor.FieldDescriptorProto_Type) *descriptor.FieldDescriptorProto_Type {
	return &t
}

func TestTree(t *testing.T) {
	fieldVals := [][]string{
		{"a.b", "1"},
		{"a.c", "xyz"},
		{"a.d", "2.718281828"},
		{"x", "true"},
	}

	info := pbinfo.Info{
		Type: map[string]pbinfo.ProtoType{
			"RootType": &descriptor.DescriptorProto{
				Field: []*descriptor.FieldDescriptorProto{
					{Name: proto.String("a"), TypeName: proto.String("AType")},
					{Name: proto.String("x"), Type: typep(descriptor.FieldDescriptorProto_TYPE_BOOL)},
				},
			},
			"AType": &descriptor.DescriptorProto{
				Field: []*descriptor.FieldDescriptorProto{
					{Name: proto.String("b"), Type: typep(descriptor.FieldDescriptorProto_TYPE_INT64)},
					{Name: proto.String("c"), Type: typep(descriptor.FieldDescriptorProto_TYPE_STRING)},
					{Name: proto.String("d"), Type: typep(descriptor.FieldDescriptorProto_TYPE_DOUBLE)},
				},
			},
		},
	}

	root := initTree{
		typ: initType{desc: info.Type["RootType"]},
	}
	for _, fv := range fieldVals {
		if err := root.parseInit(fv[0], fv[1], "", info); err != nil {
			t.Fatal(err)
		}
	}

	for _, tst := range []struct {
		path []string
		val  string
	}{
		{[]string{"a", "b"}, "1"},
		{[]string{"a", "c"}, `"xyz"`},
		{[]string{"a", "d"}, "2.718281828"},
		{[]string{"x"}, "true"},
	} {
		node := &root
		for _, p := range tst.path {
			var err error
			node, err = node.get(p, info)
			if err != nil {
				t.Error(err)
			}
		}
		if node.leafVal != tst.val {
			t.Errorf("%s = %q, want %q", strings.Join(tst.path, "->"), node.leafVal, tst.val)
		}
	}
}

func TestTreeErrors(t *testing.T) {
	info := pbinfo.Info{
		Type: map[string]pbinfo.ProtoType{
			"RootType": &descriptor.DescriptorProto{
				Field: []*descriptor.FieldDescriptorProto{
					{Name: proto.String("a"), Type: typep(descriptor.FieldDescriptorProto_TYPE_INT64)},
				},
			},
		},
	}

testcase:
	for _, tst := range [][][]string{
		{{"3", "4"}}, // bad field name

		{{"a", "1"}, {"a", "2"}}, // sets same node twice
		{{"a", "abc"}},           // type is int64, value is string
		{{"unknown_field", "3"}}, // field doesn't exist
	} {
		root := initTree{typ: initType{desc: info.Type["RootType"]}}
		for _, txt := range tst {
			if root.parseInit(txt[0], txt[1], "", info) != nil {
				continue testcase
			}
		}
		t.Errorf("initTree.parse succeeded, want error: %v", tst)
	}
}
