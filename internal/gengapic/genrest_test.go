// Copyright (C) 2021  Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gengapic

import (
	"reflect"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"google.golang.org/genproto/googleapis/api/annotations"
)

// This is a hack to provide an lvalue for a constant.
func typep(typ descriptor.FieldDescriptorProto_Type) *descriptor.FieldDescriptorProto_Type {
	return &typ
}

// Another hack, this time for field numbers.
func idxAddr(i int) *int32 {
	j := int32(i)
	return &j
}

func setupMethod(t *testing.T, g *generator, url string, allFieldNames []string) *descriptor.MethodDescriptorProto {
	msg := &descriptor.DescriptorProto{
		Name: proto.String("IdentifyRequest"),
	}
	for i, name := range allFieldNames {
		j := int32(i)
		msg.Field = append(msg.Field,
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
	if err := proto.SetExtension(mthd.Options, annotations.E_Http, &annotations.HttpRule{
		Pattern: &annotations.HttpRule_Get{
			Get: url,
		},
	}); err != nil {
		t.Fatal(err)
	}

	srv := &descriptor.ServiceDescriptorProto{
		Name:    proto.String("IdentifyMolluscService"),
		Method:  []*descriptor.MethodDescriptorProto{mthd},
		Options: &descriptor.ServiceOptions{},
	}
	if err := proto.SetExtension(srv.Options, annotations.E_DefaultHost, proto.String("linnaean.taxonomy.com")); err != nil {
		t.Fatal(err)
	}

	fds := []*descriptor.FileDescriptorProto{
		&descriptor.FileDescriptorProto{
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

	return mthd
}

func TestPathParams(t *testing.T) {
	var g generator
	g.apiName = "Awesome Mollusc API"
	g.imports = map[pbinfo.ImportSpec]bool{}
	g.opts = &options{transports: []transport{rest}}

	for tstNum, tst := range []struct {
		url           string
		allFieldNames []string
		expected      map[string]*descriptor.FieldDescriptorProto
	}{
		// Null case: no params
		{
			url:           "/kingdom",
			allFieldNames: []string{"name", "mass_kg"},
			expected:      map[string]*descriptor.FieldDescriptorProto{},
		},
		// Reasonable case: params are a subset of fields
		{
			url:           "/kingdom/{kingdom}/phylum/{phylum}",
			allFieldNames: []string{"name", "mass_kg", "kingdom", "phylum"},
			expected: map[string]*descriptor.FieldDescriptorProto{
				"kingdom": &descriptor.FieldDescriptorProto{
					Name:   proto.String("kingdom"),
					Number: idxAddr(2),
					Type:   typep(descriptor.FieldDescriptorProto_TYPE_INT32),
				},
				"phylum": &descriptor.FieldDescriptorProto{
					Name:   proto.String("phylum"),
					Number: idxAddr(3),
					Type:   typep(descriptor.FieldDescriptorProto_TYPE_INT32),
				},
			},
		},
		// Degenerate case 1: params and fields are disjoint
		{
			url:           "/kingdom/{kingdom}",
			allFieldNames: []string{"name", "mass_kg"},
			expected:      map[string]*descriptor.FieldDescriptorProto{},
		},
		// Degenerate case 2: params and fields intersect but are not subsets
		{
			url:           "/kingdom/{kingdom}/phylum/{phylum}",
			allFieldNames: []string{"name", "mass_kg", "kingdom"},
			expected: map[string]*descriptor.FieldDescriptorProto{
				"kingdom": &descriptor.FieldDescriptorProto{
					Name:   proto.String("kingdom"),
					Number: idxAddr(2),
					Type:   typep(descriptor.FieldDescriptorProto_TYPE_INT32),
				},
			},
		},
		// Degenerate case 3: fields are a subset of params
		{
			url:           "/kingdom/{kingdom}/phylum/{phylum}/class/{class}",
			allFieldNames: []string{"kingdom", "phylum"},
			expected: map[string]*descriptor.FieldDescriptorProto{
				"kingdom": &descriptor.FieldDescriptorProto{
					Name:   proto.String("kingdom"),
					Number: idxAddr(0),
					Type:   typep(descriptor.FieldDescriptorProto_TYPE_INT32),
				},
				"phylum": &descriptor.FieldDescriptorProto{
					Name:   proto.String("phylum"),
					Number: idxAddr(1),
					Type:   typep(descriptor.FieldDescriptorProto_TYPE_INT32),
				},
			},
		},
	} {
		mthd := setupMethod(t, &g, tst.url, tst.allFieldNames)

		actual := g.pathParams(mthd)
		if !reflect.DeepEqual(actual, tst.expected) {
			t.Errorf("test %d, pathParams(%q) = %q, want %q", tstNum, tst.url, actual, tst.expected)
		}
	}
}

func TestQueryParams(t *testing.T) {
	var g generator
	g.apiName = "Awesome Mollusc API"
	g.imports = map[pbinfo.ImportSpec]bool{}
	g.opts = &options{transports: []transport{rest}}
	for _, tst := range []struct {
		url           string
		allFieldNames []string
		expected      map[string]*descriptor.FieldDescriptorProto
	}{
		// All params are path params
		{
			url:           "/kingdom/{kingdom}",
			allFieldNames: []string{"kingdom"},
			expected:      map[string]*descriptor.FieldDescriptorProto{},
		},
		// Degenerate case: no fields
		{
			url:           "/kingdom/{kingdom}",
			allFieldNames: []string{},
			expected:      map[string]*descriptor.FieldDescriptorProto{},
		},
		// No path params
		{
			url:           "/kingdom",
			allFieldNames: []string{"mass_kg"},
			expected: map[string]*descriptor.FieldDescriptorProto{
				"mass_kg": &descriptor.FieldDescriptorProto{
					Name:   proto.String("mass_kg"),
					Number: idxAddr(0),
					Type:   typep(descriptor.FieldDescriptorProto_TYPE_INT32),
				},
			},
		},
		// Interesting case: some params are path and some are query
		{
			url:           "/kingdom/{kingdom}/phylum/{phylum}",
			allFieldNames: []string{"kingdom", "phylum", "mass_kg", "guess"},
			expected: map[string]*descriptor.FieldDescriptorProto{
				"mass_kg": &descriptor.FieldDescriptorProto{
					Name:   proto.String("mass_kg"),
					Number: idxAddr(2),
					Type:   typep(descriptor.FieldDescriptorProto_TYPE_INT32),
				},
				"guess": &descriptor.FieldDescriptorProto{
					Name:   proto.String("guess"),
					Number: idxAddr(3),
					Type:   typep(descriptor.FieldDescriptorProto_TYPE_INT32),
				},
			},
		},
	} {
		mthd := setupMethod(t, &g, tst.url, tst.allFieldNames)

		actual := g.queryParams(mthd)
		if !reflect.DeepEqual(actual, tst.expected) {
			t.Errorf("queryParams(%q) = %q, want %q", tst.url, actual, tst.expected)
		}
	}
}
