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

// Package printer implements auto-indenting printer.
package printer

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

// P represents a buffer for print code that tracks indentation.
type P struct {
	buf    bytes.Buffer
	indent int
}

// Reset resets p but retains the underlying storage for use by future Printfs.
func (p *P) Reset() {
	p.buf.Reset()
	p.indent = 0
}

// Printf format-writes to p's buffer. The formatting is similar to the fmt package,
// and the written content can be retrieved with Bytes. Printf is line-oriented: each call
// should correspond to a new line in the output. Printf automatically inserts a newline
// at the end of each call.
//
// Printf automatically keeps track of curly-braces and indents lines accordingly.
// The indentation is not "gofmt-correct" but should print code in "curly brace languages"
// well enough that callers don't need to keep track of indentation level.
// Printf ignores leading and trailing whitespaces in s; callers can use them for visual aid
// without interfering with automatic indentation.
func (p *P) Printf(s string, args ...interface{}) {
	s = strings.TrimSpace(s)
	if s == "" {
		p.buf.WriteByte('\n')
		return
	}

	for i := 0; i < len(s) && s[i] == '}'; i++ {
		p.indent--
	}

	for i := 0; i < p.indent; i++ {
		p.buf.WriteByte('\t')
	}

	fmt.Fprintf(&p.buf, s, args...)
	p.buf.WriteByte('\n')

	for i := len(s) - 1; i >= 0 && s[i] == '{'; i-- {
		p.indent++
	}
}

// Writer returns a writer that writes to p's underlying buffer without performing indentation.
func (p *P) Writer() io.Writer {
	return &p.buf
}

// Bytes returns the bytes written by Printf. It is valid up to the next call to Printf or Reset.
func (p *P) Bytes() []byte {
	return p.buf.Bytes()
}

func (p *P) String() string {
	return p.buf.String()
}

// Len returns the length of the printer buffer.
func (p *P) Len() int {
	return p.buf.Len()
}
