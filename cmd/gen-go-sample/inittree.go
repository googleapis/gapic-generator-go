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
	"io"
	"strings"
	"text/scanner"
	"unicode"

	"github.com/googleapis/gapic-generator-go/internal/errors"
)

// initTree represents a node in the initialization tree.
type initTree struct {
	// initTree is either a composite value (struct/array/map), where keys and vals are set,
	// or a simple value, where leafVal is set.

	// Use array representation; we need order, and we probably won't
	// have many pairs anyway.
	keys []string
	vals []*initTree

	// Text of the literal. If the literal is a string, it's already quoted.
	leafVal string
}

func (t *initTree) get(k string) *initTree {
	for i, key := range t.keys {
		if k == key {
			return t.vals[i]
		}
	}

	v := new(initTree)
	t.keys = append(t.keys, k)
	t.vals = append(t.vals, v)
	return v
}

func (t *initTree) Parse(txt string) error {
	var sc scanner.Scanner
	sc.Init(strings.NewReader(txt))
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

	// Just do simple structs for now.
	// TODO(pongad): allow map and array index.
	// spec = ident { '.' ident } [ '=' value ] .

	if sc.Scan() != scanner.Ident {
		return report(errors.E(nil, "expected ident, found %q", sc.TokenText()))
	}
	t = t.get(sc.TokenText())

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
			t = t.get(sc.TokenText())

		default:
			return report(errors.E(nil, "unexpected %q", sc.TokenText()))
		}
	}

	// TODO(pongad): check that value and type match
	// TODO(pongad): check that we don't write two values onto the same node
	// TODO(pongad): handle resource names
	// TODO(pongad): properly validate and print enums
equal:
	switch r := sc.Scan(); r {
	case scanner.Int, scanner.Float, scanner.String, scanner.Ident:
		t.leafVal = sc.TokenText()
	default:
		return report(errors.E(nil, "expected value, found %q", sc.TokenText()))
	}

	if sc.Scan() != scanner.EOF {
		return report(errors.E(nil, "expected EOF, found %q", sc.TokenText()))
	}
	return err
}

func (t *initTree) Print(w io.Writer) error {
	bw := bufio.NewWriter(w)
	t.print(bw, 1)
	return bw.Flush()
}

func (t *initTree) print(w *bufio.Writer, ind int) {
	if v := t.leafVal; v != "" {
		w.WriteString(v)
		return
	}

	// TODO(pongad): Figure out how to print type.
	w.WriteString("TYPE{\n")
	for i, k := range t.keys {
		for i := 0; i < ind; i++ {
			w.WriteByte('\t')
		}
		snakeToCapCamel(w, k)
		w.WriteString(": ")
		t.vals[i].print(w, ind+1)
		w.WriteString(",\n")
	}
	w.WriteString("}")
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
