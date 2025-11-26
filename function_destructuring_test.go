package main

import (
	"testing"

	"github.com/sambeau/parsley/pkg/evaluator"
	"github.com/sambeau/parsley/pkg/lexer"
	"github.com/sambeau/parsley/pkg/parser"
)

// Helper function to evaluate Parsley code
func evalCode(input string) evaluator.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		return &evaluator.Error{Message: p.Errors()[0]}
	}

	env := evaluator.NewEnvironment()
	return evaluator.Eval(program, env)
}

func TestArrayDestructuringInFunctionParams(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "extract first element",
			input: `
				let head = fn([first, rest]) { first };
				head([1, 2, 3, 4, 5])
			`,
			expected: "1",
		},
		{
			name: "extract rest elements",
			input: `
				let tail = fn([first, rest]) { rest };
				tail([1, 2, 3, 4, 5])
			`,
			expected: "2, 3, 4, 5",
		},
		{
			name: "multiple parameters",
			input: `
				let swap = fn([a, b]) { [b, a] };
				swap([10, 20])
			`,
			expected: "20, 10",
		},
		{
			name: "empty array",
			input: `
				let getFirst = fn([x, y]) { x };
				getFirst([])
			`,
			expected: "null",
		},
		{
			name: "single element array",
			input: `
				let extract = fn([a, b, c]) { a };
				extract([42])
			`,
			expected: "42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalCode(tt.input)
			if result.Inspect() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.Inspect())
			}
		})
	}
}

func TestDictDestructuringInFunctionParams(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "basic dict destructuring",
			input: `
				let add = fn({a, b}) { a + b };
				add({a: 10, b: 20})
			`,
			expected: "30",
		},
		{
			name: "extract string values",
			input: `
				let greet = fn({name, age}) { name };
				greet({name: "Alice", age: 30})
			`,
			expected: "Alice",
		},
		{
			name: "missing keys become null",
			input: `
				let test = fn({a, b, c}) { c };
				test({a: 1, b: 2})
			`,
			expected: "null",
		},
		{
			name: "with alias",
			input: `
				let process = fn({name, age as years}) { years };
				process({name: "Bob", age: 25})
			`,
			expected: "25",
		},
		{
			name: "with rest operator",
			input: `
				let extract = fn({a, ...rest}) { rest };
				extract({a: 1, b: 2, c: 3})
			`,
			expected: "", // Don't check exact string due to map iteration order
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalCode(tt.input)
			// Skip empty expected values (tests that check behavior without exact output)
			if tt.expected == "" {
				// For rest operator test, just verify it's a dictionary
				if tt.name == "with rest operator" {
					dict, ok := result.(*evaluator.Dictionary)
					if !ok {
						t.Errorf("expected dictionary, got %T", result)
					} else if len(dict.Pairs) != 2 {
						t.Errorf("expected dictionary with 2 keys, got %d", len(dict.Pairs))
					}
				}
			} else if result.Inspect() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.Inspect())
			}
		})
	}
}

func TestNestedDestructuringInFunctionParams(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "nested dict destructuring",
			input: `
				let getEmail = fn({profile: {email}}) { email };
				getEmail({profile: {email: "test@example.com", name: "Test"}})
			`,
			expected: "test@example.com",
		},
		{
			name: "nested with multiple keys",
			input: `
				let process = fn({data: {x, y}}) { x + y };
				process({data: {x: 5, y: 10}})
			`,
			expected: "15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalCode(tt.input)
			if result.Inspect() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.Inspect())
			}
		})
	}
}

func TestMixedParametersDestructuring(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "simple param and array destructuring",
			input: `
				let test = fn(x, [a, b]) { x + a + b };
				test(1, [2, 3])
			`,
			expected: "6",
		},
		{
			name: "simple param and dict destructuring",
			input: `
				let test = fn(x, {a, b}) { x + a + b };
				test(5, {a: 10, b: 20})
			`,
			expected: "35",
		},
		{
			name: "array and dict destructuring",
			input: `
				let test = fn([x, y], {a, b}) { x + y + a + b };
				test([1, 2], {a: 3, b: 4})
			`,
			expected: "10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalCode(tt.input)
			if result.Inspect() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.Inspect())
			}
		})
	}
}

func TestDestructuringInClosures(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "closure with array destructuring",
			input: `
				let makeAdder = fn(x) {
					fn([a, b]) { x + a + b }
				};
				let add10 = makeAdder(10);
				add10([5, 3])
			`,
			expected: "18",
		},
		{
			name: "closure with dict destructuring",
			input: `
				let makeMultiplier = fn(factor) {
					fn({value}) { factor * value }
				};
				let double = makeMultiplier(2);
				double({value: 21})
			`,
			expected: "42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalCode(tt.input)
			if result.Inspect() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.Inspect())
			}
		})
	}
}

func TestHigherOrderFunctionsWithDestructuring(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "map-like with array destructuring",
			input: `
				let process = fn(f, arr) {
					let helper = fn(idx) {
						if(idx >= len(arr)) { [] } else {
							let item = arr[idx];
							let rest = helper(idx + 1);
							[f(item)] ++ rest
						}
					};
					helper(0)
				};
				let extractFirst = fn([x, y]) { x };
				process(extractFirst, [[1, 2], [3, 4], [5, 6]])
			`,
			expected: "1, 3, 5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalCode(tt.input)
			if result.Inspect() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.Inspect())
			}
		})
	}
}

func TestPracticalExamples(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "calculate distance from point",
			input: `
				let distance = fn({x, y}) { sqrt(x * x + y * y) };
				distance({x: 3, y: 4})
			`,
			expected: "5",
		},
		{
			name: "format person info",
			input: `
				let formatPerson = fn({name, age, city}) {
					name + " (" + age + ") from " + city
				};
				formatPerson({name: "Alice", age: 30, city: "NYC"})
			`,
			expected: "Alice (30) from NYC",
		},
		/* TODO: This test requires rest operator support in array destructuring
		{
			name: "array sum with destructuring",
			input: `
				let sum = fn([first, rest]) {
					if(first == false) {
						0
					} else {
						if(rest == false) {
							first
						} else {
							first + sum(rest)
						}
					}
				};
				sum([1, 2, 3, 4, 5])
			`,
			expected: "15",
		},
		*/
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalCode(tt.input)
			if result.Inspect() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.Inspect())
			}
		})
	}
}
