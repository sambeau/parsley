package main

import (
	"testing"

	"github.com/sambeau/parsley/pkg/evaluator"
	"github.com/sambeau/parsley/pkg/lexer"
	"github.com/sambeau/parsley/pkg/parser"
)

func TestOpenEndedArraySlicing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Open-ended from start (arr[n:])
		{"[1,2,3,4,5][2:]", "[3, 4, 5]"},
		{"[1,2,3,4,5][0:]", "[1, 2, 3, 4, 5]"},
		{"[1,2,3,4,5][4:]", "[5]"},
		{"[1,2,3,4,5][5:]", "[]"},
		{"[1,2,3,4,5][-2:]", "[4, 5]"},
		{"[1,2,3,4,5][-1:]", "[5]"},

		// Open-ended from beginning (arr[:n])
		{"[1,2,3,4,5][:3]", "[1, 2, 3]"},
		{"[1,2,3,4,5][:0]", "[]"},
		{"[1,2,3,4,5][:5]", "[1, 2, 3, 4, 5]"},
		{"[1,2,3,4,5][:-2]", "[1, 2, 3]"},
		{"[1,2,3,4,5][:-1]", "[1, 2, 3, 4]"},

		// Both ends open (full copy)
		{"[1,2,3,4,5][:]", "[1, 2, 3, 4, 5]"},
		{"[][:]", "[]"},

		// With variables
		{"let arr = [10, 20, 30, 40]; arr[1:]", "[20, 30, 40]"},
		{"let arr = [10, 20, 30, 40]; arr[:2]", "[10, 20]"},
		{"let arr = [10, 20, 30, 40]; arr[:]", "[10, 20, 30, 40]"},

		// Nested arrays
		{"[[1,2],[3,4],[5,6]][1:]", "[[3, 4], [5, 6]]"},
		{"[[1,2],[3,4],[5,6]][:2]", "[[1, 2], [3, 4]]" },
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
				t.Errorf("Expected %q, got %q", tt.expected, result.Inspect())
			}
		})
	}
}

func TestOpenEndedStringSlicing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Open-ended from start
		{`"hello"[2:]`, "llo"},
		{`"hello"[0:]`, "hello"},
		{`"hello"[5:]`, ""},
		{`"hello"[-2:]`, "lo"},

		// Open-ended from beginning
		{`"hello"[:3]`, "hel"},
		{`"hello"[:0]`, ""},
		{`"hello"[:5]`, "hello"},
		{`"hello"[:-2]`, "hel"},

		// Both ends open
		{`"hello"[:]`, "hello"},
		{`""[:]`, ""},

		// With variables
		{`let s = "Parsley"; s[3:]`, "sley"},
		{`let s = "Parsley"; s[:4]`, "Pars"},
		{`let s = "Parsley"; s[:]`, "Parsley"},
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
				t.Errorf("Expected %q, got %q", tt.expected, result.Inspect())
			}
		})
	}
}

func TestSlicingInExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// In function calls
		{`len([1,2,3,4,5][2:])`, "3"},
		{`len("hello"[:3])`, "3"},

		// In for loops
		{`for (x in [1,2,3,4,5][2:]) { x }`, "[3, 4, 5]"},
		{`for (x in [1,2,3,4,5][:3]) { x }`, "[1, 2, 3]"},

		// Chained operations
		{"[1,2,3,4,5][1:][1:]", "[3, 4, 5]"},
		{`"hello"[1:][:3]`, "ell"},

		// With concatenation
		{"[1,2,3][1:] ++ [4,5]", "[2, 3, 4, 5]"},
		{"[1,2] ++ [3,4,5][:2]", "[1, 2, 3, 4]"},
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
			var result evaluator.Object

			for _, stmt := range program.Statements {
				result = evaluator.Eval(stmt, env)
				if result.Type() == evaluator.ERROR_OBJ {
					t.Fatalf("Evaluation error: %s", result.Inspect())
				}
			}

			if result == nil {
				t.Fatalf("Eval returned nil")
			}

			if result.Inspect() != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result.Inspect())
			}
		})
	}
}

func TestSlicingEdgeCases(t *testing.T) {
	tests := []struct {
		input       string
		shouldError bool
		errorMsg    string
	}{
		// Valid edge cases
		{"[][:]", false, ""},
		{`""[:]`, false, ""},
		{"[1][:100]", false, ""}, // end beyond length should work (clamped)
		{"[1][0:]", false, ""},
		{"[1,2,3][10:20]", false, ""}, // both beyond length - returns empty (clamped)

		// Invalid cases (negative after adjustment)
		{"[1,2,3][-10:]", true, "out of range"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) != 0 {
				if !tt.shouldError {
					t.Fatalf("Unexpected parser errors: %v", p.Errors())
				}
				return
			}

			env := evaluator.NewEnvironment()
			result := evaluator.Eval(program, env)

			if tt.shouldError {
				if result.Type() != evaluator.ERROR_OBJ {
					t.Fatalf("Expected error but got: %s", result.Inspect())
				}
			} else {
				if result.Type() == evaluator.ERROR_OBJ {
					t.Fatalf("Unexpected error: %s", result.Inspect())
				}
			}
		})
	}
}
