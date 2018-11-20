# Use Execution Log Parser

[参考](https://github.com/bazelbuild/bazel/blob/master/src/tools/execlog/README.md)

## 构建

1. go get github.com/bazelbuild/bazel
2. modify .bazelrc
    add --experimental_execution_log_file=/tmp/exec.log
3. cd $GOHOME/src/bazel
4. bazel build src/tools/execlog:all
5. bazel-bin/src/tools/execlog/parser --log_path=/tmp/exec.log
6. python3 main.py


```
OUT: path: "bazel-out/darwin-fastbuild/bin/src/golang/darwin_amd64_stripped/helloworld"
     hash: "82518b5d209c461f5fc519322cbdc6979de296ddc1e69ec9f702855d2e29751b"
     remote_cache_hit: true
```