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
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/google/go-cmp/cmp"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
)

func TestPagingField(t *testing.T) {
	resField := &descriptor.FieldDescriptorProto{
		Name:   proto.String("resource"),
		Type:   typep(descriptor.FieldDescriptorProto_TYPE_STRING),
		Label:  labelp(descriptor.FieldDescriptorProto_LABEL_REPEATED),
		Number: proto.Int32(1),
	}
	mapField := &descriptor.FieldDescriptorProto{
		Name:     proto.String("map"),
		Type:     typep(descriptor.FieldDescriptorProto_TYPE_MESSAGE),
		Label:    labelp(descriptor.FieldDescriptorProto_LABEL_REPEATED),
		TypeName: proto.String("MapEntry"),
		Number:   proto.Int32(1),
	}

	otherRepField := &descriptor.FieldDescriptorProto{
		Name:   proto.String("unreachable"),
		Type:   typep(descriptor.FieldDescriptorProto_TYPE_STRING),
		Label:  labelp(descriptor.FieldDescriptorProto_LABEL_REPEATED),
		Number: proto.Int32(0),
	}

	g := &generator{}
	g.descInfo.Type = map[string]pbinfo.ProtoType{
		"Foo": &descriptor.DescriptorProto{
			Name: proto.String("Foo"),
		},
		"MapEntry": &descriptor.DescriptorProto{
			Name: proto.String("MapEntry"),
			Field: []*descriptor.FieldDescriptorProto{
				{
					Name:   proto.String("key"),
					Type:   typep(descriptor.FieldDescriptorProto_TYPE_STRING),
					Number: proto.Int32(int32(1)),
				},
				{
					Name:   proto.String("value"),
					Type:   typep(descriptor.FieldDescriptorProto_TYPE_STRING),
					Number: proto.Int32(int32(2)),
				},
			},
			Options: &descriptor.MessageOptions{MapEntry: proto.Bool(bool(true))},
		},
		"PageIn": &descriptor.DescriptorProto{
			Name: proto.String("PageIn"),
			Field: []*descriptor.FieldDescriptorProto{
				{
					Name:  proto.String("page_size"),
					Type:  typep(descriptor.FieldDescriptorProto_TYPE_INT32),
					Label: labelp(descriptor.FieldDescriptorProto_LABEL_OPTIONAL),
				},
				{
					Name:  proto.String("page_token"),
					Type:  typep(descriptor.FieldDescriptorProto_TYPE_STRING),
					Label: labelp(descriptor.FieldDescriptorProto_LABEL_OPTIONAL),
				},
			},
		},
		"PageOut": &descriptor.DescriptorProto{
			Name: proto.String("PageOut"),
			Field: []*descriptor.FieldDescriptorProto{
				{
					Name:  proto.String("next_page_token"),
					Type:  typep(descriptor.FieldDescriptorProto_TYPE_STRING),
					Label: labelp(descriptor.FieldDescriptorProto_LABEL_OPTIONAL),
				},
				resField,
			},
		},
		"MapOut": &descriptor.DescriptorProto{
			Name: proto.String("MapOut"),
			Field: []*descriptor.FieldDescriptorProto{
				{
					Name:  proto.String("next_page_token"),
					Type:  typep(descriptor.FieldDescriptorProto_TYPE_STRING),
					Label: labelp(descriptor.FieldDescriptorProto_LABEL_OPTIONAL),
				},
				mapField,
			},
		},
		"NoRepeatedOut": &descriptor.DescriptorProto{
			Name: proto.String("NoRepeatedOut"),
			Field: []*descriptor.FieldDescriptorProto{
				{
					Name:  proto.String("next_page_token"),
					Type:  typep(descriptor.FieldDescriptorProto_TYPE_STRING),
					Label: labelp(descriptor.FieldDescriptorProto_LABEL_OPTIONAL),
				},
				// No repeated field
			},
		},
		"BadPageOut": &descriptor.DescriptorProto{
			Name: proto.String("BadPageOut"),
			Field: []*descriptor.FieldDescriptorProto{
				{
					Name:  proto.String("next_page_token"),
					Type:  typep(descriptor.FieldDescriptorProto_TYPE_STRING),
					Label: labelp(descriptor.FieldDescriptorProto_LABEL_OPTIONAL),
				},
				// resource is not first repeated field in message
				resField,
				otherRepField,
			},
		},
	}

	for _, tst := range []struct {
		in, out string
		field   *descriptor.FieldDescriptorProto
		err     bool
	}{
		{
			in:  "Foo",
			out: "Foo",
		},
		{
			in:  "PageIn",
			out: "Foo",
		},
		{
			in:  "Foo",
			out: "PageOut",
		},
		{
			in:    "PageIn",
			out:   "PageOut",
			field: resField,
		},
		{
			in:  "PageIn",
			out: "NoRepeatedOut",
		},
		{
			in:  "PageIn",
			out: "MapOut",
		},
		{
			in:  "PageIn",
			out: "BadPageOut",
			err: true,
		},
	} {
		meth := &descriptor.MethodDescriptorProto{
			Name:       proto.String("TestPagingField"),
			InputType:  proto.String(tst.in),
			OutputType: proto.String(tst.out),
		}
		f, err := g.pagingField(meth)
		if tst.err && err == nil {
			t.Errorf("pagingField(%v)=%v, expected error", meth, f)
		} else if !tst.err && err != nil {
			t.Errorf("pagingField(%v) errors %q, expected %v", meth, err, tst.field)
		} else if f != tst.field {
			t.Errorf("pagingField(%v)=%v, want %v", meth, f, tst.field)
		}
	}
}

func TestIterTypeOf(t *testing.T) {
	msgType := &descriptor.DescriptorProto{
		Name: proto.String("Foo"),
	}
	g := &generator{
		aux: &auxTypes{
			iters: map[string]*iterType{},
		},
		descInfo: pbinfo.Info{
			Type: map[string]pbinfo.ProtoType{
				msgType.GetName(): msgType,
			},
			ParentElement: map[pbinfo.ProtoType]pbinfo.ProtoType{},
			ParentFile: map[proto.Message]*descriptor.FileDescriptorProto{
				msgType: {
					Options: &descriptor.FileOptions{
						GoPackage: proto.String("path/to/foo;foo"),
					},
				},
			},
		},
	}

	for i, tst := range []struct {
		field *descriptor.FieldDescriptorProto
		want  iterType
	}{
		{
			field: &descriptor.FieldDescriptorProto{
				Type: typep(descriptor.FieldDescriptorProto_TYPE_STRING),
			},
			want: iterType{
				iterTypeName: "StringIterator",
				elemTypeName: "string",
			},
		},
		{
			field: &descriptor.FieldDescriptorProto{
				Type: typep(descriptor.FieldDescriptorProto_TYPE_BYTES),
			},
			want: iterType{
				iterTypeName: "BytesIterator",
				elemTypeName: "[]byte",
			},
		},
		{
			field: &descriptor.FieldDescriptorProto{
				Type:     typep(descriptor.FieldDescriptorProto_TYPE_MESSAGE),
				TypeName: proto.String(msgType.GetName()),
			},
			want: iterType{
				iterTypeName: "FooIterator",
				elemTypeName: "*foopb.Foo",
				elemImports:  []pbinfo.ImportSpec{{Name: "foopb", Path: "path/to/foo"}},
			},
		},
	} {
		g.descInfo.ParentElement[tst.field] = msgType
		got, err := g.iterTypeOf(tst.field)
		if err != nil {
			t.Error(err)
		} else if diff := cmp.Diff(tst.want, *got, cmp.AllowUnexported(*got)); diff != "" {
			t.Errorf("%d: (got=-, want=+):\n%s", i, diff)
		}
	}
}

func TestDiregapicPagingField(t *testing.T) {
	g := generator{
		apiName: "Awesome API",
		imports: map[pbinfo.ImportSpec]bool{},
		opts:    &options{transports: []transport{rest}},
	}

	// The predicate is simple:
	// * No streaming in the method
	// * Request has a int32 page_size field XOR a int32 max_results field
	// * Request has a string page_token field
	// * Response has a string next_page_token field
	// * Response has one and only one repeated or map<string, *> field

	// Messages
	validPageSize := &descriptor.DescriptorProto{
		Name: proto.String("ValidPageSizeRequest"),
		Field: []*descriptor.FieldDescriptorProto{
			&descriptor.FieldDescriptorProto{
				Name:   proto.String("page_size"),
				Number: proto.Int32(int32(1)),
				Type:   typep(descriptor.FieldDescriptorProto_TYPE_INT32),
			},
			&descriptor.FieldDescriptorProto{
				Name:   proto.String("page_token"),
				Number: proto.Int32(int32(2)),
				Type:   typep(descriptor.FieldDescriptorProto_TYPE_STRING),
			},
		},
	}
	validMaxResults := &descriptor.DescriptorProto{
		Name: proto.String("ValidMaxResultsRequest"),
		Field: []*descriptor.FieldDescriptorProto{
			&descriptor.FieldDescriptorProto{
				Name:   proto.String("max_results"),
				Number: proto.Int32(int32(1)),
				Type:   typep(descriptor.FieldDescriptorProto_TYPE_INT32),
			},
			&descriptor.FieldDescriptorProto{
				Name:   proto.String("page_token"),
				Number: proto.Int32(int32(2)),
				Type:   typep(descriptor.FieldDescriptorProto_TYPE_STRING),
			},
		},
	}
	randomMessage := &descriptor.DescriptorProto{Name: proto.String("RandomMessage")}
	validRepeated := &descriptor.DescriptorProto{
		Name: proto.String("ValidRepeatedResponse"),
		Field: []*descriptor.FieldDescriptorProto{
			&descriptor.FieldDescriptorProto{
				Name:   proto.String("next_page_token"),
				Number: proto.Int32(int32(1)),
				Type:   typep(descriptor.FieldDescriptorProto_TYPE_STRING),
			},
			&descriptor.FieldDescriptorProto{
				Name:     proto.String("items"),
				Number:   proto.Int32(int32(2)),
				Type:     typep(descriptor.FieldDescriptorProto_TYPE_MESSAGE),
				TypeName: proto.String(".diregapic.paging.RandomMessage"),
				Label:    labelp(descriptor.FieldDescriptorProto_LABEL_REPEATED),
			},
		},
	}
	mapEntry := &descriptor.DescriptorProto{
		Name:    proto.String("ItemsEntry"),
		Options: &descriptor.MessageOptions{MapEntry: proto.Bool(bool(true))},
		Field: []*descriptor.FieldDescriptorProto{
			&descriptor.FieldDescriptorProto{
				Name:   proto.String("key"),
				Number: proto.Int32(int32(1)),
				Type:   typep(descriptor.FieldDescriptorProto_TYPE_STRING),
			},
			&descriptor.FieldDescriptorProto{
				Name:     proto.String("value"),
				Number:   proto.Int32(int32(2)),
				Type:     typep(descriptor.FieldDescriptorProto_TYPE_MESSAGE),
				TypeName: proto.String(".diregapic.paging.RandomMessage"),
			},
		},
	}
	validMap := &descriptor.DescriptorProto{
		Name: proto.String("ValidMapResponse"),
		Field: []*descriptor.FieldDescriptorProto{
			&descriptor.FieldDescriptorProto{
				Name:   proto.String("next_page_token"),
				Number: proto.Int32(int32(1)),
				Type:   typep(descriptor.FieldDescriptorProto_TYPE_STRING),
			},
			&descriptor.FieldDescriptorProto{
				Name:     proto.String("items"),
				Number:   proto.Int32(int32(2)),
				Type:     typep(descriptor.FieldDescriptorProto_TYPE_MESSAGE),
				Label:    labelp(descriptor.FieldDescriptorProto_LABEL_REPEATED),
				TypeName: proto.String(".diregapic.paging.ItemsEntry"),
			},
		},
	}
	invalidRsp := &descriptor.DescriptorProto{Name: proto.String("InvalidResponse")}
	tooManyRepeated := &descriptor.DescriptorProto{
		Name: proto.String("TooManyRepeatedResponse"),
		Field: []*descriptor.FieldDescriptorProto{
			&descriptor.FieldDescriptorProto{
				Name:   proto.String("next_page_token"),
				Number: proto.Int32(int32(1)),
				Type:   typep(descriptor.FieldDescriptorProto_TYPE_STRING),
			},
			&descriptor.FieldDescriptorProto{
				Name:     proto.String("items"),
				Number:   proto.Int32(int32(2)),
				Type:     typep(descriptor.FieldDescriptorProto_TYPE_MESSAGE),
				TypeName: proto.String(".diregapic.paging.RandomMessage"),
				Label:    labelp(descriptor.FieldDescriptorProto_LABEL_REPEATED),
			},
			&descriptor.FieldDescriptorProto{
				Name:     proto.String("items_2"),
				Number:   proto.Int32(int32(3)),
				Type:     typep(descriptor.FieldDescriptorProto_TYPE_MESSAGE),
				TypeName: proto.String(".diregapic.paging.RandomMessage"),
				Label:    labelp(descriptor.FieldDescriptorProto_LABEL_REPEATED),
			},
		},
	}
	tooManyMap := &descriptor.DescriptorProto{
		Name: proto.String("TooManyMapResponse"),
		Field: []*descriptor.FieldDescriptorProto{
			&descriptor.FieldDescriptorProto{
				Name:   proto.String("next_page_token"),
				Number: proto.Int32(int32(1)),
				Type:   typep(descriptor.FieldDescriptorProto_TYPE_STRING),
			},
			&descriptor.FieldDescriptorProto{
				Name:     proto.String("items"),
				Number:   proto.Int32(int32(2)),
				Type:     typep(descriptor.FieldDescriptorProto_TYPE_MESSAGE),
				Label:    labelp(descriptor.FieldDescriptorProto_LABEL_REPEATED),
				TypeName: proto.String(".diregapic.paging.ItemsEntry"),
			},
			&descriptor.FieldDescriptorProto{
				Name:     proto.String("items_2"),
				Number:   proto.Int32(int32(3)),
				Type:     typep(descriptor.FieldDescriptorProto_TYPE_MESSAGE),
				Label:    labelp(descriptor.FieldDescriptorProto_LABEL_REPEATED),
				TypeName: proto.String(".diregapic.paging.ItemsEntry"),
			},
		},
	}

	// Methods
	validPageSizeMthd := &descriptor.MethodDescriptorProto{
		Name:       proto.String("ValidPageSize"),
		InputType:  proto.String(".diregapic.paging.ValidPageSizeRequest"),
		OutputType: proto.String(".diregapic.paging.ValidRepeatedResponse"),
	}
	validMaxResultsMthd := &descriptor.MethodDescriptorProto{
		Name:       proto.String("ValidMaxResults"),
		InputType:  proto.String(".diregapic.paging.ValidMaxResultsRequest"),
		OutputType: proto.String(".diregapic.paging.ValidMapResponse"),
	}
	clientStreamingMthd := &descriptor.MethodDescriptorProto{
		Name:            proto.String("ClientStreaming"),
		InputType:       proto.String(".diregapic.paging.ValidPageSizeRequest"),
		OutputType:      proto.String(".diregapic.paging.ValidRepeatedResponse"),
		ClientStreaming: proto.Bool(bool(true)),
	}
	serverStreamingMthd := &descriptor.MethodDescriptorProto{
		Name:            proto.String("ServerStreaming"),
		InputType:       proto.String(".diregapic.paging.ValidPageSizeRequest"),
		OutputType:      proto.String(".diregapic.paging.ValidRepeatedResponse"),
		ServerStreaming: proto.Bool(bool(true)),
	}
	tooManyRepeatedMthd := &descriptor.MethodDescriptorProto{
		Name:       proto.String("TooManyRepeated"),
		InputType:  proto.String(".diregapic.paging.ValidPageSizeRequest"),
		OutputType: proto.String(".diregapic.paging.TooManyRepeatedResponse"),
	}
	tooManyMapMthd := &descriptor.MethodDescriptorProto{
		Name:       proto.String("TooManyMap"),
		InputType:  proto.String(".diregapic.paging.ValidPageSizeRequest"),
		OutputType: proto.String(".diregapic.paging.TooManyMapResponse"),
	}

	file := &descriptor.FileDescriptorProto{
		Package: proto.String("diregapic.paging"),
		Options: &descriptor.FileOptions{
			GoPackage: proto.String("mypackage"),
		},
		MessageType: []*descriptor.DescriptorProto{
			invalidRsp,
			mapEntry,
			randomMessage,
			tooManyMap,
			tooManyRepeated,
			validMap,
			validMaxResults,
			validPageSize,
			validRepeated,
		},
		Service: []*descriptor.ServiceDescriptorProto{
			&descriptor.ServiceDescriptorProto{
				Name: proto.String("TestService"),
				Method: []*descriptor.MethodDescriptorProto{
					clientStreamingMthd,
					serverStreamingMthd,
					tooManyMapMthd,
					tooManyRepeatedMthd,
					validMaxResultsMthd,
					validPageSizeMthd,
				},
			},
		},
	}

	req := plugin.CodeGeneratorRequest{
		Parameter: proto.String("go-gapic-package=path;mypackage,transport=rest"),
		ProtoFile: []*descriptor.FileDescriptorProto{file},
	}
	g.init(&req)

	for _, tst := range []struct {
		mthd      *descriptor.MethodDescriptorProto
		sizeField *descriptor.FieldDescriptorProto // A nil field means this is not a Diregapic paged method
		iterField *descriptor.FieldDescriptorProto // A nil field means this is not a Diregapic paged method
	}{
		{mthd: clientStreamingMthd},
		{mthd: serverStreamingMthd},
		{mthd: tooManyMapMthd},
		{mthd: tooManyRepeatedMthd},
		{mthd: validPageSizeMthd, sizeField: validPageSize.GetField()[0], iterField: validRepeated.GetField()[1]},
		{mthd: validMaxResultsMthd, sizeField: validMaxResults.GetField()[0], iterField: validMap.GetField()[1]},
	} {
		actualIter, actualSize, err := g.diregapicPagingField(tst.mthd)
		if actualSize != tst.sizeField {
			t.Errorf("test %s page size field: got %s, want %s, err %q", tst.mthd.GetName(), actualSize, tst.sizeField, err)
		}

		if actualIter != tst.iterField {
			t.Errorf("test %s iter field: got %s, want %s, err %q", tst.mthd.GetName(), actualIter, tst.iterField, err)
		}
	}
}
