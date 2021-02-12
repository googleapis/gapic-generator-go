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
	"google.golang.org/genproto/googleapis/gapic/metadata"
	"google.golang.org/protobuf/encoding/protojson"
)

func (g *generator) genGapicMetadataFile() error {
	data, err := protojson.MarshalOptions{Multiline: true}.Marshal(g.metadata)
	if err != nil {
		return err
	}

	g.pt.Printf("%s", data)
	return nil
}

func (g *generator) addMetadataServiceForTransport(service, transport, lib string) {
	s, ok := g.metadata.GetServices()[service]
	if !ok {
		s = &metadata.GapicMetadata_ServiceForTransport{
			Clients: make(map[string]*metadata.GapicMetadata_ServiceAsClient),
		}
		g.metadata.Services[service] = s
	}

	s.Clients[transport] = &metadata.GapicMetadata_ServiceAsClient{
		LibraryClient: lib,
		Rpcs:          make(map[string]*metadata.GapicMetadata_MethodList),
	}
}

func (g *generator) addMetadataMethod(service, transport, rpc string) {
	// There is only one method per RPC on a generated Go client, with the same name as the RPC.
	g.metadata.GetServices()[service].GetClients()[transport].GetRpcs()[rpc] = &metadata.GapicMetadata_MethodList{
		Methods: []string{rpc},
	}
}
