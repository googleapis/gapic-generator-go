// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package snippets

//go:generate protoc --go_out=. --go_opt=paths=source_relative metadata/metadata.proto

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/googleapis/gapic-generator-go/internal/snippets/metadata"
	"google.golang.org/protobuf/encoding/protojson"
)

type ApiInfo struct {
	// ProtoPkg is the proto namespace for the API package.
	ProtoPkg string
	// LibPkg is the gapic import path.
	LibPkg string
	// ProtoServices is a map of gapic client short names to service structs.
	ProtoServices map[string]*Service
	// Version is the Go module version for the gapic client.
	Version string
	// ShortName for the service.
	ShortName string
}

// RegionTags gets the region tags keyed by client name and method name.
func (ai *ApiInfo) RegionTags() map[string]map[string]string {
	regionTags := map[string]map[string]string{}
	for svcName, svc := range ai.ProtoServices {
		regionTags[svcName] = map[string]string{}
		for mName, m := range svc.Methods {
			regionTags[svcName][mName] = m.RegionTag
		}
	}
	return regionTags
}

// RegionTags gets the region tags keyed by client name and method name.
func (ai *ApiInfo) toSnippetMetadata() *metadata.Index {
	index := &metadata.Index{
		ClientLibrary: &metadata.ClientLibrary{
			Name:     ai.LibPkg,
			Version:  ai.Version,
			Language: metadata.Language_GO,
			Apis: []*metadata.Api{
				{
					Id:      ai.ProtoPkg,
					Version: ai.protoVersion(),
				},
			},
		},
	}

	// Sorting keys to stabilize output
	var svcKeys []string
	for k := range ai.ProtoServices {
		svcKeys = append(svcKeys, k)
	}
	sort.StringSlice(svcKeys).Sort()
	for _, clientShortName := range svcKeys {
		service := ai.ProtoServices[clientShortName]
		var methodKeys []string
		for k := range service.Methods {
			methodKeys = append(methodKeys, k)
		}
		sort.StringSlice(methodKeys).Sort()
		for _, methodShortName := range methodKeys {
			method := service.Methods[methodShortName]
			snip := &metadata.Snippet{
				RegionTag:   method.RegionTag,
				Title:       fmt.Sprintf("%s %s Sample", ai.ShortName, methodShortName),
				Description: strings.TrimSpace(method.Doc),
				File:        fmt.Sprintf("%s/%s/main.go", clientShortName, methodShortName),
				Language:    metadata.Language_GO,
				Canonical:   false,
				Origin:      *metadata.Snippet_API_DEFINITION.Enum(),
				ClientMethod: &metadata.ClientMethod{
					ShortName:  methodShortName,
					FullName:   fmt.Sprintf("%s.%s.%s", ai.ProtoPkg, clientShortName, methodShortName),
					Async:      false,
					ResultType: method.Result,
					Client: &metadata.ServiceClient{
						ShortName: clientShortName,
						FullName:  fmt.Sprintf("%s.%s", ai.ProtoPkg, clientShortName),
					},
					Method: &metadata.Method{
						ShortName: methodShortName,
						FullName:  fmt.Sprintf("%s.%s.%s", ai.ProtoPkg, service.ProtoName, methodShortName),
						Service: &metadata.Service{
							ShortName: service.ProtoName,
							FullName:  fmt.Sprintf("%s.%s", ai.ProtoPkg, service.ProtoName),
						},
					},
				},
			}
			segment := &metadata.Snippet_Segment{
				Start: int32(method.RegionTagStart + 1),
				End:   int32(method.RegionTagEnd - 1),
				Type:  metadata.Snippet_Segment_FULL,
			}
			snip.Segments = append(snip.Segments, segment)
			for _, param := range method.Params {
				methParam := &metadata.ClientMethod_Parameter{
					Type: param.PType,
					Name: param.Name,
				}
				snip.ClientMethod.Parameters = append(snip.ClientMethod.Parameters, methParam)
			}
			index.Snippets = append(index.Snippets, snip)
		}
	}
	return index
}

func (ai *ApiInfo) protoVersion() string {
	ss := strings.Split(ai.ProtoPkg, ".")
	return ss[len(ss)-1]
}

var spaceSanitizerRegex = regexp.MustCompile(`:\s*`)

func (ai *ApiInfo) ToMetadataJSON() ([]byte, error) {
	m := ai.toSnippetMetadata()
	b, err := protojson.MarshalOptions{Multiline: true}.Marshal(m)
	if err != nil {
		return nil, err
	}
	// Hack to standardize output from protojson which is currently non-deterministic
	// with spacing after json keys.
	return spaceSanitizerRegex.ReplaceAll(b, []byte(": ")), nil
}

// service associates a proto service from gapic metadata with gapic client and its methods
type Service struct {
	// protoName is the name of the proto service.
	ProtoName string
	// methods is a map of gapic method short names to method structs.
	Methods map[string]*Method
}

// Method associates elements of gapic client methods (docs, params and return types)
// with snippet file details such as the region tag string and line numbers.
type Method struct {
	// Doc is the documention for the methods.
	Doc string
	// RegionTag is the region tag that will be used for the generated snippet.
	RegionTag string
	// RegionTagStart is the line number of the START region tag in the snippet file.
	RegionTagStart int
	// RegionTagEnd is the line number of the END region tag in the snippet file.
	RegionTagEnd int
	// Params are the input parameters for the gapic method.
	Params []*param
	// Result is the return value for the method.
	Result string
}

type param struct {
	// Name of the parameter.
	Name string
	// PType is the Go type for the parameter.
	PType string
}
