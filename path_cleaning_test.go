package main

import (
	"testing"

	"github.com/sambeau/parsley/pkg/evaluator"
	"github.com/sambeau/parsley/pkg/lexer"
	"github.com/sambeau/parsley/pkg/parser"
)

// Helper function for path cleaning tests
func evalPathCleaningTest(t *testing.T, input string) evaluator.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	env := evaluator.NewEnvironment()
	return evaluator.Eval(program, env)
}

// TestPathCleaning tests Rob Pike's cleanname algorithm implementation
// See: https://9p.io/sys/doc/lexnames.html
func TestPathCleaning(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Rule 2: Eliminate . (current directory)
		{
			name:     "eliminate_single_dot",
			input:    `@/./foo/./bar`,
			expected: `{__type: path, absolute: true, components: , foo, bar}`,
		},
		{
			name:     "eliminate_single_dot_relative",
			input:    `@./a/./b`,
			expected: `{__type: path, absolute: false, components: ., a, b}`,
		},

		// Rule 3: Eliminate .. and the preceding element
		{
			name:     "eliminate_dotdot_with_preceding",
			input:    `@/foo/../bar`,
			expected: `{__type: path, absolute: true, components: , bar}`,
		},
		{
			name:     "eliminate_multiple_dotdot",
			input:    `@/a/b/../../c`,
			expected: `{__type: path, absolute: true, components: , c}`,
		},
		{
			name:     "eliminate_dotdot_relative",
			input:    `@./a/b/../c`,
			expected: `{__type: path, absolute: false, components: ., a, c}`,
		},

		// Rule 4: Eliminate .. at beginning of rooted path
		{
			name:     "eliminate_dotdot_at_root",
			input:    `@/../foo`,
			expected: `{__type: path, absolute: true, components: , foo}`,
		},
		{
			name:     "eliminate_multiple_dotdot_at_root",
			input:    `@/../../../foo`,
			expected: `{__type: path, absolute: true, components: , foo}`,
		},

		// Rule 5: Leave .. intact at beginning of non-rooted path
		{
			name:     "preserve_dotdot_at_start_relative",
			input:    `@../foo`,
			expected: `{__type: path, absolute: false, components: .., foo}`,
		},
		{
			name:     "preserve_multiple_dotdot_at_start",
			input:    `@../../foo`,
			expected: `{__type: path, absolute: false, components: .., .., foo}`,
		},
		{
			name:     "preserve_dotdot_then_eliminate",
			input:    `@../a/../b`,
			expected: `{__type: path, absolute: false, components: .., b}`,
		},

		// Multiple slashes (handled by split)
		{
			name:     "multiple_slashes",
			input:    `@/a//b///c`,
			expected: `{__type: path, absolute: true, components: , a, b, c}`,
		},

		// Complex combinations
		{
			name:     "complex_cleaning",
			input:    `@./a/b/../../c/./d`,
			expected: `{__type: path, absolute: false, components: ., c, d}`,
		},
		{
			name:     "deeply_nested_cleanup",
			input:    `@/a/b/c/../../d/../e`,
			expected: `{__type: path, absolute: true, components: , a, e}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalPathCleaningTest(t, tt.input)
			if result.Inspect() != tt.expected {
				t.Errorf("for input %s:\nexpected: %s\n     got: %s", tt.input, tt.expected, result.Inspect())
			}
		})
	}
}

// TestPathCleaningToString tests that cleaned paths convert back to clean strings
func TestPathCleaningToString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "absolute_cleaned",
			input:    `let p = @/foo/../bar; p.string`,
			expected: `/bar`,
		},
		{
			name:     "relative_cleaned",
			input:    `let p = @./a/b/../c; p.string`,
			expected: `./a/c`,
		},
		{
			name:     "root_dotdot_eliminated",
			input:    `let p = @/../../../foo; p.string`,
			expected: `/foo`,
		},
		{
			name:     "relative_dotdot_preserved",
			input:    `let p = @../../../foo; p.string`,
			expected: `../../../foo`,
		},
		{
			name:     "complex_cleaning_to_string",
			input:    `let p = @./a/b/../../c/./d; p.string`,
			expected: `./c/d`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalPathCleaningTest(t, tt.input)
			if result.Inspect() != tt.expected {
				t.Errorf("for input %s:\nexpected: %s\n     got: %s", tt.input, tt.expected, result.Inspect())
			}
		})
	}
}

// TestURLPathCleaning tests that URL paths are also cleaned
func TestURLPathCleaning(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "url_path_cleaned",
			input:    `let u = @https://example.com/a/../b; u.string`,
			expected: `https://example.com/b`,
		},
		{
			name:     "url_path_dot_eliminated",
			input:    `let u = @https://example.com/./a/./b; u.string`,
			expected: `https://example.com/a/b`,
		},
		{
			name:     "url_path_complex",
			input:    `let u = @https://example.com/a/b/../../c; u.string`,
			expected: `https://example.com/c`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalPathCleaningTest(t, tt.input)
			if result.Inspect() != tt.expected {
				t.Errorf("for input %s:\nexpected: %s\n     got: %s", tt.input, tt.expected, result.Inspect())
			}
		})
	}
}
