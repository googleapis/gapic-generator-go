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
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
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

func grpcClientCall(reducedServName, methName string) string {
	return fmt.Sprintf("c.%s.%s(ctx, req, settings.GRPC...)", grpcClientField(reducedServName), methName)
}
