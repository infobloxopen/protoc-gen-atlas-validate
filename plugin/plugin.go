package plugin

import (
	"fmt"
	"strings"

	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"

	av_opts "github.com/infobloxopen/protoc-gen-atlas-validate/options"
)

const (
	// PluginName is name of the plugin specified for protoc
	PluginName = "atlas-validate"
)

type Plugin struct {
	*generator.Generator
	file    *generator.FileDescriptor
	imports []string
	methods []methodDescriptor
}

func (p *Plugin) Name() string {
	return PluginName
}

func (p *Plugin) Init(g *generator.Generator) {
	p.Generator = g
	p.methods = []methodDescriptor{}
	p.imports = []string{}
}

func (p *Plugin) GenerateImports(file *generator.FileDescriptor) {
	p.PrintImport("fmt", "fmt")
	p.PrintImport("http", "net/http")
	p.PrintImport("json", "encoding/json")
	p.PrintImport("ioutil", "io/ioutil")
	p.PrintImport("bytes", "bytes")
	p.PrintImport("context", "golang.org/x/net/context")
	p.PrintImport("metadata", "google.golang.org/grpc/metadata")
	p.PrintImport("runtime", "github.com/grpc-ecosystem/grpc-gateway/runtime")
	p.PrintImport("validate_runtime", "github.com/infobloxopen/protoc-gen-atlas-validate/runtime")

	for i, v := range p.imports {
		n := fmt.Sprintf("google_protobuf%d", i+1)
		p.PrintImport(generator.GoPackageName(n), generator.GoImportPath(v))
	}
}

func (p *Plugin) Generate(file *generator.FileDescriptor) {
	p.file = file
	p.Init(p.Generator)

	p.gatherMethods()

	p.renderValidatorMethods()
	p.renderValidatorObjectMethods()
	p.renderMethodDescriptors()
	p.renderAnnotator()
}

// getAllowUnknown function picks up correct allowUnknown option from file/service/method
// hierarchy.
func (p *Plugin) getAllowUnknown(file proto.Message, svc proto.Message, method proto.Message) bool {
	var gavOpt *av_opts.AtlasValidateFileOption
	if aExt, err := proto.GetExtension(file, av_opts.E_File); err == nil && aExt != nil {
		gavOpt = aExt.(*av_opts.AtlasValidateFileOption)
	}

	var savOpt *av_opts.AtlasValidateServiceOption
	if aExt, err := proto.GetExtension(svc, av_opts.E_Service); err == nil && aExt != nil {
		savOpt = aExt.(*av_opts.AtlasValidateServiceOption)
	}

	var mavOpt *av_opts.AtlasValidateMethodOption
	if aExt, err := proto.GetExtension(method, av_opts.E_Method); err == nil && aExt != nil {
		mavOpt = aExt.(*av_opts.AtlasValidateMethodOption)
	}

	if mavOpt != nil {
		return mavOpt.GetAllowUnknownFields()
	} else if savOpt != nil {
		return savOpt.GetAllowUnknownFields()
	}

	return gavOpt.GetAllowUnknownFields()
}

type methodDescriptor struct {
	svc                  string
	method               string
	idx                  int
	httpBody, httpMethod string
	gwPattern            string
	protoInputType       string
	allowUnknown         bool
	inputType            string
}

// gatherMethods function walks through services and methods and extracts
// google.api.http options and renders different handlers for HTTP request/pattern pair.
func (p *Plugin) gatherMethods() {
	for _, svc := range p.file.GetService() {
		for _, method := range svc.GetMethod() {
			for i, opt := range extractHTTPOpts(method) {
				p.methods = append(p.methods, methodDescriptor{
					svc:            svc.GetName(),
					method:         method.GetName(),
					idx:            i,
					httpBody:       opt.body,
					httpMethod:     opt.method,
					gwPattern:      fmt.Sprintf("%s_%s_%d", svc.GetName(), method.GetName(), i),
					protoInputType: p.getProtoType(method.GetInputType()),
					inputType:      method.GetInputType(),
					allowUnknown:   p.getAllowUnknown(p.file.Options, svc.Options, method.Options),
				})
			}
		}
	}
}

func (p *Plugin) pkgPrefix() string {
	return "." + p.file.GetPackage() + "."
}

func (p *Plugin) trimPkgPrefix(t string) string {
	return strings.TrimPrefix(t, p.pkgPrefix())
}

// renderMethodDescriptors renders array of structs that are used to trigger validation
// function on correct HTTP request according to HTTP method and grpc-gateway/runtime.Pattern.
func (p *Plugin) renderMethodDescriptors() {
	p.P(`var validate_Patterns = []struct{`)
	p.P(`pattern runtime.Pattern`)
	p.P(`httpMethod string`)
	p.P(`validator func(context.Context, json.RawMessage) error`)
	p.P(`// Included for introspection purpose.`)
	p.P(`allowUnknown bool`)
	p.P(`} {`)
	for _, m := range p.methods {
		p.P(`{`)
		// NOTE: pattern reiles on code generated by protoc-gen-grpc-gateway.
		p.P(`pattern: `, "pattern_"+m.gwPattern, `,`)
		p.P(`httpMethod: "`, m.httpMethod, `",`)
		p.P(`validator: `, "validate_"+m.gwPattern, `,`)
		p.P(`allowUnknown: `, m.allowUnknown, `,`)
		p.P(`},`)
	}
	p.P(`}`)
	p.P()
}

// renderValidatorMethods function generates entrypoints for validator one per each
// HTTP request (and HTTP request additional_bindings).
func (p *Plugin) renderValidatorMethods() {
	for _, m := range p.methods {
		p.P(`// validate_`, m.gwPattern, ` is an entrypoint for validating "`, m.httpMethod, `" HTTP request `)
		p.P(`// that match *.pb.gw.go/pattern_`, m.gwPattern, `.`)
		p.P(`func validate_`, m.gwPattern, `(ctx context.Context, r json.RawMessage) (err error) {`)
		switch m.httpBody {
		case "":
			p.P(`if len(r) != 0 {`)
			p.P(`return fmt.Errorf("Body is not allowed")`)
			p.P(`}`)
			p.P(`return nil`)
		case "*":
			if m.protoInputType != "" {
				p.P(`return validate_Object_`, p.getGoType(m.protoInputType), `(ctx, r, "")`)
			} else {
				p.P(`obj := `, p.importedType(m.inputType), `{}`)
				p.P(`if validator, ok := interface{}(obj).(interface{ AtlasValidateJSON(context.Context, json.RawMessage, string) error }); ok {`)
				p.P(`return validator.AtlasValidateJSON(ctx, r, "")`)
				p.P(`}`)
				p.P(`return nil`)
			}
		default:
			msg := p.getMessage(m.inputType)
			f := msg.GetFieldDescriptor(m.httpBody)
			if p.getProtoType(f.GetTypeName()) != "" {
				if gt := p.getGoType(f.GetTypeName()); gt != "" {
					p.P(`return validate_Object_`, gt, `(ctx, r, "")`)
				}
			} else {
				p.P(`obj := `, p.importedType(f.GetTypeName()), `{}`)
				p.P(`if validator, ok := interface{}(obj).(interface{ AtlasValidateJSON(context.Context, json.RawMessage, string) error }); ok {`)
				p.P(`return validator.AtlasValidateJSON(ctx, r, "")`)
				p.P(`}`)
				p.P(`return nil`)
			}

		}
		p.P(`}`)
		p.P()
	}
}

func (p *Plugin) renderValidatorObjectMethods() {
	for _, o := range p.file.GetMessageType() {
		otype := p.getGoType(o.GetName())
		p.renderValidatorObjectMethod(o, otype)
		for _, no := range o.GetNestedType() {
			if no.GetOptions().GetMapEntry() {
				continue
			}
			p.renderValidatorObjectMethod(no, otype+"_"+p.getGoType(no.GetName()))
		}
	}
}

func (p *Plugin) renderValidatorObjectMethod(o *descriptor.DescriptorProto, t string) {
	p.P(`// validate_Object_`, t, ` function validates a JSON for a given object.`)
	p.P(`func validate_Object_`, t, `(ctx context.Context,r json.RawMessage, path string) (err error) {`)
	p.P(`obj := &`, t, `{}`)
	p.P(`if hook, ok := interface{}(obj).(interface { AtlasJSONValidate(context.Context, json.RawMessage, string) (json.RawMessage, error) }); ok {`)
	p.P(`if r, err = hook.AtlasJSONValidate(ctx, r, path); err != nil {`)
	p.P(`return err`)
	p.P(`}`)
	p.P(`}`)
	p.P(`var v map[string]json.RawMessage`)
	p.P(`if err = json.Unmarshal(r, &v); err != nil {`)
	p.P(`return fmt.Errorf("Invalid value for %q: expected object.", path)`)
	p.P(`}`)
	p.P("allowUnknown := validate_runtime.GetAllowUnknownFromContext(ctx)")

	p.P(`for k, _ := range v {`)

	p.P(`switch k {`)
	for _, f := range o.GetField() {
		p.P(`case "`, f.GetName(), `":`)
		gt := p.getGoType(f.GetTypeName())
		if p.IsMap(f) {
			continue
		}
		if fExt, err := proto.GetExtension(f.Options, av_opts.E_Field); err == nil && fExt != nil {
			favOpt := fExt.(*av_opts.AtlasValidateFieldOption)
			if favOpt.ReadOnly {
				p.P(`return fmt.Errorf("Field %s has readonly access", k)`)
			} else if  favOpt.GetDeny() == av_opts.AtlasValidateFieldOption_create  {
				p.P(`method := validate_runtime.HTTPMethodFromContext(ctx)`)
				p.P(fmt.Sprintf(`if "POST" == method {`, ))
				p.P(`return fmt.Errorf("Field %s unsupported for create method", k)`)
				p.P("}")
			} else if favOpt.GetDeny() == av_opts.AtlasValidateFieldOption_update {
				p.P(`method := validate_runtime.HTTPMethodFromContext(ctx)`)
				p.P(`if "PUT" == method || "PATCH" == method {`)
				p.P(`return fmt.Errorf("Field %s unsupported for update method", k)`)
				p.P("}")
			}
		}

		if f.IsMessage() && f.IsRepeated() {
			p.P(`if v[k] == nil {`)
			p.P(`continue`)
			p.P(`}`)
			p.P(`var vArr []json.RawMessage`)
			p.P(`vArrPath := validate_runtime.JoinPath(path, k)`)
			p.P(`if err = json.Unmarshal(v[k], &vArr); err != nil {`)
			p.P(`return fmt.Errorf("Invalid value for %q: expected array.", vArrPath)`)
			p.P(`}`)
			if gt == "" {
				p.P(`validator, ok := interface{}(&`, p.importedType(f.GetTypeName()), `{}).(interface{ AtlasValidateJSON(context.Context,json.RawMessage, string) error })`)
				p.P(`if !ok {`)
				p.P(`continue`)
				p.P(`}`)
			}
			p.P(`for i, vv := range vArr {`)
			p.P(`vvPath := fmt.Sprintf("%s.[%d]", vArrPath, i)`)
			if gt == "" {
				p.P(`if err = validator.AtlasValidateJSON(ctx, vv, vvPath); err != nil {`)
				p.P(`return err`)
				p.P(`}`)
			} else {
				p.P(`if err = validate_Object_`, gt, `(ctx, vv, vvPath); err != nil {`)
				p.P(`return err`)
				p.P(`}`)
			}
			p.P(`}`)
		} else if f.IsMessage() {
			p.P(`if v[k] == nil {`)
			p.P(`continue`)
			p.P(`}`)
			p.P(`vv := v[k]`)
			p.P(`vvPath := validate_runtime.JoinPath(path, k)`)
			if gt == "" {
				p.P(`validator, ok := interface{}(&`, p.importedType(f.GetTypeName()), `{}).(interface{ AtlasValidateJSON(context.Context, json.RawMessage, string) error })`)
				p.P(`if !ok {`)
				p.P(`continue`)
				p.P(`}`)
				p.P(`if err = validator.AtlasValidateJSON(ctx, vv, vvPath); err != nil {`)
				p.P(`return err`)
				p.P(`}`)
			} else {
				p.P(`if err = validate_Object_`, gt, `(ctx, vv, vvPath); err != nil {`)
				p.P(`return err`)
				p.P(`}`)
			}
		}
	}
	p.P(`default:`)
	p.P(`if !allowUnknown {`)
	p.P(`return fmt.Errorf("Unknown field %q", validate_runtime.JoinPath(path, k))`)
	p.P(`}`)
	p.P(`}`)
	p.P(`}`)
	p.P(`return nil`)
	p.P(`}`)
	p.P()

	p.P(`// AtlasValidateJSON function validates a JSON for object `, t, `.`)
	p.P(`func (o *`, t, `) AtlasValidateJSON(ctx context.Context, r json.RawMessage, path string) (err error) {`)
	p.P(`if hook, ok := interface{}(o).(interface { AtlasJSONValidate(context.Context,json.RawMessage, string) (json.RawMessage, error) }); ok {`)
	p.P(`if r, err = hook.AtlasJSONValidate(ctx, r, path); err != nil {`)
	p.P(`return err`)
	p.P(`}`)
	p.P(`}`)
	p.P(`return validate_Object_`, t, `(ctx, r, path)`)
	p.P(`}`)
	p.P()
}

func (p *Plugin) renderAnnotator() {
	p.P(`// AtlasValidateAnnotator parses JSON input and validates unknown fields`)
	p.P(`// based on 'allow_unknown_fields' options specified in proto file.`)
	p.P(`func AtlasValidateAnnotator(ctx context.Context, r *http.Request) metadata.MD {`)
	p.P(`md := make(metadata.MD)`)

	p.P(`for _, v := range validate_Patterns {`)
	p.P(`if r.Method == v.httpMethod && validate_runtime.PatternMatch(v.pattern, r.URL.Path) {`)
	p.P(`var b []byte`)
	p.P(`var err error`)
	p.P(`if b, err = ioutil.ReadAll(r.Body); err != nil {`)
	p.P(`md.Set("Atlas-Validation-Error", "Invalid value: unable to parse body")`)
	p.P(`return md`)
	p.P(`}`)
	p.P(`r.Body = ioutil.NopCloser(bytes.NewReader(b))`)
	p.P(`ctx := context.WithValue(context.WithValue(context.Background(), "http-method", r.Method), "allow-unknown", v.allowUnknown)`)
	p.P(`if err = v.validator(ctx, b); err != nil {`)
	p.P(`md.Set("Atlas-Validation-Error", err.Error())`)
	p.P(`}`)
	p.P(`break`)
	p.P(`}`)
	p.P(`}`)
	p.P(`return md`)
	p.P(`}`)
	p.P()
}

func (p *Plugin) getMessage(t string) *descriptor.DescriptorProto {
	var local bool

	if strings.HasPrefix(t, p.pkgPrefix()) {
		local = true
	}

	if msg := p.file.GetMessage(p.trimPkgPrefix(t)); msg == nil && local {
		return nil
	} else if msg != nil {
		return msg
	}

	file := p.ObjectNamed(t).File()
	return file.GetMessage(strings.TrimPrefix(t, "."+file.GetPackage()+"."))
}

