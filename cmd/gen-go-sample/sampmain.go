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

package main

import (
	"fmt"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"strings"
)

func writeMain(g *generator, argNames []string, flagNames []string, argTrees []*initTree, methName string) error {
	p := g.pt.Printf
	st := newSymTab(nil)

	// We only care about symbols here, not the types
	dummyType := initType{}
	for _, arg := range argNames {
		if err := st.put(arg, dummyType); err != nil {
			return err
		}
	}

	p("func main() {")

	for i := range argNames {
		if _, ok := argTrees[i].typ.desc.(*descriptor.EnumDescriptorProto); ok {
			p(`%s := flag.String(%q, %q, "")`, argNames[i], flagNames[i], argTrees[i].leafVal)
			continue
		}

		// TODO(pongad): some types, like int32, are not supported by flag package.
		// We have to convert.
		typ := pbinfo.GoTypeForPrim[argTrees[i].typ.prim]
		p(`%s := flag.%s(%q, %s, "")`, argNames[i], snakeToPascal(typ), flagNames[i], argTrees[i].leafVal)
	}

	p("  flag.Parse()")

	var fArgs []string

	for i := range argNames {
		if e, ok := argTrees[i].typ.desc.(*descriptor.EnumDescriptorProto); ok {
			v := st.disambiguate(argNames[i], dummyType)
			tn, err := enumType(g.descInfo, e, g)
			if err != nil {
				return err
			}
			vMap := tn + "_name"
			st.put(v, dummyType)
			fArgs = append(fArgs, fmt.Sprintf("%s(%s)", tn, v))

			p(`%s, ok := %s[*%s]`, v, vMap, argNames[i])
			p("if !ok {")
			p(`log.fatal("enum type %s does not have value %%s", *%s)`, tn, v)
			p("}")
		} else {
			fArgs = append(fArgs, argNames[i])
		}
	}

	p("  if err := sample%s(%s); err != nil {", methName, flagArgs(fArgs))
	p("    log.Fatal(err)")
	p("  }")
	p("}")
	p("")

	g.imports[pbinfo.ImportSpec{Path: "flag"}] = true
	g.imports[pbinfo.ImportSpec{Path: "fmt"}] = true
	g.imports[pbinfo.ImportSpec{Path: "log"}] = true
	return nil
}

func flagArgs(names []string) string {
	if len(names) == 0 {
		return ""
	}

	var sb strings.Builder
	for _, n := range names {
		fmt.Fprintf(&sb, ", *%s", n)
	}
	return sb.String()[2:]
}
