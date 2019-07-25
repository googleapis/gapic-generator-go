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
	"bytes"
	"fmt"
	"strings"
	"text/scanner"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/gapic-generator-go/internal/errors"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
)

// symTab implements a symbol table.
//
// Some languages scope variables at function level. We track this with universe table.
// Some languages scope variables to the curly-braces. We track this with parent and scope tables.
//
// So that our samples work with the first kind of languages, we only allow declaration
// iff the name hasn't been declared in the universe scope.
//
// So that our samples work with the second kind of languages, we only allow variable references
// in the recursive closure of the current scope.
type symTab struct {
	universe map[string]bool
	parent   *symTab
	scope    map[string]initType
}

func newSymTab(parent *symTab) *symTab {
	if parent == nil {
		return &symTab{
			universe: map[string]bool{},
			scope:    map[string]initType{},
		}
	}

	return &symTab{
		universe: parent.universe,
		parent:   parent,
		scope:    map[string]initType{},
	}
}

func (st *symTab) get(s string) (initType, bool) {
	for st != nil {
		if it, ok := st.scope[s]; ok {
			return it, true
		}
		st = st.parent
	}
	return initType{}, false
}

func (st *symTab) put(ident string, typ initType) error {
	if st.universe[ident] {
		return errors.E(nil, "variable already defined: %s", ident)
	}
	st.universe[ident] = true
	st.scope[ident] = typ
	return nil
}

func (st *symTab) disambiguate(ident string, typ initType) string {
	base := ident
	sf := 1
	_, ok := st.scope[ident]
	for ok {
		sf++
		ident = fmt.Sprintf("%s%d", base, sf)
		_, ok = st.scope[ident]
	}
	if err := st.put(ident, typ); err != nil {
		panic("bad state: ident shouldn't have existed")
	}
	return ident
}

func writeOutputSpec(out OutputSpec, st *symTab, gen *generator) error {
	used := 0
	var err error

	if d := out.Define; d != "" {
		used++
		err = writeDefine(d, st, gen)
	}
	if p := out.Print; len(p) > 0 {
		used++
		err = writePrint(p[0], p[1:], st, gen)
	}
	if l := out.Loop; l != nil {
		used++
		err = errors.E(nil, "")
		if l.Collection != "" && l.Map != "" {
			err = errors.E(nil, "only one of collection and map should be set")
		} else if l.Collection != "" {
			err = writeLoop(l, st, gen)
		} else if l.Map != "" {
			err = writeMap(l, st, gen)
		}
	}
	if dp := out.WriteFile; dp != nil {
		used++
		err = writeDump(dp.FileName[0], dp.FileName[1:], dp.Contents, st, gen)
	}
	if c := out.Comment; len(c) > 0 {
		used++
		err = writeComment(c[0], c[1:], gen)
	}

	if used == 0 {
		return errors.E(nil, "OutputSpec not defined")
	}
	if used > 1 {
		return errors.E(nil, "OutputSpec cannot be defined multiple times: %v", out)
	}
	return err
}

// define = ident '=' path .
func writeDefine(txt string, st *symTab, gen *generator) error {
	sc, report := initScanner(txt)

	if sc.Scan() != scanner.Ident {
		return report(errors.E(nil, "expecting ident, got %s", sc.TokenText()))
	}

	lhs := sc.TokenText()
	if st.universe[lhs] {
		return report(errors.E(nil, "variable already defined: %q", lhs))
	}

	if sc.Scan() != '=' {
		return report(errors.E(nil, "expecting '=', got %s", sc.TokenText()))
	}

	path, typ, err := writePath(sc, st, gen.descInfo)
	if err != nil {
		return report(err)
	}
	gen.pt.Printf("%s := %s", snakeToCamel(lhs), path)
	return st.put(lhs, typ)
}

var fmtStrReplacer = strings.NewReplacer("%%", "%%", "%s", "%v")

func writePrint(pFmt string, pArgs []string, st *symTab, gen *generator) error {
	argList, err := writePaths(pArgs, st, gen)
	if err != nil {
		return err
	}

	gen.pt.Printf("fmt.Printf(%q%s)", fmtStrReplacer.Replace(pFmt)+"\n", argList)
	gen.imports[pbinfo.ImportSpec{Path: "fmt"}] = true

	return nil
}

func writeDump(fnFmt string, fnArgs []string, contPath string, st *symTab, gen *generator) error {
	argList, err := writePaths(fnArgs, st, gen)
	if err != nil {
		return err
	}

	fn := fnFmt
	if len(argList) != 0 {
		fn = fmt.Sprintf("fmt.Sprintf(%q%s)", fmtStrReplacer.Replace(fnFmt), argList)
	}

	sc, report := initScanner(contPath)
	cont, typ, err := writePath(sc, st, gen.descInfo)
	if err != nil {
		return report(err)
	}

	if typ.prim == descriptor.FieldDescriptorProto_TYPE_BYTES {
		gen.pt.Printf("if err := ioutil.WriteFile(%q, %s, 0644); err != nil {\n", fn, cont)
	} else if typ.prim == descriptor.FieldDescriptorProto_TYPE_STRING && !typ.repeated {
		gen.pt.Printf("if err := ioutil.WriteFile(%s, bytes[](%s), 0644), err != nil {\n", fn, cont)
	} else {
		return errors.E(nil, "illegal type for file contents, expecting string or bytes, got %v", typ)
	}

	gen.pt.Printf("  return err")
	gen.pt.Printf("}")

	gen.imports[pbinfo.ImportSpec{Path: "ioutil"}] = true
	return nil
}

func writeLoop(l *LoopSpec, st *symTab, gen *generator) error {
	if l.Variable == "" {
		return errors.E(nil, "variable not specified for looping over arrays")
	}

	p := gen.pt.Printf

	sc, report := initScanner(l.Collection)
	path, typ, err := writePath(sc, st, gen.descInfo)
	if err = report(err); err != nil {
		return err
	}

	p("for _, %s := range %s {", snakeToCamel(l.Variable), path)

	typ.repeated = false
	stInner := newSymTab(st)
	stInner.put(l.Variable, typ)

	for _, b := range l.Body {
		if err := writeOutputSpec(b, stInner, gen); err != nil {
			return err
		}
	}
	p("}")
	return nil
}

func writeMap(l *LoopSpec, st *symTab, gen *generator) error {
	if l.Key == "" && l.Value == "" {
		return errors.E(nil, "at least one of key and value should be specified for looping over maps")
	}

	p := gen.pt.Printf

	sc, report := initScanner(l.Map)
	path, typ, err := writePath(sc, st, gen.descInfo)
	if err = report(err); err != nil {
		return err
	}

	if typ.keyType == nil || typ.valueType == nil {
		return errors.E(nil, "%s is not a map field", l.Map)
	}
	keyType := *typ.keyType
	valueType := *typ.valueType

	stInner := newSymTab(st)

	if l.Value == "" {
		p("for %s := range %s {", snakeToCamel(l.Key), path)
		stInner.put(l.Key, keyType)
	} else if l.Key == "" {
		p("for _, %s := range %s {", snakeToCamel(l.Value), path)
		stInner.put(l.Value, valueType)
	} else {
		p("for %s, %s := range %s {", snakeToCamel(l.Key), snakeToCamel(l.Value), path)
		stInner.put(l.Key, keyType)
		stInner.put(l.Value, valueType)
	}

	for _, b := range l.Body {
		if err := writeOutputSpec(b, stInner, gen); err != nil {
			return err
		}
	}
	p("}")
	return nil
}

func writeComment(cmtFmt string, cmtArgs []string, gen *generator) error {
	var buf bytes.Buffer
	args := make([]interface{}, len(cmtArgs))
	for i := range cmtArgs {
		args[i] = snakeToCamel(cmtArgs[i])
	}

	if _, err := fmt.Fprintf(&buf, cmtFmt, args...); err != nil {
		return errors.E(err, "comment spec: bad format")
	}
	buf.WriteString("\n")
	prependLines(&buf, "// ", false)
	cmts := strings.Split(buf.String(), "\n")
	for i, c := range cmts {
		if i == len(cmts)-1 {
			continue
		}
		gen.pt.Printf(c)
	}
	return nil
}

func writePath(sc *scanner.Scanner, st *symTab, info pbinfo.Info) (string, initType, error) {
	if sc.Scan() != scanner.Ident {
		return "", initType{}, errors.E(nil, "expecting ident, got %s", sc.TokenText())
	}

	rootVar := sc.TokenText()
	rootTyp, ok := st.get(rootVar)
	if !ok {
		return "", initType{}, errors.E(nil, "variable not found: %q", rootVar)
	}

	itRoot := &initTree{typ: rootTyp}
	itLeaf, _, err := itRoot.parsePathRest(sc, info)
	if err != nil {
		return "", initType{}, err
	}
	if sc.Scan() != scanner.EOF {
		return "", initType{}, errors.E(nil, "expected EOF, found %q", sc.TokenText())
	}

	if rootVar == "$resp" {
		rootVar = "resp"
	} else {
		rootVar = snakeToCamel(rootVar)
	}

	var sb strings.Builder
	sb.WriteString(rootVar)

	// TODO(pongad): This doesn't handle oneofs properly.
	for it := itRoot; len(it.keys) > 0; it = it.vals[0] {
		// Use Get method instead of direct field access so we properly deal with unset messages.
		sb.WriteString(".Get")
		sb.WriteString(snakeToPascal(it.keys[0]))
		sb.WriteString("()")
	}

	return sb.String(), itLeaf.typ, nil
}

// writePaths translates each path into go field accessors, and joins them by commas.
func writePaths(args []string, st *symTab, gen *generator) (string, error) {
	var sb strings.Builder
	for _, arg := range args {
		sb.WriteString(", ")

		sc, report := initScanner(arg)
		path, _, err := writePath(sc, st, gen.descInfo)
		if err != nil {
			return "", report(err)
		}
		sb.WriteString(path)
	}
	return sb.String(), nil
}
