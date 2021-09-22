package examplepb

import (
	context "context"
	json "encoding/json"
	fmt "fmt"
	external "github.com/infobloxopen/protoc-gen-atlas-validate/example/external"
	runtime "github.com/infobloxopen/protoc-gen-atlas-validate/runtime"
)

// validate_Users_Create_0 is an entrypoint for validating "POST" HTTP request
// that match *.pb.gw.go/pattern_Users_Create_0.
func validate_Users_Create_0(ctx context.Context, r json.RawMessage) (err error) {
	return validate_Object_User(ctx, r, "")
}

// validate_Users_Update_0 is an entrypoint for validating "PUT" HTTP request
// that match *.pb.gw.go/pattern_Users_Update_0.
func validate_Users_Update_0(ctx context.Context, r json.RawMessage) (err error) {
	return validate_Object_User(ctx, r, "")
}

// validate_Users_Update_1 is an entrypoint for validating "PATCH" HTTP request
// that match *.pb.gw.go/pattern_Users_Update_1.
func validate_Users_Update_1(ctx context.Context, r json.RawMessage) (err error) {
	return validate_Object_User(ctx, r, "")
}

// validate_Users_List_0 is an entrypoint for validating "GET" HTTP request
// that match *.pb.gw.go/pattern_Users_List_0.
func validate_Users_List_0(ctx context.Context, r json.RawMessage) (err error) {
	if len(r) != 0 {
		return fmt.Errorf("body is not allowed")
	}
	return nil
}

// validate_Users_List_1 is an entrypoint for validating "GET" HTTP request
// that match *.pb.gw.go/pattern_Users_List_1.
func validate_Users_List_1(ctx context.Context, r json.RawMessage) (err error) {
	if len(r) != 0 {
		return fmt.Errorf("body is not allowed")
	}
	return nil
}

// validate_Users_UpdateExternalUser_0 is an entrypoint for validating "PUT" HTTP request
// that match *.pb.gw.go/pattern_Users_UpdateExternalUser_0.
func validate_Users_UpdateExternalUser_0(ctx context.Context, r json.RawMessage) (err error) {
	if validator, ok := interface{}(&external.ExternalUser{}).(interface {
		AtlasValidateJSON(context.Context, json.RawMessage, string) error
	}); ok {
		return validator.AtlasValidateJSON(ctx, r, "")
	}
	return nil
}

// validate_Users_UpdateExternalUser2_0 is an entrypoint for validating "PUT" HTTP request
// that match *.pb.gw.go/pattern_Users_UpdateExternalUser2_0.
func validate_Users_UpdateExternalUser2_0(ctx context.Context, r json.RawMessage) (err error) {
	if validator, ok := interface{}(&external.ExternalUser{}).(interface {
		AtlasValidateJSON(context.Context, json.RawMessage, string) error
	}); ok {
		return validator.AtlasValidateJSON(ctx, r, "")
	}
	return nil
}

// validate_Profiles_Create_0 is an entrypoint for validating "POST" HTTP request
// that match *.pb.gw.go/pattern_Profiles_Create_0.
func validate_Profiles_Create_0(ctx context.Context, r json.RawMessage) (err error) {
	return validate_Object_Profile(ctx, r, "")
}

// validate_Profiles_Update_0 is an entrypoint for validating "PUT" HTTP request
// that match *.pb.gw.go/pattern_Profiles_Update_0.
func validate_Profiles_Update_0(ctx context.Context, r json.RawMessage) (err error) {
	return validate_Object_Profile(ctx, r, "")
}

// validate_Groups_Create_0 is an entrypoint for validating "POST" HTTP request
// that match *.pb.gw.go/pattern_Groups_Create_0.
func validate_Groups_Create_0(ctx context.Context, r json.RawMessage) (err error) {
	return validate_Object_Group(ctx, r, "")
}

// validate_Groups_Update_0 is an entrypoint for validating "PUT" HTTP request
// that match *.pb.gw.go/pattern_Groups_Update_0.
func validate_Groups_Update_0(ctx context.Context, r json.RawMessage) (err error) {
	return validate_Object_Group(ctx, r, "")
}

// validate_Groups_ValidatedList_0 is an entrypoint for validating "GET" HTTP request
// that match *.pb.gw.go/pattern_Groups_ValidatedList_0.
func validate_Groups_ValidatedList_0(ctx context.Context, r json.RawMessage) (err error) {
	if len(r) != 0 {
		return fmt.Errorf("body is not allowed")
	}
	return nil
}

// validate_Groups_ValidatedList_1 is an entrypoint for validating "GET" HTTP request
// that match *.pb.gw.go/pattern_Groups_ValidatedList_1.
func validate_Groups_ValidatedList_1(ctx context.Context, r json.RawMessage) (err error) {
	if len(r) != 0 {
		return fmt.Errorf("body is not allowed")
	}
	return nil
}

// validate_Groups_ValidateWKT_0 is an entrypoint for validating "PUT" HTTP request
// that match *.pb.gw.go/pattern_Groups_ValidateWKT_0.
func validate_Groups_ValidateWKT_0(ctx context.Context, r json.RawMessage) (err error) {
	return nil
}

// validate_Groups_ValidateWKT_1 is an entrypoint for validating "PUT" HTTP request
// that match *.pb.gw.go/pattern_Groups_ValidateWKT_1.
func validate_Groups_ValidateWKT_1(ctx context.Context, r json.RawMessage) (err error) {
	return nil
}

// validate_Object_User function validates a JSON for a given object.
func validate_Object_User(ctx context.Context, r json.RawMessage, path string) (err error) {
	if hook, ok := interface{}(&User{}).(interface {
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

	if err = validate_required_Object_User(ctx, v, path); err != nil {
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
		case "profile":
			if v[k] == nil {
				continue
			}
			vv := v[k]
			vvPath := runtime.JoinPath(path, k)
			if err = validate_Object_Profile(ctx, vv, vvPath); err != nil {
				return err
			}
		case "address":
			if v[k] == nil {
				continue
			}
			vv := v[k]
			vvPath := runtime.JoinPath(path, k)
			if err = validate_Object_Address(ctx, vv, vvPath); err != nil {
				return err
			}
		case "groups":
			if v[k] == nil {
				continue
			}
			var vArr []json.RawMessage
			vArrPath := runtime.JoinPath(path, k)
			if err = json.Unmarshal(v[k], &vArr); err != nil {
				return fmt.Errorf("invalid value for %q: expected array.", vArrPath)
			}
			for i, vv := range vArr {
				vvPath := fmt.Sprintf("%s.[%d]", vArrPath, i)
				if err = validate_Object_Group(ctx, vv, vvPath); err != nil {
					return err
				}
			}
		case "parents":
			if v[k] == nil {
				continue
			}
			var vArr []json.RawMessage
			vArrPath := runtime.JoinPath(path, k)
			if err = json.Unmarshal(v[k], &vArr); err != nil {
				return fmt.Errorf("invalid value for %q: expected array.", vArrPath)
			}
			for i, vv := range vArr {
				vvPath := fmt.Sprintf("%s.[%d]", vArrPath, i)
				if err = validate_Object_User_Parent(ctx, vv, vvPath); err != nil {
					return err
				}
			}
		case "external_user":
			if v[k] == nil {
				continue
			}
			vv := v[k]
			vvPath := runtime.JoinPath(path, k)
			validator, ok := interface{}(&external.ExternalUser{}).(interface {
				AtlasValidateJSON(context.Context, json.RawMessage, string) error
			})
			if !ok {
				continue
			}
			if err = validator.AtlasValidateJSON(ctx, vv, vvPath); err != nil {
				return err
			}
		case "empty_list":
			if v[k] == nil {
				continue
			}
			var vArr []json.RawMessage
			vArrPath := runtime.JoinPath(path, k)
			if err = json.Unmarshal(v[k], &vArr); err != nil {
				return fmt.Errorf("invalid value for %q: expected array.", vArrPath)
			}
		case "timestamp":
		default:
			if !allowUnknown {
				return fmt.Errorf("unknown field %q.", runtime.JoinPath(path, k))
			}
		}
	}
	return nil
}

// AtlasValidateJSON function validates a JSON for object User.
func (_ *User) AtlasValidateJSON(ctx context.Context, r json.RawMessage, path string) (err error) {
	if hook, ok := interface{}(&User{}).(interface {
		AtlasJSONValidate(context.Context, json.RawMessage, string) (json.RawMessage, error)
	}); ok {
		if r, err = hook.AtlasJSONValidate(ctx, r, path); err != nil {
			return err
		}
	}
	return validate_Object_User(ctx, r, path)
}

func validate_required_Object_User(ctx context.Context, v map[string]json.RawMessage, path string) error {
	method := runtime.HTTPMethodFromContext(ctx)
	_ = method
	if _, ok := v["name"]; !ok {
		path = runtime.JoinPath(path, "name")
		return fmt.Errorf("field %q is required for %q operation.", path, method)
	}
	return nil
}

// validate_Object_User_Parent function validates a JSON for a given object.
func validate_Object_User_Parent(ctx context.Context, r json.RawMessage, path string) (err error) {
	if hook, ok := interface{}(&User_Parent{}).(interface {
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

	if err = validate_required_Object_User_Parent(ctx, v, path); err != nil {
		return err
	}

	allowUnknown := runtime.AllowUnknownFromContext(ctx)

	for k, _ := range v {
		switch k {
		case "name":
		default:
			if !allowUnknown {
				return fmt.Errorf("unknown field %q.", runtime.JoinPath(path, k))
			}
		}
	}
	return nil
}

// AtlasValidateJSON function validates a JSON for object User_Parent.
func (_ *User_Parent) AtlasValidateJSON(ctx context.Context, r json.RawMessage, path string) (err error) {
	if hook, ok := interface{}(&User_Parent{}).(interface {
		AtlasJSONValidate(context.Context, json.RawMessage, string) (json.RawMessage, error)
	}); ok {
		if r, err = hook.AtlasJSONValidate(ctx, r, path); err != nil {
			return err
		}
	}
	return validate_Object_User_Parent(ctx, r, path)
}

func validate_required_Object_User_Parent(ctx context.Context, v map[string]json.RawMessage, path string) error {
	method := runtime.HTTPMethodFromContext(ctx)
	_ = method
	return nil
}

// validate_Object_Address function validates a JSON for a given object.
func validate_Object_Address(ctx context.Context, r json.RawMessage, path string) (err error) {
	if hook, ok := interface{}(&Address{}).(interface {
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

	if err = validate_required_Object_Address(ctx, v, path); err != nil {
		return err
	}

	allowUnknown := runtime.AllowUnknownFromContext(ctx)

	for k, _ := range v {
		switch k {
		case "country":
		case "state":
			method := runtime.HTTPMethodFromContext(ctx)
			if method == "PATCH" || method == "POST" || method == "PUT" {
				return fmt.Errorf("field %q is unsupported for %q operation.", k, method)
			}
		case "city":
		case "zip":
		case "tags":
		default:
			if !allowUnknown {
				return fmt.Errorf("unknown field %q.", runtime.JoinPath(path, k))
			}
		}
	}
	return nil
}

// AtlasValidateJSON function validates a JSON for object Address.
func (_ *Address) AtlasValidateJSON(ctx context.Context, r json.RawMessage, path string) (err error) {
	if hook, ok := interface{}(&Address{}).(interface {
		AtlasJSONValidate(context.Context, json.RawMessage, string) (json.RawMessage, error)
	}); ok {
		if r, err = hook.AtlasJSONValidate(ctx, r, path); err != nil {
			return err
		}
	}
	return validate_Object_Address(ctx, r, path)
}

func validate_required_Object_Address(ctx context.Context, v map[string]json.RawMessage, path string) error {
	method := runtime.HTTPMethodFromContext(ctx)
	_ = method
	return nil
}

// validate_Object_Group function validates a JSON for a given object.
func validate_Object_Group(ctx context.Context, r json.RawMessage, path string) (err error) {
	if hook, ok := interface{}(&Group{}).(interface {
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

	if err = validate_required_Object_Group(ctx, v, path); err != nil {
		return err
	}

	allowUnknown := runtime.AllowUnknownFromContext(ctx)

	for k, _ := range v {
		switch k {
		case "id":
		case "name":
		case "notes":
		default:
			if !allowUnknown {
				return fmt.Errorf("unknown field %q.", runtime.JoinPath(path, k))
			}
		}
	}
	return nil
}

// AtlasValidateJSON function validates a JSON for object Group.
func (_ *Group) AtlasValidateJSON(ctx context.Context, r json.RawMessage, path string) (err error) {
	if hook, ok := interface{}(&Group{}).(interface {
		AtlasJSONValidate(context.Context, json.RawMessage, string) (json.RawMessage, error)
	}); ok {
		if r, err = hook.AtlasJSONValidate(ctx, r, path); err != nil {
			return err
		}
	}
	return validate_Object_Group(ctx, r, path)
}

func validate_required_Object_Group(ctx context.Context, v map[string]json.RawMessage, path string) error {
	method := runtime.HTTPMethodFromContext(ctx)
	_ = method
	if _, ok := v["id"]; !ok && (method == "PATCH" || method == "PUT") {
		path = runtime.JoinPath(path, "id")
		return fmt.Errorf("field %q is required for %q operation.", path, method)
	}
	if _, ok := v["name"]; !ok && (method == "POST") {
		path = runtime.JoinPath(path, "name")
		return fmt.Errorf("field %q is required for %q operation.", path, method)
	}
	return nil
}

// validate_Object_CreateUserRequest function validates a JSON for a given object.
func validate_Object_CreateUserRequest(ctx context.Context, r json.RawMessage, path string) (err error) {
	if hook, ok := interface{}(&CreateUserRequest{}).(interface {
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

	if err = validate_required_Object_CreateUserRequest(ctx, v, path); err != nil {
		return err
	}

	allowUnknown := runtime.AllowUnknownFromContext(ctx)

	for k, _ := range v {
		switch k {
		case "payload":
			if v[k] == nil {
				continue
			}
			vv := v[k]
			vvPath := runtime.JoinPath(path, k)
			if err = validate_Object_User(ctx, vv, vvPath); err != nil {
				return err
			}
		default:
			if !allowUnknown {
				return fmt.Errorf("unknown field %q.", runtime.JoinPath(path, k))
			}
		}
	}
	return nil
}

// AtlasValidateJSON function validates a JSON for object CreateUserRequest.
func (_ *CreateUserRequest) AtlasValidateJSON(ctx context.Context, r json.RawMessage, path string) (err error) {
	if hook, ok := interface{}(&CreateUserRequest{}).(interface {
		AtlasJSONValidate(context.Context, json.RawMessage, string) (json.RawMessage, error)
	}); ok {
		if r, err = hook.AtlasJSONValidate(ctx, r, path); err != nil {
			return err
		}
	}
	return validate_Object_CreateUserRequest(ctx, r, path)
}

func validate_required_Object_CreateUserRequest(ctx context.Context, v map[string]json.RawMessage, path string) error {
	method := runtime.HTTPMethodFromContext(ctx)
	_ = method
	return nil
}

// validate_Object_UpdateUserRequest function validates a JSON for a given object.
func validate_Object_UpdateUserRequest(ctx context.Context, r json.RawMessage, path string) (err error) {
	if hook, ok := interface{}(&UpdateUserRequest{}).(interface {
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

	if err = validate_required_Object_UpdateUserRequest(ctx, v, path); err != nil {
		return err
	}

	allowUnknown := runtime.AllowUnknownFromContext(ctx)

	for k, _ := range v {
		switch k {
		case "payload":
			if v[k] == nil {
				continue
			}
			vv := v[k]
			vvPath := runtime.JoinPath(path, k)
			if err = validate_Object_User(ctx, vv, vvPath); err != nil {
				return err
			}
		default:
			if !allowUnknown {
				return fmt.Errorf("unknown field %q.", runtime.JoinPath(path, k))
			}
		}
	}
	return nil
}

// AtlasValidateJSON function validates a JSON for object UpdateUserRequest.
func (_ *UpdateUserRequest) AtlasValidateJSON(ctx context.Context, r json.RawMessage, path string) (err error) {
	if hook, ok := interface{}(&UpdateUserRequest{}).(interface {
		AtlasJSONValidate(context.Context, json.RawMessage, string) (json.RawMessage, error)
	}); ok {
		if r, err = hook.AtlasJSONValidate(ctx, r, path); err != nil {
			return err
		}
	}
	return validate_Object_UpdateUserRequest(ctx, r, path)
}

func validate_required_Object_UpdateUserRequest(ctx context.Context, v map[string]json.RawMessage, path string) error {
	method := runtime.HTTPMethodFromContext(ctx)
	_ = method
	return nil
}

// validate_Object_EmptyRequest function validates a JSON for a given object.
func validate_Object_EmptyRequest(ctx context.Context, r json.RawMessage, path string) (err error) {
	if hook, ok := interface{}(&EmptyRequest{}).(interface {
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

	if err = validate_required_Object_EmptyRequest(ctx, v, path); err != nil {
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

// AtlasValidateJSON function validates a JSON for object EmptyRequest.
func (_ *EmptyRequest) AtlasValidateJSON(ctx context.Context, r json.RawMessage, path string) (err error) {
	if hook, ok := interface{}(&EmptyRequest{}).(interface {
		AtlasJSONValidate(context.Context, json.RawMessage, string) (json.RawMessage, error)
	}); ok {
		if r, err = hook.AtlasJSONValidate(ctx, r, path); err != nil {
			return err
		}
	}
	return validate_Object_EmptyRequest(ctx, r, path)
}

func validate_required_Object_EmptyRequest(ctx context.Context, v map[string]json.RawMessage, path string) error {
	method := runtime.HTTPMethodFromContext(ctx)
	_ = method
	return nil
}

// validate_Object_EmptyResponse function validates a JSON for a given object.
func validate_Object_EmptyResponse(ctx context.Context, r json.RawMessage, path string) (err error) {
	if hook, ok := interface{}(&EmptyResponse{}).(interface {
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

	if err = validate_required_Object_EmptyResponse(ctx, v, path); err != nil {
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

// AtlasValidateJSON function validates a JSON for object EmptyResponse.
func (_ *EmptyResponse) AtlasValidateJSON(ctx context.Context, r json.RawMessage, path string) (err error) {
	if hook, ok := interface{}(&EmptyResponse{}).(interface {
		AtlasJSONValidate(context.Context, json.RawMessage, string) (json.RawMessage, error)
	}); ok {
		if r, err = hook.AtlasJSONValidate(ctx, r, path); err != nil {
			return err
		}
	}
	return validate_Object_EmptyResponse(ctx, r, path)
}

func validate_required_Object_EmptyResponse(ctx context.Context, v map[string]json.RawMessage, path string) error {
	method := runtime.HTTPMethodFromContext(ctx)
	_ = method
	return nil
}

// validate_Object_Profile function validates a JSON for a given object.
func validate_Object_Profile(ctx context.Context, r json.RawMessage, path string) (err error) {
	if hook, ok := interface{}(&Profile{}).(interface {
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

	if err = validate_required_Object_Profile(ctx, v, path); err != nil {
		return err
	}

	allowUnknown := runtime.AllowUnknownFromContext(ctx)

	for k, _ := range v {
		switch k {
		case "id":
		case "name":
			method := runtime.HTTPMethodFromContext(ctx)
			if method == "PATCH" || method == "PUT" {
				return fmt.Errorf("field %q is unsupported for %q operation.", k, method)
			}
		case "notes":
		default:
			if !allowUnknown {
				return fmt.Errorf("unknown field %q.", runtime.JoinPath(path, k))
			}
		}
	}
	return nil
}

// AtlasValidateJSON function validates a JSON for object Profile.
func (_ *Profile) AtlasValidateJSON(ctx context.Context, r json.RawMessage, path string) (err error) {
	if hook, ok := interface{}(&Profile{}).(interface {
		AtlasJSONValidate(context.Context, json.RawMessage, string) (json.RawMessage, error)
	}); ok {
		if r, err = hook.AtlasJSONValidate(ctx, r, path); err != nil {
			return err
		}
	}
	return validate_Object_Profile(ctx, r, path)
}

func validate_required_Object_Profile(ctx context.Context, v map[string]json.RawMessage, path string) error {
	method := runtime.HTTPMethodFromContext(ctx)
	_ = method
	return nil
}

// validate_Object_UpdateProfileRequest function validates a JSON for a given object.
func validate_Object_UpdateProfileRequest(ctx context.Context, r json.RawMessage, path string) (err error) {
	if hook, ok := interface{}(&UpdateProfileRequest{}).(interface {
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

	if err = validate_required_Object_UpdateProfileRequest(ctx, v, path); err != nil {
		return err
	}

	allowUnknown := runtime.AllowUnknownFromContext(ctx)

	for k, _ := range v {
		switch k {
		case "payload":
			if v[k] == nil {
				continue
			}
			vv := v[k]
			vvPath := runtime.JoinPath(path, k)
			if err = validate_Object_Profile(ctx, vv, vvPath); err != nil {
				return err
			}
		default:
			if !allowUnknown {
				return fmt.Errorf("unknown field %q.", runtime.JoinPath(path, k))
			}
		}
	}
	return nil
}

// AtlasValidateJSON function validates a JSON for object UpdateProfileRequest.
func (_ *UpdateProfileRequest) AtlasValidateJSON(ctx context.Context, r json.RawMessage, path string) (err error) {
	if hook, ok := interface{}(&UpdateProfileRequest{}).(interface {
		AtlasJSONValidate(context.Context, json.RawMessage, string) (json.RawMessage, error)
	}); ok {
		if r, err = hook.AtlasJSONValidate(ctx, r, path); err != nil {
			return err
		}
	}
	return validate_Object_UpdateProfileRequest(ctx, r, path)
}

func validate_required_Object_UpdateProfileRequest(ctx context.Context, v map[string]json.RawMessage, path string) error {
	method := runtime.HTTPMethodFromContext(ctx)
	_ = method
	return nil
}
