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
        sum = "h1:boJj011Hh+874zpIySeApCX4GeOjPl9qhRF3QuIZq+Q=",
        version = "v0.0.0-20241223141626-cff3c89139a3",
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
        name = "com_github_go_jose_go_jose_v4",
        importpath = "github.com/go-jose/go-jose/v4",
        sum = "h1:VsjPI33J0SB9vQM6PLmNjoHqMQNGPiZ0rHL7Ni7Q6/E=",
        version = "v4.0.4",
    )
    go_repository(
        name = "com_github_go_logr_logr",
        importpath = "github.com/go-logr/logr",
        sum = "h1:6pFjapn8bFcIbiKo3XT4j/BhANplGihG6tvd+8rYgrY=",
        version = "v1.4.2",
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
        sum = "h1:CNNw5U8lSiiBk7druxtSHHTsRWcxKoac6kZKm2peBBc=",
        version = "v1.2.4",
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
        sum = "h1:XYIDZApgAnrN1c855gTgghdIA6Stxb52D5RnLI1SLyw=",
        version = "v0.3.4",
    )
    go_repository(
        name = "com_github_googleapis_gapic_showcase",
        importpath = "github.com/googleapis/gapic-showcase",
        sum = "h1:Or8k6GQXYdHFPeQeSoRtq2oTk27VAxiwfh6IBJYhf3A=",
        version = "v0.35.5",
    )
    go_repository(
        name = "com_github_googleapis_gax_go_v2",
        build_directives = [
            "gazelle:resolve go google.golang.org/genproto/googleapis/rpc/errdetails @org_golang_google_genproto_googleapis_rpc//errdetails",
            "gazelle:resolve proto go google/rpc/code.proto @com_google_googleapis//google/rpc:code_go_proto",
            "gazelle:resolve proto proto google/rpc/code.proto @com_google_googleapis//google/rpc:code_proto",
        ],
        importpath = "github.com/googleapis/gax-go/v2",
        sum = "h1:hb0FFeiPaQskmvakKu5EbCbpntQn48jyHuvrkurSS/Q=",
        version = "v2.14.1",
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
        sum = "h1:3c8yed4lgqTt+oTQ+JNMDo+F4xprBf+O/il4ZC0nRLw=",
        version = "v1.25.0",
    )
    go_repository(
        name = "com_github_googlecloudplatform_opentelemetry_operations_go_exporter_metric",
        importpath = "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric",
        sum = "h1:o90wcURuxekmXrtxmYWTyNla0+ZEHhud6DI1ZTxd1vI=",
        version = "v0.49.0",
    )
    go_repository(
        name = "com_github_googlecloudplatform_opentelemetry_operations_go_internal_resourcemapping",
        importpath = "github.com/GoogleCloudPlatform/opentelemetry-operations-go/internal/resourcemapping",
        sum = "h1:GYUJLfvd++4DMuMhCFLgLXvFwofIxh/qOwoGuS/LTew=",
        version = "v0.49.0",
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
        sum = "h1:aYUidT7k73Pcl9nb2gScu7NSrKCSHIDE89b3+6Wq+LM=",
        version = "v2.2.2",
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
        sum = "h1:KvO1DLK/DRN07sQ1LQKScxyZJuNnedQ5/wKSR38lUII=",
        version = "v1.13.1",
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
        sum = "h1:e5/vxKd/rZsfSJMUX1agtjeTDf+qv1/JdBF8gg5k9ZM=",
        version = "v1.8.1",
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
        sum = "h1:RWq5SEjt8o25SROyN3z2OrDB9l7RPd3lwTWU8EcEdcI=",
        version = "v1.19.0",
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
        sum = "h1:Xv5erBjTwe/5IxqUQTdXv5kgmIvbHo3QQyRwhJsOfJA=",
        version = "v1.10.0",
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
        sum = "h1:bKXO7RXMFDkniAAvvuMrAPtQ/VHrs9e7J5UT3yrGdTY=",
        version = "v0.118.2",
    )
    go_repository(
        name = "com_google_cloud_go_accessapproval",
        importpath = "cloud.google.com/go/accessapproval",
        sum = "h1:axlU03FRiXDNupsmPG7LKzuS4Enk1gf598M62lWVB74=",
        version = "v1.8.3",
    )
    go_repository(
        name = "com_google_cloud_go_accesscontextmanager",
        importpath = "cloud.google.com/go/accesscontextmanager",
        sum = "h1:8zVoeiBa4erMCLEXltOcqVEsZhS26JZ5/Vrgs59eQiI=",
        version = "v1.9.3",
    )
    go_repository(
        name = "com_google_cloud_go_aiplatform",
        importpath = "cloud.google.com/go/aiplatform",
        sum = "h1:vnqsPkgcwlDEpWl9t6C3/HLfHeweuGXs2gcYTzH6dMs=",
        version = "v1.70.0",
    )
    go_repository(
        name = "com_google_cloud_go_analytics",
        importpath = "cloud.google.com/go/analytics",
        sum = "h1:hX6JAsNbXd2uVjqjIuMcKpmhIybKrEunBiGxK4SwEFI=",
        version = "v0.25.3",
    )
    go_repository(
        name = "com_google_cloud_go_apigateway",
        importpath = "cloud.google.com/go/apigateway",
        sum = "h1:Mn7cC5iWJz+cSMS/Hb+N2410CpZ6c8XpJKaexBl0Gxs=",
        version = "v1.7.3",
    )
    go_repository(
        name = "com_google_cloud_go_apigeeconnect",
        importpath = "cloud.google.com/go/apigeeconnect",
        sum = "h1:Wlr+30Tha0SMCvQYZKdrh+HkpOyl0CQFSlzeY/Gg1gs=",
        version = "v1.7.3",
    )
    go_repository(
        name = "com_google_cloud_go_apigeeregistry",
        importpath = "cloud.google.com/go/apigeeregistry",
        sum = "h1:j9CJg/oC884OX5cDpiwNt1ZlDXNV6Zb9Mp1YmRrOG0k=",
        version = "v0.9.3",
    )
    go_repository(
        name = "com_google_cloud_go_appengine",
        importpath = "cloud.google.com/go/appengine",
        sum = "h1:jrcanSzj9J1erevZuxldvsDwY+0k/DeFFzlnSfPGfL8=",
        version = "v1.9.3",
    )
    go_repository(
        name = "com_google_cloud_go_area120",
        importpath = "cloud.google.com/go/area120",
        sum = "h1:dPQ07rW4eku8OgNWDOaQaVGcE4+XfhH8BSbVwdVQ+wU=",
        version = "v0.9.3",
    )
    go_repository(
        name = "com_google_cloud_go_artifactregistry",
        importpath = "cloud.google.com/go/artifactregistry",
        sum = "h1:ZNXGB6+T7VmWdf6//VqxLdZ/sk0no8W0ujanHeJwDRw=",
        version = "v1.16.1",
    )
    go_repository(
        name = "com_google_cloud_go_asset",
        importpath = "cloud.google.com/go/asset",
        sum = "h1:6oNgjcs5KCPGBD71G0IccK6TfeFsEtBTyQ3Q+Dn09bs=",
        version = "v1.20.4",
    )
    go_repository(
        name = "com_google_cloud_go_assuredworkloads",
        importpath = "cloud.google.com/go/assuredworkloads",
        sum = "h1:RU1WhF1zMggdXAZ+ezYTn4Eh/FdiX7sz8lLXGERn4Po=",
        version = "v1.12.3",
    )
    go_repository(
        name = "com_google_cloud_go_auth",
        importpath = "cloud.google.com/go/auth",
        sum = "h1:AwoJbzUdxA/whv1qj3TLKwh3XX5sikny2fc40wUl+h0=",
        version = "v0.14.1",
    )
    go_repository(
        name = "com_google_cloud_go_auth_oauth2adapt",
        importpath = "cloud.google.com/go/auth/oauth2adapt",
        sum = "h1:/Lc7xODdqcEw8IrZ9SvwnlLX6j9FHQM74z6cBk9Rw6M=",
        version = "v0.2.7",
    )
    go_repository(
        name = "com_google_cloud_go_automl",
        importpath = "cloud.google.com/go/automl",
        sum = "h1:vkD+hQ75SMINMgJBT/KDpFYvfQLzJbtIQZdw0AWq8Rs=",
        version = "v1.14.4",
    )
    go_repository(
        name = "com_google_cloud_go_baremetalsolution",
        importpath = "cloud.google.com/go/baremetalsolution",
        sum = "h1:OL+KT+wCumdDhG44aeqGAdkwdT8Wa4Lh+o4INM+CQjw=",
        version = "v1.3.3",
    )
    go_repository(
        name = "com_google_cloud_go_batch",
        importpath = "cloud.google.com/go/batch",
        sum = "h1:TLfFZJXu+89CGbDK2mMql8f6HHFXarr8uUsaQ6wKatU=",
        version = "v1.11.5",
    )
    go_repository(
        name = "com_google_cloud_go_beyondcorp",
        importpath = "cloud.google.com/go/beyondcorp",
        sum = "h1:ezavJc0Gzh4N8zBskO/DnUVMWPa8lqH/tmQSyaknmCA=",
        version = "v1.1.3",
    )
    go_repository(
        name = "com_google_cloud_go_bigquery",
        importpath = "cloud.google.com/go/bigquery",
        sum = "h1:cDM3xEUUTf6RDepFEvNZokCysGFYoivHHTIZOWXbV2E=",
        version = "v1.66.0",
    )
    go_repository(
        name = "com_google_cloud_go_bigtable",
        importpath = "cloud.google.com/go/bigtable",
        sum = "h1:eIgi3QLcN4aq8p6n9U/zPgmHeBP34sm9FiKq4ik/ZoY=",
        version = "v1.34.0",
    )
    go_repository(
        name = "com_google_cloud_go_billing",
        importpath = "cloud.google.com/go/billing",
        sum = "h1:xMlO3hc5BI0s23tRB40bL40xSpxUR1x3E07Y5/VWcjU=",
        version = "v1.20.1",
    )
    go_repository(
        name = "com_google_cloud_go_binaryauthorization",
        importpath = "cloud.google.com/go/binaryauthorization",
        sum = "h1:X8JRfmk0/vyRqLusEyAPr0nZCK6RKae9omB4lrit0XI=",
        version = "v1.9.3",
    )
    go_repository(
        name = "com_google_cloud_go_certificatemanager",
        importpath = "cloud.google.com/go/certificatemanager",
        sum = "h1:2UP31fg7b+y3F0OmNbPHOKPEJ+6LOMfxAXX4p8xGCy4=",
        version = "v1.9.3",
    )
    go_repository(
        name = "com_google_cloud_go_channel",
        importpath = "cloud.google.com/go/channel",
        sum = "h1:oHyO3QAZ6kdf6SwqnUTBz50ND6Nk2rxZtboUiF4dgLE=",
        version = "v1.19.2",
    )
    go_repository(
        name = "com_google_cloud_go_cloudbuild",
        importpath = "cloud.google.com/go/cloudbuild",
        sum = "h1:0BRKyrCnWMHlnkwtNKdEwcvpgPm3OA3NqQhzDS5c7ek=",
        version = "v1.20.0",
    )
    go_repository(
        name = "com_google_cloud_go_clouddms",
        importpath = "cloud.google.com/go/clouddms",
        sum = "h1:T/rkkKE0KhQFMcO3+QWL82xakA9kRumLXY1lq5adIts=",
        version = "v1.8.3",
    )
    go_repository(
        name = "com_google_cloud_go_cloudtasks",
        importpath = "cloud.google.com/go/cloudtasks",
        sum = "h1:rXdznKjCa7WpzmvR2plrn2KJ+RZC1oYxPiRWNQjjf3k=",
        version = "v1.13.3",
    )
    go_repository(
        name = "com_google_cloud_go_compute",
        importpath = "cloud.google.com/go/compute",
        sum = "h1:SObuy8Fs6woazArpXp1fsHCw+ZH4iJ/8dGGTxUhHZQA=",
        version = "v1.31.1",
    )
    go_repository(
        name = "com_google_cloud_go_compute_metadata",
        importpath = "cloud.google.com/go/compute/metadata",
        sum = "h1:A6hENjEsCDtC1k8byVsgwvVcioamEHvZ4j01OwKxG9I=",
        version = "v0.6.0",
    )
    go_repository(
        name = "com_google_cloud_go_contactcenterinsights",
        importpath = "cloud.google.com/go/contactcenterinsights",
        sum = "h1:xJoZbX0HM1zht8KxAB38hs2v4Hcl+vXGLo454LrdwxA=",
        version = "v1.17.1",
    )
    go_repository(
        name = "com_google_cloud_go_container",
        importpath = "cloud.google.com/go/container",
        sum = "h1:eaMrgOl6NCk+Blhh29GgUVe3QGo7IiJQlP0w/EwLoV0=",
        version = "v1.42.1",
    )
    go_repository(
        name = "com_google_cloud_go_containeranalysis",
        importpath = "cloud.google.com/go/containeranalysis",
        sum = "h1:1D8U75BeotZxrG4jR6NYBtOt+uAeBsWhpBZmSYLakQw=",
        version = "v0.13.3",
    )
    go_repository(
        name = "com_google_cloud_go_datacatalog",
        importpath = "cloud.google.com/go/datacatalog",
        sum = "h1:3bAfstDB6rlHyK0TvqxEwaeOvoN9UgCs2bn03+VXmss=",
        version = "v1.24.3",
    )
    go_repository(
        name = "com_google_cloud_go_dataflow",
        importpath = "cloud.google.com/go/dataflow",
        sum = "h1:+7IfIXzYWSybIIDGK9FN2uqBsP/5b/Y0pBYzNhcmKSU=",
        version = "v0.10.3",
    )
    go_repository(
        name = "com_google_cloud_go_dataform",
        importpath = "cloud.google.com/go/dataform",
        sum = "h1:ZpGkZV8OyhUhvN/tfLffU2ki5ERTtqOunkIaiVAhmw0=",
        version = "v0.10.3",
    )
    go_repository(
        name = "com_google_cloud_go_datafusion",
        importpath = "cloud.google.com/go/datafusion",
        sum = "h1:FTMtsf2nfGGlDCuE84/RvVaCcTIYE7WQSB0noeO0cwI=",
        version = "v1.8.3",
    )
    go_repository(
        name = "com_google_cloud_go_datalabeling",
        importpath = "cloud.google.com/go/datalabeling",
        sum = "h1:PqoA3gnOWaLcHCnqoZe4jh3jmiv6+Z7W2xUUkw/j4jE=",
        version = "v0.9.3",
    )
    go_repository(
        name = "com_google_cloud_go_dataplex",
        importpath = "cloud.google.com/go/dataplex",
        sum = "h1:oswf105Cr2EwHrW2n7wk3nRZQf7hCe3apE/GqJ8yjvY=",
        version = "v1.21.0",
    )
    go_repository(
        name = "com_google_cloud_go_dataproc_v2",
        importpath = "cloud.google.com/go/dataproc/v2",
        sum = "h1:2vOv471LrcSn91VNzijcH+OkDRLa3kdyymOfKqbwZ4c=",
        version = "v2.10.1",
    )
    go_repository(
        name = "com_google_cloud_go_dataqna",
        importpath = "cloud.google.com/go/dataqna",
        sum = "h1:lGUj2FYs650EUPDMV6plWBAoh8qH9Bu1KCz1PUYF2VY=",
        version = "v0.9.3",
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
        sum = "h1:j5cIRYJHjx/058aHa4Slip7fl62UTGHCJc4GL9bxQLQ=",
        version = "v1.12.1",
    )
    go_repository(
        name = "com_google_cloud_go_deploy",
        importpath = "cloud.google.com/go/deploy",
        sum = "h1:Hm3pXBzMFJFPOdwtDkg5e/LP53bXqIpwQpjwsVasjhU=",
        version = "v1.26.1",
    )
    go_repository(
        name = "com_google_cloud_go_dialogflow",
        importpath = "cloud.google.com/go/dialogflow",
        sum = "h1:6fU4IKLpvgpXqiUCE8gUp8eV5u629SCtiyXMudXtZSg=",
        version = "v1.64.1",
    )
    go_repository(
        name = "com_google_cloud_go_dlp",
        importpath = "cloud.google.com/go/dlp",
        sum = "h1:qAEGTTtC97zuDm6YPBozNvy4BLBszVCJah3efNytl3g=",
        version = "v1.20.1",
    )
    go_repository(
        name = "com_google_cloud_go_documentai",
        importpath = "cloud.google.com/go/documentai",
        sum = "h1:52RfiUsoblXcE57CfKJGnITWLxRM30BcqNk/BKZl2LI=",
        version = "v1.35.1",
    )
    go_repository(
        name = "com_google_cloud_go_domains",
        importpath = "cloud.google.com/go/domains",
        sum = "h1:wnqN5YwMrtLSjn+HB2sChgmZ6iocOta4Q41giQsiRjY=",
        version = "v0.10.3",
    )
    go_repository(
        name = "com_google_cloud_go_edgecontainer",
        importpath = "cloud.google.com/go/edgecontainer",
        sum = "h1:SwQuHQiheVfL7b5ar/AXDberiaqr/yiue8X55AdWnZU=",
        version = "v1.4.1",
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
        sum = "h1:Paw495vxVyKuAgcQ2NQk09iRZBhPYRytknydEnvzcv4=",
        version = "v1.7.3",
    )
    go_repository(
        name = "com_google_cloud_go_eventarc",
        importpath = "cloud.google.com/go/eventarc",
        sum = "h1:RMymT7R87LaxKugOKwooOoheWXUm1NMeOfh3CVU9g54=",
        version = "v1.15.1",
    )
    go_repository(
        name = "com_google_cloud_go_filestore",
        importpath = "cloud.google.com/go/filestore",
        sum = "h1:vTXQI5qYKZ8dmCyHN+zVfaMyXCYbyZNM0CkPzpPUn7Q=",
        version = "v1.9.3",
    )
    go_repository(
        name = "com_google_cloud_go_firestore",
        importpath = "cloud.google.com/go/firestore",
        sum = "h1:cuydCaLS7Vl2SatAeivXyhbhDEIR8BDmtn4egDhIn2s=",
        version = "v1.18.0",
    )
    go_repository(
        name = "com_google_cloud_go_functions",
        importpath = "cloud.google.com/go/functions",
        sum = "h1:V0vCHSgFTUqKn57+PUXp1UfQY0/aMkveAw7wXeM3Lq0=",
        version = "v1.19.3",
    )
    go_repository(
        name = "com_google_cloud_go_gkebackup",
        importpath = "cloud.google.com/go/gkebackup",
        sum = "h1:djdExe/QgoKdp1gnIO1G5BoO1o/yGQOQJJEZ4QKTEXQ=",
        version = "v1.6.3",
    )
    go_repository(
        name = "com_google_cloud_go_gkeconnect",
        importpath = "cloud.google.com/go/gkeconnect",
        sum = "h1:YVpR0vlHSP/wD74PXEbKua4Aamud+wiYm4TiewNjD3M=",
        version = "v0.12.1",
    )
    go_repository(
        name = "com_google_cloud_go_gkehub",
        importpath = "cloud.google.com/go/gkehub",
        sum = "h1:yZ6lNJ9rNIoQmWrG14dB3+BFjS/EIRBf7Bo6jc5QWlE=",
        version = "v0.15.3",
    )
    go_repository(
        name = "com_google_cloud_go_gkemulticloud",
        importpath = "cloud.google.com/go/gkemulticloud",
        sum = "h1:JWe6PDNpNU88ZYvQkTd7w28fgeIs/gg6i0hcjUkgZ3M=",
        version = "v1.5.1",
    )
    go_repository(
        name = "com_google_cloud_go_gsuiteaddons",
        importpath = "cloud.google.com/go/gsuiteaddons",
        sum = "h1:QafYhVhyFGpidBUUlVhy6lUHFogFOycVYm9DV7MinhA=",
        version = "v1.7.3",
    )
    go_repository(
        name = "com_google_cloud_go_iam",
        importpath = "cloud.google.com/go/iam",
        sum = "h1:ZNfy/TYfn2uh/ukvhp783WhnbVluqf/tzOaqVUPlIPA=",
        version = "v1.4.0",
    )
    go_repository(
        name = "com_google_cloud_go_iap",
        importpath = "cloud.google.com/go/iap",
        sum = "h1:OWNYFHPyIBNHEAEFdVKOltYWe0g3izSrpFJW6Iidovk=",
        version = "v1.10.3",
    )
    go_repository(
        name = "com_google_cloud_go_ids",
        importpath = "cloud.google.com/go/ids",
        sum = "h1:wbFF7twu0XScFr+dtsVxTTttbFIRYt/SJjZiHFidtYE=",
        version = "v1.5.3",
    )
    go_repository(
        name = "com_google_cloud_go_iot",
        importpath = "cloud.google.com/go/iot",
        sum = "h1:aPWYQ+A1NX6ou/5U0nFAiXWdVT8OBxZYVZt2fBl2gWA=",
        version = "v1.8.3",
    )
    go_repository(
        name = "com_google_cloud_go_kms",
        importpath = "cloud.google.com/go/kms",
        sum = "h1:aQQ8esAIVZ1atdJRxihhdxGQ64/zEbJoJnCz/ydSmKg=",
        version = "v1.20.5",
    )
    go_repository(
        name = "com_google_cloud_go_language",
        importpath = "cloud.google.com/go/language",
        sum = "h1:8hmFMiS3wjjj3TX/U1zZYTgzwZoUjDbo9PaqcYEmuB4=",
        version = "v1.14.3",
    )
    go_repository(
        name = "com_google_cloud_go_lifesciences",
        importpath = "cloud.google.com/go/lifesciences",
        sum = "h1:Z05C+Ui953f0EQx9hJ1la6+QQl8ADrIs3iNwP5Elkpg=",
        version = "v0.10.3",
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
        sum = "h1:3tyw9rO3E2XVXzSApn1gyEEnH2K9SynNQjMlBi3uHLg=",
        version = "v0.6.4",
    )
    go_repository(
        name = "com_google_cloud_go_managedidentities",
        importpath = "cloud.google.com/go/managedidentities",
        sum = "h1:b9xGs24BIjfyvLgCtJoClOZpPi8d8owPgWe5JEINgaY=",
        version = "v1.7.3",
    )
    go_repository(
        name = "com_google_cloud_go_maps",
        importpath = "cloud.google.com/go/maps",
        sum = "h1:u7U/DieTxYYMDyvHQ00la5ayXLjDImTfnhdAsyPZXyY=",
        version = "v1.17.1",
    )
    go_repository(
        name = "com_google_cloud_go_mediatranslation",
        importpath = "cloud.google.com/go/mediatranslation",
        sum = "h1:nRBjeaMLipw05Br+qDAlSCcCQAAlat4mvpafztbEVgc=",
        version = "v0.9.3",
    )
    go_repository(
        name = "com_google_cloud_go_memcache",
        importpath = "cloud.google.com/go/memcache",
        sum = "h1:XH/qT3GbbSH//R0JTqR77lRpBxaa0N9sHgAzfwbTrv0=",
        version = "v1.11.3",
    )
    go_repository(
        name = "com_google_cloud_go_metastore",
        importpath = "cloud.google.com/go/metastore",
        sum = "h1:jDqeCw6NGDRAPT9+2Y/EjnWAB0BfCcUfmPLOyhB0eHs=",
        version = "v1.14.3",
    )
    go_repository(
        name = "com_google_cloud_go_monitoring",
        importpath = "cloud.google.com/go/monitoring",
        sum = "h1:M3nXww2gn9oZ/qWN2bZ35CjolnVHM3qnSbu6srCPgjk=",
        version = "v1.23.0",
    )
    go_repository(
        name = "com_google_cloud_go_networkconnectivity",
        importpath = "cloud.google.com/go/networkconnectivity",
        sum = "h1:YsVhG71ZC4FkqCP2oCI55x/JeGFyd7738Lt8iNTrzJw=",
        version = "v1.16.1",
    )
    go_repository(
        name = "com_google_cloud_go_networkmanagement",
        importpath = "cloud.google.com/go/networkmanagement",
        sum = "h1:oEoFGPYxTBsY47h0zdoE2ojV5aU/541D83UmxfjHWaE=",
        version = "v1.18.0",
    )
    go_repository(
        name = "com_google_cloud_go_networksecurity",
        importpath = "cloud.google.com/go/networksecurity",
        sum = "h1:JLJBFbxc8D7/OS81MyRoKhc2OvnVJxy5VMoQqqAhA7k=",
        version = "v0.10.3",
    )
    go_repository(
        name = "com_google_cloud_go_notebooks",
        importpath = "cloud.google.com/go/notebooks",
        sum = "h1:+9DrGJcZhCu6B2t0JJorekjIUBvg/KvBmXJYGmfvVvA=",
        version = "v1.12.3",
    )
    go_repository(
        name = "com_google_cloud_go_optimization",
        importpath = "cloud.google.com/go/optimization",
        sum = "h1:JwQjjoBZJpsoMQe/3mhVBMVZuSdagHg2pGOnwh2Jk+E=",
        version = "v1.7.3",
    )
    go_repository(
        name = "com_google_cloud_go_orchestration",
        importpath = "cloud.google.com/go/orchestration",
        sum = "h1:SFAsKyqvtS8VFcsq+JgXAeRkrksB9UH+AH7iFamkmlc=",
        version = "v1.11.4",
    )
    go_repository(
        name = "com_google_cloud_go_orgpolicy",
        importpath = "cloud.google.com/go/orgpolicy",
        sum = "h1:WFvgmjq/FO5GiXlhebltA9N14KdbLMcgG88ME+SWeBo=",
        version = "v1.14.2",
    )
    go_repository(
        name = "com_google_cloud_go_osconfig",
        importpath = "cloud.google.com/go/osconfig",
        sum = "h1:cyf1PMK5c2/WOIr5r2lxjH/XBJMA9P4zC8Tm10i0z3M=",
        version = "v1.14.3",
    )
    go_repository(
        name = "com_google_cloud_go_oslogin",
        importpath = "cloud.google.com/go/oslogin",
        sum = "h1:yomxnFPk+ye0zd0mJ15nn9fH4Ns7ex4xA3ll+u2q59A=",
        version = "v1.14.3",
    )
    go_repository(
        name = "com_google_cloud_go_phishingprotection",
        importpath = "cloud.google.com/go/phishingprotection",
        sum = "h1:T5mGFV0ggBKg3qt9myFRiGJu+nIUucuHLAtVpAuQ08I=",
        version = "v0.9.3",
    )
    go_repository(
        name = "com_google_cloud_go_policytroubleshooter",
        importpath = "cloud.google.com/go/policytroubleshooter",
        sum = "h1:ekIWI8JbKkpOfrgH/THGamQE/D16tcVBYJyrkseVcYI=",
        version = "v1.11.3",
    )
    go_repository(
        name = "com_google_cloud_go_privatecatalog",
        importpath = "cloud.google.com/go/privatecatalog",
        sum = "h1:fu2LABMi7CgZORQ2oNGbc0hoZ0FTqLkjGqIgAV/Kc7U=",
        version = "v0.10.4",
    )
    go_repository(
        name = "com_google_cloud_go_pubsub",
        importpath = "cloud.google.com/go/pubsub",
        sum = "h1:prYj8EEAAAwkp6WNoGTE4ahe0DgHoyJd5Pbop931zow=",
        version = "v1.45.3",
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
        sum = "h1:T5YGzaXwTesHaPDNTAuU3neDwZEnfjce70zufPFUwno=",
        version = "v2.19.4",
    )
    go_repository(
        name = "com_google_cloud_go_recommendationengine",
        importpath = "cloud.google.com/go/recommendationengine",
        sum = "h1:kBpcYPx4ys4lrDGKp4OhP2uy8h7UjlmLW/qoO5Xb2bY=",
        version = "v0.9.3",
    )
    go_repository(
        name = "com_google_cloud_go_recommender",
        importpath = "cloud.google.com/go/recommender",
        sum = "h1:dVlOjxsbjuhlwu4MIcyPWe09qVcDqc419iOjdPl5RHk=",
        version = "v1.13.3",
    )
    go_repository(
        name = "com_google_cloud_go_redis",
        importpath = "cloud.google.com/go/redis",
        sum = "h1:ROQXi5dCDSJCVezt/2nD1g+Ym0T6sio3DIzZ56NgMZI=",
        version = "v1.17.3",
    )
    go_repository(
        name = "com_google_cloud_go_resourcemanager",
        importpath = "cloud.google.com/go/resourcemanager",
        sum = "h1:SHOMw0kX0xWratC5Vb5VULBeWiGlPYAs82kiZqNtWpM=",
        version = "v1.10.3",
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
        sum = "h1:PT6CUlazIFIOLLJnV+bPBtiSH8iusKZ+FZRzZYFt2vk=",
        version = "v1.19.2",
    )
    go_repository(
        name = "com_google_cloud_go_run",
        importpath = "cloud.google.com/go/run",
        sum = "h1:aeVLygw0BGLH+Zbj8v3K3nEHvKlgoq+j8fcRJaYZtxY=",
        version = "v1.8.1",
    )
    go_repository(
        name = "com_google_cloud_go_scheduler",
        importpath = "cloud.google.com/go/scheduler",
        sum = "h1:p6+h8BoYJC+TvUijGBfORN6nuhOvJ3EwZ2H84CZ1ZEU=",
        version = "v1.11.3",
    )
    go_repository(
        name = "com_google_cloud_go_secretmanager",
        importpath = "cloud.google.com/go/secretmanager",
        sum = "h1:XVGHbcXEsbrgi4XHzgK5np81l1eO7O72WOXHhXUemrM=",
        version = "v1.14.3",
    )
    go_repository(
        name = "com_google_cloud_go_security",
        importpath = "cloud.google.com/go/security",
        sum = "h1:ya9gfY1ign6Yy25VMMMgZ9xy7D/TczDB0ElXcyWmEVE=",
        version = "v1.18.3",
    )
    go_repository(
        name = "com_google_cloud_go_securitycenter",
        importpath = "cloud.google.com/go/securitycenter",
        sum = "h1:H8UvBpcvs1OjI4jZuXX8xsN1IZo88a9PezHXkU2sGps=",
        version = "v1.35.3",
    )
    go_repository(
        name = "com_google_cloud_go_servicedirectory",
        importpath = "cloud.google.com/go/servicedirectory",
        sum = "h1:oFkCp6ti7fc7hzeROmOPQuPBHFqwyhcsv3Yrma28+uc=",
        version = "v1.12.3",
    )
    go_repository(
        name = "com_google_cloud_go_shell",
        importpath = "cloud.google.com/go/shell",
        sum = "h1:mjYgUsOtV3jl9xvDmcvlRRmA64deEPf52zOfuc68b/g=",
        version = "v1.8.3",
    )
    go_repository(
        name = "com_google_cloud_go_spanner",
        importpath = "cloud.google.com/go/spanner",
        sum = "h1:0bab8QDn6MNj9lNK6XyGAVFhMlhMU2waePPa6GZNoi8=",
        version = "v1.73.0",
    )
    go_repository(
        name = "com_google_cloud_go_speech",
        importpath = "cloud.google.com/go/speech",
        sum = "h1:qvURtJs7BQzQhbxWxwai0pT79S8KLVKJ/4W8igVkt1Y=",
        version = "v1.26.0",
    )
    go_repository(
        name = "com_google_cloud_go_storage",
        importpath = "cloud.google.com/go/storage",
        sum = "h1:3TbVkzTooBvnZsk7WaAQfOsNrdoM8QHusXA1cpk6QJs=",
        version = "v1.50.0",
    )
    go_repository(
        name = "com_google_cloud_go_storagetransfer",
        importpath = "cloud.google.com/go/storagetransfer",
        sum = "h1:W3v9A7MGBN7H9sAFstyciwP/1XEQhUhZfrjclmDnpMs=",
        version = "v1.12.1",
    )
    go_repository(
        name = "com_google_cloud_go_talent",
        importpath = "cloud.google.com/go/talent",
        sum = "h1:olv+s2g+LGXeJi+MYF1wI44/TwHaVnO0N7PiucVf5ZQ=",
        version = "v1.8.0",
    )
    go_repository(
        name = "com_google_cloud_go_texttospeech",
        importpath = "cloud.google.com/go/texttospeech",
        sum = "h1:YF/RdNb+jUEp22cIZCvqiFjfA5OxGE+Dxss3mhXU7oQ=",
        version = "v1.11.0",
    )
    go_repository(
        name = "com_google_cloud_go_tpu",
        importpath = "cloud.google.com/go/tpu",
        sum = "h1:BvMNijOb6Vd46Rr/SR5jWv1MPosOhVsi0UaeAGNjeds=",
        version = "v1.8.0",
    )
    go_repository(
        name = "com_google_cloud_go_trace",
        importpath = "cloud.google.com/go/trace",
        sum = "h1:c+I4YFjxRQjvAhRmSsmjpASUKq88chOX854ied0K/pE=",
        version = "v1.11.3",
    )
    go_repository(
        name = "com_google_cloud_go_translate",
        importpath = "cloud.google.com/go/translate",
        sum = "h1:XJ7LipYJi80BCgVk2lx1fwc7DIYM6oV2qx1G4IAGQ5w=",
        version = "v1.12.3",
    )
    go_repository(
        name = "com_google_cloud_go_video",
        importpath = "cloud.google.com/go/video",
        sum = "h1:C2FH+6yr6LCZC4fP0gm9FwJB/SRh5Ul88O5Sc/bL83I=",
        version = "v1.23.3",
    )
    go_repository(
        name = "com_google_cloud_go_videointelligence",
        importpath = "cloud.google.com/go/videointelligence",
        sum = "h1:zNTOUQyatGQtnCJ2dR3faRtpWQOlC8wszJqwG5CtwVM=",
        version = "v1.12.3",
    )
    go_repository(
        name = "com_google_cloud_go_vision_v2",
        importpath = "cloud.google.com/go/vision/v2",
        sum = "h1:dPvfDuPqPH+Yscf0f2f1RprvKkoo+N/j0a+IbLYX7Cs=",
        version = "v2.9.3",
    )
    go_repository(
        name = "com_google_cloud_go_vmmigration",
        importpath = "cloud.google.com/go/vmmigration",
        sum = "h1:dpCQq3pj2HnKdbvGTftdWymm3r4ovF7JW5z8xBcO2x4=",
        version = "v1.8.3",
    )
    go_repository(
        name = "com_google_cloud_go_vmwareengine",
        importpath = "cloud.google.com/go/vmwareengine",
        sum = "h1:TfuQr5j7qriINulUMotaC/+27SQaW2thIkF3Gb6VJ38=",
        version = "v1.3.3",
    )
    go_repository(
        name = "com_google_cloud_go_vpcaccess",
        importpath = "cloud.google.com/go/vpcaccess",
        sum = "h1:vxVaoFM64M/ht619c4wZNF0iq0QPaMWElOh7Ns4r41A=",
        version = "v1.8.3",
    )
    go_repository(
        name = "com_google_cloud_go_webrisk",
        importpath = "cloud.google.com/go/webrisk",
        sum = "h1:yh0v/5n49VO4/i9pYfDm1gLJUj1Ph3Xzegn8WvK9YRA=",
        version = "v1.10.3",
    )
    go_repository(
        name = "com_google_cloud_go_websecurityscanner",
        importpath = "cloud.google.com/go/websecurityscanner",
        sum = "h1:/uxhVCWKXzPw5pVfnBOVjaSiQ6Bm0tDExDOCLV40thw=",
        version = "v1.7.3",
    )
    go_repository(
        name = "com_google_cloud_go_workflows",
        importpath = "cloud.google.com/go/workflows",
        sum = "h1:lNFDMranJymDEB7cTI7DI9czbc1WU0RWY9KCEv9zuDY=",
        version = "v1.13.3",
    )
    go_repository(
        name = "dev_cel_expr",
        importpath = "cel.dev/expr",
        sum = "h1:NciYrtDRIR0lNCnH1LFJegdjspNx9fI59O7TWcua/W4=",
        version = "v0.19.1",
    )
    go_repository(
        name = "in_gopkg_check_v1",
        importpath = "gopkg.in/check.v1",
        sum = "h1:Hei/4ADfdWqJk1ZMxUNpqntNwaWcugrBjAiHlqqRiVk=",
        version = "v1.0.0-20201130134442-10cb98267c6c",
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
        name = "io_opentelemetry_go_auto_sdk",
        importpath = "go.opentelemetry.io/auto/sdk",
        sum = "h1:cH53jehLUN6UFLY71z+NDOiNJqDdPRaXzTel0sJySYA=",
        version = "v1.1.0",
    )
    go_repository(
        name = "io_opentelemetry_go_contrib_detectors_gcp",
        importpath = "go.opentelemetry.io/contrib/detectors/gcp",
        sum = "h1:JRxssobiPg23otYU5SbWtQC//snGVIM3Tx6QRzlQBao=",
        version = "v1.34.0",
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
        sum = "h1:x7wzEgXfnzJcHDwStJT+mxOz4etr2EcexjqhBvmoakw=",
        version = "v0.60.0",
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
        sum = "h1:sbiXRNDSWJOTobXh5HyQKjq6wUC5tNybqjIqDpAY4CU=",
        version = "v0.60.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel",
        build_directives = [
            "gazelle:resolve go go.opentelemetry.io/otel/trace @io_opentelemetry_go_otel_trace//:go_default_library",
            "gazelle:resolve go go.opentelemetry.io/otel/metric @io_opentelemetry_go_otel_metric//:go_default_library",
            "gazelle:resolve go go.opentelemetry.io/otel/metric/embedded @io_opentelemetry_go_otel_metric//embedded:go_default_library",
        ],
        importpath = "go.opentelemetry.io/otel",
        sum = "h1:xKWKPxrxB6OtMCbmMY021CqC45J+3Onta9MqjhnusiQ=",
        version = "v1.35.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_metric",
        build_directives = [
            "gazelle:resolve go go.opentelemetry.io/otel/attribute @io_opentelemetry_go_otel//attribute:go_default_library",
        ],
        importpath = "go.opentelemetry.io/otel/metric",
        sum = "h1:0znxYu2SNyuMSQT4Y9WDWej0VpcsxkuklLa4/siN90M=",
        version = "v1.35.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_sdk",
        importpath = "go.opentelemetry.io/otel/sdk",
        sum = "h1:iPctf8iprVySXSKJffSS79eOjl9pvxV9ZqOWT0QejKY=",
        version = "v1.35.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_sdk_metric",
        importpath = "go.opentelemetry.io/otel/sdk/metric",
        sum = "h1:1RriWBmCKgkeHEhM7a2uMjMUfP7MsOF5JpUCaEqEI9o=",
        version = "v1.35.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_trace",
        build_directives = [
            "gazelle:resolve go go.opentelemetry.io/otel/attribute @io_opentelemetry_go_otel//attribute:go_default_library",
        ],
        importpath = "go.opentelemetry.io/otel/trace",
        sum = "h1:dPpEfJu1sDIqruz7BHFG3c7528f6ddfSWfFDVt/xgMs=",
        version = "v1.35.0",
    )
    go_repository(
        name = "org_golang_google_api",
        build_directives = [
            "gazelle:resolve go go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp @io_opentelemetry_go_contrib_instrumentation_net_http_otelhttp//:go_default_library",
            "gazelle:resolve go go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc @io_opentelemetry_go_contrib_instrumentation_google_golang_org_grpc_otelgrpc//:go_default_library",
        ],
        importpath = "google.golang.org/api",
        sum = "h1:qzaJfLhDsbMeFee8zBRdt/Nc+xmOuafD/dbdgGfutOU=",
        version = "v0.221.0",
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
        sum = "h1:Pw6WnI9W/LIdRxqK7T6XGugGbHIRl5Q7q3BssH6xk4s=",
        version = "v0.0.0-20250122153221-138b5a5a4fd4",
    )
    go_repository(
        name = "org_golang_google_genproto_googleapis_api",
        importpath = "google.golang.org/genproto/googleapis/api",
        sum = "h1:L9JNMl/plZH9wmzQUHleO/ZZDSN+9Gh41wPczNy+5Fk=",
        version = "v0.0.0-20250207221924-e9438ea467c6",
    )
    go_repository(
        name = "org_golang_google_genproto_googleapis_bytestream",
        importpath = "google.golang.org/genproto/googleapis/bytestream",
        sum = "h1:2FBsm1NN1gB068gQXGnpWH7WDYlHz9ITquoAeTMy4AI=",
        version = "v0.0.0-20250207221924-e9438ea467c6",
    )
    go_repository(
        name = "org_golang_google_genproto_googleapis_rpc",
        importpath = "google.golang.org/genproto/googleapis/rpc",
        sum = "h1:51aaUVRocpvUOSQKM6Q7VuoaktNIaMCLuhZB6DKksq4=",
        version = "v0.0.0-20250218202821-56aae31c358a",
    )
    go_repository(
        name = "org_golang_google_grpc",
        importpath = "google.golang.org/grpc",
        sum = "h1:KnzCueW4s+8ojAPZ+NnyZAELjsIMJGteKjKejieEC7M=",
        version = "v1.71.2",
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
        sum = "h1:z1NpPI8ku2WgiWnf+t9wTPsn6eP1L7ksHUlkfLvd9xY=",
        # TODO(https://github.com/googleapis/gapic-generator-go/issues/1608): Don't hard-code old version
        version = "v1.36.6",
    )
    go_repository(
        name = "org_golang_x_crypto",
        importpath = "golang.org/x/crypto",
        sum = "h1:IOBPskki6Lysi0lo9qQvbxiQ+FvsCC/YWOecCHAixus=",
        version = "v0.33.0",
    )
    go_repository(
        name = "org_golang_x_exp",
        importpath = "golang.org/x/exp",
        sum = "h1:GoHiUyI/Tp2nVkLI2mCxVkOjsbSXD66ic0XW0js0R9g=",
        version = "v0.0.0-20230905200255-921286631fa9",
    )
    go_repository(
        name = "org_golang_x_mod",
        importpath = "golang.org/x/mod",
        sum = "h1:zY54UmvipHiNd+pm+m0x9KhZ9hl1/7QNMyxXbc6ICqA=",
        version = "v0.17.0",
    )
    go_repository(
        name = "org_golang_x_net",
        importpath = "golang.org/x/net",
        sum = "h1:T5GQRQb2y08kTAByq9L4/bz8cipCdA8FbRTXewonqY8=",
        version = "v0.35.0",
    )
    go_repository(
        name = "org_golang_x_oauth2",
        importpath = "golang.org/x/oauth2",
        sum = "h1:da9Vo7/tDv5RH/7nZDz1eMGS/q1Vv1N/7FCrBhI9I3M=",
        version = "v0.27.0",
    )
    go_repository(
        name = "org_golang_x_sync",
        importpath = "golang.org/x/sync",
        sum = "h1:GGz8+XQP4FvTTrjZPzNKTMFtSXH80RAzG+5ghFPgK9w=",
        version = "v0.11.0",
    )
    go_repository(
        name = "org_golang_x_sys",
        importpath = "golang.org/x/sys",
        sum = "h1:QjkSwP/36a20jFYWkSue1YwXzLmsV5Gfq7Eiy72C1uc=",
        version = "v0.30.0",
    )
    go_repository(
        name = "org_golang_x_term",
        importpath = "golang.org/x/term",
        sum = "h1:L6pJp37ocefwRRtYPKSWOWzOtWSxVajvz2ldH/xi3iU=",
        version = "v0.29.0",
    )
    go_repository(
        name = "org_golang_x_text",
        importpath = "golang.org/x/text",
        sum = "h1:bofq7m3/HAFvbF51jz3Q9wLg3jkvSPuiZu/pD1XwgtM=",
        version = "v0.22.0",
    )
    go_repository(
        name = "org_golang_x_time",
        importpath = "golang.org/x/time",
        sum = "h1:3usCWA8tQn0L8+hFJQNgzpWbd89begxN66o1Ojdn5L4=",
        version = "v0.10.0",
    )
    go_repository(
        name = "org_golang_x_tools",
        importpath = "golang.org/x/tools",
        sum = "h1:vU5i/LfpvrRCpgM/VPfJLg5KjxD3E+hfT1SH+d9zLwg=",
        version = "v0.21.1-0.20240508182429-e35e4ccd0d2d",
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
