package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/sambeau/parsley/pkg/evaluator"
	"github.com/sambeau/parsley/pkg/lexer"
	"github.com/sambeau/parsley/pkg/parser"
)

// Helper function to evaluate Parsley code and check for errors
func testDatetimeCode(input string) (evaluator.Object, bool) {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		return &evaluator.Error{Message: strings.Join(p.Errors(), "; ")}, true
	}

	env := evaluator.NewEnvironment()
	result := evaluator.Eval(program, env)

	if _, ok := result.(*evaluator.Error); ok {
		return result, true
	}

	return result, false
}

func TestDatetimeNow(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		wantErr bool
	}{
		{
			name:    "now() with no arguments",
			code:    `let dt = now(); log(dt.year, dt.month, dt.day);`,
			wantErr: false,
		},
		{
			name:    "now() with arguments should error",
			code:    `let dt = now(123);`,
			wantErr: true,
		},
		{
			name:    "access datetime fields",
			code:    `let dt = now(); log(dt.year, dt.month, dt.day, dt.hour, dt.minute, dt.second, dt.weekday, dt.unix, dt.iso);`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, hasErr := testDatetimeCode(tt.code)
			if hasErr != tt.wantErr {
				t.Errorf("testDatetimeCode() hasErr = %v, wantErr %v", hasErr, tt.wantErr)
			}
		})
	}
}

func TestDatetimeFromString(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		wantErr bool
	}{
		{
			name:    "parse ISO 8601 datetime",
			code:    `let dt = time("2024-01-15T10:30:00Z"); log(dt.year, dt.month, dt.day);`,
			wantErr: false,
		},
		{
			name:    "parse date only",
			code:    `let dt = time("2024-12-25"); log(dt.year, dt.month, dt.day);`,
			wantErr: false,
		},
		{
			name:    "invalid datetime string",
			code:    `let dt = time("not a date");`,
			wantErr: true,
		},
		{
			name:    "parse without timezone",
			code:    `let dt = time("2024-01-15T10:30:00"); log(dt.year);`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, hasErr := testDatetimeCode(tt.code)
			if hasErr != tt.wantErr {
				t.Errorf("testDatetimeCode() hasErr = %v, wantErr %v", hasErr, tt.wantErr)
			}
		})
	}
}

func TestDatetimeFromUnix(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		wantErr bool
	}{
		{
			name:    "parse Unix timestamp",
			code:    `let dt = time(1704110400); log(dt.year, dt.month, dt.day);`,
			wantErr: false,
		},
		{
			name:    "zero timestamp",
			code:    `let dt = time(0); log(dt.year);`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, hasErr := testDatetimeCode(tt.code)
			if hasErr != tt.wantErr {
				t.Errorf("testDatetimeCode() hasErr = %v, wantErr %v", hasErr, tt.wantErr)
			}
		})
	}
}

func TestDatetimeFromDictionary(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		wantErr bool
	}{
		{
			name:    "create from full dictionary",
			code:    `let dt = time({year: 2024, month: 7, day: 4, hour: 12, minute: 30, second: 45}); log(dt.iso);`,
			wantErr: false,
		},
		{
			name:    "create from partial dictionary (date only)",
			code:    `let dt = time({year: 2024, month: 12, day: 25}); log(dt.iso);`,
			wantErr: false,
		},
		{
			name:    "missing year field should error",
			code:    `let dt = time({month: 12, day: 25}); log(dt.iso);`,
			wantErr: true,
		},
		{
			name:    "missing month field should error",
			code:    `let dt = time({year: 2024, day: 25}); log(dt.iso);`,
			wantErr: true,
		},
		{
			name:    "missing day field should error",
			code:    `let dt = time({year: 2024, month: 12}); log(dt.iso);`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, hasErr := testDatetimeCode(tt.code)
			if hasErr != tt.wantErr {
				t.Errorf("testDatetimeCode() hasErr = %v, wantErr %v", hasErr, tt.wantErr)
			}
		})
	}
}

func TestDatetimeDelta(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		wantErr bool
	}{
		{
			name:    "add days",
			code:    `let dt = time("2024-01-01T00:00:00Z", {days: 7}); log(dt.iso);`,
			wantErr: false,
		},
		{
			name:    "subtract days",
			code:    `let dt = time("2024-01-15T00:00:00Z", {days: -10}); log(dt.iso);`,
			wantErr: false,
		},
		{
			name:    "add months",
			code:    `let dt = time("2024-01-01T00:00:00Z", {months: 3}); log(dt.iso);`,
			wantErr: false,
		},
		{
			name:    "add years",
			code:    `let dt = time("2024-01-01T00:00:00Z", {years: 1}); log(dt.iso);`,
			wantErr: false,
		},
		{
			name:    "complex delta",
			code:    `let dt = time("2024-01-01T12:30:00Z", {years: 1, months: 2, days: 15, hours: 3, minutes: 45, seconds: 30}); log(dt.iso);`,
			wantErr: false,
		},
		{
			name:    "delta with now()",
			code:    `let dt = time(now(), {days: 7}); log(dt.year);`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, hasErr := testDatetimeCode(tt.code)
			if hasErr != tt.wantErr {
				t.Errorf("testDatetimeCode() hasErr = %v, wantErr %v", hasErr, tt.wantErr)
			}
		})
	}
}

func TestDatetimeFieldAccess(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		wantErr bool
	}{
		{
			name:    "access year field",
			code:    `let dt = time("2024-01-15T10:30:00Z"); log(dt.year);`,
			wantErr: false,
		},
		{
			name:    "access month field",
			code:    `let dt = time("2024-01-15T10:30:00Z"); log(dt.month);`,
			wantErr: false,
		},
		{
			name:    "access day field",
			code:    `let dt = time("2024-01-15T10:30:00Z"); log(dt.day);`,
			wantErr: false,
		},
		{
			name:    "access hour field",
			code:    `let dt = time("2024-01-15T10:30:00Z"); log(dt.hour);`,
			wantErr: false,
		},
		{
			name:    "access minute field",
			code:    `let dt = time("2024-01-15T10:30:00Z"); log(dt.minute);`,
			wantErr: false,
		},
		{
			name:    "access second field",
			code:    `let dt = time("2024-01-15T10:30:00Z"); log(dt.second);`,
			wantErr: false,
		},
		{
			name:    "access weekday field",
			code:    `let dt = time("2024-01-15T10:30:00Z"); log(dt.weekday);`,
			wantErr: false,
		},
		{
			name:    "access unix field",
			code:    `let dt = time("2024-01-15T10:30:00Z"); log(dt.unix);`,
			wantErr: false,
		},
		{
			name:    "access iso field",
			code:    `let dt = time("2024-01-15T10:30:00Z"); log(dt.iso);`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, hasErr := testDatetimeCode(tt.code)
			if hasErr != tt.wantErr {
				if hasErr {
					t.Errorf("testDatetimeCode() got error: %v", result)
				} else {
					t.Errorf("testDatetimeCode() hasErr = %v, wantErr %v", hasErr, tt.wantErr)
				}
			}
		})
	}
}

func TestDatetimeTemplateFormatting(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		wantErr bool
	}{
		{
			name:    "format with template literal",
			code:    "let dt = time(\"2024-01-15T10:30:00Z\"); let formatted = `{dt.year}-{dt.month}-{dt.day}`; log(formatted);",
			wantErr: false,
		},
		{
			name:    "format datetime with time",
			code:    "let dt = time(\"2024-01-15T10:30:00Z\"); let formatted = `{dt.year}-{dt.month}-{dt.day} {dt.hour}:{dt.minute}`; log(formatted);",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, hasErr := testDatetimeCode(tt.code)
			if hasErr != tt.wantErr {
				t.Errorf("testDatetimeCode() hasErr = %v, wantErr %v", hasErr, tt.wantErr)
			}
		})
	}
}

func TestDatetimeErrors(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "now() with arguments",
			code:    `let dt = now(123);`,
			wantErr: true,
			errMsg:  "wrong number of arguments",
		},
		{
			name:    "time() with no arguments",
			code:    `let dt = time();`,
			wantErr: true,
			errMsg:  "wrong number of arguments",
		},
		{
			name:    "time() with too many arguments",
			code:    `let dt = time("2024-01-01", {days: 1}, "extra");`,
			wantErr: true,
			errMsg:  "wrong number of arguments",
		},
		{
			name:    "time() with invalid type",
			code:    `let dt = time(true);`,
			wantErr: true,
			errMsg:  "must be a string, integer, or dictionary",
		},
		{
			name:    "time() with invalid string",
			code:    `let dt = time("not a datetime");`,
			wantErr: true,
			errMsg:  "invalid datetime string",
		},
		{
			name:    "time() with non-dictionary delta",
			code:    `let dt = time("2024-01-01", "not a dict");`,
			wantErr: true,
			errMsg:  "must be a dictionary",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, hasErr := testDatetimeCode(tt.code)
			if hasErr != tt.wantErr {
				t.Errorf("testDatetimeCode() hasErr = %v, wantErr %v", hasErr, tt.wantErr)
			}
			if hasErr && tt.errMsg != "" {
				errObj, ok := result.(*evaluator.Error)
				if !ok {
					t.Errorf("expected Error object, got %T", result)
				} else if !strings.Contains(errObj.Message, tt.errMsg) {
					t.Errorf("error message %q does not contain %q", errObj.Message, tt.errMsg)
				}
			}
		})
	}
}

func TestDatetimeSpecificValues(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		checkFn func(evaluator.Object) error
	}{
		{
			name: "parse specific date",
			code: `let dt = time("2024-01-15T10:30:00Z"); dt.year`,
			checkFn: func(obj evaluator.Object) error {
				intObj, ok := obj.(*evaluator.Integer)
				if !ok {
					return fmt.Errorf("expected Integer, got %T", obj)
				}
				if intObj.Value != 2024 {
					return fmt.Errorf("expected year 2024, got %d", intObj.Value)
				}
				return nil
			},
		},
		{
			name: "unix timestamp conversion",
			code: `let dt = time(1704110400); dt.year`,
			checkFn: func(obj evaluator.Object) error {
				intObj, ok := obj.(*evaluator.Integer)
				if !ok {
					return fmt.Errorf("expected Integer, got %T", obj)
				}
				if intObj.Value != 2024 {
					return fmt.Errorf("expected year 2024, got %d", intObj.Value)
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, hasErr := testDatetimeCode(tt.code)
			if hasErr {
				t.Fatalf("testDatetimeCode() unexpected error: %v", result)
			}
			if err := tt.checkFn(result); err != nil {
				t.Errorf("checkFn() error: %v", err)
			}
		})
	}
}

func TestDatetimeComparisons(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name:     "less than true",
			code:     `time("2024-01-15T00:00:00Z") < time("2024-01-20T00:00:00Z")`,
			expected: "true",
		},
		{
			name:     "less than false",
			code:     `time("2024-01-20T00:00:00Z") < time("2024-01-15T00:00:00Z")`,
			expected: "false",
		},
		{
			name:     "greater than true",
			code:     `time("2024-01-20T00:00:00Z") > time("2024-01-15T00:00:00Z")`,
			expected: "true",
		},
		{
			name:     "greater than false",
			code:     `time("2024-01-15T00:00:00Z") > time("2024-01-20T00:00:00Z")`,
			expected: "false",
		},
		{
			name:     "less than or equal true (less)",
			code:     `time("2024-01-15T00:00:00Z") <= time("2024-01-20T00:00:00Z")`,
			expected: "true",
		},
		{
			name:     "less than or equal true (equal)",
			code:     `time("2024-01-15T00:00:00Z") <= time("2024-01-15T00:00:00Z")`,
			expected: "true",
		},
		{
			name:     "greater than or equal true (greater)",
			code:     `time("2024-01-20T00:00:00Z") >= time("2024-01-15T00:00:00Z")`,
			expected: "true",
		},
		{
			name:     "greater than or equal true (equal)",
			code:     `time("2024-01-15T00:00:00Z") >= time("2024-01-15T00:00:00Z")`,
			expected: "true",
		},
		{
			name:     "equal true",
			code:     `time("2024-01-15T10:30:00Z") == time("2024-01-15T10:30:00Z")`,
			expected: "true",
		},
		{
			name:     "equal false",
			code:     `time("2024-01-15T10:30:00Z") == time("2024-01-15T10:31:00Z")`,
			expected: "false",
		},
		{
			name:     "not equal true",
			code:     `time("2024-01-15T10:30:00Z") != time("2024-01-15T10:31:00Z")`,
			expected: "true",
		},
		{
			name:     "not equal false",
			code:     `time("2024-01-15T10:30:00Z") != time("2024-01-15T10:30:00Z")`,
			expected: "false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, hasErr := testDatetimeCode(tt.code)
			if hasErr {
				t.Fatalf("testDatetimeCode() unexpected error: %v", result)
			}
			if result.Inspect() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result.Inspect())
			}
		})
	}
}

func TestDatetimeArithmetic(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name:     "subtract datetimes (returns Duration)",
			code:     `let diff = time("2024-01-20T00:00:00Z") - time("2024-01-15T00:00:00Z"); diff.seconds`,
			expected: "432000", // 5 days * 86400 seconds/day
		},
		{
			name:     "add seconds to datetime",
			code:     `let dt = time("2024-01-15T00:00:00Z") + 86400; dt.day`,
			expected: "16", // Next day
		},
		{
			name:     "subtract seconds from datetime",
			code:     `let dt = time("2024-01-15T00:00:00Z") - 86400; dt.day`,
			expected: "14", // Previous day
		},
		{
			name:     "add seconds to datetime (commutative)",
			code:     `let dt = 86400 + time("2024-01-15T00:00:00Z"); dt.day`,
			expected: "16",
		},
		{
			name:     "add week to datetime",
			code:     `let dt = time("2024-01-01T12:00:00Z") + 604800; dt.day`,
			expected: "8", // 7 days later
		},
		{
			name:     "subtract week from datetime",
			code:     `let dt = time("2024-01-15T12:00:00Z") - 604800; dt.day`,
			expected: "8", // 7 days earlier
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, hasErr := testDatetimeCode(tt.code)
			if hasErr {
				t.Fatalf("testDatetimeCode() unexpected error: %v", result)
			}
			if result.Inspect() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result.Inspect())
			}
		})
	}
}

func TestDatetimeTypeField(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name:     "check __type field exists",
			code:     `let dt = now(); dt.__type`,
			expected: "datetime",
		},
		{
			name:     "check __type on time() result",
			code:     `time("2024-01-15T00:00:00Z").__type`,
			expected: "datetime",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, hasErr := testDatetimeCode(tt.code)
			if hasErr {
				t.Fatalf("testDatetimeCode() unexpected error: %v", result)
			}
			if result.Inspect() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result.Inspect())
			}
		})
	}
}
