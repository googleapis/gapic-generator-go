workspace(name = "com_googleapis_gapic_generator_go")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

# Workaround for https://github.com/bazelbuild/bazel-gazelle/issues/1285. Ideally,
# we can remove this if gazelle ships a fix since we didn't need it before.
http_archive(
    name = "bazel_skylib",
    sha256 = "37cdfbc6faefea94f7b37760a305c98c08981116c2bc9e821e3b423221fad8c8",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-skylib/releases/download/1.9.2/bazel-skylib-1.9.2.tar.gz",
        "https://github.com/bazelbuild/bazel-skylib/releases/download/1.9.2/bazel-skylib-1.9.2.tar.gz",
    ],
)

load("@bazel_skylib//:workspace.bzl", "bazel_skylib_workspace")

bazel_skylib_workspace()

http_archive(
    name = "rules_java",
    urls = [
        "https://github.com/bazelbuild/rules_java/releases/download/9.8.0/rules_java-9.8.0.tar.gz",
    ],
    sha256 = "88537048d7d07589ff4939c99889300c7063e43480ddd3cfa94d6c2cc899aa2e",
)

http_archive(
    name = "bazel_features",
    sha256 = "5450bfb2c8b4bc961c75368838f86156f563cc9adef1be7d504fc5619d54daab",
    strip_prefix = "bazel_features-1.51.0",
    url = "https://github.com/bazel-contrib/bazel_features/releases/download/v1.51.0/bazel_features-v1.51.0.tar.gz",
)

load("@bazel_features//:deps.bzl", "bazel_features_deps")
bazel_features_deps()

load("@rules_java//java:rules_java_deps.bzl", "rules_java_dependencies")
rules_java_dependencies()

http_archive(
    name = "com_google_protobuf",
    sha256 = "29372fc269325b491126d204931eda436bf9a593f46b3467fac42c40f9b32944",
    strip_prefix = "protobuf-31.0",
    urls = ["https://github.com/protocolbuffers/protobuf/releases/download/v31.0/protobuf-31.0.tar.gz"],
)

load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

protobuf_deps()

http_archive(
    name = "rules_proto",
    sha256 = "14a225870ab4e91869652cfd69ef2028277fc1dc4910d65d353b62d6e0ae21f4",
    strip_prefix = "rules_proto-7.1.0",
    url = "https://github.com/bazelbuild/rules_proto/releases/download/7.1.0/rules_proto-7.1.0.tar.gz",
)

load("@rules_proto//proto:repositories.bzl", "rules_proto_dependencies")

rules_proto_dependencies()

load("@rules_proto//proto:setup.bzl", "rules_proto_setup")

rules_proto_setup()

load("@rules_proto//proto:toolchains.bzl", "rules_proto_toolchains")

rules_proto_toolchains()

http_archive(
    name = "com_google_googleapis",
    # Use `master` because googleapis isn't semantically versioned and the protos
    # this repo cares about (the annotation definitions) do not have breaking
    # changes, so we can live on HEAD. Pinning to commit is cumbersome to maintain.
    strip_prefix = "googleapis-master",
    urls = ["https://github.com/googleapis/googleapis/archive/master.tar.gz"],
)

load("@com_google_googleapis//:repository_rules.bzl", "switched_rules_by_language")

switched_rules_by_language(
    name = "com_google_googleapis_imports",
    go = True,
    grpc = True,
)

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "86d3dc8f59d253524f933aaf2f3c05896cb0b605fc35b460c0b4b039996124c6",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.60.0/rules_go-v0.60.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.60.0/rules_go-v0.60.0.zip",
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "675114d8b433d0a9f54d81171833be96ebc4113115664b791e6f204d58e93446",
    urls = [
        "https://storage.googleapis.com/bazel-mirror/github.com/bazelbuild/bazel-gazelle/releases/download/v0.47.0/bazel-gazelle-v0.47.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.47.0/bazel-gazelle-v0.47.0.tar.gz",
    ],
)

load("//:repositories.bzl", "com_googleapis_gapic_generator_go_repositories")

# gazelle:repository_macro repositories.bzl%com_googleapis_gapic_generator_go_repositories
com_googleapis_gapic_generator_go_repositories()

http_archive(
    name = "rules_python",
    sha256 = "678358df0f4ae7f00a33ce33cced97c301f3c523bdd71cb2573d20a93efbf6c6",
    strip_prefix = "rules_python-2.3.1",
    url = "https://github.com/bazelbuild/rules_python/releases/download/2.3.1/rules_python-2.3.1.tar.gz",
)

load("@rules_python//python:repositories.bzl", "py_repositories")
py_repositories()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_register_toolchains(version = "1.25.8")

go_rules_dependencies()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()
