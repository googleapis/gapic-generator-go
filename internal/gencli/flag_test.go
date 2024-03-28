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

package gencli

import (
	"testing"

	"google.golang.org/protobuf/types/descriptorpb"
)

func TestGenFlag(t *testing.T) {
	for _, tst := range []struct {
		f    *Flag
		want string
	}{
		{
			f: &Flag{
				Name:      "field",
				FieldName: "Field",
				VarName:   "ClientInput",
				Type:      descriptorpb.FieldDescriptorProto_TYPE_STRING,
				Usage:     "this is the usage",
			},
			want: `StringVar(&ClientInput.Field, "field", "", "this is the usage")`,
		},
		{
			f: &Flag{
				Name:      "field",
				FieldName: "Field",
				VarName:   "ClientInput",
				Type:      descriptorpb.FieldDescriptorProto_TYPE_BOOL,
				Usage:     "this is the usage",
			},
			want: `BoolVar(&ClientInput.Field, "field", false, "this is the usage")`,
		},
		{
			f: &Flag{
				Name:      "field",
				FieldName: "Field",
				VarName:   "ClientInput",
				Type:      descriptorpb.FieldDescriptorProto_TYPE_INT32,
				Usage:     "this is the usage",
			},
			want: `Int32Var(&ClientInput.Field, "field", 0, "this is the usage")`,
		},
		{
			f: &Flag{
				Name:      "field",
				FieldName: "Field",
				VarName:   "ClientInput",
				Type:      descriptorpb.FieldDescriptorProto_TYPE_FLOAT,
				Usage:     "this is the usage",
			},
			want: `Float32Var(&ClientInput.Field, "field", 0.0, "this is the usage")`,
		},
		{
			f: &Flag{
				Name:      "field",
				FieldName: "Field",
				VarName:   "ClientInput",
				Type:      descriptorpb.FieldDescriptorProto_TYPE_DOUBLE,
				Usage:     "this is the usage",
			},
			want: `Float64Var(&ClientInput.Field, "field", 0.0, "this is the usage")`,
		},
		{
			f: &Flag{
				Name:      "field",
				FieldName: "Field",
				VarName:   "ClientInput",
				Type:      descriptorpb.FieldDescriptorProto_TYPE_BYTES,
				Usage:     "this is the usage",
			},
			want: `BytesHexVar(&ClientInput.Field, "field", []byte{}, "this is the usage")`,
		},
		{
			f: &Flag{
				Name:      "field",
				FieldName: "Field",
				VarName:   "ClientInput",
				Type:      descriptorpb.FieldDescriptorProto_TYPE_STRING,
				Usage:     "this is the usage",
				Repeated:  true,
			},
			want: `StringSliceVar(&ClientInput.Field, "field", []string{}, "this is the usage")`,
		},
		{
			f: &Flag{
				Name:      "field",
				FieldName: "Field",
				Type:      descriptorpb.FieldDescriptorProto_TYPE_MESSAGE,
				Usage:     "this is the usage",
				Repeated:  true,
				VarName:   "ClientInputField",
			},
			want: `StringArrayVar(&ClientInputField, "field", []string{}, "this is the usage")`,
		},
		{
			f: &Flag{
				Name:      "field",
				FieldName: "Field",
				Type:      descriptorpb.FieldDescriptorProto_TYPE_ENUM,
				Usage:     "this is the usage",
				VarName:   "ClientInputField",
			},
			want: `StringVar(&ClientInputField, "field", "", "this is the usage")`,
		},
		{
			f: &Flag{
				Name:         "oneof.field",
				FieldName:    "Field",
				Type:         descriptorpb.FieldDescriptorProto_TYPE_STRING,
				Usage:        "this is the usage",
				VarName:      "ClientInputOneofField",
				IsOneOfField: true,
			},
			want: `StringVar(&ClientInputOneofField.Field, "oneof.field", "", "this is the usage")`,
		},
		{
			f: &Flag{
				Name:    "oneof_selector",
				VarName: "ClientInputOneofSelector",
				Type:    descriptorpb.FieldDescriptorProto_TYPE_STRING,
				Usage:   "this is the usage",
				OneOfs:  map[string]*Flag{"test": &Flag{}},
			},
			want: `StringVar(&ClientInputOneofSelector, "oneof_selector", "", "this is the usage")`,
		},
	} {
		if got := tst.f.GenFlag(); got != tst.want {
			t.Errorf("(%+v).GenFlag() = %q, want %q", tst.f, got, tst.want)
		}
	}
}

func TestIsMessage(t *testing.T) {
	for _, tst := range []struct {
		f    *Flag
		want bool
	}{
		{
			f:    &Flag{Type: descriptorpb.FieldDescriptorProto_TYPE_MESSAGE},
			want: true,
		},
		{
			f:    &Flag{Type: descriptorpb.FieldDescriptorProto_TYPE_STRING},
			want: false,
		},
	} {
		if got := tst.f.IsMessage(); got != tst.want {
			t.Errorf("(%v).IsMessage() = %v, want %v", tst.f.Type, got, tst.want)
		}
	}
}

func TestIsEnum(t *testing.T) {
	for _, tst := range []struct {
		f    *Flag
		want bool
	}{
		{
			f:    &Flag{Type: descriptorpb.FieldDescriptorProto_TYPE_ENUM},
			want: true,
		},
		{
			f:    &Flag{Type: descriptorpb.FieldDescriptorProto_TYPE_STRING},
			want: false,
		},
	} {
		if got := tst.f.IsEnum(); got != tst.want {
			t.Errorf("(%v).IsEnum() = %v, want %v", tst.f.Type, got, tst.want)
		}
	}
}
