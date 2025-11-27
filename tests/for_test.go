package main

import (
	"github.com/sambeau/parsley/pkg/evaluator"
	"github.com/sambeau/parsley/pkg/lexer"
	"github.com/sambeau/parsley/pkg/parser"
	"testing"
)

func TestForSimpleSyntax(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"double = fn(x) { x * 2 }; for(1,2,3) double", "2, 4, 6"},
		{"square = fn(x) { x * x }; for(1,2,3,4) square", "1, 4, 9, 16"},
		{"inc = fn(x) { x + 1 }; for(10,20,30) inc", "11, 21, 31"},
		{"for(1,2,3) fn(x){x*2}", "2, 4, 6"},
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

func TestForInSyntax(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"for(x in 1,2,3) { x * 2 }", "2, 4, 6"},
		{"for(x in 1,2,3,4) { x * x }", "1, 4, 9, 16"},
		{"for(n in 10,20,30) { n + 1 }", "11, 21, 31"},
		{"for(x in 5,15,25) { if (x > 10) { x } }", "15, 25"},
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

func TestForWithStrings(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`for("Sam") toUpper`, "S, A, M"},
		{`for("abc") toUpper`, "A, B, C"},
		{`for("XYZ") toLower`, "x, y, z"},
		{`for("Sam","Phillips") toUpper`, "SAM, PHILLIPS"},
		{`for(name in "SAM","PHILLIPS") { toLower(name) }`, "sam, phillips"},
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

func TestForEquivalenceWithMap(t *testing.T) {
	tests := []struct {
		forSyntax string
		mapSyntax string
		expected  string
	}{
		{
			"double = fn(x) { x * 2 }; for(1,2,3) double",
			"double = fn(x) { x * 2 }; map(double, 1,2,3)",
			"2, 4, 6",
		},
		{
			"for(x in 1,2,3) { x * 2 }",
			"map(fn(x){x*2}, 1,2,3)",
			"2, 4, 6",
		},
		{
			"for(x in 5,10,15) { x + 1 }",
			"map(fn(x){x+1}, 5,10,15)",
			"6, 11, 16",
		},
	}

	for _, tt := range tests {
		t.Run(tt.forSyntax, func(t *testing.T) {
			// Test for syntax
			l1 := lexer.New(tt.forSyntax)
			p1 := parser.New(l1)
			env1 := evaluator.NewEnvironment()

			program1 := p1.ParseProgram()
			if len(p1.Errors()) != 0 {
				t.Fatalf("Parser errors (for): %v", p1.Errors())
			}

			var result1 evaluator.Object
			for _, stmt := range program1.Statements {
				result1 = evaluator.Eval(stmt, env1)
			}

			// Test map syntax
			l2 := lexer.New(tt.mapSyntax)
			p2 := parser.New(l2)
			env2 := evaluator.NewEnvironment()

			program2 := p2.ParseProgram()
			if len(p2.Errors()) != 0 {
				t.Fatalf("Parser errors (map): %v", p2.Errors())
			}

			var result2 evaluator.Object
			for _, stmt := range program2.Statements {
				result2 = evaluator.Eval(stmt, env2)
			}

			if result1.Inspect() != tt.expected {
				t.Errorf("For syntax: expected %q, got %q", tt.expected, result1.Inspect())
			}

			if result2.Inspect() != tt.expected {
				t.Errorf("Map syntax: expected %q, got %q", tt.expected, result2.Inspect())
			}

			if result1.Inspect() != result2.Inspect() {
				t.Errorf("Results don't match: for=%q, map=%q", result1.Inspect(), result2.Inspect())
			}
		})
	}
}

func TestForWithBuiltins(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`for("hello") toUpper`, "H, E, L, L, O"},
		{`for("WORLD") toLower`, "w, o, r, l, d"},
		{`for("a","b","c") toUpper`, "A, B, C"},
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

func TestToUpperToLower(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`toUpper("hello")`, "HELLO"},
		{`toLower("WORLD")`, "world"},
		{`toUpper("MiXeD")`, "MIXED"},
		{`toLower("MiXeD")`, "mixed"},
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
