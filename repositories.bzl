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
        name = "com_github_cespare_xxhash_v2",
        importpath = "github.com/cespare/xxhash/v2",
        sum = "h1:UL815xU9SqsFlibzuggzjXhog7bL6oX9BbNZnL2UFvs=",
        version = "v2.3.0",
    )
    go_repository(
        name = "com_github_cncf_xds_go",
        importpath = "github.com/cncf/xds/go",
        sum = "h1:aBangftG7EVZoUb69Os8IaYg++6uMOdKK83QtkkvJik=",
        version = "v0.0.0-20260202195803-dba9d589def2",
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
        sum = "h1:U9qPSI2PIWSS1VwoXQT9A3Wy9MM3WgvqSxFWenqJduM=",
        version = "v1.1.2-0.20180830191138-d8f796af33cc",
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
        sum = "h1:u3riX6BoYRfF4Dr7dwSOroNfdSbEPe9Yyl09/B6wBrQ=",
        version = "v1.37.0",
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
        sum = "h1:MVQghNeW+LZcmXe7SY1V36Z+WFMDjpqGAGacLe2T0ds=",
        version = "v1.3.3",
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
        sum = "h1:moDMcTHmvE6Groj34emNPLs/qtYXRVcd6S7NHbHz3kA=",
        version = "v4.1.4",
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
        sum = "h1:73NfMHdiqo9JFU9+7a5ExpVa10/R29pXfZIaW559nrg=",
        version = "v0.3.17",
    )
    go_repository(
        name = "com_github_googleapis_gapic_showcase",
        importpath = "github.com/googleapis/gapic-showcase",
        sum = "h1:OEAUCT7DNrAcL9sNNLTVnA/kKDeUCw7ITU8jHFK99zo=",
        version = "v0.41.1",
    )
    go_repository(
        name = "com_github_googleapis_gax_go_v2",
        build_directives = [
            "gazelle:resolve go google.golang.org/genproto/googleapis/rpc/errdetails @org_golang_google_genproto_googleapis_rpc//errdetails",
            "gazelle:resolve proto go google/rpc/code.proto @com_google_googleapis//google/rpc:code_go_proto",
            "gazelle:resolve proto proto google/rpc/code.proto @com_google_googleapis//google/rpc:code_proto",
        ],
        importpath = "github.com/googleapis/gax-go/v2",
        sum = "h1:Tchl7qkvE7Ip3y+ztvNufYFvkfqTe7NfLTYGIdJRLuE=",
        version = "v2.23.0",
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
        sum = "h1:rIkQfkCOVKc1OiRCNcSDD8ml5RJlZbH/Xsq7lbpynwc=",
        version = "v1.32.0",
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
        sum = "h1:Jamvg5psRIccs7FGNTlIRMkT8wgtp5eCXdBlqhYGL6U=",
        version = "v1.0.1-0.20181226105442-5d4384ee4fb2",
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
        sum = "h1:kx3RQSS0VglTngQTPfywgVj+xFgre/vJouGETBlPU1E=",
        version = "v1.13.0",
    )
    go_repository(
        name = "com_google_cloud_go_accesscontextmanager",
        importpath = "cloud.google.com/go/accesscontextmanager",
        sum = "h1:0baVgug9IFV8y5cIoD4iN/m/QtJzk2k+j3MW+zJWCVw=",
        version = "v1.15.0",
    )
    go_repository(
        name = "com_google_cloud_go_aiplatform",
        importpath = "cloud.google.com/go/aiplatform",
        sum = "h1:5PxeWpQfkAyN8mtgtVp4H5nP+ntwPLFMjNREVdiolR0=",
        version = "v1.126.0",
    )
    go_repository(
        name = "com_google_cloud_go_analytics",
        importpath = "cloud.google.com/go/analytics",
        sum = "h1:GuwmzJHIaQRvtko6g4wW1ngQ7rpDvDLZ8M3iP3kiflU=",
        version = "v0.35.0",
    )
    go_repository(
        name = "com_google_cloud_go_apigateway",
        importpath = "cloud.google.com/go/apigateway",
        sum = "h1:IENNmFQlMsl58RlvYHf6wxBNQFeP+lAnVNcyxtmByTE=",
        version = "v1.13.0",
    )
    go_repository(
        name = "com_google_cloud_go_apigeeconnect",
        importpath = "cloud.google.com/go/apigeeconnect",
        sum = "h1:YvZhi0QkHobBRMrN6ri9PUPnnhgEd4VpW0+yo1WI7r0=",
        version = "v1.12.0",
    )
    go_repository(
        name = "com_google_cloud_go_apigeeregistry",
        importpath = "cloud.google.com/go/apigeeregistry",
        sum = "h1:2zIogaYq7dF4E1dAGC4MUnasRuSmohsMl7wVuJ963/4=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_google_cloud_go_appengine",
        importpath = "cloud.google.com/go/appengine",
        sum = "h1:0MiBM2KGk1WAhh489WoELQ7IeW2qw6HTE6nZIlrOw5U=",
        version = "v1.15.0",
    )
    go_repository(
        name = "com_google_cloud_go_area120",
        importpath = "cloud.google.com/go/area120",
        sum = "h1:9v6HeQsBpfRqoZtNeLI4V8SQ7Jcmt6s2gmgbjl0fVFc=",
        version = "v0.15.0",
    )
    go_repository(
        name = "com_google_cloud_go_artifactregistry",
        importpath = "cloud.google.com/go/artifactregistry",
        sum = "h1:iq5kkdY2FJY8RkkNghUelb3EOcqyy7dVhfDGVk/Xw9g=",
        version = "v1.26.0",
    )
    go_repository(
        name = "com_google_cloud_go_asset",
        importpath = "cloud.google.com/go/asset",
        sum = "h1:M6YE1exBuZQhTi8wfIKrwjKK6b2ySfHKgrqXjPZA4EM=",
        version = "v1.28.0",
    )
    go_repository(
        name = "com_google_cloud_go_assuredworkloads",
        importpath = "cloud.google.com/go/assuredworkloads",
        sum = "h1:XyNQ9SBX0eUkZYCZhC8dOY6ZFeSF6/OIjHn5dQs78Jo=",
        version = "v1.19.0",
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
        sum = "h1:brKXr4i9AMVfgbCKOYSu4rnIoz1ILIgcIAy3nerzEGc=",
        version = "v1.21.0",
    )
    go_repository(
        name = "com_google_cloud_go_baremetalsolution",
        importpath = "cloud.google.com/go/baremetalsolution",
        sum = "h1:6MaZilXzGZJ3mGqbSCdDezs9yrze6kcmBABsRUZ/aYI=",
        version = "v1.10.0",
    )
    go_repository(
        name = "com_google_cloud_go_batch",
        importpath = "cloud.google.com/go/batch",
        sum = "h1:nTL5x9HA1yVdPIdeqcbLjNnV9iYDYvMNmP6AnZ+qbHU=",
        version = "v1.20.0",
    )
    go_repository(
        name = "com_google_cloud_go_beyondcorp",
        importpath = "cloud.google.com/go/beyondcorp",
        sum = "h1:rsyle6zjxce2s57BpULw5u170xS+KSjFENSY7C0MZbs=",
        version = "v1.8.0",
    )
    go_repository(
        name = "com_google_cloud_go_bigquery",
        importpath = "cloud.google.com/go/bigquery",
        sum = "h1:+tW6oRuP/4dOgdl/p32on0+3OktVQtU+veOMxhjoOhE=",
        version = "v1.79.0",
    )
    go_repository(
        name = "com_google_cloud_go_bigtable",
        importpath = "cloud.google.com/go/bigtable",
        sum = "h1:NGLgDSr/i79BTGCjxH/maPKxyvl5q8/SsBsyLK52kdI=",
        version = "v1.47.0",
    )
    go_repository(
        name = "com_google_cloud_go_billing",
        importpath = "cloud.google.com/go/billing",
        sum = "h1:6RRjbRd6iZKZFb7/MgRvmXKq/Ism02ckkZLJazj4CQ0=",
        version = "v1.26.0",
    )
    go_repository(
        name = "com_google_cloud_go_binaryauthorization",
        importpath = "cloud.google.com/go/binaryauthorization",
        sum = "h1:f1iNUaHWP9XPdh6grftTlIRuJSBytud4j51l1Jce94E=",
        version = "v1.16.0",
    )
    go_repository(
        name = "com_google_cloud_go_certificatemanager",
        importpath = "cloud.google.com/go/certificatemanager",
        sum = "h1:PwljEZlI3lZgoFUG0pnIBNOwrQUaAP+gS0ZahTwUZBA=",
        version = "v1.15.0",
    )
    go_repository(
        name = "com_google_cloud_go_channel",
        importpath = "cloud.google.com/go/channel",
        sum = "h1:fqDZDzVw61c4iIUe7rQ+cBYlBdSZ3yqnzw5WrZTPjqs=",
        version = "v1.27.0",
    )
    go_repository(
        name = "com_google_cloud_go_cloudbuild",
        importpath = "cloud.google.com/go/cloudbuild",
        sum = "h1:tlF+KSIJJ0kaKkxACpk9htZ1euDpyOqPuO1aHR4j3oI=",
        version = "v1.32.0",
    )
    go_repository(
        name = "com_google_cloud_go_clouddms",
        importpath = "cloud.google.com/go/clouddms",
        sum = "h1:oaktwVyUeKTuzjpde3xbpmc9bAQPQq6WJWOBr7K6YEA=",
        version = "v1.14.0",
    )
    go_repository(
        name = "com_google_cloud_go_cloudtasks",
        importpath = "cloud.google.com/go/cloudtasks",
        sum = "h1:KzT7hfix/9/xAf20tNPIxwX59XGpRF0Lun2t8LHOj9E=",
        version = "v1.18.0",
    )
    go_repository(
        name = "com_google_cloud_go_compute",
        importpath = "cloud.google.com/go/compute",
        sum = "h1:tJ7lKJ8YEVa6vZX03Jc8o1YePbjKDOQhDw1BscMZ1bs=",
        version = "v1.62.0",
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
        sum = "h1:EtrRf8T7ajR+dR8nseDxwKy3o1VaWP1GFAu7lgMmT9A=",
        version = "v1.23.0",
    )
    go_repository(
        name = "com_google_cloud_go_container",
        importpath = "cloud.google.com/go/container",
        sum = "h1:KdeWHqwlOvABdhkhSTSb4QFySMOuXnc3UUoQMmq9lic=",
        version = "v1.51.0",
    )
    go_repository(
        name = "com_google_cloud_go_containeranalysis",
        importpath = "cloud.google.com/go/containeranalysis",
        sum = "h1:89ZhFvJHWzX9/jUJy/IEVbtSa2hv4+NtsvwWMEsUYnY=",
        version = "v0.19.0",
    )
    go_repository(
        name = "com_google_cloud_go_datacatalog",
        importpath = "cloud.google.com/go/datacatalog",
        sum = "h1:8V80PpoAGdOOr2QhBrp4wZ66MDCbATdAB/fmVmo5rlU=",
        version = "v1.33.0",
    )
    go_repository(
        name = "com_google_cloud_go_dataflow",
        importpath = "cloud.google.com/go/dataflow",
        sum = "h1:BchGCAl9QIZ/pyTGokv5V3daxBTrXsOmrU4ARXoTNd4=",
        version = "v0.16.0",
    )
    go_repository(
        name = "com_google_cloud_go_dataform",
        importpath = "cloud.google.com/go/dataform",
        sum = "h1:eAIsWhr6AbUIX+wOdfctO6hAr9UoQnTR2qQW739XH1I=",
        version = "v1.2.0",
    )
    go_repository(
        name = "com_google_cloud_go_datafusion",
        importpath = "cloud.google.com/go/datafusion",
        sum = "h1:56rjOW8xFBnTOLgTEWdU09KhVC2/X+y/ukzux7zYqBE=",
        version = "v1.14.0",
    )
    go_repository(
        name = "com_google_cloud_go_datalabeling",
        importpath = "cloud.google.com/go/datalabeling",
        sum = "h1:9TL2kwlO/ODD1fne3uvgTLa53SttwrQi92Bn7KWnKhg=",
        version = "v0.15.0",
    )
    go_repository(
        name = "com_google_cloud_go_dataplex",
        importpath = "cloud.google.com/go/dataplex",
        sum = "h1:r4O6RHEkwHr41ItsTUym0ZgH1iDf48R+W8fI7KoMNl8=",
        version = "v1.36.0",
    )
    go_repository(
        name = "com_google_cloud_go_dataproc_v2",
        importpath = "cloud.google.com/go/dataproc/v2",
        sum = "h1:IF1wxkXvcgoTP92NVqX6QNDf2SdKNZFPa2IWeSKP310=",
        version = "v2.25.0",
    )
    go_repository(
        name = "com_google_cloud_go_dataqna",
        importpath = "cloud.google.com/go/dataqna",
        sum = "h1:VR1y/NuN+RvPkmfmlowbmk60BLINR4MgZAusVFIGQjU=",
        version = "v0.13.0",
    )
    go_repository(
        name = "com_google_cloud_go_datastore",
        importpath = "cloud.google.com/go/datastore",
        sum = "h1:zUjMnCLCcRZVDSdQIXsbnNCl1SVRNw5Jm0J77gPaPKs=",
        version = "v1.25.0",
    )
    go_repository(
        name = "com_google_cloud_go_datastream",
        importpath = "cloud.google.com/go/datastream",
        sum = "h1:hrWGASBj85xs6wpnx9ifdSBYiRLhKDgQrA1HjeNYBhc=",
        version = "v1.21.0",
    )
    go_repository(
        name = "com_google_cloud_go_deploy",
        importpath = "cloud.google.com/go/deploy",
        sum = "h1:2PdZ8kbIztLVRMrIN/wHbalZeZZikQlxBptBFK+eEtM=",
        version = "v1.33.0",
    )
    go_repository(
        name = "com_google_cloud_go_dialogflow",
        importpath = "cloud.google.com/go/dialogflow",
        sum = "h1:BymJ6nPotDcnphK2/D0VHMAfakjJlIV5GqxLHQg/Tos=",
        version = "v1.84.0",
    )
    go_repository(
        name = "com_google_cloud_go_dlp",
        importpath = "cloud.google.com/go/dlp",
        sum = "h1:SI7KSLOtAfhEAw3af8NxpKHGHbJ9BkoM7805TMaU6m4=",
        version = "v1.34.0",
    )
    go_repository(
        name = "com_google_cloud_go_documentai",
        importpath = "cloud.google.com/go/documentai",
        sum = "h1:WLZOaWvK4NH7xMLxANsDgIdXugeqmVyYRCbVN1m87DQ=",
        version = "v1.49.0",
    )
    go_repository(
        name = "com_google_cloud_go_domains",
        importpath = "cloud.google.com/go/domains",
        sum = "h1:h+lcYPxyEulj5SYuH3+OwOsLZb4DTW11QX6nlxxXQ1I=",
        version = "v0.16.0",
    )
    go_repository(
        name = "com_google_cloud_go_edgecontainer",
        importpath = "cloud.google.com/go/edgecontainer",
        sum = "h1:4WOjcIZRCB4ynxH9Zs1tteIcm2LD58VUCLLmIwLYUAc=",
        version = "v1.10.0",
    )
    go_repository(
        name = "com_google_cloud_go_errorreporting",
        importpath = "cloud.google.com/go/errorreporting",
        sum = "h1:LlE2SVIbz0k+OSeNTksk34inr3Fy62JMhHUvNaS8f7c=",
        version = "v0.9.0",
    )
    go_repository(
        name = "com_google_cloud_go_essentialcontacts",
        importpath = "cloud.google.com/go/essentialcontacts",
        sum = "h1:+AtGn+hYLdr62sMqP0ZGMRzp9A4t+MSWe8eZoD/ho+M=",
        version = "v1.12.0",
    )
    go_repository(
        name = "com_google_cloud_go_eventarc",
        importpath = "cloud.google.com/go/eventarc",
        sum = "h1:+XAJNEmxJIPZLMtzLYdUZN61Kmqx62hWLry4AsWaTwo=",
        version = "v1.25.0",
    )
    go_repository(
        name = "com_google_cloud_go_filestore",
        importpath = "cloud.google.com/go/filestore",
        sum = "h1:2GXy5jvq2oUX6HnCpx3GjzSqlcMz8K9gkHaF9HPcYNE=",
        version = "v1.16.0",
    )
    go_repository(
        name = "com_google_cloud_go_firestore",
        importpath = "cloud.google.com/go/firestore",
        sum = "h1:x0Z3hrgjYgo2wI9whuBRQcNc2hYwzZDQy/7pkUXbXcs=",
        version = "v1.24.0",
    )
    go_repository(
        name = "com_google_cloud_go_functions",
        importpath = "cloud.google.com/go/functions",
        sum = "h1:ndUtLkam3XF9b0t2zVACH9D/EgFBISbVuQFqRz/X58k=",
        version = "v1.25.0",
    )
    go_repository(
        name = "com_google_cloud_go_gkebackup",
        importpath = "cloud.google.com/go/gkebackup",
        sum = "h1:li3BtGRis1QYrkLo8+Iq2wf5WbP9v3sz9VoUw8WqgaA=",
        version = "v1.14.0",
    )
    go_repository(
        name = "com_google_cloud_go_gkeconnect",
        importpath = "cloud.google.com/go/gkeconnect",
        sum = "h1:IDuAD/w0Xph8k/31vdt6pBkDdRx55lVXz7gW8ABOfCs=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_google_cloud_go_gkehub",
        importpath = "cloud.google.com/go/gkehub",
        sum = "h1:gHKPoWQuWpd9dXLoC58dxwBvSYTNc5/fdYYcs3JUn1s=",
        version = "v0.22.0",
    )
    go_repository(
        name = "com_google_cloud_go_gkemulticloud",
        importpath = "cloud.google.com/go/gkemulticloud",
        sum = "h1:3LPj8ro7bPdZ0BlQJwby25G32rRezGLRdhNM6U3AxGs=",
        version = "v1.12.0",
    )
    go_repository(
        name = "com_google_cloud_go_gsuiteaddons",
        importpath = "cloud.google.com/go/gsuiteaddons",
        sum = "h1:kz9DyBk84wZCda2Rdha1MOIf7/9x6R+N+LuH9B3zYFs=",
        version = "v1.12.0",
    )
    go_repository(
        name = "com_google_cloud_go_iam",
        importpath = "cloud.google.com/go/iam",
        sum = "h1:Aki3bX9aHUDKPHfnRJfDcTdVedvy6quGBQcTqx3DRXk=",
        version = "v1.12.0",
    )
    go_repository(
        name = "com_google_cloud_go_iap",
        importpath = "cloud.google.com/go/iap",
        sum = "h1:BDAVJy+juq7cMRumIx9toc4pt1K7zXoZdAI3lDD6D3g=",
        version = "v1.17.0",
    )
    go_repository(
        name = "com_google_cloud_go_ids",
        importpath = "cloud.google.com/go/ids",
        sum = "h1:394LkKIavhv/+oQVMtbXe9WQqCeY9wSFcHktqC9ZkLE=",
        version = "v1.11.0",
    )
    go_repository(
        name = "com_google_cloud_go_iot",
        importpath = "cloud.google.com/go/iot",
        sum = "h1:pyt1EuMFpbV/2BlmYlXr+HBnUNSPKlP9dVK/dpFTu1U=",
        version = "v1.13.0",
    )
    go_repository(
        name = "com_google_cloud_go_kms",
        importpath = "cloud.google.com/go/kms",
        sum = "h1:s+rEluaaZKhLVjrIWG7uNBsnWbiitElzNzFGyp6+nIg=",
        version = "v1.32.0",
    )
    go_repository(
        name = "com_google_cloud_go_language",
        importpath = "cloud.google.com/go/language",
        sum = "h1:q58bL7rmxvw6Q6VHt+wjFsAE1Tj/JkuAEIR4+84rx9U=",
        version = "v1.18.0",
    )
    go_repository(
        name = "com_google_cloud_go_lifesciences",
        importpath = "cloud.google.com/go/lifesciences",
        sum = "h1:T+7N3oVHK8a7BhX6HJLrBW8xEbMV9iY8Mb+aOl7r6+A=",
        version = "v0.16.0",
    )
    go_repository(
        name = "com_google_cloud_go_logging",
        importpath = "cloud.google.com/go/logging",
        sum = "h1:NCqhdVUg3wQ8Cobdf16FDSuTGi3+6+hdSBHrY5TsR6Q=",
        version = "v1.19.0",
    )
    go_repository(
        name = "com_google_cloud_go_longrunning",
        importpath = "cloud.google.com/go/longrunning",
        sum = "h1:WjYH3YHBGCxGJP9M4dWGHBfXr/cFIjMkNgWcJj7/iMM=",
        version = "v1.2.0",
    )
    go_repository(
        name = "com_google_cloud_go_managedidentities",
        importpath = "cloud.google.com/go/managedidentities",
        sum = "h1:ZzWkg3LSrIi9OdzCK3bxutexvw1dDDg9+tOMitbkOXw=",
        version = "v1.13.0",
    )
    go_repository(
        name = "com_google_cloud_go_maps",
        importpath = "cloud.google.com/go/maps",
        sum = "h1:YoTRWohpNKeKxdQynorwHH69NZuUM06divZtqGPYjUA=",
        version = "v1.37.0",
    )
    go_repository(
        name = "com_google_cloud_go_mediatranslation",
        importpath = "cloud.google.com/go/mediatranslation",
        sum = "h1:qqMRAK0mhc9M+lM6pb791oh3DYUbwLcIy58jOyXtadE=",
        version = "v0.13.0",
    )
    go_repository(
        name = "com_google_cloud_go_memcache",
        importpath = "cloud.google.com/go/memcache",
        sum = "h1:l/co1jsmVLkNyl8uMqj1u1HIwd3AhCvGdHbEAIPb6jc=",
        version = "v1.17.0",
    )
    go_repository(
        name = "com_google_cloud_go_metastore",
        importpath = "cloud.google.com/go/metastore",
        sum = "h1:kz0XBkP531MB1IR9jTY6/DB33UFP5Df1vFpPNgW1iPk=",
        version = "v1.20.0",
    )
    go_repository(
        name = "com_google_cloud_go_monitoring",
        importpath = "cloud.google.com/go/monitoring",
        sum = "h1:r/d+JUbyKmJ8b07iznuKfzVzrIXTWxHQ3lBRm3x2LlY=",
        version = "v1.30.0",
    )
    go_repository(
        name = "com_google_cloud_go_networkconnectivity",
        importpath = "cloud.google.com/go/networkconnectivity",
        sum = "h1:ieYJjbUn2L9ZSzEB9+qKgtoD3/l1/qAnId/Z6KE4ZBo=",
        version = "v1.27.0",
    )
    go_repository(
        name = "com_google_cloud_go_networkmanagement",
        importpath = "cloud.google.com/go/networkmanagement",
        sum = "h1:LqYm91vSQk3NM6fHjthseRh4xrDvhXIWQBsQY/Ksus8=",
        version = "v1.30.0",
    )
    go_repository(
        name = "com_google_cloud_go_networksecurity",
        importpath = "cloud.google.com/go/networksecurity",
        sum = "h1:PzJnbVd0sS2+xDR9hnk31NAyqFt9qB8ohLn7XIsKIBA=",
        version = "v0.19.0",
    )
    go_repository(
        name = "com_google_cloud_go_notebooks",
        importpath = "cloud.google.com/go/notebooks",
        sum = "h1:BgxdQcoRSSjiyszWNOTu3fmXG/6Ks4W2sHgjx0BAke0=",
        version = "v1.18.0",
    )
    go_repository(
        name = "com_google_cloud_go_optimization",
        importpath = "cloud.google.com/go/optimization",
        sum = "h1:Cd4DgujNhEMcjfaVLGqbjfDwVDnmS1Bc93cx9TE8ZYw=",
        version = "v1.12.0",
    )
    go_repository(
        name = "com_google_cloud_go_orchestration",
        importpath = "cloud.google.com/go/orchestration",
        sum = "h1:dqf2HzqUi4CaaDiNfv6jVoZjz9lF0HWjyIEgccxPHMk=",
        version = "v1.17.0",
    )
    go_repository(
        name = "com_google_cloud_go_orgpolicy",
        importpath = "cloud.google.com/go/orgpolicy",
        sum = "h1:kpVcE/OsC5aAzHCsAiuQSg3+s6ILzgPTuPZyS7n7ejA=",
        version = "v1.20.0",
    )
    go_repository(
        name = "com_google_cloud_go_osconfig",
        importpath = "cloud.google.com/go/osconfig",
        sum = "h1:r5lzneR9GNixJ96lZ6oIfJK9Cd5sKpnzZJe1LcGZf4w=",
        version = "v1.22.0",
    )
    go_repository(
        name = "com_google_cloud_go_oslogin",
        importpath = "cloud.google.com/go/oslogin",
        sum = "h1:OkccORMcY2XEHwN5+AP0XH/SCUNBl39GGjaQuSCVcIw=",
        version = "v1.18.0",
    )
    go_repository(
        name = "com_google_cloud_go_phishingprotection",
        importpath = "cloud.google.com/go/phishingprotection",
        sum = "h1:6WJF1z3Ie8fZAiCRpeSGB81h6s86XWS2SBC9iQUZxX8=",
        version = "v0.13.0",
    )
    go_repository(
        name = "com_google_cloud_go_policytroubleshooter",
        importpath = "cloud.google.com/go/policytroubleshooter",
        sum = "h1:V3CcqzY6FhY2c2idRGZv0Ef4mFinLnXnfDFtAhxyTfM=",
        version = "v1.16.0",
    )
    go_repository(
        name = "com_google_cloud_go_privatecatalog",
        importpath = "cloud.google.com/go/privatecatalog",
        sum = "h1:SBhcVNVhgEGT1edcfSFY/KSiy2xSQTSQUpY12f10rmw=",
        version = "v0.16.0",
    )
    go_repository(
        name = "com_google_cloud_go_pubsub",
        importpath = "cloud.google.com/go/pubsub",
        sum = "h1:XOaCejsqX7EEtUdQz+WPag66wWsUUGliyCOfGPKfo90=",
        version = "v1.51.0",
    )
    go_repository(
        name = "com_google_cloud_go_pubsub_v2",
        importpath = "cloud.google.com/go/pubsub/v2",
        sum = "h1:8pjR0id+GTB+krKx5G6AGJoYrHog58w2Q89PCOrfM64=",
        version = "v2.6.0",
    )
    go_repository(
        name = "com_google_cloud_go_pubsublite",
        importpath = "cloud.google.com/go/pubsublite",
        sum = "h1:kqNk9e0gt0/Eu6lChjDlgu4B8k/399FZFpAU4LNq4Co=",
        version = "v1.10.0",
    )
    go_repository(
        name = "com_google_cloud_go_recaptchaenterprise_v2",
        importpath = "cloud.google.com/go/recaptchaenterprise/v2",
        sum = "h1:9qDHoQtUIZ8FDN5uFAfnNz8jyyXdpAJf6rMN4ZMmMcU=",
        version = "v2.26.0",
    )
    go_repository(
        name = "com_google_cloud_go_recommendationengine",
        importpath = "cloud.google.com/go/recommendationengine",
        sum = "h1:G0tmQXp67YDaduvde7AX+cIdTtDCo2PBMeqv8SBOpnw=",
        version = "v0.15.0",
    )
    go_repository(
        name = "com_google_cloud_go_recommender",
        importpath = "cloud.google.com/go/recommender",
        sum = "h1:fJ6oO/7Ta/yzfRHcuJUln9iCeo6FDb5yIi2L7eLldzs=",
        version = "v1.19.0",
    )
    go_repository(
        name = "com_google_cloud_go_redis",
        importpath = "cloud.google.com/go/redis",
        sum = "h1:lV6xFUF7t3fm2ZRFjB4IVNkznG9EVwetYVJYwzmW8pc=",
        version = "v1.24.0",
    )
    go_repository(
        name = "com_google_cloud_go_resourcemanager",
        importpath = "cloud.google.com/go/resourcemanager",
        sum = "h1:aDf2RuuQ1Z+kbHHDzHyJNMHPGRFK+cMYfZIhWRaeRX0=",
        version = "v1.16.0",
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
        sum = "h1:u27GLYlnNsQrqzcvSlpZK0L6Ig8k3PEuXGkUdNFCOn4=",
        version = "v1.32.0",
    )
    go_repository(
        name = "com_google_cloud_go_run",
        importpath = "cloud.google.com/go/run",
        sum = "h1:U56fxJWdrT+yjo4S/Vrtw5m69NdNL11Cyv9jX2JOi1s=",
        version = "v1.22.0",
    )
    go_repository(
        name = "com_google_cloud_go_scheduler",
        importpath = "cloud.google.com/go/scheduler",
        sum = "h1:EPdChptxnvCasdMixuu58247qhKMO1iAlrVaQhvuRyE=",
        version = "v1.16.0",
    )
    go_repository(
        name = "com_google_cloud_go_secretmanager",
        importpath = "cloud.google.com/go/secretmanager",
        sum = "h1:GjE3NoyFXo7ipRPy26PMmg4oRX1Ra8fswH45r16rWV0=",
        version = "v1.20.0",
    )
    go_repository(
        name = "com_google_cloud_go_security",
        importpath = "cloud.google.com/go/security",
        sum = "h1:xp/htSPHpyqo225Ju6PjNS9jxfcOEQMxnF229Oze62g=",
        version = "v1.26.0",
    )
    go_repository(
        name = "com_google_cloud_go_securitycenter",
        importpath = "cloud.google.com/go/securitycenter",
        sum = "h1:k5uzLtTSxnh4fwG+8nIcuwDjEE3J+6xe2wcZ72cys7s=",
        version = "v1.45.0",
    )
    go_repository(
        name = "com_google_cloud_go_servicedirectory",
        importpath = "cloud.google.com/go/servicedirectory",
        sum = "h1:yrohWkwM8t5JfEFCmmlyksKnpMpvxM8XbRpwg+yIo64=",
        version = "v1.17.0",
    )
    go_repository(
        name = "com_google_cloud_go_shell",
        importpath = "cloud.google.com/go/shell",
        sum = "h1:vJ/g4BXCwBRMmUxHoNjdbE0DRh5Aysw/+vrZ/aPl0yc=",
        version = "v1.13.0",
    )
    go_repository(
        name = "com_google_cloud_go_spanner",
        importpath = "cloud.google.com/go/spanner",
        sum = "h1:XwXfcZ0kc1NT9Uu2IsThFiWtYptB+WgLn/KZEZcyzRg=",
        version = "v1.91.0",
    )
    go_repository(
        name = "com_google_cloud_go_speech",
        importpath = "cloud.google.com/go/speech",
        sum = "h1:ZcmRKcY02vkF3o/Eaa9yI3i1baNFTF4ctyDZVRfnNS0=",
        version = "v1.36.0",
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
        sum = "h1:jh+3SegMUulEhuJbkADCv8lmfRC8/pCpxja4oex2Ns8=",
        version = "v1.19.0",
    )
    go_repository(
        name = "com_google_cloud_go_talent",
        importpath = "cloud.google.com/go/talent",
        sum = "h1:dUNSwgBUFlDEMHJpmntEvDeLkPQAoEzBrvDbNfS/UeQ=",
        version = "v1.14.0",
    )
    go_repository(
        name = "com_google_cloud_go_texttospeech",
        importpath = "cloud.google.com/go/texttospeech",
        sum = "h1:xGqQv2LB4nttaBNMepCm19EQ0Sk0YdjzPvklDU+LYsI=",
        version = "v1.22.0",
    )
    go_repository(
        name = "com_google_cloud_go_tpu",
        importpath = "cloud.google.com/go/tpu",
        sum = "h1:u2f8Lhezy6piu0BB3+4pTLw9xyOeU5n8PCgWjbjCitk=",
        version = "v1.14.0",
    )
    go_repository(
        name = "com_google_cloud_go_trace",
        importpath = "cloud.google.com/go/trace",
        sum = "h1:GmQovzFc5F0CNfl0VLgL64aoTtu7xsM0YajW2GlG9+E=",
        version = "v1.16.0",
    )
    go_repository(
        name = "com_google_cloud_go_translate",
        importpath = "cloud.google.com/go/translate",
        sum = "h1:2dxWfGMG7y5vobxD2pPyngnfPPIk4MYdNB9QK5sEbuc=",
        version = "v1.18.0",
    )
    go_repository(
        name = "com_google_cloud_go_video",
        importpath = "cloud.google.com/go/video",
        sum = "h1:zPA41kNscRrJeRo/xD73TjNCFAnRDnP9ViNPuWARVeg=",
        version = "v1.33.0",
    )
    go_repository(
        name = "com_google_cloud_go_videointelligence",
        importpath = "cloud.google.com/go/videointelligence",
        sum = "h1:cd+s4jLMS59xiYLT6eQ6PVKP8R0PskgifXDa7YwMUUU=",
        version = "v1.17.0",
    )
    go_repository(
        name = "com_google_cloud_go_vision_v2",
        importpath = "cloud.google.com/go/vision/v2",
        sum = "h1:aTR1vj4++WtS9HD6YdGuoaYygMTJ873WaoV9sYjlQCc=",
        version = "v2.15.0",
    )
    go_repository(
        name = "com_google_cloud_go_vmmigration",
        importpath = "cloud.google.com/go/vmmigration",
        sum = "h1:OPGMxx73owRHMve5B4L6W+0KQUwlunyE9kfbYzFHqug=",
        version = "v1.16.0",
    )
    go_repository(
        name = "com_google_cloud_go_vmwareengine",
        importpath = "cloud.google.com/go/vmwareengine",
        sum = "h1:vdSutihjzH5iYEFlN+1x+gDm4xzj6Md2yMzXADKI5cQ=",
        version = "v1.9.0",
    )
    go_repository(
        name = "com_google_cloud_go_vpcaccess",
        importpath = "cloud.google.com/go/vpcaccess",
        sum = "h1:fFa1i67AjC/3MuG5GPlvx36sIMeffpsIs3kI1C7S3Jk=",
        version = "v1.14.0",
    )
    go_repository(
        name = "com_google_cloud_go_webrisk",
        importpath = "cloud.google.com/go/webrisk",
        sum = "h1:t2WMlBo3xXxH6k5kTSkyIYUl24AHgmBI76F/5OqAPRY=",
        version = "v1.17.0",
    )
    go_repository(
        name = "com_google_cloud_go_websecurityscanner",
        importpath = "cloud.google.com/go/websecurityscanner",
        sum = "h1:iV+hAXeEo8kKtFWDZWPp9Z0fsziZnuOt4nCCYp68RP0=",
        version = "v1.12.0",
    )
    go_repository(
        name = "com_google_cloud_go_workflows",
        importpath = "cloud.google.com/go/workflows",
        sum = "h1:qROpn1zRDdeIjUP9w9O3MHXJEg8v/pzVMicAcw9vLvQ=",
        version = "v1.20.0",
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
        sum = "h1:62yY3dT7/ShwOxzA0RsKRgshBmfElKI4d/Myu2OxDFU=",
        version = "v1.43.0",
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
        sum = "h1:JjwHmHpA4iZ3wBxluu2fbbE7j4kqlE8jXyAyPXH7HqU=",
        version = "v1.44.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_metric",
        build_directives = [
            "gazelle:resolve go go.opentelemetry.io/otel/attribute @io_opentelemetry_go_otel//attribute:go_default_library",
        ],
        importpath = "go.opentelemetry.io/otel/metric",
        sum = "h1:1w0gILTcHdr3YI+ixLyjemwrVnsMURbTZFrSYCdDdmc=",
        version = "v1.44.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_sdk",
        importpath = "go.opentelemetry.io/otel/sdk",
        sum = "h1:nHYwb9lK+fJPU/dnT6s7W7Z8itMWyqrnVfbheVYrZ58=",
        version = "v1.44.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_sdk_metric",
        importpath = "go.opentelemetry.io/otel/sdk/metric",
        sum = "h1:3LlKgI+VjbVsjNRFZJZAJ30WjXC5VkNRks6si09iEfI=",
        version = "v1.44.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_trace",
        build_directives = [
            "gazelle:resolve go go.opentelemetry.io/otel/attribute @io_opentelemetry_go_otel//attribute:go_default_library",
        ],
        importpath = "go.opentelemetry.io/otel/trace",
        sum = "h1:jxF5CsGYCe74MCRx2X4g7WsY/VBKRqqpNvXlX/6gtIk=",
        version = "v1.44.0",
    )
    go_repository(
        name = "org_golang_google_api",
        build_directives = [
            "gazelle:resolve go go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp @io_opentelemetry_go_contrib_instrumentation_net_http_otelhttp//:go_default_library",
            "gazelle:resolve go go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc @io_opentelemetry_go_contrib_instrumentation_google_golang_org_grpc_otelgrpc//:go_default_library",
        ],
        importpath = "google.golang.org/api",
        sum = "h1:glhO/J88obKP5I269W3hB73dvBKrjU56ZfmNlNXpgTU=",
        version = "v0.288.0",
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
        sum = "h1:CzairfKVa2xaI2sTm9v7YLXGCI9kTyi6Rw1jWcYWGnE=",
        version = "v0.0.0-20260715203245-bcc9394bd25e",
    )
    go_repository(
        name = "org_golang_google_genproto_googleapis_api",
        importpath = "google.golang.org/genproto/googleapis/api",
        sum = "h1:jQ9p21COKWjP3VwuFrNRiiOTMh3mPpN45R7SLrH/HUU=",
        version = "v0.0.0-20260630182238-925bb5da69e7",
    )
    go_repository(
        name = "org_golang_google_genproto_googleapis_bytestream",
        importpath = "google.golang.org/genproto/googleapis/bytestream",
        sum = "h1:OQRcQe8VeJabR6sVUO9FYZ9kwNFZhBAH/g0bYlaM/4g=",
        version = "v0.0.0-20260630182238-925bb5da69e7",
    )
    go_repository(
        name = "org_golang_google_genproto_googleapis_rpc",
        importpath = "google.golang.org/genproto/googleapis/rpc",
        sum = "h1:JRXr1W/pX290ntNGMipfzkE7fYo3x22P78fqXxquSNw=",
        version = "v0.0.0-20260715203245-bcc9394bd25e",
    )
    go_repository(
        name = "org_golang_google_grpc",
        importpath = "google.golang.org/grpc",
        sum = "h1:vguDnZUPjE26w09A63VoxZPnvPjB5Riyc0mkXPFmAIU=",
        version = "v1.82.0",
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
        sum = "h1:QZ4Muo8THX6CizN2vPPd5fBGHyogrdK9fG4wLPFUsto=",
        version = "v0.53.0",
    )
    go_repository(
        name = "org_golang_x_mod",
        importpath = "golang.org/x/mod",
        sum = "h1:vF1DjpVEshcIqoEaauuHebaLk1O1forxjxBaVn884JQ=",
        version = "v0.37.0",
    )
    go_repository(
        name = "org_golang_x_net",
        importpath = "golang.org/x/net",
        sum = "h1:Rw8j/hFzGvJUZwNBXnAtf5sVDVt+65SK2C7IxCxZt5o=",
        version = "v0.56.0",
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
        sum = "h1:SZjpbeLmrCk4xhRSZFNZW5gFUeCeFgjekvI/+gfScek=",
        version = "v0.22.0",
    )
    go_repository(
        name = "org_golang_x_sys",
        importpath = "golang.org/x/sys",
        sum = "h1:o7XGOvZQCADBQQ4Y7VNq2dRWQR7JmOUW8Kxx4ZsNgWs=",
        version = "v0.47.0",
    )
    go_repository(
        name = "org_golang_x_term",
        importpath = "golang.org/x/term",
        sum = "h1:0rLvDRCtNj0gZkyIXhCyOb2OAzEhLVqc4B+hrsBhrmc=",
        version = "v0.44.0",
    )
    go_repository(
        name = "org_golang_x_text",
        importpath = "golang.org/x/text",
        sum = "h1:Ub2Z6/xjgF1WrYQz2nuITOEegKFtiIy+rieRJ5lHZKs=",
        version = "v0.40.0",
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
        sum = "h1:7Kn5x/d1svx/PzryTsqeoZN4TZwqeH5pGWjefhLi/1Q=",
        version = "v0.47.0",
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
