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

	"github.com/googleapis/gapic-generator-go/internal/license"
	"github.com/googleapis/gapic-generator-go/internal/snippets/metadata"
	"google.golang.org/protobuf/encoding/protojson"
)

// VersionPlaceholder is the string value $VERSION, intended to be replaced with
// the actual module version by a generator post-processing script.
var VersionPlaceholder = "$VERSION"

// headerLen is the length of the Apache license header including trailing newlines.
var headerLen = len(strings.Split(license.Apache, "\n"))

var spaceSanitizerRegex = regexp.MustCompile(`:\s*`)

var ctxParam = &param{
	name:  "ctx",
	pType: "context.Context",
}

var optsParam = &param{
	name:  "opts",
	pType: "...gax.CallOption",
}

// SnippetMetadata is a model for capturing snippet details and writing them to
// a snippet_metadata.*.json file.
type SnippetMetadata struct {
	// protoPkg is the proto namespace for the API package.
	protoPkg string
	// libPkg is the Go import path for the GAPIC client.
	libPkg string
	// protoServices is a map of gapic service short names to service structs.
	protoServices map[string]*service
	// apiVersion is the gapic service version. (e.g. "v1", "v1beta1")
	apiVersion string
	// shortName the first element of API service config DNS-like Name field. (e.g. "bigquerymigration" from "bigquerymigration.googleapis.com")
	shortName string
}

// NewMetadata initializes the model that will collect snippet metadata, from:
// protoPkg - dot-separated, without final type name element (e.g. "google.cloud.bigquery.migration.v2")
// libPkg - the Go import path for the GAPIC client, per libraryPackage in gapic_metadata.json (e.g. "cloud.google.com/go/bigquery/migration/apiv2")
// serviceConfigName - the API service config DNS-like Name field. (e.g. "bigquerymigration.googleapis.com")
func NewMetadata(protoPkg, libPkg, serviceConfigName string) (*SnippetMetadata, error) {
	protoParts := strings.Split(protoPkg, ".")
	apiVersion := protoParts[len(protoParts)-1]
	shortName := strings.Split(serviceConfigName, ".")[0]
	if shortName == "" {
		return nil, fmt.Errorf("snippets: api-service-config is required and must contain Name for %s", protoPkg)
	}
	return &SnippetMetadata{
		protoPkg:      protoPkg,
		libPkg:        libPkg,
		shortName:     shortName,
		apiVersion:    apiVersion,
		protoServices: make(map[string]*service),
	}, nil
}

// AddService uses the service short name (e.g. "AutoscalingPolicyService") identifier
// to add a service entry.
func (sm *SnippetMetadata) AddService(serviceName string) {
	s := &service{
		protoName: serviceName,
		methods:   make(map[string]*method),
	}
	sm.protoServices[serviceName] = s
}

// AddMethod uses the service short name (e.g. "AutoscalingPolicyService") and method name
// identifiers to add an incomplete method entry that will be updated via UpdateMethodDoc
// and UpdateMethodResult.
func (sm *SnippetMetadata) AddMethod(serviceName, methodName string, regionTagEnd int) {
	m := &method{
		regionTag:      sm.RegionTag(serviceName, methodName),
		regionTagStart: headerLen,
		regionTagEnd:   regionTagEnd,
	}
	sm.protoServices[serviceName].methods[methodName] = m
}

// UpdateMethodDoc uses service short name (e.g. "AutoscalingPolicyService") and
// and method name identifiers to add a doc method comment.
func (sm *SnippetMetadata) UpdateMethodDoc(serviceName, methodName, doc string) {
	m := sm.protoServices[serviceName].methods[methodName]
	m.doc = doc
}

// UpdateMethodResult uses service short name (e.g. "AutoscalingPolicyService") and
// and method name identifiers to add a method result type.
func (sm *SnippetMetadata) UpdateMethodResult(serviceName, methodName, result string) {
	m := sm.protoServices[serviceName].methods[methodName]
	m.result = result
}

// AddParams adds a slice of 3 params to the method: ctx context.Context, req <requestType>, opts ...gax.CallOption,
// ctx and opts params are hardcoded since these are currently the same in all client wrapper methods.
// The req param will be omitted if empty requestType is given.
func (sm *SnippetMetadata) AddParams(serviceName, methodName, requestType string) {
	m := sm.protoServices[serviceName].methods[methodName]
	m.params = []*param{ctxParam}
	if requestType != "" {
		m.params = append(m.params,
			&param{
				name:  "req",
				pType: requestType,
			})
	}
	m.params = append(m.params, optsParam)
}

// RegionTag generates a snippet region tag from shortName, apiVersion, and the given full serviceName and method name.
func (sm *SnippetMetadata) RegionTag(serviceName, methodName string) string {
	return fmt.Sprintf("%s_%s_generated_%s_%s_sync", sm.shortName, sm.apiVersion, serviceName, methodName)
}

// ToMetadataJSON marshals the completed SnippetMetadata to a []byte containing
// the protojson output.
func (sm *SnippetMetadata) ToMetadataJSON() ([]byte, error) {
	m := sm.ToMetadataIndex()
	b, err := protojson.MarshalOptions{Multiline: true}.Marshal(m)
	if err != nil {
		return nil, err
	}
	// Hack to standardize output from protojson which is currently non-deterministic
	// with spacing after json keys.
	return spaceSanitizerRegex.ReplaceAll(b, []byte(": ")), nil
}

// ToMetadataIndex creates a metadata.Index from the SnippetMetadata.
func (sm *SnippetMetadata) ToMetadataIndex() *metadata.Index {
	index := &metadata.Index{
		ClientLibrary: &metadata.ClientLibrary{
			Name:     sm.libPkg,
			Version:  VersionPlaceholder,
			Language: metadata.Language_GO,
			Apis: []*metadata.Api{
				{
					Id:      sm.protoPkg,
					Version: sm.protoVersion(),
				},
			},
		},
	}

	// Sort keys to stabilize output
	var svcKeys []string
	for k := range sm.protoServices {
		svcKeys = append(svcKeys, k)
	}
	sort.StringSlice(svcKeys).Sort()
	for _, serviceShortName := range svcKeys {
		clientShortName := serviceShortName + "Client"
		service := sm.protoServices[serviceShortName]
		var methodKeys []string
		for k := range service.methods {
			methodKeys = append(methodKeys, k)
		}
		sort.StringSlice(methodKeys).Sort()
		for _, methodShortName := range methodKeys {
			method := service.methods[methodShortName]
			snp := &metadata.Snippet{
				RegionTag:   method.regionTag,
				Title:       fmt.Sprintf("%s %s Sample", sm.shortName, methodShortName),
				Description: strings.TrimSpace(method.doc),
				File:        fmt.Sprintf("%s/%s/main.go", clientShortName, methodShortName),
				Language:    metadata.Language_GO,
				Canonical:   false,
				Origin:      *metadata.Snippet_API_DEFINITION.Enum(),
				ClientMethod: &metadata.ClientMethod{
					ShortName:  methodShortName,
					FullName:   fmt.Sprintf("%s.%s.%s", sm.protoPkg, clientShortName, methodShortName),
					Async:      false,
					ResultType: method.result,
					Client: &metadata.ServiceClient{
						ShortName: clientShortName,
						FullName:  fmt.Sprintf("%s.%s", sm.protoPkg, clientShortName),
					},
					Method: &metadata.Method{
						ShortName: methodShortName,
						FullName:  fmt.Sprintf("%s.%s.%s", sm.protoPkg, service.protoName, methodShortName),
						Service: &metadata.Service{
							ShortName: service.protoName,
							FullName:  fmt.Sprintf("%s.%s", sm.protoPkg, service.protoName),
						},
					},
				},
			}
			segment := &metadata.Snippet_Segment{
				Start: int32(method.regionTagStart + 1),
				End:   int32(method.regionTagEnd - 1),
				Type:  metadata.Snippet_Segment_FULL,
			}
			snp.Segments = append(snp.Segments, segment)
			for _, param := range method.params {
				methParam := &metadata.ClientMethod_Parameter{
					Type: param.pType,
					Name: param.name,
				}
				snp.ClientMethod.Parameters = append(snp.ClientMethod.Parameters, methParam)
			}
			index.Snippets = append(index.Snippets, snp)
		}
	}
	return index
}

func (sm *SnippetMetadata) protoVersion() string {
	ss := strings.Split(sm.protoPkg, ".")
	return ss[len(ss)-1]
}

// service associates a proto service from gapic metadata with gapic client and its methods.
type service struct {
	// protoName is the service short name.
	protoName string
	// methods is a map of gapic method short names to method structs.
	methods map[string]*method
}

// method associates elements of gapic client methods (docs, params and return types)
// with snippet file details such as the region tag string and line numbers.
type method struct {
	// doc is the documention for the methods.
	doc string
	// regionTag is the region tag that will be used for the generated snippet.
	regionTag string
	// regionTagStart is the number of the line AFTER the START region tag in the snippet file.
	regionTagStart int
	// regionTagEnd is the line number of the END region tag in the snippet file.
	regionTagEnd int
	// params are the input parameters for the gapic method.
	params []*param
	// result is the return value for the method.
	result string
}

// param contains the details of a method parameter.
type param struct {
	// name of the parameter.
	name string
	// pType is the Go type for the parameter.
	pType string
}
