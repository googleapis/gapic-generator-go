// Copyright 2023 Google LLC
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
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
)

func TestSortImports(t *testing.T) {
	for _, tst := range []struct {
		name    string
		imports []pbinfo.ImportSpec
		sorted  []pbinfo.ImportSpec
		want    int
	}{
		{
			name: "foo_foopb",
			imports: []pbinfo.ImportSpec{
				{Name: "foopb", Path: "cloud.google.com/go/foopb"},
				{Name: "foo", Path: "cloud.google.com/go/foo"},
			},
			sorted: []pbinfo.ImportSpec{
				{Name: "foo", Path: "cloud.google.com/go/foo"},
				{Name: "foopb", Path: "cloud.google.com/go/foopb"},
			},
			want: 0,
		},
		{
			name: "context_foo_foopb",
			imports: []pbinfo.ImportSpec{
				{Name: "foopb", Path: "cloud.google.com/go/foopb"},
				{Path: "context"},
				{Name: "foo", Path: "cloud.google.com/go/foo"},
			},
			sorted: []pbinfo.ImportSpec{
				{Path: "context"},
				{Name: "foo", Path: "cloud.google.com/go/foo"},
				{Name: "foopb", Path: "cloud.google.com/go/foopb"},
			},
			want: 1,
		},
		{
			name: "context_datacatalog_iampb",
			imports: []pbinfo.ImportSpec{
				{Name: "iampb", Path: "google.golang.org/genproto/googleapis/iam/v1"},
				{Path: "context"},
				{Name: "datacatalog", Path: "cloud.google.com/go/datacatalog/apiv1beta1"},
			},
			sorted: []pbinfo.ImportSpec{
				{Path: "context"},
				{Name: "datacatalog", Path: "cloud.google.com/go/datacatalog/apiv1beta1"},
				{Name: "iampb", Path: "google.golang.org/genproto/googleapis/iam/v1"},
			},
			want: 1,
		},
		{
			name: "context_strings_datacatalog_iampb",
			imports: []pbinfo.ImportSpec{
				{Path: "strings"},
				{Name: "iampb", Path: "google.golang.org/genproto/googleapis/iam/v1"},
				{Path: "context"},
				{Name: "datacatalog", Path: "cloud.google.com/go/datacatalog/apiv1beta1"},
			},
			sorted: []pbinfo.ImportSpec{
				{Path: "context"},
				{Path: "strings"},
				{Name: "datacatalog", Path: "cloud.google.com/go/datacatalog/apiv1beta1"},
				{Name: "iampb", Path: "google.golang.org/genproto/googleapis/iam/v1"},
			},
			want: 2,
		},
	} {
		in := tst.imports
		got := sortImports(in)
		if diff := cmp.Diff(tst.imports, tst.sorted); diff != "" {
			t.Errorf("test %s got(-), want(+):\n%s", tst.name, diff)
		}
		if got != tst.want {
			t.Errorf("test %s got %d, want %d", tst.name, got, tst.want)
		}
	}
}
