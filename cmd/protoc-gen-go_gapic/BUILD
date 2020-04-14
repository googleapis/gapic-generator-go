load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

# gazelle:proto disable_global
go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importpath = "github.com/googleapis/gapic-generator-go/cmd/protoc-gen-go_gapic",
    visibility = ["//visibility:private"],
    deps = [
        "//internal/gengapic:go_default_library",
        "//internal/gensample:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
        "@io_bazel_rules_go//proto/wkt:compiler_plugin_go_proto",
    ],
)

go_binary(
    name = "protoc-gen-go_gapic",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)
