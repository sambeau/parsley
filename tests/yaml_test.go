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

// Helper function for evaluating Parsley code with a filename context
func testEvalYAMLWithFilename(input string, filename string) evaluator.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := evaluator.NewEnvironment()
	env.Filename = filename
	// Enable write access for tests
	env.Security = &evaluator.SecurityPolicy{
		AllowWriteAll: true,
	}
	return evaluator.Eval(program, env)
}

// TestYAMLBasic tests basic YAML file reading
func TestYAMLBasic(t *testing.T) {
	// Create a temp directory for test files
	tmpDir, err := os.MkdirTemp("", "parsley-yaml-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test YAML file
	yamlContent := `name: John Doe
age: 30
active: true
tags:
  - developer
  - golang
`
	yamlPath := filepath.Join(tmpDir, "test.yaml")
	if err := os.WriteFile(yamlPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to write YAML file: %v", err)
	}

	// Test file path - used for relative path resolution
	testFilePath := filepath.Join(tmpDir, "test.pars")

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Read name field",
			input:    `let config <== YAML(@./test.yaml); config.name`,
			expected: "John Doe",
		},
		{
			name:     "Read age field",
			input:    `let config <== YAML(@./test.yaml); config.age`,
			expected: "30",
		},
		{
			name:     "Read boolean field",
			input:    `let config <== YAML(@./test.yaml); config.active`,
			expected: "true",
		},
		{
			name:     "Read array element",
			input:    `let config <== YAML(@./test.yaml); config.tags[0]`,
			expected: "developer",
		},
		{
			name:     "Read second array element",
			input:    `let config <== YAML(@./test.yaml); config.tags[1]`,
			expected: "golang",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testEvalYAMLWithFilename(tt.input, testFilePath)
			if result == nil {
				t.Fatalf("Result is nil")
			}

			resultStr := result.Inspect()
			if resultStr != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, resultStr)
			}
		})
	}
}

// TestYAMLNestedObjects tests YAML with nested objects
func TestYAMLNestedObjects(t *testing.T) {
	// Create a temp directory for test files
	tmpDir, err := os.MkdirTemp("", "parsley-yaml-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test YAML file with nested objects
	yamlContent := `server:
  host: localhost
  port: 8080
  ssl:
    enabled: true
    cert: /path/to/cert
database:
  connection: postgres://localhost/db
  pool_size: 10
`
	yamlPath := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(yamlPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to write YAML file: %v", err)
	}

	testFilePath := filepath.Join(tmpDir, "test.pars")

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Read nested host",
			input:    `let config <== YAML(@./config.yaml); config.server.host`,
			expected: "localhost",
		},
		{
			name:     "Read nested port",
			input:    `let config <== YAML(@./config.yaml); config.server.port`,
			expected: "8080",
		},
		{
			name:     "Read deeply nested boolean",
			input:    `let config <== YAML(@./config.yaml); config.server.ssl.enabled`,
			expected: "true",
		},
		{
			name:     "Read database pool_size",
			input:    `let config <== YAML(@./config.yaml); config.database.pool_size`,
			expected: "10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testEvalYAMLWithFilename(tt.input, testFilePath)
			if result == nil {
				t.Fatalf("Result is nil")
			}

			resultStr := result.Inspect()
			if resultStr != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, resultStr)
			}
		})
	}
}

// TestYAMLDateParsing tests date parsing in YAML
func TestYAMLDateParsing(t *testing.T) {
	// Create a temp directory for test files
	tmpDir, err := os.MkdirTemp("", "parsley-yaml-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test YAML file with dates
	yamlContent := `title: Test Event
date: 2024-06-15
created_at: 2024-01-01T10:30:00Z
`
	yamlPath := filepath.Join(tmpDir, "event.yaml")
	if err := os.WriteFile(yamlPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to write YAML file: %v", err)
	}

	testFilePath := filepath.Join(tmpDir, "test.pars")

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Read title",
			input:    `let event <== YAML(@./event.yaml); event.title`,
			expected: "Test Event",
		},
		{
			name:     "Read date year",
			input:    `let event <== YAML(@./event.yaml); event.date.year`,
			expected: "2024",
		},
		{
			name:     "Read created_at month",
			input:    `let event <== YAML(@./event.yaml); event.created_at.month`,
			expected: "1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testEvalYAMLWithFilename(tt.input, testFilePath)
			if result == nil {
				t.Fatalf("Result is nil")
			}

			resultStr := result.Inspect()
			if resultStr != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, resultStr)
			}
		})
	}
}

// TestYAMLArray tests YAML with top-level arrays
func TestYAMLArray(t *testing.T) {
	// Create a temp directory for test files
	tmpDir, err := os.MkdirTemp("", "parsley-yaml-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a YAML file with a top-level array
	yamlContent := `- name: Alice
  role: admin
- name: Bob
  role: user
- name: Charlie
  role: user
`
	yamlPath := filepath.Join(tmpDir, "users.yaml")
	if err := os.WriteFile(yamlPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to write YAML file: %v", err)
	}

	testFilePath := filepath.Join(tmpDir, "test.pars")

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Read first user name",
			input:    `let users <== YAML(@./users.yaml); users[0].name`,
			expected: "Alice",
		},
		{
			name:     "Read second user role",
			input:    `let users <== YAML(@./users.yaml); users[1].role`,
			expected: "user",
		},
		{
			name:     "Read third user name",
			input:    `let users <== YAML(@./users.yaml); users[2].name`,
			expected: "Charlie",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testEvalYAMLWithFilename(tt.input, testFilePath)
			if result == nil {
				t.Fatalf("Result is nil")
			}

			resultStr := result.Inspect()
			if resultStr != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, resultStr)
			}
		})
	}
}

// TestYAMLWriteRead tests writing and reading YAML files
func TestYAMLWriteRead(t *testing.T) {
	// Create a temp directory for test files
	tmpDir, err := os.MkdirTemp("", "parsley-yaml-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFilePath := filepath.Join(tmpDir, "test.pars")
	yamlPath := filepath.Join(tmpDir, "output.yaml")

	// Test writing a dictionary to YAML using absolute path
	input := `let data = {
	name: "Test App",
	version: "1.0.0",
	enabled: true
}
data ==> YAML(@` + yamlPath + `)
"done"`

	result := testEvalYAMLWithFilename(input, testFilePath)
	if result == nil {
		t.Fatalf("Result is nil")
	}

	resultStr := result.Inspect()
	if resultStr != "done" {
		t.Errorf("Expected 'done', got %q", resultStr)
	}

	// Verify the file was created with proper YAML content
	content, err := os.ReadFile(yamlPath)
	if err != nil {
		t.Fatalf("Failed to read output YAML: %v", err)
	}

	// Check that it contains expected YAML content
	contentStr := string(content)
	if !strings.Contains(contentStr, "name:") {
		t.Errorf("Expected YAML to contain 'name:', got %s", contentStr)
	}
	if !strings.Contains(contentStr, "Test App") {
		t.Errorf("Expected YAML to contain 'Test App', got %s", contentStr)
	}

	// Now read it back
	input2 := `let loaded <== YAML(@` + yamlPath + `)
loaded.name`

	result2 := testEvalYAMLWithFilename(input2, testFilePath)
	if result2 == nil {
		t.Fatalf("Result2 is nil")
	}

	expected := "Test App"
	if result2.Inspect() != expected {
		t.Errorf("Expected %q, got %q", expected, result2.Inspect())
	}
}

// TestYAMLEmptyFile tests reading an empty YAML file
func TestYAMLEmptyFile(t *testing.T) {
	// Create a temp directory for test files
	tmpDir, err := os.MkdirTemp("", "parsley-yaml-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create an empty YAML file
	yamlPath := filepath.Join(tmpDir, "empty.yaml")
	if err := os.WriteFile(yamlPath, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to write YAML file: %v", err)
	}

	testFilePath := filepath.Join(tmpDir, "test.pars")

	// Empty YAML should parse to null/nil
	input := `let config <== YAML(@./empty.yaml); config`
	result := testEvalYAMLWithFilename(input, testFilePath)

	if result == nil {
		t.Fatalf("Result is nil")
	}

	// Empty YAML should return null
	expected := "null"
	if result.Inspect() != expected {
		t.Errorf("Expected %q, got %q", expected, result.Inspect())
	}
}
