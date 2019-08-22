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

package main

import (
	"path/filepath"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/gapic-generator-go/cmd/gen-go-sample/schema_v1p2"
	"github.com/googleapis/gapic-generator-go/internal/errors"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"github.com/googleapis/gapic-generator-go/internal/txtdiff"
)

func TestUnary(t *testing.T) {
	t.Parallel()

	g := initTestGenerator()

	sp := schema_v1p2.Sample{
		ID:        "my_sample_config",
		Rpc:       "UnaryMethod",
		Service:   "foo.FooService",
		RegionTag: "awesome_region",
		Request: []schema_v1p2.RequestConfig{
			{Field: "a.x", Value: "42", InputParameter: "the_x"},
			{Field: "a.y", Value: "3.14159"},
			{Field: "b", Value: "foobar", InputParameter: "the_b"},
			{Field: "e", Value: "BANANA"},
			{Field: "f", Value: "in a oneof"},
			{Field: "bytes", Value: "mybytes"},
			{Field: "data_alice", Value: "path/to/local/file/alice.txt", ValueIsFile: true},
			{Field: "a_array[0].x", Value: "0"},
			{Field: "a_array[0].y", Value: "1"},
			{Field: "a_array[1].x", Value: "2"},
			{Field: "a_array[1].y", Value: "3"},
			{Field: "resource_field%foo", Value: "myfoo", InputParameter: "the_foo"},
			{Field: "data_bob", Value: "path/to/local/file/bob.txt", ValueIsFile: true, InputParameter: "bob_file"},
			{Field: "resource_field%bar", Value: "mybar"},
		},
		Response: []schema_v1p2.ResponseConfig{
			{Define: "out_a = $resp.a"},
			{Print: []string{"x = %s", "$resp.a.x"}},
			{Comment: []string{
				"%s contains the field %s,\n%%%% is a single percent,\nand '\\n' specifies a newline in comment too",
				"out_a",
				"y"}},
			{Print: []string{"y = %s", "out_a.y"}},
			{
				Loop: &schema_v1p2.LoopSpec{
					Variable:   "r",
					Collection: "$resp.r",
					Body: []schema_v1p2.ResponseConfig{
						{Print: []string{"resp.r contains %s", "r"}},
					},
				},
			},
		},
	}

	methConf := GAPICMethod{
		Name:              "UnaryMethod",
		FieldNamePatterns: map[string]string{"resource_field": "foobar_thing"},
	}
	if err := g.genSample(sp, methConf); err != nil {
		t.Fatal(err)
	}

	compare(t, g, filepath.Join("testdata", "sample_unary.want"))
}

func TestSample_InitError(t *testing.T) {
	g := initTestGenerator()
	for _, fieldVal := range [][]string{
		// bad type
		{"a", "string"},
		{"a.x", "string"},
		{"a_array[0].x", "string"},

		// try to access array
		{"a_array.x", "0"},
		// try to index singular
		{"a[0].x", "0"},
		// missing ']'
		{"a_array[0.x", "0"},
		// index can only be ints
		{"a_array[0.0].x", "0"},
		{`a_array["0"].x`, "0"},

		// bad enum variant
		{"e", "DERP"},
	} {
		g.reset()

		sp := schema_v1p2.Sample{
			ID:        "my_sample_config",
			Service:   "foo.FooService",
			Rpc:       "LroMethod",
			RegionTag: "awesome_region",
			Request: []schema_v1p2.RequestConfig{
				{Field: fieldVal[0], Value: fieldVal[1]},
			},
		}

		if err := g.genSample(sp, GAPICMethod{Name: "UnaryMethod"}); err == nil {
			t.Errorf("expected error from init config: %s", fieldVal)
		}
	}

	// missing LRO config
	g.reset()

	sp := schema_v1p2.Sample{
		ID:        "my_sample_config",
		Service:   "foo.FooService",
		Rpc:       "LroMethod",
		RegionTag: "awesome_region",
	}

	if err := g.genSample(sp, GAPICMethod{Name: "LroMethod"}); err == nil {
		t.Errorf("expected error from missing config")
	}

}

func TestPaging(t *testing.T) {
	t.Parallel()

	g := initTestGenerator()

	sp := schema_v1p2.Sample{
		ID:        "my_sample_config",
		Service:   "foo.FooService",
		Rpc:       "PagingMethod",
		RegionTag: "awesome_region",
	}
	if err := g.genSample(sp, GAPICMethod{Name: "PagingMethod"}); err != nil {
		t.Fatal(err)
	}
	compare(t, g, filepath.Join("testdata", "sample_paging.want"))
}

func TestLro(t *testing.T) {
	t.Parallel()

	g := initTestGenerator()

	sp := schema_v1p2.Sample{
		ID:        "my_sample_config",
		Service:   "foo.FooService",
		Rpc:       "LroMethod",
		RegionTag: "awesome_region",
		Response: []schema_v1p2.ResponseConfig{
			{Print: []string{"x = %s", "$resp.a.x"}},
		},
	}

	methConf := GAPICMethod{
		Name: "LroMethod",
		LongRunning: LongRunningConfig{
			ReturnType:   "foo.lroReturnType",
			MetadataType: "foo.lroMetadataType",
		},
	}

	if err := g.genSample(sp, methConf); err != nil {
		t.Fatal(err)
	}
	compare(t, g, filepath.Join("testdata", "sample_lro.want"))
}

func TestEmptyObject(t *testing.T) {
	t.Parallel()

	g := initTestGenerator()
	sp := schema_v1p2.Sample{
		ID:        "my_value_set",
		Service:   "foo.FooService",
		Rpc:       "UnaryMethod",
		RegionTag: "awesome_region",
		Request: []schema_v1p2.RequestConfig{
			{Field: "f2", Value: "{}"},
			{Field: "resource_field%foo", Value: "myfoo"},
			{Field: "resource_field%bar", Value: "mybar"},
		},
	}

	methConf := GAPICMethod{
		Name:              "UnaryMethod",
		FieldNamePatterns: map[string]string{"resource_field": "foobar_thing"},
	}

	if err := g.genSample(sp, methConf); err != nil {
		t.Fatal(err)
	}

	compare(t, g, filepath.Join("testdata", "sample_empty_object.want"))
}

func TestEmpty(t *testing.T) {
	t.Parallel()

	g := initTestGenerator()
	sp := schema_v1p2.Sample{
		ID:        "my_value_set",
		Service:   "foo.FooService",
		Rpc:       "EmptyMethod",
		RegionTag: "awesome_region",
	}

	if err := g.genSample(sp, GAPICMethod{Name: "EmptyMethod"}); err != nil {
		t.Fatal(err)
	}
	compare(t, g, filepath.Join("testdata", "sample_empty.want"))
}

func TestMapOut(t *testing.T) {
	t.Parallel()

	g := initTestGenerator()
	sp := schema_v1p2.Sample{
		ID:        "my_sample_config",
		Service:   "foo.FooService",
		Rpc:       "UnaryMethod",
		RegionTag: "awesome_region",
		Response: []schema_v1p2.ResponseConfig{
			{
				Loop: &schema_v1p2.LoopSpec{
					Key: "just_key",
					Map: "$resp.mappy_map",
					Body: []schema_v1p2.ResponseConfig{
						{Print: []string{"key: %s", "just_key"}},
					},
				},
			},
			{
				Loop: &schema_v1p2.LoopSpec{
					Key:   "k",
					Value: "v",
					Map:   "$resp.mappy_map",
					Body: []schema_v1p2.ResponseConfig{
						{Print: []string{"key: %s, value: %s", "k", "v"}},
					},
				},
			},
			{
				Loop: &schema_v1p2.LoopSpec{
					Value: "only_value",
					Map:   "$resp.mappy_map",
					Body: []schema_v1p2.ResponseConfig{
						{Print: []string{"value: %s", "only_value"}},
					},
				},
			},
		},
	}
	if err := g.genSample(sp, GAPICMethod{Name: "UnaryMethod"}); err != nil {
		t.Fatal(err)
	}
	compare(t, g, filepath.Join("testdata", "sample_map_out.want"))
}

func TestAccessMapKeyValueInResponse(t *testing.T) {
	t.Parallel()

	g := initTestGenerator()
	sp := schema_v1p2.Sample{
		ID:        "my_sample_config",
		Service:   "foo.FooService",
		Rpc:       "UnaryMethod",
		RegionTag: "awesome_region",
		Response: []schema_v1p2.ResponseConfig{
			{Print: []string{"The value associated with some_key is: ", `$resp.mappy_map{"some_key"}`}},
		},
	}

	if err := g.genSample(sp, GAPICMethod{Name: "UnaryMethod"}); err != nil {
		t.Fatal(err)
	}
	compare(t, g, filepath.Join("testdata", "sample_response_map_field.want"))
}

func TestAccessRepeatedFieldInResponse(t *testing.T) {
	t.Parallel()

	g := initTestGenerator()
	sp := schema_v1p2.Sample{
		ID:        "my_sample_config",
		Service:   "foo.FooService",
		Rpc:       "UnaryMethod",
		RegionTag: "awesome_region",
		Response: []schema_v1p2.ResponseConfig{
			{Print: []string{"The first element of the array is: %s", "$resp.a_array[0].x"}},
		},
	}

	if err := g.genSample(sp, GAPICMethod{Name: "UnaryMethod"}); err != nil {
		t.Fatal(err)
	}
	compare(t, g, filepath.Join("testdata", "sample_response_repeated_field.want"))
}

func TestWriteFile(t *testing.T) {
	t.Parallel()

	g := initTestGenerator()
	sp := schema_v1p2.Sample{
		ID:        "my_sample_config",
		Service:   "foo.FooService",
		Rpc:       "UnaryMethod",
		RegionTag: "awesome_region",
		Response: []schema_v1p2.ResponseConfig{
			{
				WriteFile: &schema_v1p2.WriteFileSpec{
					FileName: []string{"my_bob.mp3"},
					Contents: "$resp.data_bob",
				},
			},
			{
				WriteFile: &schema_v1p2.WriteFileSpec{
					FileName: []string{"my_alice_%s.mp3", "$resp.a"},
					Contents: "$resp.b",
				},
			},
		},
	}

	if err := g.genSample(sp, GAPICMethod{Name: "UnaryMethod"}); err != nil {
		t.Fatal(err)
	}
	compare(t, g, filepath.Join("testdata", "sample_write_file.want"))

}

func TestEnum(t *testing.T) {
	t.Parallel()

	g := initTestGenerator()
	sp := schema_v1p2.Sample{
		ID:        "my_sample_config",
		Service:   "foo.FooService",
		Rpc:       "UnaryMethod",
		RegionTag: "awesome_region",
		Request: []schema_v1p2.RequestConfig{
			{Field: "e", Value: "APPLE", InputParameter: "my_fruit"},
		},
		Response: []schema_v1p2.ResponseConfig{
			{Print: []string{"call finished."}},
		},
	}

	if err := g.genSample(sp, GAPICMethod{Name: "UnaryMethod"}); err != nil {
		t.Fatal(err)
	}

	compare(t, g, filepath.Join("testdata", "sample_enum.want"))
}

func initTestGenerator() *generator {
	eType := &descriptor.EnumDescriptorProto{
		Name: proto.String("FruitEnum"),
		Value: []*descriptor.EnumValueDescriptorProto{
			{Name: proto.String("APPLE")},
			{Name: proto.String("BANANA")},
			{Name: proto.String("CHERRY")},
		},
	}

	aType := &descriptor.DescriptorProto{
		Name: proto.String("AType"),
		Field: []*descriptor.FieldDescriptorProto{
			{Name: proto.String("x"), Type: typep(descriptor.FieldDescriptorProto_TYPE_INT64)},
			{Name: proto.String("y"), Type: typep(descriptor.FieldDescriptorProto_TYPE_FLOAT)},
		},
		EnumType: []*descriptor.EnumDescriptorProto{eType},
	}

	mapType, mapField :=
		createMapTypeAndField(
			"mappy_map",
			".foo.InputType",
			typep(descriptor.FieldDescriptorProto_TYPE_STRING),
			typep(descriptor.FieldDescriptorProto_TYPE_MESSAGE),
			".foo.AType")

	inType := &descriptor.DescriptorProto{
		Name: proto.String("InputType"),
		OneofDecl: []*descriptor.OneofDescriptorProto{
			{Name: proto.String("Group")},
		},
		NestedType: []*descriptor.DescriptorProto{
			mapType,
		},
		Field: []*descriptor.FieldDescriptorProto{
			{Name: proto.String("a"), TypeName: proto.String(".foo.AType")},
			{Name: proto.String("a_array"), TypeName: proto.String(".foo.AType"), Label: labelp(descriptor.FieldDescriptorProto_LABEL_REPEATED)},
			{Name: proto.String("b"), Type: typep(descriptor.FieldDescriptorProto_TYPE_STRING)},
			{Name: proto.String("e"), TypeName: proto.String(".foo.AType.FruitEnum")},
			{Name: proto.String("f"), Type: typep(descriptor.FieldDescriptorProto_TYPE_STRING), OneofIndex: proto.Int32(0)},
			{Name: proto.String("f2"), TypeName: proto.String(".foo.AType"), OneofIndex: proto.Int32(0)},
			{Name: proto.String("data_alice"), Type: typep(descriptor.FieldDescriptorProto_TYPE_BYTES)},
			{Name: proto.String("data_bob"), Type: typep(descriptor.FieldDescriptorProto_TYPE_BYTES)},
			{Name: proto.String("r"), Type: typep(descriptor.FieldDescriptorProto_TYPE_STRING), Label: labelp(descriptor.FieldDescriptorProto_LABEL_REPEATED)},
			{Name: proto.String("resource_field"), Type: typep(descriptor.FieldDescriptorProto_TYPE_STRING)},
			{Name: proto.String("bytes"), Type: typep(descriptor.FieldDescriptorProto_TYPE_BYTES)},
			mapField,
		},
	}

	pageInType := &descriptor.DescriptorProto{
		Name: proto.String("PageInType"),
		Field: []*descriptor.FieldDescriptorProto{
			{Name: proto.String("a"), TypeName: proto.String("AType")},
			{Name: proto.String("page_size"), Type: typep(descriptor.FieldDescriptorProto_TYPE_INT32)},
			{Name: proto.String("page_token"), Type: typep(descriptor.FieldDescriptorProto_TYPE_STRING)},
		},
	}
	pageOutType := &descriptor.DescriptorProto{
		Name: proto.String("PageOutType"),
		Field: []*descriptor.FieldDescriptorProto{
			{Name: proto.String("a"), TypeName: proto.String("AType"), Label: labelp(descriptor.FieldDescriptorProto_LABEL_REPEATED)},
			{Name: proto.String("next_page_token"), Type: typep(descriptor.FieldDescriptorProto_TYPE_STRING)},
		},
	}

	lroInType := &descriptor.DescriptorProto{
		Name: proto.String("LroInType"),
	}
	lroReturnType := &descriptor.DescriptorProto{
		Name: proto.String("lroReturnType"),
		Field: []*descriptor.FieldDescriptorProto{
			{Name: proto.String("a"), TypeName: proto.String(".foo.AType")},
		},
	}
	lroMetadataType := &descriptor.DescriptorProto{
		Name: proto.String("lroMetadataType"),
	}

	serv := &descriptor.ServiceDescriptorProto{
		Name: proto.String("FooService"),
		Method: []*descriptor.MethodDescriptorProto{
			{
				Name:       proto.String("UnaryMethod"),
				InputType:  proto.String(".foo.InputType"),
				OutputType: proto.String(".foo.InputType"),
			},
			{
				Name:       proto.String("PagingMethod"),
				InputType:  proto.String(".foo.PageInType"),
				OutputType: proto.String(".foo.PageOutType"),
			},
			{
				Name:       proto.String("EmptyMethod"),
				InputType:  proto.String(".foo.PageInType"),
				OutputType: proto.String(".google.protobuf.Empty"),
			},
			{
				Name:       proto.String("LroMethod"),
				InputType:  proto.String(".foo.LroInType"),
				OutputType: proto.String(".google.longrunning.Operation"),
			},
		},
	}
	file := &descriptor.FileDescriptorProto{
		Options: &descriptor.FileOptions{
			GoPackage: proto.String("path.to/pb/foo;foo"),
		},
		Package:     proto.String("foo"),
		Service:     []*descriptor.ServiceDescriptorProto{serv},
		MessageType: []*descriptor.DescriptorProto{inType, aType, pageInType, pageOutType, lroInType, lroReturnType, lroMetadataType},
	}

	return &generator{
		clientPkg: pbinfo.ImportSpec{Path: "path.to/client/foo", Name: "foo"},
		imports:   map[pbinfo.ImportSpec]bool{},
		descInfo:  pbinfo.Of([]*descriptor.FileDescriptorProto{file}),
		gapic: GAPICConfig{
			Collections: []ResourceName{
				{EntityName: "foobar_thing", NamePattern: "foos/{foo}/bars/{bar}"},
			},
		},
	}
}

func labelp(l descriptor.FieldDescriptorProto_Label) *descriptor.FieldDescriptorProto_Label {
	return &l
}

// createMapTypeAndField creates the generated MapEntry protobuf message and the actual map field.
// If valueType is enum or message, vTypeName must not be empty.
func createMapTypeAndField(
	fieldName string,
	parentTyp string,
	keyType *descriptor.FieldDescriptorProto_Type,
	valueType *descriptor.FieldDescriptorProto_Type,
	vTypeName string,
) (*descriptor.DescriptorProto, *descriptor.FieldDescriptorProto) {

	var vf *descriptor.FieldDescriptorProto
	if vTypeName == "" {
		if *valueType == descriptor.FieldDescriptorProto_TYPE_ENUM {
			panic(errors.E(nil, "expecting non-empty enum type name"))
		}
		if *valueType == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
			panic(errors.E(nil, "expecting non-empty message type name"))
		}
		vf = &descriptor.FieldDescriptorProto{
			Name: proto.String("value"),
			Type: valueType,
		}
	} else {
		vf = &descriptor.FieldDescriptorProto{
			Name:     proto.String("value"),
			Type:     valueType,
			TypeName: proto.String(vTypeName),
		}
	}

	mapEntry := &descriptor.DescriptorProto{
		Name: proto.String("MapFieldEntry"),
		Field: []*descriptor.FieldDescriptorProto{
			{Name: proto.String("key"), Type: keyType},
			vf,
		},
		Options: &descriptor.MessageOptions{
			MapEntry: proto.Bool(true),
		},
	}

	mapField := &descriptor.FieldDescriptorProto{
		Name:     proto.String(fieldName),
		Type:     typep(descriptor.FieldDescriptorProto_TYPE_MESSAGE),
		TypeName: proto.String(parentTyp + ".MapFieldEntry"),
		Label:    labelp(descriptor.FieldDescriptorProto_LABEL_REPEATED),
	}

	return mapEntry, mapField
}

func compare(t *testing.T, g *generator, goldenPath string) {
	t.Helper()

	// Don't format. Format can change with Go version.
	gofmt := false
	year := 2018
	content, err := g.commit(gofmt, year)
	if err != nil {
		t.Fatal(err)
	}
	txtdiff.Diff(t, t.Name(), string(content), goldenPath)
}
