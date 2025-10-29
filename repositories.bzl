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
        sum = "h1:aQ3y1lwWyqYPiWZThqv1aFbZMiM9vblcSArJRf2Irls=",
        version = "v0.0.0-20250501225837-2ac532fd4443",
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
        sum = "h1:zEqyPVyku6IvWCFwux4x9RxkLOMUL+1vC9xUFv5l2/M=",
        version = "v0.13.4",
    )
    go_repository(
        name = "com_github_envoyproxy_go_control_plane_envoy",
        importpath = "github.com/envoyproxy/go-control-plane/envoy",
        sum = "h1:jb83lalDRZSpPWW2Z7Mck/8kXZ5CQAFYVjQcdVIr83A=",
        version = "v1.32.4",
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
        sum = "h1:DEo3O99U8j4hBFwbJfrz9VtgcDfUKS7KJ7spH3d86P8=",
        version = "v1.2.1",
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
        sum = "h1:dAwr6QBTBZIkG8roQaJjGof0pp0EeF+tNV7YBP3F/8M=",
        version = "v1.8.0",
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
        sum = "h1:TK/7NqRQZfgAh+Td8AlsrvtPoUyiHh0LqVvokh+1vHI=",
        version = "v4.1.2",
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
        sum = "h1:ZAaOCxANMuZx5RCeg0mBdEZk7DZasvvZIxtHqx8aGss=",
        version = "v2.2.1",
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
        sum = "h1:GW/XbdyBFQ8Qe+YAmFU9uHLo7OnF5tL52HFAgMmyrf4=",
        version = "v0.3.6",
    )
    go_repository(
        name = "com_github_googleapis_gapic_showcase",
        importpath = "github.com/googleapis/gapic-showcase",
        sum = "h1:smxzCetYS6PD0ME1ryJeQ6zvFFxcpMRbzoWI20/y+hU=",
        version = "v0.37.0",
    )
    go_repository(
        name = "com_github_googleapis_gax_go_v2",
        build_directives = [
            "gazelle:resolve go google.golang.org/genproto/googleapis/rpc/errdetails @org_golang_google_genproto_googleapis_rpc//errdetails",
            "gazelle:resolve proto go google/rpc/code.proto @com_google_googleapis//google/rpc:code_go_proto",
            "gazelle:resolve proto proto google/rpc/code.proto @com_google_googleapis//google/rpc:code_proto",
        ],
        importpath = "github.com/googleapis/gax-go/v2",
        sum = "h1:SyjDc1mGgZU5LncH8gimWo9lW1DtIfPibOG81vgd/bo=",
        version = "v2.15.0",
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
        sum = "h1:UQUsRi8WTzhZntp5313l+CHIAT95ojUI2lpP/ExlZa4=",
        version = "v1.29.0",
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
        sum = "h1:qOEr613fac2lOuTgWN4tPAtLL7fUSbuJL5X5XumQh94=",
        version = "v1.17.0",
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
        sum = "h1:YmeHyLY8mFWbdkNWwpr+qIL2bEqT0o95WSdkNHvL12M=",
        version = "v2.2.3",
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
        sum = "h1:5MqpDsTGNDhY8sGp0Aowyf0qKsPrhewaLSsFaodPcyo=",
        version = "v0.7.0",
    )
    go_repository(
        name = "com_github_soheilhy_cmux",
        importpath = "github.com/soheilhy/cmux",
        sum = "h1:jjzc5WVemNEDTLwv9tlmemhC73tI08BNOIGwBOo10Js=",
        version = "v0.1.5",
    )
    go_repository(
        name = "com_github_sourcegraph_conc",
        importpath = "github.com/sourcegraph/conc",
        sum = "h1:OQTbbt6P72L20UqAkXXuLOj79LfEanQ+YQFNpLA9ySo=",
        version = "v0.3.0",
    )
    go_repository(
        name = "com_github_spf13_afero",
        importpath = "github.com/spf13/afero",
        sum = "h1:UcOPyRBYczmFn6yvphxkn9ZEOY65cpwGKb5mL36mrqs=",
        version = "v1.12.0",
    )
    go_repository(
        name = "com_github_spf13_cast",
        importpath = "github.com/spf13/cast",
        sum = "h1:cuNEagBQEHWN1FnbGEjCXL2szYEXqfJPbP2HNUaca9Y=",
        version = "v1.7.1",
    )
    go_repository(
        name = "com_github_spf13_cobra",
        importpath = "github.com/spf13/cobra",
        sum = "h1:CXSaggrXdbHK9CF+8ywj8Amf7PBRmPCOJugH954Nnlo=",
        version = "v1.9.1",
    )
    go_repository(
        name = "com_github_spf13_pflag",
        importpath = "github.com/spf13/pflag",
        sum = "h1:jFzHGLGAlb3ruxLB8MhbI6A8+AQX/2eW4qeyNZXNp2o=",
        version = "v1.0.6",
    )
    go_repository(
        name = "com_github_spf13_viper",
        importpath = "github.com/spf13/viper",
        sum = "h1:ZMi+z/lvLyPSCoNtFCpqjy0S4kPbirhpTMwl8BkW9X4=",
        version = "v1.20.1",
    )
    go_repository(
        name = "com_github_spiffe_go_spiffe_v2",
        importpath = "github.com/spiffe/go-spiffe/v2",
        sum = "h1:N2I01KCUkv1FAjZXJMwh95KK1ZIQLYbPfhaxw8WS0hE=",
        version = "v2.5.0",
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
        sum = "h1:gq8OS+rQWgGRo91D2qztN+ion6AZ2T1CxBIu0ifCmVo=",
        version = "v1.8.8",
    )
    go_repository(
        name = "com_google_cloud_go_accesscontextmanager",
        importpath = "cloud.google.com/go/accesscontextmanager",
        sum = "h1:aKIfg7Jyc73pe8bzx0zypNdS5gfFdSvFvB8YNA9k2kA=",
        version = "v1.9.7",
    )
    go_repository(
        name = "com_google_cloud_go_aiplatform",
        importpath = "cloud.google.com/go/aiplatform",
        sum = "h1:XDoc+xCpBiBaeplrCLoPBSSPwozOmogPYkGuBUZKamY=",
        version = "v1.107.0",
    )
    go_repository(
        name = "com_google_cloud_go_analytics",
        importpath = "cloud.google.com/go/analytics",
        sum = "h1:souLxu9tQHzF+0NDpKoIw4pl2WQ9K2JfkdPPs36BfXw=",
        version = "v0.30.1",
    )
    go_repository(
        name = "com_google_cloud_go_apigateway",
        importpath = "cloud.google.com/go/apigateway",
        sum = "h1:ehKUTy+QFsb3n07fEi18S2dpDDjCV4UlRyrbwfZV3Zk=",
        version = "v1.7.7",
    )
    go_repository(
        name = "com_google_cloud_go_apigeeconnect",
        importpath = "cloud.google.com/go/apigeeconnect",
        sum = "h1:S6s2zojwMymx0fyZYKm0eK1TdDxrriIBAlNVvRAOzug=",
        version = "v1.7.7",
    )
    go_repository(
        name = "com_google_cloud_go_apigeeregistry",
        importpath = "cloud.google.com/go/apigeeregistry",
        sum = "h1:QziFVsuPU2lhy40Ht9uWEyciV23SH9GETWiwcu3qzdg=",
        version = "v0.10.0",
    )
    go_repository(
        name = "com_google_cloud_go_appengine",
        importpath = "cloud.google.com/go/appengine",
        sum = "h1:IxGz6j5xv0nTJX285wu95Vn6KEi2CeV9vbyRgCSEAoU=",
        version = "v1.9.7",
    )
    go_repository(
        name = "com_google_cloud_go_area120",
        importpath = "cloud.google.com/go/area120",
        sum = "h1:BbpzLwaIXVPorrrzTH+ni7P5mLemmPPfSZ7o39k7zQc=",
        version = "v0.9.7",
    )
    go_repository(
        name = "com_google_cloud_go_artifactregistry",
        importpath = "cloud.google.com/go/artifactregistry",
        sum = "h1:Gx5vsnIFEx+obM1VdtMF2AuTraYESRCrBxvc9+6jBZg=",
        version = "v1.17.2",
    )
    go_repository(
        name = "com_google_cloud_go_asset",
        importpath = "cloud.google.com/go/asset",
        sum = "h1:81Ru5hjHfiGtk+u/Ix69eaWieKpvm7Ce7UHtcZhOLbk=",
        version = "v1.22.0",
    )
    go_repository(
        name = "com_google_cloud_go_assuredworkloads",
        importpath = "cloud.google.com/go/assuredworkloads",
        sum = "h1:NQXyyGLksPmiapE1Oc64a3cMwYIBAoDBg6cWR+B3eaY=",
        version = "v1.13.0",
    )
    go_repository(
        name = "com_google_cloud_go_auth",
        importpath = "cloud.google.com/go/auth",
        sum = "h1:74yCm7hCj2rUyyAocqnFzsAYXgJhrG26XCFimrc/Kz4=",
        version = "v0.17.0",
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
        sum = "h1:YRwLbsBv4yApX64pkrdyy4emhWE6lHEnljX4b1aTQC4=",
        version = "v1.15.0",
    )
    go_repository(
        name = "com_google_cloud_go_baremetalsolution",
        importpath = "cloud.google.com/go/baremetalsolution",
        sum = "h1:g67fjVdrNCHZl8jDWdZvo+6zGTTMMuvNWO7HSgG8lnI=",
        version = "v1.4.0",
    )
    go_repository(
        name = "com_google_cloud_go_batch",
        importpath = "cloud.google.com/go/batch",
        sum = "h1:6gmnNhJhm+Y598CZE6x+CeNZNIPCkYa4vkF58S5abHo=",
        version = "v1.13.0",
    )
    go_repository(
        name = "com_google_cloud_go_beyondcorp",
        importpath = "cloud.google.com/go/beyondcorp",
        sum = "h1:mre997ya7QHFWSU+O5cT/FhBKTMy6Riqf1EXFxN46zw=",
        version = "v1.2.0",
    )
    go_repository(
        name = "com_google_cloud_go_bigquery",
        importpath = "cloud.google.com/go/bigquery",
        sum = "h1:NvSZvXU1Hyb+YiRVKQPuQXGeZaw/0NP6M/WOrBqSx3g=",
        version = "v1.71.0",
    )
    go_repository(
        name = "com_google_cloud_go_bigtable",
        importpath = "cloud.google.com/go/bigtable",
        sum = "h1:k8HfpUOvn7sQwc6oNKqjvD/yjkwynf4qBuyKwh5cU08=",
        version = "v1.40.1",
    )
    go_repository(
        name = "com_google_cloud_go_billing",
        importpath = "cloud.google.com/go/billing",
        sum = "h1:nbQjTXkpgB/E4XnYZQwcZnR63QFsbFwJ9DGsNg61Ghg=",
        version = "v1.21.0",
    )
    go_repository(
        name = "com_google_cloud_go_binaryauthorization",
        importpath = "cloud.google.com/go/binaryauthorization",
        sum = "h1:YYK0BwiZv9uA6z+Ict908AykX4OBfDECMTE476OnS3A=",
        version = "v1.10.0",
    )
    go_repository(
        name = "com_google_cloud_go_certificatemanager",
        importpath = "cloud.google.com/go/certificatemanager",
        sum = "h1:v5X8X+THKrS9OFZb6k0GRDP1WQxLXTdMko7OInBliw4=",
        version = "v1.9.6",
    )
    go_repository(
        name = "com_google_cloud_go_channel",
        importpath = "cloud.google.com/go/channel",
        sum = "h1:EeUa6SnD3+EL9B06G6N9Ud5/p/NtT6PC7lv5kmaUiHs=",
        version = "v1.20.0",
    )
    go_repository(
        name = "com_google_cloud_go_cloudbuild",
        importpath = "cloud.google.com/go/cloudbuild",
        sum = "h1:Kl4QBrOPXcHVTic6XeRMp9YgLCy3b/ifGx8i29A3pYs=",
        version = "v1.23.1",
    )
    go_repository(
        name = "com_google_cloud_go_clouddms",
        importpath = "cloud.google.com/go/clouddms",
        sum = "h1:YWsmRXTyK6Ba0hm4qTBak5g1oLhryuM8rSBxHWC8iq4=",
        version = "v1.8.8",
    )
    go_repository(
        name = "com_google_cloud_go_cloudtasks",
        importpath = "cloud.google.com/go/cloudtasks",
        sum = "h1:H2v8GEolNtMFfYzUpZBaZbydqU7drpyo99GtAgA+m4I=",
        version = "v1.13.7",
    )
    go_repository(
        name = "com_google_cloud_go_compute",
        importpath = "cloud.google.com/go/compute",
        sum = "h1:KYKIG0+pfpAWaAYayFkE/KPrAVCge0Hu82bPraAmsCk=",
        version = "v1.49.1",
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
        sum = "h1:wA4j99BhsoeYlLx6xEIqrNN1aOTtUme0wimHZegg80s=",
        version = "v1.17.4",
    )
    go_repository(
        name = "com_google_cloud_go_container",
        importpath = "cloud.google.com/go/container",
        sum = "h1:5h6FiGFXq7GUSwJH7TYE2Sy0xiEqlhiD8i14DICGiC4=",
        version = "v1.44.1",
    )
    go_repository(
        name = "com_google_cloud_go_containeranalysis",
        importpath = "cloud.google.com/go/containeranalysis",
        sum = "h1:OW2dlMPtR5VnjQGyAP+uJlZahc1l+JFxFlH/J3+l7gw=",
        version = "v0.14.2",
    )
    go_repository(
        name = "com_google_cloud_go_datacatalog",
        importpath = "cloud.google.com/go/datacatalog",
        sum = "h1:bCRKA8uSQN8wGW3Tw0gwko4E9a64GRmbW1nCblhgC2k=",
        version = "v1.26.1",
    )
    go_repository(
        name = "com_google_cloud_go_dataflow",
        importpath = "cloud.google.com/go/dataflow",
        sum = "h1:Z+UYlGrE+IoB+5IAN4/qWdPKO0IpIK9bs2Dy40HK6lg=",
        version = "v0.11.1",
    )
    go_repository(
        name = "com_google_cloud_go_dataform",
        importpath = "cloud.google.com/go/dataform",
        sum = "h1:yf6Up6m1FUt+YB5CBgNtIZfz2OzjuNBdzZWV3SLSVNE=",
        version = "v0.12.1",
    )
    go_repository(
        name = "com_google_cloud_go_datafusion",
        importpath = "cloud.google.com/go/datafusion",
        sum = "h1:tLCV+xYuOrSjdrRTkc9Cqsb5mBSQEsNfFmuTNYl5/rA=",
        version = "v1.8.7",
    )
    go_repository(
        name = "com_google_cloud_go_datalabeling",
        importpath = "cloud.google.com/go/datalabeling",
        sum = "h1:wwoct7mw38s75XvEmLoItQ2TY0RFsGiRDb0iNbXUcX4=",
        version = "v0.9.7",
    )
    go_repository(
        name = "com_google_cloud_go_dataplex",
        importpath = "cloud.google.com/go/dataplex",
        sum = "h1:renSEYTQZMQ3ag7lM0BDmSj4FWqaTGW60YQ/lvAE5iA=",
        version = "v1.27.1",
    )
    go_repository(
        name = "com_google_cloud_go_dataproc_v2",
        importpath = "cloud.google.com/go/dataproc/v2",
        sum = "h1:I/Yux/d8uaxf3W+d59kolGTOc52+VZaL6RzJw7oDOeg=",
        version = "v2.15.0",
    )
    go_repository(
        name = "com_google_cloud_go_dataqna",
        importpath = "cloud.google.com/go/dataqna",
        sum = "h1:3FREvU+sjaEHSjlKrKF6KjUmafdOvM8CbZ897rttxNs=",
        version = "v0.9.8",
    )
    go_repository(
        name = "com_google_cloud_go_datastore",
        importpath = "cloud.google.com/go/datastore",
        sum = "h1:NNpXoyEqIJmZFc0ACcwBEaXnmscUpcG4NkKnbCePmiM=",
        version = "v1.20.0",
    )
    go_repository(
        name = "com_google_cloud_go_datastream",
        importpath = "cloud.google.com/go/datastream",
        sum = "h1:7PKeDpksi8nbOR4gspmNokzsr0q/uRzDIt20bR3BtRs=",
        version = "v1.15.1",
    )
    go_repository(
        name = "com_google_cloud_go_deploy",
        importpath = "cloud.google.com/go/deploy",
        sum = "h1:QU8gLXsXDRqLyEWNrI6zJiVzuuOBX/WpMi4p0oexV+c=",
        version = "v1.27.3",
    )
    go_repository(
        name = "com_google_cloud_go_dialogflow",
        importpath = "cloud.google.com/go/dialogflow",
        sum = "h1:76UcutgQuV8s1t2yIz6rdUCkosDT6nT3LCsybPzJI8k=",
        version = "v1.70.0",
    )
    go_repository(
        name = "com_google_cloud_go_dlp",
        importpath = "cloud.google.com/go/dlp",
        sum = "h1:rBQCAcIdWxyLRbFsOkyLig0OFhcKhCThtnro+b3dvEc=",
        version = "v1.27.0",
    )
    go_repository(
        name = "com_google_cloud_go_documentai",
        importpath = "cloud.google.com/go/documentai",
        sum = "h1:ngPNmVUb3xxj6/j7ruK4vmOMm+XTJO0Vi7Jkm0R3AdU=",
        version = "v1.39.0",
    )
    go_repository(
        name = "com_google_cloud_go_domains",
        importpath = "cloud.google.com/go/domains",
        sum = "h1:G3kUq0vKBMhyOj5GqAfEYbVuez05U+ENHZUAtrEp/pI=",
        version = "v0.10.7",
    )
    go_repository(
        name = "com_google_cloud_go_edgecontainer",
        importpath = "cloud.google.com/go/edgecontainer",
        sum = "h1:6KTQo6Qf0iEtfPVotlG7orazEO1I93Ham0PMlkHYpdQ=",
        version = "v1.4.4",
    )
    go_repository(
        name = "com_google_cloud_go_errorreporting",
        importpath = "cloud.google.com/go/errorreporting",
        sum = "h1:isaoPwWX8kbAOea4qahcmttoS79+gQhvKsfg5L5AgH8=",
        version = "v0.3.2",
    )
    go_repository(
        name = "com_google_cloud_go_essentialcontacts",
        importpath = "cloud.google.com/go/essentialcontacts",
        sum = "h1:v9sO4IHFuwplaOuDnEXZFtfOrjw2bi11TSIVp5PnAU4=",
        version = "v1.7.7",
    )
    go_repository(
        name = "com_google_cloud_go_eventarc",
        importpath = "cloud.google.com/go/eventarc",
        sum = "h1:cog4hh0RNbVC7KVwibCqBPhBKklj4wnIJVz8Qrf2hVI=",
        version = "v1.17.0",
    )
    go_repository(
        name = "com_google_cloud_go_filestore",
        importpath = "cloud.google.com/go/filestore",
        sum = "h1:3KZifUVTqGhNNv6MLeONYth1HjlVM4vDhaH+xrdPljU=",
        version = "v1.10.3",
    )
    go_repository(
        name = "com_google_cloud_go_firestore",
        importpath = "cloud.google.com/go/firestore",
        sum = "h1:E3FiRsWfZKwZ6W+Lsp1YqTzZ9H6jP+QsKW40KR21C8I=",
        version = "v1.19.0",
    )
    go_repository(
        name = "com_google_cloud_go_functions",
        importpath = "cloud.google.com/go/functions",
        sum = "h1:7LcOD18euIVGRUPaeCmgO6vfWSLNIsi6STWRQcdANG8=",
        version = "v1.19.7",
    )
    go_repository(
        name = "com_google_cloud_go_gkebackup",
        importpath = "cloud.google.com/go/gkebackup",
        sum = "h1:gUgI3lZJYALZsHXE7YJOKI8bMpoAX/tF6jnNugvzT1g=",
        version = "v1.8.1",
    )
    go_repository(
        name = "com_google_cloud_go_gkeconnect",
        importpath = "cloud.google.com/go/gkeconnect",
        sum = "h1:EFql3zRaFw74yATt5lf+mcPDqPZ4EeLvoIJ+0NaEkag=",
        version = "v0.12.5",
    )
    go_repository(
        name = "com_google_cloud_go_gkehub",
        importpath = "cloud.google.com/go/gkehub",
        sum = "h1:Jk5pAXG54FlQzTRXhuKyym/NzOgS8oWRs0XNatZYDf4=",
        version = "v0.16.0",
    )
    go_repository(
        name = "com_google_cloud_go_gkemulticloud",
        importpath = "cloud.google.com/go/gkemulticloud",
        sum = "h1:AKuOmr5QBCPLJCyZhBABP3lIz+h3jxAy53LVMrEuvlg=",
        version = "v1.5.4",
    )
    go_repository(
        name = "com_google_cloud_go_gsuiteaddons",
        importpath = "cloud.google.com/go/gsuiteaddons",
        sum = "h1:Dayrv57XW8kZIvmQjAc89Tp7Kr3O9Am/hf6pXkTjYFY=",
        version = "v1.7.8",
    )
    go_repository(
        name = "com_google_cloud_go_iam",
        importpath = "cloud.google.com/go/iam",
        sum = "h1:+vMINPiDF2ognBJ97ABAYYwRgsaqxPbQDlMnbHMjolc=",
        version = "v1.5.3",
    )
    go_repository(
        name = "com_google_cloud_go_iap",
        importpath = "cloud.google.com/go/iap",
        sum = "h1:Nheb77nO0/pECm/thoE3wHVAbkQSI+G8KBWviqBepiA=",
        version = "v1.11.3",
    )
    go_repository(
        name = "com_google_cloud_go_ids",
        importpath = "cloud.google.com/go/ids",
        sum = "h1:V0pSk+KKW+5/AVpeQMhM9D1VI7aMZkayj5jddNETJos=",
        version = "v1.5.7",
    )
    go_repository(
        name = "com_google_cloud_go_iot",
        importpath = "cloud.google.com/go/iot",
        sum = "h1:PDUtxCzlFwFHODEFAgaGJy/Zv4tdvLbZ+lvZ1mKQXE4=",
        version = "v1.8.7",
    )
    go_repository(
        name = "com_google_cloud_go_kms",
        importpath = "cloud.google.com/go/kms",
        sum = "h1:4IYDQL5hG4L+HzJBhzejUySoUOheh3Lk5YT4PCyyW6k=",
        version = "v1.23.2",
    )
    go_repository(
        name = "com_google_cloud_go_language",
        importpath = "cloud.google.com/go/language",
        sum = "h1:/0Fbd3/T4oNmpPqIq5/hrWdHc/eoYGtVH5lDNkuHH3k=",
        version = "v1.14.6",
    )
    go_repository(
        name = "com_google_cloud_go_lifesciences",
        importpath = "cloud.google.com/go/lifesciences",
        sum = "h1:MO5aBahcYv7JeuCpHbg/11h7KL/BYt1+PpgHhleLDbI=",
        version = "v0.10.7",
    )
    go_repository(
        name = "com_google_cloud_go_logging",
        importpath = "cloud.google.com/go/logging",
        sum = "h1:7j0HgAp0B94o1YRDqiqm26w4q1rDMH7XNRU34lJXHYc=",
        version = "v1.13.0",
    )
    go_repository(
        name = "com_google_cloud_go_longrunning",
        importpath = "cloud.google.com/go/longrunning",
        sum = "h1:FV0+SYF1RIj59gyoWDRi45GiYUMM3K1qO51qoboQT1E=",
        version = "v0.7.0",
    )
    go_repository(
        name = "com_google_cloud_go_managedidentities",
        importpath = "cloud.google.com/go/managedidentities",
        sum = "h1:vC/q7D+97PZfb0UNf7r/+/clHauuaf1PqWwP7neuaeg=",
        version = "v1.7.7",
    )
    go_repository(
        name = "com_google_cloud_go_maps",
        importpath = "cloud.google.com/go/maps",
        sum = "h1:5G4ZtU8Wjy2grR3tbqMD6IW6uTeGHYIgt/glLmrupcA=",
        version = "v1.24.0",
    )
    go_repository(
        name = "com_google_cloud_go_mediatranslation",
        importpath = "cloud.google.com/go/mediatranslation",
        sum = "h1:JXbjms+JxgaWkj/YuaQm1OeCzuF+IZCDV17uUcZgFOU=",
        version = "v0.9.7",
    )
    go_repository(
        name = "com_google_cloud_go_memcache",
        importpath = "cloud.google.com/go/memcache",
        sum = "h1:ZDIfIMZsKKPzwdbvTMOL1il0shX24J7B9DC+sEt4Yj4=",
        version = "v1.11.7",
    )
    go_repository(
        name = "com_google_cloud_go_metastore",
        importpath = "cloud.google.com/go/metastore",
        sum = "h1:nfyUDD9AeKIs6btY5buQ1No0OVco20WpX9wIruL8UOA=",
        version = "v1.14.8",
    )
    go_repository(
        name = "com_google_cloud_go_monitoring",
        importpath = "cloud.google.com/go/monitoring",
        sum = "h1:dde+gMNc0UhPZD1Azu6at2e79bfdztVDS5lvhOdsgaE=",
        version = "v1.24.3",
    )
    go_repository(
        name = "com_google_cloud_go_networkconnectivity",
        importpath = "cloud.google.com/go/networkconnectivity",
        sum = "h1:n0IzhdgSNzIKQygWwDV8yKRXkZpX3FsjCYFbO9iNHPU=",
        version = "v1.19.1",
    )
    go_repository(
        name = "com_google_cloud_go_networkmanagement",
        importpath = "cloud.google.com/go/networkmanagement",
        sum = "h1:W5zdnH332yJFADTXlsHWLVkiXLKGWHjrsg7vXLGd+Ws=",
        version = "v1.20.1",
    )
    go_repository(
        name = "com_google_cloud_go_networksecurity",
        importpath = "cloud.google.com/go/networksecurity",
        sum = "h1:J5gdG7mHdRLrsyM7yy4nKFgbN8+geaOo/4Zpeh4DWrg=",
        version = "v0.10.7",
    )
    go_repository(
        name = "com_google_cloud_go_notebooks",
        importpath = "cloud.google.com/go/notebooks",
        sum = "h1:g5LTI1LHa/86abDTWd8nrq7/4qq8oFhVx1SmnNpZLVg=",
        version = "v1.12.7",
    )
    go_repository(
        name = "com_google_cloud_go_optimization",
        importpath = "cloud.google.com/go/optimization",
        sum = "h1:dMtxINB6G7wULbdm8nZ/x1NMa579Q+GfJc5gaN8VeDw=",
        version = "v1.7.7",
    )
    go_repository(
        name = "com_google_cloud_go_orchestration",
        importpath = "cloud.google.com/go/orchestration",
        sum = "h1:TVWDiZyvcflLFeTQH2GexHmtJ6iUSjzr0zsSiT338dA=",
        version = "v1.11.10",
    )
    go_repository(
        name = "com_google_cloud_go_orgpolicy",
        importpath = "cloud.google.com/go/orgpolicy",
        sum = "h1:0hq12wxNwcfUMojr5j3EjWECSInIuyYDhkAWXTomRhc=",
        version = "v1.15.1",
    )
    go_repository(
        name = "com_google_cloud_go_osconfig",
        importpath = "cloud.google.com/go/osconfig",
        sum = "h1:QQzK5njfsfO2rdOWYVDyLQktqSq9gKf2ohRYeKUuA10=",
        version = "v1.15.1",
    )
    go_repository(
        name = "com_google_cloud_go_oslogin",
        importpath = "cloud.google.com/go/oslogin",
        sum = "h1:YQ8P/+MLwH0tpENYU9QOgwKQxe8DYfAKxIfm6y+OBtA=",
        version = "v1.14.7",
    )
    go_repository(
        name = "com_google_cloud_go_phishingprotection",
        importpath = "cloud.google.com/go/phishingprotection",
        sum = "h1:ZJqHirY2/H6s+uTq1y1iiVASzm3ZuDiMglT5NXywPBE=",
        version = "v0.9.7",
    )
    go_repository(
        name = "com_google_cloud_go_policytroubleshooter",
        importpath = "cloud.google.com/go/policytroubleshooter",
        sum = "h1:Bbj1EiVh96u9mfO2p+JNoHrvvyC0Ms6zP+vxqQnsaG8=",
        version = "v1.11.7",
    )
    go_repository(
        name = "com_google_cloud_go_privatecatalog",
        importpath = "cloud.google.com/go/privatecatalog",
        sum = "h1:yOdy85WDvSCPxAMixkhs5X0Z96D74kosgOTp7aJEYvU=",
        version = "v0.10.8",
    )
    go_repository(
        name = "com_google_cloud_go_pubsub",
        importpath = "cloud.google.com/go/pubsub",
        sum = "h1:fzbXpPyJnSGvWXF1jabhQeXyxdbCIkXTpjXHy7xviBM=",
        version = "v1.50.1",
    )
    go_repository(
        name = "com_google_cloud_go_pubsub_v2",
        importpath = "cloud.google.com/go/pubsub/v2",
        sum = "h1:0qS6mRJ41gD1lNmM/vdm6bR7DQu6coQcVwD+VPf0Bz0=",
        version = "v2.0.0",
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
        sum = "h1:Q2CcYGxcvnvng2q3o1SaOpV+rjE/AbFVYGTJomxlG4g=",
        version = "v2.20.5",
    )
    go_repository(
        name = "com_google_cloud_go_recommendationengine",
        importpath = "cloud.google.com/go/recommendationengine",
        sum = "h1:NH89CyKQP8e98kpdKLwV0jXkQGzSEEZia0V867vkoy8=",
        version = "v0.9.7",
    )
    go_repository(
        name = "com_google_cloud_go_recommender",
        importpath = "cloud.google.com/go/recommender",
        sum = "h1:ZVZg4wr1G7yzjIPcYUNSUJAaz9+2o78rmBU4QJgC7kg=",
        version = "v1.13.6",
    )
    go_repository(
        name = "com_google_cloud_go_redis",
        importpath = "cloud.google.com/go/redis",
        sum = "h1:6LI8zSt+vmE3WQ7hE5GsJ13CbJBLV1qUw6B7CY31Wcw=",
        version = "v1.18.3",
    )
    go_repository(
        name = "com_google_cloud_go_resourcemanager",
        importpath = "cloud.google.com/go/resourcemanager",
        sum = "h1:oPZKIdjyVTuag+D4HF7HO0mnSqcqgjcuA18xblwA0V0=",
        version = "v1.10.7",
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
        sum = "h1:S9Rl1sB6Dmdak5y9I2Pwn9hcWA5ubk33Vnevxw+CSrI=",
        version = "v1.25.1",
    )
    go_repository(
        name = "com_google_cloud_go_run",
        importpath = "cloud.google.com/go/run",
        sum = "h1:zoXZ+vavS6k8wzEPlxMuh5rGkhQb5CAzrcfSFBInlS4=",
        version = "v1.12.1",
    )
    go_repository(
        name = "com_google_cloud_go_scheduler",
        importpath = "cloud.google.com/go/scheduler",
        sum = "h1:BoXY2BvBsaRw3ggVMzC9tborZqJBu+NcJcD9PqeC5Kc=",
        version = "v1.11.8",
    )
    go_repository(
        name = "com_google_cloud_go_secretmanager",
        importpath = "cloud.google.com/go/secretmanager",
        sum = "h1:19QT7ZsLJ8FSP1k+4esQvuCD7npMJml6hYzilxVyT+k=",
        version = "v1.16.0",
    )
    go_repository(
        name = "com_google_cloud_go_security",
        importpath = "cloud.google.com/go/security",
        sum = "h1:cF3FkCRRbRC1oXuaGZFl3qU2sdu2gP3iOAHKzL5y04Y=",
        version = "v1.19.2",
    )
    go_repository(
        name = "com_google_cloud_go_securitycenter",
        importpath = "cloud.google.com/go/securitycenter",
        sum = "h1:D9zpeguY4frQU35GBw8+M6Gw79CiuTF9iVs4sFm3FDY=",
        version = "v1.38.1",
    )
    go_repository(
        name = "com_google_cloud_go_servicedirectory",
        importpath = "cloud.google.com/go/servicedirectory",
        sum = "h1:je2yZlVcVFI/TshPXjjF9ZAlWedj0s5EbO2kozJrzBo=",
        version = "v1.12.7",
    )
    go_repository(
        name = "com_google_cloud_go_shell",
        importpath = "cloud.google.com/go/shell",
        sum = "h1:K1C9sh9EuNNhGpyCoqRdeudcU9zmfYTA95bhF5cokK8=",
        version = "v1.8.7",
    )
    go_repository(
        name = "com_google_cloud_go_spanner",
        importpath = "cloud.google.com/go/spanner",
        sum = "h1:jlNWusBol1Jxa9PmYGknUBzLwvD1cebuEenzqebZ9xs=",
        version = "v1.86.0",
    )
    go_repository(
        name = "com_google_cloud_go_speech",
        importpath = "cloud.google.com/go/speech",
        sum = "h1:L8kq/CypGn6y/FbipyAPyn1L9JDW4CK3zkcAe4oDN7U=",
        version = "v1.28.1",
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
        sum = "h1:Sjukr1LtUt7vLTHNvGc2gaAqlXNFeDFRIRmWGrFaJlY=",
        version = "v1.13.1",
    )
    go_repository(
        name = "com_google_cloud_go_talent",
        importpath = "cloud.google.com/go/talent",
        sum = "h1:1kJJ+WCY5LZ1A4rCa32zKh3N2xT3I8koiS63+vV0WC4=",
        version = "v1.8.4",
    )
    go_repository(
        name = "com_google_cloud_go_texttospeech",
        importpath = "cloud.google.com/go/texttospeech",
        sum = "h1:v9By5kPtOPJT5ozw7+hLzSzfYdQnAAxOTSdh9vo27T4=",
        version = "v1.15.1",
    )
    go_repository(
        name = "com_google_cloud_go_tpu",
        importpath = "cloud.google.com/go/tpu",
        sum = "h1:5DDheA1f7yZ/KUbVT/9lL+Yhgd3IqHDSVVrSqDVkAFY=",
        version = "v1.8.4",
    )
    go_repository(
        name = "com_google_cloud_go_trace",
        importpath = "cloud.google.com/go/trace",
        sum = "h1:kDNDX8JkaAG3R2nq1lIdkb7FCSi1rCmsEtKVsty7p+U=",
        version = "v1.11.7",
    )
    go_repository(
        name = "com_google_cloud_go_translate",
        importpath = "cloud.google.com/go/translate",
        sum = "h1:aSxMbfJ3MVmEdQzu5jGXmPPxCAb1ySsor2yBMCI5MT4=",
        version = "v1.12.7",
    )
    go_repository(
        name = "com_google_cloud_go_video",
        importpath = "cloud.google.com/go/video",
        sum = "h1:Hp+2AeM7b3AagdHcyh2820UTzSbGyqpFJVMu0nHbBcw=",
        version = "v1.27.1",
    )
    go_repository(
        name = "com_google_cloud_go_videointelligence",
        importpath = "cloud.google.com/go/videointelligence",
        sum = "h1:FisUrSZ+y3oLuGdlFQQgZoNTDm7FAfb2hwSTsSqX+9g=",
        version = "v1.12.7",
    )
    go_repository(
        name = "com_google_cloud_go_vision_v2",
        importpath = "cloud.google.com/go/vision/v2",
        sum = "h1:9UtOINPF8p9VACQ6KAyR/ZtkpuBHGmJsprutYupDcN0=",
        version = "v2.9.6",
    )
    go_repository(
        name = "com_google_cloud_go_vmmigration",
        importpath = "cloud.google.com/go/vmmigration",
        sum = "h1:7IED5P1BmZIDy0U9Co0OYGQkFK/19/7Y6xlCaitmIcs=",
        version = "v1.9.1",
    )
    go_repository(
        name = "com_google_cloud_go_vmwareengine",
        importpath = "cloud.google.com/go/vmwareengine",
        sum = "h1:TKvULKbk44QrIx674cnoVjcZueXhyCAm2sNAJu/S1ds=",
        version = "v1.3.6",
    )
    go_repository(
        name = "com_google_cloud_go_vpcaccess",
        importpath = "cloud.google.com/go/vpcaccess",
        sum = "h1:K6siDR1T4HgSTv6sy6CAwupY7UGza6TQ1O8jtvEYoX4=",
        version = "v1.8.7",
    )
    go_repository(
        name = "com_google_cloud_go_webrisk",
        importpath = "cloud.google.com/go/webrisk",
        sum = "h1:q6zEdVgD8Ka+4fQl3azDcSNRug8clNnQ9iVS2iLh+MM=",
        version = "v1.11.2",
    )
    go_repository(
        name = "com_google_cloud_go_websecurityscanner",
        importpath = "cloud.google.com/go/websecurityscanner",
        sum = "h1:udhvvDDRryM3nrITJk/eQe74D06KK2N3SF60/FH2njQ=",
        version = "v1.7.7",
    )
    go_repository(
        name = "com_google_cloud_go_workflows",
        importpath = "cloud.google.com/go/workflows",
        sum = "h1:FGF6QEl3rtOSIHPOMZofWRVy3KNx26jDdgoYzJZ6ZhY=",
        version = "v1.14.3",
    )
    go_repository(
        name = "dev_cel_expr",
        importpath = "cel.dev/expr",
        sum = "h1:56OvJKSH3hDGL0ml5uSxZmz3/3Pq4tJ+fb1unVLAFcY=",
        version = "v0.24.0",
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
        sum = "h1:F7q2tNlCaHY9nMKHR6XH9/qkp8FktLnIcy6jJNyOCQw=",
        version = "v1.36.0",
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
        sum = "h1:YH4g8lQroajqUwWbq/tr2QX1JFmEXaDLgG+ew9bLMWo=",
        version = "v0.63.0",
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
        sum = "h1:RbKq8BG0FI8OiXhBfcRtqqHcZcka+gU3cskNuf05R18=",
        version = "v0.63.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel",
        build_directives = [
            "gazelle:resolve go go.opentelemetry.io/otel/trace @io_opentelemetry_go_otel_trace//:go_default_library",
            "gazelle:resolve go go.opentelemetry.io/otel/metric @io_opentelemetry_go_otel_metric//:go_default_library",
            "gazelle:resolve go go.opentelemetry.io/otel/metric/embedded @io_opentelemetry_go_otel_metric//embedded:go_default_library",
        ],
        importpath = "go.opentelemetry.io/otel",
        sum = "h1:RkfdswUDRimDg0m2Az18RKOsnI8UDzppJAtj01/Ymk8=",
        version = "v1.38.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_metric",
        build_directives = [
            "gazelle:resolve go go.opentelemetry.io/otel/attribute @io_opentelemetry_go_otel//attribute:go_default_library",
        ],
        importpath = "go.opentelemetry.io/otel/metric",
        sum = "h1:Kl6lzIYGAh5M159u9NgiRkmoMKjvbsKtYRwgfrA6WpA=",
        version = "v1.38.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_sdk",
        importpath = "go.opentelemetry.io/otel/sdk",
        sum = "h1:l48sr5YbNf2hpCUj/FoGhW9yDkl+Ma+LrVl8qaM5b+E=",
        version = "v1.38.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_sdk_metric",
        importpath = "go.opentelemetry.io/otel/sdk/metric",
        sum = "h1:aSH66iL0aZqo//xXzQLYozmWrXxyFkBJ6qT5wthqPoM=",
        version = "v1.38.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_trace",
        build_directives = [
            "gazelle:resolve go go.opentelemetry.io/otel/attribute @io_opentelemetry_go_otel//attribute:go_default_library",
        ],
        importpath = "go.opentelemetry.io/otel/trace",
        sum = "h1:Fxk5bKrDZJUH+AMyyIXGcFAPah0oRcT+LuNtJrmcNLE=",
        version = "v1.38.0",
    )
    go_repository(
        name = "org_golang_google_api",
        build_directives = [
            "gazelle:resolve go go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp @io_opentelemetry_go_contrib_instrumentation_net_http_otelhttp//:go_default_library",
            "gazelle:resolve go go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc @io_opentelemetry_go_contrib_instrumentation_google_golang_org_grpc_otelgrpc//:go_default_library",
        ],
        importpath = "google.golang.org/api",
        sum = "h1:jl3XrGj7lRjnlUvZAbAdhINTLbsg5dbjmR90+pTQvt4=",
        version = "v0.254.0",
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
        sum = "h1:fQ3VVQ11pb84nu0o/8wD6oZq13Q6+HK30P+9GSRlrqk=",
        version = "v0.0.0-20251029180050-ab9386a59fda",
    )
    go_repository(
        name = "org_golang_google_genproto_googleapis_api",
        importpath = "google.golang.org/genproto/googleapis/api",
        sum = "h1:+2XxjfsAu6vqFxwGBRcHiMaDCuZiqXGDUDVWVtrFAnE=",
        version = "v0.0.0-20251029180050-ab9386a59fda",
    )
    go_repository(
        name = "org_golang_google_genproto_googleapis_bytestream",
        importpath = "google.golang.org/genproto/googleapis/bytestream",
        sum = "h1:13Hpt6ugz7+eg3ZH5Vf1HmBDd7LxspYkhNtHONkvDHE=",
        version = "v0.0.0-20251022142026-3a174f9686a8",
    )
    go_repository(
        name = "org_golang_google_genproto_googleapis_rpc",
        importpath = "google.golang.org/genproto/googleapis/rpc",
        sum = "h1:i/Q+bfisr7gq6feoJnS/DlpdwEL4ihp41fvRiM3Ork0=",
        version = "v0.0.0-20251029180050-ab9386a59fda",
    )
    go_repository(
        name = "org_golang_google_grpc",
        importpath = "google.golang.org/grpc",
        sum = "h1:UnVkv1+uMLYXoIz6o7chp59WfQUYA2ex/BXQ9rHZu7A=",
        version = "v1.76.0",
    )
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
        sum = "h1:dduJYIi3A3KOfdGOHX8AVZ/jGiyPa3IbBozJ5kNuE04=",
        version = "v0.43.0",
    )
    go_repository(
        name = "org_golang_x_mod",
        importpath = "golang.org/x/mod",
        sum = "h1:gQBtGhjxykdjY9YhZpSlZIsbnaE2+PgjfLWUQTnoZ1U=",
        version = "v0.28.0",
    )
    go_repository(
        name = "org_golang_x_net",
        importpath = "golang.org/x/net",
        sum = "h1:giFlY12I07fugqwPuWJi68oOnpfqFnJIJzaIIm2JVV4=",
        version = "v0.46.0",
    )
    go_repository(
        name = "org_golang_x_oauth2",
        importpath = "golang.org/x/oauth2",
        sum = "h1:jsCblLleRMDrxMN29H3z/k1KliIvpLgCkE6R8FXXNgY=",
        version = "v0.32.0",
    )
    go_repository(
        name = "org_golang_x_sync",
        importpath = "golang.org/x/sync",
        sum = "h1:l60nONMj9l5drqw6jlhIELNv9I0A4OFgRsG9k2oT9Ug=",
        version = "v0.17.0",
    )
    go_repository(
        name = "org_golang_x_sys",
        importpath = "golang.org/x/sys",
        sum = "h1:fdNQudmxPjkdUTPnLn5mdQv7Zwvbvpaxqs831goi9kQ=",
        version = "v0.37.0",
    )
    go_repository(
        name = "org_golang_x_term",
        importpath = "golang.org/x/term",
        sum = "h1:zMPR+aF8gfksFprF/Nc/rd1wRS1EI6nDBGyWAvDzx2Q=",
        version = "v0.36.0",
    )
    go_repository(
        name = "org_golang_x_text",
        importpath = "golang.org/x/text",
        sum = "h1:yznKA/E9zq54KzlzBEAWn1NXSQ8DIp/NYMy88xJjl4k=",
        version = "v0.30.0",
    )
    go_repository(
        name = "org_golang_x_time",
        importpath = "golang.org/x/time",
        sum = "h1:MRx4UaLrDotUKUdCIqzPC48t1Y9hANFKIRpNx+Te8PI=",
        version = "v0.14.0",
    )
    go_repository(
        name = "org_golang_x_tools",
        importpath = "golang.org/x/tools",
        sum = "h1:DVSRzp7FwePZW356yEAChSdNcQo6Nsp+fex1SUW09lE=",
        version = "v0.37.0",
    )
    go_repository(
        name = "org_gonum_v1_gonum",
        importpath = "gonum.org/v1/gonum",
        sum = "h1:5+ul4Swaf3ESvrOnidPp4GZbzf0mxVQpDCYUQE7OJfk=",
        version = "v0.16.0",
    )
    go_repository(
        name = "org_uber_go_atomic",
        importpath = "go.uber.org/atomic",
        sum = "h1:ECmE8Bn/WFTYwEW/bpKD3M8VtR/zQVbavAoalC1PYyE=",
        version = "v1.9.0",
    )
    go_repository(
        name = "org_uber_go_multierr",
        importpath = "go.uber.org/multierr",
        sum = "h1:7fIwc/ZtS0q++VgcfqFDxSBZVv/Xo49/SYnDFupUwlI=",
        version = "v1.9.0",
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
