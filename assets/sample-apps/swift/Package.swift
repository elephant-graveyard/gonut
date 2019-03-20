// swift-tools-version:4.0
// The swift-tools-version declares the minimum version of Swift required to build this package.

import PackageDescription

let package = Package(
    name: "johnny-5",
    dependencies: [
        .package(url: "https://github.com/IBM-Swift/Kitura.git", from: "2.6.2")
    ],
    targets: [
        .target(
            name: "App",
            dependencies: ["Kitura"]),
    ]
)
