load("//third_party/repos:tenncor.bzl", "tenncor_repository")
load("//third_party/repos:protobuf.bzl", "protobuf_rules_repository")
load("//third_party/repos:grpc.bzl", "grpc_rules_repository")
load("//third_party/repos:google_api.bzl", "google_apis_repository")

def dependencies(excludes = []):
    ignores = native.existing_rules().keys() + excludes

    if "com_github_mingkaic_tenncor" not in ignores:
        tenncor_repository(name = "com_github_mingkaic_tenncor")

    if "protobuf_rules" not in ignores:
        protobuf_rules_repository(name = "protobuf_rules")

    if "com_github_grpc_grpc" not in ignores:
        grpc_rules_repository(name = "com_github_grpc_grpc")

    if "com_github_googleapis" not in ignores:
        google_apis_repository(name = "com_github_googleapis")
