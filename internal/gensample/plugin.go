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

package gensample

import (
	"path"
	"path/filepath"
	"strings"

	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/googleapis/gapic-generator-go/internal/errors"
)

const (
	paramError = "need parameter in format: go-gapic-package=client/import/path;packageName"
)

// PluginEntry is the entry point of SampleGen as a protoc plugin. If gapic-generator-go
// is called as a protoc plugin with the intention to generate samples,
// it will eventually call this function to do so.
func PluginEntry(genReq *plugin.CodeGeneratorRequest) (*plugin.CodeGeneratorResponse, error) {
	var (
		gapicPkg     string
		gapicFname   string
		sampleFnames []string
	)

	// Always formats the output code if runs as a protoc plugin
	nofmt := false

	resp := plugin.CodeGeneratorResponse{}
	if genReq.Parameter == nil {
		return &resp, errors.E(nil, paramError)
	}
	for _, s := range strings.Split(*genReq.Parameter, ",") {
		if e := strings.IndexByte(s, '='); e > 0 {
			switch s[:e] {
			case "gapic-config":
				gapicFname = s[e+1:]

			case "sample":
				sample := s[e+1:]
				sampleFnames = append(sampleFnames, sample)

			case "go-gapic-package":
				gapicPkg = s[e+1:]
			}
		}
	}
	gen, err := InitGen(genReq.GetProtoFile(), sampleFnames, gapicFname, gapicPkg, nofmt)
	if err != nil {
		return &resp, err
	}

	p := strings.IndexByte(gapicPkg, ';')
	if p < 0 {
		return &resp, errors.E(nil, paramError)
	}
	outDir := filepath.FromSlash(gapicPkg[:p])

	gen.GenMethodSamples()
	for fname, content := range gen.Outputs {
		fullPath := path.Join(outDir, "samples", fname)
		contentStr := string(content)
		resp.File = append(resp.File, &plugin.CodeGeneratorResponse_File{
			Name:    &fullPath,
			Content: &contentStr,
		})
	}
	return &resp, nil
}
