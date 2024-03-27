// Copyright 2023 Google LLC
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

package gengapic

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"github.com/googleapis/gapic-generator-go/internal/snippets"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

// newSnippetsMetadata initializes the model that will collect snippet metadata.
// Does nothing and returns nil if opts.omitSnippets is true.
func (g *generator) newSnippetsMetadata(protoPkg string) *snippets.SnippetMetadata {
	if g.opts.omitSnippets {
		return nil
	}
	return snippets.NewMetadata(protoPkg, g.metadata.LibraryPackage, g.opts.pkgName)
}

// addSnippetsMetadataDoc sets the documentation for a method in the snippet metadata.
// Does nothing and returns nil if opts.omitSnippets is true, or if the streaming type is not supported in example.go.
func (g *generator) addSnippetsMetadataDoc(m *descriptorpb.MethodDescriptorProto, servName, doc string) {
	if g.opts.omitSnippets || m.GetClientStreaming() != m.GetServerStreaming() {
		// TODO(chrisdsmith): implement streaming examples correctly, see example.go TODO(pongad).
		return
	}
	g.snippetMetadata.UpdateMethodDoc(servName, m.GetName(), doc)
}

// addSnippetsMetadataParams sets the parameters for a method in the snippet metadata.
// Does nothing and returns nil if opts.omitSnippets is true, or if the streaming type is not supported in example.go.
func (g *generator) addSnippetsMetadataParams(m *descriptorpb.MethodDescriptorProto, servName, requestType string) {
	if g.opts.omitSnippets || m.GetClientStreaming() != m.GetServerStreaming() {
		// TODO(chrisdsmith): implement streaming examples correctly, see example.go TODO(pongad).
		return
	}
	g.snippetMetadata.AddParams(servName, m.GetName(), requestType)
}

// addSnippetsMetadataResult sets the result type for a method in the snippet metadata.
// Does nothing and returns nil if opts.omitSnippets is true, or if the streaming type is not supported in example.go.
func (g *generator) addSnippetsMetadataResult(m *descriptorpb.MethodDescriptorProto, servName, resultType string) {
	if g.opts.omitSnippets || m.GetClientStreaming() != m.GetServerStreaming() {
		// TODO(chrisdsmith): implement streaming examples correctly, see example.go TODO(pongad).
		return
	}
	g.snippetMetadata.UpdateMethodResult(servName, m.GetName(), resultType)
}

// genAndCommitSnippets generates and commits a snippet file for each method in a client.
// Does nothing and returns nil if opts.omitSnippets is true.
func (g *generator) genAndCommitSnippets(s *descriptorpb.ServiceDescriptorProto) error {
	if g.opts.omitSnippets {
		return nil
	}
	defaultHost := proto.GetExtension(s.Options, annotations.E_DefaultHost).(string)
	g.snippetMetadata.AddService(s.GetName(), defaultHost)
	methods := append(s.GetMethod(), g.getMixinMethods()...)
	for _, m := range methods {
		if m.GetClientStreaming() != m.GetServerStreaming() {
			// TODO(chrisdsmith): implement streaming examples correctly, see example.go TODOs.
			continue
		}
		// For each method, reset the generator in order to write a
		// separate main.go snippet file.
		g.reset()
		if err := g.genSnippetFile(s, m); err != nil {
			return err
		}
		g.imports[pbinfo.ImportSpec{Name: g.opts.pkgName, Path: g.opts.pkgPath}] = true
		// Use the client short name in this filepath.
		// E.g. the client for LoggingServiceV2 is just "Client".
		clientName := pbinfo.ReduceServName(s.GetName(), g.opts.pkgName) + "Client"
		// Get the original proto namespace for the method (different from `s` only for mixins).
		f := g.descInfo.ParentFile[m]
		// Get the original proto service for the method (different from `s` only for mixins).
		methodServ := (g.descInfo.ParentElement[m]).(*descriptorpb.ServiceDescriptorProto)
		lineCount := g.commit(filepath.Join(g.snippetsOutDir(), clientName, m.GetName(), "main.go"), "main")
		g.snippetMetadata.AddMethod(s.GetName(), m.GetName(), f.GetPackage(), methodServ.GetName(), lineCount-1)
	}
	return nil
}

// genSnippetFile generates a single RPC snippet by leveraging exampleMethodBody in gengapic/example.go.
func (g *generator) genSnippetFile(s *descriptorpb.ServiceDescriptorProto, m *descriptorpb.MethodDescriptorProto) error {
	regionTag := g.snippetMetadata.RegionTag(s.GetName(), m.GetName())
	g.headerComment(fmt.Sprintf("[START %s]", regionTag))
	pkgName := g.opts.pkgName
	servName := pbinfo.ReduceServName(s.GetName(), pkgName)

	p := g.printf
	p("func main() {")
	if err := g.exampleMethodBody(pkgName, servName, m); err != nil {
		return err
	}
	p("}")
	p("")
	g.comment(fmt.Sprintf("[END %s]\n", regionTag))
	return nil
}

// genAndCommitSnippetMetadata generates and commits the snippet metadata to the generator response.
// Does nothing and returns nil if opts.omitSnippets is true.
func (g *generator) genAndCommitSnippetMetadata(protoPkg string) error {
	if g.opts.omitSnippets {
		return nil
	}
	g.reset()
	json, err := g.snippetMetadata.ToMetadataJSON()
	if err != nil {
		return err
	}
	file := filepath.Join(g.snippetsOutDir(), fmt.Sprintf("snippet_metadata.%s.json", protoPkg))
	g.resp.File = append(g.resp.File, &pluginpb.CodeGeneratorResponse_File{
		Name:    proto.String(file),
		Content: proto.String(string(json[:])),
	})
	return nil
}

func (g *generator) snippetsOutDir() string {
	if strings.Contains(g.opts.pkgPath, "cloud.google.com/go/") {
		// Write snippet metadata at the top level of the google-cloud-go namespace, not at the client package.
		// This matches the destination directory structure in google-cloud-go.
		pkg := strings.TrimPrefix(g.opts.pkgPath, "cloud.google.com/go/")
		return filepath.Join("cloud.google.com/go", "internal", "generated", "snippets", filepath.FromSlash(pkg))
	}
	return filepath.Join(g.opts.outDir, "internal", "snippets")
}
