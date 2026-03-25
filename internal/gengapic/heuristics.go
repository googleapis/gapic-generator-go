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
	"regexp"
	"strings"

	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

var literalVarRegex = regexp.MustCompile(`([^/]+)/{[^}]+\}`)

// BuildHeuristicVocabulary builds a map of valid resource tokens
// based on the last literal before a variable in CRUD-like methods.
func BuildHeuristicVocabulary(methods []*descriptorpb.MethodDescriptorProto) map[string]bool {
	tokens := make(map[string]bool)

	// Add standard infrastructure tokens
	tokens["projects"] = true
	tokens["locations"] = true
	tokens["folders"] = true
	tokens["organizations"] = true
	tokens["billingAccounts"] = true

	discoveryExactVerbs := map[string]struct{}{
		"get":            {},
		"list":           {},
		"aggregatedlist": {},
		"create":         {},
		"update":         {},
		"delete":         {},
		"patch":          {},
		"insert":         {},
	}
	discoverySuffixes := []string{
		".get", ".list", ".create", ".update", ".delete", ".patch", ".insert",
	}

	crudPrefixes := []string{
		"get", "list", "create", "update", "delete", "patch", "insert",
	}

	for _, m := range methods {
		nameLower := strings.ToLower(m.GetName())

		var isCRUDPrefix bool
		for _, prefix := range crudPrefixes {
			if strings.HasPrefix(nameLower, prefix) {
				isCRUDPrefix = true
				break
			}
		}

		_, isDiscoveryExact := discoveryExactVerbs[nameLower]

		var isDiscoverySuffix bool
		for _, suffix := range discoverySuffixes {
			if strings.HasSuffix(nameLower, suffix) {
				isDiscoverySuffix = true
				break
			}
		}

		if !isCRUDPrefix && !isDiscoveryExact && !isDiscoverySuffix {
			continue
		}

		if m.GetOptions() == nil {
			continue
		}

		eHTTP := proto.GetExtension(m.GetOptions(), annotations.E_Http)
		if eHTTP == nil {
			continue
		}

		h, ok := eHTTP.(*annotations.HttpRule)
		if !ok || h == nil {
			continue
		}

		patterns := getHttpPatterns(h)
		for _, pattern := range patterns {
			matches := literalVarRegex.FindAllStringSubmatch(pattern, -1)
			for _, match := range matches {
				if len(match) < 2 {
					continue
				}
				literal := match[1]
				if isVersionString(literal) {
					continue
				}
				tokens[literal] = true
			}
		}
	}

	return tokens
}

func getHttpPatterns(h *annotations.HttpRule) []string {
	var patterns []string

	extract := func(pattern string) {
		if pattern != "" {
			patterns = append(patterns, pattern)
		}
	}

	switch p := h.GetPattern().(type) {
	case *annotations.HttpRule_Get:
		extract(p.Get)
	case *annotations.HttpRule_Put:
		extract(p.Put)
	case *annotations.HttpRule_Post:
		extract(p.Post)
	case *annotations.HttpRule_Delete:
		extract(p.Delete)
	case *annotations.HttpRule_Patch:
		extract(p.Patch)
	}

	for _, rule := range h.GetAdditionalBindings() {
		patterns = append(patterns, getHttpPatterns(rule)...)
	}

	return patterns
}

func isVersionString(s string) bool {
	if !strings.HasPrefix(s, "v") {
		return false
	}
	if len(s) < 2 {
		return false
	}
	if s[1] >= '0' && s[1] <= '9' {
		return true
	}
	return false
}
