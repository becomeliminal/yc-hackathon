subinclude("//infrastructure/build_defs:proto_extensions")

grpc_gateway_library(
    name = "api",
    srcs = ["rickroll.proto"],
    visibility = ["PUBLIC"],
    deps = [
        "//third_party/proto:annotations",
        "//third_party/go:genproto_api",
    ],
)
