package main

import (
	"strings"
	"testing"

	"github.com/sambeau/parsley/pkg/evaluator"
	"github.com/sambeau/parsley/pkg/lexer"
	"github.com/sambeau/parsley/pkg/parser"
)

// Helper function for datetime template tests
func evalDatetimeTemplateTest(t *testing.T, input string) evaluator.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	env := evaluator.NewEnvironment()
	return evaluator.Eval(program, env)
}

// TestDatetimeTemplateBasic tests basic datetime template interpolation
func TestDatetimeTemplateBasic(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Basic date interpolation
		{
			name:     "interpolate month in date",
			input:    `let month = "06"; @(2024-{month}-15).month`,
			expected: "6",
		},
		{
			name:     "interpolate day in date",
			input:    `let day = "25"; @(2024-12-{day}).day`,
			expected: "25",
		},
		{
			name:     "interpolate year in date",
			input:    `let year = "2025"; @({year}-01-01).year`,
			expected: "2025",
		},
		{
			name:     "interpolate multiple parts in date",
			input:    `let y = "2024"; let m = "07"; let d = "04"; @({y}-{m}-{d}).year`,
			expected: "2024",
		},

		// Full datetime interpolation
		{
			name:     "interpolate hour in datetime",
			input:    `let hour = "14"; @(2024-12-25T{hour}:30:00).hour`,
			expected: "14",
		},
		{
			name:     "interpolate minute in datetime",
			input:    `let min = "45"; @(2024-12-25T10:{min}:00).minute`,
			expected: "45",
		},

		// Using numeric values (will be converted to strings)
		{
			name:     "interpolate with numeric month",
			input:    `let month = 6; @(2024-0{month}-15).month`,
			expected: "6",
		},

		// String formatting within interpolation
		{
			name:     "date with padded values",
			input:    `let m = "03"; @(2024-{m}-01).month`,
			expected: "3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalDatetimeTemplateTest(t, tt.input)

			if result == nil {
				t.Fatalf("evaluation returned nil")
			}

			if errObj, ok := result.(*evaluator.Error); ok {
				t.Fatalf("evaluation error: %s", errObj.Message)
			}

			resultStr := result.Inspect()
			if resultStr != tt.expected {
				t.Errorf("expected=%q, got=%q", tt.expected, resultStr)
			}
		})
	}
}

// TestDatetimeTemplateTimeOnly tests time-only template interpolation
func TestDatetimeTemplateTimeOnly(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "interpolate hour in time",
			input:    `let h = "14"; @({h}:30).hour`,
			expected: "14",
		},
		{
			name:     "interpolate minute in time",
			input:    `let m = "45"; @(10:{m}).minute`,
			expected: "45",
		},
		{
			name:     "interpolate hour and minute",
			input:    `let h = "09"; let m = "15"; @({h}:{m}).hour`,
			expected: "9",
		},
		{
			name:     "interpolate time with seconds",
			input:    `let s = "30"; @(12:00:{s}).second`,
			expected: "30",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalDatetimeTemplateTest(t, tt.input)

			if result == nil {
				t.Fatalf("evaluation returned nil")
			}

			if errObj, ok := result.(*evaluator.Error); ok {
				t.Fatalf("evaluation error: %s", errObj.Message)
			}

			resultStr := result.Inspect()
			if resultStr != tt.expected {
				t.Errorf("expected=%q, got=%q", tt.expected, resultStr)
			}
		})
	}
}

// TestDatetimeTemplateKind tests that the correct kind is assigned
func TestDatetimeTemplateKind(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "date-only template has date kind",
			input:    `let m = "06"; @(2024-{m}-15).kind`,
			expected: "date",
		},
		{
			name:     "datetime template has datetime kind",
			input:    `let h = "14"; @(2024-12-25T{h}:30:00).kind`,
			expected: "datetime",
		},
		{
			name:     "time-only template has time kind",
			input:    `let m = "30"; @(12:{m}).kind`,
			expected: "time",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalDatetimeTemplateTest(t, tt.input)

			if result == nil {
				t.Fatalf("evaluation returned nil")
			}

			if errObj, ok := result.(*evaluator.Error); ok {
				t.Fatalf("evaluation error: %s", errObj.Message)
			}

			resultStr := result.Inspect()
			if resultStr != tt.expected {
				t.Errorf("expected=%q, got=%q", tt.expected, resultStr)
			}
		})
	}
}

// TestDatetimeTemplateWithTimezone tests datetime templates with timezone
func TestDatetimeTemplateWithTimezone(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "datetime with UTC timezone",
			input:    `let tz = "Z"; @(2024-12-25T14:30:00{tz}).hour`,
			expected: "14",
		},
		{
			name:     "datetime with positive offset",
			input:    `let h = "14"; @(2024-12-25T{h}:30:00+05:00).hour`,
			expected: "14",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalDatetimeTemplateTest(t, tt.input)

			if result == nil {
				t.Fatalf("evaluation returned nil")
			}

			if errObj, ok := result.(*evaluator.Error); ok {
				t.Fatalf("evaluation error: %s", errObj.Message)
			}

			resultStr := result.Inspect()
			if resultStr != tt.expected {
				t.Errorf("expected=%q, got=%q", tt.expected, resultStr)
			}
		})
	}
}

// TestDatetimeTemplateExpressions tests more complex expressions in interpolations
func TestDatetimeTemplateExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "arithmetic in interpolation",
			input:    `let base = 10; @(2024-12-{base + 5}).day`,
			expected: "15",
		},
		{
			name:     "dictionary access in interpolation",
			input:    `let date = { month: "06", day: "15" }; @(2024-{date.month}-{date.day}).day`,
			expected: "15",
		},
		{
			name:     "array access in interpolation",
			input:    `let parts = ["2024", "12", "25"]; @({parts[0]}-{parts[1]}-{parts[2]}).year`,
			expected: "2024",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalDatetimeTemplateTest(t, tt.input)

			if result == nil {
				t.Fatalf("evaluation returned nil")
			}

			if errObj, ok := result.(*evaluator.Error); ok {
				t.Fatalf("evaluation error: %s", errObj.Message)
			}

			resultStr := result.Inspect()
			if resultStr != tt.expected {
				t.Errorf("expected=%q, got=%q", tt.expected, resultStr)
			}
		})
	}
}

// TestDatetimeTemplateString tests the .iso property for datetime templates
func TestDatetimeTemplateString(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldMatch string // substring that should be present
	}{
		{
			name:        "date template iso output",
			input:       `let m = "06"; @(2024-{m}-15).iso`,
			shouldMatch: "2024-06-15",
		},
		{
			name:        "datetime template iso output",
			input:       `let h = "14"; @(2024-12-25T{h}:30:00Z).iso`,
			shouldMatch: "2024-12-25T14:30:00Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalDatetimeTemplateTest(t, tt.input)

			if result == nil {
				t.Fatalf("evaluation returned nil")
			}

			if errObj, ok := result.(*evaluator.Error); ok {
				t.Fatalf("evaluation error: %s", errObj.Message)
			}

			resultStr := result.Inspect()
			if !strings.Contains(resultStr, tt.shouldMatch) {
				t.Errorf("expected result to contain %q, got=%q", tt.shouldMatch, resultStr)
			}
		})
	}
}

// TestDatetimeTemplateErrors tests error handling in datetime templates
func TestDatetimeTemplateErrors(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedError string
	}{
		{
			name:          "invalid date format",
			input:         `let m = "invalid"; @(2024-{m}-15)`,
			expectedError: "invalid datetime",
		},
		{
			name:          "empty interpolation",
			input:         `@(2024-{}-15)`,
			expectedError: "empty interpolation",
		},
		{
			name:          "undefined variable",
			input:         `@(2024-{undefined_var}-15)`,
			expectedError: "identifier not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalDatetimeTemplateTest(t, tt.input)

			if result == nil {
				t.Fatalf("evaluation returned nil")
			}

			resultStr := result.Inspect()
			if !strings.Contains(resultStr, tt.expectedError) {
				t.Errorf("expected error containing %q, got=%q", tt.expectedError, resultStr)
			}
		})
	}
}

// TestDatetimeTemplateWithNoInterpolation tests that templates work with no interpolations
func TestDatetimeTemplateWithNoInterpolation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "date template with no interpolation",
			input:    `@(2024-12-25).year`,
			expected: "2024",
		},
		{
			name:     "time template with no interpolation",
			input:    `@(14:30).hour`,
			expected: "14",
		},
		{
			name:     "datetime template with no interpolation",
			input:    `@(2024-12-25T14:30:00).hour`,
			expected: "14",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalDatetimeTemplateTest(t, tt.input)

			if result == nil {
				t.Fatalf("evaluation returned nil")
			}

			if errObj, ok := result.(*evaluator.Error); ok {
				t.Fatalf("evaluation error: %s", errObj.Message)
			}

			resultStr := result.Inspect()
			if resultStr != tt.expected {
				t.Errorf("expected=%q, got=%q", tt.expected, resultStr)
			}
		})
	}
}

// TestDatetimeTemplateVsPathTemplate tests that datetime and path templates are distinguished
func TestDatetimeTemplateVsPathTemplate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "datetime template detected correctly",
			input:    `let m = "06"; @(2024-{m}-15).year`,
			expected: "2024",
		},
		{
			name:     "path template still works",
			input:    `let name = "config"; @(./{name}.json).ext`,
			expected: "json",
		},
		{
			name:     "url template still works",
			input:    `let ver = "v1"; @(https://api.com/{ver}/users).host`,
			expected: "api.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalDatetimeTemplateTest(t, tt.input)

			if result == nil {
				t.Fatalf("evaluation returned nil")
			}

			if errObj, ok := result.(*evaluator.Error); ok {
				t.Fatalf("evaluation error: %s", errObj.Message)
			}

			resultStr := result.Inspect()
			if resultStr != tt.expected {
				t.Errorf("expected=%q, got=%q", tt.expected, resultStr)
			}
		})
	}
}
