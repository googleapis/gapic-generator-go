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
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

type transport int

const (
	grpc transport = iota
	rest
)

var errInvalidParam = errors.New("need parameter in format: go-gapic-package=client/import/path;packageName")

type options struct {
	pkgPath           string
	pkgName           string
	outDir            string
	relLvl            string
	modulePrefix      string
	grpcConfPath      string
	serviceConfigPath string
	transports        []transport
	metadata          bool
	diregapic         bool
	restNumericEnum   bool
	pkgOverrides      map[string]string
	omitSnippets      bool
}

// parseOptions takes a string and parses it into a struct defining
// customizations on the target gapic surface.
// Options are comma-separated key/value pairs which are in turn delimited with '='.
// Valid options include:
// * go-gapic-package (package and module naming info)
// * api-service-config (filepath)
// * grpc-service-config (filepath)
// * module (name)
// * Mfile=import (e.g. Mgoogle/storage/v2/storage.proto=cloud.google.com/go/storage/internal/apiv2/stubs)
// * release-level (one of 'alpha', 'beta', 'deprecated', or empty)
// * transport ('+' separated list of transport backends to generate)
// * metadata (enable GAPIC metadata generation)
// * omit-snippets (skip example code snippets generation to the `internal/generated/snippets` path)
// The only required option is 'go-gapic-package'.
//
// Valid parameter example:
// go-gapic-package=path/to/out;pkg,module=path,transport=rest+grpc,api-service-config=api_v1.yaml,release-level=alpha
//
// It returns a pointer to a populated options if no errors were encountered while parsing.
// If errors were encountered, it returns a nil pointer and the first error.
func parseOptions(parameter *string) (*options, error) {
	opts := options{}

	if parameter == nil {
		return nil, fmt.Errorf("empty options parameter")
	}

	// parse plugin params, ignoring unknown values
	for _, s := range strings.Split(*parameter, ",") {
		// skip empty --go_gapic_opt flags
		if s == "" {
			continue
		}

		// Check for boolean flags.
		switch s {
		case "metadata":
			opts.metadata = true
			continue
		case "diregapic":
			opts.diregapic = true
			continue
		case "rest-numeric-enums":
			opts.restNumericEnum = true
			continue
		case "omit-snippets":
			opts.omitSnippets = true
			continue
		}

		e := strings.IndexByte(s, '=')
		if e < 0 {
			return nil, fmt.Errorf("invalid plugin option format, must be key=value: %q", s)
		}

		key, val := s[:e], s[e+1:]
		if val == "" {
			return nil, fmt.Errorf("invalid plugin option value, missing value in key=value: %q", s)
		}

		switch key {
		case "go-gapic-package":
			p := strings.IndexByte(s, ';')

			if p < 0 {
				return nil, errInvalidParam
			}

			opts.pkgPath = s[e+1 : p]
			opts.pkgName = s[p+1:]
			opts.outDir = filepath.FromSlash(opts.pkgPath)
		case "gapic-service-config":
			// Deprecated: this option is deprecated and will be removed in a
			// later release.
			fallthrough
		case "api-service-config":
			opts.serviceConfigPath = val
		case "grpc-service-config":
			opts.grpcConfPath = val
		case "module":
			opts.modulePrefix = val
		case "release-level":
			opts.relLvl = strings.ToLower(val)
		case "transport":
			// Prevent duplicates
			transports := map[transport]bool{}
			for _, t := range strings.Split(val, "+") {
				switch t {
				case "grpc":
					transports[grpc] = true
				case "rest":
					transports[rest] = true
				default:
					return nil, fmt.Errorf("invalid transport option: %q", t)
				}
			}
			for t := range transports {
				opts.transports = append(opts.transports, t)
			}
			sort.Slice(opts.transports, func(i, j int) bool {
				return opts.transports[i] < opts.transports[j]
			})
		default:
			// go_package override for the protobuf/grpc stubs.
			// Mgoogle/storage/v2/storage.proto=cloud.google.com/go/storage/internal/apiv2/stubs
			if key[0] == 'M' {
				file := key[1:]
				if opts.pkgOverrides == nil {
					opts.pkgOverrides = make(map[string]string)
				}
				opts.pkgOverrides[file] = val
			}
		}
	}

	if opts.pkgPath == "" || opts.pkgName == "" || opts.outDir == "" {
		return nil, errInvalidParam
	}

	if opts.modulePrefix != "" {
		if !strings.HasPrefix(opts.outDir, opts.modulePrefix) {
			return nil, fmt.Errorf("go-gapic-package %q does not match prefix %q", opts.outDir, opts.modulePrefix)
		}
		opts.outDir = strings.TrimPrefix(opts.outDir, opts.modulePrefix+"/")
	}

	// Default is just grpc for now.
	if opts.transports == nil {
		opts.transports = []transport{grpc}
	}

	return &opts, nil
}

// Utility function for stringifying the Transport enum
func (t transport) String() string {
	switch t {
	case grpc:
		return "grpc"
	case rest:
		return "rest"
	default:
		// Add new transport variants as need be.
		return fmt.Sprintf("%d", int(t))
	}
}
