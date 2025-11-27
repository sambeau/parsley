package main

import (
	"testing"

	"github.com/sambeau/parsley/pkg/evaluator"
	"github.com/sambeau/parsley/pkg/lexer"
	"github.com/sambeau/parsley/pkg/parser"
)

func testDurationCode(code string) (evaluator.Object, bool) {
	l := lexer.New(code)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		return &evaluator.String{Value: p.Errors()[0]}, true
	}

	env := evaluator.NewEnvironment()
	result := evaluator.Eval(program, env)
	if result != nil && result.Type() == evaluator.ERROR_OBJ {
		return result, true
	}

	return result, false
}

func TestDurationLiterals(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name:     "simple seconds duration",
			code:     `let d = @30s; d.seconds`,
			expected: "30",
		},
		{
			name:     "simple minutes duration",
			code:     `let d = @5m; d.seconds`,
			expected: "300", // 5 * 60
		},
		{
			name:     "simple hours duration",
			code:     `let d = @2h; d.seconds`,
			expected: "7200", // 2 * 3600
		},
		{
			name:     "simple days duration",
			code:     `let d = @7d; d.seconds`,
			expected: "604800", // 7 * 86400
		},
		{
			name:     "simple weeks duration",
			code:     `let d = @2w; d.seconds`,
			expected: "1209600", // 2 * 7 * 86400
		},
		{
			name:     "simple months duration",
			code:     `let d = @6mo; d.months`,
			expected: "6",
		},
		{
			name:     "simple years duration",
			code:     `let d = @1y; d.months`,
			expected: "12",
		},
		{
			name:     "compound duration hours and minutes",
			code:     `let d = @2h30m; d.seconds`,
			expected: "9000", // 2*3600 + 30*60
		},
		{
			name:     "compound duration days and hours",
			code:     `let d = @3d12h; d.seconds`,
			expected: "302400", // 3*86400 + 12*3600
		},
		{
			name:     "compound duration with all units",
			code:     `let d = @1y2mo3w4d5h6m7s; d.months`,
			expected: "14", // 1*12 + 2
		},
		{
			name:     "compound duration seconds calculation",
			code:     `let d = @1y2mo3w4d5h6m7s; d.seconds`,
			expected: "2178367", // 3*604800 + 4*86400 + 5*3600 + 6*60 + 7
		},
		{
			name:     "totalSeconds exists for pure seconds durations",
			code:     `let d = @2h30m; d.totalSeconds`,
			expected: "9000",
		},
		{
			name:     "totalSeconds is null for month-based durations",
			code:     `let d = @1y; d.totalSeconds`,
			expected: "null",
		},
		{
			name:     "duration type field",
			code:     `let d = @1h; d.__type`,
			expected: "duration",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, hasErr := testDurationCode(tt.code)
			if hasErr {
				t.Fatalf("testDurationCode() unexpected error: %v", result)
			}
			if result.Inspect() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result.Inspect())
			}
		})
	}
}

func TestDurationArithmetic(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name:     "add two durations (seconds only)",
			code:     `let d = @2h + @30m; d.seconds`,
			expected: "9000",
		},
		{
			name:     "subtract two durations (seconds only)",
			code:     `let d = @1d - @6h; d.seconds`,
			expected: "64800", // 86400 - 21600
		},
		{
			name:     "add durations with months",
			code:     `let d = @1y + @6mo; d.months`,
			expected: "18",
		},
		{
			name:     "subtract durations with months",
			code:     `let d = @2y - @3mo; d.months`,
			expected: "21", // 24 - 3
		},
		{
			name:     "add mixed durations (months and seconds)",
			code:     `let d1 = @1y2mo; let d2 = @3d4h; let d = d1 + d2; d.months`,
			expected: "14",
		},
		{
			name:     "add mixed durations seconds part",
			code:     `let d1 = @1y2mo; let d2 = @3d4h; let d = d1 + d2; d.seconds`,
			expected: "273600", // 3*86400 + 4*3600
		},
		{
			name:     "multiply duration by integer",
			code:     `let d = @2h * 3; d.seconds`,
			expected: "21600", // 7200 * 3
		},
		{
			name:     "divide duration by integer",
			code:     `let d = @1d / 2; d.seconds`,
			expected: "43200", // 86400 / 2
		},
		{
			name:     "multiply month-based duration",
			code:     `let d = @1y * 2; d.months`,
			expected: "24",
		},
		{
			name:     "divide month-based duration",
			code:     `let d = @2y / 4; d.months`,
			expected: "6",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, hasErr := testDurationCode(tt.code)
			if hasErr {
				t.Fatalf("testDurationCode() unexpected error: %v", result)
			}
			if result.Inspect() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result.Inspect())
			}
		})
	}
}

func TestDatetimeDurationOperations(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name:     "add duration to datetime",
			code:     `let dt = time("2024-01-15T00:00:00Z") + @2d; dt.day`,
			expected: "17",
		},
		{
			name:     "add hours to datetime",
			code:     `let dt = time("2024-01-15T10:00:00Z") + @3h; dt.hour`,
			expected: "13",
		},
		{
			name:     "subtract duration from datetime",
			code:     `let dt = time("2024-01-15T00:00:00Z") - @2d; dt.day`,
			expected: "13",
		},
		{
			name:     "add months to datetime",
			code:     `let dt = time("2024-01-31T00:00:00Z") + @1mo; dt.month`,
			expected: "3", // Jan 31 + 1 month = Mar 2 (Feb 31 doesn't exist)
		},
		{
			name:     "add months to datetime (day normalization)",
			code:     `let dt = time("2024-01-31T00:00:00Z") + @1mo; dt.day`,
			expected: "2", // Normalized from Feb 31 to Mar 2
		},
		{
			name:     "add years to datetime",
			code:     `let dt = time("2024-06-15T00:00:00Z") + @1y; dt.year`,
			expected: "2025",
		},
		{
			name:     "add compound duration to datetime",
			code:     `let dt = time("2024-01-01T00:00:00Z") + @1y2mo3d; dt.month`,
			expected: "3", // Jan + 14 months = March of next year
		},
		{
			name:     "datetime minus datetime returns duration",
			code:     `let diff = time("2024-01-20T00:00:00Z") - time("2024-01-15T00:00:00Z"); diff.seconds`,
			expected: "432000", // 5 days
		},
		{
			name:     "datetime minus datetime has no months",
			code:     `let diff = time("2024-01-20T00:00:00Z") - time("2024-01-15T00:00:00Z"); diff.months`,
			expected: "0",
		},
		{
			name:     "datetime minus datetime totalSeconds",
			code:     `let diff = time("2024-01-20T00:00:00Z") - time("2024-01-15T00:00:00Z"); diff.totalSeconds`,
			expected: "432000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, hasErr := testDurationCode(tt.code)
			if hasErr {
				t.Fatalf("testDurationCode() unexpected error: %v", result)
			}
			if result.Inspect() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result.Inspect())
			}
		})
	}
}

func TestDurationComparison(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name:     "compare equal durations",
			code:     `@2h == @7200s`,
			expected: "true",
		},
		{
			name:     "compare unequal durations",
			code:     `@2h != @3h`,
			expected: "true",
		},
		{
			name:     "less than comparison",
			code:     `@1h < @2h`,
			expected: "true",
		},
		{
			name:     "greater than comparison",
			code:     `@3h > @2h`,
			expected: "true",
		},
		{
			name:     "less than or equal",
			code:     `@2h <= @2h`,
			expected: "true",
		},
		{
			name:     "greater than or equal",
			code:     `@3h >= @2h`,
			expected: "true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, hasErr := testDurationCode(tt.code)
			if hasErr {
				t.Fatalf("testDurationCode() unexpected error: %v", result)
			}
			if result.Inspect() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result.Inspect())
			}
		})
	}
}

func TestDurationErrors(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		expectError bool
	}{
		{
			name:        "cannot compare durations with months",
			code:        `@1y < @12mo`,
			expectError: true,
		},
		{
			name:        "cannot add datetime to duration",
			code:        `@1d + time("2024-01-01T00:00:00Z")`,
			expectError: true,
		},
		{
			name:        "division by zero",
			code:        `@1d / 0`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, hasErr := testDurationCode(tt.code)
			if tt.expectError && !hasErr {
				t.Errorf("Expected error but got result: %s", result.Inspect())
			}
			if !tt.expectError && hasErr {
				t.Errorf("Unexpected error: %v", result)
			}
		})
	}
}
