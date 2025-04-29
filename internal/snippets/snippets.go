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
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
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
	// pkgName is the short package name after the semi-colon in go-gapic-package.
	pkgName string
}

// NewMetadata initializes the model that will collect snippet metadata, from:
// protoPkg - dot-separated, without final type name element (e.g. "google.cloud.bigquery.migration.v2")
// libPkg - the Go import path for the GAPIC client, per libraryPackage in gapic_metadata.json (e.g. "cloud.google.com/go/bigquery/migration/apiv2")
// pkgName - stored in g.opts.pkgName, used as an argument to pbinfo.ReduceServName
func NewMetadata(protoPkg, libPkg, pkgName string) *SnippetMetadata {
	lastDot := strings.LastIndex(protoPkg, ".")
	apiVersion := protoPkg[lastDot+1:]
	return &SnippetMetadata{
		protoPkg:      protoPkg,
		libPkg:        libPkg,
		apiVersion:    apiVersion,
		protoServices: make(map[string]*service),
		pkgName:       pkgName,
	}
}

// AddService creates a service entry from:
// servName - the service short name (e.g. "AutoscalingPolicyService") identifier
// defaultHost - the DNS hostname for the service, available from annotations.E_DefaultHost. (e.g. "bigquerymigration.googleapis.com")
func (sm *SnippetMetadata) AddService(servName, defaultHost string) {
	shortName := strings.Split(defaultHost, ".")[0]
	s := &service{
		protoName: servName,
		methods:   make(map[string]*method),
		shortName: shortName,
	}
	sm.protoServices[servName] = s
}

// AddMethod uses the service short name (e.g. "AutoscalingPolicyService") and method name
// to add an incomplete method entry that will be updated via UpdateMethodDoc and UpdateMethodResult.
// parentProtoPkg and parentName are the original proto namespace and service for the method.
// (In mixin methods, these are different from the protoPkg and service into which it has been mixed.)
func (sm *SnippetMetadata) AddMethod(servName, methodName, parentProtoPkg, parentName string, regionTagEnd int) {
	m := &method{
		regionTag:      sm.RegionTag(servName, methodName),
		regionTagStart: headerLen,
		regionTagEnd:   regionTagEnd,
		parentProtoPkg: parentProtoPkg,
		parentName:     parentName,
	}
	sm.protoServices[servName].methods[methodName] = m
}

// UpdateMethodDoc uses service short name (e.g. "AutoscalingPolicyService") and
// and method name identifiers to add a doc method comment.
func (sm *SnippetMetadata) UpdateMethodDoc(servName, methodName, doc string) {
	m := sm.protoServices[servName].methods[methodName]
	lines := strings.Split(doc, "\n")
	var b strings.Builder
	for _, l := range lines {
		b.WriteString(strings.TrimSpace(l) + "\n")
	}
	m.doc = b.String()
}

// UpdateMethodResult uses service short name (e.g. "AutoscalingPolicyService") and
// and method name identifiers to add a method result type.
func (sm *SnippetMetadata) UpdateMethodResult(servName, methodName, result string) {
	m := sm.protoServices[servName].methods[methodName]
	m.result = result
}

// AddParams adds a slice of 3 params to the method: ctx context.Context, req <requestType>, opts ...gax.CallOption,
// ctx and opts params are hardcoded since these are currently the same in all client wrapper methods.
// The req param will be omitted if empty requestType is given.
func (sm *SnippetMetadata) AddParams(servName, methodName, requestType string) {
	m := sm.protoServices[servName].methods[methodName]
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

// RegionTag generates a snippet region tag from service.shortName(defaultHost),
// apiVersion, and the given full servName and method name.
func (sm *SnippetMetadata) RegionTag(servName, methodName string) string {
	s := sm.protoServices[servName]
	return fmt.Sprintf("%s_%s_generated_%s_%s_sync", s.shortName, sm.apiVersion, servName, methodName)
}

// ToMetadataJSON marshals the completed SnippetMetadata to a []byte containing
// the protojson output.
func (sm *SnippetMetadata) ToMetadataJSON() ([]byte, error) {
	m := sm.ToMetadataIndex()
	b, err := protojson.MarshalOptions{Multiline: true}.Marshal(m)
	if err != nil {
		return nil, err
	}
	// Hack to standardize output from protojson which is currently
	// non-deterministic with spacing after json keys.
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
		reducedServName := pbinfo.ReduceServName(serviceShortName, sm.pkgName)
		clientShortName := reducedServName + "Client"
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
				Title:       fmt.Sprintf("%s %s Sample", service.shortName, methodShortName),
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
						FullName:  fmt.Sprintf("%s.%s.%s", method.parentProtoPkg, method.parentName, methodShortName),
						Service: &metadata.Service{
							ShortName: method.parentName,
							FullName:  fmt.Sprintf("%s.%s", method.parentProtoPkg, method.parentName),
						},
					},
				},
			}
			segment := &metadata.Snippet_Segment{
				// The line where this segment begins, inclusive.
				// For the FULL segment, this will be the START region tag line + 1.
				Start: int32(method.regionTagStart + 1),
				// The line where this segment ends, inclusive.
				// For the FULL segment, this will be the END region tag line - 1.
				End:  int32(method.regionTagEnd - 1),
				Type: metadata.Snippet_Segment_FULL,
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
	// shortName the first element of the default DNS hostname for the service.
	// (e.g. "bigquerymigration" from "bigquerymigration.googleapis.com")
	shortName string
}

// method associates elements of gapic client methods (docs, params and return types)
// with snippet file details such as the region tag string and line numbers.
type method struct {
	// doc is the documentation for the methods.
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
	// parentProtoPkg is the original proto namespace for the method.
	// In mixin methods, this namespace is different from the protoPkg namespace into which it has been mixed.
	parentProtoPkg string
	// parentName is the proto name of the method's original parent.
	// In mixin methods, this parent is different from the service into which it has been mixed.
	parentName string
}

// param contains the details of a method parameter.
type param struct {
	// name of the parameter.
	name string
	// pType is the Go type for the parameter.
	pType string
}
