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
	"path/filepath"
	"testing"

	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"github.com/googleapis/gapic-generator-go/internal/testing/sample"
	"github.com/googleapis/gapic-generator-go/internal/txtdiff"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func TestDocFile(t *testing.T) {
	g := generator{
		apiName: sample.ServiceTitle,
		imports: map[pbinfo.ImportSpec]bool{},
		cfg: &generatorConfig{
			pkgPath:          sample.GoPackagePath,
			pkgName:          sample.GoPackageName,
			transports:       []transport{grpc, rest},
			APIServiceConfig: sample.ServiceConfig(),
		},
	}

	inputType := sample.InputType(sample.CreateRequest)
	outputType := sample.OutputType(sample.Resource)
	file := sample.File()

	commonTypes(&g)
	for _, typ := range []*descriptorpb.DescriptorProto{
		inputType, outputType,
	} {
		typName := sample.DescriptorInfoTypeName(typ.GetName())
		g.descInfo.Type[typName] = typ
		g.descInfo.ParentFile[typ] = file
	}

	serv := sample.Service()
	for _, tst := range []struct {
		relLvl, want string
	}{
		{
			want: filepath.Join("testdata", "doc_file.want"),
		},
		{
			relLvl: alpha,
			want:   filepath.Join("testdata", "doc_file_alpha.want"),
		},
		{
			relLvl: beta,
			want:   filepath.Join("testdata", "doc_file_beta.want"),
		},
	} {
		t.Run(tst.want, func(t *testing.T) {
			g.cfg.relLvl = tst.relLvl
			g.genDocFile(sample.Year, []*descriptorpb.ServiceDescriptorProto{serv})
			txtdiff.Diff(t, g.pt.String(), tst.want)
			g.reset()
		})
	}
}

func TestDocFile_APIVersionSection(t *testing.T) {
	g := generator{
		apiName: sample.ServiceTitle,
		imports: map[pbinfo.ImportSpec]bool{},
		cfg: &generatorConfig{
			pkgPath:          sample.GoPackagePath,
			pkgName:          sample.GoPackageName,
			transports:       []transport{grpc, rest},
			APIServiceConfig: sample.ServiceConfig(),
		},
	}

	inputType := sample.InputType(sample.CreateRequest)
	outputType := sample.OutputType(sample.Resource)
	file := sample.File()

	commonTypes(&g)
	for _, typ := range []*descriptorpb.DescriptorProto{
		inputType, outputType,
	} {
		typName := sample.DescriptorInfoTypeName(typ.GetName())
		g.descInfo.Type[typName] = typ
		g.descInfo.ParentFile[typ] = file
	}

	serv := sample.Service()
	serv.Options = sample.APIVersionOptions("2024-09-14")

	g.genDocFile(sample.Year, []*descriptorpb.ServiceDescriptorProto{serv})
	txtdiff.Diff(t, g.pt.String(), filepath.Join("testdata", "doc_file_api_version_section.want"))
	g.reset()
}

func TestDocFileEmptyService(t *testing.T) {
	g := generator{
		apiName: sample.ServiceTitle,
		imports: map[pbinfo.ImportSpec]bool{},
		cfg: &generatorConfig{
			pkgPath:          sample.GoPackagePath,
			pkgName:          sample.GoPackageName,
			transports:       []transport{grpc, rest},
			APIServiceConfig: sample.ServiceConfig(),
		},
	}
	inputType := sample.InputType(sample.CreateRequest)
	outputType := sample.OutputType(sample.Resource)
	file := sample.File()

	commonTypes(&g)
	for _, typ := range []*descriptorpb.DescriptorProto{
		inputType, outputType,
	} {
		typName := sample.DescriptorInfoTypeName(typ.GetName())
		g.descInfo.Type[typName] = typ
		g.descInfo.ParentFile[typ] = file
	}

	serv := sample.Service()
	serv.Method = nil
	for _, tst := range []struct {
		relLvl, want string
	}{
		{
			want: filepath.Join("testdata", "doc_file_emptyservice.want"),
		},
		{
			relLvl: alpha,
			want:   filepath.Join("testdata", "doc_file_alpha_emptyservice.want"),
		},
		{
			relLvl: beta,
			want:   filepath.Join("testdata", "doc_file_beta_emptyservice.want"),
		},
		{
			relLvl: deprecated,
			want:   filepath.Join("testdata", "doc_file_deprecated_emptyservice.want"),
		},
	} {
		t.Run(tst.want, func(t *testing.T) {
			g.cfg.relLvl = tst.relLvl
			g.genDocFile(sample.Year, []*descriptorpb.ServiceDescriptorProto{serv})
			txtdiff.Diff(t, g.pt.String(), tst.want)
			g.reset()
		})
	}
}

func TestApiVersionSection(t *testing.T) {
	g := generator{
		cfg: &generatorConfig{
			pkgName:          sample.GoPackageName,
			APIServiceConfig: sample.ServiceConfig(),
		},
	}

	for _, tst := range []struct {
		services  []*descriptorpb.ServiceDescriptorProto
		want      string
		wantEmpty bool
	}{
		{
			want: filepath.Join("testdata", "api_versions_section_list.want"),
			services: []*descriptorpb.ServiceDescriptorProto{
				{
					Name:    proto.String("Foo"),
					Options: sample.APIVersionOptions("2024-09-14"),
				},
				{
					Name:    proto.String("Bar"),
					Options: sample.APIVersionOptions("2024-04-04"),
				}},
		},
		{
			want: filepath.Join("testdata", "api_versions_section_same.want"),
			services: []*descriptorpb.ServiceDescriptorProto{
				{
					Name:    proto.String("Foo"),
					Options: sample.APIVersionOptions("2024-09-14"),
				},
				{
					Name:    proto.String("Bar"),
					Options: sample.APIVersionOptions("2024-09-14"),
				}},
		},
		{
			wantEmpty: true, // no API versions
			services: []*descriptorpb.ServiceDescriptorProto{
				{
					Name: proto.String("Foo"),
				},
				{
					Name: proto.String("Bar"),
				}},
		},
		{
			wantEmpty: true, // no services
		},
	} {
		t.Run(tst.want, func(t *testing.T) {
			g.apiVersionSection(tst.services)
			got := g.pt.String()

			if tst.wantEmpty {
				if got != "" {
					t.Errorf("apiVersionSection(): expected empty API version section got %s", got)
				}
			} else {
				txtdiff.Diff(t, got, tst.want)
			}
			g.reset()
		})
	}
}
