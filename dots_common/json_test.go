package dots_common_test

import (
	"testing"
	"github.com/nttdots/go-dots/dots_common"
)

func Test_JSON(t *testing.T) {
	var expects interface{}
	validData := "{\"message\": \"dots_client\"}"
	validSchemaPath := "test"

	result := dots_common.ValidateJson(validSchemaPath, validData)

	if result != nil {
		t.Errorf("ValidateJson got %s, want %s", result, expects)
	}

	invalidData := "{\"message\": 1}"
	validSchemaPath = "test"

	result = dots_common.ValidateJson(validSchemaPath, invalidData)

	if result == nil {
		t.Errorf("ValidateJson got %s, want %s", result, expects)
	}

	validData = "{\"message\": \"dots_client\"}"
	invalidSchemaPath := "invalid"

	result = dots_common.ValidateJson(invalidSchemaPath, validData)

	if result == nil {
		t.Errorf("ValidateJson got %s, want %s", result, expects)
	}
}
