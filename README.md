# Parsley

```
█▀█ ▄▀█ █▀█ █▀ █░░ █▀▀ █▄█
█▀▀ █▀█ █▀▄ ▄█ █▄▄ ██▄ ░█░ v 0.9.8
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
// Array iteration (single parameter - element only)
for (item in items) {
    log(item)
}

// Array iteration with index (two parameters - index and element)
for (i, item in items) {
    log(i, ":", item)  // 0 : first, 1 : second, etc.
}

// String iteration with index
for (i, char in "hello") {
    log(i, "=", char)  // 0 = h, 1 = e, 2 = l, etc.
}

// Dictionary iteration
for (key, value in dict) {
    log(key, "=", value)
}

// Map pattern (returns array)
let doubled = for (n in [1, 2, 3]) { n * 2 }

// Map with index - enumerate pattern
let numbered = for (i, item in ["apple", "banana"]) {
    (i + 1) + ". " + item  // ["1. apple", "2. banana"]
}

// Filter pattern (if returns null, item is excluded)
let evens = for (n in numbers) {
    if (n % 2 == 0) { n }
}

// Filter with index - take first 3 items
let firstThree = for (i, item in items) {
    if (i < 3) { item }
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

// Computed properties for easy formatting
current.date      // "2025-11-27" (date only)
current.time      // "14:30" or "14:30:45" (time only)
current.format    // "November 27, 2025 at 14:30" (human-readable)
current.timestamp // Same as .unix (more intuitive)
current.dayOfYear // 331 (day number in year, 1-366)
current.week      // 48 (ISO week number, 1-53)

// Using in templates
let event = time({year: 2024, month: 12, day: 25, hour: 18, minute: 0})
log(<p>Event on {event.format}</p>)
// Output: <p>Event on December 25, 2024 at 18:00</p>

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

#### Web Components (Hyphenated Tags)

```parsley
<my-component>content</my-component>
<custom-element id="test">text</custom-element>
<my-icon name="star" />
```

#### XML Comments

XML comments are skipped during parsing:

```parsley
<div>hello<!-- this is a comment -->world</div>
// Output: <div>helloworld</div>
```

#### CDATA Sections

CDATA sections preserve literal content:

```parsley
<div><![CDATA[literal <b>text</b>]]></div>
// Output: <div>literal <b>text</b></div>
```

#### XML Processing Instructions

XML processing instructions (`<?...?>`) are passed through as strings:

```parsley
<?xml version="1.0" encoding="UTF-8"?>
// Output: <?xml version="1.0" encoding="UTF-8"?>

// Concatenate with HTML:
<?xml version="1.0"?> + <html><body>content</body></html>
// Output: <?xml version="1.0"?><html><body>content</body></html>
```

#### DOCTYPE Declarations

DOCTYPE declarations are passed through as strings:

```parsley
<!DOCTYPE html>
// Output: <!DOCTYPE html>

// Full HTML5 document:
<!DOCTYPE html> + <html><head></head><body></body></html>
```

#### Raw Text Mode (Style/Script Tags)

Inside `<style>` and `<script>` tags, braces `{}` are treated as literal characters (for CSS rules and JavaScript code). Use `@{}` for interpolation:

```parsley
<style>body { color: red; }</style>
// Output: <style>body { color: red; }</style>

// Interpolation with @{}:
color = "blue"
<style>.class { color: @{color}; }</style>
// Output: <style>.class { color: blue; }</style>

// JavaScript example:
value = 42
<script>var x = @{value};</script>
// Output: <script>var x = 42;</script>
```

Outside of style/script tags, `{}` works as normal interpolation.

#### Programmatic Tag Creation

The `tag()` function creates tag dictionaries for programmatic manipulation:

```parsley
// Create a tag programmatically
tag("div")
// Returns: {__type: tag, name: div, attrs: {}, contents: null}

// With attributes
tag("a", {href: "/home"})
// Returns: {__type: tag, name: a, attrs: {href: /home}, contents: null}

// With content
tag("p", {class: "intro"}, "Hello world")
// Returns: {__type: tag, name: p, attrs: {class: intro}, contents: Hello world}

// Convert back to HTML string
toString(tag("div", {class: "container"}, "Hello"))
// Output: <div class="container">Hello</div>

// Self-closing tags
toString(tag("br"))
// Output: <br />
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
- `trim(str)` - Remove leading/trailing whitespace
- `len(str)` - String length

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

#### Tag Operations
- `tag(name)` - Create tag dictionary
- `tag(name, attrs)` - Create tag with attributes dictionary
- `tag(name, attrs, contents)` - Create tag with attributes and content
- Note: `toString()` converts tag dictionaries back to HTML strings

#### Mathematical
- `sqrt(x)`, `round(x)`, `pow(base, exp)`
- `pi()` - Returns π
- `sin(x)`, `cos(x)`, `tan(x)`
- `asin(x)`, `acos(x)`, `atan(x)`

#### Date/Time
- `now()` - Current time
- `time(input)` - Parse/create datetime
- `time(input, delta)` - Apply time delta

#### Localization
- `formatNumber(num, locale?)` - Format number with locale (e.g., `formatNumber(1234.5, "de-DE")` → "1.234,5")
- `formatCurrency(num, currencyCode, locale?)` - Format currency (e.g., `formatCurrency(99.99, "EUR", "de-DE")` → "99,99 €")
- `formatPercent(num, locale?)` - Format percentage (e.g., `formatPercent(0.1234, "de-DE")` → "12,34 %")
- `formatDate(datetime, style?, locale?)` - Format date (style: "short", "medium", "long", "full")
- `format(duration, locale?)` - Format duration as relative time (e.g., `format(@-1d, "de-DE")` → "gestern")
- `format(array, style?, locale?)` - Format list with locale (style: "and", "or", "unit")

Example locale-aware formatting:
```parsley
// Numbers
formatNumber(1234567.89, "de-DE")     // "1.234.567,89"
formatNumber(1234567.89, "fr-FR")     // "1 234 567,89"

// Currency
formatCurrency(99.99, "USD", "en-US") // "$ 99.99"
formatCurrency(99.99, "EUR", "de-DE") // "99,99 €"

// Dates
let d = time({year: 2024, month: 12, day: 25})
formatDate(d, "long", "en-US")        // "December 25, 2024"
formatDate(d, "long", "de-DE")        // "25. Dezember 2024"
formatDate(d, "long", "fr-FR")        // "25 décembre 2024"
formatDate(d, "long", "ja-JP")        // "2024年12月25日"
formatDate(d, "full", "es-ES")        // "miércoles, 25 de diciembre de 2024"

// Relative Time (durations)
format(@1d)                           // "tomorrow"
format(@-1d)                          // "yesterday"
format(@-2d, "de-DE")                 // "vorgestern"
format(@3h, "fr-FR")                  // "dans 3 heures"

// Relative time with datetime arithmetic
let christmas = @2025-12-25
format(christmas - now())             // "in 4 weeks" (varies by current date)

// List Formatting
format(["apple", "banana", "cherry"]) // "apple, banana, and cherry"
format(["a", "b", "c"], "or")         // "a, b, or c"
format(["a", "b", "c"], "and", "en-GB") // "a, b and c" (no Oxford comma)
format(["Apfel", "Banane"], "and", "de-DE") // "Apfel und Banane"
format(["りんご", "バナナ"], "and", "ja-JP") // "りんごとバナナ"
format(["5 feet", "6 inches"], "unit") // "5 feet, 6 inches"
```

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
let darkColor = "#333"
let brightColor = "#007bff"
let date = time(now().date).format
let author = "Sam Phillips"

let Page = fn({title, content}) {
	<!DOCTYPE html> + <html>
		<head>
			<meta charset="utf-8" />
			<title>{title}</title>
			// style tags can contain raw CSS without parsing
			// only @{} sections are interpolated
			<style>
				/* updated: @{date} by: @{author} */
				body {
					font-family: sans-serif;
					margin: 2em;
					line-height: 1.6;
				}
				h1 {
					color: @{darkColor};
					border-bottom: 2px solid @{brightColor};
					padding-bottom: 0.5em;
				}
				.container {
					max-width: 800px;
					margin: 0 auto;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<h1>{title}</h1>
				{content}
			</div>
		</body>
	</html>
}

// contents of Page component get passed as "content" prop
<Page title="My Blog">
    <h1>Welcome to my blog!</h1>
    <p>This is generated with Parsley.</p>
</Page>
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
