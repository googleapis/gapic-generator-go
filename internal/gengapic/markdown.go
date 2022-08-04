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

	"gitlab.com/golang-commonmark/markdown"
)

var linkParser = regexp.MustCompile(`<a\s+href=["'](.+)["']`)
var referenceParser = regexp.MustCompile(`\[([a-zA-Z1-9._]+)\]\[[a-zA-Z1-9._]*\]`)

func mdPlain(s string) string {
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
	listLevel   int
	// linkOpen indicates where there is an HTML link tag opened in order to prevent
	// printing softbreaks that would be empty lines in the tag value.
	linkOpen bool
}

func (m *mdRenderer) plain(t markdown.Token) {
	switch t := t.(type) {
	case *markdown.Inline:
		for _, c := range t.Children {
			m.plain(c)
		}
	case *markdown.Text:
		// Strip reference links, like [Foo][bar.Foo], that are invalid markdown
		// in the context of individual protobuf comments, down to just the
		// link text, which in this case is "Foo".
		content := referenceParser.ReplaceAllString(t.Content, "$1")
		m.sb.WriteString(content)
	case *markdown.CodeInline:
		m.sb.WriteString(t.Content)
	case *markdown.Softbreak:
		// Softbreaks in a link aren't supported
		if m.linkOpen {
			return
		}

		m.sb.WriteByte('\n')
		// indent multiple line list items according to list level
		m.indent()

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

	case *markdown.BulletListOpen:
		m.listLevel++
	case *markdown.BulletListClose:
		m.listLevel--

	case *markdown.ListItemOpen:
		// indent item content according to list level
		m.indent()
	case *markdown.ListItemClose:

	// ignore # headings
	case *markdown.HeadingOpen:
	case *markdown.HeadingClose:

	// ignore ** ** bold markers
	case *markdown.StrongOpen:
	case *markdown.StrongClose:

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
		m.linkOpen = true
	} else if t.Content == "</a>" {
		// print link
		l := len(m.linkTargets)
		fmt.Fprintf(&m.sb, " (at %s)", m.linkTargets[l-1])
		m.linkTargets = m.linkTargets[:l-1]
		m.linkOpen = false
	}
}

func (m *mdRenderer) indent() {
	if m.listLevel > 0 {
		for l := 0; l < m.listLevel; l++ {
			m.sb.WriteString("  ")
		}
	}
}
