// Copyright 2019 Google LLC
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
	"text/scanner"

	"github.com/googleapis/gapic-generator-go/internal/errors"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
)

func bytesFmt() func(*generator, string) (string, error) {
	return func(_ *generator, s string) (string, error) {
		return "[]byte(" + s + ")", nil
	}
}

// fileInfo keeps track of information we need to generate reading all bytes of a file
// and assign them to a local variable.
type fileInfo struct {
	// fileName is the text of the file name. If it's a string literal, it's already quoted.
	fileName string
	// varName is the name of the local variable to hold the bytes of the file.
	varName string
}

// fileVarName generates a name for the local variable that hold bytes from a local file. It is Only
// used when SampleArgumentName is not provided.
func fileVarName(param string) (string, error) {
	// TODO: if there are multiple fields whose values are to be copied from local files,
	// this method cannot generate a unique variable name for each of them when:
	// 1) SampleArgumentName are not provided for some or all of thse fields
	// 2) these fields have the same name
	// Considering the collision only happens in very rare cases, they are not handled for now
	sc, _ := initScanner("." + param)

	var varName string
	for {
		switch r := sc.Scan(); r {
		case '%':
			return "", errors.E(nil, "bad format for bytes field, unexpected '%%'")
		case '.':
			if r := sc.Scan(); r != scanner.Ident {
				return "", errors.E(nil, "expected ident, found %q", sc.TokenText())
			}
			varName = sc.TokenText()
		case '[':
			if r := sc.Scan(); r != scanner.Int {
				return "", errors.E(nil, "expected int, found %q", sc.TokenText())
			}
			if r := sc.Scan(); r != ']' {
				return "", errors.E(nil, "expected ']', found %q", sc.TokenText())
			}
		case scanner.EOF:
			return snakeToCamel(varName) + "Content", nil
		default:
			return "", errors.E(nil, "unhandled rune: %c", r)
		}
	}
}

func file(info *fileInfo, buf *bytes.Buffer, g *generator) {
	vn := info.varName
	fn := info.fileName
	g.imports[pbinfo.ImportSpec{Path: "io"}] = true
	g.imports[pbinfo.ImportSpec{Path: "io/ioutil"}] = true

	fmt.Fprintf(buf, "file, err = os.Open(%s)\n", fn)
	fmt.Fprintf(buf, "if err != nil {\n")
	fmt.Fprintf(buf, "\tlog.Fatalf(\"Failed to read file: %%v\", err)\n")
	fmt.Fprintf(buf, "}\n")
	fmt.Fprintf(buf, "defer file.Close()\n")
	fmt.Fprintf(buf, "var %s []byte\n", vn)
	fmt.Fprintf(buf, "%s, err = ioutil.ReadAll(file)\n", vn)
	fmt.Fprintf(buf, "if err != nil {\n")
	fmt.Fprintf(buf, "\tlog.Fatalf(\"Failed to read file: %%v\", err)\n")
	fmt.Fprintf(buf, "}\n\n")
}
