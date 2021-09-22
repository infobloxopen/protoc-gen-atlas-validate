package examplepb

import (
	context "context"
	json "encoding/json"
	fmt "fmt"
	runtime "github.com/infobloxopen/protoc-gen-atlas-validate/runtime"
)

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

	allowUnknown := runtime.AllowUnknownFromContext(ctx)

	for k, _ := range v {
		switch k {
		case "id":
			method := runtime.HTTPMethodFromContext(ctx)
			if method == "POST" {
				return fmt.Errorf("field %q is unsupported for %q operation.", k, method)
			}
		case "name":
		default:
			if !allowUnknown {
				return fmt.Errorf("unknown field %q.", runtime.JoinPath(path, k))
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
	method := runtime.HTTPMethodFromContext(ctx)
	_ = method
	if _, ok := v["name"]; !ok {
		path = runtime.JoinPath(path, "name")
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

	allowUnknown := runtime.AllowUnknownFromContext(ctx)

	for k, _ := range v {
		switch k {
		default:
			if !allowUnknown {
				return fmt.Errorf("unknown field %q.", runtime.JoinPath(path, k))
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
	method := runtime.HTTPMethodFromContext(ctx)
	_ = method
	return nil
}
