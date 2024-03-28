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
	"fmt"

	"google.golang.org/protobuf/types/descriptorpb"
)

// Used for both bidi and client streaming.
func (g *generator) noRequestStreamCall(servName string, s *descriptorpb.ServiceDescriptorProto, m *descriptorpb.MethodDescriptorProto) error {
	p := g.printf

	servSpec, err := g.descInfo.ImportSpec(s)
	if err != nil {
		return err
	}
	g.imports[servSpec] = true

	// We DON'T want to export the transport layers.
	lowcaseServName := lowerFirst(servName + "GRPCClient")

	retTyp := fmt.Sprintf("%s.%s_%sClient", servSpec.Name, s.GetName(), m.GetName())
	p("func (c *%s) %s(ctx context.Context, opts ...gax.CallOption) (%s, error) {",
		lowcaseServName, m.GetName(), retTyp)
	g.insertRequestHeaders(nil, grpc)
	p("  var resp %s", retTyp)

	g.appendCallOpts(m)

	p("  err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {")
	p("    var err error")
	p("    resp, err = c.%s.%s(ctx, settings.GRPC...)", grpcClientField(servName), m.GetName())
	p("    return err")
	p("  }, opts...)")
	p("  if err != nil {")
	p("    return nil, err")
	p("  }")
	p("  return resp, nil")
	p("}")
	p("")
	return nil
}

func (g *generator) serverStreamCall(servName string, s *descriptorpb.ServiceDescriptorProto, m *descriptorpb.MethodDescriptorProto) error {
	inType := g.descInfo.Type[*m.InputType]

	inSpec, err := g.descInfo.ImportSpec(inType)
	if err != nil {
		return err
	}
	g.imports[inSpec] = true

	servSpec, err := g.descInfo.ImportSpec(s)
	if err != nil {
		return err
	}
	g.imports[servSpec] = true

	p := g.printf
	lowcaseServName := lowerFirst(servName + "GRPCClient")

	retTyp := fmt.Sprintf("%s.%s_%sClient", servSpec.Name, s.GetName(), m.GetName())
	p("func (c *%s) %s(ctx context.Context, req *%s.%s, opts ...gax.CallOption) (%s, error) {",
		lowcaseServName, m.GetName(), inSpec.Name, inType.GetName(), retTyp)

	g.insertRequestHeaders(m, grpc)
	g.appendCallOpts(m)

	p("  var resp %s", retTyp)
	p("err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {")
	p("  var err error")
	p("  resp, err = %s", g.grpcStubCall(m))
	p("  return err")
	p("}, opts...)")
	p("if err != nil {")
	p("  return nil, err")
	p("}")
	p("return resp, nil")

	p("}")
	p("")

	return nil
}
