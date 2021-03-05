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
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
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

func (g *generator) collectMixins() {
	for _, api := range g.serviceConfig.GetApis() {
		if _, ok := mixinFiles[api.GetName()]; ok {
			g.mixins[api.GetName()] = true
		}
	}
}

func (g *generator) getMixinFiles() []*descriptor.FileDescriptorProto {
	files := []*descriptor.FileDescriptorProto{}
	for key := range g.mixins {
		files = append(files, mixinFiles[key]...)
	}
	return files
}

func (g *generator) getMixinMethods(serv *descriptor.ServiceDescriptorProto) []*descriptor.MethodDescriptorProto {
	methods := []*descriptor.MethodDescriptorProto{}
	if g.hasLocationMixin() {
		methods = append(methods, getLocationsMethods()...)
	}
	if g.hasIAMPolicyMixin() && !hasIAMPolicyOverrides(serv) {
		methods = append(methods, getIAMPolicyMethods()...)
	}
	if g.hasLROMixin() {
		methods = append(methods, getOperationsMethods()...)
	}

	return methods
}

func (g *generator) hasLROMixin() bool {
	return g.mixins["google.longrunning.Operations"]
}

func (g *generator) hasIAMPolicyMixin() bool {
	return g.mixins["google.iam.v1.IAMPolicy"]
}

func (g *generator) hasLocationMixin() bool {
	return g.mixins["google.cloud.location.Locations"]
}

func hasIAMPolicyOverrides(serv *descriptor.ServiceDescriptorProto) bool {
	for _, iamMethod := range getIAMPolicyMethods() {
		if hasMethod(serv, iamMethod.GetName()) {
			return true
		}
	}

	return false
}

func getLocationsDescriptor() *descriptor.FileDescriptorProto {
	return mixinFiles["google.cloud.location.Locations"][0]
}

func getLocationsMethods() []*descriptor.MethodDescriptorProto {
	return getLocationsDescriptor().GetService()[0].GetMethod()
}

func getIAMPolicyDescriptors() []*descriptor.FileDescriptorProto {
	return mixinFiles["google.iam.v1.IAMPolicy"]
}

func getIAMPolicyMethods() []*descriptor.MethodDescriptorProto {
	return getIAMPolicyDescriptors()[0].GetService()[0].GetMethod()
}

func getOperationsDescriptor() *descriptor.FileDescriptorProto {
	return mixinFiles["google.longrunning.Operations"][0]
}

func getOperationsMethods() []*descriptor.MethodDescriptorProto {
	return getOperationsDescriptor().GetService()[0].GetMethod()
}
