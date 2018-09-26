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
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/gapic-generator-go/internal/errors"
	"github.com/googleapis/gapic-generator-go/internal/license"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"github.com/googleapis/gapic-generator-go/internal/printer"
	yaml "gopkg.in/yaml.v2"
)

func main() {
	descFname := flag.String("desc", "", "proto descriptor")
	gapicFname := flag.String("gapic", "", "gapic config")
	clientPkg := flag.String("clientpkg", "", "the package of the client, in format 'url/to/client/pkg;name'")
	nofmt := flag.Bool("nofmt", false, "skip gofmt, useful for debugging code with syntax error")
	flag.Parse()

	gen := generator{
		imports: map[pbinfo.ImportSpec]bool{},
	}

	if p := strings.IndexByte(*clientPkg, ';'); p >= 0 {
		gen.clientPkg = pbinfo.ImportSpec{Path: (*clientPkg)[:p], Name: (*clientPkg)[p+1:]}
	} else {
		log.Fatalf("need -clientPkg in 'url/to/client/pkg;name' format, got %q", *clientPkg)
	}

	donec := make(chan struct{})
	go func() {
		f, err := os.Open(*gapicFname)
		if err != nil {
			log.Fatal(errors.E(err, "cannot read GAPIC config file"))
		}
		defer f.Close()

		if err := yaml.NewDecoder(f).Decode(&gen.gapic); err != nil {
			log.Fatal(errors.E(err, "error reading GAPIC config file"))
		}
		donec <- struct{}{}
	}()
	go func() {
		descBytes, err := ioutil.ReadFile(*descFname)
		if err != nil {
			log.Fatal(errors.E(err, "cannot read proto descriptor file"))
		}

		if err := proto.Unmarshal(descBytes, &gen.desc); err != nil {
			log.Fatal(errors.E(err, "error reading proto descriptor file"))
		}

		gen.descInfo = pbinfo.Of(gen.desc.GetFile())
		donec <- struct{}{}
	}()

	<-donec
	<-donec

	// TODO: split some parts of this loop into functions, so we can idiomatically
	// chain errors.
	for _, iface := range gen.gapic.Interfaces {
		for _, meth := range iface.Methods {
			valSets := map[string]SampleValueSet{}
			for _, vs := range meth.SampleValueSets {
				valSets[vs.ID] = vs
			}

			for _, sam := range meth.Samples.Standalone {
				for _, vsID := range sam.ValueSets {
					vs, ok := valSets[vsID]
					if !ok {
						log.Fatal(errors.E(nil, "value set not found: %q", vsID))
					}

					gen.reset()
					gen.imports[gen.clientPkg] = true
					if err := gen.genSample(iface.Name, meth.Name, sam.RegionTag, vs); err != nil {
						log.Fatal(errors.E(err, "value set: %q", vs))
					}
					if err := gen.commit(!*nofmt); err != nil {
						log.Fatal(errors.E(err, "can't commit value set: %q", vs))
					}
				}
			}
		}
	}
}

type generator struct {
	desc     descriptor.FileDescriptorSet
	descInfo pbinfo.Info
	gapic    GAPICConfig

	clientPkg pbinfo.ImportSpec

	pt      printer.P
	imports map[pbinfo.ImportSpec]bool
}

func (g *generator) reset() {
	g.pt.Reset()
	for imp := range g.imports {
		delete(g.imports, imp)
	}
}

func (g *generator) commit(gofmt bool) error {
	// We'll gofmt unless user asks us to not, so no need to think too hard about sorting
	// "correctly". We just want a deterministic output here.
	var imports []pbinfo.ImportSpec
	for imp := range g.imports {
		imports = append(imports, imp)
	}
	sort.Slice(imports, func(i, j int) bool {
		// All non-stdlib imports come after stdlib imports.
		iDot := strings.IndexByte(imports[i].Path, '.') >= 0
		jDot := strings.IndexByte(imports[j].Path, '.') >= 0
		if iDot != jDot {
			return jDot
		}

		if imports[i].Path != imports[j].Path {
			return imports[i].Path < imports[j].Path
		}
		return imports[i].Name < imports[j].Name
	})

	firstNonStd := sort.Search(len(imports), func(i int) bool { return strings.IndexByte(imports[i].Path, '.') >= 0 })

	var file bytes.Buffer
	fmt.Fprintf(&file, license.Apache, time.Now().Year())
	file.WriteString("package main\n")
	file.WriteString("import(\n")
	for i, imp := range imports {
		if i == firstNonStd {
			file.WriteByte('\n')
		}
		fmt.Fprintf(&file, "%s %q\n", imp.Name, imp.Path)
	}
	file.WriteString(")\n")
	file.Write(g.pt.Bytes())

	b := file.Bytes()
	if gofmt {
		b2, err := format.Source(b)
		if err != nil {
			return errors.E(err, "syntax error, run with -nofmt to find out why?")
		}
		b = b2
	}
	_, err := os.Stdout.Write(b)
	return err
}

func (g *generator) genSample(ifaceName, methName, regTag string, valSet SampleValueSet) error {
	serv := g.descInfo.Serv["."+ifaceName]
	if serv == nil {
		return errors.E(nil, "can't find service: %q", ifaceName)
	}

	var meth *descriptor.MethodDescriptorProto
	for _, m := range serv.GetMethod() {
		if m.GetName() == methName {
			meth = m
			break
		}
	}
	if meth == nil {
		return errors.E(nil, "service %q doesn't have method %q", serv.GetName(), methName)
	}

	p := g.pt.Printf

	p("// [START %s]", regTag)
	p("")

	p("func sample%s() {", methName)
	p("  ctx := context.Background()")
	p("  c := %s.New%sClient(ctx)", g.clientPkg.Name, pbinfo.ReduceServName(serv.GetName(), g.clientPkg.Name))
	p("")
	g.imports[pbinfo.ImportSpec{Path: "context"}] = true

	// TODO(pongad): properly create the request.
	inType := g.descInfo.Type[meth.GetInputType()]
	if inType == nil {
		return errors.E(nil, "can't find input type %q", meth.GetInputType())
	}
	inSpec, err := g.descInfo.ImportSpec(inType)
	if err != nil {
		return errors.E(err, "can't import input type: %q", inType)
	}

	itree := initTree{
		typ: initType{desc: inType},
	}
	for _, def := range valSet.Parameters.Defaults {
		if err := itree.Parse(def, g.descInfo); err != nil {
			return errors.E(err, "can't set default value: %q", def)
		}
	}
	{
		w := g.pt.Writer()

		if _, err := w.Write([]byte("req := ")); err != nil {
			return err
		}
		if err := itree.Print(g.pt.Writer(), g); err != nil {
			return err
		}
		if _, err := w.Write([]byte{'\n'}); err != nil {
			return err
		}
	}

	// TODO(pongad): handle non-unary
	p("  resp, err := c.%s(ctx, req)", methName)
	p("  if err != nil {")
	p("    // TODO: Handle error.")
	p("  }")

	// TODO(pongad): handle output printing
	p("  fmt.Println(resp)")

	p("}")
	p("")
	p("// [END %s]", regTag)
	p("")

	p("func main() {")
	p("  sample%s()", methName)
	p("}")
	p("")

	g.imports[inSpec] = true
	g.imports[pbinfo.ImportSpec{Path: "fmt"}] = true
	return nil
}
