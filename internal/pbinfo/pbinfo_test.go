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

package pbinfo

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

func TestNameSpec(t *testing.T) {
	t.Parallel()

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

	info := Of([]*descriptor.FileDescriptorProto{file})

	for _, tst := range []struct {
		e    ProtoType
		name string
	}{
		{msg, "Message"},
		{subMsg, "Message_SubMessage"},
	} {
		name, imp, err := info.NameSpec(tst.e)
		if err != nil {
			t.Error(err)
			continue
		}
		if name != tst.name {
			t.Errorf("NameSpec(%v).name = %q, want %q", tst.e, name, tst.name)
		}
		wantImp := ImportSpec{Path: "path.to/pb/foo", Name: "foopb"}
		if imp != wantImp {
			t.Errorf("NameSpec(%v).imp = %v, want %v", tst.e, imp, wantImp)
		}
	}
}
