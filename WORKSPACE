workspace(name = "com_googleapis_gapic_generator_go")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "com_google_protobuf",
    strip_prefix = "protobuf-3.12.0-rc1",
    urls = ["https://github.com/protocolbuffers/protobuf/archive/v3.12.0-rc1.tar.gz"],
    sha256 = "f9c6433581be5209085d7a347fbaf16477f423969922bdc2f31faa66e7e3c6d6",
)

load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

# TODO(ndietz) replace this with the next version (v0.22.5?) of rules_go
# that contains github.com/golang/protobuf v1.4.1 and google.golang.org/protobuf
# v1.22.0
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
        "https://storage.googleapis.com/bazel-mirror/github.com/bazelbuild/bazel-gazelle/releases/download/v0.20.0/bazel-gazelle-v0.20.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.20.0/bazel-gazelle-v0.20.0.tar.gz",
    ],
    sha256 = "d8c45ee70ec39a57e7a05e5027c32b1576cc7f16d9dd37135b0eddde45cf1b10",
)

# gazelle:repo bazel_gazelle
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

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

protobuf_deps()

go_rules_dependencies()

go_register_toolchains()

gazelle_dependencies()

load("//:repositories.bzl", "com_googleapis_gapic_generator_go_repositories")

com_googleapis_gapic_generator_go_repositories()
