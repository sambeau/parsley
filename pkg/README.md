# Parsley Language Packages

This directory contains the core packages that implement the Parsley programming language interpreter.

## Package Overview

### `ast/` - Abstract Syntax Tree
Defines all AST node types that represent the parsed structure of Parsley programs.

**Key Types:**
- Expression nodes: identifiers, literals, operators
- Statement nodes: let, return, assignment
- Function literals and calls
- Array and dictionary literals
- Tag expressions (singleton and paired)
- Destructuring patterns (array and dictionary)

### `evaluator/` - Program Evaluator
Evaluates the AST and executes Parsley programs.

**Features:**
- Expression evaluation
- Environment management for variable scoping
- Built-in functions (len, first, last, rest, push, map, filter, etc.)
- Tag evaluation (HTML/SVG generation)
- Template literal interpolation
- Destructuring assignment support
- Error handling with position information

### `formatter/` - Output Formatting
Formats program output for better readability.

**Capabilities:**
- HTML pretty-printing with auto-detection
- Smart indentation (2 spaces per level)
- Inline vs block element handling
- CSS/JavaScript content preservation

### `lexer/` - Lexical Analysis
Tokenizes Parsley source code into a stream of tokens.

**Handles:**
- Keywords and identifiers
- Literals (strings, numbers, booleans, null)
- Operators and delimiters
- Template literals with interpolation
- Tag syntax (singleton and paired)
- Multi-line strings
- Escape sequences

### `parser/` - Syntax Analysis
Converts token stream into an Abstract Syntax Tree using Pratt parsing.

**Supports:**
- Operator precedence
- Expression parsing
- Statement parsing
- Function definitions
- Array and dictionary literals
- Tag expressions
- Destructuring patterns
- Template literals

### `repl/` - Read-Eval-Print Loop
Interactive shell for executing Parsley code.

**Features:**
- Multi-line input support
- Immediate expression evaluation
- Error reporting
- Version display

## Testing

Each package includes comprehensive test files:
- Unit tests for individual functions
- Integration tests for end-to-end scenarios
- Table-driven tests for thorough coverage

Run all tests:
```bash
go test ./...
```

Run tests for a specific package:
```bash
go test ./pkg/lexer
go test ./pkg/parser
go test ./pkg/evaluator
```

## Usage

These packages are designed to work together to implement the Parsley interpreter:

```go
import (
    "github.com/sambeau/parsley/pkg/lexer"
    "github.com/sambeau/parsley/pkg/parser"
    "github.com/sambeau/parsley/pkg/evaluator"
)

// Tokenize input
l := lexer.New(input)

// Parse into AST
p := parser.New(l)
program := p.ParseProgram()

// Evaluate
env := evaluator.NewEnvironment()
result := evaluator.Eval(program, env)
```

## Code Organization

- Each package is self-contained with minimal dependencies
- All packages follow standard Go conventions
- Public APIs are well-documented
- Error handling includes line/column information
- Test coverage ensures reliability
