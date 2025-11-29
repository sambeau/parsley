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

func testEvalDirWrite(input string) evaluator.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := evaluator.NewEnvironment()
	// Enable write access for directory operations
	env.Security = &evaluator.SecurityPolicy{
		AllowWriteAll: true,
	}
	return evaluator.Eval(program, env)
}

// Test file().mkdir() method
func TestFileMethodMkdir(t *testing.T) {
	// Create a temp directory for testing
	tmpDir, err := os.MkdirTemp("", "parsley-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testPath := filepath.Join(tmpDir, "test-dir")

	script := `file(@` + testPath + `).mkdir()`

	result := testEvalDirWrite(script)
	if result.Type() != evaluator.NULL_OBJ {
		t.Errorf("Expected NULL, got=%s", result.Inspect())
	}

	// Verify directory was created
	if _, err := os.Stat(testPath); os.IsNotExist(err) {
		t.Error("Directory was not created")
	}
}

// Test file().mkdir({parents: true}) for nested paths
func TestFileMethodMkdirParents(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "parsley-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testPath := filepath.Join(tmpDir, "parent", "child", "grandchild")

	script := `file(@` + testPath + `).mkdir({parents: true})`

	result := testEvalDirWrite(script)
	if result.Type() != evaluator.NULL_OBJ {
		t.Errorf("Expected NULL, got=%s", result.Inspect())
	}

	// Verify nested directories were created
	if _, err := os.Stat(testPath); os.IsNotExist(err) {
		t.Error("Nested directories were not created")
	}
}

// Test file().rmdir() on empty directory
func TestFileMethodRmdir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "parsley-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testPath := filepath.Join(tmpDir, "to-remove")
	if err := os.Mkdir(testPath, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	script := `file(@` + testPath + `).rmdir()`

	result := testEvalDirWrite(script)
	if result.Type() != evaluator.NULL_OBJ {
		t.Errorf("Expected NULL, got=%s", result.Inspect())
	}

	// Verify directory was removed
	if _, err := os.Stat(testPath); !os.IsNotExist(err) {
		t.Error("Directory was not removed")
	}
}

// Test file().rmdir({recursive: true}) on non-empty directory
func TestFileMethodRmdirRecursive(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "parsley-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testPath := filepath.Join(tmpDir, "to-remove")
	if err := os.Mkdir(testPath, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create a file inside the directory
	testFile := filepath.Join(testPath, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	script := `file(@` + testPath + `).rmdir({recursive: true})`

	result := testEvalDirWrite(script)
	if result.Type() != evaluator.NULL_OBJ {
		t.Errorf("Expected NULL, got=%s", result.Inspect())
	}

	// Verify directory and contents were removed
	if _, err := os.Stat(testPath); !os.IsNotExist(err) {
		t.Error("Directory was not removed")
	}
}

// Test dir().mkdir() method
func TestDirMethodMkdir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "parsley-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testPath := filepath.Join(tmpDir, "test-dir-2")

	script := `dir(@` + testPath + `).mkdir()`

	result := testEvalDirWrite(script)
	if result.Type() != evaluator.NULL_OBJ {
		t.Errorf("Expected NULL, got=%s", result.Inspect())
	}

	// Verify directory was created
	if _, err := os.Stat(testPath); os.IsNotExist(err) {
		t.Error("Directory was not created")
	}
}

// Test dir().rmdir() method
func TestDirMethodRmdir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "parsley-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testPath := filepath.Join(tmpDir, "to-remove-2")
	if err := os.Mkdir(testPath, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	script := `dir(@` + testPath + `).rmdir()`

	result := testEvalDirWrite(script)
	if result.Type() != evaluator.NULL_OBJ {
		t.Errorf("Expected NULL, got=%s", result.Inspect())
	}

	// Verify directory was removed
	if _, err := os.Stat(testPath); !os.IsNotExist(err) {
		t.Error("Directory was not removed")
	}
}

// Test error handling for mkdir without parents option
func TestFileMethodMkdirErrorNoParents(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "parsley-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testPath := filepath.Join(tmpDir, "nonexistent", "nested")

	script := `file(@` + testPath + `).mkdir()`

	result := testEvalDirWrite(script)
	if result.Type() != evaluator.ERROR_OBJ {
		t.Errorf("Expected ERROR, got=%s", result.Type())
	}
	if !strings.Contains(result.Inspect(), "no such file or directory") {
		t.Errorf("Expected 'no such file or directory' error, got: %s", result.Inspect())
	}
}

// Test error handling for rmdir on non-empty directory without recursive option
func TestFileMethodRmdirErrorNotEmpty(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "parsley-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testPath := filepath.Join(tmpDir, "not-empty")
	if err := os.Mkdir(testPath, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create a file inside
	testFile := filepath.Join(testPath, "file.txt")
	if err := os.WriteFile(testFile, []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	script := `file(@` + testPath + `).rmdir()`

	result := testEvalDirWrite(script)
	if result.Type() != evaluator.ERROR_OBJ {
		t.Errorf("Expected ERROR, got=%s", result.Type())
	}
	if !strings.Contains(result.Inspect(), "directory not empty") {
		t.Errorf("Expected 'directory not empty' error, got: %s", result.Inspect())
	}
}
