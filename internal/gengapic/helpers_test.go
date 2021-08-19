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
	"testing"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"
)

func TestCamelToSnake(t *testing.T) {
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

func TestSnakeToCamel(t *testing.T) {
	for _, tst := range []struct {
		in, want string
	}{
		// Note: we cannot determine if a sequence of lower case letters
		// represent an acronym, so they are treated like discrete words.
		{"iam_credentials", "IamCredentials"},
		{"dlp", "Dlp"},
		{"os_config", "OsConfig"},
		{"display_video_360_advertiser_links", "DisplayVideo_360AdvertiserLinks"},
	} {
		got := snakeToCamel(tst.in)
		if diff := cmp.Diff(got, tst.want); diff != "" {
			t.Errorf("snakeToCamel(%q) = got(-),want(+):\n%s", tst.in, diff)
		}
	}
}

func TestLowerFirst(t *testing.T) {
	for _, tst := range []struct {
		in, want string
	}{
		{"Foo", "foo"},
		{"", ""},
		{"BarBaz", "barBaz"},
	} {
		got := lowerFirst(tst.in)
		if diff := cmp.Diff(got, tst.want); diff != "" {
			t.Errorf("Test_lowerFirst: got(-),want(+):\n%s", diff)
		}
	}
}

func TestUpperFirst(t *testing.T) {
	for _, tst := range []struct {
		in, want string
	}{
		{"foo", "Foo"},
		{"", ""},
		{"barBaz", "BarBaz"},
	} {
		got := upperFirst(tst.in)
		if diff := cmp.Diff(got, tst.want); diff != "" {
			t.Errorf("upperFirst(%q): got(-),want(+):\n%s", tst.in, diff)
		}
	}
}

func TestIsOptional(t *testing.T) {
	msg := &descriptor.DescriptorProto{
		Field: []*descriptor.FieldDescriptorProto{
			{
				Name:           proto.String("opt"),
				Proto3Optional: proto.Bool(true),
			},
			{
				Name: proto.String("not_opt"),
			},
		},
	}
	for _, tst := range []struct {
		field string
		want  bool
	}{
		{
			field: "opt",
			want:  true,
		},
		{
			field: "not_opt",
		},
		{
			field: "no_such_field",
		},
	} {
		if got := isOptional(msg, tst.field); got != tst.want {
			t.Errorf("isOptional(%q) = got %v, want %v", tst.field, got, tst.want)
		}
	}
}

func TestStrContains(t *testing.T) {
	set := []string{
		"foo",
		"bar",
		"exists",
		"baz",
	}
	for _, tst := range []struct {
		in   string
		want bool
	}{
		{
			in:   "exists",
			want: true,
		},
		{
			in: "doesn't exist",
		},
	} {
		if got := strContains(set, tst.in); got != tst.want {
			t.Errorf("strContains(%v, %q): got %v, want %v", set, tst.in, got, tst.want)
		}
	}
}

func TestHasMethod(t *testing.T) {
	serv := &descriptor.ServiceDescriptorProto{
		Method: []*descriptor.MethodDescriptorProto{
			{Name: proto.String("ListFoos")},
			{Name: proto.String("GetFoo")},
			{Name: proto.String("CreateFoo")},
		},
	}
	for _, tst := range []struct {
		in   string
		want bool
	}{
		{in: "GetFoo", want: true},
		{in: "DeleteBar"},
	} {
		if got := hasMethod(serv, tst.in); !cmp.Equal(got, tst.want) {
			t.Errorf("TestHasMethod got %v want %v", got, tst.want)
		}
	}
}
