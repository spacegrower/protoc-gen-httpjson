package main

import (
	"github.com/spacegrower/protoc-gen-httpjson/plugin"
)

func main() {
	// tsGen := plugin.NewHttpJsonGen("httpjson", plugin.TS{
	// 	ResponseTypeName: "GrpcGatewayResponse",
	// 	ResponseTypeStruct: `{
	// 		meta: {
	// 			code: number;
	// 			message: string;
	// 			rid: string;
	// 		};
	// 		data: any;
	// 	}`,
	// 	ImportTsProtoPackageName: "pb",
	// })
	// var flags flag.FlagSet
	// protogen.Options{
	// 	ParamFunc: flags.Set,
	// 	ImportRewriteFunc: func(gip protogen.GoImportPath) protogen.GoImportPath {
	// 		return "filtered"
	// 	},
	// }.Run(tsGen.Generate)

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
	tsGen.Run(plugin.Gen)
}
