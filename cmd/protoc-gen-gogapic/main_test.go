package main

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

func TestComment(t *testing.T) {
	var g generator

	for _, tst := range []struct {
		in, want string
	}{
		{
			in:   "",
			want: "",
		},
		{
			in:   "abc\ndef\n",
			want: "// abc\n// def\n",
		},
	} {
		g.sb.Reset()
		g.comment(tst.in)
		if got := g.sb.String(); got != tst.want {
			t.Errorf("comment(%q) = %q, want %q", tst.in, got, tst.want)
		}
	}
}

func TestMethodDoc(t *testing.T) {
	m := &descriptor.MethodDescriptorProto{
		Name: proto.String("MyMethod"),
	}

	var g generator
	g.comments = make(map[proto.Message]string)

	for _, tst := range []struct {
		in, want string
	}{
		{
			in:   "",
			want: "",
		},
		{
			in:   "Does stuff.\n It also does other stuffs.",
			want: "// MyMethod does stuff.\n// It also does other stuffs.\n",
		},
	} {
		g.comments[m] = tst.in
		g.sb.Reset()
		g.methodDoc(m)
		if got := g.sb.String(); got != tst.want {
			t.Errorf("comment(%q) = %q, want %q", tst.in, got, tst.want)
		}
	}
}

func TestReduceServName(t *testing.T) {
	for _, tst := range []struct {
		in, want string
	}{
		{"Foo", "Foo"},
		{"FooV2", "Foo"},
		{"FooService", "Foo"},
		{"FooServiceV2", "Foo"},
		{"FooV2Bar", "FooV2Bar"},
	} {
		if got := reduceServName(tst.in); got != tst.want {
			t.Errorf("reduceServName(%q) = %q, want %q", tst.in, got, tst.want)
		}
	}
}
