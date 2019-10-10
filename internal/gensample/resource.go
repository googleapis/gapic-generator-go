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

package gensample

import (
	"strings"

	"github.com/googleapis/gapic-generator-go/internal/errors"
)

// namePattern represents a specially formatted string used in API definitions.
// The format is written in string form like this: "shelves/{shelf}/books/{book}".
// The curly braces delimit placeholders.
type namePattern struct {
	// Even positions are literal strings.
	// Odd positions are placeholders with braces already stripped.
	pieces []string

	// maps placeholder to the location in pieces
	pos map[string]int
}

func parseNamePattern(s string) (namePattern, error) {
	origPat := s
	var pieces []string

	for {
		p := strings.IndexByte(s, '{')
		if p < 0 {
			if s != "" {
				pieces = append(pieces, s)
			}
			break
		}
		pieces = append(pieces, s[:p])
		s = s[p:]

		p = strings.IndexByte(s, '}')
		if p < 0 {
			return namePattern{}, errors.E(nil, "unclosed brace: %q", origPat)
		}
		pieces = append(pieces, s[1:p])
		s = s[p+1:]
	}

	pos := map[string]int{}
	for i := 1; i < len(pieces); i += 2 {
		piece := pieces[i]
		if _, exist := pos[piece]; exist {
			return namePattern{}, errors.E(nil, "duplicate placeholder %q: %q", piece, origPat)
		}
		pos[piece] = i
	}
	return namePattern{pieces, pos}, nil
}

func (np namePattern) fmtSpec() string {
	var sb strings.Builder
	for i := range np.pieces {
		if i%2 == 0 {
			sb.WriteString(np.pieces[i])
		} else {
			sb.WriteString("%s")
		}
	}
	return sb.String()
}
