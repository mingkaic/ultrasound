workspace(name = "com_github_mingkaic_ultrasound")

# gazelle:repo bazel_gazelle

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

http_archive(
    name = "io_bazel_rules_go",
    urls = ["https://github.com/bazelbuild/rules_go/releases/download/0.18.5/rules_go-0.18.5.tar.gz"],
    sha256 = "a82a352bffae6bee4e95f68a8d80a70e87f42c4741e6a448bec11998fcc82329",
)

load("@io_bazel_rules_go//go:deps.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains()

http_archive(
    name = "bazel_gazelle",
    urls = ["https://github.com/bazelbuild/bazel-gazelle/releases/download/0.17.0/bazel-gazelle-0.17.0.tar.gz"],
    sha256 = "3c681998538231a2d24d0c07ed5a7658cb72bfb5fd4bf9911157c0e9ac6a2687",
)

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

gazelle_dependencies()

http_archive(
    name = "io_bazel_rules_docker",
    sha256 = "29d109605e0d6f9c892584f07275b8c9260803bf0c6fcb7de2623b2bedc910bd",
    strip_prefix = "rules_docker-0.5.1",
    urls = ["https://github.com/bazelbuild/rules_docker/archive/v0.5.1.tar.gz"],
)

load(
    "@io_bazel_rules_docker//go:image.bzl",
    _go_image_repos = "repositories",
)

_go_image_repos()

# === external dependencies ===
load("//:third_party/all.bzl", "dependencies")

dependencies()

load("@com_github_mingkaic_tenncor//:third_party/all.bzl", _tenncor_repos = "dependencies")

_tenncor_repos()

load("@com_github_stackb_rules_proto//go:deps.bzl", "go_grpc_library")

go_grpc_library()

load("@com_github_stackb_rules_proto//cpp:deps.bzl", "cpp_proto_library", "cpp_grpc_library")

cpp_proto_library()

cpp_grpc_library()

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()

# === golang external dependencies ===

go_repository(
    name = "com_github_sirupsen_logrus",
    commit = "2a22dbedbad1fd454910cd1f44f210ef90c28464",
    importpath = "github.com/sirupsen/logrus",
)

go_repository(
    name = "com_github_lib_pq",
    commit = "4ded0e9383f75c197b3a2aaa6d590ac52df6fd79",
    importpath = "github.com/lib/pq",
)

go_repository(
    name = "com_github_grpc_ecosystem_grpc_gateway",
    commit = "e6f18d33a7b3bfa5b94f3d5fb513252184ce2d90",
    importpath = "github.com/grpc-ecosystem/grpc-gateway",
)

go_repository(
    name = "com_github_ghodss_yaml",
    commit = "0ca9ea5df5451ffdf184b4428c902747c2c11cd7",
    importpath = "github.com/ghodss/yaml",
)

go_repository(
    name = "in_gopkg_yaml_v2",
    commit = "eb3733d160e74a9c7e442f435eb3bea458e1d19f",
    importpath = "gopkg.in/yaml.v2",
)

http_archive(
    name = "com_github_bazelbuild_buildtools",
    strip_prefix = "buildtools-bf564b4925ab5876a3f64d8b90fab7f769013d42",
    url = "https://github.com/bazelbuild/buildtools/archive/bf564b4925ab5876a3f64d8b90fab7f769013d42.zip",
)

load("@com_github_bazelbuild_buildtools//buildifier:deps.bzl", "buildifier_dependencies")

buildifier_dependencies()
