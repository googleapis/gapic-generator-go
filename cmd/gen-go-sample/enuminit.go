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

func describeEnum(info pbinfo.Info, enum *descriptor.EnumDescriptorProto) (valid func(string) bool, format func(*generator, string) (string, error)) {
	valid = func(s string) bool {
		for _, enumVal := range enum.Value {
			if s == enumVal.GetName() {
				return true
			}
		}
		return false
	}

	// TODO(pongad): we probably need something like this for nested messages too.
	format = func(g *generator, s string) (string, error) {
		var element pbinfo.ProtoType = enum
		var topElement pbinfo.ProtoType
		var parts []string

		for ; element != nil; element = info.ParentElement[element] {
			parts = append(parts, element.GetName())
			topElement = element
		}

		// An enum is scoped using C++ rule. If MessageType contains EnumType,
		// then elements are accessed as MessageType.Element, not MessageType.EnumType.Element.
		// TODO(pongad): figure out what happens with top-level enum.
		parts = parts[1:]

		// We made array in [child, parent, grandparent] order, we want grandparent fist.
		for i := 0; i < len(parts)/2; i++ {
			parts[i], parts[len(parts)-1-i] = parts[len(parts)-1-i], parts[i]
		}

		impSpec, err := info.ImportSpec(topElement)
		if err != nil {
			return "", err
		}
		g.imports[impSpec] = true

		return impSpec.Name + "." + strings.Join(parts, "_") + "_" + s, nil
	}

	return
}
