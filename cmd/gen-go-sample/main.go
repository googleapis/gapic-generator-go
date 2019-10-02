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
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/gapic-generator-go/internal/errors"
	"github.com/googleapis/gapic-generator-go/internal/gensample"
)

const (
	paramError = "need parameter in format: go-gapic-package=client/import/path;packageName"
)

type SampleValue []string

func (v *SampleValue) String() string {
	return strings.Join(*v, ",")
}

func (v *SampleValue) Set(value string) error {
	*v = append(*v, value)
	return nil
}

func main() {
	descFname := flag.String("desc", "", "proto descriptor")
	gapicFname := flag.String("gapic", "", "gapic config")
	clientPkg := flag.String("clientpkg", "", "the package of the client, in format 'url/to/client/pkg;name'")
	nofmt := flag.Bool("nofmt", false, "skip gofmt, useful for debugging code with syntax error")
	outDir := flag.String("o", ".", "directory to write samples to")

	var sampleFnames SampleValue
	flag.Var(&sampleFnames, "sample", "path to a sample config file. There can be more than one --sample flag.")

	flag.Parse()

	descBytes, err := ioutil.ReadFile(*descFname)
	if err != nil {
		log.Fatal(errors.E(err, "cannot read proto descriptor file"))
	}

	var desc descriptor.FileDescriptorSet
	if err := proto.Unmarshal(descBytes, &desc); err != nil {
		log.Fatal(errors.E(err, "error reading proto descriptor file"))
	}

	gen, err := gensample.InitGen(desc.GetFile(), sampleFnames, *gapicFname, *clientPkg, *nofmt)
	if err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(*outDir, 0755); err != nil {
		log.Fatal(err)
	}

	gen.GenMethodSamples()
	for fname, content := range gen.Outputs {
		if err := ioutil.WriteFile(filepath.Join(*outDir, fname), content, 0644); err != nil {
			log.Fatal(err)
		}
	}
}
