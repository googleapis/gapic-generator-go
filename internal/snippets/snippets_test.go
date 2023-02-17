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

	"github.com/googleapis/gapic-generator-go/internal/snippets/metadata"
)

var protoPkg = "google.cloud.bigquery.migration.v2"
var libPkg = "cloud.google.com/go/bigquery/migration/apiv2"
var defaultHost = "bigquerymigration.googleapis.com"
var version = "v2"
var pkgName = "migration"

func TestNewMetadata(t *testing.T) {
	sm, err := NewMetadata(protoPkg, libPkg, pkgName)
	if err != nil {
		t.Fatal(err)
	}

	if sm.protoPkg != protoPkg {
		t.Errorf("%s: wanted %s, got %s", t.Name(), protoPkg, sm.protoPkg)
	}
	if sm.libPkg != libPkg {
		t.Errorf("%s: wanted %s, got %s", t.Name(), libPkg, sm.libPkg)
	}
	if got := len(sm.protoServices); got != 0 {
		t.Errorf("%s: wanted empty, got %d", t.Name(), len(sm.protoServices))
	}
	if sm.apiVersion != version {
		t.Errorf("%s: wanted %s, got %s", t.Name(), version, sm.apiVersion)
	}
}

func TestToMetadataJSON(t *testing.T) {
	sm, err := NewMetadata(protoPkg, libPkg, pkgName)
	if err != nil {
		t.Fatal(err)
	}
	regionTagStart := 18
	regionTagEnd := 50
	for i := 0; i < 2; i++ {
		serviceName := fmt.Sprintf("Foo%dService", i)
		methodName := fmt.Sprintf("Bar%dMethod", i)
		sm.AddService(serviceName, defaultHost)
		sm.AddMethod(serviceName, methodName, regionTagEnd)
		sm.UpdateMethodDoc(serviceName, methodName, methodName+" doc")
		sm.UpdateMethodResult(serviceName, methodName, "mypackage."+methodName+"Result")
		sm.AddParams(serviceName, methodName, "mypackage."+methodName+"Request")
	}

	mi := sm.ToMetadataIndex()
	// TODO(chrisdsmith): replace assertions with go-cmp(..., proto.EQUAL)
	cl := mi.ClientLibrary
	if cl.Name != libPkg {
		t.Errorf("%s: wanted %s, got %s", t.Name(), libPkg, cl.Name)
	}
	if cl.Version != VersionPlaceholder {
		t.Errorf("%s: wanted %s, got %s", t.Name(), VersionPlaceholder, cl.Version)
	}
	if cl.Language != metadata.Language_GO {
		t.Errorf("%s: wanted %s, got %s", t.Name(), metadata.Language_GO, cl.Language)
	}
	if got := len(cl.Apis); got != 1 {
		t.Errorf("%s: wanted len 1 Apis, got %d", t.Name(), got)
	}
	if got := cl.Apis[0].Id; got != protoPkg {
		t.Errorf("%s: wanted %s, got %s", t.Name(), protoPkg, got)
	}
	if got := cl.Apis[0].Version; got != version {
		t.Errorf("%s: wanted %s, got %s", t.Name(), version, got)
	}
	if got := len(mi.Snippets); got != 2 {
		t.Errorf("%s: wanted len 2 Snippets, got %d", t.Name(), got)
	}
	for i, snp := range mi.Snippets {
		want := fmt.Sprintf("bigquerymigration_v2_generated_Foo%dService_Bar%dMethod_sync", i, i)
		if got := snp.RegionTag; got != want {
			t.Errorf("%s: wanted %s, got %s", t.Name(), want, got)
		}
		want = fmt.Sprintf("bigquerymigration Bar%dMethod Sample", i)
		if got := snp.Title; got != want {
			t.Errorf("%s: wanted %s, got %s", t.Name(), want, got)
		}
		want = fmt.Sprintf("Bar%dMethod doc", i)
		if got := snp.Description; got != want {
			t.Errorf("%s: wanted %s, got %s", t.Name(), want, got)
		}
		want = fmt.Sprintf("Foo%dClient/Bar%dMethod/main.go", i, i)
		if got := snp.File; got != want {
			t.Errorf("%s: wanted %s, got %s", t.Name(), want, got)
		}
		if snp.Language != metadata.Language_GO {
			t.Errorf("%s: wanted %s, got %s", t.Name(), metadata.Language_GO, snp.Language)
		}
		if snp.Canonical {
			t.Errorf("%s: wanted Canonical false, got true", t.Name())
		}
		cm := snp.ClientMethod
		want = fmt.Sprintf("Bar%dMethod", i)
		if got := cm.ShortName; got != want {
			t.Errorf("%s: wanted %s, got %s", t.Name(), want, got)
		}
		want = fmt.Sprintf("google.cloud.bigquery.migration.v2.Foo%dClient.Bar%dMethod", i, i)
		if got := cm.FullName; got != want {
			t.Errorf("%s: wanted %s, got %s", t.Name(), want, got)
		}
		if cm.Async {
			t.Errorf("%s: wanted Async false, got true", t.Name())
		}
		want = fmt.Sprintf("mypackage.Bar%dMethodResult", i)
		if got := cm.ResultType; got != want {
			t.Errorf("%s: wanted %s, got %s", t.Name(), want, got)
		}
		want = fmt.Sprintf("Foo%dClient", i)
		if got := cm.Client.ShortName; got != want {
			t.Errorf("%s: wanted %s, got %s", t.Name(), want, got)
		}
		want = fmt.Sprintf("google.cloud.bigquery.migration.v2.Foo%dClient", i)
		if got := cm.Client.FullName; got != want {
			t.Errorf("%s: wanted %s, got %s", t.Name(), want, got)
		}
		want = fmt.Sprintf("Bar%dMethod", i)
		if got := cm.Method.ShortName; got != want {
			t.Errorf("%s: wanted %s, got %s", t.Name(), want, got)
		}
		want = fmt.Sprintf("google.cloud.bigquery.migration.v2.Foo%dService.Bar%dMethod", i, i)
		if got := cm.Method.FullName; got != want {
			t.Errorf("%s: wanted %s, got %s", t.Name(), want, got)
		}
		want = fmt.Sprintf("Foo%dService", i)
		if got := cm.Method.Service.ShortName; got != want {
			t.Errorf("%s: wanted %s, got %s", t.Name(), want, got)
		}
		want = fmt.Sprintf("google.cloud.bigquery.migration.v2.Foo%dService", i)
		if got := cm.Method.Service.FullName; got != want {
			t.Errorf("%s: wanted %s, got %s", t.Name(), want, got)
		}
		if got := len(snp.Segments); got != 1 {
			t.Errorf("%s: wanted len 1 Segments, got %d", t.Name(), got)
		}
		if got := snp.Segments[0].Start; got != int32(regionTagStart) {
			t.Errorf("%s: wanted %d, got %d", t.Name(), regionTagStart, got)
		}
		if got := int(snp.Segments[0].End); got != regionTagEnd-1 {
			t.Errorf("%s: wanted 1, got %d", t.Name(), got)
		}
		if got := snp.Segments[0].Type; got != metadata.Snippet_Segment_FULL {
			t.Errorf("%s: wanted metadata.Snippet_Segment_FULL, got %d", t.Name(), got)
		}
		if got := len(cm.Parameters); got != 3 {
			t.Errorf("%s: wanted len 3 Parameters, got %d", t.Name(), got)
		}
		if got := cm.Parameters[0].Type; got != "context.Context" {
			t.Errorf("%s: wanted context.Context, got %s", t.Name(), got)
		}
		if got := cm.Parameters[0].Name; got != "ctx" {
			t.Errorf("%s: wanted ctx, got %s", t.Name(), got)
		}
		want = fmt.Sprintf("mypackage.Bar%dMethodRequest", i)
		if got := cm.Parameters[1].Type; got != want {
			t.Errorf("%s: wanted %s, got %s", t.Name(), want, got)
		}
		if got := cm.Parameters[1].Name; got != "req" {
			t.Errorf("%s: wanted req, got %s", t.Name(), got)
		}
		if got := cm.Parameters[2].Type; got != "...gax.CallOption" {
			t.Errorf("%s: wanted ...gax.CallOption, got %s", t.Name(), got)
		}
		if got := cm.Parameters[2].Name; got != "opts" {
			t.Errorf("%s: wanted opts, got %s", t.Name(), got)
		}
	}

	json, err := sm.ToMetadataJSON()
	if err != nil {
		t.Fatal(err)
	}

	if got := len(json); got == 0 {
		t.Errorf("%s: wanted non-empty []byte, got len 0", t.Name())
	}
}

func TestRegionTag(t *testing.T) {
	protoPkg := "google.cloud.bigquery.migration.v2"
	libPkg := "google.golang.org/genproto/googleapis/cloud/bigquery/migration/v2"
	sm, err := NewMetadata(protoPkg, libPkg, pkgName)
	if err != nil {
		t.Fatal(err)
	}
	serviceName := "MigrationService"
	defaultHost := "bigquerymigration.googleapis.com"
	sm.AddService(serviceName, defaultHost)
	methodName := "GetMigrationWorkflow"
	want := "bigquerymigration_v2_generated_MigrationService_GetMigrationWorkflow_sync"
	if got := sm.RegionTag(serviceName, methodName); got != want {
		t.Errorf("%s: wanted %s, got %s", t.Name(), want, got)
	}
}
