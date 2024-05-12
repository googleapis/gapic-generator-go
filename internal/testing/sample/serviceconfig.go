// Copyright 2024 Google LLC
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

// Package sample provides functionality for generating sample values of
// the types contained in the internal package for testing purposes.
package sample

import (
	"fmt"

	"google.golang.org/genproto/googleapis/api"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/genproto/googleapis/api/serviceconfig"
	"google.golang.org/protobuf/types/known/apipb"
)

// ServiceConfig returns service config information.
func ServiceConfig() *serviceconfig.Service {
	return &serviceconfig.Service{
		Name:  ServiceURL,
		Title: ServiceTitle,
		Apis: []*apipb.Api{
			{
				Name: "google.cloud.location.Locations",
			},
			{
				Name: fmt.Sprintf("%s.%s", ProtoPackagePath, ServiceName),
			},
		},
		Documentation: &serviceconfig.Documentation{
			Summary:  "Stores sensitive data such as API keys, passwords, and certificates. Provides convenience while improving security.",
			Overview: "Secret Manager Overview",
			Rules: []*serviceconfig.DocumentationRule{
				{
					Selector:    "google.cloud.location.Locations.GetLocation",
					Description: "Gets information about a location.",
				},
				{
					Selector:    "google.cloud.location.Locations.ListLocations",
					Description: "Lists information about the supported locations for this service.",
				},
			},
		},
		Http: &annotations.Http{
			Rules: []*annotations.HttpRule{
				{
					Selector: "google.cloud.location.Locations.GetLocation",
					Pattern: &annotations.HttpRule_Get{
						Get: "/v1/{name=projects/*/locations/*}",
					},
				},
				{
					Selector: "google.cloud.location.Locations.ListLocation",
					Pattern: &annotations.HttpRule_Get{
						Get: "/v1/{name=projects/*}/locations",
					},
				},
			},
		},
		Authentication: &serviceconfig.Authentication{
			Rules: []*serviceconfig.AuthenticationRule{
				{
					Selector: "google.cloud.location.Locations.GetLocation",
					Oauth: &serviceconfig.OAuthRequirements{
						CanonicalScopes: ServiceOAuthScope,
					},
				},
				{
					Selector: "google.cloud.location.Locations.ListLocations",
					Oauth: &serviceconfig.OAuthRequirements{
						CanonicalScopes: ServiceOAuthScope,
					},
				},
				{
					Selector: "'google.cloud.secretmanager.v1.SecretManagerService.*'",
					Oauth: &serviceconfig.OAuthRequirements{
						CanonicalScopes: ServiceOAuthScope,
					},
				},
			},
		},
		Publishing: &annotations.Publishing{
			NewIssueUri:      "https://issuetracker.google.com/issues/new?component=784854&template=1380926",
			DocumentationUri: "https://cloud.google.com/secret-manager/docs/overview",
			ApiShortName:     "secretmanager",
			GithubLabel:      "'api: secretmanager'",
			DocTagPrefix:     "secretmanager",
			Organization:     annotations.ClientLibraryOrganization_CLOUD,
			LibrarySettings: []*annotations.ClientLibrarySettings{
				{
					Version:     "google.cloud.secretmanager.v1",
					LaunchStage: api.LaunchStage_GA,
					GoSettings: &annotations.GoSettings{
						Common: &annotations.CommonLanguageSettings{
							Destinations: []annotations.ClientLibraryDestination{
								annotations.ClientLibraryDestination_PACKAGE_MANAGER,
							},
						},
					},
				},
			},
			ProtoReferenceDocumentationUri: "https://cloud.google.com/secret-manager/docs/reference/rpc",

			// The data below is fake. method_settings is not used in
			// https://github.com/googleapis/googleapis/blob/f7df662a24c56ecaab79cb7d808fed4d2bb4981d/google/cloud/secretmanager/v1/secretmanager_v1.yaml.
			MethodSettings: []*annotations.MethodSettings{
				{
					Selector: fmt.Sprintf("%s.%s.%s", ProtoPackagePath, ServiceName, CreateMethodWithSettings),
					AutoPopulatedFields: []string{
						"request_id",
					},
				},
			},
		},
	}
}
