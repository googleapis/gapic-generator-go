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
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/googleapis/gapic-generator-go/internal/gengapic"
	"github.com/googleapis/gapic-generator-go/internal/gensample"
)

func main() {
	reqBytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	var genReq plugin.CodeGeneratorRequest
	if err := genReq.Unmarshal(reqBytes); err != nil {
		log.Fatal(err)
	}

	genResp, err := gengapic.Gen(&genReq)
	if err != nil {
		genResp.Error = proto.String(err.Error())
	}

	sampleResp, err := gensample.PluginEntry(&genReq)
	if err != nil {
		sampleResp.Error = proto.String(err.Error())
	}

	genResp = merge(genResp, sampleResp)

	outBytes, err := proto.Marshal(genResp)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stdout.Write(outBytes); err != nil {
		log.Fatal(err)
	}
}

func merge(gapicResp *plugin.CodeGeneratorResponse, sampleResp *plugin.CodeGeneratorResponse) *plugin.CodeGeneratorResponse {
	if gapicResp.GetError() != "" {
		return gapicResp
	}
	if sampleResp.GetError() != "" {
		return sampleResp
	}
	resp := plugin.CodeGeneratorResponse{}
	resp.File = append(resp.File, gapicResp.GetFile()...)
	resp.File = append(resp.File, sampleResp.GetFile()...)
	return &resp
}
