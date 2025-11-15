package(default_visibility = ["PUBLIC"])

# Go toolchain
go_toolchain(
    name = "gotool",
    version = "1.23.4",
)

#### Core Go modules ####

go_mod_download(
    name = "google.golang.org.protobuf_dl",
    module = "google.golang.org/protobuf",
    version = "v1.34.1",
)

go_module(
    name = "google.golang.org.protobuf",
    download = ":google.golang.org.protobuf_dl",
    install = ["..."],
    module = "google.golang.org/protobuf",
    deps = [":cmp"],
)

go_module(
    name = "cmp",
    install = [
        "cmp",
        "cmp/...",
    ],
    module = "github.com/google/go-cmp",
    version = "v0.7.0",
)

go_module(
    name = "protobuf",
    install = ["proto"],
    module = "github.com/golang/protobuf",
    version = "v1.5.4",
    deps = [":google.golang.org.protobuf"],
)

go_module(
    name = "google-uuid",
    module = "github.com/google/uuid",
    version = "v1.6.0",
)

go_module(
    name = "grpc",
    install = [
        ".",
        "attributes",
        "backoff",
        "balancer",
        "balancer/base",
        "balancer/grpclb/...",
        "balancer/grpclb/state",
        "balancer/pickfirst",
        "balancer/rls/...",
        "balancer/roundrobin",
        "binarylog/grpc_binarylog_v1",
        "channelz",
        "codes",
        "connectivity",
        "credentials",
        "credentials/alts",
        "credentials/alts/internal",
        "credentials/alts/internal/authinfo",
        "credentials/alts/internal/conn",
        "credentials/alts/internal/handshaker",
        "credentials/alts/internal/handshaker/service",
        "credentials/alts/internal/proto/grpc_gcp",
        "credentials/google",
        "credentials/insecure",
        "credentials/oauth",
        "encoding",
        "encoding/gzip",
        "encoding/proto",
        "experimental/stats",
        "keepalive",
        "grpclog",
        "grpclog/internal",
        "health",
        "mem",
        "health/grpc_health_v1",
        "internal",
        "internal/backoff",
        "internal/balancer/...",
        "internal/balancerload",
        "internal/binarylog",
        "internal/buffer",
        "internal/channelz",
        "internal/credentials",
        "internal/envconfig",
        "internal/googlecloud",
        "internal/grpclog",
        "internal/grpcsync",
        "internal/grpcutil",
        "internal/idle",
        "internal/metadata",
        "internal/pretty",
        "internal/resolver",
        "internal/resolver/dns",
        "internal/resolver/dns/internal",
        "internal/resolver/passthrough",
        "internal/resolver/unix",
        "internal/serviceconfig",
        "internal/stats",
        "internal/status",
        "internal/syscall",
        "internal/transport",
        "internal/transport/networktype",
        "internal/xds",
        "interop/grpc_testing",
        "keepalive",
        "metadata",
        "peer",
        "resolver",
        "resolver/dns",
        "resolver/passthrough",
        "resolver/manual",
        "reflection",
        "serviceconfig",
        "stats",
        "status",
        "tap",
    ],
    module = "google.golang.org/grpc",
    version = "v1.67.1",
    visibility = ["PUBLIC"],
    deps = [
        ":genproto_rpc",
        ":google.golang.org.protobuf",
        ":google-uuid",
        ":protobuf",
        ":x_net",
        ":x_oauth2",
        ":x_sync",
        ":x_sys",
    ],
)

go_mod_download(
    name = "grpc-gateway-download",
    module = "github.com/grpc-ecosystem/grpc-gateway/v2",
    strip = ["examples"],
    version = "v2.20.0",
)

go_module(
    name = "yaml.v3",
    module = "gopkg.in/yaml.v3",
    version = "v3.0.1",
    visibility = ["PUBLIC"],
)

go_module(
    name = "grpc_gateway",
    download = ":grpc-gateway-download",
    install = ["..."],
    module = "github.com/grpc-ecosystem/grpc-gateway/v2",
    deps = [
        ":genproto_api",
        ":google.golang.org.protobuf",
        ":grpc",
        ":yaml.v3",
    ],
)

go_module(
    name = "protoc-gen-go",
    binary = True,
    download = ":google.golang.org.protobuf_dl",
    install = ["cmd/protoc-gen-go"],
    module = "google.golang.org/protobuf",
    visibility = ["PUBLIC"],
    deps = [":google.golang.org.protobuf"],
)

go_module(
    name = "protoc-gen-go-grpc",
    binary = True,
    module = "google.golang.org/grpc/cmd/protoc-gen-go-grpc",
    version = "v1.1.0",
    visibility = ["PUBLIC"],
    deps = [":google.golang.org.protobuf"],
)

go_module(
    name = "protoc-gen-grpc-gateway",
    binary = True,
    download = ":grpc-gateway-download",
    install = ["protoc-gen-grpc-gateway"],
    module = "github.com/grpc-ecosystem/grpc-gateway/v2",
    deps = [
        ":genproto_api",
        ":google.golang.org.protobuf",
        ":grpc",
        ":yaml.v3",
    ],
)

go_module(
    name = "grpc_gateway_x402",
    install = ["..."],
    module = "github.com/becomeliminal/grpc-gateway-x402",
    version = "v0.1.0",
    deps = [
        ":grpc",
        ":grpc_gateway",
    ],
)

#### Google genproto ####

go_module(
    name = "genproto_api",
    install = [
        ".",
        "visibility",
        "annotations",
        "httpbody",
        "expr/...",
        "distribution",
        "label",
        "metric",
        "monitoredres",
    ],
    module = "google.golang.org/genproto/googleapis/api",
    version = "v0.0.0-20240610135401-a8a62080eff3",
    visibility = ["PUBLIC"],
    deps = [
        ":genproto_rpc",
        ":google.golang.org.protobuf",
        ":grpc",
        ":protobuf",
    ],
)

go_module(
    name = "genproto_rpc",
    install = [
        "status",
        "code",
        "errdetails",
    ],
    module = "google.golang.org/genproto/googleapis/rpc",
    version = "v0.0.0-20240610135401-a8a62080eff3",
    deps = [
        ":google.golang.org.protobuf",
        ":protobuf",
    ],
)

#### golang.org/x dependencies ####

go_module(
    name = "x_net",
    install = ["..."],
    module = "golang.org/x/net",
    strip = ["http2/h2demo"],
    version = "v0.0.0-20210614182718-04defd469f4e",
    deps = [
        ":x_crypto",
        ":x_text",
    ],
)

go_module(
    name = "x_sys",
    install = ["..."],
    module = "golang.org/x/sys",
    version = "v0.28.0",
)

go_module(
    name = "x_text",
    install = ["..."],
    module = "golang.org/x/text",
    version = "v0.21.0",
    strip = [
        "cmd",
        "message/pipeline",
    ],
)

go_module(
    name = "x_crypto",
    install = ["..."],
    module = "golang.org/x/crypto",
    strip = [
        "acme",
        "acme/autocert",
    ],
    version = "v0.32.0",
    deps = [
        ":x_sys",
        ":x_term",
    ],
)

go_module(
    name = "x_term",
    module = "golang.org/x/term",
    version = "v0.27.0",
    deps = [":x_sys"],
)

go_module(
    name = "x_sync",
    install = ["..."],
    module = "golang.org/x/sync",
    version = "v0.6.0",
)

go_module(
    name = "gcp_metadata",
    install = ["..."],
    module = "cloud.google.com/go/compute/metadata",
    version = "v0.6.0",
    deps = [
        ":x_net",
    ],
)

go_module(
    name = "x_oauth2",
    install = ["..."],
    module = "golang.org/x/oauth2",
    version = "v0.27.0",
    deps = [
        ":gcp_metadata",
        ":x_net",
    ],
)
