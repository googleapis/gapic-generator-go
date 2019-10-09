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
	"bytes"
	"crypto/sha256"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"go/format"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/gapic-generator-go/internal/errors"
	"github.com/googleapis/gapic-generator-go/internal/gensample/schema_v1p2"
	"github.com/googleapis/gapic-generator-go/internal/license"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"github.com/googleapis/gapic-generator-go/internal/printer"
	yaml "gopkg.in/yaml.v2"
)

const expectedSampleConfigType = "com.google.api.codegen.samplegen.v1p2.SampleConfigProto"
const expectedSampleConfigVersion = "1.2.0"

// InitGen creates a new sample generator.
// sampleFnames is the filenames of all sample config files.
// gapicFname is the filename of gapic config.
// clientPkg is the Go package of the generated gapic client library.
// nofmt set to true will instruct the generator not to format the generated code. This could be useful for debugging purposes.
func InitGen(desc []*descriptor.FileDescriptorProto, sampleFnames []string, gapicFname string, clientPkg string, nofmt bool) (*generator, error) {

	gen := generator{
		imports:      map[pbinfo.ImportSpec]bool{},
		desc:         desc,
		descInfo:     pbinfo.Of(desc),
		sampleConfig: schema_v1p2.SampleConfig{},
	}

	if p := strings.IndexByte(clientPkg, ';'); p >= 0 {
		gen.clientPkg = pbinfo.ImportSpec{Path: (clientPkg)[:p], Name: (clientPkg)[p+1:]}
	} else {
		return nil, errors.E(nil, "need -clientPkg in 'url/to/client/pkg;name' format, got %q", clientPkg)
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		errChan <- gen.readGapicConfigFile(gapicFname)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		errChan <- gen.readSampleConfigFiles(sampleFnames)
	}()
	wg.Wait()

	close(errChan)
	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}
	return &gen, nil
}

// GenMethodSamples generators samples from protos and configurations stored in the generator,
// and writes the generated samples to generator.Outputs.
func (gen *generator) GenMethodSamples() error {
	gen.disambiguateSampleIDs()
	gen.Outputs = make(map[string][]byte)

	for _, samp := range gen.sampleConfig.Samples {
		var iface GAPICInterface
		var method GAPICMethod
		for _, iface = range gen.gapic.Interfaces {
			if iface.Name == samp.Service {
				for _, method = range iface.Methods {
					if method.Name == samp.Rpc {
						break
					}
				}
				break
			}
		}

		if method.Name != samp.Rpc {
			return errors.E(nil, "generating sample %q: rpc %q not found", samp.ID, samp.Rpc)
		}
		if iface.Name != samp.Service {
			return errors.E(nil, "generating sample %q: service %q not found", samp.ID, samp.Service)
		}

		gen.reset()
		if err := gen.genSample(*samp, method); err != nil {
			err = errors.E(err, "generating: %s.%s:%s", iface.Name, method.Name, samp.ID)
			log.Fatal(err)
		}

		content, err := gen.commit(!gen.nofmt, time.Now().Year())
		if err != nil {
			return err
		}

		fname := samp.ID + ".go"
		gen.Outputs[fname] = content
	}
	return nil
}

type generator struct {
	// desc is the set of proto descriptors of the API
	desc []*descriptor.FileDescriptorProto

	// descInfo has some pre-processed information for the proto descriptors,
	// so that looking up message or method configs is easier
	descInfo pbinfo.Info
	gapic    GAPICConfig

	// sampleConfig is the sample configurations, in v1.2 schema
	sampleConfig schema_v1p2.SampleConfig

	// clientPkg is the go package of the generated Gapic client
	clientPkg pbinfo.ImportSpec

	// if set to true, the generator will not format the generated code
	nofmt bool

	pt      printer.P
	imports map[pbinfo.ImportSpec]bool
	Outputs map[string][]byte
}

// readSampleConfigFiles loads sample configs from local files and stores
// them in generator.sampleConfig.
func (gen *generator) readSampleConfigFiles(paths []string) error {
	for _, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			return errors.E(err, "cannot open sample config file: %s", path)
		}
		decoder := yaml.NewDecoder(f)
		config, err := readSampleConfig(decoder, path)
		if err != nil {
			return err
		}
		gen.sampleConfig.Samples = append(gen.sampleConfig.Samples, config.Samples...)
	}
	return nil
}

// readGapicConfigFile loads gapic config from a local file and stores it in
// generator.gapic.
func (gen *generator) readGapicConfigFile(gapicFname string) error {
	// ignore gapicFname if unspecified
	if gapicFname == "" {
		return nil
	}
	f, err := os.Open(gapicFname)
	if err != nil {
		return errors.E(err, "cannot read GAPIC config file: %q", gapicFname)
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(&gen.gapic); err != nil {
		return errors.E(err, "error reading GAPIC config file: %q", gapicFname)
	}
	return nil
}

func (g *generator) reset() {
	g.pt.Reset()
	for imp := range g.imports {
		delete(g.imports, imp)
	}
}

// commit writes the internal data representation of the generator to bytes.
// It does the following:
// 1) dumps license with the given year
// 2) dumps package and imports
// 3) dumps generated code in generator.printer
// 4) formats the code if gofmt is true
func (g *generator) commit(gofmt bool, year int) ([]byte, error) {
	// We'll gofmt unless user asks us not to, so no need to think too hard about sorting
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
			return nil, errors.E(err, "syntax error, run with -nofmt to find out why")
		}
		b = b2
	}
	return b, nil
}

// disambiguateSampleIDs assigns unique sample IDs to each sample config.
// User-specified sample IDs (or region tags) are used if they are unique.
// Otherwise, the generator will generate a unique ID from the content
// of the sample, or fail if it is unable to do so (when multiple samples
// have identical contents).
func (g *generator) disambiguateSampleIDs() error {
	idCount := make(map[string]int)
	hashes := make(map[string]bool)
	samples := g.sampleConfig.Samples
	for i := range samples {
		// default ID to region tag
		if samples[i].ID == "" {
			samples[i].ID = samples[i].RegionTag
		}
		if samples[i].ID != "" {
			idCount[samples[i].ID]++
		}
	}

	for i := range samples {
		// Generate a sample ID when the user-provided ID is empty or not unique
		if samples[i].ID == "" || idCount[samples[i].ID] > 1 {
			jsonStr, err := json.Marshal(samples[i])
			if err != nil {
				return err
			}
			checkSum := sha256.Sum256([]byte(jsonStr))
			encodedStr := base32.StdEncoding.EncodeToString(checkSum[:])
			suffix := string([]rune(encodedStr)[0:8])
			if _, found := hashes[suffix]; found {
				return errors.E(nil, "unable to get a unique hash: multiple samples with identical contents")
			}
			hashes[suffix] = true
			samples[i].ID += suffix
		}
	}
	return nil
}

// genSample generates one sample from sample config and gapic config.
// TODO(hzyi): this method is getting long. Split it up.
func (g *generator) genSample(sampConf schema_v1p2.Sample, methConf GAPICMethod) error {
	ifaceName := sampConf.Service

	// Preparation
	g.imports[g.clientPkg] = true
	serv := g.descInfo.Serv["."+ifaceName]
	if serv == nil {
		return errors.E(nil, "can't find service: %q", ifaceName)
	}

	// We still need method-level GAPIC config for LRO types and Resource name bindings
	// until we start to parse proto annotations
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

	initInfo, err := g.getInitInfo(inType, methConf, sampConf.Request)
	if err != nil {
		return err
	}

	p := g.pt.Printf

	// Region tag
	argStr, err := argListStr(initInfo, g)
	if err != nil {
		return err
	}
	p("")
	p("// [START %s]", sampConf.RegionTag)
	p("")

	// comments above the sample function
	requiresNewLine := false

	writeCommentLines := func(comment string) {
		comment = wrapComment(strings.TrimSpace(comment))
		writeComment(comment, nil, g)
	}

	if sampConf.Description != "" {
		writeCommentLines(fmt.Sprintf("sample%s: %s", meth.GetName(), sampConf.Description))
		requiresNewLine = true
	}

	for i, argName := range initInfo.argNames {
		comment := initInfo.argTrees[i].comment
		if comment == "" {
			continue
		}
		if requiresNewLine {
			writeCommentLines("")
		}
		requiresNewLine = true
		writeCommentLines(fmt.Sprintf("%s: %s", argName, comment))
	}

	// function signature and client initialization
	p("func sample%s(%s) error {", meth.GetName(), argStr)
	p("  ctx := context.Background()")
	p("  c, err := %s.New%sClient(ctx)", g.clientPkg.Name, pbinfo.ReduceServName(serv.GetName(), g.clientPkg.Name))
	p("  if err != nil {")
	p("    return err")
	p("  }")
	p("")
	g.imports[pbinfo.ImportSpec{Path: "context"}] = true

	// Set up the request object
	if err := g.handleRequest(initInfo); err != nil {
		return err
	}

	// Make the RPC call and handle output
	if meth.GetOutputType() == ".google.protobuf.Empty" {
		err = g.emptyOut(meth, sampConf.Response)
	} else if meth.GetOutputType() == ".google.longrunning.Operation" {
		err = g.lro(meth, methConf, sampConf.Response)
	} else if meth.GetServerStreaming() || meth.GetClientStreaming() {
		// TODO(hzyi): github.com/googleapis/gapic-generator-go/issues/177
		err = errors.E(nil, "streaming methods not supported yet")
	} else if pf, err2 := pagingField(g.descInfo, meth); err2 != nil {
		err = errors.E(err2, "can't determine whether method is paging")
	} else if pf != nil {
		err = g.paging(meth, pf, sampConf.Response)
	} else {
		err = g.unary(meth, sampConf.Response)
	}

	if err != nil {
		return err
	}

	p("return nil")
	p("}")
	p("")
	p("// [END %s]", sampConf.RegionTag)
	p("")

	// main
	if err := writeMain(g, initInfo.argNames, initInfo.flagNames, initInfo.argTrees, meth.GetName()); err != nil {
		return err
	}

	return nil
}

// getInitInfo extracts information from request configs to prepare for generation of the request setup part of the sample.
func (g *generator) getInitInfo(inType pbinfo.ProtoType, methConf GAPICMethod, fieldConfs []schema_v1p2.RequestConfig) (initInfo, error) {
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
	for _, fieldConf := range fieldConfs {
		if err := itree.parseInit(fieldConf.Field, fieldConf.Value, fieldConf.Comment, g.descInfo); err != nil {
			return initInfo{}, errors.E(err, "can't set default value: %s=%s", fieldConf.Field, fieldConf.Value)
		}
	}

	// Some parts of request object are from arguments.
	for _, fieldConf := range fieldConfs {
		isInputParam := fieldConf.InputParameter != ""
		if fieldConf.ValueIsFile {
			var varName string
			var err error
			if isInputParam {
				varName = snakeToCamel(fieldConf.InputParameter) + fileContentSuffix
			} else {
				varName, err = fileVarName(fieldConf.Field)
				if err != nil {
					return initInfo{}, errors.E(err, "can't determine variable name to store bytes from local file")
				}
			}

			subTree, err := itree.parseSampleArgPath(
				fieldConf.Field,
				g.descInfo,
				varName,
			)
			if err != nil {
				return initInfo{}, errors.E(err, "can't set sample function argument: %q", fieldConf.Field)
			}
			if subTree.typ.prim != descriptor.FieldDescriptorProto_TYPE_BYTES {
				return initInfo{}, errors.E(nil, "can only assign file contents to bytes field")
			}
			subTree.typ.prim = descriptor.FieldDescriptorProto_TYPE_STRING
			subTree.typ.valFmt = nil
			if subTree.leafVal == "" {
				return initInfo{}, errors.E(nil, "default value not given: %q", fieldConf.Field)
			}
			fileName := subTree.leafVal
			if isInputParam {
				fileName = snakeToCamel(fieldConf.InputParameter)
				argNames = append(argNames, fileName)
				flagNames = append(flagNames, fieldConf.InputParameter)
				argTrees = append(argTrees, subTree)
			}
			files = append(files, &fileInfo{fileName, varName, fieldConf.Comment})
			continue
		}

		if !isInputParam {
			continue
		}
		name := snakeToCamel(fieldConf.InputParameter)
		subTree, err := itree.parseSampleArgPath(fieldConf.Field, g.descInfo, name)
		if err != nil {
			return initInfo{}, errors.E(err, "can't set sample function argument: %q", fieldConf.Field)
		}

		argNames = append(argNames, name)
		flagNames = append(flagNames, fieldConf.InputParameter)
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

// argListStr returns a comma-separated string of the list of arguments that the sample function takes.
func argListStr(init initInfo, g *generator) (string, error) {
	if len(init.argNames) > 0 {
		var sb strings.Builder
		for i, name := range init.argNames {
			typ, err := g.getGoTypeName(init.argTrees[i].typ)
			if err != nil {
				return "", err
			}
			fmt.Fprintf(&sb, ", %s %s", name, typ)
		}
		return sb.String()[2:], nil
	}
	return "", nil
}

func (g *generator) handleRequest(initInfo initInfo) error {
	var buf bytes.Buffer

	for i, name := range initInfo.argNames {
		fmt.Fprintf(&buf, "%s := ", name)
		if err := initInfo.argTrees[i].Print(&buf, g); err != nil {
			return errors.E(err, "can't initialize parameter: %s", name)
		}
		buf.WriteByte('\n')
	}
	if err := prependLines(&buf, "// ", false); err != nil {
		return err
	}

	for _, info := range initInfo.files {
		handleReadFile(info, &buf, g)
	}

	buf.WriteString("req := ")
	if err := initInfo.reqTree.Print(&buf, g); err != nil {
		return errors.E(err, "can't initialize request object")
	}
	buf.WriteByte('\n')
	if err := prependLines(&buf, "\t", true); err != nil {
		return err
	}

	if _, err := buf.WriteTo(g.pt.Writer()); err != nil {
		return err
	}
	return nil
}

func (g *generator) unary(meth *descriptor.MethodDescriptorProto, respConfs []schema_v1p2.ResponseConfig) error {
	p := g.pt.Printf

	p("resp, err := c.%s(ctx, req)", meth.GetName())
	p("if err != nil {")
	p("  return err")
	p("}")
	p("")

	return g.handleOut(meth, respConfs, &initType{desc: g.descInfo.Type[meth.GetOutputType()]})
}

func (g *generator) emptyOut(meth *descriptor.MethodDescriptorProto, respConfs []schema_v1p2.ResponseConfig) error {
	p := g.pt.Printf

	p("if err := c.%s(ctx, req); err != nil {", meth.GetName())
	p("  return err")
	p("}")
	p("")

	return g.handleOut(meth, respConfs, nil)
}

func (g *generator) paging(meth *descriptor.MethodDescriptorProto, pf *descriptor.FieldDescriptorProto, respConfs []schema_v1p2.ResponseConfig) error {
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

	err := g.handleOut(meth, respConfs, &typ)

	p("}")
	p("")
	g.imports[pbinfo.ImportSpec{Path: "google.golang.org/api/iterator"}] = true

	return err
}

func (g *generator) lro(meth *descriptor.MethodDescriptorProto, methConf GAPICMethod, respConfs []schema_v1p2.ResponseConfig) error {
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
	return g.handleOut(meth, respConfs, &typ)
}

func (g *generator) handleOut(meth *descriptor.MethodDescriptorProto, respConfs []schema_v1p2.ResponseConfig, respTyp *initType) error {
	st := newSymTab(nil)
	if respTyp != nil {
		st.universe["$resp"] = true
		st.scope["$resp"] = *respTyp
	}

	for _, out := range respConfs {
		if err := writeOutputSpec(out, st, g); err != nil {
			return errors.E(err, "cannot write output handling code")
		}
	}

	if respTyp != nil && len(respConfs) == 0 {
		g.pt.Printf("fmt.Println(resp)")
		g.imports[pbinfo.ImportSpec{Path: "fmt"}] = true
	}
	return nil
}

// readSampleConfig loads sample configs from a file.
func readSampleConfig(decoder *yaml.Decoder, fname string) (schema_v1p2.SampleConfig, error) {
	var config schema_v1p2.SampleConfig
	for {
		var sc schema_v1p2.SampleConfig
		err := decoder.Decode(&sc)
		if err != nil && err.Error() == "EOF" {
			// last YAML document, all done
			break
		}
		if err != nil {
			return config, err
		}
		if sc.Type != expectedSampleConfigType {
			// ignore non sample config YAMLs
			continue
		}
		if sc.Version != expectedSampleConfigVersion {
			// ignore unsupported versions
			continue
		}
		config.Samples = append(config.Samples, sc.Samples...)
	}
	if len(config.Samples) == 0 {
		return config, errors.E(nil, "Found no valid sample config in %q", fname)
	}
	return config, nil
}

// prependLines adds prefix to every line in b. A line is defined as a possibly empty run
// of non-newlines terminated by a newline character.
// If b doesn't end with a newline, prependLines returns an error.
// If skipEmptyLine is true, prependLines does not prepend prefix to empty lines.
func prependLines(b *bytes.Buffer, prefix string, skipEmptyLine bool) error {
	if b.Len() == 0 {
		return nil
	}
	if b.Bytes()[b.Len()-1] != '\n' {
		return errors.E(nil, "prependLines: must end with newline")
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
	return nil
}

// wrapComment wraps comment at 100 characters, iff comment does not contain any newline characters
// and comment has more than 110 characters.
//
// comment cannot have leading or trailing white spaces.
func wrapComment(comment string) string {
	if strings.ContainsRune(comment, '\n') || len(comment) < 110 {
		return comment
	}

	var output strings.Builder
	s := 0
	prev := -1

	for true {
		p := strings.IndexByte(comment[prev+1:], ' ')
		// we reached the end of comment
		if p < 0 {
			output.WriteString(comment[s:])
			return output.String()
		}

		p = p + prev + 1
		if p-s > 100 {
			// a single word has more than 100 characters
			if prev == s-1 {
				prev = p
			}
			// break the line at prev
			output.WriteString(comment[s:prev])
			output.WriteByte('\n')
			s = prev + 1
		}
		prev = p
	}

	return ""
}
