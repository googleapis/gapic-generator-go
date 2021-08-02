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

package gensample

import (
	"testing"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"google.golang.org/protobuf/proto"
)

func TestEnumFmt(t *testing.T) {
	subMsg := &descriptor.DescriptorProto{
		Name: proto.String("SubMessage"),
	}
	msg := &descriptor.DescriptorProto{
		Name:       proto.String("Message"),
		NestedType: []*descriptor.DescriptorProto{subMsg},
	}
	file := &descriptor.FileDescriptorProto{
		Options: &descriptor.FileOptions{
			GoPackage: proto.String("path.to/pb/foo;foo"),
		},
		MessageType: []*descriptor.DescriptorProto{msg},
	}

	enumDesc := &descriptor.EnumDescriptorProto{
		Name: proto.String("EnumType"),
	}

	for _, tst := range []struct {
		enumPlace *[]*descriptor.EnumDescriptorProto
		want      string
	}{
		{&file.EnumType, "foopb.EnumType_VAL"},
		{&msg.EnumType, "foopb.Message_VAL"},
		{&subMsg.EnumType, "foopb.Message_SubMessage_VAL"},
	} {
		*tst.enumPlace = []*descriptor.EnumDescriptorProto{enumDesc}

		gen := &generator{
			imports:  map[pbinfo.ImportSpec]bool{},
			descInfo: pbinfo.Of([]*descriptor.FileDescriptorProto{file}),
		}
		ef := enumFmt(gen.descInfo, enumDesc)
		got, err := ef(gen, "VAL")
		if err != nil {
			t.Error(err)
		} else if got != tst.want {
			t.Errorf("enumFmt = %q, want %q", got, tst.want)
		}

		*tst.enumPlace = nil
	}
}
