package examplepb

import (
	"encoding/json"
	"testing"
)

func TestCreateProfile(t *testing.T) {
	tests := []struct {
		input            json.RawMessage
		validateFunction func(message json.RawMessage) error
		negative         bool
	}{
		{
			input:            json.RawMessage([]byte(`{"id": 1, "name": "first", "notes": "some notes"}`)),
			validateFunction: validate_Profiles_Create_0,
			negative:         true,
		},
		{
			input:            json.RawMessage([]byte(`{"id": 1, "name": "first", "notes": "some notes"}`)),
			validateFunction: validate_Profiles_Update_0,
			negative:         true,
		},
		{
			input:            json.RawMessage([]byte(`{"id": 1, "name": "first", "profile": {"name":"some name"}}`)),
			validateFunction: validate_Users_Create_0,
			negative:         true,
		},
		{
			input:            json.RawMessage([]byte(`{"id": 1, "name": "first", "profile": {"id": 1}}`)),
			validateFunction: validate_Users_Create_0,
			negative:         false,
		},
		{
			input:            json.RawMessage([]byte(`{"id": 1, "name": "first", "profile": {"id": 1}, "address": {"country": "USA","state":"NY", "city": "New York", "zip": "21412412"}}`)),
			validateFunction: validate_Users_Create_0,
			negative:         false,
		},
		{
			input:            json.RawMessage([]byte(`{"id": 1, "name": "first", "profile": {"id": 1}, "address": {"country": "USA","state":"NY", "city": "New York", "zip": "21412412"}}`)),
			validateFunction: validate_Users_Update_0,
			negative:         true,
		},
	}

	for n, test := range tests {
		err := test.validateFunction(test.input)
		if err != nil && !test.negative {
			t.Errorf(" %d test failed, error %s \n", n, err.Error())
		}

	}
}
