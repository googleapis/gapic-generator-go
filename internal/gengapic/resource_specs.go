// Copyright 2026 Google LLC
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
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

var paramRegexp = regexp.MustCompile(`\{([a-zA-Z0-9_\-]+)(?:=[^\}]+)?\}`)

// resourceSpec holds metadata extracted from google.api.resource annotations.
type resourceSpec struct {
	Type        string
	Collection  string
	Pattern     string
	AltPatterns []string
	Params      []string
}

func (g *generator) genResourceSpecsFile(files []*descriptorpb.FileDescriptorProto) error {
	specs := g.collectResourceSpecs(files)
	if len(specs) == 0 {
		return nil
	}

	g.reset()
	g.genResourceSpecs(specs)
	g.commit(filepath.Join(g.cfg.outDir, "resource_specs.go"), g.cfg.pkgName)
	return nil
}

func (g *generator) collectResourceSpecs(genReqFiles []*descriptorpb.FileDescriptorProto) map[string]resourceSpec {
	specs := make(map[string]resourceSpec)

	for _, file := range genReqFiles {
		// 1. Check file-level resource_definition annotations
		if opts := file.GetOptions(); opts != nil && proto.HasExtension(opts, annotations.E_ResourceDefinition) {
			if resDefs, ok := proto.GetExtension(opts, annotations.E_ResourceDefinition).([]*annotations.ResourceDescriptor); ok {
				for _, rd := range resDefs {
					addResourceDescriptor(specs, rd)
				}
			}
		}

		// 2. Check message-level resource annotations
		for _, msg := range file.GetMessageType() {
			collectMessageResourceSpecs(specs, msg)
		}
	}

	return specs
}

func collectMessageResourceSpecs(specs map[string]resourceSpec, msg *descriptorpb.DescriptorProto) {
	if opts := msg.GetOptions(); opts != nil && proto.HasExtension(opts, annotations.E_Resource) {
		if rd, ok := proto.GetExtension(opts, annotations.E_Resource).(*annotations.ResourceDescriptor); ok {
			addResourceDescriptor(specs, rd)
		}
	}

	for _, nested := range msg.GetNestedType() {
		collectMessageResourceSpecs(specs, nested)
	}
}

func addResourceDescriptor(specs map[string]resourceSpec, rd *annotations.ResourceDescriptor) {
	if rd == nil || rd.GetType() == "" || len(rd.GetPattern()) == 0 {
		return
	}

	resType := rd.GetType()
	patterns := rd.GetPattern()
	primaryPattern := patterns[0]
	var altPatterns []string
	if len(patterns) > 1 {
		altPatterns = patterns[1:]
	}

	collection := extractCollectionName(resType, primaryPattern)
	params := extractParams(primaryPattern)

	specs[resType] = resourceSpec{
		Type:        resType,
		Collection:  collection,
		Pattern:     primaryPattern,
		AltPatterns: altPatterns,
		Params:      params,
	}
}

func extractCollectionName(resType, pattern string) string {
	parts := strings.Split(pattern, "/")
	if len(parts) >= 2 {
		last := parts[len(parts)-1]
		if strings.HasPrefix(last, "{") && len(parts) >= 2 {
			return parts[len(parts)-2]
		}
	}
	if slash := strings.LastIndex(resType, "/"); slash >= 0 {
		return strings.ToLower(resType[slash+1:])
	}
	return resType
}

func extractParams(pattern string) []string {
	matches := paramRegexp.FindAllStringSubmatch(pattern, -1)
	var params []string
	seen := make(map[string]bool)
	for _, m := range matches {
		if len(m) >= 2 && !seen[m[1]] {
			seen[m[1]] = true
			params = append(params, m[1])
		}
	}
	return params
}

func (g *generator) genResourceSpecs(specs map[string]resourceSpec) {
	p := g.pt.Printf

	p("// ResourceSpec defines metadata for a Google API resource.")
	p("type ResourceSpec struct {")
	p("	Type        string")
	p("	Collection  string")
	p("	Pattern     string")
	p("	AltPatterns []string")
	p("	Params      []string")
	p("}")
	p("")

	var types []string
	for t := range specs {
		types = append(types, t)
	}
	sort.Strings(types)

	p("// ResourceSpecs maps resource types to their declarative specifications.")
	p("var ResourceSpecs = map[string]ResourceSpec{")
	for _, t := range types {
		spec := specs[t]
		p("	%q: {", spec.Type)
		p("		Type:        %q,", spec.Type)
		p("		Collection:  %q,", spec.Collection)
		p("		Pattern:     %q,", spec.Pattern)
		if len(spec.AltPatterns) > 0 {
			p("		AltPatterns: []string{")
			for _, alt := range spec.AltPatterns {
				p("			%q,", alt)
			}
			p("		},")
		}
		p("		Params: []string{")
		for _, param := range spec.Params {
			p("			%q,", param)
		}
		p("		},")
		p("	},")
	}
	p("}")
}
