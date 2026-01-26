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
	"testing"

	"google.golang.org/genproto/googleapis/api/serviceconfig"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func TestProcessLegacyEnablements(t *testing.T) {
	for _, tc := range []struct {
		desc                string
		protoPkg            string
		apiName             string
		wantFeaturesEnabled []featureID
		wantFeatureDisabled []featureID
	}{
		{
			desc: "default",
			wantFeatureDisabled: []featureID{
				EnableMTLSHardBoundTokens,
				EnableOrderedRoutingHeaders,
				EnableWrapperTypesForPageSize,
			},
		},
		{
			desc:                "bigquery package",
			protoPkg:            "google.cloud.bigquery.v2",
			wantFeaturesEnabled: []featureID{EnableWrapperTypesForPageSize},
		},
		{
			desc:                "bigquery package only",
			protoPkg:            "google.cloud.bigquery.v2",
			wantFeaturesEnabled: []featureID{EnableWrapperTypesForPageSize},
			wantFeatureDisabled: []featureID{EnableOrderedRoutingHeaders, EnableMTLSHardBoundTokens},
		},
		{
			desc:                "firestore admin package only",
			protoPkg:            "google.firestore.admin.v1",
			wantFeaturesEnabled: []featureID{EnableOrderedRoutingHeaders},
			wantFeatureDisabled: []featureID{EnableWrapperTypesForPageSize, EnableMTLSHardBoundTokens},
		},
		{
			desc:                "cloud kms api name",
			apiName:             "cloudkms.googleapis.com",
			wantFeaturesEnabled: []featureID{EnableMTLSHardBoundTokens},
			wantFeatureDisabled: []featureID{EnableWrapperTypesForPageSize, EnableOrderedRoutingHeaders},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			req := &pluginpb.CodeGeneratorRequest{
				Parameter: proto.String("go-gapic-package=foo/bar/baz;baz"),
				ProtoFile: []*descriptorpb.FileDescriptorProto{
					{
						Package: proto.String("foo"),
					},
				},
			}
			if tc.protoPkg != "" {
				req.ProtoFile[0].Package = proto.String(tc.protoPkg)
			}
			cfg, err := configFromRequest(req.Parameter)
			if err != nil {
				t.Fatalf("configFromRequest err: %v", err)
			}
			if tc.apiName != "" {
				cfg.APIServiceConfig = &serviceconfig.Service{
					Name: tc.apiName,
				}
			}
			processLegacyEnablements(cfg, req)
			for _, f := range tc.wantFeaturesEnabled {
				if !cfg.FeatureEnabled(f) {
					t.Errorf("expected feature %q enabled, was not", f)
				}
			}
			for _, f := range tc.wantFeatureDisabled {
				if cfg.FeatureEnabled(f) {
					t.Errorf("expected feature %q to be disabled, was enabled", f)
				}
			}
		})

	}
}
