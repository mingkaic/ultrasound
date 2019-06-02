load("//third_party/repos:tenncor.bzl", "tenncor_repository")
load("//third_party/repos:google_api.bzl", "google_apis_repository")

def dependencies(excludes = []):
    ignores = native.existing_rules().keys() + excludes

    if "com_github_mingkaic_tenncor" not in ignores:
        tenncor_repository()

    if "com_github_googleapis" not in ignores:
        google_apis_repository()
