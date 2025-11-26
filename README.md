# Parsley

```
█▀█ ▄▀█ █▀█ █▀ █░░ █▀▀ █▄█
█▀▀ █▀█ █▀▄ ▄█ █▄▄ ██▄ ░█░ v 0.2.2
```

A concatenative programming language interpreter.
- Writen in Go
- Similar to JSX

## Features

### Core Language Features
- Variable declarations with `let`
- Direct variable assignment (e.g., `x = 5`)
- Array destructuring assignment (e.g., `x,y,z = 1,2,3`)
- Functions with `fn`
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
- Special `_` variable (write-only, always returns `null`)

### Data Types
- **Integers:** `42`, `-15`
- **Floats:** `3.14159`, `2.718`
- **Strings:** `"hello world"`
- **Booleans:** `true`, `false`
- **Arrays:** `1,2,3`, `[1,2,3]`, `[[1,2],[3,4]]`, mixed types allowed
- **Dictionaries:** `{ name: "Sam", age: 57 }`, key-value pairs with lazy evaluation

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

Supported escape sequences in strings:
- `\n` - newline
- `\t` - tab
- `\\` - backslash
- `\"` - double quote

Template literals support additional escape sequences:
- `\`` - backtick (literal backtick character)
- `\{` - left brace (prevents interpolation)
- `\}` - right brace (literal closing brace)

#### Template Literals

Template literals use backticks and support expression interpolation:

```
>> `Hello, World!`
Hello, World!
>> name = "Sam"
Sam
>> `Welcome, {name}!`
Welcome, Sam!
```

Interpolate any expression with `{}`:

```
>> a = 5
5
>> b = 10
10
>> `Sum: {a + b}`
Sum: 15
>> `Result: {a * 2 + b}`
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

Escape special characters in templates:

```
>> `Literal backtick: \``
Literal backtick: `
>> `Not interpolated: \{variable}`
Not interpolated: {variable}
```

String concatenation with automatic type conversion:

```
>> "Count: " + 42
Count: 42
>> "Result: " + (5 + 3)
Result: 8
```

### Singleton Tags

Singleton tags provide a convenient syntax for generating HTML/XML markup. Tags are self-closing and use the `<tagname ... />` syntax.

#### Standard Tags (Lowercase)

Standard tags have lowercase names and are evaluated as interpolated strings that produce HTML/XML output:

```
>> <br/>
<br />
>> <meta charset="utf-8" />
<meta charset="utf-8"  />
>> <img src="photo.jpg" width="300" height="200" />
<img src="photo.jpg" width="300" height="200"  />
```

#### Tag Interpolation

Tags support expression interpolation using `{expr}` syntax:

```
>> charset = "utf-8"
utf-8
>> <meta charset="{charset}" />
<meta charset="utf-8"  />

>> width = 300
300
>> height = 200
200
>> <img width="{width}" height="{height}" />
<img width="300" height="200"  />
```

Interpolate any expression:

```
>> x = 10
10
>> <div data-value="{x * 2}" />
<div data-value="20"  />

>> disabled = true
true
>> <button disabled="{if(disabled){"disabled"}}" />
<button disabled="disabled"  />
```

#### Boolean Props

Standalone props without values are treated as boolean attributes:

```
>> <input type="checkbox" checked />
<input type="checkbox" checked  />
>> <button type="submit" disabled />
<button type="submit" disabled  />
```

#### Multiline Tags

Tags can span multiple lines for better readability:

```
>> <img
   src="https://example.com/image.png"
   width="{300}"
   height="{200}"
   alt="Example Image" />
<img
   src="https://example.com/image.png"
   width="300"
   height="200"
   alt="Example Image"  />
```

#### Custom Tags (Uppercase/TitleCase)

Custom tags have names starting with an uppercase letter and are treated as function calls. The tag props are passed as a dictionary to the function:

```
>> Dog = fn(props) {
   name = props.name
   age = props.age
   toString("Dog: ", name, ", Age: ", age)
}
>> <Dog name="Rover" age="5" />
Dog: Rover, Age: 5
```

Interpolate expressions in custom tag props:

```
>> Card = fn(props) {
   title = props.title
   content = props.content
   `<div class="card">
     <h2>{title}</h2>
     <p>{content}</p>
   </div>`
}
>> <Card title="Welcome" content="Hello World" />
<div class="card">
     <h2>Welcome</h2>
     <p>Hello World</p>
   </div>
```

Custom tags with computed values:

```
>> Double = fn(props) {
   value = props.value
   value * 2
}
>> <Double value="{10 + 5}" />
30
```

Boolean props in custom tags:

```
>> Button = fn(props) {
   isDisabled = has(props, "disabled")
   if (isDisabled) {
     "Button is disabled"
   } else {
     "Button is enabled"
   }
}
>> <Button disabled />
Button is disabled
>> <Button type="submit" />
Button is enabled
```

#### Practical Tag Examples

Generate HTML components:

```
>> Link = fn(props) {
   url = props.url
   text = props.text
   `<a href="{url}">{text}</a>`
}
>> <Link url="https://example.com" text="Click here" />
<a href="https://example.com">Click here</a>
```

Build reusable UI components:

```
>> Alert = fn(props) {
   type = props.type
   message = props.message
   `<div class="alert alert-{type}">{message}</div>`
}
>> <Alert type="warning" message="Please save your work" />
<div class="alert alert-warning">Please save your work</div>
```

Tags work seamlessly with other Parsley features:

```
>> tags = [<br/>, <hr/>]
>> tags[0]
<br />

>> toString(<br/>, <hr/>, <br/>)
<br /><hr /><br />
```

### Tag Pairs

Tag pairs provide opening and closing tags with content between them, enabling the creation of complete HTML documents and components.

#### Basic Tag Pairs

Tag pairs use `<tag>content</tag>` syntax with text, interpolations, and nested tags:

```
>> <div>Hello, World!</div>
<div>Hello, World!</div>

>> <p>This is a paragraph.</p>
<p>This is a paragraph.</p>

>> name = "Alice"
Alice
>> <h1>Welcome, {name}!</h1>
<h1>Welcome, Alice!</h1>
```

#### Nested Tags

Tags can be nested to create complex structures:

```
>> <div><p>Nested content</p></div>
<div><p>Nested content</p></div>

>> <article><h1>Title</h1><p>Content goes here</p></article>
<article><h1>Title</h1><p>Content goes here</p></article>
```

#### Interpolation in Tag Content

Use `{expr}` to interpolate expressions within tag content:

```
>> x = "First"
First
>> y = "Second"
Second
>> <div>{x} - {y}</div>
<div>First - Second</div>

>> count = 5
5
>> <p>You have {count} new messages.</p>
<p>You have 5 new messages.</p>
```

#### Empty Grouping Tags

Use `<>...</>` to group content without adding wrapper tags:

```
>> <>Hello</>
Hello

>> <><div>First</div><div>Second</div></>
<div>First</div><div>Second</div>
```

#### Creating HTML Documents

Build complete, valid HTML pages:

```
>> Page = fn(props) {
   title = props.title
   content = props.content
   <html>
     <head>
       <title>{title}</title>
       <meta charset="utf-8" />
     </head>
     <body>
       <h1>{title}</h1>
       <div>{content}</div>
     </body>
   </html>
}

>> <Page title="My Site" content="Welcome!" />
<html>
  <head>
    <title>My Site</title>
    <meta charset="utf-8" />
  </head>
  <body>
    <h1>My Site</h1>
    <div>Welcome!</div>
  </body>
</html>
```

#### HTML Components with Contents

Components can receive content via `props.contents`:

```
>> Card = fn(props) {
   <div class="card">
     <div class="card-body">
       {props.contents}
     </div>
   </div>
}

>> <Card><h3>Title</h3><p>Description</p></Card>
<div class="card">
  <div class="card-body">
    <h3>Title</h3><p>Description</p>
  </div>
</div>
```

Navigation menu component:

```
>> Nav = fn(props) {
   <nav>
     <ul>
       {props.contents}
     </ul>
   </nav>
}

>> <Nav>
   <li><a href="/">Home</a></li>
   <li><a href="/about">About</a></li>
   <li><a href="/contact">Contact</a></li>
</Nav>
<nav>
  <ul>
    <li><a href="/">Home</a></li>
    <li><a href="/about">About</a></li>
    <li><a href="/contact">Contact</a></li>
  </ul>
</nav>
```

#### HTML Components with Props and Contents

Combine props and contents for flexible components:

```
>> Section = fn(props) {
   <section class="{props.theme}">
     <h2>{props.title}</h2>
     <div class="content">
       {props.contents}
     </div>
   </section>
}

>> <Section title="Welcome" theme="dark">
   <p>This is the section content.</p>
   <p>It can contain multiple paragraphs.</p>
</Section>
<section class="dark">
  <h2>Welcome</h2>
  <div class="content">
    <p>This is the section content.</p>
    <p>It can contain multiple paragraphs.</p>
  </div>
</section>
```

Article with metadata:

```
>> Article = fn(props) {
   <article>
     <header>
       <h1>{props.title}</h1>
       <p class="meta">By {props.author} on {props.date}</p>
     </header>
     <div class="body">
       {props.contents}
     </div>
   </article>
}

>> <Article title="Getting Started" author="Sam" date="2025-11-26">
   <p>This is the first paragraph.</p>
   <p>This is the second paragraph.</p>
</Article>
<article>
  <header>
    <h1>Getting Started</h1>
    <p class="meta">By Sam on 2025-11-26</p>
  </header>
  <div class="body">
    <p>This is the first paragraph.</p>
    <p>This is the second paragraph.</p>
  </div>
</article>
```

#### SVG Components

Create scalable vector graphics with tag pairs:

```
>> Circle = fn(props) {
   <circle 
     cx="{props.x}" 
     cy="{props.y}" 
     r="{props.radius}" 
     fill="{props.color}" />
}

>> <svg width="200" height="200">
   <Circle x="100" y="100" radius="50" color="blue" />
   <Circle x="150" y="100" radius="30" color="red" />
</svg>
<svg width="200" height="200">
  <circle 
    cx="100" 
    cy="100" 
    r="50" 
    fill="blue" />
  <circle 
    cx="150" 
    cy="100" 
    r="30" 
    fill="red" />
</svg>
```

Complete SVG icon component:

```
>> Icon = fn(props) {
   size = props.size
   <svg width="{size}" height="{size}" viewBox="0 0 24 24">
     <path d="{props.path}" fill="{props.color}" />
   </svg>
}

>> heartPath = "M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"
M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z

>> <Icon size="24" color="red" path="{heartPath}" />
<svg width="24" height="24" viewBox="0 0 24 24">
  <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z" fill="red" />
</svg>
```

Simple chart component:

```
>> BarChart = fn(props) {
   values = props.values
   max = props.max
   
   <svg width="400" height="200">
     <rect x="10" y="{200 - (values[0] / max * 180)}" 
           width="80" height="{values[0] / max * 180}" fill="blue" />
     <rect x="110" y="{200 - (values[1] / max * 180)}" 
           width="80" height="{values[1] / max * 180}" fill="green" />
     <rect x="210" y="{200 - (values[2] / max * 180)}" 
           width="80" height="{values[2] / max * 180}" fill="orange" />
   </svg>
}

>> <BarChart values="[75, 120, 90]" max="150" />
<svg width="400" height="200">
  <rect x="10" y="110" 
        width="80" height="90" fill="blue" />
  <rect x="110" y="56" 
        width="80" height="144" fill="green" />
  <rect x="210" y="92" 
        width="80" height="108" fill="orange" />
</svg>
```

#### Combining Features

Use loops to generate repeated content:

```
>> items = ["Apple", "Banana", "Cherry"]
["Apple", "Banana", "Cherry"]

>> List = fn(props) {
   <ul>
     {for(item in props.items) {
       <li>{item}</li>
     }}
   </ul>
}

>> <List items="{items}" />
<ul>
  <li>Apple</li>
  <li>Banana</li>
  <li>Cherry</li>
</ul>
```

Conditional rendering:

```
>> UserGreeting = fn(props) {
   <div>
     {if(props.loggedIn) {
       <p>Welcome back, {props.name}!</p>
     } else {
       <p>Please log in.</p>
     }}
   </div>
}

>> <UserGreeting loggedIn="{true}" name="Alice" />
<div>
  <p>Welcome back, Alice!</p>
</div>

>> <UserGreeting loggedIn="{false}" name="Alice" />
<div>
  <p>Please log in.</p>
</div>
```

### Type Conversions

Convert strings to numbers:

```
>> toInt("42")
42
>> toFloat("3.14")
3.14
>> toNumber("42")
42
>> toNumber("3.14")
3.14
```

The `toNumber()` function automatically detects whether to return an integer or float based on the string content (presence of decimal point).

### Comments

Single-line comments start with `//`:

```pars
// This is a comment
let x = 5  // Inline comment
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
>> gt10 = fn(x) { if (x > 10) { return x } }
>> map(gt10, 5,15,25,8,3,12)
15, 25, 12
```

Using `for` as a filter by omitting return values:

```
>> for(x in 5,15,25,8,3,12) { if (x > 10) { x } }
15, 25, 12
>> numbers = 1,2,3,4,5,6,7,8,9,10
>> for(n in numbers) { if (n % 2 == 0) { n } }
2, 4, 6, 8, 10
```

When an if statement has no else clause and the condition is false, it returns `null`. Since `for` loops filter out `null` values, this provides a concise way to filter arrays.

### toString() Function

The `toString()` function converts values to strings and joins them without any whitespace:

```
>> toString(1, 2, 3)
123
>> toString("Hello", "World")
HelloWorld
>> xs = (1,2,3)
>> toString(xs)
123
>> toString("Result:", 42)
Result:42
```

### toDebug() Function

The `toDebug()` function converts values to a debug representation with proper formatting:

```
>> toDebug(1, 2.5, "hello", true)
1, 2.5, "hello", true
>> here = "HERE!"
>> xs = 1, 2.0, "Sam", "was", here
>> toDebug(xs)
[1, 2, "Sam", "was", "HERE!"]
>> nested = [[1, 2], ["a", "b"]]
>> toDebug(nested)
[[1, 2], ["a", "b"]]
```

### log() Function

The `log()` function outputs values in debug format immediately to stdout, useful for debugging:

```
log("Starting computation...")
x = 5
log("x is:", x)
// Output: x is: 5

arr = [10, 4, 16, 2, 18]
log("Final result:", arr)
// Output: Final result: [10, 4, 16, 2, 18]

for (item in ["apple", "banana", "cherry"]) {
	log("Processing:", item)
	item
}
// Output during loop execution:
// Processing: "apple"
// Processing: "banana"
// Processing: "cherry"
```

**Special behavior:** If the first argument is a string, it's displayed without quotes and has no comma after it, making it ideal for labels.

**Note:** `log()` returns `null` and outputs immediately, making it ideal for debugging loops and tracking execution flow.

### logLine() Function

The `logLine()` function outputs values with the filename and line number prefix, useful for tracking execution location:

```
logLine("Starting program")
// Output: program.pars:1: Starting program

x = 5
logLine("x is:", x)
// Output: program.pars:3: x is: 5

for (item in ["apple", "banana", "cherry"]) {
	logLine("Processing:", item)
}
// Output during loop execution:
// program.pars:6: Processing: "apple"
// program.pars:6: Processing: "banana"
// program.pars:6: Processing: "cherry"
```

**Special behavior:** Like `log()`, if the first argument is a string, it's displayed without quotes for clean label output.

**Note:** `logLine()` returns `null` and is particularly useful for debugging to understand where in your code execution is happening, especially in loops and nested functions.

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
>> for(x in 5,15,25) { if (x > 10) { x } }
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

### Array Destructuring

Array destructuring allows you to assign multiple variables at once from an array or comma-separated values:

#### Basic Destructuring

Assign multiple variables from values:
```
>> x,y,z = 1,2,3
1, 2, 3
>> x
1
>> y
2
>> z
3
```

Destructure from an existing array:
```
>> xs = 10,20,30
10, 20, 30
>> a,b,c = xs
10, 20, 30
>> a
10
>> b
20
```

#### Tail Collection

When there are more values than variables, the last variable receives all remaining values as an array:
```
>> p,q,r = 1,2,3,4,5,6
1, 2, 3, 4, 5, 6
>> p
1
>> q
2
>> r
3, 4, 5, 6
```

#### Destructuring with `let`

Works with `let` statements:
```
>> let m,n,o = 100,200,300
100, 200, 300
>> m
100
```

#### Using `_` to Ignore Values

Combine destructuring with the `_` variable to ignore unwanted values:
```
>> first,_,third = "A","B","C"
A, B, C
>> first
A
>> third
C
```

#### Head and Tail Functions

Common pattern for list processing:
```
>> head = fn(list) { h,_ = list; h }
>> tail = fn(list) { _,t = list; t }
>> numbers = 1,2,3,4,5
1, 2, 3, 4, 5
>> head(numbers)
1
>> tail(numbers)
2, 3, 4, 5
```

**Note:** When using comma-separated variables in an expression context (not assignment), wrap them in parentheses to avoid ambiguity: `(x,y)` instead of `x,y`.

### Special Variables

#### The `_` Variable

The `_` variable is a special write-only variable that discards any value assigned to it and always evaluates to `null`. This is useful when you need to execute an expression for its side effects but don't care about the result:

```
>> _ = 100
100
>> _
null
>> let _ = "hello"
hello
>> _
null
```

Useful for ignoring values:
```
>> x = 10
>> _ = x * 2  // Calculate but don't store
20
>> _
null
```

### Functions

Functions are first-class objects in pars, meaning they can be assigned to variables, stored in arrays, and passed as arguments:

```
>> let circleArea = fn(r) { pi() * pow(r, 2) }
>> circleArea(10)
314.1592653589793
```

#### Functions as Array Elements

Since functions are first-class objects, they can be stored in arrays and called using indexing:

```
>> double = fn(x) { x + x }
>> square = fn(x) { x * x }
>> funs = double, square
>> funs[0](3)
6
>> funs[1](3)
9
```

You can also use this to create lookup tables of operations:

```
>> ops = fn(a,b){a+b}, fn(a,b){a-b}, fn(a,b){a*b}, fn(a,b){a/b}
>> ops[0](10, 5)
15
>> ops[2](10, 5)
50
```

### Array Sorting

The `sort()` function provides natural sorting for arrays, treating consecutive digits in strings as numbers and comparing them numerically. This creates more intuitive ordering for humans.

#### Sorting Numbers
```
>> xs = 3,2,4,10,1
>> sort(xs)
1, 2, 3, 4, 10
```

#### Natural String Sorting
Natural sort properly handles numbers within strings:
```
>> xs = "z11", "z1", "z2"
>> sort(xs)
z1, z2, z11

>> items = "item 20", "item 3", "item 100", "item 1"
>> sort(items)
item 1, item 3, item 20, item 100
```

#### Mixed Type Sorting
Numbers are sorted before strings:
```
>> xs = "a", 10, "b", 2, 1
>> sort(xs)
1, 2, 10, a, b

>> files = "file 1.txt", "file 10.txt", "file 2.txt", 5, "file 20.txt"
>> sort(files)
5, file 1.txt, file 2.txt, file 10.txt, file 20.txt
```

**Note:** The `sort()` function returns a new sorted array and does not modify the original.

### Array Reversal

The `reverse()` function returns a new array with elements in reverse order:

```
>> xs = 1,2,3,4,5
>> reverse(xs)
5, 4, 3, 2, 1

>> words = "apple", "banana", "cherry"
>> reverse(words)
cherry, banana, apple

>> original = 1,2,3
>> reversed = reverse(original)
>> original
1, 2, 3
>> reversed
3, 2, 1
```

**Note:** The `reverse()` function returns a new reversed array and does not modify the original.

### Custom Sorting with sortBy

The `sortBy()` function allows custom sorting using a comparison function. The comparison function takes two values and returns them in the desired order as a 2-element array:

```
>> normalOrder = fn(a,b){ sort([a,b]) }
>> reverseOrder = fn(a,b){ reverse(sort([a,b])) }

>> normalOrder(20, 10)
10, 20

>> reverseOrder(300, 500)
500, 300

>> sortBy([1,2,3,4,5], reverseOrder)
5, 4, 3, 2, 1

>> sortBy([5, 50, 10, 100, 6, 60, 7, 70], normalOrder)
5, 6, 7, 10, 50, 60, 70, 100

>> sortBy([1,2,3,4,5], fn(a,b){ reverse(sort([a,b])) })
5, 4, 3, 2, 1
```

The comparison function receives two elements and should return a 2-element array with those elements in the desired order. If the first element of the returned array equals the first input, it's considered to come before the second.

**Note:** The `sortBy()` function returns a new sorted array and does not modify the original.

## Error Reporting

Pars provides clear, helpful error messages with:

- **Filename** in the error message
- **Line and column numbers** for precise error location
- **Human-readable descriptions** instead of technical token types
- **Visual pointer** (^) showing the exact error position
- **Source code context** displaying the problematic line

### Error Message Example

**Source file with error:**
```pars
let x = 5
let y =
let z = 10
```

**Error output:**
```
Error in 'example.pars':
  line 2, column 8: unexpected 'let'
    let y =
           ^
```

See [ERROR_DEMO.md](examples/ERROR_DEMO.md) for more examples of error messages.

## Development

This project is structured to easily add new modules and packages:

- `main.go` - Entry point of the application
- `pkg/` - Public packages that can be used by external applications
  - `lexer/` - Tokenizes input into lexical tokens
  - `parser/` - Converts tokens into an Abstract Syntax Tree
  - `ast/` - Defines the Abstract Syntax Tree nodes
  - `evaluator/` - Evaluates the AST and executes the program
  - `repl/` - Read-Eval-Print Loop for interactive usage

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
