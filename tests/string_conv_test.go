package main

import (
	"strings"
	"testing"

	"github.com/sambeau/parsley/pkg/evaluator"
	"github.com/sambeau/parsley/pkg/lexer"
	"github.com/sambeau/parsley/pkg/parser"
)

func evalStringConv(input string) string {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := evaluator.NewEnvironment()
	result := evaluator.Eval(program, env)
	return evaluator.ObjectToPrintString(result)
}

func evalStringConvInspect(input string) string {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := evaluator.NewEnvironment()
	result := evaluator.Eval(program, env)
	if result != nil {
		return result.Inspect()
	}
	return ""
}

func TestDurationString(t *testing.T) {
	result := evalStringConv("@1d")
	expected := "1 day"
	if result != expected {
		t.Errorf("got %q, want %q", result, expected)
	}
}

func TestDurationMultiUnit(t *testing.T) {
	result := evalStringConv("@2d12h")
	expected := "2 days 12 hours"
	if result != expected {
		t.Errorf("got %q, want %q", result, expected)
	}
}

func TestDurationYearsMonths(t *testing.T) {
	result := evalStringConv("@1y6mo")
	expected := "1 year 6 months"
	if result != expected {
		t.Errorf("got %q, want %q", result, expected)
	}
}

func TestRegexString(t *testing.T) {
	result := evalStringConv("/hello.*/i")
	expected := "/hello.*/i"
	if result != expected {
		t.Errorf("got %q, want %q", result, expected)
	}
}

func TestRegexFormatPattern(t *testing.T) {
	result := evalStringConv(`/hello.*/i.format("pattern")`)
	expected := "hello.*"
	if result != expected {
		t.Errorf("got %q, want %q", result, expected)
	}
}

func TestRegexFormatVerbose(t *testing.T) {
	result := evalStringConv(`/hello.*/i.format("verbose")`)
	expected := "pattern: hello.*, flags: i"
	if result != expected {
		t.Errorf("got %q, want %q", result, expected)
	}
}

func TestRegexTest(t *testing.T) {
	result := evalStringConv(`/hello/i.test("Hello World")`)
	if result != "true" {
		t.Errorf("got %q, want true", result)
	}
}

func TestRegexTestNoMatch(t *testing.T) {
	result := evalStringConv(`/goodbye/i.test("Hello World")`)
	if result != "false" {
		t.Errorf("got %q, want false", result)
	}
}

func TestPathString(t *testing.T) {
	result := evalStringConv("@./src/main.go")
	expected := "./src/main.go"
	if result != expected {
		t.Errorf("got %q, want %q", result, expected)
	}
}

func TestURLString(t *testing.T) {
	result := evalStringConv("@https://example.com:8080/api/users")
	expected := "https://example.com:8080/api/users"
	if result != expected {
		t.Errorf("got %q, want %q", result, expected)
	}
}

func TestPathToDict(t *testing.T) {
	// toDict returns the dictionary itself, which when inspected shows the dict format
	result := evalStringConvInspect(`let p = @./src/main.go
p.toDict()`)
	if !strings.Contains(result, "__type") || !strings.Contains(result, "path") {
		t.Errorf("toDict should return dict with __type: path, got %q", result)
	}
}

func TestURLToDict(t *testing.T) {
	result := evalStringConvInspect(`let u = @https://example.com
u.toDict()`)
	if !strings.Contains(result, "__type") || !strings.Contains(result, "url") {
		t.Errorf("toDict should return dict with __type: url, got %q", result)
	}
}

func TestDatetimeToDict(t *testing.T) {
	result := evalStringConvInspect(`let d = @2024-12-25
d.toDict()`)
	if !strings.Contains(result, "__type") || !strings.Contains(result, "datetime") {
		t.Errorf("toDict should return dict with __type: datetime, got %q", result)
	}
}

func TestDurationToDict(t *testing.T) {
	result := evalStringConvInspect(`let d = @1d
d.toDict()`)
	if !strings.Contains(result, "__type") || !strings.Contains(result, "duration") {
		t.Errorf("toDict should return dict with __type: duration, got %q", result)
	}
}

func TestRegexToDict(t *testing.T) {
	result := evalStringConvInspect(`let r = /hello.*/i
r.toDict()`)
	if !strings.Contains(result, "__type") || !strings.Contains(result, "regex") {
		t.Errorf("toDict should return dict with __type: regex, got %q", result)
	}
}

func TestReprDuration(t *testing.T) {
	result := evalStringConv("repr(@1d)")
	if !strings.Contains(result, "__type") || !strings.Contains(result, "duration") {
		t.Errorf("repr should return dict string with __type: duration, got %q", result)
	}
}

func TestReprPath(t *testing.T) {
	result := evalStringConv("repr(@./src/main.go)")
	if !strings.Contains(result, "__type") || !strings.Contains(result, "path") {
		t.Errorf("repr should return dict string with __type: path, got %q", result)
	}
}

func TestReprURL(t *testing.T) {
	result := evalStringConv("repr(@https://example.com)")
	if !strings.Contains(result, "__type") || !strings.Contains(result, "url") {
		t.Errorf("repr should return dict string with __type: url, got %q", result)
	}
}

func TestReprDatetime(t *testing.T) {
	result := evalStringConv("repr(@2024-12-25)")
	if !strings.Contains(result, "__type") || !strings.Contains(result, "datetime") {
		t.Errorf("repr should return dict string with __type: datetime, got %q", result)
	}
}

func TestReprRegex(t *testing.T) {
	result := evalStringConv("repr(/test/i)")
	if !strings.Contains(result, "__type") || !strings.Contains(result, "regex") {
		t.Errorf("repr should return dict string with __type: regex, got %q", result)
	}
}

// Array join() tests
func TestArrayJoinNoSeparator(t *testing.T) {
	result := evalStringConv(`["a", "b", "c"].join()`)
	expected := "abc"
	if result != expected {
		t.Errorf("got %q, want %q", result, expected)
	}
}

func TestArrayJoinWithSeparator(t *testing.T) {
	result := evalStringConv(`["a", "b", "c"].join("-")`)
	expected := "a-b-c"
	if result != expected {
		t.Errorf("got %q, want %q", result, expected)
	}
}

func TestArrayJoinWithComma(t *testing.T) {
	result := evalStringConv(`["apple", "banana", "cherry"].join(", ")`)
	expected := "apple, banana, cherry"
	if result != expected {
		t.Errorf("got %q, want %q", result, expected)
	}
}

func TestArrayJoinPathComponents(t *testing.T) {
	result := evalStringConv(`let p = path("/usr/local/bin"); p.components.join("/")`)
	expected := "/usr/local/bin"
	if result != expected {
		t.Errorf("got %q, want %q", result, expected)
	}
}
