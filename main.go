package main

import "github.com/spacegrower/protoc-gen-httpjson/plugin"

func main() {
	plugin.Run("jsonhttp", plugin.TS{
		ResponseTypeName: "GrpcGatewayResponse",
		ResponseTypeStruct: `{
			meta: {
				code: number;
				message: string;
				request_id: string;
			};
			data: any;
		}`,
		ImportTsProtoPackageName: "pb",
	})
}
