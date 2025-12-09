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

	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

const (
	// ServiceURL is the hostname of the service.
	//
	// Example:
	// https://github.com/googleapis/googleapis/blob/f7df662a24c56ecaab79cb7d808fed4d2bb4981d/google/cloud/secretmanager/v1/secretmanager_v1.yaml#L3
	ServiceURL = "secretmanager.googleapis.com"

	// ServiceName is the name of the service.
	//
	// Example:
	// https://github.com/googleapis/googleapis/blob/f7df662a24c56ecaab79cb7d808fed4d2bb4981d/google/cloud/secretmanager/v1/service.proto#L46
	ServiceName = "SecretManagerService"

	// ServiceTitle is the name of the service provided in the service YAML
	// file.
	//
	// Example:
	// https://github.com/googleapis/googleapis/blob/f7df662a24c56ecaab79cb7d808fed4d2bb4981d/google/cloud/secretmanager/v1/secretmanager_v1.yaml#L4
	ServiceTitle = "Secret Manager API"

	// ServiceDocumentationSummary is the description of the service provided in the
	// service YAML file.
	//
	// Example:
	// https://github.com/googleapis/googleapis/blob/f7df662a24c56ecaab79cb7d808fed4d2bb4981d/google/cloud/secretmanager/v1/secretmanager_v1.yaml#L11
	ServiceDocumentationSummary = "Stores sensitive data such as API keys, passwords, and certificates. Provides convenience while improving security."

	// ServiceOAuthScope is the OAuth 2.0 scope(s) that is needed to request
	// access to the service's API, as defined in the service YAML file.
	//
	// Example:
	// https://github.com/googleapis/googleapis/blob/f7df662a24c56ecaab79cb7d808fed4d2bb4981d/google/cloud/secretmanager/v1/secretmanager_v1.yaml#L34
	ServiceOAuthScope = "https://www.googleapis.com/auth/cloud-platform"

	// CreateMethod is the name of the RPC method for creating a resource.
	// The same name is used for the proto RPC method and the Go method.
	//
	// Example:
	// https://pkg.go.dev/cloud.google.com/go/secretmanager@v1.13.0/apiv1#Client.CreateSecret
	// https://github.com/googleapis/googleapis/blob/f7df662a24c56ecaab79cb7d808fed4d2bb4981d/google/cloud/secretmanager/v1/service.proto#L62
	CreateMethod = "CreateSecret"

	// CreateRequest is the name of the request for creating a resource.
	// The same name is used for the proto message and the Go type.
	//
	// Example:
	// https://pkg.go.dev/cloud.google.com/go/secretmanager/apiv1/secretmanagerpb#CreateSecretRequest
	// https://github.com/googleapis/googleapis/blob/f7df662a24c56ecaab79cb7d808fed4d2bb4981d/google/cloud/secretmanager/v1/service.proto#L330
	CreateRequest = "CreateSecretRequest"

	// GetMethod is the name of the RPC method used to fetch a resource.
	// The same name is used for the proto RPC method and the Go method.
	//
	// Example:
	// https://pkg.go.dev/cloud.google.com/go/secretmanager/apiv1#Client.GetSecret
	// https://github.com/googleapis/googleapis/blob/f7df662a24c56ecaab79cb7d808fed4d2bb4981d/google/cloud/secretmanager/v1/service.proto#L90
	GetMethod = "GetSecret"

	// GetRequest is the name of the request for fetching a resource.
	// The same name is used for the proto message and the Go type.
	//
	// A GetRequest often contains `google.api.resource_reference`, in order to
	// reference the name of the resource (see https://aip.dev/4231#referencing-other-resources).
	//
	// Example:
	// https://pkg.go.dev/cloud.google.com/go/secretmanager/apiv1/secretmanagerpb#GetSecretRequest
	// https://github.com/googleapis/googleapis/blob/f7df662a24c56ecaab79cb7d808fed4d2bb4981d/google/cloud/secretmanager/v1/service.proto#L367
	GetRequest = "GetSecretRequest"

	// Resource is the name of the resource returned by a Get or Create request.
	//
	// A resource message often contains a `google.api.resource` option with a
	// type and pattern (see https://aip.dev/4231#resource-messages).
	//
	// Example:
	// https://pkg.go.dev/cloud.google.com/go/secretmanager/apiv1/secretmanagerpb#Secret
	// https://github.com/googleapis/googleapis/blob/f7df662a24c56ecaab79cb7d808fed4d2bb4981d/google/cloud/secretmanager/v1/resources.proto#L40
	Resource = "Secret"

	// ListMethod is the name of the method for listing resources.
	//
	// Example:
	// https://pkg.go.dev/cloud.google.com/go/secretmanager/apiv1/secretmanagerpb#ListSecrets
	// https://github.com/googleapis/googleapis/blob/f7df662a24c56ecaab79cb7d808fed4d2bb4981d/google/cloud/secretmanager/v1/service.proto#L52
	ListMethod = "ListSecrets"

	// ListRequest is the name of the request for listing resources.
	//
	// A ListRequest often contains `page_size` and `page_token` fields (see
	// https://aip.dev/158).
	//
	// Example:
	// https://pkg.go.dev/cloud.google.com/go/secretmanager/apiv1/secretmanagerpb#ListSecretsRequest
	// https://github.com/googleapis/googleapis/blob/f7df662a24c56ecaab79cb7d808fed4d2bb4981d/google/cloud/secretmanager/v1/service.proto#L274
	ListRequest = "ListSecretsRequest"

	// ListResponse is the name of the resource returned by a List request.
	//
	// A ListResponse message often contains the Resource, a `next_page_token`
	// and sometimes a `total_size` field (https://aip.dev/158).
	// type and pattern (see https://aip.dev/4231#resource-messages).
	//
	// Example:
	// https://pkg.go.dev/cloud.google.com/go/secretmanager/apiv1/secretmanagerpb#ListSecretsResponse
	// https://github.com/googleapis/googleapis/blob/f7df662a24c56ecaab79cb7d808fed4d2bb4981d/google/cloud/secretmanager/v1/service.proto#L304
	ListResponse = "ListSecretsResponse"

	// CreateMethodWithSettings is a fake method for the purpose of testing
	// the method_settings functionality in the publishing YAML.
	//
	// The actual secretmanager publishing does not use this field:
	// https://github.com/googleapis/googleapis/blob/f7df662a24c56ecaab79cb7d808fed4d2bb4981d/google/cloud/secretmanager/v1/secretmanager_v1.yaml#L44
	CreateMethodWithSettings = "CreateSecretWithMethodSettings"
)

const (
	// Year is the year used in the copyright header for package documentation
	// and samples.
	Year = 2024
)

const (
	// ProtoPackagePath is the package path of the proto file.
	//
	// Example:
	// https://github.com/googleapis/googleapis/blob/f7df662a24c56ecaab79cb7d808fed4d2bb4981d/google/cloud/secretmanager/v1/resources.proto#L17
	ProtoPackagePath = "google.cloud.secretmanager.v1"

	// ProtoVersion is the major version as defined in the protofile.
	ProtoVersion = "v1"
)

const (
	// GoPackageName is the package name for the auto-generated Go package.
	//
	// Example:
	// https://pkg.go.dev/cloud.google.com/go/secretmanager/apiv1
	GoPackageName = "secretmanager"

	// GoPackagePath is the package import path for the auto-generated Go
	// package.
	//
	// Example:
	// https://pkg.go.dev/cloud.google.com/go/secretmanager/apiv1
	GoPackagePath = "cloud.google.com/go/secretmanager/apiv1"

	// GoProtoPackageName is the package name of the auto-generated proto
	// package, which is imported by package at GoPackagePath. This name is
	// derived from the value following the ";" `go_package` in the proto file.
	//
	// Example:
	// https://pkg.go.dev/cloud.google.com/go/secretmanager/apiv1/secretmanagerpb
	// https://github.com/googleapis/googleapis/blob/f7df662a24c56ecaab79cb7d808fed4d2bb4981d/google/cloud/secretmanager/v1/resources.proto#L26
	GoProtoPackageName = "secretmanagerpb"

	// GoProtoPackagePath is the package import path of the auto-generated proto
	// package.  This name is derived from the value before the ";"
	// `go_package` in the proto file.
	//
	// Example:
	// https://pkg.go.dev/cloud.google.com/go/secretmanager/apiv1/secretmanagerpb
	// https://github.com/googleapis/googleapis/blob/f7df662a24c56ecaab79cb7d808fed4d2bb4981d/google/cloud/secretmanager/v1/resources.proto#L26
	GoProtoPackagePath = "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"

	// GoVersion is the version used in the package path for versioning the Go
	// module containing the package.
	GoVersion = "apiv1"
)

const (
	// SnippetsDirectory is the directory path containing generated snippets.
	//
	// Example:
	// https://github.com/googleapis/google-cloud-go/tree/main/internal/generated/snippets/secretmanager/apiv1
	SnippetsDirectory = "cloud.google.com/go/internal/generated/snippets/secretmanager/apiv1"
)

// DescriptorInfoTypeName constructs the name format used by g.descInfo.Type.
func DescriptorInfoTypeName(typ string) string {
	return fmt.Sprintf(".%s.%s", ProtoPackagePath, typ)
}

// Service returns a service descriptor using the sample values.
func Service() *descriptorpb.ServiceDescriptorProto {
	return &descriptorpb.ServiceDescriptorProto{
		Name: proto.String(ServiceName),
		Method: []*descriptorpb.MethodDescriptorProto{
			{
				Name:       proto.String(CreateMethod),
				InputType:  proto.String(DescriptorInfoTypeName(CreateRequest)),
				OutputType: proto.String(DescriptorInfoTypeName(Resource)),
			},
			{
				Name:       proto.String(GetMethod),
				InputType:  proto.String(DescriptorInfoTypeName(GetRequest)),
				OutputType: proto.String(DescriptorInfoTypeName(Resource)),
			},
			{
				Name:       proto.String(ListMethod),
				InputType:  proto.String(DescriptorInfoTypeName(ListRequest)),
				OutputType: proto.String(DescriptorInfoTypeName(ListResponse)),
			},
		},
	}
}

// InputType returns an input type for a method.
func InputType(input string) *descriptorpb.DescriptorProto {
	return &descriptorpb.DescriptorProto{
		Name: proto.String(input),
	}
}

// OutputType returns an output type for a method.
func OutputType(output string) *descriptorpb.DescriptorProto {
	return &descriptorpb.DescriptorProto{
		Name: proto.String(output),
	}
}

// File returns a proto file.
func File() *descriptorpb.FileDescriptorProto {
	return &descriptorpb.FileDescriptorProto{
		Options: &descriptorpb.FileOptions{
			GoPackage: proto.String(GoProtoPackagePath),
		},
		Package: proto.String(ProtoPackagePath),
	}
}

// ApiVersionOptions returns a ServiceOptions instance with the
// google.api.api_version extension set to the provided version string.
func ApiVersionOptions(v string) *descriptorpb.ServiceOptions {
	o := &descriptorpb.ServiceOptions{}
	proto.SetExtension(o, annotations.E_ApiVersion, v)
	return o
}
