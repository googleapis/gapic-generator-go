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
	"strconv"
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
	// If the type is a protobuf "compound type", desc describes the type.
	// If the type is a primitive, prim tells us what the type is.
	desc pbinfo.ProtoType
	prim descriptor.FieldDescriptorProto_Type

	repeated bool

	// valFmt, if not nil, post-processes values to be included into init struct.
	// NOTE(pongad): This func signature might seem too general. I think it is just general enough
	// to deal with enums and bytes. Time will tell.
	valFmt func(*generator, string) (string, error)
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

	var fields []*descriptor.FieldDescriptorProto
	if msg, ok := t.typ.desc.(*descriptor.DescriptorProto); ok {
		fields = msg.Field
	} else {
		return nil, errors.E(nil, "type does not have fields: %T", t.typ.desc)
	}

	v := new(initTree)
	for _, f := range fields {
		if f.GetName() != k {
			continue
		}

		if tn := f.GetTypeName(); tn == "" {
			// We're a primitive type.
			v.typ.prim = f.GetType()
		} else if typ := info.Type[tn]; typ == nil {
			return nil, errors.E(nil, "cannot find descriptor of %q", f.GetTypeName())
		} else {
			// type is a message
			v.typ.desc = typ
		}
		v.typ.repeated = f.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED

		break
	}
	if v.typ.desc == nil && v.typ.prim == 0 {
		return nil, errors.E(nil, "type %q does not have field %q", t.typ.desc.GetName(), k)
	}

	t.keys = append(t.keys, k)
	t.vals = append(t.vals, v)
	return v, nil
}

func (t *initTree) index(k string) *initTree {
	for i, key := range t.keys {
		if k == key {
			return t.vals[i]
		}
	}

	v := new(initTree)
	v.typ = t.typ
	v.typ.repeated = false

	t.keys = append(t.keys, k)
	t.vals = append(t.vals, v)
	return v
}

// Just do simple structs for now.
// TODO(pongad): allow map and array index.
//
// spec = path [ '=' value ] .
func (t *initTree) parseInit(txt string, info pbinfo.Info) error {
	// The first ident is treated specially to be a member of the root.
	// Since we know the root is a struct, a dot is the only legal token anyway,
	// so just insert a dot here so we don't need to treat the first token specially.
	sc, report := initScanner("." + txt)

	t, r, err := t.parsePathRest(sc, info)
	if err != nil {
		return report(err)
	}
	if r != '=' {
		return report(errors.E(nil, "expected '=', got %q", r))
	}

	// TODO(pongad): handle resource names
	switch r := sc.Scan(); r {
	case scanner.Int, scanner.Float, scanner.String, scanner.Ident:
		if lv := t.leafVal; lv != "" {
			return report(errors.E(nil, "value already set to %q", lv))
		}

		tok := sc.TokenText()

		if enum, ok := t.typ.desc.(*descriptor.EnumDescriptorProto); ok {
			validEnum, eFmt := describeEnum(info, enum)
			if !validEnum(tok) {
				return report(errors.E(nil, "invalid value for type %q: %q", enum.GetName(), tok))
			}
			t.typ.valFmt = eFmt
		} else {
			pType := t.typ.prim
			validPrim := validPrims[pType]
			if validPrim == nil {
				return report(errors.E(nil, "not a primitive type? %q", pType))
			}
			if !validPrim(tok) {
				return report(errors.E(nil, "invalid value for type %q: %q", pType, tok))
			}
		}
		t.leafVal = tok

	default:
		return report(errors.E(nil, "expected value, found %q", sc.TokenText()))
	}

	if sc.Scan() != scanner.EOF {
		return report(errors.E(nil, "expected EOF, found %q", sc.TokenText()))
	}
	return report(nil)
}

func (t *initTree) parseSampleArgPath(txt string, info pbinfo.Info, varName string) (*initTree, error) {
	sc, report := initScanner("." + txt)
	t, _, err := t.parsePathRest(sc, info)
	if err != nil {
		return nil, report(err)
	}

	var cp initTree
	cp = *t
	*t = initTree{leafVal: varName}

	return &cp, report(nil)
}

// parsePathRest parses pathRest, an "unrooted" path, as defined below.
// It returns the subtree specified by the path, the last scanned token, and any error.
//
// path = ident pathRest .
// pathRest = { '.' ident } .
func (t *initTree) parsePathRest(sc *scanner.Scanner, info pbinfo.Info) (*initTree, rune, error) {
	for {
		switch r := sc.Scan(); r {
		case '.':
			if t.typ.repeated {
				return nil, 0, errors.E(nil, "cannot access member of repeated field")
			}
			if r := sc.Scan(); r != scanner.Ident {
				return nil, r, errors.E(nil, "expected ident, found %q", sc.TokenText())
			}
			if t2, err := t.get(sc.TokenText(), info); err != nil {
				return nil, 0, err
			} else {
				t = t2
			}

		case '[':
			if !t.typ.repeated {
				return nil, 0, errors.E(nil, "cannot index into singular field")
			}
			if r := sc.Scan(); r != scanner.Int {
				return nil, r, errors.E(nil, "expected int, found %q", sc.TokenText())
			}
			indVal := sc.TokenText()
			if r := sc.Scan(); r != ']' {
				return nil, r, errors.E(nil, "expected ']', found %q", sc.TokenText())
			}

			t = t.index(indVal)

		default:
			return t, r, nil
		}
	}
}

func initScanner(s string) (*scanner.Scanner, func(error) error) {
	var sc scanner.Scanner
	sc.Init(strings.NewReader(s))
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
		e = errors.E(e, "while scaning: %q", s)
		report(e)
	}

	sc.IsIdentRune = func(r rune, i int) bool {
		return unicode.IsLetter(r) || r == '_' || i == 0 && r == '$' || i > 0 && unicode.IsDigit(r)
	}

	return &sc, report
}

func (t *initTree) Print(w io.Writer, g *generator) error {
	bw := bufio.NewWriter(w)
	err := t.print(bw, g, 0)
	if err2 := bw.Flush(); err == nil {
		err = err2
	}
	return err
}

func (t *initTree) print(w *bufio.Writer, g *generator, ind int) error {
	if prim := t.typ.prim; prim != 0 && t.leafVal == "" {
		return errors.E(nil, "init value not defined for primitive type: %s", prim)
	}
	if v := t.leafVal; v != "" {
		if vf := t.typ.valFmt; vf != nil {
			v2, err := vf(g, v)
			if err != nil {
				return err
			}
			v = v2
		}
		w.WriteString(v)
		return nil
	}

	desc := t.typ.desc
	if desc == nil {
		return errors.E(nil, "internal error? value neither primitive nor compound type")
	}

	impSpec, err := g.descInfo.ImportSpec(desc)
	if err != nil {
		return err
	}
	g.imports[impSpec] = true

	// map field name to oneof name
	var oneofs map[string]string
	if msg, ok := desc.(*descriptor.DescriptorProto); ok {
		oneofs = map[string]string{}
		for _, f := range msg.Field {
			if f.OneofIndex != nil {
				oneofs[f.GetName()] = msg.OneofDecl[*f.OneofIndex].GetName()
			}
		}
	}

	indent := func(j int) {
		for i := 0; i < j; i++ {
			w.WriteByte('\t')
		}
	}

	tvals := t.vals
	writeKey := true
	if t.typ.repeated {
		// We allow array elements to be declared in any order,
		// we reorder it here so we don't need to uglily write the numbers in the
		// composite literal.
		tlen := len(t.vals)
		tvals = make([]*initTree, tlen)

		for i, k := range t.keys {
			n, err := strconv.Atoi(k)
			if err != nil {
				return errors.E(err, "not an array index: %q", k)
			}
			if n < 0 {
				return errors.E(nil, "array index cannot be negative: %d", n)
			}
			if n >= tlen {
				return errors.E(nil, "holes in arrays not allowed; got index %d but length %d", n, tlen)
			}
			tvals[n] = t.vals[i]
		}
		writeKey = false
	}

	// TODO(pongad): handle primitive array
	var typPrefix string
	switch {
	case t.typ.repeated:
		typPrefix = "[]*"
	default:
		typPrefix = "&"
	}

	fmt.Fprintf(w, "%s%s.%s{\n", typPrefix, impSpec.Name, desc.GetName())
	for i, k := range t.keys {
		indent(ind + 1)

		var closeBrace bool
		if oneof, ok := oneofs[k]; ok {
			fmt.Fprintf(w, "%s: &%s.%s_%s{\n", snakeToPascal(oneof), impSpec.Name, desc.GetName(), snakeToPascal(k))
			closeBrace = true
			indent(ind + 2)
		}

		if writeKey {
			w.WriteString(snakeToPascal(k))
			w.WriteString(": ")
		}

		if err := tvals[i].print(w, g, ind+1); err != nil {
			return err
		}

		if closeBrace {
			w.WriteString(",\n")
			indent(ind + 1)
			w.WriteByte('}')
		}
		w.WriteString(",\n")
	}
	indent(ind)
	w.WriteString("}")
	return nil
}

func snakeToPascal(s string) string {
	var sb strings.Builder
	cap := true
	for _, r := range s {
		if r == '_' {
			cap = true
		} else if cap {
			cap = false
			sb.WriteRune(unicode.ToUpper(r))
		} else {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

func snakeToCamel(s string) string {
	var sb strings.Builder
	cap := false
	for _, r := range s {
		if r == '_' {
			cap = true
		} else if cap {
			cap = false
			sb.WriteRune(unicode.ToUpper(r))
		} else {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}
