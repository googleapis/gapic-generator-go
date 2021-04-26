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
	"fmt"
	"path/filepath"
	"strings"

	"github.com/googleapis/gapic-generator-go/internal/errors"
)

type transport int

const (
	grpc transport = iota
	rest
)

const paramError = "need parameter in format: go-gapic-package=client/import/path;packageName"

type options struct {
	pkgPath           string
	pkgName           string
	outDir            string
	relLvl            string
	modulePrefix      string
	grpcConfPath      string
	serviceConfigPath string
	sampleOnly        bool
	transports        []transport
	metadata          bool
}

// parseOptions takes a string and parses it into a struct defining
// customizations on the target gapic surface.
// Options are comma-separated key/value pairs which are in turn delimited with '='.
// Valid options include:
// * go-gapic-package (package and module naming info)
// * sample-only (only checked for presence)
// * api-service-config (filepath)
// * grpc-service-config (filepath)
// * module (name)
// * release-level (one of 'alpha', 'beta', or empty)
// * transport ('+' separated list of transport backends to generate)
// * metadata (enable GAPIC metadata generation)
// The only required option is 'go-gapic-package'.
//
// Valid parameter example:
// go-gapic-package=path/to/out;pkg,module=path,transport=rest+grpc,api-service-config=api_v1.yaml,release-level=alpha
//
// It returns a pointer to a populated options if no errors were encountered while parsing.
// If errors were encountered, it returns a nil pointer and the first error.
func parseOptions(parameter *string) (*options, error) {
	opts := options{sampleOnly: false}

	if parameter == nil {
		return nil, errors.E(nil, "empty options parameter")
	}

	// parse plugin params, ignoring unknown values
	for _, s := range strings.Split(*parameter, ",") {
		// skip empty --go_gapic_opt flags
		if s == "" {
			continue
		}
		// check for the boolean flag, sample-only, that disables client generation
		if s == "sample-only" {
			return &options{sampleOnly: true}, nil
		}
		if s == "metadata" {
			opts.metadata = true
			continue
		}

		e := strings.IndexByte(s, '=')
		if e < 0 {
			return nil, errors.E(nil, "invalid plugin option format, must be key=value: %s", s)
		}

		key, val := s[:e], s[e+1:]
		if val == "" {
			return nil, errors.E(nil, "invalid plugin option value, missing value in key=value: %s", s)
		}

		switch key {
		case "go-gapic-package":
			p := strings.IndexByte(s, ';')

			if p < 0 {
				return nil, errors.E(nil, paramError)
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
			for _, t := range strings.Split(val, "+") {
				switch t {
				case "grpc":
					opts.transports = append(opts.transports, grpc)
				case "rest":
					opts.transports = append(opts.transports, rest)
				default:
					return nil, errors.E(nil, "invalid transport option: %s", t)
				}
			}
		}
	}

	if opts.pkgPath == "" || opts.pkgName == "" || opts.outDir == "" {
		return nil, errors.E(nil, paramError)
	}

	if opts.modulePrefix != "" {
		if !strings.HasPrefix(opts.outDir, opts.modulePrefix) {
			return nil, errors.E(nil, "go-gapic-package %q does not match prefix %q", opts.outDir, opts.modulePrefix)
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
