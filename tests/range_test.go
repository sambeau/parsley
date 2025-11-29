package main

import (
	"testing"

	"github.com/sambeau/parsley/pkg/evaluator"
	"github.com/sambeau/parsley/pkg/lexer"
	"github.com/sambeau/parsley/pkg/parser"
)

func TestRangeOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Basic forward ranges
		{"1..5", "[1, 2, 3, 4, 5]"},
		{"0..3", "[0, 1, 2, 3]"},
		{"10..15", "[10, 11, 12, 13, 14, 15]"},

		// Reverse ranges
		{"5..1", "[5, 4, 3, 2, 1]"},
		{"3..0", "[3, 2, 1, 0]"},
		{"10..5", "[10, 9, 8, 7, 6, 5]"},

		// Negative numbers
		{"-2..2", "[-2, -1, 0, 1, 2]"},
		{"-5..-1", "[-5, -4, -3, -2, -1]"},
		{"2..-2", "[2, 1, 0, -1, -2]"},
		{"-1..-5", "[-1, -2, -3, -4, -5]"},

		// Single element ranges
		{"5..5", "[5]"},
		{"0..0", "[0]"},
		{"-1..-1", "[-1]"},

		// Large ranges
		{"1..100", "[1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100]"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) != 0 {
				t.Fatalf("Parser errors: %v", p.Errors())
			}

			env := evaluator.NewEnvironment()
			result := evaluator.Eval(program, env)

			if result == nil {
				t.Fatalf("Eval returned nil")
			}

			if result.Inspect() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Inspect())
			}
		})
	}
}

func TestRangeInExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// In for loops
		{"for (i in 1..5) { i }", "[1, 2, 3, 4, 5]"},
		{"for (i in 5..1) { i }", "[5, 4, 3, 2, 1]"},
		{"for (i in 0..3) { i * 2 }", "[0, 2, 4, 6]"},

		// With array methods
		{"(1..5).length()", "5"},
		{"(1..10).filter(fn(x) { x % 2 == 0 })", "[2, 4, 6, 8, 10]"},
		{"(1..5).map(fn(x) { x * x })", "[1, 4, 9, 16, 25]"},
		{"(1..3).reverse()", "[3, 2, 1]"},

		// With variables
		{"let a = 1..5; a", "[1, 2, 3, 4, 5]"},
		{"let start = 2; let end = 6; start..end", "[2, 3, 4, 5, 6]"},
		{"let n = 3; (1..n).length()", "3"},

		// Indexing and slicing
		{"(1..10)[0]", "1"},
		{"(1..10)[-1]", "10"},
		{"(1..10)[2:5]", "[3, 4, 5]"},
		{"(5..1)[1]", "4"},

		// Array operations
		{"(1..3) ++ (4..6)", "[1, 2, 3, 4, 5, 6]"},
		{"(1..5) && (3..7)", "[3, 4, 5]"},
		{"(1..5) || (4..8)", "[1, 2, 3, 4, 5, 6, 7, 8]"},
		{"(1..5) - (2..3)", "[1, 4, 5]"},

		// In function arguments
		{"len(1..10)", "10"},
		{"(1..5).join(\", \")", "1, 2, 3, 4, 5"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) != 0 {
				t.Fatalf("Parser errors: %v", p.Errors())
			}

			env := evaluator.NewEnvironment()
			result := evaluator.Eval(program, env)

			if result == nil {
				t.Fatalf("Eval returned nil")
			}

			if result.Type() == evaluator.ERROR_OBJ {
				t.Fatalf("Evaluation error: %s", result.Inspect())
			}

			if result.Inspect() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Inspect())
			}
		})
	}
}

func TestRangeErrors(t *testing.T) {
	tests := []struct {
		input       string
		expectedErr string
	}{
		{"1.5..5", "range start must be an integer"},
		{"1..5.5", "range end must be an integer"},
		{"\"a\"..\"z\"", "range start must be an integer"},
		{"1..[5]", "range end must be an integer"},
		{"true..false", "range start must be an integer"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()

			// Some errors might be caught by parser
			if len(p.Errors()) > 0 {
				return // Parser error is acceptable
			}

			env := evaluator.NewEnvironment()
			result := evaluator.Eval(program, env)

			if result == nil {
				t.Fatalf("Eval returned nil")
			}

			if result.Type() != evaluator.ERROR_OBJ {
				t.Fatalf("expected error, got %s", result.Type())
			}

			errMsg := result.Inspect()
			if len(tt.expectedErr) > 0 && !contains(errMsg, tt.expectedErr) {
				t.Errorf("expected error containing %q, got %q", tt.expectedErr, errMsg)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || findInString(s, substr)))
}

func findInString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
