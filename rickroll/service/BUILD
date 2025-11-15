go_binary(
    name = "service",
    srcs = [
        "main.go",
        "handler.go",
    ],
    visibility = ["PUBLIC"],
    deps = [
        "//rickroll/proto:api",
        "//third_party/go:grpc",
        "//third_party/go:grpc_gateway",
        "//third_party/go:grpc_gateway_x402",
    ],
)
