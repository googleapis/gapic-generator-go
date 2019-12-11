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
	"flag"
	"io"
	"log"
	"os"
)

var sentinels = [][]byte{
	[]byte("Copyright"),
	[]byte("Google LLC"),
	[]byte("Apache License"),
	[]byte("limitations under the License."),
}

func main() {
	flag.Parse()

	var buf [600]byte
	exitCode := 0

	for _, fname := range flag.Args() {
		f, err := os.Open(fname)
		if err != nil {
			log.Fatal(err)
		}

		n, err := io.ReadFull(f, buf[:])
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Fatal(err)
		}

		f.Close()

		for _, s := range sentinels {
			if bytes.Index(buf[:n], s) < 0 {
				log.Printf("file doesn't have license header, can't find %q: %s", s, fname)
				exitCode = 1
			}
		}
	}

	os.Exit(exitCode)
}
