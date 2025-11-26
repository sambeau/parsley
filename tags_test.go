package main

import (
	"strings"
	"testing"

	"pars/pkg/evaluator"
	"pars/pkg/lexer"
	"pars/pkg/parser"
)

func TestStandardTags(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple br tag",
			input:    `<br/>`,
			expected: `<br />`,
		},
		{
			name:     "meta tag with attribute",
			input:    `<meta charset="utf-8" />`,
			expected: `<meta charset="utf-8"  />`,
		},
		{
			name:     "input with multiple attributes",
			input:    `<input type="text" name="username" />`,
			expected: `<input type="text" name="username"  />`,
		},
		{
			name:     "tag with boolean attribute",
			input:    `<input type="checkbox" checked />`,
			expected: `<input type="checkbox" checked  />`,
		},
		{
			name:     "multiline tag",
			input:    `<img src="test.png" width="100" height="200" />`,
			expected: `<img src="test.png" width="100" height="200"  />`,
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

func TestTagInterpolation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "interpolate string variable",
			input: `charset = "utf-8"
<meta charset="{charset}" />`,
			expected: `<meta charset="utf-8"  />`,
		},
		{
			name: "interpolate number variable",
			input: `width = 300
height = 200
<img width="{width}" height="{height}" />`,
			expected: `<img width="300" height="200"  />`,
		},
		{
			name: "interpolate expression",
			input: `x = 10
<div data-value="{x * 2}" />`,
			expected: `<div data-value="20"  />`,
		},
		{
			name: "conditional interpolation",
			input: `disabled = true
<button disabled="{if(disabled){"disabled"}}" />`,
			expected: `<button disabled="disabled"  />`,
		},
		{
			name: "conditional interpolation false",
			input: `disabled = false
<button disabled="{if(disabled){"disabled"}}" />`,
			expected: `<button disabled=""  />`,
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
			var result evaluator.Object
			for _, stmt := range program.Statements {
				result = evaluator.Eval(stmt, env)
			}

			if result == nil {
				t.Fatalf("Eval returned nil")
			}

			if result.Inspect() != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result.Inspect())
			}
		})
	}
}

func TestCustomTags(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "custom tag with string props",
			input: `Dog = fn(props) {
  name = props.name
  age = props.age
  toString("Dog: ", name, ", Age: ", age)
}
<Dog name="Rover" age="5" />`,
			expected: `Dog: Rover, Age: 5`,
		},
		{
			name: "custom tag with interpolated props",
			input: `Dog = fn(props) {
  name = props.name
  age = props.age
  toString("Dog: ", name, ", Age: ", age)
}
dogAge = 7
<Dog name="Max" age="{dogAge}" />`,
			expected: `Dog: Max, Age: 7`,
		},
		{
			name: "custom tag returning template",
			input: `Link = fn(props) {
  url = props.url
  text = props.text
  toString("<a href=\"", url, "\">", text, "</a>")
}
<Link url="https://example.com" text="Click here" />`,
			expected: `<a href="https://example.com">Click here</a>`,
		},
		{
			name: "custom tag with expression prop",
			input: `Double = fn(props) {
  value = props.value
  value * 2
}
<Double value="{10 + 5}" />`,
			expected: `30`,
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
			var result evaluator.Object
			for _, stmt := range program.Statements {
				result = evaluator.Eval(stmt, env)
			}

			if result == nil {
				t.Fatalf("Eval returned nil")
			}

			if result.Inspect() != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result.Inspect())
			}
		})
	}
}

func TestCustomTagsWithBooleanProps(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "custom tag with boolean prop",
			input: `Button = fn(props) {
  isDisabled = has(props, "disabled")
  if (isDisabled) {
    "disabled"
  } else {
    "enabled"
  }
}
<Button disabled />`,
			expected: `disabled`,
		},
		{
			name: "custom tag without boolean prop",
			input: `Button = fn(props) {
  isDisabled = has(props, "disabled")
  if (isDisabled) {
    "disabled"
  } else {
    "enabled"
  }
}
<Button type="submit" />`,
			expected: `enabled`,
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
			var result evaluator.Object
			for _, stmt := range program.Statements {
				result = evaluator.Eval(stmt, env)
			}

			if result == nil {
				t.Fatalf("Eval returned nil")
			}

			if result.Inspect() != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result.Inspect())
			}
		})
	}
}

func TestTagsInExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "tag in array",
			input: `tags = [<br/>, <hr/>]
tags[0]`,
			expected: `<br />`,
		},
		{
			name: "tag in variable assignment",
			input: `x = <input type="text" />
x`,
			expected: `<input type="text"  />`,
		},
		{
			name:     "multiple tags concatenated",
			input:    `toString(<br/>, <hr/>)`,
			expected: `<br /><hr />`,
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
			var result evaluator.Object
			for _, stmt := range program.Statements {
				result = evaluator.Eval(stmt, env)
			}

			if result == nil {
				t.Fatalf("Eval returned nil")
			}

			if result.Inspect() != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result.Inspect())
			}
		})
	}
}

func TestTagErrors(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "custom tag function not found",
			input:       `<NonExistent name="test" />`,
			expectError: true,
		},
		{
			name: "unclosed brace in interpolation",
			input: `x = 5
<div data="{x" />`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) != 0 {
				if !tt.expectError {
					t.Fatalf("Unexpected parser errors: %v", p.Errors())
				}
				return
			}

			env := evaluator.NewEnvironment()
			var result evaluator.Object
			for _, stmt := range program.Statements {
				result = evaluator.Eval(stmt, env)
			}

			if tt.expectError {
				if result == nil {
					t.Fatal("Expected error but got nil")
				}
				if result.Type() != evaluator.ERROR_OBJ {
					t.Errorf("Expected error object, got %s", result.Type())
				}
			}
		})
	}
}

func TestLexerTagTokenization(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []lexer.TokenType
	}{
		{
			name:  "simple tag",
			input: `<br/>`,
			expected: []lexer.TokenType{
				lexer.TAG,
				lexer.EOF,
			},
		},
		{
			name:  "tag with attributes",
			input: `<input type="text" />`,
			expected: []lexer.TokenType{
				lexer.TAG,
				lexer.EOF,
			},
		},
		{
			name:  "tag not confused with less than",
			input: `x < 5`,
			expected: []lexer.TokenType{
				lexer.IDENT,
				lexer.LT,
				lexer.INT,
				lexer.EOF,
			},
		},
		{
			name:  "tag with variable",
			input: `x = <br/>`,
			expected: []lexer.TokenType{
				lexer.IDENT,
				lexer.ASSIGN,
				lexer.TAG,
				lexer.EOF,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			tokens := []lexer.TokenType{}

			for {
				tok := l.NextToken()
				tokens = append(tokens, tok.Type)
				if tok.Type == lexer.EOF {
					break
				}
			}

			if len(tokens) != len(tt.expected) {
				t.Fatalf("Expected %d tokens, got %d", len(tt.expected), len(tokens))
			}

			for i, expectedType := range tt.expected {
				if tokens[i] != expectedType {
					t.Errorf("Token %d: expected %s, got %s", i, expectedType, tokens[i])
				}
			}
		})
	}
}

func TestMultilineTagsPreserveWhitespace(t *testing.T) {
	input := `<img
  src="test.png"
  width="100"
  height="200" />`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	env := evaluator.NewEnvironment()
	result := evaluator.Eval(program, env)

	if result == nil {
		t.Fatal("Eval returned nil")
	}

	// Check that newlines are preserved
	if !strings.Contains(result.Inspect(), "\n") {
		t.Error("Expected multiline tag to preserve newlines")
	}
}
