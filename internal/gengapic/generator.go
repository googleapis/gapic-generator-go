// Copyright 2021 Google LLC
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
	"strings"
	"time"

	"github.com/googleapis/gapic-generator-go/internal/license"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"github.com/googleapis/gapic-generator-go/internal/printer"
	"github.com/googleapis/gapic-generator-go/internal/snippets"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/genproto/googleapis/gapic/metadata"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoiface"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

type generator struct {
	pt printer.P

	// Protobuf descriptor properties
	descInfo pbinfo.Info

	// Maps proto elements to their comments
	comments map[protoiface.MessageV1]string

	resp pluginpb.CodeGeneratorResponse

	// Comments to appear after the license header and before the package declaration.
	headerComments printer.P

	imports map[pbinfo.ImportSpec]bool

	// Human-readable name of the API used in docs
	apiName string

	// Auxiliary types to be generated in the package
	aux *auxTypes

	// Options for the generator determining module names, transports,
	// config file paths, etc.  This should be treated as immutable once
	// configured.
	cfg *generatorConfig

	// GapicMetadata for recording proto-to-code mappings in a
	// gapic_metadata.json file.
	metadata *metadata.GapicMetadata

	// Model for capturing snippet details in a snippet_metadata.*.json file.
	snippetMetadata *snippets.SnippetMetadata

	mixins mixins

	hasIAMPolicyOverrides bool

	// customOpServices is a map of service descriptors with methods that create custom operations
	// to the service descriptors of the custom operation services that manage those custom operation instances.
	customOpServices map[*descriptorpb.ServiceDescriptorProto]*descriptorpb.ServiceDescriptorProto

	// learned vocabulary for heuristic path templates
	vocabulary map[string]bool

	// clientProtoPkg is the proto package of the service currently being generated.
	clientProtoPkg string
}

func newGenerator(req *pluginpb.CodeGeneratorRequest) (*generator, error) {
	g := &generator{
		metadata: &metadata.GapicMetadata{
			Schema:   "1.0",
			Language: "go",
			Comment:  "This file maps proto services/RPCs to the corresponding library clients/methods.",
			Services: make(map[string]*metadata.GapicMetadata_ServiceForTransport),
		},
		mixins:           make(mixins),
		comments:         map[protoiface.MessageV1]string{},
		imports:          map[pbinfo.ImportSpec]bool{},
		customOpServices: map[*descriptorpb.ServiceDescriptorProto]*descriptorpb.ServiceDescriptorProto{},
		aux: &auxTypes{
			iters:           map[string]*iterType{},
			methodToWrapper: map[*descriptorpb.MethodDescriptorProto]operationWrapper{},
			opWrappers:      map[string]operationWrapper{},
		},
	}

	// Build and validate the immutable configuration from the CodeGeneratorRequest plugin args.
	cfg, err := configFromRequest(req.Parameter)
	if err != nil {
		return nil, err
	}

	// attach config to generator.
	g.cfg = cfg

	var methods []*descriptorpb.MethodDescriptorProto
	for _, f := range req.GetProtoFile() {
		for _, s := range f.GetService() {
			methods = append(methods, s.GetMethod()...)
		}
	}
	g.vocabulary = buildHeuristicVocabulary(methods)

	files := req.GetProtoFile()
	files = append(files, wellKnownTypeFiles...)

	g.collectMixins()
	files = append(files, g.getMixinFiles()...)

	g.descInfo = pbinfo.Of(files)
	if len(g.cfg.pkgOverrides) > 0 {
		g.descInfo.PkgOverrides = g.cfg.pkgOverrides
	}

	for _, f := range files {
		for _, loc := range f.GetSourceCodeInfo().GetLocation() {
			if loc.LeadingComments == nil {
				continue
			}

			// p is an array with format [f1, i1, f2, i2, ...]
			// - f1 refers to the protobuf field tag
			// - if field refer to by f1 is a slice, i1 refers to an element in that slice
			// - f2 and i2 works recursively.
			// So, [6, x] refers to the xth service defined in the file,
			// since the field tag of Service is 6.
			// [6, x, 2, y] refers to the yth method in that service,
			// since the field tag of Method is 2.
			p := loc.Path
			switch {
			case len(p) == 2 && p[0] == 6:
				g.comments[f.Service[p[1]]] = *loc.LeadingComments
			case len(p) == 4 && p[0] == 6 && p[2] == 2:
				g.comments[f.Service[p[1]].Method[p[3]]] = *loc.LeadingComments
			}
		}
	}

	return g, nil
}

// featureEnabled is a simple boolean checker for probing if a given feature has been enabled.
func (g *generator) featureEnabled(f featureID) bool {
	if g.cfg == nil {
		return false
	}
	if _, ok := g.cfg.featureEnablement[f]; ok {
		return true
	}
	return false
}

// printf formatted-prints to sb, using the print syntax from fmt package.
//
// It automatically keeps track of indentation caused by curly-braces.
// To make nested blocks easier to write elsewhere in the code,
// leading and trailing whitespaces in s are ignored.
// These spaces are for humans reading the code, not machines.
//
// Currently it's not terribly difficult to confuse the auto-indenter.
// To fix-up, manipulate g.in or write to g.sb directly.
func (g *generator) printf(s string, a ...interface{}) {
	g.pt.Printf(s, a...)
}

// TODO(chrisdsmith): Add generator_test.go with TestCommit

func (g *generator) commit(fileName, pkgName string) int {
	return g.commitWithBuildTag(fileName, pkgName, "")
}

// commit adds header, etc to current pt and returns the line length of the
// final file output.
func (g *generator) commitWithBuildTag(fileName, pkgName, buildTag string) int {
	var header strings.Builder
	fmt.Fprintf(&header, license.Apache, time.Now().Year())
	header.WriteString(g.headerComments.String() + "\n")
	if buildTag != "" {
		fmt.Fprintf(&header, "//go:build %s\n\n", buildTag)
	}
	fmt.Fprintf(&header, "package %s\n\n", pkgName)

	var imps []pbinfo.ImportSpec
	dupCheck := map[string]bool{}
	for imp := range g.imports {
		// TODO(codyoss): This if can be removed once the public protos
		// have been migrated to their new package. This should be soon after this
		// code is merged.
		if imp.Path == "google.golang.org/genproto/googleapis/longrunning" {
			imp.Path = "cloud.google.com/go/longrunning/autogen/longrunningpb"
		}
		if imp.Path == "google.golang.org/genproto/googleapis/iam/v1" {
			imp.Path = "cloud.google.com/go/iam/apiv1/iampb"
		}
		if exists := dupCheck[imp.Path]; !exists {
			dupCheck[imp.Path] = true
			imps = append(imps, imp)
		}
	}
	impDiv := sortImports(imps)

	writeImp := func(is pbinfo.ImportSpec) {
		s := "\t%[2]q\n"
		if is.Name != "" {
			s = "\t%s %q\n"
		}
		fmt.Fprintf(&header, s, is.Name, is.Path)
	}

	header.WriteString("import (\n")
	for _, imp := range imps[:impDiv] {
		writeImp(imp)
	}
	if impDiv != 0 && impDiv != len(imps) {
		header.WriteByte('\n')
	}
	for _, imp := range imps[impDiv:] {
		writeImp(imp)
	}
	header.WriteString(")\n\n")
	lineCount := len(strings.Split(header.String(), "\n"))
	g.resp.File = append(g.resp.File, &pluginpb.CodeGeneratorResponse_File{
		Name:    &fileName,
		Content: proto.String(header.String()),
	})

	// Trim trailing newlines so we have only one.
	// NOTE(pongad): This might be an overkill since we have gofmt,
	// but the rest of the file already conforms to gofmt, so we might as well?
	body := g.pt.String()
	if !strings.HasSuffix(body, "\n") {
		body += "\n"
	}
	for i := len(body) - 1; i >= 0; i-- {
		if body[i] != '\n' {
			body = body[:i+2]
			break
		}
	}

	g.resp.File = append(g.resp.File, &pluginpb.CodeGeneratorResponse_File{
		Content: proto.String(body),
	})

	return lineCount + len(strings.Split(body, "\n"))
}

func (g *generator) reset() {
	g.pt.Reset()
	g.headerComments.Reset()
	for k := range g.imports {
		delete(g.imports, k)
	}
}

// fqn recursively builds the fully qualified proto element name,
// but omits the leading ".". For example, google.foo.v1.FooMessage.
func (g *generator) fqn(p pbinfo.ProtoType) string {
	// Base case. Use proto package instead of relative file name.
	if f, isFile := p.(*descriptorpb.FileDescriptorProto); isFile {
		return f.GetPackage()
	}

	parent := g.descInfo.ParentElement[p]
	if parent == nil {
		parent = g.descInfo.ParentFile[p]
	}
	return fmt.Sprintf("%s.%s", g.fqn(parent), p.GetName())
}

func (g *generator) nestedName(nested pbinfo.ProtoType) string {
	name := nested.GetName()

	parent, hasParent := g.descInfo.ParentElement[nested]
	for hasParent {
		name = fmt.Sprintf("%s_%s", parent.GetName(), name)
		parent, hasParent = g.descInfo.ParentElement[parent]
	}

	return name
}

// autoPopulatedFields returns an array of FieldDescriptorProto pointers for the
// given MethodDescriptorProto that are specified for auto-population per the
// following restrictions:
//
// * The field is a top-level string field of a unary method's request message.
// * The field is not annotated with google.api.field_behavior = REQUIRED.
// * The field name is listed in google.api.publishing.method_settings.auto_populated_fields.
// * The field is annotated with google.api.field_info.format = UUID4.
func (g *generator) autoPopulatedFields(_ string, m *descriptorpb.MethodDescriptorProto) []*descriptorpb.FieldDescriptorProto {
	var apfs []string
	// Find the service config's AutoPopulatedFields entry by method name.
	mfqn := g.fqn(m)
	for _, s := range g.cfg.APIServiceConfig.GetPublishing().GetMethodSettings() {
		if s.GetSelector() == mfqn {
			apfs = s.AutoPopulatedFields
			break
		}
	}
	inType := g.descInfo.Type[m.GetInputType()].(*descriptorpb.DescriptorProto)
	var validated []*descriptorpb.FieldDescriptorProto
	for _, apf := range apfs {
		field := getField(inType, apf)
		// Do nothing and continue iterating unless all conditions above are met.
		switch {
		case field == nil:
		case field.GetType() != fieldTypeString:
		case isRequired(field):
		case proto.GetExtension(field.GetOptions(), annotations.E_FieldInfo).(*annotations.FieldInfo).GetFormat() == annotations.FieldInfo_UUID4:
			validated = append(validated, field)
		}
	}
	return validated
}

// getServiceNameOverride checks to see if the service has a defined service name override.
func (g *generator) getServiceNameOverride(s *descriptorpb.ServiceDescriptorProto) string {
	if g.cfg == nil || g.cfg.APIServiceConfig == nil || g.cfg.APIServiceConfig.GetPublishing() == nil {
		return ""
	}
	ls := g.cfg.APIServiceConfig.GetPublishing().GetLibrarySettings()

	protoPkg := g.descInfo.ParentFile[s].GetPackage()

	for _, setting := range ls {
		if setting.GetVersion() != "" && setting.GetVersion() != protoPkg {
			continue
		}
		if goSettings := setting.GetGoSettings(); goSettings != nil {
			if renamedServices := goSettings.GetRenamedServices(); renamedServices != nil {
				if v, ok := renamedServices[s.GetName()]; ok {
					return v
				}
			}
		}
	}

	return ""
}

// getSelectiveGapicGeneration returns the SelectiveGapicGeneration config if it is present and the feature is enabled.
func (g *generator) getSelectiveGapicGeneration(protoPkg string) *annotations.SelectiveGapicGeneration {
	if g.cfg == nil {
		return nil
	}
	if _, ok := g.cfg.featureEnablement[SelectiveGapicGenerationFeature]; !ok {
		return nil
	}
	if g.cfg.APIServiceConfig == nil || g.cfg.APIServiceConfig.GetPublishing() == nil {
		return nil
	}
	ls := g.cfg.APIServiceConfig.GetPublishing().GetLibrarySettings()

	for _, setting := range ls {
		if setting.GetVersion() != "" && setting.GetVersion() != protoPkg {
			continue
		}
		if goSettings := setting.GetGoSettings(); goSettings != nil && goSettings.GetCommon() != nil && goSettings.GetCommon().GetSelectiveGapicGeneration() != nil {
			return goSettings.GetCommon().GetSelectiveGapicGeneration()
		}
	}

	return nil
}

// isInternalMethod determines if a method should be generated as internal based on SelectiveGapicGeneration config.
func (g *generator) isInternalMethod(protoPkg, fqn string) bool {
	sgg := g.getSelectiveGapicGeneration(protoPkg)
	if sgg == nil || !sgg.GetGenerateOmittedAsInternal() {
		return false
	}
	for _, m := range sgg.GetMethods() {
		if m == fqn {
			return false
		}
	}
	return true
}

// isInternalService determines if a service should be treated as internal (i.e. has at least one internal method).
func (g *generator) isInternalService(s *descriptorpb.ServiceDescriptorProto) bool {
	protoPkg := g.descInfo.ParentFile[s].GetPackage()
	sgg := g.getSelectiveGapicGeneration(protoPkg)
	if sgg == nil || !sgg.GetGenerateOmittedAsInternal() {
		return false
	}

	// A service is internal if it has any internal methods
	for _, m := range append(s.GetMethod(), g.getMixinMethods()...) {
		parent := g.descInfo.ParentElement[m].(*descriptorpb.ServiceDescriptorProto)
		fqn := fmt.Sprintf("%s.%s.%s", g.descInfo.ParentFile[parent].GetPackage(), parent.GetName(), m.GetName())
		if g.isInternalMethod(protoPkg, fqn) {
			return true
		}
	}
	return false
}

// methodName returns the appropriate Go method name (exported or unexported) for a given RPC.
func (g *generator) methodName(m *descriptorpb.MethodDescriptorProto) string {
	// 1. Determine the logical API package context.
	//
	// When we generate mixin methods (like IAM's GetIamPolicy or LRO's GetOperation),
	// their raw Protobuf definitions live in generic packages like `google.iam.v1`
	// rather than the specific service package being generated (e.g. `google.example.v1`).
	//
	// However, the `service.yaml` publishing settings (including the SGG allowlist)
	// that govern these mixins are defined under the *current service's* version block
	// (e.g. `version: google.example.v1`).
	//
	// Therefore, to look up the correct SGG config, we must use `g.clientProtoPkg`
	// (the package of the service we are *currently generating a client for*),
	// rather than the package where the raw method was originally defined.
	protoPkg := g.clientProtoPkg

	// 2. Fallback for testing environments.
	//
	// In some narrow unit testing contexts (like `TestDocFile`), `g.clientProtoPkg`
	// might not be explicitly populated. In these cases, we fall back to reading
	// the package directly from the method's parent service.
	if protoPkg == "" {
		if parent, ok := g.descInfo.ParentElement[m].(*descriptorpb.ServiceDescriptorProto); ok && parent != nil {
			protoPkg = g.descInfo.ParentFile[parent].GetPackage()
		}
	}

	// 3. Look up the Selective GAPIC config.
	sgg := g.getSelectiveGapicGeneration(protoPkg)
	if sgg == nil || !sgg.GetGenerateOmittedAsInternal() {
		return m.GetName()
	}

	// 4. Construct the Fully Qualified Name (FQN) of the method.
	// We MUST use the method's *original* package and service name for the FQN,
	// because that is how it is referenced in the `service.yaml` allowlist
	// (e.g. `methods: ["google.iam.v1.IAMPolicy.GetIamPolicy"]`).
	parent, ok := g.descInfo.ParentElement[m].(*descriptorpb.ServiceDescriptorProto)
	if !ok || parent == nil {
		return m.GetName()
	}
	fqn := fmt.Sprintf("%s.%s.%s", g.descInfo.ParentFile[parent].GetPackage(), parent.GetName(), m.GetName())

	// 5. Evaluate visibility.
	if g.isInternalMethod(protoPkg, fqn) {
		return lowerFirst(m.GetName())
	}
	return m.GetName()
}

// clientName returns the appropriate Go client name for a given service.
func (g *generator) clientName(s *descriptorpb.ServiceDescriptorProto, pkgName string) string {
	override := g.getServiceNameOverride(s)
	servName := pbinfo.ReduceServNameWithOverride(s.GetName(), pkgName, override)
	if g.isInternalService(s) {
		servName = "Base" + servName
	}
	return servName
}
