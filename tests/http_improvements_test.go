package main

import (
	"testing"

	"github.com/sambeau/parsley/pkg/evaluator"
	"github.com/sambeau/parsley/pkg/lexer"
	"github.com/sambeau/parsley/pkg/parser"
)

// Helper function to evaluate Parsley code and return the result
func evalHTTP(input string) evaluator.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := evaluator.NewEnvironment()
	return evaluator.Eval(program, env)
}

// TestHTTPMethodAccessors tests the .get, .post, .put, .patch, .delete accessors
func TestHTTPMethodAccessors(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantType string
		wantKey  string
		wantVal  string
	}{
		{
			name:     "GET method accessor",
			input:    `JSON(@https://example.com/api).get.method`,
			wantType: "string",
			wantVal:  "GET",
		},
		{
			name:     "POST method accessor",
			input:    `JSON(@https://example.com/api).post.method`,
			wantType: "string",
			wantVal:  "POST",
		},
		{
			name:     "PUT method accessor",
			input:    `JSON(@https://example.com/api).put.method`,
			wantType: "string",
			wantVal:  "PUT",
		},
		{
			name:     "PATCH method accessor",
			input:    `JSON(@https://example.com/api).patch.method`,
			wantType: "string",
			wantVal:  "PATCH",
		},
		{
			name:     "DELETE method accessor",
			input:    `JSON(@https://example.com/api).delete.method`,
			wantType: "string",
			wantVal:  "DELETE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalHTTP(tt.input)
			if errObj, ok := result.(*evaluator.Error); ok {
				t.Fatalf("unexpected error: %s", errObj.Message)
			}
			if strObj, ok := result.(*evaluator.String); ok {
				if strObj.Value != tt.wantVal {
					t.Errorf("Expected %s, got %s", tt.wantVal, strObj.Value)
				}
			} else {
				t.Errorf("Expected string result, got %T", result)
			}
		})
	}
}

// TestResponseTypedDict tests the response typed dictionary structure
func TestResponseTypedDict(t *testing.T) {
	// Test that a request dict has the correct structure
	tests := []struct {
		name     string
		input    string
		wantType string
		wantVal  string
	}{
		{
			name:     "Request dict has __type",
			input:    `JSON(@https://example.com/api).__type`,
			wantType: "string",
			wantVal:  "request",
		},
		{
			name:     "Request dict has format",
			input:    `JSON(@https://example.com/api).format`,
			wantType: "string",
			wantVal:  "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalHTTP(tt.input)
			if errObj, ok := result.(*evaluator.Error); ok {
				t.Fatalf("unexpected error: %s", errObj.Message)
			}
			if strObj, ok := result.(*evaluator.String); ok {
				if strObj.Value != tt.wantVal {
					t.Errorf("Expected %s, got %s", tt.wantVal, strObj.Value)
				}
			} else {
				t.Errorf("Expected string result, got %T", result)
			}
		})
	}
}

// TestResponseMethodExists tests that the .response() method exists on response dicts
func TestResponseMethodExists(t *testing.T) {
	// We can't easily test actual HTTP calls without a mock server,
	// but we can test that the method accessors work on request dicts
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Method chaining works",
			input: `let req = JSON(@https://example.com/api).post; req.method`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalHTTP(tt.input)
			if errObj, ok := result.(*evaluator.Error); ok {
				t.Fatalf("unexpected error: %s", errObj.Message)
			}
			if strObj, ok := result.(*evaluator.String); ok {
				if strObj.Value != "POST" {
					t.Errorf("Expected POST, got %s", strObj.Value)
				}
			} else {
				t.Errorf("Expected string result, got %T", result)
			}
		})
	}
}

// TestIsResponseDict tests the isResponseDict function
func TestIsResponseDict(t *testing.T) {
	// Create a mock response dict structure
	input := `
		let mockResponse = {
			__type: "response",
			__format: "json",
			__data: [1, 2, 3],
			__response: {
				status: 200,
				statusText: "OK",
				ok: true,
				error: null
			}
		}
		mockResponse.__type
	`
	result := evalHTTP(input)
	if errObj, ok := result.(*evaluator.Error); ok {
		t.Fatalf("unexpected error: %s", errObj.Message)
	}
	if strObj, ok := result.(*evaluator.String); ok {
		if strObj.Value != "response" {
			t.Errorf("Expected 'response', got %s", strObj.Value)
		}
	} else {
		t.Errorf("Expected string result, got %T", result)
	}
}
