package main

import (
	"testing"

	"github.com/sambeau/parsley/pkg/evaluator"
	"github.com/sambeau/parsley/pkg/lexer"
	"github.com/sambeau/parsley/pkg/parser"
)

// Helper functions for testing
func testEvalDatetime(input string) evaluator.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := evaluator.NewEnvironment()
	return evaluator.Eval(program, env)
}

func testExpectedDatetime(t *testing.T, input string, obj evaluator.Object, expected string) {
	if obj == nil {
		t.Errorf("For input '%s': got nil object", input)
		return
	}

	if err, ok := obj.(*evaluator.Error); ok {
		t.Errorf("For input '%s': got error: %s", input, err.Message)
		return
	}

	actual := obj.Inspect()
	if actual != expected {
		t.Errorf("For input '%s': expected %s, got %s", input, expected, actual)
	}
}

func TestDatetimeLiteralBasic(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`@2024-12-25`, `{__type: "datetime", year: 2024, month: 12, day: 25, hour: 0, minute: 0, second: 0, unix: 1735084800, weekday: "Wednesday", iso: "2024-12-25T00:00:00Z"}`},
		{`@2024-01-01`, `{__type: "datetime", year: 2024, month: 1, day: 1, hour: 0, minute: 0, second: 0, unix: 1704067200, weekday: "Monday", iso: "2024-01-01T00:00:00Z"}`},
		{`@2024-06-15`, `{__type: "datetime", year: 2024, month: 6, day: 15, hour: 0, minute: 0, second: 0, unix: 1718409600, weekday: "Saturday", iso: "2024-06-15T00:00:00Z"}`},
	}

	for _, tt := range tests {
		evaluated := testEvalDatetime(tt.input)
		testExpectedDatetime(t, tt.input, evaluated, tt.expected)
	}
}

func TestDatetimeLiteralWithTime(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`@2024-12-25T14:30:00`, `{__type: "datetime", year: 2024, month: 12, day: 25, hour: 14, minute: 30, second: 0, unix: 1735136400, weekday: "Wednesday", iso: "2024-12-25T14:30:00Z"}`},
		{`@2024-01-15T09:45:30`, `{__type: "datetime", year: 2024, month: 1, day: 15, hour: 9, minute: 45, second: 30, unix: 1705314330, weekday: "Monday", iso: "2024-01-15T09:45:30Z"}`},
		{`@2024-06-01T23:59:59`, `{__type: "datetime", year: 2024, month: 6, day: 1, hour: 23, minute: 59, second: 59, unix: 1717286399, weekday: "Saturday", iso: "2024-06-01T23:59:59Z"}`},
	}

	for _, tt := range tests {
		evaluated := testEvalDatetime(tt.input)
		testExpectedDatetime(t, tt.input, evaluated, tt.expected)
	}
}

func TestDatetimeLiteralWithTimezone(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`@2024-12-25T14:30:00Z`, `{__type: "datetime", year: 2024, month: 12, day: 25, hour: 14, minute: 30, second: 0, unix: 1735136400, weekday: "Wednesday", iso: "2024-12-25T14:30:00Z"}`},
		{`@2024-12-25T14:30:00-05:00`, `{__type: "datetime", year: 2024, month: 12, day: 25, hour: 19, minute: 30, second: 0, unix: 1735154400, weekday: "Wednesday", iso: "2024-12-25T19:30:00Z"}`},
		{`@2024-06-15T08:00:00+08:00`, `{__type: "datetime", year: 2024, month: 6, day: 15, hour: 0, minute: 0, second: 0, unix: 1718380800, weekday: "Saturday", iso: "2024-06-15T00:00:00Z"}`},
	}

	for _, tt := range tests {
		evaluated := testEvalDatetime(tt.input)
		testExpectedDatetime(t, tt.input, evaluated, tt.expected)
	}
}

func TestDatetimeLiteralFieldAccess(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`@2024-12-25.year`, `2024`},
		{`@2024-12-25.month`, `12`},
		{`@2024-12-25.day`, `25`},
		{`@2024-12-25T14:30:00.hour`, `14`},
		{`@2024-12-25T14:30:00.minute`, `30`},
		{`@2024-12-25T14:30:00.second`, `0`},
		{`@2024-12-25.weekday`, `"Wednesday"`},
		{`@2024-12-25.iso`, `"2024-12-25T00:00:00Z"`},
	}

	for _, tt := range tests {
		evaluated := testEvalDatetime(tt.input)
		testExpectedDatetime(t, tt.input, evaluated, tt.expected)
	}
}

func TestDatetimeLiteralComparison(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`@2024-12-25 < @2024-12-26`, `true`},
		{`@2024-12-25 > @2024-12-24`, `true`},
		{`@2024-12-25 == @2024-12-25`, `true`},
		{`@2024-12-25 != @2024-12-26`, `true`},
		{`@2024-12-25T14:30:00 < @2024-12-25T14:30:01`, `true`},
		{`@2024-12-25T14:30:00 >= @2024-12-25T14:30:00`, `true`},
		{`@2024-01-01 <= @2024-12-31`, `true`},
	}

	for _, tt := range tests {
		evaluated := testEvalDatetime(tt.input)
		testExpectedDatetime(t, tt.input, evaluated, tt.expected)
	}
}

func TestDatetimeLiteralArithmetic(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Add seconds to datetime
		{`@2024-12-25 + 86400`, `{__type: "datetime", year: 2024, month: 12, day: 26, hour: 0, minute: 0, second: 0, unix: 1735171200, weekday: "Thursday", iso: "2024-12-26T00:00:00Z"}`},
		{`@2024-12-25T14:30:00 + 3600`, `{__type: "datetime", year: 2024, month: 12, day: 25, hour: 15, minute: 30, second: 0, unix: 1735140000, weekday: "Wednesday", iso: "2024-12-25T15:30:00Z"}`},
		// Subtract seconds from datetime
		{`@2024-12-25 - 86400`, `{__type: "datetime", year: 2024, month: 12, day: 24, hour: 0, minute: 0, second: 0, unix: 1734998400, weekday: "Tuesday", iso: "2024-12-24T00:00:00Z"}`},
		// Commutative addition
		{`86400 + @2024-12-25`, `{__type: "datetime", year: 2024, month: 12, day: 26, hour: 0, minute: 0, second: 0, unix: 1735171200, weekday: "Thursday", iso: "2024-12-26T00:00:00Z"}`},
	}

	for _, tt := range tests {
		evaluated := testEvalDatetime(tt.input)
		testExpectedDatetime(t, tt.input, evaluated, tt.expected)
	}
}

func TestDatetimeLiteralSubtraction(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// BREAKING CHANGE: Datetime - datetime now returns Duration
		{`let d = @2024-12-26 - @2024-12-25; d.seconds`, `86400`},
		{`let d = @2024-12-25 - @2024-12-24; d.seconds`, `86400`},
		{`let d = @2024-12-25T15:30:00 - @2024-12-25T14:30:00; d.seconds`, `3600`},
		{`let d = @2024-01-01 - @2024-01-01; d.seconds`, `0`},
	}

	for _, tt := range tests {
		evaluated := testEvalDatetime(tt.input)
		testExpectedDatetime(t, tt.input, evaluated, tt.expected)
	}
}

func TestDatetimeLiteralEquivalence(t *testing.T) {
	// Test that @syntax produces same result as time() function
	tests := []struct {
		literal  string
		function string
	}{
		{`@2024-12-25`, `time("2024-12-25")`},
		{`@2024-12-25T14:30:00`, `time("2024-12-25T14:30:00")`},
		{`@2024-01-15T09:45:30Z`, `time("2024-01-15T09:45:30Z")`},
	}

	for _, tt := range tests {
		literalResult := testEvalDatetime(tt.literal)
		functionResult := testEvalDatetime(tt.function)

		litInspect := literalResult.Inspect()
		funcInspect := functionResult.Inspect()

		if litInspect != funcInspect {
			t.Errorf("Datetime literal and function results don't match.\nLiteral:  %s\nFunction: %s",
				litInspect, funcInspect)
		}
	}
}

func TestDatetimeLiteralInExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// In variables
		{`let christmas = @2024-12-25; christmas.day`, `25`},
		{`let birthday = @2024-03-15T14:30:00; birthday.hour`, `14`},
		// In conditionals
		{`if(@2024-12-25 < @2024-12-26) { "true" } else { "false" }`, `true`},
		// In arrays
		{`[@2024-01-01, @2024-06-15, @2024-12-31][1].month`, `6`},
		// In function calls
		{`let getYear = fn(dt) { dt.year }; getYear(@2024-12-25)`, `2024`},
	}

	for _, tt := range tests {
		evaluated := testEvalDatetime(tt.input)
		testExpectedDatetime(t, tt.input, evaluated, tt.expected)
	}
}

func TestDatetimeLiteralEdgeCases(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Leap year
		{`@2024-02-29`, `{__type: "datetime", year: 2024, month: 2, day: 29, hour: 0, minute: 0, second: 0, unix: 1709164800, weekday: "Thursday", iso: "2024-02-29T00:00:00Z"}`},
		// New Year's Day
		{`@2024-01-01T00:00:00`, `{__type: "datetime", year: 2024, month: 1, day: 1, hour: 0, minute: 0, second: 0, unix: 1704067200, weekday: "Monday", iso: "2024-01-01T00:00:00Z"}`},
		// New Year's Eve
		{`@2024-12-31T23:59:59`, `{__type: "datetime", year: 2024, month: 12, day: 31, hour: 23, minute: 59, second: 59, unix: 1735689599, weekday: "Tuesday", iso: "2024-12-31T23:59:59Z"}`},
	}

	for _, tt := range tests {
		evaluated := testEvalDatetime(tt.input)
		testExpectedDatetime(t, tt.input, evaluated, tt.expected)
	}
}

func TestDatetimeLiteralErrors(t *testing.T) {
	tests := []struct {
		input         string
		expectedError string
	}{
		{`@2024-13-01`, "invalid datetime literal"}, // Invalid month
		{`@2024-02-30`, "invalid datetime literal"}, // Invalid day for February
		{`@2024-04-31`, "invalid datetime literal"}, // Invalid day for April
		{`@not-a-date`, "identifier not found"},     // Lexed as @ followed by identifier
		{`@2024`, "invalid duration literal"},       // Incomplete - parsed as duration, not datetime
	}

	for _, tt := range tests {
		evaluated := testEvalDatetime(tt.input)
		errObj, ok := evaluated.(*evaluator.Error)
		if !ok {
			t.Errorf("Expected error for input %q, got %T (%+v)", tt.input, evaluated, evaluated)
			continue
		}
		if len(errObj.Message) < len(tt.expectedError) || errObj.Message[:len(tt.expectedError)] != tt.expectedError {
			t.Errorf("Wrong error message. Expected to start with %q, got %q",
				tt.expectedError, errObj.Message)
		}
	}
}
