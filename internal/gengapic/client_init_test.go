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
	"time"

	"github.com/google/go-cmp/cmp"
	conf "github.com/googleapis/gapic-generator-go/internal/grpc_service_config"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"github.com/googleapis/gapic-generator-go/internal/snippets"
	"github.com/googleapis/gapic-generator-go/internal/txtdiff"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/genproto/googleapis/api/serviceconfig"
	"google.golang.org/genproto/googleapis/cloud/extendedops"
	code "google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoiface"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/apipb"
	duration "google.golang.org/protobuf/types/known/durationpb"
	wrappers "google.golang.org/protobuf/types/known/wrapperspb"
	"google.golang.org/protobuf/types/pluginpb"
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
	cpb := &conf.ServiceConfig{
		MethodConfig: []*conf.MethodConfig{
			{
				Name: []*conf.MethodConfig_Name{
					{
						Service: "bar.FooService",
						Method:  "Zip",
					},
				},
				Timeout:                 duration.New(10 * time.Second),
				MaxRequestMessageBytes:  &wrappers.UInt32Value{Value: 123456},
				MaxResponseMessageBytes: &wrappers.UInt32Value{Value: 123456},
				RetryOrHedgingPolicy: &conf.MethodConfig_RetryPolicy_{
					RetryPolicy: &conf.MethodConfig_RetryPolicy{
						InitialBackoff:    &duration.Duration{Nanos: 100000000},
						MaxBackoff:        &duration.Duration{Seconds: 60},
						BackoffMultiplier: 1.3,
						RetryableStatusCodes: []code.Code{
							code.Code_UNKNOWN,
							code.Code_UNAVAILABLE,
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
				Timeout:                 duration.New(5 * time.Second),
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
		t.Fatal(err)
	}
	in := bytes.NewReader(data)
	grpcConf, err := conf.New(in)
	if err != nil {
		t.Fatal(err)
	}

	g := generator{
		imports: map[pbinfo.ImportSpec]bool{},
		mixins: mixins{
			"google.longrunning.Operations":   operationsMethods(),
			"google.cloud.location.Locations": locationMethods(),
			"google.iam.v1.IAMPolicy":         iamPolicyMethods(),
		},
		serviceConfig: &serviceconfig.Service{
			Apis: []*apipb.Api{
				{Name: "foo.bar.Baz"},
				{Name: "google.iam.v1.IAMPolicy"},
				{Name: "google.cloud.location.Locations"},
				{Name: "google.longrunning.Operations"},
			},
		},
		opts:     &options{transports: []transport{grpc, rest}},
		grpcConf: grpcConf,
	}

	serv := &descriptorpb.ServiceDescriptorProto{
		Name: proto.String("FooService"),
		Method: []*descriptorpb.MethodDescriptorProto{
			{Name: proto.String("Zip"), Options: &descriptorpb.MethodOptions{}},
			{Name: proto.String("Zap"), Options: &descriptorpb.MethodOptions{}},
			{Name: proto.String("Smack")},
		},
		Options: &descriptorpb.ServiceOptions{},
	}
	proto.SetExtension(serv.Options, annotations.E_DefaultHost, "foo.googleapis.com")

	// Test some annotations
	proto.SetExtension(serv.Method[0].Options, annotations.E_Http, &annotations.HttpRule{
		Pattern: &annotations.HttpRule_Get{
			Get: "/zip",
		},
	})
	proto.SetExtension(serv.Method[1].Options, annotations.E_Http, &annotations.HttpRule{
		Pattern: &annotations.HttpRule_Post{
			Post: "/zap",
		},
	})

	servHostPort := &descriptorpb.ServiceDescriptorProto{
		Name: proto.String("ServHostPort"),
		Method: []*descriptorpb.MethodDescriptorProto{
			{Name: proto.String("Smack")},
		},
		Options: &descriptorpb.ServiceOptions{},
	}
	proto.SetExtension(servHostPort.Options, annotations.E_DefaultHost, "foo.googleapis.com:1234")

	servIAMOverride := &descriptorpb.ServiceDescriptorProto{
		Name: proto.String("ServIamOverride"),
		Method: []*descriptorpb.MethodDescriptorProto{
			{Name: proto.String("GetIamPolicy")},
			{Name: proto.String("SetIamPolicy")},
			{Name: proto.String("TestIamPermissions")},
		},
		Options: &descriptorpb.ServiceOptions{},
	}
	proto.SetExtension(servIAMOverride.Options, annotations.E_DefaultHost, "foo.googleapis.com:1234")

	f := &descriptorpb.FileDescriptorProto{
		Package: proto.String("bar"),
		Service: []*descriptorpb.ServiceDescriptorProto{serv, servHostPort, servIAMOverride},
	}
	files := append(g.getMixinFiles(), f)
	g.descInfo = pbinfo.Of(files)

	for _, tst := range []struct {
		tstName, servName string
		serv              *descriptorpb.ServiceDescriptorProto
		hasOverride       bool
		imports           map[pbinfo.ImportSpec]bool
	}{
		{
			tstName:  "foo_opt",
			servName: "Foo",
			serv:     serv,
			imports: map[pbinfo.ImportSpec]bool{
				{Path: "google.golang.org/api/option"}:                 true,
				{Path: "google.golang.org/api/option/internaloption"}:  true,
				{Path: "google.golang.org/grpc/codes"}:                 true,
				{Path: "math"}:                                         true,
				{Path: "time"}:                                         true,
				{Name: "gax", Path: "github.com/googleapis/gax-go/v2"}: true,
			},
		},
		{
			tstName:  "empty_opt",
			servName: "",
			serv:     serv,
			imports: map[pbinfo.ImportSpec]bool{
				{Path: "google.golang.org/api/option"}:                 true,
				{Path: "google.golang.org/api/option/internaloption"}:  true,
				{Path: "google.golang.org/grpc/codes"}:                 true,
				{Path: "math"}:                                         true,
				{Path: "time"}:                                         true,
				{Name: "gax", Path: "github.com/googleapis/gax-go/v2"}: true,
			},
		},
		{
			tstName:  "host_port_opt",
			servName: "Bar",
			serv:     servHostPort,
			imports: map[pbinfo.ImportSpec]bool{
				{Path: "google.golang.org/api/option"}:                 true,
				{Path: "google.golang.org/api/option/internaloption"}:  true,
				{Path: "google.golang.org/grpc/codes"}:                 true,
				{Path: "math"}:                                         true,
				{Path: "time"}:                                         true,
				{Name: "gax", Path: "github.com/googleapis/gax-go/v2"}: true,
			},
		},
		{
			tstName:     "iam_override_opt",
			servName:    "Baz",
			serv:        servIAMOverride,
			hasOverride: true,
			imports: map[pbinfo.ImportSpec]bool{
				{Path: "google.golang.org/api/option"}:                 true,
				{Path: "google.golang.org/api/option/internaloption"}:  true,
				{Path: "google.golang.org/grpc/codes"}:                 true,
				{Path: "math"}:                                         true,
				{Path: "time"}:                                         true,
				{Name: "gax", Path: "github.com/googleapis/gax-go/v2"}: true,
			},
		},
	} {
		t.Run(tst.tstName, func(t *testing.T) {
			g.reset()
			g.hasIAMPolicyOverrides = tst.hasOverride
			if err := g.clientOptions(tst.serv, tst.servName); err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(g.imports, tst.imports); diff != "" {
				t.Errorf("TestClientOpt(%s): imports got(-),want(+):\n%s", tst.tstName, diff)
			}
			txtdiff.Diff(t, g.pt.String(), filepath.Join("testdata", tst.tstName+".want"))
		})
	}
}

func TestServiceDoc(t *testing.T) {
	serv := &descriptorpb.ServiceDescriptorProto{
		Name: proto.String("MyService"),
	}

	var g generator
	g.comments = make(map[protoiface.MessageV1]string)

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
		serv.Options = &descriptorpb.ServiceOptions{
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

	cop := &descriptorpb.DescriptorProto{
		Name: proto.String("Operation"),
	}

	servPlain := &descriptorpb.ServiceDescriptorProto{
		Name: proto.String("Foo"),
		Method: []*descriptorpb.MethodDescriptorProto{
			{
				Name:       proto.String("Zip"),
				InputType:  proto.String(".mypackage.Bar"),
				OutputType: proto.String(".mypackage.Foo"),
				Options:    &descriptorpb.MethodOptions{},
			},
		},
	}
	servLRO := &descriptorpb.ServiceDescriptorProto{
		Name: proto.String("Foo"),
		Method: []*descriptorpb.MethodDescriptorProto{
			{
				Name:       proto.String("Zip"),
				InputType:  proto.String(".mypackage.Bar"),
				OutputType: proto.String(".google.longrunning.Operation"),
				Options:    &descriptorpb.MethodOptions{},
			},
		},
	}
	servDeprecated := &descriptorpb.ServiceDescriptorProto{
		Name: proto.String("Foo"),
		Method: []*descriptorpb.MethodDescriptorProto{
			{
				Name:       proto.String("Zip"),
				InputType:  proto.String(".mypackage.Bar"),
				OutputType: proto.String(".mypackage.Foo"),
				Options:    &descriptorpb.MethodOptions{},
			},
		},
		Options: &descriptorpb.ServiceOptions{
			Deprecated: proto.Bool(true),
		},
	}

	opS := &descriptorpb.ServiceDescriptorProto{
		Name:   proto.String("FooOperationService"),
		Method: []*descriptorpb.MethodDescriptorProto{},
	}

	customOpOpts := &descriptorpb.MethodOptions{}
	servCustomOp := &descriptorpb.ServiceDescriptorProto{
		Name: proto.String("Foo"),
		Method: []*descriptorpb.MethodDescriptorProto{
			{
				Name:       proto.String("Zip"),
				InputType:  proto.String(".mypackage.Bar"),
				OutputType: proto.String(".mypackage.Operation"),
				Options:    customOpOpts,
			},
		},
	}
	for _, s := range []*descriptorpb.ServiceDescriptorProto{servPlain, servDeprecated, servLRO} {
		proto.SetExtension(s.Method[0].GetOptions(), annotations.E_Http, &annotations.HttpRule{
			Pattern: &annotations.HttpRule_Get{
				Get: "/zip",
			},
		})
	}

	for _, tst := range []struct {
		tstName      string
		servName     string
		mixins       mixins
		serv         *descriptorpb.ServiceDescriptorProto
		customOpServ *descriptorpb.ServiceDescriptorProto
		parameter    *string
		imports      map[pbinfo.ImportSpec]bool
		wantNumSnps  int
		setExt       func() (protoreflect.ExtensionType, interface{})
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
			imports: map[pbinfo.ImportSpec]bool{
				{Path: "context"}:                                                                  true,
				{Path: "google.golang.org/grpc"}:                                                   true,
				{Path: "google.golang.org/api/option"}:                                             true,
				{Name: "gtransport", Path: "google.golang.org/api/transport/grpc"}:                 true,
				{Name: "iampb", Path: "cloud.google.com/go/iam/apiv1/iampb"}:                       true,
				{Name: "locationpb", Path: "google.golang.org/genproto/googleapis/cloud/location"}: true,
				{Name: "mypackagepb", Path: "github.com/googleapis/mypackage"}:                     true,
			},
			wantNumSnps: 6,
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
			imports: map[pbinfo.ImportSpec]bool{
				{Path: "context"}:                                                                  true,
				{Path: "google.golang.org/api/option"}:                                             true,
				{Path: "google.golang.org/api/option/internaloption"}:                              true,
				{Path: "google.golang.org/grpc"}:                                                   true,
				{Path: "net/http"}:                                                                 true,
				{Name: "httptransport", Path: "google.golang.org/api/transport/http"}:              true,
				{Name: "iampb", Path: "cloud.google.com/go/iam/apiv1/iampb"}:                       true,
				{Name: "locationpb", Path: "google.golang.org/genproto/googleapis/cloud/location"}: true,
				{Name: "mypackagepb", Path: "github.com/googleapis/mypackage"}:                     true,
			},
			wantNumSnps: 6,
		},
		{
			tstName:   "empty_client_init",
			servName:  "",
			serv:      servPlain,
			parameter: proto.String("go-gapic-package=path;mypackage,transport=grpc+rest"),
			imports: map[pbinfo.ImportSpec]bool{
				{Path: "google.golang.org/api/option"}:                                true,
				{Path: "google.golang.org/api/option/internaloption"}:                 true,
				{Path: "net/http"}:                                                    true,
				{Path: "context"}:                                                     true,
				{Path: "google.golang.org/grpc"}:                                      true,
				{Name: "gtransport", Path: "google.golang.org/api/transport/grpc"}:    true,
				{Name: "mypackagepb", Path: "github.com/googleapis/mypackage"}:        true,
				{Name: "httptransport", Path: "google.golang.org/api/transport/http"}: true,
			},
			wantNumSnps: 1,
		},
		{
			tstName: "lro_client_init",
			mixins: mixins{
				"google.longrunning.Operations": operationsMethods(),
			},
			servName:  "Foo",
			serv:      servLRO,
			parameter: proto.String("go-gapic-package=path;mypackage"),
			imports: map[pbinfo.ImportSpec]bool{
				{Name: "gtransport", Path: "google.golang.org/api/transport/grpc"}:                     true,
				{Name: "longrunningpb", Path: "cloud.google.com/go/longrunning/autogen/longrunningpb"}: true,
				{Name: "lroauto", Path: "cloud.google.com/go/longrunning/autogen"}:                     true,
				{Name: "mypackagepb", Path: "github.com/googleapis/mypackage"}:                         true,
				{Path: "context"}:                      true,
				{Path: "google.golang.org/api/option"}: true,
				{Path: "google.golang.org/grpc"}:       true,
			},
			wantNumSnps: 6,
		},
		{
			tstName:   "deprecated_client_init",
			servName:  "",
			serv:      servDeprecated,
			parameter: proto.String("go-gapic-package=path;mypackage,transport=grpc+rest"),
			imports: map[pbinfo.ImportSpec]bool{
				{Name: "gtransport", Path: "google.golang.org/api/transport/grpc"}:    true,
				{Name: "httptransport", Path: "google.golang.org/api/transport/http"}: true,
				{Name: "mypackagepb", Path: "github.com/googleapis/mypackage"}:        true,
				{Path: "context"}:                                     true,
				{Path: "google.golang.org/api/option"}:                true,
				{Path: "google.golang.org/api/option/internaloption"}: true,
				{Path: "google.golang.org/grpc"}:                      true,
				{Path: "net/http"}:                                    true,
			},
			wantNumSnps: 1,
		},
		{
			tstName:      "custom_op_init",
			servName:     "",
			serv:         servCustomOp,
			customOpServ: opS,
			parameter:    proto.String("go-gapic-package=path;mypackage,transport=rest,diregapic"),
			imports: map[pbinfo.ImportSpec]bool{
				{Path: "context"}:                                                     true,
				{Path: "google.golang.org/api/option"}:                                true,
				{Path: "google.golang.org/api/option/internaloption"}:                 true,
				{Path: "google.golang.org/grpc"}:                                      true,
				{Path: "net/http"}:                                                    true,
				{Name: "httptransport", Path: "google.golang.org/api/transport/http"}: true,
				{Name: "mypackagepb", Path: "github.com/googleapis/mypackage"}:        true,
			},
			wantNumSnps: 1,
			setExt: func() (protoreflect.ExtensionType, interface{}) {
				return extendedops.E_OperationService, opS.GetName()
			},
		},
	} {
		t.Run(tst.tstName, func(t *testing.T) {
			setExt := tst.setExt
			if setExt == nil {
				setExt = func() (protoreflect.ExtensionType, interface{}) {
					return annotations.E_Http, &annotations.HttpRule{
						Pattern: &annotations.HttpRule_Get{
							Get: "/zip",
						},
					}
				}
			}
			ext, value := setExt()
			proto.SetExtension(tst.serv.Method[0].GetOptions(), ext, value)

			fds := append(mixinDescriptors(), &descriptorpb.FileDescriptorProto{
				Package: proto.String("mypackage"),
				Options: &descriptorpb.FileOptions{
					GoPackage: proto.String("github.com/googleapis/mypackage/v1"),
				},
				Service: []*descriptorpb.ServiceDescriptorProto{tst.serv, tst.customOpServ},
				MessageType: []*descriptorpb.DescriptorProto{
					{
						Name: proto.String("Bar"),
					},
					{
						Name: proto.String("Foo"),
					},
					cop,
				},
			})
			request := pluginpb.CodeGeneratorRequest{
				Parameter: tst.parameter,
				ProtoFile: fds,
			}
			g.init(&request)
			g.comments = map[protoiface.MessageV1]string{
				tst.serv:                "Foo service does stuff.",
				tst.serv.GetMethod()[0]: "Does some stuff.",
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
			g.aux.customOp = &customOp{
				message: cop,
			}
			if tst.customOpServ != nil {
				g.customOpServices = map[*descriptorpb.ServiceDescriptorProto]*descriptorpb.ServiceDescriptorProto{
					tst.serv: tst.customOpServ,
				}
			}

			g.reset()
			sm := snippets.NewMetadata("mypackage", "github.com/googleapis/mypackage", "mypackagego")
			sm.AddService(tst.serv.GetName(), "mypackage.googleapis.com")
			for _, m := range tst.serv.GetMethod() {
				sm.AddMethod(tst.serv.GetName(), m.GetName(), "mypackage", tst.serv.GetName(), 50)
			}
			for _, m := range g.getMixinMethods() {
				sm.AddMethod(tst.serv.GetName(), m.GetName(), "mypackage", tst.serv.GetName(), 50)
			}
			g.snippetMetadata = sm
			g.makeClients(tst.serv, tst.servName)

			if diff := cmp.Diff(g.imports, tst.imports); diff != "" {
				t.Errorf("ClientInit(%s) imports got(-),want(+):\n%s", tst.tstName, diff)
			}

			txtdiff.Diff(t, g.pt.String(), filepath.Join("testdata", tst.tstName+".want"))
			mi := g.snippetMetadata.ToMetadataIndex()
			if got := len(mi.Snippets); got != tst.wantNumSnps {
				t.Errorf("%s: got %d want len %d", t.Name(), tst.wantNumSnps, got)
			}
			for _, snp := range mi.Snippets {
				if got := snp.ClientMethod.Parameters[0].Name; got != "ctx" {
					t.Errorf("%s: got %s want ctx,", t.Name(), got)
				}
				if got := snp.ClientMethod.Parameters[1].Name; got != "req" {
					t.Errorf("%s: got %s want req,", t.Name(), got)
				}
				if got := snp.ClientMethod.Parameters[2].Name; got != "opts" {
					t.Errorf("%s: got %s want opts,", t.Name(), got)
				}
				if snp.ClientMethod.ShortName != "CancelOperation" && snp.ClientMethod.ShortName != "DeleteOperation" && snp.ClientMethod.ResultType == "" {
					t.Errorf("%s: got empty string, want ResultType for %s", t.Name(), snp.ClientMethod.ShortName)
				}
			}
		})
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

// mixinDescriptors is used for testing purposes only.
func mixinDescriptors() []*descriptorpb.FileDescriptorProto {
	files := []*descriptorpb.FileDescriptorProto{}
	for _, fds := range mixinFiles {
		files = append(files, fds...)
	}
	return files
}
