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

	"github.com/google/go-cmp/cmp"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/genproto/googleapis/api/serviceconfig"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoiface"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/apipb"
)

func TestCollectMixins(t *testing.T) {
	operationsHTTP := &annotations.HttpRule{
		Selector: "google.longrunning.Operations.GetOperation",
		Pattern: &annotations.HttpRule_Get{
			Get: "/v1/{operation=projects/*/locations/*/operations/*}",
		},
	}
	locationHTTP := &annotations.HttpRule{
		Selector: "google.cloud.location.Locations.GetLocation",
		Pattern: &annotations.HttpRule_Get{
			Get: "/v1/{location=projects/*/locations/*}",
		},
	}
	iamHTTP := &annotations.HttpRule{
		Selector: "google.iam.v1.IAMPolicy.GetIamPolicy",
		Pattern: &annotations.HttpRule_Get{
			Get: "/v1/{resource=projects/*/locations/*/foos/*}",
		},
	}
	iamDescription := "Gets the access control policy for a resource. Returns an empty policy if the resource exists and does not have a policy set."
	g := generator{
		comments: make(map[protoiface.MessageV1]string),
		mixins:   make(mixins),
		serviceConfig: &serviceconfig.Service{
			Apis: []*apipb.Api{
				{Name: "google.example.library.v1.Library"},
				{Name: "google.longrunning.Operations"},
				{Name: "google.cloud.location.Locations"},
				{Name: "google.iam.v1.IAMPolicy"},
			},
			Http: &annotations.Http{
				Rules: []*annotations.HttpRule{
					operationsHTTP,
					locationHTTP,
					iamHTTP,
				},
			},
			Documentation: &serviceconfig.Documentation{
				Rules: []*serviceconfig.DocumentationRule{
					{
						Selector:    "google.iam.v1.IAMPolicy.GetIamPolicy",
						Description: iamDescription,
					},
				},
			},
		},
	}

	g.collectMixins()

	for _, want := range []struct {
		api, comment string
		len          int
		ext          *annotations.HttpRule
	}{
		{
			api:     "google.longrunning.Operations",
			comment: "is a utility method from google.longrunning.Operations.",
			len:     1,
			ext:     operationsHTTP,
		},
		{
			api:     "google.cloud.location.Locations",
			comment: "is a utility method from google.cloud.location.Locations.",
			len:     1,
			ext:     locationHTTP,
		},
		{
			api:     "google.iam.v1.IAMPolicy",
			comment: iamDescription,
			len:     1,
			ext:     iamHTTP,
		},
	} {
		if got := len(g.mixins[want.api]); got != want.len {
			t.Errorf("TestCollectMixins(%q) got %d method(s), want %d method(s)\n", want.api, got, want.len)
		} else if got := proto.GetExtension(g.mixins[want.api][0].GetOptions(), annotations.E_Http); !cmp.Equal(got, want.ext, cmp.Comparer(proto.Equal)) {
			t.Errorf("TestCollectMixins(%q) got %v, want %v\n", want.api, got, want.ext)
		} else if diff := cmp.Diff(g.comments[g.mixins[want.api][0]], want.comment); diff != "" {
			t.Errorf("TestCollectMixins(%q) got(-),want(+):\n%s", want.api, diff)
		}
	}
}

func TestGetMixinFiles(t *testing.T) {
	g := generator{
		mixins: mixins{
			"google.longrunning.Operations":   operationsMethods(),
			"google.cloud.location.Locations": locationMethods(),
			"google.iam.v1.IAMPolicy":         iamPolicyMethods(),
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
		mixins: mixins{
			"google.longrunning.Operations":   operationsMethods(),
			"google.cloud.location.Locations": locationMethods(),
		},
		serviceConfig: &serviceconfig.Service{
			Apis: []*apipb.Api{
				{Name: "foo.bar.Baz"},
				{Name: "google.iam.v1.IAMPolicy"},
			},
		},
	}

	var want bool
	if got := g.hasIAMPolicyMixin(); !cmp.Equal(got, want) {
		t.Errorf("TestHasIAMPolicyMixin wanted %v but got %v", want, got)
	}

	want = true
	g.mixins["google.iam.v1.IAMPolicy"] = iamPolicyMethods()
	if got := g.hasIAMPolicyMixin(); !cmp.Equal(got, want) {
		t.Errorf("TestHasIAMPolicyMixin wanted %v but got %v", want, got)
	}

	want = false
	g.serviceConfig.Apis = []*apipb.Api{{Name: "google.iam.v1.IAMPolicy"}}
	if got := g.hasIAMPolicyMixin(); !cmp.Equal(got, want) {
		t.Errorf("TestHasIAMPolicyMixin wanted %v but got %v", want, got)
	}

	g.hasIAMPolicyOverrides = true
	g.serviceConfig.Apis = append(g.serviceConfig.Apis, &apipb.Api{Name: "foo.bar.Baz"})
	if got := g.hasIAMPolicyMixin(); !cmp.Equal(got, want) {
		t.Errorf("TestHasIAMPolicyMixin wanted %v but got %v", want, got)
	}
}

func TestCheckIAMPolicyOverrides(t *testing.T) {
	g := &generator{
		mixins: make(mixins),
	}
	serv := &descriptorpb.ServiceDescriptorProto{
		Method: []*descriptorpb.MethodDescriptorProto{
			{Name: proto.String("ListFoos")},
			{Name: proto.String("GetFoo")},
		},
	}
	other := &descriptorpb.ServiceDescriptorProto{
		Method: []*descriptorpb.MethodDescriptorProto{
			{Name: proto.String("ListBars")},
			{Name: proto.String("GetBar")},
		},
	}
	servs := []*descriptorpb.ServiceDescriptorProto{serv, other}
	var want bool
	g.checkIAMPolicyOverrides(servs)
	if got := g.hasIAMPolicyOverrides; !cmp.Equal(got, want) {
		t.Errorf("TestCheckIAMPolicyOverrides wanted %v but got %v", want, got)
	}

	want = true
	g.mixins["google.iam.v1.IAMPolicy"] = iamPolicyMethods()
	serv.Method = append(serv.Method, &descriptorpb.MethodDescriptorProto{Name: proto.String("GetIamPolicy")})
	g.checkIAMPolicyOverrides(servs)
	if got := g.hasIAMPolicyOverrides; !cmp.Equal(got, want) {
		t.Errorf("TestCheckIAMPolicyOverrides wanted %v but got %v", want, got)
	}
}

func TestHasLocationMixin(t *testing.T) {
	g := generator{
		mixins: mixins{
			"google.longrunning.Operations": operationsMethods(),
			"google.iam.v1.IAMPolicy":       iamPolicyMethods(),
		},
		serviceConfig: &serviceconfig.Service{
			Apis: []*apipb.Api{
				{Name: "foo.bar.Baz"},
				{Name: "google.cloud.location.Locations"},
			},
		},
	}

	var want bool
	if got := g.hasLocationMixin(); !cmp.Equal(got, want) {
		t.Errorf("TestHasLocationMixin wanted %v but got %v", want, got)
	}

	want = true
	g.mixins["google.cloud.location.Locations"] = locationMethods()
	if got := g.hasLocationMixin(); !cmp.Equal(got, want) {
		t.Errorf("TestHasLocationMixin wanted %v but got %v", want, got)
	}

	want = false
	g.serviceConfig.Apis = []*apipb.Api{{Name: "google.cloud.location.Locations"}}
	if got := g.hasLocationMixin(); !cmp.Equal(got, want) {
		t.Errorf("TestHasLocationMixin wanted %v but got %v", want, got)
	}
}

func TestHasLROMixin(t *testing.T) {
	g := generator{
		mixins: mixins{
			"google.cloud.location.Locations": locationMethods(),
			"google.iam.v1.IAMPolicy":         iamPolicyMethods(),
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
	g.mixins["google.longrunning.Operations"] = operationsMethods()
	if got := g.hasLROMixin(); !cmp.Equal(got, want) {
		t.Errorf("TestHasLROMixin wanted %v but got %v", want, got)
	}

	want = false
	g.serviceConfig.Apis = []*apipb.Api{{Name: "google.longrunning.Operations"}}
	if got := g.hasLROMixin(); !cmp.Equal(got, want) {
		t.Errorf("TestHasLROMixin wanted %v but got %v", want, got)
	}
}

func TestGetOperationPathOverride(t *testing.T) {
	for _, tc := range []struct {
		pkg, want string
		http      *annotations.Http
	}{
		{
			pkg:  "google.example.library.v1beta2",
			want: "/v1beta2/%s",
		},
		{
			pkg:  "google.example.library.v1",
			want: "/v1/%s",
		},
		{
			pkg:  "google.example.library.v3alpha1p1",
			want: "/v3alpha1p1/%s",
		},
		{
			pkg:  "google.example.library.v2",
			want: "/v2/%s",
			http: &annotations.Http{
				Rules: []*annotations.HttpRule{
					&annotations.HttpRule{
						Selector: "google.longrunning.Operations.GetOperation",
						Pattern: &annotations.HttpRule_Get{
							Get: "/v2/{operation=projects/*/locations/*/operations/*}",
						},
					},
				},
			},
		},
	} {
		initMixinFiles()
		g := generator{
			comments: make(map[protoiface.MessageV1]string),
			mixins:   make(mixins),
			serviceConfig: &serviceconfig.Service{
				Http: tc.http,
			},
		}
		if got := g.getOperationPathOverride(tc.pkg); !cmp.Equal(got, tc.want) {
			t.Errorf("TestGetOperationPathOverrideMissing wanted %v but got %v", tc.want, got)
		}
	}
}

func TestIncludeMixinInputFile(t *testing.T) {
	for _, tst := range []struct {
		name, file, pkgPath string
		want                bool
	}{
		{
			name:    "non-mixin",
			file:    "google/cloud/library/v1/library.proto",
			pkgPath: "cloud.google.com/go/library/apiv1",
			want:    true,
		},
		{
			name:    "exclude-iam",
			file:    "google/iam/v1/iam_policy.proto",
			pkgPath: "cloud.google.com/go/library/apiv1",
			want:    false,
		},
		{
			name:    "exclude-locations",
			file:    "google/cloud/location/locations.proto",
			pkgPath: "cloud.google.com/go/library/apiv1",
			want:    false,
		},
		{
			name:    "exclude-longrunning",
			file:    "google/longrunning/operations.proto",
			pkgPath: "cloud.google.com/go/library/apiv1",
			want:    false,
		},
		{
			name:    "include-iam",
			file:    "google/iam/v1/iam_policy.proto",
			pkgPath: "cloud.google.com/go/iam/apiv1",
			want:    true,
		},
		{
			name:    "include-locations",
			file:    "google/cloud/location/locations.proto",
			pkgPath: "cloud.google.com/go/locations/apiv1",
			want:    true,
		},
		{
			name:    "include-longrunning",
			file:    "google/longrunning/operations.proto",
			pkgPath: "cloud.google.com/go/longrunning/apiv1",
			want:    true,
		},
	} {
		g := &generator{opts: &options{pkgPath: tst.pkgPath}}
		if got := g.includeMixinInputFile(tst.file); got != tst.want {
			t.Errorf("%s: got %v, want %v", tst.name, got, tst.want)
		}
	}
}

// locationMethods is just used for testing.
func locationMethods() []*descriptorpb.MethodDescriptorProto {
	return mixinFiles["google.cloud.location.Locations"][0].GetService()[0].GetMethod()
}

// iamPolicyMethods is just used for testing.
func iamPolicyMethods() []*descriptorpb.MethodDescriptorProto {
	return mixinFiles["google.iam.v1.IAMPolicy"][0].GetService()[0].GetMethod()
}

// operationsMethods is just used for testing.
func operationsMethods() []*descriptorpb.MethodDescriptorProto {
	return mixinFiles["google.longrunning.Operations"][0].GetService()[0].GetMethod()
}
