# Parsley

```
█▀█ ▄▀█ █▀█ █▀ █░░ █▀▀ █▄█
█▀▀ █▀█ █▀▄ ▄█ █▄▄ ██▄ ░█░ v 0.9.1
```

A minimalist concatenative programming language interpreter.
- Written in Go
- If JSX and PHP had a cool baby
- Based on Basil from 2001

## Table of Contents

- [Quick Start](#quick-start)
- [Language Overview](#language-overview)
- [Language Guide](#language-guide)
  - [Variables and Functions](#variables-and-functions)
  - [Arrays](#arrays)
  - [Dictionaries](#dictionaries)
  - [Strings](#strings)
  - [Control Flow](#control-flow)
  - [Regular Expressions](#regular-expressions)
  - [Dates and Times](#dates-and-times)
  - [Paths and URLs](#paths-and-urls)
  - [Module System](#module-system)
  - [HTML/XML Tags](#htmlxml-tags)
- [Reference](#reference)
  - [Data Types](#data-types)
  - [Operators](#operators)
  - [Built-in Functions](#built-in-functions)
- [Development](#development)
- [Examples](#examples)
- [License](#license)

## Quick Start

### Installation

```bash
git clone https://github.com/sambeau/parsley.git
cd parsley
go build -o pars .
```

### Hello World

```bash
# Interactive REPL
./pars

# Run a file
echo 'log("Hello, World!")' > hello.pars
./pars hello.pars
```

### Your First Program

```parsley
// Variables and functions
let name = "Alice"
let greet = fn(who) { "Hello, " + who + "!" }
log(greet(name))  // "Hello, Alice!"

// Arrays and iteration
let numbers = [1, 2, 3, 4, 5]
let doubled = for (n in numbers) { n * 2 }
log(doubled)  // [2, 4, 6, 8, 10]

// Generate HTML
let page = <html>
    <body>
        <h1>Welcome to Parsley!</h1>
    </body>
</html>
log(page)
```

## Language Overview

Parsley is a concatenative language with:
- **First-class functions** with closures
- **Pattern matching** via destructuring
- **Module system** for code reuse
- **Native HTML/XML** tag syntax
- **Regular expressions** as first-class values
- **Rich data types**: arrays, dictionaries, dates, paths, URLs

### Core Concepts

```parsley
// Everything is an expression
let x = if (true) 42 else 0

// Functions are values
let add = fn(a, b) { a + b }
let operations = [add, fn(a, b) { a - b }]

// Destructuring everywhere
let {name, age} = {name: "Sam", age: 57}
let first, rest = [1, 2, 3, 4]

// Template interpolation
let greeting = "Hello, {name}!"

// Module imports
let {add, multiply} = import(@./math.pars)
```

## Language Guide

### Variables and Functions

#### Variable Declaration

```parsley
// Using 'let' (exported from modules)
let x = 42
let name = "Alice"

// Direct assignment (private in modules)
count = 0

// Destructuring
let a, b, c = 1, 2, 3
let {x, y} = {x: 10, y: 20}

// Special underscore (write-only)
let _, value = [99, 100]  // Ignores 99
```

#### Functions

```parsley
// Basic function
let square = fn(x) { x * x }

// Multiple parameters
let add = fn(a, b) { a + b }

// Implicit return (last expression)
let double = fn(x) { x * 2 }

// Closures
let makeCounter = fn() {
    count = 0
    fn() { count = count + 1; count }
}
let counter = makeCounter()
counter()  // 1
counter()  // 2

// Destructuring parameters
let getX = fn({x, y}) { x }
let sum = fn(arr) { arr[0] + arr[1] }
```

### Arrays

```parsley
// Creation
let nums = [1, 2, 3]
let mixed = [1, "two", true, [4, 5]]

// Indexing (0-based)
nums[0]     // 1
nums[-1]    // 3 (last element)

// Slicing
nums[0:2]   // [1, 2]  - elements 0 and 1
nums[1:3]   // [2, 3]  - elements 1 and 2
nums[2:]    // [3]     - from index 2 to end
nums[:2]    // [1, 2]  - from start to index 2
nums[:]     // [1, 2, 3] - full copy

// Concatenation
[1, 2] ++ [3, 4]  // [1, 2, 3, 4]

// Iteration
for (n in nums) { n * 2 }  // [2, 4, 6]

// Common operations
len([1, 2, 3])              // 3
sort([3, 1, 2])             // [1, 2, 3]
reverse([1, 2, 3])          // [3, 2, 1]
map(fn(x) { x * 2 }, nums)  // [2, 4, 6]
```

### Dictionaries

```parsley
// Creation
let person = {
    name: "Sam",
    age: 57,
    greet: fn() { "Hello, " + this.name }
}

// Access
person.name        // "Sam"
person["age"]      // 57
person.greet()     // "Hello, Sam" (this binding)

// Lazy evaluation
let config = {
    width: 100,
    height: 200,
    area: this.width * this.height  // Computed on access
}

// Iteration
for (key, value in person) {
    log(key, ":", value)
}

// Operations
keys(person)        // ["name", "age", "greet"]
values(person)      // ["Sam", 57, fn]
has(person, "name") // true

// Merging
{a: 1} ++ {b: 2}   // {a: 1, b: 2}
```

### Strings

```parsley
// Basic strings
let text = "hello"

// Template interpolation
let name = "World"
let greeting = "Hello, {name}!"  // "Hello, World!"

// Multi-line
let poem = "
    Roses are red
    Violets are blue
"

// Indexing and slicing
"hello"[0]      // "h"
"hello"[1:4]    // "ell"
"hello"[-1]     // "o"
"hello"[2:]     // "llo"
"hello"[:3]     // "hel"

// Concatenation
"hello" + " " + "world"  // "hello world"

// Operations
len("hello")         // 5
toUpper("hello")     // "HELLO"
toLower("WORLD")     // "world"
split("a,b,c", ",")  // ["a", "b", "c"]
replace("hello", "h", "H")  // "Hello"
```

### Control Flow

#### If-Else

```parsley
// Expression form
let result = if (x > 10) "big" else "small"

// Block form
if (user.age >= 18) {
    log("Adult")
} else {
    log("Minor")
}

// Chaining
let grade = if (score >= 90) "A"
           else if (score >= 80) "B"
           else if (score >= 70) "C"
           else "F"
```

#### For Loops

```parsley
// Array iteration
for (item in items) {
    log(item)
}

// Dictionary iteration
for (key, value in dict) {
    log(key, "=", value)
}

// Map pattern (returns array)
let doubled = for (n in [1, 2, 3]) { n * 2 }

// Filter pattern (if returns null, item is excluded)
let evens = for (n in numbers) {
    if (n % 2 == 0) { n }
}

// Reduce pattern (accumulate values)
let sum = 0
for (n in numbers) {
    sum = sum + n
}
```

### Regular Expressions

```parsley
// Regex literals
let emailPattern = /^[\w.+-]+@[\w.-]+\.[a-zA-Z]{2,}$/
let phonePattern = /\((\d{3})\) (\d{3})-(\d{4})/

// Match operator (~ returns array or null)
let result = "test@example.com" ~ emailPattern
if (result) {
    log("Valid email:", result[0])
}

"invalid" ~ emailPattern  // null

// Not-match operator (!~)
"hello" !~ /\d+/  // true

// Capture groups
let text = "Phone: (555) 123-4567"
let match = text ~ phonePattern
if (match) {
    log("Area:", match[1])   // "555"
    log("Prefix:", match[2]) // "123"
    log("Line:", match[3])   // "4567"
}

// Dynamic regex
let pattern = regex("\\d+", "i")

// Replace
replace("hello world", /world/, "Parsley")  // "hello Parsley"

// Split
split("a1b2c3", /\d+/)  // ["a", "b", "c"]
```

### Dates and Times

```parsley
// Current time
let current = now()

// Parsing
let date = time("2024-11-26")
let withTime = time("2024-11-26T15:30:00")
let timestamp = time(1732579200)  // Unix timestamp

// Access components (dictionary properties)
current.year      // 2025
current.month     // 11
current.day       // 27
current.hour      // 14
current.minute    // 30
current.second    // 0
current.weekday   // "Thursday"
current.iso       // "2025-11-27T14:30:00Z"
current.unix      // Unix timestamp

// Comparison
let date1 = time("2024-01-01")
let date2 = time("2024-06-01")
date1 < date2  // true
```

### Paths and URLs

#### Paths

```parsley
// Path literals
let config = @./config.json
let binary = @/usr/local/bin/tool

// Properties (not functions)
config.basename    // "config.json"
config.ext         // "json"
config.stem        // "config"
config.dirname     // Directory path

// Dynamic paths
let p = path("/usr/local/bin/tool")
p.basename  // "tool"
p.dirname   // "/usr/local/bin"
p.ext       // ""
```

#### URLs

```parsley
// URL literals
let api = @https://api.example.com/users
let local = @http://localhost:8080/api

// Properties (not functions)
api.scheme    // "https"
api.host      // "api.example.com"
api.path      // "/users"

// Query parameters
let search = @https://example.com?q=test&page=2
search.query.q     // "test"
search.query.page  // "2"

// Dynamic URLs
let u = url("https://example.com:8080/path")
u.scheme  // "https"
u.host    // "example.com"
u.port    // "8080"
u.path    // "/path"
```

### Module System

Modules are regular Parsley scripts. Only `let` bindings are exported.

#### Creating a Module

**math.pars:**
```parsley
let PI = 3.14159

let add = fn(a, b) { a + b }
let multiply = fn(a, b) { a * b }

// Private (not exported - no 'let')
helper = fn(x) { x * 2 }
```

#### Importing Modules

```parsley
// Import entire module
let math = import(@./math.pars)
math.add(2, 3)  // 5

// Destructure imports
let {add, multiply} = import(@./math.pars)
add(10, 5)  // 15
```

#### Module Features

- **Caching**: Modules loaded once and cached
- **Circular dependency detection**: Prevents import cycles
- **Relative paths**: `@./file.pars`, `@../lib/utils.pars`
- **Private state**: Variables without `let` are module-private

#### Example: Counter with Private State

**counter.pars:**
```parsley
count = 0  // Private (no 'let')

let increment = fn() {
    count = count + 1
    count
}

let getCount = fn() { count }
```

**Usage:**
```parsley
let counter = import(@./counter.pars)
counter.increment()  // 1
counter.increment()  // 2
counter.count        // null (not exported)
counter.getCount()   // 2 (live access)
```

### HTML/XML Tags

#### Singleton Tags

```parsley
<br/>
<img src="photo.jpg" width="300" />
<meta charset="utf-8" />
```

#### Tag Pairs

```parsley
<div>
    <h1>Welcome</h1>
    <p>Hello, World!</p>
</div>
```

#### Components (Uppercase)

```parsley
let Card = fn({title, body}) {
    <div class="card">
        <h2>{title}</h2>
        <p>{body}</p>
    </div>
}

<Card title="Hello" body="This is a card" />
```

#### Fragments

```parsley
<>
    <p>First paragraph</p>
    <p>Second paragraph</p>
</>
```

## Reference

### Data Types

| Type | Example | Description |
|------|---------|-------------|
| Integer | `42`, `-15` | Whole numbers |
| Float | `3.14`, `2.718` | Decimal numbers |
| String | `"hello"`, `"world"` | Text with interpolation |
| Boolean | `true`, `false` | Logical values |
| Null | `null` | Absence of value |
| Array | `[1, 2, 3]` | Ordered collections |
| Dictionary | `{x: 1, y: 2}` | Key-value pairs |
| Function | `fn(x) { x * 2 }` | First-class functions |
| Regex | `/pattern/flags` | Regular expressions |
| Date/Time | `@2024-11-26` | Temporal values |
| Duration | `@1d`, `@2h` | Time spans |
| Path | `@./file.pars` | File paths |
| URL | `@https://example.com` | Web addresses |

### Operators

#### Arithmetic
- `+` Addition
- `-` Subtraction
- `*` Multiplication
- `/` Division
- `%` Modulo
- `++` Concatenation (arrays, dictionaries, strings)

#### Comparison
- `==` Equal
- `!=` Not equal
- `<` Less than
- `<=` Less than or equal
- `>` Greater than
- `>=` Greater than or equal

#### Logical
- `&&` AND
- `||` OR
- `!` NOT

#### Pattern Matching
- `~` Regex match (returns array or null)
- `!~` Regex not-match (returns boolean)

#### Special
- `=` Assignment
- `:` Dictionary key-value separator
- `.` Dictionary access
- `[]` Indexing and slicing
- `...` Spread operator (in progress)

### Built-in Functions

#### Type Conversion
- `toInt(str)` - String to integer
- `toFloat(str)` - String to float
- `toNumber(str)` - Auto-detect int/float
- `toString(...)` - Convert to string
- `toDebug(...)` - Debug representation

#### String Operations
- `toUpper(str)` - Convert to uppercase
- `toLower(str)` - Convert to lowercase
- `split(str, delim)` - Split by string or regex delimiter
- `replace(str, pattern, replacement)` - Replace matches (string or regex)
- `len(str)` - String length
- `contains(str, substr)` - Check if contains substring
- `starts_with(str, prefix)` - Check if starts with prefix
- `ends_with(str, suffix)` - Check if ends with suffix
- `trim(str)` - Remove leading/trailing whitespace

#### Array Operations
- `len(array)` - Array length
- `map(fn, array)` - Apply function to each element (or use for loops)
- `sort(array)` - Natural sort (returns new array)
- `sortBy(array, compareFn)` - Custom sort
- `reverse(array)` - Reverse copy

Note: `filter()` and `reduce()` can be implemented using for loops

#### Dictionary Operations
- `keys(dict)` - All keys
- `values(dict)` - All values (evaluated)
- `has(dict, key)` - Check key exists
- `toArray(dict)` - Convert to `[key, value]` pairs
- `toDict(array)` - Convert pairs to dictionary

#### Mathematical
- `sqrt(x)`, `round(x)`, `pow(base, exp)`
- `pi()` - Returns π
- `sin(x)`, `cos(x)`, `tan(x)`
- `asin(x)`, `acos(x)`, `atan(x)`

#### Date/Time
- `now()` - Current time
- `time(input)` - Parse/create datetime
- `time(input, delta)` - Apply time delta

#### Regular Expressions
- `regex(pattern, flags?)` - Create regex
- `replace(text, pattern, replacement)` - Replace
- `split(text, delimiter)` - Split

#### Modules
- `import(path)` - Import module

#### Paths/URLs
- `path(str)` - Parse file path
- `url(str)` - Parse URL

#### Debugging
- `log(...)` - Output to stdout
- `logLine(...)` - Output with file:line prefix

## Development

### Building from Source

```bash
# Using Make
make build    # Build binary
make test     # Run tests
make clean    # Remove binary
make install  # Install to $GOPATH/bin

# Manual build
go build -ldflags "-X main.Version=$(cat VERSION)" -o pars .
```

### Running

```bash
# Interactive REPL
./pars

# Execute file
./pars script.pars

# Pretty-print HTML output
./pars --pretty page.pars

# Show version
./pars --version
```

### Testing

```bash
# All tests
go test ./...

# Specific package
go test ./pkg/lexer -v
go test ./pkg/parser -v
go test ./pkg/evaluator -v

# With coverage
go test -cover ./...
```

## Examples

### Web Page Generator

```parsley
let Page = fn({title, content}) {
    <html>
        <head>
            <title>{title}</title>
            <style>{"
                body { font-family: sans-serif; margin: 2em; }
                h1 { color: #333; }
            "}</style>
        </head>
        <body>
            <h1>{title}</h1>
            {content}
        </body>
    </html>
}

<Page 
    title="My Blog" 
    content=<>
        <p>Welcome to my blog!</p>
        <p>This is generated with Parsley.</p>
    </>
/>
```

### Data Processing

```parsley
let data = [
    {name: "Alice", score: 95},
    {name: "Bob", score: 82},
    {name: "Carol", score: 91}
]

// Filter using for loop
let topStudents = for (student in data) {
    if (student.score >= 90) { student }
}

// Sort
let sorted = sortBy(topStudents, fn(a, b) {
    if (a.score > b.score) { [a, b] } else { [b, a] }
})

// Display results
for (student in sorted) {
    log(student.name, ":", student.score)
}

// Reduce pattern - calculate average
let total = 0
for (student in data) {
    total = total + student.score
}
let average = total / len(data)
log("Average:", average)
```

### Module Example

**validators.pars:**
```parsley
let emailRegex = /^[\w.+-]+@[\w.-]+\.[a-zA-Z]{2,}$/
let phoneRegex = /^\(\d{3}\) \d{3}-\d{4}$/

let isEmail = fn(str) { str ~ emailRegex != null }
let isPhone = fn(str) { str ~ phoneRegex != null }
let isStrongPassword = fn(str) {
    len(str) >= 8 && str ~ /[A-Z]/ && str ~ /[0-9]/
}
```

**app.pars:**
```parsley
let {isEmail, isStrongPassword} = import(@./validators.pars)

let validateUser = fn(email, password) {
    if (!isEmail(email)) {
        "Invalid email"
    } else if (!isStrongPassword(password)) {
        "Password must be 8+ chars with uppercase and number"
    } else {
        "Valid"
    }
}

log(validateUser("test@example.com", "Secret123"))  // "Valid"
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
