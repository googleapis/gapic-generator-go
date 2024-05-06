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
        name = "co_honnef_go_tools",
        importpath = "honnef.co/go/tools",
        sum = "h1:/hemPrYIhOhy8zYrNj+069zDB68us2sMGsfkFJO0iZs=",
        version = "v0.0.0-20190523083050-ea95bdfd59fc",
    )

    go_repository(
        name = "com_github_bufbuild_protocompile",
        importpath = "github.com/bufbuild/protocompile",
        sum = "h1:+jW/wnLMLxaCEG8AX9lD0bQ5v9h1RUiMKOBOT5ll9dM=",
        version = "v0.10.0",
    )
    go_repository(
        name = "com_github_burntsushi_toml",
        importpath = "github.com/BurntSushi/toml",
        sum = "h1:WXkYYl6Yr3qBf1K79EBnL4mak0OimBfB0XUf9Vl28OQ=",
        version = "v0.3.1",
    )
    go_repository(
        name = "com_github_census_instrumentation_opencensus_proto",
        importpath = "github.com/census-instrumentation/opencensus-proto",
        sum = "h1:iKLQ0xPNFxR/2hzXZMrBo8f1j86j5WHzznCCQxV/b8g=",
        version = "v0.4.1",
    )
    go_repository(
        name = "com_github_cespare_xxhash_v2",
        importpath = "github.com/cespare/xxhash/v2",
        sum = "h1:DC2CZ1Ep5Y4k3ZQ899DldepgrayRUGE6BBZ/cd9Cj44=",
        version = "v2.2.0",
    )
    go_repository(
        name = "com_github_client9_misspell",
        importpath = "github.com/client9/misspell",
        sum = "h1:ta993UF76GwbvJcIo3Y68y/M3WxlpEHPWIGDkJYwzJI=",
        version = "v0.3.4",
    )
    go_repository(
        name = "com_github_cncf_udpa_go",
        importpath = "github.com/cncf/udpa/go",
        sum = "h1:WBZRG4aNOuI15bLRrCgN8fCq8E5Xuty6jGbmSNEvSsU=",
        version = "v0.0.0-20191209042840-269d4d468f6f",
    )
    go_repository(
        name = "com_github_cncf_xds_go",
        importpath = "github.com/cncf/xds/go",
        sum = "h1:jQCWAUqqlij9Pgj2i/PB79y4KOPYVyFYdROxgaCwdTQ=",
        version = "v0.0.0-20231128003011-0fa0005c9caa",
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
        sum = "h1:4X+VP1GHd1Mhj6IB5mMeGbLCleqxjletLK6K0rbxyZI=",
        version = "v0.12.0",
    )
    go_repository(
        name = "com_github_envoyproxy_protoc_gen_validate",
        importpath = "github.com/envoyproxy/protoc-gen-validate",
        sum = "h1:gVPz/FMfvh57HdSJQyvBtF00j8JU4zdyUgIUNhlgg0A=",
        version = "v1.0.4",
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
        sum = "h1:8JEhPFa5W2WU7YfeZzPNqzMP6Lwt7L2715Ggo0nosvA=",
        version = "v1.7.0",
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
        name = "com_github_go_logr_logr",
        importpath = "github.com/go-logr/logr",
        sum = "h1:pKouT5E8xu9zeFC39JXRDukb6JFQPXM5p5I91188VAQ=",
        version = "v1.4.1",
    )
    go_repository(
        name = "com_github_go_logr_stdr",
        importpath = "github.com/go-logr/stdr",
        sum = "h1:hSWxHoqTgW2S2qGc0LTAI563KZ5YKYRhT3MFKZMbjag=",
        version = "v1.2.2",
    )
    go_repository(
        name = "com_github_golang_glog",
        importpath = "github.com/golang/glog",
        sum = "h1:uCdmnmatrKCgMBlM4rMuJZWOkPDqdbZPnrMXDY4gI68=",
        version = "v1.2.0",
    )
    go_repository(
        name = "com_github_golang_groupcache",
        importpath = "github.com/golang/groupcache",
        sum = "h1:oI5xCqsCo564l8iNU+DwB5epxmsaqB+rhGL0m5jtYqE=",
        version = "v0.0.0-20210331224755-41bb18bfe9da",
    )
    go_repository(
        name = "com_github_golang_mock",
        importpath = "github.com/golang/mock",
        sum = "h1:G5FRp8JnTd7RQH5kemVNlMeyXQAztQ3mOWV95KxsXH8=",
        version = "v1.1.1",
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
        sum = "h1:ofyhxvXcZhMsU5ulbFiLKl/XBFqE1GSq7atu8tAmTRI=",
        version = "v0.6.0",
    )
    go_repository(
        name = "com_github_google_go_pkcs11",
        importpath = "github.com/google/go-pkcs11",
        sum = "h1:OF1IPgv+F4NmqmJ98KTjdN97Vs1JxDPB3vbmYzV2dpk=",
        version = "v0.2.1-0.20230907215043-c6f79328ddf9",
    )
    go_repository(
        name = "com_github_google_martian_v3",
        importpath = "github.com/google/martian/v3",
        sum = "h1:IqNFLAmvJOgVlpdEBiQbDc2EwKW77amAycfTuWKdfvw=",
        version = "v3.3.2",
    )
    go_repository(
        name = "com_github_google_s2a_go",
        importpath = "github.com/google/s2a-go",
        sum = "h1:60BLSyTrOV4/haCDW4zb1guZItoSq8foHCXrAnjBo/o=",
        version = "v0.1.7",
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
        sum = "h1:Vie5ybvEvT75RniqhfFxPRy3Bf7vr3h0cechB90XaQs=",
        version = "v0.3.2",
    )
    go_repository(
        name = "com_github_googleapis_gapic_showcase",
        importpath = "github.com/googleapis/gapic-showcase",
        sum = "h1:uRNh7SNmfFh18XfmU6p1caJVAaP2ZVIrLMCMSaKefs4=",
        version = "v0.35.0",
    )
    go_repository(
        name = "com_github_googleapis_gax_go_v2",
        build_directives = [
            "gazelle:resolve go google.golang.org/genproto/googleapis/rpc/errdetails @org_golang_google_genproto_googleapis_rpc//errdetails",
            "gazelle:resolve proto go google/rpc/code.proto @com_google_googleapis//google/rpc:code_go_proto",
            "gazelle:resolve proto proto google/rpc/code.proto @com_google_googleapis//google/rpc:code_proto",
        ],
        importpath = "github.com/googleapis/gax-go/v2",
        sum = "h1:9gWcmF85Wvq4ryPFvGFaOgPIs1AQX0d0bcbGw4Z96qg=",
        version = "v2.12.4",
    )
    go_repository(
        name = "com_github_googleapis_grpc_fallback_go",
        importpath = "github.com/googleapis/grpc-fallback-go",
        sum = "h1:tEDqZnKGKQpYrmEuu3VVBTw3pijHJsKz/Lu2U0L9AV0=",
        version = "v0.1.4",
    )
    go_repository(
        name = "com_github_gorilla_mux",
        importpath = "github.com/gorilla/mux",
        sum = "h1:TuBL49tXwgrFYWhqrNgrUNEY92u81SPhu7sTdzQEiWY=",
        version = "v1.8.1",
    )
    go_repository(
        name = "com_github_hashicorp_hcl",
        importpath = "github.com/hashicorp/hcl",
        sum = "h1:0Anlzjpi4vEasTeNFn2mLJgTSwt0+6sfsiTG8qcWGx4=",
        version = "v1.0.0",
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
        name = "com_github_magiconair_properties",
        importpath = "github.com/magiconair/properties",
        sum = "h1:IeQXZAiQcpL9mgcAe1Nu6cX9LLw6ExEHKjN0VQdvPDY=",
        version = "v1.8.7",
    )
    go_repository(
        name = "com_github_mitchellh_mapstructure",
        importpath = "github.com/mitchellh/mapstructure",
        sum = "h1:jeMsZIYE/09sWLaz43PL7Gy6RuMjD2eJVyuac5Z2hdY=",
        version = "v1.5.0",
    )
    go_repository(
        name = "com_github_pelletier_go_toml_v2",
        importpath = "github.com/pelletier/go-toml/v2",
        sum = "h1:FnwAJ4oYMvbT/34k9zzHuZNrhlz48GB3/s6at6/MHO4=",
        version = "v2.1.0",
    )
    go_repository(
        name = "com_github_pmezard_go_difflib",
        importpath = "github.com/pmezard/go-difflib",
        sum = "h1:4DBwDE0NGyQoBHbLQYPwSUPoCMWR5BEzIk/f1lZbAQM=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_prometheus_client_model",
        importpath = "github.com/prometheus/client_model",
        sum = "h1:gQz4mCbXsO+nc9n1hCxHcGA3Zx3Eo+UHZoInFGUIXNM=",
        version = "v0.0.0-20190812154241-14fe0d1b01d4",
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
        sum = "h1:HApY1R9zGo4DBgr7dqsTH/JJxLTTsOt7u6keLGt6kNQ=",
        version = "v0.4.0",
    )
    go_repository(
        name = "com_github_sagikazarmark_slog_shim",
        importpath = "github.com/sagikazarmark/slog-shim",
        sum = "h1:diDBnUNK9N/354PgrxMywXnAwEr1QZcOr6gto+ugjYE=",
        version = "v0.1.0",
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
        sum = "h1:WJQKhtpdm3v2IzqG8VMqrr6Rf3UYpEF239Jy9wNepM8=",
        version = "v1.11.0",
    )
    go_repository(
        name = "com_github_spf13_cast",
        importpath = "github.com/spf13/cast",
        sum = "h1:GEiTHELF+vaR5dhz3VqZfFSzZjYbgeKDpBxQVS4GYJ0=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_github_spf13_cobra",
        importpath = "github.com/spf13/cobra",
        sum = "h1:7aJaZx1B85qltLMc546zn58BxxfZdR/W22ej9CFoEf0=",
        version = "v1.8.0",
    )
    go_repository(
        name = "com_github_spf13_pflag",
        importpath = "github.com/spf13/pflag",
        sum = "h1:iy+VFUOCP1a+8yFto/drg2CJ5u0yRoB7fZw3DKv/JXA=",
        version = "v1.0.5",
    )
    go_repository(
        name = "com_github_spf13_viper",
        importpath = "github.com/spf13/viper",
        sum = "h1:LUXCnvUvSM6FXAsj6nnfc8Q2tp1dIgUfY9Kc8GsSOiQ=",
        version = "v1.18.2",
    )
    go_repository(
        name = "com_github_stretchr_objx",
        importpath = "github.com/stretchr/objx",
        sum = "h1:1zr/of2m5FGMsad5YfcqgdqdWrIhu+EBEJRhR1U7z/c=",
        version = "v0.5.0",
    )
    go_repository(
        name = "com_github_stretchr_testify",
        importpath = "github.com/stretchr/testify",
        sum = "h1:HtqpIVDClZ4nwg75+f6Lvsy/wHu+3BoSGCbBAcpTsTg=",
        version = "v1.9.0",
    )
    go_repository(
        name = "com_github_subosito_gotenv",
        importpath = "github.com/subosito/gotenv",
        sum = "h1:9NlTDc1FTs4qu0DDq7AEtTPNw6SVm7uBMsUCUjABIf8=",
        version = "v1.6.0",
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
        sum = "h1:oYrL81N608MLZhma3ruL8qTM4xcpYECGut8KSxRY59g=",
        version = "v0.0.0-20191026162114-a0c2df6c8f82",
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
        sum = "h1:ZaGT6LiG7dBzi6zNOvVZwacaXlmf3lRqnC4DQzqyRQw=",
        version = "v0.112.2",
    )
    go_repository(
        name = "com_google_cloud_go_accessapproval",
        importpath = "cloud.google.com/go/accessapproval",
        sum = "h1:fMbP4cJX/926h+kwGtABmcG83PXsjkB+q7nSBzZpJoo=",
        version = "v1.7.6",
    )
    go_repository(
        name = "com_google_cloud_go_accesscontextmanager",
        importpath = "cloud.google.com/go/accesscontextmanager",
        sum = "h1:NipmPd3BCzwa/mr40SK8pWRkbzv9Th5Azhi4dBYazlM=",
        version = "v1.8.6",
    )
    go_repository(
        name = "com_google_cloud_go_aiplatform",
        importpath = "cloud.google.com/go/aiplatform",
        sum = "h1:YWeqD4BjYwrmY4fa+isGcw0P81lJ3dKVxbWxdBchoiU=",
        version = "v1.67.0",
    )
    go_repository(
        name = "com_google_cloud_go_analytics",
        importpath = "cloud.google.com/go/analytics",
        sum = "h1:UH/PWBcPxHKolWxaS3hO+aj+wDTuq7XKvdtpqR3lyOI=",
        version = "v0.23.1",
    )
    go_repository(
        name = "com_google_cloud_go_apigateway",
        importpath = "cloud.google.com/go/apigateway",
        sum = "h1:60GMRN1JFwq9MldvEVMdR3gDJ0vI0C/BwgkImG6bx/M=",
        version = "v1.6.6",
    )
    go_repository(
        name = "com_google_cloud_go_apigeeconnect",
        importpath = "cloud.google.com/go/apigeeconnect",
        sum = "h1:ObsKNGtdu0ckkCSUpCN5fVAVg+DmzFg89JlqsW4OAWk=",
        version = "v1.6.6",
    )
    go_repository(
        name = "com_google_cloud_go_apigeeregistry",
        importpath = "cloud.google.com/go/apigeeregistry",
        sum = "h1:l8VFHdNMC+9Q4EHKye2eOZBu5IwddXF6ufAXI7D+PB8=",
        version = "v0.8.4",
    )
    go_repository(
        name = "com_google_cloud_go_appengine",
        importpath = "cloud.google.com/go/appengine",
        sum = "h1:cM+Lj0A0nCtujM9lqRId5L8hz7bHRfu3q3KdKtpn+38=",
        version = "v1.8.6",
    )
    go_repository(
        name = "com_google_cloud_go_area120",
        importpath = "cloud.google.com/go/area120",
        sum = "h1:7QJ4ZzqLOwh0pHD3UAa+VwiyWmDI7bdmkYKVkte8/wg=",
        version = "v0.8.6",
    )
    go_repository(
        name = "com_google_cloud_go_artifactregistry",
        importpath = "cloud.google.com/go/artifactregistry",
        sum = "h1:icIyRzJ1Ag6EOafuDuFFJ/AdStcOFRVfSGURn27/7Pk=",
        version = "v1.14.8",
    )
    go_repository(
        name = "com_google_cloud_go_asset",
        importpath = "cloud.google.com/go/asset",
        sum = "h1:D1fzO6/fdMkGA7Lxrv/in2ahVm59DAe1EMMaDt33v9k=",
        version = "v1.19.0",
    )
    go_repository(
        name = "com_google_cloud_go_assuredworkloads",
        importpath = "cloud.google.com/go/assuredworkloads",
        sum = "h1:3NlUes0xLN2kcSU24qQADFYsOaetCPg0HSA302AyV5s=",
        version = "v1.11.6",
    )
    go_repository(
        name = "com_google_cloud_go_auth",
        importpath = "cloud.google.com/go/auth",
        sum = "h1:PRyzEpGfx/Z9e8+lHsbkoUVXD0gnu4MNmm7Gp8TQNIs=",
        version = "v0.3.0",
    )
    go_repository(
        name = "com_google_cloud_go_auth_oauth2adapt",
        importpath = "cloud.google.com/go/auth/oauth2adapt",
        sum = "h1:+TTV8aXpjeChS9M+aTtN/TjdQnzJvmzKFt//oWu7HX4=",
        version = "v0.2.2",
    )
    go_repository(
        name = "com_google_cloud_go_automl",
        importpath = "cloud.google.com/go/automl",
        sum = "h1:NHBO5cjo2IgwaJ5qlez/iA35XI1db87PPlOB0Kjt5RM=",
        version = "v1.13.6",
    )
    go_repository(
        name = "com_google_cloud_go_baremetalsolution",
        importpath = "cloud.google.com/go/baremetalsolution",
        sum = "h1:jCR4rnVsG6ocK6ngFr2Z6ugKZfTENmMZkiV6Ma2tEeE=",
        version = "v1.2.5",
    )
    go_repository(
        name = "com_google_cloud_go_batch",
        importpath = "cloud.google.com/go/batch",
        sum = "h1:54IA0n5e9PyYq9baH7ULYjyvubxHFT8Ys6LR2Om/TuY=",
        version = "v1.8.4",
    )
    go_repository(
        name = "com_google_cloud_go_beyondcorp",
        importpath = "cloud.google.com/go/beyondcorp",
        sum = "h1:fnil8viEdcAJJiwiJPCT2fl3Grx3XFwXxTq7n9g/8QM=",
        version = "v1.0.5",
    )
    go_repository(
        name = "com_google_cloud_go_bigquery",
        importpath = "cloud.google.com/go/bigquery",
        sum = "h1:w2Goy9n6gh91LVi6B2Sc+HpBl8WbWhIyzdvVvrAuEIw=",
        version = "v1.61.0",
    )
    go_repository(
        name = "com_google_cloud_go_billing",
        importpath = "cloud.google.com/go/billing",
        sum = "h1:XcYB8aKj97d4/0kh+LQgrxPxOo/0jgPNp5Q1nyzCyvs=",
        version = "v1.18.4",
    )
    go_repository(
        name = "com_google_cloud_go_binaryauthorization",
        importpath = "cloud.google.com/go/binaryauthorization",
        sum = "h1:XiAdW5HAWtk9WEjEA5MXYiRJ28Tjp1xGPPAMRFI5bnc=",
        version = "v1.8.2",
    )
    go_repository(
        name = "com_google_cloud_go_certificatemanager",
        importpath = "cloud.google.com/go/certificatemanager",
        sum = "h1:oc15T+leZ2Wm5oocvGTbDXwonka0chpJTEiHIVsBiEA=",
        version = "v1.8.0",
    )
    go_repository(
        name = "com_google_cloud_go_channel",
        importpath = "cloud.google.com/go/channel",
        sum = "h1:rBnTls9G5nC/jOrb9V/tZFHFXt6kBMNlIQKg6DjAlRM=",
        version = "v1.17.6",
    )
    go_repository(
        name = "com_google_cloud_go_cloudbuild",
        importpath = "cloud.google.com/go/cloudbuild",
        sum = "h1:66hY1gXV2cdn4gquy5TieaJwaZmRzrQ5cK++pVgnCkQ=",
        version = "v1.16.0",
    )
    go_repository(
        name = "com_google_cloud_go_clouddms",
        importpath = "cloud.google.com/go/clouddms",
        sum = "h1:t1nc49kRzEU2vrDxQQIEc5rZ4Hr187YrdOZPMAgMwMI=",
        version = "v1.7.5",
    )
    go_repository(
        name = "com_google_cloud_go_cloudtasks",
        importpath = "cloud.google.com/go/cloudtasks",
        sum = "h1:Ev+poxwb7pudBhiF0ObwAWT7Dh9BZAcsvAfFTWg0MPc=",
        version = "v1.12.7",
    )
    go_repository(
        name = "com_google_cloud_go_compute",
        importpath = "cloud.google.com/go/compute",
        sum = "h1:ZRpHJedLtTpKgr3RV1Fx23NuaAEN1Zfx9hw1u4aJdjU=",
        version = "v1.25.1",
    )
    go_repository(
        name = "com_google_cloud_go_compute_metadata",
        importpath = "cloud.google.com/go/compute/metadata",
        sum = "h1:Tz+eQXMEqDIKRsmY3cHTL6FVaynIjX2QxYC4trgAKZc=",
        version = "v0.3.0",
    )
    go_repository(
        name = "com_google_cloud_go_contactcenterinsights",
        importpath = "cloud.google.com/go/contactcenterinsights",
        sum = "h1:sCDKUmDj9Tfd6Qj7x4XbwC43oYzEBwSDLC1tReQWS/Y=",
        version = "v1.13.1",
    )
    go_repository(
        name = "com_google_cloud_go_container",
        importpath = "cloud.google.com/go/container",
        sum = "h1:y5gmgrMMhTrLnQQdMCw0t/Yis9Ps7jvAG4JYcRWxR8g=",
        version = "v1.35.0",
    )
    go_repository(
        name = "com_google_cloud_go_containeranalysis",
        importpath = "cloud.google.com/go/containeranalysis",
        sum = "h1:yzohQ0HDoZq2TtCJkbUBsJs9RIR5WbKZlHrD7ilp2yg=",
        version = "v0.11.5",
    )
    go_repository(
        name = "com_google_cloud_go_datacatalog",
        importpath = "cloud.google.com/go/datacatalog",
        sum = "h1:BGDsEjqpAo0Ka+b9yDLXnE5k+jU3lXGMh//NsEeDMIg=",
        version = "v1.20.0",
    )
    go_repository(
        name = "com_google_cloud_go_dataflow",
        importpath = "cloud.google.com/go/dataflow",
        sum = "h1:GuZJgkOL64cYySwYEYqQkggdxwoZTk8cvekQW0t3KRM=",
        version = "v0.9.6",
    )
    go_repository(
        name = "com_google_cloud_go_dataform",
        importpath = "cloud.google.com/go/dataform",
        sum = "h1:0EzWf+c2R5V/ujZBb22H/EL5wpzD/bDFYPA2f3czB1g=",
        version = "v0.9.3",
    )
    go_repository(
        name = "com_google_cloud_go_datafusion",
        importpath = "cloud.google.com/go/datafusion",
        sum = "h1:zSmMj/qZ0Yk+q/v5Wg40FTJ0WYPCtanYYekRt7cSKJ0=",
        version = "v1.7.6",
    )
    go_repository(
        name = "com_google_cloud_go_datalabeling",
        importpath = "cloud.google.com/go/datalabeling",
        sum = "h1:2zz44bPbDMHsPanQ89SybqhHDQBH1EZk8icGotyrvSU=",
        version = "v0.8.6",
    )
    go_repository(
        name = "com_google_cloud_go_dataplex",
        importpath = "cloud.google.com/go/dataplex",
        sum = "h1:Ob8NPT1UcB4kDaDx7/UdsRfZ8xUvUggZshXUlGWDahk=",
        version = "v1.15.0",
    )
    go_repository(
        name = "com_google_cloud_go_dataproc_v2",
        importpath = "cloud.google.com/go/dataproc/v2",
        sum = "h1:+cM8p/R6FdTuQYlriJOSUCvAZfMDgBKf0/ph9bMIjaY=",
        version = "v2.4.1",
    )
    go_repository(
        name = "com_google_cloud_go_dataqna",
        importpath = "cloud.google.com/go/dataqna",
        sum = "h1:FI/1q7VnANchQR9ug+nzujfiusLMfC3BHT7/fHi+BVU=",
        version = "v0.8.6",
    )
    go_repository(
        name = "com_google_cloud_go_datastore",
        importpath = "cloud.google.com/go/datastore",
        sum = "h1:0P9WcsQeTWjuD1H14JIY7XQscIPQ4Laje8ti96IC5vg=",
        version = "v1.15.0",
    )
    go_repository(
        name = "com_google_cloud_go_datastream",
        importpath = "cloud.google.com/go/datastream",
        sum = "h1:nHdOKbFmKJ4tPjGtNNIO0//G7QAht6eHTUnREWPn2yI=",
        version = "v1.10.5",
    )
    go_repository(
        name = "com_google_cloud_go_deploy",
        importpath = "cloud.google.com/go/deploy",
        sum = "h1:UxcxzjwxGPkT7RBdMmoc5a7QxHQVdpZllD6el2VC3JA=",
        version = "v1.17.2",
    )
    go_repository(
        name = "com_google_cloud_go_dialogflow",
        importpath = "cloud.google.com/go/dialogflow",
        sum = "h1:B8Y5j4/QsDirX136OoPm62kG3y5gd8rzBpHSR/FW9vI=",
        version = "v1.52.0",
    )
    go_repository(
        name = "com_google_cloud_go_dlp",
        importpath = "cloud.google.com/go/dlp",
        sum = "h1:dTsEN6r1BoplUACz7teOmE6lRG1swREiwXkfkY7bi6c=",
        version = "v1.12.1",
    )
    go_repository(
        name = "com_google_cloud_go_documentai",
        importpath = "cloud.google.com/go/documentai",
        sum = "h1:pnUL6P5y6gfDRBi4JxR0FZl3ckB7TC7fKdV0ZxhIjpM=",
        version = "v1.27.0",
    )
    go_repository(
        name = "com_google_cloud_go_domains",
        importpath = "cloud.google.com/go/domains",
        sum = "h1:NHqZk4XzHFlmXM3LMGwDVET4NKr60W2jaNCRGYod5Ic=",
        version = "v0.9.6",
    )
    go_repository(
        name = "com_google_cloud_go_edgecontainer",
        importpath = "cloud.google.com/go/edgecontainer",
        sum = "h1:a++vBi1J00NP1ifVP5mV/3j1/EJKWPj0h6NfUPLfuCQ=",
        version = "v1.2.0",
    )
    go_repository(
        name = "com_google_cloud_go_errorreporting",
        importpath = "cloud.google.com/go/errorreporting",
        sum = "h1:kj1XEWMu8P0qlLhm3FwcaFsUvXChV/OraZwA70trRR0=",
        version = "v0.3.0",
    )
    go_repository(
        name = "com_google_cloud_go_essentialcontacts",
        importpath = "cloud.google.com/go/essentialcontacts",
        sum = "h1:FDdJGJEXK4RxvT6gdRBqGaCQVpi96RRB7MTyRHUcb20=",
        version = "v1.6.7",
    )
    go_repository(
        name = "com_google_cloud_go_eventarc",
        importpath = "cloud.google.com/go/eventarc",
        sum = "h1:JMUiLYzxkxr7BqnCPkyJ6Ycgrs6YQlZT44H0rHg7jQY=",
        version = "v1.13.5",
    )
    go_repository(
        name = "com_google_cloud_go_filestore",
        importpath = "cloud.google.com/go/filestore",
        sum = "h1:BpaB7bxICPUTntAV+yVUK9bxAUOv7uHRSEibSKMBJVA=",
        version = "v1.8.2",
    )
    go_repository(
        name = "com_google_cloud_go_firestore",
        importpath = "cloud.google.com/go/firestore",
        sum = "h1:/k8ppuWOtNuDHt2tsRV42yI21uaGnKDEQnRFeBpbFF8=",
        version = "v1.15.0",
    )
    go_repository(
        name = "com_google_cloud_go_functions",
        importpath = "cloud.google.com/go/functions",
        sum = "h1:0kcko/2AKwm4USnWcGs/W/k++PAYPA3dYaQw1y5Xg3M=",
        version = "v1.16.1",
    )
    go_repository(
        name = "com_google_cloud_go_gkebackup",
        importpath = "cloud.google.com/go/gkebackup",
        sum = "h1:SATJwsF8PjErP7GwHE+xK8gJ7f7hULuqtazV19ylPgg=",
        version = "v1.4.0",
    )
    go_repository(
        name = "com_google_cloud_go_gkeconnect",
        importpath = "cloud.google.com/go/gkeconnect",
        sum = "h1:7X9P6lGkOF/nJRYZBQBG2XPhunqWbNMacy9AXN7qUcU=",
        version = "v0.8.6",
    )
    go_repository(
        name = "com_google_cloud_go_gkehub",
        importpath = "cloud.google.com/go/gkehub",
        sum = "h1:kKreFf+097KfW+Tz/SqZKeXs/eFOjs1NDrsVjRPI18o=",
        version = "v0.14.6",
    )
    go_repository(
        name = "com_google_cloud_go_gkemulticloud",
        importpath = "cloud.google.com/go/gkemulticloud",
        sum = "h1:CFBoDcQi9zLOkzM6xqmRzljZhF4A6A47QaQ0WtNd+DA=",
        version = "v1.1.2",
    )
    go_repository(
        name = "com_google_cloud_go_gsuiteaddons",
        importpath = "cloud.google.com/go/gsuiteaddons",
        sum = "h1:q3x2NE0je/tSVL66MAht5YVbGGHjTV9BxFD2lyDQ0dU=",
        version = "v1.6.6",
    )
    go_repository(
        name = "com_google_cloud_go_iam",
        importpath = "cloud.google.com/go/iam",
        sum = "h1:r7umDwhj+BQyz0ScZMp4QrGXjSTI3ZINnpgU2nlB/K0=",
        version = "v1.1.8",
    )
    go_repository(
        name = "com_google_cloud_go_iap",
        importpath = "cloud.google.com/go/iap",
        sum = "h1:FrLAtgXzWPwe8rNp7AD+2Lgg4LqyhgXvEdiGK+jtd9g=",
        version = "v1.9.5",
    )
    go_repository(
        name = "com_google_cloud_go_ids",
        importpath = "cloud.google.com/go/ids",
        sum = "h1:tNc3NpIp2LUmFJxP2CBlzYw0FTnd68r73mIzg8UlM3Q=",
        version = "v1.4.6",
    )
    go_repository(
        name = "com_google_cloud_go_iot",
        importpath = "cloud.google.com/go/iot",
        sum = "h1:nRV/e1e3lEjsVoD5mW99JERwL8MKohyQyY63+lfBMA4=",
        version = "v1.7.6",
    )
    go_repository(
        name = "com_google_cloud_go_kms",
        importpath = "cloud.google.com/go/kms",
        sum = "h1:szIeDCowID8th2i8XE4uRev5PMxQFqW+JjwYxL9h6xs=",
        version = "v1.15.8",
    )
    go_repository(
        name = "com_google_cloud_go_language",
        importpath = "cloud.google.com/go/language",
        sum = "h1:srkreCxnVa5+a5PXUri/K+VWxG50wGvz5+PEYjgaENQ=",
        version = "v1.12.4",
    )
    go_repository(
        name = "com_google_cloud_go_lifesciences",
        importpath = "cloud.google.com/go/lifesciences",
        sum = "h1:8w3edjRiSN6GCxT0uJoXr6Zo2XNYD+6TxPZa7uIIOaU=",
        version = "v0.9.6",
    )
    go_repository(
        name = "com_google_cloud_go_logging",
        importpath = "cloud.google.com/go/logging",
        sum = "h1:iEIOXFO9EmSiTjDmfpbRjOxECO7R8C7b8IXUGOj7xZw=",
        version = "v1.9.0",
    )
    go_repository(
        name = "com_google_cloud_go_longrunning",
        importpath = "cloud.google.com/go/longrunning",
        sum = "h1:xAe8+0YaWoCKr9t1+aWe+OeQgN/iJK1fEgZSXmjuEaE=",
        version = "v0.5.6",
    )
    go_repository(
        name = "com_google_cloud_go_managedidentities",
        importpath = "cloud.google.com/go/managedidentities",
        sum = "h1:7+hGPQSojhnYNZCg3fG2mQIF7XMfvNpCpi2Zg5/Qx1g=",
        version = "v1.6.6",
    )
    go_repository(
        name = "com_google_cloud_go_maps",
        importpath = "cloud.google.com/go/maps",
        sum = "h1:r8I8GHUnZzaMn06ZDE1dtVTBQMep3nXIho/13uEQd80=",
        version = "v1.7.2",
    )
    go_repository(
        name = "com_google_cloud_go_mediatranslation",
        importpath = "cloud.google.com/go/mediatranslation",
        sum = "h1:EVW0wCQ7asoMjw8fm8oUe6pQWBaVQth/iquk7ysidy0=",
        version = "v0.8.6",
    )
    go_repository(
        name = "com_google_cloud_go_memcache",
        importpath = "cloud.google.com/go/memcache",
        sum = "h1:rqDPCIUfVBvv7ojOGx5PRkPgWeWSKpOht2w6plaxklY=",
        version = "v1.10.6",
    )
    go_repository(
        name = "com_google_cloud_go_metastore",
        importpath = "cloud.google.com/go/metastore",
        sum = "h1:K7gyYoqPvQgCc82tiB0CQkXOpg8AZxJtRGMVdN5B83U=",
        version = "v1.13.5",
    )
    go_repository(
        name = "com_google_cloud_go_monitoring",
        importpath = "cloud.google.com/go/monitoring",
        sum = "h1:nIQdZdf1/9M0cEmXSlLB2jMq/k3CRh9p3oUzS06VDG8=",
        version = "v1.18.2",
    )
    go_repository(
        name = "com_google_cloud_go_networkconnectivity",
        importpath = "cloud.google.com/go/networkconnectivity",
        sum = "h1:t67aEKwmO+SXvQC5ncOjm3vTwnsbO/mTzlCWdK0nwqs=",
        version = "v1.14.5",
    )
    go_repository(
        name = "com_google_cloud_go_networkmanagement",
        importpath = "cloud.google.com/go/networkmanagement",
        sum = "h1:eRgeqeZFVu/WeECx0ChoDAwz/KVhIAWTFS53rGahzYI=",
        version = "v1.13.1",
    )
    go_repository(
        name = "com_google_cloud_go_networksecurity",
        importpath = "cloud.google.com/go/networksecurity",
        sum = "h1:3ggPKshcFs1oRh5qI+Gq1s2CIU9BL99MKtYSBG4Z8s0=",
        version = "v0.9.6",
    )
    go_repository(
        name = "com_google_cloud_go_notebooks",
        importpath = "cloud.google.com/go/notebooks",
        sum = "h1:A9jxIdxEccgL9iJLqvU4j5HT3/13YluoF2IbiV+hAN4=",
        version = "v1.11.4",
    )
    go_repository(
        name = "com_google_cloud_go_optimization",
        importpath = "cloud.google.com/go/optimization",
        sum = "h1:T/j8xyIkmHGjU6kxeUjK3UTqiXlbvpZQ2A+F5vnH21Y=",
        version = "v1.6.4",
    )
    go_repository(
        name = "com_google_cloud_go_orchestration",
        importpath = "cloud.google.com/go/orchestration",
        sum = "h1:i5iSxsu1Cx1itTQEEY/YvsAo1OO8gosGGXhnOjBjgJA=",
        version = "v1.9.1",
    )
    go_repository(
        name = "com_google_cloud_go_orgpolicy",
        importpath = "cloud.google.com/go/orgpolicy",
        sum = "h1:x9GttuUZXXeKcJgHSGxYoPn2hOJhhuaN5YYJKfAfmLo=",
        version = "v1.12.2",
    )
    go_repository(
        name = "com_google_cloud_go_osconfig",
        importpath = "cloud.google.com/go/osconfig",
        sum = "h1:wIOhgzklE0hHZsho02rRVXYBHSfsAwYZYIaxFaUBIjs=",
        version = "v1.12.6",
    )
    go_repository(
        name = "com_google_cloud_go_oslogin",
        importpath = "cloud.google.com/go/oslogin",
        sum = "h1:v71OrrkKyqr5Mfnt345GqSOURzByv08qfrtvfhOVcnc=",
        version = "v1.13.2",
    )
    go_repository(
        name = "com_google_cloud_go_phishingprotection",
        importpath = "cloud.google.com/go/phishingprotection",
        sum = "h1:DcAre1psFwJM/FBA/MkDj0H6uxZhACE5IW/xF9ssHDQ=",
        version = "v0.8.6",
    )
    go_repository(
        name = "com_google_cloud_go_policytroubleshooter",
        importpath = "cloud.google.com/go/policytroubleshooter",
        sum = "h1:wxBRfNoMy7rnoEeaTOHIEHCUEdUIQIwQGUqfBWH6cyQ=",
        version = "v1.10.4",
    )
    go_repository(
        name = "com_google_cloud_go_privatecatalog",
        importpath = "cloud.google.com/go/privatecatalog",
        sum = "h1:bcIABOUmpnzQip83OVv+Ju/NxXjUTRLUSP+FVLFG6kk=",
        version = "v0.9.6",
    )
    go_repository(
        name = "com_google_cloud_go_pubsub",
        importpath = "cloud.google.com/go/pubsub",
        sum = "h1:0uEEfaB1VIJzabPpwpZf44zWAKAme3zwKKxHk7vJQxQ=",
        version = "v1.37.0",
    )
    go_repository(
        name = "com_google_cloud_go_pubsublite",
        importpath = "cloud.google.com/go/pubsublite",
        sum = "h1:pX+idpWMIH30/K7c0epN6V703xpIcMXWRjKJsz0tYGY=",
        version = "v1.8.1",
    )
    go_repository(
        name = "com_google_cloud_go_recaptchaenterprise_v2",
        importpath = "cloud.google.com/go/recaptchaenterprise/v2",
        sum = "h1:nykUP2WD/914jui/IldiCOuoTn6T8ha1Ys6/N9sAqJY=",
        version = "v2.12.0",
    )
    go_repository(
        name = "com_google_cloud_go_recommendationengine",
        importpath = "cloud.google.com/go/recommendationengine",
        sum = "h1:m0eQtYCToxMSbDKOnpJ2YGdQhyjOPffg4Y8lM2RWzao=",
        version = "v0.8.6",
    )
    go_repository(
        name = "com_google_cloud_go_recommender",
        importpath = "cloud.google.com/go/recommender",
        sum = "h1:3M6lD39/GlOMYOikeF5wflSa4EP5pGFthoIASbyhIXE=",
        version = "v1.12.2",
    )
    go_repository(
        name = "com_google_cloud_go_redis",
        importpath = "cloud.google.com/go/redis",
        sum = "h1:zlGxeAsiwcPU+Cta76ALduhdBAVhuYpEjv59V5L/ves=",
        version = "v1.14.3",
    )
    go_repository(
        name = "com_google_cloud_go_resourcemanager",
        importpath = "cloud.google.com/go/resourcemanager",
        sum = "h1:VPfJFbWxrTYQzEXCDbJNpcvSB8eZhTSM0YHH146fIB8=",
        version = "v1.9.6",
    )
    go_repository(
        name = "com_google_cloud_go_resourcesettings",
        importpath = "cloud.google.com/go/resourcesettings",
        sum = "h1:l/IbRDDmGJFlR4bRZGtfYvix1Pu0jAKGLr7wgUtixHQ=",
        version = "v1.6.6",
    )
    go_repository(
        name = "com_google_cloud_go_retail",
        importpath = "cloud.google.com/go/retail",
        sum = "h1:AyVdElkdIU3JedWpX/qENbt8iUmKD+kiyj7ZpzguhTg=",
        version = "v1.16.1",
    )
    go_repository(
        name = "com_google_cloud_go_run",
        importpath = "cloud.google.com/go/run",
        sum = "h1:xQND6EJn1LgouCLPSfykkzagyr4gq4FKiRexNxXixV0=",
        version = "v1.3.6",
    )
    go_repository(
        name = "com_google_cloud_go_scheduler",
        importpath = "cloud.google.com/go/scheduler",
        sum = "h1:h1/VZk0XdkSh/jI7dDNp3V0Qi8yTkclOljDVPelXvAw=",
        version = "v1.10.7",
    )
    go_repository(
        name = "com_google_cloud_go_secretmanager",
        importpath = "cloud.google.com/go/secretmanager",
        sum = "h1:e5pIo/QEgiFiHPVJPxM5jbtUr4O/u5h2zLHYtkFQr24=",
        version = "v1.12.0",
    )
    go_repository(
        name = "com_google_cloud_go_security",
        importpath = "cloud.google.com/go/security",
        sum = "h1:dzc3oYxFG/9+uMwmxpnG+te4ZEEkbzvBoAFR1Va36N4=",
        version = "v1.16.0",
    )
    go_repository(
        name = "com_google_cloud_go_securitycenter",
        importpath = "cloud.google.com/go/securitycenter",
        sum = "h1:jKqAqF4iLwAzCe9vzncKuxGKStmxX/HL0ItvQJdFqEA=",
        version = "v1.29.0",
    )
    go_repository(
        name = "com_google_cloud_go_servicedirectory",
        importpath = "cloud.google.com/go/servicedirectory",
        sum = "h1:gkzx9Cd+OTOD+zY4u5vtbdvOx7vrvHYdeDiNdC6vKyw=",
        version = "v1.11.5",
    )
    go_repository(
        name = "com_google_cloud_go_shell",
        importpath = "cloud.google.com/go/shell",
        sum = "h1:/oJf9sboa2FfHWCmHXy+XfTRnZy8lC7O5zFyfE1EA6s=",
        version = "v1.7.6",
    )
    go_repository(
        name = "com_google_cloud_go_spanner",
        importpath = "cloud.google.com/go/spanner",
        sum = "h1:O9kf49dfaDRzPpKJNChHUJ+Bao02WPedZb8ZPyi02lI=",
        version = "v1.60.0",
    )
    go_repository(
        name = "com_google_cloud_go_speech",
        importpath = "cloud.google.com/go/speech",
        sum = "h1:evZ9B5S6OBGP5PeJXnYWW1UpgIYV4jbsXwsQr5x7fPM=",
        version = "v1.23.0",
    )
    go_repository(
        name = "com_google_cloud_go_storage",
        importpath = "cloud.google.com/go/storage",
        sum = "h1:MvraqHKhogCOTXTlct/9C3K3+Uy2jBmFYb3/Sp6dVtY=",
        version = "v1.39.1",
    )
    go_repository(
        name = "com_google_cloud_go_storagetransfer",
        importpath = "cloud.google.com/go/storagetransfer",
        sum = "h1:BawJo/u0P21cdxc2gB878qIFDC80COq2i0qWZeNevSw=",
        version = "v1.10.5",
    )
    go_repository(
        name = "com_google_cloud_go_talent",
        importpath = "cloud.google.com/go/talent",
        sum = "h1:4xgDFfOcgcSY0dUzaSc2tQCSRoLDEJ5CfbW5jfcgNJk=",
        version = "v1.6.7",
    )
    go_repository(
        name = "com_google_cloud_go_texttospeech",
        importpath = "cloud.google.com/go/texttospeech",
        sum = "h1:gLEyDoJeFGdoX7jSKbf+nJy7CTgjsSbCZXwzzkXgH9w=",
        version = "v1.7.6",
    )
    go_repository(
        name = "com_google_cloud_go_tpu",
        importpath = "cloud.google.com/go/tpu",
        sum = "h1:Cb1txkZYbKlGIZ4tQr9EjEB9snAOU6qyjvNezGXDunI=",
        version = "v1.6.6",
    )
    go_repository(
        name = "com_google_cloud_go_trace",
        importpath = "cloud.google.com/go/trace",
        sum = "h1:XF0Ejdw0NpRfAvuZUeQe3ClAG4R/9w5JYICo7l2weaw=",
        version = "v1.10.6",
    )
    go_repository(
        name = "com_google_cloud_go_translate",
        importpath = "cloud.google.com/go/translate",
        sum = "h1:SXOtKYnT7ZkeMtPwujaBOBt5Ph4kf6LIuMpAgu/WON0=",
        version = "v1.10.2",
    )
    go_repository(
        name = "com_google_cloud_go_video",
        importpath = "cloud.google.com/go/video",
        sum = "h1:y4jgUqDiWMfX+beJnlrnloBQxEIa9v+KrlkD2QJVpeE=",
        version = "v1.20.5",
    )
    go_repository(
        name = "com_google_cloud_go_videointelligence",
        importpath = "cloud.google.com/go/videointelligence",
        sum = "h1:P0Sa8+5KOEAVk/fazUNjVPzRCijCheZWJ8wL8xBn9Uk=",
        version = "v1.11.6",
    )
    go_repository(
        name = "com_google_cloud_go_vision_v2",
        importpath = "cloud.google.com/go/vision/v2",
        sum = "h1:kvR1sHcuPYat1wI3BYY7CEX2xLAcUHPYL6dOzV2Xf4Q=",
        version = "v2.8.1",
    )
    go_repository(
        name = "com_google_cloud_go_vmmigration",
        importpath = "cloud.google.com/go/vmmigration",
        sum = "h1:sbaWK76csqtk0TGPGCiJqZi7tfrU0LnrhUjZHI5YVdc=",
        version = "v1.7.6",
    )
    go_repository(
        name = "com_google_cloud_go_vmwareengine",
        importpath = "cloud.google.com/go/vmwareengine",
        sum = "h1:Mf8abigBstvjfSGq9twhtbMTCONL0Cjds+tGbc2pV0M=",
        version = "v1.1.2",
    )
    go_repository(
        name = "com_google_cloud_go_vpcaccess",
        importpath = "cloud.google.com/go/vpcaccess",
        sum = "h1:wbMTRdZ9P5+3D6oQWWqB/YxDCFR5m5OJ+b+hHwaBBQQ=",
        version = "v1.7.6",
    )
    go_repository(
        name = "com_google_cloud_go_webrisk",
        importpath = "cloud.google.com/go/webrisk",
        sum = "h1:rVhi2WOHcZF72X7spXVTFTmRGeFN4NFeW7/Ku7kgeug=",
        version = "v1.9.6",
    )
    go_repository(
        name = "com_google_cloud_go_websecurityscanner",
        importpath = "cloud.google.com/go/websecurityscanner",
        sum = "h1:YAwNB/HjKOVAy9D7W8Bkv8OQ9G2lqIqFOuJbyH5Xo4Q=",
        version = "v1.6.6",
    )
    go_repository(
        name = "com_google_cloud_go_workflows",
        importpath = "cloud.google.com/go/workflows",
        sum = "h1:hH511zmS93oE6j64m/eiGWnfgqailh/S8+f3MVNLcE8=",
        version = "v1.12.5",
    )
    go_repository(
        name = "in_gopkg_check_v1",
        importpath = "gopkg.in/check.v1",
        sum = "h1:yhCVgyC4o1eVCa2tZl7eS0r+SDo693bJlVdllGtEeKM=",
        version = "v0.0.0-20161208181325-20d25e280405",
    )
    go_repository(
        name = "in_gopkg_ini_v1",
        importpath = "gopkg.in/ini.v1",
        sum = "h1:Dgnx+6+nfE+IfzjUEISNeydPJh9AXNNsWbGP9KzCsOA=",
        version = "v1.67.0",
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
        sum = "h1:4Pp6oUg3+e/6M4C0A/3kJ2VYa++dsWVTtGgLVj5xtHg=",
        version = "v0.49.0",
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
        sum = "h1:jq9TW8u3so/bN+JPT166wjOI6/vQPF6Xe7nMNIltagk=",
        version = "v0.49.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel",
        build_directives = [
            "gazelle:resolve go go.opentelemetry.io/otel/trace @io_opentelemetry_go_otel_trace//:go_default_library",
            "gazelle:resolve go go.opentelemetry.io/otel/metric @io_opentelemetry_go_otel_metric//:go_default_library",
            "gazelle:resolve go go.opentelemetry.io/otel/metric/embedded @io_opentelemetry_go_otel_metric//embedded:go_default_library",
        ],
        importpath = "go.opentelemetry.io/otel",
        sum = "h1:0LAOdjNmQeSTzGBzduGe/rU4tZhMwL5rWgtp9Ku5Jfo=",
        version = "v1.24.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_metric",
        build_directives = [
            "gazelle:resolve go go.opentelemetry.io/otel/attribute @io_opentelemetry_go_otel//attribute:go_default_library",
        ],
        importpath = "go.opentelemetry.io/otel/metric",
        sum = "h1:6EhoGWWK28x1fbpA4tYTOWBkPefTDQnb8WSGXlc88kI=",
        version = "v1.24.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_sdk",
        importpath = "go.opentelemetry.io/otel/sdk",
        sum = "h1:YMPPDNymmQN3ZgczicBY3B6sf9n62Dlj9pWD3ucgoDw=",
        version = "v1.24.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_trace",
        build_directives = [
            "gazelle:resolve go go.opentelemetry.io/otel/attribute @io_opentelemetry_go_otel//attribute:go_default_library",
        ],
        importpath = "go.opentelemetry.io/otel/trace",
        sum = "h1:CsKnnL4dUAr/0llH9FKuc698G04IrpWV0MQA/Y1YELI=",
        version = "v1.24.0",
    )
    go_repository(
        name = "org_golang_google_api",
        build_directives = [
            "gazelle:resolve go go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp @io_opentelemetry_go_contrib_instrumentation_net_http_otelhttp//:go_default_library",
            "gazelle:resolve go go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc @io_opentelemetry_go_contrib_instrumentation_google_golang_org_grpc_otelgrpc//:go_default_library",
        ],
        importpath = "google.golang.org/api",
        sum = "h1:8a0p/BbPa65GlqGWtUKxot4p0TV8OGOfyTjtmkXNXmk=",
        version = "v0.177.0",
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
        sum = "h1:MTmrc2F5TZKDKXigcZetYkH04YwqtOPEQJwh4PPOgfk=",
        version = "v0.0.0-20240429193739-8cf5692501f6",
    )
    go_repository(
        name = "org_golang_google_genproto_googleapis_api",
        importpath = "google.golang.org/genproto/googleapis/api",
        sum = "h1:DTJM0R8LECCgFeUwApvcEJHz85HLagW8uRENYxHh1ww=",
        version = "v0.0.0-20240429193739-8cf5692501f6",
    )
    go_repository(
        name = "org_golang_google_genproto_googleapis_bytestream",
        importpath = "google.golang.org/genproto/googleapis/bytestream",
        sum = "h1:GtsRfMHDREQPg/snOM0QudeC54kX7UqodmmK4uELHLQ=",
        version = "v0.0.0-20240429193739-8cf5692501f6",
    )
    go_repository(
        name = "org_golang_google_genproto_googleapis_rpc",
        importpath = "google.golang.org/genproto/googleapis/rpc",
        sum = "h1:DujSIu+2tC9Ht0aPNA7jgj23Iq8Ewi5sgkQ++wdvonE=",
        version = "v0.0.0-20240429193739-8cf5692501f6",
    )
    go_repository(
        name = "org_golang_google_grpc",
        importpath = "google.golang.org/grpc",
        sum = "h1:MUeiw1B2maTVZthpU5xvASfTh3LDbxHd6IJ6QQVU+xM=",
        version = "v1.63.2",
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
        sum = "h1:9ddQBjfCyZPOHPUiPxpYESBLc+T8P3E+Vo4IbKZgFWg=",
        version = "v1.34.1",
    )
    go_repository(
        name = "org_golang_x_crypto",
        importpath = "golang.org/x/crypto",
        sum = "h1:g1v0xeRhjcugydODzvb3mEM9SQ0HGp9s/nh3COQ/C30=",
        version = "v0.22.0",
    )
    go_repository(
        name = "org_golang_x_exp",
        importpath = "golang.org/x/exp",
        sum = "h1:GoHiUyI/Tp2nVkLI2mCxVkOjsbSXD66ic0XW0js0R9g=",
        version = "v0.0.0-20230905200255-921286631fa9",
    )
    go_repository(
        name = "org_golang_x_lint",
        importpath = "golang.org/x/lint",
        sum = "h1:XQyxROzUlZH+WIQwySDgnISgOivlhjIEwaQaJEJrrN0=",
        version = "v0.0.0-20190313153728-d0100b6bd8b3",
    )
    go_repository(
        name = "org_golang_x_mod",
        importpath = "golang.org/x/mod",
        sum = "h1:LUYupSeNrTNCGzR/hVBk2NHZO4hXcVaW1k4Qx7rjPx8=",
        version = "v0.8.0",
    )
    go_repository(
        name = "org_golang_x_net",
        importpath = "golang.org/x/net",
        sum = "h1:1PcaxkF854Fu3+lvBIx5SYn9wRlBzzcnHZSiaFFAb0w=",
        version = "v0.24.0",
    )
    go_repository(
        name = "org_golang_x_oauth2",
        importpath = "golang.org/x/oauth2",
        sum = "h1:9+E/EZBCbTLNrbN35fHv/a/d/mOBatymz1zbtQrXpIg=",
        version = "v0.19.0",
    )
    go_repository(
        name = "org_golang_x_sync",
        importpath = "golang.org/x/sync",
        sum = "h1:YsImfSBoP9QPYL0xyKJPq0gcaJdG3rInoqxTWbfQu9M=",
        version = "v0.7.0",
    )
    go_repository(
        name = "org_golang_x_sys",
        importpath = "golang.org/x/sys",
        sum = "h1:q5f1RH2jigJ1MoAWp2KTp3gm5zAGFUTarQZ5U386+4o=",
        version = "v0.19.0",
    )
    go_repository(
        name = "org_golang_x_term",
        importpath = "golang.org/x/term",
        sum = "h1:+ThwsDv+tYfnJFhF4L8jITxu1tdTWRTZpdsWgEgjL6Q=",
        version = "v0.19.0",
    )
    go_repository(
        name = "org_golang_x_text",
        importpath = "golang.org/x/text",
        sum = "h1:ScX5w1eTa3QqT8oi6+ziP7dTV1S2+ALU0bI+0zXKWiQ=",
        version = "v0.14.0",
    )
    go_repository(
        name = "org_golang_x_time",
        importpath = "golang.org/x/time",
        sum = "h1:o7cqy6amK/52YcAKIPlM3a+Fpj35zvRj2TP+e1xFSfk=",
        version = "v0.5.0",
    )
    go_repository(
        name = "org_golang_x_tools",
        importpath = "golang.org/x/tools",
        sum = "h1:BOw41kyTf3PuCW1pVQf8+Cyg8pMlkYB1oo9iJ6D/lKM=",
        version = "v0.6.0",
    )
    go_repository(
        name = "org_golang_x_xerrors",
        importpath = "golang.org/x/xerrors",
        sum = "h1:E7g+9GITq07hpfrRu66IVDexMakfv52eLZ2CXBWiKr4=",
        version = "v0.0.0-20191204190536-9bdfabe68543",
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
