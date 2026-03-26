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
	"reflect"
	"testing"

	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func TestBuildHeuristicVocabulary(t *testing.T) {
	tests := []struct {
		name    string
		methods []*descriptorpb.MethodDescriptorProto
		want    map[string]bool
	}{
		{
			name: "Standard CRUD paths",
			methods: []*descriptorpb.MethodDescriptorProto{
				{
					Name: proto.String("GetTopic"),
					Options: func() *descriptorpb.MethodOptions {
						opts := &descriptorpb.MethodOptions{}
						proto.SetExtension(opts, annotations.E_Http, &annotations.HttpRule{
							Pattern: &annotations.HttpRule_Get{Get: "v1/projects/{project}/topics/{topic}"},
						})
						return opts
					}(),
				},
				{
					Name: proto.String("ListSubscriptions"),
					Options: func() *descriptorpb.MethodOptions {
						opts := &descriptorpb.MethodOptions{}
						proto.SetExtension(opts, annotations.E_Http, &annotations.HttpRule{
							Pattern: &annotations.HttpRule_Get{Get: "v1/projects/{project}/subscriptions/{subscription}"},
						})
						return opts
					}(),
				},
			},
			want: map[string]bool{
				"projects":        true,
				"locations":       true,
				"folders":         true,
				"organizations":   true,
				"billingAccounts": true,
				"topics":          true,
				"subscriptions":   true,
			},
		},
		{
			name: "Version string filtering",
			methods: []*descriptorpb.MethodDescriptorProto{
				{
					Name: proto.String("GetProject"),
					Options: func() *descriptorpb.MethodOptions {
						opts := &descriptorpb.MethodOptions{}
						proto.SetExtension(opts, annotations.E_Http, &annotations.HttpRule{
							Pattern: &annotations.HttpRule_Get{Get: "v1/{project}"},
						})
						return opts
					}(),
				},
			},
			want: map[string]bool{
				"projects":        true,
				"locations":       true,
				"folders":         true,
				"organizations":   true,
				"billingAccounts": true,
			},
		},
		{
			name: "Literals nested inside path variables are ignored",
			methods: []*descriptorpb.MethodDescriptorProto{
				{
					Name: proto.String("UpdateInstance"),
					Options: func() *descriptorpb.MethodOptions {
						opts := &descriptorpb.MethodOptions{}
						proto.SetExtension(opts, annotations.E_Http, &annotations.HttpRule{
							Pattern: &annotations.HttpRule_Patch{Patch: "v1/{name=projects/*/instances/*}"},
							AdditionalBindings: []*annotations.HttpRule{
								{
									Pattern: &annotations.HttpRule_Post{Post: "v1/{instance=projects/*/zones/*/instances/*}:insert"},
								},
							},
						})
						return opts
					}(),
				},
			},
			want: map[string]bool{
				"projects":        true,
				"locations":       true,
				"folders":         true,
				"organizations":   true,
				"billingAccounts": true,
			},
		},
		{
			name: "Verifying custom verb stripping on literals",
			methods: []*descriptorpb.MethodDescriptorProto{
				{
					Name: proto.String("GetTopic"),
					Options: func() *descriptorpb.MethodOptions {
						opts := &descriptorpb.MethodOptions{}
						proto.SetExtension(opts, annotations.E_Http, &annotations.HttpRule{
							Pattern: &annotations.HttpRule_Get{Get: "v1/projects/{project}/topics:cancel/{topic}"},
						})
						return opts
					}(),
				},
			},
			want: map[string]bool{
				"projects":        true,
				"locations":       true,
				"folders":         true,
				"organizations":   true,
				"billingAccounts": true,
				"topics":          true,
			},
		},
		{
			name: "Verifying we ignore non-CRUD methods (:verb pattern)",
			methods: []*descriptorpb.MethodDescriptorProto{
				{
					Name: proto.String("PublishTopic"),
					Options: func() *descriptorpb.MethodOptions {
						opts := &descriptorpb.MethodOptions{}
						proto.SetExtension(opts, annotations.E_Http, &annotations.HttpRule{
							Pattern: &annotations.HttpRule_Post{Post: "v1/projects/{project}/topics/{topic}:publish"},
						})
						return opts
					}(),
				},
			},
			want: map[string]bool{
				"projects":        true,
				"locations":       true,
				"folders":         true,
				"organizations":   true,
				"billingAccounts": true,
			},
		},
		{
			name: "Multiple collections in a single path",
			methods: []*descriptorpb.MethodDescriptorProto{
				{
					Name: proto.String("GetTopic"),
					Options: func() *descriptorpb.MethodOptions {
						opts := &descriptorpb.MethodOptions{}
						proto.SetExtension(opts, annotations.E_Http, &annotations.HttpRule{
							Pattern: &annotations.HttpRule_Get{Get: "v1/projects/{project}/topics/{topic}/snapshots/{snapshot}"},
						})
						return opts
					}(),
				},
			},
			want: map[string]bool{
				"projects":        true,
				"locations":       true,
				"folders":         true,
				"organizations":   true,
				"billingAccounts": true,
				"topics":          true,
				"snapshots":       true,
			},
		},
	}

	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			got := BuildHeuristicVocabulary(tst.methods)
			if !reflect.DeepEqual(got, tst.want) {
				t.Errorf("BuildHeuristicVocabulary(%v): got %v, want %v", tst.name, got, tst.want)
			}
		})
	}
}

func TestIdentifyHeuristicTarget(t *testing.T) {
	vocabulary := map[string]bool{
		"projects":        true,
		"locations":       true,
		"folders":         true,
		"organizations":   true,
		"billingAccounts": true,
		"topics":          true,
		"subscriptions":   true,
	}

	tests := []struct {
		name       string
		methodName string
		pattern    string
		want       *HeuristicTarget
	}{
		{
			name:       "Standard AIP pattern",
			methodName: "GetTopic",
			pattern:    "v1/projects/{project}/topics/{topic}",
			want: &HeuristicTarget{
				Format:     "projects/%v/topics/%v",
				FieldNames: []string{"project", "topic"},
			},
		},
		{
			name:       "Skips unrecognized collections to find the closest known parent",
			methodName: "GetVolume",
			pattern:    "v1/projects/{project}/unsupported/{unsupported}/volumes/{volume}",
			want: &HeuristicTarget{
				Format:     "projects/%v",
				FieldNames: []string{"project"},
			},
		},
		{
			name:       "Deduces parent resource for List endpoints (ends in a literal)",
			methodName: "ListTopics",
			pattern:    "v1/projects/{project}/topics",
			want: &HeuristicTarget{
				Format:     "projects/%v",
				FieldNames: []string{"project"},
			},
		},

		{
			name:       "Base case with a single-segment resource (GetProject)",
			methodName: "GetProject",
			pattern:    "v1/projects/{project}",
			want: &HeuristicTarget{
				Format:     "projects/%v",
				FieldNames: []string{"project"},
			},
		},
		{
			name:       "Compute: interstitial literal 'global' with unknown collection",
			methodName: "GetCrossSiteNetwork",
			pattern:    "v1/projects/{project}/global/crossSiteNetworks/{cross_site_network}",
			want: &HeuristicTarget{
				Format:     "projects/%v",
				FieldNames: []string{"project"},
			},
		},
		{
			name:       "Compute: leading 'locations/global' with known collection (topics)",
			methodName: "ListLocationGlobalTopics",
			pattern:    "v1/locations/global/topics/{topic}",
			want: &HeuristicTarget{
				Format:     "locations/global/topics/%v",
				FieldNames: []string{"topic"},
			},
		},
		{
			name:       "Custom verb stripping on literals (topics:cancel)",
			methodName: "GetTopic",
			pattern:    "v1/projects/{project}/topics:cancel/{topic}",
			want: &HeuristicTarget{
				Format:     "projects/%v/topics/%v",
				FieldNames: []string{"project", "topic"},
			},
		},
	}

	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			m := &descriptorpb.MethodDescriptorProto{Name: proto.String(tst.methodName)}
			opts := &descriptorpb.MethodOptions{}
			proto.SetExtension(opts, annotations.E_Http, &annotations.HttpRule{
				Pattern: &annotations.HttpRule_Get{Get: tst.pattern},
			})
			m.Options = opts

			got, err := IdentifyHeuristicTarget(m, opts.ProtoReflect().Get(annotations.E_Http.TypeDescriptor()).Message().Interface().(*annotations.HttpRule), vocabulary)
			if err != nil {
				t.Fatalf("IdentifyHeuristicTarget failed: %v", err)
			}

			if !reflect.DeepEqual(got, tst.want) {
				t.Errorf("IdentifyHeuristicTarget(%s): got %v, want %v", tst.name, got, tst.want)
			}
		})
	}
}

func TestSplitPathSegments(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		want    []string
	}{
		{
			name:    "Standard path",
			pattern: "v1/projects/{project}/topics/{topic}",
			want:    []string{"v1", "projects", "{project}", "topics", "{topic}"},
		},
		{
			name:    "Glob heavy variable",
			pattern: "v1/{name=projects/*/locations/*/topics/*}",
			want:    []string{"v1", "{name=projects/*/locations/*/topics/*}"},
		},
		{
			name:    "Leading slash behavior (matching strings.Split)",
			pattern: "/v1/projects/{project}",
			want:    []string{"", "v1", "projects", "{project}"},
		},
	}

	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			got := splitPathSegments(tst.pattern)
			if !reflect.DeepEqual(got, tst.want) {
				t.Errorf("splitPathSegments(%s): got %v, want %v", tst.name, got, tst.want)
			}
		})
	}
}
