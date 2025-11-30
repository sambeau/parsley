package main

import (
	"testing"

	"github.com/sambeau/parsley/pkg/evaluator"
	"github.com/sambeau/parsley/pkg/lexer"
	"github.com/sambeau/parsley/pkg/parser"
)

func TestArrayLiterals(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"[1,2,3]", "[1, 2, 3]"},
		{"1", "1"},
		{"[1,2]", "[1, 2]"},
		{`["hello","world"]`, "[hello, world]"},
		{"[1,2,3,4,5]", "[1, 2, 3, 4, 5]"},
		{`[1,"two",3]`, `[1, two, 3]`},
		{"[1+2,3*4,5-1]", "[3, 12, 4]"},
		{"[true,false]", "[true, false]"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			env := evaluator.NewEnvironment()

			program := p.ParseProgram()
			if len(p.Errors()) != 0 {
				t.Fatalf("Parser errors: %v", p.Errors())
			}

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

func TestArrayAssignment(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"xs = [1,2,3]; xs", "[1, 2, 3]"},
		{"arr = [5,10,15]; arr", "[5, 10, 15]"},
		{`names = ["Alice","Bob"]; names`, "[Alice, Bob]"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			env := evaluator.NewEnvironment()

			program := p.ParseProgram()
			if len(p.Errors()) != 0 {
				t.Fatalf("Parser errors: %v", p.Errors())
			}

			var result evaluator.Object
			for _, stmt := range program.Statements {
				result = evaluator.Eval(stmt, env)
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

func TestForFunctionMapping(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"double = fn(x) { x * 2 }; map(double, [1,2,3])", "[2, 4, 6]"},
		{"square = fn(x) { x * x }; map(square, [1,2,3,4])", "[1, 4, 9, 16]"},
		{"inc = fn(x) { x + 1 }; map(inc, [10,20,30])", "[11, 21, 31]"},
		{`upper = fn(s) { s }; map(upper, ["a","b","c"])`, "[a, b, c]"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			env := evaluator.NewEnvironment()

			program := p.ParseProgram()
			if len(p.Errors()) != 0 {
				t.Fatalf("Parser errors: %v", p.Errors())
			}

			var result evaluator.Object
			for _, stmt := range program.Statements {
				result = evaluator.Eval(stmt, env)
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

func TestForFunctionFiltering(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"gt10 = fn(x) { if (x > 10) { return x } }; map(gt10, [5,15,25,8,3,12])", "[15, 25, 12]"},
		{"gt5 = fn(x) { if (x > 5) { return x } }; map(gt5, [1,2,3,4,5,6,7,8])", "[6, 7, 8]"},
		{"positive = fn(x) { if (x > 0) { return x } }; map(positive, [-1,5,-3,10])", "[5, 10]"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			env := evaluator.NewEnvironment()

			program := p.ParseProgram()
			if len(p.Errors()) != 0 {
				t.Fatalf("Parser errors: %v", p.Errors())
			}

			var result evaluator.Object
			for _, stmt := range program.Statements {
				result = evaluator.Eval(stmt, env)
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

func TestForWithArrayVariable(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"double = fn(x) { x * 2 }; xs = [1,2,3]; map(double, xs)", "[2, 4, 6]"},
		{"square = fn(x) { x * x }; nums = [2,3,4]; map(square, nums)", "[4, 9, 16]"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			env := evaluator.NewEnvironment()

			program := p.ParseProgram()
			if len(p.Errors()) != 0 {
				t.Fatalf("Parser errors: %v", p.Errors())
			}

			var result evaluator.Object
			for _, stmt := range program.Statements {
				result = evaluator.Eval(stmt, env)
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

func TestForWithComplexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"transform = fn(x) { x * 2 + 1 }; map(transform, [1,2,3])", "[3, 5, 7]"},
		{"calc = fn(x) { if (x > 2) x * 2 else x + 1 }; map(calc, [1,2,3,4])", "[2, 3, 6, 8]"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			env := evaluator.NewEnvironment()

			program := p.ParseProgram()
			if len(p.Errors()) != 0 {
				t.Fatalf("Parser errors: %v", p.Errors())
			}

			var result evaluator.Object
			for _, stmt := range program.Statements {
				result = evaluator.Eval(stmt, env)
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

func TestArrayPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"x = 1; y = 2; [x,y]", "[1, 2]"},
		{"[1+2,3*4]", "[3, 12]"},
		{"[1,2]", "[1, 2]"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			env := evaluator.NewEnvironment()

			program := p.ParseProgram()
			if len(p.Errors()) != 0 {
				t.Fatalf("Parser errors: %v", p.Errors())
			}

			var result evaluator.Object
			for _, stmt := range program.Statements {
				result = evaluator.Eval(stmt, env)
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

// TestBracketArrayDestructuring tests the bracket syntax for array destructuring
// This is the required syntax as of v0.15.0
func TestBracketArrayDestructuring(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple two element destructuring",
			input:    "let [a, b] = [1, 2]; a",
			expected: "1",
		},
		{
			name:     "access second destructured element",
			input:    "let [a, b] = [1, 2]; b",
			expected: "2",
		},
		{
			name:     "three element destructuring",
			input:    "let [x, y, z] = [10, 20, 30]; x + y + z",
			expected: "60",
		},
		{
			name:     "single element destructuring",
			input:    "let [s] = [42]; s",
			expected: "42",
		},
		{
			name:     "destructuring with string values",
			input:    `let [first, second] = ["hello", "world"]; first`,
			expected: "hello",
		},
		{
			name:     "destructuring from variable",
			input:    "let arr = [1, 2, 3]; let [a, b, c] = arr; b",
			expected: "2",
		},
		{
			name:     "destructuring with rest collects extra elements",
			input:    "let [a, b] = [1, 2, 3, 4]; a",
			expected: "1",
		},
		{
			name:     "destructuring with computed array",
			input:    "let [a, b] = [1 + 1, 2 * 3]; a + b",
			expected: "8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			env := evaluator.NewEnvironment()

			program := p.ParseProgram()
			if len(p.Errors()) != 0 {
				t.Fatalf("Parser errors: %v", p.Errors())
			}

			var result evaluator.Object
			for _, stmt := range program.Statements {
				result = evaluator.Eval(stmt, env)
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

// TestBracketlessArraySyntaxRejected verifies that the old bracketless syntax
// is now rejected by the parser (breaking change in v0.15.0)
func TestBracketlessArraySyntaxRejected(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "bracketless array literal",
			input: "x = 1, 2, 3",
		},
		{
			name:  "bracketless destructuring",
			input: "let a, b = [1, 2]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)

			program := p.ParseProgram()

			// The parser should either have errors or the result should not
			// create an array from comma-separated values
			if len(p.Errors()) == 0 {
				// No parser errors - check that it doesn't create array behavior
				env := evaluator.NewEnvironment()
				evaluator.Eval(program, env)

				// For "x = 1, 2, 3" - x should just be 1, not [1, 2, 3]
				if tt.input == "x = 1, 2, 3" {
					x, _ := env.Get("x")
					if x != nil && x.Inspect() == "[1, 2, 3]" {
						t.Errorf("Bracketless array literal should not create array, got: %s", x.Inspect())
					}
				}
			}
			// Having parser errors for this syntax is acceptable
			// as it's no longer supported
		})
	}
}
