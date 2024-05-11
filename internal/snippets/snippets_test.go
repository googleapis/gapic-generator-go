// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package snippets

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/googleapis/gapic-generator-go/internal/snippets/metadata"
	"github.com/googleapis/gapic-generator-go/internal/testing/sample"
	"google.golang.org/protobuf/proto"
)

func TestNewMetadata(t *testing.T) {
	sm := NewMetadata(sample.ProtoPackagePath, sample.GoPackagePath, sample.GoPackageName)

	if sm.protoPkg != sample.ProtoPackagePath {
		t.Errorf("%s: got %s want %s,", t.Name(), sm.protoPkg, sample.ProtoPackagePath)
	}
	if sm.libPkg != sample.GoPackagePath {
		t.Errorf("%s: got %s want %s", t.Name(), sm.libPkg, sample.GoPackagePath)
	}
	if got := len(sm.protoServices); got != 0 {
		t.Errorf("%s: got %d want empty", t.Name(), len(sm.protoServices))
	}
	if sm.apiVersion != sample.ProtoVersion {
		t.Errorf("%s: got %s want %s", t.Name(), sm.apiVersion, sample.ProtoVersion)
	}
}

func TestToMetadataJSON(t *testing.T) {
	// Build fixture
	sm := NewMetadata(sample.ProtoPackagePath, sample.GoPackagePath, sample.GoPackageName)
	regionTagStart := 18
	regionTagEnd := 50
	for i := 0; i < 2; i++ {
		serviceName := fmt.Sprintf("Foo%dService", i)
		methodName := fmt.Sprintf("Bar%dMethod", i)
		sm.AddService(serviceName, sample.ServiceURL)
		sm.AddMethod(serviceName, methodName, sample.ProtoPackagePath, serviceName, regionTagEnd)
		sm.UpdateMethodDoc(serviceName, methodName, methodName+" doc\n New line.")
		sm.UpdateMethodResult(serviceName, methodName, "mypackage."+methodName+"Result")
		sm.AddParams(serviceName, methodName, "mypackage."+methodName+"Request")
	}

	// Build expectation
	want := &metadata.Index{
		ClientLibrary: &metadata.ClientLibrary{
			Name:     sample.GoPackagePath,
			Version:  VersionPlaceholder,
			Language: metadata.Language_GO,
			Apis: []*metadata.Api{
				{
					Id:      sample.ProtoPackagePath,
					Version: "v2",
				},
			},
		},
	}
	for i := 0; i < 2; i++ {
		snp := &metadata.Snippet{
			RegionTag:   fmt.Sprintf("bigquerymigration_v2_generated_Foo%dService_Bar%dMethod_sync", i, i),
			Title:       fmt.Sprintf("bigquerymigration Bar%dMethod Sample", i),
			Description: fmt.Sprintf("Bar%dMethod doc\nNew line.", i),
			File:        fmt.Sprintf("Foo%dClient/Bar%dMethod/main.go", i, i),
			Language:    metadata.Language_GO,
			Canonical:   false,
			Origin:      *metadata.Snippet_API_DEFINITION.Enum(),
			ClientMethod: &metadata.ClientMethod{
				ShortName:  fmt.Sprintf("Bar%dMethod", i),
				FullName:   fmt.Sprintf("google.cloud.bigquery.migration.v2.Foo%dClient.Bar%dMethod", i, i),
				Async:      false,
				ResultType: fmt.Sprintf("mypackage.Bar%dMethodResult", i),
				Client: &metadata.ServiceClient{
					ShortName: fmt.Sprintf("Foo%dClient", i),
					FullName:  fmt.Sprintf("google.cloud.bigquery.migration.v2.Foo%dClient", i),
				},
				Method: &metadata.Method{
					ShortName: fmt.Sprintf("Bar%dMethod", i),
					FullName:  fmt.Sprintf("google.cloud.bigquery.migration.v2.Foo%dService.Bar%dMethod", i, i),
					Service: &metadata.Service{
						ShortName: fmt.Sprintf("Foo%dService", i),
						FullName:  fmt.Sprintf("google.cloud.bigquery.migration.v2.Foo%dService", i),
					},
				},
				Parameters: []*metadata.ClientMethod_Parameter{
					{
						Type: "context.Context",
						Name: "ctx",
					},
					{
						Type: fmt.Sprintf("mypackage.Bar%dMethodRequest", i),
						Name: "req",
					},
					{
						Type: "...gax.CallOption",
						Name: "opts",
					},
				},
			},
			Segments: []*metadata.Snippet_Segment{
				{
					Start: int32(regionTagStart),
					End:   int32(regionTagEnd - 1),
					Type:  metadata.Snippet_Segment_FULL,
				},
			},
		}
		want.Snippets = append(want.Snippets, snp)
	}

	mi := sm.ToMetadataIndex()
	if diff := cmp.Diff(mi, want, cmp.Comparer(proto.Equal)); diff != "" {
		t.Errorf("ToMetadataIndex(): got(-),want(+):\n%s", diff)
	}

	json, err := sm.ToMetadataJSON()
	if err != nil {
		t.Fatal(err)
	}

	if got := len(json); got == 0 {
		t.Errorf("%s: got len 0, want non-empty []byte", t.Name())
	}
}

func TestRegionTag(t *testing.T) {
	libPkg := "google.golang.org/genproto/googleapis/cloud/bigquery/migration/v2"
	sm := NewMetadata(sample.ProtoPackagePath, libPkg, sample.GoPackagePath)
	sm.AddService(sample.ServiceName, sample.ServiceURL)
	want := "bigquerymigration_v2_generated_MigrationService_GetMigrationWorkflow_sync"
	if got := sm.RegionTag(sample.ServiceName, sample.GetMethod); got != want {
		t.Errorf("%s: got %s want %s", t.Name(), got, want)
	}
}
