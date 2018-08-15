package plugin

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"strings"
	// goproto "github.com/golang/protobuf/proto"
	av_opts "github.com/askurydzin/protoc-gen-atlas-validate/options"
	http_annotations "github.com/gogo/googleapis/google/api"
	_ "github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
)

const (
	// PluginName is name of the plugin specified for protoc
	PluginName = "atlas-validate"
)

type requestDescriptor struct {
	body         string
	name         string
	pattern      string
	method       string
	allowUnknown bool
}

type Plugin struct {
	*generator.Generator
	requests       []requestDescriptor
	seen           map[string]bool
	seenValidators map[string]bool
	file           *generator.FileDescriptor
}

func (p *Plugin) Name() string {
	return PluginName
}

func (p *Plugin) Init(g *generator.Generator) {
	p.Generator = g
	p.requests = []requestDescriptor{}
	p.seen = map[string]bool{}
	p.seenValidators = map[string]bool{}
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
	p.PrintImport("validate_runtime", "github.com/askurydzin/protoc-gen-atlas-validate/runtime")

}

func (p *Plugin) Generate(file *generator.FileDescriptor) {
	p.file = file

	var gavOpt *av_opts.AtlasValidateFileOption
	if proto.HasExtension(file.GetOptions(), av_opts.E_File) {
		if aExt, err := proto.GetExtension(file.Options, av_opts.E_File); err == nil && aExt != nil {
			gavOpt = aExt.(*av_opts.AtlasValidateFileOption)
		}
	}

	for _, svc := range file.GetService() {

		var savOpt *av_opts.AtlasValidateServiceOption
		if proto.HasExtension(svc.GetOptions(), av_opts.E_Service) {
			if aExt, err := proto.GetExtension(svc.Options, av_opts.E_Service); err == nil && aExt != nil {
				savOpt = aExt.(*av_opts.AtlasValidateServiceOption)
			}
		}

		for _, method := range svc.GetMethod() {
			if method.GetOptions() == nil {
				continue
			}

			if !proto.HasExtension(method.GetOptions(), http_annotations.E_Http) {
				continue
			}

			ext, err := proto.GetExtension(method.Options, http_annotations.E_Http)
			if err != nil {
				continue
			}

			var mavOpt *av_opts.AtlasValidateMethodOption
			if aExt, err := proto.GetExtension(method.Options, av_opts.E_Method); err == nil && aExt != nil {
				mavOpt = aExt.(*av_opts.AtlasValidateMethodOption)
			}

			if t := p.extractType(method.GetInputType()); t != "" {
				if httpRule, ok := ext.(*http_annotations.HttpRule); ok {
					rd := requestDescriptor{
						name:    t,
						body:    httpRule.Body,
						method:  p.getHttpMethod(httpRule),
						pattern: fmt.Sprintf("pattern_%s_%s_0", svc.GetName(), method.GetName()),
					}

					if mavOpt != nil {
						rd.allowUnknown = mavOpt.GetAllowUnknownFields()
					} else if savOpt != nil {
						rd.allowUnknown = savOpt.GetAllowUnknownFields()
					} else {
						rd.allowUnknown = gavOpt.GetAllowUnknownFields()
					}

					p.requests = append(p.requests, rd)

					for i, httpRule := range httpRule.GetAdditionalBindings() {
						rd := requestDescriptor{
							name:    t,
							body:    httpRule.Body,
							method:  p.getHttpMethod(httpRule),
							pattern: fmt.Sprintf("pattern_%s_%s_%d", svc.GetName(), method.GetName(), i+1),
						}

						if mavOpt != nil {
							rd.allowUnknown = mavOpt.GetAllowUnknownFields()
						} else if savOpt != nil {
							rd.allowUnknown = savOpt.GetAllowUnknownFields()
						} else {
							rd.allowUnknown = gavOpt.GetAllowUnknownFields()
						}

						p.requests = append(p.requests, rd)
					}

				}
			}
		}
	}

	for _, v := range p.requests {
		p.renderChildren(v.name)
		p.renderValidateJson(v.name)
	}

	p.renderPatterns()

	p.renderAnnotator()
}

func (p *Plugin) renderChildren(t string) {
	if p.seen[t] {
		return
	}

	msg := p.file.GetMessage(t)
	if msg == nil {
		return
	} else {
		p.seen[t] = true
	}

	for _, field := range msg.GetField() {
		if t := p.extractType(field.GetTypeName()); t != "" {
			if !p.seen[t] {
				p.renderChildren(t)
				p.renderValidateJson(t)
				if p.file.GetMessage(t) != nil {
					p.seen[t] = true
				}
			}
		}
	}
}

func (p *Plugin) extractType(t string) string {
	if strings.HasPrefix(t, ".") {
		s := strings.Split(t, ".")
		return s[len(s)-1]
	}

	return ""
}

func (p *Plugin) findReq(n string) (requestDescriptor, bool) {
	for _, v := range p.requests {
		if n == v.name {
			return v, true
		}
	}

	return requestDescriptor{}, false
}

func (p *Plugin) renderValidateJson(msgType string) {
	if p.seenValidators[msgType] {
		return
	} else {
		p.seenValidators[msgType] = true
	}
	msg := p.file.GetMessage(msgType)
	if msg == nil {
		return
	}
	req, ok := p.findReq(msgType)
	body := req.body
	p.P(`// ValidateJSON validates JSON values for `, msg.GetName())
	p.P(`func Default`, msg.GetName(), `ValidateJSON(v map[string]interface{}, path string) (err error) {`)
	p.P()
	if ok && body == "" {
		p.P(`if v != nil {`)
		p.P(`err = fmt.Errorf("Body is not allowed")`)
		p.P(`return err`)
		p.P(`}`)
	} else if !ok || body == "*" {
		p.P(`for k, _ := range v {`)
		p.P(`switch k {`)
		for _, field := range msg.GetField() {
			p.P(`case "`, field.GetName(), `":`)
			if tn := p.extractType(field.GetTypeName()); tn != "" && p.seen[tn] {
				if field.IsRepeated() && field.IsMessage() {
					p.P(`if v[k] == nil {`)
					p.P(`continue`)
					p.P(`}`)
					p.P(`vv := v[k]`)
					p.P(`if vArr, ok := vv.([]interface{}); ok {`)
					p.P(`if validator, ok := interface{}(&`, tn, `{}).(interface{ ValidateJSON(map[string]interface{}, string) (error) }); ok {`)
					p.P(`for i, vVal := range vArr {`)
					p.P(`if vVal == nil {`)
					p.P(`continue`)
					p.P(`}`)
					p.P(`aPath := fmt.Sprintf("%s.[%d]", validate_runtime.JoinPath(path, k), i)`)
					p.P(`if v, ok := vVal.(map[string]interface{}); ok {`)
					p.P(`if err = validator.ValidateJSON(v, aPath); err != nil {`)
					p.P(`return err`)
					p.P(`}`)
					p.P(`} else {`)
					p.P(`return fmt.Errorf("Invalid value for %q: expected object", aPath)`)
					p.P(`}`)
					p.P(`}`)
					p.P(`} else {`)
					p.P(`for i, vVal := range vArr {`)
					p.P(`if vVal == nil {`)
					p.P(`continue`)
					p.P(`}`)
					p.P(`aPath := fmt.Sprintf("%s.[%d]", validate_runtime.JoinPath(path, k), i)`)
					p.P(`if v, ok := vVal.(map[string]interface{}); ok {`)
					p.P(`if err = Default`, tn, `ValidateJSON(v, aPath); err != nil {`)
					p.P(`return err`)
					p.P(`}`)
					p.P(`} else {`)
					p.P(`return fmt.Errorf("Invalid value for %q: expected object", aPath)`)
					p.P(`}`)
					p.P(`}`)
					p.P(`}`)
					p.P(`} else {`)
					p.P(`return fmt.Errorf("Invalid value for %q: expected array", validate_runtime.JoinPath(path, k))`)
					p.P(`}`)
				} else if field.IsMessage() {
					p.P(`if v[k] == nil {`)
					p.P(`continue`)
					p.P(`}`)
					p.P(`vv := v[k]`)
					p.P(`if v, ok := vv.(map[string]interface{}); ok {`)
					p.P(`if validator, ok := interface{}(&`, tn, `{}).(interface{ ValidateJSON(map[string]interface{}, string) (error) }); ok {`)
					p.P(`if err = validator.ValidateJSON(v, validate_runtime.JoinPath(path, k)); err != nil {`)
					p.P(`return err`)
					p.P(`}`)
					p.P(`} else {`)
					p.P(`if err = Default`, tn, `ValidateJSON(v, validate_runtime.JoinPath(path, k)); err != nil {`)
					p.P(`return err`)
					p.P(`}`)
					p.P(`}`)
					p.P(`} else {`)
					p.P(`return fmt.Errorf("Invalid value for %q: expected object", validate_runtime.JoinPath(path, k))`)
					p.P(`}`)
				}
			}
		}
		p.P(`default:`)
		if req.allowUnknown {
			p.P(`continue`)
		} else {
			p.P(`return fmt.Errorf("Unknown field %q", validate_runtime.JoinPath(path, k))`)
		}
		p.P(`}`)
		p.P(`}`)
	} else if body != "" {
		if tn := p.extractType(msg.GetFieldDescriptor(body).GetTypeName()); tn != "" && p.seen[tn] {
			p.P(`if validator, ok := interface{}(&(`, tn, `{})).(interface{ ValidateJSON(map[string]interface{}, string) (error) }); ok {`)
			p.P(`if err = validator.ValidateJSON(v, path); err != nil {`)
			p.P(`return err`)
			p.P(`}`)
			p.P(`} else {`)
			p.P(`if err = Default`, tn, `ValidateJSON(v, path); err != nil {`)
			p.P(`return err`)
			p.P(`}`)
			p.P(`}`)
		}
	}

	p.P(`return err`)
	p.P(`}`)
	p.P()
}

func (p *Plugin) renderPatterns() {
	p.P(`var patterns = []struct{ `)
	p.P(`method string`)
	p.P(`pattern runtime.Pattern`)
	p.P(`validator func(map[string]interface{}, string) error`)
	p.P(`}{`)
	for _, v := range p.requests {
		if v.allowUnknown {
			p.P(`// {`)
			p.P(`// method: "`, v.method, `",`)
			p.P(`// pattern: `, v.pattern, `,`)
			p.P(`// validator: Default`, v.name, `ValidateJSON,`)
			p.P(`// },`)
		} else {
			p.P(`{`)
			p.P(`method: "`, v.method, `",`)
			p.P(`pattern: `, v.pattern, `,`)
			p.P(`validator: Default`, v.name, `ValidateJSON,`)
			p.P(`},`)
		}
	}
	p.P(`}`)
}

func (p *Plugin) renderAnnotator() {
	p.P(`// ValidationAnnotator function validates JSON.`)
	p.P(`func ValidationAnnotator(ctx context.Context, r *http.Request) metadata.MD {`)
	p.P(`var jv map[string]interface{}`)
	p.P()

	p.P(`md := make(metadata.MD)`)
	p.P(`if len(patterns) == 0 {`)
	p.P(`return md`)
	p.P(`}`)
	p.P(`b, err := ioutil.ReadAll(r.Body)`)
	p.P(`r.Body = ioutil.NopCloser(bytes.NewReader(b))`)
	p.P(`if err != nil {`)
	p.P(`md.Set("Atlas-Validation-Error", fmt.Sprintf("Unable to read JSON request"))`)
	p.P(`return md`)
	p.P(`} else if err := json.Unmarshal(b, &jv); err != nil {`)
	p.P(`if len(b) != 0 {`)
	p.P(`md.Set("Atlas-Validation-Error", fmt.Sprintf("Unable to parse JSON request"))`)
	p.P(`return md`)
	p.P(`}`)
	p.P(`}`)
	p.P()

	p.P(`for _, v := range patterns {`)
	p.P(`if r.Method == v.method && validate_runtime.PatternMatch(v.pattern, r.URL.Path) {`)
	p.P(`if err := v.validator(jv, ""); err != nil {`)
	p.P(`md.Set("Atlas-Validation-Error", err.Error())`)
	p.P(`}`)
	p.P(`break`)
	p.P(`}`)
	p.P(`}`)
	p.P(`return md`)
	p.P(`}`)
	p.P()
}

func (p *Plugin) getHttpMethod(r *http_annotations.HttpRule) string {
	switch r.GetPattern().(type) {
	case *http_annotations.HttpRule_Get:
		return "GET"
	case *http_annotations.HttpRule_Post:
		return "POST"
	case *http_annotations.HttpRule_Put:
		return "PUT"
	case *http_annotations.HttpRule_Delete:
		return "DELETE"
	case *http_annotations.HttpRule_Patch:
		return "PATCH"
	}

	return ""
}
