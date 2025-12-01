package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sambeau/parsley/pkg/ast"
	"github.com/sambeau/parsley/pkg/evaluator"
	"github.com/sambeau/parsley/pkg/lexer"
	"github.com/sambeau/parsley/pkg/parser"
)

func evalExport(input string, filename string) evaluator.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := evaluator.NewEnvironment()
	env.Filename = filename
	// Enable execute-all for module import tests
	env.Security = &evaluator.SecurityPolicy{
		AllowExecuteAll: true,
	}
	return evaluator.Eval(program, env)
}

func parseExportProgram(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

// TestExportKeyword tests the export keyword functionality
func TestExportKeyword(t *testing.T) {
	// Create a temporary directory for module files
	tmpDir, err := os.MkdirTemp("", "parsley-export-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name           string
		moduleCode     string
		mainCode       string
		expectedOutput string
		expectError    bool
	}{
		{
			name: "export let - explicit export",
			moduleCode: `
export let add = fn(a, b) { a + b }
let multiply = fn(a, b) { a * b }
`,
			mainCode:       `let mod = import("%s"); mod.add(2, 3)`,
			expectedOutput: "5",
		},
		{
			name: "export without let",
			moduleCode: `
export add = fn(a, b) { a + b }
`,
			mainCode:       `let mod = import("%s"); mod.add(2, 3)`,
			expectedOutput: "5",
		},
		{
			name: "let without export - backward compat",
			moduleCode: `
let add = fn(a, b) { a + b }
`,
			mainCode:       `let mod = import("%s"); mod.add(2, 3)`,
			expectedOutput: "5",
		},
		{
			name: "bare assignment not exported",
			moduleCode: `
add = fn(a, b) { a + b }
let greet = fn(name) { name }
`,
			mainCode:       `let mod = import("%s"); mod.greet("World")`,
			expectedOutput: "World",
		},
		{
			name: "exported variable value",
			moduleCode: `
export let version = "1.0.0"
`,
			mainCode:       `let mod = import("%s"); mod.version`,
			expectedOutput: "1.0.0",
		},
		{
			name: "export multiple variables",
			moduleCode: `
export let x = 10
export let y = 20
`,
			mainCode:       `let mod = import("%s"); mod.x + mod.y`,
			expectedOutput: "30",
		},
		{
			name: "export single numeric value",
			moduleCode: `
export Pi = 3.141592653589793
`,
			mainCode:       `let mod = import("%s"); mod.Pi`,
			expectedOutput: "3.141592653589793",
		},
		{
			name: "export single tag value",
			moduleCode: `
export Logo = <img src="logo.png" alt="Logo"/>
`,
			mainCode:       `let mod = import("%s"); mod.Logo`,
			expectedOutput: `<img src="logo.png" alt="Logo" />`,
		},
		{
			name: "destructure single export",
			moduleCode: `
export Pi = 3.14159
`,
			mainCode:       `let {Pi} = import("%s"); Pi`,
			expectedOutput: "3.14159",
		},
		{
			name: "export tag and use in template",
			moduleCode: `
export Header = <header><h1>Welcome</h1></header>
`,
			mainCode:       `let mod = import("%s"); <div>{mod.Header}</div>`,
			expectedOutput: `<div><header><h1>Welcome</h1></header></div>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Write module file
			moduleFile := filepath.Join(tmpDir, tt.name+".pars")
			err := os.WriteFile(moduleFile, []byte(tt.moduleCode), 0644)
			if err != nil {
				t.Fatalf("Failed to write module file: %v", err)
			}

			// Format main code with module path
			mainCode := formatString(tt.mainCode, moduleFile)

			// Create a dummy main file path in the same directory
			mainFile := filepath.Join(tmpDir, "main.pars")
			result := evalExport(mainCode, mainFile)

			if tt.expectError {
				if result.Type() != evaluator.ERROR_OBJ {
					t.Errorf("expected error, got %s (%s)", result.Type(), result.Inspect())
				}
				return
			}

			if result.Type() == evaluator.ERROR_OBJ {
				t.Errorf("unexpected error: %s", result.Inspect())
				return
			}

			if result.Inspect() != tt.expectedOutput {
				t.Errorf("expected %s, got %s", tt.expectedOutput, result.Inspect())
			}
		})
	}
}

// TestExportDestructuring tests export with destructuring
func TestExportDestructuring(t *testing.T) {
	// Create a temporary directory for module files
	tmpDir, err := os.MkdirTemp("", "parsley-export-destruct-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name           string
		moduleCode     string
		mainCode       string
		expectedOutput string
	}{
		{
			name: "export let array destructuring",
			moduleCode: `
export let [a, b] = [1, 2]
`,
			mainCode:       `let mod = import("%s"); mod.a + mod.b`,
			expectedOutput: "3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Write module file
			moduleFile := filepath.Join(tmpDir, tt.name+".pars")
			err := os.WriteFile(moduleFile, []byte(tt.moduleCode), 0644)
			if err != nil {
				t.Fatalf("Failed to write module file: %v", err)
			}

			// Format main code with module path
			mainCode := formatString(tt.mainCode, moduleFile)

			// Create a dummy main file path
			mainFile := filepath.Join(tmpDir, "main.pars")
			result := evalExport(mainCode, mainFile)

			if result.Type() == evaluator.ERROR_OBJ {
				t.Errorf("unexpected error: %s", result.Inspect())
				return
			}

			if result.Inspect() != tt.expectedOutput {
				t.Errorf("expected %s, got %s", tt.expectedOutput, result.Inspect())
			}
		})
	}
}

// TestExportStatementString tests that export statements stringify correctly
func TestExportStatementString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "export let x = 5",
			expected: "export let x = 5;",
		},
		{
			input:    "export x = 10",
			expected: "export x = 10;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			program := parseExportProgram(tt.input)
			if program == nil {
				t.Fatalf("failed to parse program")
			}
			if len(program.Statements) != 1 {
				t.Fatalf("expected 1 statement, got %d", len(program.Statements))
			}
			if program.Statements[0].String() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, program.Statements[0].String())
			}
		})
	}
}

// TestBareAssignmentNotExported tests that bare assignments are NOT exported
func TestBareAssignmentNotExported(t *testing.T) {
	// Create a temporary directory for module files
	tmpDir, err := os.MkdirTemp("", "parsley-bare-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Module with only bare assignment (no let or export)
	moduleCode := `
internal = fn(a, b) { a + b }
`
	moduleFile := filepath.Join(tmpDir, "bare.pars")
	err = os.WriteFile(moduleFile, []byte(moduleCode), 0644)
	if err != nil {
		t.Fatalf("Failed to write module file: %v", err)
	}

	// Try to access the bare assignment from imported module
	mainCode := formatString(`let mod = import("%s"); mod.internal(1, 2)`, moduleFile)
	mainFile := filepath.Join(tmpDir, "main.pars")
	result := evalExport(mainCode, mainFile)

	// Accessing a non-exported variable should return null or error
	// The function won't exist on the module, so calling it should fail
	if result.Type() != evaluator.ERROR_OBJ && result.Inspect() != "null" {
		t.Errorf("expected error or null when accessing non-exported var, got %s (%s)",
			result.Type(), result.Inspect())
	}
}

// Helper function to format strings with %s
func formatString(format string, args ...string) string {
	result := format
	for _, arg := range args {
		for i := 0; i < len(result)-1; i++ {
			if result[i] == '%' && result[i+1] == 's' {
				result = result[:i] + arg + result[i+2:]
				break
			}
		}
	}
	return result
}
