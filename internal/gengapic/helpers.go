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
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func lowerFirst(s string) string {
	if s == "" {
		return ""
	}
	r, w := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[w:]
}

func upperFirst(s string) string {
	if s == "" {
		return ""
	}
	r, w := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[w:]
}

func camelToSnake(s string) string {
	var sb strings.Builder
	runes := []rune(s)

	for i, r := range runes {
		if unicode.IsUpper(r) && i != 0 {
			// An uppercase rune followed by a lowercase
			// rune indicates the start of a word,
			// keeping uppercase acronyms together.
			next := i + 1
			if len(runes) > next && !unicode.IsUpper(runes[next]) {
				sb.WriteByte('_')
			}
		}
		sb.WriteRune(unicode.ToLower(r))
	}
	return sb.String()
}

// snakeToCamel converts snake_case and SNAKE_CASE to CamelCase.
func snakeToCamel(s string) string {
	var sb strings.Builder
	up := true
	for _, r := range s {
		if r == '_' {
			up = true
		} else if up && unicode.IsDigit(r) {
			sb.WriteRune('_')
			sb.WriteRune(r)
			up = false
		} else if up {
			sb.WriteRune(unicode.ToUpper(r))
			up = false
		} else {
			sb.WriteRune(unicode.ToLower(r))
		}
	}
	return sb.String()
}

// isOptional returns true if the named Field in the given Message
// is proto3_optional.
func isOptional(m *descriptorpb.DescriptorProto, n string) bool {
	for _, f := range m.GetField() {
		if f.GetName() == n {
			return f.GetProto3Optional()
		}
	}

	return false
}

func strContains(a []string, s string) bool {
	for _, as := range a {
		if as == s {
			return true
		}
	}
	return false
}

// grpcClientField reports the field name to store gRPC client.
func grpcClientField(reducedServName string) string {
	// Not the same as pbinfo.ReduceServName(*serv.Name, pkg)+"Client".
	// If the service name is reduced to empty string, we should
	// lower-case "client" so that the field is not exported.
	return lowerFirst(reducedServName + "Client")
}

// getField returns a FieldDescriptorProto pointer if the target
// DescriptorProto has the given field, otherwise it returns nil.
func getField(m *descriptorpb.DescriptorProto, field string) *descriptorpb.FieldDescriptorProto {
	for _, f := range m.GetField() {
		if f.GetName() == field {
			return f
		}
	}
	return nil
}

// hasField returns true if the target DescriptorProto has the given field,
// otherwise it returns false.
func hasField(m *descriptorpb.DescriptorProto, field string) bool {
	return getField(m, field) != nil
}

// hasMethod reports if the given service defines an RPC with the same name as
// the given simple method name.
func hasMethod(service *descriptorpb.ServiceDescriptorProto, method string) bool {
	for _, m := range service.GetMethod() {
		if m.GetName() == method {
			return true
		}
	}

	return false
}

// hasRESTMethod reports if there is at least one RPC on the Service that
// has a gRPC-HTTP transcoding, or REST, annotation on it.
func hasRESTMethod(service *descriptorpb.ServiceDescriptorProto) bool {
	for _, m := range service.GetMethod() {
		eHTTP := proto.GetExtension(m.GetOptions(), annotations.E_Http)
		if h := eHTTP.(*annotations.HttpRule); h.GetPattern() != nil {
			return true
		}
	}

	return false
}

// getMethod returns the MethodDescriptorProto for the given service RPC and simple method name.
func getMethod(service *descriptorpb.ServiceDescriptorProto, method string) *descriptorpb.MethodDescriptorProto {
	for _, m := range service.GetMethod() {
		if m.GetName() == method {
			return m
		}
	}

	return nil
}

// containsTransport determines if a set of transports contains a specific
// transport.
func containsTransport(t []transport, tr transport) bool {
	for _, x := range t {
		if x == tr {
			return true
		}
	}

	return false
}

// containsService determines if a set of services contains a specific service,
// by simple name.
func containsService(s []*descriptorpb.ServiceDescriptorProto, srv *descriptorpb.ServiceDescriptorProto) bool {
	for _, x := range s {
		if x.GetName() == srv.GetName() {
			return true
		}
	}

	return false
}

// isRequired returns if a field is annotated as REQUIRED or not.
func isRequired(field *descriptorpb.FieldDescriptorProto) bool {
	if field.GetOptions() == nil {
		return false
	}

	eBehav := proto.GetExtension(field.GetOptions(), annotations.E_FieldBehavior)

	behaviors := eBehav.([]annotations.FieldBehavior)
	for _, b := range behaviors {
		if b == annotations.FieldBehavior_REQUIRED {
			return true
		}
	}

	return false
}

// This takes in a path template from a routing annotation and converts it into a regex string.
// The named capture is the named segment portion for the header itself.
func convertPathTemplateToRegex(pattern string) string {
	// If path template doesn't exist, then use a wildcard.
	if pattern == "" {
		return "(.*)"
	}
	// Replace name of header to named capture.
	regexPattern := strings.ReplaceAll(pattern, "{", "(?P<")
	regexPattern = strings.ReplaceAll(regexPattern, "}", ")")
	// If not named, then entire segment is a wildcard.
	if !strings.Contains(pattern, "=") || !strings.Contains(pattern, "/") {
		regexPattern = strings.ReplaceAll(regexPattern, "*", "")
		regexPattern = strings.ReplaceAll(regexPattern, "=", "")
		regexPattern = strings.ReplaceAll(regexPattern, ")", ">.*)")
		return regexPattern
	}
	// Replace segment wildcards with regex equivalent
	regexPattern = strings.ReplaceAll(regexPattern, "/**", "(?:/.*)?")
	regexPattern = strings.ReplaceAll(regexPattern, "/*", "/[^/]+")
	regexPattern = strings.ReplaceAll(regexPattern, "=**", ">.*")
	regexPattern = strings.ReplaceAll(regexPattern, "=*", ">[^/]+")
	regexPattern = strings.ReplaceAll(regexPattern, "=", ">")
	regexPattern = strings.ReplaceAll(regexPattern, "**", ".*")
	return regexPattern
}

// This intakes a path template and returns the name of the header to be returned.
func getHeaderName(pattern string) string {
	curlyBraceRegex := regexp.MustCompile(`{([^}]+)\}`)
	// Path template should only contain one name (e.g. at most one `=`) or
	// a collectionId, and should be contained within curly braces.
	if strings.Count(pattern, "=") > 1 || !curlyBraceRegex.MatchString(pattern) {
		return ""
	}
	// curlyBraceSegment returns the named capture within the path template that is within the curly braces.
	curlyBraceSegment := curlyBraceRegex.FindStringSubmatch(pattern)[1]
	// If there is no equal sign, then the path template is a collectionId which is its own name.
	// and both the named segment and path template are wildcards.
	if strings.Count(pattern, "=") < 1 {
		return curlyBraceSegment
	}
	getBeforeEqualsSign := regexp.MustCompile("(?P<before>[^=]*)=.*")
	matches := getBeforeEqualsSign.FindStringSubmatch(curlyBraceSegment)
	return matches[1]
}
