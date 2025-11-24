# Pars

A Go-based toy concatenative programming language interpreter.

## Features

### Core Language Features
- Variable declarations with `let`
- Direct variable assignment (e.g., `x = 5`)
- Functions with `fn`
- If-then-else expressions (ternary-style conditionals)
- Arrays with comma separator
- Array indexing and slicing with `[]`
- Array concatenation with `++`
- String indexing and slicing
- String concatenation with `+`
- String escape sequences (`\n`, `\t`, etc.)
- Template literals with `${}` interpolation
- Integer and floating-point arithmetic
- Boolean logic

### Data Types
- **Integers:** `42`, `-15`
- **Floats:** `3.14159`, `2.718`
- **Strings:** `"hello world"`
- **Booleans:** `true`, `false`
- **Arrays:** `1,2,3`, `"Sam","Phillips"`, mixed types allowed

### Built-in Functions

- **String Functions:**
  - `toUpper(str)` - Convert string to uppercase
  - `toLower(str)` - Convert string to lowercase
  - `len(str)` - Get the length of a string

- **Array Functions:**
  - `map(func, elements...)` - Apply function to each element, filter out nulls
  - `for(array) func` - Sugar syntax for map with function
  - `for(var in array) { body }` - Sugar syntax for map with inline function
  - `len(array)` - Get the length of an array

- **Mathematical Functions:**
  - `sqrt(x)` - Square root
  - `pow(base, exp)` - Power function
  - `pi()` - Returns the value of π
  - `sin(x)` - Sine function
  - `cos(x)` - Cosine function
  - `tan(x)` - Tangent function
  - `asin(x)` - Arcsine function
  - `acos(x)` - Arccosine function
  - `atan(x)` - Arctangent function


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

Or after building:

```bash
./pars
```

To execute a pars source file:

```bash
./pars filename.pars
```

To build the project:

```bash
go build -o pars
./pars
```

To see the version:

```bash
./pars --version
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

### Arrays

Arrays are created by separating values with commas:

```
>> xs = 1,2,3
1, 2, 3
>> names = "Sam","Phillips"
Sam, Phillips
>> mixed = 1,"two",3.0,true
1, two, 3, true
```

#### Array Indexing

Arrays use zero-based indexing with square brackets:

```
>> xs = 1,2,3
1, 2, 3
>> xs[0]
1
>> xs[2]
3
>> (10,20,30)[1]
20
```

Negative indices count from the end of the array:

```
>> xs = 1,2,3
1, 2, 3
>> xs[-1]
3
>> xs[-2]
2
```

#### Array Slicing

Create sub-arrays using slice notation `[start:end]` (half-open range):

```
>> arr = 10,20,30,40,50
10, 20, 30, 40, 50
>> arr[1:4]
20, 30, 40
>> arr[0:2]
10, 20
```

#### Array Concatenation

Use the `++` operator to concatenate arrays:

```
>> 1,2,3 ++ 4,5,6
1, 2, 3, 4, 5, 6
>> arr = 10,20,30
>> arr ++ 40,50
10, 20, 30, 40, 50
```

Single values are treated as single-element arrays:

```
>> 1 ++ 2 ++ 3
1, 2, 3
>> 1,2,3 ++ 4
1, 2, 3, 4
```

#### Array Length

Get the number of elements in an array:

```
>> arr = 10,20,30,40,50
10, 20, 30, 40, 50
>> len(arr)
5
```

### Strings

#### String Concatenation

Use the `+` operator to join strings:

```
>> "Hello, " + "world!"
Hello, world!
>> name = "Sam"
Sam
>> "Hello, " + name + "!"
Hello, Sam!
```

#### String Indexing

Strings can be indexed like arrays (zero-based):

```
>> "Hello"[1]
e
>> "World"[-1]
d
```

#### String Slicing

Extract substrings using slice notation:

```
>> "Pars"[0:2]
Pa
>> "Concatenation"[3:7]
cate
>> str = "Hello, World!"
Hello, World!
>> str[7:12]
World
```

#### String Length

Get the number of characters in a string:

```
>> len("Hello")
5
>> len("")
0
>> str = "Hello, Pars!"
Hello, Pars!
>> len(str)
13
```

#### Escape Sequences

Strings support escape sequences for special characters:

```
>> "Line 1\nLine 2"
Line 1
Line 2
>> "Column1\tColumn2"
Column1	Column2
>> "Quote: \"Hello\""
Quote: "Hello"
>> "Backslash: \\"
Backslash: \
```

Supported escape sequences:
- `\n` - newline
- `\t` - tab
- `\\` - backslash
- `\"` - double quote

#### Template Literals

Template literals use backticks and support expression interpolation:

```
>> `Hello, World!`
Hello, World!
>> name = "Sam"
Sam
>> `Welcome, ${name}!`
Welcome, Sam!
```

Interpolate any expression with `${}`:

```
>> a = 5
5
>> b = 10
10
>> `Sum: ${a + b}`
Sum: 15
>> `Result: ${a * 2 + b}`
Result: 20
```

Template literals are multiline and preserve whitespace:

```
>> `Line 1
   Line 2
   Line 3`
Line 1
Line 2
Line 3
```

Arrays in templates are joined without commas:

```
>> `Items: ${"A","B","C"}`
Items: ABC
```

Type coercion in templates:

```
>> `Number: ${42}`
Number: 42
>> `Boolean: ${true}`
Boolean: true
>> `Expression: ${10 > 5}`
Expression: true
```

Escape special characters in templates:

```
>> `Literal backtick: \``
Literal backtick: `
>> `Not interpolated: \${variable}`
Not interpolated: ${variable}
```

String concatenation with automatic type conversion:

```
>> "Count: " + 42
Count: 42
>> "Result: " + (5 + 3)
Result: 8
```

### Array Operations with map

The `map` function applies a function to each element of an array, filtering out null values:

```
>> double = fn(x) { x * 2 }
>> map(double, 1,2,3)
2, 4, 6
>> square = fn(x) { x * x }
>> nums = 2,3,4
>> map(square, nums)
4, 9, 16
```

Filtering with map (null values are skipped):

```
>> gt10 = fn(x) { if (x > 10) then x }
>> map(gt10, 5,15,25,8,3,12)
15, 25, 12
```

### For Loops

Simple form - apply function to array:

```
>> double = fn(x) { x * 2 }
>> for(1,2,3) double
2, 4, 6
```

For-in form - inline function body:

```
>> for(x in 1,2,3) { x * 2 }
2, 4, 6
>> for(x in 5,15,25) { if (x > 10) then x }
15, 25
```

Strings are automatically converted to arrays of characters:

```
>> for("Sam") toUpper
S, A, M
>> for("Sam","Phillips") toUpper
SAM, PHILLIPS
>> for(name in "SAM","PHILLIPS") { toLower(name) }
sam, phillips
```

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
