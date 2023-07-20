package plugin

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

type genOptions = protogen.Options

func New(suffix string, tscfg TS) (*HttpJsonGen, error) {
	// if len(os.Args) > 1 {
	// 	return nil, fmt.Errorf("unknown argument %q (this program should be run by protoc, not directly)", os.Args[1])
	// }
	in, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return nil, err
	}
	req := &pluginpb.CodeGeneratorRequest{}
	if err := proto.Unmarshal(in, req); err != nil {
		return nil, err
	}

	return &HttpJsonGen{
		suffix: suffix,
		ts:     tscfg,

		req: req,
	}, nil
}

type HttpJsonGen struct {
	suffix string
	ts     TS

	req *pluginpb.CodeGeneratorRequest

	generatedFile []*File
}

type GeneratedFile struct {
	skip     bool
	fileName string
	filePath string

	buf              bytes.Buffer
	packageNames     map[string]string
	usedPackageNames map[string]bool
	manualImports    map[string]bool
	resultFileName   string
}

func (g *GeneratedFile) P(v ...interface{}) {
	for _, x := range v {
		switch x := x.(type) {
		default:
			fmt.Fprint(&g.buf, x)
		}
	}
	fmt.Fprintln(&g.buf)
}

var (
	fileReg = new(protoregistry.Files)
)

type File struct {
	Desc  protoreflect.FileDescriptor
	Proto *descriptorpb.FileDescriptorProto

	gen *GeneratedFile
}
