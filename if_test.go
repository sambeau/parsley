package main

import (
	"strings"
	"testing"

	"pars/pkg/evaluator"
	"pars/pkg/lexer"
	"pars/pkg/parser"
)

func TestIfThenElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Basic if-then-else
		{"x = if 1 > 0 then true else false", "true"},
		{"y = if 1 < 0 then true else false", "false"},
		{"z = if 5 > 3 then 100 else 50", "100"},
		{"w = if 2 < 1 then 100 else 50", "50"},
		
		// If-then without else (should return null when false)
		{"a = if 1 > 0 then 42", "42"},
		{"b = if 1 < 0 then 42", "null"},
		
		// Complex conditions
		{"foo = if 3 * 4 > 10 then 100 else 0", "100"},
		{"bar = if 2 + 2 == 5 then 1 else 2", "2"},
		
		// Using variables in conditions
		{"x = 10; result = if x > 5 then \"big\" else \"small\"", "big"},
		{"x = 3; result = if x > 5 then \"big\" else \"small\"", "small"},
		
		// Nested if expressions
		{"x = if 1 > 0 then if 2 > 1 then 3 else 4 else 5", "3"},
		{"y = if 1 < 0 then 1 else if 2 > 1 then 2 else 3", "2"},
		
		// Using trigonometric functions in conditions
		{"a = if sin(0) == 0 then true else false", "true"},
		{"b = if cos(0) == 1 then 10 else 20", "10"},
		
		// Mathematical expressions in branches
		{"val = if 1 > 0 then pow(2, 3) else sqrt(16)", "8"},
		{"val2 = if 1 < 0 then pow(2, 3) else sqrt(16)", "4"},
		
		// Complex expression example from user request
		{"bar = 15; foo = if bar * 20 > 100 then 100 else bar", "100"},
		{"bar = 3; foo = if bar * 20 > 100 then 100 else bar", "3"},
		
		// Boolean conditions
		{"x = true; result = if x then 1 else 0", "1"},
		{"x = false; result = if x then 1 else 0", "0"},
	}

	for _, tt := range tests {
		env := evaluator.NewEnvironment()
		var result evaluator.Object
		
		// Handle multiple statements separated by semicolon
		statements := strings.Split(tt.input, ";")
		
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			
			l := lexer.New(stmt)
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) != 0 {
				t.Fatalf("parser errors for %q: %v", stmt, p.Errors())
			}

			result = evaluator.Eval(program, env)
			
			if result != nil && result.Type() == evaluator.ERROR_OBJ {
				t.Fatalf("evaluation error for %q: %s", stmt, result.Inspect())
			}
		}

		if result == nil {
			t.Fatalf("Eval returned nil for input: %s", tt.input)
		}

		if result.Inspect() != tt.expected {
			t.Errorf("Expected %s for input %q, got %s", tt.expected, tt.input, result.Inspect())
		}

		t.Logf("✓ Input: %s, Result: %s", tt.input, result.Inspect())
	}
}
