use grpc_build::build;

fn main() {
    build(
        "./src",       // protobuf files input dir
        "autogen/lang/rust", // output directory
        true,           // --build_server=true
        true,           // --build_client=true
        true,           // --force
    )
        .unwrap();
}