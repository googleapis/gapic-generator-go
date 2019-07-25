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
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
)

// enumFmt returns a function that transforms a protobuf enum value literal into it's corresponding Go const.
func enumFmt(info pbinfo.Info, enum *descriptor.EnumDescriptorProto) func(*generator, string) (string, error) {
	return func(g *generator, s string) (string, error) {
		parts, impSpec, err := ancestors(info, enum, g)
		if err != nil {
			return "", err
		}
		if len(parts) > 1 {
			// An enum is scoped using C++ rule. If MessageType contains EnumType,
			// then elements are accessed as MessageType.Element, not MessageType.EnumType.Element.
			//
			// An exception to this rule is top-level enums, where len(parts)==1; they are qualified by
			// EnumType.
			parts = parts[:len(parts)-1]
		}

		return impSpec.Name + "." + strings.Join(parts, "_") + "_" + s, nil
	}
}

// enumType returns the Go Type name of a protobuf enum type.
func enumType(info pbinfo.Info, enum *descriptor.EnumDescriptorProto, g *generator) (string, error) {
	parts, impSpec, err := ancestors(info, enum, g)
	if err != nil {
		return "", err
	}
	return impSpec.Name + "." + strings.Join(parts, "_"), nil
}

// ancestors returns the names of all the types down to `enum`, included and the ImportSpec.
// of the slice.
func ancestors(info pbinfo.Info, enum *descriptor.EnumDescriptorProto, g *generator) ([]string, pbinfo.ImportSpec, error) {
	var element pbinfo.ProtoType = enum
	var topElement pbinfo.ProtoType
	var parts []string

	for ; element != nil; element = info.ParentElement[element] {
		parts = append(parts, element.GetName())
		topElement = element
	}

	// We made array in [child, parent, grandparent] order, we want grandparent fist.
	for i := 0; i < len(parts)/2; i++ {
		parts[i], parts[len(parts)-1-i] = parts[len(parts)-1-i], parts[i]
	}

	impSpec, err := info.ImportSpec(topElement)
	if err != nil {
		return []string{}, pbinfo.ImportSpec{}, err
	}
	return parts, impSpec, nil
}
