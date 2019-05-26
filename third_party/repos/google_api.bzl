load("@bazel_tools//tools/build_defs/repo:git.bzl", "new_git_repository")

GOOGLEAPIS_BUILD = """
load("@build_stack_rules_proto//cpp:rules.bzl", "cc_proto_library")
load("@build_stack_rules_proto//python:rules.bzl", "py_proto_library")

package(
    default_visibility = ["//visibility:public"],
)

filegroup(
    name = "annotations_proto",
    srcs = [
        "google/api/annotations.proto",
        "google/api/http.proto",
    ],
)

cc_proto_library(
    name = "annotations_cc_proto",
    protos = [":annotations_proto"],
    imports = ["external/com_google_protobuf/src"],
    inputs = ["@com_google_protobuf//:well_known_protos"],
)

py_proto_library(
    name = "annotations_py_proto",
    protos = [":annotations_proto"],
    imports = ["external/com_google_protobuf/src"],
    inputs = ["@com_google_protobuf//:well_known_protos"],
)
"""

def google_apis_repository(name):
    new_git_repository(
        name = name,
        remote = "https://github.com/googleapis/googleapis",
        commit = "8f1de3d40e2835d30f4c0bc861b4e8e8ec551138",
        build_file_content = GOOGLEAPIS_BUILD,
    )
