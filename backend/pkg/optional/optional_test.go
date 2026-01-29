package optional

import (
	"encoding/json"
	"strings"
	"testing"
)

type TestStruct struct {
	Name     string           `json:"name"`
	Optional Optional[string] `json:"optional,omitempty"`
}

const resultName = "test"

func TestOptional_UnmarshalJSON_MissingField(t *testing.T) {
	// JSON field missing → Optional.None
	jsonData := `{"name":"test"}`

	var result TestStruct
	if err := json.Unmarshal([]byte(jsonData), &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if !result.Optional.IsNone() {
		t.Errorf("Expected Optional to be None when field is missing, but got IsSome()=%v", result.Optional.IsSome())
	}

	if result.Name != resultName {
		t.Errorf("Expected name to be 'test', got '%s'", result.Name)
	}
}

func TestOptional_UnmarshalJSON_NullField(t *testing.T) {
	// JSON field = null → Optional.None
	jsonData := `{"name":"test","optional":null}`

	var result TestStruct
	if err := json.Unmarshal([]byte(jsonData), &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if !result.Optional.IsNone() {
		t.Errorf("Expected Optional to be None when field is null, but got IsSome()=%v", result.Optional.IsSome())
	}

	if result.Name != resultName {
		t.Errorf("Expected name to be 'test', got '%s'", result.Name)
	}
}

func TestOptional_UnmarshalJSON_WithValue(t *testing.T) {
	// JSON field = value → Optional.Some(value)
	jsonData := `{"name":"test","optional":"some-value"}`

	var result TestStruct
	if err := json.Unmarshal([]byte(jsonData), &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if !result.Optional.IsSome() {
		t.Errorf("Expected Optional to be Some when field has value, but got IsNone()=%v", result.Optional.IsNone())
	}

	value := result.Optional.Unwrap()
	if value != "some-value" {
		t.Errorf("Expected unwrapped value to be 'some-value', got '%s'", value)
	}

	if result.Name != resultName {
		t.Errorf("Expected name to be 'test', got '%s'", result.Name)
	}
}

func TestOptional_MarshalJSON_None(t *testing.T) {
	// Optional.None should marshal to null
	testStruct := TestStruct{
		Name:     resultName,
		Optional: None[string](),
	}

	jsonData, err := json.Marshal(testStruct)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	jsonStr := string(jsonData)
	// Optional.None should marshal to null (omitempty doesn't work with custom types that aren't the zero value)
	if !contains(jsonStr, `"optional":null`) {
		t.Errorf("Expected Optional.None to marshal to null, got: %s", jsonStr)
	}

	// Verify unmarshalling back produces None
	var unmarshalled TestStruct
	if err = json.Unmarshal(jsonData, &unmarshalled); err != nil {
		t.Fatalf("Failed to unmarshal back: %v", err)
	}
	if !unmarshalled.Optional.IsNone() {
		t.Error("Expected unmarshalled Optional to be None")
	}
	if unmarshalled.Name != resultName {
		t.Errorf("Expected name to be 'test', got '%s'", unmarshalled.Name)
	}
}

func TestOptional_MarshalJSON_Some(t *testing.T) {
	// Optional.Some(value) should marshal to the value
	testStruct := TestStruct{
		Name:     resultName,
		Optional: Some("some-value"),
	}

	jsonData, err := json.Marshal(testStruct)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	jsonStr := string(jsonData)
	// Verify the JSON contains the optional field with the correct value
	if !contains(jsonStr, `"optional":"some-value"`) {
		t.Errorf("Expected JSON to contain '\"optional\":\"some-value\"', got: %s", jsonStr)
	}

	// Verify unmarshalling back produces Some
	var unmarshalled TestStruct
	err = json.Unmarshal(jsonData, &unmarshalled)
	if err != nil {
		t.Fatalf("Failed to unmarshal back: %v", err)
	}
	if !unmarshalled.Optional.IsSome() || unmarshalled.Optional.Unwrap() != "some-value" {
		t.Errorf(
			"Expected unmarshalled Optional to be Some('some-value'), got IsSome()=%v, value=%v",
			unmarshalled.Optional.IsSome(),
			unmarshalled.Optional.Unwrap(),
		)
	}
	if unmarshalled.Name != resultName {
		t.Errorf("Expected name to be 'test', got '%s'", unmarshalled.Name)
	}
}

func TestOptional_MarshalUnmarshalRoundTrip(t *testing.T) {
	// Test round-trip: marshal then unmarshal
	tests := []struct {
		name     string
		optional Optional[string]
	}{
		{
			name:     "Some value",
			optional: Some("test-value"),
		},
		{
			name:     "None value",
			optional: None[string](),
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				// Marshal
				testStruct := TestStruct{
					Name:     resultName,
					Optional: tt.optional,
				}

				jsonData, err := json.Marshal(testStruct)
				if err != nil {
					t.Fatalf("Failed to marshal: %v", err)
				}

				// Unmarshal back
				var result TestStruct
				err = json.Unmarshal(jsonData, &result)
				if err != nil {
					t.Fatalf("Failed to unmarshal: %v", err)
				}

				// Verify state matches
				if tt.optional.IsSome() != result.Optional.IsSome() {
					t.Errorf(
						"IsSome() mismatch: original=%v, result=%v",
						tt.optional.IsSome(),
						result.Optional.IsSome(),
					)
				}

				if tt.optional.IsNone() != result.Optional.IsNone() {
					t.Errorf(
						"IsNone() mismatch: original=%v, result=%v",
						tt.optional.IsNone(),
						result.Optional.IsNone(),
					)
				}

				if tt.optional.IsSome() {
					if tt.optional.Unwrap() != result.Optional.Unwrap() {
						t.Errorf(
							"Value mismatch: original=%v, result=%v",
							tt.optional.Unwrap(),
							result.Optional.Unwrap(),
						)
					}
				}
			},
		)
	}
}

func TestOptional_UnmarshalJSON_WithDifferentTypes(t *testing.T) {
	// Test with different types to ensure generic works correctly
	type IntTestStruct struct {
		Value    int           `json:"value"`
		Optional Optional[int] `json:"optional,omitempty"`
	}

	type BoolTestStruct struct {
		Flag     bool           `json:"flag"`
		Optional Optional[bool] `json:"optional,omitempty"`
	}

	tests := []struct {
		name     string
		jsonData string
		check    func(t *testing.T, result interface{})
	}{
		{
			name:     "int missing",
			jsonData: `{"value":42}`,
			check: func(t *testing.T, result interface{}) {
				r := result.(*IntTestStruct)
				if !r.Optional.IsNone() {
					t.Error("Expected Optional[int] to be None when missing")
				}
			},
		},
		{
			name:     "int null",
			jsonData: `{"value":42,"optional":null}`,
			check: func(t *testing.T, result interface{}) {
				r := result.(*IntTestStruct)
				if !r.Optional.IsNone() {
					t.Error("Expected Optional[int] to be None when null")
				}
			},
		},
		{
			name:     "int value",
			jsonData: `{"value":42,"optional":100}`,
			check: func(t *testing.T, result interface{}) {
				r := result.(*IntTestStruct)
				if !r.Optional.IsSome() || r.Optional.Unwrap() != 100 {
					t.Errorf(
						"Expected Optional[int] to be Some(100), got IsSome()=%v, value=%d",
						r.Optional.IsSome(),
						r.Optional.Unwrap(),
					)
				}
			},
		},
		{
			name:     "bool missing",
			jsonData: `{"flag":true}`,
			check: func(t *testing.T, result interface{}) {
				r := result.(*BoolTestStruct)
				if !r.Optional.IsNone() {
					t.Error("Expected Optional[bool] to be None when missing")
				}
			},
		},
		{
			name:     "bool null",
			jsonData: `{"flag":true,"optional":null}`,
			check: func(t *testing.T, result interface{}) {
				r := result.(*BoolTestStruct)
				if !r.Optional.IsNone() {
					t.Error("Expected Optional[bool] to be None when null")
				}
			},
		},
		{
			name:     "bool value",
			jsonData: `{"flag":true,"optional":false}`,
			check: func(t *testing.T, result interface{}) {
				r := result.(*BoolTestStruct)
				if !r.Optional.IsSome() || r.Optional.Unwrap() != false {
					t.Errorf(
						"Expected Optional[bool] to be Some(false), got IsSome()=%v, value=%v",
						r.Optional.IsSome(),
						r.Optional.Unwrap(),
					)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if tt.name[:3] == "int" {
					var result IntTestStruct
					err := json.Unmarshal([]byte(tt.jsonData), &result)
					if err != nil {
						t.Fatalf("Failed to unmarshal: %v", err)
					}
					tt.check(t, &result)
				} else {
					var result BoolTestStruct
					err := json.Unmarshal([]byte(tt.jsonData), &result)
					if err != nil {
						t.Fatalf("Failed to unmarshal: %v", err)
					}
					tt.check(t, &result)
				}
			},
		)
	}
}

// contains is a helper function to check if a string contains a substring.
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
