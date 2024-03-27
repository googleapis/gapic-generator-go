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
	"io"
	"log"
	"os"

	"github.com/googleapis/gapic-generator-go/internal/gencli"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	reqBytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	var genReq pluginpb.CodeGeneratorRequest
	if err := proto.Unmarshal(reqBytes, &genReq); err != nil {
		log.Fatal(err)
	}

	genResp, err := gencli.Gen(&genReq)
	if err != nil {
		genResp.Error = proto.String(err.Error())
	}
	genResp.SupportedFeatures = proto.Uint64(uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL))

	outBytes, err := proto.Marshal(genResp)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stdout.Write(outBytes); err != nil {
		log.Fatal(err)
	}
}
