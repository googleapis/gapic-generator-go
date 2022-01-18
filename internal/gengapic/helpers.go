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
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
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
func isOptional(m *descriptor.DescriptorProto, n string) bool {
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

// hasMethod reports if the given service defines an RPC with the same name as
// the given simple method name.
func hasMethod(service *descriptor.ServiceDescriptorProto, method string) bool {
	for _, m := range service.GetMethod() {
		if m.GetName() == method {
			return true
		}
	}

	return false
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
func containsService(s []*descriptor.ServiceDescriptorProto, srv *descriptor.ServiceDescriptorProto) bool {
	for _, x := range s {
		if x.GetName() == srv.GetName() {
			return true
		}
	}

	return false
}

// isRequired returns if a field is annotated as REQUIRED or not.
func isRequired(field *descriptor.FieldDescriptorProto) bool {
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
