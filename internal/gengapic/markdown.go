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

// TODO(pongad): maybe move markdown stuff into its own package
// once we have a proper import path.

package gengapic

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/golang-commonmark/markdown"
)

var linkParser = regexp.MustCompile(`<a href=["'](.+)["']`)

func MDPlain(s string) string {
	var mdr mdRenderer
	for _, tok := range markdown.New(markdown.HTML(true)).Parse([]byte(s)) {
		mdr.plain(tok)
	}

	return strings.TrimSpace(mdr.sb.String())
}

type mdRenderer struct {
	sb strings.Builder

	// NOTE(pongad): the parser parses links into a slice [linkopen, text, linkclose].
	// The ref is contained in the linkopen, so we need to save it here.
	// Because the data structure is an array, it's technically possible for the links
	// to nest, though I'm not sure if that'd be a valid Markdown.
	linkTargets []string
}

func (m *mdRenderer) plain(t markdown.Token) {
	switch t := t.(type) {
	case *markdown.Inline:
		for _, c := range t.Children {
			m.plain(c)
		}
	case *markdown.Text:
		m.sb.WriteString(t.Content)
	case *markdown.CodeInline:
		m.sb.WriteString(t.Content)
	case *markdown.Softbreak:
		m.sb.WriteByte('\n')

	case *markdown.ParagraphOpen:
	case *markdown.ParagraphClose:
		m.sb.WriteString("\n\n")

	case *markdown.LinkOpen:
		m.linkTargets = append(m.linkTargets, t.Href)
	case *markdown.LinkClose:
		l := len(m.linkTargets)
		fmt.Fprintf(&m.sb, " (at %s)", m.linkTargets[l-1])
		m.linkTargets = m.linkTargets[:l-1]

	case *markdown.HTMLInline:
		m.html(t)

	default:
		log.Printf("unhandled type: %T", t)
	}
}

func (m *mdRenderer) html(t *markdown.HTMLInline) {
	// font-based tags like <b> and most closing tags are just ignored entirely

	// write a new line for break tags
	if t.Content == "<br>" {
		m.sb.WriteByte('\n')
	} else if matches := linkParser.FindStringSubmatch(t.Content); len(matches) > 1 {
		// parse out href
		m.linkTargets = append(m.linkTargets, matches[1])
	} else if t.Content == "</a>" {
		// print link
		l := len(m.linkTargets)
		fmt.Fprintf(&m.sb, " (at %s)", m.linkTargets[l-1])
		m.linkTargets = m.linkTargets[:l-1]
	}
}
