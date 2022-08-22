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
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/pluginpb"
)

const (
	metadataPkgPath  = "google.golang.org/grpc/metadata"
	gwruntimePkgPath = "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	runtimePkgPath   = "github.com/infobloxopen/protoc-gen-atlas-validate/runtime"
)

type validateBuilder struct {
	plugin      *protogen.Plugin
	methods     map[string][]*methodDescriptor
	packageName string
	renderOnce  bool
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
		methods: make(map[string][]*methodDescriptor),
		plugin:  plugin,
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
	var lastG *protogen.GeneratedFile

	for _, protoFile := range plugin.Files {
		if !protoFile.Generate {
			continue
		}

		fileName := protoFile.GeneratedFilenamePrefix + ".pb.atlas.validate.go"
		g := b.plugin.NewGeneratedFile(fileName, ".")
		lastG = g
		g.P("package ", protoFile.GoPackageName)

		b.packageName = string(protoFile.Desc.Package())
		b.renderValidatorMethods(protoFile, g)
		b.renderValidatorObjectMethods(protoFile, g)

		if !b.renderOnce && sameName(b.packageName, *protoFile.Proto.Name) {
			b.renderOnce = true
			b.renderMethodDescriptors(g)
			b.renderAnnotator(g)
		}
	}

	if !b.renderOnce {
		b.renderMethodDescriptors(lastG)
		b.renderAnnotator(lastG)
	}

	return plugin.Response()
}

func sameName(packageName string, fileName string) bool {
	fs := strings.Split(fileName, "/")
	return strings.HasPrefix(fs[len(fs)-1], packageName+".")
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
func (b *validateBuilder) renderValidatorMethods(protoFile *protogen.File, g *protogen.GeneratedFile) {
	for _, m := range b.methods[*protoFile.Proto.Name] {
		g.P(`// validate_`, m.gwPattern, ` is an entrypoint for validating "`, m.httpMethod, `" HTTP request `)
		g.P(`// that match *.pb.gw.go/pattern_`, m.gwPattern, `.`)
		g.P(`func validate_`, m.gwPattern, `(ctx `, generateImport("Context", "context", g), `, r `, generateImport("RawMessage", "encoding/json", g), `) (err error) {`)

		if m.httpBody == "" {
			g.P(`if len(r) != 0 {`)
			g.P(`return `, generateImport("Errorf", "fmt", g), `("body is not allowed")`)
			g.P(`}`)
			g.P(`return nil`)
		} else if b.isWKT(m.inputTypeMessage.Desc) {
			g.P(`return nil`)
		} else {
			typeName := string(m.inputTypeMessage.Desc.Name())
			msg := m.inputTypeMessage.Desc
			goImportPath := string(m.inputTypeMessage.GoIdent.GoImportPath)

			if m.httpBody != "*" {
				for _, field := range m.inputTypeMessage.Fields {
					if string(field.Desc.Name()) == m.httpBody {
						typeName = string(field.Message.Desc.Name())
						msg = field.Message.Desc
						goImportPath = string(field.Message.GoIdent.GoImportPath)
						break
					}
				}
			}

			if b.isLocal(msg) {
				g.P(`return validate_Object_`, typeName, `(ctx, r, "")`)
			} else {
				nonLocalName := generateImport(typeName, goImportPath, g)
				g.P(`if validator, ok := `, b.generateAtlasValidateJSONInterfaceSignature(nonLocalName, g), `; ok {`)
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

func (b *validateBuilder) isLocal(msgType protoreflect.MessageDescriptor) bool {
	return strings.HasPrefix(string(msgType.FullName()), b.packageName)
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

func (b *validateBuilder) isWKT(msgDesc protoreflect.MessageDescriptor) bool {
	return wkt[string(msgDesc.FullName())]
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

func (b *validateBuilder) renderValidatorObjectMethods(protoFile *protogen.File, g *protogen.GeneratedFile) {
	for _, message := range protoFile.Messages {
		b.renderValidatorObjectMethod(message, g)
		b.generateValidateRequired(message, g)

		for _, innerMessage := range message.Messages {
			if innerMessage.Desc.IsMapEntry() {
				continue
			}

			b.renderValidatorObjectMethod(innerMessage, g)
			b.generateValidateRequired(innerMessage, g)
		}
	}
}

func (b *validateBuilder) renderValidatorObjectMethod(message *protogen.Message, g *protogen.GeneratedFile) {
	msgName := message.GoIdent.GoName

	g.P(`// validate_Object_`, msgName, ` function validates a JSON for a given object.`)
	g.P(`func validate_Object_`, msgName, `(ctx `, generateImport("Context", "context", g), `, r `, generateImport("RawMessage", "encoding/json", g), `, path string) (err error) {`)
	g.P(`if hook, ok := `, b.generateAtlasJSONValidateInterfaceSignature(msgName, g), `; ok {`)
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
	g.P(`if err = validate_required_Object_`, msgName, `(ctx, v, path); err != nil {`)
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

			msg := f.Desc.Message()
			msgName := f.Message.GoIdent.GoName

			if b.isWKT(msg) {
				continue
			}

			if !b.isLocal(msg) {
				nonLocalName := generateImport(msgName, string(f.Message.GoIdent.GoImportPath), g)
				g.P(`validator, ok := `, b.generateAtlasValidateJSONInterfaceSignature(nonLocalName, g))
				g.P(`if !ok {`)
				g.P(`continue`)
				g.P(`}`)
			}
			g.P(`for i, vv := range vArr {`)
			g.P(`vvPath := `, generateImport("Sprintf", "fmt", g), `("%s.[%d]", vArrPath, i)`)
			if b.isLocal(msg) {
				g.P(`if err = validate_Object_`, msgName, `(ctx, vv, vvPath); err != nil {`)
				g.P(`return err`)
				g.P(`}`)
			} else {
				g.P(`if err = validator.AtlasValidateJSON(ctx, vv, vvPath); err != nil {`)
				g.P(`return err`)
				g.P(`}`)
			}
			g.P(`}`)

		} else if f.Message != nil {
			msg := f.Desc.Message()
			msgName := f.Message.GoIdent.GoName

			if b.isWKT(msg) {
				continue
			}

			g.P(`if v[k] == nil {`)
			g.P(`continue`)
			g.P(`}`)
			g.P(`vv := v[k]`)
			g.P(`vvPath := `, generateImport("JoinPath", runtimePkgPath, g), `(path, k)`)
			if b.isLocal(msg) {
				g.P(`if err = validate_Object_`, msgName, `(ctx, vv, vvPath); err != nil {`)
				g.P(`return err`)
				g.P(`}`)
			} else {
				nonLocalName := generateImport(msgName, string(f.Message.GoIdent.GoImportPath), g)
				g.P(`validator, ok := `, b.generateAtlasValidateJSONInterfaceSignature(nonLocalName, g))
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

	g.P(`// AtlasValidateJSON function validates a JSON for object `, msgName, `.`)
	g.P(`func (_ *`, msgName, `) AtlasValidateJSON(ctx `, generateImport("Context", "context", g), `, r `, generateImport("RawMessage", "encoding/json", g), `, path string) (err error) {`)
	g.P(`if hook, ok := `, b.generateAtlasJSONValidateInterfaceSignature(msgName, g), `; ok {`)
	g.P(`if r, err = hook.AtlasJSONValidate(ctx, r, path); err != nil {`)
	g.P(`return err`)
	g.P(`}`)
	g.P(`}`)
	g.P(`return validate_Object_`, msgName, `(ctx, r, path)`)
	g.P(`}`)
	g.P()
}

//Return methods to which field marked as denied
func (b *validateBuilder) GetDeniedMethods(options []av_opts.AtlasValidateFieldOption_Operation) []string {
	httpMethods := make(map[string]struct{})
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

//Return methods to which field marked as required
func (b *validateBuilder) GetRequiredMethods(options []av_opts.AtlasValidateFieldOption_Operation) []string {
	requiredMethods := make(map[string]struct{})
	for _, op := range options {
		switch op {
		case av_opts.AtlasValidateFieldOption_create:
			requiredMethods["POST"] = struct{}{}
		case av_opts.AtlasValidateFieldOption_update:
			requiredMethods["PATCH"] = struct{}{}
		case av_opts.AtlasValidateFieldOption_replace:
			requiredMethods["PUT"] = struct{}{}
		}
	}

	uniqueMethods := make([]string, 0)
	for m := range requiredMethods {
		uniqueMethods = append(uniqueMethods, m)
	}

	sort.StringSlice(uniqueMethods).Sort()
	return uniqueMethods
}

func (b *validateBuilder) generateValidateRequired(message *protogen.Message, g *protogen.GeneratedFile) {
	requiredFields := make(map[string][]string)
	msgName := message.GoIdent.GoName

	for _, f := range message.Fields {
		fExt := proto.GetExtension(f.Desc.Options(), av_opts.E_Field)
		if fExt != nil {
			favOpt := fExt.(*av_opts.AtlasValidateFieldOption)
			methods := b.GetRequiredMethods(favOpt.GetRequired())
			if len(methods) == 0 {
				continue
			}

			requiredFields[string(f.Desc.Name())] = methods
		}
	}

	g.P(`func validate_required_Object_`, msgName, `(ctx `, generateImport("Context", "context", g), `, v map[string]`, generateImport("RawMessage", "encoding/json", g), `, path string) error {`)
	g.P(`method := `, generateImport("HTTPMethodFromContext", runtimePkgPath, g), `(ctx)`)
	g.P(`_ = method`)

	var fields []string
	for v := range requiredFields {
		fields = append(fields, v)
	}

	sort.StringSlice(fields).Sort()

	for _, field := range fields {
		methods := requiredFields[field]
		if len(methods) == 3 {
			g.P(`if _, ok := v["`, field, `"]; !ok {`)
			g.P(`path = `, generateImport("JoinPath", runtimePkgPath, g), `(path, "`, field, `")`)
			g.P(`return `, generateImport("Errorf", "fmt", g), `("field %q is required for %q operation.", path, method)`)
			g.P(`}`)
		} else {
			cond := strings.Join(methods, `" || method == "`)
			g.P(`if _, ok := v["`, field, `"]; !ok && (method == "`, cond, `") {`)
			g.P(`path = `, generateImport("JoinPath", runtimePkgPath, g), `(path, "`, field, `")`)
			g.P(`return `, generateImport("Errorf", "fmt", g), `("field %q is required for %q operation.", path, method)`)
			g.P(`}`)
		}
	}
	g.P(`return nil`)
	g.P(`}`)
}

// renderMethodDescriptors renders array of structs that are used to trigger validation
// function on correct HTTP request according to HTTP method and grpc-gateway/runtime.Pattern.
func (b *validateBuilder) renderMethodDescriptors(g *protogen.GeneratedFile) {
	g.P(`var validate_Patterns = []struct{`)
	g.P(`pattern `, generateImport("Pattern", gwruntimePkgPath, g))
	g.P(`httpMethod string`)
	g.P(`validator func(`, generateImport("Context", "context", g), `, `, generateImport("RawMessage", "encoding/json", g), `) error`)
	g.P(`// Included for introspection purpose.`)
	g.P(`allowUnknown bool`)
	g.P(`} {`)

	var files []string
	for f := range b.methods {
		files = append(files, f)
	}

	sort.StringSlice(files).Sort()

	for _, f := range files {
		methods := b.methods[f]
		g.P(`// patterns for file `, f)
		for _, m := range methods {
			g.P(`{`)
			// NOTE: pattern reiles on code generated by protoc-gen-grpc-gateway.
			g.P(`pattern: `, "pattern_"+m.gwPattern, `,`)
			g.P(`httpMethod: "`, m.httpMethod, `",`)
			g.P(`validator: `, "validate_"+m.gwPattern, `,`)
			g.P(`allowUnknown: `, m.allowUnknown, `,`)
			g.P(`},`)
		}
		g.P()
	}
	g.P(`}`)
	g.P()
}

func (b *validateBuilder) renderAnnotator(g *protogen.GeneratedFile) {
	g.P(`// AtlasValidateAnnotator parses JSON input and validates unknown fields`)
	g.P(`// based on 'allow_unknown_fields' options specified in proto file.`)
	g.P(`func AtlasValidateAnnotator(ctx `, generateImport("Context", "context", g), `, r *`,
		generateImport("Request", "net/http", g), `) `, generateImport("MD", metadataPkgPath, g), ` {`)
	g.P(`md := make(`, generateImport("MD", metadataPkgPath, g), `)`)

	g.P(`for _, v := range validate_Patterns {`)
	g.P(`if r.Method == v.httpMethod && `, generateImport("PatternMatch", runtimePkgPath, g), `(v.pattern, r.URL.Path) {`)
	g.P(`var b []byte`)
	g.P(`var err error`)
	g.P(`if b, err = `, generateImport("ReadAll", "io/ioutil", g), `(r.Body); err != nil {`)
	g.P(`md.Set("Atlas-Validation-Error", "invalid value: unable to parse body")`)
	g.P(`return md`)
	g.P(`}`)
	g.P(`r.Body = `, generateImport("NopCloser", "io/ioutil", g), `(`, generateImport("NewReader", "bytes", g), `(b))`)
	g.P(`ctx := `, generateImport("WithValue", "context", g), `(`, generateImport("WithValue", "context", g), `(`, generateImport("Background", "context", g), `(), `,
		generateImport("HTTPMethodContextKey", runtimePkgPath, g), `, r.Method), `, generateImport("AllowUnknownContextKey", runtimePkgPath, g), `, v.allowUnknown)`)
	g.P(`if err = v.validator(ctx, b); err != nil {`)
	g.P(`md.Set("Atlas-Validation-Error", err.Error())`)
	g.P(`}`)
	g.P(`break`)
	g.P(`}`)
	g.P(`}`)
	g.P(`return md`)
	g.P(`}`)
	g.P()
}
