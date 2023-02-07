protoc-gen-httpjson 是一个 protoc 代码生成插件  
负责生成 gRPC Json Web 客户端，需要结合 [ts-proto](https://github.com/stephenh/ts-proto) 使用  
`ts-proto` 负责生成 ts 类型文件  
`protoc-gen-httpjson` 负责生成 web 客户端框架，该客户端可以将 protobuf 定义的数据结构转为 json 格式的 HTTP Post 请求

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

### protobuf

```protobuf
syntax = "proto3";

package srv;

service Platform {
    rpc Login(LoginRequest) returns (LoginReply) {}
}

message LoginRequest {
    string account = 1;
    string password = 2;
}
message LoginReply {
    int64 id = 1;
    string account = 2;
    string nickname = 3;
    bool is_new = 4;
}
```

### typescript

```ts
export type GrpcGatewayResponse = {
  meta: {
    code: number;
    message: string;
    rid: string;
  };
  data: any;
};

type CallHandler<T> = (
  path: string,
  body: any,
  cfg: T | undefined
) => Promise<GrpcGatewayResponse>;

export interface CallOptions<T> {
  handler?: CallHandler<T>;
  cfg: T;
}

// for example of axios:
//
// async function handler(
//     path: string,
//	   body: any,
//	   cfg?: AxiosRequestConfig,
// ): Promise<GrpcGatewayResponse> {
//	   const resp = await axios.post(path, body, cfg);
//	   return resp.data as GrpcGatewayResponse;
// }
//
// const client = new XxxClient(handler);

export class PlatformClient<T> {
  private _baseURL: string;
  private _handler: CallHandler<T>;

  constructor(handler: CallHandler<T>, baseURL?: string) {
    if (
      baseURL !== undefined &&
      baseURL.substr(baseURL.length - 1, 1) === "/"
    ) {
      baseURL = baseURL.substring(baseURL.length - 1, 1);
    }
    this._handler = handler;
    this._baseURL = baseURL || "";
  }

  async Login(
    req: pbtype.LoginRequest,
    callOptions?: CallOptions<T>
  ): Promise<pb.LoginReply> {
    const resp = await this._handler(
      this._baseURL + "/srv.Platform/Login",
      req,
      callOptions?.cfg
    );
    if (
      !resp.meta ||
      resp.meta.code === undefined ||
      resp.meta.message === undefined
    ) {
      throw new Error("unknown response type");
    }
    if (resp.meta.code !== 0) {
      throw new Error(resp.meta.message);
    }
    return pb.LoginReply.create(resp.data as pb.LoginReply);
  }
}
```
