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
        sum = "h1:l2Nfbl2GPXdWorv+dT2XfinX2jOOw4zv1VhLstx+6rE=",
        version = "v0.36.0",
    )
    _maybe(
        go_repository,
        name = "org_golang_x_oauth2",
        importpath = "golang.org/x/oauth2",
        sum = "h1:ld7aEMNHoBnnDAX15v1T6z31v8HwR2A9FYOuAhWqkwc=",
        version = "v0.0.0-20200902213428-5d25da1a8d43",
    )
    _maybe(
        go_repository,
        name = "com_github_google_go_cmp",
        importpath = "github.com/google/go-cmp/cmp",
        sum = "h1:L8R9j+yAqZuZjsqh/z+F1NCffTKKLShY6zXTItVIZ8M=",
        version = "v0.5.4",
    )
    _maybe(
        go_repository,
        name = "io_opencensus_go",
        importpath = "go.opencensus.io",
        sum = "h1:dntmOdLpSpHlVqbW5Eay97DelsZHe+55D+xC6i0dDS0=",
        version = "v0.22.5",
    )
    _maybe(
        go_repository,
        name = "com_github_golang_groupcache",
        importpath = "github.com/golang/groupcache",
        sum = "h1:1r7pUrabqp18hOBcwBwiTsbnFeTZHV9eER/QT5JVZxY=",
        version = "v0.0.0-20200121045136-8c9f03a8e57e",
    )
    _maybe(
        go_repository,
        name = "com_github_googleapis_gax_go",
        importpath = "github.com/googleapis/gax-go",
        sum = "h1:9dMLqhaibYONnDRcnHdUs9P8Mw64jLlZTYlDe3leBtQ=",
        version = "v1.0.3",
    )
    # This is a pseudo-cycle, because Go GAPIC's depend on
    # common packages in the same mono-repo, particularly the
    # longrunning package.
    _maybe(
        go_repository,
        name = "com_google_cloud_go",
        importpath = "cloud.google.com/go",
        sum = "h1:ujhG1RejZYi+HYfJNlgBh3j/bVKD8DewM7AkJ5UPyBc=",
        version = "v0.70.0",
        # This is part of a fix for https://github.com/googleapis/gapic-generator-go/issues/387.
        build_extra_args = ["-exclude=longrunning/autogen/info.go"],
    )

def _maybe(repo_rule, name, strip_repo_prefix = "", **kwargs):
    if not name.startswith(strip_repo_prefix):
        return
    repo_name = name[len(strip_repo_prefix):]
    if repo_name in native.existing_rules():
        return
    repo_rule(name = repo_name, **kwargs)
