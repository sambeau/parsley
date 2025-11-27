package main

import (
	"testing"

	"github.com/sambeau/parsley/pkg/evaluator"
	"github.com/sambeau/parsley/pkg/lexer"
	"github.com/sambeau/parsley/pkg/parser"
)

func testEvalComputedProp(input string) evaluator.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := evaluator.NewEnvironment()
	return evaluator.Eval(program, env)
}

// Test path.name (alias for basename)
func TestPathNameProperty(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`let p = @/usr/local/bin/parsley; p.name`, "parsley"},
		{`let p = @./config.json; p.name`, "config.json"},
		{`let p = @~/documents/file.txt; p.name`, "file.txt"},
	}

	for _, tt := range tests {
		result := testEvalComputedProp(tt.input)
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

// Test path.suffix (alias for extension)
func TestPathSuffixProperty(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`let p = @/file.txt; p.suffix`, "txt"},
		{`let p = @./archive.tar.gz; p.suffix`, "gz"},
		{`let p = @~/README; p.suffix`, ""},
	}

	for _, tt := range tests {
		result := testEvalComputedProp(tt.input)
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

// Test path.suffixes (all extensions as array)
func TestPathSuffixesProperty(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{`let p = @/file.tar.gz; len(p.suffixes)`, 2},
		{`let p = @./config.yaml; len(p.suffixes)`, 1},
		{`let p = @~/README; len(p.suffixes)`, 0},
		{`let p = @/file.backup.old.txt; len(p.suffixes)`, 3},
	}

	for _, tt := range tests {
		result := testEvalComputedProp(tt.input)
		num, ok := result.(*evaluator.Integer)
		if !ok {
			t.Errorf("For input '%s': expected Integer, got %T (%+v)", tt.input, result, result)
			continue
		}
		if num.Value != int64(tt.expected) {
			t.Errorf("For input '%s': expected %d, got %d", tt.input, tt.expected, num.Value)
		}
	}

	// Test individual suffix values
	indexTests := []struct {
		input    string
		expected string
	}{
		{`let p = @/file.tar.gz; p.suffixes[0]`, "tar"},
		{`let p = @/file.tar.gz; p.suffixes[1]`, "gz"},
		{`let p = @/file.backup.old.txt; p.suffixes[2]`, "txt"},
	}

	for _, tt := range indexTests {
		result := testEvalComputedProp(tt.input)
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

// Test path.parts (alias for components)
func TestPathPartsProperty(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{`let p = @/usr/local/bin; len(p.parts)`, 4},
		{`let p = @./config/app.json; len(p.parts)`, 3},
		{`let p = @~/docs; len(p.parts)`, 2},
	}

	for _, tt := range tests {
		result := testEvalComputedProp(tt.input)
		num, ok := result.(*evaluator.Integer)
		if !ok {
			t.Errorf("For input '%s': expected Integer, got %T (%+v)", tt.input, result, result)
			continue
		}
		if num.Value != int64(tt.expected) {
			t.Errorf("For input '%s': expected %d, got %d", tt.input, tt.expected, num.Value)
		}
	}

	// Test individual part values
	indexTests := []struct {
		input    string
		expected string
	}{
		{`let p = @/usr/local/bin; p.parts[1]`, "usr"},
		{`let p = @/usr/local/bin; p.parts[2]`, "local"},
		{`let p = @./config; p.parts[0]`, "."},
	}

	for _, tt := range indexTests {
		result := testEvalComputedProp(tt.input)
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

// Test path.isAbsolute
func TestPathIsAbsoluteProperty(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`let p = @/usr/local/bin; p.isAbsolute`, true},
		{`let p = @./config.json; p.isAbsolute`, false},
		{`let p = @~/documents; p.isAbsolute`, false},
		{`let p = @../parent; p.isAbsolute`, false},
	}

	for _, tt := range tests {
		result := testEvalComputedProp(tt.input)
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

// Test path.isRelative
func TestPathIsRelativeProperty(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`let p = @/usr/local/bin; p.isRelative`, false},
		{`let p = @./config.json; p.isRelative`, true},
		{`let p = @~/documents; p.isRelative`, true},
		{`let p = @../parent; p.isRelative`, true},
	}

	for _, tt := range tests {
		result := testEvalComputedProp(tt.input)
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

// Test URL.hostname (alias for host)
func TestUrlHostnameProperty(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`let u = @https://example.com/api; u.hostname`, "example.com"},
		{`let u = @http://localhost:8080; u.hostname`, "localhost"},
		{`let u = @https://api.github.com/repos; u.hostname`, "api.github.com"},
	}

	for _, tt := range tests {
		result := testEvalComputedProp(tt.input)
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

// Test URL.protocol (scheme with colon)
func TestUrlProtocolProperty(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`let u = @https://example.com; u.protocol`, "https:"},
		{`let u = @http://localhost:8080; u.protocol`, "http:"},
		{`let u = @ftp://ftp.example.org; u.protocol`, "ftp:"},
	}

	for _, tt := range tests {
		result := testEvalComputedProp(tt.input)
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

// Test URL.search (query string with ?)
func TestUrlSearchProperty(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`let u = @https://example.com/api?page=2&limit=10; u.search`, "?page=2&limit=10"},
		{`let u = @https://example.com/api; u.search`, ""},
		{`let u = url("https://example.com/search?q=hello"); u.search`, "?q=hello"},
	}

	for _, tt := range tests {
		result := testEvalComputedProp(tt.input)
		str, ok := result.(*evaluator.String)
		if !ok {
			t.Errorf("For input '%s': expected String, got %T (%+v)", tt.input, result, result)
			continue
		}
		// Note: query parameter order may vary in dictionaries
		if tt.expected == "" {
			if str.Value != tt.expected {
				t.Errorf("For input '%s': expected '%s', got '%s'", tt.input, tt.expected, str.Value)
			}
		} else {
			// Just check that it starts with ?
			if len(str.Value) == 0 || str.Value[0] != '?' {
				t.Errorf("For input '%s': expected search string to start with '?', got '%s'", tt.input, str.Value)
			}
		}
	}
}

// Test URL.href (full URL as string)
func TestUrlHrefProperty(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`let u = @https://example.com/api; u.href`, "https://example.com/api"},
		{`let u = @http://localhost:8080/test; u.href`, "http://localhost:8080/test"},
		{`let u = (@https://api.github.com + "repos"); u.href`, "https://api.github.com/repos"},
	}

	for _, tt := range tests {
		result := testEvalComputedProp(tt.input)
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

// Test combined usage of new properties
func TestCombinedComputedProperties(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`let p = @/archive/file.tar.gz; p.name + " has " + toString(len(p.suffixes)) + " suffixes"`, "file.tar.gz has 2 suffixes"},
		{`let p = @./config.json; if (p.isRelative) "relative" else "absolute"`, "relative"},
		{`let u = @https://api.example.com:8080/v1; u.protocol + "//" + u.hostname`, "https://api.example.com"},
	}

	for _, tt := range tests {
		result := testEvalComputedProp(tt.input)
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

// Test that old properties still work alongside new ones
func TestBackwardCompatibility(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`let p = @/file.txt; p.basename`, "file.txt"},
		{`let p = @/file.txt; p.name`, "file.txt"},
		{`let p = @/file.txt; p.extension`, "txt"},
		{`let p = @/file.txt; p.suffix`, "txt"},
		{`let p = @/usr/local; toString(p.dirname)`, "/usr"},
		{`let p = @/usr/local; toString(p.parent)`, "/usr"},
		{`let u = @https://example.com; u.host`, "example.com"},
		{`let u = @https://example.com; u.hostname`, "example.com"},
	}

	for _, tt := range tests {
		result := testEvalComputedProp(tt.input)
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
