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

	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"google.golang.org/protobuf/types/descriptorpb"
)

// iterType describes iterators used by paging RPCs.
type iterType struct {
	iterTypeName, elemTypeName, mapValueTypeName string

	// If the elem type is a message, elemImports contains pbinfo.ImportSpec for the type.
	// Otherwise, len(elemImports)==0.
	elemImports []pbinfo.ImportSpec
}

// iterTypeOf deduces iterType from a field to be iterated over.
// elemField should be the "resource" of a paginating RPC.
// TODO(dovs): augment with paged map iterators
func (g *generator) iterTypeOf(elemField *descriptorpb.FieldDescriptorProto) (*iterType, error) {
	var pt iterType

	switch t := *elemField.Type; {
	case t == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
		eType := g.descInfo.Type[elemField.GetTypeName()]

		imp, err := g.descInfo.ImportSpec(eType)
		if err != nil {
			return &iterType{}, err
		}

		// Prepend parent Message name for nested Messages
		// to match the generated Go type name.
		typeName := g.nestedName(eType)

		eMsg, ok := eType.(*descriptorpb.DescriptorProto)
		if !ok {
			return nil, fmt.Errorf("cannot find message type %q, malformed descriptor", eType)
		}

		// Most repeated fields are not maps, so handle maps separately
		// and override these defaults.
		pt.elemTypeName = fmt.Sprintf("*%s.%s", imp.Name, typeName)
		pt.iterTypeName = typeName + "Iterator"

		if eMsg.GetOptions().GetMapEntry() {
			var valueField *descriptorpb.FieldDescriptorProto
			for _, f := range eMsg.GetField() {
				if f.GetName() == "value" {
					valueField = f
					break
				}
			}
			if valueField == nil {
				return nil, fmt.Errorf("unusual map entry message: %q", eMsg)
			}

			// The most common case is mapping to messages,
			// but check in case it's a primitive.
			if valueField.GetType() == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE {
				vType := g.descInfo.Type[valueField.GetTypeName()]
				n, imp, err := g.descInfo.NameSpec(vType)
				if err != nil {
					return nil, err
				}

				pt.mapValueTypeName = fmt.Sprintf("*%s.%s", imp.Name, n)
				pt.elemTypeName = fmt.Sprintf("%sPair", n)
			} else {
				pt.mapValueTypeName = pbinfo.GoTypeForPrim[valueField.GetType()]
				pt.elemTypeName = fmt.Sprintf("%sPair", upperFirst(pt.mapValueTypeName))
			}
			pt.iterTypeName = pt.elemTypeName + "Iterator"
		}

		pt.elemImports = []pbinfo.ImportSpec{imp}

	case t == descriptorpb.FieldDescriptorProto_TYPE_ENUM:
		log.Panic("iterating enum not supported yet")

	case t == descriptorpb.FieldDescriptorProto_TYPE_BYTES:
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

// getPagingFields reports the "resource field" to be iterated over by paginating method m
// and the "num elements" field that tells the server the maximum number of elements to return per page.
// Makes particular allowance for diregapic idioms: maps can be paginated over,
// and either 'page_size' XOR 'max_results' are allowable fields in the request.
func (g *generator) getPagingFields(m *descriptorpb.MethodDescriptorProto) (repeatedField, pageSizeField *descriptorpb.FieldDescriptorProto, e error) {
	// TODO: remove this once the next version of the Talent API is published.
	//
	// This is a workaround to disable auto-pagination for specifc RPCs in
	// Talent v4beta1. The API team will make their API non-conforming in the
	// next version.
	//
	// This should not be done for any other API.
	if g.descInfo.ParentFile[m].GetPackage() == "google.cloud.talent.v4beta1" &&
		(m.GetName() == "SearchProfiles" || m.GetName() == "SearchJobs") {
		return nil, nil, nil
	}

	if m.GetClientStreaming() || m.GetServerStreaming() {
		return nil, nil, nil
	}

	inType := g.descInfo.Type[m.GetInputType()]
	if inType == nil {
		return nil, nil, fmt.Errorf("expected %q to be message type, found %T", m.GetInputType(), inType)
	}
	inMsg, ok := inType.(*descriptorpb.DescriptorProto)
	if !ok {
		return nil, nil, fmt.Errorf("cannot find message type %q, malformed descriptor", m.GetInputType())
	}

	outType := g.descInfo.Type[m.GetOutputType()]
	if outType == nil {
		return nil, nil, fmt.Errorf("expected %q to be message type, found %T", m.GetOutputType(), outType)
	}
	outMsg, ok := outType.(*descriptorpb.DescriptorProto)
	if !ok {
		return nil, nil, fmt.Errorf("cannot find message type %q, malformed descriptor", m.GetOutputType())
	}

	hasPageToken := false
	for _, f := range inMsg.GetField() {
		isInt32 := f.GetType() == descriptorpb.FieldDescriptorProto_TYPE_INT32 || f.GetType() == descriptorpb.FieldDescriptorProto_TYPE_UINT32
		if (f.GetName() == "page_size" || f.GetName() == "max_results") && isInt32 {
			if pageSizeField != nil {
				return nil, nil, fmt.Errorf("found both page_size and max_results fields in message %q", m.GetInputType())
			}
			pageSizeField = f
			continue
		}

		hasPageToken = hasPageToken || (f.GetName() == "page_token" && f.GetType() == descriptorpb.FieldDescriptorProto_TYPE_STRING)
	}

	if !hasPageToken || pageSizeField == nil {
		// Not an error, just not paginated
		return nil, nil, nil
	}

	hasNextPageToken := false
	for _, f := range outMsg.GetField() {
		if f.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED {

			if repeatedField != nil {
				// Multiple repeated fields are tacitly okay as long as the
				// first listed repeated field has the lowest field number.
				// In this case, subsequent repeated fields are ignored.
				// See https://aip.dev/4233 for details.
				if repeatedField.GetNumber() > f.GetNumber() {
					return nil, nil, fmt.Errorf("found multiple repeated or map fields in message %q", m.GetOutputType())
				}
				// We want the _first_ repeated field to be the one paged over.
				continue
			}
			repeatedField = f
		}

		hasNextPageToken = hasNextPageToken || (f.GetName() == "next_page_token" && f.GetType() == descriptorpb.FieldDescriptorProto_TYPE_STRING)
	}

	if !hasNextPageToken || repeatedField == nil {
		return nil, nil, nil
	}

	return repeatedField, pageSizeField, nil
}

func (g *generator) maybeSortMapPage(elemField *descriptorpb.FieldDescriptorProto, pt *iterType) string {
	p := g.printf

	repeatedField, elems := fmt.Sprintf("resp%s", fieldGetter(elemField.GetName())), ""
	// Most paged methods have a normal repeated field and not a map, so use that as a default.
	elems = repeatedField
	if pt.mapValueTypeName != "" {
		elems = "elems"
		p("")
		p("    elems := make([]%s, 0, len(%s))", pt.elemTypeName, repeatedField)
		p("    for k, v := range %s {", repeatedField)
		p("        elems = append(elems, %s{k, v})", pt.elemTypeName)
		p("    }")
		p("    sort.Slice(elems, func(i, j int) bool { return elems[i].Key < elems[j].Key } )")
		p("")
		g.imports[pbinfo.ImportSpec{Path: "sort"}] = true
	}
	return elems
}

func (g *generator) makeFetchAndIterUpdate(pageSizeFieldName string) {
	p := g.printf

	p("fetch := func(pageSize int, pageToken string) (string, error) {")
	p("  items, nextPageToken, err := it.InternalFetch(pageSize, pageToken)")
	p("  if err != nil {")
	p(`    return "", err`)
	p("  }")
	p("  it.items = append(it.items, items...)")
	p("  return nextPageToken, nil")
	p("}")
	p("")
	p("it.pageInfo, it.nextFunc = iterator.NewPageInfo(fetch, it.bufLen, it.takeBuf)")
	p("it.pageInfo.MaxSize = int(req.Get%s())", pageSizeFieldName)
	p("it.pageInfo.Token = req.GetPageToken()")
	p("")
	p("return it")
}

func (g *generator) internalFetchSetup(outType *descriptorpb.DescriptorProto, outSpec pbinfo.ImportSpec, tok, pageSizeFieldName, max, ps string) {
	p := g.printf

	p("  resp := &%s.%s{}", outSpec.Name, outType.GetName())
	p(`  if pageToken != "" {`)
	p("    req.PageToken = %s", tok)
	p("  }")
	p("  if pageSize > math.MaxInt32 {")
	p("    req.%s = %s", pageSizeFieldName, max)
	p("  } else if pageSize != 0 {")
	p("    req.%s = %s", pageSizeFieldName, ps)
	p("  }")
}

func (g *generator) pagingCall(servName string, m *descriptorpb.MethodDescriptorProto, elemField, pageSize *descriptorpb.FieldDescriptorProto, pt *iterType) error {
	inType := g.descInfo.Type[m.GetInputType()].(*descriptorpb.DescriptorProto)
	outType := g.descInfo.Type[m.GetOutputType()].(*descriptorpb.DescriptorProto)

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

	g.insertRequestHeaders(m, grpc)
	g.appendCallOpts(m)
	pageSizeFieldName := snakeToCamel(pageSize.GetName())
	p("it := &%s{}", pt.iterTypeName)
	p("req = proto.Clone(req).(*%s.%s)", inSpec.Name, inType.GetName())
	p("it.InternalFetch = func(pageSize int, pageToken string) ([]%s, string, error) {", pt.elemTypeName)
	g.internalFetchSetup(outType, outSpec, tok, pageSizeFieldName, max, ps)
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
	elems := g.maybeSortMapPage(elemField, pt)
	p("  return %s, resp.GetNextPageToken(), nil", elems)
	p("}")
	g.makeFetchAndIterUpdate(pageSizeFieldName)
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
	if pt.mapValueTypeName != "" {
		p("// %s is a holder type for string/%s map entries", pt.elemTypeName, pt.mapValueTypeName)
		p("type %s struct {", pt.elemTypeName)
		p("  Key string")
		p("  Value %s", pt.mapValueTypeName)
		p("}")
	}

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

	g.imports[pbinfo.ImportSpec{Path: "google.golang.org/api/iterator"}] = true
	for _, spec := range pt.elemImports {
		g.imports[spec] = true
	}
}
