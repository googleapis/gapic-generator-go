// Copyright 2026 Google LLC
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
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/types/descriptorpb"
)

func (g *generator) getSelectiveGapicGeneration(protoPkg string) *annotations.SelectiveGapicGeneration {
	if g.cfg == nil || g.cfg.APIServiceConfig == nil || g.cfg.APIServiceConfig.GetPublishing() == nil {
		return nil
	}
	if !g.featureEnabled(SelectiveGapicGenerationFeature) {
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

func (g *generator) getMethods(s *descriptorpb.ServiceDescriptorProto) []*descriptorpb.MethodDescriptorProto {
	methods := append(s.GetMethod(), g.getMixinMethods()...)
	var filtered []*descriptorpb.MethodDescriptorProto
	for _, m := range methods {
		if g.shouldGenerateMethod(s, m) {
			filtered = append(filtered, m)
		}
	}
	return filtered
}

func (g *generator) shouldGenerateMethod(s *descriptorpb.ServiceDescriptorProto, m *descriptorpb.MethodDescriptorProto) bool {
	protoPkg := g.descInfo.ParentFile[s].GetPackage()
	sgg := g.getSelectiveGapicGeneration(protoPkg)
	if sgg == nil {
		return true
	}

	methods := sgg.GetMethods()
	generateOmittedAsInternal := sgg.GetGenerateOmittedAsInternal()

	if len(methods) == 0 && !generateOmittedAsInternal {
		return true
	}

	mfqn := g.fqn(m)
	for _, inc := range methods {
		if inc == mfqn {
			return true
		}
	}

	return generateOmittedAsInternal
}

func (g *generator) isMethodInternal(m *descriptorpb.MethodDescriptorProto) bool {
	protoPkg := g.clientProtoPkg
	if protoPkg == "" {
		if parent, ok := g.descInfo.ParentElement[m].(*descriptorpb.ServiceDescriptorProto); ok && parent != nil {
			protoPkg = g.descInfo.ParentFile[parent].GetPackage()
		}
	}

	sgg := g.getSelectiveGapicGeneration(protoPkg)
	if sgg == nil || !sgg.GetGenerateOmittedAsInternal() {
		return false
	}

	mfqn := g.fqn(m)
	for _, inc := range sgg.GetMethods() {
		if inc == mfqn {
			return false
		}
	}

	return true
}

func (g *generator) methodName(m *descriptorpb.MethodDescriptorProto) string {
	if g.isMethodInternal(m) {
		return lowerFirst(m.GetName())
	}
	return m.GetName()
}

func (g *generator) isInternalService(s *descriptorpb.ServiceDescriptorProto) bool {
	protoPkg := g.descInfo.ParentFile[s].GetPackage()
	sgg := g.getSelectiveGapicGeneration(protoPkg)
	if sgg == nil || !sgg.GetGenerateOmittedAsInternal() {
		return false
	}

	// A service is internal if it has any internal methods.
	// We check the original methods (including mixins) before filtering.
	for _, m := range append(s.GetMethod(), g.getMixinMethods()...) {
		if g.isMethodInternal(m) {
			return true
		}
	}
	return false
}
