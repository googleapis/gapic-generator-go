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

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/gapic-generator-go/internal/errors"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"github.com/googleapis/gapic-generator-go/internal/printer"
	yaml "gopkg.in/yaml.v2"
)

func main() {
	descFname := flag.String("desc", "", "proto descriptor")
	gapicFname := flag.String("gapic", "", "gapic config")
	nofmt := flag.Bool("nofmt", false, "skip gofmt, useful for debugging code with syntax error")
	flag.Parse()

	var (
		gen = generator{
			imps: map[pbinfo.ImportSpec]bool{},
		}

		donec = make(chan struct{})
	)

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

		gen.dinfo = pbinfo.Of(gen.desc.GetFile())
		donec <- struct{}{}
	}()

	<-donec
	<-donec

	for _, iface := range gen.gapic.Interfaces {
		for _, meth := range iface.Methods {
			valSets := map[string]GapicValueSet{}
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
					if err := gen.genSample(iface.Name, meth.Name, vs); err != nil {
						log.Fatal(errors.E(err, "value set: %q", vs))
					}
					if err := gen.commit(!*nofmt); err != nil {
						log.Fatal(errors.E(err, "value set: %q", vs))
					}
				}
			}
		}
	}
}

type generator struct {
	desc  descriptor.FileDescriptorSet
	dinfo pbinfo.Info
	gapic GapicConfig

	pt   printer.P
	imps map[pbinfo.ImportSpec]bool
}

func (g *generator) reset() {
	g.pt.Reset()
	for imp := range g.imps {
		delete(g.imps, imp)
	}
}

func (g *generator) commit(gofmt bool) error {
	// We'll gofmt unless user asks us to not, so no need to think too hard about sorting
	// "correctly". We just want a deterministic output here.
	var imps []pbinfo.ImportSpec
	for imp := range g.imps {
		imps = append(imps, imp)
	}
	sort.Slice(imps, func(i, j int) bool {
		if imps[i].Path != imps[j].Path {
			return imps[i].Path < imps[j].Path
		}
		return imps[i].Name < imps[j].Name
	})

	var file bytes.Buffer
	file.WriteString("package main\n")
	file.WriteString("import(\n")
	for _, imp := range imps {
		fmt.Fprintf(&file, "%s %q", imp.Name, imp.Path)
	}
	file.WriteString(")\n")
	file.Write(g.pt.Bytes())

	b := file.Bytes()
	if gofmt {
		b2, err := format.Source(b)
		if err != nil {
			return err
		}
		b = b2
	}
	_, err := os.Stdout.Write(b)
	return err
}

func (g *generator) genSample(ifaceName, methName string, valSet GapicValueSet) error {
	serv := g.dinfo.Serv["."+ifaceName]
	if serv == nil {
		return errors.E(nil, "can't find service: %q", ifaceName)
	}

	servSpec, err := g.dinfo.ImportSpec(serv)
	if err != nil {
		return errors.E(err, "can't import service: %q", ifaceName)
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

	p("// [START %s]", valSet.ID)
	p("")

	// TODO(pongad): properly reduce service name instead of hardcoding "NewClient"
	p("func sample%s() {", methName)
	p("  ctx := context.Background()")
	p("  c := %s.%s(ctx)", servSpec.Name, "NewClient")
	p("")

	// TODO(pongad): properly create the request.
	inType := g.dinfo.Type[meth.GetInputType()]
	if inType == nil {
		return errors.E(nil, "can't find input type %q", meth.GetInputType())
	}
	inSpec, err := g.dinfo.ImportSpec(inType)
	if err != nil {
		return errors.E(err, "can't import input type: %q", inType)
	}

	p("  req := %s.%s{", inSpec.Name, inType.GetName())
	for _, def := range valSet.Parameters.Defaults {
		p("// %s", def)
	}
	p("  }")

	// TODO(pongad): handle non-unary
	p("  resp, err := c.%s(ctx, req)", methName)
	p("  if err != nil {")
	p("    // TODO: Handle error.")
	p("  }")

	// TODO(pongad): handle output printing
	p("  fmt.Println(resp)")

	p("}")
	p("")
	p("// [END %s]", valSet.ID)
	p("")

	p("func main() {")
	p("  sample%s()", methName)
	p("}")
	p("")

	g.imps[servSpec] = true
	g.imps[inSpec] = true
	return nil
}
