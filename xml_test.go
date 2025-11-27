package main

import (
	"testing"

	"github.com/sambeau/parsley/pkg/evaluator"
	"github.com/sambeau/parsley/pkg/lexer"
	"github.com/sambeau/parsley/pkg/parser"
)

// TestXMLComments tests that XML comments are properly skipped
func TestXMLComments(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "comment in tag content",
			input:    `<div>hello<!-- this is a comment -->world</div>`,
			expected: "<div>helloworld</div>",
		},
		{
			name:     "comment at start",
			input:    `<!-- comment --><p>text</p>`,
			expected: "<p>text</p>",
		},
		{
			name: "comment with newlines",
			input: `<div><!-- 
multiline
comment
-->content</div>`,
			expected: "<div>content</div>",
		},
		{
			name:     "multiple comments",
			input:    `<div><!-- one -->hello<!-- two -->world<!-- three --></div>`,
			expected: "<div>helloworld</div>",
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
				t.Errorf("Expected %q, got %q", tt.expected, result.Inspect())
			}
		})
	}
}

// TestCDATASections tests CDATA section handling
func TestCDATASections(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic CDATA",
			input:    `<![CDATA[hello world]]>`,
			expected: "hello world",
		},
		{
			name:     "CDATA in tag content",
			input:    `<div><![CDATA[literal <b>text</b>]]></div>`,
			expected: "<div>literal <b>text</b></div>",
		},
		{
			name:     "CDATA with special chars",
			input:    `<![CDATA[<>&"']]>`,
			expected: `<>&"'`,
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
				t.Errorf("Expected %q, got %q", tt.expected, result.Inspect())
			}
		})
	}
}

// TestWebComponentTags tests hyphenated tag names for web components
func TestWebComponentTags(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple web component",
			input:    `<my-component>content</my-component>`,
			expected: "<my-component>content</my-component>",
		},
		{
			name:     "nested web components",
			input:    `<my-app><my-header>Title</my-header></my-app>`,
			expected: "<my-app><my-header>Title</my-header></my-app>",
		},
		{
			name:     "web component with attributes",
			input:    `<custom-element id="test">text</custom-element>`,
			expected: `<custom-element id="test">text</custom-element>`,
		},
		{
			name:     "self-closing web component",
			input:    `<my-icon name="star" />`,
			expected: `<my-icon name="star"  />`,
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
				t.Errorf("Expected %q, got %q", tt.expected, result.Inspect())
			}
		})
	}
}

// TestRawTextTags tests style/script tags with literal {} and @{} interpolation
func TestRawTextTags(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "literal braces in style",
			input:    `<style>body { color: red; }</style>`,
			expected: `<style>body { color: red; }</style>`,
		},
		{
			name:     "interpolation with @{} in style",
			input:    `color = "blue"; <style>.class { color: @{color}; }</style>`,
			expected: `<style>.class { color: blue; }</style>`,
		},
		{
			name:     "multiple rules with literal braces",
			input:    `<style>h1 { font-size: 2em; } p { margin: 10px; }</style>`,
			expected: `<style>h1 { font-size: 2em; } p { margin: 10px; }</style>`,
		},
		{
			name:     "script with literal braces",
			input:    `<script>function test() { return 42; }</script>`,
			expected: `<script>function test() { return 42; }</script>`,
		},
		{
			name:     "script with @{} interpolation",
			input:    `value = 100; <script>var x = @{value};</script>`,
			expected: `<script>var x = 100;</script>`,
		},
		{
			name:     "complex CSS with interpolation",
			input:    `primary = "#007bff"; <style>.btn { background: @{primary}; padding: 10px; }</style>`,
			expected: `<style>.btn { background: #007bff; padding: 10px; }</style>`,
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
				t.Errorf("Expected %q, got %q", tt.expected, result.Inspect())
			}
		})
	}
}

// TestTagBuiltin tests the tag() built-in function for programmatic tag creation
func TestTagBuiltin(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check func(t *testing.T, result evaluator.Object)
	}{
		{
			name:  "tag with name only",
			input: `tag("div")`,
			check: func(t *testing.T, result evaluator.Object) {
				dict, ok := result.(*evaluator.Dictionary)
				if !ok {
					t.Fatalf("Expected Dictionary, got %T", result)
				}
				// Check __type is "tag"
				typeExpr, ok := dict.Pairs["__type"]
				if !ok {
					t.Fatal("Missing __type in tag dictionary")
				}
				typeObj := evaluator.Eval(typeExpr, evaluator.NewEnvironment())
				if typeObj.Inspect() != "tag" {
					t.Errorf("Expected __type='tag', got %s", typeObj.Inspect())
				}
				// Check name is "div"
				nameExpr, ok := dict.Pairs["name"]
				if !ok {
					t.Fatal("Missing name in tag dictionary")
				}
				nameObj := evaluator.Eval(nameExpr, evaluator.NewEnvironment())
				if nameObj.Inspect() != "div" {
					t.Errorf("Expected name='div', got %s", nameObj.Inspect())
				}
			},
		},
		{
			name:  "tag with attributes",
			input: `tag("div", {class: "container", id: "main"})`,
			check: func(t *testing.T, result evaluator.Object) {
				dict, ok := result.(*evaluator.Dictionary)
				if !ok {
					t.Fatalf("Expected Dictionary, got %T", result)
				}
				// Check attrs exists
				attrsExpr, ok := dict.Pairs["attrs"]
				if !ok {
					t.Fatal("Missing attrs in tag dictionary")
				}
				attrsObj := evaluator.Eval(attrsExpr, evaluator.NewEnvironment())
				attrsDict, ok := attrsObj.(*evaluator.Dictionary)
				if !ok {
					t.Fatalf("Expected attrs to be Dictionary, got %T", attrsObj)
				}
				// Check class attribute
				classExpr, ok := attrsDict.Pairs["class"]
				if !ok {
					t.Fatal("Missing class in attrs")
				}
				classObj := evaluator.Eval(classExpr, evaluator.NewEnvironment())
				if classObj.Inspect() != "container" {
					t.Errorf("Expected class='container', got %s", classObj.Inspect())
				}
			},
		},
		{
			name:  "tag with string contents",
			input: `tag("p", {}, "Hello world")`,
			check: func(t *testing.T, result evaluator.Object) {
				dict, ok := result.(*evaluator.Dictionary)
				if !ok {
					t.Fatalf("Expected Dictionary, got %T", result)
				}
				contentsExpr, ok := dict.Pairs["contents"]
				if !ok {
					t.Fatal("Missing contents in tag dictionary")
				}
				contentsObj := evaluator.Eval(contentsExpr, evaluator.NewEnvironment())
				if contentsObj.Inspect() != "Hello world" {
					t.Errorf("Expected contents='Hello world', got %s", contentsObj.Inspect())
				}
			},
		},
		{
			name:  "tag with all parameters",
			input: `tag("a", {href: "/home"}, "Click here")`,
			check: func(t *testing.T, result evaluator.Object) {
				dict, ok := result.(*evaluator.Dictionary)
				if !ok {
					t.Fatalf("Expected Dictionary, got %T", result)
				}
				// Verify name
				nameExpr, _ := dict.Pairs["name"]
				nameObj := evaluator.Eval(nameExpr, evaluator.NewEnvironment())
				if nameObj.Inspect() != "a" {
					t.Errorf("Expected name='a', got %s", nameObj.Inspect())
				}
				// Verify attrs has href
				attrsExpr, _ := dict.Pairs["attrs"]
				attrsObj := evaluator.Eval(attrsExpr, evaluator.NewEnvironment())
				attrsDict := attrsObj.(*evaluator.Dictionary)
				hrefExpr, _ := attrsDict.Pairs["href"]
				hrefObj := evaluator.Eval(hrefExpr, evaluator.NewEnvironment())
				if hrefObj.Inspect() != "/home" {
					t.Errorf("Expected href='/home', got %s", hrefObj.Inspect())
				}
				// Verify contents
				contentsExpr, _ := dict.Pairs["contents"]
				contentsObj := evaluator.Eval(contentsExpr, evaluator.NewEnvironment())
				if contentsObj.Inspect() != "Click here" {
					t.Errorf("Expected contents='Click here', got %s", contentsObj.Inspect())
				}
			},
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

			if errObj, ok := result.(*evaluator.Error); ok {
				t.Fatalf("Evaluation error: %s", errObj.Message)
			}

			tt.check(t, result)
		})
	}
}

// TestTagToString tests converting tag dictionaries back to HTML strings
func TestTagToString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string // Allow multiple valid outputs for ordering
	}{
		{
			name:     "simple tag",
			input:    `toString(tag("div", {}, "Hello"))`,
			expected: []string{`<div>Hello</div>`},
		},
		{
			name:     "self-closing tag",
			input:    `toString(tag("br"))`,
			expected: []string{`<br />`},
		},
		{
			name:     "tag with attributes",
			input:    `toString(tag("a", {href: "/home"}, "Link"))`,
			expected: []string{`<a href="/home">Link</a>`},
		},
		{
			name:  "tag with multiple attributes",
			input: `toString(tag("img", {src: "test.png", alt: "Test"}))`,
			// Dictionary key order is not guaranteed
			expected: []string{
				`<img src="test.png" alt="Test" />`,
				`<img alt="Test" src="test.png" />`,
			},
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

			if errObj, ok := result.(*evaluator.Error); ok {
				t.Fatalf("Evaluation error: %s", errObj.Message)
			}

			str, ok := result.(*evaluator.String)
			if !ok {
				t.Fatalf("Expected String, got %T", result)
			}

			// Check if result matches any of the expected values
			found := false
			for _, exp := range tt.expected {
				if str.Value == exp {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected one of %v, got %q", tt.expected, str.Value)
			}
		})
	}
}

// TestProcessingInstructions tests <?xml ... ?> handling
func TestProcessingInstructions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "xml declaration",
			input:    `<?xml version="1.0" encoding="UTF-8"?>`,
			expected: `<?xml version="1.0" encoding="UTF-8"?>`,
		},
		{
			name:     "xml declaration concatenated with html",
			input:    `<?xml version="1.0"?> + <html><body>content</body></html>`,
			expected: `<?xml version="1.0"?><html><body>content</body></html>`,
		},
		{
			name:     "stylesheet processing instruction",
			input:    `<?xml-stylesheet type="text/xsl" href="style.xsl"?>`,
			expected: `<?xml-stylesheet type="text/xsl" href="style.xsl"?>`,
		},
		{
			name:     "php processing instruction",
			input:    `<?php echo "hello"; ?>`,
			expected: `<?php echo "hello"; ?>`,
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
				t.Errorf("Expected %q, got %q", tt.expected, result.Inspect())
			}
		})
	}
}

// TestDoctypeDeclarations tests <!DOCTYPE ...> handling
func TestDoctypeDeclarations(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "html5 doctype",
			input:    `<!DOCTYPE html>`,
			expected: `<!DOCTYPE html>`,
		},
		{
			name:     "doctype concatenated with html",
			input:    `<!DOCTYPE html> + <html><head></head></html>`,
			expected: `<!DOCTYPE html><html><head></head></html>`,
		},
		{
			name:     "xhtml doctype",
			input:    `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">`,
			expected: `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">`,
		},
		{
			name:     "svg doctype",
			input:    `<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">`,
			expected: `<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">`,
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
				t.Errorf("Expected %q, got %q", tt.expected, result.Inspect())
			}
		})
	}
}
