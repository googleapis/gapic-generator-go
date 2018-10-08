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
	"strings"
	"text/scanner"

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

func writeOutSpec(out OutSpec, st *symTab, info pbinfo.Info, w *bufio.Writer) error {
	used := 0
	var err error

	if d := out.Define; d != "" {
		used++
		err = writeDefine(d, st, info, w)
	}
	if p := out.Print; len(p) > 0 {
		used++
		err = writePrint(p[0], p[1:], st, info, w)
	}

	if used == 0 {
		return errors.E(nil, "OutSpec not defined")
	}
	if used > 1 {
		return errors.E(nil, "OutSpec cannot be defined multiple times: %v", out)
	}
	return err
}

// define = ident '=' ident path .
func writeDefine(txt string, st *symTab, info pbinfo.Info, w *bufio.Writer) error {
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
	if sc.Scan() != scanner.Ident {
		return report(errors.E(nil, "expecting ident, got %s", sc.TokenText()))
	}

	rootVar := sc.TokenText()
	rootTyp, ok := st.get(rootVar)
	if !ok {
		return report(errors.E(nil, "variable not found: %q", rootVar))
	}

	itRoot := &initTree{typ: rootTyp}
	itLeaf, _, err := itRoot.parsePathRest(sc, info)
	if err != nil {
		return report(err)
	}
	if sc.Scan() != scanner.EOF {
		return report(errors.E(nil, "expected EOF, found %q", sc.TokenText()))
	}

	st.scope[lhs] = itLeaf.typ
	st.universe[lhs] = true

	if rootVar == "$resp" {
		rootVar = "resp"
	}
	fmt.Fprintf(w, "%s := %s", snakeToCamel(lhs), snakeToCamel(rootVar))
	writePathRest(itRoot, w)
	w.WriteByte('\n')

	return report(nil)
}

var fmtStrReplacer = strings.NewReplacer("%%", "%%", "%s", "%v")

func writePrint(pFmt string, pArgs []string, st *symTab, info pbinfo.Info, w *bufio.Writer) error {
	fmt.Fprintf(w, "fmt.Printf(%q", fmtStrReplacer.Replace(pFmt)+"\n")

	for _, arg := range pArgs {
		w.WriteString(", ")

		sc, report := initScanner(arg)
		if sc.Scan() != scanner.Ident {
			return report(errors.E(nil, "expecting ident, got %s", sc.TokenText()))
		}

		rootVar := sc.TokenText()
		rootTyp, ok := st.get(rootVar)
		if !ok {
			return report(errors.E(nil, "variable not found: %q", rootVar))
		}

		itRoot := &initTree{typ: rootTyp}
		_, _, err := itRoot.parsePathRest(sc, info)
		if err != nil {
			return report(err)
		}
		if sc.Scan() != scanner.EOF {
			return report(errors.E(nil, "expected EOF, found %q", sc.TokenText()))
		}

		if rootVar == "$resp" {
			rootVar = "resp"
		}
		w.WriteString(snakeToCamel(rootVar))
		writePathRest(itRoot, w)
	}
	w.WriteString(")\n")
	return nil
}

func writePathRest(it *initTree, w *bufio.Writer) {
	// TODO(pongad): This doesn't handle oneofs properly.
	for len(it.keys) > 0 {
		w.WriteRune('.')
		w.WriteString(snakeToPascal(it.keys[0]))
		it = it.vals[0]
	}
}
