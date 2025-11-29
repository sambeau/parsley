package main

import (
	"strings"
	"testing"

	"github.com/sambeau/parsley/pkg/evaluator"
	"github.com/sambeau/parsley/pkg/lexer"
	"github.com/sambeau/parsley/pkg/parser"
)

// TestPathTemplateInterpolation tests path templates with variable interpolation
func TestPathTemplateInterpolation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple variable interpolation",
			input:    `name = "config"; p = @(./data/{name}.json); p.string`,
			expected: "./data/config.json",
		},
		{
			name:     "multiple interpolations",
			input:    `dir = "src"; file = "main"; p = @(./{dir}/{file}.go); p.string`,
			expected: "./src/main.go",
		},
		{
			name:     "expression interpolation",
			input:    `n = 1; p = @(./file{n + 1}.txt); p.string`,
			expected: "./file2.txt",
		},
		{
			name:     "absolute path with interpolation",
			input:    `user = "sam"; p = @(/home/{user}/docs); p.string`,
			expected: "/home/sam/docs",
		},
		{
			name:     "home path with interpolation",
			input:    `folder = "projects"; p = @(~/{folder}/code); p.string`,
			expected: "~/projects/code",
		},
		{
			name:     "nested property access",
			input:    `config = {dir: "build"}; p = @(./{config.dir}/output); p.string`,
			expected: "./build/output",
		},
		{
			name:     "array index interpolation",
			input:    `dirs = ["src", "lib"]; p = @(./{dirs[0]}/main.go); p.string`,
			expected: "./src/main.go",
		},
		{
			name:     "function call in interpolation",
			input:    `p = @(./{"hello".toUpper()}.txt); p.string`,
			expected: "./HELLO.txt",
		},
		{
			name:     "number conversion",
			input:    `version = 2; p = @(./v{version}/api); p.string`,
			expected: "./v2/api",
		},
		{
			name:     "string concatenation in interpolation",
			input:    `prefix = "test"; suffix = "data"; p = @(./{prefix + "_" + suffix}.json); p.string`,
			expected: "./test_data.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) > 0 {
				t.Fatalf("parser errors: %v", p.Errors())
			}

			env := evaluator.NewEnvironment()
			result := evaluator.Eval(program, env)

			if result == nil {
				t.Fatalf("evaluation returned nil")
			}

			if errObj, ok := result.(*evaluator.Error); ok {
				t.Fatalf("evaluation error: %s", errObj.Message)
			}

			strResult := result.Inspect()
			if strResult != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, strResult)
			}
		})
	}
}

// TestUrlTemplateInterpolation tests URL templates with variable interpolation
func TestUrlTemplateInterpolation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "path segment interpolation",
			input:    `version = "v2"; u = @(https://api.example.com/{version}/users); u.string`,
			expected: "https://api.example.com/v2/users",
		},
		{
			name:     "host interpolation",
			input:    `host = "api.test.com"; u = @(https://{host}/data); u.string`,
			expected: "https://api.test.com/data",
		},
		{
			name:     "port interpolation",
			input:    `port = 8080; u = @(http://localhost:{port}/api); u.port`,
			expected: "8080",
		},
		{
			name:     "fragment interpolation",
			input:    `section = "intro"; u = @(https://docs.com/guide#{section}); u.fragment`,
			expected: "intro",
		},
		{
			name:     "expression in URL",
			input:    `n = 1; u = @(https://api.com/v{n + 1}/data); u.string`,
			expected: "https://api.com/v2/data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) > 0 {
				t.Fatalf("parser errors: %v", p.Errors())
			}

			env := evaluator.NewEnvironment()
			result := evaluator.Eval(program, env)

			if result == nil {
				t.Fatalf("evaluation returned nil")
			}

			if errObj, ok := result.(*evaluator.Error); ok {
				t.Fatalf("evaluation error: %s", errObj.Message)
			}

			strResult := result.Inspect()
			if strResult != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, strResult)
			}
		})
	}
}

// TestPathTemplateComponents tests that path template results have correct components
func TestPathTemplateComponents(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		checkField string
		expected   string
	}{
		{
			name:       "basename from template",
			input:      `name = "main"; p = @(./src/{name}.go); p.basename`,
			checkField: "basename",
			expected:   "main.go",
		},
		{
			name:       "ext from template",
			input:      `ext = "json"; p = @(./config.{ext}); p.ext`,
			checkField: "ext",
			expected:   "json",
		},
		{
			name:       "dir from template",
			input:      `folder = "lib"; p = @(./{folder}/utils.js); p.dir`,
			checkField: "dir",
			expected:   "./lib",
		},
		{
			name:       "stem from template",
			input:      `name = "data"; p = @(./{name}.csv); p.stem`,
			checkField: "stem",
			expected:   "data",
		},
		{
			name:       "isAbsolute from absolute template",
			input:      `dir = "usr"; p = @(/{dir}/local); p.isAbsolute`,
			checkField: "isAbsolute",
			expected:   "true",
		},
		{
			name:       "isAbsolute from relative template",
			input:      `dir = "src"; p = @(./{dir}/main.go); p.isAbsolute`,
			checkField: "isAbsolute",
			expected:   "false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) > 0 {
				t.Fatalf("parser errors: %v", p.Errors())
			}

			env := evaluator.NewEnvironment()
			result := evaluator.Eval(program, env)

			if result == nil {
				t.Fatalf("evaluation returned nil")
			}

			if errObj, ok := result.(*evaluator.Error); ok {
				t.Fatalf("evaluation error: %s", errObj.Message)
			}

			strResult := result.Inspect()
			if strResult != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, strResult)
			}
		})
	}
}

// TestUrlTemplateComponents tests that URL template results have correct components
func TestUrlTemplateComponents(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "scheme preserved",
			input:    `path = "api"; u = @(https://example.com/{path}); u.scheme`,
			expected: "https",
		},
		{
			name:     "host preserved",
			input:    `version = "v1"; u = @(https://api.example.com/{version}/users); u.host`,
			expected: "api.example.com",
		},
		{
			name:     "full url string",
			input:    `id = 123; u = @(https://api.com/users/{id}/profile); u.string`,
			expected: "https://api.com/users/123/profile",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) > 0 {
				t.Fatalf("parser errors: %v", p.Errors())
			}

			env := evaluator.NewEnvironment()
			result := evaluator.Eval(program, env)

			if result == nil {
				t.Fatalf("evaluation returned nil")
			}

			if errObj, ok := result.(*evaluator.Error); ok {
				t.Fatalf("evaluation error: %s", errObj.Message)
			}

			strResult := result.Inspect()
			if strResult != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, strResult)
			}
		})
	}
}

// TestPathTemplateErrors tests error handling in path templates
func TestPathTemplateErrors(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedError string
	}{
		{
			name:          "unclosed brace",
			input:         `name = "test"; p = @(./data/{name.json)`,
			expectedError: "unclosed",
		},
		{
			name:          "empty interpolation",
			input:         `p = @(./data/{}.json)`,
			expectedError: "empty interpolation",
		},
		{
			name:          "undefined variable",
			input:         `p = @(./data/{undefined}/file)`,
			expectedError: "identifier not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) > 0 {
				// Parser error is acceptable for some tests
				return
			}

			env := evaluator.NewEnvironment()
			result := evaluator.Eval(program, env)

			if result == nil {
				t.Fatalf("evaluation returned nil, expected error")
			}

			errObj, ok := result.(*evaluator.Error)
			if !ok {
				t.Fatalf("expected error, got %T: %v", result, result)
			}

			if !strings.Contains(strings.ToLower(errObj.Message), strings.ToLower(tt.expectedError)) {
				t.Errorf("expected error containing %q, got %q", tt.expectedError, errObj.Message)
			}
		})
	}
}

// TestLexerPathUrlTemplateTokens tests that the lexer correctly tokenizes path/URL templates
func TestLexerPathUrlTemplateTokens(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected lexer.TokenType
		literal  string
	}{
		{
			name:     "path template relative",
			input:    `@(./path/{name}/file)`,
			expected: lexer.PATH_TEMPLATE,
			literal:  "./path/{name}/file",
		},
		{
			name:     "path template absolute",
			input:    `@(/usr/{user}/local)`,
			expected: lexer.PATH_TEMPLATE,
			literal:  "/usr/{user}/local",
		},
		{
			name:     "path template home",
			input:    `@(~/{dir}/config)`,
			expected: lexer.PATH_TEMPLATE,
			literal:  "~/{dir}/config",
		},
		{
			name:     "url template https",
			input:    `@(https://api.com/{v}/users)`,
			expected: lexer.URL_TEMPLATE,
			literal:  "https://api.com/{v}/users",
		},
		{
			name:     "url template http",
			input:    `@(http://localhost:{port}/api)`,
			expected: lexer.URL_TEMPLATE,
			literal:  "http://localhost:{port}/api",
		},
		{
			name:     "static path still works",
			input:    `@./static/path`,
			expected: lexer.PATH_LITERAL,
			literal:  "./static/path",
		},
		{
			name:     "static url still works",
			input:    `@https://example.com`,
			expected: lexer.URL_LITERAL,
			literal:  "https://example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			tok := l.NextToken()

			if tok.Type != tt.expected {
				t.Errorf("expected token type %s, got %s", tt.expected, tok.Type)
			}

			if tok.Literal != tt.literal {
				t.Errorf("expected literal %q, got %q", tt.literal, tok.Literal)
			}
		})
	}
}

// TestPathTemplateInExpressions tests path templates used in various expression contexts
func TestPathTemplateInExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "path template in function",
			input:    `makePath = fn(name) { @(./{name}.json) }; makePath("config").string`,
			expected: "./config.json",
		},
		{
			name:     "path template assigned to variable",
			input:    `name = "test"; p = @(./data/{name}.txt); p.basename`,
			expected: "test.txt",
		},
		{
			name:     "path template with computed string",
			input:    `base = "config"; ext = "yaml"; p = @(./{base}.{ext}); p.string`,
			expected: "./config.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) > 0 {
				t.Fatalf("parser errors: %v", p.Errors())
			}

			env := evaluator.NewEnvironment()
			result := evaluator.Eval(program, env)

			if result == nil {
				t.Fatalf("evaluation returned nil")
			}

			if errObj, ok := result.(*evaluator.Error); ok {
				t.Fatalf("evaluation error: %s", errObj.Message)
			}

			strResult := result.Inspect()
			if strResult != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, strResult)
			}
		})
	}
}
