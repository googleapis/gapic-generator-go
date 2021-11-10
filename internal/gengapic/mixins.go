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
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/googleapis/gapic-generator-go/internal/pbinfo"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/genproto/googleapis/cloud/location"
	iam "google.golang.org/genproto/googleapis/iam/v1"
	"google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/protobuf/reflect/protodesc"
)

func init() {
	mixinFiles = map[string][]*descriptor.FileDescriptorProto{
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

var mixinFiles map[string][]*descriptor.FileDescriptorProto

type mixins map[string][]*descriptor.MethodDescriptorProto

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
func (g *generator) collectMixinMethods(api string) []*descriptor.MethodDescriptorProto {
	methods := map[string]*descriptor.MethodDescriptorProto{}
	methodsToGenerate := []*descriptor.MethodDescriptorProto{}

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

		if err := proto.SetExtension(m.Options, annotations.E_Http, rule); err != nil {
			log.Println("Encountered error setting HTTP annotations:", err)
		}
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
func (g *generator) getMixinFiles() []*descriptor.FileDescriptorProto {
	files := []*descriptor.FileDescriptorProto{}
	for key := range g.mixins {
		files = append(files, mixinFiles[key]...)
	}
	return files
}

// getMixinMethods is a convenience method to collect the method descriptors of
// those methods to be generated based on if they should be included or not.
func (g *generator) getMixinMethods() []*descriptor.MethodDescriptorProto {
	methods := []*descriptor.MethodDescriptorProto{}
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

		g.imports[pbinfo.ImportSpec{Name: "longrunningpb", Path: "google.golang.org/genproto/googleapis/longrunning"}] = true
	}

	if g.hasIAMPolicyMixin() {

		p("iamPolicyClient iampb.IAMPolicyClient")
		p("")

		g.imports[pbinfo.ImportSpec{Name: "iampb", Path: "google.golang.org/genproto/googleapis/iam/v1"}] = true
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

// checkIAMPolicyOverrides determines if any of the given services define an
// IAMPolicy RPC and sets the hasIAMpolicyOverrides generator flag if so. If set
// to true, the IAMPolicy mixin will not be generated on any service client. This
// is for backwards compatibility with existing IAMPolicy redefinitions.
func (g *generator) checkIAMPolicyOverrides(servs []*descriptor.ServiceDescriptorProto) {
	iam, hasMixin := g.mixins["google.iam.v1.IAMPolicy"]
	if !hasMixin {
		return
	}

	for _, s := range servs {
		for _, iamMethod := range iam {
			if hasMethod(s, iamMethod.GetName()) {
				g.hasIAMPolicyOverrides = true
				return
			}
		}
	}
}

// includeMixinInputFile determines if the given proto file name matches
// a known mixin file and indicates if it should be included in the
// protos-to-be-generated file set based on if the package is using it for
// mixins or not.
func (g *generator) includeMixinInputFile(file string) bool {
	if file == "google/cloud/location/locations.proto" && g.hasLocationMixin() {
		return false
	}
	if file == "google/iam/v1/iam_policy.proto" && g.hasIAMPolicyMixin() {
		return false
	}
	if file == "google/longrunning/operations.proto" && g.hasLROMixin() {
		return false
	}
	return true
}
