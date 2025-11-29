// filehandle_test.go - Tests for file handle objects and format factories

package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sambeau/parsley/pkg/evaluator"
	"github.com/sambeau/parsley/pkg/lexer"
	"github.com/sambeau/parsley/pkg/parser"
)

func testEvalFileHandle(input string) evaluator.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := evaluator.NewEnvironment()
	return evaluator.Eval(program, env)
}

// Test file() builtin basic creation
func TestFileBuiltinBasic(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Format auto-detection from extension
		{`let f = file(@./config.json); f.format`, "json"},
		{`let f = file(@./data.csv); f.format`, "csv"},
		{`let f = file(@./readme.txt); f.format`, "text"},
		{`let f = file(@./app.log); f.format`, "lines"},
		{`let f = file(@./readme.md); f.format`, "text"},
		{`let f = file(@./index.html); f.format`, "text"},
		// __type field
		{`let f = file(@./test.json); f.__type`, "file"},
	}

	for _, tt := range tests {
		result := testEvalFileHandle(tt.input)
		if result == nil {
			t.Errorf("For input '%s': evaluation returned nil", tt.input)
			continue
		}
		str, ok := result.(*evaluator.String)
		if !ok {
			t.Errorf("For input '%s': expected String, got %T (%v)", tt.input, result, result.Inspect())
			continue
		}
		if str.Value != tt.expected {
			t.Errorf("For input '%s': expected '%s', got '%s'", tt.input, tt.expected, str.Value)
		}
	}
}

// Test format-specific factories
func TestFormatFactories(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`let f = JSON(@./config.json); f.format`, "json"},
		{`let f = CSV(@./data.csv); f.format`, "csv"},
		{`let f = lines(@./app.log); f.format`, "lines"},
		{`let f = text(@./readme.txt); f.format`, "text"},
		{`let f = bytes(@./image.png); f.format`, "bytes"},
		// Format factories override extension inference
		{`let f = JSON(@./config.txt); f.format`, "json"},
		{`let f = text(@./data.json); f.format`, "text"},
	}

	for _, tt := range tests {
		result := testEvalFileHandle(tt.input)
		if result == nil {
			t.Errorf("For input '%s': evaluation returned nil", tt.input)
			continue
		}
		str, ok := result.(*evaluator.String)
		if !ok {
			t.Errorf("For input '%s': expected String, got %T (%v)", tt.input, result, result.Inspect())
			continue
		}
		if str.Value != tt.expected {
			t.Errorf("For input '%s': expected '%s', got '%s'", tt.input, tt.expected, str.Value)
		}
	}
}

// Test file handle computed properties with real files
func TestFileComputedPropertiesWithRealFiles(t *testing.T) {
	// Create a temp file for testing
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(testFile, []byte("hello world"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create a temp subdirectory
	subDir := filepath.Join(tmpDir, "subdir")
	err = os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"exists true", `let f = file(@` + testFile + `); f.exists`, "true"},
		{"exists false", `let f = file(@` + filepath.Join(tmpDir, "nonexistent.txt") + `); f.exists`, "false"},
		{"size", `let f = file(@` + testFile + `); f.size`, "11"},
		{"isFile", `let f = file(@` + testFile + `); f.isFile`, "true"},
		{"isDir file", `let f = file(@` + testFile + `); f.isDir`, "false"},
		{"isDir dir", `let f = file(@` + subDir + `); f.isDir`, "true"},
		{"basename", `let f = file(@` + testFile + `); f.basename`, "test.txt"},
		{"ext", `let f = file(@` + testFile + `); f.ext`, "txt"},
		{"stem", `let f = file(@` + testFile + `); f.stem`, "test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testEvalFileHandle(tt.input)
			if result == nil {
				t.Fatalf("evaluation returned nil")
			}
			resultStr := result.Inspect()
			if resultStr != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, resultStr)
			}
		})
	}
}

// Test file handle modified datetime property
func TestFileModifiedDatetime(t *testing.T) {
	// Create a temp file for testing
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(testFile, []byte("hello world"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test that modified returns a datetime dictionary
	input := `let f = file(@` + testFile + `); f.modified.__type`
	result := testEvalFileHandle(input)
	if result == nil {
		t.Fatalf("evaluation returned nil")
	}
	str, ok := result.(*evaluator.String)
	if !ok {
		t.Fatalf("expected String, got %T (%v)", result, result.Inspect())
	}
	if str.Value != "datetime" {
		t.Errorf("expected 'datetime', got '%s'", str.Value)
	}
}

// Test file handle path property
func TestFilePathProperty(t *testing.T) {
	// Test that path returns a path dictionary
	input := `let f = file(@./test.json); f.path.__type`
	result := testEvalFileHandle(input)
	if result == nil {
		t.Fatalf("evaluation returned nil")
	}
	str, ok := result.(*evaluator.String)
	if !ok {
		t.Fatalf("expected String, got %T (%v)", result, result.Inspect())
	}
	if str.Value != "path" {
		t.Errorf("expected 'path', got '%s'", str.Value)
	}
}

// Test file handle with string argument (instead of path literal)
func TestFileWithStringArgument(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`let f = file("./config.json"); f.format`, "json"},
		{`let f = JSON("./data.json"); f.format`, "json"},
		{`let f = CSV("./data.csv"); f.format`, "csv"},
	}

	for _, tt := range tests {
		result := testEvalFileHandle(tt.input)
		if result == nil {
			t.Errorf("For input '%s': evaluation returned nil", tt.input)
			continue
		}
		str, ok := result.(*evaluator.String)
		if !ok {
			t.Errorf("For input '%s': expected String, got %T (%v)", tt.input, result, result.Inspect())
			continue
		}
		if str.Value != tt.expected {
			t.Errorf("For input '%s': expected '%s', got '%s'", tt.input, tt.expected, str.Value)
		}
	}
}

// Test CSV with header option
func TestCSVWithOptions(t *testing.T) {
	// Test that options are stored
	input := `let f = CSV(@./data.csv, {header: true}); f.options`
	result := testEvalFileHandle(input)
	if result == nil {
		t.Fatalf("evaluation returned nil")
	}
	// Options should be a dictionary
	dict, ok := result.(*evaluator.Dictionary)
	if !ok {
		t.Fatalf("expected Dictionary, got %T (%v)", result, result.Inspect())
	}
	// Check that header is in options
	if _, ok := dict.Pairs["header"]; !ok {
		t.Errorf("expected 'header' key in options, got %v", dict.Inspect())
	}
}

// Test error handling for invalid arguments
func TestFileBuiltinErrors(t *testing.T) {
	tests := []struct {
		input string
	}{
		{`file()`},    // No arguments
		{`file(123)`}, // Wrong type
		{`JSON()`},    // No arguments
		{`CSV()`},     // No arguments
		{`lines()`},   // No arguments
		{`text()`},    // No arguments
		{`bytes()`},   // No arguments
	}

	for _, tt := range tests {
		result := testEvalFileHandle(tt.input)
		if result == nil {
			t.Errorf("For input '%s': evaluation returned nil", tt.input)
			continue
		}
		err, ok := result.(*evaluator.Error)
		if !ok {
			t.Errorf("For input '%s': expected Error, got %T (%v)", tt.input, result, result.Inspect())
		} else if err.Message == "" {
			t.Errorf("For input '%s': expected error message, got empty string", tt.input)
		}
	}
}

// Test file remove() method
func TestFileRemoveMethod(t *testing.T) {
	// Create a temp directory and test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "remove_test.txt")
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Fatalf("Test file was not created")
	}

	// Remove the file using Parsley (use string path instead of path literal)
	input := `let f = file("` + testFile + `"); f.remove()`
	result := testEvalFileHandle(input)
	if result == nil {
		t.Fatalf("evaluation returned nil")
	}

	// Check that remove() returns NULL on success
	if result != evaluator.NULL {
		if err, ok := result.(*evaluator.Error); ok {
			t.Fatalf("remove() returned error: %s", err.Message)
		} else {
			t.Fatalf("expected NULL, got %T (%v)", result, result.Inspect())
		}
	}

	// Verify file no longer exists
	if _, err := os.Stat(testFile); !os.IsNotExist(err) {
		t.Errorf("File still exists after remove()")
	}
}

// Test file remove() method with non-existent file
func TestFileRemoveNonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "nonexistent.txt")

	// Try to remove a file that doesn't exist (use string path)
	input := `let f = file("` + testFile + `"); f.remove()`
	result := testEvalFileHandle(input)
	if result == nil {
		t.Fatalf("evaluation returned nil")
	}

	// Should return an error
	err, ok := result.(*evaluator.Error)
	if !ok {
		t.Fatalf("expected Error for non-existent file, got %T (%v)", result, result.Inspect())
	}
	if err.Message == "" {
		t.Errorf("expected error message, got empty string")
	}
}

// Test file remove() method with wrong number of arguments
func TestFileRemoveWrongArgs(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Try to call remove with an argument (use string path)
	input := `let f = file("` + testFile + `"); f.remove(true)`
	result := testEvalFileHandle(input)
	if result == nil {
		t.Fatalf("evaluation returned nil")
	}

	// Should return an error
	errObj, ok := result.(*evaluator.Error)
	if !ok {
		t.Fatalf("expected Error for wrong number of args, got %T (%v)", result, result.Inspect())
	}
	if errObj.Message == "" {
		t.Errorf("expected error message, got empty string")
	}
}
