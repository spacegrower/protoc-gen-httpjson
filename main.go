package main

import (
	"github.com/spacegrower/protoc-gen-httpjson/plugin"
)

func main() {
	tsGen, err := plugin.New("httpjson", plugin.TS{
		ResponseTypeName: "GrpcGatewayResponse",
		ResponseTypeStruct: `{
			meta: {
				code: number;
				message: string;
				rid: string;
			};
			data: any;
		}`,
		ImportTsProtoPackageName: "pb",
	})

	if err != nil {
		panic(err)
	}
	tsGen.Run()
}
