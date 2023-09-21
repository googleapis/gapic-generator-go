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
	"strings"

	longrunning "cloud.google.com/go/longrunning/autogen/longrunningpb"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
)

func (g *generator) lroCall(servName string, m *descriptor.MethodDescriptorProto) error {
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

func (g *generator) lroType(servName string, serv *descriptor.ServiceDescriptorProto, m *descriptor.MethodDescriptorProto) error {
	protoPkg := g.descInfo.ParentFile[serv].GetPackage()
	mFQN := fmt.Sprintf("%s.%s.%s", protoPkg, serv.GetName(), m.GetName())
	lroType := lroTypeName(m)
	p := g.printf
	hasREST := containsTransport(g.opts.transports, rest)

	eLRO := proto.GetExtension(m.Options, longrunning.E_OperationInfo)
	opInfo := eLRO.(*longrunning.OperationInfo)
	fullName := opInfo.GetResponseType()
	if fullName == "" {
		return fmt.Errorf("rpc %q has google.longrunning.operation_info but is missing option google.longrunning.operation_info.response_type", mFQN)
	}

	var respType string
	{
		// eLRO.ResponseType is either fully-qualified or top-level in the same package as the method.
		//
		// TODO(ndietz) this won't work with nested message types in the same package;
		// migrating to protoreflect will help remove from semantic meaning in the names.
		if strings.IndexByte(fullName, '.') < 0 {
			fullName = protoPkg + "." + fullName
		}

		// When we build a map[name]Type in pbinfo, we prefix names with '.' to signify that they are fully qualified.
		// The string in ResponseType does not have the prefix, so we add it.
		fullName = "." + fullName

		typ := g.descInfo.Type[fullName]
		name, respSpec, err := g.descInfo.NameSpec(typ)
		if err != nil {
			return fmt.Errorf("unable to resolve google.longrunning.operation_info.response_type value %q in rpc %q", opInfo.GetResponseType(), mFQN)
		}

		if fullName != emptyType {
			g.imports[respSpec] = true

			respType = fmt.Sprintf("%s.%s", respSpec.Name, name)
		}
	}

	hasMeta := opInfo.GetMetadataType() != ""
	var metaType string
	if hasMeta {
		fullName := opInfo.GetMetadataType()
		// TODO(ndietz) this won't work with nested message types in the same package;
		// migrating to protoreflect will help remove from semantic meaning in the names.
		if strings.IndexByte(fullName, '.') < 0 {
			fullName = protoPkg + "." + fullName
		}
		fullName = "." + fullName

		typ := g.descInfo.Type[fullName]
		name, meta, err := g.descInfo.NameSpec(typ)
		if err != nil {
			return fmt.Errorf("unable to resolve google.longrunning.operation_info.metadata_type value %q in rpc %q", opInfo.GetMetadataType(), mFQN)
		}
		g.imports[meta] = true

		metaType = fmt.Sprintf("%s.%s", meta.Name, name)
	}

	// Type definition
	{
		p("// %s manages a long-running operation from %s.", lroType, m.GetName())
		p("type %s struct {", lroType)
		p("  lro *longrunning.Operation")
		if hasREST {
			p("  pollPath string")
		}
		p("}")
		p("")
	}

	// LRO from name
	{
		for _, t := range g.opts.transports {
			p("// %[1]s returns a new %[1]s from a given name.", lroType)
			p("// The name must be that of a previously created %s, possibly from a different process.", lroType)

			switch t {
			case grpc:
				receiver := lowcaseGRPCClientName(servName)
				p("func (c *%s) %[2]s(name string) *%[2]s {", receiver, lroType)
				p("  return &%s{", lroType)
				p("    lro: longrunning.InternalNewOperation(*c.LROClient, &longrunningpb.Operation{Name: name}),")
				p("  }")
				p("}")
				p("")
			case rest:
				receiver := lowcaseRestClientName(servName)
				override := g.getOperationPathOverride(protoPkg)
				p("func (c *%s) %[2]s(name string) *%[2]s {", receiver, lroType)
				p("  override := fmt.Sprintf(%q, name)", override)
				p("  return &%s{", lroType)
				p("    lro: longrunning.InternalNewOperation(*c.LROClient, &longrunningpb.Operation{Name: name}),")
				p("    pollPath: override,")
				p("  }")
				p("}")
				p("")
			}
		}
		g.imports[pbinfo.ImportSpec{Name: "longrunningpb", Path: "cloud.google.com/go/longrunning/autogen/longrunningpb"}] = true
	}

	// Wait
	{
		p("// Wait blocks until the long-running operation is completed, returning the response and any errors encountered.")
		p("//")
		p("// See documentation of Poll for error-handling information.")
		if opInfo.GetResponseType() == emptyValue {
			p("func (op *%s) Wait(ctx context.Context, opts ...gax.CallOption) error {", lroType)
			if hasREST {
				p("opts = append([]gax.CallOption{gax.WithPath(op.pollPath)}, opts...)")
			}
			p("  return op.lro.WaitWithInterval(ctx, nil, %s, opts...)", defaultPollMaxDelay)
		} else {
			p("func (op *%s) Wait(ctx context.Context, opts ...gax.CallOption) (*%s, error) {", lroType, respType)
			if hasREST {
				p("opts = append([]gax.CallOption{gax.WithPath(op.pollPath)}, opts...)")
			}
			p("  var resp %s", respType)
			p("  if err := op.lro.WaitWithInterval(ctx, &resp, %s, opts...); err != nil {", defaultPollMaxDelay)
			p("    return nil, err")
			p("  }")
			p("  return &resp, nil")
		}
		p("}")
		p("")

		g.imports[pbinfo.ImportSpec{Path: "time"}] = true
	}

	// Poll
	{
		p("// Poll fetches the latest state of the long-running operation.")
		p("//")
		if hasMeta {
			p("// Poll also fetches the latest metadata, which can be retrieved by Metadata.")
			p("//")
		}
		p("// If Poll fails, the error is returned and op is unmodified. If Poll succeeds and")
		p("// the operation has completed with failure, the error is returned and op.Done will return true.")
		p("// If Poll succeeds and the operation has completed successfully,")
		p("// op.Done will return true, and the response of the operation is returned.")
		p("// If Poll succeeds and the operation has not completed, the returned response and error are both nil.")
		if opInfo.GetResponseType() == emptyValue {
			p("func (op *%s) Poll(ctx context.Context, opts ...gax.CallOption) error {", lroType)
			if hasREST {
				p("opts = append([]gax.CallOption{gax.WithPath(op.pollPath)}, opts...)")
			}
			p("  return op.lro.Poll(ctx, nil, opts...)")
		} else {
			p("func (op *%s) Poll(ctx context.Context, opts ...gax.CallOption) (*%s, error) {", lroType, respType)
			if hasREST {
				p("opts = append([]gax.CallOption{gax.WithPath(op.pollPath)}, opts...)")
			}
			p("  var resp %s", respType)
			p("  if err := op.lro.Poll(ctx, &resp, opts...); err != nil {")
			p("    return nil, err")
			p("  }")
			p("  if !op.Done() {")
			p("    return nil, nil")
			p("  }")
			p("  return &resp, nil")
		}
		p("}")
		p("")
	}

	// Metadata
	if hasMeta {
		p("// Metadata returns metadata associated with the long-running operation.")
		p("// Metadata itself does not contact the server, but Poll does.")
		p("// To get the latest metadata, call this method after a successful call to Poll.")
		p("// If the metadata is not available, the returned metadata and error are both nil.")
		p("func (op *%s) Metadata() (*%s, error) {", lroType, metaType)
		p("  var meta %s", metaType)
		p("  if err := op.lro.Metadata(&meta); err == longrunning.ErrNoMetadata {")
		p("    return nil, nil")
		p("  } else if err != nil {")
		p("    return nil, err")
		p("  }")
		p("  return &meta, nil")
		p("}")
		p("")
	}

	// Done
	{
		p("// Done reports whether the long-running operation has completed.")
		p("func (op *%s) Done() bool {", lroType)
		p("return op.lro.Done()")
		p("}")
		p("")
	}

	// Name
	{
		p("// Name returns the name of the long-running operation.")
		p("// The name is assigned by the server and is unique within the service from which the operation is created.")
		p("func (op *%s) Name() string {", lroType)
		p("return op.lro.Name()")
		p("}")
		p("")
	}
	return nil
}

func lroTypeName(m *descriptor.MethodDescriptorProto) string {
	// This whole if block is a hack to workaround a operation handler namespace
	// collision. We should remove this in the future if the design is fixed for
	// the v1 api. This is for aiplatform.featureregistryservice.createfeature.
	if eHTTP, ok := proto.GetExtension(m.GetOptions(), annotations.E_Http).(*annotations.HttpRule); ok && eHTTP != nil && eHTTP.Pattern != nil {
		switch t := eHTTP.Pattern.(type) {
		case *annotations.HttpRule_Post:
			if t.Post == "/v1beta1/{parent=projects/*/locations/*/featureGroups/*}/features" {
				return m.GetName() + "RegistryOperation"
			}
		}
	}
	return m.GetName() + "Operation"
}
