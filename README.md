# Pars

A Go-based programming language interpreter with comprehensive mathematical and trigonometric function support.

## Features

### Core Language Features
- Variable declarations with `let`
- Direct variable assignment (e.g., `x = 5`)
- Functions with `fn`
- If-then-else expressions (ternary-style conditionals)
- Integer and floating-point arithmetic
- String operations
- Boolean logic

### Built-in Mathematical Functions
- **Trigonometric Functions:**
  - `sin(x)` - Sine function
  - `cos(x)` - Cosine function
  - `tan(x)` - Tangent function
  - `asin(x)` - Arcsine function
  - `acos(x)` - Arccosine function
  - `atan(x)` - Arctangent function

- **Mathematical Functions:**
  - `sqrt(x)` - Square root
  - `pow(base, exp)` - Power function
  - `pi()` - Returns the value of π

### Data Types
- **Integers:** `42`, `-15`
- **Floats:** `3.14159`, `2.718`
- **Strings:** `"hello world"`
- **Booleans:** `true`, `false`

## Getting Started

### Prerequisites

- Go 1.19 or higher

### Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd pars
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

### Running the Project

To start the interactive REPL:

```bash
go run main.go
```

To build the project:

```bash
go build -o pars
./pars
```

## Usage Examples

### Basic Arithmetic
```
>> 2 + 3
5
>> 10 * 4.5
45
```

### Trigonometric Functions
```
>> sin(0)
0
>> cos(0)
1
>> sin(pi() / 2)
1
>> tan(pi() / 4)
1
```

### Mathematical Operations
```
>> sqrt(16)
4
>> pow(2, 8)
256
>> pi()
3.141592653589793
```

### Variable Assignment
```
>> x = 5
5
>> y = sin(x)
-0.9589242746631385
>> z = x + y
4.041075725336861
```

### Variable Updates
```
>> radius = 3
3
>> area = pi() * pow(radius, 2)
28.274333882308138
>> radius = 5
5
>> area = pi() * pow(radius, 2)
78.53981633974483
```

### Using Variables in Complex Expressions
```
>> a = 3
>> b = 4  
>> c = sqrt(pow(a, 2) + pow(b, 2))
5
>> angle = atan(b / a)
0.9272952180016122
```

### If-Then-Else Expressions
```
>> x = if 5 > 0 then true else false
true
>> y = if 1 < 0 then 0
null
>> a = 10
>> result = if a > 5 then "big" else "small"
big
>> bar = 15
>> foo = if bar * 20 > 100 then 100 else bar
100
>> nested = if 1 > 0 then if 2 > 1 then 3 else 4 else 5
3
```

If-then-else expressions work like ternary operators and can be used anywhere an expression is expected. The `else` clause is optional - if omitted and the condition is false, the expression evaluates to `null`.

### Functions
```
>> let circleArea = fn(r) { pi() * pow(r, 2) }
>> circleArea(10)
314.1592653589793
```

## Development

This project is structured to easily add new modules and packages:

- `main.go` - Entry point of the application
- `pkg/` - Public packages that can be used by external applications
  - `lexer/` - Tokenizes input into lexical tokens
  - `parser/` - Converts tokens into an Abstract Syntax Tree
  - `ast/` - Defines the Abstract Syntax Tree nodes
  - `evaluator/` - Evaluates the AST and executes the program
  - `repl/` - Read-Eval-Print Loop for interactive usage
- `internal/` - Private packages that are specific to this application
- `cmd/` - Command-line applications
- `api/` - API definitions (OpenAPI/Swagger specs, protocol definition files)
- `configs/` - Configuration file templates or default configs

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Testing

Run all tests:
```bash
go test ./...
```

Run specific package tests:
```bash
go test ./pkg/lexer -v
go test ./pkg/parser -v
go test ./pkg/evaluator -v
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
