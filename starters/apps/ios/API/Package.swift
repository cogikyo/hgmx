// swift-tools-version: 6.0
import PackageDescription

let package = Package(
    name: "API",
    platforms: [.iOS(.v17)],
    products: [.library(name: "API", targets: ["API"])],
    dependencies: [
        .package(url: "https://github.com/connectrpc/connect-swift.git", exact: "1.2.3"),
        .package(url: "https://github.com/apple/swift-protobuf.git", exact: "1.39.0"),
    ],
    targets: [
        .target(
            name: "API",
            dependencies: [
                .product(name: "Connect", package: "connect-swift"),
                .product(name: "SwiftProtobuf", package: "swift-protobuf"),
            ]
        ),
    ]
)
