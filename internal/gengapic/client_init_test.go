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
	"bytes"
	"path/filepath"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/golang/protobuf/ptypes/duration"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	"github.com/google/go-cmp/cmp"
	conf "github.com/googleapis/gapic-generator-go/internal/grpc_service_config"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"github.com/googleapis/gapic-generator-go/internal/txtdiff"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/genproto/googleapis/api/serviceconfig"
	code "google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/apipb"
)

func TestClientHook(t *testing.T) {
	var g generator

	g.clientHook("Foo")
	got := g.pt.String()
	want := "var newFooClientHook clientHook\n\n"

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("clientHook() (-got,+want): %s", diff)
	}
}

func TestClientOpt(t *testing.T) {
	var g generator
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
	g.opts = &options{transports: []transport{grpc}}
	cpb := &conf.ServiceConfig{
		MethodConfig: []*conf.MethodConfig{
			{
				Name: []*conf.MethodConfig_Name{
					{
						Service: "bar.FooService",
						Method:  "Zip",
					},
				},
				MaxRequestMessageBytes:  &wrappers.UInt32Value{Value: 123456},
				MaxResponseMessageBytes: &wrappers.UInt32Value{Value: 123456},
				RetryOrHedgingPolicy: &conf.MethodConfig_RetryPolicy_{
					RetryPolicy: &conf.MethodConfig_RetryPolicy{
						InitialBackoff:    &duration.Duration{Nanos: 100000000},
						MaxBackoff:        &duration.Duration{Seconds: 60},
						BackoffMultiplier: 1.3,
						RetryableStatusCodes: []code.Code{
							code.Code_UNKNOWN,
						},
					},
				},
			},
			{
				Name: []*conf.MethodConfig_Name{
					{
						Service: "bar.FooService",
					},
					{
						Service: "bar.ServHostPort",
					},
					{
						Service: "bar.ServIamOverride",
					},
				},
				MaxRequestMessageBytes:  &wrappers.UInt32Value{Value: 654321},
				MaxResponseMessageBytes: &wrappers.UInt32Value{Value: 654321},
				RetryOrHedgingPolicy: &conf.MethodConfig_RetryPolicy_{
					RetryPolicy: &conf.MethodConfig_RetryPolicy{
						InitialBackoff:    &duration.Duration{Nanos: 10000000},
						MaxBackoff:        &duration.Duration{Seconds: 7},
						BackoffMultiplier: 1.1,
						RetryableStatusCodes: []code.Code{
							code.Code_UNKNOWN,
						},
					},
				},
			},
		},
	}
	data, err := protojson.Marshal(cpb)
	if err != nil {
		t.Error(err)
	}
	in := bytes.NewReader(data)
	g.grpcConf, err = conf.New(in)
	if err != nil {
		t.Error(err)
	}

	serv := &descriptor.ServiceDescriptorProto{
		Name: proto.String("FooService"),
		Method: []*descriptor.MethodDescriptorProto{
			{Name: proto.String("Zip"), Options: &descriptor.MethodOptions{}},
			{Name: proto.String("Zap"), Options: &descriptor.MethodOptions{}},
			{Name: proto.String("Smack")},
		},
		Options: &descriptor.ServiceOptions{},
	}
	if err := proto.SetExtension(serv.Options, annotations.E_DefaultHost, proto.String("foo.googleapis.com")); err != nil {
		t.Fatal(err)
	}

	// Test some annotations
	if err := proto.SetExtension(serv.Method[0].Options, annotations.E_Http, &annotations.HttpRule{
		Pattern: &annotations.HttpRule_Get{
			Get: "/zip",
		},
	}); err != nil {
		t.Fatal(err)
	}
	if err := proto.SetExtension(serv.Method[1].Options, annotations.E_Http, &annotations.HttpRule{
		Pattern: &annotations.HttpRule_Post{
			Post: "/zap",
		},
	}); err != nil {
		t.Fatal(err)
	}

	servHostPort := &descriptor.ServiceDescriptorProto{
		Name: proto.String("ServHostPort"),
		Method: []*descriptor.MethodDescriptorProto{
			{Name: proto.String("Smack")},
		},
		Options: &descriptor.ServiceOptions{},
	}
	if err := proto.SetExtension(servHostPort.Options, annotations.E_DefaultHost, proto.String("foo.googleapis.com:1234")); err != nil {
		t.Fatal(err)
	}

	servIAMOverride := &descriptor.ServiceDescriptorProto{
		Name: proto.String("ServIamOverride"),
		Method: []*descriptor.MethodDescriptorProto{
			{Name: proto.String("GetIamPolicy")},
			{Name: proto.String("SetIamPolicy")},
			{Name: proto.String("TestIamPermissions")},
		},
		Options: &descriptor.ServiceOptions{},
	}
	if err := proto.SetExtension(servIAMOverride.Options, annotations.E_DefaultHost, proto.String("foo.googleapis.com:1234")); err != nil {
		t.Fatal(err)
	}

	f := &descriptor.FileDescriptorProto{
		Package: proto.String("bar"),
		Service: []*descriptor.ServiceDescriptorProto{serv, servHostPort, servIAMOverride},
	}
	files := append(g.getMixinFiles(), f)
	g.descInfo = pbinfo.Of(files)

	for _, tst := range []struct {
		tstName, servName string
		serv              *descriptor.ServiceDescriptorProto
		hasOverride       bool
	}{
		{tstName: "foo_opt", servName: "Foo", serv: serv},
		{tstName: "empty_opt", servName: "", serv: serv},
		{tstName: "host_port_opt", servName: "Bar", serv: servHostPort},
		{tstName: "iam_override_opt", servName: "Baz", serv: servIAMOverride, hasOverride: true},
	} {
		g.reset()
		g.hasIAMPolicyOverrides = tst.hasOverride
		if err := g.clientOptions(tst.serv, tst.servName); err != nil {
			t.Error(err)
			continue
		}
		txtdiff.Diff(t, tst.tstName, g.pt.String(), filepath.Join("testdata", tst.tstName+".want"))
	}
}

func TestServiceDoc(t *testing.T) {
	serv := &descriptor.ServiceDescriptorProto{
		Name: proto.String("MyService"),
	}

	var g generator
	g.comments = make(map[proto.Message]string)

	for _, tst := range []struct {
		in, want   string
		deprecated bool
	}{
		{
			in:   "",
			want: "",
		},
		{
			in:   "Does stuff.\n It also does other stuffs.",
			want: "//\n// Does stuff.\n// It also does other stuffs.\n",
		},
		{
			in:         "This is deprecated.\n It does not have a proper comment.",
			want:       "//\n// This is deprecated.\n// It does not have a proper comment.\n//\n// Deprecated: MyService may be removed in a future version.\n",
			deprecated: true,
		},
		{
			in:         "Deprecated: this is a proper deprecation notice.",
			want:       "//\n// MyService is deprecated.\n//\n// Deprecated: this is a proper deprecation notice.\n",
			deprecated: true,
		},
		{
			in:         "Does my thing.\nDeprecated: this is a proper deprecation notice.",
			want:       "//\n// Does my thing.\n// Deprecated: this is a proper deprecation notice.\n",
			deprecated: true,
		},
		{
			in:         "",
			want:       "//\n// MyService is deprecated.\n//\n// Deprecated: MyService may be removed in a future version.\n",
			deprecated: true,
		},
	} {
		g.comments[serv] = tst.in
		serv.Options = &descriptor.ServiceOptions{
			Deprecated: proto.Bool(tst.deprecated),
		}
		g.pt.Reset()
		g.serviceDoc(serv)
		if diff := cmp.Diff(g.pt.String(), tst.want); diff != "" {
			t.Errorf("comment() got(-),want(+):\n%s", diff)
		}
	}
}

func TestClientInit(t *testing.T) {
	var g generator
	g.apiName = "Awesome Foo API"
	g.imports = map[pbinfo.ImportSpec]bool{}
	g.opts = &options{transports: []transport{grpc}}

	servPlain := &descriptor.ServiceDescriptorProto{
		Name: proto.String("Foo"),
		Method: []*descriptor.MethodDescriptorProto{
			{
				Name:       proto.String("Zip"),
				InputType:  proto.String(".mypackage.Bar"),
				OutputType: proto.String(".mypackage.Foo"),
				Options:    &descriptor.MethodOptions{},
			},
		},
	}
	servLRO := &descriptor.ServiceDescriptorProto{
		Name: proto.String("Foo"),
		Method: []*descriptor.MethodDescriptorProto{
			{
				Name:       proto.String("Zip"),
				InputType:  proto.String(".mypackage.Bar"),
				OutputType: proto.String(".google.longrunning.Operation"),
				Options:    &descriptor.MethodOptions{},
			},
		},
	}
	servDeprecated := &descriptor.ServiceDescriptorProto{
		Name: proto.String("Foo"),
		Method: []*descriptor.MethodDescriptorProto{
			{
				Name:       proto.String("Zip"),
				InputType:  proto.String(".mypackage.Bar"),
				OutputType: proto.String(".mypackage.Foo"),
				Options:    &descriptor.MethodOptions{},
			},
		},
		Options: &descriptor.ServiceOptions{
			Deprecated: proto.Bool(true),
		},
	}
	for _, s := range []*descriptor.ServiceDescriptorProto{servPlain, servDeprecated, servLRO} {
		proto.SetExtension(s.Method[0].Options, annotations.E_Http, &annotations.HttpRule{
			Pattern: &annotations.HttpRule_Get{
				Get: "/zip",
			},
		})
	}

	for _, tst := range []struct {
		tstName   string
		servName  string
		mixins    mixins
		serv      *descriptor.ServiceDescriptorProto
		parameter *string
		httpRule  *annotations.HttpRule
	}{
		{
			tstName: "foo_client_init",
			mixins: mixins{
				"google.cloud.location.Locations": locationMethods(),
				"google.iam.v1.IAMPolicy":         iamPolicyMethods(),
			},
			servName:  "Foo",
			serv:      servPlain,
			parameter: proto.String("go-gapic-package=path;mypackage"),
		},
		{
			tstName: "foo_rest_client_init",
			mixins: mixins{
				"google.cloud.location.Locations": locationMethods(),
				"google.iam.v1.IAMPolicy":         iamPolicyMethods(),
			},
			servName:  "Foo",
			serv:      servPlain,
			parameter: proto.String("go-gapic-package=path;mypackage,transport=rest"),
		},
		{
			tstName:   "empty_client_init",
			servName:  "",
			serv:      servPlain,
			parameter: proto.String("go-gapic-package=path;mypackage,transport=grpc+rest"),
		},
		{
			tstName: "lro_client_init",
			mixins: mixins{
				"google.longrunning.Operations": operationsMethods(),
			},
			servName:  "Foo",
			serv:      servLRO,
			parameter: proto.String("go-gapic-package=path;mypackage"),
		},
		{
			tstName:   "deprecated_client_init",
			servName:  "",
			serv:      servDeprecated,
			parameter: proto.String("go-gapic-package=path;mypackage,transport=grpc+rest"),
		},
	} {
		request := plugin.CodeGeneratorRequest{
			Parameter: tst.parameter,
			ProtoFile: []*descriptor.FileDescriptorProto{
				{
					Package: proto.String("mypackage"),
					Options: &descriptor.FileOptions{
						GoPackage: proto.String("mypackage"),
					},
					Service: []*descriptor.ServiceDescriptorProto{tst.serv},
					MessageType: []*descriptor.DescriptorProto{
						{
							Name: proto.String("Bar"),
						},
						{
							Name: proto.String("Foo"),
						},
					},
				},
				{
					Package: proto.String("google.longrunning"),
					Options: &descriptor.FileOptions{
						GoPackage: proto.String("google.golang.org/genproto/googleapis/longrunning;longrunning"),
					},
					MessageType: []*descriptor.DescriptorProto{
						{
							Name: proto.String("Operation"),
						},
						{
							Name: proto.String("GetOperationRequest"),
						},
						{
							Name: proto.String("DeleteOperationRequest"),
						},
						{
							Name: proto.String("WaitOperationRequest"),
						},
						{
							Name: proto.String("ListOperationsRequest"),
						},
						{
							Name: proto.String("ListOperationsResponse"),
						},
						{
							Name: proto.String("CancelOperationRequest"),
						},
					},
				},
				{
					Package: proto.String("google.cloud.location"),
					Options: &descriptor.FileOptions{
						GoPackage: proto.String("google.golang.org/genproto/googleapis/cloud/location;location"),
					},
					MessageType: []*descriptor.DescriptorProto{
						{
							Name: proto.String("ListLocationsRequest"),
						},
						{
							Name: proto.String("ListLocationsResponse"),
						},
						{
							Name: proto.String("Location"),
						},
						{
							Name: proto.String("GetLocationRequest"),
						},
						{
							Name: proto.String("GetLocationResponse"),
						},
					},
				},
				{
					Package: proto.String("google.iam.v1"),
					Options: &descriptor.FileOptions{
						GoPackage: proto.String("google.golang.org/genproto/googleapis/iam/v1;iam"),
					},
					MessageType: []*descriptor.DescriptorProto{
						{
							Name: proto.String("Policy"),
						},
						{
							Name: proto.String("TestIamPermissionsRequest"),
						},
						{
							Name: proto.String("TestIamPermissionsResponse"),
						},
						{
							Name: proto.String("SetIamPolicyRequest"),
						},
						{
							Name: proto.String("SetIamPolicyResponse"),
						},
						{
							Name: proto.String("GetIamPolicyRequest"),
						},
					},
				},
			}}
		g.init(&request)
		g.comments = map[proto.Message]string{
			tst.serv: "Foo service does stuff.",
		}
		g.mixins = tst.mixins
		g.serviceConfig = &serviceconfig.Service{
			Apis: []*apipb.Api{
				{Name: "foo.bar.Baz"},
				{Name: "google.iam.v1.IAMPolicy"},
				{Name: "google.cloud.location.Locations"},
				{Name: "google.longrunning.Operations"},
			},
		}

		g.reset()
		g.makeClients(tst.serv, tst.servName)

		txtdiff.Diff(t, tst.tstName, g.pt.String(), filepath.Join("testdata", tst.tstName+".want"))
	}
}

func TestGenerateDefaultAudience(t *testing.T) {
	tests := []struct {
		name string
		host string
		want string
	}{
		{name: "plain host", host: "foo.googleapis.com", want: "https://foo.googleapis.com/"},
		{name: "host with port", host: "foo.googleapis.com:443", want: "https://foo.googleapis.com/"},
		{name: "host with scheme", host: "https://foo.googleapis.com", want: "https://foo.googleapis.com/"},
		{name: "host with scheme and port", host: "https://foo.googleapis.com:1234", want: "https://foo.googleapis.com/"},
		{name: "host is a proper audience", host: "https://foo.googleapis.com/", want: "https://foo.googleapis.com/"},
		{name: "host with non-http scheme", host: "ftp://foo.googleapis.com", want: "ftp://foo.googleapis.com/"},
		{name: "host with path", host: "foo.googleapis.com:443/extra/path", want: "https://foo.googleapis.com/"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := generateDefaultAudience(tc.host); got != tc.want {
				t.Errorf("generateDefaultAudience(%q) = %q, want %q", tc.host, got, tc.want)
			}
		})
	}
}
