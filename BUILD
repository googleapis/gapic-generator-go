load("@bazel_gazelle//:def.bzl", "gazelle")

# gazelle:prefix github.com/googleapis/gapic-generator-go
gazelle(name = "gazelle")

# gazelle:resolve proto proto google/rpc/code.proto @googleapis//google/rpc:code_proto
# gazelle:resolve proto go google/rpc/code.proto  @org_golang_google_genproto_googleapis_rpc//code
# gazelle:resolve proto proto google/api/annotations.proto @googleapis//google/api:annotations_proto
# gazelle:resolve proto go google/api/annotations.proto  @org_golang_google_genproto//googleapis/api/annotations
# gazelle:resolve proto proto google/longrunning/operations.proto @googleapis//google/longrunning:operations_proto
# gazelle:resolve proto go google/longrunning/operations.proto @googleapis//google/longrunning:longrunning_go_proto
