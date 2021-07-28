// Copyright 2021 Google LLC
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
	"fmt"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

type customOp struct {
	proto     *descriptor.DescriptorProto
	generated bool
}

func (g *generator) customOpProtoName() string {
	f := g.descInfo.ParentFile[g.aux.customOp.proto]
	return fmt.Sprintf(".%s.%s", f.GetPackage(), g.aux.customOp.proto.GetName())
}

func (g *generator) customOpPointerTyp() (string, error) {
	op := g.aux.customOp
	if op == nil {
		return "", nil
	}

	opName, imp, err := g.descInfo.NameSpec(op.proto)
	if err != nil {
		return "", err
	}

	s := fmt.Sprintf("*%s.%s", imp.Name, opName)

	return s, nil
}

func (g *generator) customOpInit(p string) string {
	opName := g.aux.customOp.proto.GetName()

	s := fmt.Sprintf("&%s{proto: %s}", opName, p)

	return s
}

func (g *generator) customOperationType() error {
	op := g.aux.customOp
	if op == nil {
		return nil
	}
	opName := op.proto.GetName()

	ptyp, err := g.customOpPointerTyp()
	if err != nil {
		return err
	}

	p := g.printf

	p("// %s represents a long running operation for this API.", opName)
	p("type %s struct {", opName)
	p("  proto %s", ptyp)
	p("}")
	p("")
	p("// Proto returns the raw type this wraps.")
	p("func (o *%s) Proto() %s {", opName, ptyp)
	p("  return o.proto")
	p("}")

	return nil
}
