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

	"github.com/google/go-cmp/cmp"
	"github.com/googleapis/gapic-generator-go/internal/txtdiff"
	"google.golang.org/genproto/googleapis/gapic/metadata"
	"google.golang.org/protobuf/proto"
)

func TestAddMetadataServiceForTransport(t *testing.T) {
	for _, tst := range []struct {
		service, lib string
		init, want   *metadata.GapicMetadata
	}{
		{
			service: "LibraryService",
			lib:     "LibraryService",
			init: &metadata.GapicMetadata{
				Services: make(map[string]*metadata.GapicMetadata_ServiceForTransport),
			},
			want: &metadata.GapicMetadata{
				Services: map[string]*metadata.GapicMetadata_ServiceForTransport{
					"LibraryService": {
						Clients: map[string]*metadata.GapicMetadata_ServiceAsClient{
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
			init: &metadata.GapicMetadata{
				Services: map[string]*metadata.GapicMetadata_ServiceForTransport{
					"LibraryService": {
						Clients: map[string]*metadata.GapicMetadata_ServiceAsClient{
							"rest": {
								LibraryClient: "LibraryServiceRESTClient",
								Rpcs:          make(map[string]*metadata.GapicMetadata_MethodList),
							},
						},
					},
				},
			},
			want: &metadata.GapicMetadata{
				Services: map[string]*metadata.GapicMetadata_ServiceForTransport{
					"LibraryService": {
						Clients: map[string]*metadata.GapicMetadata_ServiceAsClient{
							"grpc": {
								LibraryClient: "LibraryServiceClient",
								Rpcs:          make(map[string]*metadata.GapicMetadata_MethodList),
							},
							"rest": {
								LibraryClient: "LibraryServiceRESTClient",
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
			init: &metadata.GapicMetadata{
				Services: map[string]*metadata.GapicMetadata_ServiceForTransport{
					"LibraryService": {
						Clients: map[string]*metadata.GapicMetadata_ServiceAsClient{
							"rest": {
								LibraryClient: "LibraryServiceRESTClient",
								Rpcs:          make(map[string]*metadata.GapicMetadata_MethodList),
							},
						},
					},
				},
			},
			want: &metadata.GapicMetadata{
				Services: map[string]*metadata.GapicMetadata_ServiceForTransport{
					"LibraryService": {
						Clients: map[string]*metadata.GapicMetadata_ServiceAsClient{
							"grpc": {
								LibraryClient: "Client",
								Rpcs:          make(map[string]*metadata.GapicMetadata_MethodList),
							},
							"rest": {
								LibraryClient: "LibraryServiceRESTClient",
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
		init, want   *metadata.GapicMetadata
	}{
		{
			service: "LibraryService",
			rpc:     "GetBook",
			init: &metadata.GapicMetadata{
				Services: map[string]*metadata.GapicMetadata_ServiceForTransport{
					"LibraryService": {
						Clients: map[string]*metadata.GapicMetadata_ServiceAsClient{
							"grpc": {
								LibraryClient: "LibraryServiceClient",
								Rpcs:          make(map[string]*metadata.GapicMetadata_MethodList),
							},
						},
					},
				},
			},
			want: &metadata.GapicMetadata{
				Services: map[string]*metadata.GapicMetadata_ServiceForTransport{
					"LibraryService": {
						Clients: map[string]*metadata.GapicMetadata_ServiceAsClient{
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
			init: &metadata.GapicMetadata{
				Services: map[string]*metadata.GapicMetadata_ServiceForTransport{
					"LibraryService": {
						Clients: map[string]*metadata.GapicMetadata_ServiceAsClient{
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
			want: &metadata.GapicMetadata{
				Services: map[string]*metadata.GapicMetadata_ServiceForTransport{
					"LibraryService": {
						Clients: map[string]*metadata.GapicMetadata_ServiceAsClient{
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

func TestGenGapicMetadataFile_standardized(t *testing.T) {
	g := generator{
		metadata: &metadata.GapicMetadata{
			Schema:         "schema",
			Comment:        "comment",
			Language:       "language",
			ProtoPackage:   "packagename",
			LibraryPackage: "lib",
			Services: map[string]*metadata.GapicMetadata_ServiceForTransport{
				"FooService": {
					Clients: map[string]*metadata.GapicMetadata_ServiceAsClient{
						"grpc": {
							LibraryClient: "libClient",
							Rpcs: map[string]*metadata.GapicMetadata_MethodList{
								"GetBook": {Methods: []string{"GetBook"}},
							},
						},
					},
				},
			},
		},
	}
	if err := g.genGapicMetadataFile(); err != nil {
		t.Fatalf("got genGapicMetadataFile() = %v, want nil", err)
	}
	txtdiff.Diff(t, g.pt.String(), "testdata/metadata.want")
}
