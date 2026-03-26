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

// literalVarRegex matches a path segment literal followed immediately
// by a curly-brace variable ({...}).
var literalVarRegex = regexp.MustCompile(`([^/]+)/{[^}]+\}`)

// HeuristicTarget represents a resolved resource name template pattern and its
// associated variable field names.
type HeuristicTarget struct {
	// Format is the resource template string (e.g. projects/%v/topics/%v).
	Format string
	// FieldNames are the names of variables extracted from the template.
	FieldNames []string
}


// BuildHeuristicVocabulary builds a map of valid resource collections
// by examining standard CRUD-like patterns in routes.
func BuildHeuristicVocabulary(methods []*descriptorpb.MethodDescriptorProto) map[string]bool {
	resourceCollections := make(map[string]bool)

	// Step 1: Seed standard infrastructure resource collections.
	resourceCollections["projects"] = true
	resourceCollections["locations"] = true
	resourceCollections["folders"] = true
	resourceCollections["organizations"] = true
	resourceCollections["billingAccounts"] = true

	// Step 2: Define "CRUD-like" patterns for vocabulary learning.
	// Why do we filter for CRUD? Non-CRUD methods (e.g., `CancelOperation`, `CheckHealth`)
	// often have path literals that are random verbs or actions, not resource collections.
	// If we learned from those, we might pollute our vocabulary with non-resource nouns.
	discoveryExactVerbs := map[string]struct{}{
		"get":  {},
		"list": {},
		// aggregatedlist is used in Compute Engine to list across zones in a project.
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

	// Step 3: Walk methods and learn valid resource collection nouns.
	// By "learn" and "valid", we mean finding standard verbs (Get, List, Create),
	// reading their paths (e.g., `/v1/projects/{project}/topics/{topic}`), and
	// extracting the static literal nouns that sit immediately before a `{variable}`
	// (`projects` and `topics`). We add these to our `resourceCollections` map so we can
	// validate unannotated field patterns later, in IdentifyHeuristicTarget.
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

		// We only learn from standard verbs.
		if !isCRUDPrefix && !isDiscoveryExact && !isDiscoverySuffix {
			continue
		}

		// m.GetOptions() retrieves protobuf options
		// (e.g. `(google.api.http) = { get: "/v1/..." }` attached to the method).
		if m.GetOptions() == nil {
			continue
		}

		// proto.GetExtension extracts standard Protobuf Option Extensions. In this case,
		// we ask for the HTTP mapping rule (annotations.E_Http / google.api.http).
		// If we have unannotated methods, we skip them as they are probably either:
		// - Pure gRPC (no REST binding).
		// - Legacy REST using custom non-standard transports.
		eHTTP := proto.GetExtension(m.GetOptions(), annotations.E_Http)
		if eHTTP == nil {
			continue
		}

		// Cast to HttpRule for use in getHttpPatterns, below.
		h, ok := eHTTP.(*annotations.HttpRule)
		if !ok || h == nil {
			continue
		}

		// Step 4: Extract literals that appear before variables.
		patterns := getHttpPatterns(h)
		for _, pattern := range patterns {
			// Find all instances of a literal right before a variable `{...}`
			// Trace example: `/v1/projects/{project}/topics/{topic}`
			// Matches:
			//   match[0] = "projects/{project}" (with groups [ projects ])
			//   match[1] = "topics/{topic}" (with groups [ topics ])
			matches := literalVarRegex.FindAllStringSubmatch(pattern, -1)
			for _, match := range matches {
				if len(match) < 2 {
					continue
				}
				literal := match[1]

				// We should discard verbs (anything after and including the ':' character)
				// so we do not learn generic collection nouns like "projects:cancel".
				if idx := strings.Index(literal, ":"); idx != -1 {
					literal = literal[:idx]
				}

				if isVersionString(literal) {
					continue
				}
				resourceCollections[literal] = true
			}
		}
	}

	return resourceCollections
}

// getHttpPatterns flattens an HttpRule (and any recursive AdditionalBindings)
// into a flat slice of string patterns.
//
// We use standard method name verbs (Get, List, Create) to filter for where to look.
// But the actual vocabulary we want to learn is the collection nouns (literals)
// in those method's URI paths. An HTTP rule can define a primary path but often
// has secondary endpoints via `additional_bindings` to support legacy paths or
// alternative styles. To learn the full set of valid collection nouns for an API,
// we must process all possible paths.
func getHttpPatterns(h *annotations.HttpRule) []string {
	var patterns []string

	extract := func(pattern string) {
		if pattern != "" {
			patterns = append(patterns, pattern)
		}
	}

	// Step 1: Extract the primary verb's pattern.
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

	// Step 2: Recursively process any additional bindings attached to this rule.
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

// IdentifyHeuristicTarget reconstructs a canonical resource name template for an
// unannotated API endpoint.
//
// It parses a single HTTP rule and attempts to infer its method's resource format
// (and request field parameters) using standard literal/variable pairs. The vocabulary
// param is a map of known valid collection nouns (e.g., "topics", "subscriptions")
// learned by BuildHeuristicVocabulary. It is used to validate whether a literal
// sitting before a variable is a true resource collection, or just an interstitial
// literal (like "global" in `v1/global/firewalls`).
//
// This should only be used if `google.api.resource` or `google.api.resource_reference`
// are not present. Modern services define resources using proto annotations:
// - `google.api.resource` on Message types
// - `google.api.resource_reference` on Field types
//
// For unannotated legacy services (like Compute), these annotations are missing.
// However, if the HTTP path follows standard conventions (like `projects/{project}/topics/{topic}`),
// we can still deduce the vocabulary.
func IdentifyHeuristicTarget(m *descriptorpb.MethodDescriptorProto, h *annotations.HttpRule, vocabulary map[string]bool) (*HeuristicTarget, error) {
	patterns := getHttpPatterns(h)
	if len(patterns) == 0 {
		return nil, nil
	}

	pattern := patterns[0]

	// Step 1: Split the path into atomic segments.
	segments := splitPathSegments(pattern)

	varNameRegex := regexp.MustCompile(`{([^=}\s]+)`)

	// Step 2: Walk BACKWARDS from the end of the HTTP path in order to find the most
	// specific resource (the child) at the end of the URI.
	for i := len(segments) - 1; i >= 0; i-- {
		seg := segments[i]

		// If it's a leading slash segment or cannot be a child variable segment, skip.
		if !strings.Contains(seg, "{") || i == 0 || strings.Contains(segments[i-1], "{") {
			continue
		}

		token := segments[i-1]
		if idx := strings.Index(token, ":"); idx != -1 {
			token = token[:idx]
		}

		// The segment immediately preceding our variable must be a known vocabulary token
		// (e.g., learned from CRUD or standard seeds) or a version string. Discard verbs
		// attached to collections (e.g. `projects/{project}/topics/{topic}:publish`)
		// when checking against the vocabulary.
		if !vocabulary[token] && !isVersionString(token) {
			continue
		}

		firstIndex := i
		if vocabulary[token] {
			// Known collection literal. Walk backward to find the root of the chain.
			firstIndex = i - 1
			for firstIndex > 0 {
				prevSeg := segments[firstIndex-1]
				if !strings.Contains(prevSeg, "{") {
					// Skip standalone literals (e.g., "global").
					if isVersionString(prevSeg) {
						break // Version strings bound the chain
					}
					firstIndex--
					continue
				}

				// If it's a variable segment preceded by a known literal,
				// we pair them up and keep walking.
				if firstIndex < 2 || strings.Contains(segments[firstIndex-2], "{") {
					break // No preceding literal to pair with the variable, broken chain.
				}
				prevLiteralVal := segments[firstIndex-2]
				if vocabulary[prevLiteralVal] {
					firstIndex -= 2 // pair consumed
					continue
				}
				if isVersionString(prevLiteralVal) {
					firstIndex--
				}
				break
			}
		}

		// Step 3: Reject partial matches with unhandled variables on the left.
		//
		// Example: `/v1/users/{user}/topics/{topic}`
		//
		// We have unhandled variables if `{user}` is missing from our vocabulary.
		// If this API package doesn't define a standard `GetUser` or `ListUsers`,
		// we never learn `users` as a valid collection noun.
		//
		// We discard these partial matches to avoid returning an incomplete or
		// disconnected resource name.
		disconnected := false
		for k := 0; k < firstIndex; k++ {
			if strings.Contains(segments[k], "{") {
				disconnected = true
				break
			}
		}
		if disconnected {
			continue
		}

		var fields []string
		var formatSegments []string

		// Step 4: Extract fields and build the format template from the validated sub-path.
		//
		// Example: For `projects/{project}/topics/{topic}`, the template will be
		// `projects/%v/topics/%v` and the field names will be `["project", "topic"]`.
		targetSegments := segments[firstIndex : i+1]
		for _, s := range targetSegments {
			if strings.Contains(s, "{") {
				match := varNameRegex.FindStringSubmatch(s)
				if len(match) > 1 {
					fields = append(fields, match[1])
				}
				formatSegments = append(formatSegments, "%v")
			} else {
				token := s
				if idx := strings.Index(token, ":"); idx != -1 {
					token = token[:idx]
				}
				if !isVersionString(token) {
					formatSegments = append(formatSegments, token)
				}
			}
		}

		return &HeuristicTarget{
			Format:     strings.Join(formatSegments, "/"),
			FieldNames: fields,
		}, nil
	}

	return nil, nil
}

// splitPathSegments is a custom character-by-character scanner for splitting
// a path by `/` while ignoring slashes inside curly braces like `{name=projects/*/locations/*/topics/*}`.
func splitPathSegments(pattern string) []string {
	var segments []string
	var current strings.Builder
	inVar := false

	for _, r := range pattern {
		switch r {
		case '{':
			inVar = true
			current.WriteRune(r)
		case '}':
			inVar = false
			current.WriteRune(r)
		case '/':
			if !inVar {
				// We've hit a true segment boundary. Yield the segment.
				segments = append(segments, current.String())
				current.Reset()
			} else {
				current.WriteRune(r)
			}
		default:
			current.WriteRune(r)
		}
	}
	segments = append(segments, current.String())
	return segments
}
