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

load("@bazel_gazelle//:deps.bzl", "go_repository")

def com_googleapis_gapic_go_mod():
    go_repository(
        name = "co_honnef_go_tools",
        importpath = "honnef.co/go/tools",
        sum = "h1:XJP7lxbSxWLOMNdBE4B/STaqVy6L73o0knwj2vIlxnw=",
        version = "v0.0.0-20190102054323-c2f93a96b099",
    )
    go_repository(
        name = "com_github_burntsushi_toml",
        importpath = "github.com/BurntSushi/toml",
        sum = "h1:WXkYYl6Yr3qBf1K79EBnL4mak0OimBfB0XUf9Vl28OQ=",
        version = "v0.3.1",
    )
    go_repository(
        name = "com_github_client9_misspell",
        importpath = "github.com/client9/misspell",
        sum = "h1:ta993UF76GwbvJcIo3Y68y/M3WxlpEHPWIGDkJYwzJI=",
        version = "v0.3.4",
    )
    go_repository(
        name = "com_github_golang_commonmark_html",
        importpath = "github.com/golang-commonmark/html",
        sum = "h1:FeNEDxIy7XouGTJKiJ9Ze5vUbcAIW/FRhQbtKBNmEz8=",
        version = "v0.0.0-20180910111043-7d7c804e1d46",
    )
    go_repository(
        name = "com_github_golang_commonmark_linkify",
        importpath = "github.com/golang-commonmark/linkify",
        sum = "h1:TkuRzcq232K5ytXtQ+BPicsjYWZgt/lS6gJ5HqcUifQ=",
        version = "v0.0.0-20180910111149-f05efb453a0e",
    )
    go_repository(
        name = "com_github_golang_commonmark_markdown",
        importpath = "github.com/golang-commonmark/markdown",
        sum = "h1:YaQaotRjMcVth1VzHUEQlD2oeyQAglA7CXdxp9QLvKM=",
        version = "v0.0.0-20180910011815-a8f139058164",
    )
    go_repository(
        name = "com_github_golang_commonmark_mdurl",
        importpath = "github.com/golang-commonmark/mdurl",
        sum = "h1:XkgfhPs5AotQfcu3EfDEjyAUx91KdtjrxHXYGnZJhoU=",
        version = "v0.0.0-20180910110917-8d018c6567d6",
    )
    go_repository(
        name = "com_github_golang_commonmark_puny",
        importpath = "github.com/golang-commonmark/puny",
        sum = "h1:DUgQdQmDg4sk4SfNR+qOkXcopGz36BL02vp/V7WbPQI=",
        version = "v0.0.0-20180910110745-050be392d8b8",
    )
    go_repository(
        name = "com_github_golang_glog",
        importpath = "github.com/golang/glog",
        sum = "h1:VKtxabqXZkF25pY9ekfRL6a582T4P37/31XEstQ5p58=",
        version = "v0.0.0-20160126235308-23def4e6c14b",
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
        sum = "h1:6nsPYzhq5kReh6QImI3k5qWzO4PEbvbIW2cwSfR/6xs=",
        version = "v1.3.2",
    )
    go_repository(
        name = "com_github_google_go_cmp",
        importpath = "github.com/google/go-cmp",
        sum = "h1:Xye71clBPdm5HgqGwUkwhbynsUJZhDbS20FvLhQ2izg=",
        version = "v0.3.1",
    )
    go_repository(
        name = "com_github_jhump_protoreflect",
        importpath = "github.com/jhump/protoreflect",
        sum = "h1:NgpVT+dX71c8hZnxHof2M7QDK7QtohIJ7DYycjnkyfc=",
        version = "v1.5.0",
    )
    go_repository(
        name = "com_github_opennota_wd",
        importpath = "github.com/opennota/wd",
        sum = "h1:cVQhwfBgiKTMAdYPbVeuIiTkdY59qZ3sp5RpyO8CNtg=",
        version = "v0.0.0-20180911144301-b446539ab1e7",
    )
    go_repository(
        name = "com_github_russross_blackfriday",
        importpath = "github.com/russross/blackfriday",
        sum = "h1:cBXrhZNUf9C+La9/YpS+UHpUT8YD6Td9ZMSU9APFcsk=",
        version = "v2.0.0+incompatible",
    )
    go_repository(
        name = "com_github_shurcool_sanitized_anchor_name",
        importpath = "github.com/shurcooL/sanitized_anchor_name",
        sum = "h1:PdmoCO6wvbs+7yrJyMORt4/BmY5IYyJwS/kOiWx8mHo=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_google_cloud_go",
        importpath = "cloud.google.com/go",
        sum = "h1:e0WKqKTd5BnrG8aKH3J3h+QvEIQtSUcf2n5UZ5ZgLtQ=",
        version = "v0.26.0",
    )
    go_repository(
        name = "in_gopkg_check_v1",
        importpath = "gopkg.in/check.v1",
        sum = "h1:yhCVgyC4o1eVCa2tZl7eS0r+SDo693bJlVdllGtEeKM=",
        version = "v0.0.0-20161208181325-20d25e280405",
    )
    go_repository(
        name = "in_gopkg_yaml_v2",
        importpath = "gopkg.in/yaml.v2",
        sum = "h1:ZCJp+EgiOT7lHqUV2J862kp8Qj64Jo6az82+3Td9dZw=",
        version = "v2.2.2",
    )
    go_repository(
        name = "org_golang_google_appengine",
        importpath = "google.golang.org/appengine",
        sum = "h1:/wp5JvzpHIxhs/dumFmF7BXTf3Z+dd4uXta4kVyO508=",
        version = "v1.4.0",
    )
    go_repository(
        name = "org_golang_google_genproto",
        importpath = "google.golang.org/genproto",
        sum = "h1:dnhIt8HlkxWTUDwTelzFnLXLXxbuq1AbZhTc7BDUGTo=",
        version = "v0.0.0-20190819204819-24fa4b261c55",
    )
    go_repository(
        name = "org_golang_google_grpc",
        importpath = "google.golang.org/grpc",
        sum = "h1:cfg4PD8YEdSFnm7qLV4++93WcmhH2nIUhMjhdCvl3j8=",
        version = "v1.19.0",
    )
    go_repository(
        name = "org_golang_x_exp",
        importpath = "golang.org/x/exp",
        sum = "h1:c2HOrn5iMezYjSlGPncknSEr/8x5LELb/ilJbXi9DEA=",
        version = "v0.0.0-20190121172915-509febef88a4",
    )
    go_repository(
        name = "org_golang_x_lint",
        importpath = "golang.org/x/lint",
        sum = "h1:GmgasJE571dBGXS7E282h2rIZj+KvCLV8z5I6QXbKNI=",
        version = "v0.0.0-20190227174305-5b3e6a55c961",
    )
    go_repository(
        name = "org_golang_x_net",
        importpath = "golang.org/x/net",
        sum = "h1:HuTn7WObtcDo9uEEU7rEqL0jYthdXAmZ6PP+meazmaU=",
        version = "v0.0.0-20190213061140-3a22650c66bd",
    )
    go_repository(
        name = "org_golang_x_oauth2",
        importpath = "golang.org/x/oauth2",
        sum = "h1:vEDujvNQGv4jgYKudGeI/+DAX4Jffq6hpD55MmoEvKs=",
        version = "v0.0.0-20180821212333-d2e6202438be",
    )
    go_repository(
        name = "org_golang_x_sync",
        importpath = "golang.org/x/sync",
        sum = "h1:Bl/8QSvNqXvPGPGXa2z5xUTmV7VDcZyvRZ+QQXkXTZQ=",
        version = "v0.0.0-20181108010431-42b317875d0f",
    )
    go_repository(
        name = "org_golang_x_sys",
        importpath = "golang.org/x/sys",
        sum = "h1:Ve1ORMCxvRmSXBwJK+t3Oy+V2vRW2OetUQBq4rJIkZE=",
        version = "v0.0.0-20180830151530-49385e6e1522",
    )
    go_repository(
        name = "org_golang_x_text",
        importpath = "golang.org/x/text",
        sum = "h1:g61tztE5qeGQ89tm6NTjjM9VPIm088od1l6aSorWRWg=",
        version = "v0.3.0",
    )
    go_repository(
        name = "org_golang_x_tools",
        importpath = "golang.org/x/tools",
        sum = "h1:vamGzbGri8IKo20MQncCuljcQ5uAO6kaCeawQPVblAI=",
        version = "v0.0.0-20190226205152-f727befe758c",
    )
