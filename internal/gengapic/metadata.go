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

// Initializes the service metadata entry with the API version, if set.
// Per-transport clients will be added as they are generated.
func (g *generator) addMetadataServiceEntry(service, apiVersion string) {
	if g.metadata.Services == nil {
		g.metadata.Services = make(map[string]*metadata.GapicMetadata_ServiceForTransport)
	}
	_, ok := g.metadata.GetServices()[service]
	if !ok {
		s := &metadata.GapicMetadata_ServiceForTransport{
			Clients:    make(map[string]*metadata.GapicMetadata_ServiceAsClient),
			ApiVersion: apiVersion,
		}
		g.metadata.Services[service] = s
	}
}

// Adds a metadata structure for the (service, transport) combination.
// Will exit early if addMetadataServiceEntry is not called prior to this.
// This method is idempotent.
func (g *generator) addMetadataServiceForTransport(service, transport, lib string) {
	if g.metadata.Services == nil {
		return
	}

	s, ok := g.metadata.GetServices()[service]
	if !ok {
		return
	}

	if _, ok := s.Clients[transport]; !ok {
		s.Clients[transport] = &metadata.GapicMetadata_ServiceAsClient{
			// The "Client" part of the generated type's name is hard-coded in the
			// generator so we need to append it to the lib name.
			LibraryClient: lib + "Client",
			Rpcs:          make(map[string]*metadata.GapicMetadata_MethodList),
		}
	}
}

// Adds a metadata service transport client method entry for the given RPC.
// Will exit early if addMetadataServiceEntry or addMetadaServiceForTransport
// are not called prior to this.
func (g *generator) addMetadataMethod(service, transport, rpc string) {
	if g.metadata.Services == nil {
		return
	}
	s, ok := g.metadata.GetServices()[service]
	if !ok {
		return
	}
	c, ok := s.GetClients()[transport]
	if !ok {
		return
	}
	if c.GetRpcs() == nil {
		c.Rpcs = make(map[string]*metadata.GapicMetadata_MethodList)
	}
	// There is only one method per RPC on a generated Go client, with the same name as the RPC.
	c.GetRpcs()[rpc] = &metadata.GapicMetadata_MethodList{
		Methods: []string{rpc},
	}
}
