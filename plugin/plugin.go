package plugin

import (
	"fmt"
	"strings"
        "github.com/gogo/protobuf/proto"
	// goproto "github.com/golang/protobuf/proto"
        _ "github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
        "github.com/gogo/protobuf/protoc-gen-gogo/generator"
	http_annotations "github.com/gogo/googleapis/google/api"
)

const (
	// PluginName is name of the plugin specified for protoc
	PluginName = "atlas-validate"
)

type requestDescriptor struct {
	body string
	pattern string
	method string
}


type Plugin struct {
	*generator.Generator
	requests map[string]requestDescriptor
	seen map[string]bool
	file *generator.FileDescriptor
}

func (p *Plugin) Name() string {
	return PluginName
}

func (p *Plugin) Init(g *generator.Generator) {
	p.Generator = g
	p.requests = map[string]requestDescriptor{}
	p.seen = map[string]bool{}
}

func (p *Plugin) GenerateImports(file *generator.FileDescriptor) {
	p.PrintImport("fmt", "fmt")
	p.PrintImport("strings", "strings")
	p.PrintImport("http", "net/http")
	p.PrintImport("json", "encoding/json")
	p.PrintImport("ioutil", "io/ioutil")
	p.PrintImport("context", "golang.org/x/net/context")
	p.PrintImport("metadata", "google.golang.org/grpc/metadata")

}

func (p *Plugin) Generate(file *generator.FileDescriptor) {
	p.file = file

	for _, svc := range file.GetService() {
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

			if t := p.extractType(method.GetInputType()); t != "" {
				if httpRule, ok := ext.(*http_annotations.HttpRule); ok {
					rd := requestDescriptor{
						body: httpRule.Body,
						method: p.getHttpMethod(httpRule),
						pattern: fmt.Sprintf("pattern_%s_%s_0", svc.GetName(), method.GetName()),
					}

					p.requests[t] = rd
				}
			}
		}
	}

	for t, _ := range p.requests {
		p.renderChildren(t)
		p.renderValidateJson(t)
	}

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

func (p *Plugin) renderValidateJson(msgType string) {
	msg := p.file.GetMessage(msgType)
	if msg == nil {
		return
	}
	req, ok := p.requests[msgType]
	body := req.body
	p.P(`// ValidateJSON validates JSON values for `, msg.GetName())
	p.P(`func (m *`, msg.GetName(), `) ValidateJSON(v map[string]interface{}, path string) error {`)
	p.P(`var err error`)
	p.P()
	if ok && body == "" {
		p.P(`for range v {`)
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
						p.P(`var ePath string`)
						p.P(`if path == "" {`)
						p.P(`ePath = k`)
						p.P(`} else {`)
						p.P(`ePath = path + "." + k`)
						p.P(`}`)
						//p.P(`if validator, ok := interface{}(m.Get`, generator.CamelCase(field.GetName()), `()).(interface{ ValidateJSON(map[string]interface{}, string) (error) }); ok {`)
						p.P(`if validator, ok := interface{}(&`, tn, `{}).(interface{ ValidateJSON(map[string]interface{}, string) (error) }); ok {`)
						p.P(`if vArr, ok := vv.([]interface{}); ok {`)
						p.P(`for i, vVal := range vArr {`)
						p.P(`if vVal == nil {`)
						p.P(`continue`)
						p.P(`}`)
						p.P(`aPath := fmt.Sprintf("%s.[%d]", ePath, i)`)
						p.P(`if v, ok := vVal.(map[string]interface{}); ok {`)
						p.P(`if err = validator.ValidateJSON(v, aPath); err != nil {`)
						p.P(`return err`)
						p.P(`}`)
						p.P(`} else {`)
						p.P(`return fmt.Errorf("Invalid value for %s: expected object", aPath)`)
						p.P(`}`)
						p.P(`}`)
						p.P(`} else {`)
						p.P(`return fmt.Errorf("Invalid value for %s: expected array", ePath)`)
						p.P(`}`)
						p.P(`}`)
					} else if field.IsMessage() {
						p.P(`if v[k] == nil {`)
						p.P(`continue`)
						p.P(`}`)
						p.P(`vv := v[k]`)
						p.P(`var ePath string`)
						p.P(`if path == "" {`)
						p.P(`ePath = k`)
						p.P(`} else {`)
						p.P(`ePath = path + "." + k`)
						p.P(`}`)
						p.P(`if v, ok := vv.(map[string]interface{}); ok {`)
						//p.P(`if validator, ok := interface{}(m.Get`, generator.CamelCase(field.GetName()), `()).(interface{ ValidateJSON(map[string]interface{}, string) (error) }); ok {`)
						p.P(`if validator, ok := interface{}(&`, tn, `{}).(interface{ ValidateJSON(map[string]interface{}, string) (error) }); ok {`)
						p.P(`if err = validator.ValidateJSON(v, ePath); err != nil {`)
						p.P(`return err`)
						p.P(`}`)
						p.P(`}`)
						p.P(`} else {`)
						p.P(`return fmt.Errorf("Invalid value for %s: expected object", ePath)`)
						p.P(`}`)
					}
				}
		}
		p.P(`default: return fmt.Errorf("Unknown field '%s.%s'", path, k)`)
		p.P(`}`)
		p.P(`}`)
	} else if body != "" {
		//p.P(`m.Get`, generator.CamelCase(body), `().ValidateJSON(v.(map[string]interface{}))`)
		p.P(`if validator, ok := interface{}(m.Get`, generator.CamelCase(body), `()).(interface{ ValidateJSON(map[string]interface{}, string) (error) }); ok {`)
		p.P(`if err = validator.ValidateJSON(v, path); err != nil {`)
		p.P(`return err`)
		p.P(`}`)
		p.P(`}`)
	}

	p.P(`return err`)
	p.P(`}`)
	p.P()
}

func (p *Plugin) renderAnnotator() {
	p.P(`// ValidationAnnotator function validates JSON.`)
	p.P(`func ValidationAnnotator(ctx context.Context, r *http.Request) metadata.MD {`)
	p.P(`var err error`)
	p.P(`var v map[string]interface{}`)
	p.P(`var components []string`)
	p.P(`var idx, l int`)
        p.P(`var c, verb string`)
	p.P()

	p.P(`components = strings.Split(r.URL.Path[1:], "/")`)
	p.P(`l = len(components)`)
	p.P(`if idx = strings.LastIndex(components[l-1], ":"); idx > 0 {`)
	p.P(`c = components[l-1]`)
        p.P(`components[l-1], verb = c[:idx], c[idx+1:]`)
	p.P(`}`)
	p.P()

	p.P(`md := make(metadata.MD)`)
	p.P(`defer r.Body.Close()`)
	p.P(`b, readErr := ioutil.ReadAll(r.Body)`)
	p.P(`if err != nil {`)
	p.P(`err = fmt.Errorf("Unable to parse body: %v", readErr)`)
	p.P(`goto End`)
	p.P(`}`)
	p.P()

	p.P(`if marshalErr := json.Unmarshal(b, &v); err != nil {`)
	p.P(`err = fmt.Errorf("Unable to unmarshal JSON: %v", marshalErr)`)
	p.P(`goto End`)
	p.P(`}`)

	for n, v := range p.requests {
		p.P(`if r.Method == "`, v.method, `" {`)
		p.P(`if _, matchErr := `, v.pattern, `.Match(components, verb); matchErr == nil {`)
		p.P(`err = (&`, n, `{}).ValidateJSON(v, "")`)
		p.P(`goto End`)
		p.P(`}`)
		p.P(`}`)
	}

	p.P(`End:`)
	p.P(`if err != nil {`)
	p.P(`md.Set("Validation-Error", err.Error())`)
	p.P(`}`)
	p.P(`return md`)
	p.P(`}`)
	p.P()
}

func (p *Plugin) getHttpMethod(r *http_annotations.HttpRule) string {
	switch r.GetPattern().(type) {
	case *http_annotations.HttpRule_Get: return "GET"
	case *http_annotations.HttpRule_Post: return "POST"
	case *http_annotations.HttpRule_Put: return "PUT"
	case *http_annotations.HttpRule_Delete: return "DELETE"
	}

	return ""
}
