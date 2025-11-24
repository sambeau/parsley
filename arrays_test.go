package main

import (
	"pars/pkg/evaluator"
	"pars/pkg/lexer"
	"pars/pkg/parser"
	"testing"
)

func TestArrayLiterals(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1,2,3", "1, 2, 3"},
		{"1", "1"},
		{"1,2", "1, 2"},
		{`"hello","world"`, "hello, world"},
		{"1,2,3,4,5", "1, 2, 3, 4, 5"},
		{`1,"two",3`, `1, two, 3`},
		{"1+2,3*4,5-1", "3, 12, 4"},
		{"true,false", "true, false"},
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
		{"xs = 1,2,3; xs", "1, 2, 3"},
		{"arr = 5,10,15; arr", "5, 10, 15"},
		{`names = "Alice","Bob"; names`, "Alice, Bob"},
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
		{"double = fn(x) { x * 2 }; map(double, 1,2,3)", "2, 4, 6"},
		{"square = fn(x) { x * x }; map(square, 1,2,3,4)", "1, 4, 9, 16"},
		{"inc = fn(x) { x + 1 }; map(inc, 10,20,30)", "11, 21, 31"},
		{`upper = fn(s) { s }; map(upper, "a","b","c")`, "a, b, c"},
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
		{"gt10 = fn(x) { if (x > 10) then return x }; map(gt10, 5,15,25,8,3,12)", "15, 25, 12"},
		{"gt5 = fn(x) { if (x > 5) then return x }; map(gt5, 1,2,3,4,5,6,7,8)", "6, 7, 8"},
		{"positive = fn(x) { if (x > 0) then return x }; map(positive, -1,5,-3,10)", "5, 10"},
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
		{"double = fn(x) { x * 2 }; xs = 1,2,3; map(double, xs)", "2, 4, 6"},
		{"square = fn(x) { x * x }; nums = 2,3,4; map(square, nums)", "4, 9, 16"},
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
		{"transform = fn(x) { x * 2 + 1 }; map(transform, 1,2,3)", "3, 5, 7"},
		{"calc = fn(x) { if (x > 2) then x * 2 else x + 1 }; map(calc, 1,2,3,4)", "2, 3, 6, 8"},
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
		{"x = 1; y = 2; x,y", "1, 2"},
		{"1+2,3*4", "3, 12"},
		{"(1,2)", "1, 2"},
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
