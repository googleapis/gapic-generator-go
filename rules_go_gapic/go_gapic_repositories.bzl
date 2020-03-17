# Copyright 2020 Google LLC
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

def go_gapic_repositories():
    _maybe(
        go_repository,
        name = "com_github_googleapis_gax_go_v2",
        importpath = "github.com/googleapis/gax-go/v2",
        sum = "h1:sjZBwGj9Jlw33ImPtvFviGYvseOtDM7hkSKB7+Tv3SM=",
        version = "v2.0.5",
    )
    _maybe(
        go_repository,
        name = "org_golang_google_api",
        importpath = "google.golang.org/api",
        sum = "h1:jz2KixHX7EcCPiQrySzPdnYT7DbINAypCqKZ1Z7GM40=",
        version = "v0.20.0",
    )
    _maybe(
        go_repository,
        name = "org_golang_x_oauth2",
        importpath = "golang.org/x/oauth2",
        sum = "h1:TzXSXBo42m9gQenoE3b9BGiEpg5IG2JkU5FkPIawgtw=",
        version = "v0.0.0-20200107190931-bf48bf16ab8d",
    )
    _maybe(
        go_repository,
        name = "com_github_google_go_cmp",
        importpath = "github.com/google/go-cmp/cmp",
        sum = "h1:xsAVV57WRhGj6kEIi8ReJzQlHHqcBYCElAvkovg3B/4=",
        version = "v0.4.0",
    )
    _maybe(
        go_repository,
        name = "io_opencensus_go",
        importpath = "go.opencensus.io",
        sum = "h1:8sGtKOrtQqkN1bp2AtX+misvLIlOmsEsNd+9NIcPEm8=",
        version = "v0.22.3",
    )
    _maybe(
        go_repository,
        name = "com_github_golang_groupcache",
        importpath = "github.com/golang/groupcache",
        sum = "h1:ZgQEtGgCBiWRM39fZuwSd1LwSqqSW0hOdXCYYDX0R3I=",
        version = "v0.0.0-20190702054246-869f871628b6",
    )
    _maybe(
        go_repository,
        name = "com_github_googleapis_gax_go",
        importpath = "github.com/googleapis/gax-go",
        sum = "h1:9dMLqhaibYONnDRcnHdUs9P8Mw64jLlZTYlDe3leBtQ=",
        version = "v1.0.3",
    )

def _maybe(repo_rule, name, strip_repo_prefix = "", **kwargs):
    if not name.startswith(strip_repo_prefix):
        return
    repo_name = name[len(strip_repo_prefix):]
    if repo_name in native.existing_rules():
        return
    repo_rule(name = repo_name, **kwargs)
