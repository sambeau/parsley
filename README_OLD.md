# Parsley

```
█▀█ ▄▀█ █▀█ █▀ █░░ █▀▀ █▄█
█▀▀ █▀█ █▀▄ ▄█ █▄▄ ██▄ ░█░ v 0.8.0
```

A concatenative programming language interpreter.
- Writen in Go
- If JSX and PHP had a cool baby
- Based on Basil from 2001

## Features

### Core Language Features
- Variable declarations with `let`
- Direct variable assignment (e.g., `x = 5`)
- Array destructuring assignment (e.g., `x,y,z = 1,2,3`)
- Dictionary destructuring in assignments and function parameters
- Functions with `fn`
- Function parameter destructuring for arrays and dictionaries
- If-else expressions with block or expression forms
- Single return statements allowed after if without braces
- Arrays with comma separator or square bracket notation `[...]`
- Multi-dimensional arrays (arrays containing arrays)
- Array indexing and slicing with `[]`
- Chained indexing for nested arrays (e.g., `arr[0][1][2]`)
- Array concatenation with `++`
- Natural sorting with `sort()` function
- Dictionary objects with key-value pairs
- Lazy evaluation of dictionary values with `this` binding
- Dictionary access via dot notation and bracket indexing
- Dictionary concatenation and merging with `++`
- Dictionary iteration with `for(key, value in dict)`
- Dictionary manipulation: `keys()`, `values()`, `has()`, `toArray()`, `toDict()`, `delete`
- Module system with `import()` for code reuse across files
- Module caching and circular dependency detection
- String indexing and slicing
- String concatenation with `+`
- String escape sequences (`\n`, `\t`, etc.)
- Template literals with `{}` interpolation
- Singleton tags for HTML/XML markup (`<tag attr="value" />`)
- Tag pairs for structured content (`<tag>content</tag>`)
- Component system with props and contents
- Empty grouping tags (`<>...</>`) for fragments
- Integer and floating-point arithmetic
- Modulo operator (`%`) for remainder calculations
- Boolean logic
- Single-line comments with `//`
- Regular expression literals (`/pattern/flags`)
- Regex matching with `~` operator (returns array or null)
- Regex non-matching with `!~` operator (returns boolean)
- Special `_` variable (write-only, always returns `null`)

### Data Types
- **Integers:** `42`, `-15`
- **Floats:** `3.14159`, `2.718`
- **Strings:** `"hello world"`, multi-line strings supported
- **Booleans:** `true`, `false`
- **Null:** `null` - represents absence of a value
- **Arrays:** `1,2,3`, `[1,2,3]`, `[[1,2],[3,4]]`, mixed types allowed
- **Dictionaries:** `{ name: "Sam", age: 57 }`, key-value pairs with lazy evaluation
- **Functions:** `fn(x) { x * 2 }`, first-class functions with closures
- **Regular Expressions:** `/pattern/flags` - regex literals with `~` match operator
- **Paths:** `@/usr/local/bin`, `@./config.json` - file path literals with component access
- **URLs:** `@https://example.com/api` - URL literals with parsed components

**Note:** Tags (`<div>`, `<Component />`) are syntactic constructs that evaluate to strings, not a separate data type. Paths and URLs are dictionary-based types with special computed properties.

### Built-in Functions

- **Type Conversion Functions:**
  - `toInt(str)` - Convert string to integer
  - `toFloat(str)` - Convert string to float
  - `toNumber(str)` - Convert string to integer or float (auto-detects)
  - `toString(values...)` - Convert values to strings and join without whitespace
  - `toDebug(values...)` - Convert values to debug representation (arrays in `[...]`, strings in `"quotes"`)

- **Debugging Functions:**
  - `log(values...)` - Output values in debug format immediately to stdout (returns `null`)
  - `logLine(values...)` - Output values with filename and line number prefix (returns `null`)

- **String Functions:**
  - `toUpper(str)` - Convert string to uppercase
  - `toLower(str)` - Convert string to lowercase
  - `len(str)` - Get the length of a string

- **Array Functions:**
  - `map(func, elements...)` - Apply function to each element, filter out nulls
  - `for(array) func` - Sugar syntax for map with function
  - `for(var in array) { body }` - Sugar syntax for map with inline function
  - `len(array)` - Get the length of an array
  - `sort(array)` - Return a naturally sorted copy of the array
  - `sortBy(array, compareFunc)` - Return a sorted copy using a custom comparison function
  - `reverse(array)` - Return a reversed copy of the array

- **Dictionary Functions:**
  - `keys(dict)` - Return an array of all dictionary keys
  - `values(dict)` - Return an array of all dictionary values (evaluated)
  - `has(dict, key)` - Check if dictionary contains a key (returns boolean)
  - `toArray(dict)` - Convert dictionary to array of `[key, value]` pairs
  - `toDict(array)` - Convert array of `[key, value]` pairs to dictionary

- **Mathematical Functions:**
  - `sqrt(x)` - Square root
  - `round(x)` - Round to nearest integer
  - `pow(base, exp)` - Power function
  - `pi()` - Returns the value of π
  - `sin(x)` - Sine function
  - `cos(x)` - Cosine function
  - `tan(x)` - Tangent function
  - `asin(x)` - Arcsine function
  - `acos(x)` - Arccosine function
  - `atan(x)` - Arctangent function

- **Date/Time Functions:**
  - `now()` - Returns current time as a dictionary
  - `time(input)` - Parse/create datetime from string, integer (Unix timestamp), or dictionary
  - `time(input, delta)` - Parse/create datetime and apply time delta

- **Regular Expression Functions:**
  - `regex(pattern, flags?)` - Create regex from string pattern (flags optional)
  - `replace(text, pattern, replacement)` - Replace matches (pattern can be string or regex)
  - `split(text, delimiter)` - Split string by delimiter (can be string or regex)

- **Module Functions:**
  - `import(path)` - Import a Parsley module and return its exported scope as a dictionary

- **Path Functions:**
  - `path(str)` - Parse a file path string into a path dictionary with components

- **URL Functions:**
  - `url(str)` - Parse a URL string into a URL dictionary with parsed components


## Getting Started

### Prerequisites

- Go 1.19 or higher

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/sambeau/parsley.git
   cd parsley
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Build the project:
   ```bash
   make build
   # or manually:
   go build -ldflags "-X main.Version=$(cat VERSION)" -o pars .
   ```

The version number is read from the `VERSION` file at the repository root and compiled into the binary.

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

To pretty-print HTML output:

```bash
./pars -pp filename.pars
# or
./pars --pretty filename.pars
```

The `-pp` or `--pretty` flag auto-detects HTML output and formats it with proper indentation, making it easier to read and debug. Non-HTML output is left unchanged.

To see the version:

```bash
./pars --version
```

### Build Commands

Using Make:
```bash
make build    # Build the binary
make test     # Run tests
make clean    # Remove built binary
make install  # Install to $GOPATH/bin
```

Manual build:
```bash
go build -ldflags "-X main.Version=$(cat VERSION)" -o pars .
```

## Usage Examples

### Basic Arithmetic
```
>> 2 + 3
5
>> 10 * 4.5
45
>> 17 % 5
2
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
>> 10 % 3
1
>> 15 % 4
3
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

### If-Else Expressions
```
>> x = if (5 > 0) true else false
true
>> y = if (1 < 0) 0
null
>> a = 10
>> result = if (a > 5) "big" else "small"
big
>> bar = 15
>> foo = if (bar * 20 > 100) 100 else bar
100
>> nested = if (1 > 0) if (2 > 1) 3 else 4 else 5
3
```

If-else expressions can be used anywhere an expression is expected. The `else` clause is optional - if omitted and the condition is false, the expression evaluates to `null`.

#### If Statement Forms

Parsley supports three forms of if statements:

**Block form** (for multiple statements):
```pars
if (x > 10) {
    y = x * 2
    return y
}
```

**Expression form** (single expression):
```pars
if (x > 10) x * 2 else x + 1
```

**Single return form** (return without braces):
```pars
if (x > 10)
    return x * 2
```

This is particularly useful in for loops:
```pars
for (x in items) {
    if (x > 10)
        return x
}
```

### Arrays

Arrays can be created using comma notation or square bracket notation:

```
>> xs = 1,2,3
1, 2, 3
>> ys = [1,2,3]
1, 2, 3
>> names = "Sam","Phillips"
Sam, Phillips
>> mixed = 1,"two",3.0,true
1, two, 3, true
```

#### Multi-dimensional Arrays

Use square brackets `[...]` to create nested arrays:

```
>> matrix = [[1,2,3],[4,5,6],[7,8,9]]
1, 2, 3, 4, 5, 6, 7, 8, 9
>> matrix[0]
1, 2, 3
>> matrix[1][2]
6
>> tensor = [[[1,2],[3,4]],[[5,6],[7,8]]]
1, 2, 3, 4, 5, 6, 7, 8
>> tensor[1][0][1]
7
```

Empty arrays and nested empty arrays are supported:

```
>> []

>> [[]]

```

See [MULTIDIM_ARRAYS.md](examples/MULTIDIM_ARRAYS.md) for more examples.

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

### Dictionaries

Dictionaries are key-value data structures with lazy evaluation semantics. Values in dictionaries are stored as unevaluated expressions and only computed when accessed, enabling powerful patterns like self-referential objects and computed properties.

#### Creating Dictionaries

Dictionaries are created using curly braces with `key: value` syntax. Keys must be identifiers, and values can be any expression:

```
>> person = { name: "Sam", age: 57 }
{name: Sam, age: 57}
>> person
{name: Sam, age: 57}
```

Single-line and multi-line formats are supported:

```
>> point = { x: 10, y: 20 }
{x: 10, y: 20}

>> config = {
   timeout: 30
   retries: 3
   endpoint: "api.example.com"
}
{timeout: 30, retries: 3, endpoint: api.example.com}
```

#### Accessing Dictionary Values

Access values using dot notation or bracket indexing:

```
>> person = { name: "Alice", age: 30, city: "NYC" }
{name: Alice, age: 30, city: NYC}
>> person.name
Alice
>> person["age"]
30
>> person.city
NYC
```

Both forms evaluate the value expression when accessed.

#### Lazy Evaluation

Dictionary values are stored as expressions and only evaluated when accessed. This allows self-referential dictionaries using the special `this` variable:

```
>> circle = {
   radius: 5
   area: pi() * pow(this.radius, 2)
   circumference: 2 * pi() * this.radius
}
{radius: 5, area: pi() * pow(this.radius, 2), circumference: 2 * pi() * this.radius}
>> circle.area
78.53981633974483
>> circle.circumference
31.41592653589793
```

The `this` variable always refers to the dictionary being accessed, enabling computed properties:

```
>> rectangle = {
   width: 10
   height: 5
   area: this.width * this.height
   perimeter: 2 * (this.width + this.height)
}
{width: 10, height: 5, area: this.width * this.height, perimeter: 2 * (this.width + this.height)}
>> rectangle.area
50
>> rectangle.perimeter
30
```

#### Functions in Dictionaries

Dictionary values can be functions that reference other properties via `this`:

```
>> calculator = {
   x: 10
   y: 5
   add: fn() { this.x + this.y }
   multiply: fn() { this.x * this.y }
}
{x: 10, y: 5, add: fn() { this.x + this.y }, multiply: fn() { this.x * this.y }}
>> calculator.add()
15
>> calculator.multiply()
50
```

No-argument functions are automatically called when accessed via dictionary builtins like `values()` and `toArray()`:

```
>> obj = {
   name: "Greeter"
   getMessage: fn() { "Hello, " + this.name + "!" }
}
{name: Greeter, getMessage: fn() { "Hello, " + this.name + "!" }}
>> values(obj)
Greeter, Hello, Greeter!
```

#### Deleting Dictionary Keys

Remove keys from dictionaries using the `delete` statement:

```
>> user = { name: "Bob", age: 25, email: "bob@example.com" }
{name: Bob, age: 25, email: bob@example.com}
>> delete user.email
null
>> user
{name: Bob, age: 25}
```

Delete also works with bracket notation:

```
>> delete user["age"]
null
>> user
{name: Bob}
```

The `delete` statement returns `null` and modifies the dictionary in place.

#### Dictionary Concatenation

Merge dictionaries using the `++` operator. The right dictionary's values override the left's on key collision:

```
>> defaults = { timeout: 30, retries: 3, debug: false }
{timeout: 30, retries: 3, debug: false}
>> custom = { retries: 5, debug: true }
{retries: 5, debug: true}
>> config = defaults ++ custom
{timeout: 30, retries: 5, debug: true}
```

This is useful for configuration merging and object composition:

```
>> base = { a: 1, b: 2 }
{a: 1, b: 2}
>> override = { b: 20, c: 30 }
{b: 20, c: 30}
>> result = base ++ override
{a: 1, b: 20, c: 30}
>> result.b
20
```

#### Iterating Over Dictionaries

Use `for(key, value in dict)` to iterate over dictionary entries:

```
>> data = { name: "Alice", age: 30, city: "NYC" }
{name: Alice, age: 30, city: NYC}
>> for(key, value in data) {
   log(key + ":", value)
}
name: "Alice"
age: 30
city: "NYC"
null
```

The loop evaluates each value expression with `this` bound to the dictionary:

```
>> temps = {
   celsius: 25
   fahrenheit: this.celsius * 9/5 + 32
   kelvin: this.celsius + 273.15
}
{celsius: 25, fahrenheit: this.celsius * 9/5 + 32, kelvin: this.celsius + 273.15}
>> for(key, value in temps) {
   key + " = " + toString(value)
}
celsius = 25, fahrenheit = 77, kelvin = 298.15
```

#### Dictionary Built-in Functions

**`keys(dict)`** - Returns an array of all dictionary keys:

```
>> person = { name: "Sam", age: 57, city: "NYC" }
{name: Sam, age: 57, city: NYC}
>> keys(person)
name, age, city
```

**`values(dict)`** - Returns an array of all evaluated dictionary values:

```
>> point = { x: 10, y: 20 }
{x: 10, y: 20}
>> values(point)
10, 20
```

Values are evaluated with `this` bound to the dictionary:

```
>> obj = {
   radius: 5
   area: pi() * pow(this.radius, 2)
}
{radius: 5, area: pi() * pow(this.radius, 2)}
>> values(obj)
5, 78.53981633974483
```

**`has(dict, key)`** - Checks if a dictionary contains a key:

```
>> user = { name: "Alice", email: "alice@example.com" }
{name: Alice, email: alice@example.com}
>> has(user, "name")
true
>> has(user, "age")
false
```

**`toArray(dict)`** - Converts a dictionary to an array of `[key, value]` pairs:

```
>> person = { name: "Sam", age: 57 }
{name: Sam, age: 57}
>> toArray(person)
[["name", "Sam"], ["age", 57]]
```

No-argument functions are automatically called:

```
>> obj = {
   x: 10
   doubled: fn() { this.x * 2 }
}
{x: 10, doubled: fn() { this.x * 2 }}
>> toArray(obj)
[["x", 10], ["doubled", 20]]
```

**`toDict(array)`** - Converts an array of `[key, value]` pairs to a dictionary:

```
>> pairs = [["name", "Bob"], ["age", 25]]
[["name", "Bob"], ["age", 25]]
>> toDict(pairs)
{name: Bob, age: 25}
```

Round-trip conversion preserves structure:

```
>> original = { x: 100, y: 200 }
{x: 100, y: 200}
>> reconstructed = toDict(toArray(original))
{x: 100, y: 200}
>> reconstructed.x
100
```

#### Practical Dictionary Examples

**Configuration management:**

```
>> defaults = { host: "localhost", port: 8080, debug: false }
>> userConfig = { port: 3000, debug: true }
>> finalConfig = defaults ++ userConfig
{host: localhost, port: 3000, debug: true}
```

**Data transformation:**

```
>> greetings = { en: "Hello", es: "Hola", fr: "Salut" }
{en: Hello, es: Hola, fr: Salut}
>> for(lang, greeting in greetings) {
   lang + ": " + greeting
}
en: Hello, es: Hola, fr: Salut
```

**Computed properties with `this`:**

```
>> invoice = {
   items: 3
   pricePerItem: 25
   subtotal: this.items * this.pricePerItem
   tax: this.subtotal * 0.1
   total: this.subtotal + this.tax
}
{items: 3, pricePerItem: 25, subtotal: this.items * this.pricePerItem, tax: this.subtotal * 0.1, total: this.subtotal + this.tax}
>> invoice.total
82.5
```

**Methods using functions:**

```
>> counter = {
   count: 0
   increment: fn() { this.count + 1 }
   decrement: fn() { this.count - 1 }
}
{count: 0, increment: fn() { this.count + 1 }, decrement: fn() { this.count - 1 }}
>> counter.increment()
1
>> counter.decrement()
0
```

**Filtering dictionary entries:**

```
>> scores = { alice: 85, bob: 92, charlie: 78, diana: 95 }
{alice: 85, bob: 92, charlie: 78, diana: 95}
>> highScores = for(name, score in scores) {
   if (score >= 90) { [name, score] }
}
[["bob", 92], ["diana", 95]]
>> toDict(highScores)
{bob: 92, diana: 95}
```

### Strings

Parsley supports two types of strings:

#### String Literals (Double Quotes)

String literals use double quotes and can span multiple lines:

```
>> "Hello, World!"
Hello, World!
>> name = "Alice"
Alice
```

Multi-line string literals preserve newlines and formatting, making them ideal for embedding CSS, HTML, SQL, or other formatted text:

```
>> css = "
body {
  margin: 0;
  padding: 0;
}
.container {
  max-width: 1200px;
}
"
body {
  margin: 0;
  padding: 0;
}
.container {
  max-width: 1200px;
}

>> html = "<div class='card'>
  <h1>Title</h1>
  <p>Content</p>
</div>"
<div class='card'>
  <h1>Title</h1>
  <p>Content</p>
</div>
```

#### Template Literals (Backticks)

Template literals use backticks and support expression interpolation with `{expression}` syntax:

```
>> name = "Sam"
Sam
>> `Welcome, {name}!`
Welcome, Sam!
>> a = 5
5
>> b = 10
10
>> `Sum: {a + b}`
Sum: 15
```

Template literals also support multi-line text:

```
>> card = `
<div class="card">
  <h2>{title}</h2>
  <p>{description}</p>
</div>
`
<div class="card">
  <h2>{title}</h2>
  <p>{description}</p>
</div>
```

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
>> str = "Hello, Parsley!"
Hello, Parsley!
>> len(str)
15
```

#### Escape Sequences

**In String Literals (double quotes):**

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
- `\r` - carriage return
- `\\` - backslash
- `\"` - double quote

Note: Multi-line string literals preserve newlines literally, so `\n` is typically only needed when constructing strings programmatically.

**In Template Literals (backticks):**

```
>> `Literal backtick: \``
Literal backtick: `
>> `Not interpolated: \{variable}`
Not interpolated: {variable}
>> css = `body \{ margin: 0; \}`
body { margin: 0; }
```

Supported escape sequences:
- `\n` - newline
- `\t` - tab
- `\r` - carriage return
- `\\` - backslash
- `\`` - backtick
- `\{` - literal left brace (prevents interpolation)
- `\}` - literal right brace (prevents interpolation)

Template literals also preserve literal newlines, making them ideal for multi-line formatted output.

#### Practical Example: Inline CSS

Multi-line string literals are particularly useful for embedding CSS directly in components:

```javascript
Page = fn({title, contents}) {
  <html lang="en">
    <head>
      <title>{title}</title>
      <style>{"
        *, *:before, *:after { box-sizing: border-box; }
        html, body { margin: 0; padding: 0; }
        .container { max-width: 1200px; margin: 0 auto; }
      "}</style>
    </head>
    <body>
      {contents}
    </body>
  </html>
}
```

The string literal preserves formatting and indentation, making the CSS easy to read and maintain.

#### Template Literal Interpolation Details

Arrays in templates are joined without commas:

```
>> `Items: {"A","B","C"}`
Items: ABC
```

Type coercion in templates:

```
>> `Number: {42}`
Number: 42
>> `Boolean: {true}`
Boolean: true
>> `Expression: {10 > 5}`
Expression: true
```

String concatenation with automatic type conversion:

```
>> "Count: " + 42
Count: 42
>> "Result: " + (5 + 3)
Result: 8
```

### Date and Time

Parsley provides minimal, composable datetime support through two built-in functions that work with dictionaries. Datetimes are represented as dictionaries with standard fields, making them transparent and easy to manipulate.

#### Getting the Current Time

Use `now()` to get the current time as a dictionary:

```
>> now()
{year: 2025, month: 11, day: 26, hour: 21, minute: 15, second: 42, weekday: "Wednesday", unix: 1764191742, iso: "2025-11-26T21:15:42Z"}
>> let dt = now()
{year: 2025, month: 11, day: 26, hour: 21, minute: 15, second: 42, weekday: "Wednesday", unix: 1764191742, iso: "2025-11-26T21:15:42Z"}
>> dt.year
2025
>> dt.weekday
Wednesday
```

#### Creating and Parsing Datetimes

The `time()` function creates datetimes from multiple input types:

**From ISO 8601 string:**
```
>> time("2024-01-15T10:30:00Z")
{year: 2024, month: 1, day: 15, hour: 10, minute: 30, second: 0, weekday: "Monday", unix: 1705316400, iso: "2024-01-15T10:30:00Z"}
>> time("2024-12-25")  // Date only (time defaults to 00:00:00)
{year: 2024, month: 12, day: 25, hour: 0, minute: 0, second: 0, weekday: "Wednesday", unix: 1735081200, iso: "2024-12-25T00:00:00Z"}
```

**From Unix timestamp:**
```
>> time(1704110400)
{year: 2024, month: 1, day: 1, hour: 12, minute: 0, second: 0, weekday: "Monday", unix: 1704110400, iso: "2024-01-01T12:00:00Z"}
```

**From dictionary:**
```
>> time({year: 2024, month: 7, day: 4, hour: 12, minute: 30, second: 0})
{year: 2024, month: 7, day: 4, hour: 12, minute: 30, second: 0, weekday: "Thursday", unix: 1720097400, iso: "2024-07-04T12:30:00Z"}
>> time({year: 2024, month: 12, day: 25})  // Time fields optional (default to 0)
{year: 2024, month: 12, day: 25, hour: 0, minute: 0, second: 0, weekday: "Wednesday", unix: 1735081200, iso: "2024-12-25T00:00:00Z"}
```

#### Datetime Literal Syntax

For cleaner code, use the `@` prefix with ISO-8601 format to create datetime literals directly:

```
>> @2024-12-25
{year: 2024, month: 12, day: 25, hour: 0, minute: 0, second: 0, weekday: "Wednesday", unix: 1735084800, iso: "2024-12-25T00:00:00Z"}
>> @2024-12-25T14:30:00
{year: 2024, month: 12, day: 25, hour: 14, minute: 30, second: 0, weekday: "Wednesday", iso: "2024-12-25T14:30:00Z"}
>> @2024-12-25T14:30:00Z
{year: 2024, month: 12, day: 25, hour: 14, minute: 30, second: 0, weekday: "Wednesday", iso: "2024-12-25T14:30:00Z"}
```

With timezone offsets:
```
>> @2024-12-25T14:30:00-05:00  // EST timezone
>> @2024-06-15T08:00:00+08:00  // Singapore timezone
```

Datetime literals work anywhere a datetime dictionary is expected:

```
>> let christmas = @2024-12-25;
>> christmas.day
25
>> if @2024-12-25 < now() { "Past" } else { "Future" }
>> [@2024-01-01, @2024-06-15, @2024-12-31]
```

Equivalent to `time()` function:
```
>> @2024-12-25 == time("2024-12-25")
true
```

#### Duration Literals

**Version 0.7.0+** Duration literals use the `@` prefix with time units to create duration values:

**Supported units:**
- `s` - seconds
- `m` - minutes  
- `h` - hours
- `d` - days
- `w` - weeks
- `mo` - months (note: two letters to distinguish from minutes)
- `y` - years

**Basic durations:**
```
>> @30s
{__type: duration, months: 0, seconds: 30, totalSeconds: 30}
>> @5m
{__type: duration, months: 0, seconds: 300, totalSeconds: 300}
>> @2h
{__type: duration, months: 0, seconds: 7200, totalSeconds: 7200}
>> @7d
{__type: duration, months: 0, seconds: 604800, totalSeconds: 604800}
```

**Compound durations:**
```
>> @2h30m
{__type: duration, months: 0, seconds: 9000, totalSeconds: 9000}
>> @1y6mo
{__type: duration, months: 18, seconds: 0, totalSeconds: null}
>> @1y2mo3w4d5h6m7s
{__type: duration, months: 14, seconds: 2178367, totalSeconds: null}
```

Duration dictionaries contain:
- `__type`: Always "duration"
- `months`: Number of months (variable-length, 12 months = 1 year)
- `seconds`: Number of seconds (fixed-length)
- `totalSeconds`: Total seconds if no months component, otherwise `null`

**Duration arithmetic:**
```
>> @2h + @30m
{__type: duration, months: 0, seconds: 9000, totalSeconds: 9000}
>> @1d - @6h
{__type: duration, months: 0, seconds: 64800, totalSeconds: 64800}
>> @2h * 3
{__type: duration, months: 0, seconds: 21600, totalSeconds: 21600}
>> @1d / 2
{__type: duration, months: 0, seconds: 43200, totalSeconds: 43200}
```

**Adding durations to datetimes:**
```
>> @2024-01-15 + @2d
{year: 2024, month: 1, day: 17, ...}
>> @2024-01-31 + @1mo  
{year: 2024, month: 3, day: 2, ...}  // Feb 31 normalizes to Mar 2
>> @2024-06-15T10:00:00 + @1h30m
{year: 2024, month: 6, day: 15, hour: 11, minute: 30, ...}
```

**⚠️ BREAKING CHANGE (v0.7.0):** Datetime subtraction now returns a Duration:
```
>> @2024-12-26 - @2024-12-25
{__type: duration, months: 0, seconds: 86400, totalSeconds: 86400}
>> let diff = @2024-01-20 - @2024-01-15
>> diff.seconds
432000
>> diff.seconds / 86400  // Convert to days
5
```

**Duration comparison** (seconds-only durations):
```
>> @1h < @2h
true
>> @2h == @120m
true
>> @30d > @4w  // 30*86400 > 4*7*86400
true
```

**Note:** Cannot compare durations with month components (months have variable length):
```
>> @1y < @12mo
ERROR: cannot compare durations with month components
```

**Migration from v0.6.x:**
If your code uses datetime subtraction expecting an integer result, update it:
```
// Old (v0.6.x):
let diff = dt1 - dt2  // Returns integer seconds

// New (v0.7.0+):
let diff = dt1 - dt2
let seconds = diff.seconds  // Access seconds field
```

#### Datetime Arithmetic

Apply time deltas using a second dictionary argument with `years`, `months`, `days`, `hours`, `minutes`, or `seconds` fields:

```
>> time("2024-01-01T00:00:00Z", {days: 7})  // Add 7 days
{year: 2024, month: 1, day: 8, hour: 0, minute: 0, second: 0, weekday: "Monday", unix: 1704672000, iso: "2024-01-08T00:00:00Z"}
>> time("2024-01-15T00:00:00Z", {days: -10})  // Subtract 10 days
{year: 2024, month: 1, day: 5, hour: 0, minute: 0, second: 0, weekday: "Friday", unix: 1704412800, iso: "2024-01-05T00:00:00Z"}
>> time("2024-01-01T00:00:00Z", {months: 3})  // Add 3 months
{year: 2024, month: 4, day: 1, hour: 0, minute: 0, second: 0, weekday: "Monday", unix: 1711958400, iso: "2024-04-01T00:00:00Z"}
```

Combine multiple deltas:
```
>> time("2024-01-01T12:30:00Z", {years: 1, months: 2, days: 15, hours: 3, minutes: 45, seconds: 30})
{year: 2025, month: 3, day: 16, hour: 16, minute: 15, second: 30, weekday: "Sunday", unix: 1742141730, iso: "2025-03-16T16:15:30Z"}
```

Use with `now()`:
```
>> time(now(), {days: 7})  // One week from now
>> time(now(), {days: -30})  // 30 days ago
```

#### Formatting Datetimes

Use template literals to format datetimes as needed:

```
>> let dt = time("2024-01-15T10:30:00Z")
{year: 2024, month: 1, day: 15, hour: 10, minute: 30, second: 0, weekday: "Monday", unix: 1705316400, iso: "2024-01-15T10:30:00Z"}
>> `{dt.year}-{dt.month}-{dt.day}`
2024-1-15
>> `{dt.year}-{dt.month}-{dt.day} {dt.hour}:{dt.minute}`
2024-1-15 10:30
>> `{dt.weekday}, {dt.month}/{dt.day}/{dt.year}`
Monday, 1/15/2024
```

Use ISO format for standardized output:
```
>> dt.iso
2024-01-15T10:30:00Z
```

#### Dictionary Fields

Datetime dictionaries contain the following fields:
- `year` - Four-digit year (e.g., 2024)
- `month` - Month number (1-12)
- `day` - Day of month (1-31)
- `hour` - Hour (0-23)
- `minute` - Minute (0-59)
- `second` - Second (0-59)
- `weekday` - Day name ("Monday", "Tuesday", etc.)
- `unix` - Unix timestamp (seconds since 1970-01-01)
- `iso` - ISO 8601 string (e.g., "2024-01-15T10:30:00Z")
- `__type` - Internal type tag ("datetime") for operator overloading

All times are in UTC.

#### Datetime Comparisons

Datetime dictionaries support standard comparison operators:

```
>> let dt1 = time("2024-01-15T10:30:00Z")
{year: 2024, month: 1, day: 15, hour: 10, minute: 30, second: 0, weekday: "Monday", unix: 1705316400, iso: "2024-01-15T10:30:00Z", __type: "datetime"}
>> let dt2 = time("2024-01-20T10:30:00Z")
{year: 2024, month: 1, day: 20, hour: 10, minute: 30, second: 0, weekday: "Saturday", unix: 1705748400, iso: "2024-01-20T10:30:00Z", __type: "datetime"}
>> dt1 < dt2
true
>> dt1 > dt2
false
>> dt1 == dt2
false
>> dt1 != dt2
true
```

All comparison operators work: `<`, `>`, `<=`, `>=`, `==`, `!=`

**Practical example:**
```
>> let deadline = time("2024-12-31T23:59:59Z")
>> let today = now()
>> today > deadline
false
>> log("Deadline passed:", today > deadline)
```

#### Datetime Arithmetic

**Difference between datetimes** returns seconds:
```
>> let dt1 = time("2024-01-15T00:00:00Z")
>> let dt2 = time("2024-01-20T00:00:00Z")
>> dt2 - dt1
432000
>> (dt2 - dt1) / 86400  // Convert to days
5
```

**Add/subtract seconds** from datetimes:
```
>> let dt = time("2024-01-15T12:00:00Z")
>> dt + 86400  // Add 1 day (86400 seconds)
{year: 2024, month: 1, day: 16, hour: 12, minute: 0, second: 0, weekday: "Tuesday", unix: 1705411200, iso: "2024-01-16T12:00:00Z", __type: "datetime"}
>> dt - 86400  // Subtract 1 day
{year: 2024, month: 1, day: 14, hour: 12, minute: 0, second: 0, weekday: "Sunday", unix: 1705238400, iso: "2024-01-14T12:00:00Z", __type: "datetime"}
>> 604800 + dt  // Addition is commutative (7 days)
{year: 2024, month: 1, day: 22, hour: 12, minute: 0, second: 0, weekday: "Monday", unix: 1705929600, iso: "2024-01-22T12:00:00Z", __type: "datetime"}
```

**Common time intervals:**
- 1 hour: `3600`
- 1 day: `86400`
- 1 week: `604800`
- 30 days: `2592000`

**Practical examples:**
```
>> let now_dt = now()
>> let in_one_week = now_dt + 604800
>> let days_ago = now_dt - (30 * 86400)

>> // Check if date is within range
>> let start = time("2024-01-01T00:00:00Z")
>> let end = time("2024-12-31T23:59:59Z")
>> let check = time("2024-06-15T12:00:00Z")
>> check >= start & check <= end
true
```

### Durations

Parsley provides first-class duration support through `@` prefix literals (like `@2h30m`, `@7d`, `@1y6mo`). Durations are dictionary-based with `__type: "duration"` and separate `months` and `seconds` components to handle variable-length months correctly.

#### Duration Literal Syntax

Duration literals start with `@` followed by number-unit pairs:

**Supported units:**
- `y` = years (12 months)
- `mo` = months (variable length, converted to months)
- `w` = weeks (7 days, converted to seconds)
- `d` = days (24 hours, converted to seconds)
- `h` = hours (3600 seconds)
- `m` = minutes (60 seconds)
- `s` = seconds

**Examples:**
```
>> @30s
{__type: duration, months: 0, seconds: 30, totalSeconds: 30}

>> @2h30m
{__type: duration, months: 0, seconds: 9000, totalSeconds: 9000}

>> @7d
{__type: duration, months: 0, seconds: 604800, totalSeconds: 604800}

>> @1y6mo
{__type: duration, months: 18, seconds: 0, totalSeconds: null}

>> @1y2mo3w4d5h6m7s
{__type: duration, months: 14, seconds: 2178367, totalSeconds: null}
```

**Note:** The `totalSeconds` field is `null` for durations with month components (since months have variable lengths). For pure seconds-based durations, `totalSeconds` equals `seconds`.

#### Duration Fields

Access duration components via dictionary fields:

```
>> let d = @2h30m
>> d.months
0
>> d.seconds
9000
>> d.totalSeconds
9000

>> let longDuration = @1y6mo
>> longDuration.months
18
>> longDuration.totalSeconds
null
```

#### Duration Arithmetic

**Add/subtract durations:**
```
>> @2h + @30m
{__type: duration, months: 0, seconds: 9000, totalSeconds: 9000}

>> @1d - @6h
{__type: duration, months: 0, seconds: 64800, totalSeconds: 64800}

>> @1y + @6mo
{__type: duration, months: 18, seconds: 0, totalSeconds: null}
```

**Multiply/divide by numbers:**
```
>> @2h * 3
{__type: duration, months: 0, seconds: 21600, totalSeconds: 21600}

>> @1d / 2
{__type: duration, months: 0, seconds: 43200, totalSeconds: 43200}

>> @1y * 2
{__type: duration, months: 24, seconds: 0, totalSeconds: null}
```

#### Duration Comparisons

Comparison operators work **only for pure seconds-based durations** (no month components):

```
>> @1h < @2h
true

>> @2h == @7200s
true

>> @1d >= @12h
true
```

**Error on month-based comparisons:**
```
>> @1y < @12mo
ERROR: cannot compare durations with month components (months have variable length)
```

#### Datetime + Duration Operations

Add or subtract durations from datetimes:

```
>> let start = @2024-01-15
>> start + @2d
{year: 2024, month: 1, day: 17, ...}

>> let meeting = @2024-06-15T10:00:00
>> meeting + @1h30m
{year: 2024, month: 6, day: 15, hour: 11, minute: 30, ...}

>> @2024-06-15 + @1y
{year: 2025, month: 6, day: 15, ...}

>> @2024-01-31 + @1mo
{year: 2024, month: 3, day: 2, ...}  // Normalized from Feb 31
```

**Month arithmetic:** When adding months to dates like Jan 31, if the result would be invalid (like Feb 31), it normalizes to the next valid date (Mar 2/3).

#### Breaking Change: Datetime Subtraction

**BREAKING:** Subtracting datetimes now returns a Duration instead of seconds:

**Before (v0.5.x):**
```
>> @2024-01-20 - @2024-01-15
432000  // Just an integer
```

**Now (v0.6.0+):**
```
>> @2024-01-20 - @2024-01-15
{__type: duration, months: 0, seconds: 432000, totalSeconds: 432000}

>> let diff = @2024-01-20 - @2024-01-15
>> diff.seconds / 86400  // Get days
5
```

**Migration:** Update code that uses `datetime - datetime` to access the `.seconds` field.

#### Practical Duration Examples

```
>> // Calculate project timeline
>> let sprint = @2w
>> let total = sprint * 6
>> total.seconds / 86400
84  // days

>> // Meeting time calculation
>> let daily_standup = @15m
>> let weekly_total = daily_standup * 5
>> weekly_total.seconds / 60
75  // minutes

>> // Date arithmetic with durations
>> let deadline = @2024-01-01 + @3mo2w
>> deadline.month
3

>> // Time until event
>> let event = @2024-12-25T00:00:00
>> let now_time = @2024-11-26T00:00:00
>> let remaining = event - now_time
>> remaining.seconds / 86400
29  // days
```

### Regular Expressions

Parsley provides first-class regex support through `/pattern/flags` literals and the `~` match operator. Regular expressions are dictionary-based (like datetimes and durations) with `__type: "regex"`, making them transparent and composable.

#### Regex Literals

Create regular expressions using familiar `/pattern/flags` syntax:

```
>> /\d+/
{pattern: "\d+", flags: "", __type: "regex"}
>> /hello/i
{pattern: "hello", flags: "i", __type: "regex"}
>> /test/gim
{pattern: "test", flags: "gim", __type: "regex"}
```

Access regex components:
```
>> let rx = /\w+@\w+/
>> rx.pattern
"\w+@\w+"
>> rx.flags
""
```

#### Match Operator (~)

The `~` operator matches a string against a regex, returning an array with the full match and capture groups, or `null` if no match:

```
>> "hello 123" ~ /\d+/
["123"]
>> "no numbers" ~ /\d+/
null
>> "user@example.com" ~ /(\w+)@([\w.]+)/
["user@example.com", "user", "example.com"]
```

The result is **truthy** (array) or **falsy** (null), perfect for conditionals:

```
>> let match = "Order #12345" ~ /Order #(\d+)/
>> if (match) {
     log("Order number:", match[1])
   }
Order number: "12345"
```

#### Not-Match Operator (!~)

The `!~` operator returns a boolean: `true` if the pattern does NOT match:

```
>> "hello world" !~ /\d+/
true
>> "hello 123" !~ /\d+/
false
```

#### Destructuring Captures

Use array destructuring to extract capture groups elegantly:

```
>> let email = "john@test.org"
>> let full, name, domain = email ~ /(\w+)@([\w.]+)/
>> log("Name:", name, "Domain:", domain)
Name: "john", "Domain:", "test.org"
```

#### regex() Builtin

Create regexes dynamically from strings:

```
>> let pattern = regex("\\d+")
>> "count: 42" ~ pattern
["42"]
>> let rx = regex("test", "i")  // with flags
>> rx.flags
"i"
```

#### replace() Function

Replace text using strings or regex patterns:

```
>> replace("hello world", "world", "Parsley")
"hello Parsley"
>> replace("test123test456", /\d+/, "XXX")
"testXXXtestXXX"
>> replace("HELLO", /hello/i, "hi")
"hi"
```

#### split() Function

Split strings by delimiter or pattern:

```
>> split("a,b,c", ",")
["a", "b", "c"]
>> split("one1two2three", /\d+/)
["one", "two", "three"]
>> split("hello  world", /\s+/)
["hello", "world"]
```

#### Regex Flags

Parsley supports common regex flags:

- **`i`** - Case-insensitive matching
- **`m`** - Multi-line mode (^ and $ match line boundaries)
- **`s`** - Dot matches newline
- **`g`** - Global (used internally by `replace` and `split`)

Examples:
```
>> "Hello" ~ /hello/
null
>> "Hello" ~ /hello/i
["Hello"]
```

#### Practical Examples

**Email validation:**
```
>> let emailRegex = /^[\w.+-]+@[\w.-]+\.[a-zA-Z]{2,}$/
>> "user@example.com" ~ emailRegex
["user@example.com"]
>> "invalid@" ~ emailRegex
null
```

**URL parsing:**
```
>> let url = "https://example.com/path"
>> let protocol, host, path = url ~ /^(https?):\/\/([^\/]+)(\/.*)?$/
>> log("Protocol:", protocol, "Host:", host)
Protocol: "https", "Host:", "example.com"
```

**Phone number extraction:**
```
>> let phone = "Call (555) 123-4567"
>> let match = phone ~ /\((\d{3})\) (\d{3})-(\d{4})/
>> log("Area:", match[1], "Number:", match[2] + "-" + match[3])
Area: "555", "Number:", "123-4567"
```

**Date parsing:**
```
>> let dateStr = "2024-11-26"
>> let full, year, month, day = dateStr ~ /(\d{4})-(\d{2})-(\d{2})/
>> log("Year:", year, "Month:", month, "Day:", day)
Year: "2024", "Month:", "11", "Day:", "26"
```

**CSV processing:**
```
>> let csv = "apple,banana,cherry"
>> let fruits = split(csv, ",")
>> log("Count:", len(fruits), "Items:", fruits)
Count:3, "Items:", ["apple", "banana", "cherry"]
```

### Module System

Parsley supports a minimalist module system that enables code reuse across files. Modules are just normal Parsley scripts—no special syntax required. All `let` bindings in a module are automatically exported.

#### Basic Module Import

**math.pars:**
```parsley
let PI = 3.14159
let add = fn(a, b) { a + b }
let square = fn(x) { x * x }
```

**Using the module:**
```parsley
let math = import(@./math.pars)
log(math.PI)           // 3.14159
log(math.add(2, 3))    // 5
log(math.square(4))    // 16
```

The `import()` function:
- Takes a path as a string or path literal (`@./file.pars`)
- Returns a dictionary containing all `let` bindings from the module
- Paths are resolved relative to the importing file
- Modules are cached (loaded once, even if imported multiple times)

#### Dictionary Destructuring

Import specific functions or values using destructuring:

```parsley
let {add, square} = import(@./math.pars)
log(add(10, 5))      // 15
log(square(7))       // 49
```

#### Aliasing with `as`

Rename imported items to avoid naming conflicts:

```parsley
let {square as sq, add as plus} = import(@./math.pars)
log(sq(5))           // 25
log(plus(1, 2))      // 3
```

#### Module Caching

Modules are loaded once and cached. Multiple imports return the same module dictionary:

```parsley
let mod1 = import(@./math.pars)
let mod2 = import(@./math.pars)
log(mod1 == mod2)    // true
```

This ensures:
- Efficient loading (files read once)
- Consistent state across imports
- Fast subsequent imports

#### Circular Dependency Detection

Parsley detects circular dependencies and reports errors:

**a.pars:**
```parsley
let b = import(@./b.pars)
let valueA = 1
```

**b.pars:**
```parsley
let a = import(@./a.pars)  // Error: circular dependency
let valueB = 2
```

#### Module Scope Isolation

Each module executes in its own isolated environment. **Only `let` bindings are exported** - variables created with direct assignment are private to the module:

**counter.pars:**
```parsley
count = 0  // Private - not exported (no 'let')
let increment = fn() {
    count = count + 1
    count
}
let getCount = fn() { count }
```

**Using the module:**
```parsley
let counter = import(@./counter.pars)
log(counter.increment())  // 1
log(counter.increment())  // 2
log(counter.count)        // null (not exported - no 'let')
log(counter.getCount())   // 2 (live access)
```

#### Module State and Closures

Exported values are **snapshots** at import time. If you need live access to module state, use functions:

```parsley
// counter.pars
let count = 0

let increment = fn() {
    count = count + 1
    count
}

let getCount = fn() { count }  // Live access through closure

// main.pars
let counter = import(@./counter.pars)
counter.increment()  // 1
counter.increment()  // 2
counter.count        // 0 (snapshot at import time)
counter.getCount()   // 2 (live access)
```

#### Practical Examples

1. **String utilities (strings.pars):**
```parsley
let isEmpty = fn(str) { len(str) == 0 }
let capitalize = fn(str) { toUpper(str[0]) + str[1:] }
let repeat = fn(str, n) {
    if (n <= 0) "" else str + repeat(str, n - 1)
}
```

2. **Email validator (validators.pars):**
```parsley
let emailRegex = /^[\w.+-]+@[\w.-]+\.[a-zA-Z]{2,}$/
let isEmail = fn(str) { str ~ emailRegex }
let isStrongPassword = fn(str) {
    len(str) >= 8 && str ~ /[A-Z]/ && str ~ /[0-9]/
}
```

3. **Using multiple modules:**
```parsley
let {isEmpty, capitalize} = import(@./strings.pars)
let {isEmail} = import(@./validators.pars)

let processEmail = fn(email) {
    if (isEmpty(email)) "Empty email"
    else if (!isEmail(email)) "Invalid email"
    else capitalize(email)
}

log(processEmail("alice@example.com"))  // "Alice@example.com"
log(processEmail(""))                   // "Empty email"
log(processEmail("invalid@"))           // "Invalid email"
```

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
