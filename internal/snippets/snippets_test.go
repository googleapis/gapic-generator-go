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
	"testing"
)

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
