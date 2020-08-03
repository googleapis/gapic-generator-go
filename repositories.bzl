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

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@bazel_gazelle//:deps.bzl", "go_repository")

def com_googleapis_gapic_generator_go_repositories():
    _maybe(
        http_archive,
        name = "com_google_api_codegen",
        strip_prefix = "gapic-generator-a79e9ea3fcf686a80d92461a4788c5bcf55cea5a",
        urls = ["https://github.com/googleapis/gapic-generator/archive/a79e9ea3fcf686a80d92461a4788c5bcf55cea5a.zip"],
        sha256 = "c6a13fd221189458ad9eeb1de1f40e21bd80f0063bf05b9fa243722c18577f17",
    )
    _maybe(
        go_repository,
        name = "co_honnef_go_tools",
        importpath = "honnef.co/go/tools",
        sum = "h1:UoveltGrhghAA7ePc+e+QYDHXrBps2PqFZiHkGR/xK8=",
        version = "v0.0.1-2020.1.4",
    )
    _maybe(
        go_repository,
        name = "com_github_burntsushi_toml",
        importpath = "github.com/BurntSushi/toml",
        sum = "h1:WXkYYl6Yr3qBf1K79EBnL4mak0OimBfB0XUf9Vl28OQ=",
        version = "v0.3.1",
    )
    _maybe(
        go_repository,
        name = "com_github_census_instrumentation_opencensus_proto",
        importpath = "github.com/census-instrumentation/opencensus-proto",
        sum = "h1:glEXhBS5PSLLv4IXzLA5yPRVX4bilULVyxxbrfOtDAk=",
        version = "v0.2.1",
    )
    _maybe(
        go_repository,
        name = "com_github_client9_misspell",
        importpath = "github.com/client9/misspell",
        sum = "h1:ta993UF76GwbvJcIo3Y68y/M3WxlpEHPWIGDkJYwzJI=",
        version = "v0.3.4",
    )
    _maybe(
        go_repository,
        name = "com_github_envoyproxy_go_control_plane",
        importpath = "github.com/envoyproxy/go-control-plane",
        sum = "h1:4cmBvAEBNJaGARUEs3/suWRyfyBfhf7I60WBZq+bv2w=",
        version = "v0.9.1-0.20191026205805-5f8ba28d4473",
    )
    _maybe(
        go_repository,
        name = "com_github_envoyproxy_protoc_gen_validate",
        importpath = "github.com/envoyproxy/protoc-gen-validate",
        sum = "h1:EQciDnbrYxy13PgWoY8AqoxGiPrpgBZ1R8UNe3ddc+A=",
        version = "v0.1.0",
    )
    _maybe(
        go_repository,
        name = "com_github_golang_glog",
        importpath = "github.com/golang/glog",
        sum = "h1:VKtxabqXZkF25pY9ekfRL6a582T4P37/31XEstQ5p58=",
        version = "v0.0.0-20160126235308-23def4e6c14b",
    )
    _maybe(
        go_repository,
        name = "com_github_golang_mock",
        importpath = "github.com/golang/mock",
        sum = "h1:G5FRp8JnTd7RQH5kemVNlMeyXQAztQ3mOWV95KxsXH8=",
        version = "v1.1.1",
    )
    _maybe(
        go_repository,
        name = "com_github_golang_protobuf",
        importpath = "github.com/golang/protobuf",
        sum = "h1:+Z5KGCizgyZCbGh1KZqA0fcLLkwbsjIzS4aV2v7wJX0=",
        version = "v1.4.2",
    )
    _maybe(
        go_repository,
        name = "com_github_google_go_cmp",
        importpath = "github.com/google/go-cmp",
        sum = "h1:/QaMHBdZ26BB3SSst0Iwl10Epc+xhTquomWX0oZEB6w=",
        version = "v0.5.0",
    )
    _maybe(
        go_repository,
        name = "com_github_google_renameio",
        importpath = "github.com/google/renameio",
        sum = "h1:GOZbcHa3HfsPKPlmyPyN2KEohoMXOhdMbHrvbpl2QaA=",
        version = "v0.1.0",
    )
    _maybe(
        go_repository,
        name = "com_github_gordonklaus_ineffassign",
        importpath = "github.com/gordonklaus/ineffassign",
        sum = "h1:vc7Dmrk4JwS0ZPS6WZvWlwDflgDTA26jItmbSj83nug=",
        version = "v0.0.0-20200309095847-7953dde2c7bf",
    )
    _maybe(
        go_repository,
        name = "com_github_jhump_protoreflect",
        importpath = "github.com/jhump/protoreflect",
        sum = "h1:qJ7piXPrjP3mDrfHf5ATkxfLix8ANs226vpo0aACOn0=",
        version = "v1.7.0",
        # Manually added in order to disable testproto BUILD file generation.
        build_file_proto_mode = "disable",
    )
    _maybe(
        go_repository,
        name = "com_github_kisielk_gotool",
        importpath = "github.com/kisielk/gotool",
        sum = "h1:AV2c/EiW3KqPNT9ZKl07ehoAGi4C5/01Cfbblndcapg=",
        version = "v1.0.0",
    )
    _maybe(
        go_repository,
        name = "com_github_kr_pretty",
        importpath = "github.com/kr/pretty",
        sum = "h1:L/CwN0zerZDmRFUapSPitk6f+Q3+0za1rQkzVuMiMFI=",
        version = "v0.1.0",
    )
    _maybe(
        go_repository,
        name = "com_github_kr_pty",
        importpath = "github.com/kr/pty",
        sum = "h1:VkoXIwSboBpnk99O/KFauAEILuNHv5DVFKZMBN/gUgw=",
        version = "v1.1.1",
    )
    _maybe(
        go_repository,
        name = "com_github_kr_text",
        importpath = "github.com/kr/text",
        sum = "h1:45sCR5RtlFHMR4UwH9sdQ5TC8v0qDQCHnXt+kaKSTVE=",
        version = "v0.1.0",
    )
    _maybe(
        go_repository,
        name = "com_github_nishanths_predeclared",
        importpath = "github.com/nishanths/predeclared",
        sum = "h1:3f0nxAmdj/VoCGN/ijdMy7bj6SBagaqYg1B0hu8clMA=",
        version = "v0.0.0-20200524104333-86fad755b4d3",
    )
    _maybe(
        go_repository,
        name = "com_github_prometheus_client_model",
        importpath = "github.com/prometheus/client_model",
        sum = "h1:gQz4mCbXsO+nc9n1hCxHcGA3Zx3Eo+UHZoInFGUIXNM=",
        version = "v0.0.0-20190812154241-14fe0d1b01d4",
    )
    _maybe(
        go_repository,
        name = "com_github_rogpeppe_go_internal",
        importpath = "github.com/rogpeppe/go-internal",
        sum = "h1:RR9dF3JtopPvtkroDZuVD7qquD0bnHlKSqaQhgwt8yk=",
        version = "v1.3.0",
    )
    _maybe(
        go_repository,
        name = "com_github_russross_blackfriday",
        importpath = "github.com/russross/blackfriday",
        sum = "h1:cBXrhZNUf9C+La9/YpS+UHpUT8YD6Td9ZMSU9APFcsk=",
        version = "v2.0.0+incompatible",
    )
    _maybe(
        go_repository,
        name = "com_github_shurcool_sanitized_anchor_name",
        importpath = "github.com/shurcooL/sanitized_anchor_name",
        sum = "h1:PdmoCO6wvbs+7yrJyMORt4/BmY5IYyJwS/kOiWx8mHo=",
        version = "v1.0.0",
    )
    _maybe(
        go_repository,
        name = "com_github_yuin_goldmark",
        importpath = "github.com/yuin/goldmark",
        sum = "h1:nqDD4MMMQA0lmWq03Z2/myGPYLQoXtmi0rGVs95ntbo=",
        version = "v1.1.27",
    )
    _maybe(
        go_repository,
        name = "com_gitlab_golang_commonmark_html",
        importpath = "gitlab.com/golang-commonmark/html",
        sum = "h1:K+bMSIx9A7mLES1rtG+qKduLIXq40DAzYHtb0XuCukA=",
        version = "v0.0.0-20191124015941-a22733972181",
    )
    _maybe(
        go_repository,
        name = "com_gitlab_golang_commonmark_linkify",
        importpath = "gitlab.com/golang-commonmark/linkify",
        sum = "h1:oYrL81N608MLZhma3ruL8qTM4xcpYECGut8KSxRY59g=",
        version = "v0.0.0-20191026162114-a0c2df6c8f82",
    )
    _maybe(
        go_repository,
        name = "com_gitlab_golang_commonmark_markdown",
        importpath = "gitlab.com/golang-commonmark/markdown",
        sum = "h1:HsZm6XaTpEgZiZqcXZkUbG6BNtSZE3XyCTfo52YBoDY=",
        version = "v0.0.0-20191127184510-91b5b3c99c19",
    )
    _maybe(
        go_repository,
        name = "com_gitlab_golang_commonmark_mdurl",
        importpath = "gitlab.com/golang-commonmark/mdurl",
        sum = "h1:qqjvoVXdWIcZCLPMlzgA7P9FZWdPGPvP/l3ef8GzV6o=",
        version = "v0.0.0-20191124015652-932350d1cb84",
    )
    _maybe(
        go_repository,
        name = "com_gitlab_golang_commonmark_puny",
        importpath = "gitlab.com/golang-commonmark/puny",
        sum = "h1:Wku8eEdeJqIOFHtrfkYUByc4bCaTeA6fL0UJgfEiFMI=",
        version = "v0.0.0-20191124015043-9f83538fa04f",
    )
    _maybe(
        go_repository,
        name = "com_gitlab_opennota_wd",
        importpath = "gitlab.com/opennota/wd",
        sum = "h1:uPZaMiz6Sz0PZs3IZJWpU5qHKGNy///1pacZC9txiUI=",
        version = "v0.0.0-20180912061657-c5d65f63c638",
    )
    _maybe(
        go_repository,
        name = "com_google_cloud_go",
        importpath = "cloud.google.com/go",
        sum = "h1:e0WKqKTd5BnrG8aKH3J3h+QvEIQtSUcf2n5UZ5ZgLtQ=",
        version = "v0.26.0",
    )
    _maybe(
        go_repository,
        name = "in_gopkg_check_v1",
        importpath = "gopkg.in/check.v1",
        sum = "h1:qIbj1fsPNlZgppZ+VLlY7N33q108Sa+fhmuc+sWQYwY=",
        version = "v1.0.0-20180628173108-788fd7840127",
    )
    _maybe(
        go_repository,
        name = "in_gopkg_errgo_v2",
        importpath = "gopkg.in/errgo.v2",
        sum = "h1:0vLT13EuvQ0hNvakwLuFZ/jYrLp5F3kcWHXdRggjCE8=",
        version = "v2.1.0",
    )
    _maybe(
        go_repository,
        name = "in_gopkg_yaml_v2",
        importpath = "gopkg.in/yaml.v2",
        sum = "h1:clyUAQHOM3G0M3f5vQj7LuJrETvjVot3Z5el9nffUtU=",
        version = "v2.3.0",
    )
    _maybe(
        go_repository,
        name = "org_golang_google_appengine",
        importpath = "google.golang.org/appengine",
        sum = "h1:/wp5JvzpHIxhs/dumFmF7BXTf3Z+dd4uXta4kVyO508=",
        version = "v1.4.0",
    )
    _maybe(
        go_repository,
        name = "org_golang_google_genproto",
        importpath = "google.golang.org/genproto",
        sum = "h1:L8iMEkPD/F6k/2Ui4/kfwRn5edLwHtfuj/+KK0nlU/E=",
        version = "v0.0.0-20200706141556-5779274c8e96",
    )
    _maybe(
        go_repository,
        name = "org_golang_google_grpc",
        importpath = "google.golang.org/grpc",
        sum = "h1:rRYRFMVgRv6E0D70Skyfsr28tDXIuuPZyWGMPdMcnXg=",
        version = "v1.27.0",
    )
    _maybe(
        go_repository,
        name = "org_golang_google_protobuf",
        importpath = "google.golang.org/protobuf",
        sum = "h1:UhZDfRO8JRQru4/+LlLE0BRKGF8L+PICnvYZmx/fEGA=",
        version = "v1.24.0",
    )
    _maybe(
        go_repository,
        name = "org_golang_x_crypto",
        importpath = "golang.org/x/crypto",
        sum = "h1:ObdrDkeb4kJdCP557AjRjq69pTHfNouLtWZG7j9rPN8=",
        version = "v0.0.0-20191011191535-87dc89f01550",
    )
    _maybe(
        go_repository,
        name = "org_golang_x_exp",
        importpath = "golang.org/x/exp",
        sum = "h1:c2HOrn5iMezYjSlGPncknSEr/8x5LELb/ilJbXi9DEA=",
        version = "v0.0.0-20190121172915-509febef88a4",
    )
    _maybe(
        go_repository,
        name = "org_golang_x_lint",
        importpath = "golang.org/x/lint",
        sum = "h1:XQyxROzUlZH+WIQwySDgnISgOivlhjIEwaQaJEJrrN0=",
        version = "v0.0.0-20190313153728-d0100b6bd8b3",
    )
    _maybe(
        go_repository,
        name = "org_golang_x_mod",
        importpath = "golang.org/x/mod",
        sum = "h1:KU7oHjnv3XNWfa5COkzUifxZmxp1TyI7ImMXqFxLwvQ=",
        version = "v0.2.0",
    )
    _maybe(
        go_repository,
        name = "org_golang_x_net",
        importpath = "golang.org/x/net",
        sum = "h1:0mm1VjtFUOIlE1SbDlwjYaDxZVDP2S5ou6y0gSgXHu8=",
        version = "v0.0.0-20200226121028-0de0cce0169b",
    )
    _maybe(
        go_repository,
        name = "org_golang_x_oauth2",
        importpath = "golang.org/x/oauth2",
        sum = "h1:vEDujvNQGv4jgYKudGeI/+DAX4Jffq6hpD55MmoEvKs=",
        version = "v0.0.0-20180821212333-d2e6202438be",
    )
    _maybe(
        go_repository,
        name = "org_golang_x_sync",
        importpath = "golang.org/x/sync",
        sum = "h1:vcxGaoTs7kV8m5Np9uUNQin4BrLOthgV7252N8V+FwY=",
        version = "v0.0.0-20190911185100-cd5d95a43a6e",
    )
    _maybe(
        go_repository,
        name = "org_golang_x_sys",
        importpath = "golang.org/x/sys",
        sum = "h1:+R4KGOnez64A81RvjARKc4UT5/tI9ujCIVX+P5KiHuI=",
        version = "v0.0.0-20190412213103-97732733099d",
    )
    _maybe(
        go_repository,
        name = "org_golang_x_text",
        importpath = "golang.org/x/text",
        sum = "h1:tW2bmiBqwgJj/UpqtC8EpXEZVYOwU0yG4iWbprSVAcs=",
        version = "v0.3.2",
    )
    _maybe(
        go_repository,
        name = "org_golang_x_tools",
        importpath = "golang.org/x/tools",
        sum = "h1:SjQ2+AKWgZLc1xej6WSzL+Dfs5Uyd5xcZH1mGC411IA=",
        version = "v0.0.0-20200522201501-cb1345f3a375",
    )
    _maybe(
        go_repository,
        name = "org_golang_x_xerrors",
        importpath = "golang.org/x/xerrors",
        sum = "h1:E7g+9GITq07hpfrRu66IVDexMakfv52eLZ2CXBWiKr4=",
        version = "v0.0.0-20191204190536-9bdfabe68543",
    )

def _maybe(repo_rule, name, strip_repo_prefix = "", **kwargs):
    if not name.startswith(strip_repo_prefix):
        return
    repo_name = name[len(strip_repo_prefix):]
    if repo_name in native.existing_rules():
        return
    repo_rule(name = repo_name, **kwargs)
