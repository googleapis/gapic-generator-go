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
