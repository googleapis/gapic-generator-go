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

	"github.com/googleapis/gapic-generator-go/internal/txtdiff"
	"google.golang.org/genproto/googleapis/api/serviceconfig"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"google.golang.org/protobuf/proto"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
)

func TestDocFile(t *testing.T) {
	var g generator
	g.apiName = "Awesome Foo API"
	g.serviceConfig = &serviceconfig.Service{
		Documentation: &serviceconfig.Documentation{
			Summary: "The Awesome Foo API is really really awesome. It enables the use of [Foo](https://api.foo.com) with [Buz](https://api.buz.com) and [Baz](https://api.baz.com) to acclerate `bar`.",
		},
	}
	g.opts = &options{pkgPath: "path/to/awesome", pkgName: "awesome", transports: []transport{grpc}}
	g.imports = map[pbinfo.ImportSpec]bool{}

	inputType := &descriptor.DescriptorProto{
		Name: proto.String("InputType"),
	}
	outputType := &descriptor.DescriptorProto{
		Name: proto.String("OutputType"),
	}

	file := &descriptor.FileDescriptorProto{
		Options: &descriptor.FileOptions{
			GoPackage: proto.String("mypackage"),
		},
	}

	commonTypes(&g)
	for _, typ := range []*descriptor.DescriptorProto{
		inputType, outputType,
	} {
		g.descInfo.Type[".my.pkg."+typ.GetName()] = typ
		g.descInfo.ParentFile[typ] = file
	}

	serv := &descriptor.ServiceDescriptorProto{
		Name: proto.String("Foo"),
		Method: []*descriptor.MethodDescriptorProto{
			{
				Name:       proto.String("GetOneThing"),
				InputType:  proto.String(".my.pkg.InputType"),
				OutputType: proto.String(".my.pkg.OutputType"),
			},
		},
	}

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
		g.opts.relLvl = tst.relLvl
		g.genDocFile(42, []string{"https://foo.bar.com/auth", "https://zip.zap.com/auth"}, serv)
		txtdiff.Diff(t, "doc_file", g.pt.String(), tst.want)
		g.reset()
	}
}
