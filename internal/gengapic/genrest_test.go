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
	"testing"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/google/go-cmp/cmp"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
)

// This is a hack to provide an lvalue for a constant.
func typep(typ descriptor.FieldDescriptorProto_Type) *descriptor.FieldDescriptorProto_Type {
	return &typ
}

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
				"kingdom": &descriptor.FieldDescriptorProto{
					Name:   proto.String("kingdom"),
					Number: proto.Int32(int32(2)),
					Type:   typep(descriptor.FieldDescriptorProto_TYPE_INT32),
				},
				"phylum": &descriptor.FieldDescriptorProto{
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
				"kingdom": &descriptor.FieldDescriptorProto{
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
				"kingdom": &descriptor.FieldDescriptorProto{
					Name:   proto.String("kingdom"),
					Number: proto.Int32(int32(0)),
					Type:   typep(descriptor.FieldDescriptorProto_TYPE_INT32),
				},
				"phylum": &descriptor.FieldDescriptorProto{
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
				"mass_kg": &descriptor.FieldDescriptorProto{
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
				"mass_kg": &descriptor.FieldDescriptorProto{
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
			&descriptor.FieldDescriptorProto{
				Name:   proto.String("mass_kg"),
				Number: proto.Int32(int32(0)),
				Type:   typep(descriptor.FieldDescriptorProto_TYPE_INT32),
			},
			&descriptor.FieldDescriptorProto{
				Name:   proto.String("saltwater_p"),
				Number: proto.Int32(int32(1)),
				Type:   typep(descriptor.FieldDescriptorProto_TYPE_BOOL),
			},
		},
	}

	innermostMsg := &descriptor.DescriptorProto{
		Name: proto.String("Chromatophore"),
		Field: []*descriptor.FieldDescriptorProto{
			&descriptor.FieldDescriptorProto{
				Name:   proto.String("color_code"),
				Number: proto.Int32(int32(0)),
				Type:   typep(descriptor.FieldDescriptorProto_TYPE_INT32),
			},
		},
	}
	nestedMsg := &descriptor.DescriptorProto{
		Name: proto.String("Mantle"),
		Field: []*descriptor.FieldDescriptorProto{
			&descriptor.FieldDescriptorProto{
				Name:   proto.String("mass_kg"),
				Number: proto.Int32(int32(0)),
				Type:   typep(descriptor.FieldDescriptorProto_TYPE_INT32),
			},
			&descriptor.FieldDescriptorProto{
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
			&descriptor.FieldDescriptorProto{
				Name:   proto.String("length_m"),
				Number: proto.Int32(int32(0)),
				Type:   typep(descriptor.FieldDescriptorProto_TYPE_INT32),
			},
			&descriptor.FieldDescriptorProto{
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
			&descriptor.FieldDescriptorProto{
				Name:   proto.String("mass_kg"),
				Number: proto.Int32(int32(0)),
				Type:   typep(descriptor.FieldDescriptorProto_TYPE_INT32),
			},
			&descriptor.FieldDescriptorProto{
				Name:     proto.String("whelk"),
				Number:   proto.Int32(int32(1)),
				Type:     typep(descriptor.FieldDescriptorProto_TYPE_MESSAGE),
				TypeName: proto.String(".animalia.mollusca.Whelk"),
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
		},
	}
	req := plugin.CodeGeneratorRequest{
		Parameter: proto.String("go-gapic-package=path;mypackage,transport=rest"),
		ProtoFile: []*descriptor.FileDescriptorProto{file},
	}
	g.init(&req)

	for _, tst := range []struct {
		name     string
		msg      *descriptor.DescriptorProto
		expected map[string]*descriptor.FieldDescriptorProto
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
			name: "recursive_message",
			msg:  recursiveMsg,
			expected: map[string]*descriptor.FieldDescriptorProto{
				"mass_kg":       recursiveMsg.GetField()[0],
				"whelk.mass_kg": recursiveMsg.GetField()[0],
			},
		},
	} {
		actual := g.getLeafs(tst.msg)
		if diff := cmp.Diff(actual, tst.expected, cmp.Comparer(proto.Equal)); diff != "" {
			t.Errorf("test %s, got(-),want(+):\n%s", tst.name, diff)
		}
	}
}
