load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

# gazelle:proto disable_global
go_library(
    name = "protoc-gen-go_gapic_lib",
    srcs = ["main.go"],
    importpath = "github.com/googleapis/gapic-generator-go/cmd/protoc-gen-go_gapic",
    visibility = ["//visibility:private"],
    deps = [
        "//internal/gengapic",
        "@com_github_golang_protobuf//protoc-gen-go/plugin",
        "@org_golang_google_protobuf//proto",
    ],
)

go_binary(
    name = "protoc-gen-go_gapic",
    embed = [":protoc-gen-go_gapic_lib"],
    visibility = ["//visibility:public"],
)
