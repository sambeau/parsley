// write_operator_test.go - Tests for write operators ==> and ==>>

package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sambeau/parsley/pkg/evaluator"
	"github.com/sambeau/parsley/pkg/lexer"
	"github.com/sambeau/parsley/pkg/parser"
)

func testEvalWriteOp(input string) evaluator.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := evaluator.NewEnvironment()
	// Enable write-all for write operator tests
	env.Security = &evaluator.SecurityPolicy{
		AllowWriteAll: true,
	}
	return evaluator.Eval(program, env)
}

// TestWriteOperatorText tests writing text files
func TestWriteOperatorText(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "parsley_write_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name     string
		code     string
		file     string
		expected string
	}{
		{
			name:     "write string to text file",
			code:     `"Hello, World!" ==> text("` + filepath.Join(tmpDir, "test1.txt") + `")`,
			file:     "test1.txt",
			expected: "Hello, World!",
		},
		{
			name:     "write variable to text file",
			code:     `let msg = "Greetings!"; msg ==> text("` + filepath.Join(tmpDir, "test2.txt") + `")`,
			file:     "test2.txt",
			expected: "Greetings!",
		},
		{
			name:     "write number to text file",
			code:     `42 ==> text("` + filepath.Join(tmpDir, "test3.txt") + `")`,
			file:     "test3.txt",
			expected: "42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testEvalWriteOp(tt.code)
			if result != nil && result.Type() == "ERROR" {
				t.Errorf("Evaluation error: %s", result.Inspect())
				return
			}

			filePath := filepath.Join(tmpDir, tt.file)
			content, err := os.ReadFile(filePath)
			if err != nil {
				t.Errorf("Failed to read file: %v", err)
				return
			}
			if string(content) != tt.expected {
				t.Errorf("Expected file content %q, got %q", tt.expected, string(content))
			}
		})
	}
}

// TestWriteOperatorJSON tests writing JSON files
func TestWriteOperatorJSON(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "parsley_write_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name         string
		code         string
		file         string
		expectedJSON string
	}{
		{
			name:         "write dict to JSON file",
			code:         `let d = {name: "Alice", age: 30}; d ==> JSON("` + filepath.Join(tmpDir, "test1.json") + `")`,
			file:         "test1.json",
			expectedJSON: `"name": "Alice"`,
		},
		{
			name:         "write array to JSON file",
			code:         `[1, 2, 3] ==> JSON("` + filepath.Join(tmpDir, "test2.json") + `")`,
			file:         "test2.json",
			expectedJSON: `[`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testEvalWriteOp(tt.code)
			if result != nil && result.Type() == "ERROR" {
				t.Errorf("Evaluation error: %s", result.Inspect())
				return
			}

			filePath := filepath.Join(tmpDir, tt.file)
			content, err := os.ReadFile(filePath)
			if err != nil {
				t.Errorf("Failed to read file: %v", err)
				return
			}
			if !strings.Contains(string(content), tt.expectedJSON) {
				t.Errorf("Expected file to contain %q, got %q", tt.expectedJSON, string(content))
			}
		})
	}
}

// TestWriteOperatorLines tests writing line-based files
func TestWriteOperatorLines(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "parsley_write_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name     string
		code     string
		file     string
		expected string
	}{
		{
			name:     "write array of strings as lines",
			code:     `["line1", "line2", "line3"] ==> lines("` + filepath.Join(tmpDir, "test1.log") + `")`,
			file:     "test1.log",
			expected: "line1\nline2\nline3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testEvalWriteOp(tt.code)
			if result != nil && result.Type() == "ERROR" {
				t.Errorf("Evaluation error: %s", result.Inspect())
				return
			}

			filePath := filepath.Join(tmpDir, tt.file)
			content, err := os.ReadFile(filePath)
			if err != nil {
				t.Errorf("Failed to read file: %v", err)
				return
			}
			if string(content) != tt.expected {
				t.Errorf("Expected file content %q, got %q", tt.expected, string(content))
			}
		})
	}
}

// TestWriteOperatorCSV tests writing CSV files
func TestWriteOperatorCSV(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "parsley_write_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name     string
		code     string
		file     string
		contains string
	}{
		{
			name:     "write array of arrays as CSV",
			code:     `[["a", "b"], ["c", "d"]] ==> CSV("` + filepath.Join(tmpDir, "test1.csv") + `")`,
			file:     "test1.csv",
			contains: "a,b",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testEvalWriteOp(tt.code)
			if result != nil && result.Type() == "ERROR" {
				t.Errorf("Evaluation error: %s", result.Inspect())
				return
			}

			filePath := filepath.Join(tmpDir, tt.file)
			content, err := os.ReadFile(filePath)
			if err != nil {
				t.Errorf("Failed to read file: %v", err)
				return
			}
			if !strings.Contains(string(content), tt.contains) {
				t.Errorf("Expected file to contain %q, got %q", tt.contains, string(content))
			}
		})
	}
}

// TestAppendOperator tests the ==>> append operator
func TestAppendOperator(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "parsley_append_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create initial file
	initialFile := filepath.Join(tmpDir, "append.txt")
	if err := os.WriteFile(initialFile, []byte("initial"), 0644); err != nil {
		t.Fatalf("Failed to create initial file: %v", err)
	}

	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name:     "append to existing file",
			code:     `"-appended" ==>> text("` + initialFile + `")`,
			expected: "initial-appended",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testEvalWriteOp(tt.code)
			if result != nil && result.Type() == "ERROR" {
				t.Errorf("Evaluation error: %s", result.Inspect())
				return
			}

			content, err := os.ReadFile(initialFile)
			if err != nil {
				t.Errorf("Failed to read file: %v", err)
				return
			}
			if string(content) != tt.expected {
				t.Errorf("Expected file content %q, got %q", tt.expected, string(content))
			}
		})
	}
}

// TestAppendToNewFile tests appending to a non-existent file (creates it)
func TestAppendToNewFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "parsley_append_new_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	newFile := filepath.Join(tmpDir, "new.txt")
	code := `"created by append" ==>> text("` + newFile + `")`

	result := testEvalWriteOp(code)
	if result != nil && result.Type() == "ERROR" {
		t.Errorf("Evaluation error: %s", result.Inspect())
		return
	}

	content, err := os.ReadFile(newFile)
	if err != nil {
		t.Errorf("Failed to read file: %v", err)
		return
	}
	if string(content) != "created by append" {
		t.Errorf("Expected file content %q, got %q", "created by append", string(content))
	}
}

// TestWriteOperatorBytes tests writing raw bytes
func TestWriteOperatorBytes(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "parsley_bytes_write_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	bytesFile := filepath.Join(tmpDir, "test.bin")
	code := `[72, 101, 108, 108, 111] ==> bytes("` + bytesFile + `")`

	result := testEvalWriteOp(code)
	if result != nil && result.Type() == "ERROR" {
		t.Errorf("Evaluation error: %s", result.Inspect())
		return
	}

	content, err := os.ReadFile(bytesFile)
	if err != nil {
		t.Errorf("Failed to read file: %v", err)
		return
	}
	if string(content) != "Hello" {
		t.Errorf("Expected file content 'Hello', got %q", string(content))
	}
}

// TestWriteOperatorErrors tests error handling for write operators
func TestWriteOperatorErrors(t *testing.T) {
	tests := []struct {
		name          string
		code          string
		errorContains string
	}{
		{
			name:          "write to non-file handle",
			code:          `"test" ==> "not a file"`,
			errorContains: "requires a file handle",
		},
		{
			name:          "write bytes with non-array",
			code:          `"not bytes" ==> bytes("/tmp/test.bin")`,
			errorContains: "requires an array",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testEvalWriteOp(tt.code)
			if result == nil || result.Type() != "ERROR" {
				t.Errorf("Expected error, got %v", result)
				return
			}
			if !strings.Contains(result.Inspect(), tt.errorContains) {
				t.Errorf("Expected error containing %q, got %q", tt.errorContains, result.Inspect())
			}
		})
	}
}

// TestWriteReadRoundtrip tests writing and reading back data
func TestWriteReadRoundtrip(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "parsley_roundtrip_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name:     "JSON roundtrip",
			code:     `let d = {name: "Bob", score: 95}; d ==> JSON("` + filepath.Join(tmpDir, "rt.json") + `"); let data <== JSON("` + filepath.Join(tmpDir, "rt.json") + `"); data.name`,
			expected: "Bob",
		},
		{
			name:     "text roundtrip",
			code:     `"round trip test" ==> text("` + filepath.Join(tmpDir, "rt.txt") + `"); let content <== text("` + filepath.Join(tmpDir, "rt.txt") + `"); content`,
			expected: "round trip test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testEvalWriteOp(tt.code)
			if result == nil {
				t.Error("Expected result, got nil")
				return
			}
			if result.Type() == "ERROR" {
				t.Errorf("Evaluation error: %s", result.Inspect())
				return
			}
			if result.Inspect() != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result.Inspect())
			}
		})
	}
}
