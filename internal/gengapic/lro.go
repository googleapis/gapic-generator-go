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
	"sort"

	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"google.golang.org/protobuf/types/descriptorpb"
)

func (g *generator) lroCall(servName string, m *descriptorpb.MethodDescriptorProto) error {
	inType := g.descInfo.Type[m.GetInputType()]
	outType := g.descInfo.Type[m.GetOutputType()]

	inSpec, err := g.descInfo.ImportSpec(inType)
	if err != nil {
		return err
	}

	outSpec, err := g.descInfo.ImportSpec(outType)
	if err != nil {
		return err
	}

	lroType := lroTypeName(m)
	p := g.printf

	lowcaseServName := lowerFirst(servName + "GRPCClient")

	p("func (c *%s) %s(ctx context.Context, req *%s.%s, opts ...gax.CallOption) (*%s, error) {",
		lowcaseServName, m.GetName(), inSpec.Name, inType.GetName(), lroType)

	g.insertRequestHeaders(m, grpc)
	g.appendCallOpts(m)

	p("  var resp *%s.%s", outSpec.Name, outType.GetName())
	p("  err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {")
	p("    var err error")
	p("    resp, err = %s", g.grpcStubCall(m))
	p("    return err")
	p("  }, opts...)")
	p("  if err != nil {")
	p("    return nil, err")
	p("  }")
	p("  return &%s{", lroType)
	p("    lro: longrunning.InternalNewOperation(*c.LROClient, resp),")
	p("  }, nil")

	p("}")
	p("")

	g.imports[pbinfo.ImportSpec{Path: "cloud.google.com/go/longrunning"}] = true
	g.imports[inSpec] = true
	g.imports[outSpec] = true
	return nil
}

// genOperationBuilders generates the name-based builder methods for each
// known RPC-specific wrapper type used in the service that the serice client
// needs to expose. These builders will construct an operation wrapper given
// only the operation name.
func (g *generator) genOperationBuilders(s *descriptorpb.ServiceDescriptorProto, servName string) error {
	methods := s.GetMethod()
	sort.Slice(methods, func(i, j int) bool {
		return methods[i].GetName() < methods[j].GetName()
	})
	for _, m := range methods {
		if _, ok := g.aux.methodToWrapper[m]; !ok {
			continue
		}
		if err := g.genOperationBuilder(servName, m); err != nil {
			return err
		}
	}
	return nil
}

// genOperationBuilder generates the code for the builder method specific to a
// given method's operation wrapper.
func (g *generator) genOperationBuilder(servName string, m *descriptorpb.MethodDescriptorProto) error {
	protoPkg := g.descInfo.ParentFile[m].GetPackage()
	ow := g.aux.methodToWrapper[m]
	p := g.printf

	// LRO from name
	{
		for _, t := range g.opts.transports {
			p("// %[1]s returns a new %[1]s from a given name.", ow.name)
			p("// The name must be that of a previously created %s, possibly from a different process.", ow.name)

			switch t {
			case grpc:
				receiver := lowcaseGRPCClientName(servName)
				p("func (c *%s) %[2]s(name string) *%[2]s {", receiver, ow.name)
				p("  return &%s{", ow.name)
				p("    lro: longrunning.InternalNewOperation(*c.LROClient, &longrunningpb.Operation{Name: name}),")
				p("  }")
				p("}")
				p("")
			case rest:
				receiver := lowcaseRestClientName(servName)
				override := g.getOperationPathOverride(protoPkg)
				p("func (c *%s) %[2]s(name string) *%[2]s {", receiver, ow.name)
				p("  override := fmt.Sprintf(%q, name)", override)
				p("  return &%s{", ow.name)
				p("    lro: longrunning.InternalNewOperation(*c.LROClient, &longrunningpb.Operation{Name: name}),")
				p("    pollPath: override,")
				p("  }")
				p("}")
				p("")
			}
		}
		g.imports[pbinfo.ImportSpec{Name: "longrunningpb", Path: "cloud.google.com/go/longrunning/autogen/longrunningpb"}] = true
		g.imports[pbinfo.ImportSpec{Path: "cloud.google.com/go/longrunning"}] = true
	}

	return nil
}
