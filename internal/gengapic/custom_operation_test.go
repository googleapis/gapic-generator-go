// Copyright 2021 Google LLC
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
	"fmt"
	"path/filepath"
	"testing"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/google/go-cmp/cmp"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"github.com/googleapis/gapic-generator-go/internal/txtdiff"
	"google.golang.org/genproto/googleapis/cloud/extendedops"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoiface"
)

func TestCustomOpProtoName(t *testing.T) {
	pkg := "google.cloud.foo.v1"
	op := &descriptor.DescriptorProto{
		Name: proto.String("Operation"),
	}
	g := &generator{
		aux: &auxTypes{
			customOp: &customOp{
				message: op,
			},
		},
		descInfo: pbinfo.Info{
			ParentFile: map[protoiface.MessageV1]*descriptor.FileDescriptorProto{
				op: {
					Package: proto.String(pkg),
				},
			},
		},
	}
	got := g.customOpProtoName()
	want := fmt.Sprintf(".%s.%s", pkg, op.GetName())
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("got(-),want(+):\n%s", diff)
	}
}

func TestCustomPointerTyp(t *testing.T) {
	op := &descriptor.DescriptorProto{
		Name: proto.String("Operation"),
	}
	g := &generator{
		aux: &auxTypes{
			customOp: &customOp{
				message: op,
			},
		},
		descInfo: pbinfo.Info{
			ParentFile: map[protoiface.MessageV1]*descriptor.FileDescriptorProto{
				op: {
					Package: proto.String("google.cloud.foo.v1"),
					Options: &descriptor.FileOptions{
						GoPackage: proto.String("google.golang.org/genproto/cloud/foo/v1;foo"),
					},
				},
			},
		},
	}
	got, err := g.customOpPointerType()
	if err != nil {
		t.Fatal(err)
	}
	want := "*foopb.Operation"
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("got(-),want(+):\n%s", diff)
	}
}

func TestCustomOpInit(t *testing.T) {
	op := &descriptor.DescriptorProto{
		Name: proto.String("Operation"),
	}
	projFieldOpts := &descriptor.FieldOptions{}
	proto.SetExtension(projFieldOpts, extendedops.E_OperationRequestField, "project")
	projField := &descriptor.FieldDescriptorProto{
		Name:    proto.String("request_project"),
		Options: projFieldOpts,
	}
	zoneFieldOpts := &descriptor.FieldOptions{}
	proto.SetExtension(zoneFieldOpts, extendedops.E_OperationRequestField, "zone")
	zoneField := &descriptor.FieldDescriptorProto{
		Name:    proto.String("request_zone"),
		Options: zoneFieldOpts,
	}
	req := &descriptor.DescriptorProto{
		Field: []*descriptor.FieldDescriptorProto{projField, zoneField},
	}
	opServ := &descriptor.ServiceDescriptorProto{Name: proto.String("FooOperationService")}
	g := &generator{
		aux: &auxTypes{
			customOp: &customOp{
				message: op,
				pollingParams: map[*descriptor.ServiceDescriptorProto][]string{
					opServ: {"project", "zone"},
				},
			},
		},
		opts: &options{pkgName: "bar"},
	}
	g.customOpInit("foo", "req", "op", req, opServ)
	txtdiff.Diff(t, g.pt.String(), filepath.Join("testdata", "custom_op_init_helper.want"))
}

func TestCustomOperationType(t *testing.T) {
	errorType := &descriptor.DescriptorProto{
		Name: proto.String("Error"),
		Field: []*descriptor.FieldDescriptorProto{
			{
				Name:  proto.String("nested"),
				Type:  typep(descriptor.FieldDescriptorProto_TYPE_STRING),
				Label: labelp(descriptor.FieldDescriptorProto_LABEL_OPTIONAL),
			},
		},
	}
	errorField := &descriptor.FieldDescriptorProto{
		Name:     proto.String("error"),
		Type:     typep(descriptor.FieldDescriptorProto_TYPE_MESSAGE),
		TypeName: proto.String("Error"),
	}

	errorCodeOpts := &descriptor.FieldOptions{}
	proto.SetExtension(errorCodeOpts, extendedops.E_OperationField, extendedops.OperationResponseMapping_ERROR_CODE)
	errorCodeField := &descriptor.FieldDescriptorProto{
		Name:    proto.String("http_error_status_code"),
		Type:    descriptor.FieldDescriptorProto_TYPE_STRING.Enum(),
		Options: errorCodeOpts,
	}

	errorMessageOpts := &descriptor.FieldOptions{}
	proto.SetExtension(errorMessageOpts, extendedops.E_OperationField, extendedops.OperationResponseMapping_ERROR_MESSAGE)
	errorMessageField := &descriptor.FieldDescriptorProto{
		Name:    proto.String("http_error_message"),
		Type:    descriptor.FieldDescriptorProto_TYPE_STRING.Enum(),
		Options: errorMessageOpts,
	}

	nameOpts := &descriptor.FieldOptions{}
	proto.SetExtension(nameOpts, extendedops.E_OperationField, extendedops.OperationResponseMapping_NAME)
	nameField := &descriptor.FieldDescriptorProto{
		Name:    proto.String("name"),
		Type:    descriptor.FieldDescriptorProto_TYPE_STRING.Enum(),
		Options: nameOpts,
	}

	statusEnum := &descriptor.EnumDescriptorProto{
		Name: proto.String("Status"),
		Value: []*descriptor.EnumValueDescriptorProto{
			{
				Name:   proto.String("DONE"),
				Number: proto.Int32(0),
			},
		},
	}

	statusOpts := &descriptor.FieldOptions{}
	proto.SetExtension(statusOpts, extendedops.E_OperationField, extendedops.OperationResponseMapping_STATUS)
	statusEnumField := &descriptor.FieldDescriptorProto{
		Name:     proto.String("status"),
		Type:     descriptor.FieldDescriptorProto_TYPE_ENUM.Enum(),
		TypeName: proto.String(".google.cloud.foo.v1.Operation.Status"),
		Options:  statusOpts,
	}

	statusBoolField := &descriptor.FieldDescriptorProto{
		Name:    proto.String("status"),
		Type:    descriptor.FieldDescriptorProto_TYPE_BOOL.Enum(),
		Options: statusOpts,
	}

	op := &descriptor.DescriptorProto{
		Name:     proto.String("Operation"),
		EnumType: []*descriptor.EnumDescriptorProto{statusEnum},
	}

	inNameOpts := &descriptor.FieldOptions{}
	proto.SetExtension(inNameOpts, extendedops.E_OperationResponseField, nameField.GetName())
	inNameField := &descriptor.FieldDescriptorProto{
		Name:    proto.String("operation"),
		Type:    descriptor.FieldDescriptorProto_TYPE_STRING.Enum(),
		Options: inNameOpts,
	}

	getOpts := &descriptor.MethodOptions{}
	proto.SetExtension(getOpts, extendedops.E_OperationPollingMethod, true)
	getInput := &descriptor.DescriptorProto{
		Name:  proto.String("GetFooOperationRequest"),
		Field: []*descriptor.FieldDescriptorProto{inNameField},
	}

	fooOpServ := &descriptor.ServiceDescriptorProto{
		Name: proto.String("FooOperationsService"),
		Method: []*descriptor.MethodDescriptorProto{
			{
				Name:      proto.String("Get"),
				InputType: proto.String(".google.cloud.foo.v1.GetFooOperationRequest"),
				Options:   getOpts,
			},
		},
	}

	f := &descriptor.FileDescriptorProto{
		Package: proto.String("google.cloud.foo.v1"),
		Options: &descriptor.FileOptions{
			GoPackage: proto.String("google.golang.org/genproto/cloud/foo/v1;foo"),
		},
		Service: []*descriptor.ServiceDescriptorProto{fooOpServ},
	}
	g := &generator{
		aux: &auxTypes{
			customOp: &customOp{
				message: op,
				handles: []*descriptor.ServiceDescriptorProto{fooOpServ},
				pollingParams: map[*descriptor.ServiceDescriptorProto][]string{
					fooOpServ: {"project", "zone"},
				},
			},
		},
		descInfo: pbinfo.Info{
			ParentFile: map[protoiface.MessageV1]*descriptor.FileDescriptorProto{
				op:        f,
				fooOpServ: f,
				getInput:  f,
			},
			ParentElement: map[pbinfo.ProtoType]pbinfo.ProtoType{
				statusEnum: op,
			},
			Type: map[string]pbinfo.ProtoType{
				statusEnumField.GetTypeName():                 statusEnum,
				".google.cloud.foo.v1.Error":                  errorType,
				".google.cloud.foo.v1.GetFooOperationRequest": getInput,
			},
		},
		imports: map[pbinfo.ImportSpec]bool{},
	}
	for _, tst := range []struct {
		name       string
		st         *descriptor.FieldDescriptorProto
		errorField bool
	}{
		{
			name: "enum",
			st:   statusEnumField,
		},
		{
			name:       "bool",
			st:         statusBoolField,
			errorField: true,
		},
	} {
		t.Run(tst.name, func(t *testing.T) {
			op.Field = []*descriptor.FieldDescriptorProto{errorCodeField, errorMessageField, nameField, tst.st}
			if tst.errorField {
				op.Field = append(op.Field, errorField)
			}
			err := g.customOperationType()
			if err != nil {
				t.Fatal(err)
			}
			tn := "custom_op_type_" + tst.name
			txtdiff.Diff(t, g.pt.String(), filepath.Join("testdata", tn+".want"))
			g.reset()
		})
	}
}
