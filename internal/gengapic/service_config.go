// Copyright 2019 Google LLC
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

// serviceConfig represents a gapic service config
// Deprecated: workaround for not having annotations yet; to be removed
type serviceConfig struct {
	Name           string
	Title          string
	Documentation  *configDocumentation
	Authentication *configAuthentication
}

// configDocumentation represents gapic service config documentation section
// Deprecated: workaround for not having annotations yet; to be removed
type configDocumentation struct {
	Summary string
}

// configAuthentication represents gapic service config authentication section
// Deprecated: workaround for not having annotations yet; to be removed
type configAuthentication struct {
	Rules []*configAuthRules
}

// configAuthRules represents gapic service config auth rules
// Deprecated: workaround for not having annotations yet; to be removed
type configAuthRules struct {
	Selector string
	Oauth    *configAuthRulesOauth
}

// configAuthRulesOauth represents a gapic service config oauth rule
// Deprecated: workaround for not having annotations yet; to be removed
type configAuthRulesOauth struct {
	CanonicalScopes string `yaml:"canonical_scopes"`
}
