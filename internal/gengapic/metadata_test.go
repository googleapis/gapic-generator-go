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
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/genproto/googleapis/gapic/metadata"
	metadatapb "google.golang.org/genproto/googleapis/gapic/metadata"
)

func TestAddMetadataServiceForTransport(t *testing.T) {
	for _, tst := range []struct {
		service, lib string
		init, want   *metadatapb.GapicMetadata
	}{
		{
			service: "LibraryService",
			lib:     "LibraryService",
			init: &metadatapb.GapicMetadata{
				Services: make(map[string]*metadatapb.GapicMetadata_ServiceForTransport),
			},
			want: &metadatapb.GapicMetadata{
				Services: map[string]*metadatapb.GapicMetadata_ServiceForTransport{
					"LibraryService": {
						Clients: map[string]*metadatapb.GapicMetadata_ServiceAsClient{
							"grpc": {
								LibraryClient: "LibraryServiceClient",
								Rpcs:          make(map[string]*metadata.GapicMetadata_MethodList),
							},
						},
					},
				},
			},
		},
		{
			service: "LibraryService",
			lib:     "LibraryService",
			init: &metadatapb.GapicMetadata{
				Services: map[string]*metadatapb.GapicMetadata_ServiceForTransport{
					"LibraryService": {
						Clients: map[string]*metadatapb.GapicMetadata_ServiceAsClient{
							"rest": {
								LibraryClient: "LibraryServiceRestClient",
								Rpcs:          make(map[string]*metadata.GapicMetadata_MethodList),
							},
						},
					},
				},
			},
			want: &metadatapb.GapicMetadata{
				Services: map[string]*metadatapb.GapicMetadata_ServiceForTransport{
					"LibraryService": {
						Clients: map[string]*metadatapb.GapicMetadata_ServiceAsClient{
							"grpc": {
								LibraryClient: "LibraryServiceClient",
								Rpcs:          make(map[string]*metadata.GapicMetadata_MethodList),
							},
							"rest": {
								LibraryClient: "LibraryServiceRestClient",
								Rpcs:          make(map[string]*metadata.GapicMetadata_MethodList),
							},
						},
					},
				},
			},
		},
		{
			service: "LibraryService",
			lib:     "",
			init: &metadatapb.GapicMetadata{
				Services: map[string]*metadatapb.GapicMetadata_ServiceForTransport{
					"LibraryService": {
						Clients: map[string]*metadatapb.GapicMetadata_ServiceAsClient{
							"rest": {
								LibraryClient: "LibraryServiceRestClient",
								Rpcs:          make(map[string]*metadata.GapicMetadata_MethodList),
							},
						},
					},
				},
			},
			want: &metadatapb.GapicMetadata{
				Services: map[string]*metadatapb.GapicMetadata_ServiceForTransport{
					"LibraryService": {
						Clients: map[string]*metadatapb.GapicMetadata_ServiceAsClient{
							"grpc": {
								LibraryClient: "Client",
								Rpcs:          make(map[string]*metadata.GapicMetadata_MethodList),
							},
							"rest": {
								LibraryClient: "LibraryServiceRestClient",
								Rpcs:          make(map[string]*metadata.GapicMetadata_MethodList),
							},
						},
					},
				},
			},
		},
	} {
		g := generator{
			metadata: tst.init,
		}
		g.addMetadataServiceForTransport(tst.service, "grpc", tst.lib)

		if diff := cmp.Diff(g.metadata, tst.want, cmp.Comparer(proto.Equal)); diff != "" {
			t.Errorf("addMetadataServiceForTransport(%q, %q, %q): got(-),want(+):\n%s", tst.service, "grpc", tst.lib, diff)
		}
	}
}

func TestAddMetadataMethod(t *testing.T) {
	for _, tst := range []struct {
		service, rpc string
		init, want   *metadatapb.GapicMetadata
	}{
		{
			service: "LibraryService",
			rpc:     "GetBook",
			init: &metadatapb.GapicMetadata{
				Services: map[string]*metadatapb.GapicMetadata_ServiceForTransport{
					"LibraryService": {
						Clients: map[string]*metadatapb.GapicMetadata_ServiceAsClient{
							"grpc": {
								LibraryClient: "LibraryServiceClient",
								Rpcs:          make(map[string]*metadata.GapicMetadata_MethodList),
							},
						},
					},
				},
			},
			want: &metadatapb.GapicMetadata{
				Services: map[string]*metadatapb.GapicMetadata_ServiceForTransport{
					"LibraryService": {
						Clients: map[string]*metadatapb.GapicMetadata_ServiceAsClient{
							"grpc": {
								LibraryClient: "LibraryServiceClient",
								Rpcs: map[string]*metadata.GapicMetadata_MethodList{
									"GetBook": {Methods: []string{"GetBook"}},
								},
							},
						},
					},
				},
			},
		},
		{
			service: "LibraryService",
			rpc:     "GetBook",
			init: &metadatapb.GapicMetadata{
				Services: map[string]*metadatapb.GapicMetadata_ServiceForTransport{
					"LibraryService": {
						Clients: map[string]*metadatapb.GapicMetadata_ServiceAsClient{
							"grpc": {
								LibraryClient: "LibraryServiceClient",
								Rpcs: map[string]*metadata.GapicMetadata_MethodList{
									"ListBooks": {Methods: []string{"ListBooks"}},
								},
							},
						},
					},
				},
			},
			want: &metadatapb.GapicMetadata{
				Services: map[string]*metadatapb.GapicMetadata_ServiceForTransport{
					"LibraryService": {
						Clients: map[string]*metadatapb.GapicMetadata_ServiceAsClient{
							"grpc": {
								LibraryClient: "LibraryServiceClient",
								Rpcs: map[string]*metadata.GapicMetadata_MethodList{
									"GetBook":   {Methods: []string{"GetBook"}},
									"ListBooks": {Methods: []string{"ListBooks"}},
								},
							},
						},
					},
				},
			},
		},
	} {
		g := generator{
			metadata: tst.init,
		}
		g.addMetadataMethod(tst.service, "grpc", tst.rpc)

		if diff := cmp.Diff(g.metadata, tst.want, cmp.Comparer(proto.Equal)); diff != "" {
			t.Errorf("addMetadataMethod(%q, %q, %q): got(-),want(+):\n%s", tst.service, "grpc", tst.rpc, diff)
		}
	}
}
