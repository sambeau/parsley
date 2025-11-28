package main

import (
	"testing"

	"github.com/sambeau/parsley/pkg/evaluator"
	"github.com/sambeau/parsley/pkg/lexer"
	"github.com/sambeau/parsley/pkg/parser"
)

func testPathUrl(input string) (evaluator.Object, bool) {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		return nil, true
	}

	env := evaluator.NewEnvironment()
	result := evaluator.Eval(program, env)

	if err, ok := result.(*evaluator.Error); ok {
		return err, true
	}

	return result, false
}

func TestPathParsing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "absolute unix path",
			input:    `path("/usr/local/bin")`,
			expected: `{__type: path, absolute: true, components: , usr, local, bin}`,
		},
		{
			name:     "relative path",
			input:    `path("./config/app.json")`,
			expected: `{__type: path, absolute: false, components: ., config, app.json}`,
		},
		{
			name:     "path components access",
			input:    `let p = path("/usr/local/bin"); p.components`,
			expected: `, usr, local, bin`,
		},
		{
			name:     "path components index",
			input:    `let p = path("/usr/local/bin"); p.components[-1]`,
			expected: `bin`,
		},
		{
			name:     "path absolute flag",
			input:    `let p = path("/usr/local/bin"); p.absolute`,
			expected: `true`,
		},
		{
			name:     "relative path absolute flag",
			input:    `let p = path("./config/app.json"); p.absolute`,
			expected: `false`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, hasErr := testPathUrl(tt.input)
			if hasErr {
				t.Fatalf("unexpected error: %v", result)
			}
			if result.Inspect() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.Inspect())
			}
		})
	}
}

func TestPathComputedProperties(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basename",
			input:    `let p = path("/usr/local/bin/tool.exe"); p.basename`,
			expected: `tool.exe`,
		},
		{
			name:     "extension",
			input:    `let p = path("/usr/local/bin/tool.exe"); p.extension`,
			expected: `exe`,
		},
		{
			name:     "ext alias",
			input:    `let p = path("/usr/local/bin/tool.exe"); p.ext`,
			expected: `exe`,
		},
		{
			name:     "stem",
			input:    `let p = path("/usr/local/bin/tool.exe"); p.stem`,
			expected: `tool`,
		},
		{
			name:     "dirname as string",
			input:    `let p = path("/usr/local/bin/tool.exe"); toString(p.dirname)`,
			expected: `/usr/local/bin`,
		},
		{
			name:     "parent alias",
			input:    `let p = path("/usr/local/bin/tool.exe"); toString(p.parent)`,
			expected: `/usr/local/bin`,
		},
		{
			name:     "no extension",
			input:    `let p = path("/usr/local/bin/README"); p.extension`,
			expected: ``,
		},
		{
			name:     "multi-part extension",
			input:    `let p = path("/archive/file.tar.gz"); p.extension`,
			expected: `gz`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, hasErr := testPathUrl(tt.input)
			if hasErr {
				t.Fatalf("unexpected error: %v", result)
			}
			if result.Inspect() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.Inspect())
			}
		})
	}
}

func TestPathToString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "absolute path to string",
			input:    `let p = path("/usr/local/bin"); toString(p)`,
			expected: `/usr/local/bin`,
		},
		{
			name:     "relative path to string",
			input:    `let p = path("./config/app.json"); toString(p)`,
			expected: `./config/app.json`,
		},
		{
			name:     "roundtrip path",
			input:    `let p = path("/usr/local/bin"); toString(path(toString(p)))`,
			expected: `/usr/local/bin`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, hasErr := testPathUrl(tt.input)
			if hasErr {
				t.Fatalf("unexpected error: %v", result)
			}
			if result.Inspect() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.Inspect())
			}
		})
	}
}

func TestUrlParsing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic https URL",
			input:    `let u = url("https://example.com"); u.scheme`,
			expected: `https`,
		},
		{
			name:     "URL host",
			input:    `let u = url("https://example.com"); u.host`,
			expected: `example.com`,
		},
		{
			name:     "URL with port",
			input:    `let u = url("https://example.com:8080"); u.port`,
			expected: `8080`,
		},
		{
			name:     "URL path",
			input:    `let u = url("https://example.com/api/users"); u.path`,
			expected: `, api, users`,
		},
		{
			name:     "URL query param",
			input:    `let u = url("https://example.com/api?limit=10"); u.query.limit`,
			expected: `10`,
		},
		{
			name:     "URL fragment",
			input:    `let u = url("https://example.com/page#section"); u.fragment`,
			expected: `section`,
		},
		{
			name:     "URL with username",
			input:    `let u = url("https://user@example.com"); u.username`,
			expected: `user`,
		},
		{
			name:     "URL with username and password",
			input:    `let u = url("https://user:pass@example.com"); u.password`,
			expected: `pass`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, hasErr := testPathUrl(tt.input)
			if hasErr {
				t.Fatalf("unexpected error: %v", result)
			}
			if result.Inspect() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.Inspect())
			}
		})
	}
}

func TestUrlComputedProperties(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "origin without port",
			input:    `let u = url("https://example.com/api"); u.origin`,
			expected: `https://example.com`,
		},
		{
			name:     "origin with port",
			input:    `let u = url("https://example.com:8080/api"); u.origin`,
			expected: `https://example.com:8080`,
		},
		{
			name:     "pathname",
			input:    `let u = url("https://example.com/api/users"); u.pathname`,
			expected: `/api/users`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, hasErr := testPathUrl(tt.input)
			if hasErr {
				t.Fatalf("unexpected error: %v", result)
			}
			if result.Inspect() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.Inspect())
			}
		})
	}
}

func TestUrlToString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple URL to string",
			input:    `let u = url("https://example.com"); toString(u)`,
			expected: `https://example.com`,
		},
		{
			name:     "URL with path to string",
			input:    `let u = url("https://example.com/api/users"); toString(u)`,
			expected: `https://example.com/api/users`,
		},
		{
			name:     "URL with query to string",
			input:    `let u = url("https://example.com/api?limit=10"); toString(u)`,
			expected: `https://example.com/api?limit=10`,
		},
		{
			name:     "URL with fragment to string",
			input:    `let u = url("https://example.com/page#section"); toString(u)`,
			expected: `https://example.com/page#section`,
		},
		{
			name:     "complex URL roundtrip",
			input:    `let u = url("https://user@example.com:8080/api?limit=10#top"); toString(u)`,
			expected: `https://user@example.com:8080/api?limit=10#top`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, hasErr := testPathUrl(tt.input)
			if hasErr {
				t.Fatalf("unexpected error: %v", result)
			}
			if result.Inspect() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.Inspect())
			}
		})
	}
}

func TestPathArrayComposition(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "filter path components using join",
			input:    `let p = path("/usr/local/bin"); p.components.filter(fn(c) { c != "local" }).join("/")`,
			expected: `/usr/bin`,
		},
		{
			name:     "map path components",
			input:    `let p = path("/usr/local/bin"); p.components.map(fn(c) { len(c) })`,
			expected: `0, 3, 5, 3`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, hasErr := testPathUrl(tt.input)
			if hasErr {
				t.Fatalf("unexpected error: %v", result)
			}
			if result.Inspect() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.Inspect())
			}
		})
	}
}
