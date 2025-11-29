package main

import (
	"strings"
	"testing"

	"github.com/sambeau/parsley/pkg/evaluator"
	"github.com/sambeau/parsley/pkg/lexer"
	"github.com/sambeau/parsley/pkg/parser"
)

// ============================================================================
// String Method Tests
// ============================================================================

func TestStringMethods(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// upper()
		{`"hello".toUpper()`, "HELLO"},
		{`"Hello World".toUpper()`, "HELLO WORLD"},
		{`"ALREADY UPPER".toUpper()`, "ALREADY UPPER"},

		// lower()
		{`"HELLO".toLower()`, "hello"},
		{`"Hello World".toLower()`, "hello world"},
		{`"already lower".toLower()`, "already lower"},

		// trim()
		{`"  hello  ".trim()`, "hello"},
		{`"  hello world  ".trim()`, "hello world"},
		{`"no trim needed".trim()`, "no trim needed"},

		// split()
		{`"a,b,c".split(",")`, []string{"a", "b", "c"}},
		{`"hello world".split(" ")`, []string{"hello", "world"}},
		{`"one".split(",")`, []string{"one"}},

		// replace()
		{`"hello world".replace("world", "there")`, "hello there"},
		{`"aaa".replace("a", "b")`, "bbb"},

		// length()
		{`"hello".length()`, int64(5)},
		{`"".length()`, int64(0)},
		{`"日本語".length()`, int64(3)}, // Unicode rune count
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()
			env := evaluator.NewEnvironment()
			result := evaluator.Eval(program, env)

			switch expected := tt.expected.(type) {
			case string:
				str, ok := result.(*evaluator.String)
				if !ok {
					t.Fatalf("expected String, got %T (%+v)", result, result)
				}
				if str.Value != expected {
					t.Errorf("expected %q, got %q", expected, str.Value)
				}
			case int64:
				num, ok := result.(*evaluator.Integer)
				if !ok {
					t.Fatalf("expected Integer, got %T (%+v)", result, result)
				}
				if num.Value != expected {
					t.Errorf("expected %d, got %d", expected, num.Value)
				}
			case []string:
				arr, ok := result.(*evaluator.Array)
				if !ok {
					t.Fatalf("expected Array, got %T (%+v)", result, result)
				}
				if len(arr.Elements) != len(expected) {
					t.Fatalf("expected %d elements, got %d", len(expected), len(arr.Elements))
				}
				for i, exp := range expected {
					str, ok := arr.Elements[i].(*evaluator.String)
					if !ok {
						t.Fatalf("expected element %d to be String, got %T", i, arr.Elements[i])
					}
					if str.Value != exp {
						t.Errorf("element %d: expected %q, got %q", i, exp, str.Value)
					}
				}
			}
		})
	}
}

// ============================================================================
// Array Method Tests
// ============================================================================

func TestArrayMethods(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// length()
		{`[1, 2, 3].length()`, int64(3)},
		{`[].length()`, int64(0)},

		// reverse()
		{`[1, 2, 3].reverse()`, []int64{3, 2, 1}},
		{`["a", "b", "c"].reverse()`, []string{"c", "b", "a"}},

		// sort()
		{`[3, 1, 2].sort()`, []int64{1, 2, 3}},
		{`["banana", "apple", "cherry"].sort()`, []string{"apple", "banana", "cherry"}},

		// map()
		{`[1, 2, 3].map(fn(x) { x * 2 })`, []int64{2, 4, 6}},

		// filter()
		{`[1, 2, 3, 4, 5].filter(fn(x) { x > 2 })`, []int64{3, 4, 5}},
		{`["a", "bb", "ccc"].filter(fn(s) { len(s) > 1 })`, []string{"bb", "ccc"}},
		{`[1, 2, 3].filter(fn(x) { x > 10 })`, []int64{}},

		// format()
		{`["apple", "banana", "cherry"].format()`, "apple, banana, and cherry"},
		{`["a", "b"].format("or")`, "a or b"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()
			env := evaluator.NewEnvironment()
			result := evaluator.Eval(program, env)

			switch expected := tt.expected.(type) {
			case int64:
				num, ok := result.(*evaluator.Integer)
				if !ok {
					t.Fatalf("expected Integer, got %T (%+v)", result, result)
				}
				if num.Value != expected {
					t.Errorf("expected %d, got %d", expected, num.Value)
				}
			case []int64:
				arr, ok := result.(*evaluator.Array)
				if !ok {
					t.Fatalf("expected Array, got %T (%+v)", result, result)
				}
				if len(arr.Elements) != len(expected) {
					t.Fatalf("expected %d elements, got %d", len(expected), len(arr.Elements))
				}
				for i, exp := range expected {
					num, ok := arr.Elements[i].(*evaluator.Integer)
					if !ok {
						t.Fatalf("expected element %d to be Integer, got %T", i, arr.Elements[i])
					}
					if num.Value != exp {
						t.Errorf("element %d: expected %d, got %d", i, exp, num.Value)
					}
				}
			case []string:
				arr, ok := result.(*evaluator.Array)
				if !ok {
					t.Fatalf("expected Array, got %T (%+v)", result, result)
				}
				if len(arr.Elements) != len(expected) {
					t.Fatalf("expected %d elements, got %d", len(expected), len(arr.Elements))
				}
				for i, exp := range expected {
					str, ok := arr.Elements[i].(*evaluator.String)
					if !ok {
						t.Fatalf("expected element %d to be String, got %T", i, arr.Elements[i])
					}
					if str.Value != exp {
						t.Errorf("element %d: expected %q, got %q", i, exp, str.Value)
					}
				}
			case string:
				str, ok := result.(*evaluator.String)
				if !ok {
					t.Fatalf("expected String, got %T (%+v)", result, result)
				}
				if str.Value != expected {
					t.Errorf("expected %q, got %q", expected, str.Value)
				}
			}
		})
	}
}

// ============================================================================
// Dictionary Method Tests
// ============================================================================

func TestDictionaryMethods(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
		check    string // "bool" for boolean, "contains" for array containment
	}{
		// keys()
		{`let d = {a: 1, b: 2}; d.keys()`, []string{"a", "b"}, "contains"},

		// values()
		{`let d = {a: 1, b: 2}; d.values()`, []int64{1, 2}, "contains"},

		// has()
		{`let d = {a: 1, b: 2}; d.has("a")`, true, "bool"},
		{`let d = {a: 1, b: 2}; d.has("c")`, false, "bool"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()
			env := evaluator.NewEnvironment()
			result := evaluator.Eval(program, env)

			switch tt.check {
			case "bool":
				expected := tt.expected.(bool)
				b, ok := result.(*evaluator.Boolean)
				if !ok {
					t.Fatalf("expected Boolean, got %T (%+v)", result, result)
				}
				if b.Value != expected {
					t.Errorf("expected %v, got %v", expected, b.Value)
				}
			case "contains":
				arr, ok := result.(*evaluator.Array)
				if !ok {
					t.Fatalf("expected Array, got %T (%+v)", result, result)
				}
				// Check that all expected values are present (order may vary)
				switch expected := tt.expected.(type) {
				case []string:
					if len(arr.Elements) != len(expected) {
						t.Fatalf("expected %d elements, got %d", len(expected), len(arr.Elements))
					}
					for _, exp := range expected {
						found := false
						for _, elem := range arr.Elements {
							if s, ok := elem.(*evaluator.String); ok && s.Value == exp {
								found = true
								break
							}
						}
						if !found {
							t.Errorf("expected to find %q in array", exp)
						}
					}
				case []int64:
					if len(arr.Elements) != len(expected) {
						t.Fatalf("expected %d elements, got %d", len(expected), len(arr.Elements))
					}
					for _, exp := range expected {
						found := false
						for _, elem := range arr.Elements {
							if n, ok := elem.(*evaluator.Integer); ok && n.Value == exp {
								found = true
								break
							}
						}
						if !found {
							t.Errorf("expected to find %d in array", exp)
						}
					}
				}
			}
		})
	}
}

// ============================================================================
// Path Method Tests
// ============================================================================

func TestPathMethods(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		// isAbsolute()
		{`let p = path("/usr/local"); p.isAbsolute()`, true},
		{`let p = path("relative/path"); p.isAbsolute()`, false},

		// isRelative()
		{`let p = path("/usr/local"); p.isRelative()`, false},
		{`let p = path("relative/path"); p.isRelative()`, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()
			env := evaluator.NewEnvironment()
			result := evaluator.Eval(program, env)

			b, ok := result.(*evaluator.Boolean)
			if !ok {
				t.Fatalf("expected Boolean, got %T (%+v)", result, result)
			}
			if b.Value != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, b.Value)
			}
		})
	}
}

// ============================================================================
// URL Method Tests
// ============================================================================

func TestUrlMethods(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// origin()
		{`let u = url("https://example.com/api"); u.origin()`, "https://example.com"},
		{`let u = url("https://example.com:8080/api"); u.origin()`, "https://example.com:8080"},

		// pathname()
		{`let u = url("https://example.com/api/users"); u.pathname()`, "/api/users"},
		{`let u = url("https://example.com"); u.pathname()`, "/"},

		// href()
		{`let u = url("https://example.com/api"); u.href()`, "https://example.com/api"},
		{`let u = url("https://example.com:8080/api?limit=10#top"); u.href()`, "https://example.com:8080/api?limit=10#top"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()
			env := evaluator.NewEnvironment()
			result := evaluator.Eval(program, env)

			s, ok := result.(*evaluator.String)
			if !ok {
				t.Fatalf("expected String, got %T (%+v)", result, result)
			}
			if s.Value != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, s.Value)
			}
		})
	}
}

// ============================================================================
// Number Method Tests
// ============================================================================

func TestNumberMethods(t *testing.T) {
	tests := []struct {
		input    string
		contains string
	}{
		// Integer.format()
		{`1234567.format()`, "1,234,567"},
		{`1234567.format("de-DE")`, "1.234.567"},

		// Float.format()
		{`1234.56.format()`, "1,234.56"},
		{`1234.56.format("de-DE")`, "1.234,56"},

		// Integer.currency()
		{`99.currency("USD")`, "$"},
		{`99.currency("EUR", "de-DE")`, "€"},

		// Float.currency()
		{`99.99.currency("USD")`, "$"},
		{`99.99.currency("GBP", "en-GB")`, "£"},

		// Integer.percent()
		{`0.percent()`, "%"},
		{`1.percent()`, "100"},

		// Float.percent()
		{`0.125.percent()`, "12"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()
			env := evaluator.NewEnvironment()
			result := evaluator.Eval(program, env)

			str, ok := result.(*evaluator.String)
			if !ok {
				t.Fatalf("expected String, got %T (%+v)", result, result)
			}
			if !strings.Contains(str.Value, tt.contains) {
				t.Errorf("expected to contain %q, got %q", tt.contains, str.Value)
			}
		})
	}
}

// ============================================================================
// Datetime Method Tests
// ============================================================================

func TestDatetimeMethods(t *testing.T) {
	tests := []struct {
		input    string
		contains string
	}{
		// format()
		{`let d = time({year: 2024, month: 12, day: 25}); d.format()`, "December 25, 2024"},
		{`let d = time({year: 2024, month: 12, day: 25}); d.format("short")`, "12/25/24"},
		{`let d = time({year: 2024, month: 12, day: 25}); d.format("long", "de-DE")`, "25. Dezember 2024"},
		{`let d = time({year: 2024, month: 12, day: 25}); d.format("long", "fr-FR")`, "25 décembre 2024"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()
			env := evaluator.NewEnvironment()
			result := evaluator.Eval(program, env)

			str, ok := result.(*evaluator.String)
			if !ok {
				t.Fatalf("expected String, got %T (%+v)", result, result)
			}
			if !strings.Contains(str.Value, tt.contains) {
				t.Errorf("expected to contain %q, got %q", tt.contains, str.Value)
			}
		})
	}
}

// ============================================================================
// Duration Method Tests
// ============================================================================

func TestDurationMethods(t *testing.T) {
	tests := []struct {
		input    string
		contains string
	}{
		// format()
		{`@1d.format()`, "tomorrow"},
		{`@-1d.format()`, "yesterday"},
		{`@1d.format("de-DE")`, "morgen"},
		{`@-1d.format("de-DE")`, "gestern"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()
			env := evaluator.NewEnvironment()
			result := evaluator.Eval(program, env)

			str, ok := result.(*evaluator.String)
			if !ok {
				t.Fatalf("expected String, got %T (%+v)", result, result)
			}
			if !strings.Contains(str.Value, tt.contains) {
				t.Errorf("expected to contain %q, got %q", tt.contains, str.Value)
			}
		})
	}
}

// ============================================================================
// Method Chaining Tests
// ============================================================================

func TestMethodChaining(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// String chaining
		{`"  hello world  ".trim().toUpper()`, "HELLO WORLD"},
		{`"HELLO".toLower().toUpper()`, "HELLO"},

		// String → Array chaining
		{`"a,b,c".split(",").length()`, int64(3)},
		{`"c,a,b".split(",").sort()`, []string{"a", "b", "c"}},
		{`"c,b,a".split(",").reverse()`, []string{"a", "b", "c"}},

		// Array → String chaining
		{`["hello", "world"].format().toUpper()`, "HELLO AND WORLD"},

		// Array chaining
		{`[3, 1, 2].sort().reverse()`, []int64{3, 2, 1}},
		{`[1, 2, 3].map(fn(x) { x * 2 }).reverse()`, []int64{6, 4, 2}},

		// Number → String chaining
		// 1234567 formats to "1,234,567" which splits to ["1", "234", "567"] (3 parts)
		{`1234567.format().split(",").length()`, int64(3)},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()
			env := evaluator.NewEnvironment()
			result := evaluator.Eval(program, env)

			switch expected := tt.expected.(type) {
			case string:
				str, ok := result.(*evaluator.String)
				if !ok {
					t.Fatalf("expected String, got %T (%+v)", result, result)
				}
				if str.Value != expected {
					t.Errorf("expected %q, got %q", expected, str.Value)
				}
			case int64:
				num, ok := result.(*evaluator.Integer)
				if !ok {
					t.Fatalf("expected Integer, got %T (%+v)", result, result)
				}
				if num.Value != expected {
					t.Errorf("expected %d, got %d", expected, num.Value)
				}
			case []int64:
				arr, ok := result.(*evaluator.Array)
				if !ok {
					t.Fatalf("expected Array, got %T (%+v)", result, result)
				}
				if len(arr.Elements) != len(expected) {
					t.Fatalf("expected %d elements, got %d", len(expected), len(arr.Elements))
				}
				for i, exp := range expected {
					num, ok := arr.Elements[i].(*evaluator.Integer)
					if !ok {
						t.Fatalf("expected element %d to be Integer, got %T", i, arr.Elements[i])
					}
					if num.Value != exp {
						t.Errorf("element %d: expected %d, got %d", i, exp, num.Value)
					}
				}
			case []string:
				arr, ok := result.(*evaluator.Array)
				if !ok {
					t.Fatalf("expected Array, got %T (%+v)", result, result)
				}
				if len(arr.Elements) != len(expected) {
					t.Fatalf("expected %d elements, got %d", len(expected), len(arr.Elements))
				}
				for i, exp := range expected {
					str, ok := arr.Elements[i].(*evaluator.String)
					if !ok {
						t.Fatalf("expected element %d to be String, got %T", i, arr.Elements[i])
					}
					if str.Value != exp {
						t.Errorf("element %d: expected %q, got %q", i, exp, str.Value)
					}
				}
			}
		})
	}
}

// ============================================================================
// Null Propagation Tests
// ============================================================================

func TestNullPropagation(t *testing.T) {
	tests := []struct {
		input string
	}{
		// Method calls on null return null (using missing dictionary key to get null)
		{`let d = {a: 1}; d.b.toUpper()`},
		{`let d = {a: 1}; d.b.toLower()`},
		{`let d = {a: 1}; d.b.split(",")`},
		{`let d = {a: 1}; d.b.length()`},
		{`let d = {a: 1}; d.b.format()`},

		// Chained null propagation
		{`let d = {a: 1}; d.b.toUpper().toLower()`},
		{`let d = {a: 1}; d.b.split(",").reverse()`},

		// Property access on null returns null
		{`let d = {a: 1}; let x = d.b; x.foo`},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()
			env := evaluator.NewEnvironment()
			result := evaluator.Eval(program, env)

			if result != evaluator.NULL {
				t.Errorf("expected NULL, got %T (%+v)", result, result)
			}
		})
	}
}

// ============================================================================
// Method Error Tests
// ============================================================================

func TestMethodErrors(t *testing.T) {
	tests := []struct {
		input       string
		errContains string
	}{
		// Wrong argument count
		{`"hello".toUpper("arg")`, "wrong number of arguments"},
		{`"hello".split()`, "wrong number of arguments"},
		{`"hello".replace("a")`, "wrong number of arguments"},

		// Wrong argument type
		{`"hello".split(123)`, "must be a string"},
		{`"hello".replace(1, 2)`, "must be a string"},

		// Unknown method
		{`"hello".unknown()`, "unknown method"},
		{`[1, 2, 3].unknown()`, "unknown method"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()
			env := evaluator.NewEnvironment()
			result := evaluator.Eval(program, env)

			err, ok := result.(*evaluator.Error)
			if !ok {
				t.Fatalf("expected Error, got %T (%+v)", result, result)
			}
			if !strings.Contains(err.Message, tt.errContains) {
				t.Errorf("expected error to contain %q, got %q", tt.errContains, err.Message)
			}
		})
	}
}
