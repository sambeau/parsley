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

func testEvalReadOp(input string) evaluator.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := evaluator.NewEnvironment()
	return evaluator.Eval(program, env)
}

func testEvalReadOpString(input string) string {
	result := testEvalReadOp(input)
	if result == nil {
		return "<nil>"
	}
	return result.Inspect()
}

// TestReadOperatorBasic tests the basic <== operator functionality
func TestReadOperatorBasic(t *testing.T) {
	// Create a temp directory for test files
	tmpDir, err := os.MkdirTemp("", "parsley_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test text file
	textPath := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(textPath, []byte("Hello, World!"), 0644); err != nil {
		t.Fatalf("Failed to write text file: %v", err)
	}

	// Create test JSON file
	jsonPath := filepath.Join(tmpDir, "test.json")
	if err := os.WriteFile(jsonPath, []byte(`{"name": "Alice", "age": 30}`), 0644); err != nil {
		t.Fatalf("Failed to write JSON file: %v", err)
	}

	// Create test lines file
	linesPath := filepath.Join(tmpDir, "test.lines")
	if err := os.WriteFile(linesPath, []byte("line1\nline2\nline3"), 0644); err != nil {
		t.Fatalf("Failed to write lines file: %v", err)
	}

	// Create test CSV file with header
	csvPath := filepath.Join(tmpDir, "test.csv")
	if err := os.WriteFile(csvPath, []byte("name,age\nAlice,30\nBob,25"), 0644); err != nil {
		t.Fatalf("Failed to write CSV file: %v", err)
	}

	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name:     "read text file with let",
			code:     `let content <== text("` + textPath + `"); content`,
			expected: "Hello, World!",
		},
		{
			name:     "read JSON file with let",
			code:     `let data <== JSON("` + jsonPath + `"); data.name`,
			expected: "Alice",
		},
		{
			name:     "read JSON file age",
			code:     `let data <== JSON("` + jsonPath + `"); data.age`,
			expected: "30",
		},
		{
			name:     "read lines file",
			code:     `let lines <== lines("` + linesPath + `"); lines.length()`,
			expected: "3",
		},
		{
			name:     "read lines file first line",
			code:     `let lines <== lines("` + linesPath + `"); lines[0]`,
			expected: "line1",
		},
		{
			name:     "read CSV file",
			code:     `let rows <== CSV("` + csvPath + `"); rows.length()`,
			expected: "2",
		},
		{
			name:     "read CSV file first row name",
			code:     `let rows <== CSV("` + csvPath + `"); rows[0].name`,
			expected: "Alice",
		},
		{
			name:     "read CSV file second row age",
			code:     `let rows <== CSV("` + csvPath + `"); rows[1].age`,
			expected: "25",
		},
		{
			name:     "read with dict destructuring",
			code:     `let {name, age} <== JSON("` + jsonPath + `"); name + " is " + age`,
			expected: "Alice is 30",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testEvalReadOpString(tt.code)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestReadOperatorReassignment tests <== for variable reassignment
func TestReadOperatorReassignment(t *testing.T) {
	// Create a temp directory for test files
	tmpDir, err := os.MkdirTemp("", "parsley_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	file1 := filepath.Join(tmpDir, "file1.txt")
	if err := os.WriteFile(file1, []byte("First"), 0644); err != nil {
		t.Fatalf("Failed to write file1: %v", err)
	}

	file2 := filepath.Join(tmpDir, "file2.txt")
	if err := os.WriteFile(file2, []byte("Second"), 0644); err != nil {
		t.Fatalf("Failed to write file2: %v", err)
	}

	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name:     "reassign variable with <==",
			code:     `let x = "initial"; x <== text("` + file1 + `"); x`,
			expected: "First",
		},
		{
			name:     "read multiple times",
			code:     `let x <== text("` + file1 + `"); let y <== text("` + file2 + `"); x + " " + y`,
			expected: "First Second",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testEvalReadOpString(tt.code)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestReadOperatorBytes tests reading raw bytes
func TestReadOperatorBytes(t *testing.T) {
	// Create a temp directory for test files
	tmpDir, err := os.MkdirTemp("", "parsley_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test binary file
	binPath := filepath.Join(tmpDir, "test.bin")
	if err := os.WriteFile(binPath, []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f}, 0644); err != nil {
		t.Fatalf("Failed to write binary file: %v", err)
	}

	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name:     "read bytes length",
			code:     `let data <== bytes("` + binPath + `"); data.length()`,
			expected: "5",
		},
		{
			name:     "read bytes first byte",
			code:     `let data <== bytes("` + binPath + `"); data[0]`,
			expected: "72", // 0x48 = 72 = 'H'
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testEvalReadOpString(tt.code)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestReadOperatorErrors tests error handling for <==
func TestReadOperatorErrors(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		shouldError bool
		errorContains string
	}{
		{
			name:        "read non-existent file",
			code:        `let x <== text("/nonexistent/file.txt"); x`,
			shouldError: true,
			errorContains: "failed to read file",
		},
		{
			name:        "read without file handle",
			code:        `let x <== "not a file"; x`,
			shouldError: true,
			errorContains: "requires a file handle",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testEvalReadOp(tt.code)
			if tt.shouldError {
				if result == nil {
					t.Error("Expected error, got nil")
					return
				}
				if result.Type() != "ERROR" {
					t.Errorf("Expected ERROR, got %s", result.Type())
					return
				}
				if tt.errorContains != "" {
					if !containsSubstr(result.Inspect(), tt.errorContains) {
						t.Errorf("Expected error containing %q, got %q", tt.errorContains, result.Inspect())
					}
				}
			}
		})
	}
}

// TestReadOperatorCSVNoHeader tests reading CSV without headers
// NOTE: Currently skipped because dictionary property assignment (dict.key = value)
// is not yet supported in Parsley. Will be enabled when that feature is added.
func TestReadOperatorCSVNoHeader(t *testing.T) {
	t.Skip("Dictionary property assignment not yet supported - skipping CSV no header tests")

	// Create a temp directory for test files
	tmpDir, err := os.MkdirTemp("", "parsley_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test CSV file without header
	csvPath := filepath.Join(tmpDir, "data.csv")
	if err := os.WriteFile(csvPath, []byte("Alice,30\nBob,25"), 0644); err != nil {
		t.Fatalf("Failed to write CSV file: %v", err)
	}

	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name:     "read CSV no header - row count",
			code:     `let handle = file("` + csvPath + `"); handle.format = "csv-noheader"; let rows <== handle; rows.length()`,
			expected: "2",
		},
		{
			name:     "read CSV no header - first row first element",
			code:     `let handle = file("` + csvPath + `"); handle.format = "csv-noheader"; let rows <== handle; rows[0][0]`,
			expected: "Alice",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testEvalReadOpString(tt.code)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// containsString checks if a string contains a substring
func containsSubstr(s, substr string) bool {
	return strings.Contains(s, substr)
}
