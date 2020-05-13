workspace(name = "com_googleapis_gapic_generator_go")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

# We depend on a release candidate because it contains the proto3_optional
# feature.
#
# TODO(ndietz) Replace this with stable 3.12.0 version when it is released so as
# to not depend on a release candidate.
http_archive(
    name = "com_google_protobuf",
    strip_prefix = "protobuf-3.12.0-rc2",
    urls = ["https://github.com/protocolbuffers/protobuf/archive/v3.12.0-rc2.tar.gz"],
    sha256 = "afaa4f65e7e97adb10b32b7c699b7b6be4090912b471028ef0f40ccfb271f96a",
)

load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

protobuf_deps()

# Depend on this rules_go commit because it contains unreleased patches that
# allow us to override com_github_golang_protobuf to v1.4.1.
#
# TODO(ndietz) replace this with the next version (v0.23.0) of rules_go
# that contains github.com/golang/protobuf v1.4.1 and google.golang.org/protobuf
# v1.22.0: https://github.com/bazelbuild/rules_go/issues/2471
http_archive(
    name = "io_bazel_rules_go",
    # master as of 5pm PST on 5/04/2020
    urls = [
        "https://github.com/bazelbuild/rules_go/archive/695da5906684ab96e558c3159f7d88a20a6a1703.zip",
    ],
    strip_prefix = "rules_go-695da5906684ab96e558c3159f7d88a20a6a1703",
    sha256 = "b5f8747b60f5b3d02125b5c9a69ed5ad59a46fce637b3073ff0576f152e2354c",
)

load("@io_bazel_rules_go//go:deps.bzl", "go_rules_dependencies", "go_register_toolchains")

http_archive(
    name = "bazel_gazelle",
    urls = [
        "https://storage.googleapis.com/bazel-mirror/github.com/bazelbuild/bazel-gazelle/releases/download/v0.21.0/bazel-gazelle-v0.21.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.21.0/bazel-gazelle-v0.21.0.tar.gz",
    ],
    sha256 = "bfd86b3cbe855d6c16c6fce60d76bd51f5c8dbc9cfcaef7a2bb5c1aafd0710e8",
)

# gazelle:repo bazel_gazelle
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

# These golang protobuf dependencies are added here to override those provided
# by rules_go. This is done to capture the proto3_optional changes necessary
# to build the go_gapic plugin.
#
# TODO(ndietz) remove these overrides once rules_go v0.23.0 is released
# https://github.com/bazelbuild/rules_go/issues/2471
go_repository(
    name = "org_golang_google_protobuf",
    build_file_proto_mode = "disable_global",
    importpath = "google.golang.org/protobuf",
    sum = "h1:cJv5/xdbk1NnMPR1VP9+HU6gupuG9MLBoH1r6RHZ2MY=",
    version = "v1.22.0",
)

go_repository(
    name = "com_github_golang_protobuf",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/golang/protobuf",
    patch_args = ["-p1"],
    patches = ["@io_bazel_rules_go//third_party:com_github_golang_protobuf-extras.patch"],
    sum = "h1:ZFgWrT+bLgsYPirOnRfKLYJLvssAegOj/hgyMFdJZe0=",
    version = "v1.4.1",
)

# These repository macros are invoked after the golang protobuf dependencies
# so as to override them with the desired version for this repo.
#
# TODO(ndietz) move these back to where they were next to their corresponding
# dependency
go_rules_dependencies()

go_register_toolchains()

gazelle_dependencies()

load("//:repositories.bzl", "com_googleapis_gapic_generator_go_repositories")

com_googleapis_gapic_generator_go_repositories()
