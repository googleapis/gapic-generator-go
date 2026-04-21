package gengapic

import (
	"testing"

	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/genproto/googleapis/api/serviceconfig"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
)

func TestMethodSelectiveGeneration(t *testing.T) {
	fooMethod := &descriptorpb.MethodDescriptorProto{
		Name: proto.String("Foo"),
	}
	barMethod := &descriptorpb.MethodDescriptorProto{
		Name: proto.String("Bar"),
	}
	mixinMethod := &descriptorpb.MethodDescriptorProto{
		Name: proto.String("GetMixin"),
	}

	srv := &descriptorpb.ServiceDescriptorProto{
		Name:   proto.String("MyService"),
		Method: []*descriptorpb.MethodDescriptorProto{fooMethod, barMethod},
	}
	mixinSrv := &descriptorpb.ServiceDescriptorProto{
		Name: proto.String("MixinService"),
		Method: []*descriptorpb.MethodDescriptorProto{mixinMethod},
	}

	file := &descriptorpb.FileDescriptorProto{
		Package: proto.String("my.pkg"),
		Service: []*descriptorpb.ServiceDescriptorProto{srv},
	}
	mixinFile := &descriptorpb.FileDescriptorProto{
		Package: proto.String("google.mixin.v1"),
		Service: []*descriptorpb.ServiceDescriptorProto{mixinSrv},
	}

	descInfo := pbinfo.Info{
		ParentFile: map[proto.Message]*descriptorpb.FileDescriptorProto{
			srv: file,
			fooMethod: file,
			barMethod: file,
			mixinMethod: mixinFile,
			mixinSrv: mixinFile,
		},
		ParentElement: map[pbinfo.ProtoType]pbinfo.ProtoType{
			fooMethod: srv,
			barMethod: srv,
			mixinMethod: mixinSrv,
		},
	}

	tests := []struct {
		name           string
		cfg            *annotations.SelectiveGapicGeneration
		wantFooGen     bool
		wantFooName    string
		wantBarGen     bool
		wantBarName    string
		wantMixinGen   bool
		wantMixinName  string
		wantBasePrefix bool
	}{
		{
			name: "No selective gen",
			cfg:  nil,
			wantFooGen:  true,
			wantFooName: "Foo",
			wantBarGen:  true,
			wantBarName: "Bar",
			wantMixinGen: true,
			wantMixinName: "GetMixin",
			wantBasePrefix: false,
		},
		{
			name: "Empty methods, no generate_omitted",
			cfg: &annotations.SelectiveGapicGeneration{
				Methods: []string{},
				GenerateOmittedAsInternal: false,
			},
			wantFooGen:  true,
			wantFooName: "Foo",
			wantBarGen:  true,
			wantBarName: "Bar",
			wantMixinGen: true,
			wantMixinName: "GetMixin",
			wantBasePrefix: false,
		},
		{
			name: "Empty methods, generate_omitted",
			cfg: &annotations.SelectiveGapicGeneration{
				Methods: []string{},
				GenerateOmittedAsInternal: true,
			},
			wantFooGen:  true,
			wantFooName: "foo",
			wantBarGen:  true,
			wantBarName: "bar",
			wantMixinGen: true,
			wantMixinName: "getMixin",
			wantBasePrefix: true,
		},
		{
			name: "Foo only, no generate_omitted",
			cfg: &annotations.SelectiveGapicGeneration{
				Methods: []string{"my.pkg.MyService.Foo"},
				GenerateOmittedAsInternal: false,
			},
			wantFooGen:  true,
			wantFooName: "Foo",
			wantBarGen:  false,
			wantBarName: "Bar",
			wantMixinGen: false,
			wantMixinName: "GetMixin",
			wantBasePrefix: false,
		},
		{
			name: "Foo and Mixin, generate_omitted",
			cfg: &annotations.SelectiveGapicGeneration{
				Methods: []string{"my.pkg.MyService.Foo", "google.mixin.v1.MixinService.GetMixin"},
				GenerateOmittedAsInternal: true,
			},
			wantFooGen:  true,
			wantFooName: "Foo",
			wantBarGen:  true,
			wantBarName: "bar",
			wantMixinGen: true,
			wantMixinName: "GetMixin",
			wantBasePrefix: true,
		},
	}

	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			var ls []*annotations.ClientLibrarySettings
			if tst.cfg != nil {
				ls = []*annotations.ClientLibrarySettings{
					{
						Version: "my.pkg",
						GoSettings: &annotations.GoSettings{
							Common: &annotations.CommonLanguageSettings{
								SelectiveGapicGeneration: tst.cfg,
							},
						},
					},
				}
			}

			g := &generator{
				descInfo: descInfo,
				cfg: &generatorConfig{
					featureEnablement: map[featureID]struct{}{
						SelectiveGapicGenerationFeature: {},
					},
					APIServiceConfig: &serviceconfig.Service{
						Publishing: &annotations.Publishing{
							LibrarySettings: ls,
						},
					},
				},
				mixins: mixins{
					"google.mixin.v1.MixinService": mixinSrv.GetMethod(),
				},
				clientProtoPkg: "my.pkg",
			}

			fooGen := g.shouldGenerateMethod(srv, fooMethod)
			barGen := g.shouldGenerateMethod(srv, barMethod)
			mixinGen := g.shouldGenerateMethod(srv, mixinMethod)
			fooName := g.methodName(fooMethod)
			barName := g.methodName(barMethod)
			mixinName := g.methodName(mixinMethod)
			isInternal := g.isInternalService(srv)

			if fooGen != tst.wantFooGen {
				t.Errorf("Foo generate = %v, want %v", fooGen, tst.wantFooGen)
			}
			if barGen != tst.wantBarGen {
				t.Errorf("Bar generate = %v, want %v", barGen, tst.wantBarGen)
			}
			if mixinGen != tst.wantMixinGen {
				t.Errorf("Mixin generate = %v, want %v", mixinGen, tst.wantMixinGen)
			}
			if fooName != tst.wantFooName {
				t.Errorf("Foo name = %q, want %q", fooName, tst.wantFooName)
			}
			if barName != tst.wantBarName {
				t.Errorf("Bar name = %q, want %q", barName, tst.wantBarName)
			}
			if mixinName != tst.wantMixinName {
				t.Errorf("Mixin name = %q, want %q", mixinName, tst.wantMixinName)
			}
			if isInternal != tst.wantBasePrefix {
				t.Errorf("isInternalService = %v, want %v", isInternal, tst.wantBasePrefix)
			}
		})
	}
}
