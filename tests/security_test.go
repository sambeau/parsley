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

// TestSecurityWriteDefault tests that writes are denied by default
func TestSecurityWriteDefault(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")

	code := `"hello" ==> text("` + testFile + `")`

	env := evaluator.NewEnvironment()
	env.Security = &evaluator.SecurityPolicy{} // Default policy
	env.Filename = "test.pars"

	l := lexer.New(code)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Fatalf("Parse errors: %v", p.Errors())
	}

	result := evaluator.Eval(program, env)

	// Should get error about write access denied
	if result.Type() != evaluator.ERROR_OBJ {
		t.Errorf("Expected error, got %s", result.Type())
	}

	errObj := result.(*evaluator.Error)
	if !strings.Contains(errObj.Message, "write") && !strings.Contains(errObj.Message, "file write not allowed") {
		t.Errorf("Expected write-related error, got: %s", errObj.Message)
	}
}

// TestSecurityWriteAllowed tests that writes work when allowed
func TestSecurityWriteAllowed(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")

	code := `"hello" ==> text("` + testFile + `")`

	env := evaluator.NewEnvironment()
	env.Security = &evaluator.SecurityPolicy{
		AllowWrite: []string{tempDir},
	}
	env.Filename = "test.pars"

	l := lexer.New(code)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Fatalf("Parse errors: %v", p.Errors())
	}

	result := evaluator.Eval(program, env)

	// Should succeed
	if result.Type() == evaluator.ERROR_OBJ {
		t.Errorf("Unexpected error: %s", result.(*evaluator.Error).Message)
	}

	// Verify file was written
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Errorf("File not created: %v", err)
	}
	if string(content) != "hello" {
		t.Errorf("Expected 'hello', got '%s'", string(content))
	}
}

// TestSecurityWriteAll tests unrestricted writes
func TestSecurityWriteAll(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")

	code := `"hello" ==> text("` + testFile + `")`

	env := evaluator.NewEnvironment()
	env.Security = &evaluator.SecurityPolicy{
		AllowWriteAll: true,
	}
	env.Filename = "test.pars"

	l := lexer.New(code)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Fatalf("Parse errors: %v", p.Errors())
	}

	result := evaluator.Eval(program, env)

	// Should succeed
	if result.Type() == evaluator.ERROR_OBJ {
		t.Errorf("Unexpected error: %s", result.(*evaluator.Error).Message)
	}

	// Verify file was written
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Errorf("File not created: %v", err)
	}
	if string(content) != "hello" {
		t.Errorf("Expected 'hello', got '%s'", string(content))
	}
}

// TestSecurityReadRestricted tests read restrictions
func TestSecurityReadRestricted(t *testing.T) {
	// Create temporary directory and file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	os.WriteFile(testFile, []byte("hello"), 0644)

	code := `content <== text("` + testFile + `")`

	env := evaluator.NewEnvironment()
	env.Security = &evaluator.SecurityPolicy{
		RestrictRead: []string{tempDir},
	}
	env.Filename = "test.pars"

	l := lexer.New(code)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Fatalf("Parse errors: %v", p.Errors())
	}

	result := evaluator.Eval(program, env)

	// Should get error about read restricted
	if result.Type() != evaluator.ERROR_OBJ {
		t.Errorf("Expected error, got %s", result.Type())
	}

	errObj := result.(*evaluator.Error)
	if !strings.Contains(errObj.Message, "file read restricted") {
		t.Errorf("Expected 'file read restricted' error, got: %s", errObj.Message)
	}
}

// TestSecurityNoRead tests no-read policy
func TestSecurityNoRead(t *testing.T) {
	// Create temporary directory and file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	os.WriteFile(testFile, []byte("hello"), 0644)

	code := `content <== text("` + testFile + `")`

	env := evaluator.NewEnvironment()
	env.Security = &evaluator.SecurityPolicy{
		NoRead: true,
	}
	env.Filename = "test.pars"

	l := lexer.New(code)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Fatalf("Parse errors: %v", p.Errors())
	}

	result := evaluator.Eval(program, env)

	// Should get error about read access denied
	if result.Type() != evaluator.ERROR_OBJ {
		t.Errorf("Expected error, got %s", result.Type())
	}

	errObj := result.(*evaluator.Error)
	if !strings.Contains(errObj.Message, "file read access denied") {
		t.Errorf("Expected 'file read access denied' error, got: %s", errObj.Message)
	}
}

// TestSecurityExecuteDefault tests that module imports are denied by default
func TestSecurityExecuteDefault(t *testing.T) {
	// Create temporary module file
	tempDir := t.TempDir()
	moduleFile := filepath.Join(tempDir, "module.pars")
	os.WriteFile(moduleFile, []byte("let x = 42"), 0644)

	code := `import("` + moduleFile + `")`

	env := evaluator.NewEnvironment()
	env.Security = &evaluator.SecurityPolicy{} // Default policy
	env.Filename = "test.pars"

	l := lexer.New(code)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Fatalf("Parse errors: %v", p.Errors())
	}

	result := evaluator.Eval(program, env)

	// Should get error about execute access denied
	if result.Type() != evaluator.ERROR_OBJ {
		t.Errorf("Expected error, got %s", result.Type())
	}

	errObj := result.(*evaluator.Error)
	if !strings.Contains(errObj.Message, "execute") && !strings.Contains(errObj.Message, "script execution") {
		t.Errorf("Expected execute-related error, got: %s", errObj.Message)
	}
}

// TestSecurityExecuteAllowed tests that module imports work when allowed
func TestSecurityExecuteAllowed(t *testing.T) {
	// Create temporary module file
	tempDir := t.TempDir()
	moduleFile := filepath.Join(tempDir, "module.pars")
	os.WriteFile(moduleFile, []byte("let x = 42"), 0644)

	code := `import("` + moduleFile + `")`

	env := evaluator.NewEnvironment()
	env.Security = &evaluator.SecurityPolicy{
		AllowExecute: []string{tempDir},
	}
	env.Filename = "test.pars"

	l := lexer.New(code)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Fatalf("Parse errors: %v", p.Errors())
	}

	result := evaluator.Eval(program, env)

	// Should succeed
	if result.Type() == evaluator.ERROR_OBJ {
		t.Errorf("Unexpected error: %s", result.(*evaluator.Error).Message)
	}
}

// TestSecurityRemoveRequiresWrite tests that file removal requires write access
func TestSecurityRemoveRequiresWrite(t *testing.T) {
	// Create temporary file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	os.WriteFile(testFile, []byte("hello"), 0644)

	code := `file("` + testFile + `").remove()`

	env := evaluator.NewEnvironment()
	env.Security = &evaluator.SecurityPolicy{} // Default policy
	env.Filename = "test.pars"

	l := lexer.New(code)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Fatalf("Parse errors: %v", p.Errors())
	}

	result := evaluator.Eval(program, env)

	// Should get error about write access denied
	if result.Type() != evaluator.ERROR_OBJ {
		t.Errorf("Expected error, got %s", result.Type())
	}

	errObj := result.(*evaluator.Error)
	if !strings.Contains(errObj.Message, "write") {
		t.Errorf("Expected write-related error, got: %s", errObj.Message)
	}

	// Verify file still exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Error("File should not have been deleted")
	}
}

// TestSecurityDirListRequiresRead tests that directory listing requires read access
func TestSecurityDirListRequiresRead(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()

	code := `dir("` + tempDir + `").files`

	env := evaluator.NewEnvironment()
	env.Security = &evaluator.SecurityPolicy{
		RestrictRead: []string{tempDir},
	}
	env.Filename = "test.pars"

	l := lexer.New(code)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Fatalf("Parse errors: %v", p.Errors())
	}

	result := evaluator.Eval(program, env)

	// Should get error about read restricted
	if result.Type() != evaluator.ERROR_OBJ {
		t.Errorf("Expected error, got %s", result.Type())
	}

	errObj := result.(*evaluator.Error)
	if !strings.Contains(errObj.Message, "read") {
		t.Errorf("Expected read-related error, got: %s", errObj.Message)
	}
}
