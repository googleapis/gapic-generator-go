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

	"google.golang.org/genproto/googleapis/api/serviceconfig"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/apipb"
)

const (
	// ServiceURL is the hostname of the service.
	//
	// Example:
	// https://github.com/googleapis/googleapis/blob/f7df662a24c56ecaab79cb7d808fed4d2bb4981d/google/cloud/bigquery/migration/v2/migration_service.proto#L36
	ServiceURL = "bigquerymigration.googleapis.com"

	// ServiceName is the name of the service.
	//
	// Example:
	// https://github.com/googleapis/googleapis/blob/f7df662a24c56ecaab79cb7d808fed4d2bb4981d/google/cloud/bigquery/migration/v2/migration_service.proto#L35
	ServiceName = "MigrationService"

	// CreateMethod is the name of the RPC method for creating a resource.
	// The same name is used for the proto RPC method and the Go method.
	//
	// Example:
	// https://pkg.go.dev/cloud.google.com/go/bigquery/migration/apiv2/migrationpb#MigrationServiceClient.CreateMigrationWorkflow
	// https://github.com/googleapis/googleapis/blob/f7df662a24c56ecaab79cb7d808fed4d2bb4981d/google/cloud/bigquery/migration/v2/migration_service.proto#L41
	CreateMethod = "CreateMigrationWorkflow"

	// CreateRequest is the name of the request for creating a resource.
	// The same name is used for the proto message and the Go type.
	//
	// Example:
	// https://pkg.go.dev/cloud.google.com/go/bigquery/migration/apiv2/migrationpb#CreateMigrationWorkflowRequest
	// https://github.com/googleapis/googleapis/blob/f7df662a24c56ecaab79cb7d808fed4d2bb4981d/google/cloud/bigquery/migration/v2/migration_service.proto#L110
	CreateRequest = "CreateMigrationWorkflowRequest"

	// GetMethod is the name of the RPC method used to fetch a resource.
	// The same name is used for the proto RPC method and the Go method.
	//
	// Example:
	// https://pkg.go.dev/cloud.google.com/go/bigquery/migration/apiv2/migrationpb#MigrationServiceClient.GetMigrationWorkflow
	// https://github.com/googleapis/googleapis/blob/f7df662a24c56ecaab79cb7d808fed4d2bb4981d/google/cloud/bigquery/migration/v2/migration_service.proto#L51
	GetMethod = "GetMigrationWorkflow"

	// GetRequest is the name of the request for fetching a resource.
	// The same name is used for the proto message and the Go type.
	//
	// A GetRequest often contains `google.api.resource_reference`, in order to
	// reference the name of the resource (see https://aip.dev/4231#referencing-other-resources).
	//
	// Example:
	// https://pkg.go.dev/cloud.google.com/go/bigquery/migration/apiv2/migrationpb#GetMigrationWorkflowRequest
	// https://github.com/googleapis/googleapis/blob/f7df662a24c56ecaab79cb7d808fed4d2bb4981d/google/cloud/bigquery/migration/v2/migration_service.proto#L126
	GetRequest = "GetMigrationWorkflowRequest"

	// Resource is the name of the resource returned by a Get or Create request.
	//
	// A resource message often contains a `google.api.resource` option with a
	// type and pattern (see https://aip.dev/4231#resource-messages).
	//
	// Example:
	// https://pkg.go.dev/cloud.google.com/go/bigquery/migration/apiv2/migrationpb#MigrationWorkflow
	// https://github.com/googleapis/googleapis/blob/master/google/cloud/bigquery/migration/v2alpha/migration_entities.proto#L38
	Resource = "MigrationWorkflow"
)

const (

	// ProtoServiceName is the fully qualified name of service.
	//
	// Example:
	// https://github.com/googleapis/googleapis/blob/f7df662a24c56ecaab79cb7d808fed4d2bb4981d/google/cloud/bigquery/migration/v2/migration_service.proto#L35.
	ProtoServiceName = "google.cloud.bigquery.migration.v2.MigrationService"

	// ProtoPackagePath is the package path of the proto file.
	//
	// Example:
	// https://github.com/googleapis/googleapis/blob/f7df662a24c56ecaab79cb7d808fed4d2bb4981d/google/cloud/bigquery/migration/v2/migration_service.proto#L17
	ProtoPackagePath = "google.cloud.bigquery.migration.v2"

	// ProtoVersion is the major version as defined in the protofile.
	ProtoVersion = "v2"
)

const (
	// GoPackageName is the package name for the auto-generated Go package.
	//
	// Example:
	// https://pkg.go.dev/cloud.google.com/go/bigquery/migration/apiv2
	GoPackageName = "migration"

	// GoPackagePath is the package import path for the auto-generated Go
	// package.
	//
	// Example:
	// https://pkg.go.dev/cloud.google.com/go/bigquery/migration/apiv2
	// https://github.com/googleapis/googleapis/blob/f7df662a24c56ecaab79cb7d808fed4d2bb4981d/google/cloud/bigquery/migration/v2/migration_service.proto#L28
	GoPackagePath = "cloud.google.com/go/bigquery/migration/apiv2"

	// GoProtoPackageName is the package name of the auto-generated proto
	// package, which is imported by package at GoPackagePath. This name is
	// derived from the value following the ";" `go_package` in the proto file.
	//
	// Example:
	// https://pkg.go.dev/cloud.google.com/go/bigquery/migration/apiv2/migrationpb
	// https://github.com/googleapis/googleapis/blob/f7df662a24c56ecaab79cb7d808fed4d2bb4981d/google/cloud/bigquery/migration/v2/migration_service.proto#L28
	GoProtoPackageName = "migrationpb"

	// GoProtoPackagePath is the package import path of the auto-generated proto
	// package.  This name is derived from the value before the ";"
	// `go_package` in the proto file.
	//
	// Example:
	// https://pkg.go.dev/cloud.google.com/go/bigquery/migration/apiv2/migrationpb.
	// https://github.com/googleapis/googleapis/blob/f7df662a24c56ecaab79cb7d808fed4d2bb4981d/google/cloud/bigquery/migration/v2/migration_service.proto#L28
	GoProtoPackagePath = "cloud.google.com/go/bigquery/migration/apiv2/migrationpb"

	// GoVersion is the version used in the package path for versioning the Go
	// module containing the package.
	GoVersion = "apiv2"
)

// DescriptorInfoTypeName constructs the name format used by g.descInfo.Type.
func DescriptorInfoTypeName(typ string) string {
	return fmt.Sprintf(".%s.%s", ProtoPackagePath, typ)
}

// ServiceConfig returns service config information.
func ServiceConfig() *serviceconfig.Service {
	return &serviceconfig.Service{
		Apis: []*apipb.Api{
			{Name: ProtoServiceName},
		},
	}
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
