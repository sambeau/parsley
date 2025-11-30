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
func testEvalSVGWithFilename(input string, filename string) evaluator.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := evaluator.NewEnvironment()
	env.Filename = filename
	// Enable write operations for SVG tests (reads are allowed by default)
	env.Security = &evaluator.SecurityPolicy{
		AllowWriteAll: true,
	}
	return evaluator.Eval(program, env)
}

// TestSVGFormat tests the SVG file format for reading SVG files
func TestSVGFormat(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "parsley-svg-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test SVG file with XML prolog and DOCTYPE
	svgContent := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
  <circle cx="12" cy="12" r="10"/>
</svg>`
	svgPath := filepath.Join(tmpDir, "test.svg")
	if err := os.WriteFile(svgPath, []byte(svgContent), 0644); err != nil {
		t.Fatalf("Failed to write test SVG: %v", err)
	}

	// Test file path - used for relative path resolution
	testFilePath := filepath.Join(tmpDir, "test.pars")

	tests := []struct {
		name     string
		input    string
		contains string
	}{
		{
			name:     "SVG stripping XML prolog",
			input:    `let Icon <== SVG(@./test.svg); Icon`,
			contains: `<svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">`,
		},
		{
			name:     "SVG as custom component",
			input:    `let MyIcon <== SVG(@./test.svg)` + "\n" + `<div><MyIcon/></div>`,
			contains: `<div><svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testEvalSVGWithFilename(tt.input, testFilePath)
			if result == nil {
				t.Fatalf("Result is nil")
			}

			resultStr := result.Inspect()

			if !strings.Contains(resultStr, tt.contains) {
				t.Errorf("Expected result to contain:\n%s\n\nGot:\n%s", tt.contains, resultStr)
			}

			// Ensure XML prolog is stripped
			if strings.Contains(resultStr, "<?xml") {
				t.Errorf("XML prolog should be stripped, but found in: %s", resultStr)
			}
			if strings.Contains(resultStr, "<!DOCTYPE") {
				t.Errorf("DOCTYPE should be stripped, but found in: %s", resultStr)
			}
		})
	}
}

// TestSVGStripXMLProlog tests the XML prolog stripping
func TestSVGStripXMLProlog(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Strip XML prolog",
			input:    "<?xml version=\"1.0\"?>\n<svg></svg>",
			expected: "<svg></svg>",
		},
		{
			name:     "Strip XML prolog with encoding",
			input:    "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<svg></svg>",
			expected: "<svg></svg>",
		},
		{
			name:     "Strip DOCTYPE",
			input:    "<!DOCTYPE svg PUBLIC \"-//W3C//DTD SVG 1.1//EN\" \"...\">\n<svg></svg>",
			expected: "<svg></svg>",
		},
		{
			name:     "Strip both prolog and DOCTYPE",
			input:    "<?xml version=\"1.0\"?>\n<!DOCTYPE svg PUBLIC \"...\" \"...\">\n<svg></svg>",
			expected: "<svg></svg>",
		},
		{
			name:     "No prolog - unchanged",
			input:    "<svg></svg>",
			expected: "<svg></svg>",
		},
		{
			name:     "Only whitespace around svg",
			input:    "  \n  <svg></svg>  \n  ",
			expected: "<svg></svg>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temp file with the input SVG
			tmpDir, err := os.MkdirTemp("", "parsley-svg-strip-test")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			svgPath := filepath.Join(tmpDir, "test.svg")
			if err := os.WriteFile(svgPath, []byte(tt.input), 0644); err != nil {
				t.Fatalf("Failed to write test SVG: %v", err)
			}

			// Test file path for relative path resolution
			testFilePath := filepath.Join(tmpDir, "test.pars")

			// Create test .pars code
			parsCode := `let Icon <== SVG(@./test.svg); Icon`

			result := testEvalSVGWithFilename(parsCode, testFilePath)
			if result == nil {
				t.Fatalf("Result is nil")
			}

			resultStr := strings.TrimSpace(result.Inspect())

			if resultStr != tt.expected {
				t.Errorf("Expected:\n%s\n\nGot:\n%s", tt.expected, resultStr)
			}
		})
	}
}

// TestSVGWrite tests writing SVG content to files
func TestSVGWrite(t *testing.T) {
	// Create a temp directory for test files
	tmpDir, err := os.MkdirTemp("", "parsley-svg-write-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test file path for relative path resolution
	testFilePath := filepath.Join(tmpDir, "test.pars")

	// Test writing SVG content (correct syntax: value ==> file handle)
	code := `let icon = "<svg><circle/></svg>"; icon ==> SVG(@./output.svg); "done"`
	result := testEvalSVGWithFilename(code, testFilePath)
	if result == nil {
		t.Fatalf("Result is nil")
	}

	// Check if there's an error in the result
	if result.Type() == evaluator.ERROR_OBJ {
		t.Fatalf("Evaluation error: %s", result.Inspect())
	}

	// Verify the file was written
	outputPath := filepath.Join(tmpDir, "output.svg")
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	expected := "<svg><circle/></svg>"
	if string(content) != expected {
		t.Errorf("Expected file content:\n%s\n\nGot:\n%s", expected, string(content))
	}
}
