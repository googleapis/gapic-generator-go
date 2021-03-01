workspace(name = "com_googleapis_gapic_generator_go")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "com_google_protobuf",
    sha256 = "b10bf4e2d1a7586f54e64a5d9e7837e5188fc75ae69e36f215eb01def4f9721b",
    strip_prefix = "protobuf-3.15.3",
    urls = ["https://github.com/protocolbuffers/protobuf/archive/v3.15.3.tar.gz"],
)

load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

protobuf_deps()

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "7904dbecbaffd068651916dce77ff3437679f9d20e1a7956bff43826e7645fcc",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.25.1/rules_go-v0.25.1.tar.gz",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.25.1/rules_go-v0.25.1.tar.gz",
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "b85f48fa105c4403326e9525ad2b2cc437babaa6e15a3fc0b1dbab0ab064bc7c",
    urls = [
        "https://storage.googleapis.com/bazel-mirror/github.com/bazelbuild/bazel-gazelle/releases/download/v0.22.2/bazel-gazelle-v0.22.2.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.22.2/bazel-gazelle-v0.22.2.tar.gz",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

# gazelle:repo bazel_gazelle
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

# TODO(noahdietz): remove with next rules_go release.
# https://github.com/googleapis/gapic-generator-go/issues/529
go_repository(
    name = "org_golang_google_genproto",
    build_file_proto_mode = "disable_global",
    importpath = "google.golang.org/genproto",
    sum = "h1:SHQi2osmct8Y+ngNVppVlyB/WdW+XA9gHs8wPEE2xFY=",
    version = "v0.0.0-20210222212404-3e1e516060db",
)

go_rules_dependencies()

go_register_toolchains(
    version = "1.15.8",
)

gazelle_dependencies()

load("//:repositories.bzl", "com_googleapis_gapic_generator_go_repositories")

# gazelle:repository_macro repositories.bzl%com_googleapis_gapic_generator_go_repositories
com_googleapis_gapic_generator_go_repositories()
