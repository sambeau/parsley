package main

import (
	"strings"
	"testing"

	"github.com/sambeau/parsley/pkg/evaluator"
	"github.com/sambeau/parsley/pkg/lexer"
	"github.com/sambeau/parsley/pkg/parser"
)

func evalDictOp(input string) evaluator.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		return &evaluator.Error{Message: p.Errors()[0]}
	}

	env := evaluator.NewEnvironment()
	return evaluator.Eval(program, env)
}

// TestDictionaryIntersection tests dictionary intersection operator
func TestDictionaryIntersection(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		checkKeys []string
		checkNotKeys []string
	}{
		{
			name:         "partial overlap",
			input:        "{a: 1, b: 2, c: 3} && {b: 5, c: 6, d: 7}",
			checkKeys:    []string{"b", "c"},
			checkNotKeys: []string{"a", "d"},
		},
		{
			name:         "no overlap",
			input:        "{a: 1, b: 2} && {c: 3, d: 4}",
			checkKeys:    []string{},
			checkNotKeys: []string{"a", "b", "c", "d"},
		},
		{
			name:         "all keys match",
			input:        "{a: 1, b: 2} && {a: 5, b: 6}",
			checkKeys:    []string{"a", "b"},
			checkNotKeys: []string{},
		},
		{
			name:         "left empty",
			input:        "{} && {a: 1, b: 2}",
			checkKeys:    []string{},
			checkNotKeys: []string{"a", "b"},
		},
		{
			name:         "right empty",
			input:        "{a: 1, b: 2} && {}",
			checkKeys:    []string{},
			checkNotKeys: []string{"a", "b"},
		},
		{
			name:         "both empty",
			input:        "{} && {}",
			checkKeys:    []string{},
			checkNotKeys: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalDictOp(tt.input)
			resultStr := result.Inspect()

			// Check expected keys are present
			for _, key := range tt.checkKeys {
				if !strings.Contains(resultStr, key+":") {
					t.Errorf("expected key %q to be present in %q", key, resultStr)
				}
			}

			// Check unexpected keys are not present
			for _, key := range tt.checkNotKeys {
				if strings.Contains(resultStr, key+":") {
					t.Errorf("expected key %q to NOT be present in %q", key, resultStr)
				}
			}
		})
	}
}

// TestDictionaryIntersectionValues tests that intersection preserves left values
func TestDictionaryIntersectionValues(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		getValue string
		expected string
	}{
		{
			name:     "left values preserved",
			input:    "let result = {a: 1, b: 2} && {a: 99, b: 88}; result.a",
			getValue: "a",
			expected: "1",
		},
		{
			name:     "left values preserved for b",
			input:    "let result = {a: 1, b: 2} && {a: 99, b: 88}; result.b",
			getValue: "b",
			expected: "2",
		},
		{
			name:     "string values from left",
			input:    `let result = {name: "Alice", age: 30} && {name: "Bob", age: 25}; result.name`,
			getValue: "name",
			expected: "Alice",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalDictOp(tt.input)
			if result.Inspect() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Inspect())
			}
		})
	}
}

// TestDictionarySubtraction tests dictionary subtraction operator
func TestDictionarySubtraction(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		checkKeys    []string
		checkNotKeys []string
	}{
		{
			name:         "remove some keys",
			input:        "{a: 1, b: 2, c: 3} - {b: 0}",
			checkKeys:    []string{"a", "c"},
			checkNotKeys: []string{"b"},
		},
		{
			name:         "remove multiple keys",
			input:        "{a: 1, b: 2, c: 3, d: 4} - {b: 99, d: 88}",
			checkKeys:    []string{"a", "c"},
			checkNotKeys: []string{"b", "d"},
		},
		{
			name:         "non-existent keys",
			input:        "{a: 1, b: 2} - {c: 3, d: 4}",
			checkKeys:    []string{"a", "b"},
			checkNotKeys: []string{"c", "d"},
		},
		{
			name:         "remove all keys",
			input:        "{a: 1, b: 2} - {a: 0, b: 0, c: 0}",
			checkKeys:    []string{},
			checkNotKeys: []string{"a", "b", "c"},
		},
		{
			name:         "empty right",
			input:        "{a: 1, b: 2} - {}",
			checkKeys:    []string{"a", "b"},
			checkNotKeys: []string{},
		},
		{
			name:         "empty left",
			input:        "{} - {a: 1, b: 2}",
			checkKeys:    []string{},
			checkNotKeys: []string{"a", "b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalDictOp(tt.input)
			resultStr := result.Inspect()

			// Check expected keys are present
			for _, key := range tt.checkKeys {
				if !strings.Contains(resultStr, key+":") {
					t.Errorf("expected key %q to be present in %q", key, resultStr)
				}
			}

			// Check unexpected keys are not present
			for _, key := range tt.checkNotKeys {
				if strings.Contains(resultStr, key+":") {
					t.Errorf("expected key %q to NOT be present in %q", key, resultStr)
				}
			}
		})
	}
}

// TestDictionarySubtractionIgnoresValues tests that right values are ignored
func TestDictionarySubtractionIgnoresValues(t *testing.T) {
	input := `{a: 1, b: 2, c: 3} - {b: 999, d: 888}`
	result := evalDictOp(input)
	resultStr := result.Inspect()

	// b should be removed regardless of its value in right dict
	if strings.Contains(resultStr, "b:") {
		t.Errorf("expected key 'b' to be removed, got %q", resultStr)
	}

	// a and c should remain
	if !strings.Contains(resultStr, "a:") || !strings.Contains(resultStr, "c:") {
		t.Errorf("expected keys 'a' and 'c' to remain, got %q", resultStr)
	}

	// d should not be added
	if strings.Contains(resultStr, "d:") {
		t.Errorf("expected key 'd' to NOT be present, got %q", resultStr)
	}
}

// TestDictionaryMergeInteraction tests interaction with existing ++ merge
func TestDictionaryMergeInteraction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "merge then intersect",
			input:    `let merged = {a: 1} ++ {b: 2}; let other = {a: 99, b: 88}; let result = merged && other; result.a`,
			expected: "1",
		},
		{
			name:     "intersect then merge",
			input:    `let intersected = {a: 1, b: 2, c: 3} && {a: 5, b: 6}; let result = intersected ++ {d: 4}; result.d`,
			expected: "4",
		},
		{
			name:     "merge then subtract",
			input:    `let merged = {a: 1, b: 2} ++ {c: 3}; let result = merged - {b: 0}; result.a`,
			expected: "1",
		},
		{
			name:     "subtract then merge",
			input:    `let subtracted = {a: 1, b: 2, c: 3} - {b: 0}; let result = subtracted ++ {b: 99}; result.b`,
			expected: "99",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalDictOp(tt.input)
			if result.Inspect() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Inspect())
			}
		})
	}
}

// TestDictionaryOperatorsWithNestedDicts tests shallow operation
func TestDictionaryOperatorsWithNestedDicts(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check func(string) bool
	}{
		{
			name:  "intersection with nested dicts",
			input: `{a: {x: 1}, b: 2} && {a: {y: 2}, c: 3}`,
			check: func(result string) bool {
				// Should have key 'a' (shallow check)
				return strings.Contains(result, "a:")
			},
		},
		{
			name:  "subtraction with nested dicts",
			input: `{a: {x: 1}, b: 2, c: 3} - {a: {}}`,
			check: func(result string) bool {
				// 'a' should be removed (shallow operation)
				return !strings.Contains(result, "a:") && strings.Contains(result, "b:")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalDictOp(tt.input)
			if !tt.check(result.Inspect()) {
				t.Errorf("check failed for input %q, got %q", tt.input, result.Inspect())
			}
		})
	}
}

// TestDictionaryOperatorsWithVariables tests operators with variables
func TestDictionaryOperatorsWithVariables(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check func(string) bool
	}{
		{
			name:  "intersection with vars",
			input: `let d1 = {a: 1, b: 2}; let d2 = {b: 3, c: 4}; d1 && d2`,
			check: func(result string) bool {
				return strings.Contains(result, "b:") && !strings.Contains(result, "a:") && !strings.Contains(result, "c:")
			},
		},
		{
			name:  "subtraction with vars",
			input: `let d1 = {a: 1, b: 2, c: 3}; let d2 = {b: 0}; d1 - d2`,
			check: func(result string) bool {
				return strings.Contains(result, "a:") && strings.Contains(result, "c:") && !strings.Contains(result, "b:")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalDictOp(tt.input)
			if !tt.check(result.Inspect()) {
				t.Errorf("check failed for %q, got %q", tt.name, result.Inspect())
			}
		})
	}
}

// TestDictionaryOperatorsPracticalUseCase tests real-world scenarios
func TestDictionaryOperatorsPracticalUseCase(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check func(string) bool
	}{
		{
			name: "filter config to common keys",
			input: `let defaults = {timeout: 30, retries: 3, logging: true};
					let userConfig = {timeout: 60, cache: true};
					let commonSettings = defaults && userConfig;
					commonSettings`,
			check: func(result string) bool {
				// Should only have 'timeout' (common key)
				return strings.Contains(result, "timeout:") && 
					   !strings.Contains(result, "retries:") && 
					   !strings.Contains(result, "logging:") &&
					   !strings.Contains(result, "cache:")
			},
		},
		{
			name: "remove sensitive keys",
			input: `let data = {id: 1, name: "Alice", password: "secret", token: "xyz"};
					let cleaned = data - {password: null, token: null};
					cleaned`,
			check: func(result string) bool {
				// Should have id and name, not password or token
				return strings.Contains(result, "id:") && 
					   strings.Contains(result, "name:") &&
					   !strings.Contains(result, "password:") &&
					   !strings.Contains(result, "token:")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalDictOp(tt.input)
			if !tt.check(result.Inspect()) {
				t.Errorf("check failed for %q, got %q", tt.name, result.Inspect())
			}
		})
	}
}
