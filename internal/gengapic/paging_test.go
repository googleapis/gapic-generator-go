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

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/google/go-cmp/cmp"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoiface"
)

// TODO(dovs): Augment with map iterator
func TestIterTypeOf(t *testing.T) {
	msgType := &descriptor.DescriptorProto{
		Name: proto.String("Foo"),
	}
	mapEntry := &descriptor.DescriptorProto{
		Name:    proto.String("FooEntry"),
		Options: &descriptor.MessageOptions{MapEntry: proto.Bool(bool(true))},
	}
	g := &generator{
		aux: &auxTypes{
			iters: map[string]*iterType{},
		},
		descInfo: pbinfo.Info{
			Type: map[string]pbinfo.ProtoType{
				msgType.GetName():  msgType,
				mapEntry.GetName(): mapEntry,
			},
			ParentElement: map[pbinfo.ProtoType]pbinfo.ProtoType{},
			ParentFile: map[protoiface.MessageV1]*descriptor.FileDescriptorProto{
				msgType: {
					Options: &descriptor.FileOptions{
						GoPackage: proto.String("path/to/foo;foo"),
					},
				},
				mapEntry: {
					Options: &descriptor.FileOptions{
						GoPackage: proto.String("path/to/foo;foo"),
					},
				},
			},
		},
	}

	for i, tst := range []struct {
		field     *descriptor.FieldDescriptorProto
		want      iterType
		shouldErr bool
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
		{
			field: &descriptor.FieldDescriptorProto{
				Type:     typep(descriptor.FieldDescriptorProto_TYPE_MESSAGE),
				TypeName: proto.String(mapEntry.GetName()),
			},
			shouldErr: true,
		},
	} {
		g.descInfo.ParentElement[tst.field] = msgType
		got, err := g.iterTypeOf(tst.field)
		if tst.shouldErr {
			if err == nil {
				t.Errorf("field %v should error", tst.field)
			}
			continue
		}
		if err != nil {
			t.Error(err)
		} else if diff := cmp.Diff(tst.want, *got, cmp.AllowUnexported(*got)); diff != "" {
			t.Errorf("%d: (got=-, want=+):\n%s", i, diff)
		}
	}
}

func TestPagingField(t *testing.T) {
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
			{
				Name:   proto.String("page_size"),
				Number: proto.Int32(int32(1)),
				Type:   typep(descriptor.FieldDescriptorProto_TYPE_INT32),
			},
			{
				Name:   proto.String("page_token"),
				Number: proto.Int32(int32(2)),
				Type:   typep(descriptor.FieldDescriptorProto_TYPE_STRING),
			},
		},
	}
	validMaxResults := &descriptor.DescriptorProto{
		Name: proto.String("ValidMaxResultsRequest"),
		Field: []*descriptor.FieldDescriptorProto{
			{
				Name:   proto.String("max_results"),
				Number: proto.Int32(int32(1)),
				Type:   typep(descriptor.FieldDescriptorProto_TYPE_INT32),
			},
			{
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
			{
				Name:   proto.String("next_page_token"),
				Number: proto.Int32(int32(1)),
				Type:   typep(descriptor.FieldDescriptorProto_TYPE_STRING),
			},
			{
				Name:     proto.String("items"),
				Number:   proto.Int32(int32(2)),
				Type:     typep(descriptor.FieldDescriptorProto_TYPE_MESSAGE),
				TypeName: proto.String(".paging.RandomMessage"),
				Label:    labelp(descriptor.FieldDescriptorProto_LABEL_REPEATED),
			},
		},
	}
	mapEntry := &descriptor.DescriptorProto{
		Name:    proto.String("ItemsEntry"),
		Options: &descriptor.MessageOptions{MapEntry: proto.Bool(bool(true))},
		Field: []*descriptor.FieldDescriptorProto{
			{
				Name:   proto.String("key"),
				Number: proto.Int32(int32(1)),
				Type:   typep(descriptor.FieldDescriptorProto_TYPE_STRING),
			},
			{
				Name:     proto.String("value"),
				Number:   proto.Int32(int32(2)),
				Type:     typep(descriptor.FieldDescriptorProto_TYPE_MESSAGE),
				TypeName: proto.String(".paging.RandomMessage"),
			},
		},
	}
	validMap := &descriptor.DescriptorProto{
		Name: proto.String("ValidMapResponse"),
		Field: []*descriptor.FieldDescriptorProto{
			{
				Name:   proto.String("next_page_token"),
				Number: proto.Int32(int32(1)),
				Type:   typep(descriptor.FieldDescriptorProto_TYPE_STRING),
			},
			{
				Name:     proto.String("items"),
				Number:   proto.Int32(int32(2)),
				Type:     typep(descriptor.FieldDescriptorProto_TYPE_MESSAGE),
				Label:    labelp(descriptor.FieldDescriptorProto_LABEL_REPEATED),
				TypeName: proto.String(".paging.ItemsEntry"),
			},
		},
	}
	invalidRsp := &descriptor.DescriptorProto{Name: proto.String("InvalidResponse")}
	multipleRepeated := &descriptor.DescriptorProto{
		Name: proto.String("MultipleRepeatedResponse"),
		Field: []*descriptor.FieldDescriptorProto{
			{
				Name:   proto.String("next_page_token"),
				Number: proto.Int32(int32(1)),
				Type:   typep(descriptor.FieldDescriptorProto_TYPE_STRING),
			},
			{
				Name:     proto.String("items"),
				Number:   proto.Int32(int32(2)),
				Type:     typep(descriptor.FieldDescriptorProto_TYPE_MESSAGE),
				TypeName: proto.String(".paging.RandomMessage"),
				Label:    labelp(descriptor.FieldDescriptorProto_LABEL_REPEATED),
			},
			{
				Name:     proto.String("items_2"),
				Number:   proto.Int32(int32(3)),
				Type:     typep(descriptor.FieldDescriptorProto_TYPE_MESSAGE),
				TypeName: proto.String(".paging.RandomMessage"),
				Label:    labelp(descriptor.FieldDescriptorProto_LABEL_REPEATED),
			},
		},
	}
	tooManyRepeated := &descriptor.DescriptorProto{
		Name: proto.String("TooManyRepeatedResponse"),
		Field: []*descriptor.FieldDescriptorProto{
			{
				Name:   proto.String("next_page_token"),
				Number: proto.Int32(int32(1)),
				Type:   typep(descriptor.FieldDescriptorProto_TYPE_STRING),
			},
			{
				Name:     proto.String("items"),
				Number:   proto.Int32(int32(3)), // Note that the "first" repeated field has a higher field number.
				Type:     typep(descriptor.FieldDescriptorProto_TYPE_MESSAGE),
				TypeName: proto.String(".paging.RandomMessage"),
				Label:    labelp(descriptor.FieldDescriptorProto_LABEL_REPEATED),
			},
			{
				Name:     proto.String("items_2"),
				Number:   proto.Int32(int32(2)),
				Type:     typep(descriptor.FieldDescriptorProto_TYPE_MESSAGE),
				TypeName: proto.String(".paging.RandomMessage"),
				Label:    labelp(descriptor.FieldDescriptorProto_LABEL_REPEATED),
			},
		},
	}
	tooManyMap := &descriptor.DescriptorProto{
		Name: proto.String("TooManyMapResponse"),
		Field: []*descriptor.FieldDescriptorProto{
			{
				Name:   proto.String("next_page_token"),
				Number: proto.Int32(int32(1)),
				Type:   typep(descriptor.FieldDescriptorProto_TYPE_STRING),
			},
			{
				Name:     proto.String("items"),
				Number:   proto.Int32(int32(3)), // Note that the "first" repeated field has a higher field number.
				Type:     typep(descriptor.FieldDescriptorProto_TYPE_MESSAGE),
				Label:    labelp(descriptor.FieldDescriptorProto_LABEL_REPEATED),
				TypeName: proto.String(".paging.ItemsEntry"),
			},
			{
				Name:     proto.String("items_2"),
				Number:   proto.Int32(int32(2)),
				Type:     typep(descriptor.FieldDescriptorProto_TYPE_MESSAGE),
				Label:    labelp(descriptor.FieldDescriptorProto_LABEL_REPEATED),
				TypeName: proto.String(".paging.ItemsEntry"),
			},
		},
	}
	noRepeatedField := &descriptor.DescriptorProto{
		Name: proto.String("NoRepeatedFieldResponse"),
		Field: []*descriptor.FieldDescriptorProto{
			{
				Name:   proto.String("next_page_token"),
				Number: proto.Int32(int32(1)),
				Type:   typep(descriptor.FieldDescriptorProto_TYPE_STRING),
			},
		},
	}
	noNextPageToken := &descriptor.DescriptorProto{
		Name: proto.String("NoNextPageTokenResponse"),
		Field: []*descriptor.FieldDescriptorProto{
			{
				Name:     proto.String("items"),
				Number:   proto.Int32(int32(1)),
				Type:     typep(descriptor.FieldDescriptorProto_TYPE_MESSAGE),
				TypeName: proto.String(".paging.RandomMessage"),
				Label:    labelp(descriptor.FieldDescriptorProto_LABEL_REPEATED),
			},
		},
	}

	// Methods
	validPageSizeMthd := &descriptor.MethodDescriptorProto{
		Name:       proto.String("ValidPageSize"),
		InputType:  proto.String(".paging.ValidPageSizeRequest"),
		OutputType: proto.String(".paging.ValidRepeatedResponse"),
	}
	validPageSizeMultipleMthd := &descriptor.MethodDescriptorProto{
		Name:       proto.String("ValidPageSizeMultiple"),
		InputType:  proto.String(".paging.ValidPageSizeRequest"),
		OutputType: proto.String(".paging.MultipleRepeatedResponse"),
	}
	validMaxResultsRepeatedMthd := &descriptor.MethodDescriptorProto{
		Name:       proto.String("ValidMaxResults"),
		InputType:  proto.String(".paging.ValidMaxResultsRequest"),
		OutputType: proto.String(".paging.ValidRepeatedResponse"),
	}
	validMaxResultsMthd := &descriptor.MethodDescriptorProto{
		Name:       proto.String("ValidMaxResults"),
		InputType:  proto.String(".paging.ValidMaxResultsRequest"),
		OutputType: proto.String(".paging.ValidMapResponse"),
	}
	validPageSizeMapMthd := &descriptor.MethodDescriptorProto{
		Name:       proto.String("ValidPageSize"),
		InputType:  proto.String(".paging.ValidPageSizeRequest"),
		OutputType: proto.String(".paging.ValidMapResponse"),
	}
	clientStreamingMthd := &descriptor.MethodDescriptorProto{
		Name:            proto.String("ClientStreaming"),
		InputType:       proto.String(".paging.ValidPageSizeRequest"),
		OutputType:      proto.String(".paging.ValidRepeatedResponse"),
		ClientStreaming: proto.Bool(bool(true)),
	}
	serverStreamingMthd := &descriptor.MethodDescriptorProto{
		Name:            proto.String("ServerStreaming"),
		InputType:       proto.String(".paging.ValidPageSizeRequest"),
		OutputType:      proto.String(".paging.ValidRepeatedResponse"),
		ServerStreaming: proto.Bool(bool(true)),
	}
	tooManyRepeatedMthd := &descriptor.MethodDescriptorProto{
		Name:       proto.String("TooManyRepeated"),
		InputType:  proto.String(".paging.ValidPageSizeRequest"),
		OutputType: proto.String(".paging.TooManyRepeatedResponse"),
	}
	tooManyMapMthd := &descriptor.MethodDescriptorProto{
		Name:       proto.String("TooManyMap"),
		InputType:  proto.String(".paging.ValidPageSizeRequest"),
		OutputType: proto.String(".paging.TooManyMapResponse"),
	}
	noNextPageTokenMthd := &descriptor.MethodDescriptorProto{
		Name:       proto.String("NoNextPageToken"),
		InputType:  proto.String(".paging.ValidPageSizeRequest"),
		OutputType: proto.String(".paging.NoNextPageTokenResponse"),
	}
	noRepeatedFieldMthd := &descriptor.MethodDescriptorProto{
		Name:       proto.String("NoRepeatedField"),
		InputType:  proto.String(".paging.ValidPageSizeRequest"),
		OutputType: proto.String(".paging.NoRepeatedFieldResponse"),
	}

	file := &descriptor.FileDescriptorProto{
		Package: proto.String("paging"),
		Options: &descriptor.FileOptions{
			GoPackage: proto.String("mypackage"),
		},
		MessageType: []*descriptor.DescriptorProto{
			invalidRsp,
			mapEntry,
			multipleRepeated,
			noNextPageToken,
			noRepeatedField,
			randomMessage,
			tooManyMap,
			tooManyRepeated,
			validMap,
			validMaxResults,
			validPageSize,
			validRepeated,
		},
		Service: []*descriptor.ServiceDescriptorProto{
			{
				Name: proto.String("TestService"),
				Method: []*descriptor.MethodDescriptorProto{
					clientStreamingMthd,
					noNextPageTokenMthd,
					noRepeatedFieldMthd,
					serverStreamingMthd,
					tooManyMapMthd,
					tooManyRepeatedMthd,
					validMaxResultsMthd,
					validMaxResultsRepeatedMthd,
					validPageSizeMapMthd,
					validPageSizeMthd,
					validPageSizeMultipleMthd,
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
		sizeField *descriptor.FieldDescriptorProto // A nil field means this is not a paged method
		iterField *descriptor.FieldDescriptorProto // A nil field means this is not a paged method
	}{
		{mthd: clientStreamingMthd},
		{mthd: serverStreamingMthd},
		{mthd: tooManyMapMthd},
		{mthd: tooManyRepeatedMthd},
		{mthd: noNextPageTokenMthd},
		{mthd: noRepeatedFieldMthd},
		{mthd: validMaxResultsRepeatedMthd, sizeField: validMaxResults.GetField()[0], iterField: validRepeated.GetField()[1]},
		{mthd: validPageSizeMapMthd, sizeField: validPageSize.GetField()[0], iterField: validMap.GetField()[1]},
		{mthd: validPageSizeMthd, sizeField: validPageSize.GetField()[0], iterField: validRepeated.GetField()[1]},
		{mthd: validPageSizeMultipleMthd, sizeField: validPageSize.GetField()[0], iterField: multipleRepeated.GetField()[1]},
		{mthd: validMaxResultsMthd, sizeField: validMaxResults.GetField()[0], iterField: validMap.GetField()[1]},
	} {
		actualIter, actualSize, err := g.getPagingFields(tst.mthd)
		if actualSize != tst.sizeField {
			t.Errorf("test %s page size field: got %s, want %s, err %v", tst.mthd.GetName(), actualSize, tst.sizeField, err)
		}

		if actualIter != tst.iterField {
			t.Errorf("test %s iter field: got %s, want %s, err %v", tst.mthd.GetName(), actualIter, tst.iterField, err)
		}
	}
}
