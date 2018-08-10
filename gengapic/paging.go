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

	"gapic-generator-go/internal/errors"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

var primitiveFieldToGoType = [...]string{
	descriptor.FieldDescriptorProto_TYPE_DOUBLE:   "float64",
	descriptor.FieldDescriptorProto_TYPE_FLOAT:    "float32",
	descriptor.FieldDescriptorProto_TYPE_INT64:    "int64",
	descriptor.FieldDescriptorProto_TYPE_UINT64:   "uint64",
	descriptor.FieldDescriptorProto_TYPE_INT32:    "int32",
	descriptor.FieldDescriptorProto_TYPE_FIXED64:  "uint64",
	descriptor.FieldDescriptorProto_TYPE_FIXED32:  "uint32",
	descriptor.FieldDescriptorProto_TYPE_BOOL:     "bool",
	descriptor.FieldDescriptorProto_TYPE_STRING:   "string",
	descriptor.FieldDescriptorProto_TYPE_BYTES:    "[]byte",
	descriptor.FieldDescriptorProto_TYPE_UINT32:   "uint32",
	descriptor.FieldDescriptorProto_TYPE_SFIXED32: "int32",
	descriptor.FieldDescriptorProto_TYPE_SFIXED64: "int64",
	descriptor.FieldDescriptorProto_TYPE_SINT32:   "int32",
	descriptor.FieldDescriptorProto_TYPE_SINT64:   "int64",
}

// iterType describes iterators used by paging RPCs.
type iterType struct {
	iterTypeName, elemTypeName string

	// If the elem type is a message, elemImports contains importSpec for the type.
	// Otherwise, len(elemImports)==0.
	elemImports []importSpec
}

// iterTypeOf deduces iterType from a field to be iterated over.
// elemField should be the "resource" of a paginating RPC.
func (g *generator) iterTypeOf(elemField *descriptor.FieldDescriptorProto) (iterType, error) {
	var pt iterType

	switch t := *elemField.Type; {
	case t == descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		eType := g.types[elemField.GetTypeName()]

		imp, err := g.importSpec(eType)
		if err != nil {
			return iterType{}, err
		}

		pt.elemTypeName = fmt.Sprintf("*%s.%s", imp.name, eType.GetName())
		pt.iterTypeName = *eType.Name + "Iterator"

		pt.elemImports = []importSpec{imp}

	case t == descriptor.FieldDescriptorProto_TYPE_ENUM:
		log.Panic("iterating enum not supported yet")

	case t == descriptor.FieldDescriptorProto_TYPE_BYTES:
		pt.elemTypeName = "[]byte"
		pt.iterTypeName = "BytesIterator"

	case t < 0 || int(t) >= len(primitiveFieldToGoType) || primitiveFieldToGoType[t] == "":
		log.Panicf("unrecognized type: %v", t)

	default:
		pt.elemTypeName = primitiveFieldToGoType[t]
		pt.iterTypeName = upperFirst(pt.elemTypeName) + "Iterator"
	}
	return pt, nil
}

// TODO(pongad): this will probably need to read from annotations later.

// pagingField reports the "resource field" to be iterated over by paginating method m.
// If the method is not a paging method, pagingField returns (nil, nil).
// If the method looks like a paging method, but the field cannot be determined, pagingField errors.
func (g *generator) pagingField(m *descriptor.MethodDescriptorProto) (*descriptor.FieldDescriptorProto, error) {
	var (
		hasSize, hasToken, hasNextToken bool
		elemFields                      []*descriptor.FieldDescriptorProto
	)

	inType := g.types[m.GetInputType()]
	if inType == nil {
		return nil, errors.E(nil, "cannot find message type %q, malformed descriptor?", m.GetInputType())
	}
	outType := g.types[m.GetOutputType()]
	if outType == nil {
		return nil, errors.E(nil, "cannot find message type %q, malformed descriptor?", m.GetOutputType())
	}

	for _, f := range inType.Field {
		if f.GetName() == "page_size" && f.GetType() == descriptor.FieldDescriptorProto_TYPE_INT32 {
			hasSize = true
		}
		if f.GetName() == "page_token" && f.GetType() == descriptor.FieldDescriptorProto_TYPE_STRING {
			hasToken = true
		}
	}
	for _, f := range outType.Field {
		if f.GetName() == "next_page_token" && f.GetType() == descriptor.FieldDescriptorProto_TYPE_STRING {
			hasNextToken = true
		}
		if f.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			elemFields = append(elemFields, f)
		}
	}
	if !hasSize || !hasToken || !hasNextToken {
		return nil, nil
	}
	if len(elemFields) == 0 {
		return nil, fmt.Errorf("%s looks like paging method, but can't find repeated field in %s", *m.Name, *outType.Name)
	}
	if len(elemFields) > 1 {
		return nil, fmt.Errorf("%s looks like paging method, but too many repeated fields in %s", *m.Name, *outType.Name)
	}
	return elemFields[0], nil
}

func (g *generator) pagingCall(servName string, m *descriptor.MethodDescriptorProto, elemField *descriptor.FieldDescriptorProto, pt iterType) error {
	inType := g.types[*m.InputType]
	outType := g.types[*m.OutputType]

	inSpec, err := g.importSpec(inType)
	if err != nil {
		return err
	}
	outSpec, err := g.importSpec(outType)
	if err != nil {
		return err
	}

	p := g.printf
	p("func (c *%sClient) %s(ctx context.Context, req *%s.%s, opts ...gax.CallOption) *%s {",
		servName, *m.Name, inSpec.name, *inType.Name, pt.iterTypeName)

	g.insertMetadata()
	g.appendCallOpts(m)

	p("it := &%s{}", pt.iterTypeName)
	p("req = proto.Clone(req).(*%s.%s)", inSpec.name, *inType.Name)
	p("it.InternalFetch = func(pageSize int, pageToken string) ([]%s, string, error) {", pt.elemTypeName)
	p("  var resp *%s.%s", outSpec.name, *outType.Name)
	p("  req.PageToken = pageToken")
	p("  if pageSize > math.MaxInt32 {")
	p("    req.PageSize = math.MaxInt32")
	p("  } else {")
	p("    req.PageSize = int32(pageSize)")
	p("  }")
	p("  err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {")
	p("    var err error")
	p("    resp, err = %s", grpcClientCall(servName, *m.Name))
	p("    return err")
	p("  }, opts...)")
	p("  if err != nil {")
	p("    return nil, \"\", err")
	p("  }")
	p("  return resp.%s, resp.NextPageToken, nil", snakeToCamel(*elemField.Name))
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
	p("it.pageInfo.MaxSize = int(req.PageSize)")
	p("return it")

	p("}")
	p("")

	g.imports[importSpec{path: "math"}] = true
	g.imports[importSpec{path: "github.com/golang/protobuf/proto"}] = true
	g.imports[importSpec{path: "google.golang.org/api/iterator"}] = true
	g.imports[inSpec] = true
	g.imports[outSpec] = true
	for _, spec := range pt.elemImports {
		g.imports[spec] = true
	}
	return nil
}

func (g *generator) pagingIter(pt iterType) {
	p := g.printf

	p("// %s manages a stream of %s.", pt.iterTypeName, pt.elemTypeName)
	p("type %s struct {", pt.iterTypeName)
	p("  items    []%s", pt.elemTypeName)
	p("  pageInfo *iterator.PageInfo")
	p("  nextFunc func() error")
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
