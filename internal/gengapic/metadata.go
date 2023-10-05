// Copyright 2021 Google LLC
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

package gengapic

import (
	"regexp"

	"google.golang.org/genproto/googleapis/gapic/metadata"
	"google.golang.org/protobuf/encoding/protojson"
)

var spaceSanitizerRegex = regexp.MustCompile(`:\s*`)

func (g *generator) genGapicMetadataFile() error {
	data, err := protojson.MarshalOptions{Multiline: true}.Marshal(g.metadata)
	if err != nil {
		return err
	}
	// Hack to standardize output from protojson which is currently non-deterministic
	// with spacing after json keys.
	data = spaceSanitizerRegex.ReplaceAll(data, []byte(": "))
	g.pt.Printf("%s", data)
	return nil
}

// Adds a metadata structure for the (service, transport) combination.
// This method is idempotent.
func (g *generator) addMetadataServiceForTransport(service, transport, lib string) {
	s, ok := g.metadata.GetServices()[service]
	if !ok {
		s = &metadata.GapicMetadata_ServiceForTransport{
			Clients: make(map[string]*metadata.GapicMetadata_ServiceAsClient),
		}
		g.metadata.Services[service] = s
	}

	if _, ok := s.Clients[transport]; !ok {
		s.Clients[transport] = &metadata.GapicMetadata_ServiceAsClient{
			// The "Client" part of the generated type's name is hard-coded in the
			// generator so we need to append it to the lib name.
			//
			// TODO(noahdietz): when REGAPIC is added we may need to special-case based
			// on transport.
			LibraryClient: lib + "Client",
			Rpcs:          make(map[string]*metadata.GapicMetadata_MethodList),
		}
	}
}

func (g *generator) addMetadataMethod(service, transport, rpc string) {
	// There is only one method per RPC on a generated Go client, with the same name as the RPC.
	g.metadata.GetServices()[service].GetClients()[transport].GetRpcs()[rpc] = &metadata.GapicMetadata_MethodList{
		Methods: []string{rpc},
	}
}
