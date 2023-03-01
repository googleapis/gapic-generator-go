// Copyright 2021 Google LLC
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
	"testing"

	longrunning "cloud.google.com/go/longrunning/autogen/longrunningpb"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/google/go-cmp/cmp"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"github.com/googleapis/gapic-generator-go/internal/txtdiff"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/genproto/googleapis/api/serviceconfig"
	"google.golang.org/genproto/googleapis/cloud/extendedops"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/runtime/protoiface"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Note: the fields parameter contains the names of _all_ the request message's fields,
// not just those that are path or query params.
func setupMethod(g *generator, url, body string, fields []string) (*descriptor.MethodDescriptorProto, error) {
	msg := &descriptor.DescriptorProto{
		Name: proto.String("IdentifyRequest"),
	}
	for i, name := range fields {
		j := int32(i)
		msg.Field = append(msg.GetField(),
			&descriptor.FieldDescriptorProto{
				Name:   proto.String(name),
				Number: &j,
				Type:   typep(descriptor.FieldDescriptorProto_TYPE_INT32),
			},
		)
	}
	mthd := &descriptor.MethodDescriptorProto{
		Name:      proto.String("Identify"),
		InputType: proto.String(".identify.IdentifyRequest"),
		Options:   &descriptor.MethodOptions{},
	}

	// Just use Get for everything and assume parsing general verbs is tested elsewhere.
	proto.SetExtension(mthd.GetOptions(), annotations.E_Http, &annotations.HttpRule{
		Body: body,
		Pattern: &annotations.HttpRule_Get{
			Get: url,
		},
	})

	srv := &descriptor.ServiceDescriptorProto{
		Name:    proto.String("IdentifyMolluscService"),
		Method:  []*descriptor.MethodDescriptorProto{mthd},
		Options: &descriptor.ServiceOptions{},
	}
	proto.SetExtension(srv.GetOptions(), annotations.E_DefaultHost, "linnaean.taxonomy.com")

	fds := []*descriptor.FileDescriptorProto{
		{
			Package:     proto.String("identify"),
			Service:     []*descriptor.ServiceDescriptorProto{srv},
			MessageType: []*descriptor.DescriptorProto{msg},
		},
	}
	req := plugin.CodeGeneratorRequest{
		Parameter: proto.String("go-gapic-package=path;mypackage,transport=rest"),
		ProtoFile: fds,
	}
	g.init(&req)

	return mthd, nil
}

func TestPathParams(t *testing.T) {
	var g generator
	g.apiName = "Awesome Mollusc API"
	g.imports = map[pbinfo.ImportSpec]bool{}
	g.opts = &options{transports: []transport{rest}}

	for _, tst := range []struct {
		name     string
		body     string
		url      string
		fields   []string
		expected map[string]*descriptor.FieldDescriptorProto
	}{
		{
			name:     "no_params",
			url:      "/kingdom",
			fields:   []string{"name", "mass_kg"},
			expected: map[string]*descriptor.FieldDescriptorProto{},
		},
		{
			name:   "params_subset_of_fields",
			url:    "/kingdom/{kingdom}/phylum/{phylum}",
			fields: []string{"name", "mass_kg", "kingdom", "phylum"},
			expected: map[string]*descriptor.FieldDescriptorProto{
				"kingdom": {
					Name:   proto.String("kingdom"),
					Number: proto.Int32(int32(2)),
					Type:   typep(descriptor.FieldDescriptorProto_TYPE_INT32),
				},
				"phylum": {
					Name:   proto.String("phylum"),
					Number: proto.Int32(int32(3)),
					Type:   typep(descriptor.FieldDescriptorProto_TYPE_INT32),
				},
			},
		},
		{
			name:     "disjoint_fields_and_params",
			url:      "/kingdom/{kingdom}",
			fields:   []string{"name", "mass_kg"},
			expected: map[string]*descriptor.FieldDescriptorProto{},
		},
		{
			name:   "params_and_fields_intersect",
			url:    "/kingdom/{kingdom}/phylum/{phylum}",
			fields: []string{"name", "mass_kg", "kingdom"},
			expected: map[string]*descriptor.FieldDescriptorProto{
				"kingdom": {
					Name:   proto.String("kingdom"),
					Number: proto.Int32(int32(2)),
					Type:   typep(descriptor.FieldDescriptorProto_TYPE_INT32),
				},
			},
		},
		{
			name:   "fields_subset_of_params",
			url:    "/kingdom/{kingdom}/phylum/{phylum}/class/{class}",
			fields: []string{"kingdom", "phylum"},
			expected: map[string]*descriptor.FieldDescriptorProto{
				"kingdom": {
					Name:   proto.String("kingdom"),
					Number: proto.Int32(int32(0)),
					Type:   typep(descriptor.FieldDescriptorProto_TYPE_INT32),
				},
				"phylum": {
					Name:   proto.String("phylum"),
					Number: proto.Int32(int32(1)),
					Type:   typep(descriptor.FieldDescriptorProto_TYPE_INT32),
				},
			},
		},
	} {
		mthd, err := setupMethod(&g, tst.url, tst.body, tst.fields)
		if err != nil {
			t.Errorf("test %s setup got error: %s", tst.name, err.Error())
		}

		actual := g.pathParams(mthd)
		if diff := cmp.Diff(actual, tst.expected, cmp.Comparer(proto.Equal)); diff != "" {
			t.Errorf("test %s, got(-),want(+):\n%s", tst.name, diff)
		}
	}
}

func TestQueryParams(t *testing.T) {
	var g generator
	g.apiName = "Awesome Mollusc API"
	g.imports = map[pbinfo.ImportSpec]bool{}
	g.opts = &options{transports: []transport{rest}}
	for _, tst := range []struct {
		name     string
		body     string
		url      string
		fields   []string
		expected map[string]*descriptor.FieldDescriptorProto
	}{
		{
			name:     "all_params_are_path",
			url:      "/kingdom/{kingdom}",
			fields:   []string{"kingdom"},
			expected: map[string]*descriptor.FieldDescriptorProto{},
		},
		{
			name:     "no_fields",
			url:      "/kingdom/{kingdom}",
			fields:   []string{},
			expected: map[string]*descriptor.FieldDescriptorProto{},
		},
		{
			name:   "no_path_params",
			body:   "guess",
			url:    "/kingdom",
			fields: []string{"mass_kg", "guess"},
			expected: map[string]*descriptor.FieldDescriptorProto{
				"mass_kg": {
					Name:   proto.String("mass_kg"),
					Number: proto.Int32(int32(0)),
					Type:   typep(descriptor.FieldDescriptorProto_TYPE_INT32),
				},
			},
		},
		{
			name:   "path_query_param_mix",
			body:   "guess",
			url:    "/kingdom/{kingdom}/phylum/{phylum}",
			fields: []string{"kingdom", "phylum", "mass_kg", "guess"},
			expected: map[string]*descriptor.FieldDescriptorProto{
				"mass_kg": {
					Name:   proto.String("mass_kg"),
					Number: proto.Int32(int32(2)),
					Type:   typep(descriptor.FieldDescriptorProto_TYPE_INT32),
				},
			},
		},
	} {
		mthd, err := setupMethod(&g, tst.url, tst.body, tst.fields)
		if err != nil {
			t.Errorf("test %s setup got error: %s", tst.name, err.Error())
		}

		actual := g.queryParams(mthd)
		if diff := cmp.Diff(actual, tst.expected, cmp.Comparer(proto.Equal)); diff != "" {
			t.Errorf("test %s, got(-),want(+):\n%s", tst.name, diff)
		}
	}
}

func TestLeafFields(t *testing.T) {
	var g generator
	g.apiName = "Awesome Mollusc API"
	g.imports = map[pbinfo.ImportSpec]bool{}
	g.opts = &options{transports: []transport{rest}}

	basicMsg := &descriptor.DescriptorProto{
		Name: proto.String("Clam"),
		Field: []*descriptor.FieldDescriptorProto{
			{
				Name:   proto.String("mass_kg"),
				Number: proto.Int32(int32(0)),
				Type:   typep(descriptor.FieldDescriptorProto_TYPE_INT32),
			},
			{
				Name:   proto.String("saltwater_p"),
				Number: proto.Int32(int32(1)),
				Type:   typep(descriptor.FieldDescriptorProto_TYPE_BOOL),
			},
		},
	}

	innermostMsg := &descriptor.DescriptorProto{
		Name: proto.String("Chromatophore"),
		Field: []*descriptor.FieldDescriptorProto{
			{
				Name:   proto.String("color_code"),
				Number: proto.Int32(int32(0)),
				Type:   typep(descriptor.FieldDescriptorProto_TYPE_INT32),
			},
		},
	}
	nestedMsg := &descriptor.DescriptorProto{
		Name: proto.String("Mantle"),
		Field: []*descriptor.FieldDescriptorProto{
			{
				Name:   proto.String("mass_kg"),
				Number: proto.Int32(int32(0)),
				Type:   typep(descriptor.FieldDescriptorProto_TYPE_INT32),
			},
			{
				Name:     proto.String("chromatophore"),
				Number:   proto.Int32(int32(1)),
				Type:     typep(descriptor.FieldDescriptorProto_TYPE_MESSAGE),
				TypeName: proto.String(".animalia.mollusca.Chromatophore"),
			},
		},
	}
	complexMsg := &descriptor.DescriptorProto{
		Name: proto.String("Squid"),
		Field: []*descriptor.FieldDescriptorProto{
			{
				Name:   proto.String("length_m"),
				Number: proto.Int32(int32(0)),
				Type:   typep(descriptor.FieldDescriptorProto_TYPE_INT32),
			},
			{
				Name:     proto.String("mantle"),
				Number:   proto.Int32(int32(1)),
				Type:     typep(descriptor.FieldDescriptorProto_TYPE_MESSAGE),
				TypeName: proto.String(".animalia.mollusca.Mantle"),
			},
		},
	}

	recursiveMsg := &descriptor.DescriptorProto{
		// Usually it's turtles all the way down, but here it's whelks
		Name: proto.String("Whelk"),
		Field: []*descriptor.FieldDescriptorProto{
			{
				Name:   proto.String("mass_kg"),
				Number: proto.Int32(int32(0)),
				Type:   typep(descriptor.FieldDescriptorProto_TYPE_INT32),
			},
			{
				Name:     proto.String("whelk"),
				Number:   proto.Int32(int32(1)),
				Type:     typep(descriptor.FieldDescriptorProto_TYPE_MESSAGE),
				TypeName: proto.String(".animalia.mollusca.Whelk"),
			},
		},
	}

	overarchingMsg := &descriptor.DescriptorProto{
		Name: proto.String("Trawl"),
		Field: []*descriptor.FieldDescriptorProto{
			{
				Name:     proto.String("clams"),
				Number:   proto.Int32(int32(0)),
				Label:    labelp(descriptor.FieldDescriptorProto_LABEL_REPEATED),
				Type:     typep(descriptor.FieldDescriptorProto_TYPE_MESSAGE),
				TypeName: proto.String(".animalia.mollusca"),
			},
			{
				Name:   proto.String("mass_kg"),
				Number: proto.Int32(int32(1)),
				Type:   typep(descriptor.FieldDescriptorProto_TYPE_INT32),
			},
		},
	}

	wellKnownMsg := &descriptor.DescriptorProto{
		Name: proto.String("Update"),
		Field: []*descriptor.FieldDescriptorProto{
			{
				Name:     proto.String("update_mask"),
				Number:   proto.Int32(int32(0)),
				Type:     typep(descriptor.FieldDescriptorProto_TYPE_MESSAGE),
				TypeName: proto.String(".google.protobuf.FieldMask"),
			},
		},
	}

	file := &descriptor.FileDescriptorProto{
		Package: proto.String("animalia.mollusca"),
		Options: &descriptor.FileOptions{
			GoPackage: proto.String("mypackage"),
		},
		MessageType: []*descriptor.DescriptorProto{
			basicMsg,
			innermostMsg,
			nestedMsg,
			complexMsg,
			recursiveMsg,
			overarchingMsg,
			wellKnownMsg,
		},
	}
	req := plugin.CodeGeneratorRequest{
		Parameter: proto.String("go-gapic-package=path;mypackage,transport=rest"),
		ProtoFile: []*descriptor.FileDescriptorProto{file},
	}
	g.init(&req)

	for _, tst := range []struct {
		name           string
		msg            *descriptor.DescriptorProto
		expected       map[string]*descriptor.FieldDescriptorProto
		excludedFields []*descriptor.FieldDescriptorProto
	}{
		{
			name: "basic_message_test",
			msg:  basicMsg,
			expected: map[string]*descriptor.FieldDescriptorProto{
				"mass_kg":     basicMsg.GetField()[0],
				"saltwater_p": basicMsg.GetField()[1],
			},
		},
		{
			name: "complex_message_test",
			msg:  complexMsg,
			expected: map[string]*descriptor.FieldDescriptorProto{
				"length_m":                        complexMsg.GetField()[0],
				"mantle.mass_kg":                  nestedMsg.GetField()[0],
				"mantle.chromatophore.color_code": innermostMsg.GetField()[0],
			},
		},
		{
			name: "excluded_message_test",
			msg:  complexMsg,
			expected: map[string]*descriptor.FieldDescriptorProto{
				"length_m":       complexMsg.GetField()[0],
				"mantle.mass_kg": nestedMsg.GetField()[0],
			},
			excludedFields: []*descriptor.FieldDescriptorProto{
				nestedMsg.GetField()[1],
			},
		},
		{
			name: "recursive_message_test",
			msg:  recursiveMsg,
			expected: map[string]*descriptor.FieldDescriptorProto{
				"mass_kg":       recursiveMsg.GetField()[0],
				"whelk.mass_kg": recursiveMsg.GetField()[0],
			},
		},
		{
			name: "repeated_message_test",
			msg:  overarchingMsg,
			expected: map[string]*descriptor.FieldDescriptorProto{
				"mass_kg": overarchingMsg.GetField()[1],
			},
		},
		{
			name: "well_known_message_test",
			msg:  wellKnownMsg,
			expected: map[string]*descriptor.FieldDescriptorProto{
				"update_mask": wellKnownMsg.GetField()[0],
			},
		},
	} {
		actual := g.getLeafs(tst.msg, tst.excludedFields...)
		if diff := cmp.Diff(actual, tst.expected, cmp.Comparer(proto.Equal)); diff != "" {
			t.Errorf("test %s, got(-),want(+):\n%s", tst.name, diff)
		}
	}
}

func TestGenRestMethod(t *testing.T) {
	pkg := "google.cloud.foo.v1"

	sizeOpts := &descriptor.FieldOptions{}
	proto.SetExtension(sizeOpts, annotations.E_FieldBehavior, []annotations.FieldBehavior{annotations.FieldBehavior_REQUIRED})
	sizeField := &descriptor.FieldDescriptorProto{
		Name:    proto.String("size"),
		Type:    descriptor.FieldDescriptorProto_TYPE_INT32.Enum(),
		Options: sizeOpts,
	}
	otherField := &descriptor.FieldDescriptorProto{
		Name:           proto.String("other"),
		Type:           descriptor.FieldDescriptorProto_TYPE_STRING.Enum(),
		Proto3Optional: proto.Bool(true),
	}

	foo := &descriptor.DescriptorProto{
		Name:  proto.String("Foo"),
		Field: []*descriptor.FieldDescriptorProto{sizeField, otherField},
	}
	foofqn := fmt.Sprintf(".%s.Foo", pkg)

	pageSizeField := &descriptor.FieldDescriptorProto{
		Name: proto.String("page_size"),
		Type: descriptor.FieldDescriptorProto_TYPE_INT32.Enum(),
	}
	pageTokenField := &descriptor.FieldDescriptorProto{
		Name: proto.String("page_token"),
		Type: descriptor.FieldDescriptorProto_TYPE_STRING.Enum(),
	}
	pagedFooReq := &descriptor.DescriptorProto{
		Name:  proto.String("PagedFooRequest"),
		Field: []*descriptor.FieldDescriptorProto{pageSizeField, pageTokenField},
	}
	pagedFooReqFQN := fmt.Sprintf(".%s.PagedFooRequest", pkg)

	foosField := &descriptor.FieldDescriptorProto{
		Name:     proto.String("foos"),
		Type:     descriptor.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
		TypeName: proto.String(foofqn),
		Label:    descriptor.FieldDescriptorProto_LABEL_REPEATED.Enum(),
	}
	nextPageTokenField := &descriptor.FieldDescriptorProto{
		Name: proto.String("next_page_token"),
		Type: descriptor.FieldDescriptorProto_TYPE_STRING.Enum(),
	}
	pagedFooRes := &descriptor.DescriptorProto{
		Name:  proto.String("PagedFooResponse"),
		Field: []*descriptor.FieldDescriptorProto{foosField, nextPageTokenField},
	}
	pagedFooResFQN := fmt.Sprintf(".%s.PagedFooResponse", pkg)

	fooField := &descriptor.FieldDescriptorProto{
		Name:     proto.String("foo"),
		Type:     descriptor.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
		TypeName: proto.String(foofqn),
	}

	maskField := &descriptor.FieldDescriptorProto{
		Name:     proto.String("update_mask"),
		JsonName: proto.String("updateMask"),
		Type:     descriptor.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
		TypeName: proto.String(".google.protobuf.FieldMask"),
	}

	repeatedPrimField := &descriptor.FieldDescriptorProto{
		Name:     proto.String("primitives"),
		Type:     descriptor.FieldDescriptorProto_TYPE_STRING.Enum().Enum(),
		TypeName: proto.String(foofqn),
		Label:    descriptor.FieldDescriptorProto_LABEL_REPEATED.Enum(),
	}

	updateReq := &descriptor.DescriptorProto{
		Name:  proto.String("UpdateRequest"),
		Field: []*descriptor.FieldDescriptorProto{fooField, maskField, repeatedPrimField},
	}
	updateReqFqn := fmt.Sprintf(".%s.UpdateRequest", pkg)

	nameOpts := &descriptor.FieldOptions{}
	proto.SetExtension(nameOpts, extendedops.E_OperationField, extendedops.OperationResponseMapping_NAME)
	nameField := &descriptor.FieldDescriptorProto{
		Name:    proto.String("name"),
		Type:    descriptor.FieldDescriptorProto_TYPE_STRING.Enum(),
		Options: nameOpts,
	}
	op := &descriptor.DescriptorProto{
		Name:  proto.String("Operation"),
		Field: []*descriptor.FieldDescriptorProto{nameField},
	}
	opfqn := fmt.Sprintf(".%s.Operation", pkg)

	opRPCOpt := &descriptor.MethodOptions{}
	proto.SetExtension(opRPCOpt, annotations.E_Http, &annotations.HttpRule{
		Pattern: &annotations.HttpRule_Post{
			Post: "/v1/foo",
		},
	})
	proto.SetExtension(opRPCOpt, extendedops.E_OperationService, "FooOperationService")

	opRPC := &descriptor.MethodDescriptorProto{
		Name:       proto.String("CustomOp"),
		InputType:  proto.String(foofqn),
		OutputType: proto.String(opfqn),
		Options:    opRPCOpt,
	}

	emptyRPCOpt := &descriptor.MethodOptions{}
	proto.SetExtension(emptyRPCOpt, annotations.E_Http, &annotations.HttpRule{
		Pattern: &annotations.HttpRule_Delete{
			Delete: "/v1/foo/{other=*}",
		},
	})

	emptyRPC := &descriptor.MethodDescriptorProto{
		Name:       proto.String("EmptyRPC"),
		InputType:  proto.String(foofqn),
		OutputType: proto.String(emptyType),
		Options:    emptyRPCOpt,
	}

	unaryRPCOpt := &descriptor.MethodOptions{}
	proto.SetExtension(unaryRPCOpt, annotations.E_Http, &annotations.HttpRule{
		Pattern: &annotations.HttpRule_Post{
			Post: "/v1/foo",
		},
		Body: "*",
	})
	proto.SetExtension(unaryRPCOpt, annotations.E_Routing, &annotations.RoutingRule{
		RoutingParameters: []*annotations.RoutingParameter{
			{Field: "other"},
		},
	})

	unaryRPC := &descriptor.MethodDescriptorProto{
		Name:       proto.String("UnaryRPC"),
		InputType:  proto.String(foofqn),
		OutputType: proto.String(foofqn),
		Options:    unaryRPCOpt,
	}

	pagingRPCOpt := &descriptor.MethodOptions{}
	proto.SetExtension(pagingRPCOpt, annotations.E_Http, &annotations.HttpRule{
		Pattern: &annotations.HttpRule_Get{
			Get: "/v1/foo",
		},
	})

	pagingRPC := &descriptor.MethodDescriptorProto{
		Name:       proto.String("PagingRPC"),
		InputType:  proto.String(pagedFooReqFQN),
		OutputType: proto.String(pagedFooResFQN),
		Options:    pagingRPCOpt,
	}

	serverStreamRPC := &descriptor.MethodDescriptorProto{
		Name:            proto.String("ServerStreamRPC"),
		InputType:       proto.String(foofqn),
		OutputType:      proto.String(foofqn),
		ServerStreaming: proto.Bool(true),
		// Reuse the unary RPC options because it's basically the same.
		Options: unaryRPCOpt,
	}

	clientStreamRPC := &descriptor.MethodDescriptorProto{
		Name:            proto.String("ClientStreamRPC"),
		InputType:       proto.String(foofqn),
		OutputType:      proto.String(foofqn),
		ClientStreaming: proto.Bool(true),
		// Reuse the unary RPC options because it's basically the same.
		Options: unaryRPCOpt,
	}

	lroRPCOpt := &descriptor.MethodOptions{}
	proto.SetExtension(lroRPCOpt, annotations.E_Http, &annotations.HttpRule{
		Pattern: &annotations.HttpRule_Post{
			Post: "/v1/foo",
		},
		Body: "*",
	})
	proto.SetExtension(lroRPCOpt, longrunning.E_OperationInfo, &longrunning.OperationInfo{
		// Need to trim the leading "." as the annotation value shouldn't have it.
		ResponseType: foofqn[1:],
		MetadataType: foofqn[1:],
	})
	lroDesc := protodesc.ToDescriptorProto((&longrunning.Operation{}).ProtoReflect().Descriptor())
	lroRPC := &descriptor.MethodDescriptorProto{
		Name:       proto.String("LongrunningRPC"),
		InputType:  proto.String(foofqn),
		OutputType: proto.String(lroType),
		Options:    lroRPCOpt,
	}

	httpBodyDesc := protodesc.ToDescriptorProto((&httpbody.HttpBody{}).ProtoReflect().Descriptor())
	httpBodyRPCOpt := &descriptor.MethodOptions{}
	proto.SetExtension(httpBodyRPCOpt, annotations.E_Http, &annotations.HttpRule{
		Pattern: &annotations.HttpRule_Post{
			Post: "/v1/foo",
		},
		Body: "*",
	})
	proto.SetExtension(httpBodyRPCOpt, annotations.E_Routing, &annotations.RoutingRule{
		RoutingParameters: []*annotations.RoutingParameter{
			{Field: "other"},
		},
	})
	httpBodyRPC := &descriptor.MethodDescriptorProto{
		Name:       proto.String("HttpBodyRPC"),
		InputType:  proto.String(foofqn),
		OutputType: proto.String(httpBodyType),
		Options:    httpBodyRPCOpt,
	}

	updateRPCOpt := &descriptor.MethodOptions{}
	proto.SetExtension(updateRPCOpt, annotations.E_Http, &annotations.HttpRule{
		Pattern: &annotations.HttpRule_Post{
			Post: "/v1/foo",
		},
		Body: "foo",
	})

	updateRPC := &descriptor.MethodDescriptorProto{
		Name:       proto.String("UpdateRPC"),
		InputType:  proto.String(updateReqFqn),
		OutputType: proto.String(foofqn),
		Options:    updateRPCOpt,
	}

	s := &descriptor.ServiceDescriptorProto{
		Name: proto.String("FooService"),
	}
	opS := &descriptor.ServiceDescriptorProto{
		Name: proto.String("FooOperationService"),
	}

	f := &descriptor.FileDescriptorProto{
		Package: proto.String(pkg),
		Options: &descriptor.FileOptions{
			GoPackage: proto.String("google.golang.org/genproto/cloud/foo/v1;foo"),
		},
		Service: []*descriptor.ServiceDescriptorProto{s, opS},
	}

	g := &generator{
		aux: &auxTypes{
			customOp: &customOp{
				message: op,
			},
			iters: map[string]*iterType{},
			lros:  map[*descriptor.MethodDescriptorProto]bool{},
		},
		opts: &options{},
		customOpServices: map[*descriptor.ServiceDescriptorProto]*descriptor.ServiceDescriptorProto{
			s: opS,
		},
		descInfo: pbinfo.Info{
			ParentFile: map[protoiface.MessageV1]*descriptor.FileDescriptorProto{
				op:           f,
				opS:          f,
				opRPC:        f,
				lroRPC:       f,
				updateRPC:    f,
				foo:          f,
				s:            f,
				pagedFooReq:  f,
				pagedFooRes:  f,
				lroDesc:      protodesc.ToFileDescriptorProto(longrunning.File_google_longrunning_operations_proto),
				httpBodyDesc: protodesc.ToFileDescriptorProto(httpbody.File_google_api_httpbody_proto),
				updateReq:    f,
			},
			ParentElement: map[pbinfo.ProtoType]pbinfo.ProtoType{
				opRPC:           s,
				emptyRPC:        s,
				unaryRPC:        s,
				pagingRPC:       s,
				serverStreamRPC: s,
				clientStreamRPC: s,
				lroRPC:          s,
				httpBodyRPC:     s,
				updateRPC:       s,
				nameField:       op,
				sizeField:       foo,
				otherField:      foo,
				maskField:       updateReq,
			},
			Type: map[string]pbinfo.ProtoType{
				opfqn:          op,
				foofqn:         foo,
				emptyType:      protodesc.ToDescriptorProto((&emptypb.Empty{}).ProtoReflect().Descriptor()),
				pagedFooReqFQN: pagedFooReq,
				pagedFooResFQN: pagedFooRes,
				lroType:        lroDesc,
				httpBodyType:   httpBodyDesc,
				updateReqFqn:   updateReq,
			},
		},
	}

	for _, tst := range []struct {
		name    string
		method  *descriptor.MethodDescriptorProto
		options *options
		imports map[pbinfo.ImportSpec]bool
	}{
		{
			name:    "custom_op",
			method:  opRPC,
			options: &options{diregapic: true},
			imports: map[pbinfo.ImportSpec]bool{
				{Path: "google.golang.org/protobuf/encoding/protojson"}: true,
				{Path: "io/ioutil"}: true,
				{Path: "net/url"}:   true,
				{Path: "fmt"}:       true,
				{Path: "google.golang.org/api/googleapi"}:                        true,
				{Name: "foopb", Path: "google.golang.org/genproto/cloud/foo/v1"}: true,
			},
		},
		{
			name:    "empty_rpc",
			method:  emptyRPC,
			options: &options{},
			imports: map[pbinfo.ImportSpec]bool{
				{Path: "fmt"}: true,
				{Path: "google.golang.org/api/googleapi"}:                        true,
				{Path: "net/url"}:                                                true,
				{Name: "foopb", Path: "google.golang.org/genproto/cloud/foo/v1"}: true,
			},
		},
		{
			name:    "unary_rpc",
			method:  unaryRPC,
			options: &options{restNumericEnum: true},
			imports: map[pbinfo.ImportSpec]bool{
				{Path: "bytes"}: true,
				{Path: "fmt"}:   true,
				{Path: "google.golang.org/protobuf/encoding/protojson"}: true,
				{Path: "io/ioutil"}:                       true,
				{Path: "google.golang.org/api/googleapi"}: true,
				{Path: "net/url"}:                         true,
				{Path: "regexp"}:                          true,
				{Path: "strings"}:                         true,
				{Name: "foopb", Path: "google.golang.org/genproto/cloud/foo/v1"}: true,
			},
		},
		{
			name:    "paging_rpc",
			method:  pagingRPC,
			options: &options{},
			imports: map[pbinfo.ImportSpec]bool{
				{Path: "math"}:    true,
				{Path: "net/url"}: true,
				{Path: "google.golang.org/protobuf/encoding/protojson"}: true,
				{Path: "fmt"}: true,
				{Path: "google.golang.org/api/googleapi"}:                        true,
				{Path: "google.golang.org/api/iterator"}:                         true,
				{Path: "google.golang.org/protobuf/proto"}:                       true,
				{Path: "io/ioutil"}:                                              true,
				{Name: "foopb", Path: "google.golang.org/genproto/cloud/foo/v1"}: true,
			},
		},
		{
			name:    "server_stream_rpc",
			method:  serverStreamRPC,
			options: &options{},
			imports: map[pbinfo.ImportSpec]bool{
				{Path: "bytes"}:   true,
				{Path: "context"}: true,
				{Path: "fmt"}:     true,
				{Path: "google.golang.org/api/googleapi"}:                        true,
				{Path: "net/url"}:                                                true,
				{Path: "regexp"}:                                                 true,
				{Path: "strings"}:                                                true,
				{Path: "google.golang.org/grpc/metadata"}:                        true,
				{Path: "google.golang.org/protobuf/encoding/protojson"}:          true,
				{Name: "foopb", Path: "google.golang.org/genproto/cloud/foo/v1"}: true,
			},
		},
		{
			name:    "no_request_stream_rpc",
			method:  clientStreamRPC,
			options: &options{},
			imports: map[pbinfo.ImportSpec]bool{
				{Path: "context"}: true,
				{Path: "fmt"}:     true,
				{Name: "foopb", Path: "google.golang.org/genproto/cloud/foo/v1"}: true,
			},
		},
		{
			name:    "lro_rpc",
			method:  lroRPC,
			options: &options{transports: []transport{rest}},
			imports: map[pbinfo.ImportSpec]bool{
				{Path: "bytes"}: true,
				{Path: "cloud.google.com/go/longrunning"}: true,
				{Path: "fmt"}: true,
				{Path: "google.golang.org/api/googleapi"}:                        true,
				{Path: "google.golang.org/grpc/metadata"}:                        true,
				{Path: "google.golang.org/protobuf/encoding/protojson"}:          true,
				{Path: "io/ioutil"}:                                              true,
				{Path: "net/url"}:                                                true,
				{Path: "time"}:                                                   true,
				{Name: "foopb", Path: "google.golang.org/genproto/cloud/foo/v1"}: true,
				{Name: "longrunningpb", Path: "google.golang.org/genproto/googleapis/longrunning"}: true,
			},
		},
		{
			name:    "http_body_rpc",
			method:  httpBodyRPC,
			options: &options{},
			imports: map[pbinfo.ImportSpec]bool{
				{Path: "bytes"}: true,
				{Path: "fmt"}:   true,
				{Path: "google.golang.org/protobuf/encoding/protojson"}: true,
				{Path: "io/ioutil"}:                       true,
				{Path: "google.golang.org/api/googleapi"}: true,
				{Path: "net/url"}:                         true,
				{Path: "regexp"}:                          true,
				{Path: "strings"}:                         true,
				{Name: "foopb", Path: "google.golang.org/genproto/cloud/foo/v1"}:                 true,
				{Name: "httpbodypb", Path: "google.golang.org/genproto/googleapis/api/httpbody"}: true,
			},
		},
		{
			name:    "update_rpc",
			method:  updateRPC,
			options: &options{restNumericEnum: true},
			imports: map[pbinfo.ImportSpec]bool{
				{Path: "bytes"}: true,
				{Path: "fmt"}:   true,
				{Path: "google.golang.org/protobuf/encoding/protojson"}: true,
				{Path: "io/ioutil"}:                       true,
				{Path: "google.golang.org/api/googleapi"}: true,
				{Path: "net/url"}:                         true,
				{Name: "foopb", Path: "google.golang.org/genproto/cloud/foo/v1"}: true,
			},
		},
	} {
		s.Method = []*descriptor.MethodDescriptorProto{tst.method}
		g.opts = tst.options
		g.imports = make(map[pbinfo.ImportSpec]bool)
		g.serviceConfig = &serviceconfig.Service{
			Http: &annotations.Http{
				Rules: []*annotations.HttpRule{
					{
						Selector: "google.longrunning.Operations.GetOperation",
						Pattern: &annotations.HttpRule_Get{
							Get: "/v1beta1/{name=projects/*/locations/*/operations/*}",
						},
					},
				},
			},
		}

		if err := g.genRESTMethod("Foo", s, tst.method); err != nil {
			t.Fatal(err)
		}

		var genLros []*descriptor.MethodDescriptorProto
		for m := range g.aux.lros {
			genLros = append(genLros, m)
		}
		sort.Slice(genLros, func(i, j int) bool {
			return genLros[i].GetName() < genLros[j].GetName()
		})
		for _, m := range genLros {
			if err := g.lroType("Foo", s, m); err != nil {
				t.Error(err)
			}
		}

		if diff := cmp.Diff(g.imports, tst.imports); diff != "" {
			t.Errorf("TestGenRESTMethod(%s): imports got(-),want(+):\n%s", tst.name, diff)
		}

		txtdiff.Diff(t, fmt.Sprintf("%s_%s", t.Name(), tst.name), g.pt.String(), filepath.Join("testdata", fmt.Sprintf("rest_%s.want", tst.method.GetName())))
		g.reset()
		g.aux.lros = make(map[*descriptor.MethodDescriptorProto]bool)
	}
}
