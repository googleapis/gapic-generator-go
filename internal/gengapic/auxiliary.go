// Copyright 2023 Google LLC
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
	"path/filepath"
	"sort"
	"strings"

	longrunning "cloud.google.com/go/longrunning/autogen/longrunningpb"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

// auxTypes gathers details of types we need to generate along with the client
type auxTypes struct {
	// Map of RPC descriptor to assigned operationWrapper.
	methodToWrapper map[*descriptorpb.MethodDescriptorProto]operationWrapper

	opWrappers map[string]operationWrapper

	// "List" of iterator types. We use these to generate FooIterator returned by paging methods.
	// Since multiple methods can page over the same type, we dedupe by the name of the iterator,
	// which is in turn determined by the element type name.
	iters map[string]*iterType

	customOp *customOp
}

type operationWrapper struct {
	name                       string
	response, metadata         *descriptorpb.DescriptorProto
	responseName, metadataName protoreflect.FullName
}

func (g *generator) genAuxFile() error {
	if g.aux.customOp != nil {
		if err := g.customOperationType(); err != nil {
			return err
		}
		g.commit(filepath.Join(g.opts.outDir, "operations.go"), g.opts.pkgName)
		g.reset()
	}

	if err := g.genOperations(); err != nil {
		return err
	}

	if err := g.genIterators(); err != nil {
		return err
	}

	g.commit(filepath.Join(g.opts.outDir, "aux.go"), g.opts.pkgName)
	g.reset()

	return nil
}

func (g *generator) genOperations() error {
	// Sort operation wrappers-to-generate by type
	// name to avoid spurious regenerations created
	// by non-deterministic map traversal order.
	var wrappers []operationWrapper
	for _, ow := range g.aux.opWrappers {
		wrappers = append(wrappers, ow)
	}
	sort.Slice(wrappers, func(i, j int) bool {
		return wrappers[i].name < wrappers[j].name
	})
	for _, ow := range wrappers {
		if err := g.genOperationWrapperType(ow); err != nil {
			return err
		}
	}
	return nil
}

func (g *generator) genIterators() error {
	// Sort iterators to generate by type name to
	// avoid spurious regenerations created by
	// non-deterministic map traversal order.
	var iters []*iterType
	for _, iter := range g.aux.iters {
		iters = append(iters, iter)
	}
	sort.Slice(iters, func(i, j int) bool {
		return iters[i].iterTypeName < iters[j].iterTypeName
	})
	for _, iter := range iters {
		g.pagingIter(iter)
	}

	return nil
}

func (a auxTypes) wrapperExists(ow operationWrapper) (bool, error) {
	ew, exists := a.opWrappers[ow.name]
	if !exists {
		return false, nil
	}

	if ow.responseName != ew.responseName {
		return true, fmt.Errorf("duplicate operation wrapper types %q have mismatched response_types: %s v. %s", ow.name, ew.responseName, ow.responseName)
	}

	if ow.metadataName != ew.metadataName {
		return true, fmt.Errorf("duplicate operation wrapper types %q have mismatched metadata_types: %s v. %s", ow.name, ew.metadataName, ow.responseName)
	}

	return true, nil
}

func (g *generator) maybeAddOperationWrapper(m *descriptorpb.MethodDescriptorProto) error {
	if !proto.HasExtension(m.GetOptions(), longrunning.E_OperationInfo) {
		return fmt.Errorf("%s missing google.longrunning.operation_info", m.GetName())
	}

	protoPkg := protoreflect.FullName(g.descInfo.ParentFile[m].GetPackage())
	eLRO := proto.GetExtension(m.GetOptions(), longrunning.E_OperationInfo)
	opInfo := eLRO.(*longrunning.OperationInfo)
	ow := operationWrapper{
		name: lroTypeName(m),
	}

	// Response type resolution.
	{
		var respType protoreflect.FullName
		rawResp := opInfo.GetResponseType()
		if rawResp == "" {
			return fmt.Errorf("rpc %q has google.longrunning.operation_info but is missing option google.longrunning.operation_info.response_type", m.GetName())
		}
		if !strings.Contains(rawResp, ".") {
			respType = protoPkg.Append(protoreflect.Name(rawResp))
		} else {
			respType = protoreflect.FullName(rawResp)
		}

		// When we build a map[name]Type in pbinfo, we prefix names with '.' to signify that they are fully qualified.
		// The string in ResponseType does not have the prefix, so we add it.
		typ, ok := g.descInfo.Type["."+string(respType)]
		if !ok {
			return fmt.Errorf("unable to resolve google.longrunning.operation_info.response_type value %q in rpc %q", opInfo.GetResponseType(), m.GetName())
		}

		ow.response = typ.(*descriptorpb.DescriptorProto)
		ow.responseName = respType
	}

	// Metadata type resolution.
	{
		rawMeta := opInfo.GetMetadataType()
		if rawMeta == "" {
			return fmt.Errorf("rpc %q has google.longrunning.operation_info but is missing option google.longrunning.operation_info.metadata_type", m.GetName())
		}

		var metaType protoreflect.FullName
		if !strings.Contains(rawMeta, ".") {
			metaType = protoPkg.Append(protoreflect.Name(rawMeta))
		} else {
			metaType = protoreflect.FullName(rawMeta)
		}

		typ, ok := g.descInfo.Type["."+string(metaType)]

		if !ok {
			return fmt.Errorf("unable to resolve google.longrunning.operation_info.metadata_type value %q in rpc %q", opInfo.GetMetadataType(), m.GetName())
		}

		ow.metadata = typ.(*descriptorpb.DescriptorProto)
		ow.metadataName = metaType
	}

	if exists, err := g.aux.wrapperExists(ow); err != nil {
		return err
	} else if !exists {
		g.aux.opWrappers[ow.name] = ow
	}

	g.aux.methodToWrapper[m] = g.aux.opWrappers[ow.name]

	return nil
}

func (g *generator) genOperationWrapperType(ow operationWrapper) error {
	p := g.pt.Printf
	hasREST := containsTransport(g.opts.transports, rest)
	isEmpty := ow.responseName == emptyValue

	// Response Go type resolution.
	var respType string
	if !isEmpty {
		name, respSpec, err := g.descInfo.NameSpec(ow.response)
		if err != nil {
			return err
		}
		g.imports[respSpec] = true

		respType = fmt.Sprintf("%s.%s", respSpec.Name, name)
	}

	// Metadata Go type resolution.
	name, meta, err := g.descInfo.NameSpec(ow.metadata)
	if err != nil {
		return err
	}
	g.imports[meta] = true
	metaType := fmt.Sprintf("%s.%s", meta.Name, name)

	// Operation wrapper type definition
	{
		p("// %s manages a long-running operation from %s.", ow.name, strings.TrimSuffix(ow.name, "Operation"))
		p("type %s struct {", ow.name)
		p("  lro *longrunning.Operation")
		if hasREST {
			p("  pollPath string")
		}
		p("}")
		p("")
		g.imports[pbinfo.ImportSpec{Path: "cloud.google.com/go/longrunning"}] = true
	}

	// Wait
	{
		p("// Wait blocks until the long-running operation is completed, returning the response and any errors encountered.")
		p("//")
		p("// See documentation of Poll for error-handling information.")
		if isEmpty {
			p("func (op *%s) Wait(ctx context.Context, opts ...gax.CallOption) error {", ow.name)
			if hasREST {
				p("opts = append([]gax.CallOption{gax.WithPath(op.pollPath)}, opts...)")
			}
			p("  return op.lro.WaitWithInterval(ctx, nil, %s, opts...)", defaultPollMaxDelay)
		} else {
			p("func (op *%s) Wait(ctx context.Context, opts ...gax.CallOption) (*%s, error) {", ow.name, respType)
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

		g.imports[pbinfo.ImportSpec{Path: "context"}] = true
		g.imports[pbinfo.ImportSpec{Path: "time"}] = true
		g.imports[pbinfo.ImportSpec{Name: "gax", Path: "github.com/googleapis/gax-go/v2"}] = true
	}

	// Poll
	{
		p("// Poll fetches the latest state of the long-running operation.")
		p("//")
		p("// Poll also fetches the latest metadata, which can be retrieved by Metadata.")
		p("//")
		p("// If Poll fails, the error is returned and op is unmodified. If Poll succeeds and")
		p("// the operation has completed with failure, the error is returned and op.Done will return true.")
		p("// If Poll succeeds and the operation has completed successfully,")
		p("// op.Done will return true, and the response of the operation is returned.")
		p("// If Poll succeeds and the operation has not completed, the returned response and error are both nil.")
		if isEmpty {
			p("func (op *%s) Poll(ctx context.Context, opts ...gax.CallOption) error {", ow.name)
			if hasREST {
				p("opts = append([]gax.CallOption{gax.WithPath(op.pollPath)}, opts...)")
			}
			p("  return op.lro.Poll(ctx, nil, opts...)")
		} else {
			p("func (op *%s) Poll(ctx context.Context, opts ...gax.CallOption) (*%s, error) {", ow.name, respType)
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
	{
		p("// Metadata returns metadata associated with the long-running operation.")
		p("// Metadata itself does not contact the server, but Poll does.")
		p("// To get the latest metadata, call this method after a successful call to Poll.")
		p("// If the metadata is not available, the returned metadata and error are both nil.")
		p("func (op *%s) Metadata() (*%s, error) {", ow.name, metaType)
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
		p("func (op *%s) Done() bool {", ow.name)
		p("return op.lro.Done()")
		p("}")
		p("")
	}

	// Name
	{
		p("// Name returns the name of the long-running operation.")
		p("// The name is assigned by the server and is unique within the service from which the operation is created.")
		p("func (op *%s) Name() string {", ow.name)
		p("return op.lro.Name()")
		p("}")
		p("")
	}
	return nil
}

func lroTypeName(m *descriptorpb.MethodDescriptorProto) string {
	return m.GetName() + "Operation"
}
