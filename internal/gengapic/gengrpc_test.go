// Copyright 2024 Google LLC
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

	conf "github.com/googleapis/gapic-generator-go/internal/grpc_service_config"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"github.com/googleapis/gapic-generator-go/internal/txtdiff"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/genproto/googleapis/api/serviceconfig"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	duration "google.golang.org/protobuf/types/known/durationpb"
)

func TestServiceRenaming(t *testing.T) {
	inputType := &descriptorpb.DescriptorProto{
		Name: proto.String("InputType"),
	}
	outputType := &descriptorpb.DescriptorProto{
		Name: proto.String("OutputType"),
	}
	metadataType := &descriptorpb.DescriptorProto{
		Name: proto.String("MetadataType"),
	}

	file := &descriptorpb.FileDescriptorProto{
		Package: proto.String("my.pkg"),
		Options: &descriptorpb.FileOptions{
			GoPackage: proto.String("mypackage"),
		},
	}
	serv := &descriptorpb.ServiceDescriptorProto{
		Name: proto.String("Foo"),
	}

	var g generator
	g.imports = map[pbinfo.ImportSpec]bool{}
	g.cfg = &generatorConfig{
		pkgName:           "pkg",
		transports:        []transport{grpc},
		featureEnablement: map[featureID]struct{}{OpenTelemetryTracingFeature: {}},
		APIServiceConfig: &serviceconfig.Service{
			Name: "foo.googleapis.com",
			Publishing: &annotations.Publishing{
				LibrarySettings: []*annotations.ClientLibrarySettings{
					{
						GoSettings: &annotations.GoSettings{
							RenamedServices: map[string]string{"Foo": "Bar"},
						},
					},
				},
			},
		},
	}
	cpb := &conf.ServiceConfig{
		MethodConfig: []*conf.MethodConfig{
			{
				Name: []*conf.MethodConfig_Name{
					{
						Service: "my.pkg.Foo",
					},
				},
				Timeout: &duration.Duration{Seconds: 10},
			},
		},
	}
	data, err := protojson.Marshal(cpb)
	if err != nil {
		t.Error(err)
	}
	in := bytes.NewReader(data)
	g.cfg.gRPCServiceConfig, err = conf.New(in)
	if err != nil {
		t.Error(err)
	}

	commonTypes(&g)
	for _, typ := range []*descriptorpb.DescriptorProto{
		inputType, outputType, metadataType,
	} {
		g.descInfo.Type[".my.pkg."+*typ.Name] = typ
		g.descInfo.ParentFile[typ] = file
	}
	g.descInfo.ParentFile[serv] = file
	g.descInfo.ParentElement = map[pbinfo.ProtoType]pbinfo.ProtoType{}

	m := &descriptorpb.MethodDescriptorProto{
		Name:       proto.String("Baz"),
		InputType:  proto.String(".my.pkg.InputType"),
		OutputType: proto.String(emptyType),
	}

	serv.Method = []*descriptorpb.MethodDescriptorProto{
		m,
	}
	g.descInfo.ParentFile[serv] = file
	g.descInfo.ParentElement[m] = serv
	imp, err := g.descInfo.ImportSpec(serv)
	if err != nil {
		t.Fatal(err)
	}

	// Test the generation of the client boilerplate and single gRPC method rename.
	g.grpcClientInit(serv, "Bar", imp, false)
	if err := g.genGRPCMethod("Bar", serv, m); err != nil {
		t.Fatal(err)
	}

	txtdiff.Diff(t, g.pt.String(), filepath.Join("testdata", "service_rename.want"))
}
