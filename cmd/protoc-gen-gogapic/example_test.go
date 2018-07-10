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
	"path/filepath"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

func TestExample(t *testing.T) {
	var g generator
	g.imports = map[importSpec]bool{}

	inputType := &descriptor.DescriptorProto{
		Name: proto.String("InputType"),
	}
	outputType := &descriptor.DescriptorProto{
		Name: proto.String("OutputType"),
	}

	g.types = map[string]*descriptor.DescriptorProto{
		".my.pkg.InputType":  inputType,
		".my.pkg.OutputType": outputType,
	}

	file := &descriptor.FileDescriptorProto{
		Options: &descriptor.FileOptions{
			GoPackage: proto.String("mypackage"),
		},
	}

	g.parentFile = map[proto.Message]*descriptor.FileDescriptorProto{
		inputType:  file,
		outputType: file,
	}

	serv := &descriptor.ServiceDescriptorProto{
		Name: proto.String("Foo"),
		Method: []*descriptor.MethodDescriptorProto{
			{
				Name:       proto.String("GetEmptyThing"),
				InputType:  proto.String(".my.pkg.InputType"),
				OutputType: proto.String(emptyType),
			},
			{
				Name:       proto.String("GetOneThing"),
				InputType:  proto.String(".my.pkg.InputType"),
				OutputType: proto.String(".my.pkg.OutputType"),
			},
			{
				Name:       proto.String("GetBigThing"),
				InputType:  proto.String(".my.pkg.InputType"),
				OutputType: proto.String(".google.longrunning.Operation"),
			},
		},
	}
	for _, tst := range []struct {
		tstName, pkgName string
	}{
		{tstName: "empty_example", pkgName: "Foo"},
		{tstName: "foo_example", pkgName: "Bar"},
	} {
		g.reset()
		g.genExampleFile(serv, tst.pkgName)
		diff(t, tst.tstName, []byte(g.sb.String()), filepath.Join("testdata", tst.tstName+".want"))
	}
}
