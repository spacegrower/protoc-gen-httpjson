protoc-gen-jsonhttp 是一个 protoc 代码生成插件  
负责生成 gRPC Json web 客户端，需要结合 [ts-proto](https://github.com/stephenh/ts-proto) 使用  
`ts-proto` 负责生成 ts 类型文件  
`protoc-gen-jsonhttp` 负责生成 web 客户端框架

## 示例

假设 web 项目工程目录如下

```shell
...
protobuf/
src/
index.html
package.json
...
```

生成后的 ts 文件放置 src 目录下的 protobuf 文件夹

在根目录下执行：  

```shell
protoc --plugin=./node_modules/.bin/protoc-gen-ts_proto --jsonhttp_out=:./src/protobuf/ --ts_proto_out=:./src/ --ts_proto_opt=snakeToCamel=false ./protobuf/*
```
