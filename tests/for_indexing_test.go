package main

import (
	"sort"
	"strings"
	"testing"

	"github.com/sambeau/parsley/pkg/evaluator"
	"github.com/sambeau/parsley/pkg/lexer"
	"github.com/sambeau/parsley/pkg/parser"
)

// TestForArrayWithIndex tests for loops with index parameter on arrays
func TestForArrayWithIndex(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`for(i, x in [10, 20, 30]) { i }`,
			`[0, 1, 2]`,
		},
		{
			`for(i, x in [10, 20, 30]) { x }`,
			`[10, 20, 30]`,
		},
		{
			`for(i, x in [10, 20, 30]) { i * 10 + x }`,
			`[10, 30, 50]`,
		},
		{
			`for(idx, val in ["a", "b", "c"]) { idx + ":" + val }`,
			`[0:a, 1:b, 2:c]`,
		},
		{
			// Empty array
			`for(i, x in []) { i }`,
			``,
		},
		{
			// Single element
			`for(i, x in [42]) { [i, x] }`,
			`[[0, 42]]`,
		},
		{
			// Verify index starts at 0
			`let result = for(i, x in [5, 6, 7]) { if (i == 0) { x } }; result`,
			`[5]`,
		},
		{
			// Create array of [index, value] pairs
			`for(i, x in [100, 200, 300]) { [i, x] }`,
			`[[0, 100], [1, 200], [2, 300]]`,
		},
	}

	for _, tt := range tests {
		evaluated := testEvalForIndexing(tt.input)
		testArrayOutputForIndexing(t, evaluated, tt.expected, tt.input)
	}
}

// TestForStringWithIndex tests for loops with index parameter on strings
func TestForStringWithIndex(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`for(i, c in "abc") { i }`,
			`[0, 1, 2]`,
		},
		{
			`for(i, c in "abc") { c }`,
			`[a, b, c]`,
		},
		{
			`for(i, c in "hello") { i + ":" + c }`,
			`[0:h, 1:e, 2:l, 3:l, 4:o]`,
		},
		{
			// Empty string
			`for(i, c in "") { i }`,
			``,
		},
		{
			// Single character
			`for(i, c in "x") { [i, c] }`,
			`[[0, x]]`,
		},
		{
			// Unicode characters
			`for(i, c in "🎉🎊") { i }`,
			`[0, 1]`,
		},
	}

	for _, tt := range tests {
		evaluated := testEvalForIndexing(tt.input)
		testArrayOutputForIndexing(t, evaluated, tt.expected, tt.input)
	}
}

// TestForBackwardCompatibility ensures single-parameter for loops still work
func TestForBackwardCompatibility(t *testing.T) {
	tests := []struct {
		input     string
		expected  string
		unordered bool // Set to true for dictionary iteration where order is non-deterministic
	}{
		{
			// Single parameter - element only
			input:    `for(x in [1, 2, 3]) { x * 2 }`,
			expected: `[2, 4, 6]`,
		},
		{
			// Simple form still works
			input:    `for([1, 2, 3]) fn(x) { x + 10 }`,
			expected: `[11, 12, 13]`,
		},
		{
			// Dictionary iteration unchanged (key, value)
			// Note: Dictionary iteration order is non-deterministic in Go
			input:     `for(k, v in {a: 1, b: 2}) { k }`,
			expected:  `[a, b]`,
			unordered: true,
		},
		{
			// String iteration with single param
			input:    `for(c in "hi") { toUpper(c) }`,
			expected: `[H, I]`,
		},
	}

	for _, tt := range tests {
		evaluated := testEvalForIndexing(tt.input)
		if tt.unordered {
			testArrayOutputUnordered(t, evaluated, tt.expected, tt.input)
		} else {
			testArrayOutputForIndexing(t, evaluated, tt.expected, tt.input)
		}
	}
}

// TestForIndexEdgeCases tests edge cases for indexed iteration
func TestForIndexEdgeCases(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			// Filter with index - only even indices
			`for(i, x in [10, 20, 30, 40]) { if (i % 2 == 0) { x } }`,
			`[10, 30]`,
		},
		{
			// Filter with index - only odd indices
			`for(i, x in [10, 20, 30, 40]) { if (i % 2 == 1) { x } }`,
			`[20, 40]`,
		},
		{
			// Use both index and value in computation
			`for(i, x in [5, 5, 5]) { x + i }`,
			`[5, 6, 7]`,
		},
		{
			// Nested arrays with indexing
			`for(i, row in [[1, 2], [3, 4]]) { i }`,
			`[0, 1]`,
		},
		{
			// Large array (verify index increments correctly)
			`let arr = [0, 0, 0, 0, 0]; for(i, x in arr) { i }`,
			`[0, 1, 2, 3, 4]`,
		},
	}

	for _, tt := range tests {
		evaluated := testEvalForIndexing(tt.input)
		testArrayOutputForIndexing(t, evaluated, tt.expected, tt.input)
	}
}

// TestForIndexWithVariableNames tests various parameter naming
func TestForIndexWithVariableNames(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			// Traditional names
			`for(index, value in [7, 8, 9]) { index }`,
			`[0, 1, 2]`,
		},
		{
			// Short names
			`for(i, v in [7, 8, 9]) { v }`,
			`[7, 8, 9]`,
		},
		{
			// Descriptive names
			`for(position, item in ["first", "second"]) { position + 1 }`,
			`[1, 2]`,
		},
		{
			// Underscore for unused index
			`for(_, x in [10, 20, 30]) { x }`,
			`[10, 20, 30]`,
		},
		{
			// Underscore for unused value
			`for(i, _ in [10, 20, 30]) { i * 100 }`,
			`[0, 100, 200]`,
		},
	}

	for _, tt := range tests {
		evaluated := testEvalForIndexing(tt.input)
		testArrayOutputForIndexing(t, evaluated, tt.expected, tt.input)
	}
}

// TestForIndexErrorCases tests error handling
func TestForIndexErrorCases(t *testing.T) {
	tests := []struct {
		input         string
		expectedError string
	}{
		{
			// Function with 3 parameters not supported
			`for([1, 2, 3]) fn(a, b, c) { a }`,
			`function passed to for must take 1 or 2 parameters, got 3`,
		},
	}

	for _, tt := range tests {
		evaluated := testEvalForIndexing(tt.input)
		errObj, ok := evaluated.(*evaluator.Error)
		if !ok {
			t.Errorf("Expected error for input: %s, got: %T (%+v)", tt.input, evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expectedError {
			t.Errorf("Wrong error message.\nExpected: %s\nGot: %s\nInput: %s",
				tt.expectedError, errObj.Message, tt.input)
		}
	}
}

// TestForIndexPracticalExamples tests real-world use cases
func TestForIndexPracticalExamples(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			// Enumerate pattern - create numbered list
			`for(i, item in ["apple", "banana", "cherry"]) { (i + 1) + ". " + item }`,
			`[1. apple, 2. banana, 3. cherry]`,
		},
		{
			// Find index of element
			`let items = ["a", "b", "c"]; for(i, x in items) { if (x == "b") { i } }`,
			`[1]`,
		},
		{
			// Skip first element using index
			`for(i, x in [10, 20, 30, 40]) { if (i > 0) { x } }`,
			`[20, 30, 40]`,
		},
		{
			// Take first N elements using index
			`for(i, x in [1, 2, 3, 4, 5]) { if (i < 3) { x } }`,
			`[1, 2, 3]`,
		},
	}

	for _, tt := range tests {
		evaluated := testEvalForIndexing(tt.input)
		testArrayOutputForIndexing(t, evaluated, tt.expected, tt.input)
	}
}

// Helper function to evaluate Parsley code
func testEvalForIndexing(input string) evaluator.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := evaluator.NewEnvironment()
	return evaluator.Eval(program, env)
}

// Helper function to test array output
func testArrayOutputForIndexing(t *testing.T, obj evaluator.Object, expected string, input string) {
	arr, ok := obj.(*evaluator.Array)
	if !ok {
		// Check if it's an error
		if errObj, isErr := obj.(*evaluator.Error); isErr {
			t.Errorf("Evaluation error for input: %s\nError: %s", input, errObj.Message)
			return
		}
		// If expected is empty string, obj might be empty array or null
		if expected == "" {
			if obj == evaluator.NULL {
				return
			}
			if arr, ok := obj.(*evaluator.Array); ok && len(arr.Elements) == 0 {
				return
			}
		}
		t.Errorf("Expected Array, got %T (%+v) for input: %s", obj, obj, input)
		return
	}

	if expected == "" {
		if len(arr.Elements) != 0 {
			t.Errorf("Expected empty array, got %s for input: %s", arr.Inspect(), input)
		}
		return
	}

	result := arr.Inspect()
	if result != expected {
		t.Errorf("Wrong result.\nExpected: %s\nGot: %s\nInput: %s", expected, result, input)
	}
}

// Helper function to test array output ignoring order (for dictionary iteration)
func testArrayOutputUnordered(t *testing.T, obj evaluator.Object, expected string, input string) {
	arr, ok := obj.(*evaluator.Array)
	if !ok {
		if errObj, isErr := obj.(*evaluator.Error); isErr {
			t.Errorf("Evaluation error for input: %s\nError: %s", input, errObj.Message)
			return
		}
		t.Errorf("Expected Array, got %T (%+v) for input: %s", obj, obj, input)
		return
	}

	// Parse expected string like "[a, b]" into sorted elements
	expectedStr := strings.TrimPrefix(strings.TrimSuffix(expected, "]"), "[")
	expectedParts := strings.Split(expectedStr, ", ")
	sort.Strings(expectedParts)

	// Get actual elements and sort them
	var actualParts []string
	for _, elem := range arr.Elements {
		actualParts = append(actualParts, elem.Inspect())
	}
	sort.Strings(actualParts)

	// Compare sorted slices
	if len(expectedParts) != len(actualParts) {
		t.Errorf("Wrong number of elements.\nExpected: %v\nGot: %v\nInput: %s", expectedParts, actualParts, input)
		return
	}
	for i := range expectedParts {
		if expectedParts[i] != actualParts[i] {
			t.Errorf("Wrong elements (order-independent).\nExpected: %v\nGot: %v\nInput: %s", expectedParts, actualParts, input)
			return
		}
	}
}
