package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

type validateBuilder struct {
	methods map[string][]*methodDescriptor
}

func main() {
	// response := command.GeneratePlugin(command.Read(), plugin, ".pb.atlas.validate.go")

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

	builder := &validateBuilder{}

	for _, protoFile := range plugin.Files {
		methods := builder.gatherMethods(protoFile)
		fmt.Fprintf(os.Stderr, "methods: %#v\n", methods)
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
	return &pluginpb.CodeGeneratorResponse{}
}

func (b *validateBuilder) gatherMethods(file *protogen.File) []*methodDescriptor {
	var methods []*methodDescriptor

	for _, service := range file.Services {
		for _, method := range service.Methods {
			for i, opt := range extractHTTPOpts(method) {
				methods = append(methods, &methodDescriptor{
					svc:          string(service.Desc.Name()),
					method:       string(method.Desc.Name()),
					idx:          i,
					httpBody:     opt.body,
					httpMethod:   opt.method,
					gwPattern:    fmt.Sprintf("%s_%s_%d", service.Desc.Name(), method.Desc.Name(), i),
					inputType:    string(method.Input.Desc.Name()),
					allowUnknown: false,
				})
			}
		}
	}

	return methods
}

type methodDescriptor struct {
	svc          string
	method       string
	httpBody     string
	httpMethod   string
	gwPattern    string
	inputType    string
	idx          int
	allowUnknown bool
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
