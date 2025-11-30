package main

import (
	"testing"

	"github.com/sambeau/parsley/pkg/evaluator"
)

func TestRegexLiterals(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`let rx = /\d+/; rx.pattern`, `"\d+"`},
		{`let rx = /\d+/; rx.flags`, `""`},
		{`let rx = /hello/i; rx.flags`, `"i"`},
		{`let rx = /test/gim; rx.flags`, `"gim"`},
		{`let rx = /[a-z]+/; rx.__type`, `"regex"`},
	}

	for _, tt := range tests {
		evaluated := testEvalHelper(tt.input)
		testExpectedObject(t, tt.input, evaluated, tt.expected)
	}
}

func TestRegexMatch(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello 123" ~ /\d+/`, `["123"]`},
		{`"no numbers" ~ /\d+/`, `null`},
		{`"user@example.com" ~ /(\w+)@([\w.]+)/`, `["user@example.com", "user", "example.com"]`},
		{`"2024-01-15" ~ /(\d+)-(\d+)-(\d+)/`, `["2024-01-15", "2024", "01", "15"]`},
		{`"Test" ~ /test/`, `null`},      // case-sensitive
		{`"Test" ~ /test/i`, `["Test"]`}, // case-insensitive with flag
	}

	for _, tt := range tests {
		evaluated := testEvalHelper(tt.input)
		testExpectedObject(t, tt.input, evaluated, tt.expected)
	}
}

func TestRegexNotMatch(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello world" !~ /\d+/`, `true`},
		{`"hello 123" !~ /\d+/`, `false`},
		{`"test" !~ /TEST/`, `true`},
		{`"test" !~ /TEST/i`, `false`},
	}

	for _, tt := range tests {
		evaluated := testEvalHelper(tt.input)
		testExpectedObject(t, tt.input, evaluated, tt.expected)
	}
}

func TestRegexConditional(t *testing.T) {
	input := `
		let email = "john@example.com";
		let match = email ~ /(\w+)@([\w.]+)/;
		if (match) {
			match[1]
		} else {
			"no match"
		}
	`
	evaluated := testEvalHelper(input)
	testExpectedObject(t, input, evaluated, `"john"`)
}

func TestRegexDestructuring(t *testing.T) {
	input := `
		let email = "jane@test.org";
		let [full, name, domain] = email ~ /(\w+)@([\w.]+)/;
		name + ":" + domain
	`
	evaluated := testEvalHelper(input)
	testExpectedObject(t, input, evaluated, `"jane:test.org"`)
}

func TestRegexBuiltin(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`let rx = regex("\\d+"); rx.pattern`, `"\d+"`},
		{`let rx = regex("test", "i"); rx.flags`, `"i"`},
		{`let rx = regex("\\w+"); "hello" ~ rx`, `["hello"]`},
	}

	for _, tt := range tests {
		evaluated := testEvalHelper(tt.input)
		testExpectedObject(t, tt.input, evaluated, tt.expected)
	}
}

func TestReplaceFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`replace("hello world", "world", "Go")`, `"hello Go"`},
		{`replace("one1two2three", /\d/, "X")`, `"oneXtwoXthree"`},
		{`replace("Test", /test/i, "SUCCESS")`, `"SUCCESS"`},
		{`replace("a,b,c", ",", ";")`, `"a;b;c"`},
	}

	for _, tt := range tests {
		evaluated := testEvalHelper(tt.input)
		testExpectedObject(t, tt.input, evaluated, tt.expected)
	}
}

func TestSplitFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`split("a,b,c", ",")`, `["a", "b", "c"]`},
		{`split("one1two2three", /\d/)`, `["one", "two", "three"]`},
		{`split("hello world", " ")`, `["hello", "world"]`},
		{`len(split("a:b:c:d", ":"))`, `4`},
	}

	for _, tt := range tests {
		evaluated := testEvalHelper(tt.input)
		testExpectedObject(t, tt.input, evaluated, tt.expected)
	}
}

func TestRegexFlags(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Case-insensitive flag
		{`"Hello" ~ /hello/i`, `["Hello"]`},
		{`"Hello" ~ /hello/`, `null`},

		// Multi-line flag - test simpler pattern
		{`"test" ~ /test/m`, `["test"]`},
	}

	for _, tt := range tests {
		evaluated := testEvalHelper(tt.input)
		testExpectedObject(t, tt.input, evaluated, tt.expected)
	}
}

func TestRegexErrors(t *testing.T) {
	tests := []struct {
		input       string
		expectedErr string
	}{
		{`/[/`, "invalid regex pattern"},
		{`regex("(unclosed")`, "invalid regex pattern"},
		{`123 ~ /\d+/`, "left operand of ~ must be a string"},
		{`"test" ~ "not regex"`, "right operand of ~ must be a regex"},
	}

	for _, tt := range tests {
		evaluated := testEvalHelper(tt.input)
		errObj, ok := evaluated.(*evaluator.Error)
		if !ok {
			t.Errorf("Expected error for input '%s', got %T", tt.input, evaluated)
			continue
		}
		if errObj.Message == "" {
			t.Errorf("Expected error message for input '%s', got empty", tt.input)
		}
	}
}

func TestRegexComplexPatterns(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Email validation
		{`"valid@email.com" ~ /^[\w.+-]+@[\w.-]+\.[a-zA-Z]{2,}$/`, `["valid@email.com"]`},
		{`"invalid@" ~ /^[\w.+-]+@[\w.-]+\.[a-zA-Z]{2,}$/`, `null`},

		// URL parsing
		{`"https://example.com" ~ /^(https?):\/\/([^\/]+)/`, `["https://example.com", "https", "example.com"]`},

		// Phone number
		{`"(123) 456-7890" ~ /\((\d{3})\) (\d{3})-(\d{4})/`, `["(123) 456-7890", "123", "456", "7890"]`},
	}

	for _, tt := range tests {
		evaluated := testEvalHelper(tt.input)
		testExpectedObject(t, tt.input, evaluated, tt.expected)
	}
}

// Helper to test expected output
func testExpectedObject(t *testing.T, input string, obj evaluator.Object, expected string) {
	if obj == nil {
		t.Errorf("For input '%s': got nil object", input)
		return
	}

	var actual string
	switch v := obj.(type) {
	case *evaluator.Integer:
		actual = v.Inspect()
	case *evaluator.String:
		actual = `"` + v.Inspect() + `"`
	case *evaluator.Boolean:
		actual = v.Inspect()
	case *evaluator.Null:
		actual = "null"
	case *evaluator.Array:
		actual = "["
		for i, elem := range v.Elements {
			if i > 0 {
				actual += ", "
			}
			if str, ok := elem.(*evaluator.String); ok {
				actual += `"` + str.Inspect() + `"`
			} else {
				actual += elem.Inspect()
			}
		}
		actual += "]"
	case *evaluator.Error:
		t.Errorf("For input '%s': got error: %s", input, v.Message)
		return
	default:
		actual = obj.Inspect()
	}

	if actual != expected {
		t.Errorf("For input '%s': expected %s, got %s", input, expected, actual)
	}
}
