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

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func TestNameSpec(t *testing.T) {
	t.Parallel()

	subMsg := &descriptorpb.DescriptorProto{
		Name: proto.String("SubMessage"),
	}
	msg := &descriptorpb.DescriptorProto{
		Name:       proto.String("Message"),
		NestedType: []*descriptorpb.DescriptorProto{subMsg},
	}
	anotherMsg := &descriptorpb.DescriptorProto{
		Name: proto.String("AnotherMessage"),
	}
	file := &descriptorpb.FileDescriptorProto{
		Options: &descriptorpb.FileOptions{
			GoPackage: proto.String("path.to/pb/foo;foo"),
		},
		MessageType: []*descriptorpb.DescriptorProto{msg},
	}
	anotherFile := &descriptorpb.FileDescriptorProto{
		Name: proto.String("bar.proto"),
		Options: &descriptorpb.FileOptions{
			GoPackage: proto.String("path.to/pb/bar;bar"),
		},
		MessageType: []*descriptorpb.DescriptorProto{anotherMsg},
	}

	info := Of([]*descriptorpb.FileDescriptorProto{file, anotherFile})
	info.PkgOverrides = map[string]string{
		anotherFile.GetName(): "path.to/pb/foo;foo",
	}

	for _, tst := range []struct {
		e    ProtoType
		name string
	}{
		{msg, "Message"},
		{subMsg, "Message_SubMessage"},
		{anotherMsg, "AnotherMessage"},
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
