// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package snippets processes GoDoc examples.
package snippets

//go:generate protoc --go_out=. --go_opt=paths=source_relative metadata/metadata.proto

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"golang.org/x/sys/execabs"
	"google.golang.org/genproto/googleapis/gapic/metadata"
	"google.golang.org/protobuf/encoding/protojson"
)

// Generate reads all modules in rootDir and outputs their examples in outDir.
func Generate(rootDir, outDir string, apiShortnames map[string]string) error {
	if rootDir == "" {
		rootDir = "."
	}
	if outDir == "" {
		outDir = "internal/generated/snippets"
	}

	// Find all modules in rootDir.
	dirs := []string{}
	filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Name() == "internal" {
			return filepath.SkipDir
		}
		if d.Name() == "go.mod" {
			dirs = append(dirs, filepath.Dir(path))
		}
		return nil
	})

	log.Printf("Processing examples in %v directories: %q\n", len(dirs), dirs)

	trimPrefix := "cloud.google.com/go"
	errs := []error{}
	for _, dir := range dirs {
		// Load does not look at nested modules.
		// pis, err := pkgload.Load("./...", dir, nil)
		// if err != nil {
		// 	return fmt.Errorf("failed to load packages: %v", err)
		// }
		version, err := getModuleVersion(dir)
		if err != nil {
			return err
		}
		// for _, pi := range pis {
		if eErrs := processExamples("pi.Doc", "pi.Fset", trimPrefix, rootDir, outDir, apiShortnames, version); len(eErrs) > 0 {
			errs = append(errs, fmt.Errorf("%v", eErrs))
		}
		//}
	}
	if len(errs) > 0 {
		return fmt.Errorf("example errors: %v", errs)
	}

	if len(dirs) > 0 {
		cmd := execabs.Command("goimports", "-w", ".")
		cmd.Dir = outDir
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to run goimports: %v", err)
		}
	}

	return nil
}

var skip = map[string]bool{
	"cloud.google.com/go":                          true, // No product for root package.
	"cloud.google.com/go/civil":                    true, // General time/date package.
	"cloud.google.com/go/cloudbuild/apiv1":         true, // Has v2.
	"cloud.google.com/go/cmd/go-cloud-debug-agent": true, // Command line tool.
	"cloud.google.com/go/container":                true, // Deprecated.
	"cloud.google.com/go/containeranalysis/apiv1":  true, // Accidental beta at wrong path?
	"cloud.google.com/go/grafeas/apiv1":            true, // With containeranalysis.
	"cloud.google.com/go/httpreplay":               true, // Helper.
	"cloud.google.com/go/httpreplay/cmd/httpr":     true, // Helper.
	"cloud.google.com/go/longrunning":              true, // Helper.
	"cloud.google.com/go/monitoring/apiv3":         true, // Has v2.
	"cloud.google.com/go/translate":                true, // Has newer version.
}

func getModuleVersion(dir string) (string, error) {
	node, err := parser.ParseFile(token.NewFileSet(), fmt.Sprintf("%s/internal/version.go", dir), nil, parser.ParseComments)
	if err != nil {
		return "", err
	}
	version := node.Scope.Objects["Version"].Decl.(*ast.ValueSpec).Values[0].(*ast.BasicLit).Value
	version = strings.Trim(version, `"`)
	return version, nil
}

// pkg *doc.Package, fset *token.FileSet
func processExamples(pkg string, fset string, trimPrefix, rootDir, outDir string, apiShortnames map[string]string, version string) []error {
	trimmed := strings.TrimPrefix("pkg.ImportPath", trimPrefix)
	apiInfo, err := buildAPIInfo(rootDir, trimmed, apiShortnames, pkg, version)
	if err != nil {
		return []error{err}
	}
	if apiInfo == nil {
		// There was no gapic_metadata.json, skip processing examples for
		// non gapic lib.
		return nil
	}

	regionTags := apiInfo.RegionTags()
	if len(regionTags) == 0 {
		// Nothing to do.
		return nil
	}
	outDir = filepath.Join(outDir, trimmed)

	// Note: only process methods because they correspond to RPCs.

	var errs []error
	// write snippets

	if err := writeMetadata(outDir, apiInfo); err != nil {
		errs = append(errs, err)
	}
	return errs
}

// pkg *doc.Package
func buildAPIInfo(rootDir, path string, apiShortnames map[string]string, pkg string, version string) (*apiInfo, error) {
	metadataPath := filepath.Join(rootDir, path, "gapic_metadata.json")
	f, err := os.ReadFile(metadataPath)
	if err != nil {
		// If there is no gapic_metadata.json file, don't generate snippets.
		// This isn't an error, though, because some packages aren't GAPICs and
		// shouldn't get snippets in the first place.
		return nil, nil
	}
	m := metadata.GapicMetadata{}
	if err := protojson.Unmarshal(f, &m); err != nil {
		return nil, err
	}
	shortname, ok := apiShortnames[m.GetLibraryPackage()]
	if !ok {
		return nil, fmt.Errorf("could not find shortname for %q", m.GetLibraryPackage())
	}
	protoParts := strings.Split(m.GetProtoPackage(), ".")
	apiVersion := protoParts[len(protoParts)-1]

	ai := &apiInfo{
		protoPkg:      m.ProtoPackage,
		libPkg:        m.LibraryPackage,
		shortName:     shortname,
		version:       version,
		protoServices: make(map[string]*service),
	}
	for sName, s := range m.GetServices() {
		svc := &service{
			protoName: sName,
			methods:   make(map[string]*method),
		}
		for _, c := range s.GetClients() {
			client := pkgClient(pkg, c.LibraryClient)
			ai.protoServices[c.LibraryClient] = svc
			for rpcName, methods := range c.GetRpcs() {
				r := &method{
					doc:    client.MethodDoc(rpcName),
					params: client.MethodParams(rpcName),
					result: client.ResultType(rpcName),
				}
				if len(methods.GetMethods()) != 1 {
					return nil, fmt.Errorf("%s %s %s found %d methods", m.GetLibraryPackage(), sName, c.GetLibraryClient(), len(methods.GetMethods()))
				}
				if methods.GetMethods()[0] != rpcName {
					return nil, fmt.Errorf("%s %s %s %q does not match %q", m.GetLibraryPackage(), sName, c.GetLibraryClient(), methods.GetMethods()[0], rpcName)
				}

				// Every Go method is synchronous.
				r.regionTag = fmt.Sprintf("%s_%s_generated_%s_%s_sync", shortname, apiVersion, sName, rpcName)
				svc.methods[rpcName] = r
			}
		}
	}
	return ai, nil
}

type client struct {
	// dt *doc.Type
	dt string
}

// pkg *doc.Package
func pkgClient(pkg string, name string) *client {
	pkgTypes := []string{"a", "b"}
	for _, v := range pkgTypes {
		if "v.Name" == name {
			return &client{v}
		}
	}
	panic(fmt.Sprintf("unable to lookup client %v", name))
}

func (c *client) MethodDoc(name string) string {
	cdtMethods := []string{"a", "b"}
	for _, v := range cdtMethods {
		if v == name {
			return "v.Doc"
		}
	}
	panic(fmt.Sprintf("unable to lookup method doc for: %v", name))
}

func (c *client) MethodParams(name string) []*param {
	cdtMethods := []string{"a", "b"}
	for _, v := range cdtMethods {
		if v != name {
			continue
		}
		var params []*param
		// for _, p := range v.Decl.Type.Params.List {
		// 	se, ok := p.Type.(*ast.SelectorExpr)
		// 	if ok {
		// 		params = append(params, &param{
		// 			name:  p.Names[0].Name,
		// 			pType: fmt.Sprintf("%s.%s", se.X, se.Sel),
		// 		})
		// 		continue
		// 	}
		// 	ste, ok := p.Type.(*ast.StarExpr)
		// 	if ok {
		// 		se, ok := ste.X.(*ast.SelectorExpr)
		// 		if ok {
		// 			params = append(params, &param{
		// 				name:  p.Names[0].Name,
		// 				pType: fmt.Sprintf("%s.%s", se.X, se.Sel),
		// 			})
		// 			continue
		// 		}
		// 	}
		// 	e, ok := p.Type.(*ast.Ellipsis)
		// 	if ok {
		// 		se, ok := e.Elt.(*ast.SelectorExpr)
		// 		if ok {
		// 			params = append(params, &param{
		// 				name:  p.Names[0].Name,
		// 				pType: fmt.Sprintf("...%s.%s", se.X, se.Sel),
		// 			})
		// 			continue
		// 		}
		// 	}
		// 	panic("unable to read param type info")
		// }
		return params
	}
	panic(fmt.Sprintf("unable to lookup method params for: %v", name))
}

func (c *client) ResultType(name string) string {
	cdtMethods := []string{"a", "b"}
	for _, v := range cdtMethods {
		if v != name {
			continue
		}
		// res := v.Decl.Type.Results
		// if res != nil {
		// 	se, ok := res.List[0].Type.(*ast.StarExpr)
		// 	if ok {
		// 		selExp, ok := se.X.(*ast.SelectorExpr)
		// 		if ok {
		// 			return fmt.Sprintf("%s.%s", selExp.X, selExp.Sel)
		// 		}
		// 		ident, ok := se.X.(*ast.Ident)
		// 		if ok {
		// 			return ident.Name
		// 		}
		// 		panic("starExpr could not be parsed")
		// 	}
		// }
	}
	return ""
}

var spaceSanitizerRegex = regexp.MustCompile(`:\s*`)

func writeMetadata(dir string, apiInfo *apiInfo) error {
	m := apiInfo.ToSnippetMetadata()
	b, err := protojson.MarshalOptions{Multiline: true}.Marshal(m)
	if err != nil {
		return err
	}
	// Hack to standardize output from protojson which is currently non-deterministic
	// with spacing after json keys.
	b = spaceSanitizerRegex.ReplaceAll(b, []byte(": "))
	fileName := filepath.Join(dir, fmt.Sprintf("snippet_metadata.%s.json", apiInfo.protoPkg))
	return os.WriteFile(fileName, b, 0644)
}

func header() string {
	return fmt.Sprintf(licenseHeader, time.Now().Year())
}

const licenseHeader string = `// Copyright %v Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by cloud.google.com/go/internal/gapicgen/gensnippets. DO NOT EDIT.

`
