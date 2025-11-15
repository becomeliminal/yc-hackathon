package(default_visibility = ["PUBLIC"])

subinclude("///proto//build_defs:proto")

protoc_binary(
    name = "protoc",
    version = "3.20.3",
    visibility = ["PUBLIC"],
)

subinclude("//infrastructure/build_defs:github_release_asset", "//infrastructure/build_defs:remote", "//infrastructure/build_defs:proto_extensions")

def google_api(
        name,
        repo="googleapis/googleapis",
        files=[],
        dir,
        deps=[],
        revision="master",
        grpc=False,
        grpc_gateway=False):
    """Defines a rule to fetch the googleapis dependencies.
    Args:
      name (str): Name of the rule.
      repo (str): Repository to fetch files from. Defaults to grpc-gateway's.
      dir (str): Subdirectory in the repo to fetch files from.
      deps (list): Any other dependencies.
      visibility (list): Visibility spec of the rule.
    """
    srcs = []
    sub_directory = "/".join(dir.split("/")[1:])

    # Package name will be folder name. So we don't want to remove it.
    sub_directory = "/" + sub_directory if sub_directory != name else ""
    for f in files:
        sanitized_filename = f.replace("/", "_")
        src = github_file(
            name = f"_{name}#{sanitized_filename}",
            out = f,
            file = join_path(dir, f),
            repo = repo,
            revision = revision,
        )
        new_src = genrule(
            name = f"_{name}#{sanitized_filename}_correct_packages",
            srcs = [src],
            out = join_path(dir, f),
            cmd = "cat $SRCS > $OUT",
            deps = ["//third_party/proto:_protoc#wkt"],
        )
        srcs += [new_src]
    proto_def = grpc_library if grpc else proto_library
    proto_def = grpc_gateway_library if grpc_gateway else proto_def

    filegroup(
        name = name,
        srcs = srcs,
        deps = deps,
    )

google_api(
    name = "api",
    dir = "google/api",
    files = [
        "launch_stage.proto",
    ],
)

google_api(
    name = "httpbody",
    dir = "google/api",
    files = ["httpbody.proto"],
)

google_api(
    name = "annotations",
    dir = "google/api",
    files = [
        "annotations.proto",
        "client.proto",
        "field_behavior.proto",
        "http.proto",
        "resource.proto",
    ],
    deps = [
        ":api",
    ],
)

google_api(
    name = "status",
    dir = "google/rpc",
    files = ["status.proto"],
)

google_api(
    name = "code",
    dir = "google/rpc",
    files = ["code.proto"],
)

google_api(
    name = "longrunning",
    dir = "google/longrunning",
    files = ["operations.proto"],
    grpc_gateway = True,
    deps = [
        ":annotations",
        ":api",
        ":status",
    ],
)

google_api(
    name = "visibility",
    dir = "google/api",
    files = ["visibility.proto"],
)
