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
	"log"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/gapic-generator-go/internal/errors"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
)

// iterType describes iterators used by paging RPCs.
type iterType struct {
	iterTypeName, elemTypeName string

	// If the elem type is a message, elemImports contains pbinfo.ImportSpec for the type.
	// Otherwise, len(elemImports)==0.
	elemImports []pbinfo.ImportSpec

	generated bool
}

// iterTypeOf deduces iterType from a field to be iterated over.
// elemField should be the "resource" of a paginating RPC.
func (g *generator) iterTypeOf(elemField *descriptor.FieldDescriptorProto) (*iterType, error) {
	var pt iterType

	switch t := *elemField.Type; {
	case t == descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		eType := g.descInfo.Type[elemField.GetTypeName()]

		imp, err := g.descInfo.ImportSpec(eType)
		if err != nil {
			return &iterType{}, err
		}

		// Prepend parent Message name for nested Messages
		// to match the generated Go type name.
		typ := eType
		typeName := typ.GetName()
		for parent, ok := g.descInfo.ParentElement[typ]; ok; parent, ok = g.descInfo.ParentElement[typ] {
			typeName = fmt.Sprintf("%s_%s", parent.GetName(), typeName)
			typ = parent
		}

		pt.elemTypeName = fmt.Sprintf("*%s.%s", imp.Name, typeName)
		pt.iterTypeName = typeName + "Iterator"

		pt.elemImports = []pbinfo.ImportSpec{imp}

	case t == descriptor.FieldDescriptorProto_TYPE_ENUM:
		log.Panic("iterating enum not supported yet")

	case t == descriptor.FieldDescriptorProto_TYPE_BYTES:
		pt.elemTypeName = "[]byte"
		pt.iterTypeName = "BytesIterator"

	default:
		pType := pbinfo.GoTypeForPrim[t]
		if pType == "" {
			log.Panicf("unrecognized type: %v", t)
		}
		pt.elemTypeName = pType
		pt.iterTypeName = upperFirst(pt.elemTypeName) + "Iterator"
	}

	if iter, ok := g.aux.iters[pt.iterTypeName]; ok {
		return iter, nil
	}
	g.aux.iters[pt.iterTypeName] = &pt

	return &pt, nil
}

// pagingField reports the "resource field" to be iterated over by paginating method m.
// If the method is not a paging method, pagingField returns (nil, nil).
// If the method looks like a paging method, but the field cannot be determined, pagingField errors.
func (g *generator) pagingField(m *descriptor.MethodDescriptorProto) (*descriptor.FieldDescriptorProto, error) {
	var (
		hasSize, hasToken, hasNextToken bool
		elemFields                      []*descriptor.FieldDescriptorProto
	)

	// TODO: remove this once the next version of the Talent API is published.
	//
	// This is a workaround to disable auto-pagination for specifc RPCs in
	// Talent v4beta1. The API team will make their API non-conforming in the
	// next version.
	//
	// This should not be done for any other API.
	if g.descInfo.ParentFile[m].GetPackage() == "google.cloud.talent.v4beta1" &&
		(m.GetName() == "SearchProfiles" || m.GetName() == "SearchJobs") {
		return nil, nil
	}

	inType := g.descInfo.Type[m.GetInputType()]
	if inType == nil {
		return nil, errors.E(nil, "cannot find message type %q, malformed descriptor?", m.GetInputType())
	}
	inMsg, ok := inType.(*descriptor.DescriptorProto)
	if !ok {
		return nil, errors.E(nil, "expected %q to be message type, found %T", m.GetInputType(), inType)
	}

	outType := g.descInfo.Type[m.GetOutputType()]
	if outType == nil {
		return nil, errors.E(nil, "cannot find message type %q, malformed descriptor?", m.GetOutputType())
	}
	outMsg, ok := outType.(*descriptor.DescriptorProto)
	if !ok {
		return nil, errors.E(nil, "expected %q to be message type, found %T", m.GetOutputType(), outType)
	}

	for _, f := range inMsg.Field {
		if f.GetName() == "page_size" && f.GetType() == descriptor.FieldDescriptorProto_TYPE_INT32 {
			hasSize = true
		}
		if f.GetName() == "page_token" && f.GetType() == descriptor.FieldDescriptorProto_TYPE_STRING {
			hasToken = true
		}
	}
	for _, f := range outMsg.Field {
		if f.GetName() == "next_page_token" && f.GetType() == descriptor.FieldDescriptorProto_TYPE_STRING {
			hasNextToken = true
		}
		if f.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			// Ignore repeated fields that are actually synthetic MapEntry types.
			if eType, ok := g.descInfo.Type[f.GetTypeName()]; ok {
				if msg, ok := eType.(*descriptor.DescriptorProto); ok && msg.GetOptions().GetMapEntry() {
					continue
				}
			}
			elemFields = append(elemFields, f)
		}
	}

	// If any of the required pagination fields are missing or if the response
	// has no repeated fields, it is not a fully conforming paginated method.
	if !hasSize || !hasToken || !hasNextToken || len(elemFields) == 0 {
		return nil, nil
	}

	if len(elemFields) > 1 {
		// ensure that the first repeated field visited has the lowest field number
		// according to https://aip.dev/4233, and use that one.
		min := elemFields[0].GetNumber()
		for _, elem := range elemFields {
			if elem.GetNumber() < min {
				return nil, fmt.Errorf("%s: can't pick repeated field in %s based on aip.dev/4233", *m.Name, outType.GetName())
			}
		}
	}
	return elemFields[0], nil
}

func (g *generator) pagingCall(servName string, m *descriptor.MethodDescriptorProto, elemField *descriptor.FieldDescriptorProto, pt *iterType) error {
	inType := g.descInfo.Type[m.GetInputType()].(*descriptor.DescriptorProto)
	outType := g.descInfo.Type[m.GetOutputType()].(*descriptor.DescriptorProto)

	// We DON'T want to export the transport layers.
	lowcaseServName := lowerFirst(servName + "GRPCClient")

	inSpec, err := g.descInfo.ImportSpec(inType)
	if err != nil {
		return err
	}

	outSpec, err := g.descInfo.ImportSpec(outType)
	if err != nil {
		return err
	}

	max := "math.MaxInt32"
	ps := "int32(pageSize)"
	if isOptional(inType, "page_size") {
		max = fmt.Sprintf("proto.Int32(%s)", max)
		ps = fmt.Sprintf("proto.Int32(%s)", ps)
	}

	tok := "pageToken"
	if isOptional(inType, "page_token") {
		tok = fmt.Sprintf("proto.String(%s)", tok)
	}

	p := g.printf
	p("func (c *%s) %s(ctx context.Context, req *%s.%s, opts ...gax.CallOption) *%s {",
		lowcaseServName, *m.Name, inSpec.Name, inType.GetName(), pt.iterTypeName)

	err = g.insertMetadata(m)
	if err != nil {
		return err
	}
	g.appendCallOpts(m)

	p("it := &%s{}", pt.iterTypeName)
	p("req = proto.Clone(req).(*%s.%s)", inSpec.Name, inType.GetName())
	p("it.InternalFetch = func(pageSize int, pageToken string) ([]%s, string, error) {", pt.elemTypeName)
	p("  var resp *%s.%s", outSpec.Name, outType.GetName())
	p("  req.PageToken = %s", tok)
	p("  if pageSize > math.MaxInt32 {")
	p("    req.PageSize = %s", max)
	p("  } else {")
	p("    req.PageSize = %s", ps)
	p("  }")
	p("  err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {")
	p("    var err error")
	p("    resp, err = %s", g.grpcStubCall(m))
	p("    return err")
	p("  }, opts...)")
	p("  if err != nil {")
	p("    return nil, \"\", err")
	p("  }")
	p("")
	p("  it.Response = resp")
	p("  return resp.Get%s(), resp.GetNextPageToken(), nil", snakeToCamel(elemField.GetName()))
	p("}")

	p("fetch := func(pageSize int, pageToken string) (string, error) {")
	p("  items, nextPageToken, err := it.InternalFetch(pageSize, pageToken)")
	p("  if err != nil {")
	p("    return \"\", err")
	p("  }")
	p("  it.items = append(it.items, items...)")
	p("  return nextPageToken, nil")
	p("}")

	p("it.pageInfo, it.nextFunc = iterator.NewPageInfo(fetch, it.bufLen, it.takeBuf)")
	p("it.pageInfo.MaxSize = int(req.GetPageSize())")
	p("it.pageInfo.Token = req.GetPageToken()")
	p("return it")

	p("}")
	p("")

	g.imports[pbinfo.ImportSpec{Path: "google.golang.org/protobuf/proto"}] = true
	g.imports[pbinfo.ImportSpec{Path: "google.golang.org/api/iterator"}] = true
	g.imports[inSpec] = true
	g.imports[outSpec] = true
	for _, spec := range pt.elemImports {
		g.imports[spec] = true
	}
	return nil
}

func (g *generator) pagingIter(pt *iterType) {
	p := g.printf

	p("// %s manages a stream of %s.", pt.iterTypeName, pt.elemTypeName)
	p("type %s struct {", pt.iterTypeName)
	p("  items    []%s", pt.elemTypeName)
	p("  pageInfo *iterator.PageInfo")
	p("  nextFunc func() error")
	p("")
	p("  // Response is the raw response for the current page.")
	p("  // It must be cast to the RPC response type.")
	p("  // Calling Next() or InternalFetch() updates this value.")
	p("  Response interface{}")
	p("")
	p("  // InternalFetch is for use by the Google Cloud Libraries only.")
	p("  // It is not part of the stable interface of this package.")
	p("  //")
	p("  // InternalFetch returns results from a single call to the underlying RPC.")
	p("  // The number of results is no greater than pageSize.")
	p("  // If there are no more results, nextPageToken is empty and err is nil.")
	p("  InternalFetch func(pageSize int, pageToken string) (results []%s, nextPageToken string, err error)", pt.elemTypeName)
	p("}")
	p("")

	p("// PageInfo supports pagination. See the google.golang.org/api/iterator package for details.")
	p("func (it *%s) PageInfo() *iterator.PageInfo {", pt.iterTypeName)
	p("  return it.pageInfo")
	p("}")
	p("")

	p("// Next returns the next result. Its second return value is iterator.Done if there are no more")
	p("// results. Once Next returns Done, all subsequent calls will return Done.")
	p("func (it *%s) Next() (%s, error) {", pt.iterTypeName, pt.elemTypeName)
	p("  var item %s", pt.elemTypeName)
	p("  if err := it.nextFunc(); err != nil {")
	p("    return item, err")
	p("  }")
	p("  item = it.items[0]")
	p("  it.items = it.items[1:]")
	p("  return item, nil")
	p("}")
	p("")

	p("func (it *%s) bufLen() int {", pt.iterTypeName)
	p("  return len(it.items)")
	p("}")
	p("")

	p("func (it *%s) takeBuf() interface{} {", pt.iterTypeName)
	p("  b := it.items")
	p("  it.items = nil")
	p("  return b")
	p("}")
	p("")
}
