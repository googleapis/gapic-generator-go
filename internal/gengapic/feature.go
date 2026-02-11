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
	WrapperTypesForPageSizeFeature  featureID = "wrapper_types_for_page_size"
	OrderedRoutingHeadersFeature    featureID = "ordered_routing_headers"
	MTLSHardBoundTokensFeature      featureID = "mtls_hard_bound_tokens"
	OpenTelemetryTracingFeature     featureID = "open_telemetry_tracing"
	SelectiveGapicGenerationFeature featureID = "selectivegapicgenerationfeature"
)

// featureRegistry contains the registry of defined features.
// Introducing a new capability to the generator generally starts here, as features
// must be registered to be enabled.  This should not be modified at runtime.  Those
// who attempt to do so will be given a stern talking to.
var featureRegistry = map[featureID]*featureInfo{

	SelectiveGapicGenerationFeature: {
		Description: "Enable selective GAPIC generation using google.api.ClientLibrarySettings.",
		TrackingID:  "b/4483092298",
	},
	OpenTelemetryTracingFeature: {
		Description: "Enable OpenTelemetry tracing support (Service Identity, Resource Names, URL Templates).",
		TrackingID:  "b/467342602,b/467403185",
	},
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
