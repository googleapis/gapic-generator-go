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

load(
    "@bazel_gazelle//:deps.bzl",

    gazelle_go_repository = "go_repository",
)
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

_rules_gapic_version = "0.14.1"

def com_googleapis_gapic_generator_go_repositories():
    _maybe(
        http_archive,
        name = "rules_gapic",
        strip_prefix = "rules_gapic-%s" % _rules_gapic_version,
        urls = ["https://github.com/googleapis/rules_gapic/archive/v%s.tar.gz" % _rules_gapic_version],
    )
    go_repository(
        name = "com_github_bufbuild_protocompile",
        importpath = "github.com/bufbuild/protocompile",
        sum = "h1:iA73zAf/fyljNjQKwYzUHD6AD4R8KMasmwa/FBatYVw=",
        version = "v0.14.1",
    )
    go_repository(
        name = "com_github_cespare_xxhash_v2",
        importpath = "github.com/cespare/xxhash/v2",
        sum = "h1:UL815xU9SqsFlibzuggzjXhog7bL6oX9BbNZnL2UFvs=",
        version = "v2.3.0",
    )
    go_repository(
        name = "com_github_cncf_xds_go",
        importpath = "github.com/cncf/xds/go",
        sum = "h1:6xNmx7iTtyBRev0+D/Tv1FZd4SCg8axKApyNyRsAt/w=",
        version = "v0.0.0-20251210132809-ee656c7534f5",
    )
    go_repository(
        name = "com_github_creack_pty",
        importpath = "github.com/creack/pty",
        sum = "h1:uDmaGzcdjhF4i/plgjmEsriH11Y0o7RKapEf/LDaM3w=",
        version = "v1.1.9",
    )
    go_repository(
        name = "com_github_davecgh_go_spew",
        importpath = "github.com/davecgh/go-spew",
        sum = "h1:vj9j/u1bqnvCEfJOwUhtlOARqs3+rkHYY13jYWTU97c=",
        version = "v1.1.1",
    )
    go_repository(
        name = "com_github_envoyproxy_go_control_plane",
        importpath = "github.com/envoyproxy/go-control-plane",
        sum = "h1:hbG2kr4RuFj222B6+7T83thSPqLjwBIfQawTkC++2HA=",
        version = "v0.14.0",
    )
    go_repository(
        name = "com_github_envoyproxy_go_control_plane_envoy",
        importpath = "github.com/envoyproxy/go-control-plane/envoy",
        sum = "h1:yg/JjO5E7ubRyKX3m07GF3reDNEnfOboJ0QySbH736g=",
        version = "v1.36.0",
    )
    go_repository(
        name = "com_github_envoyproxy_go_control_plane_ratelimit",
        importpath = "github.com/envoyproxy/go-control-plane/ratelimit",
        sum = "h1:/G9QYbddjL25KvtKTv3an9lx6VBE2cnb8wp1vEGNYGI=",
        version = "v0.1.0",
    )
    go_repository(
        name = "com_github_envoyproxy_protoc_gen_validate",
        importpath = "github.com/envoyproxy/protoc-gen-validate",
        sum = "h1:TvGH1wof4H33rezVKWSpqKz5NXWg5VPuZ0uONDT6eb4=",
        version = "v1.3.0",
    )
    go_repository(
        name = "com_github_felixge_httpsnoop",
        importpath = "github.com/felixge/httpsnoop",
        sum = "h1:NFTV2Zj1bL4mc9sqWACXbQFVBBg2W3GPvqp8/ESS2Wg=",
        version = "v1.0.4",
    )
    go_repository(
        name = "com_github_fsnotify_fsnotify",
        importpath = "github.com/fsnotify/fsnotify",
        sum = "h1:2Ml+OJNzbYCTzsxtv8vKSFD9PbJjmhYF14k/jKC7S9k=",
        version = "v1.9.0",
    )
    go_repository(
        name = "com_github_ghodss_yaml",
        build_directives = [
            "gazelle:resolve go go gopkg.in/yaml.v2 @in_gopkg_yaml_v2//:go_default_library",
        ],
        importpath = "github.com/ghodss/yaml",
        sum = "h1:wQHKEahhL6wmXdzwWG11gIVCkOv05bNOh+Rxn0yngAk=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_go_jose_go_jose_v4",
        importpath = "github.com/go-jose/go-jose/v4",
        sum = "h1:CVLmWDhDVRa6Mi/IgCgaopNosCaHz7zrMeF9MlZRkrs=",
        version = "v4.1.3",
    )
    go_repository(
        name = "com_github_go_logr_logr",
        importpath = "github.com/go-logr/logr",
        sum = "h1:CjnDlHq8ikf6E492q6eKboGOC0T8CDaOvkHCIg8idEI=",
        version = "v1.4.3",
    )
    go_repository(
        name = "com_github_go_logr_stdr",
        importpath = "github.com/go-logr/stdr",
        sum = "h1:hSWxHoqTgW2S2qGc0LTAI563KZ5YKYRhT3MFKZMbjag=",
        version = "v1.2.2",
    )
    go_repository(
        name = "com_github_go_viper_mapstructure_v2",
        importpath = "github.com/go-viper/mapstructure/v2",
        sum = "h1:EBsztssimR/CONLSZZ04E8qAkxNYq4Qp9LvH92wZUgs=",
        version = "v2.4.0",
    )
    go_repository(
        name = "com_github_golang_glog",
        importpath = "github.com/golang/glog",
        sum = "h1:DrW6hGnjIhtvhOIiAKT6Psh/Kd/ldepEa81DKeiRJ5I=",
        version = "v1.2.5",
    )
    go_repository(
        name = "com_github_golang_groupcache",
        importpath = "github.com/golang/groupcache",
        sum = "h1:oI5xCqsCo564l8iNU+DwB5epxmsaqB+rhGL0m5jtYqE=",
        version = "v0.0.0-20210331224755-41bb18bfe9da",
    )
    go_repository(
        name = "com_github_golang_protobuf",
        importpath = "github.com/golang/protobuf",
        sum = "h1:i7eJL8qZTpSEXOPTxNKhASYpMn+8e5Q6AdndVa1dWek=",
        version = "v1.5.4",
    )
    go_repository(
        name = "com_github_golang_snappy",
        importpath = "github.com/golang/snappy",
        sum = "h1:yAGX7huGHXlcLOEtBnF4w7FQwA26wojNCwOYAEhLjQM=",
        version = "v0.0.4",
    )
    go_repository(
        name = "com_github_google_go_cmp",
        importpath = "github.com/google/go-cmp",
        sum = "h1:wk8382ETsv4JYUZwIsn6YpYiWiBsYLSJiTsyBybVuN8=",
        version = "v0.7.0",
    )
    go_repository(
        name = "com_github_google_go_pkcs11",
        importpath = "github.com/google/go-pkcs11",
        sum = "h1:PVRnTgtArZ3QQqTGtbtjtnIkzl2iY2kt24yqbrf7td8=",
        version = "v0.3.0",
    )
    go_repository(
        name = "com_github_google_martian_v3",
        importpath = "github.com/google/martian/v3",
        sum = "h1:DIhPTQrbPkgs2yJYdXU/eNACCG5DVQjySNRNlflZ9Fc=",
        version = "v3.3.3",
    )
    go_repository(
        name = "com_github_google_s2a_go",
        importpath = "github.com/google/s2a-go",
        sum = "h1:LGD7gtMgezd8a/Xak7mEWL0PjoTQFvpRudN895yqKW0=",
        version = "v0.1.9",
    )
    go_repository(
        name = "com_github_google_uuid",
        importpath = "github.com/google/uuid",
        sum = "h1:NIvaJDMOsjHA8n1jAhLSgzrAzy1Hgr+hNrb57e+94F0=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_github_googleapis_enterprise_certificate_proxy",
        importpath = "github.com/googleapis/enterprise-certificate-proxy",
        sum = "h1:yh8ncqsbUY4shRD5dA6RlzjJaT4hi3kII+zYw8wmLb8=",
        version = "v0.3.14",
    )
    go_repository(
        name = "com_github_googleapis_gapic_showcase",
        importpath = "github.com/googleapis/gapic-showcase",
        sum = "h1:v7LN2yRsE/VZu9LbAMpTTvmTLKZVKMVAydzJtYlK74E=",
        version = "v0.39.0",
    )
    go_repository(
        name = "com_github_googleapis_gax_go_v2",
        build_directives = [
            "gazelle:resolve go google.golang.org/genproto/googleapis/rpc/errdetails @org_golang_google_genproto_googleapis_rpc//errdetails",
            "gazelle:resolve proto go google/rpc/code.proto @com_google_googleapis//google/rpc:code_go_proto",
            "gazelle:resolve proto proto google/rpc/code.proto @com_google_googleapis//google/rpc:code_proto",
        ],
        importpath = "github.com/googleapis/gax-go/v2",
        sum = "h1:PjIWBpgGIVKGoCXuiCoP64altEJCj3/Ei+kSU5vlZD4=",
        version = "v2.22.0",
    )
    go_repository(
        name = "com_github_googleapis_grpc_fallback_go",
        importpath = "github.com/googleapis/grpc-fallback-go",
        sum = "h1:tEDqZnKGKQpYrmEuu3VVBTw3pijHJsKz/Lu2U0L9AV0=",
        version = "v0.1.4",
    )
    go_repository(
        name = "com_github_googlecloudplatform_opentelemetry_operations_go_detectors_gcp",
        importpath = "github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp",
        sum = "h1:DHa2U07rk8syqvCge0QIGMCE1WxGj9njT44GH7zNJLQ=",
        version = "v1.31.0",
    )
    go_repository(
        name = "com_github_googlecloudplatform_opentelemetry_operations_go_exporter_metric",
        importpath = "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric",
        sum = "h1:owcC2UnmsZycprQ5RfRgjydWhuoxg71LUfyiQdijZuM=",
        version = "v0.53.0",
    )
    go_repository(
        name = "com_github_googlecloudplatform_opentelemetry_operations_go_internal_resourcemapping",
        importpath = "github.com/GoogleCloudPlatform/opentelemetry-operations-go/internal/resourcemapping",
        sum = "h1:Ron4zCA/yk6U7WOBXhTJcDpsUBG9npumK6xw2auFltQ=",
        version = "v0.53.0",
    )
    go_repository(
        name = "com_github_gorilla_mux",
        importpath = "github.com/gorilla/mux",
        sum = "h1:TuBL49tXwgrFYWhqrNgrUNEY92u81SPhu7sTdzQEiWY=",
        version = "v1.8.1",
    )
    go_repository(
        name = "com_github_iancoleman_strcase",
        importpath = "github.com/iancoleman/strcase",
        sum = "h1:nTXanmYxhfFAMjZL34Ov6gkzEsSJZ5DbhxWjvSASxEI=",
        version = "v0.3.0",
    )
    go_repository(
        name = "com_github_inconshreveable_mousetrap",
        importpath = "github.com/inconshreveable/mousetrap",
        sum = "h1:wN+x4NVGpMsO7ErUn/mUI3vEoE6Jt13X2s0bqwp9tc8=",
        version = "v1.1.0",
    )

    go_repository(
        name = "com_github_jhump_gopoet",
        importpath = "github.com/jhump/gopoet",
        sum = "h1:gYjOPnzHd2nzB37xYQZxj4EIQNpBrBskRqQQ3q4ZgSg=",
        version = "v0.1.0",
    )
    go_repository(
        name = "com_github_jhump_goprotoc",
        importpath = "github.com/jhump/goprotoc",
        sum = "h1:Y1UgUX+txUznfqcGdDef8ZOVlyQvnV0pKWZH08RmZuo=",
        version = "v0.5.0",
    )
    go_repository(
        name = "com_github_jhump_protoreflect",
        # Added in order to disable testproto BUILD file generation.
        # This should be retained by gazelle.
        build_file_proto_mode = "disable",
        importpath = "github.com/jhump/protoreflect",
        sum = "h1:54fZg+49widqXYQ0b+usAFHbMkBGR4PpXrsHc8+TBDg=",
        version = "v1.16.0",
    )
    go_repository(
        name = "com_github_kr_pretty",
        importpath = "github.com/kr/pretty",
        sum = "h1:flRD4NNwYAUpkphVc1HcthR4KEIFJ65n8Mw5qdRn3LE=",
        version = "v0.3.1",
    )
    go_repository(
        name = "com_github_kr_pty",
        importpath = "github.com/kr/pty",
        sum = "h1:VkoXIwSboBpnk99O/KFauAEILuNHv5DVFKZMBN/gUgw=",
        version = "v1.1.1",
    )
    go_repository(
        name = "com_github_kr_text",
        importpath = "github.com/kr/text",
        sum = "h1:5Nx0Ya0ZqY2ygV366QzturHI13Jq95ApcVaJBhpS+AY=",
        version = "v0.2.0",
    )
    go_repository(
        name = "com_github_pelletier_go_toml_v2",
        importpath = "github.com/pelletier/go-toml/v2",
        sum = "h1:mye9XuhQ6gvn5h28+VilKrrPoQVanw5PMw/TB0t5Ec4=",
        version = "v2.2.4",
    )

    go_repository(
        name = "com_github_pkg_diff",
        importpath = "github.com/pkg/diff",
        sum = "h1:aoZm08cpOy4WuID//EZDgcC4zIxODThtZNPirFr42+A=",
        version = "v0.0.0-20210226163009-20ebb0f2a09e",
    )
    go_repository(
        name = "com_github_planetscale_vtprotobuf",
        importpath = "github.com/planetscale/vtprotobuf",
        sum = "h1:GFCKgmp0tecUJ0sJuv4pzYCqS9+RGSn52M3FUwPs+uo=",
        version = "v0.6.1-0.20240319094008-0393e58bdf10",
    )
    go_repository(
        name = "com_github_pmezard_go_difflib",
        importpath = "github.com/pmezard/go-difflib",
        sum = "h1:4DBwDE0NGyQoBHbLQYPwSUPoCMWR5BEzIk/f1lZbAQM=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_rogpeppe_go_internal",
        importpath = "github.com/rogpeppe/go-internal",
        sum = "h1:UQB4HGPB6osV0SQTLymcB4TgvyWu6ZyliaW0tI/otEQ=",
        version = "v1.14.1",
    )
    go_repository(
        name = "com_github_russross_blackfriday_v2",
        importpath = "github.com/russross/blackfriday/v2",
        sum = "h1:JIOH55/0cWyOuilr9/qlrm0BSXldqnqwMsf35Ld67mk=",
        version = "v2.1.0",
    )
    go_repository(
        name = "com_github_sagikazarmark_locafero",
        importpath = "github.com/sagikazarmark/locafero",
        sum = "h1:/NQhBAkUb4+fH1jivKHWusDYFjMOOKU88eegjfxfHb4=",
        version = "v0.12.0",
    )
    go_repository(
        name = "com_github_soheilhy_cmux",
        importpath = "github.com/soheilhy/cmux",
        sum = "h1:jjzc5WVemNEDTLwv9tlmemhC73tI08BNOIGwBOo10Js=",
        version = "v0.1.5",
    )
    go_repository(
        name = "com_github_spf13_afero",
        importpath = "github.com/spf13/afero",
        sum = "h1:b/YBCLWAJdFWJTN9cLhiXXcD7mzKn9Dm86dNnfyQw1I=",
        version = "v1.15.0",
    )
    go_repository(
        name = "com_github_spf13_cast",
        importpath = "github.com/spf13/cast",
        sum = "h1:h2x0u2shc1QuLHfxi+cTJvs30+ZAHOGRic8uyGTDWxY=",
        version = "v1.10.0",
    )
    go_repository(
        name = "com_github_spf13_cobra",
        importpath = "github.com/spf13/cobra",
        sum = "h1:DMTTonx5m65Ic0GOoRY2c16WCbHxOOw6xxezuLaBpcU=",
        version = "v1.10.2",
    )
    go_repository(
        name = "com_github_spf13_pflag",
        importpath = "github.com/spf13/pflag",
        sum = "h1:4EBh2KAYBwaONj6b2Ye1GiHfwjqyROoF4RwYO+vPwFk=",
        version = "v1.0.10",
    )
    go_repository(
        name = "com_github_spf13_viper",
        importpath = "github.com/spf13/viper",
        sum = "h1:x5S+0EU27Lbphp4UKm1C+1oQO+rKx36vfCoaVebLFSU=",
        version = "v1.21.0",
    )
    go_repository(
        name = "com_github_spiffe_go_spiffe_v2",
        importpath = "github.com/spiffe/go-spiffe/v2",
        sum = "h1:l+DolpxNWYgruGQVV0xsfeya3CsC7m8iBzDnMpsbLuo=",
        version = "v2.6.0",
    )
    go_repository(
        name = "com_github_stretchr_objx",
        importpath = "github.com/stretchr/objx",
        sum = "h1:xuMeJ0Sdp5ZMRXx/aWO6RZxdr3beISkG5/G/aIRr3pY=",
        version = "v0.5.2",
    )
    go_repository(
        name = "com_github_stretchr_testify",
        importpath = "github.com/stretchr/testify",
        sum = "h1:7s2iGBzp5EwR7/aIZr8ao5+dra3wiQyKjjFuvgVKu7U=",
        version = "v1.11.1",
    )
    go_repository(
        name = "com_github_subosito_gotenv",
        importpath = "github.com/subosito/gotenv",
        sum = "h1:9NlTDc1FTs4qu0DDq7AEtTPNw6SVm7uBMsUCUjABIf8=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_github_zeebo_errs",
        importpath = "github.com/zeebo/errs",
        sum = "h1:XNdoD/RRMKP7HD0UhJnIzUy74ISdGGxURlYG8HSWSfM=",
        version = "v1.4.0",
    )

    go_repository(
        name = "com_gitlab_golang_commonmark_html",
        importpath = "gitlab.com/golang-commonmark/html",
        sum = "h1:K+bMSIx9A7mLES1rtG+qKduLIXq40DAzYHtb0XuCukA=",
        version = "v0.0.0-20191124015941-a22733972181",
    )
    go_repository(
        name = "com_gitlab_golang_commonmark_linkify",
        build_directives = [
            "gazelle:resolve go go golang.org/x/text/unicode/rangetable @org_golang_x_text//unicode/rangetable:go_default_library",
        ],
        importpath = "gitlab.com/golang-commonmark/linkify",
        sum = "h1:1Coh5BsUBlXoEJmIEaNzVAWrtg9k7/eJzailMQr1grw=",
        version = "v0.0.0-20200225224916-64bca66f6ad3",
    )
    go_repository(
        name = "com_gitlab_golang_commonmark_markdown",
        importpath = "gitlab.com/golang-commonmark/markdown",
        sum = "h1:O85GKETcmnCNAfv4Aym9tepU8OE0NmcZNqPlXcsBKBs=",
        version = "v0.0.0-20211110145824-bf3e522c626a",
    )
    go_repository(
        name = "com_gitlab_golang_commonmark_mdurl",
        importpath = "gitlab.com/golang-commonmark/mdurl",
        sum = "h1:qqjvoVXdWIcZCLPMlzgA7P9FZWdPGPvP/l3ef8GzV6o=",
        version = "v0.0.0-20191124015652-932350d1cb84",
    )
    go_repository(
        name = "com_gitlab_golang_commonmark_puny",
        importpath = "gitlab.com/golang-commonmark/puny",
        sum = "h1:Wku8eEdeJqIOFHtrfkYUByc4bCaTeA6fL0UJgfEiFMI=",
        version = "v0.0.0-20191124015043-9f83538fa04f",
    )
    go_repository(
        name = "com_gitlab_opennota_wd",
        importpath = "gitlab.com/opennota/wd",
        sum = "h1:uPZaMiz6Sz0PZs3IZJWpU5qHKGNy///1pacZC9txiUI=",
        version = "v0.0.0-20180912061657-c5d65f63c638",
    )
    go_repository(
        name = "com_google_cloud_go",
        # This is part of a fix for https://github.com/googleapis/gapic-generator-go/issues/387.
        build_extra_args = ["-exclude=longrunning/autogen/info.go"],
        importpath = "cloud.google.com/go",
        sum = "h1:2NAUJwPR47q+E35uaJeYoNhuNEM9kM8SjgRgdeOJUSE=",
        version = "v0.123.0",
    )
    go_repository(
        name = "com_google_cloud_go_accessapproval",
        importpath = "cloud.google.com/go/accessapproval",
        sum = "h1:OP0GXiVIXUxGf0ZxlfYHTqOvRdvxiQcn/l0ZWQ+6LRA=",
        version = "v1.11.0",
    )
    go_repository(
        name = "com_google_cloud_go_accesscontextmanager",
        importpath = "cloud.google.com/go/accesscontextmanager",
        sum = "h1:IEQzenTC6z/AvkRLxWw9vb51yty3VnNT+tJX6MJGWUI=",
        version = "v1.12.0",
    )
    go_repository(
        name = "com_google_cloud_go_aiplatform",
        importpath = "cloud.google.com/go/aiplatform",
        sum = "h1:77uy+1+G11yP98axztHeVaW85NcXVLA7WP6ynUXmOlE=",
        version = "v1.124.0",
    )
    go_repository(
        name = "com_google_cloud_go_analytics",
        importpath = "cloud.google.com/go/analytics",
        sum = "h1:f5bm5lIXyokBFBvRPEkGDEiuUmdci9NzQlQSFNabLrg=",
        version = "v0.33.0",
    )
    go_repository(
        name = "com_google_cloud_go_apigateway",
        importpath = "cloud.google.com/go/apigateway",
        sum = "h1:HKUf9+ltcoL7uQVTrb3izB1HY4YKpOGiDoIwLGzaCmY=",
        version = "v1.10.0",
    )
    go_repository(
        name = "com_google_cloud_go_apigeeconnect",
        importpath = "cloud.google.com/go/apigeeconnect",
        sum = "h1:y9ozXuug+lsdHxFvyrTBmhYKiX8f3JIijq5Gjsr0aww=",
        version = "v1.10.0",
    )
    go_repository(
        name = "com_google_cloud_go_apigeeregistry",
        importpath = "cloud.google.com/go/apigeeregistry",
        sum = "h1:RaK9SAcqYxoqKur2YASVr8eYmylBWUEwHGt2u++Gll0=",
        version = "v0.13.0",
    )
    go_repository(
        name = "com_google_cloud_go_appengine",
        importpath = "cloud.google.com/go/appengine",
        sum = "h1:GlchNBW7pn8WCVqH44DTd2gk62p84mOTVoiAEeVqZuo=",
        version = "v1.12.0",
    )
    go_repository(
        name = "com_google_cloud_go_area120",
        importpath = "cloud.google.com/go/area120",
        sum = "h1:vlK71l9DU1Pm8txF4M3LjGu8my+8jwuprQkf0y/yxcc=",
        version = "v0.13.0",
    )
    go_repository(
        name = "com_google_cloud_go_artifactregistry",
        importpath = "cloud.google.com/go/artifactregistry",
        sum = "h1:Isjhp2tw/7sQf3d0I7modlkIPHTC9TVZbHcWHd1JizY=",
        version = "v1.23.0",
    )
    go_repository(
        name = "com_google_cloud_go_asset",
        importpath = "cloud.google.com/go/asset",
        sum = "h1:aonULd4vfBAMVYUMO6fWAAq+gkNiWpdo248iTgGsESs=",
        version = "v1.25.0",
    )
    go_repository(
        name = "com_google_cloud_go_assuredworkloads",
        importpath = "cloud.google.com/go/assuredworkloads",
        sum = "h1:A07Q/QdwPT31qNY1CMxTsXzZGB3lMx3mJM2Zg2FR0TM=",
        version = "v1.16.0",
    )
    go_repository(
        name = "com_google_cloud_go_auth",
        importpath = "cloud.google.com/go/auth",
        sum = "h1:kXTssoVb4azsVDoUiF8KvxAqrsQcQtB53DcSgta74CA=",
        version = "v0.20.0",
    )
    go_repository(
        name = "com_google_cloud_go_auth_oauth2adapt",
        importpath = "cloud.google.com/go/auth/oauth2adapt",
        sum = "h1:keo8NaayQZ6wimpNSmW5OPc283g65QNIiLpZnkHRbnc=",
        version = "v0.2.8",
    )
    go_repository(
        name = "com_google_cloud_go_automl",
        importpath = "cloud.google.com/go/automl",
        sum = "h1:W9YNW1DN3wqocRgxxnXVg7ZttT0avprqPNNz/fN/S3k=",
        version = "v1.18.0",
    )
    go_repository(
        name = "com_google_cloud_go_baremetalsolution",
        importpath = "cloud.google.com/go/baremetalsolution",
        sum = "h1:qui/9Ylke6cR+WYROIWMQpK+UsvMPZ4S/MIxP9wFDms=",
        version = "v1.7.0",
    )
    go_repository(
        name = "com_google_cloud_go_batch",
        importpath = "cloud.google.com/go/batch",
        sum = "h1:k4SgGgITcibw2wZMXOEah502NaA3I/6i3rx4tzXOBW0=",
        version = "v1.17.0",
    )
    go_repository(
        name = "com_google_cloud_go_beyondcorp",
        importpath = "cloud.google.com/go/beyondcorp",
        sum = "h1:30VnvhOQmfNZ1mpDO21yQBkTDpeVIO3LL2GRinWsVJ8=",
        version = "v1.5.0",
    )
    go_repository(
        name = "com_google_cloud_go_bigquery",
        importpath = "cloud.google.com/go/bigquery",
        sum = "h1:wnfVSXN6GEMlsAoHWdhzTC8NMsptOx2hsqPiI+lTs3I=",
        version = "v1.76.0",
    )
    go_repository(
        name = "com_google_cloud_go_bigtable",
        importpath = "cloud.google.com/go/bigtable",
        sum = "h1:Bd6vITx01s6gsdEPvjIGJL/oMMdKvraGI34UiixR2gk=",
        version = "v1.46.0",
    )
    go_repository(
        name = "com_google_cloud_go_billing",
        importpath = "cloud.google.com/go/billing",
        sum = "h1:N+uNaIUT3XUy/z1x9rc3QX6gS/OkGm7KD+vaIRGEpPo=",
        version = "v1.24.0",
    )
    go_repository(
        name = "com_google_cloud_go_binaryauthorization",
        importpath = "cloud.google.com/go/binaryauthorization",
        sum = "h1:Vi4EplNiT2PFaKrHX6OFucu/bUgrzWZt6OQpG/eJgoI=",
        version = "v1.13.0",
    )
    go_repository(
        name = "com_google_cloud_go_certificatemanager",
        importpath = "cloud.google.com/go/certificatemanager",
        sum = "h1:cGtIA5WPVpZDqC35E+i5FRZDziUPbIIKE4wo6mNzqCI=",
        version = "v1.12.0",
    )
    go_repository(
        name = "com_google_cloud_go_channel",
        importpath = "cloud.google.com/go/channel",
        sum = "h1:sS8ocaYXX3al9yGFhYfPzz0UzOpAt1DdbFiuoqr4rdw=",
        version = "v1.24.0",
    )
    go_repository(
        name = "com_google_cloud_go_cloudbuild",
        importpath = "cloud.google.com/go/cloudbuild",
        sum = "h1:rOfVbca71Be0MjV+PCg2Nacqdjlh8UZw4bdg5PvZZtI=",
        version = "v1.28.0",
    )
    go_repository(
        name = "com_google_cloud_go_clouddms",
        importpath = "cloud.google.com/go/clouddms",
        sum = "h1:OFNGNDRsrAr7pXITuzUbk9h+/oNJi9vQY5mAESaM54U=",
        version = "v1.11.0",
    )
    go_repository(
        name = "com_google_cloud_go_cloudtasks",
        importpath = "cloud.google.com/go/cloudtasks",
        sum = "h1:VxCvIC4zo0qO2xXsvTeLJaexPJ20s6FSVeP6n4LWpFw=",
        version = "v1.16.0",
    )
    go_repository(
        name = "com_google_cloud_go_compute",
        importpath = "cloud.google.com/go/compute",
        sum = "h1:CqGt23ysz990ZZe1vq/9aDPKKnmwM6kcC7Y1Q05H2kI=",
        version = "v1.60.0",
    )
    go_repository(
        name = "com_google_cloud_go_compute_metadata",
        importpath = "cloud.google.com/go/compute/metadata",
        sum = "h1:pDUj4QMoPejqq20dK0Pg2N4yG9zIkYGdBtwLoEkH9Zs=",
        version = "v0.9.0",
    )
    go_repository(
        name = "com_google_cloud_go_contactcenterinsights",
        importpath = "cloud.google.com/go/contactcenterinsights",
        sum = "h1:tUPpMlGktJuvVX/xp/7xQZpRg59tL+TgeIuXUBnpvyg=",
        version = "v1.20.0",
    )
    go_repository(
        name = "com_google_cloud_go_container",
        importpath = "cloud.google.com/go/container",
        sum = "h1:K4nmtmJezHOzsIyedAOv1Ok36krw1apFmo4zXBaRL1A=",
        version = "v1.49.0",
    )
    go_repository(
        name = "com_google_cloud_go_containeranalysis",
        importpath = "cloud.google.com/go/containeranalysis",
        sum = "h1:8+Buz/+pVqsLiapngXpq+NTR9y/eN8n4E6ubIfq7weA=",
        version = "v0.17.0",
    )
    go_repository(
        name = "com_google_cloud_go_datacatalog",
        importpath = "cloud.google.com/go/datacatalog",
        sum = "h1:K5zBLN4YLAPyzFt/IzAE/OhKmQUHSfzMflCB7IjAg/E=",
        version = "v1.29.0",
    )
    go_repository(
        name = "com_google_cloud_go_dataflow",
        importpath = "cloud.google.com/go/dataflow",
        sum = "h1:on1PwQQjEcZcCvb/lYJp8Zj0s4oCu9DcLC3xqIKIUnM=",
        version = "v0.14.0",
    )
    go_repository(
        name = "com_google_cloud_go_dataform",
        importpath = "cloud.google.com/go/dataform",
        sum = "h1:sJV9oGp3SutWmjWcSDOipPaRdqwZQePoHfYK1mcaXrE=",
        version = "v0.17.0",
    )
    go_repository(
        name = "com_google_cloud_go_datafusion",
        importpath = "cloud.google.com/go/datafusion",
        sum = "h1:LoLIvA54c3yuV2vQ6qiE815hMOGcuyE5ySHSHlAfF+I=",
        version = "v1.11.0",
    )
    go_repository(
        name = "com_google_cloud_go_datalabeling",
        importpath = "cloud.google.com/go/datalabeling",
        sum = "h1:KzMS/gRcZcyGwi7Mj5O5eNIbhPxwWITYZqRnQpsj6n8=",
        version = "v0.12.0",
    )
    go_repository(
        name = "com_google_cloud_go_dataplex",
        importpath = "cloud.google.com/go/dataplex",
        sum = "h1:FXcPlYhlp5jdnaIygqh+6n7HJyeSnunbIRV8Y1yfim8=",
        version = "v1.32.0",
    )
    go_repository(
        name = "com_google_cloud_go_dataproc_v2",
        importpath = "cloud.google.com/go/dataproc/v2",
        sum = "h1:nigeuU3AoKMOuPoQ/F3QDCGRDdJCkFMR6mAhXC4uDfA=",
        version = "v2.19.0",
    )
    go_repository(
        name = "com_google_cloud_go_dataqna",
        importpath = "cloud.google.com/go/dataqna",
        sum = "h1:8ylMV3/iLOCQgBgH472ILLd8O4g+z8SBX7rvWxmIUtg=",
        version = "v0.12.0",
    )
    go_repository(
        name = "com_google_cloud_go_datastore",
        importpath = "cloud.google.com/go/datastore",
        sum = "h1:FOyx2Ag6ibD2wFkz9S8EiNrmBugia8pQOfpyJxi2yqA=",
        version = "v1.22.0",
    )
    go_repository(
        name = "com_google_cloud_go_datastream",
        importpath = "cloud.google.com/go/datastream",
        sum = "h1:oAlCDBb71E8aWEaaHctTaJqUWtTJ4d7fsWjo2gwMz9o=",
        version = "v1.18.0",
    )
    go_repository(
        name = "com_google_cloud_go_deploy",
        importpath = "cloud.google.com/go/deploy",
        sum = "h1:gprUC11589FPGXFD8dB2qOneceAzOV6ZqtvWceIpyOM=",
        version = "v1.30.0",
    )
    go_repository(
        name = "com_google_cloud_go_dialogflow",
        importpath = "cloud.google.com/go/dialogflow",
        sum = "h1:6w2iWsmLSC/fB9CEy4yXDqFMomSiwVDcSqKWwpx6ULY=",
        version = "v1.80.0",
    )
    go_repository(
        name = "com_google_cloud_go_dlp",
        importpath = "cloud.google.com/go/dlp",
        sum = "h1:idIs7hz0IVBEzraaV1TkKydpg7Ou+E+mRs+FOl6b5pQ=",
        version = "v1.32.0",
    )
    go_repository(
        name = "com_google_cloud_go_documentai",
        importpath = "cloud.google.com/go/documentai",
        sum = "h1:5h/lbY90wJGWjJXliSRFJCi5h6NU1DGPjL7R2w4clrg=",
        version = "v1.46.0",
    )
    go_repository(
        name = "com_google_cloud_go_domains",
        importpath = "cloud.google.com/go/domains",
        sum = "h1:aZhTyHEXBazcK0ozBIlpocfyI4ByokQEQtKf/KAzwNQ=",
        version = "v0.13.0",
    )
    go_repository(
        name = "com_google_cloud_go_edgecontainer",
        importpath = "cloud.google.com/go/edgecontainer",
        sum = "h1:BgOHohjlwyX1W7/WgoAO7+tmyoeLo4JTY8DawcxIHt0=",
        version = "v1.7.0",
    )
    go_repository(
        name = "com_google_cloud_go_errorreporting",
        importpath = "cloud.google.com/go/errorreporting",
        sum = "h1:X6sVXp7qhJTvfUqaykmGBfy1RrshgFDfG7SXeqXysD8=",
        version = "v0.7.0",
    )
    go_repository(
        name = "com_google_cloud_go_essentialcontacts",
        importpath = "cloud.google.com/go/essentialcontacts",
        sum = "h1:5TxmLnvtKCISuh4tMbVsLAVLkSIckWNP8vLJ67gXsZM=",
        version = "v1.10.0",
    )
    go_repository(
        name = "com_google_cloud_go_eventarc",
        importpath = "cloud.google.com/go/eventarc",
        sum = "h1:BzkhBK2K/FxEnaNNJuYyHq/cW5AdFWMFNNHWqeVqgxk=",
        version = "v1.21.0",
    )
    go_repository(
        name = "com_google_cloud_go_filestore",
        importpath = "cloud.google.com/go/filestore",
        sum = "h1:7L8pvhPr6ZaTKpTinDmAuSwkYEt+B+eRl0OA/Pm882w=",
        version = "v1.13.0",
    )
    go_repository(
        name = "com_google_cloud_go_firestore",
        importpath = "cloud.google.com/go/firestore",
        sum = "h1:avooeboIq37vKXobrbPUFhFBxS/c3FqmWoX0xs8dO6E=",
        version = "v1.22.0",
    )
    go_repository(
        name = "com_google_cloud_go_functions",
        importpath = "cloud.google.com/go/functions",
        sum = "h1:rJ2bSt2KUEi0OBMsUKICI/lJYCsTOw3aMgzKxBmuyNo=",
        version = "v1.22.0",
    )
    go_repository(
        name = "com_google_cloud_go_gkebackup",
        importpath = "cloud.google.com/go/gkebackup",
        sum = "h1:SbClooIm/uYF2x+01AjlQ1nFNNZUiG1cGeOknpOUTtg=",
        version = "v1.11.0",
    )
    go_repository(
        name = "com_google_cloud_go_gkeconnect",
        importpath = "cloud.google.com/go/gkeconnect",
        sum = "h1:aTOSotFKQiOGhVsiugNXL25qLnpwg9GfK6lUzt7923s=",
        version = "v0.15.0",
    )
    go_repository(
        name = "com_google_cloud_go_gkehub",
        importpath = "cloud.google.com/go/gkehub",
        sum = "h1:bdVbZyqR8rtIyxRITTqJLir/4F9R65kQQGQEc6T/ZcU=",
        version = "v0.19.0",
    )
    go_repository(
        name = "com_google_cloud_go_gkemulticloud",
        importpath = "cloud.google.com/go/gkemulticloud",
        sum = "h1:mUmMB6aLW/RJAZ9iWkd30mdtaoAnuYa2vgyKwrPBc9o=",
        version = "v1.9.0",
    )
    go_repository(
        name = "com_google_cloud_go_gsuiteaddons",
        importpath = "cloud.google.com/go/gsuiteaddons",
        sum = "h1:xNf2+E01Pmv/xyV1l5frvAW+n9S//j/zpZ5X4UxzB/I=",
        version = "v1.10.0",
    )
    go_repository(
        name = "com_google_cloud_go_iam",
        importpath = "cloud.google.com/go/iam",
        sum = "h1:89wyjxT6DL4b5rk/Nk8eBC9DHqf+JiMstrn5IEYxFw4=",
        version = "v1.9.0",
    )
    go_repository(
        name = "com_google_cloud_go_iap",
        importpath = "cloud.google.com/go/iap",
        sum = "h1:Qc70Lj88C4A9GwKYdb0cRLR0Gg1SsUf0k8uNDeFmRMk=",
        version = "v1.15.0",
    )
    go_repository(
        name = "com_google_cloud_go_ids",
        importpath = "cloud.google.com/go/ids",
        sum = "h1:kwYIE/qZj5bFyjdpbKnVNilV4OeLxO03nzxD746ATJw=",
        version = "v1.8.0",
    )
    go_repository(
        name = "com_google_cloud_go_iot",
        importpath = "cloud.google.com/go/iot",
        sum = "h1:qrGFVW2XINbVaGjokqaGMRFSoLaM2BDdPEEOydUN22o=",
        version = "v1.11.0",
    )
    go_repository(
        name = "com_google_cloud_go_kms",
        importpath = "cloud.google.com/go/kms",
        sum = "h1:bAW1C5FQf+6GhPkywQzPlsULALCG7c16qpXLFGV9ivY=",
        version = "v1.29.0",
    )
    go_repository(
        name = "com_google_cloud_go_language",
        importpath = "cloud.google.com/go/language",
        sum = "h1:ichJvPS/41ezA3d3wGuCVQ/N7c7533xfqw8O8c5HGKs=",
        version = "v1.17.0",
    )
    go_repository(
        name = "com_google_cloud_go_lifesciences",
        importpath = "cloud.google.com/go/lifesciences",
        sum = "h1:v6JK7TSnJ0sfzp1NpPicniwOKZlZ+FXCpglGUXBk4PQ=",
        version = "v0.13.0",
    )
    go_repository(
        name = "com_google_cloud_go_logging",
        importpath = "cloud.google.com/go/logging",
        sum = "h1:MMNgYRvZ/pEwiNSkcoJTKWfAbAJDqCqAMJiarZx+/CI=",
        version = "v1.16.0",
    )
    go_repository(
        name = "com_google_cloud_go_longrunning",
        importpath = "cloud.google.com/go/longrunning",
        sum = "h1:fE4XVLJQj+gRnw1HrbDyQXXgC0aiqY3wxP7DDU4cWk0=",
        version = "v0.11.0",
    )
    go_repository(
        name = "com_google_cloud_go_managedidentities",
        importpath = "cloud.google.com/go/managedidentities",
        sum = "h1:MLLqEBZqvKLpEIrVQ53Nzy2vIOqCIU7PumbQvkHxQ4k=",
        version = "v1.10.0",
    )
    go_repository(
        name = "com_google_cloud_go_maps",
        importpath = "cloud.google.com/go/maps",
        sum = "h1:BlptNOYwjeoKGB7J2DzBX/sXB4m/b40PHczy8SAJkDI=",
        version = "v1.33.0",
    )
    go_repository(
        name = "com_google_cloud_go_mediatranslation",
        importpath = "cloud.google.com/go/mediatranslation",
        sum = "h1:5x9B8UkK7J9gfDP25pdQKDTzPmffa81FGyCi3whTMPU=",
        version = "v0.12.0",
    )
    go_repository(
        name = "com_google_cloud_go_memcache",
        importpath = "cloud.google.com/go/memcache",
        sum = "h1:g2KUZ9DLNmQ4Er1XJPaPts2zj+sFnjkGg7Oj+zdtq5E=",
        version = "v1.14.0",
    )
    go_repository(
        name = "com_google_cloud_go_metastore",
        importpath = "cloud.google.com/go/metastore",
        sum = "h1:LfWZmTiGk/nAp6WBMgOxLyy6WEeipQxWB1QEY8BcIRc=",
        version = "v1.17.0",
    )
    go_repository(
        name = "com_google_cloud_go_monitoring",
        importpath = "cloud.google.com/go/monitoring",
        sum = "h1:BhYwMqao+e5Nn7JtWMM9m6zRtKtVUK6kJWMizXChkLU=",
        version = "v1.27.0",
    )
    go_repository(
        name = "com_google_cloud_go_networkconnectivity",
        importpath = "cloud.google.com/go/networkconnectivity",
        sum = "h1:jJ1YbJyJUAggqBYvap8VxD7MGD+fDjJLQYbbjU9Dn54=",
        version = "v1.24.0",
    )
    go_repository(
        name = "com_google_cloud_go_networkmanagement",
        importpath = "cloud.google.com/go/networkmanagement",
        sum = "h1:x9LDaO0I+YhkF6j2QJC7VoQZRRpwGb92aPpdjY+W8Vo=",
        version = "v1.26.0",
    )
    go_repository(
        name = "com_google_cloud_go_networksecurity",
        importpath = "cloud.google.com/go/networksecurity",
        sum = "h1:v1NKTztDfEDIr7C5wHmdELebok4pmnaIC6nGAOv3AuY=",
        version = "v0.14.0",
    )
    go_repository(
        name = "com_google_cloud_go_notebooks",
        importpath = "cloud.google.com/go/notebooks",
        sum = "h1:HJ589RNPGVbvuttHGiD6fsVZTiK1FxU3fRfBwn0baQA=",
        version = "v1.15.0",
    )
    go_repository(
        name = "com_google_cloud_go_optimization",
        importpath = "cloud.google.com/go/optimization",
        sum = "h1:ZviinmU5P9DCxLBUXHz64PTLzetL8dNXs8aJ+Ii37O4=",
        version = "v1.10.0",
    )
    go_repository(
        name = "com_google_cloud_go_orchestration",
        importpath = "cloud.google.com/go/orchestration",
        sum = "h1:XJcy4tHp4GBDk6CbcNnHtp3KR4Lt89b+OmgS2Jz2lwI=",
        version = "v1.14.0",
    )
    go_repository(
        name = "com_google_cloud_go_orgpolicy",
        importpath = "cloud.google.com/go/orgpolicy",
        sum = "h1:/DfoElp43vYbWrNSZqaz/NPeNty4vXfvpIzT1ZqOtsY=",
        version = "v1.18.0",
    )
    go_repository(
        name = "com_google_cloud_go_osconfig",
        importpath = "cloud.google.com/go/osconfig",
        sum = "h1:sLX13w1WsKxnknfMIW+RS5lsKplNBxsgHYyTr5hPSsQ=",
        version = "v1.19.0",
    )
    go_repository(
        name = "com_google_cloud_go_oslogin",
        importpath = "cloud.google.com/go/oslogin",
        sum = "h1:B7HE//+eMtlwuTKcFkxk7+U5tfBlQ7nW9rpYSk+8CBY=",
        version = "v1.17.0",
    )
    go_repository(
        name = "com_google_cloud_go_phishingprotection",
        importpath = "cloud.google.com/go/phishingprotection",
        sum = "h1:TLdbZ1g1m/qjUHPoSU9yvTi7COASdiJ1zX3QiE5sfTA=",
        version = "v0.12.0",
    )
    go_repository(
        name = "com_google_cloud_go_policytroubleshooter",
        importpath = "cloud.google.com/go/policytroubleshooter",
        sum = "h1:BJ4RoLqUvzEeHakP3YsXcYQxvC1/RDERwZzghpdwXBA=",
        version = "v1.14.0",
    )
    go_repository(
        name = "com_google_cloud_go_privatecatalog",
        importpath = "cloud.google.com/go/privatecatalog",
        sum = "h1:7K4oV81cfMJBBrQM/tuQ4jZX4hyU28mTgRa+quvCj9g=",
        version = "v0.13.0",
    )
    go_repository(
        name = "com_google_cloud_go_pubsub",
        importpath = "cloud.google.com/go/pubsub",
        sum = "h1:54Up97HnThdP4H8jjWJSSQ/mnYG2EKon7ZSNETRq0tM=",
        version = "v1.50.2",
    )
    go_repository(
        name = "com_google_cloud_go_pubsub_v2",
        importpath = "cloud.google.com/go/pubsub/v2",
        sum = "h1:+TwXJr78P9RrMV3S8lKHIhJo2E99jI7ta65e+ujJjts=",
        version = "v2.5.1",
    )
    go_repository(
        name = "com_google_cloud_go_pubsublite",
        importpath = "cloud.google.com/go/pubsublite",
        sum = "h1:jLQozsEVr+c6tOU13vDugtnaBSUy/PD5zK6mhm+uF1Y=",
        version = "v1.8.2",
    )
    go_repository(
        name = "com_google_cloud_go_recaptchaenterprise_v2",
        importpath = "cloud.google.com/go/recaptchaenterprise/v2",
        sum = "h1:/cwXOOITPM6M5jSCIDBJ2SO9+KkIOG44K0g8d9uvBxw=",
        version = "v2.24.0",
    )
    go_repository(
        name = "com_google_cloud_go_recommendationengine",
        importpath = "cloud.google.com/go/recommendationengine",
        sum = "h1:WfodE5PFHnqsdhyv/fw2GcZIKec3kBmPvpcIbF/9PHM=",
        version = "v0.12.0",
    )
    go_repository(
        name = "com_google_cloud_go_recommender",
        importpath = "cloud.google.com/go/recommender",
        sum = "h1:5VCy6TBEAQ0zqHh5ROgzw+mOEB3IZs9b2qtNgpqhVIQ=",
        version = "v1.16.0",
    )
    go_repository(
        name = "com_google_cloud_go_redis",
        importpath = "cloud.google.com/go/redis",
        sum = "h1:2sz4rMZ/1+UwDwMeS61LqQQHMddyzmE0FRO2rbU7MWk=",
        version = "v1.21.0",
    )
    go_repository(
        name = "com_google_cloud_go_resourcemanager",
        importpath = "cloud.google.com/go/resourcemanager",
        sum = "h1:cc291PxLoKrHKVxqoJ2uMMzrxVJj+sRe+iEb1DFlDNA=",
        version = "v1.13.0",
    )
    go_repository(
        name = "com_google_cloud_go_resourcesettings",
        importpath = "cloud.google.com/go/resourcesettings",
        sum = "h1:13HOFU7v4cEvIHXSAQbinF4wp2Baybbq7q9FMctg1Ek=",
        version = "v1.8.3",
    )
    go_repository(
        name = "com_google_cloud_go_retail",
        importpath = "cloud.google.com/go/retail",
        sum = "h1:tix8f2anQixETjsgWy2KMEyVkI+3MJv29+rSsVLF3jc=",
        version = "v1.29.0",
    )
    go_repository(
        name = "com_google_cloud_go_run",
        importpath = "cloud.google.com/go/run",
        sum = "h1:kjXZKDwrUOeUYDd7/0TZ/iKsG3rJ3Lq3cyksTspcNSU=",
        version = "v1.19.0",
    )
    go_repository(
        name = "com_google_cloud_go_scheduler",
        importpath = "cloud.google.com/go/scheduler",
        sum = "h1:SZVfUC+OaGRl3DdBjhFia4vB2S9h9SBLCxxiCxuCn5k=",
        version = "v1.14.0",
    )
    go_repository(
        name = "com_google_cloud_go_secretmanager",
        importpath = "cloud.google.com/go/secretmanager",
        sum = "h1:dm9BK06xl+hrxp2unT2psjZeypPj5c6uPiABb6fmicE=",
        version = "v1.19.0",
    )
    go_repository(
        name = "com_google_cloud_go_security",
        importpath = "cloud.google.com/go/security",
        sum = "h1:9DIk8YIb2RAW29aB4u45Ir9vP0enQI9+HXwGHtqkIdg=",
        version = "v1.22.0",
    )
    go_repository(
        name = "com_google_cloud_go_securitycenter",
        importpath = "cloud.google.com/go/securitycenter",
        sum = "h1:sH9dLXwcnQE9qcElYsBwm0wbKKHa0ooyDkJgIppfQdk=",
        version = "v1.42.0",
    )
    go_repository(
        name = "com_google_cloud_go_servicedirectory",
        importpath = "cloud.google.com/go/servicedirectory",
        sum = "h1:jF7M3bxnzLbyQDyVXMfyjTEhHJ9IptDlTh/WVlyIjqA=",
        version = "v1.15.0",
    )
    go_repository(
        name = "com_google_cloud_go_shell",
        importpath = "cloud.google.com/go/shell",
        sum = "h1:N2S8p86S/qcHkOqzl93rkxvt3mNutgCwhJFUg8yamp8=",
        version = "v1.11.0",
    )
    go_repository(
        name = "com_google_cloud_go_spanner",
        importpath = "cloud.google.com/go/spanner",
        sum = "h1:tVUzYI9IZJY1SDvCmCzyuLoaAyotZgiyQq6go+lNH5A=",
        version = "v1.90.0",
    )
    go_repository(
        name = "com_google_cloud_go_speech",
        importpath = "cloud.google.com/go/speech",
        sum = "h1:555yroj4HCS7SPgfHuDU8zX+E5KrhccVWG96HNyBUAk=",
        version = "v1.33.0",
    )
    go_repository(
        name = "com_google_cloud_go_storage",
        importpath = "cloud.google.com/go/storage",
        sum = "h1:iixmq2Fse2tqxMbWhLWC9HfBj1qdxqAmiK8/eqtsLxI=",
        version = "v1.56.0",
    )
    go_repository(
        name = "com_google_cloud_go_storagetransfer",
        importpath = "cloud.google.com/go/storagetransfer",
        sum = "h1:9XRKD/zyjrbNcZdIk8Jf5D1a9MRIPCbYoBSCXiKMp3E=",
        version = "v1.16.0",
    )
    go_repository(
        name = "com_google_cloud_go_talent",
        importpath = "cloud.google.com/go/talent",
        sum = "h1:zqIv5lUYieNriMtkrzPMXnIqK2keXZynLlJ7Q/S7wGM=",
        version = "v1.11.0",
    )
    go_repository(
        name = "com_google_cloud_go_texttospeech",
        importpath = "cloud.google.com/go/texttospeech",
        sum = "h1:4AfzGgXxLCaXgoH1rlRzvsnhHXsI3Ef4gBB5EZTreKQ=",
        version = "v1.19.0",
    )
    go_repository(
        name = "com_google_cloud_go_tpu",
        importpath = "cloud.google.com/go/tpu",
        sum = "h1:XtAov+UByMUzJrID/CrmRU6IQY3fBggSh8Ps/WG4Ab4=",
        version = "v1.11.0",
    )
    go_repository(
        name = "com_google_cloud_go_trace",
        importpath = "cloud.google.com/go/trace",
        sum = "h1:jUtnmOrNcu5XJNk4Gz0fv+v5sM0weaOa3z5MPQUjRXs=",
        version = "v1.14.0",
    )
    go_repository(
        name = "com_google_cloud_go_translate",
        importpath = "cloud.google.com/go/translate",
        sum = "h1:3TYEwGIqE/9CFyX2HUaTmkCFrASRZSZu84c8aEX9tz4=",
        version = "v1.15.0",
    )
    go_repository(
        name = "com_google_cloud_go_video",
        importpath = "cloud.google.com/go/video",
        sum = "h1:DF/BtPyMGYAV1rl5v24IYcbDq+chXfVq7LqxLsJhv8w=",
        version = "v1.30.0",
    )
    go_repository(
        name = "com_google_cloud_go_videointelligence",
        importpath = "cloud.google.com/go/videointelligence",
        sum = "h1:i+vn1AHvqNxuFijMoW9i/IuGE9iSX2UfZKxu/abQfU8=",
        version = "v1.15.0",
    )
    go_repository(
        name = "com_google_cloud_go_vision_v2",
        importpath = "cloud.google.com/go/vision/v2",
        sum = "h1:GHCFhIlvpEX0IEwKCrIqS3KlOpxzINa3gucvdjoWBCE=",
        version = "v2.12.0",
    )
    go_repository(
        name = "com_google_cloud_go_vmmigration",
        importpath = "cloud.google.com/go/vmmigration",
        sum = "h1:wVzHT4l6oO2GWlL6ZH1izrFKy3GszIrmFIXcd0sK/Bg=",
        version = "v1.13.0",
    )
    go_repository(
        name = "com_google_cloud_go_vmwareengine",
        importpath = "cloud.google.com/go/vmwareengine",
        sum = "h1:KGHC9gXboFJ34AVLNmseij/Wdzr/JyKdvIlBpKSTovE=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_google_cloud_go_vpcaccess",
        importpath = "cloud.google.com/go/vpcaccess",
        sum = "h1:1blFhaghROzfNERac0vu4PQJgwqqB/bT5TpQxiVyc1c=",
        version = "v1.11.0",
    )
    go_repository(
        name = "com_google_cloud_go_webrisk",
        importpath = "cloud.google.com/go/webrisk",
        sum = "h1:Lo6y/UdNSZk8FGWa/zuhx4bwJn1/Axdi7QsDj8CGqwQ=",
        version = "v1.14.0",
    )
    go_repository(
        name = "com_google_cloud_go_websecurityscanner",
        importpath = "cloud.google.com/go/websecurityscanner",
        sum = "h1:2QtVdTKJeCivg1mtPy4LBL9Knhb7HTvw5q7aRKVenHM=",
        version = "v1.10.0",
    )
    go_repository(
        name = "com_google_cloud_go_workflows",
        importpath = "cloud.google.com/go/workflows",
        sum = "h1:pZdcmGRNe39x/yP5BOn5nYPxzcM4y6z9djOjsYb69w0=",
        version = "v1.17.0",
    )
    go_repository(
        name = "dev_cel_expr",
        importpath = "cel.dev/expr",
        sum = "h1:1KrZg61W6TWSxuNZ37Xy49ps13NUovb66QLprthtwi4=",
        version = "v0.25.1",
    )
    go_repository(
        name = "in_gopkg_check_v1",
        importpath = "gopkg.in/check.v1",
        sum = "h1:Hei/4ADfdWqJk1ZMxUNpqntNwaWcugrBjAiHlqqRiVk=",
        version = "v1.0.0-20201130134442-10cb98267c6c",
    )
    go_repository(
        name = "in_gopkg_yaml_v2",
        importpath = "gopkg.in/yaml.v2",
        sum = "h1:D8xgwECY7CYvx+Y2n4sBz93Jn9JRvxdiyyo8CTfuKaY=",
        version = "v2.4.0",
    )
    go_repository(
        name = "in_gopkg_yaml_v3",
        importpath = "gopkg.in/yaml.v3",
        sum = "h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=",
        version = "v3.0.1",
    )
    go_repository(
        name = "in_yaml_go_yaml_v3",
        importpath = "go.yaml.in/yaml/v3",
        sum = "h1:tfq32ie2Jv2UxXFdLJdh3jXuOzWiL1fo0bu/FbuKpbc=",
        version = "v3.0.4",
    )
    go_repository(
        name = "io_opencensus_go",
        importpath = "go.opencensus.io",
        sum = "h1:y73uSU6J157QMP2kn2r30vwW1A2W2WFwSCGnAVxeaD0=",
        version = "v0.24.0",
    )
    go_repository(
        name = "io_opentelemetry_go_auto_sdk",
        importpath = "go.opentelemetry.io/auto/sdk",
        sum = "h1:jXsnJ4Lmnqd11kwkBV2LgLoFMZKizbCi5fNZ/ipaZ64=",
        version = "v1.2.1",
    )
    go_repository(
        name = "io_opentelemetry_go_contrib_detectors_gcp",
        importpath = "go.opentelemetry.io/contrib/detectors/gcp",
        sum = "h1:kWRNZMsfBHZ+uHjiH4y7Etn2FK26LAGkNFw7RHv1DhE=",
        version = "v1.39.0",
    )
    go_repository(
        name = "io_opentelemetry_go_contrib_instrumentation_google_golang_org_grpc_otelgrpc",
        build_directives = [
            "gazelle:resolve go go.opentelemetry.io/otel @io_opentelemetry_go_otel//:go_default_library",
            "gazelle:resolve go go.opentelemetry.io/otel/attribute @io_opentelemetry_go_otel//attribute",
            "gazelle:resolve go go.opentelemetry.io/otel/codes @io_opentelemetry_go_otel//codes",
            "gazelle:resolve go go.opentelemetry.io/otel/metric @io_opentelemetry_go_otel_metric//:go_default_library",
            "gazelle:resolve go go.opentelemetry.io/otel/metric/noop @io_opentelemetry_go_otel_metric//noop:go_default_library",
            "gazelle:resolve go go.opentelemetry.io/otel/propagation @io_opentelemetry_go_otel//propagation",
            "gazelle:resolve go go.opentelemetry.io/otel/semconv/v1.17.0 @io_opentelemetry_go_otel//semconv/v1.17.0:v1_17_0",
            "gazelle:resolve go go.opentelemetry.io/otel/trace @io_opentelemetry_go_otel_trace//:go_default_library",
        ],
        importpath = "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc",
        sum = "h1:yI1/OhfEPy7J9eoa6Sj051C7n5dvpj0QX8g4sRchg04=",
        version = "v0.67.0",
    )
    go_repository(
        name = "io_opentelemetry_go_contrib_instrumentation_net_http_otelhttp",
        build_directives = [
            "gazelle:resolve go go.opentelemetry.io/otel @io_opentelemetry_go_otel//:go_default_library",
            "gazelle:resolve go go.opentelemetry.io/otel/attribute @io_opentelemetry_go_otel//attribute",
            "gazelle:resolve go go.opentelemetry.io/otel/metric @io_opentelemetry_go_otel_metric//:go_default_library",
            "gazelle:resolve go go.opentelemetry.io/otel/propagation @io_opentelemetry_go_otel//propagation",
            "gazelle:resolve go go.opentelemetry.io/otel/semconv/v1.20.0 @io_opentelemetry_go_otel//semconv/v1.20.0:v1_20_0",
            "gazelle:resolve go go.opentelemetry.io/otel/trace @io_opentelemetry_go_otel_trace//:go_default_library",
        ],
        importpath = "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp",
        sum = "h1:OyrsyzuttWTSur2qN/Lm0m2a8yqyIjUVBZcxFPuXq2o=",
        version = "v0.67.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel",
        build_directives = [
            "gazelle:resolve go go.opentelemetry.io/otel/trace @io_opentelemetry_go_otel_trace//:go_default_library",
            "gazelle:resolve go go.opentelemetry.io/otel/metric @io_opentelemetry_go_otel_metric//:go_default_library",
            "gazelle:resolve go go.opentelemetry.io/otel/metric/embedded @io_opentelemetry_go_otel_metric//embedded:go_default_library",
        ],
        importpath = "go.opentelemetry.io/otel",
        sum = "h1:mYIM03dnh5zfN7HautFE4ieIig9amkNANT+xcVxAj9I=",
        version = "v1.43.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_metric",
        build_directives = [
            "gazelle:resolve go go.opentelemetry.io/otel/attribute @io_opentelemetry_go_otel//attribute:go_default_library",
        ],
        importpath = "go.opentelemetry.io/otel/metric",
        sum = "h1:d7638QeInOnuwOONPp4JAOGfbCEpYb+K6DVWvdxGzgM=",
        version = "v1.43.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_sdk",
        importpath = "go.opentelemetry.io/otel/sdk",
        sum = "h1:pi5mE86i5rTeLXqoF/hhiBtUNcrAGHLKQdhg4h4V9Dg=",
        version = "v1.43.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_sdk_metric",
        importpath = "go.opentelemetry.io/otel/sdk/metric",
        sum = "h1:S88dyqXjJkuBNLeMcVPRFXpRw2fuwdvfCGLEo89fDkw=",
        version = "v1.43.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_trace",
        build_directives = [
            "gazelle:resolve go go.opentelemetry.io/otel/attribute @io_opentelemetry_go_otel//attribute:go_default_library",
        ],
        importpath = "go.opentelemetry.io/otel/trace",
        sum = "h1:BkNrHpup+4k4w+ZZ86CZoHHEkohws8AY+WTX09nk+3A=",
        version = "v1.43.0",
    )
    go_repository(
        name = "org_golang_google_api",
        build_directives = [
            "gazelle:resolve go go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp @io_opentelemetry_go_contrib_instrumentation_net_http_otelhttp//:go_default_library",
            "gazelle:resolve go go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc @io_opentelemetry_go_contrib_instrumentation_google_golang_org_grpc_otelgrpc//:go_default_library",
        ],
        importpath = "google.golang.org/api",
        sum = "h1:nVArUtfLEihtW+b0DdcqRGK1xoEm2+ltAihyztq7MKY=",
        version = "v0.276.0",
    )
    go_repository(
        name = "org_golang_google_appengine",
        importpath = "google.golang.org/appengine",
        sum = "h1:IhEN5q69dyKagZPYMSdIjS2HqprW324FRQZJcGqPAsM=",
        version = "v1.6.8",
    )
    go_repository(
        name = "org_golang_google_genproto",
        importpath = "google.golang.org/genproto",
        sum = "h1:QoMBg0moLIlB/eucPzc+ID5SgPZWuirtjAn3l8nW2Dg=",
        version = "v0.0.0-20260420184626-e10c466a9529",
    )
    go_repository(
        name = "org_golang_google_genproto_googleapis_api",
        importpath = "google.golang.org/genproto/googleapis/api",
        sum = "h1:yQugLulqltosq0B/f8l4w9VryjV+N/5gcW0jQ3N8Qec=",
        version = "v0.0.0-20260414002931-afd174a4e478",
    )
    go_repository(
        name = "org_golang_google_genproto_googleapis_bytestream",
        importpath = "google.golang.org/genproto/googleapis/bytestream",
        sum = "h1:6jJ8xNMxB3mmH3HsxEMPjJU+WeFOiwBJ+cLm60OvSZs=",
        version = "v0.0.0-20260319201613-d00831a3d3e7",
    )
    go_repository(
        name = "org_golang_google_genproto_googleapis_rpc",
        importpath = "google.golang.org/genproto/googleapis/rpc",
        sum = "h1:XF8+t6QQiS0o9ArVan/HW8Q7cycNPGsJf6GA2nXxYAg=",
        version = "v0.0.0-20260420184626-e10c466a9529",
    )
    go_repository(
        name = "org_golang_google_grpc",
        importpath = "google.golang.org/grpc",
        sum = "h1:Xr6m2WmWZLETvUNvIUmeD5OAagMw3FiKmMlTdViWsHM=",
        version = "v1.80.0",
    )

    #keep: frozen due to https://github.com/googleapis/gapic-generator-go/issues/1608
    go_repository(
        name = "org_golang_google_protobuf",
        build_directives = [
            "gazelle:proto disable",  # https://github.com/bazelbuild/rules_go/issues/3906
        ],
        build_extra_args = [
            "-exclude=**/testdata",
        ],
        importpath = "google.golang.org/protobuf",
        sum = "h1:8Ar7bF+apOIoThw1EdZl0p1oWvMqTHmpA2fRTyZO8io=",
        # TODO(https://github.com/googleapis/gapic-generator-go/issues/1608): Don't hard-code old version
        version = "v1.35.2",
    )
    go_repository(
        name = "org_golang_x_crypto",
        importpath = "golang.org/x/crypto",
        sum = "h1:+Ng2ULVvLHnJ/ZFEq4KdcDd/cfjrrjjNSXNzxg0Y4U4=",
        version = "v0.49.0",
    )
    go_repository(
        name = "org_golang_x_mod",
        importpath = "golang.org/x/mod",
        sum = "h1:tHFzIWbBifEmbwtGz65eaWyGiGZatSrT9prnU8DbVL8=",
        version = "v0.33.0",
    )
    go_repository(
        name = "org_golang_x_net",
        importpath = "golang.org/x/net",
        sum = "h1:He/TN1l0e4mmR3QqHMT2Xab3Aj3L9qjbhRm78/6jrW0=",
        version = "v0.52.0",
    )
    go_repository(
        name = "org_golang_x_oauth2",
        importpath = "golang.org/x/oauth2",
        sum = "h1:peZ/1z27fi9hUOFCAZaHyrpWG5lwe0RJEEEeH0ThlIs=",
        version = "v0.36.0",
    )
    go_repository(
        name = "org_golang_x_sync",
        importpath = "golang.org/x/sync",
        sum = "h1:e0PTpb7pjO8GAtTs2dQ6jYa5BWYlMuX047Dco/pItO4=",
        version = "v0.20.0",
    )
    go_repository(
        name = "org_golang_x_sys",
        importpath = "golang.org/x/sys",
        sum = "h1:omrd2nAlyT5ESRdCLYdm3+fMfNFE/+Rf4bDIQImRJeo=",
        version = "v0.42.0",
    )
    go_repository(
        name = "org_golang_x_term",
        importpath = "golang.org/x/term",
        sum = "h1:QCgPso/Q3RTJx2Th4bDLqML4W6iJiaXFq2/ftQF13YU=",
        version = "v0.41.0",
    )
    go_repository(
        name = "org_golang_x_text",
        importpath = "golang.org/x/text",
        sum = "h1:JOVx6vVDFokkpaq1AEptVzLTpDe9KGpj5tR4/X+ybL8=",
        version = "v0.35.0",
    )
    go_repository(
        name = "org_golang_x_time",
        importpath = "golang.org/x/time",
        sum = "h1:bbrp8t3bGUeFOx08pvsMYRTCVSMk89u4tKbNOZbp88U=",
        version = "v0.15.0",
    )
    go_repository(
        name = "org_golang_x_tools",
        importpath = "golang.org/x/tools",
        sum = "h1:uNgphsn75Tdz5Ji2q36v/nsFSfR/9BRFvqhGBaJGd5k=",
        version = "v0.42.0",
    )
    go_repository(
        name = "org_gonum_v1_gonum",
        importpath = "gonum.org/v1/gonum",
        sum = "h1:VbpOemQlsSMrYmn7T2OUvQ4dqxQXU+ouZFQsZOx50z4=",
        version = "v0.17.0",
    )

def _maybe(repo_rule, name, strip_repo_prefix = "", **kwargs):
    if not name.startswith(strip_repo_prefix):
        return
    repo_name = name[len(strip_repo_prefix):]
    if repo_name in native.existing_rules():
        return
    repo_rule(name = repo_name, **kwargs)

# Redefine gazelle's go_repository rule with our own that will wrap go_repository
# with the _maybe macro. gazelle update-repos can only generate go_reposoritory
# targets. This way we get both the _maybe functionality and still use
# gazelle update-repos. The real go_repository rule is loaded with an alias:
# gazelle_go_repository.
def go_repository(name, importpath, sum, version, build_file_proto_mode = "", build_extra_args = [], build_directives = []):
    _maybe(
        gazelle_go_repository,
        name = name,
        importpath = importpath,
        sum = sum,
        version = version,
        build_file_proto_mode = build_file_proto_mode,
        build_extra_args = build_extra_args,
        build_directives = build_directives,
    )
