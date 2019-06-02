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

git_repository(
    name = "com_github_grpc_ecosystem_grpc_gateway",
    remote = "https://github.com/grpc-ecosystem/grpc-gateway",
    commit = "aeab1d96e0f1368d243e2e5f526aa29d495517bb",
)

load("@com_github_grpc_ecosystem_grpc_gateway//:repositories.bzl", "repositories")

repositories()

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
    name = "com_github_jinzhu_gorm",
    commit = "6ed508ec6a4ecb3531899a69cbc746ccf65a4166",
    importpath = "github.com/jinzhu/gorm",
)

go_repository(
    name = "com_github_jinzhu_inflection",
    commit = "04140366298a54a039076d798123ffa108fff46c",
    importpath = "github.com/jinzhu/inflection",
)

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
