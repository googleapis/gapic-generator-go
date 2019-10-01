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
	"strings"
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

	var samplePaths SampleValue
	flag.Var(&samplePaths, "sample", "path to a sample config file. There can be more than one --sample flag.")

	flag.Parse()
	gen(*descFname, samplePaths, *gapicFname, *clientPkg, *nofmt, *outDir)
}

func Plugin(genReq *plugin.CodeGeneratorRequest) (*plugin.CodeGeneratorResponse, error) {
	
}
