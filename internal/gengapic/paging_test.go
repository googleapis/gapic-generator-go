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

	"github.com/google/go-cmp/cmp"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

// TODO(dovs): Augment with map iterator
func TestIterTypeOf(t *testing.T) {
	msgType := &descriptorpb.DescriptorProto{
		Name: proto.String("Foo"),
	}
	mapEntry := &descriptorpb.DescriptorProto{
		Name:    proto.String("FooEntry"),
		Options: &descriptorpb.MessageOptions{MapEntry: proto.Bool(bool(true))},
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
			ParentFile: map[protoreflect.ProtoMessage]*descriptorpb.FileDescriptorProto{
				msgType: {
					Options: &descriptorpb.FileOptions{
						GoPackage: proto.String("path/to/foo;foo"),
					},
				},
				mapEntry: {
					Options: &descriptorpb.FileOptions{
						GoPackage: proto.String("path/to/foo;foo"),
					},
				},
			},
		},
	}

	for i, tst := range []struct {
		field     *descriptorpb.FieldDescriptorProto
		want      iterType
		shouldErr bool
	}{
		{
			field: &descriptorpb.FieldDescriptorProto{
				Type: typep(descriptorpb.FieldDescriptorProto_TYPE_STRING),
			},
			want: iterType{
				iterTypeName: "StringIterator",
				elemTypeName: "string",
			},
		},
		{
			field: &descriptorpb.FieldDescriptorProto{
				Type: typep(descriptorpb.FieldDescriptorProto_TYPE_BYTES),
			},
			want: iterType{
				iterTypeName: "BytesIterator",
				elemTypeName: "[]byte",
			},
		},
		{
			field: &descriptorpb.FieldDescriptorProto{
				Type:     typep(descriptorpb.FieldDescriptorProto_TYPE_MESSAGE),
				TypeName: proto.String(msgType.GetName()),
			},
			want: iterType{
				iterTypeName: "FooIterator",
				elemTypeName: "*foopb.Foo",
				elemImports:  []pbinfo.ImportSpec{{Name: "foopb", Path: "path/to/foo"}},
			},
		},
		{
			field: &descriptorpb.FieldDescriptorProto{
				Type:     typep(descriptorpb.FieldDescriptorProto_TYPE_MESSAGE),
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
	// The predicate is simple:
	// * No streaming in the method
	// * Request has a int32 page_size field XOR a int32 max_results field
	// * Request has a string page_token field
	// * Response has a string next_page_token field
	// * Response has one and only one repeated or map<string, *> field

	// Messages
	validPageSize := &descriptorpb.DescriptorProto{
		Name: proto.String("ValidPageSizeRequest"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:   proto.String("page_size"),
				Number: proto.Int32(int32(1)),
				Type:   typep(descriptorpb.FieldDescriptorProto_TYPE_INT32),
			},
			{
				Name:   proto.String("page_token"),
				Number: proto.Int32(int32(2)),
				Type:   typep(descriptorpb.FieldDescriptorProto_TYPE_STRING),
			},
		},
	}
	validMaxResults := &descriptorpb.DescriptorProto{
		Name: proto.String("ValidMaxResultsRequest"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:   proto.String("max_results"),
				Number: proto.Int32(int32(1)),
				Type:   typep(descriptorpb.FieldDescriptorProto_TYPE_INT32),
			},
			{
				Name:   proto.String("page_token"),
				Number: proto.Int32(int32(2)),
				Type:   typep(descriptorpb.FieldDescriptorProto_TYPE_STRING),
			},
		},
	}
	randomMessage := &descriptorpb.DescriptorProto{Name: proto.String("RandomMessage")}
	validRepeated := &descriptorpb.DescriptorProto{
		Name: proto.String("ValidRepeatedResponse"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:   proto.String("next_page_token"),
				Number: proto.Int32(int32(1)),
				Type:   typep(descriptorpb.FieldDescriptorProto_TYPE_STRING),
			},
			{
				Name:     proto.String("items"),
				Number:   proto.Int32(int32(2)),
				Type:     typep(descriptorpb.FieldDescriptorProto_TYPE_MESSAGE),
				TypeName: proto.String(".paging.RandomMessage"),
				Label:    labelp(descriptorpb.FieldDescriptorProto_LABEL_REPEATED),
			},
		},
	}
	mapEntry := &descriptorpb.DescriptorProto{
		Name:    proto.String("ItemsEntry"),
		Options: &descriptorpb.MessageOptions{MapEntry: proto.Bool(bool(true))},
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:   proto.String("key"),
				Number: proto.Int32(int32(1)),
				Type:   typep(descriptorpb.FieldDescriptorProto_TYPE_STRING),
			},
			{
				Name:     proto.String("value"),
				Number:   proto.Int32(int32(2)),
				Type:     typep(descriptorpb.FieldDescriptorProto_TYPE_MESSAGE),
				TypeName: proto.String(".paging.RandomMessage"),
			},
		},
	}
	validMap := &descriptorpb.DescriptorProto{
		Name: proto.String("ValidMapResponse"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:   proto.String("next_page_token"),
				Number: proto.Int32(int32(1)),
				Type:   typep(descriptorpb.FieldDescriptorProto_TYPE_STRING),
			},
			{
				Name:     proto.String("items"),
				Number:   proto.Int32(int32(2)),
				Type:     typep(descriptorpb.FieldDescriptorProto_TYPE_MESSAGE),
				Label:    labelp(descriptorpb.FieldDescriptorProto_LABEL_REPEATED),
				TypeName: proto.String(".paging.ItemsEntry"),
			},
		},
	}
	invalidRsp := &descriptorpb.DescriptorProto{Name: proto.String("InvalidResponse")}
	multipleRepeated := &descriptorpb.DescriptorProto{
		Name: proto.String("MultipleRepeatedResponse"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:   proto.String("next_page_token"),
				Number: proto.Int32(int32(1)),
				Type:   typep(descriptorpb.FieldDescriptorProto_TYPE_STRING),
			},
			{
				Name:     proto.String("items"),
				Number:   proto.Int32(int32(2)),
				Type:     typep(descriptorpb.FieldDescriptorProto_TYPE_MESSAGE),
				TypeName: proto.String(".paging.RandomMessage"),
				Label:    labelp(descriptorpb.FieldDescriptorProto_LABEL_REPEATED),
			},
			{
				Name:     proto.String("items_2"),
				Number:   proto.Int32(int32(3)),
				Type:     typep(descriptorpb.FieldDescriptorProto_TYPE_MESSAGE),
				TypeName: proto.String(".paging.RandomMessage"),
				Label:    labelp(descriptorpb.FieldDescriptorProto_LABEL_REPEATED),
			},
		},
	}
	tooManyRepeated := &descriptorpb.DescriptorProto{
		Name: proto.String("TooManyRepeatedResponse"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:   proto.String("next_page_token"),
				Number: proto.Int32(int32(1)),
				Type:   typep(descriptorpb.FieldDescriptorProto_TYPE_STRING),
			},
			{
				Name:     proto.String("items"),
				Number:   proto.Int32(int32(3)), // Note that the "first" repeated field has a higher field number.
				Type:     typep(descriptorpb.FieldDescriptorProto_TYPE_MESSAGE),
				TypeName: proto.String(".paging.RandomMessage"),
				Label:    labelp(descriptorpb.FieldDescriptorProto_LABEL_REPEATED),
			},
			{
				Name:     proto.String("items_2"),
				Number:   proto.Int32(int32(2)),
				Type:     typep(descriptorpb.FieldDescriptorProto_TYPE_MESSAGE),
				TypeName: proto.String(".paging.RandomMessage"),
				Label:    labelp(descriptorpb.FieldDescriptorProto_LABEL_REPEATED),
			},
		},
	}
	tooManyMap := &descriptorpb.DescriptorProto{
		Name: proto.String("TooManyMapResponse"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:   proto.String("next_page_token"),
				Number: proto.Int32(int32(1)),
				Type:   typep(descriptorpb.FieldDescriptorProto_TYPE_STRING),
			},
			{
				Name:     proto.String("items"),
				Number:   proto.Int32(int32(3)), // Note that the "first" repeated field has a higher field number.
				Type:     typep(descriptorpb.FieldDescriptorProto_TYPE_MESSAGE),
				Label:    labelp(descriptorpb.FieldDescriptorProto_LABEL_REPEATED),
				TypeName: proto.String(".paging.ItemsEntry"),
			},
			{
				Name:     proto.String("items_2"),
				Number:   proto.Int32(int32(2)),
				Type:     typep(descriptorpb.FieldDescriptorProto_TYPE_MESSAGE),
				Label:    labelp(descriptorpb.FieldDescriptorProto_LABEL_REPEATED),
				TypeName: proto.String(".paging.ItemsEntry"),
			},
		},
	}
	noRepeatedField := &descriptorpb.DescriptorProto{
		Name: proto.String("NoRepeatedFieldResponse"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:   proto.String("next_page_token"),
				Number: proto.Int32(int32(1)),
				Type:   typep(descriptorpb.FieldDescriptorProto_TYPE_STRING),
			},
		},
	}
	noNextPageToken := &descriptorpb.DescriptorProto{
		Name: proto.String("NoNextPageTokenResponse"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:     proto.String("items"),
				Number:   proto.Int32(int32(1)),
				Type:     typep(descriptorpb.FieldDescriptorProto_TYPE_MESSAGE),
				TypeName: proto.String(".paging.RandomMessage"),
				Label:    labelp(descriptorpb.FieldDescriptorProto_LABEL_REPEATED),
			},
		},
	}

	// Methods
	validPageSizeMthd := &descriptorpb.MethodDescriptorProto{
		Name:       proto.String("ValidPageSize"),
		InputType:  proto.String(".paging.ValidPageSizeRequest"),
		OutputType: proto.String(".paging.ValidRepeatedResponse"),
	}
	validPageSizeMultipleMthd := &descriptorpb.MethodDescriptorProto{
		Name:       proto.String("ValidPageSizeMultiple"),
		InputType:  proto.String(".paging.ValidPageSizeRequest"),
		OutputType: proto.String(".paging.MultipleRepeatedResponse"),
	}
	validMaxResultsRepeatedMthd := &descriptorpb.MethodDescriptorProto{
		Name:       proto.String("ValidMaxResults"),
		InputType:  proto.String(".paging.ValidMaxResultsRequest"),
		OutputType: proto.String(".paging.ValidRepeatedResponse"),
	}
	validMaxResultsMthd := &descriptorpb.MethodDescriptorProto{
		Name:       proto.String("ValidMaxResults"),
		InputType:  proto.String(".paging.ValidMaxResultsRequest"),
		OutputType: proto.String(".paging.ValidMapResponse"),
	}
	validPageSizeMapMthd := &descriptorpb.MethodDescriptorProto{
		Name:       proto.String("ValidPageSize"),
		InputType:  proto.String(".paging.ValidPageSizeRequest"),
		OutputType: proto.String(".paging.ValidMapResponse"),
	}
	clientStreamingMthd := &descriptorpb.MethodDescriptorProto{
		Name:            proto.String("ClientStreaming"),
		InputType:       proto.String(".paging.ValidPageSizeRequest"),
		OutputType:      proto.String(".paging.ValidRepeatedResponse"),
		ClientStreaming: proto.Bool(bool(true)),
	}
	serverStreamingMthd := &descriptorpb.MethodDescriptorProto{
		Name:            proto.String("ServerStreaming"),
		InputType:       proto.String(".paging.ValidPageSizeRequest"),
		OutputType:      proto.String(".paging.ValidRepeatedResponse"),
		ServerStreaming: proto.Bool(bool(true)),
	}
	tooManyRepeatedMthd := &descriptorpb.MethodDescriptorProto{
		Name:       proto.String("TooManyRepeated"),
		InputType:  proto.String(".paging.ValidPageSizeRequest"),
		OutputType: proto.String(".paging.TooManyRepeatedResponse"),
	}
	tooManyMapMthd := &descriptorpb.MethodDescriptorProto{
		Name:       proto.String("TooManyMap"),
		InputType:  proto.String(".paging.ValidPageSizeRequest"),
		OutputType: proto.String(".paging.TooManyMapResponse"),
	}
	noNextPageTokenMthd := &descriptorpb.MethodDescriptorProto{
		Name:       proto.String("NoNextPageToken"),
		InputType:  proto.String(".paging.ValidPageSizeRequest"),
		OutputType: proto.String(".paging.NoNextPageTokenResponse"),
	}
	noRepeatedFieldMthd := &descriptorpb.MethodDescriptorProto{
		Name:       proto.String("NoRepeatedField"),
		InputType:  proto.String(".paging.ValidPageSizeRequest"),
		OutputType: proto.String(".paging.NoRepeatedFieldResponse"),
	}

	file := &descriptorpb.FileDescriptorProto{
		Package: proto.String("paging"),
		Options: &descriptorpb.FileOptions{
			GoPackage: proto.String("mypackage"),
		},
		MessageType: []*descriptorpb.DescriptorProto{
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
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: proto.String("TestService"),
				Method: []*descriptorpb.MethodDescriptorProto{
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

	req := pluginpb.CodeGeneratorRequest{
		Parameter: proto.String("go-gapic-package=path;mypackage,transport=rest"),
		ProtoFile: []*descriptorpb.FileDescriptorProto{file},
	}
	g, err := newGenerator(&req)
	if err != nil {
		t.Fatal(err)
	}
	g.apiName = "Awesome API"
	g.imports = map[pbinfo.ImportSpec]bool{}
	g.opts = &options{transports: []transport{rest}}

	for _, tst := range []struct {
		mthd      *descriptorpb.MethodDescriptorProto
		sizeField *descriptorpb.FieldDescriptorProto // A nil field means this is not a paged method
		iterField *descriptorpb.FieldDescriptorProto // A nil field means this is not a paged method
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
