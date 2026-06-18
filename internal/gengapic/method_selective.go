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

type sggConfig struct {
	allowedMethods            map[string]bool
	generateOmittedAsInternal bool
	enabled                   bool
}

func (g *generator) getSGGConfig(protoPkg string) *sggConfig {
	if g.sggConfigs == nil {
		g.sggConfigs = make(map[string]*sggConfig)
	}

	if config, ok := g.sggConfigs[protoPkg]; ok {
		return config
	}

	config := &sggConfig{
		allowedMethods: make(map[string]bool),
	}
	g.sggConfigs[protoPkg] = config

	if g.cfg == nil || g.cfg.APIServiceConfig == nil || g.cfg.APIServiceConfig.GetPublishing() == nil {
		return config
	}
	if !g.featureEnabled(SelectiveGapicGenerationFeature) {
		return config
	}
	ls := g.cfg.APIServiceConfig.GetPublishing().GetLibrarySettings()

	var sgg *annotations.SelectiveGapicGeneration
	for _, setting := range ls {
		if setting.GetVersion() != "" && setting.GetVersion() != protoPkg {
			continue
		}
		if goSettings := setting.GetGoSettings(); goSettings != nil && goSettings.GetCommon() != nil && goSettings.GetCommon().GetSelectiveGapicGeneration() != nil {
			sgg = goSettings.GetCommon().GetSelectiveGapicGeneration()
			break
		}
	}

	if sgg == nil {
		return config
	}

	config.enabled = true
	config.generateOmittedAsInternal = sgg.GetGenerateOmittedAsInternal()
	for _, m := range sgg.GetMethods() {
		config.allowedMethods[m] = true
	}

	return config
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
	sgg := g.getSGGConfig(protoPkg)
	if sgg == nil || !sgg.enabled {
		return true
	}

	if len(sgg.allowedMethods) == 0 && !sgg.generateOmittedAsInternal {
		return true
	}

	mfqn := g.fqn(m)
	if sgg.allowedMethods[mfqn] {
		return true
	}

	return sgg.generateOmittedAsInternal
}

func (g *generator) isMethodInternal(m *descriptorpb.MethodDescriptorProto) bool {
	protoPkg := g.clientProtoPkg
	if protoPkg == "" {
		if parent, ok := g.descInfo.ParentElement[m].(*descriptorpb.ServiceDescriptorProto); ok && parent != nil {
			protoPkg = g.descInfo.ParentFile[parent].GetPackage()
		}
	}

	sgg := g.getSGGConfig(protoPkg)
	if sgg == nil || !sgg.enabled || !sgg.generateOmittedAsInternal {
		return false
	}

	mfqn := g.fqn(m)
	return !sgg.allowedMethods[mfqn]
}

func (g *generator) methodName(m *descriptorpb.MethodDescriptorProto) string {
	if g.isMethodInternal(m) {
		return lowerFirst(m.GetName())
	}
	return m.GetName()
}

func (g *generator) isInternalService(s *descriptorpb.ServiceDescriptorProto) bool {
	protoPkg := g.descInfo.ParentFile[s].GetPackage()
	sgg := g.getSGGConfig(protoPkg)
	if sgg == nil || !sgg.enabled || !sgg.generateOmittedAsInternal {
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
