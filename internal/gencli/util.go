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

package gencli

import (
	"fmt"
	"strings"

	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"github.com/jhump/protoreflect/desc"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	// ShortDescMax is the maximum length accepted for
	// the Short usage docs
	ShortDescMax = 50

	// LongDescMax is the maximum length accepted for
	// the Long usage docs
	LongDescMax = 150
)

func toLongUsage(cmt string) string {
	return shorten(cmt, LongDescMax)
}

func toShortUsage(cmt string) string {
	return shorten(cmt, ShortDescMax)
}

func shorten(cmt string, limit int) string {
	if len(cmt) > limit {
		sep := strings.LastIndex(cmt[:limit], " ")
		if sep == -1 {
			sep = limit
		}
		cmt = cmt[:sep] + "..."
	}

	return cmt
}

func sanitizeComment(cmt string) string {
	cmt = strings.Replace(cmt, "\\", `\\`, -1)
	cmt = strings.Replace(cmt, "\n", " ", -1)
	cmt = strings.Replace(cmt, "\"", "'", -1)
	cmt = strings.TrimSpace(cmt)
	return cmt
}

func putImport(imports map[string]*pbinfo.ImportSpec, pkg *pbinfo.ImportSpec) {
	key := pkg.Name
	if key == "" {
		key = pkg.Path
	}

	imports[key] = pkg
}

// does not remove dot, but does split on it,
// title-casing each segment before rejoining.
func title(name string) string {
	split := strings.Split(name, "_")
	for i, s := range split {
		dotSplit := strings.Split(s, ".")
		for j, ds := range dotSplit {
			dotSplit[j] = toTitle(ds)
		}
		split[i] = strings.Join(dotSplit, ".")
	}

	return strings.Join(split, "")
}

// does not remove underscores
func dotToCamel(name string) (s string) {
	for _, tkn := range strings.Split(name, ".") {
		s += toTitle(tkn)
	}

	return
}

func toTitle(str string) string {
	return cases.Title(language.AmericanEnglish, cases.NoLower).String(str)
}

func oneofTypeName(field, inputMsgType string, flag *Flag) string {
	upperField := title(field)
	tname := fmt.Sprintf("%s_%s", inputMsgType, upperField)

	if flag.IsNested {
		imp := flag.MessageImport.Name

		// If the field is a message not defined in the same package as the
		// parent, we must use the parent's import name for the oneof wrapper
		// type.
		if flag.IsMessage() && flag.MsgDesc.GetFile().GetPackage() != flag.OneOfDesc.GetOwner().GetFile().GetPackage() {
			// There is no good way to handle an error from a template helper
			// function so we will let the generation fail because the import
			// will be empty - but this shouldn't ever happen because the Go pkg
			// option is required.
			p, _ := getImport(flag.OneOfDesc.GetOwner())
			imp = p.Name
		}
		tname = fmt.Sprintf("%s.%s_%s", imp, flag.Message, upperField)
	}

	if flag.IsMessage() {
		p := flag.MsgDesc.GetParent()
		// This is a nested message definition, check it against oneof type name
		if par, ok := p.(*desc.MessageDescriptor); ok {
			nname := par.GetName() + "_" + flag.MsgDesc.GetName()
			if strings.HasSuffix(tname, nname) {
				tname += "_"
			}
		}
	}

	return tname
}
