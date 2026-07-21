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

func TestExtractParams(t *testing.T) {
	tests := []struct {
		pattern string
		want    []string
	}{
		{
			pattern: "projects/{project}/topics/{topic}",
			want:    []string{"project", "topic"},
		},
		{
			pattern: "projects/{project}/locations/{location}/instances/{instance=*}",
			want:    []string{"project", "location", "instance"},
		},
		{
			pattern: "global/networks/{network}",
			want:    []string{"network"},
		},
	}

	for _, tt := range tests {
		got := extractParams(tt.pattern)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("extractParams(%q) = %v, want %v", tt.pattern, got, tt.want)
		}
	}
}

func TestExtractCollectionName(t *testing.T) {
	tests := []struct {
		resType string
		pattern string
		want    string
	}{
		{
			resType: "pubsub.googleapis.com/Topic",
			pattern: "projects/{project}/topics/{topic}",
			want:    "topics",
		},
		{
			resType: "compute.googleapis.com/Instance",
			pattern: "projects/{project}/zones/{zone}/instances/{instance}",
			want:    "instances",
		},
	}

	for _, tt := range tests {
		got := extractCollectionName(tt.resType, tt.pattern)
		if got != tt.want {
			t.Errorf("extractCollectionName(%q, %q) = %q, want %q", tt.resType, tt.pattern, got, tt.want)
		}
	}
}

func TestCollectResourceSpecs(t *testing.T) {
	rd := &annotations.ResourceDescriptor{
		Type:    "pubsub.googleapis.com/Topic",
		Pattern: []string{"projects/{project}/topics/{topic}"},
	}

	msgOpts := &descriptorpb.MessageOptions{}
	proto.SetExtension(msgOpts, annotations.E_Resource, rd)

	msg := &descriptorpb.DescriptorProto{
		Name:    proto.String("Topic"),
		Options: msgOpts,
	}

	file := &descriptorpb.FileDescriptorProto{
		Name:        proto.String("pubsub.proto"),
		MessageType: []*descriptorpb.DescriptorProto{msg},
	}

	g := &generator{}
	specs := g.collectResourceSpecs([]*descriptorpb.FileDescriptorProto{file})

	spec, ok := specs["pubsub.googleapis.com/Topic"]
	if !ok {
		t.Fatalf("expected resource spec for pubsub.googleapis.com/Topic, got none")
	}

	if spec.Collection != "topics" {
		t.Errorf("got collection %q, want topics", spec.Collection)
	}

	wantParams := []string{"project", "topic"}
	if !reflect.DeepEqual(spec.Params, wantParams) {
		t.Errorf("got params %v, want %v", spec.Params, wantParams)
	}
}
