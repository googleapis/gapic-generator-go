workspace(name = "com_googleapis_gapic_go")

load("//:repositories.bzl", "com_googleapis_gapic_go_repositories")

com_googleapis_gapic_go_repositories()

load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

protobuf_deps()

load("@io_bazel_rules_go//go:deps.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains()

# gazelle:repo bazel_gazelle
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()

load("//:repositories_go.bzl", "com_googleapis_gapic_go_mod")

com_googleapis_gapic_go_mod()
