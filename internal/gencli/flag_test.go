package gencli

import (
	"testing"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

func TestInputFieldName(t *testing.T) {
	for _, tst := range []struct {
		f    *Flag
		want string
	}{
		{
			f:    &Flag{Name: "name"},
			want: "Name",
		},
		{
			f:    &Flag{Name: "kiosk_id"},
			want: "KioskId",
		},
	} {
		if got := tst.f.InputFieldName(); got != tst.want {
			t.Errorf("(%s).InputFieldName() = %q, want %q", tst.f.Name, got, tst.want)
		}
	}
}

func TestGenOtherVarName(t *testing.T) {
	for _, tst := range []struct {
		f        *Flag
		in, want string
	}{
		{
			f:    &Flag{Name: "message.enum"},
			in:   "ClientInput",
			want: "ClientInputMessageEnum",
		},
	} {
		if got := tst.f.GenOtherVarName(tst.in); got != tst.want {
			t.Errorf("(%s).GenOtherVarName(%s) = %q, want %q", tst.f.Name, tst.in, got, tst.want)
		}
	}
}

func TestGenOneOfVarName(t *testing.T) {
	for _, tst := range []struct {
		f        *Flag
		in, want string
	}{
		{
			f:    &Flag{Name: "oneof.field"},
			in:   "ClientInput",
			want: "ClientInputOneofField",
		},
		{
			f:    &Flag{Name: "oneof.field_snake"},
			in:   "ClientInput",
			want: "ClientInputOneofFieldSnake",
		},
		{
			f:    &Flag{Name: "oneof.msg.field"},
			in:   "ClientInput",
			want: "ClientInputOneofMsg",
		},
		{
			f:    &Flag{Name: "oneof.msg.field", IsNested: true},
			in:   "ClientInput",
			want: "ClientInputOneofMsgField",
		},
	} {
		if got := tst.f.GenOneOfVarName(tst.in); got != tst.want {
			t.Errorf("(%s, %v).GenOneOfVarName(%s) = %q, want %q", tst.f.Name, tst.f.IsNested, tst.in, got, tst.want)
		}
	}
}

func TestOneOfInputFieldName(t *testing.T) {
	for _, tst := range []struct {
		f    *Flag
		want string
	}{
		{
			f:    &Flag{Name: "oneof.field"},
			want: "Field",
		},
		{
			f:    &Flag{Name: "oneof.field_snake"},
			want: "FieldSnake",
		},
		{
			f:    &Flag{Name: "oneof.msg.field"},
			want: "Msg.Field",
		},
		{
			f:    &Flag{Name: "oneof.nested.msg.field", IsNested: true},
			want: "Field",
		},
	} {
		if got := tst.f.OneOfInputFieldName(); got != tst.want {
			t.Errorf("(%s, %v).OneOfInputFieldName() = %q, want %q", tst.f.Name, tst.f.IsNested, got, tst.want)
		}
	}
}

func TestGenFlag(t *testing.T) {
	for _, tst := range []struct {
		f        *Flag
		in, want string
	}{
		{
			f: &Flag{
				Name:  "field",
				Type:  descriptor.FieldDescriptorProto_TYPE_STRING,
				Usage: "this is the usage",
			},
			in:   "ClientInput",
			want: `StringVar(&ClientInput.Field, "field", "", "this is the usage")`,
		},
		{
			f: &Flag{
				Name:  "field",
				Type:  descriptor.FieldDescriptorProto_TYPE_BOOL,
				Usage: "this is the usage",
			},
			in:   "ClientInput",
			want: `BoolVar(&ClientInput.Field, "field", false, "this is the usage")`,
		},
		{
			f: &Flag{
				Name:  "field",
				Type:  descriptor.FieldDescriptorProto_TYPE_INT32,
				Usage: "this is the usage",
			},
			in:   "ClientInput",
			want: `Int32Var(&ClientInput.Field, "field", 0, "this is the usage")`,
		},
		{
			f: &Flag{
				Name:  "field",
				Type:  descriptor.FieldDescriptorProto_TYPE_FLOAT,
				Usage: "this is the usage",
			},
			in:   "ClientInput",
			want: `Float32Var(&ClientInput.Field, "field", 0.0, "this is the usage")`,
		},
		{
			f: &Flag{
				Name:  "field",
				Type:  descriptor.FieldDescriptorProto_TYPE_DOUBLE,
				Usage: "this is the usage",
			},
			in:   "ClientInput",
			want: `Float64Var(&ClientInput.Field, "field", 0.0, "this is the usage")`,
		},
		{
			f: &Flag{
				Name:  "field",
				Type:  descriptor.FieldDescriptorProto_TYPE_BYTES,
				Usage: "this is the usage",
			},
			in:   "ClientInput",
			want: `BytesHexVar(&ClientInput.Field, "field", []byte{}, "this is the usage")`,
		},
		{
			f: &Flag{
				Name:     "field",
				Type:     descriptor.FieldDescriptorProto_TYPE_STRING,
				Usage:    "this is the usage",
				Repeated: true,
			},
			in:   "ClientInput",
			want: `StringSliceVar(&ClientInput.Field, "field", []string{}, "this is the usage")`,
		},
		{
			f: &Flag{
				Name:     "field",
				Type:     descriptor.FieldDescriptorProto_TYPE_MESSAGE,
				Usage:    "this is the usage",
				Repeated: true,
			},
			in:   "ClientInput",
			want: `StringArrayVar(&ClientInputField, "field", []string{}, "this is the usage")`,
		},
		{
			f: &Flag{
				Name:  "field",
				Type:  descriptor.FieldDescriptorProto_TYPE_ENUM,
				Usage: "this is the usage",
			},
			in:   "ClientInput",
			want: `StringVar(&ClientInputField, "field", "", "this is the usage")`,
		},
		{
			f: &Flag{
				Name:         "oneof.field",
				Type:         descriptor.FieldDescriptorProto_TYPE_STRING,
				Usage:        "this is the usage",
				IsOneOfField: true,
			},
			in:   "ClientInput",
			want: `StringVar(&ClientInputOneofField.Field, "oneof.field", "", "this is the usage")`,
		},
		{
			f: &Flag{
				Name:   "oneof_selector",
				Type:   descriptor.FieldDescriptorProto_TYPE_STRING,
				Usage:  "this is the usage",
				OneOfs: map[string]*Flag{"test": &Flag{}},
			},
			in:   "ClientInput",
			want: `StringVar(&ClientInputOneofSelector, "oneof_selector", "", "this is the usage")`,
		},
	} {
		if got := tst.f.GenFlag(tst.in); got != tst.want {
			t.Errorf("(%+v).GenFlag(%s) = %q, want %q", tst.f, tst.in, got, tst.want)
		}
	}
}

func TestGenRequired(t *testing.T) {
	for _, tst := range []struct {
		f    *Flag
		want string
	}{
		{
			f:    &Flag{Name: "name"},
			want: `cmd.MarkFlagRequired("name")`,
		},
		{
			f:    &Flag{Name: "kiosk_id"},
			want: `cmd.MarkFlagRequired("kiosk_id")`,
		},
	} {
		if got := tst.f.GenRequired(); got != tst.want {
			t.Errorf("(%s).GenRequired() = %q, want %q", tst.f.Name, got, tst.want)
		}
	}
}

func TestIsMessage(t *testing.T) {
	for _, tst := range []struct {
		f    *Flag
		want bool
	}{
		{
			f:    &Flag{Type: descriptor.FieldDescriptorProto_TYPE_MESSAGE},
			want: true,
		},
		{
			f:    &Flag{Type: descriptor.FieldDescriptorProto_TYPE_STRING},
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
			f:    &Flag{Type: descriptor.FieldDescriptorProto_TYPE_ENUM},
			want: true,
		},
		{
			f:    &Flag{Type: descriptor.FieldDescriptorProto_TYPE_STRING},
			want: false,
		},
	} {
		if got := tst.f.IsEnum(); got != tst.want {
			t.Errorf("(%v).IsEnum() = %v, want %v", tst.f.Type, got, tst.want)
		}
	}
}
