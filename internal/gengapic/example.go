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
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/gapic-generator-go/internal/errors"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"google.golang.org/genproto/googleapis/longrunning"
)

func (g *generator) genExampleFile(serv *descriptor.ServiceDescriptorProto) error {
	pkgName := g.opts.pkgName
	servName := pbinfo.ReduceServName(serv.GetName(), pkgName)

	g.exampleClientFactory(pkgName, servName)

	methods := append(serv.GetMethod(), g.getMixinMethods()...)

	for _, m := range methods {
		if err := g.exampleMethod(pkgName, servName, m); err != nil {
			return err
		}
	}
	return nil
}

func (g *generator) exampleClientFactory(pkgName, servName string) {
	p := g.printf
	for _, t := range g.opts.transports {
		s := servName
		if t == rest {
			s += "REST"
		}

		p("func ExampleNew%sClient() {", s)
		g.exampleInitClient(pkgName, s)
		p("")
		p("  // TODO: Use client.")
		p("  _ = c")
		p("}")
		p("")
	}

	g.imports[pbinfo.ImportSpec{Path: "context"}] = true
}

func (g *generator) exampleInitClient(pkgName, servName string) {
	p := g.printf

	p("ctx := context.Background()")
	p("c, err := %s.New%sClient(ctx)", pkgName, servName)
	p("if err != nil {")
	p("  // TODO: Handle error.")
	p("}")
	p("defer c.Close()")
}

func (g *generator) exampleMethod(pkgName, servName string, m *descriptor.MethodDescriptorProto) error {
	if m.GetClientStreaming() != m.GetServerStreaming() {
		// TODO(pongad): implement this correctly.
		return nil
	}

	p := g.printf

	p("func Example%sClient_%s() {", servName, m.GetName())
	g.exampleMethodBody(pkgName, servName, m)

	p("}")
	p("")
	return nil
}

func (g *generator) exampleMethodBody(pkgName, servName string, m *descriptor.MethodDescriptorProto) error {
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
	// Pick the first transport for simplicity. We don't need examples
	// of each method for both transports when they have the same surface.
	t := g.opts.transports[0]
	s := servName
	if t == rest {
		s += "REST"
	}
	g.exampleInitClient(pkgName, s)

	if !m.GetClientStreaming() && !m.GetServerStreaming() {
		p("")
		p("req := &%s.%s{", inSpec.Name, inType.GetName())
		p("  // TODO: Fill request struct fields.")
		p("}")
	}

	pf, _, err := g.getPagingFields(m)
	if err != nil {
		return err
	}
	if pf != nil {
		g.examplePagingCall(m)
	} else if g.isLRO(m) {
		g.exampleLROCall(m)
	} else if *m.OutputType == emptyType {
		g.exampleEmptyCall(m)
	} else if m.GetClientStreaming() && m.GetServerStreaming() {
		g.exampleBidiCall(m, inType, inSpec)
	} else {
		g.exampleUnaryCall(m)
	}

	return nil
}

func (g *generator) exampleLROCall(m *descriptor.MethodDescriptorProto) {
	p := g.printf
	retVars := "resp, err :="

	// if response_type is google.protobuf.Empty, don't generate a "resp" var
	eLRO, err := proto.GetExtension(m.Options, longrunning.E_OperationInfo)
	if err == nil {
		opInfo := eLRO.(*longrunning.OperationInfo)
		if opInfo.GetResponseType() == emptyValue {
			// no new variables when this is used
			// therefore don't attempt to delcare it
			retVars = "err ="
		}
	}

	p("op, err := c.%s(ctx, req)", *m.Name)
	p("if err != nil {")
	p("  // TODO: Handle error.")
	p("}")
	p("")

	p("%s op.Wait(ctx)", retVars)
	p("if err != nil {")
	p("  // TODO: Handle error.")
	p("}")
	// generate response handling snippet
	if strings.Contains(retVars, "resp") {
		p("// TODO: Use resp.")
		p("_ = resp")
	}
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
