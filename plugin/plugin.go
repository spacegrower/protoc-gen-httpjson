package plugin

import (
	"fmt"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

type HttpJsonGen struct {
	suffix string
	ts     TS
}

func NewHttpJsonGen(suffix string, tscfg TS) *HttpJsonGen {
	return &HttpJsonGen{
		suffix: suffix,
		ts:     tscfg,
	}
}

type TS struct {
	ResponseTypeName         string
	ResponseTypeStruct       string
	ImportTsProtoPackageName string
}

func (j *HttpJsonGen) Generate(p *protogen.Plugin) error {
	if j.ts.ImportTsProtoPackageName == "" {
		j.ts.ImportTsProtoPackageName = "pb"
	}
	if j.ts.ResponseTypeName == "" {
		j.ts.ResponseTypeName = "GrpcGatewayResponse"
	}
	if j.ts.ResponseTypeStruct == "" {
		j.ts.ResponseTypeStruct = "{data: any}"
	}
	if j.suffix == "" {
		j.suffix = "httpjson"
	}
	for _, f := range p.Files {
		if !f.Generate {
			continue
		}
		if err := j.gen(p, f); err != nil {
			return err
		}
	}

	return nil
}

func (j *HttpJsonGen) gen(p *protogen.Plugin, f *protogen.File) error {

	protoFileName := filepath.Base(f.GeneratedFilenamePrefix)
	genFileName := fmt.Sprintf("%s_%s.ts", protoFileName, j.suffix)
	g := p.NewGeneratedFile(genFileName, f.GoImportPath)
	g.P()
	g.P()

	g.P("// Code generated by protoc-gen-" + j.suffix + ". DO NOT EDIT.")
	g.P("// - protoc             ", protocVersion(p))
	if f.Proto.GetOptions().GetDeprecated() {
		g.P("// ", f.Desc.Path(), " is a deprecated file.")
	} else {
		g.P("// source: ", f.Desc.Path())
	}
	g.P()
	g.P()
	// import ts-proto generate file
	g.P(`import * as ` + j.ts.ImportTsProtoPackageName + ` from "./` + protoFileName + "\";")
	g.P(`import type * as ` + j.ts.ImportTsProtoPackageName + `type from "./` + protoFileName + "\";")
	g.P()
	g.P("export type " + j.ts.ResponseTypeName + " = {")
	g.P(strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(j.ts.ResponseTypeStruct, "{"), "}")))
	g.P("}")
	g.P()
	g.P(`
type CallHandler<T> = (
    path: string,
    body: any,
    cfg?: T,
) => Promise<` + j.ts.ResponseTypeName + `>;

export interface CallOptions<T> {
    handler?: CallHandler<T>;
    cfg?: T;
}

// for example of axios: 
//
// async function handler(
//     path: string,
//	   body: any,
//	   cfg?: AxiosRequestConfig,
// ): Promise<` + j.ts.ResponseTypeName + `> {
//	   const resp = await axios.post(path, body, cfg);
//	   return resp.data as ` + j.ts.ResponseTypeName + `;
// }
// 
// const client = new XxxClient(handler);
	`)

	for _, service := range f.Services {
		j.generateClass(g, service)
	}

	return nil
}

func (j *HttpJsonGen) generateClass(g *protogen.GeneratedFile, service *protogen.Service) {
	className := service.Desc.Name() + "Client"
	g.P()
	g.P(`
export class ` + className + `<T> {

    private _baseURL: string;
	private _handler: CallHandler<T>;

	constructor(handler: CallHandler<T>, baseURL?: string) {
		if (baseURL !== undefined && baseURL.substring(baseURL.length - 1, baseURL.length) === '/') {
			baseURL = baseURL.substring(0, baseURL.length - 1)
		}
		this._handler = handler;
		this._baseURL = baseURL || '';
	}
	`)
	g.P()

	for _, method := range service.Methods {
		if !method.Desc.IsStreamingServer() && !method.Desc.IsStreamingClient() {
			j.generateClassMethod(g, service, method)
		} else {
			g.P("	// " + method.Desc.Name() + " is not support")
		}
	}
	g.P("}")
	g.P()
}

func (j *HttpJsonGen) generateClassMethod(g *protogen.GeneratedFile, service *protogen.Service, method *protogen.Method) {
	gwi := func(i ...interface{}) {
		var a = []interface{}{"    "}
		a = append(a, i...)
		g.P(a...)
	}
	gwi("async "+string(method.Desc.Name())+"(req: ", j.ts.ImportTsProtoPackageName, "type."+string(method.Input.Desc.Name())+", callOptions?: CallOptions<T>): Promise<"+j.ts.ImportTsProtoPackageName+"."+string(method.Output.Desc.Name())+"> {")
	gwi("    const resp = await this._handler(this._baseURL + '" + formatFullMethodName(service, method) + "', req, callOptions?.cfg)")
	gwi("    if (!resp.meta || resp.meta.code === undefined || resp.meta.message === undefined) {")
	gwi("        throw new Error('unknown response type');")
	gwi("    }")
	gwi("    if (resp.meta.code !== 0) {")
	gwi("        throw new Error(resp.meta.message);")
	gwi("    }")
	gwi("    return " + j.ts.ImportTsProtoPackageName + "." + method.Output.GoIdent.GoName + ".create(resp.data as " + j.ts.ImportTsProtoPackageName + "." + string(method.Output.Desc.Name()) + ")")
	gwi("}")
	gwi()
}

func protocVersion(gen *protogen.Plugin) string {
	v := gen.Request.GetCompilerVersion()
	if v == nil {
		return "(unknown)"
	}
	var suffix string
	if s := v.GetSuffix(); s != "" {
		suffix = "-" + s
	}
	return fmt.Sprintf("v%d.%d.%d%s", v.GetMajor(), v.GetMinor(), v.GetPatch(), suffix)
}

func formatFullMethodName(service *protogen.Service, method *protogen.Method) string {
	return fmt.Sprintf("/%s/%s", service.Desc.FullName(), method.Desc.Name())
}
