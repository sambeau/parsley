package main

import (
	"strings"
	"testing"

	"github.com/sambeau/parsley/pkg/evaluator"
	"github.com/sambeau/parsley/pkg/lexer"
	"github.com/sambeau/parsley/pkg/parser"
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

func TestBasicTagPairs(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple div tag",
			input:    `<div>Hello, World!</div>`,
			expected: "<div>Hello, World!</div>",
		},
		{
			name:     "paragraph tag",
			input:    `<p>This is a paragraph.</p>`,
			expected: "<p>This is a paragraph.</p>",
		},
		{
			name:     "empty tag",
			input:    `<div></div>`,
			expected: "<div></div>",
		},
		{
			name:     "tag with trailing space",
			input:    `<div>Hello </div>`,
			expected: "<div>Hello </div>",
		},
		{
			name:     "tag with leading space",
			input:    `<div> Hello</div>`,
			expected: "<div> Hello</div>",
		},
		{
			name:     "tag with multiple spaces",
			input:    `<div>Hello   World</div>`,
			expected: "<div>Hello   World</div>",
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

			if result.Inspect() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Inspect())
			}
		})
	}
}

func TestNestedTags(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple nesting",
			input:    `<div><p>Nested</p></div>`,
			expected: "<div><p>Nested</p></div>",
		},
		{
			name:     "multiple nested tags",
			input:    `<div><h1>Title</h1><p>Content</p></div>`,
			expected: "<div><h1>Title</h1><p>Content</p></div>",
		},
		{
			name:     "deeply nested tags",
			input:    `<div><section><article><p>Deep</p></article></section></div>`,
			expected: "<div><section><article><p>Deep</p></article></section></div>",
		},
		{
			name:     "nested with text between",
			input:    `<div>Before<p>Middle</p>After</div>`,
			expected: "<div>Before<p>Middle</p>After</div>",
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

			if result.Inspect() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Inspect())
			}
		})
	}
}

func TestTagsWithInterpolation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple interpolation",
			input:    "name = \"World\"\n<div>Hello, {name}!</div>",
			expected: "<div>Hello, World!</div>",
		},
		{
			name:     "multiple interpolations",
			input:    "x = \"A\"\ny = \"B\"\n<div>{x} and {y}</div>",
			expected: "<div>A and B</div>",
		},
		{
			name:     "interpolation with spaces",
			input:    "x = \"FIRST\"\ny = \"SECOND\"\n<div>{x} - {y}</div>",
			expected: "<div>FIRST - SECOND</div>",
		},
		{
			name:     "interpolation at start",
			input:    "name = \"Start\"\n<div>{name} here</div>",
			expected: "<div>Start here</div>",
		},
		{
			name:     "interpolation at end",
			input:    "name = \"End\"\n<div>Here is {name}</div>",
			expected: "<div>Here is End</div>",
		},
		{
			name:     "only interpolation",
			input:    "name = \"Only\"\n<div>{name}</div>",
			expected: "<div>Only</div>",
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

			if result.Inspect() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Inspect())
			}
		})
	}
}

func TestEmptyGroupingTags(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple grouping",
			input:    `<>Hello</>`,
			expected: "Hello",
		},
		{
			name:     "grouping with nested tags",
			input:    `<><div>First</div><div>Second</div></>`,
			expected: "<div>First</div><div>Second</div>",
		},
		{
			name:     "grouping with interpolation",
			input:    "x = \"Test\"\n<>{x}</>",
			expected: "Test",
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

			if result.Inspect() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Inspect())
			}
		})
	}
}

func TestComponentsWithContents(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "basic component with contents",
			input: `Title = fn(props) {
				<title>{props.contents}</title>
			}
			<Title>Home Page</Title>`,
			expected: "<title>Home Page</title>",
		},
		{
			name: "component with contents and interpolation",
			input: `SiteName = "MyGroovySite"
			Title = fn(props) {
				<title>{props.contents} - {SiteName}</title>
			}
			<Title>Home Page</Title>`,
			expected: "<title>Home Page - MyGroovySite</title>",
		},
		{
			name: "component with nested tags in contents",
			input: `Card = fn(props) {
				<div>{props.contents}</div>
			}
			<Card><h2>Title</h2><p>Body</p></Card>`,
			expected: "<div><h2>Title</h2><p>Body</p></div>",
		},
		{
			name: "component with interpolation in contents",
			input: `name = "Alice"
			Wrapper = fn(props) {
				<div>{props.contents}</div>
			}
			<Wrapper>Hello, {name}!</Wrapper>`,
			expected: "<div>Hello, Alice!</div>",
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

			if result.Inspect() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Inspect())
			}
		})
	}
}

func TestComponentsWithProps(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "component with single prop",
			input: `Greeting = fn(props) {
				<h1>Hello, {props.name}!</h1>
			}
			<Greeting name="World" />`,
			expected: "<h1>Hello, World!</h1>",
		},
		{
			name: "component with multiple props",
			input: `Person = fn(props) {
				<div>{props.name} is {props.age} years old</div>
			}
			<Person name="Alice" age="30" />`,
			expected: "<div>Alice is 30 years old</div>",
		},
		{
			name: "component used in map",
			input: `Welcome = fn(name) {
				<h2>Hello, {name}</h2>
			}
			names = ["Sara", "Cahal", "Edite"]
			result = map(Welcome, names)
			<div>{result}</div>`,
			expected: "<div><h2>Hello, Sara</h2><h2>Hello, Cahal</h2><h2>Hello, Edite</h2></div>",
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

			if result.Inspect() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Inspect())
			}
		})
	}
}

func TestWhitespacePreservation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "preserve space after comma",
			input:    "name = \"Sara\"\n<h2>Hello, {name}</h2>",
			expected: "<h2>Hello, Sara</h2>",
		},
		{
			name:     "preserve spaces around dash",
			input:    "x = \"A\"\ny = \"B\"\n<div>{x} - {y}</div>",
			expected: "<div>A - B</div>",
		},
		{
			name:     "preserve trailing space before interpolation",
			input:    "name = \"World\"\n<div>Hello {name}</div>",
			expected: "<div>Hello World</div>",
		},
		{
			name:     "preserve space after interpolation",
			input:    "name = \"Alice\"\n<div>{name} here</div>",
			expected: "<div>Alice here</div>",
		},
		{
			name:     "preserve multiple spaces",
			input:    `<div>A    B</div>`,
			expected: "<div>A    B</div>",
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

			if result.Inspect() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Inspect())
			}
		})
	}
}

func TestComplexTagExamples(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "component composition",
			input: `Card = fn(props) {
				<div><h3>{props.title}</h3><p>{props.body}</p></div>
			}
			<Card title="Welcome" body="This is the content" />`,
			expected: "<div><h3>Welcome</h3><p>This is the content</p></div>",
		},
		{
			name: "nested components with contents",
			input: `Inner = fn(props) { <span>{props.contents}</span> }
			Outer = fn(props) { <div>{props.contents}</div> }
			<Outer><Inner>Hello</Inner></Outer>`,
			expected: "<div><span>Hello</span></div>",
		},
		{
			name: "tag with expressions",
			input: `x = 5
			y = 10
			<div>The sum is {x + y}</div>`,
			expected: "<div>The sum is 15</div>",
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

			if result.Inspect() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result.Inspect())
			}
		})
	}
}
