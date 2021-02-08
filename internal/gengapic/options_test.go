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
	"reflect"
	"testing"
)

func TestParseOptions(t *testing.T) {
	for _, tst := range []struct {
		param        string
		expectedOpts *options
		expectErr    bool
	}{
		{
			param: "transport=grpc,go-gapic-package=path;pkg",
			expectedOpts: &options{
				transports: []transport{grpc},
				pkgPath:    "path",
				pkgName:    "pkg",
				outDir:     "path",
			},
			expectErr: false,
		},
		{
			param: "transport=rest+grpc,go-gapic-package=path;pkg",
			expectedOpts: &options{
				transports: []transport{rest, grpc},
				pkgPath:    "path",
				pkgName:    "pkg",
				outDir:     "path",
			},
			expectErr: false,
		},
		{
			param: "go-gapic-package=path;pkg",
			expectedOpts: &options{
				transports: []transport{grpc},
				pkgPath:    "path",
				pkgName:    "pkg",
				outDir:     "path",
			},
		},
		{
			param: "module=path,go-gapic-package=path/to/out;pkg",
			expectedOpts: &options{
				transports:   []transport{grpc},
				pkgPath:      "path/to/out",
				pkgName:      "pkg",
				outDir:       "to/out",
				modulePrefix: "path",
			},
			expectErr: false,
		},
		{
			param:     "transport=tcp,go-gapic-package=path;pkg",
			expectErr: true,
		},
		{
			param:     "go-gapic-package=pkg;",
			expectErr: true,
		},
		{
			param:     "go-gapic-package=;path",
			expectErr: true,
		},
		{
			param:     "go-gapic-package=bogus",
			expectErr: true,
		},
		{
			param:     "module=different_path,go-gapic-package=path;pkg",
			expectErr: true,
		},
	} {
		opts, err := parseOptions(&tst.param)
		if tst.expectErr && err == nil {
			t.Errorf("ParseOptions(%s) expected error", tst.param)
			continue
		}

		if !tst.expectErr && err != nil {
			t.Errorf("ParseOptions(%s) got unexpected error: %v", tst.param, err)
			continue
		}

		if !reflect.DeepEqual(opts, tst.expectedOpts) {
			t.Errorf("ParseOptions(%s) = %v, expected %v", tst.param, opts, tst.expectedOpts)
			continue
		}
	}
}
