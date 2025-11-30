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
func testEvalMDWithFilename(input string, filename string) evaluator.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := evaluator.NewEnvironment()
	env.Filename = filename
	// Reads are allowed by default, no special security needed
	return evaluator.Eval(program, env)
}

// TestMarkdownBasic tests basic markdown parsing
func TestMarkdownBasic(t *testing.T) {
	// Create a temp directory for test files
	tmpDir, err := os.MkdirTemp("", "parsley-md-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a simple markdown file without frontmatter
	mdContent := `# Hello World

This is a paragraph.

- Item 1
- Item 2
`
	mdPath := filepath.Join(tmpDir, "simple.md")
	if err := os.WriteFile(mdPath, []byte(mdContent), 0644); err != nil {
		t.Fatalf("Failed to write markdown file: %v", err)
	}

	// Test file path for relative path resolution
	testFilePath := filepath.Join(tmpDir, "test.pars")

	// Test reading markdown
	code := `let post <== MD(@./simple.md); post.html`
	result := testEvalMDWithFilename(code, testFilePath)

	if result == nil {
		t.Fatalf("Result is nil")
	}

	if result.Type() == evaluator.ERROR_OBJ {
		t.Fatalf("Evaluation error: %s", result.Inspect())
	}

	html := result.Inspect()
	if !strings.Contains(html, "<h1>Hello World</h1>") {
		t.Errorf("Expected h1 tag, got: %s", html)
	}
	if !strings.Contains(html, "<li>Item 1</li>") {
		t.Errorf("Expected list items, got: %s", html)
	}
}

// TestMarkdownWithFrontmatter tests markdown with YAML frontmatter
func TestMarkdownWithFrontmatter(t *testing.T) {
	// Create a temp directory for test files
	tmpDir, err := os.MkdirTemp("", "parsley-md-frontmatter-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a markdown file with frontmatter
	mdContent := `---
title: My Blog Post
author: John Doe
tags:
  - go
  - parsley
draft: false
---
# Content

This is the blog content.
`
	mdPath := filepath.Join(tmpDir, "blog.md")
	if err := os.WriteFile(mdPath, []byte(mdContent), 0644); err != nil {
		t.Fatalf("Failed to write markdown file: %v", err)
	}

	testFilePath := filepath.Join(tmpDir, "test.pars")

	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name:     "access title",
			code:     `let post <== MD(@./blog.md); post.title`,
			expected: "My Blog Post",
		},
		{
			name:     "access author",
			code:     `let post <== MD(@./blog.md); post.author`,
			expected: "John Doe",
		},
		{
			name:     "access draft",
			code:     `let post <== MD(@./blog.md); post.draft`,
			expected: "false",
		},
		{
			name:     "access tags array",
			code:     `let post <== MD(@./blog.md); post.tags[0]`,
			expected: "go",
		},
		{
			name:     "html contains content",
			code:     `let post <== MD(@./blog.md); post.html`,
			expected: "<h1>Content</h1>",
		},
		{
			name:     "raw contains markdown",
			code:     `let post <== MD(@./blog.md); post.raw`,
			expected: "# Content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testEvalMDWithFilename(tt.code, testFilePath)

			if result == nil {
				t.Fatalf("Result is nil")
			}

			if result.Type() == evaluator.ERROR_OBJ {
				t.Fatalf("Evaluation error: %s", result.Inspect())
			}

			if !strings.Contains(result.Inspect(), tt.expected) {
				t.Errorf("Expected to contain '%s', got: %s", tt.expected, result.Inspect())
			}
		})
	}
}

// TestMarkdownAsComponent tests using markdown in templates
func TestMarkdownAsComponent(t *testing.T) {
	// Create a temp directory for test files
	tmpDir, err := os.MkdirTemp("", "parsley-md-component-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a markdown file with frontmatter
	mdContent := `---
title: Hello World
---
This is **bold** text.
`
	mdPath := filepath.Join(tmpDir, "post.md")
	if err := os.WriteFile(mdPath, []byte(mdContent), 0644); err != nil {
		t.Fatalf("Failed to write markdown file: %v", err)
	}

	testFilePath := filepath.Join(tmpDir, "test.pars")

	// Test using markdown in a template
	code := `let post <== MD(@./post.md)
<article>
  <h1>{post.title}</h1>
  <div class="content">{post.html}</div>
</article>`

	result := testEvalMDWithFilename(code, testFilePath)

	if result == nil {
		t.Fatalf("Result is nil")
	}

	if result.Type() == evaluator.ERROR_OBJ {
		t.Fatalf("Evaluation error: %s", result.Inspect())
	}

	html := result.Inspect()
	if !strings.Contains(html, "<h1>Hello World</h1>") {
		t.Errorf("Expected title in h1, got: %s", html)
	}
	if !strings.Contains(html, "<strong>bold</strong>") {
		t.Errorf("Expected bold text, got: %s", html)
	}
}
