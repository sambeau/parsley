package main

import (
	"testing"

	"github.com/sambeau/parsley/pkg/evaluator"
	"github.com/sambeau/parsley/pkg/lexer"
	"github.com/sambeau/parsley/pkg/parser"
)

func evalArrayOp(input string) evaluator.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		return &evaluator.Error{Message: p.Errors()[0]}
	}

	env := evaluator.NewEnvironment()
	return evaluator.Eval(program, env)
}

// TestArrayScalarConcat tests prepending and appending scalars to arrays
func TestArrayScalarConcat(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Prepend scalar
		{"1 ++ [2,3,4]", "[1, 2, 3, 4]"},
		{`"a" ++ ["b","c"]`, `[a, b, c]`},
		{"1 ++ []", "[1]"},
		
		// Append scalar
		{"[1,2,3] ++ 4", "[1, 2, 3, 4]"},
		{`["a","b"] ++ "c"`, `[a, b, c]`},
		{"[] ++ 1", "[1]"},
		
		// Chaining
		{"1 ++ [2] ++ 3", "[1, 2, 3]"},
		{"1 ++ 2 ++ 3", "[1, 2, 3]"},
		
		// Mixed types
		{`1 ++ ["two", 3]`, `[1, two, 3]`},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := evalArrayOp(tt.input)
			if result.Inspect() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Inspect())
			}
		})
	}
}

// TestArrayIntersection tests array intersection operator
func TestArrayIntersection(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Basic intersection
		{"[1,2,3] && [2,3,4]", "[2, 3]"},
		{"[1,2,3,4,5] && [3,4,5,6,7]", "[3, 4, 5]"},
		
		// Empty arrays
		{"[1,2,3] && []", "[]"},
		{"[] && [1,2,3]", "[]"},
		{"[] && []", "[]"},
		
		// No overlap
		{"[1,2] && [3,4]", "[]"},
		
		// Duplicates in left
		{"[1,2,2,3,3,3] && [2,3]", "[2, 3]"},
		
		// Mixed types
		{`[1, "2", 3] && [1, 3]`, "[1, 3]"},
		{`["a", "b", "c"] && ["b", "c", "d"]`, "[b, c]"},
		
		// All elements match
		{"[1,2,3] && [1,2,3]", "[1, 2, 3]"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := evalArrayOp(tt.input)
			if result.Inspect() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Inspect())
			}
		})
	}
}

// TestArrayUnion tests array union operator
func TestArrayUnion(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Basic union
		{"[1,2,3] || [3,4,5]", "[1, 2, 3, 4, 5]"},
		{"[1,2] || [3,4]", "[1, 2, 3, 4]"},
		
		// Empty arrays
		{"[1,2,3] || []", "[1, 2, 3]"},
		{"[] || [1,2,3]", "[1, 2, 3]"},
		{"[] || []", "[]"},
		
		// Duplicates
		{"[1,2,2,3] || [2,3,3,4]", "[1, 2, 3, 4]"},
		{"[1,1,1] || [1,1,1]", "[1]"},
		
		// Order preservation (left then right)
		{"[3,1,2] || [2,4,5]", "[3, 1, 2, 4, 5]"},
		
		// Mixed types
		{`[1, "2"] || ["2", 3]`, `[1, 2, 3]`},
		{`["a", "b"] || ["b", "c"]`, "[a, b, c]"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := evalArrayOp(tt.input)
			if result.Inspect() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Inspect())
			}
		})
	}
}

// TestArraySubtraction tests array subtraction operator
func TestArraySubtraction(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Basic subtraction
		{"[1,2,3,4] - [2,4]", "[1, 3]"},
		{"[1,2,3,4,5] - [3,4,5]", "[1, 2]"},
		
		// Non-existent elements
		{"[1,2,3] - [4,5]", "[1, 2, 3]"},
		
		// Empty arrays
		{"[1,2,3] - []", "[1, 2, 3]"},
		{"[] - [1,2,3]", "[]"},
		
		// All removed
		{"[1,2,3] - [1,2,3,4]", "[]"},
		
		// Duplicates
		{"[1,2,2,3,3,3] - [2,3]", "[1]"},
		
		// Mixed types
		{`[1, "2", 3] - [1, 3]`, "[2]"},
		{`["a", "b", "c"] - ["b"]`, "[a, c]"},
		
		// Order preservation
		{"[5,3,1,4,2] - [3,1]", "[5, 4, 2]"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := evalArrayOp(tt.input)
			if result.Inspect() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Inspect())
			}
		})
	}
}

// TestArrayChunking tests array chunking operator
func TestArrayChunking(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Exact division
		{"[1,2,3,4] / 2", "[[1, 2], [3, 4]]"},
		{"[1,2,3,4,5,6] / 3", "[[1, 2, 3], [4, 5, 6]]"},
		
		// Ragged last chunk
		{"[1,2,3,4,5] / 2", "[[1, 2], [3, 4], [5]]"},
		{"[1,2,3,4,5,6,7] / 3", "[[1, 2, 3], [4, 5, 6], [7]]"},
		
		// Chunk size = 1
		{"[1,2,3] / 1", "[[1], [2], [3]]"},
		
		// Chunk size > array length
		{"[1,2] / 5", "[[1, 2]]"},
		{"[1] / 10", "[[1]]"},
		
		// Empty array
		{"[] / 2", "[]"},
		
		// Large chunk
		{"[1,2,3,4,5,6,7,8,9,10] / 4", "[[1, 2, 3, 4], [5, 6, 7, 8], [9, 10]]"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := evalArrayOp(tt.input)
			if result.Inspect() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Inspect())
			}
		})
	}
}

// TestArrayChunkingErrors tests error cases for chunking
func TestArrayChunkingErrors(t *testing.T) {
	tests := []struct {
		input       string
		shouldError bool
		errorMsg    string
	}{
		{"[1,2,3] / 0", true, "chunk size must be > 0"},
		{"[1,2,3] / -1", true, "chunk size must be > 0"},
		{"[1,2,3] / -10", true, "chunk size must be > 0"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := evalArrayOp(tt.input)
			if tt.shouldError {
				if err, ok := result.(*evaluator.Error); !ok {
					t.Errorf("expected error, got %T", result)
				} else if tt.errorMsg != "" && err.Message != tt.errorMsg && 
					len(err.Message) < len(tt.errorMsg) || err.Message[:len(tt.errorMsg)] != tt.errorMsg {
					t.Errorf("expected error containing %q, got %q", tt.errorMsg, err.Message)
				}
			}
		})
	}
}

// TestStringRepetition tests string repetition operator
func TestStringRepetition(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Basic repetition
		{`"abc" * 3`, "abcabcabc"},
		{`"x" * 5`, "xxxxx"},
		{`"hello" * 2`, "hellohello"},
		
		// Zero and negative
		{`"abc" * 0`, ""},
		{`"abc" * -1`, ""},
		{`"abc" * -10`, ""},
		
		// Count = 1
		{`"test" * 1`, "test"},
		
		// Empty string
		{`"" * 5`, ""},
		{`"" * 0`, ""},
		
		// Special characters
		{`"!\n" * 3`, "!\n!\n!\n"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := evalArrayOp(tt.input)
			if result.Inspect() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Inspect())
			}
		})
	}
}

// TestArrayRepetition tests array repetition operator
func TestArrayRepetition(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Basic repetition
		{"[1,2] * 3", "[1, 2, 1, 2, 1, 2]"},
		{"[1] * 5", "[1, 1, 1, 1, 1]"},
		{`["a", "b"] * 2`, "[a, b, a, b]"},
		
		// Zero and negative
		{"[1,2,3] * 0", "[]"},
		{"[1,2,3] * -1", "[]"},
		{"[1,2,3] * -10", "[]"},
		
		// Count = 1
		{"[1,2,3] * 1", "[1, 2, 3]"},
		
		// Empty array
		{"[] * 5", "[]"},
		{"[] * 0", "[]"},
		
		// Mixed types
		{`[1, "two", 3] * 2`, `[1, two, 3, 1, two, 3]`},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := evalArrayOp(tt.input)
			if result.Inspect() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Inspect())
			}
		})
	}
}

// TestArrayOperatorChaining tests chaining multiple operators
func TestArrayOperatorChaining(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Intersection then union
		{"([1,2,3] && [2,3,4]) || [5,6]", "[2, 3, 5, 6]"},
		
		// Subtraction then concat
		{"([1,2,3,4] - [2,4]) ++ [5]", "[1, 3, 5]"},
		
		// Concat then intersection
		{"([1,2] ++ [3,4]) && [2,3,4,5]", "[2, 3, 4]"},
		
		// Union then subtraction
		{"([1,2] || [2,3]) - [2]", "[1, 3]"},
		
		// Repetition then chunking
		{"([1,2] * 3) / 2", "[[1, 2], [1, 2], [1, 2]]"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := evalArrayOp(tt.input)
			if result.Inspect() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Inspect())
			}
		})
	}
}

// TestArrayOperatorsWithVariables tests operators with variables
func TestArrayOperatorsWithVariables(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"let a = [1,2,3]; let b = [2,3,4]; a && b", "[2, 3]"},
		{"let a = [1,2]; let b = [3,4]; a || b", "[1, 2, 3, 4]"},
		{"let a = [1,2,3,4]; let b = [2,4]; a - b", "[1, 3]"},
		{"let nums = [1,2,3,4,5,6]; nums / 2", "[[1, 2], [3, 4], [5, 6]]"},
		{`let str = "ab"; str * 3`, "ababab"},
		{"let arr = [1,2]; arr * 3", "[1, 2, 1, 2, 1, 2]"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := evalArrayOp(tt.input)
			if result.Inspect() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Inspect())
			}
		})
	}
}
