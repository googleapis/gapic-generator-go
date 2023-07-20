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
    sha256 = "5f3cd52d6e5062071d99da57a96ea87e39bc020d6d25748001d919c474a4d8ed",
    strip_prefix = "protobuf-4.23.4",
    urls = ["https://github.com/protocolbuffers/protobuf/archive/v4.23.4.tar.gz"],
)

load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

protobuf_deps()

http_archive(
    name = "com_google_googleapis",
    strip_prefix = "googleapis-master",
    urls = ["https://github.com/googleapis/googleapis/archive/master.tar.gz"],
)

load("@com_google_googleapis//:repository_rules.bzl", "switched_rules_by_language")
switched_rules_by_language(name = "com_google_googleapis_imports", go = True, grpc = True)

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "278b7ff5a826f3dc10f04feaf0b70d48b68748ccd512d7f98bf442077f043fe3",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.41.0/rules_go-v0.41.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.41.0/rules_go-v0.41.0.zip",
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "29218f8e0cebe583643cbf93cae6f971be8a2484cdcfa1e45057658df8d54002",
    urls = [
        "https://storage.googleapis.com/bazel-mirror/github.com/bazelbuild/bazel-gazelle/releases/download/v0.32.0/bazel-gazelle-v0.32.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.32.0/bazel-gazelle-v0.32.0.tar.gz",
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

