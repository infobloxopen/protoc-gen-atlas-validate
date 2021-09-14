package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	av_opts "github.com/infobloxopen/protoc-gen-atlas-validate/options"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

const (
	metadataPkgPath  = "google.golang.org/grpc/metadata"
	gwruntimePkgPath = "github.com/grpc-ecosystem/grpc-gateway/runtime"

	runtimePkgPath = "github.com/infobloxopen/protoc-gen-atlas-validate/runtime"
)

type validateBuilder struct {
	plugin      *protogen.Plugin
	methods     map[string][]*methodDescriptor
	genFiles    map[string]*protogen.GeneratedFile
	packageName string
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
		b.packageName = string(protoFile.GoPackageName)
		b.renderValidatorMethods(protoFile)
		b.renderValidatorObjectMethods(protoFile)
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
					inputType:        string(method.Input.Desc.FullName()),
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
			typeName := string(m.inputTypeMessage.Desc.Name())
			fullTypeName := string(m.inputTypeMessage.Desc.FullName())

			if m.httpBody != "*" {
				for _, field := range m.inputTypeMessage.Fields {
					if string(field.Desc.Name()) == m.httpBody {
						typeName = string(field.Message.Desc.Name())
						fullTypeName = string(field.Message.Desc.FullName())
						break
					}
				}
			}

			fmt.Fprintf(os.Stderr, "typeName: %s, fullName: %s\n", typeName, fullTypeName)

			if b.isLocal(fullTypeName) {
				g.P(`return validate_Object_`, typeName, `(ctx, r, "")`)
			} else {
				g.P(`if validator, ok := `, b.generateAtlasValidateJSONInterfaceSignature(fullTypeName, g), `; ok {`)
				g.P(`return validator.AtlasValidateJSON(ctx, r, "")`)
				g.P(`}`)
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

func (b *validateBuilder) isLocal(fullTypeName string) bool {
	sp := strings.Split(fullTypeName, ".")
	return sp[0] == b.packageName
}

var wkt = map[string]bool{
	// ptypes
	"google.protobuf.Timestamp": true,
	"google.protobuf.Duration":  true,
	"google.protobuf.Empty":     true,
	"google.protobuf.Any":       true,
	"google.protobuf.Struct":    true,

	// nillable values
	"google.protobuf.StringValue": true,
	"google.protobuf.BytesValue":  true,
	"google.protobuf.Int32Value":  true,
	"google.protobuf.UInt32Value": true,
	"google.protobuf.Int64Value":  true,
	"google.protobuf.UInt64Value": true,
	"google.protobuf.FloatValue":  true,
	"google.protobuf.DoubleValue": true,
	"google.protobuf.BoolValue":   true,
}

func (b *validateBuilder) isWKT(t string) bool {
	return wkt[t]
}

func (b *validateBuilder) generateAtlasValidateJSONInterfaceSignature(t string, g *protogen.GeneratedFile) string {
	return fmt.Sprintf(`interface{}(&%s{}).(interface{ AtlasValidateJSON(%s, %s, string) error })`,
		t, generateImport("Context", "context", g), generateImport("RawMessage", "encoding/json", g))
}

func (b *validateBuilder) generateAtlasJSONValidateInterfaceSignature(t string, g *protogen.GeneratedFile) string {
	rawMessage := generateImport("RawMessage", "encoding/json", g)
	return fmt.Sprintf(`interface{}(&%s{}).(interface { AtlasJSONValidate(%s, %s, string) (%s, error) })`,
		t, generateImport("Context", "context", g), rawMessage, rawMessage)
}

func (b *validateBuilder) renderValidatorObjectMethods(protoFile *protogen.File) {
	g := b.generateFile(protoFile)
	for _, message := range protoFile.Messages {
		// ptype := "." + p.file.GetPackage() + "." + o.GetName()
		// otype := p.TypeName(p.objectNamed(ptype))

		b.renderValidatorObjectMethod(message, g)
		// b.generateValidateRequired(o, otype)

		for _, innerMessage := range message.Messages {
			if innerMessage.Desc.IsMapEntry() {
				continue
			}

			// notype := p.TypeName(p.objectNamed(ptype + "." + no.GetName()))

			// b.renderValidatorObjectMethod(message, g)
			// b.generateValidateRequired(no, notype)
		}
	}
}

func (b *validateBuilder) renderValidatorObjectMethod(message *protogen.Message, g *protogen.GeneratedFile) {
	// t should be type
	t := string(message.Desc.Name())

	g.P(`// validate_Object_`, t, ` function validates a JSON for a given object.`)
	g.P(`func validate_Object_`, t, `(ctx `, generateImport("Context", "context", g), `, r `, generateImport("RawMessage", "encoding/json", g), `, path string) (err error) {`)
	g.P(`if hook, ok := `, b.generateAtlasJSONValidateInterfaceSignature(t, g), `; ok {`)
	g.P(`if r, err = hook.AtlasJSONValidate(ctx, r, path); err != nil {`)
	g.P(`return err`)
	g.P(`}`)
	g.P(`}`)
	g.P()
	g.P(`var v map[string]`, generateImport("RawMessage", "encoding/json", g))
	g.P(`if err = `, generateImport("Unmarshal", "encoding/json", g), `(r, &v); err != nil {`)
	g.P(`return `, generateImport("Errorf", "fmt", g), `("invalid value for %q: expected object.", path)`)
	g.P(`}`)
	g.P()
	g.P(`if err = validate_required_Object_`, t, `(ctx, v, path); err != nil {`)
	g.P(`return err`)
	g.P(`}`)
	g.P()
	g.P(`allowUnknown := `, generateImport("AllowUnknownFromContext", runtimePkgPath, g), `(ctx)`)
	g.P()
	g.P(`for k, _ := range v {`)

	g.P(`switch k {`)
	for _, f := range message.Fields {
		g.P(`case "`, f.Desc.Name(), `":`)

		if f.Desc.IsMap() {
			continue
		}

		fExt := proto.GetExtension(f.Desc.Options(), av_opts.E_Field)
		if fExt != nil {
			favOpt := fExt.(*av_opts.AtlasValidateFieldOption)
			methods := b.GetDeniedMethods(favOpt.GetDeny())
			if len(methods) != 0 {
				cond := strings.Join(methods, `" || method == "`)
				g.P(`method := `, generateImport("HTTPMethodFromContext", runtimePkgPath, g), `(ctx)`)
				g.P(`if method == "`, cond, `" {`)
				g.P(`return `, generateImport("Errorf", "fmt", g), `("field %q is unsupported for %q operation.", k, method)`)
				g.P("}")
			}
		}

		if f.Message != nil && f.Desc.IsList() {
			g.P(`if v[k] == nil {`)
			g.P(`continue`)
			g.P(`}`)
			g.P(`var vArr []`, generateImport("RawMessage", "encoding/json", g))
			g.P(`vArrPath := `, generateImport("JoinPath", runtimePkgPath, g), `(path, k)`)
			g.P(`if err = `, generateImport("Unmarshal", "encoding/json", g), `(v[k], &vArr); err != nil {`)
			g.P(`return `, generateImport("Errorf", "fmt", g), `("invalid value for %q: expected array.", vArrPath)`)
			g.P(`}`)

			ft := strings.Title(string(f.Desc.Message().Name()))
			fft := string(f.Desc.Message().FullName())

			if b.isWKT(fft) {
				continue
			}

			if !b.isLocal(fft) {
				g.P(`validator, ok := `, b.generateAtlasValidateJSONInterfaceSignature(ft, g))
				g.P(`if !ok {`)
				g.P(`continue`)
				g.P(`}`)
			}
			g.P(`for i, vv := range vArr {`)
			g.P(`vvPath := `, generateImport("Sprintf", "fmt", g), `("%s.[%d]", vArrPath, i)`)
			if b.isLocal(fft) {
				g.P(`if err = validate_Object_`, objectName(fft), `(ctx, vv, vvPath); err != nil {`)
				g.P(`return err`)
				g.P(`}`)
			} else {
				g.P(`if err = validator.AtlasValidateJSON(ctx, vv, vvPath); err != nil {`)
				g.P(`return err`)
				g.P(`}`)
			}
			g.P(`}`)

		} else if f.Message != nil {
			// ft := strings.Title(string(f.Desc.Message().Name()))
			fft := string(f.Desc.Message().FullName())

			if b.isWKT(fft) {
				continue
			}

			g.P(`if v[k] == nil {`)
			g.P(`continue`)
			g.P(`}`)
			g.P(`vv := v[k]`)
			g.P(`vvPath := `, generateImport("JoinPath", runtimePkgPath, g), `(path, k)`)
			if b.isLocal(fft) {
				g.P(`if err = validate_Object_`, objectName(fft), `(ctx, vv, vvPath); err != nil {`)
				g.P(`return err`)
				g.P(`}`)
			} else {
				g.P(`validator, ok := `, b.generateAtlasValidateJSONInterfaceSignature(fft, g))
				g.P(`if !ok {`)
				g.P(`continue`)
				g.P(`}`)
				g.P(`if err = validator.AtlasValidateJSON(ctx, vv, vvPath); err != nil {`)
				g.P(`return err`)
				g.P(`}`)
			}
		}
	}

	g.P(`default:`)
	g.P(`if !allowUnknown {`)
	g.P(`return `, generateImport("Errorf", "fmt", g), `("unknown field %q.", `, generateImport("JoinPath", runtimePkgPath, g), `(path, k))`)
	g.P(`}`)
	g.P(`}`)
	g.P(`}`)
	g.P(`return nil`)
	g.P(`}`)
	g.P()

	g.P(`// AtlasValidateJSON function validates a JSON for object `, t, `.`)
	g.P(`func (_ *`, t, `) AtlasValidateJSON(ctx `, generateImport("Context", "context", g), `, r `, generateImport("RawMessage", "encoding/json", g), `, path string) (err error) {`)
	g.P(`if hook, ok := `, b.generateAtlasJSONValidateInterfaceSignature(t, g), `; ok {`)
	g.P(`if r, err = hook.AtlasJSONValidate(ctx, r, path); err != nil {`)
	g.P(`return err`)
	g.P(`}`)
	g.P(`}`)
	g.P(`return validate_Object_`, t, `(ctx, r, path)`)
	g.P(`}`)
	g.P()
}

func objectName(fullName string) string {
	sp := strings.Split(fullName, ".")
	if len(sp) == 1 {
		return fullName
	}

	return strings.Join(sp[1:], "_")
}

//Return methods to which field marked as denied
func (b *validateBuilder) GetDeniedMethods(options []av_opts.AtlasValidateFieldOption_Operation) []string {
	httpMethods := make(map[string]struct{}, 0)
	for _, op := range options {
		switch op {
		case av_opts.AtlasValidateFieldOption_create:
			httpMethods["POST"] = struct{}{}
		case av_opts.AtlasValidateFieldOption_update:
			httpMethods["PATCH"] = struct{}{}
		case av_opts.AtlasValidateFieldOption_replace:
			httpMethods["PUT"] = struct{}{}
		}
	}

	uniqueMethods := make([]string, 0)
	for m := range httpMethods {
		uniqueMethods = append(uniqueMethods, m)
	}

	sort.StringSlice(uniqueMethods).Sort()
	return uniqueMethods
}
