package plugin

import (
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"

	http_opts "github.com/gogo/googleapis/google/api"
)

type httpOpt struct {
	body   string
	method string
}

func getHttpMethod(r *http_opts.HttpRule) string {
	switch r.GetPattern().(type) {
	case *http_opts.HttpRule_Get:
		return "GET"
	case *http_opts.HttpRule_Post:
		return "POST"
	case *http_opts.HttpRule_Put:
		return "PUT"
	case *http_opts.HttpRule_Delete:
		return "DELETE"
	case *http_opts.HttpRule_Patch:
		return "PATCH"
	}

	return ""
}

func extractHTTPOpts(m *descriptor.MethodDescriptorProto) []httpOpt {
	r := []httpOpt{}

	if !proto.HasExtension(m.GetOptions(), http_opts.E_Http) {
		return nil
	}

	ext, err := proto.GetExtension(m.Options, http_opts.E_Http)
	if err != nil {
		return nil
	}

	if httpRule, ok := ext.(*http_opts.HttpRule); ok {
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
