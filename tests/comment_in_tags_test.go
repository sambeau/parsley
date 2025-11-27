package main

import (
	"testing"

	"github.com/sambeau/parsley/pkg/evaluator"
	"github.com/sambeau/parsley/pkg/lexer"
	"github.com/sambeau/parsley/pkg/parser"
)

// TestCommentsInAndAroundTags tests that Parsley comments are properly skipped
// in various contexts around and inside tags
func TestCommentsInAndAroundTags(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "comment between tags in tag content",
			input: `<div>
    <p>First</p>
    // This is a comment
    <p>Second</p>
</div>`,
			expected: "<div><p>First</p> \n    <p>Second</p></div>",
		},
		{
			name: "comment with @{} between tags",
			input: `<div>
    <p>First</p>
    // This comment has @{} in it
    <p>Second</p>
</div>`,
			expected: "<div><p>First</p> \n    <p>Second</p></div>",
		},
		{
			name: "comment inside style tag",
			input: `<style>
// Comment inside style
body {
    color: red;
}
</style>`,
			expected: "<style>\n// Comment inside style\nbody {\n    color: red;\n}\n</style>",
		},
		{
			name: "comment with @{} inside style tag",
			input: `variable = "test"
<style>
// This has @{variable} in the comment
body {
    color: red;
}
</style>`,
			expected: "<style>\n// This has test in the comment\nbody {\n    color: red;\n}\n</style>",
		},
		{
			name: "comment between title and style in nested tags",
			input: `let Page = fn() {
    <html>
        <head>
            <title>Test</title>
            // only @{} sections are interpolated
            <style>
                body { color: red; }
            </style>
        </head>
    </html>
}
<Page></Page>`,
			expected: "<html><head><title>Test</title> \n            <style>\n                body { color: red; }\n            </style></head></html>",
		},
		{
			name: "multiple comments in style tag",
			input: `<style>
// First comment
body {
    color: red;
}
// Second comment
h1 {
    font-size: 2em;
}
</style>`,
			expected: "<style>\n// First comment\nbody {\n    color: red;\n}\n// Second comment\nh1 {\n    font-size: 2em;\n}\n</style>",
		},
		{
			name: "comment inside script tag",
			input: `<script>
// JavaScript comment
function test() {
    console.log("hello");
}
</script>`,
			expected: "<script>\n// JavaScript comment\nfunction test() {\n    console.log(\"hello\");\n}\n</script>",
		},
		{
			name: "interpolation in script comment for datestamp",
			input: `version = "1.0.0"
<script>
// Version: @{version}
function getVersion() {
    return "1.0.0";
}
</script>`,
			expected: "<script>\n// Version: 1.0.0\nfunction getVersion() {\n    return \"1.0.0\";\n}\n</script>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) != 0 {
				t.Fatalf("Parser errors: %v", p.Errors())
			}

			env := evaluator.NewEnvironment()
			result := evaluator.Eval(program, env)

			if result == nil {
				t.Fatalf("Eval returned nil")
			}

			if result.Inspect() != tt.expected {
				t.Errorf("Expected:\n%s\n\nGot:\n%s", tt.expected, result.Inspect())
			}
		})
	}
}
