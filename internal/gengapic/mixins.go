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
	"regexp"
	"strings"

	iam "cloud.google.com/go/iam/apiv1/iampb"
	longrunning "cloud.google.com/go/longrunning/autogen/longrunningpb"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/genproto/googleapis/cloud/location"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/types/descriptorpb"
)

func init() {
	initMixinFiles()
}

var apiVersionRegexp = regexp.MustCompile(`v\d+[a-z]*\d*[a-z]*\d*`)

var mixinFiles map[string][]*descriptorpb.FileDescriptorProto

type mixins map[string][]*descriptorpb.MethodDescriptorProto

// initMixinFiles allows test code to re-initialize the mixinFiles global.
func initMixinFiles() {
	mixinFiles = map[string][]*descriptorpb.FileDescriptorProto{
		"google.cloud.location.Locations": {
			protodesc.ToFileDescriptorProto(location.File_google_cloud_location_locations_proto),
		},
		"google.iam.v1.IAMPolicy": {
			protodesc.ToFileDescriptorProto(iam.File_google_iam_v1_iam_policy_proto),
			protodesc.ToFileDescriptorProto(iam.File_google_iam_v1_policy_proto),
			protodesc.ToFileDescriptorProto(iam.File_google_iam_v1_options_proto),
		},
		"google.longrunning.Operations": {
			protodesc.ToFileDescriptorProto(longrunning.File_google_longrunning_operations_proto),
		},
	}
}

// collectMixins collects the configured mixin APIs from the Service config and
// gathers the appropriately configured mixin methods to generate for each.
func (g *generator) collectMixins() {
	for _, api := range g.serviceConfig.GetApis() {
		if _, ok := mixinFiles[api.GetName()]; ok {
			g.mixins[api.GetName()] = g.collectMixinMethods(api.GetName())
		}
	}
}

// collectMixinMethods collects the method descriptors for the given mixin API
// that should be generated for a client. In order for a method to be included
// for generation, it must have a google.api.http defined in the Service config
// http.rules section. The MethodDescriptorProto.options are overwritten with
// that same google.api.http binding. Furthermore, a basic leading comment is
// defined for the method to be generated.
func (g *generator) collectMixinMethods(api string) []*descriptorpb.MethodDescriptorProto {
	methods := map[string]*descriptorpb.MethodDescriptorProto{}
	methodsToGenerate := []*descriptorpb.MethodDescriptorProto{}

	// Note: Triple nested loops are nasty, but this is tightly bound and really
	// the only way to traverse proto descriptors that are backed by slices.
	for _, file := range mixinFiles[api] {
		for _, service := range file.GetService() {
			for _, method := range service.GetMethod() {
				fqn := fmt.Sprintf("%s.%s.%s", file.GetPackage(), service.GetName(), method.GetName())
				methods[fqn] = method

				// Set a default comment in case the Service does not have a DocumentationRule for it.
				// Exclude the leading method name because methodDoc adds automatically.
				g.comments[method] = fmt.Sprintf("is a utility method from %s.%s.", file.GetPackage(), service.GetName())
			}
		}
	}

	// Overwrite the google.api.http annotations with bindings from the Service config.
	for _, rule := range g.serviceConfig.GetHttp().GetRules() {
		m, match := methods[rule.GetSelector()]
		if !match {
			continue
		}

		proto.SetExtension(m.Options, annotations.E_Http, rule)
		methodsToGenerate = append(methodsToGenerate, m)
	}

	// Include any documentation from the Service config.
	for _, rule := range g.serviceConfig.GetDocumentation().GetRules() {
		m, match := methods[rule.GetSelector()]
		if !match {
			continue
		}

		g.comments[m] = rule.GetDescription()
	}

	return methodsToGenerate
}

// getMixinFiles returns a set of file descriptors for the APIs configured to be
// mixed in.
func (g *generator) getMixinFiles() []*descriptorpb.FileDescriptorProto {
	files := []*descriptorpb.FileDescriptorProto{}
	for key := range g.mixins {
		files = append(files, mixinFiles[key]...)
	}
	return files
}

// getMixinMethods is a convenience method to collect the method descriptors of
// those methods to be generated based on if they should be included or not.
func (g *generator) getMixinMethods() []*descriptorpb.MethodDescriptorProto {
	methods := []*descriptorpb.MethodDescriptorProto{}
	if g.hasLocationMixin() {
		methods = append(methods, g.mixins["google.cloud.location.Locations"]...)
	}
	if g.hasIAMPolicyMixin() {
		methods = append(methods, g.mixins["google.iam.v1.IAMPolicy"]...)
	}
	if g.hasLROMixin() {
		methods = append(methods, g.mixins["google.longrunning.Operations"]...)
	}

	return methods
}

// mixinStubs prints the field definition for the mixin gRPC stubs that are
// configured to be generated. This is used in the definition of the generated
// client type(s).
func (g *generator) mixinStubs() {
	p := g.printf

	if g.hasLROMixin() {
		p("operationsClient longrunningpb.OperationsClient")
		p("")

		g.imports[pbinfo.ImportSpec{Name: "longrunningpb", Path: "cloud.google.com/go/longrunning/autogen/longrunningpb"}] = true
	}

	if g.hasIAMPolicyMixin() {

		p("iamPolicyClient iampb.IAMPolicyClient")
		p("")

		g.imports[pbinfo.ImportSpec{Name: "iampb", Path: "cloud.google.com/go/iam/apiv1/iampb"}] = true
	}

	if g.hasLocationMixin() {

		p("locationsClient locationpb.LocationsClient")
		p("")

		g.imports[pbinfo.ImportSpec{Name: "locationpb", Path: "google.golang.org/genproto/googleapis/cloud/location"}] = true
	}
}

// mixinStubsInit prints the stub intialization for the mixin gRPC stubs that are
// configured to be generated. This is used in the factory method of a generated client.
func (g *generator) mixinStubsInit() {
	p := g.printf

	if g.hasLROMixin() {
		p("    operationsClient: longrunningpb.NewOperationsClient(connPool),")
	}
	if g.hasIAMPolicyMixin() {
		p("    iamPolicyClient: iampb.NewIAMPolicyClient(connPool),")
	}
	if g.hasLocationMixin() {
		p("    locationsClient: locationpb.NewLocationsClient(connPool),")
	}
}

// hasLROMixin is a convenience method for determining if the Operations mixin
// should be generated.
func (g *generator) hasLROMixin() bool {
	return len(g.mixins["google.longrunning.Operations"]) > 0 && len(g.serviceConfig.GetApis()) > 1
}

// hasIAMPolicyMixin is a convenience method for determining if the IAMPolicy
// mixin should be generated.
func (g *generator) hasIAMPolicyMixin() bool {
	return len(g.mixins["google.iam.v1.IAMPolicy"]) > 0 && !g.hasIAMPolicyOverrides && len(g.serviceConfig.GetApis()) > 1
}

// hasLocationMixin is a convenience method for determining if the Locations
// mixin should be generated.
func (g *generator) hasLocationMixin() bool {
	return len(g.mixins["google.cloud.location.Locations"]) > 0 && len(g.serviceConfig.GetApis()) > 1
}

// containsIAMPolicyOverrides determines if any of the given services define an
// IAMPolicy RPC and sets the hasIAMpolicyOverrides generator flag if so. If set
// to true, the IAMPolicy mixin will not be generated on any service client. This
// is for backwards compatibility with existing IAMPolicy redefinitions.
func (g *generator) containsIAMPolicyOverrides(servs []*descriptorpb.ServiceDescriptorProto) bool {
	iam, hasMixin := g.mixins["google.iam.v1.IAMPolicy"]
	if !hasMixin {
		return false
	}

	for _, s := range servs {
		for _, iamMethod := range iam {
			if hasMethod(s, iamMethod.GetName()) {
				return true
			}
		}
	}
	return false
}

// includeMixinInputFile determines if the given proto file name matches
// a known mixin file and indicates if it should be included in the
// protos-to-be-generated file set based on if the Go package to be generated
// is for one of the mixin services explicitly or not.
func (g *generator) includeMixinInputFile(file string) bool {
	if strings.HasPrefix(file, "google/cloud/location") && !strings.Contains(g.opts.pkgPath, "location") {
		return false
	}
	if strings.HasPrefix(file, "google/iam/v1") && !strings.Contains(g.opts.pkgPath, "iam") {
		return false
	}
	if strings.HasPrefix(file, "google/longrunning") && !strings.Contains(g.opts.pkgPath, "longrunning") {
		return false
	}
	// Not a mixin file or generating a mixin GAPIC explicitly so include the file.
	return true
}

// lookUpGetOperationOverride looks up the google.api.http rule defined in the
// service config for the given RPC.
func (g *generator) lookupHTTPOverride(fqn string, f func(h *annotations.HttpRule) string) string {
	for _, rule := range g.serviceConfig.GetHttp().GetRules() {
		if rule.GetSelector() == fqn {
			return f(rule)
		}
	}

	return ""
}

// getOperationPathOverride looks up the google.api.http rule for LRO GetOperation
// and returns the path override. If no value is present, it synthesizes a path
// using the proto package client version, for example, "/v1/{name=operations/**}".
func (g *generator) getOperationPathOverride(protoPkg string) string {
	get := func(h *annotations.HttpRule) string { return h.GetGet() }
	override := g.lookupHTTPOverride("google.longrunning.Operations.GetOperation", get)
	if override == "" {
		// extract httpInfo from "hot loaded" Operations.GetOperation MethodDescriptor
		// Should be "/v1/{name=operations/**}"
		file := mixinFiles["google.longrunning.Operations"][0]
		mdp := getMethod(file.GetService()[0], "GetOperation")
		getOperationPath := getHTTPInfo(mdp).url

		// extract client version from proto package with global regex
		// replace version base path in GetOperation path with proto package version segment
		version := apiVersionRegexp.FindStringSubmatch(protoPkg)
		override = apiVersionRegexp.ReplaceAllStringFunc(getOperationPath, func(s string) string { return version[0] })
	}
	override = httpPatternVarRegex.ReplaceAllStringFunc(override, func(s string) string { return "%s" })
	return override
}
