package main

import (
	"strings"
	"testing"

	"github.com/sambeau/parsley/pkg/evaluator"
	"github.com/sambeau/parsley/pkg/lexer"
	"github.com/sambeau/parsley/pkg/parser"
)

// Helper function to evaluate Parsley code
func evalInput(input string) evaluator.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		return &evaluator.Error{Message: strings.Join(p.Errors(), "; ")}
	}

	env := evaluator.NewEnvironment()
	return evaluator.Eval(program, env)
}

func TestBasicDictDestructuring(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic destructuring",
			input:    `let {a, b} = {a: 1, b: 2}; a`,
			expected: "1",
		},
		{
			name:     "access second value",
			input:    `let {a, b} = {a: 1, b: 2}; b`,
			expected: "2",
		},
		{
			name:     "multiple keys",
			input:    `let {x, y, z} = {x: 10, y: 20, z: 30}; x + y + z`,
			expected: "60",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalInput(tt.input)
			if result.Inspect() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.Inspect())
			}
		})
	}
}

func TestDictDestructuringWithAlias(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single alias",
			input:    `let {a as x} = {a: 5}; x`,
			expected: "5",
		},
		{
			name:     "multiple aliases",
			input:    `let {a as x, b as y} = {a: 1, b: 2}; x + y`,
			expected: "3",
		},
		{
			name:     "mixed alias and no alias",
			input:    `let {a, b as y} = {a: 10, b: 20}; a + y`,
			expected: "30",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalInput(tt.input)
			if result.Inspect() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.Inspect())
			}
		})
	}
}

func TestDictDestructuringWithRest(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "rest operator with no remaining keys",
			input:    `let {a, b, ...rest} = {a: 1, b: 2}; rest`,
			expected: "{}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalInput(tt.input)
			if result.Inspect() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.Inspect())
			}
		})
	}
}

func TestDictDestructuringMissingKeys(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "missing key becomes null",
			input:    `let {a, b, c} = {a: 1}; c`,
			expected: "null",
		},
		{
			name:     "all keys missing",
			input:    `let {x} = {}; x`,
			expected: "null",
		},
		{
			name:     "partial match",
			input:    `let {a, b, c} = {b: 2}; a`,
			expected: "null",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalInput(tt.input)
			if result.Inspect() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.Inspect())
			}
		})
	}
}

func TestDictDestructuringTypeError(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "destructure string",
			input: `let {a} = "not a dict"`,
		},
		{
			name:  "destructure number",
			input: `let {a} = 42`,
		},
		{
			name:  "destructure array",
			input: `let {a} = [1, 2, 3]`,
		},
		{
			name:  "destructure null",
			input: `let {a} = null`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalInput(tt.input)
			if result.Type() != evaluator.ERROR_OBJ {
				t.Errorf("expected error, got %v", result)
			}
		})
	}
}

func TestDictDestructuringAssignment(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "reassign with destructuring",
			input:    `let a = 0; let b = 0; {a, b} = {a: 10, b: 20}; a`,
			expected: "10",
		},
		{
			name:     "reassign second value",
			input:    `let a = 0; let b = 0; {a, b} = {a: 10, b: 20}; b`,
			expected: "20",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalInput(tt.input)
			if result.Inspect() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.Inspect())
			}
		})
	}
}

func TestNestedDictDestructuring(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "nested one level",
			input:    `let {a: {b}} = {a: {b: 5}}; b`,
			expected: "5",
		},
		{
			name:     "nested two levels",
			input:    `let {a: {b: {c}}} = {a: {b: {c: 10}}}; c`,
			expected: "10",
		},
		{
			name:     "nested with alias",
			input:    `let {a: {b as x}} = {a: {b: 7}}; x`,
			expected: "7",
		},
		{
			name:     "nested with multiple keys",
			input:    `let {a: {b, c}} = {a: {b: 1, c: 2}}; b + c`,
			expected: "3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalInput(tt.input)
			if result.Inspect() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.Inspect())
			}
		})
	}
}

func TestDictDestructuringComplexExamples(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "destructure object with computation",
			input: `
				let obj = {x: 10, y: 20, z: 30};
				let {x, y} = obj;
				x * y
			`,
			expected: "200",
		},
		{
			name: "multiple destructuring statements",
			input: `
				let {a} = {a: 1};
				let {b} = {b: 2};
				let {c} = {c: 3};
				a + b + c
			`,
			expected: "6",
		},
		{
			name: "destructure then modify",
			input: `
				let {a, b} = {a: 5, b: 10};
				a = a * 2;
				b = b * 2;
				a + b
			`,
			expected: "30",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalInput(tt.input)
			if result.Inspect() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.Inspect())
			}
		})
	}
}

// TestDictionaryMethodsWithThis tests that user-defined methods can use 'this'
func TestDictionaryMethodsWithThis(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic this binding",
			input:    `let user = {name: "Sam", greet: fn() { "Hi, " + this.name }}; user.greet()`,
			expected: `Hi, Sam`,
		},
		{
			name:     "this with numeric property",
			input:    `let calc = {value: 10, double: fn() { this.value * 2 }}; calc.double()`,
			expected: "20",
		},
		{
			name:     "method with arguments",
			input:    `let calc = {value: 10, add: fn(x) { this.value + x }}; calc.add(5)`,
			expected: "15",
		},
		{
			name:     "multiple this references",
			input:    `let p = {first: "John", last: "Doe", full: fn() { this.first + " " + this.last }}; p.full()`,
			expected: `John Doe`,
		},
		{
			name:     "nested method call",
			input:    `let p = {name: "Bob", greet: fn() { "Hi " + this.name }, hello: fn() { this.greet() + "!" }}; p.hello()`,
			expected: `Hi Bob!`,
		},
		{
			name:     "built-in methods still work with user method",
			input:    `let d = {a: 1, b: 2, getA: fn() { this.a }}; d.has("a")`,
			expected: `true`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalInput(tt.input)
			if result.Inspect() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.Inspect())
			}
		})
	}
}

// TestDictDeleteMethod tests the .delete(key) method for dictionaries
func TestDictDeleteMethod(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "delete single key",
			input:    `let d = {a: 1, b: 2, c: 3}; d.delete("b"); d.keys().sort()`,
			expected: `[a, c]`,
		},
		{
			name:     "delete first key",
			input:    `let d = {a: 1, b: 2}; d.delete("a"); d`,
			expected: `{b: 2}`,
		},
		{
			name:     "delete last key",
			input:    `let d = {a: 1, b: 2}; d.delete("b"); d`,
			expected: `{a: 1}`,
		},
		{
			name:     "delete non-existent key (no error)",
			input:    `let d = {a: 1}; d.delete("x"); d`,
			expected: `{a: 1}`,
		},
		{
			name:     "delete all keys",
			input:    `let d = {a: 1, b: 2}; d.delete("a"); d.delete("b"); d`,
			expected: `{}`,
		},
		{
			name:     "delete returns null",
			input:    `let d = {a: 1}; d.delete("a")`,
			expected: `null`,
		},
		{
			name:     "delete with variable key",
			input:    `let d = {a: 1, b: 2}; let key = "a"; d.delete(key); d`,
			expected: `{b: 2}`,
		},
		{
			name:     "delete does not affect user function named delete",
			input:    `let d = {a: 1, delete: fn(){"custom"}}; d.delete("a"); d["delete"]()`,
			expected: `custom`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalInput(tt.input)
			if result.Inspect() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.Inspect())
			}
		})
	}
}

// TestDeleteKeywordNoLongerReserved tests that 'delete' can be used as identifier
func TestDeleteKeywordNoLongerReserved(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "delete as variable name",
			input:    `let delete = "value"; delete`,
			expected: `value`,
		},
		{
			name:     "delete as function name",
			input:    `let delete = fn(x) { x * 2 }; delete(5)`,
			expected: `10`,
		},
		{
			name:     "delete as dict key",
			input:    `let d = {delete: 42}; d.delete`,
			expected: `42`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalInput(tt.input)
			if result.Inspect() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.Inspect())
			}
		})
	}
}
