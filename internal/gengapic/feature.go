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

import "google.golang.org/protobuf/types/pluginpb"

// We use a custom type for feature ID to clarify that features are setup explicitly.
type featureID string

// featureInfo contains basic information about features, and provides an extensible mechanism
// for adding more infomation down the line (maturity level, warning about stale experiments, etc)
type featureInfo struct {
	// Short summary of feature.
	Description string

	// Tracking ID for more information.  Github issue, internal b/ issue, etc.
	TrackingID string
}

// Define feature ID strings here.  More details about features are kept in the featureRegistry map.
const (
	WrapperTypesForPageSizeFeature featureID = "wrapper_types_for_page_size"
	OrderedRoutingHeadersFeature   featureID = "ordered_routing_headers"
	MTLSHardBoundTokensFeature     featureID = "mtls_hard_bound_tokens"
)

// featureRegistry contains the registry of defined features.
// Introducing a new capability to the generator generally starts here, as features
// must be registered to be enabled.  This should not be modified at runtime.  Those
// who attempt to do so will be given a stern talking to.
var featureRegistry = map[featureID]*featureInfo{
	MTLSHardBoundTokensFeature: {
		Description: "support MTLS hard bound tokens",
		TrackingID:  "b/327916505",
	},
	OrderedRoutingHeadersFeature: {
		Description: "Specify that routing headers are emitted in a deterministic fashion.  Primarily used for firestore.",
	},
	WrapperTypesForPageSizeFeature: {
		Description: "Allow List RPCs to generator with support for protobuf wrapper types (e.g. Int32Value, etc).",
		TrackingID:  "b/352331075",
	},
}

// legacyFeatureEnablementByPackage is temporary bridge functionality.  Features should be enabled via protoc flags, but
// to bootstrap without breaking generation we keep the legacy definitions enabled here until we can move
// configuration upstream into tools like librarian/bazel/etc as needed.
var legacyFeatureEnablementByPackage = map[featureID][]string{
	OrderedRoutingHeadersFeature: []string{
		"google.firestore.v1",
		"google.firestore.admin.v1",
	},
	WrapperTypesForPageSizeFeature: []string{
		"google.cloud.bigquery.v2",
	},
}

// similar to legacyFeatureEnablementByPackage, this is legacy feature enablement using the "name" field from the API
// service config.
var legacyFeatureEnablementByAPIName = map[featureID][]string{
	MTLSHardBoundTokensFeature: []string{
		"bigquery.googleapis.com",
		"cloudasset.googleapis.com",
		"clouderrorreporting.googleapis.com",
		"cloudkms.googleapis.com",
		"cloudresourcemanager.googleapis.com",
		"cloudtasks.googleapis.com",
		"cloudtrace.googleapis.com",
		"dataflow.googleapis.com",
		"datastore.googleapis.com",
		"essentialcontacts.googleapis.com",
		"firestore.googleapis.com",
		"iam.googleapis.com",
		"iamcredentials.googleapis.com",
		"logging.googleapis.com",
		"monitoring.googleapis.com",
		"orgpolicy.googleapis.com",
		"pubsub.googleapis.com",
		"recommender.googleapis.com",
		"secretmanager.googleapis.com",
		"showcase.googleapis.com",
	},
}

// This function consolidates legacy processing of feature enablements.
// Like the associated allowlists, it should go away once librarian and bazel can pass feature enablements directly.
func processLegacyEnablements(cfg *generatorConfig, req *pluginpb.CodeGeneratorRequest) {
	// Use the first proto file in the FileDescriptorSet to handle legacy enablement by package name.
	if len(req.GetProtoFile()) > 0 {
		probePackage := req.GetProtoFile()[0].GetPackage()
		for f, packages := range legacyFeatureEnablementByPackage {
			for _, v := range packages {
				if probePackage == v { // matched
					if cfg.featureEnablement == nil {
						cfg.featureEnablement = make(map[featureID]bool)
					}
					cfg.featureEnablement[f] = true
					break
				}
			}
		}
	}
	// Now, process legacy feature enablements based on API name.
	if cfg.APIServiceConfig != nil {
		probeName := cfg.APIServiceConfig.GetName()
		for f, apis := range legacyFeatureEnablementByAPIName {
			for _, v := range apis {
				if probeName == v { // matched
					if cfg.featureEnablement == nil {
						cfg.featureEnablement = make(map[featureID]bool)
					}
					cfg.featureEnablement[f] = true
					break
				}
			}
		}
	}
}
