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
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/gapic-generator-go/internal/errors"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
)

func (g *generator) getGoTypeName(typ initType) (string, error) {
	if typ.prim != 0 {
		goType, ok := pbinfo.GoTypeForPrim[typ.prim]
		if !ok {
			return "", errors.E(nil, "unrecognized primitive type: %s", typ.prim)
		}
		return goType, nil
	}

	if enum, ok2 := typ.desc.(*descriptor.EnumDescriptorProto); ok2 {
		goType, err := goTypeForEnum(g.descInfo, enum)
		if err != nil {
			return "", errors.E(err, "unrecognized enum type: %s", enum)
		}
		return goType, nil
	}

	return "", errors.E(nil, "internal: unhandled initType: %v", typ)
}
