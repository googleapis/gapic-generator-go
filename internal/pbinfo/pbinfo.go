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

// Package pbinfo provides convenience types for looking up protobuf elements.
package pbinfo

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/gapic-generator-go/internal/errors"
)

type Info struct {
	// Maps services and messages to the file containing them,
	// so we can figure out the import.
	ParentFile map[proto.Message]*descriptor.FileDescriptorProto

	// Maps type names to their messages.
	Type map[string]*descriptor.DescriptorProto

	// Maps service names to their descriptors.
	Serv map[string]*descriptor.ServiceDescriptorProto
}

// Of creates Info from given protobuf files.
func Of(files []*descriptor.FileDescriptorProto) Info {
	info := Info{
		ParentFile: map[proto.Message]*descriptor.FileDescriptorProto{},
		Type:       map[string]*descriptor.DescriptorProto{},
		Serv:       map[string]*descriptor.ServiceDescriptorProto{},
	}

	for _, f := range files {
		// ParentFile
		for _, m := range f.MessageType {
			info.ParentFile[m] = f
		}
		for _, s := range f.Service {
			info.ParentFile[s] = f
		}

		// Type
		for _, m := range f.MessageType {
			// In descriptors, putting the dot in front means the name is fully-qualified.
			fullyQualifiedName := fmt.Sprintf(".%s.%s", f.GetPackage(), m.GetName())
			info.Type[fullyQualifiedName] = m
		}

		// Serv
		for _, s := range f.Service {
			fullyQualifiedName := fmt.Sprintf(".%s.%s", f.GetPackage(), s.GetName())
			info.Serv[fullyQualifiedName] = s
		}
	}

	return info
}

type ImportSpec struct {
	Name, Path string
}

// ImportSpec reports the ImportSpec for package containing protobuf element e.
func (in *Info) ImportSpec(e proto.Message) (ImportSpec, error) {
	var eTxt interface{} = e
	if et, ok := eTxt.(interface{ GetName() string }); ok {
		eTxt = et.GetName()
	}

	fdesc := in.ParentFile[e]
	if fdesc == nil {
		return ImportSpec{}, errors.E(nil, "can't determine import path for %v; can't find parent file", eTxt)
	}

	pkg := fdesc.GetOptions().GetGoPackage()
	if pkg == "" {
		return ImportSpec{}, errors.E(nil, "can't determine import path for %v, file %q missing `option go_package`", eTxt, fdesc.GetName())
	}

	if p := strings.IndexByte(pkg, ';'); p >= 0 {
		return ImportSpec{Path: pkg[:p], Name: pkg[p+1:] + "pb"}, nil
	}

	for {
		p := strings.LastIndexByte(pkg, '/')
		if p < 0 {
			return ImportSpec{Path: pkg, Name: pkg + "pb"}, nil
		}
		elem := pkg[p+1:]
		if len(elem) >= 2 && elem[0] == 'v' && elem[1] >= '0' && elem[1] <= '9' {
			// It's a version number; skip so we get a more meaningful name
			pkg = pkg[:p]
			continue
		}
		return ImportSpec{Path: pkg, Name: elem + "pb"}, nil
	}
}
