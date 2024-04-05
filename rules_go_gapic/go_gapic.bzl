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

load("@io_bazel_rules_go//go:def.bzl", "go_context", "go_library")
load("@rules_gapic//:gapic.bzl", "proto_custom_library", "unzipped_srcjar")

def _go_gapic_postprocessed_srcjar_impl(ctx):
    go_ctx = go_context(ctx)

    gapic_srcjar = ctx.file.gapic_srcjar
    output_main = ctx.outputs.main
    output_test = ctx.outputs.test
    output_snippets = ctx.outputs.snippets
    output_metadata = ctx.outputs.metadata

    output_dir_name = ctx.label.name
    output_dir_path = "%s/%s" % (output_main.dirname, output_dir_name)

    formatter = _get_gofmt(go_ctx)

    script = """
    unzip -q {gapic_srcjar} -d {output_dir_path}
    {formatter} -w -l {output_dir_path}
    pushd .
    cd {output_dir_path}
    zip -q -r {output_dir_name}-test.srcjar . -i ./*_test.go
    find . -name "*_test.go" -delete
    zip -q -r {output_dir_name}-snippets.srcjar . -i ./*main.go ./*snippet_metadata.*.json
    find . -name "*main.go" -delete
    find . -name "*snippet_metadata.*.json" -delete
    zip -q -r {output_dir_name}.srcjar . -i ./*.go
    find . -name "*.go" -delete
    zip -q -r {output_dir_name}-metadata.srcjar . -i ./*.json
    popd
    mv {output_dir_path}/{output_dir_name}-test.srcjar {output_test}
    mv {output_dir_path}/{output_dir_name}-snippets.srcjar {output_snippets}
    mv {output_dir_path}/{output_dir_name}.srcjar {output_main}
    mv {output_dir_path}/{output_dir_name}-metadata.srcjar {output_metadata}
    """.format(
        gapic_srcjar = gapic_srcjar.path,
        output_dir_name = output_dir_name,
        output_dir_path = output_dir_path,
        formatter = formatter.path,
        output_main = output_main.path,
        output_test = output_test.path,
        output_snippets = output_snippets.path,
        output_metadata = output_metadata.path,
    )

    ctx.actions.run_shell(
        inputs = [gapic_srcjar],
        tools = [formatter],
        command = script,
        outputs = [output_main, output_test, output_snippets, output_metadata],
    )

_go_gapic_postprocessed_srcjar = rule(
  _go_gapic_postprocessed_srcjar_impl,
  attrs = {
      "gapic_srcjar": attr.label(mandatory = True, allow_single_file = True),
      "_go_context_data": attr.label(
          default = "@io_bazel_rules_go//:go_context_data",
      ),
  },
  outputs = {
      "main": "%{name}.srcjar",
      "test": "%{name}-test.srcjar",
      "snippets": "%{name}-snippets.srcjar",
      "metadata": "%{name}-metadata.srcjar",
  },
  toolchains = ["@io_bazel_rules_go//go:toolchain"],
)

def _get_gofmt(go_ctx):
  for tool in go_ctx.sdk.tools:
      if tool.basename == "gofmt":
          return tool
  return None

def go_gapic_library(
  name,
  srcs,
  importpath,
  deps,
  release_level = "",
  grpc_service_config = None,
  service_yaml = None,
  gapic_yaml = None,
  metadata = False,
  transport = "grpc",
  diregapic = False,
  rest_numeric_enums = False,
  omit_snippets = False,
  **kwargs):

  file_args = {}

  if grpc_service_config:
    file_args[grpc_service_config] =  "grpc-service-config"

  if service_yaml:
    file_args[service_yaml] = "api-service-config"

  plugin_args = [
    "go-gapic-package={}".format(importpath),
    "transport={}".format(transport)
  ]

  if release_level:
    plugin_args.append("release-level={}".format(release_level))

  if metadata:
    plugin_args.append("metadata")

  if diregapic:
    plugin_args.append("diregapic")

  if rest_numeric_enums:
    plugin_args.append("rest-numeric-enums")

  if omit_snippets:
    plugin_args.append("omit-snippets")

  srcjar_name = name+"_srcjar"
  raw_srcjar_name = srcjar_name+"_raw"
  output_suffix = ".srcjar"

  proto_custom_library(
    name = raw_srcjar_name,
    deps = srcs,
    plugin = Label("//cmd/protoc-gen-go_gapic"),
    plugin_args = plugin_args,
    plugin_file_args = file_args,
    output_type = "go_gapic",
    output_suffix = output_suffix,
    **kwargs
  )

  _go_gapic_postprocessed_srcjar(
    name = srcjar_name,
    gapic_srcjar = ":%s" % raw_srcjar_name,
    **kwargs
  )

  actual_deps = deps + [
    "@com_github_google_uuid//:go_default_library",
    "@com_github_googleapis_gax_go_v2//:go_default_library",
    "@com_github_googleapis_gax_go_v2//apierror:go_default_library",
    "@io_opentelemetry_go_contrib_instrumentation_net_http_otelhttp//:go_default_library",
    "@io_opentelemetry_go_contrib_instrumentation_google_golang_org_grpc_otelgrpc//:go_default_library",
    "@org_golang_google_api//googleapi:go_default_library",
    "@org_golang_google_api//option:go_default_library",
    "@org_golang_google_api//option/internaloption:go_default_library",
    "@org_golang_google_api//iterator:go_default_library",
    "@org_golang_google_api//transport/http:go_default_library",
    "@org_golang_google_api//transport/grpc:go_default_library",
    "@org_golang_google_grpc//:go_default_library",
    "@org_golang_google_grpc//codes:go_default_library",
    "@org_golang_google_grpc//metadata:go_default_library",
    "@org_golang_google_grpc//status:go_default_library",
    "@org_golang_google_protobuf//proto:go_default_library",
    "@org_golang_google_protobuf//encoding/protojson:go_default_library",
    "@org_golang_google_protobuf//types/known/emptypb:go_default_library",
    "@org_golang_google_protobuf//types/known/wrapperspb:go_default_library",
    "@org_golang_google_protobuf//types/known/fieldmaskpb:go_default_library",
    "@org_golang_google_protobuf//types/known/timestamppb:go_default_library",
    "@org_golang_google_protobuf//types/known/anypb:go_default_library",
    "@org_golang_google_protobuf//types/known/structpb:go_default_library",
    "@org_golang_google_protobuf//types/known/durationpb:go_default_library",
  ]

  main_file = ":%s" % srcjar_name + output_suffix
  main_dir = "%s_main" % srcjar_name

  unzipped_srcjar(
    name = main_dir,
    srcjar = main_file,
    extension = ".go",
  )

  # Strip the trailing package alias so that this
  # generated library can shade com_google_cloud_go
  # go_library targets.
  imp = importpath[:importpath.index(";")]
  go_library(
    name = name,
    srcs = [":%s" % main_dir],
    deps = actual_deps,
    importpath = imp,
  )

  test_file = ":%s-test.srcjar" % srcjar_name
  test_dir = "%s_test" % srcjar_name

  unzipped_srcjar(
    name = test_dir,
    srcjar = test_file,
    extension = ".go",
    **kwargs
  )

  snippets_file = ":%s-snippets.srcjar" % srcjar_name
  snippets_dir = "%s_snippets" % srcjar_name

  unzipped_srcjar(
    name = snippets_dir,
    srcjar = snippets_file,
    extension = ".go",
    **kwargs
  )

  if metadata:
    metadata_file = ":%s-metadata.srcjar" % srcjar_name
    metadata_dir = "%s_metadata" % srcjar_name

    unzipped_srcjar(
      name = metadata_dir,
      srcjar = metadata_file,
      extension = ".json",
      **kwargs
    )
