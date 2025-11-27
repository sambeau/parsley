package main

import (
	"testing"

	"github.com/sambeau/parsley/pkg/evaluator"
	"github.com/sambeau/parsley/pkg/lexer"
	"github.com/sambeau/parsley/pkg/parser"
)

func testEvalOperator(input string) evaluator.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := evaluator.NewEnvironment()
	return evaluator.Eval(program, env)
}

// Test path + operator
func TestPathPlusOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`let p = @/usr/local + "bin"; toString(p)`, "/usr/local/bin"},
		{`let p = @/var + "log"; toString(p)`, "/var/log"},
		{`let p = @~/Documents + "project"; toString(p)`, "~/Documents/project"},
		{`let p = @./config + "app.json"; toString(p)`, "./config/app.json"},
		{`let p = (@/usr + "local") + "bin"; toString(p)`, "/usr/local/bin"},
	}

	for _, tt := range tests {
		result := testEvalOperator(tt.input)
		str, ok := result.(*evaluator.String)
		if !ok {
			t.Errorf("For input '%s': expected String, got %T (%+v)", tt.input, result, result)
			continue
		}
		if str.Value != tt.expected {
			t.Errorf("For input '%s': expected '%s', got '%s'", tt.input, tt.expected, str.Value)
		}
	}
}

// Test path / operator
func TestPathSlashOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`let p = @/usr / "local"; toString(p)`, "/usr/local"},
		{`let p = @/var / "log" / "system.log"; toString(p)`, "/var/log/system.log"},
		{`let p = @~/ / "Documents"; toString(p)`, "~/Documents"},
	}

	for _, tt := range tests {
		result := testEvalOperator(tt.input)
		str, ok := result.(*evaluator.String)
		if !ok {
			t.Errorf("For input '%s': expected String, got %T (%+v)", tt.input, result, result)
			continue
		}
		if str.Value != tt.expected {
			t.Errorf("For input '%s': expected '%s', got '%s'", tt.input, tt.expected, str.Value)
		}
	}
}

// Test path comparison operators
func TestPathComparisonOperators(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`@/usr/local/bin == @/usr/local/bin`, true},
		{`@/usr/local/bin != @/usr/bin`, true},
		{`@/usr/local/bin == @/usr/bin`, false},
		{`@./config == @./config`, true},
		{`@~/docs != @~/documents`, true},
		{`let p1 = @/usr/local; let p2 = @/usr/local; p1 == p2`, true},
		{`let p1 = @/usr/local; let p2 = @/usr/bin; p1 != p2`, true},
	}

	for _, tt := range tests {
		result := testEvalOperator(tt.input)
		boolean, ok := result.(*evaluator.Boolean)
		if !ok {
			t.Errorf("For input '%s': expected Boolean, got %T (%+v)", tt.input, result, result)
			continue
		}
		if boolean.Value != tt.expected {
			t.Errorf("For input '%s': expected %v, got %v", tt.input, tt.expected, boolean.Value)
		}
	}
}

// Test path + with computed properties
func TestPathOperatorWithProperties(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`let p = @/usr + "local" + "bin"; p.basename`, "bin"},
		{`let p = @/var / "log" / "system.log"; p.extension`, "log"},
		{`let p = (@/usr + "local") / "file.txt"; p.stem`, "file"},
		{`let p = @/usr + "local"; toString(p.dirname)`, "/usr"},
	}

	for _, tt := range tests {
		result := testEvalOperator(tt.input)
		str, ok := result.(*evaluator.String)
		if !ok {
			t.Errorf("For input '%s': expected String, got %T (%+v)", tt.input, result, result)
			continue
		}
		if str.Value != tt.expected {
			t.Errorf("For input '%s': expected '%s', got '%s'", tt.input, tt.expected, str.Value)
		}
	}
}

// Test URL + operator
func TestUrlPlusOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`let u = @https://example.com + "api"; toString(u)`, "https://example.com/api"},
		{`let u = @https://api.example.com + "users"; toString(u)`, "https://api.example.com/users"},
		{`let u = (@https://api.github.com + "repos") + "sambeau"; toString(u)`, "https://api.github.com/repos/sambeau"},
		{`let u = @http://localhost:8080 + "api" + "v1"; toString(u)`, "http://localhost:8080/api/v1"},
	}

	for _, tt := range tests {
		result := testEvalOperator(tt.input)
		str, ok := result.(*evaluator.String)
		if !ok {
			t.Errorf("For input '%s': expected String, got %T (%+v)", tt.input, result, result)
			continue
		}
		if str.Value != tt.expected {
			t.Errorf("For input '%s': expected '%s', got '%s'", tt.input, tt.expected, str.Value)
		}
	}
}

// Test URL comparison operators
func TestUrlComparisonOperators(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`@https://example.com/api == @https://example.com/api`, true},
		{`@https://example.com/api != @https://example.com/docs`, true},
		{`@https://example.com == @https://example.com`, true},
		{`@https://example.com:8080 != @https://example.com`, true},
		{`let u1 = @https://api.example.com; let u2 = @https://api.example.com; u1 == u2`, true},
		{`let u1 = @https://api.example.com + "users"; let u2 = @https://api.example.com + "posts"; u1 != u2`, true},
	}

	for _, tt := range tests {
		result := testEvalOperator(tt.input)
		boolean, ok := result.(*evaluator.Boolean)
		if !ok {
			t.Errorf("For input '%s': expected Boolean, got %T (%+v)", tt.input, result, result)
			continue
		}
		if boolean.Value != tt.expected {
			t.Errorf("For input '%s': expected %v, got %v", tt.input, tt.expected, boolean.Value)
		}
	}
}

// Test URL + with computed properties
func TestUrlOperatorWithProperties(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`let u = @https://example.com + "api" + "v1"; u.pathname`, "api/v1"},
		{`let u = @https://api.github.com + "repos"; u.host`, "api.github.com"},
		{`let u = @https://example.com:8080 + "api"; u.origin`, "https://example.com:8080"},
		{`let u = (@https://example.com + "api") + "users"; u.pathname`, "api/users"},
	}

	for _, tt := range tests {
		result := testEvalOperator(tt.input)
		switch expected := tt.expected; {
		case expected == "api.github.com" || expected == "example.com":
			str, ok := result.(*evaluator.String)
			if !ok {
				t.Errorf("For input '%s': expected String, got %T (%+v)", tt.input, result, result)
				continue
			}
			if str.Value != expected {
				t.Errorf("For input '%s': expected '%s', got '%s'", tt.input, expected, str.Value)
			}
		default:
			str, ok := result.(*evaluator.String)
			if !ok {
				t.Errorf("For input '%s': expected String, got %T (%+v)", tt.input, result, result)
				continue
			}
			if str.Value != expected {
				t.Errorf("For input '%s': expected '%s', got '%s'", tt.input, expected, str.Value)
			}
		}
	}
}

// Test mixed path and URL operations
func TestMixedOperations(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`let p = (@/usr + "local") / "bin"; toString(p)`, "/usr/local/bin"},
		{`let u = (@https://api.example.com + "v1") + "users"; toString(u)`, "https://api.example.com/v1/users"},
		{`let p = @/var / "log"; let p2 = p + "system.log"; toString(p2)`, "/var/log/system.log"},
	}

	for _, tt := range tests {
		result := testEvalOperator(tt.input)
		str, ok := result.(*evaluator.String)
		if !ok {
			t.Errorf("For input '%s': expected String, got %T (%+v)", tt.input, result, result)
			continue
		}
		if str.Value != tt.expected {
			t.Errorf("For input '%s': expected '%s', got '%s'", tt.input, tt.expected, str.Value)
		}
	}
}

// Test operator precedence and associativity
func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`let p = @/usr + "local" + "bin"; p.basename`, "bin"},
		{`let u = @https://api.example.com + "v1" + "users"; u.pathname`, "v1/users"},
	}

	for _, tt := range tests {
		result := testEvalOperator(tt.input)
		str, ok := result.(*evaluator.String)
		if !ok {
			t.Errorf("For input '%s': expected String, got %T (%+v)", tt.input, result, result)
			continue
		}
		if str.Value != tt.expected {
			t.Errorf("For input '%s': expected '%s', got '%s'", tt.input, tt.expected, str.Value)
		}
	}
}
