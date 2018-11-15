# Bazel Out Files

## 背景

- [Bazel build event prototol](https://docs.bazel.build/versions/master/build-event-protocol.html) bazel 构建事件生成

- [google protobuf](https://developers.google.com/protocol-buffers/)

## 作用

将bazel构建事件的生成文件通过proto解析，然后筛选出构建生成的本地文件路径，并添加计算出该文件的sha256sum。

## 主要目录

BuildEventStream 和 ProtoBuf 是由bazel官方提供的proto文件生成的pb.go文件用于解析。

## 流程

1. 项目.bazelrc中添加build build --build_event_binary_file=/path/to/file
2. 使用bazel 编译项目生成事件文件以及生成输出文件。
3. 运行该项目。go run main.go --file=/path/to/file
