# 🦁 Parsley

```
█▀█ ▄▀█ █▀█ █▀ █░░ █▀▀ █▄█
█▀▀ █▀█ █▀▄ ▄█ █▄▄ ██▄ ░█░ v 0.14.0
```

A †minimalist concatenative language for generating HTML/XML with first-class file I/O.

- Written in Go
- If JSX and PHP … and Perl … had a cool baby
- Based on Basil from 2001

<sub><sup>†*minimalism not guaranteed…*<sup><sub>

## Table of Contents

- [Quick Start](#quick-start)
- [Language Guide](#language-guide)
  - [Variables and Functions](#variables-and-functions)
  - [Data Types](#data-types)
  - [Control Flow](#control-flow)
  - [HTML/XML Tags](#htmlxml-tags)
  - [File I/O](#file-io)
  - [SFTP](#sftp)
  - [Database](#database)
  - [HTTP Requests](#http-requests)
  - [Process Execution](#process-execution)
  - [Modules](#modules)
- [Examples](#examples)
- [Development](#development)
- [Reference](#reference)

## Quick Start

### Installation

```bash
git clone https://github.com/sambeau/parsley.git
cd parsley
go build -o pars .
```

### Hello World

```bash
./pars                              # Interactive REPL
./pars hello.pars                   # Run a file
./pars --pretty page.pars           # Pretty-print HTML output
```

### Your First Template

```parsley
let name = "World"

let Page = fn({title, contents}) {
    <!DOCTYPE html> + <html>
        <head><title>{title}</title></head>
        <body>
            <h1>{title}</h1>
            {contents}
        </body>
    </html>
}

<Page title="Hello!">
    <p>Welcome to Parsley.</p>
    <p>Generated at {now().format("long")}</p>
</Page>
```

## Language Guide

### Variables and Functions

```parsley
// Variables
let name = "Alice"
let count = 42

// Destructuring
let {x, y} = {x: 10, y: 20}
let a, b, c = 1, 2, 3

// Functions
let greet = fn(name) { "Hello, " + name + "!" }
let add = fn(a, b) { a + b }

// Closures
let makeCounter = fn() {
    count = 0
    fn() { count = count + 1; count }
}
let counter = makeCounter()
counter()  // 1
counter()  // 2
```

### Data Types

```parsley
// Primitives
42                    // Integer
3.14                  // Float
"hello {name}"        // String with interpolation
true, false           // Boolean
null                  // Null

// Collections
[1, 2, 3]             // Array
{name: "Sam", age: 57} // Dictionary

// Special types
/\w+@\w+\.\w+/        // Regex
@2024-12-25           // Date
@2024-12-25T14:30:00  // DateTime
@12:30                // Time
@1d2h30m              // Duration
@./config.json        // Path
@https://example.com  // URL
```

#### Strings

```parsley
let name = "World"
"Hello, {name}!"              // Interpolation

"hello".toUpper()               // "HELLO"
"a,b,c".split(",")            // ["a", "b", "c"]
"hello"[1:4]                  // "ell" (slicing)
```

#### Arrays

```parsley
let nums = [1, 2, 3]

nums[0]                       // 1
nums[-1]                      // 3 (last)
nums[1:]                      // [2, 3] (slice)

nums.length()                 // 3
nums.sort()                   // [1, 2, 3]
nums.map(fn(x) { x * 2 })     // [2, 4, 6]
nums.filter(fn(x) { x > 1 })  // [2, 3]

let range = 1..10             // [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
for (i in 1..5) { log(i) }    // Loop from 1 to 5

[1, 2] ++ [3, 4]              // [1, 2, 3, 4]
[1, 2, 3] && [2, 3, 4]        // [2, 3] (intersection)
[1, 2] || [2, 3]              // [1, 2, 3] (union)
[1, 2, 3] - [2]               // [1, 3] (subtraction)
["hi"] * 3                    // [hi, hi, hi]  (repetition)
```

#### Dictionaries

```parsley
let user = {
    name: "Sam",
    age: 57,
    greet: fn() { "Hi, " + this.name }
}

user.name                     // "Sam"
user["age"]                   // 57
user.greet()                  // "Hi, Sam"

user.keys()                   // ["name", "age", "greet"]
user.values()                 // ["Sam", 57, fn]
user.has("name")              // true

{a: 1} ++ {b: 2}              // {a: 1, b: 2} (merge)
{a: 1, b: 2} && {b: 3, c: 4}  // {b: 2} (intersection)
{a: 1, b: 2} - {b: 0}         // {a: 1} (subtract keys)
```

#### Numbers

```parsley
1234567.format()              // "1,234,567"
99.99.currency("USD")         // "$99.99"
0.15.percent()                // "15%"

sqrt(16)                      // 4
round(3.7)                    // 4
```

#### Dates and Durations

```parsley
let dt = now()
dt.year, dt.month, dt.day     // Components
dt.format("long")             // "November 28, 2024"
dt.format("long", "de-DE")    // "28. November 2024"

@2024-12-25                   // Date
@1d                           // Duration: 1 day
@-1d.format()                 // "yesterday"

// Interpolated datetime templates
let month = "06"
let day = "15"
let dt2 = @(2024-{month}-{day})    // Builds date from variables
dt2.month                          // 6
dt2.day                            // 15

// Time templates
let hour = "14"
let meeting = @({hour}:30)         // Creates time-only value
meeting.hour                       // 14
```

### Control Flow

```parsley
// If expression
let status = if (age >= 18) "adult" else "minor"

// For loops with map/filter
let doubled = for (n in [1, 2, 3]) { n * 2 }      // [2, 4, 6]

let evens = for (n in [1, 2, 3, 4]) {
    if (n % 2 == 0) { n }                          // [2, 4]
}

// With index
for (i, item in items) {
    <li>{i + 1}. {item}</li>
}

// Dictionary iteration
for (key, value in dict) {
    <dt>{key}</dt><dd>{value}</dd>
}
```

### HTML/XML Tags

```parsley
// Basic tags
<div class="container">
    <h1>{title}</h1>
    <p>{content}</p>
</div>

// Self-closing
<img src="photo.jpg" alt="Photo" />
<br/>

// Components (uppercase)
let Card = fn({title, contents}) {
    <article class="card">
        <h2>{title}</h2>
        <div class="body">{contents}</div>
    </article>
}

<Card title="Welcome">
    <p>This is the card content.</p>
</Card>

// Fragments
<>
    <li>Item 1</li>
    <li>Item 2</li>
</>

// Style/Script tags (use @{} for interpolation)
let accent = "#007bff"
<style>
    h1 { color: @{accent}; }
    .box { border: 1px solid @{accent}; }
</style>
```

### File I/O

Parsley has built-in file handling with format-aware reading and writing.

#### Reading Files

```parsley
// Read with format decoding
let config <== JSON(@./config.json)     // Returns dict
let users <== CSV(@./users.csv)         // Returns array of dicts
let content <== text(@./readme.md)      // Returns string

// Load SVG icons as components
let Arrow <== SVG(@./icons/arrow.svg)   // Returns cleaned SVG string
<button><Arrow/> Next</button>          // Use as component

// Load markdown with frontmatter
let post <== MD(@./blog.md)             // Returns dict with html + metadata
<article>
  <h1>{post.title}</h1>
  <time>{post.date}</time>
  {post.html}
</article>

// Destructure from file
let {name, version} <== JSON(@./package.json)

// Error handling
let {data, error} <== JSON(@./config.json)
if (error) {
    <p class="error">Failed: {error}</p>
} else {
    <pre>{data}</pre>
}

// Fallback with ??
let config <== JSON(@./config.json) ?? {theme: "light"}
```

#### Writing Files

```parsley
// Write with format encoding
userData ==> JSON(@./output.json)
records ==> CSV(@./export.csv)
"Hello" ==> text(@./greeting.txt)
"<svg>...</svg>" ==> SVG(@./icon.svg)

// Append
logEntry ==>> lines(@./app.log)
```

#### Stdin/Stdout/Stderr

Read from stdin and write to stdout/stderr for Unix pipeline integration:

```parsley
// @- is the Unix convention: stdin for reads, stdout for writes
let data <== JSON(@-)    // Read JSON from stdin
data ==> JSON(@-)        // Write JSON to stdout

// Explicit aliases
let input <== text(@stdin)
"error" ==> text(@stderr)

// Pipeline example: filter active items
let data <== JSON(@-)
let active = for (item in data.items) {
    if (item.active) { item }
}
active ==> JSON(@-)
```

Run with: `cat input.json | pars filter.pars > output.json`

#### File Operations

```parsley
// Remove/delete files
let f = file(@./temp.txt)
f.remove()  // Returns null on success

// Create directories
file(@./new-dir).mkdir()
file(@./parent/child).mkdir({parents: true})  // Create nested

// Remove directories
file(@./old-dir).rmdir()
file(@./dir-tree).rmdir({recursive: true})  // Remove with contents

// With error handling
if (f.exists) {
    f.remove()
}
```

#### Directories and File Patterns

```parsley
let d = dir(@./images)
d.exists                      // true
d.count                       // Number of files

let files <== dir(@./images)
for (f in files) {
    <p>{f.basename}: {f.size} bytes</p>
}

// File patterns
let images = files(@./images/*.jpg)
for (img in images) {
    <img src="{img.path}" />
}
```

### SFTP

Transfer files securely over SSH using SFTP with the same intuitive syntax as local file I/O.

```parsley
// Connect with SSH key (recommended)
let conn = SFTP(@sftp://user@example.com, {
    keyFile: @~/.ssh/id_rsa,
    timeout: @10s
})

// Or with password
let conn = SFTP(@sftp://user:password@example.com)

// Read remote files
let {config, error} <=/= conn(@/remote/config.json).json
if (!error) {
    log("Config loaded:", config.name)
}

let {logs, err} <=/= conn(@/var/log/app.log).lines

// Write remote files
myData =/=> conn(@/remote/output.json).json
"Hello SFTP" =/=> conn(@/messages/greeting.txt).text

// Append to remote files
logEntry =/=>> conn(@/var/log/app.log).lines

// Directory operations
let {files, dirErr} <=/= conn(@/uploads).dir
if (!dirErr) {
    for (file in files) {
        log(file.name, "-", file.size, "bytes")
    }
}

conn(@/data/archive).mkdir({parents: true})
conn(@/temp/old-file.txt).remove()

// Close when done
conn.close()
```

**Supported formats**: `.json`, `.text`, `.csv`, `.lines`, `.bytes`, `.file`, `.dir`  
**Authentication**: SSH keys (preferred), passwords  
**Features**: Connection caching, timeout support, error capture pattern

See [examples/sftp_demo.pars](examples/sftp_demo.pars) for complete examples.

### Database

Parsley has first-class SQLite database support with a clean, expressive syntax.

```parsley
// Create connection
let db = SQLITE(":memory:")  // or SQLITE(@./data.db)

// Execute DDL/DML (returns {affected, lastId})
let _ = db <=!=> "CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT, email TEXT)"
let result = db <=!=> "INSERT INTO users (name, email) VALUES ('Alice', 'alice@example.com')"

// Query single row (returns dict or null)
let user = db <=?=> "SELECT * FROM users WHERE id = 1"

// Query multiple rows (returns array of dicts)
let users = db <=??=> "SELECT * FROM users WHERE age > 25"

// Transactions
db.begin()
let _ = db <=!=> "INSERT INTO users (name) VALUES ('Bob')"
let _ = db <=!=> "INSERT INTO posts (user_id) VALUES (2)"
db.commit()  // or db.rollback()

// Connection methods
db.ping()      // Test connection
db.close()     // Close connection
```

See [examples/database_demo.pars](examples/database_demo.pars) for a complete working example.

### HTTP Requests

Fetch content from URLs using the `<=/=` operator with format factories:

```parsley
// Fetch JSON from API
let users <=/= JSON(@https://api.example.com/users)
log(users[0].name)

// Fetch plain text
let html <=/= text(@https://example.com)

// POST request with body
let response <=/= JSON(@https://api.example.com/users, {
    method: "POST",
    body: {name: "Alice", email: "alice@example.com"}
})

// Error handling with destructuring
let {data, error, status} <=/= JSON(@https://api.example.com/data)
if (error != null) {
    log("Fetch failed:", error)
} else if (status == 200) {
    log("Success:", data)
}

// Custom timeout
let data <=/= JSON(@https://slow-api.com/data, {
    timeout: 10000  // 10 seconds (default: 30000)
})
```

**Format factories**: `JSON()`, `text()`, `YAML()`, `lines()`, `bytes()`  
**HTTP methods**: GET (default), POST, PUT, PATCH, DELETE, HEAD, OPTIONS  
**Response fields**: `data`, `error`, `status`, `headers`

See [examples/fetch_demo.pars](examples/fetch_demo.pars) for complete examples.

### Process Execution

Execute external commands and process their output:

```parsley
// Simple command execution
let result = COMMAND("ls", ["-la"]) <=#=> null
log(result.stdout)    // Command output
log(result.exitCode)  // 0 for success

// Command with options
let result = COMMAND("node", ["script.js"], {
    env: {NODE_ENV: "production"},
    dir: "/path/to/project",
    timeout: @30s
}) <=#=> null

// Format conversion helpers
let data = parseJSON("{\"name\":\"Alice\"}")
log(data.name)  // "Alice"

let json = stringifyJSON({x: 1, y: 2})
log(json)  // {"x":1,"y":2}

let rows = parseCSV("a,b,c\n1,2,3", {header: true})
log(rows[0].a)  // "1"
```

**Security**: Process execution requires `--allow-execute-all` or `--allow-execute=PATHS` flag.

See [examples/process_demo.pars](examples/process_demo.pars) for complete examples.

### Modules

Parsley supports importing and organizing code with modules.

#### Creating Modules

```parsley
// validators.pars
export let emailRegex = /^[\w.+-]+@[\w.-]+\.[a-zA-Z]{2,}$/
export let isEmail = fn(str) { str ~ emailRegex != null }

// Private helper (not exported)
privateHelper = fn(x) { x * 2 }
```

Use `export` to explicitly mark exported bindings. Variables without `export` or `let` are module-private. Note: `let` bindings without `export` are still exported for backward compatibility.

#### Importing Modules

```parsley
// Import with destructuring
let {isEmail} = import(@./validators.pars)

if (isEmail(userInput)) {
    <p class="valid">Email is valid</p>
}

// Import entire module
let validators = import(@./validators.pars)
validators.isEmail(userInput)
```

#### Standard Library

Parsley includes a growing standard library in the `std/` directory.

## Examples

### Blog Template

```parsley
let posts <== JSON(@./posts.json) ?? []

let PostCard = fn({post}) {
    <article class="post">
        <h2><a href="/posts/{post.slug}">{post.title}</a></h2>
        <time>{time(post.date).format("long")}</time>
        <p>{post.excerpt}</p>
    </article>
}

let BlogPage = fn({title, contents}) {
    <!DOCTYPE html> + <html lang="en">
        <head>
            <meta charset="utf-8" />
            <title>{title}</title>
            <style>
                body { font-family: system-ui; max-width: 800px; margin: 2em auto; }
                .post { margin-bottom: 2em; padding-bottom: 1em; border-bottom: 1px solid #eee; }
                time { color: #666; font-size: 0.9em; }
            </style>
        </head>
        <body>
            <header><h1>{title}</h1></header>
            <main>{contents}</main>
        </body>
    </html>
}

<BlogPage title="My Blog">
    {for (post in posts) {
        <PostCard post={post} />
    }}
</BlogPage>
```

### Data Dashboard

```parsley
let {data, error} <== CSV(@./sales.csv)

let Dashboard = fn({title, contents}) {
    <!DOCTYPE html> + <html>
        <head>
            <title>{title}</title>
            <style>
                table { border-collapse: collapse; width: 100%; }
                th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
                th { background: #f5f5f5; }
                .total { font-weight: bold; background: #e8f4e8; }
                .error { color: red; padding: 1em; background: #fee; }
            </style>
        </head>
        <body>{contents}</body>
    </html>
}

if (error) {
    <Dashboard title="Error">
        <p class="error">Failed to load data: {error}</p>
    </Dashboard>
} else {
    // Calculate total
    let total = 0
    for (row in data) { total = total + toFloat(row.amount ?? "0") }

    <Dashboard title="Sales Report">
        <h1>Sales Report</h1>
        <p>Generated: {now().format("long")}</p>
        
        <table>
            <thead>
                <tr><th>Date</th><th>Product</th><th>Amount</th></tr>
            </thead>
            <tbody>
                for (row in data) {
                    <tr>
                        <td>{row.date}</td>
                        <td>{row.product}</td>
                        <td>{toFloat(row.amount).currency("USD")}</td>
                    </tr>
                }
                <tr class="total">
                    <td colspan="2">Total</td>
                    <td>{total.currency("USD")}</td>
                </tr>
            </tbody>
        </table>
    </Dashboard>
}
```

### Static Site Generator

```parsley
// Generate pages from markdown files
let pages = files(@./content/*.md)

for (page in pages) {
    let content <== text(page)
    let slug = page.stem
    
    let html = <html>
        <head><title>{slug}</title></head>
        <body>
            <article>{content}</article>
        </body>
    </html>
    
    toString(html) ==> text(@./dist/{slug}.html)
}

log("Generated", pages.length(), "pages")
```

## Reference

For complete API documentation, see [docs/reference.md](docs/reference.md).

### Quick Reference

| Type | Methods |
|------|---------|
| String | `.toUpper()` `.toLower()` `.trim()` `.split()` `.replace()` `.length()` |
| Array | `.length()` `.sort()` `.reverse()` `.map()` `.filter()` `.format()` |
| Dictionary | `.keys()` `.values()` `.has()` |
| Number | `.format()` `.currency()` `.percent()` |
| Datetime | `.format()` + properties: `.year` `.month` `.day` `.hour` etc. |
| Duration | `.format()` |
| Path | `.basename` `.ext` `.stem` `.dirname` `.isAbsolute()` |
| URL | `.scheme` `.host` `.path` `.query` `.origin()` |

### Operators

| Op | Description |
|----|-------------|
| `??` | Nullish coalescing: `value ?? default` |
| `~` | Regex match: `str ~ /pattern/` |
| `<==` | Read file: `let data <== JSON(@./file.json)` |
| `==>` | Write file: `data ==> JSON(@./out.json)` |
| `==>>` | Append file: `line ==>> text(@./log.txt)` |
| `++` | Concatenate: `[1] ++ [2]` or `{a:1} ++ {b:2}` |

---

## HTML/XML Tags

Generate structured markup with tag literals.

### Basic Tags

```parsley
// Self-closing tags
let icon = <img src="logo.png" alt="Logo">

// Tags with content
let heading = <h1>"Welcome"</h1>

// Nested tags
let nav = <nav class="main">
  <a href="/">"Home"</a>
  <a href="/about">"About"</a>
</nav>
```

### Dynamic Content

```parsley
let userName = "Alice"
let isAdmin = true

let cls = if (isAdmin) "admin" else "user"
let badge = <span class={cls}>{userName}</span>
```

### Generating Lists

```parsley
let items = ["Apple", "Banana", "Cherry"]

let list = <ul>
  {for (item in items) {
    <li>{item}</li>
  }}
</ul>
```

### Tag Factories

```parsley
// Create tags programmatically
let card = tag("div", { class: "card" }, "Content")
toString(card)  // <div class="card">Content</div>

// SVG elements
let svg = <svg viewBox="0 0 100 100">
  <circle cx="50" cy="50" r="40" fill="blue" />
</svg>
```

---

## Error Handling

### Error Capture Pattern

Capture errors instead of halting execution:

```parsley
// Wrap in {data, error} to capture errors
let {data, error} <== JSON(@./config.json)

if error {
  log("Failed to load config:", error)
  let data = { defaults: true }
}
```

### Validation with Regex

```parsley
fn validateEmail(email) {
  if !(email ~ /^[\w.-]+@[\w.-]+\.\w+$/) {
    { valid: false, error: "Invalid email format" }
  } else {
    { valid: true, email: email }
  }
}

let result = validateEmail("test@example.com")
```

### Nullish Coalescing

```parsley
// Use ?? for fallback values
let config <== JSON(@./config.json)
let port = config.port ?? 8080
let host = config.host ?? "localhost"
```

---

## Localization

Format numbers, currencies, and dates for different locales.

### Number Formatting

```parsley
let price = 1234567.89

// With locale
price.format("en-US")      // "1,234,567.89"
price.format("de-DE")      // "1.234.567,89"

// Currency
price.currency("USD", "en-US")  // "$1,234,567.89"
price.currency("EUR", "de-DE")  // "1.234.567,89 €"

// Percentage
let rate = 0.156
rate.percent("en-US", 1)   // "15.6%"
```

### Date Formatting

```parsley
let date = now()

date.format("short", "en-US")   // "1/15/25"
date.format("medium", "en-GB")  // "15 Jan 2025"
date.format("long", "fr-FR")    // "15 janvier 2025"
```

---

## Complete Example: Static Site Generator

```parsley
// site-generator.pars
// Generate a static blog site from markdown-style data

let site = {
  title: "My Blog",
  author: "Alice",
  baseUrl: "https://blog.example.com"
}

let posts = [
  { slug: "hello", title: "Hello World", date: "2025-01-15", body: "First post!" },
  { slug: "update", title: "Big Update", date: "2025-01-20", body: "New features..." }
]

// Generate HTML page
let renderPage = fn(title, pageContent) {
  <html lang="en">
    <head>
      <meta charset="UTF-8" />
      <title>{title} | {site.title}</title>
      <link rel="stylesheet" href="/style.css" />
    </head>
    <body>
      <header>
        <h1>{site.title}</h1>
        <nav>
          <a href="/">"Home"</a>
          <a href="/posts">"Posts"</a>
        </nav>
      </header>
      <main>{pageContent}</main>
      <footer>
        "© 2025 " {site.author}
      </footer>
    </body>
  </html>
}

// Generate post page
let renderPost = fn(post) {
  let pageContent = <article>
    <h1>{post.title}</h1>
    <time>{post.date}</time>
    <div class="body">{post.body}</div>
  </article>
  
  renderPage(post.title, pageContent)
}

// Generate index page
let renderIndex = fn() {
  let pageContent = <div>
    <h1>"Recent Posts"</h1>
    <ul class="post-list">
      {for (post in posts) {
        <li>
          <a href={"/posts/" ++ post.slug}>{post.title}</a>
          <time>{post.date}</time>
        </li>
      }}
    </ul>
  </div>
  
  renderPage("Home", pageContent)
}

// Output to stdout
renderIndex()
```

Run with: `./pars site-generator.pars`

---

## Security

Parsley provides command-line flags to restrict file system access, protecting against unauthorized file operations and malicious scripts.

### Security Model

- **Read operations**: Unrestricted by default (opt-in blacklist)
- **Write operations**: Denied by default (opt-in whitelist)
- **Execute operations**: Denied by default (opt-in whitelist for module imports)

### Security Flags

```bash
# Read restrictions
--restrict-read=PATHS    # Deny reading from paths (blacklist)
--no-read                # Deny all file reads

# Write permissions
--allow-write=PATHS      # Allow writes to specific paths
--allow-write-all, -w    # Allow unrestricted writes

# Execute permissions  
--allow-execute=PATHS    # Allow importing modules from paths
--allow-execute-all, -x  # Allow unrestricted module imports
```

### Examples

```bash
# Static site generator: read freely, write to output
./pars --allow-write=./public build.pars

# API processor: restrict sensitive reads, write results
./pars --restrict-read=/etc --allow-write=./output process.pars

# Development mode: unrestricted (old behavior)
./pars -w -x dev-script.pars

# Production: specific permissions
./pars --allow-write=./data --allow-execute=./lib app.pars
```

### Migration from v0.9.x

In v0.10.0, write and execute operations are denied by default:

```bash
# Old (v0.9.x) - everything allowed
./pars build.pars

# New (v0.10.0) - writes denied by default
./pars build.pars  # ERROR: write access denied

# Fix: Add -w for unrestricted writes (old behavior)
./pars -w build.pars

# Or: Specify allowed paths
./pars --allow-write=./output build.pars
```

---

## Development

### Building from Source

```bash
# Clone repository
git clone https://github.com/username/parsley.git
cd parsley

# Build with make
make build              # Build binary
make test               # Run tests
make install            # Install to $GOPATH/bin

# Or build manually
go build -ldflags "-X main.Version=$(cat VERSION)" -o pars .
```

### Testing

```bash
go test ./...                    # All tests
go test ./pkg/evaluator -v       # Specific package
go test -cover ./...             # With coverage
```

### Running Parsley

```bash
./pars                           # Interactive REPL
./pars script.pars               # Execute file
./pars --pretty page.pars        # Pretty-print HTML output
./pars --version                 # Show version
```

### Project Structure

```
parsley/
├── main.go              # Entry point
├── pkg/
│   ├── ast/             # Abstract Syntax Tree
│   ├── lexer/           # Tokenizer
│   ├── parser/          # Parser
│   ├── evaluator/       # Interpreter
│   └── repl/            # Interactive mode
├── std/                 # Standard library
├── examples/            # Example scripts
└── docs/                # Documentation
    └── reference.md     # Full API reference
```

---

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Write tests for new functionality
4. Ensure all tests pass (`go test ./...`)
5. Submit a pull request

---

## License

MIT License - see [LICENSE](LICENSE) file for details.
