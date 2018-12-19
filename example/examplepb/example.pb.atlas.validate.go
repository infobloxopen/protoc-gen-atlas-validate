// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: example/examplepb/example.proto

package examplepb // import "github.com/infobloxopen/protoc-gen-atlas-validate/example/examplepb"

import fmt "fmt"
import http "net/http"
import json "encoding/json"
import ioutil "io/ioutil"
import bytes "bytes"
import context "golang.org/x/net/context"
import metadata "google.golang.org/grpc/metadata"
import runtime "github.com/grpc-ecosystem/grpc-gateway/runtime"
import validate_runtime "github.com/infobloxopen/protoc-gen-atlas-validate/runtime"
import google_protobuf1 "github.com/infobloxopen/protoc-gen-atlas-validate/example/external"
import proto "github.com/gogo/protobuf/proto"
import math "math"
import _ "github.com/infobloxopen/protoc-gen-atlas-validate/example/external"
import _ "github.com/infobloxopen/protoc-gen-atlas-validate/options"
import _ "google.golang.org/genproto/googleapis/api/annotations"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

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
	obj := google_protobuf1.ExternalUser{}
	if validator, ok := interface{}(obj).(interface {
		AtlasValidateJSON(context.Context, json.RawMessage, string) error
	}); ok {
		return validator.AtlasValidateJSON(ctx, r, "")
	}
	return nil
}

// validate_Users_UpdateExternalUser2_0 is an entrypoint for validating "PUT" HTTP request
// that match *.pb.gw.go/pattern_Users_UpdateExternalUser2_0.
func validate_Users_UpdateExternalUser2_0(ctx context.Context, r json.RawMessage) (err error) {
	obj := google_protobuf1.ExternalUser{}
	if validator, ok := interface{}(obj).(interface {
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

// validate_Object_User function validates a JSON for a given object.
func validate_Object_User(ctx context.Context, r json.RawMessage, path string) (err error) {
	obj := &User{}
	if hook, ok := interface{}(obj).(interface {
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
	allowUnknown := validate_runtime.AllowUnknownFromContext(ctx)

	if err = validate_required_Object_User(ctx, v, path); err != nil {
		return err
	}

	for k, _ := range v {
		switch k {
		case "id":
			method := validate_runtime.HTTPMethodFromContext(ctx)
			if method == "POST" {
				return fmt.Errorf("field %q is unsupported for %q operation.", k, method)
			}
		case "name":
		case "profile":
			if v[k] == nil {
				continue
			}
			vv := v[k]
			vvPath := validate_runtime.JoinPath(path, k)
			if err = validate_Object_Profile(ctx, vv, vvPath); err != nil {
				return err
			}
		case "address":
			if v[k] == nil {
				continue
			}
			vv := v[k]
			vvPath := validate_runtime.JoinPath(path, k)
			if err = validate_Object_Address(ctx, vv, vvPath); err != nil {
				return err
			}
		case "groups":
			if v[k] == nil {
				continue
			}
			var vArr []json.RawMessage
			vArrPath := validate_runtime.JoinPath(path, k)
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
			vArrPath := validate_runtime.JoinPath(path, k)
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
			vvPath := validate_runtime.JoinPath(path, k)
			validator, ok := interface{}(&google_protobuf1.ExternalUser{}).(interface {
				AtlasValidateJSON(context.Context, json.RawMessage, string) error
			})
			if !ok {
				continue
			}
			if err = validator.AtlasValidateJSON(ctx, vv, vvPath); err != nil {
				return err
			}
		default:
			if !allowUnknown {
				return fmt.Errorf("unknown field %q.", validate_runtime.JoinPath(path, k))
			}
		}
	}
	return nil
}

// AtlasValidateJSON function validates a JSON for object User.
func (o *User) AtlasValidateJSON(ctx context.Context, r json.RawMessage, path string) (err error) {
	if hook, ok := interface{}(o).(interface {
		AtlasJSONValidate(context.Context, json.RawMessage, string) (json.RawMessage, error)
	}); ok {
		if r, err = hook.AtlasJSONValidate(ctx, r, path); err != nil {
			return err
		}
	}
	return validate_Object_User(ctx, r, path)
}

// ValidateRequiredFileds function return required fields of objectUser.
func (o *User) ValidateRequiredFileds() map[string][]string {
	return map[string][]string{
		"POST":  []string{"Name"},
		"PUT":   []string{"Name"},
		"PATCH": []string{"Name"},
	}
}

func validate_required_Object_User(ctx context.Context, v map[string]json.RawMessage, path string) error {
	method := validate_runtime.HTTPMethodFromContext(ctx)
	_ = method
	if _, ok := v["name"]; !ok {
		fieldPath := validate_runtime.JoinPath(path, "name")
		return fmt.Errorf("field %q is required for %q operation.", fieldPath, method)
	}
	return nil
}

// validate_Object_User_Parent function validates a JSON for a given object.
func validate_Object_User_Parent(ctx context.Context, r json.RawMessage, path string) (err error) {
	obj := &User_Parent{}
	if hook, ok := interface{}(obj).(interface {
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
	allowUnknown := validate_runtime.AllowUnknownFromContext(ctx)

	if err = validate_required_Object_User_Parent(ctx, v, path); err != nil {
		return err
	}

	for k, _ := range v {
		switch k {
		case "name":
		default:
			if !allowUnknown {
				return fmt.Errorf("unknown field %q.", validate_runtime.JoinPath(path, k))
			}
		}
	}
	return nil
}

// AtlasValidateJSON function validates a JSON for object User_Parent.
func (o *User_Parent) AtlasValidateJSON(ctx context.Context, r json.RawMessage, path string) (err error) {
	if hook, ok := interface{}(o).(interface {
		AtlasJSONValidate(context.Context, json.RawMessage, string) (json.RawMessage, error)
	}); ok {
		if r, err = hook.AtlasJSONValidate(ctx, r, path); err != nil {
			return err
		}
	}
	return validate_Object_User_Parent(ctx, r, path)
}

func validate_required_Object_User_Parent(ctx context.Context, v map[string]json.RawMessage, path string) error {
	method := validate_runtime.HTTPMethodFromContext(ctx)
	_ = method
	return nil
}

// validate_Object_Address function validates a JSON for a given object.
func validate_Object_Address(ctx context.Context, r json.RawMessage, path string) (err error) {
	obj := &Address{}
	if hook, ok := interface{}(obj).(interface {
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
	allowUnknown := validate_runtime.AllowUnknownFromContext(ctx)

	if err = validate_required_Object_Address(ctx, v, path); err != nil {
		return err
	}

	for k, _ := range v {
		switch k {
		case "country":
		case "state":
			method := validate_runtime.HTTPMethodFromContext(ctx)
			if method == "PATCH" || method == "PUT" || method == "POST" {
				return fmt.Errorf("field %q is unsupported for %q operation.", k, method)
			}
		case "city":
		case "zip":
		case "tags":
		default:
			if !allowUnknown {
				return fmt.Errorf("unknown field %q.", validate_runtime.JoinPath(path, k))
			}
		}
	}
	return nil
}

// AtlasValidateJSON function validates a JSON for object Address.
func (o *Address) AtlasValidateJSON(ctx context.Context, r json.RawMessage, path string) (err error) {
	if hook, ok := interface{}(o).(interface {
		AtlasJSONValidate(context.Context, json.RawMessage, string) (json.RawMessage, error)
	}); ok {
		if r, err = hook.AtlasJSONValidate(ctx, r, path); err != nil {
			return err
		}
	}
	return validate_Object_Address(ctx, r, path)
}

func validate_required_Object_Address(ctx context.Context, v map[string]json.RawMessage, path string) error {
	method := validate_runtime.HTTPMethodFromContext(ctx)
	_ = method
	return nil
}

// validate_Object_Group function validates a JSON for a given object.
func validate_Object_Group(ctx context.Context, r json.RawMessage, path string) (err error) {
	obj := &Group{}
	if hook, ok := interface{}(obj).(interface {
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
	allowUnknown := validate_runtime.AllowUnknownFromContext(ctx)

	if err = validate_required_Object_Group(ctx, v, path); err != nil {
		return err
	}

	for k, _ := range v {
		switch k {
		case "id":
		case "name":
		case "notes":
		default:
			if !allowUnknown {
				return fmt.Errorf("unknown field %q.", validate_runtime.JoinPath(path, k))
			}
		}
	}
	return nil
}

// AtlasValidateJSON function validates a JSON for object Group.
func (o *Group) AtlasValidateJSON(ctx context.Context, r json.RawMessage, path string) (err error) {
	if hook, ok := interface{}(o).(interface {
		AtlasJSONValidate(context.Context, json.RawMessage, string) (json.RawMessage, error)
	}); ok {
		if r, err = hook.AtlasJSONValidate(ctx, r, path); err != nil {
			return err
		}
	}
	return validate_Object_Group(ctx, r, path)
}

// ValidateRequiredFileds function return required fields of objectGroup.
func (o *Group) ValidateRequiredFileds() map[string][]string {
	return map[string][]string{
		"PATCH": []string{"Id"},
		"PUT":   []string{"Id"},
		"POST":  []string{"Name"},
	}
}

func validate_required_Object_Group(ctx context.Context, v map[string]json.RawMessage, path string) error {
	method := validate_runtime.HTTPMethodFromContext(ctx)
	_ = method
	if _, ok := v["id"]; !ok && (method == "PATCH" || method == "PUT") {
		fieldPath := validate_runtime.JoinPath(path, "id")
		return fmt.Errorf("field %q is required for %q operation.", fieldPath, method)
	}
	if _, ok := v["name"]; !ok && (method == "POST") {
		fieldPath := validate_runtime.JoinPath(path, "name")
		return fmt.Errorf("field %q is required for %q operation.", fieldPath, method)
	}
	return nil
}

// validate_Object_CreateUserRequest function validates a JSON for a given object.
func validate_Object_CreateUserRequest(ctx context.Context, r json.RawMessage, path string) (err error) {
	obj := &CreateUserRequest{}
	if hook, ok := interface{}(obj).(interface {
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
	allowUnknown := validate_runtime.AllowUnknownFromContext(ctx)

	if err = validate_required_Object_CreateUserRequest(ctx, v, path); err != nil {
		return err
	}

	for k, _ := range v {
		switch k {
		case "payload":
			if v[k] == nil {
				continue
			}
			vv := v[k]
			vvPath := validate_runtime.JoinPath(path, k)
			if err = validate_Object_User(ctx, vv, vvPath); err != nil {
				return err
			}
		default:
			if !allowUnknown {
				return fmt.Errorf("unknown field %q.", validate_runtime.JoinPath(path, k))
			}
		}
	}
	return nil
}

// AtlasValidateJSON function validates a JSON for object CreateUserRequest.
func (o *CreateUserRequest) AtlasValidateJSON(ctx context.Context, r json.RawMessage, path string) (err error) {
	if hook, ok := interface{}(o).(interface {
		AtlasJSONValidate(context.Context, json.RawMessage, string) (json.RawMessage, error)
	}); ok {
		if r, err = hook.AtlasJSONValidate(ctx, r, path); err != nil {
			return err
		}
	}
	return validate_Object_CreateUserRequest(ctx, r, path)
}

func validate_required_Object_CreateUserRequest(ctx context.Context, v map[string]json.RawMessage, path string) error {
	method := validate_runtime.HTTPMethodFromContext(ctx)
	_ = method
	return nil
}

// validate_Object_UpdateUserRequest function validates a JSON for a given object.
func validate_Object_UpdateUserRequest(ctx context.Context, r json.RawMessage, path string) (err error) {
	obj := &UpdateUserRequest{}
	if hook, ok := interface{}(obj).(interface {
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
	allowUnknown := validate_runtime.AllowUnknownFromContext(ctx)

	if err = validate_required_Object_UpdateUserRequest(ctx, v, path); err != nil {
		return err
	}

	for k, _ := range v {
		switch k {
		case "payload":
			if v[k] == nil {
				continue
			}
			vv := v[k]
			vvPath := validate_runtime.JoinPath(path, k)
			if err = validate_Object_User(ctx, vv, vvPath); err != nil {
				return err
			}
		default:
			if !allowUnknown {
				return fmt.Errorf("unknown field %q.", validate_runtime.JoinPath(path, k))
			}
		}
	}
	return nil
}

// AtlasValidateJSON function validates a JSON for object UpdateUserRequest.
func (o *UpdateUserRequest) AtlasValidateJSON(ctx context.Context, r json.RawMessage, path string) (err error) {
	if hook, ok := interface{}(o).(interface {
		AtlasJSONValidate(context.Context, json.RawMessage, string) (json.RawMessage, error)
	}); ok {
		if r, err = hook.AtlasJSONValidate(ctx, r, path); err != nil {
			return err
		}
	}
	return validate_Object_UpdateUserRequest(ctx, r, path)
}

func validate_required_Object_UpdateUserRequest(ctx context.Context, v map[string]json.RawMessage, path string) error {
	method := validate_runtime.HTTPMethodFromContext(ctx)
	_ = method
	return nil
}

// validate_Object_EmptyRequest function validates a JSON for a given object.
func validate_Object_EmptyRequest(ctx context.Context, r json.RawMessage, path string) (err error) {
	obj := &EmptyRequest{}
	if hook, ok := interface{}(obj).(interface {
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
	allowUnknown := validate_runtime.AllowUnknownFromContext(ctx)

	if err = validate_required_Object_EmptyRequest(ctx, v, path); err != nil {
		return err
	}

	for k, _ := range v {
		switch k {
		default:
			if !allowUnknown {
				return fmt.Errorf("unknown field %q.", validate_runtime.JoinPath(path, k))
			}
		}
	}
	return nil
}

// AtlasValidateJSON function validates a JSON for object EmptyRequest.
func (o *EmptyRequest) AtlasValidateJSON(ctx context.Context, r json.RawMessage, path string) (err error) {
	if hook, ok := interface{}(o).(interface {
		AtlasJSONValidate(context.Context, json.RawMessage, string) (json.RawMessage, error)
	}); ok {
		if r, err = hook.AtlasJSONValidate(ctx, r, path); err != nil {
			return err
		}
	}
	return validate_Object_EmptyRequest(ctx, r, path)
}

func validate_required_Object_EmptyRequest(ctx context.Context, v map[string]json.RawMessage, path string) error {
	method := validate_runtime.HTTPMethodFromContext(ctx)
	_ = method
	return nil
}

// validate_Object_EmptyResponse function validates a JSON for a given object.
func validate_Object_EmptyResponse(ctx context.Context, r json.RawMessage, path string) (err error) {
	obj := &EmptyResponse{}
	if hook, ok := interface{}(obj).(interface {
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
	allowUnknown := validate_runtime.AllowUnknownFromContext(ctx)

	if err = validate_required_Object_EmptyResponse(ctx, v, path); err != nil {
		return err
	}

	for k, _ := range v {
		switch k {
		default:
			if !allowUnknown {
				return fmt.Errorf("unknown field %q.", validate_runtime.JoinPath(path, k))
			}
		}
	}
	return nil
}

// AtlasValidateJSON function validates a JSON for object EmptyResponse.
func (o *EmptyResponse) AtlasValidateJSON(ctx context.Context, r json.RawMessage, path string) (err error) {
	if hook, ok := interface{}(o).(interface {
		AtlasJSONValidate(context.Context, json.RawMessage, string) (json.RawMessage, error)
	}); ok {
		if r, err = hook.AtlasJSONValidate(ctx, r, path); err != nil {
			return err
		}
	}
	return validate_Object_EmptyResponse(ctx, r, path)
}

func validate_required_Object_EmptyResponse(ctx context.Context, v map[string]json.RawMessage, path string) error {
	method := validate_runtime.HTTPMethodFromContext(ctx)
	_ = method
	return nil
}

// validate_Object_Profile function validates a JSON for a given object.
func validate_Object_Profile(ctx context.Context, r json.RawMessage, path string) (err error) {
	obj := &Profile{}
	if hook, ok := interface{}(obj).(interface {
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
	allowUnknown := validate_runtime.AllowUnknownFromContext(ctx)

	if err = validate_required_Object_Profile(ctx, v, path); err != nil {
		return err
	}

	for k, _ := range v {
		switch k {
		case "id":
		case "name":
			method := validate_runtime.HTTPMethodFromContext(ctx)
			if method == "PATCH" || method == "PUT" {
				return fmt.Errorf("field %q is unsupported for %q operation.", k, method)
			}
		case "notes":
		default:
			if !allowUnknown {
				return fmt.Errorf("unknown field %q.", validate_runtime.JoinPath(path, k))
			}
		}
	}
	return nil
}

// AtlasValidateJSON function validates a JSON for object Profile.
func (o *Profile) AtlasValidateJSON(ctx context.Context, r json.RawMessage, path string) (err error) {
	if hook, ok := interface{}(o).(interface {
		AtlasJSONValidate(context.Context, json.RawMessage, string) (json.RawMessage, error)
	}); ok {
		if r, err = hook.AtlasJSONValidate(ctx, r, path); err != nil {
			return err
		}
	}
	return validate_Object_Profile(ctx, r, path)
}

func validate_required_Object_Profile(ctx context.Context, v map[string]json.RawMessage, path string) error {
	method := validate_runtime.HTTPMethodFromContext(ctx)
	_ = method
	return nil
}

// validate_Object_UpdateProfileRequest function validates a JSON for a given object.
func validate_Object_UpdateProfileRequest(ctx context.Context, r json.RawMessage, path string) (err error) {
	obj := &UpdateProfileRequest{}
	if hook, ok := interface{}(obj).(interface {
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
	allowUnknown := validate_runtime.AllowUnknownFromContext(ctx)

	if err = validate_required_Object_UpdateProfileRequest(ctx, v, path); err != nil {
		return err
	}

	for k, _ := range v {
		switch k {
		case "payload":
			if v[k] == nil {
				continue
			}
			vv := v[k]
			vvPath := validate_runtime.JoinPath(path, k)
			if err = validate_Object_Profile(ctx, vv, vvPath); err != nil {
				return err
			}
		default:
			if !allowUnknown {
				return fmt.Errorf("unknown field %q.", validate_runtime.JoinPath(path, k))
			}
		}
	}
	return nil
}

// AtlasValidateJSON function validates a JSON for object UpdateProfileRequest.
func (o *UpdateProfileRequest) AtlasValidateJSON(ctx context.Context, r json.RawMessage, path string) (err error) {
	if hook, ok := interface{}(o).(interface {
		AtlasJSONValidate(context.Context, json.RawMessage, string) (json.RawMessage, error)
	}); ok {
		if r, err = hook.AtlasJSONValidate(ctx, r, path); err != nil {
			return err
		}
	}
	return validate_Object_UpdateProfileRequest(ctx, r, path)
}

func validate_required_Object_UpdateProfileRequest(ctx context.Context, v map[string]json.RawMessage, path string) error {
	method := validate_runtime.HTTPMethodFromContext(ctx)
	_ = method
	return nil
}

var validate_Patterns = []struct {
	pattern    runtime.Pattern
	httpMethod string
	validator  func(context.Context, json.RawMessage) error
	// Included for introspection purpose.
	allowUnknown bool
}{
	{
		pattern:      pattern_Users_Create_0,
		httpMethod:   "POST",
		validator:    validate_Users_Create_0,
		allowUnknown: false,
	},
	{
		pattern:      pattern_Users_Update_0,
		httpMethod:   "PUT",
		validator:    validate_Users_Update_0,
		allowUnknown: false,
	},
	{
		pattern:      pattern_Users_Update_1,
		httpMethod:   "PATCH",
		validator:    validate_Users_Update_1,
		allowUnknown: false,
	},
	{
		pattern:      pattern_Users_List_0,
		httpMethod:   "GET",
		validator:    validate_Users_List_0,
		allowUnknown: false,
	},
	{
		pattern:      pattern_Users_List_1,
		httpMethod:   "GET",
		validator:    validate_Users_List_1,
		allowUnknown: false,
	},
	{
		pattern:      pattern_Users_UpdateExternalUser_0,
		httpMethod:   "PUT",
		validator:    validate_Users_UpdateExternalUser_0,
		allowUnknown: false,
	},
	{
		pattern:      pattern_Users_UpdateExternalUser2_0,
		httpMethod:   "PUT",
		validator:    validate_Users_UpdateExternalUser2_0,
		allowUnknown: false,
	},
	{
		pattern:      pattern_Profiles_Create_0,
		httpMethod:   "POST",
		validator:    validate_Profiles_Create_0,
		allowUnknown: false,
	},
	{
		pattern:      pattern_Profiles_Update_0,
		httpMethod:   "PUT",
		validator:    validate_Profiles_Update_0,
		allowUnknown: true,
	},
	{
		pattern:      pattern_Groups_Create_0,
		httpMethod:   "POST",
		validator:    validate_Groups_Create_0,
		allowUnknown: true,
	},
	{
		pattern:      pattern_Groups_Update_0,
		httpMethod:   "PUT",
		validator:    validate_Groups_Update_0,
		allowUnknown: true,
	},
	{
		pattern:      pattern_Groups_ValidatedList_0,
		httpMethod:   "GET",
		validator:    validate_Groups_ValidatedList_0,
		allowUnknown: false,
	},
	{
		pattern:      pattern_Groups_ValidatedList_1,
		httpMethod:   "GET",
		validator:    validate_Groups_ValidatedList_1,
		allowUnknown: false,
	},
}

// AtlasValidateAnnotator parses JSON input and validates unknown fields
// based on 'allow_unknown_fields' options specified in proto file.
func AtlasValidateAnnotator(ctx context.Context, r *http.Request) metadata.MD {
	md := make(metadata.MD)
	for _, v := range validate_Patterns {
		if r.Method == v.httpMethod && validate_runtime.PatternMatch(v.pattern, r.URL.Path) {
			var b []byte
			var err error
			if b, err = ioutil.ReadAll(r.Body); err != nil {
				md.Set("Atlas-Validation-Error", "invalid value: unable to parse body")
				return md
			}
			r.Body = ioutil.NopCloser(bytes.NewReader(b))
			ctx := context.WithValue(context.WithValue(context.Background(), validate_runtime.HTTPMethodContextKey, r.Method), validate_runtime.AllowUnknownContextKey, v.allowUnknown)
			if err = v.validator(ctx, b); err != nil {
				md.Set("Atlas-Validation-Error", err.Error())
			}
			break
		}
	}
	return md
}
