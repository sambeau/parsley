// nullish_test.go - Tests for the nullish coalescing operator (??)

package main

import (
	"testing"

	"github.com/sambeau/parsley/pkg/evaluator"
	"github.com/sambeau/parsley/pkg/lexer"
	"github.com/sambeau/parsley/pkg/parser"
)

func testEvalNullish(input string) evaluator.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := evaluator.NewEnvironment()
	return evaluator.Eval(program, env)
}

// Test basic nullish coalescing behavior
func TestNullishCoalescing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Basic null handling
		{`null ?? "default"`, "default"},
		{`"value" ?? "default"`, "value"},
		// Empty string is NOT null
		{`"" ?? "default"`, ""},
		// Zero is NOT null
		{`0 ?? 42`, "0"},
		// False is NOT null
		{`false ?? true`, "false"},
		// Integer values
		{`42 ?? 0`, "42"},
		// String values
		{`"hello" ?? "world"`, "hello"},
		// Boolean values
		{`true ?? false`, "true"},
	}

	for _, tt := range tests {
		result := testEvalNullish(tt.input)
		if result == nil {
			t.Errorf("For input '%s': evaluation returned nil", tt.input)
			continue
		}
		resultStr := result.Inspect()
		if resultStr != tt.expected {
			t.Errorf("For input '%s': expected %q, got %q", tt.input, tt.expected, resultStr)
		}
	}
}

// Test nullish with variables
func TestNullishWithVariables(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Defined variable returns its value
		{`let x = "value"; x ?? "default"`, "value"},
		// Null variable triggers default
		{`let x = null; x ?? "default"`, "default"},
		// Variable with zero returns zero
		{`let x = 0; x ?? 42`, "0"},
		// Using _ as null
		{`_ ?? "default"`, "default"},
	}

	for _, tt := range tests {
		result := testEvalNullish(tt.input)
		if result == nil {
			t.Errorf("For input '%s': evaluation returned nil", tt.input)
			continue
		}
		resultStr := result.Inspect()
		if resultStr != tt.expected {
			t.Errorf("For input '%s': expected %q, got %q", tt.input, tt.expected, resultStr)
		}
	}
}

// Test chained nullish operations
func TestNullishChained(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`null ?? null ?? "final"`, "final"},
		{`null ?? "second" ?? "third"`, "second"},
		{`"first" ?? null ?? "third"`, "first"},
	}

	for _, tt := range tests {
		result := testEvalNullish(tt.input)
		if result == nil {
			t.Errorf("For input '%s': evaluation returned nil", tt.input)
			continue
		}
		resultStr := result.Inspect()
		if resultStr != tt.expected {
			t.Errorf("For input '%s': expected %q, got %q", tt.input, tt.expected, resultStr)
		}
	}
}

// Test nullish with dictionary access
func TestNullishWithDictAccess(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Missing key returns null, nullish provides default
		{`let d = { a: 1 }; d.b ?? "missing"`, "missing"},
		// Existing key returns value
		{`let d = { a: 1 }; d.a ?? "missing"`, "1"},
		// Nested missing key
		{`let d = { a: { b: 1 } }; d.a.c ?? "default"`, "default"},
	}

	for _, tt := range tests {
		result := testEvalNullish(tt.input)
		if result == nil {
			t.Errorf("For input '%s': evaluation returned nil", tt.input)
			continue
		}
		resultStr := result.Inspect()
		if resultStr != tt.expected {
			t.Errorf("For input '%s': expected %q, got %q", tt.input, tt.expected, resultStr)
		}
	}
}

// Test nullish with function returns
func TestNullishWithFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`let f = fn() { null }; f() ?? "default"`, "default"},
		{`let f = fn() { "result" }; f() ?? "default"`, "result"},
	}

	for _, tt := range tests {
		result := testEvalNullish(tt.input)
		if result == nil {
			t.Errorf("For input '%s': evaluation returned nil", tt.input)
			continue
		}
		resultStr := result.Inspect()
		if resultStr != tt.expected {
			t.Errorf("For input '%s': expected %q, got %q", tt.input, tt.expected, resultStr)
		}
	}
}

// Test that nullish is properly short-circuit evaluated
func TestNullishShortCircuit(t *testing.T) {
	// Test 1: Right side should NOT be evaluated when left is non-null
	input1 := `
		let counter = 0
		let inc = fn() { counter = counter + 1; "evaluated" }
		let result = "value" ?? inc()
		counter
	`
	result1 := testEvalNullish(input1)
	if result1.Inspect() != "0" {
		t.Errorf("Short-circuit failed: right side was evaluated when left was non-null. Counter = %s", result1.Inspect())
	}

	// Test 2: Right side SHOULD be evaluated when left is null
	input2 := `
		let counter = 0
		let inc = fn() { counter = counter + 1; "evaluated" }
		let result = null ?? inc()
		counter
	`
	result2 := testEvalNullish(input2)
	if result2.Inspect() != "1" {
		t.Errorf("Right side should have been evaluated when left was null. Counter = %s", result2.Inspect())
	}
}

// Test precedence of nullish operator
func TestNullishPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Nullish has lower precedence than arithmetic
		{`null ?? 1 + 2`, "3"},
		{`null ?? 2 * 3`, "6"},
		// Nullish has lower precedence than comparison
		{`null ?? 1 < 2`, "true"},
		// Parentheses override precedence
		{`(null ?? 1) + 2`, "3"},
	}

	for _, tt := range tests {
		result := testEvalNullish(tt.input)
		if result == nil {
			t.Errorf("For input '%s': evaluation returned nil", tt.input)
			continue
		}
		resultStr := result.Inspect()
		if resultStr != tt.expected {
			t.Errorf("For input '%s': expected %q, got %q", tt.input, tt.expected, resultStr)
		}
	}
}

// Test nullish with complex expressions
func TestNullishWithExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// With explicit null from function
		{`let getVal = fn() { null }; getVal() ?? "default"`, "default"},
		// With nested nullish
		{`let a = null; let b = null; (a ?? b) ?? "fallback"`, "fallback"},
	}

	for _, tt := range tests {
		result := testEvalNullish(tt.input)
		if result == nil {
			t.Errorf("For input '%s': evaluation returned nil", tt.input)
			continue
		}
		resultStr := result.Inspect()
		if resultStr != tt.expected {
			t.Errorf("For input '%s': expected %q, got %q", tt.input, tt.expected, resultStr)
		}
	}
}
