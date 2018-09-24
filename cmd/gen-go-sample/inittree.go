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

package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"text/scanner"
	"unicode"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/gapic-generator-go/internal/errors"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
)

// initType represents a type of a value in initialization tree.
//
// In protobuf, array-ness and map-ness are properties of fields, not types,
// and primitives are treated differently than message types.
// initType treats a type in a way similar to how Go treats them,
// so that it's easy to generate initialization code.
type initType struct {
	// If the type is a message, desc describes the message.
	// If the type is a primitive, primValid reports whether the token is valid for the type.
	// TODO(pongad): desc probably won't work with enums.
	desc      *descriptor.DescriptorProto
	primValid func(string) bool
}

// initTree represents a node in the initialization tree.
type initTree struct {
	typ initType

	// initTree is either a composite value (struct/array/map), where keys and vals are set,
	// or a simple value, where leafVal is set.

	// Use array representation; we need order, and we probably won't
	// have many pairs anyway.
	keys []string
	vals []*initTree

	// Text of the literal. If the literal is a string, it's already quoted.
	leafVal string
}

func (t *initTree) get(k string, info pbinfo.Info) (*initTree, error) {
	for i, key := range t.keys {
		if k == key {
			return t.vals[i], nil
		}
	}

	v := new(initTree)
	for _, f := range t.typ.desc.Field {
		if f.GetName() != k {
			continue
		}

		if tn := f.GetTypeName(); tn != "" {
			typ := info.Type[f.GetTypeName()]
			if typ == nil {
				return nil, errors.E(nil, "cannot find descriptor of %q", f.GetTypeName())
			}
			v.typ.desc = typ
		} else {
			// Since the tokens are given to us by scanner, they must already be a valid token of some type,
			// no need to check exhaustively.
			switch f.GetType() {
			case descriptor.FieldDescriptorProto_TYPE_BOOL:
				v.typ.primValid = func(s string) bool { return s == "true" || s == "false" }
			case descriptor.FieldDescriptorProto_TYPE_BYTES, descriptor.FieldDescriptorProto_TYPE_STRING:
				v.typ.primValid = func(s string) bool { return s[0] == '"' }
			case descriptor.FieldDescriptorProto_TYPE_DOUBLE, descriptor.FieldDescriptorProto_TYPE_FLOAT:
				v.typ.primValid = func(s string) bool { return s == "inf" || s == "nan" || strings.Trim(s, "+-.0123456789") == "" }
			default:
				v.typ.primValid = func(s string) bool { return strings.Trim(s, "0123456789") == "" }
			}
		}
	}
	if v.typ.desc == nil && v.typ.primValid == nil {
		return nil, errors.E(nil, "type %q does not have field %q", t.typ.desc.GetName(), k)
	}

	t.keys = append(t.keys, k)
	t.vals = append(t.vals, v)
	return v, nil
}

// Just do simple structs for now.
// TODO(pongad): allow map and array index.
//
// spec = ident { '.' ident } [ '=' value ] .
func (t *initTree) Parse(txt string, info pbinfo.Info) error {
	var sc scanner.Scanner

	// The first ident is treated specially to be a member of the root.
	// Since we know the root is a struct, a dot is the only legal token anyway,
	// so just insert a dot here so we don't need to treat the first token specially.
	sc.Init(strings.NewReader("." + txt))
	sc.Mode = scanner.ScanIdents | scanner.ScanInts | scanner.ScanFloats | scanner.ScanStrings

	// Scanner reports error by calling sc.Error; we save the error so we always report
	// the first.
	var err error
	report := func(e error) error {
		if err == nil {
			err = e
		}
		return err
	}
	sc.Error = func(_ *scanner.Scanner, msg string) {
		e := errors.E(nil, msg)
		e = errors.E(e, "while scaning: %q", txt)
		report(e)
	}

	for {
		switch sc.Scan() {
		case '=':
			goto equal
		case scanner.EOF:
			// TODO: set t.leafVal to zero value of type.
			return report(nil)

		case '.':
			if sc.Scan() != scanner.Ident {
				return report(errors.E(nil, "expected ident, found %q", sc.TokenText()))
			}
			if t2, err := t.get(sc.TokenText(), info); err != nil {
				return report(err)
			} else {
				t = t2
			}

		default:
			return report(errors.E(nil, "unexpected %q", sc.TokenText()))
		}
	}

	// TODO(pongad): handle resource names
	// TODO(pongad): properly validate and print enums
equal:
	switch r := sc.Scan(); r {
	case scanner.Int, scanner.Float, scanner.String, scanner.Ident:
		if lv := t.leafVal; lv != "" {
			return report(errors.E(nil, "value already set to %q", lv))
		}
		if t.typ.primValid == nil {
			return report(errors.E(nil, "not a primitive field"))
		}

		tok := sc.TokenText()
		if !t.typ.primValid(tok) {
			// TODO(pongad): we should probably tell user what type the field is.
			return report(errors.E(nil, "invalid value for type: %q", tok))
		}
		t.leafVal = tok
	default:
		return report(errors.E(nil, "expected value, found %q", sc.TokenText()))
	}

	if sc.Scan() != scanner.EOF {
		return report(errors.E(nil, "expected EOF, found %q", sc.TokenText()))
	}
	return err
}

func (t *initTree) Print(w io.Writer, g *generator) error {
	bw := bufio.NewWriter(w)
	err := t.print(bw, g, 1)
	if err2 := bw.Flush(); err == nil {
		err = err2
	}
	return err
}

func (t *initTree) print(w *bufio.Writer, g *generator, ind int) error {
	if v := t.leafVal; v != "" {
		w.WriteString(v)
		return nil
	}

	desc := t.typ.desc
	impSpec, err := g.descInfo.ImportSpec(desc)
	if err != nil {
		return err
	}
	g.imports[impSpec] = true

	fmt.Fprintf(w, "%s.%s{\n", impSpec.Name, desc.GetName())
	for i, k := range t.keys {
		for i := 0; i < ind; i++ {
			w.WriteByte('\t')
		}
		snakeToCapCamel(w, k)
		w.WriteString(": ")
		if err := t.vals[i].print(w, g, ind+1); err != nil {
			return err
		}
		w.WriteString(",\n")
	}
	w.WriteString("}")
	return nil
}

func snakeToCapCamel(w *bufio.Writer, s string) {
	cap := true
	for _, r := range s {
		if r == '_' {
			cap = true
		} else if cap {
			cap = false
			w.WriteRune(unicode.ToUpper(r))
		} else {
			w.WriteRune(r)
		}
	}
}
