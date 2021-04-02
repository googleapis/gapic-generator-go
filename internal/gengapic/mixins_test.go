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

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/genproto/googleapis/api/serviceconfig"
	"google.golang.org/protobuf/types/known/apipb"
)

func TestCollectMixins(t *testing.T) {
	g := generator{
		mixins: map[string]bool{},
		serviceConfig: &serviceconfig.Service{
			Apis: []*apipb.Api{
				{Name: "google.example.library.v1.Library"},
				{Name: "google.longrunning.Operations"},
				{Name: "google.cloud.location.Locations"},
				{Name: "google.iam.v1.IAMPolicy"},
			},
		},
	}
	want := map[string]bool{
		"google.longrunning.Operations":   true,
		"google.cloud.location.Locations": true,
		"google.iam.v1.IAMPolicy":         true,
	}

	g.collectMixins()
	if diff := cmp.Diff(g.mixins, want); diff != "" {
		t.Errorf("TestCollectMixins got(-),want(+):\n%s", diff)
	}
}

func TestGetMixinFiles(t *testing.T) {
	g := generator{
		mixins: map[string]bool{
			"google.longrunning.Operations":   true,
			"google.cloud.location.Locations": true,
			"google.iam.v1.IAMPolicy":         true,
		},
	}

	// This isn't a great test, but this isn't a sophisticated function.
	want := 5
	if files := g.getMixinFiles(); !cmp.Equal(len(files), want) {
		t.Errorf("TestGetMixinFiles wanted %d mixin proto files but got %d", want, len(files))
	}
}

func TestHasIAMPolicyMixin(t *testing.T) {
	g := generator{
		mixins: map[string]bool{
			"google.longrunning.Operations":   true,
			"google.cloud.location.Locations": true,
		},
	}

	var want bool
	if got := g.hasIAMPolicyMixin(); !cmp.Equal(got, want) {
		t.Errorf("TestHasIAMPolicyMixin wanted %v but got %v", want, got)
	}

	want = true
	g.mixins["google.iam.v1.IAMPolicy"] = true
	if got := g.hasIAMPolicyMixin(); !cmp.Equal(got, want) {
		t.Errorf("TestHasIAMPolicyMixin wanted %v but got %v", want, got)
	}

	want = false
	g.hasIAMPolicyOverrides = true
	if got := g.hasIAMPolicyMixin(); !cmp.Equal(got, want) {
		t.Errorf("TestHasIAMPolicyMixin wanted %v but got %v", want, got)
	}
}

func TestHasIAMPolicyOverrides(t *testing.T) {
	serv := &descriptor.ServiceDescriptorProto{
		Method: []*descriptor.MethodDescriptorProto{
			{Name: proto.String("ListFoos")},
			{Name: proto.String("GetFoo")},
		},
	}
	other := &descriptor.ServiceDescriptorProto{
		Method: []*descriptor.MethodDescriptorProto{
			{Name: proto.String("ListBars")},
			{Name: proto.String("GetBar")},
		},
	}
	servs := []*descriptor.ServiceDescriptorProto{serv, other}
	var want bool
	if got := hasIAMPolicyOverrides(servs); !cmp.Equal(got, want) {
		t.Errorf("TestHasIAMPolicyOverrides wanted %v but got %v", want, got)
	}

	want = true
	serv.Method = append(serv.Method, &descriptor.MethodDescriptorProto{Name: proto.String("GetIamPolicy")})
	if got := hasIAMPolicyOverrides(servs); !cmp.Equal(got, want) {
		t.Errorf("TestHasIAMPolicyOverrides wanted %v but got %v", want, got)
	}
}

func TestHasLocationMixin(t *testing.T) {
	g := generator{
		mixins: map[string]bool{
			"google.longrunning.Operations": true,
			"google.iam.v1.IAMPolicy":       true,
		},
	}

	var want bool
	if got := g.hasLocationMixin(); !cmp.Equal(got, want) {
		t.Errorf("TestHasLocationMixin wanted %v but got %v", want, got)
	}

	want = true
	g.mixins["google.cloud.location.Locations"] = true
	if got := g.hasLocationMixin(); !cmp.Equal(got, want) {
		t.Errorf("TestHasLocationMixin wanted %v but got %v", want, got)
	}
}

func TestHasLROMixin(t *testing.T) {
	g := generator{
		mixins: map[string]bool{
			"google.cloud.location.Locations": true,
			"google.iam.v1.IAMPolicy":         true,
		},
		serviceConfig: &serviceconfig.Service{
			Apis: []*apipb.Api{
				{Name: "foo.bar.Baz"},
				{Name: "google.iam.v1.IAMPolicy"},
				{Name: "google.cloud.location.Locations"},
			},
		},
	}

	var want bool
	if got := g.hasLROMixin(); !cmp.Equal(got, want) {
		t.Errorf("TestHasLROMixin wanted %v but got %v", want, got)
	}

	want = true
	g.mixins["google.longrunning.Operations"] = true
	if got := g.hasLROMixin(); !cmp.Equal(got, want) {
		t.Errorf("TestHasLROMixin wanted %v but got %v", want, got)
	}

	want = false
	g.serviceConfig.Apis = []*apipb.Api{{Name: "google.longrunning.Operations"}}
	if got := g.hasLROMixin(); !cmp.Equal(got, want) {
		t.Errorf("TestHasLROMixin wanted %v but got %v", want, got)
	}
}
