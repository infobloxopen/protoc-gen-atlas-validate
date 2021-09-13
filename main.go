package main

import (
	"fmt"
	"io/ioutil"
	"os"

	av_opts "github.com/infobloxopen/protoc-gen-atlas-validate/options"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

type validateBuilder struct {
	plugin   *protogen.Plugin
	methods  map[string][]*methodDescriptor
	genFiles map[string]*protogen.GeneratedFile
}

func main() {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	var request pluginpb.CodeGeneratorRequest
	err = proto.Unmarshal(input, &request)
	if err != nil {
		panic(err)
	}

	opts := protogen.Options{}

	plugin, err := opts.New(&request)
	if err != nil {
		panic(err)
	}

	builder := &validateBuilder{
		methods:  make(map[string][]*methodDescriptor),
		genFiles: make(map[string]*protogen.GeneratedFile),
		plugin:   plugin,
	}

	for _, protoFile := range plugin.Files {
		methods := builder.gatherMethods(protoFile)
		protoName := *protoFile.Proto.Name
		builder.methods[protoName] = methods
	}

	resp := builder.generate(plugin)
	out, err := proto.Marshal(resp)
	if err != nil {
		panic(err)
	}

	fmt.Fprint(os.Stdout, string(out))
}

func (b *validateBuilder) generate(plugin *protogen.Plugin) *pluginpb.CodeGeneratorResponse {
	fmt.Fprintf(os.Stderr, "running generate\n")

	for _, protoFile := range plugin.Files {
		b.renderValidatorMethods(protoFile)
	}

	return plugin.Response()
}

func (b *validateBuilder) generateFile(protoFile *protogen.File) *protogen.GeneratedFile {
	g := b.genFiles[*protoFile.Proto.Name]
	if g != nil {
		return g
	}

	// create package first and return generatedFile
	fileName := protoFile.GeneratedFilenamePrefix + ".pb.atlas.validate.go"
	g = b.plugin.NewGeneratedFile(fileName, ".")
	g.P("package ", protoFile.GoPackageName)
	b.genFiles[*protoFile.Proto.Name] = g

	return g
}

func (b *validateBuilder) gatherMethods(file *protogen.File) []*methodDescriptor {
	var methods []*methodDescriptor

	for _, service := range file.Services {
		for _, method := range service.Methods {
			for i, opt := range extractHTTPOpts(method) {
				methods = append(methods, &methodDescriptor{
					inputTypeMessage: method.Input,
					svc:              string(service.Desc.Name()),
					method:           string(method.Desc.Name()),
					idx:              i,
					httpBody:         opt.body,
					httpMethod:       opt.method,
					gwPattern:        fmt.Sprintf("%s_%s_%d", service.Desc.Name(), method.Desc.Name(), i),
					inputType:        string(method.Input.Desc.Name()),
					allowUnknown:     b.getAllowUnknown(file.Desc.Options(), service.Desc.Options(), method.Desc.Options()),
				})
			}
		}
	}

	return methods
}

type methodDescriptor struct {
	inputTypeMessage *protogen.Message
	svc              string
	method           string
	httpBody         string
	httpMethod       string
	gwPattern        string
	inputType        string
	idx              int
	allowUnknown     bool
}

type httpOpt struct {
	body   string
	method string
}

func getHttpMethod(r *annotations.HttpRule) string {
	switch r.GetPattern().(type) {
	case *annotations.HttpRule_Get:
		return "GET"
	case *annotations.HttpRule_Post:
		return "POST"
	case *annotations.HttpRule_Put:
		return "PUT"
	case *annotations.HttpRule_Delete:
		return "DELETE"
	case *annotations.HttpRule_Patch:
		return "PATCH"
	}

	return ""
}

func extractHTTPOpts(m *protogen.Method) []httpOpt {
	r := []httpOpt{}

	options := m.Desc.Options()
	if options == nil {
		return nil
	}

	if !proto.HasExtension(options, annotations.E_Http) {
		return nil
	}

	ext := proto.GetExtension(options, annotations.E_Http)
	if ext == nil {
		return nil
	}

	if httpRule, ok := ext.(*annotations.HttpRule); ok {
		r = append(r, httpOpt{
			body:   httpRule.Body,
			method: getHttpMethod(httpRule),
		})
		for _, b := range httpRule.GetAdditionalBindings() {
			r = append(r, httpOpt{
				body:   b.Body,
				method: getHttpMethod(b),
			})
		}
	} else {
		return nil
	}

	return r
}

// getAllowUnknown function picks up correct allowUnknown option from file/service/method
// hierarchy.
func (b *validateBuilder) getAllowUnknown(file proto.Message, svc proto.Message, method proto.Message) bool {
	var gavOpt *av_opts.AtlasValidateFileOption
	v := proto.GetExtension(file, av_opts.E_File)
	if v != nil {
		gavOpt = v.(*av_opts.AtlasValidateFileOption)
	}

	var savOpt *av_opts.AtlasValidateServiceOption
	v = proto.GetExtension(svc, av_opts.E_Service)
	if v != nil {
		savOpt = v.(*av_opts.AtlasValidateServiceOption)
	}

	var mavOpt *av_opts.AtlasValidateMethodOption
	v = proto.GetExtension(method, av_opts.E_Method)
	if v != nil {
		mavOpt = v.(*av_opts.AtlasValidateMethodOption)
	}

	if mavOpt != nil {
		return mavOpt.GetAllowUnknownFields()
	} else if savOpt != nil {
		return savOpt.GetAllowUnknownFields()
	}

	return gavOpt.GetAllowUnknownFields()
}

// renderValidatorMethods function generates entrypoints for validator one per each
// HTTP request (and HTTP request additional_bindings).
func (b *validateBuilder) renderValidatorMethods(protoFile *protogen.File) {
	var g *protogen.GeneratedFile
	for _, m := range b.methods[*protoFile.Proto.Name] {
		// create protoFile iff we need to generate something
		g = b.generateFile(protoFile)

		g.P(`// validate_`, m.gwPattern, ` is an entrypoint for validating "`, m.httpMethod, `" HTTP request `)
		g.P(`// that match *.pb.gw.go/pattern_`, m.gwPattern, `.`)
		g.P(`func validate_`, m.gwPattern, `(ctx `, generateImport("Context", "context", g), `, r `, generateImport("RawMessage", "encoding/json", g), `) (err error) {`)

		if m.httpBody == "" {
			g.P(`if len(r) != 0 {`)
			g.P(`return `, generateImport("Errorf", "fmt", g), `("body is not allowed")`)
			g.P(`}`)
			g.P(`return nil`)
		} else if b.isWKT(m.inputType) {
			g.P(`return nil`)
		} else {
			fmt.Fprintf(os.Stderr, "httpBody: %s\n", m.httpBody)
			fmt.Fprintf(os.Stderr, "inputType : %s\n", m.inputTypeMessage.Desc.Name())
			typeName := string(m.inputTypeMessage.Desc.Name())

			if m.httpBody != "*" {
				for _, field := range m.inputTypeMessage.Fields {
					if string(field.Desc.Name()) == m.httpBody {
						typeName = string(field.Message.Desc.Name())
					}
				}
			}

			fmt.Fprintf(os.Stderr, "typeName: %s\n", typeName)

			if b.isLocal(m.inputTypeMessage) {
				g.P(`return validate_Object_`, typeName, `(ctx, r, "")`)
			} else {
				// g.P(`if validator, ok := `, b.generateAtlasValidateJSONInterfaceSignature(t), `; ok {`)
				// g.P(`return validator.AtlasValidateJSON(ctx, r, "")`)
				// g.P(`}`)
				g.P(`return nil`)
			}
		}
		g.P(`}`)
		g.P()
	}
}

func generateImport(name string, importPath string, g *protogen.GeneratedFile) string {
	return g.QualifiedGoIdent(protogen.GoIdent{
		GoName:       name,
		GoImportPath: protogen.GoImportPath(importPath),
	})
}

func (b *validateBuilder) isLocal(message *protogen.Message) bool {
	return true
}

var wkt = map[string]bool{
	// ptypes
	".google.protobuf.Timestamp": true,
	".google.protobuf.Duration":  true,
	".google.protobuf.Empty":     true,
	".google.protobuf.Any":       true,
	".google.protobuf.Struct":    true,

	// nillable values
	".google.protobuf.StringValue": true,
	".google.protobuf.BytesValue":  true,
	".google.protobuf.Int32Value":  true,
	".google.protobuf.UInt32Value": true,
	".google.protobuf.Int64Value":  true,
	".google.protobuf.UInt64Value": true,
	".google.protobuf.FloatValue":  true,
	".google.protobuf.DoubleValue": true,
	".google.protobuf.BoolValue":   true,
}

func (b *validateBuilder) isWKT(t string) bool {
	return wkt[t]
}
