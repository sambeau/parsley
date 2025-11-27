package main

import (
	"testing"

	"github.com/sambeau/parsley/pkg/evaluator"
	"github.com/sambeau/parsley/pkg/lexer"
	"github.com/sambeau/parsley/pkg/parser"
)

func TestLogicalOperators(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// AND operator with &
		{"true & true", "true"},
		{"true & false", "false"},
		{"false & true", "false"},
		{"false & false", "false"},

		// AND operator with 'and' keyword
		{"true and true", "true"},
		{"true and false", "false"},
		{"false and true", "false"},
		{"false and false", "false"},

		// OR operator with |
		{"true | true", "true"},
		{"true | false", "true"},
		{"false | true", "true"},
		{"false | false", "false"},

		// OR operator with 'or' keyword
		{"true or true", "true"},
		{"true or false", "true"},
		{"false or true", "true"},
		{"false or false", "false"},

		// NOT operator with !
		{"!true", "false"},
		{"!false", "true"},
		{"!!true", "true"},
		{"!!false", "false"},

		// NOT operator with 'not' keyword
		{"not true", "false"},
		{"not false", "true"},
		{"not not true", "true"},
		{"not not false", "false"},

		// Combined logical operators
		{"true & true & true", "true"},
		{"true & true & false", "false"},
		{"false | false | true", "true"},
		{"false | false | false", "false"},

		// Logical operators with comparisons
		{"5 > 3 & 2 < 4", "true"},
		{"5 > 3 and 2 < 4", "true"},
		{"5 > 3 & 2 > 4", "false"},
		{"5 < 3 | 2 < 4", "true"},
		{"5 > 3 or 2 > 4", "true"},
		{"5 < 3 or 2 > 4", "false"},

		// Precedence: AND has higher precedence than OR
		{"true | false & false", "true"}, // true | (false & false) = true | false = true
		{"false & true | true", "true"},  // (false & true) | true = false | true = true

		// NOT with other operators
		{"!true & true", "false"},
		{"!(true & false)", "true"},
		{"not (5 > 3)", "false"},
		{"not (5 < 3)", "true"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			t.Fatalf("parser errors for %q: %v", tt.input, p.Errors())
		}

		env := evaluator.NewEnvironment()
		result := evaluator.Eval(program, env)

		if result == nil {
			t.Fatalf("Eval returned nil for input: %s", tt.input)
		}

		if result.Type() == evaluator.ERROR_OBJ {
			t.Fatalf("evaluation error for %q: %s", tt.input, result.Inspect())
		}

		if result.Inspect() != tt.expected {
			t.Errorf("Expected %s for input %q, got %s", tt.expected, tt.input, result.Inspect())
		}

		t.Logf("✓ Input: %s, Result: %s", tt.input, result.Inspect())
	}
}

func TestComparisonOperators(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Not equal operator
		{"5 != 3", "true"},
		{"5 != 5", "false"},
		{"3 != 5", "true"},

		// Less than or equal
		{"3 <= 5", "true"},
		{"5 <= 5", "true"},
		{"5 <= 3", "false"},

		// Greater than or equal (already tested, but let's be thorough)
		{"5 >= 3", "true"},
		{"5 >= 5", "true"},
		{"3 >= 5", "false"},

		// Combined with other operators
		{"5 != 3 & 2 <= 4", "true"},
		{"5 != 5 | 2 <= 4", "true"},
		{"3 <= 5 and 5 >= 3", "true"},

		// String comparisons
		{"\"hello\" != \"world\"", "true"},
		{"\"hello\" != \"hello\"", "false"},

		// Float comparisons
		{"3.5 >= 3.0", "true"},
		{"3.0 <= 3.5", "true"},
		{"3.5 != 3.0", "true"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			t.Fatalf("parser errors for %q: %v", tt.input, p.Errors())
		}

		env := evaluator.NewEnvironment()
		result := evaluator.Eval(program, env)

		if result == nil {
			t.Fatalf("Eval returned nil for input: %s", tt.input)
		}

		if result.Type() == evaluator.ERROR_OBJ {
			t.Fatalf("evaluation error for %q: %s", tt.input, result.Inspect())
		}

		if result.Inspect() != tt.expected {
			t.Errorf("Expected %s for input %q, got %s", tt.expected, tt.input, result.Inspect())
		}

		t.Logf("✓ Input: %s, Result: %s", tt.input, result.Inspect())
	}
}
