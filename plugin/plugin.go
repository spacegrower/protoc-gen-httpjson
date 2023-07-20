package plugin

import (
	"os"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	importPkg        = make(map[string]PkgInfo)
	importPkgFaileds = make(map[string][]protoreflect.MessageDescriptor)
	filePkg          = ""
)

type PkgInfo struct {
	Prefix  string
	PkgName string
	PkgPath string
}

func (t *TS) getInputName(desc protoreflect.MessageDescriptor) string {
	pkgName := strings.TrimSuffix(string(desc.FullName()), "."+string(desc.Name()))
	if pkgName == filePkg {
		return t.ImportTsProtoPackageName + "type." + string(desc.Name())
	}

	if pkgInfo, exist := importPkg[pkgName]; exist {
		return pkgInfo.Prefix + string(desc.Name())
	}
	return string(desc.Name())
}

func (t *TS) getOutputName(desc protoreflect.MessageDescriptor) string {
	pkgName := strings.TrimSuffix(string(desc.FullName()), "."+string(desc.Name()))
	if pkgName == filePkg {
		return t.ImportTsProtoPackageName + "." + string(desc.Name())
	}

	if pkgInfo, exist := importPkg[pkgName]; exist {
		return pkgInfo.Prefix + string(desc.Name())
	}
	return string(desc.Name())
}

func (s *HttpJsonGen) Run() error {
	for _, v := range s.req.ProtoFile {
		if err := Gen(s, v); err != nil {
			return nil
		}
	}

	resp := s.Response()
	out, err := proto.Marshal(resp)
	if err != nil {
		return err
	}
	if _, err := os.Stdout.Write(out); err != nil {
		return err
	}
	return nil
}

type TS struct {
	ResponseTypeName         string
	ResponseTypeStruct       string
	ImportTsProtoPackageName string
}
