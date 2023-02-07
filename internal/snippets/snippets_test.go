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
var serviceConfigName = "bigquerymigration.googleapis.com"
var version = "v2"

func TestNewMetadata(t *testing.T) {
	sm := NewMetadata(protoPkg, libPkg, serviceConfigName)

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
	if shortName := "bigquerymigration"; sm.shortName != shortName {
		t.Errorf("%s: wanted %s, got %s", t.Name(), shortName, sm.shortName)
	}
}

func TestToMetadataJSON(t *testing.T) {
	serviceName1 := "Service1"
	serviceName2 := "Service2"
	methodName1 := "Method1"
	methodName2 := "Method2"
	methodName3 := "Method3"

	sm := NewMetadata(protoPkg, libPkg, serviceConfigName)
	sm.AddService(serviceName1)
	sm.AddService(serviceName2)
	sm.AddMethod(serviceName1, methodName1, 51)
	sm.AddMethod(serviceName1, methodName2, 52)
	sm.AddMethod(serviceName2, methodName3, 53)
	sm.UpdateMethodDoc(serviceName1, methodName1, "methodName1 doc")
	sm.UpdateMethodDoc(serviceName1, methodName2, "methodName2 doc")
	sm.UpdateMethodDoc(serviceName2, methodName3, "methodName3 doc")
	sm.UpdateMethodResult(serviceName1, methodName1, "mypackage.MethodName1Result")
	sm.UpdateMethodResult(serviceName1, methodName2, "mypackage.MethodName2Result")
	sm.UpdateMethodResult(serviceName2, methodName3, "mypackage.MethodName3Result")
	sm.AddParams(serviceName1, methodName1, "mypackage.MethodName1Request")
	sm.AddParams(serviceName1, methodName2, "mypackage.MethodName2Request")
	sm.AddParams(serviceName2, methodName3, "mypackage.MethodName3Request")

	mi := sm.toSnippetMetadata()
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

	if got := len(mi.Snippets); got != 3 {
		t.Errorf("%s: wanted len 3 Snippets, got %d", t.Name(), got)
	}
	for i, snippet := range mi.Snippets {
		want := fmt.Sprintf("bigquerymigration Method%d Sample", i+1)
		if got := snippet.Title; got != want {
			t.Errorf("%s: wanted %s, got %s", t.Name(), want, got)
		}
	}

	json, err := sm.ToMetadataJSON()
	if err != nil {
		t.Fatal(err)
	}

	if got := len(json); got == 0 {
		t.Errorf("%s: wanted non-empty []byte, got 0", t.Name())
	}
}

func TestRegionTag(t *testing.T) {
	protoPkg := "google.cloud.bigquery.migration.v2"
	libPkg := "google.golang.org/genproto/googleapis/cloud/bigquery/migration/v2"
	serviceConfigName := "bigquerymigration.googleapis.com"
	m := NewMetadata(protoPkg, libPkg, serviceConfigName)

	serviceName := "MigrationService"
	methodName := "GetMigrationWorkflow"
	want := "bigquerymigration_v2_generated_MigrationService_GetMigrationWorkflow_sync"
	if got := m.RegionTag(serviceName, methodName); got != want {
		t.Errorf("%s: wanted %s, got %s", t.Name(), want, got)
	}
}
