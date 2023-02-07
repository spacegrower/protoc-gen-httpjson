package main

import (
	"flag"

	"github.com/spacegrower/protoc-gen-httpjson/plugin"
	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	tsGen := plugin.NewJsonHttpGen("httpjson", plugin.TS{
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
	var flags flag.FlagSet
	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(tsGen.Generate)
}
