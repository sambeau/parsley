package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sambeau/parsley/pkg/evaluator"
	"github.com/sambeau/parsley/pkg/lexer"
	"github.com/sambeau/parsley/pkg/parser"
)

func testEvalDir(input string) evaluator.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := evaluator.NewEnvironment()
	return evaluator.Eval(program, env)
}

func testEvalDirString(input string) string {
	result := testEvalDir(input)
	if result == nil {
		return "<nil>"
	}
	return result.Inspect()
}

// TestDirBasic tests the dir() function and basic properties
func TestDirBasic(t *testing.T) {
	// Create a temp directory structure for testing
	tempDir := t.TempDir()

	// Create test files
	testFiles := []string{"test1.txt", "test2.json", "test3.csv"}
	for _, name := range testFiles {
		os.WriteFile(filepath.Join(tempDir, name), []byte("content"), 0644)
	}

	// Create a subdirectory
	subDir := filepath.Join(tempDir, "subdir")
	os.Mkdir(subDir, 0755)
	os.WriteFile(filepath.Join(subDir, "nested.txt"), []byte("nested"), 0644)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "dir exists",
			input:    `let d = dir("` + tempDir + `"); d.exists`,
			expected: "true",
		},
		{
			name:     "dir isDir",
			input:    `let d = dir("` + tempDir + `"); d.isDir`,
			expected: "true",
		},
		{
			name:     "dir isFile",
			input:    `let d = dir("` + tempDir + `"); d.isFile`,
			expected: "false",
		},
		{
			name:     "dir count",
			input:    `let d = dir("` + tempDir + `"); d.count`,
			expected: "4", // 3 files + 1 subdir
		},
		{
			name:     "dir name",
			input:    `let d = dir("` + tempDir + `"); d.name`,
			expected: filepath.Base(tempDir),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testEvalDirString(tt.input)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

// TestDirRead tests reading directory contents with <==
func TestDirRead(t *testing.T) {
	// Create a temp directory structure for testing
	tempDir := t.TempDir()

	// Create test files
	testFiles := []string{"a.txt", "b.json", "c.csv"}
	for _, name := range testFiles {
		os.WriteFile(filepath.Join(tempDir, name), []byte("content"), 0644)
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "read dir returns array",
			input:    `let d = dir("` + tempDir + `"); let files <== d; files.length()`,
			expected: "3",
		},
		{
			name:     "read dir first file has basename",
			input:    `let d = dir("` + tempDir + `"); let files <== d; files[0].basename`,
			expected: "a.txt",
		},
		{
			name:     "read dir files have isFile property",
			input:    `let d = dir("` + tempDir + `"); let files <== d; files[0].isFile`,
			expected: "true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testEvalDirString(tt.input)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

// TestDirFilesProperty tests the .files property
func TestDirFilesProperty(t *testing.T) {
	// Create a temp directory
	tempDir := t.TempDir()

	// Create test files
	testFiles := []string{"file1.txt", "file2.txt"}
	for _, name := range testFiles {
		os.WriteFile(filepath.Join(tempDir, name), []byte("content"), 0644)
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "files property returns array",
			input:    `let d = dir("` + tempDir + `"); d.files.length()`,
			expected: "2",
		},
		{
			name:     "files property first file basename",
			input:    `let d = dir("` + tempDir + `"); d.files[0].basename`,
			expected: "file1.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testEvalDirString(tt.input)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

// TestGlobBasic tests the glob() function
func TestGlobBasic(t *testing.T) {
	// Create a temp directory
	tempDir := t.TempDir()

	// Create test files with different extensions
	os.WriteFile(filepath.Join(tempDir, "test1.txt"), []byte("content"), 0644)
	os.WriteFile(filepath.Join(tempDir, "test2.txt"), []byte("content"), 0644)
	os.WriteFile(filepath.Join(tempDir, "data.json"), []byte("{}"), 0644)
	os.WriteFile(filepath.Join(tempDir, "data.csv"), []byte("a,b"), 0644)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "glob matches txt files",
			input:    `let files = glob("` + tempDir + `/*.txt"); files.length()`,
			expected: "2",
		},
		{
			name:     "glob matches json files",
			input:    `let files = glob("` + tempDir + `/*.json"); files.length()`,
			expected: "1",
		},
		{
			name:     "glob matches all files",
			input:    `let files = glob("` + tempDir + `/*"); files.length()`,
			expected: "4",
		},
		{
			name:     "glob result has correct format",
			input:    `let files = glob("` + tempDir + `/*.json"); files[0].format`,
			expected: "json",
		},
		{
			name:     "glob no matches returns empty array",
			input:    `let files = glob("` + tempDir + `/*.xyz"); files.length()`,
			expected: "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testEvalDirString(tt.input)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

// TestGlobWithDirs tests glob returning both files and directories
func TestGlobWithDirs(t *testing.T) {
	// Create a temp directory
	tempDir := t.TempDir()

	// Create files and subdirectory
	os.WriteFile(filepath.Join(tempDir, "file1.txt"), []byte("content"), 0644)
	os.Mkdir(filepath.Join(tempDir, "subdir"), 0755)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "glob returns both files and dirs",
			input:    `let items = glob("` + tempDir + `/*"); items.length()`,
			expected: "2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testEvalDirString(tt.input)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

// TestDirNonExistent tests behavior with non-existent directories
func TestDirNonExistent(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "non-existent dir exists is false",
			input:    `let d = dir("/nonexistent/path/that/does/not/exist"); d.exists`,
			expected: "false",
		},
		{
			name:     "non-existent dir isDir is false",
			input:    `let d = dir("/nonexistent/path"); d.isDir`,
			expected: "false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testEvalDirString(tt.input)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}
