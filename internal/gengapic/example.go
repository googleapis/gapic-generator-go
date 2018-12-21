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
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/gapic-generator-go/internal/errors"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
)

func (g *generator) genExampleFile(serv *descriptor.ServiceDescriptorProto, pkgName string) error {
	servName := pbinfo.ReduceServName(*serv.Name, pkgName)
	p := g.printf

	p("func ExampleNew%sClient() {", servName)
	g.exampleInitClient(pkgName, servName)
	p("  // TODO: Use client.")
	p("  _ = c")
	p("}")
	p("")
	g.imports[pbinfo.ImportSpec{Path: "context"}] = true

	for _, m := range serv.Method {
		if err := g.exampleMethod(pkgName, servName, m); err != nil {
			return err
		}
	}
	return nil
}

func (g *generator) exampleInitClient(pkgName, servName string) {
	p := g.printf

	p("ctx := context.Background()")
	p("c, err := %s.New%sClient(ctx)", pkgName, servName)
	p("if err != nil {")
	p("  // TODO: Handle error.")
	p("}")
}

func (g *generator) exampleMethod(pkgName, servName string, m *descriptor.MethodDescriptorProto) error {
	if m.GetClientStreaming() != m.GetServerStreaming() {
		// TODO(pongad): implement this correctly.
		return nil
	}

	p := g.printf

	inType := g.descInfo.Type[m.GetInputType()]
	if inType == nil {
		return errors.E(nil, "cannot find type %q, malformed descriptor?", m.GetInputType())
	}

	inSpec, err := g.descInfo.ImportSpec(inType)
	if err != nil {
		return err
	}

	g.imports[inSpec] = true

	p("func Example%sClient_%s() {", servName, m.GetName())
	g.exampleInitClient(pkgName, servName)

	if !m.GetClientStreaming() && !m.GetServerStreaming() {
		p("")
		p("req := &%s.%s{", inSpec.Name, inType.GetName())
		p("  // TODO: Fill request struct fields.")
		p("}")
	}

	if pf, err := g.pagingField(m); err != nil {
		return err
	} else if pf != nil {
		g.examplePagingCall(m)
	} else if *m.OutputType == lroType {
		g.exampleLROCall(m)
	} else if *m.OutputType == emptyType {
		g.exampleEmptyCall(m)
	} else if m.GetClientStreaming() && m.GetServerStreaming() {
		g.exampleBidiCall(m, inType, inSpec)
	} else {
		g.exampleUnaryCall(m)
	}

	p("}")
	p("")
	return nil
}

func (g *generator) exampleLROCall(m *descriptor.MethodDescriptorProto) {
	p := g.printf

	p("op, err := c.%s(ctx, req)", *m.Name)
	p("if err != nil {")
	p("  // TODO: Handle error.")
	p("}")
	p("")

	p("resp, err := op.Wait(ctx)")
	p("if err != nil {")
	p("  // TODO: Handle error.")
	p("}")
	p("// TODO: Use resp.")
	p("_ = resp")
}

func (g *generator) exampleUnaryCall(m *descriptor.MethodDescriptorProto) {
	p := g.printf

	p("resp, err := c.%s(ctx, req)", *m.Name)
	p("if err != nil {")
	p("  // TODO: Handle error.")
	p("}")
	p("// TODO: Use resp.")
	p("_ = resp")
}

func (g *generator) exampleEmptyCall(m *descriptor.MethodDescriptorProto) {
	p := g.printf

	p("err = c.%s(ctx, req)", *m.Name)
	p("if err != nil {")
	p("  // TODO: Handle error.")
	p("}")
}

func (g *generator) examplePagingCall(m *descriptor.MethodDescriptorProto) {
	p := g.printf

	p("it := c.%s(ctx, req)", *m.Name)
	p("for {")
	p("  resp, err := it.Next()")
	p("  if err == iterator.Done {")
	p("    break")
	p("  }")
	p("  if err != nil {")
	p("    // TODO: Handle error.")
	p("  }")
	p("  // TODO: Use resp.")
	p("  _ = resp")
	p("}")

	g.imports[pbinfo.ImportSpec{Path: "google.golang.org/api/iterator"}] = true
}

func (g *generator) exampleBidiCall(m *descriptor.MethodDescriptorProto, inType pbinfo.ProtoType, inSpec pbinfo.ImportSpec) {
	p := g.printf

	p("stream, err := c.%s(ctx)", m.GetName())
	p("if err != nil {")
	p("  // TODO: Handle error.")
	p("}")

	p("go func() {")
	p("  reqs := []*%s.%s{", inSpec.Name, inType.GetName())
	p("    // TODO: Create requests.")
	p("  }")
	p("  for _, req := range reqs {")
	p("    if err := stream.Send(req); err != nil {")
	p("            // TODO: Handle error.")
	p("    }")
	p("  }")
	p("  stream.CloseSend()")
	p("}()")

	p("for {")
	p("  resp, err := stream.Recv()")
	p("  if err == io.EOF {")
	p("    break")
	p("  }")
	p("  if err != nil {")
	p("    // TODO: handle error.")
	p("  }")
	p("  // TODO: Use resp.")
	p("  _ = resp")
	p("}")

	g.imports[pbinfo.ImportSpec{Path: "io"}] = true
}
