package plugin

import (
	"flag"

	"google.golang.org/protobuf/compiler/protogen"
)

func Run(fileSuffix string, ts TS) {

	// flags.StringVar(&g.Prefix, "prefix", "/", "API path prefix")
	j := &JsonHttpGen{
		suffix: fileSuffix,
		TS:     ts,
	}
	var flags flag.FlagSet
	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(j.Generate)
}
