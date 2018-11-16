package examplepb

import (
	"context"
	"encoding/json"
	"github.com/infobloxopen/protoc-gen-atlas-validate/runtime"
	"testing"
)

type Test struct {
	input            json.RawMessage
	validateFunction func(ctx context.Context, message json.RawMessage) error
	context          context.Context
	negative         bool
}

func TestDenyFields(t *testing.T) {
	tests := []Test{
		{
			input:            json.RawMessage([]byte(`{"id": 1, "name": "first", "notes": "some notes"}`)),
			validateFunction: validate_Profiles_Create_0,
			context:          context.WithValue(context.Background(), runtime.HTTPMethodContextKey, "POST"),
			negative:         false,
		},
		{
			input:            json.RawMessage([]byte(`{"id": 1, "name": "first", "notes": "some notes"}`)),
			validateFunction: validate_Profiles_Update_0,
			context:          context.WithValue(context.Background(), runtime.HTTPMethodContextKey, "PATCH"),
			negative:         true,
		},
		{
			input:            json.RawMessage([]byte(`{"id": 1, "notes": "some notes"}`)),
			validateFunction: validate_Profiles_Update_0,
			context:          context.WithValue(context.Background(), runtime.HTTPMethodContextKey, "PATCH"),
			negative:         false,
		},
		{
			input:            json.RawMessage([]byte(`{"name": "first"}`)),
			validateFunction: validate_Users_Create_0,
			context:          context.WithValue(context.Background(), runtime.HTTPMethodContextKey, "POST"),
			negative:         false,
		},
		{
			input:            json.RawMessage([]byte(`{"id":1, "name": "first"}`)),
			validateFunction: validate_Users_Create_0,
			context:          context.WithValue(context.Background(), runtime.HTTPMethodContextKey, "POST"),
			negative:         true,
		},
		{
			input:            json.RawMessage([]byte(`{"id":1, "name": "first"}`)),
			validateFunction: validate_Users_Update_0,
			context:          context.WithValue(context.Background(), runtime.HTTPMethodContextKey, "PUT"),
			negative:         false,
		},
		{
			input:            json.RawMessage([]byte(`{"name": "first", "profile": {"id": 1, "name":"some name"}}`)),
			validateFunction: validate_Users_Create_0,
			context:          context.WithValue(context.Background(), runtime.HTTPMethodContextKey, "POST"),
			negative:         false,
		},
		{
			input:            json.RawMessage([]byte(`{"name": "first", "profile": {"id": 1, "name":"some name"}}`)),
			validateFunction: validate_Users_Update_0,
			context:          context.WithValue(context.Background(), runtime.HTTPMethodContextKey, "PUT"),
			negative:         true,
		},
		{
			input:            json.RawMessage([]byte(`{"id":1, "name": "first", "profile": {"id": 1}}`)),
			validateFunction: validate_Users_Update_0,
			context:          context.WithValue(context.Background(), runtime.HTTPMethodContextKey, "PUT"),
			negative:         false,
		},
		{
			input:            json.RawMessage([]byte(`{"name": "first", "profile": {"id": 1}, "address": {"country": "USA","state":"NY", "city": "New York", "zip": "21412412"}}`)),
			validateFunction: validate_Users_Create_0,
			context:          context.WithValue(context.Background(), runtime.HTTPMethodContextKey, "POST"),
			negative:         true,
		},
		{
			input:            json.RawMessage([]byte(`{"name": "first", "profile": {"id": 1}, "address": {"country": "USA","city": "New York", "zip": "21412412"}}`)),
			validateFunction: validate_Users_Create_0,
			context:          context.WithValue(context.Background(), runtime.HTTPMethodContextKey, "POST"),
			negative:         false,
		},
		{
			input:            json.RawMessage([]byte(`{"name": "first", "profile": {"id": 1}, "address": {"country": "USA","state":"NY", "city": "New York", "zip": "21412412"}}`)),
			validateFunction: validate_Users_Update_0,
			context:          context.WithValue(context.Background(), runtime.HTTPMethodContextKey, "PUT"),
			negative:         true,
		},
		{
			input:            json.RawMessage([]byte(`{"name": "first", "profile": {"id": 1}, "address": {"country": "USA","city": "New York", "zip": "21412412"}}`)),
			validateFunction: validate_Users_Update_0,
			context:          context.WithValue(context.Background(), runtime.HTTPMethodContextKey, "PUT"),
			negative:         false,
		},
	}

	for n, test := range tests {
		err := test.validateFunction(test.context, test.input)
		if err == nil && test.negative {
			t.Errorf(" %d test failed, error must be not nil \n", n+1)
		}

		if err != nil && !test.negative {
			t.Errorf(" %d test failed, error %s \n", n+1, err.Error())
		}

	}
}

func TestAllowUnknown(t *testing.T) {
	tests := []Test{
		{
			input:            json.RawMessage([]byte(`{"id": 1, "notes": "some notes", "unknown_field": "unknown_value"}`)),
			validateFunction: validate_Profiles_Update_0,
			context:          context.WithValue(context.WithValue(context.Background(), runtime.HTTPMethodContextKey, "PUT"), runtime.AllowUnknownContextKey, true),
			negative:         false,
		},
		{
			input:            json.RawMessage([]byte(`{"id": 1, "name": "first", "notes": "some notes", "unknown_field": "unknown_value"}`)),
			validateFunction: validate_Profiles_Create_0,
			context:          context.WithValue(context.Background(), runtime.HTTPMethodContextKey, "POST"),
			negative:         true,
		},
		{
			input:            json.RawMessage([]byte(`{"id": 1, "name": "first", "notes": "some notes", "unknown_field": "unknown_value"}`)),
			validateFunction: validate_Groups_Create_0,
			context:          context.WithValue(context.WithValue(context.Background(), runtime.HTTPMethodContextKey, "PUT"), runtime.AllowUnknownContextKey, true),
			negative:         false,
		},
		{
			input:            json.RawMessage([]byte(`{"id": 1, "name": "first", "notes": "some notes", "unknown_field": "unknown_value"}`)),
			validateFunction: validate_Groups_Update_0,
			context:          context.WithValue(context.WithValue(context.Background(), runtime.HTTPMethodContextKey, "PUT"), runtime.AllowUnknownContextKey, true),
			negative:         false,
		},
		{
			input:            json.RawMessage([]byte(`{"name": "first", "profile": {"name":"some name"}, "unknown_field": "unknown_value"}`)),
			validateFunction: validate_Users_Create_0,
			context:          context.WithValue(context.Background(), runtime.HTTPMethodContextKey, "POST"),
			negative:         true,
		},
	}

	for n, test := range tests {
		err := test.validateFunction(test.context, test.input)
		if err != nil && !test.negative {
			t.Errorf(" %d test failed, error %s \n", n+1, err.Error())
		}
	}
}
