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
	"path/filepath"
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
	outDir := flag.String("o", ".", "directory to write samples to")
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

	if err := os.MkdirAll(*outDir, 0755); err != nil {
		log.Fatal(err)
	}

	<-donec
	<-donec

	for _, iface := range gen.gapic.Interfaces {
		for _, meth := range iface.Methods {
			if err := genMethodSamples(&gen, iface, meth, *nofmt, *outDir); err != nil {
				err = errors.E(err, "generating: %s", iface.Name+"."+meth.Name)
				log.Fatal(err)
			}
		}
	}
}

func genMethodSamples(gen *generator, iface GAPICInterface, meth GAPICMethod, nofmt bool, outDir string) error {
	valSets := map[string]SampleValueSet{}
	for _, vs := range meth.SampleValueSets {
		valSets[vs.ID] = vs
	}

	for _, sam := range meth.Samples.Standalone {
		for _, vsID := range sam.ValueSets {
			vs, ok := valSets[vsID]
			if !ok {
				return errors.E(nil, "value set not found: %q", vsID)
			}

			gen.reset()
			if err := gen.genSample(iface.Name, meth, sam.RegionTag, vs); err != nil {
				return errors.E(err, "value set: %s", vsID)
			}
			content, err := gen.commit(!nofmt, time.Now().Year())
			if err != nil {
				return err
			}

			fname := iface.Name + "_" + meth.Name + "_" + vsID + ".go"
			if err := ioutil.WriteFile(filepath.Join(outDir, fname), content, 0644); err != nil {
				return err
			}
		}
	}
	return nil
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

func (g *generator) commit(gofmt bool, year int) ([]byte, error) {
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
	fmt.Fprintf(&file, license.Apache, year)
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
			return nil, errors.E(err, "syntax error, run with -nofmt to find out why?")
		}
		b = b2
	}
	return b, nil
}

func (g *generator) genSample(ifaceName string, methConf GAPICMethod, regTag string, valSet SampleValueSet) error {
	// TODO(pongad): This method's error cases are not well tested. Split and test.

	// Preparation
	g.imports[g.clientPkg] = true
	serv := g.descInfo.Serv["."+ifaceName]
	if serv == nil {
		return errors.E(nil, "can't find service: %q", ifaceName)
	}

	var meth *descriptor.MethodDescriptorProto
	for _, m := range serv.GetMethod() {
		if m.GetName() == methConf.Name {
			meth = m
			break
		}
	}
	if meth == nil {
		return errors.E(nil, "service %q doesn't have method %q", serv.GetName(), methConf.Name)
	}

	inType := g.descInfo.Type[meth.GetInputType()]
	if inType == nil {
		return errors.E(nil, "can't find input type %q", meth.GetInputType())
	}
	inSpec, err := g.descInfo.ImportSpec(inType)
	if err != nil {
		return errors.E(err, "can't import input type: %q", inType)
	}
	g.imports[inSpec] = true

	initInfo, err := g.getInitInfo(inType, methConf, valSet)
	argNames := initInfo.argNames
	argTrees := initInfo.argTrees
	flagNames := initInfo.flagNames
	files := initInfo.files
	itree := initInfo.reqTree

	if err := g.sampFuncAndClient(serv.GetName(), meth.GetName(), initInfo, regTag); err != nil {
		return nil
	}

	// Set up the request
	{
		var buf bytes.Buffer

		for i, name := range argNames {
			fmt.Fprintf(&buf, "%s := ", name)
			if err := argTrees[i].Print(&buf, g); err != nil {
				return errors.E(err, "can't initialize parameter: %s", name)
			}
			buf.WriteByte('\n')
		}
		prependLines(&buf, "// ", false)

		for _, info := range files {
			handleReadFile(info, &buf, g)
		}

		buf.WriteString("req := ")
		if err := itree.Print(&buf, g); err != nil {
			return errors.E(err, "can't initialize request object")
		}
		buf.WriteByte('\n')
		prependLines(&buf, "\t", true)

		if _, err := buf.WriteTo(g.pt.Writer()); err != nil {
			return err
		}
	}

	// Make the RPC call and handle output
	if meth.GetOutputType() == ".google.protobuf.Empty" {
		err = g.emptyOut(meth, valSet)
	} else if meth.GetOutputType() == ".google.longrunning.Operation" {
		err = g.lro(meth, methConf, valSet)
	} else if meth.GetServerStreaming() || meth.GetClientStreaming() {
		err = errors.E(nil, "streaming methods not supported yet")
	} else if pf, err2 := pagingField(g.descInfo, meth); err2 != nil {
		err = errors.E(err2, "can't determine whether method is paging")
	} else if pf != nil {
		err = g.paging(meth, pf, valSet)
	} else {
		err = g.unary(meth, valSet)
	}

	if err != nil {
		return err
	}

	p := g.pt.Printf
	p("return nil")
	p("}")
	p("")
	p("// [END %s]", regTag)
	p("")

	// main
	return writeMain(g, argNames, flagNames, argTrees, meth.GetName())
}

func (g *generator) getInitInfo(inType pbinfo.ProtoType, methConf GAPICMethod, valSet SampleValueSet) (initInfo, error) {
	var (
		argNames  []string
		flagNames []string
		argTrees  []*initTree
		files     []*fileInfo
	)

	itree := initTree{
		typ: initType{desc: inType},
	}

	// Set up resource names. We need this info when setting up request object.
	// TODO(hzyi): handle collection_oneofs
	for field, entName := range methConf.FieldNamePatterns {
		var pat string
		for _, col := range g.gapic.Collections {
			if col.EntityName == entName {
				pat = col.NamePattern
				break
			}
		}
		if pat == "" {
			return initInfo{}, errors.E(nil, "undefined resource name: %q", entName)
		}

		namePat, err := parseNamePattern(pat)
		if err != nil {
			return initInfo{}, err
		}

		subTree, err := itree.get(field, g.descInfo)
		if err != nil {
			return initInfo{}, errors.E(err, "cannot set up resource name: %q", entName)
		}
		if typ := subTree.typ; typ.prim != descriptor.FieldDescriptorProto_TYPE_STRING {
			return initInfo{}, errors.E(err, "cannot set up resource name for %q, need field to be string, got %v", field, typ)
		}
		subTree.typ.prim = 0
		subTree.typ.namePat = &namePat
	}

	// Set up request object.
	for _, def := range valSet.Parameters.Defaults {
		if err := itree.parseInit(def, g.descInfo); err != nil {
			return initInfo{}, errors.E(err, "can't set default value: %q", def)
		}
	}

	// Some parts of request object are from arguments.
	for _, attr := range valSet.Parameters.Attributes {
		if attr.ReadFile {
			var varName string
			var err error
			if attr.SampleArgumentName != "" {
				varName = snakeToCamel(attr.SampleArgumentName) + fileContentSuffix
			} else {
				varName, err = fileVarName(attr.Parameter)
				if err != nil {
					return initInfo{}, errors.E(err, "can't determine variable name to store bytes from local file")
				}
			}

			subTree, err := itree.parseSampleArgPath(
				attr.Parameter,
				g.descInfo,
				varName,
			)
			if err != nil {
				return initInfo{}, errors.E(err, "can't set sample function argument: %q", attr.Parameter)
			}
			if subTree.typ.prim != descriptor.FieldDescriptorProto_TYPE_BYTES {
				return initInfo{}, errors.E(nil, "can only assign file contents to bytes field")
			}
			subTree.typ.prim = descriptor.FieldDescriptorProto_TYPE_STRING
			subTree.typ.valFmt = nil
			if subTree.leafVal == "" {
				return initInfo{}, errors.E(nil, "default value not given: %q", attr.Parameter)
			}
			fileName := subTree.leafVal
			if attr.SampleArgumentName != "" {
				fileName = snakeToCamel(attr.SampleArgumentName)
				argNames = append(argNames, fileName)
				flagNames = append(flagNames, attr.SampleArgumentName)
				argTrees = append(argTrees, subTree)
			}
			files = append(files, &fileInfo{fileName, varName})
			continue
		}
		if attr.SampleArgumentName == "" {
			continue
		}
		name := snakeToCamel(attr.SampleArgumentName)
		subTree, err := itree.parseSampleArgPath(attr.Parameter, g.descInfo, name)
		if err != nil {
			return initInfo{}, errors.E(err, "can't set sample function argument: %q", attr.Parameter)
		}

		argNames = append(argNames, name)
		flagNames = append(flagNames, attr.SampleArgumentName)
		argTrees = append(argTrees, subTree)
	}

	initInfo := initInfo{
		argNames:  argNames,
		argTrees:  argTrees,
		flagNames: flagNames,
		files:     files,
		reqTree:   itree,
	}

	return initInfo, nil
}

func (g *generator) sampFuncAndClient(servName string, methName string, init initInfo, regTag string) error {
	var argStr string
	if len(init.argNames) > 0 {
		var sb strings.Builder
		for i, name := range init.argNames {
			typ, err := g.getGoTypeName(init.argTrees[i].typ)
			if err != nil {
				return err
			}
			fmt.Fprintf(&sb, ", %s %s", name, typ)
		}
		argStr = sb.String()[2:]
	}

	p := g.pt.Printf

	p("// [START %s]", regTag)
	p("")
	p("func sample%s(%s) error {", methName, argStr)
	p("  ctx := context.Background()")
	p("  c, err := %s.New%sClient(ctx)", g.clientPkg.Name, pbinfo.ReduceServName(servName, g.clientPkg.Name))
	p("  if err != nil {")
	p("    return err")
	p("  }")
	p("")
	g.imports[pbinfo.ImportSpec{Path: "context"}] = true
	return nil
}

func (g *generator) unary(meth *descriptor.MethodDescriptorProto, valSet SampleValueSet) error {
	p := g.pt.Printf

	p("resp, err := c.%s(ctx, req)", meth.GetName())
	p("if err != nil {")
	p("  return err")
	p("}")
	p("")

	return g.handleOut(meth, valSet, &initType{desc: g.descInfo.Type[meth.GetOutputType()]})
}

func (g *generator) emptyOut(meth *descriptor.MethodDescriptorProto, valSet SampleValueSet) error {
	p := g.pt.Printf

	p("if err := c.%s(ctx, req); err != nil {", meth.GetName())
	p("  return err")
	p("}")
	p("")

	return g.handleOut(meth, valSet, nil)
}

func (g *generator) paging(meth *descriptor.MethodDescriptorProto, pf *descriptor.FieldDescriptorProto, valSet SampleValueSet) error {
	p := g.pt.Printf

	p("it := c.%s(ctx, req)", meth.GetName())
	p("for {")
	p("  resp, err := it.Next()")
	p("  if err == iterator.Done {")
	p("    break")
	p("  }")
	p("  if err != nil {")
	p("    return err")
	p("  }")

	var typ initType
	if tn := pf.GetTypeName(); tn != "" {
		typ = initType{desc: g.descInfo.Type[tn]}
	} else {
		typ = initType{prim: pf.GetType()}
	}

	err := g.handleOut(meth, valSet, &typ)

	p("}")
	p("")
	g.imports[pbinfo.ImportSpec{Path: "google.golang.org/api/iterator"}] = true

	return err
}

func (g *generator) lro(meth *descriptor.MethodDescriptorProto, methConf GAPICMethod, valSet SampleValueSet) error {
	p := g.pt.Printf

	p("op, err := c.%s(ctx, req)", meth.GetName())
	p("if err != nil {")
	p("  return err")
	p("}")
	p("")
	p("resp, err := op.Wait(ctx)")
	p("if err != nil {")
	p("  return err")
	p("}")
	p("")

	retType := methConf.LongRunning.ReturnType
	if retType == "" {
		return errors.E(nil, "LRO return type not given")
	}

	typ := initType{desc: g.descInfo.Type["."+retType]}
	return g.handleOut(meth, valSet, &typ)
}

func (g *generator) handleOut(meth *descriptor.MethodDescriptorProto, valSet SampleValueSet, respTyp *initType) error {
	st := newSymTab(nil)
	if respTyp != nil {
		st.universe["$resp"] = true
		st.scope["$resp"] = *respTyp
	}

	for _, out := range valSet.OnSuccess {
		if err := writeOutputSpec(out, st, g); err != nil {
			return errors.E(err, "cannot write output handling code")
		}
	}

	if respTyp != nil && len(valSet.OnSuccess) == 0 {
		g.pt.Printf("fmt.Println(resp)")
		g.imports[pbinfo.ImportSpec{Path: "fmt"}] = true
	}
	return nil
}

// prependLines adds prefix to every line in b. A line is defined as a possibly empty run
// of non-newlines terminated by a newline character.
// If b doesn't end with a newline, prependLines panics.
// If skipEmptyLine is true, prependLines does not prepend prefix to empty lines.
func prependLines(b *bytes.Buffer, prefix string, skipEmptyLine bool) {
	if b.Len() == 0 {
		return
	}
	if b.Bytes()[b.Len()-1] != '\n' {
		panic(errors.E(nil, "prependLines: must end with newline"))
	}
	// Don't split with b.Bytes; we have to make a copy of the content since we're overwriting the buffer.
	lines := strings.SplitAfter(b.String(), "\n")
	b.Reset()
	for _, l := range lines {
		// When splitting, we have an empty string after the last newline, ignore it.
		if l == "" {
			continue
		}
		if !skipEmptyLine || l != "\n" {
			b.WriteString(prefix)
		}
		b.WriteString(l)
	}
}
