// Copyright 2018 Google LLC
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

package gencli

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
)

const (
	// ExpectedParams is the number of expected plugin parameters
	ExpectedParams = 2

	// ShortDescMax is the maximum length accepted for
	// the Short usage docs
	ShortDescMax = 50
)

func parseParameters(params *string) (rootDir string, gapicDir string, err error) {
	if params == nil {
		err = fmt.Errorf("Missing required parameters. See usage")
		return
	}

	split := strings.Split(*params, ",")
	if len(split) != ExpectedParams {
		err = fmt.Errorf("Improper number of parameters. Got %d, require %d. See usage", len(split), ExpectedParams)
		return
	}

	for _, str := range split {
		sepNdx := strings.Index(str, ":")
		if sepNdx == -1 {
			err = fmt.Errorf("Unknown parameter: %s", str)
			return
		}

		switch str[:sepNdx] {
		case "gapic":
			gapicDir = str[sepNdx+1:]
		case "root":
			rootDir = str[sepNdx+1:]
		default:
			err = fmt.Errorf("Unknown parameter: %s", str)
		}
	}

	return
}

func toShortUsage(cmt string) string {
	if len(cmt) > ShortDescMax {
		sep := strings.LastIndex(cmt[:ShortDescMax], " ")
		if sep == -1 {
			sep = ShortDescMax
		}
		cmt = cmt[:sep] + "..."
	}

	return cmt
}

func sanitizeComment(cmt string) string {
	cmt = strings.Replace(cmt, "\\", `\\`, -1)
	cmt = strings.Replace(cmt, "\n", " ", -1)
	cmt = strings.TrimSpace(cmt)
	return cmt
}

func strContains(a []string, s string) bool {
	for _, as := range a {
		if as == s {
			return true
		}
	}
	return false
}

func copyImports(from, to map[string]*pbinfo.ImportSpec) {
	for _, val := range from {
		putImport(to, val)
	}
}

func putImport(imports map[string]*pbinfo.ImportSpec, pkg *pbinfo.ImportSpec) {
	if _, ok := imports[pkg.Name]; !ok {
		imports[pkg.Name] = pkg
	}
}

func parseMessageName(field *descriptor.FieldDescriptorProto, msg *descriptor.DescriptorProto) (name string) {
	t := field.GetTypeName()
	last := strings.LastIndex(t, ".")
	name = t[last+1:]

	// check if it is a nested type
	if strings.Contains(t, msg.GetName()) {
		pre := t[:last]
		parent := pre[strings.LastIndex(pre, ".")+1:]
		name = parent + "_" + name
	}

	return
}

func title(name string) string {
	split := strings.Split(name, "_")
	for i, s := range split {
		split[i] = strings.Title(s)
	}

	return strings.Join(split, "")
}
