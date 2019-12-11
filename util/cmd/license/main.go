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
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var sentinels = [][]byte{
	[]byte("Copyright"),
	[]byte("Google LLC"),
	[]byte("Apache License"),
	[]byte("limitations under the License."),
}

func main() {
	var buf [600]byte
	exitCode := 0

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error: unable to get working dir: %+v", err)
	}

	// Recursively walk the project and check the license header of all
	// non-Protobuf Go files.
	err = filepath.Walk(pwd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, ".pb.go") {
			return nil
		}

		clean, err := check(path, buf)
		if !clean {
			exitCode = 1
		}

		return err
	})
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(exitCode)
}

func check(fname string, buf [600]byte) (bool, error) {
	clean := true
	f, err := os.Open(fname)
	if err != nil {
		return false, err
	}

	n, err := io.ReadFull(f, buf[:])
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return false, err
	}

	f.Close()

	for _, s := range sentinels {
		if bytes.Index(buf[:n], s) < 0 {
			log.Printf("invalid license header, can't find %q: %s", s, fname)
			clean = false
		}
	}

	return clean, nil
}
