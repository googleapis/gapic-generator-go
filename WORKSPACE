workspace(name = "com_googleapis_gapic_generator_go")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

# Workaround for https://github.com/bazelbuild/bazel-gazelle/issues/1285. Ideally,
# we can remove this if gazelle ships a fix since we didn't need it before.
http_archive(
    name = "bazel_skylib",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-skylib/releases/download/1.4.2/bazel-skylib-1.4.2.tar.gz",
        "https://github.com/bazelbuild/bazel-skylib/releases/download/1.4.2/bazel-skylib-1.4.2.tar.gz",
    ],
    sha256 = "66ffd9315665bfaafc96b52278f57c7e2dd09f5ede279ea6d39b2be471e7e3aa",
)

load("@bazel_skylib//:workspace.bzl", "bazel_skylib_workspace")

bazel_skylib_workspace()

http_archive(
    name = "com_google_protobuf",
    sha256 = "21fcb4b0df6a8e6279e5843af8c9f2245919cf0d3ec2021c76fccc4fc4bf9aca",
    strip_prefix = "protobuf-4.23.3",
    urls = ["https://github.com/protocolbuffers/protobuf/archive/v4.23.3.tar.gz"],
)

load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

protobuf_deps()

http_archive(
    # go_googleapis is used instead of com_google_googleapis in order to override
    # the dependency on github.com/googleapis/googleapis defined by rules_go
    # that is named go_googleapis. googleapis already has all of the necessary
    # rules, so using the rules_go patched version isn't necessary and it lags in
    # freshness which would require dependency overrides anyways.
    name = "go_googleapis",
    # Use `master` because googleapis isn't semantically versioned and the protos
    # this repo cares about (the annotation definitions) do not have breaking
    # changes, so we can live on HEAD. Pinning to commit is cumbersome to maintain.
    strip_prefix = "googleapis-master",
    urls = ["https://github.com/googleapis/googleapis/archive/master.tar.gz"],
)

load("@go_googleapis//:repository_rules.bzl", "switched_rules_by_language")
switched_rules_by_language(name = "com_google_googleapis_imports", go = True, grpc = True)

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "6dc2da7ab4cf5d7bfc7c949776b1b7c733f05e56edc4bcd9022bb249d2e2a996",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.39.1/rules_go-v0.39.1.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.39.1/rules_go-v0.39.1.zip",
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "29d5dafc2a5582995488c6735115d1d366fcd6a0fc2e2a153f02988706349825",
    urls = [
        "https://storage.googleapis.com/bazel-mirror/github.com/bazelbuild/bazel-gazelle/releases/download/v0.31.0/bazel-gazelle-v0.31.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.31.0/bazel-gazelle-v0.31.0.tar.gz",
    ],
)

load("//:repositories.bzl", "com_googleapis_gapic_generator_go_repositories")

# gazelle:repository_macro repositories.bzl%com_googleapis_gapic_generator_go_repositories
com_googleapis_gapic_generator_go_repositories()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_register_toolchains(version = "1.18.6")

go_rules_dependencies()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()

