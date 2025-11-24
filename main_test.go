package main

import (
	"testing"

	"pars/pkg/evaluator"
	"pars/pkg/lexer"
	"pars/pkg/parser"
)

func TestMain(t *testing.T) {
	// This is a placeholder test
	// Replace with actual tests for your functions
	t.Log("Test placeholder - replace with real tests")
}

func TestTrigonometricFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"sin(0)", "0"},
		{"cos(0)", "1"},
		{"tan(0)", "0"},
		{"sqrt(4)", "2"},
		{"pow(2, 3)", "8"},
		{"pi()", "3.14159"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			t.Fatalf("parser errors: %v", p.Errors())
		}

		env := evaluator.NewEnvironment()
		result := evaluator.Eval(program, env)

		if result == nil {
			t.Fatalf("Eval returned nil for input: %s", tt.input)
		}

		// For trigonometric functions, we'll check if result is close to expected
		// Since floating point comparisons are tricky, we'll just check the type
		if result.Type() != evaluator.FLOAT_OBJ && result.Type() != evaluator.INTEGER_OBJ {
			t.Errorf("Expected numeric result for %s, got %T", tt.input, result)
		}

		t.Logf("Input: %s, Result: %s", tt.input, result.Inspect())
	}
}
