package examplepb

import (
	bytes "bytes"
	context "context"
	json "encoding/json"
	runtime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	runtime1 "github.com/infobloxopen/protoc-gen-atlas-validate/runtime"
	metadata "google.golang.org/grpc/metadata"
	ioutil "io/ioutil"
	http "net/http"
)

var validate_Patterns = []struct {
	pattern    runtime.Pattern
	httpMethod string
	validator  func(context.Context, json.RawMessage) error
	// Included for introspection purpose.
	allowUnknown bool
}{
	// patterns for file atlas_validate.proto

	// patterns for file examplepb/example.proto
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
	{
		pattern:      pattern_Groups_ValidateWKT_0,
		httpMethod:   "PUT",
		validator:    validate_Groups_ValidateWKT_0,
		allowUnknown: false,
	},
	{
		pattern:      pattern_Groups_ValidateWKT_1,
		httpMethod:   "PUT",
		validator:    validate_Groups_ValidateWKT_1,
		allowUnknown: false,
	},

	// patterns for file examplepb/example_multi.proto
	{
		pattern:      pattern_Users2_Create2_0,
		httpMethod:   "POST",
		validator:    validate_Users2_Create2_0,
		allowUnknown: false,
	},

	// patterns for file examplepb/examplepb.proto

	// patterns for file external/external.proto

	// patterns for file google/api/annotations.proto

	// patterns for file google/api/http.proto

	// patterns for file google/protobuf/any.proto

	// patterns for file google/protobuf/descriptor.proto

	// patterns for file google/protobuf/empty.proto

	// patterns for file google/protobuf/timestamp.proto

	// patterns for file google/protobuf/wrappers.proto

}

// AtlasValidateAnnotator parses JSON input and validates unknown fields
// based on 'allow_unknown_fields' options specified in proto file.
func AtlasValidateAnnotator(ctx context.Context, r *http.Request) metadata.MD {
	md := make(metadata.MD)
	for _, v := range validate_Patterns {
		if r.Method == v.httpMethod && runtime1.PatternMatch(v.pattern, r.URL.Path) {
			var b []byte
			var err error
			if b, err = ioutil.ReadAll(r.Body); err != nil {
				md.Set("Atlas-Validation-Error", "invalid value: unable to parse body")
				return md
			}
			r.Body = ioutil.NopCloser(bytes.NewReader(b))
			ctx := context.WithValue(context.WithValue(context.Background(), runtime1.HTTPMethodContextKey, r.Method), runtime1.AllowUnknownContextKey, v.allowUnknown)
			if err = v.validator(ctx, b); err != nil {
				md.Set("Atlas-Validation-Error", err.Error())
			}
			break
		}
	}
	return md
}
