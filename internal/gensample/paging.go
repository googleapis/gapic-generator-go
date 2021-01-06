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
	"fmt"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/gapic-generator-go/internal/errors"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
)

// pagingField reports the "resource field" to be iterated over by paginating method m.
// This is temporarily copied from gengapic/paging.go. We probably need to read this from annoations later.
func pagingField(info pbinfo.Info, m *descriptor.MethodDescriptorProto) (*descriptor.FieldDescriptorProto, error) {
	var (
		hasSize, hasToken, hasNextToken bool
		elemFields                      []*descriptor.FieldDescriptorProto
	)

	inType := info.Type[m.GetInputType()]
	if inType == nil {
		return nil, errors.E(nil, "cannot find message type %q, malformed descriptor?", m.GetInputType())
	}
	inMsg, ok := inType.(*descriptor.DescriptorProto)
	if !ok {
		return nil, errors.E(nil, "expected %q to be message type, found %T", m.GetInputType(), inType)
	}

	outType := info.Type[m.GetOutputType()]
	if outType == nil {
		return nil, errors.E(nil, "cannot find message type %q, malformed descriptor?", m.GetOutputType())
	}
	outMsg, ok := outType.(*descriptor.DescriptorProto)
	if !ok {
		return nil, errors.E(nil, "expected %q to be message type, found %T", m.GetOutputType(), outType)
	}

	for _, f := range inMsg.Field {
		if f.GetName() == "page_size" && f.GetType() == descriptor.FieldDescriptorProto_TYPE_INT32 {
			hasSize = true
		}
		if f.GetName() == "page_token" && f.GetType() == descriptor.FieldDescriptorProto_TYPE_STRING {
			hasToken = true
		}
	}
	for _, f := range outMsg.Field {
		if f.GetName() == "next_page_token" && f.GetType() == descriptor.FieldDescriptorProto_TYPE_STRING {
			hasNextToken = true
		}
		if f.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			elemFields = append(elemFields, f)
		}
	}
	if !hasSize || !hasToken || !hasNextToken {
		return nil, nil
	}
	// TODO(noahdietz) relax this requirement and treat the method as non-paging.
	// See https://github.com/googleapis/gapic-generator-go/issues/493.
	if len(elemFields) == 0 {
		return nil, fmt.Errorf("%s looks like paging method, but can't find repeated field in %s", *m.Name, outType.GetName())
	}
	if len(elemFields) > 1 {
		return nil, fmt.Errorf("%s looks like paging method, but too many repeated fields in %s", *m.Name, outType.GetName())
	}
	return elemFields[0], nil
}
