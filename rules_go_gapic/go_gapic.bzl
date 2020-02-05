# Copyright 2019 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@com_google_api_codegen//rules_gapic:gapic.bzl", "proto_custom_library", "unzipped_srcjar")

def go_gapic_library(
  name,
  srcs,
  importpath,
  deps,
  release_level = "",
  grpc_service_config = None,
  service_yaml = None,
  samples = [],
  sample_only = False,
  **kwargs):

  output_suffix = ".srcjar"
  file_args = {}

  if grpc_service_config:
    file_args[grpc_service_config] =  "grpc-service-config"

  if service_yaml:
    file_args[service_yaml] = "gapic-service-config"

  if samples:
    for path in samples:
        file_args[path] = "sample"

  plugin_args = [
    "go-gapic-package={}".format(importpath),
    "release-level={}".format(release_level),
  ]

  if sample_only:
    plugin_args.append("sample-only={}".format(sample_only))

  proto_custom_library(
    name = name,
    deps = srcs,
    plugin = Label("//cmd/protoc-gen-go_gapic"),
    plugin_args = plugin_args,
    plugin_file_args = file_args,
    output_type = "go_gapic",
    output_suffix = output_suffix,
    **kwargs
  )

  # This dependecy list was copied directly from gapic-generator/rules_gapic/go.
  # Ideally, this should be a common list used by macros in both repos.
  #
  # TODO(ndietz) de-dupe this dep list with gapic-generator/rules_gapic/go
  actual_deps = deps + [
    "@com_github_googleapis_gax_go//v2:go_default_library",
    "@org_golang_google_api//option:go_default_library",
    "@org_golang_google_api//iterator:go_default_library",
    "@org_golang_google_api//transport:go_default_library",
    "@org_golang_google_grpc//:go_default_library",
    "@org_golang_google_grpc//codes:go_default_library",
    "@org_golang_google_grpc//metadata:go_default_library",
    "@com_github_golang_protobuf//proto:go_default_library",
    "@com_github_golang_protobuf//ptypes:go_default_library",
    "@com_github_golang_protobuf//ptypes/empty:go_default_library",
    "@com_github_golang_protobuf//ptypes/timestamp:go_default_library",
    "@org_golang_google_genproto//protobuf/field_mask:go_default_library",
    "@com_google_googleapis//google/rpc:status_go_proto",
    "@org_golang_google_grpc//status:go_default_library",
  ]

  main_file = ":%s" % name + output_suffix
  main_dir = "%s_main" % name

  unzipped_srcjar(
    name = main_dir,
    srcjar = main_file,
    extension = ".go",
  )

  go_library(
    name = name+"_pkg",
    srcs = [":%s" % main_dir],
    deps = actual_deps,
    importpath = importpath,
  )
