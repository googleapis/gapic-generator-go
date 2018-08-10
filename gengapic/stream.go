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

import "github.com/golang/protobuf/protoc-gen-go/descriptor"

func (g *generator) bidiCall(servName string, s *descriptor.ServiceDescriptorProto, m *descriptor.MethodDescriptorProto) error {
	p := g.printf

	servSpec, err := g.importSpec(s)
	if err != nil {
		return err
	}
	g.imports[servSpec] = true

	p("func (c *%sClient) %s(ctx context.Context, opts ...gax.CallOption) (%s.%s_%sClient, error) {",
		servName, m.GetName(), servSpec.name, s.GetName(), m.GetName())
	g.insertMetadata()
	g.appendCallOpts(m)
	p("  var resp %s.%s_%sClient", servSpec.name, s.GetName(), m.GetName())

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
