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

// customOp represents a custom operation type for long running operations.
type customOp struct {
	message   *descriptor.DescriptorProto
	generated bool
}

// isCustomOp determines if the given method should return a custom operation wrapper.
func (g *generator) isCustomOp(m *descriptor.MethodDescriptorProto, info *httpInfo) bool {
	return g.opts.diregapic && // Generator in DIREGAPIC mode.
		g.aux.customOp != nil && // API Defines a custom operation.
		m.GetOutputType() == g.customOpProtoName() && // Method returns the custom operation.
		info.verb != "get" && // Method is not a GET (polling methods).
		m.GetName() != "Wait" // Method is not a Wait (uses POST).
}

// customOpProtoName builds the fully-qualified proto name for the custom
// operation message type.
func (g *generator) customOpProtoName() string {
	f := g.descInfo.ParentFile[g.aux.customOp.message]
	return fmt.Sprintf(".%s.%s", f.GetPackage(), g.aux.customOp.message.GetName())
}

// customOpPointerType builds a string containing the Go code for a pointer to
// the custom operation type.
func (g *generator) customOpPointerType() (string, error) {
	op := g.aux.customOp
	if op == nil {
		return "", nil
	}

	opName, imp, err := g.descInfo.NameSpec(op.message)
	if err != nil {
		return "", err
	}

	s := fmt.Sprintf("*%s.%s", imp.Name, opName)

	return s, nil
}

// customOpInit builds a string containing the Go code for initializing the
// operation wrapper type with the Go identifier for a variable that is the
// proto-defined operation type.
func (g *generator) customOpInit(p string) string {
	opName := g.aux.customOp.message.GetName()

	s := fmt.Sprintf("&%s{proto: %s}", opName, p)

	return s
}

// customOperationType generates the custom operation wrapper type using the
// generators current printer. This should only be called once per package.
func (g *generator) customOperationType() error {
	op := g.aux.customOp
	if op == nil {
		return nil
	}
	opName := op.message.GetName()

	ptyp, err := g.customOpPointerType()
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
