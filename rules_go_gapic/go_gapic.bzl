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

load("@com_google_api_codegen//rules_gapic:gapic.bzl", "proto_custom_library", "unzipped_srcjar")
load("@io_bazel_rules_go//go:def.bzl", "go_context", "go_library")

def _go_gapic_postprocessed_srcjar_impl(ctx):
    go_ctx = go_context(ctx)

    gapic_srcjar = ctx.file.gapic_srcjar
    output_main = ctx.outputs.main
    output_test = ctx.outputs.test

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
    zip -q -r {output_dir_name}.srcjar . -i ./*.go
    popd
    mv {output_dir_path}/{output_dir_name}-test.srcjar {output_test}
    mv {output_dir_path}/{output_dir_name}.srcjar {output_main}
    """.format(
        gapic_srcjar = gapic_srcjar.path,
        output_dir_name = output_dir_name,
        output_dir_path = output_dir_path,
        formatter = formatter.path,
        output_main = output_main.path,
        output_test = output_test.path,
    )

    ctx.actions.run_shell(
        inputs = [gapic_srcjar],
        tools = [formatter],
        command = script,
        outputs = [output_main, output_test],
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
  samples = [],
  sample_only = False,
  **kwargs):

  file_args = {}

  if grpc_service_config:
    file_args[grpc_service_config] =  "grpc-service-config"

  if service_yaml:
    file_args[service_yaml] = "gapic-service-config"
  
  if gapic_yaml:
    file_args[gapic_yaml] = "gapic-config"

  if samples:
    for path in samples:
        file_args[path] = "sample"

  plugin_args = [
    "go-gapic-package={}".format(importpath),
  ]

  if release_level:
    plugin_args.append("release-level={}".format(release_level))

  if sample_only:
    plugin_args.append("sample-only")

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
    "@com_github_googleapis_gax_go_v2//:go_default_library",
    "@org_golang_google_api//option:go_default_library",
    "@org_golang_google_api//option/internaloption:go_default_library",
    "@org_golang_google_api//iterator:go_default_library",
    "@org_golang_google_api//transport/grpc:go_default_library",
    "@org_golang_google_grpc//:go_default_library",
    "@org_golang_google_grpc//codes:go_default_library",
    "@org_golang_google_grpc//metadata:go_default_library",
    "@org_golang_google_grpc//status:go_default_library",
    "@com_github_golang_protobuf//proto:go_default_library",
    "@com_github_golang_protobuf//ptypes:go_default_library_gen",
    "@io_bazel_rules_go//proto/wkt:empty_go_proto",
    "@io_bazel_rules_go//proto/wkt:field_mask_go_proto",
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
