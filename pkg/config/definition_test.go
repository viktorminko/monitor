package config

import (
	"testing"
)

func TestDefinitionPrepare(t *testing.T) {
	methodData := &Definition{
		"test_prepare",
		"/api/{{.a1}}/s1/{{.b2}}/s2/{{.c3}}",
		false,
		"GET",
		nil,
		"",
		Duration{1},
		Duration{5},
		200,
	}

	err := methodData.Prepare(&Context{
		"a1": "1a",
		"b2": "2b",
		"c3": "3c",
		"d4": "4d",
	})

	if err != nil {
		t.Fatalf("Unexpected error returned: %s", err.Error())
	}

	expected := "/api/1a/s1/2b/s2/3c"
	if methodData.URL != expected {
		t.Fatalf(
			"unexpected URL prepared. Expected %s, got %s",
			expected,
			methodData.URL,
		)
	}
}
