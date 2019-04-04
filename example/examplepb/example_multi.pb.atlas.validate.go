// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: example/examplepb/example_multi.proto

package examplepb // import "github.com/infobloxopen/protoc-gen-atlas-validate/example/examplepb"

import context "context"
import fmt "fmt"
import json "encoding/json"
import runtime1 "github.com/infobloxopen/protoc-gen-atlas-validate/runtime"
import proto "github.com/gogo/protobuf/proto"
import math "math"
import _ "github.com/infobloxopen/protoc-gen-atlas-validate/options"
import _ "google.golang.org/genproto/googleapis/api/annotations"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// validate_Users2_Create2_0 is an entrypoint for validating "POST" HTTP request
// that match *.pb.gw.go/pattern_Users2_Create2_0.
func validate_Users2_Create2_0(ctx context.Context, r json.RawMessage) (err error) {
	return validate_Object_User2(ctx, r, "")
}

// validate_Object_User2 function validates a JSON for a given object.
func validate_Object_User2(ctx context.Context, r json.RawMessage, path string) (err error) {
	if hook, ok := interface{}(&User2{}).(interface {
		AtlasJSONValidate(context.Context, json.RawMessage, string) (json.RawMessage, error)
	}); ok {
		if r, err = hook.AtlasJSONValidate(ctx, r, path); err != nil {
			return err
		}
	}

	var v map[string]json.RawMessage
	if err = json.Unmarshal(r, &v); err != nil {
		return fmt.Errorf("invalid value for %q: expected object.", path)
	}

	if err = validate_required_Object_User2(ctx, v, path); err != nil {
		return err
	}

	allowUnknown := runtime1.AllowUnknownFromContext(ctx)

	for k, _ := range v {
		switch k {
		case "id":
			method := runtime1.HTTPMethodFromContext(ctx)
			if method == "POST" {
				return fmt.Errorf("field %q is unsupported for %q operation.", k, method)
			}
		case "name":
		default:
			if !allowUnknown {
				return fmt.Errorf("unknown field %q.", runtime1.JoinPath(path, k))
			}
		}
	}
	return nil
}

// AtlasValidateJSON function validates a JSON for object User2.
func (_ *User2) AtlasValidateJSON(ctx context.Context, r json.RawMessage, path string) (err error) {
	if hook, ok := interface{}(&User2{}).(interface {
		AtlasJSONValidate(context.Context, json.RawMessage, string) (json.RawMessage, error)
	}); ok {
		if r, err = hook.AtlasJSONValidate(ctx, r, path); err != nil {
			return err
		}
	}
	return validate_Object_User2(ctx, r, path)
}

func validate_required_Object_User2(ctx context.Context, v map[string]json.RawMessage, path string) error {
	method := runtime1.HTTPMethodFromContext(ctx)
	_ = method
	if _, ok := v["name"]; !ok {
		path = runtime1.JoinPath(path, "name")
		return fmt.Errorf("field %q is required for %q operation.", path, method)
	}
	return nil
}

// validate_Object_EmptyResponse2 function validates a JSON for a given object.
func validate_Object_EmptyResponse2(ctx context.Context, r json.RawMessage, path string) (err error) {
	if hook, ok := interface{}(&EmptyResponse2{}).(interface {
		AtlasJSONValidate(context.Context, json.RawMessage, string) (json.RawMessage, error)
	}); ok {
		if r, err = hook.AtlasJSONValidate(ctx, r, path); err != nil {
			return err
		}
	}

	var v map[string]json.RawMessage
	if err = json.Unmarshal(r, &v); err != nil {
		return fmt.Errorf("invalid value for %q: expected object.", path)
	}

	if err = validate_required_Object_EmptyResponse2(ctx, v, path); err != nil {
		return err
	}

	allowUnknown := runtime1.AllowUnknownFromContext(ctx)

	for k, _ := range v {
		switch k {
		default:
			if !allowUnknown {
				return fmt.Errorf("unknown field %q.", runtime1.JoinPath(path, k))
			}
		}
	}
	return nil
}

// AtlasValidateJSON function validates a JSON for object EmptyResponse2.
func (_ *EmptyResponse2) AtlasValidateJSON(ctx context.Context, r json.RawMessage, path string) (err error) {
	if hook, ok := interface{}(&EmptyResponse2{}).(interface {
		AtlasJSONValidate(context.Context, json.RawMessage, string) (json.RawMessage, error)
	}); ok {
		if r, err = hook.AtlasJSONValidate(ctx, r, path); err != nil {
			return err
		}
	}
	return validate_Object_EmptyResponse2(ctx, r, path)
}

func validate_required_Object_EmptyResponse2(ctx context.Context, v map[string]json.RawMessage, path string) error {
	method := runtime1.HTTPMethodFromContext(ctx)
	_ = method
	return nil
}
