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
	"unicode"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// ProtoType represents a type in protobuf descriptors.
// It is an interface implemented by DescriptorProto and EnumDescriptorProto.
type ProtoType interface {
	proto.Message
	GetName() string
}

// Info provides lookup tables for various protobuf properties.
// For example, we can look up a type by name without iterating the entire
// descriptor.
type Info struct {
	// Maps services and messages to the file containing them,
	// so we can figure out the import.
	ParentFile map[proto.Message]*descriptorpb.FileDescriptorProto

	// NOTE(pongad): ParentElement and sub-types are only used in samples.
	// They are added in the shared package because they share a lot of similarities
	// with things that are already here. Maybe revisit this in the future?

	// Maps a protobuf element to the enclosing scope.
	// If enum E is defined in message M which is in file F,
	// ParentElement[E]=M, ParentElement[M]=nil, and ParentFile[M]=F
	ParentElement map[ProtoType]ProtoType

	// Maps type names to their messages.
	Type map[string]ProtoType

	// Maps service names to their descriptors.
	Serv map[string]*descriptorpb.ServiceDescriptorProto

	// PkgOverrides is file-to-import mapping used to override the
	// go_package option in the given proto file.
	PkgOverrides map[string]string
}

// Of creates Info from given protobuf files.
func Of(files []*descriptorpb.FileDescriptorProto) Info {
	info := Info{
		ParentFile:    map[proto.Message]*descriptorpb.FileDescriptorProto{},
		ParentElement: map[ProtoType]ProtoType{},
		Type:          map[string]ProtoType{},
		Serv:          map[string]*descriptorpb.ServiceDescriptorProto{},
		PkgOverrides:  map[string]string{},
	}

	for _, f := range files {
		// ParentFile
		for _, m := range f.GetMessageType() {
			info.ParentFile[m] = f
		}
		for _, e := range f.GetEnumType() {
			info.ParentFile[e] = f
		}
		for _, s := range f.GetService() {
			info.ParentFile[s] = f
			for _, m := range s.GetMethod() {
				info.ParentFile[m] = f
				info.ParentElement[m] = s
			}
		}

		// Type
		for _, m := range f.GetMessageType() {
			// In descriptors, putting the dot in front means the name is fully-qualified.
			addMessage(info.Type, info.ParentElement, "."+f.GetPackage(), m, nil)
		}
		for _, e := range f.GetEnumType() {
			info.Type["."+f.GetPackage()+"."+e.GetName()] = e
		}

		// Serv
		for _, s := range f.GetService() {
			fullyQualifiedName := fmt.Sprintf(".%s.%s", f.GetPackage(), s.GetName())
			info.Serv[fullyQualifiedName] = s
		}
	}

	return info
}

func addMessage(typMap map[string]ProtoType, parentMap map[ProtoType]ProtoType, prefix string, msg, parentMsg *descriptorpb.DescriptorProto) {
	fullName := prefix + "." + msg.GetName()
	typMap[fullName] = msg
	if parentMsg != nil {
		parentMap[msg] = parentMsg
	}

	for _, subMsg := range msg.GetNestedType() {
		addMessage(typMap, parentMap, fullName, subMsg, msg)
	}

	for _, subEnum := range msg.GetEnumType() {
		typMap[fullName+"."+subEnum.GetName()] = subEnum
		parentMap[subEnum] = msg
	}

	for _, field := range msg.GetField() {
		parentMap[field] = msg
	}
}

// ImportSpec represents a Go module import path and an optional alias.
type ImportSpec struct {
	Name, Path string
}

// NameSpec reports the name and ImportSpec of e.
//
// The reported name is the same with how protoc-gen-go refers to e.
// E.g. if type B is nested under A, then the name of type B is "A_B".
func (in *Info) NameSpec(e ProtoType) (string, ImportSpec, error) {
	appendpb := func(n string) string {
		if !strings.HasSuffix(n, "pb") {
			n += "pb"
		}
		return n
	}

	topLvl := e
	var nameParts []string
	for e2 := e; e2 != nil; e2 = in.ParentElement[e2] {
		topLvl = e2
		nameParts = append(nameParts, e2.GetName())
	}
	for i, l := 0, len(nameParts); i < l/2; i++ {
		nameParts[i], nameParts[l-i-1] = nameParts[l-i-1], nameParts[i]
	}
	name := strings.Join(nameParts, "_")

	var eTxt interface{} = e
	if et, ok := eTxt.(interface{ GetName() string }); ok {
		eTxt = et.GetName()
	}

	fdesc := in.ParentFile[topLvl]
	if fdesc == nil {
		return "", ImportSpec{}, fmt.Errorf("can't determine import path for %q; can't find parent file", eTxt)
	}

	pkg := fdesc.GetOptions().GetGoPackage()
	if pkgOverride, ok := in.PkgOverrides[fdesc.GetName()]; ok {
		pkg = pkgOverride
	}
	if pkg == "" {
		return "", ImportSpec{}, fmt.Errorf("can't determine import path for %q, file %q missing `option go_package`", eTxt, fdesc.GetName())
	}

	if p := strings.IndexByte(pkg, ';'); p >= 0 {
		return name, ImportSpec{Path: pkg[:p], Name: appendpb(pkg[p+1:])}, nil
	}

	for {
		p := strings.LastIndexByte(pkg, '/')
		if p < 0 {
			return name, ImportSpec{Path: pkg, Name: appendpb(pkg)}, nil
		}
		elem := pkg[p+1:]
		if len(elem) >= 2 && elem[0] == 'v' && elem[1] >= '0' && elem[1] <= '9' {
			// It's a version number; skip so we get a more meaningful name
			pkg = pkg[:p]
			continue
		}
		return name, ImportSpec{Path: pkg, Name: appendpb(elem)}, nil
	}
}

// ImportSpec reports the ImportSpec for package containing protobuf element e.
// Deprecated: Use NameSpec instead.
func (in *Info) ImportSpec(e ProtoType) (ImportSpec, error) {
	_, imp, err := in.NameSpec(e)
	return imp, err
}

// ReduceServName removes redundant components from the service name.
// For example, FooServiceV2 -> Foo.
// The returned name is used as part of longer names, like FooClient.
// If the package name and the service name is the same,
// ReduceServName returns empty string, so we get foo.Client instead of foo.FooClient.
func ReduceServName(svc, pkg string) string {
	// remove trailing version
	if p := strings.LastIndexByte(svc, 'V'); p >= 0 {
		isVer := true
		for _, r := range svc[p+1:] {
			if !unicode.IsDigit(r) {
				isVer = false
				break
			}
		}
		if isVer {
			svc = svc[:p]
		}
	}

	svc = strings.TrimSuffix(svc, "Service")
	if strings.EqualFold(svc, pkg) {
		svc = ""
	}

	// This is a special case for IAM and should not be
	// extended to support any new API name containing
	// an acronym.
	//
	// In order to avoid a breaking change for IAM
	// clients, we must keep consistent identifier casing.
	if strings.Contains(svc, "IAM") {
		svc = strings.ReplaceAll(svc, "IAM", "Iam")
	}

	return svc
}
