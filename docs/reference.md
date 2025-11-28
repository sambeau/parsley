# Parsley Reference

Complete reference for all Parsley types, methods, and operators.

## Table of Contents

- [Data Types](#data-types)
- [Operators](#operators)
- [String Methods](#string-methods)
- [Array Methods](#array-methods)
- [Dictionary Methods](#dictionary-methods)
- [Number Methods](#number-methods)
- [Datetime Methods](#datetime-methods)
- [Duration Methods](#duration-methods)
- [Path Methods](#path-methods)
- [URL Methods](#url-methods)
- [File I/O](#file-io)
- [Regex](#regex)
- [Modules](#modules)
- [Tags](#tags)
- [Utility Functions](#utility-functions)

---

## Data Types

| Type | Example | Description |
|------|---------|-------------|
| Integer | `42`, `-15` | Whole numbers |
| Float | `3.14`, `2.718` | Decimal numbers |
| String | `"hello"`, `"world"` | Text with `{interpolation}` |
| Boolean | `true`, `false` | Logical values |
| Null | `null` | Absence of value |
| Array | `[1, 2, 3]` | Ordered collections |
| Dictionary | `{x: 1, y: 2}` | Key-value pairs |
| Function | `fn(x) { x * 2 }` | First-class functions |
| Regex | `/pattern/flags` | Regular expressions |
| Date/Time | `@2024-11-26` | Temporal values |
| Duration | `@1d`, `@2h30m` | Time spans |
| Path | `@./file.pars` | File system paths |
| URL | `@https://example.com` | Web addresses |
| File Handle | `JSON(@./config.json)` | File with format binding |
| Directory | `dir(@./folder)` | Directory handle |

---

## Operators

### Arithmetic
| Operator | Description | Example |
|----------|-------------|---------|
| `+` | Addition | `2 + 3` → `5` |
| `-` | Subtraction | `5 - 2` → `3` |
| `*` | Multiplication | `4 * 3` → `12` |
| `/` | Division | `10 / 4` → `2.5` |
| `%` | Modulo | `10 % 3` → `1` |
| `++` | Concatenation | `[1] ++ [2]` → `[1, 2]` |

### Comparison
| Operator | Description |
|----------|-------------|
| `==` | Equal |
| `!=` | Not equal |
| `<` | Less than |
| `<=` | Less than or equal |
| `>` | Greater than |
| `>=` | Greater than or equal |

### Logical
| Operator | Description |
|----------|-------------|
| `&&` | AND |
| `\|\|` | OR |
| `!` | NOT |

### Pattern Matching
| Operator | Description | Example |
|----------|-------------|---------|
| `~` | Regex match | `"test" ~ /\w+/` → `["test"]` |
| `!~` | Regex not-match | `"abc" !~ /\d/` → `true` |

### Nullish Coalescing
| Operator | Description | Example |
|----------|-------------|---------|
| `??` | Default if null | `null ?? "default"` → `"default"` |

```parsley
value ?? default           // Returns default only if value is null
null ?? "fallback"         // "fallback"
"hello" ?? "fallback"      // "hello"
0 ?? 42                    // 0 (not null)
a ?? b ?? c ?? "default"   // First non-null value
```

### File I/O
| Operator | Description | Example |
|----------|-------------|---------|
| `<==` | Read from file | `let data <== JSON(@./file.json)` |
| `==>` | Write to file | `data ==> JSON(@./out.json)` |
| `==>>` | Append to file | `line ==>> lines(@./log.txt)` |

### Other
| Operator | Description |
|----------|-------------|
| `=` | Assignment |
| `:` | Key-value separator |
| `.` | Property/method access |
| `[]` | Indexing and slicing |

---

## String Methods

| Method | Description | Example |
|--------|-------------|---------|
| `.length()` | String length | `"hello".length()` → `5` |
| `.upper()` | Uppercase | `"hello".upper()` → `"HELLO"` |
| `.lower()` | Lowercase | `"HELLO".lower()` → `"hello"` |
| `.trim()` | Remove whitespace | `"  hi  ".trim()` → `"hi"` |
| `.split(delim)` | Split to array | `"a,b,c".split(",")` → `["a","b","c"]` |
| `.replace(old, new)` | Replace text | `"hello".replace("l", "L")` → `"heLLo"` |

### Indexing and Slicing
```parsley
"hello"[0]      // "h"
"hello"[-1]     // "o" (last)
"hello"[1:4]    // "ell"
"hello"[2:]     // "llo"
"hello"[:3]     // "hel"
```

### Interpolation
```parsley
let name = "World"
"Hello, {name}!"  // "Hello, World!"
```

---

## Array Methods

| Method | Description | Example |
|--------|-------------|---------|
| `.length()` | Array length | `[1,2,3].length()` → `3` |
| `.sort()` | Sort ascending | `[3,1,2].sort()` → `[1,2,3]` |
| `.reverse()` | Reverse order | `[1,2,3].reverse()` → `[3,2,1]` |
| `.map(fn)` | Transform each | `[1,2].map(fn(x){x*2})` → `[2,4]` |
| `.filter(fn)` | Keep matching | `[1,2,3].filter(fn(x){x>1})` → `[2,3]` |
| `.join()` | Join to string | `["a","b","c"].join()` → `"abc"` |
| `.join(sep)` | Join with separator | `["a","b","c"].join(",")` → `"a,b,c"` |
| `.format()` | List as prose | `["a","b"].format()` → `"a and b"` |
| `.format("or")` | With conjunction | `["a","b"].format("or")` → `"a or b"` |

### Indexing and Slicing
```parsley
nums[0]      // First element
nums[-1]     // Last element
nums[1:3]    // Elements 1 and 2
nums[2:]     // From index 2 to end
nums[:2]     // From start to index 2
```

### Concatenation
```parsley
[1, 2] ++ [3, 4]  // [1, 2, 3, 4]
```

---

## Dictionary Methods

| Method | Description | Example |
|--------|-------------|---------|
| `.keys()` | All keys | `{a:1}.keys()` → `["a"]` |
| `.values()` | All values | `{a:1}.values()` → `[1]` |
| `.has(key)` | Key exists | `{a:1}.has("a")` → `true` |

### Access
```parsley
dict.key        // Dot notation
dict["key"]     // Bracket notation
```

### Self-Reference with `this`
```parsley
let config = {
    width: 100,
    height: 200,
    area: this.width * this.height  // Computed on access
}
```

### Merging
```parsley
{a: 1} ++ {b: 2}  // {a: 1, b: 2}
```

---

## Number Methods

| Method | Description | Example |
|--------|-------------|---------|
| `.format()` | Locale format | `1234567.format()` → `"1,234,567"` |
| `.format(locale)` | With locale | `1234.format("de-DE")` → `"1.234"` |
| `.currency(code)` | Currency format | `99.currency("USD")` → `"$99.00"` |
| `.currency(code, locale)` | With locale | `99.currency("EUR","de-DE")` → `"99,00 €"` |
| `.percent()` | Percentage | `0.125.percent()` → `"13%"` |

### Math Functions
```parsley
sqrt(16)        // 4
round(3.7)      // 4
pow(2, 8)       // 256
pi()            // 3.14159...
sin(x), cos(x), tan(x)
asin(x), acos(x), atan(x)
```

---

## Datetime Methods

### Creation
```parsley
now()                                    // Current time
time("2024-11-26")                       // Parse ISO date
time("2024-11-26T15:30:00")              // With time
time(1732579200)                         // Unix timestamp
time({year: 2024, month: 12, day: 25})   // From components
```

### Literals
```parsley
@2024-11-26          // Date
@2024-11-26T15:30    // DateTime
```

### Properties
| Property | Description |
|----------|-------------|
| `.year` | Year number |
| `.month` | Month (1-12) |
| `.day` | Day of month |
| `.hour` | Hour (0-23) |
| `.minute` | Minute (0-59) |
| `.second` | Second (0-59) |
| `.weekday` | Day name ("Monday", etc.) |
| `.iso` | ISO 8601 string |
| `.unix` | Unix timestamp |
| `.date` | Date only ("2024-11-26") |
| `.time` | Time only ("15:30") |
| `.dayOfYear` | Day number (1-366) |
| `.week` | ISO week number (1-53) |

### Methods
| Method | Description | Example |
|--------|-------------|---------|
| `.format()` | Default format | `dt.format()` → `"11/26/2024"` |
| `.format(style)` | Style format | `dt.format("long")` → `"November 26, 2024"` |
| `.format(style, locale)` | Localized | `dt.format("long","de-DE")` → `"26. November 2024"` |
| `.toDict()` | Dictionary form | `dt.toDict()` → `{year: 2024, month: 11, ...}` |

Style options: `"short"`, `"medium"`, `"long"`, `"full"`

---

## Duration Methods

### Literals
```parsley
@1d          // 1 day
@2h          // 2 hours
@30m         // 30 minutes
@1d2h30m     // Combined
@-1d         // Negative (yesterday)
```

### Methods
| Method | Description | Example |
|--------|-------------|---------|
| `.format()` | Relative time | `@1d.format()` → `"tomorrow"` |
| `.format(locale)` | Localized | `@-1d.format("de-DE")` → `"gestern"` |
| `.toDict()` | Dictionary form | `@1d2h.toDict()` → `{__type: "duration", ...}` |

### String Conversion
Durations convert to human-readable strings in templates and print statements:
```parsley
let d = @1d2h30m
"{d}"              // "1 day, 2 hours, 30 minutes"
log(d)             // 1 day, 2 hours, 30 minutes
```

### Arithmetic
```parsley
let christmas = @2025-12-25
let daysUntil = christmas - now()
daysUntil.format()  // "in 4 weeks"
```

---

## Path Methods

### Creation
```parsley
@./config.json       // Relative path
@/usr/local/bin      // Absolute path
path("some/path")    // Dynamic path
```

### Properties
| Property | Description | Example |
|----------|-------------|---------|
| `.basename` | Filename | `"config.json"` |
| `.ext` | Extension | `"json"` |
| `.stem` | Name without ext | `"config"` |
| `.dirname` | Parent directory | Path object |

### Methods
| Method | Description |
|--------|-------------|
| `.isAbsolute()` | Is absolute path |
| `.isRelative()` | Is relative path |
| `.toDict()` | Dictionary form |

### String Conversion
Paths convert to their path string in templates:
```parsley
let p = @./src/main.go
"{p}"              // "./src/main.go"
log(p)             // ./src/main.go
```

---

## URL Methods

### Creation
```parsley
@https://api.example.com/users    // URL
url("https://example.com:8080")   // Dynamic URL
```

### Properties
| Property | Description |
|----------|-------------|
| `.scheme` | Protocol ("https") |
| `.host` | Hostname |
| `.port` | Port number |
| `.path` | Path component |
| `.query` | Query parameters dict |

### Methods
| Method | Description |
|--------|-------------|
| `.origin()` | Scheme + host + port |
| `.pathname()` | Path only |
| `.search()` | Query string with `?` |
| `.href()` | Full URL string |
| `.toDict()` | Dictionary form |

```parsley
let u = @https://example.com?q=test&page=2
u.query.q      // "test"
u.query.page   // "2"
```

### String Conversion
URLs convert to their full URL string in templates:
```parsley
let u = @https://api.example.com/v1
"{u}"              // "https://api.example.com/v1"
log(u)             // https://api.example.com/v1
```

---

## File I/O

### File Handle Factories
| Factory | Format | Read Returns | Write Accepts |
|---------|--------|--------------|---------------|
| `file(path)` | Auto-detect | Depends on ext | String |
| `JSON(path)` | JSON | Dict or Array | Dict or Array |
| `CSV(path)` | CSV | Array of Dicts | Array of Dicts |
| `MD(path)` | Markdown | Dict (html + frontmatter) | String |
| `SVG(path)` | SVG | String (prolog stripped) | String |
| `lines(path)` | Lines | Array of Strings | Array of Strings |
| `text(path)` | Text | String | String |
| `bytes(path)` | Binary | Byte Array | Byte Array |

### File Handle Properties
| Property | Description |
|----------|-------------|
| `.exists` | File exists |
| `.size` | Size in bytes |
| `.modified` | Last modified datetime |
| `.isFile` | Is a file |
| `.isDir` | Is a directory |
| `.ext` | File extension |
| `.basename` | Filename |
| `.stem` | Name without extension |

### Reading (`<==`)
```parsley
let config <== JSON(@./config.json)
let rows <== CSV(@./data.csv)
let content <== text(@./readme.txt)

// Load SVG icons as reusable components
let Arrow <== SVG(@./icons/arrow.svg)
<button><Arrow/> Next</button>

// Load markdown with YAML frontmatter
let post <== MD(@./blog.md)
post.title       // From frontmatter
post.date        // Parsed as DateTime if ISO format
post.tags        // Array from frontmatter
post.html        // Rendered HTML
post.raw         // Original markdown body

// Destructure from file
let {name, version} <== JSON(@./package.json)

// Error capture pattern
let {data, error} <== JSON(@./config.json)
if (error) {
    log("Error:", error)
}

// Fallback
let config <== JSON(@./config.json) ?? {defaults: true}
```

### Writing (`==>`)
```parsley
myDict ==> JSON(@./output.json)
records ==> CSV(@./export.csv)
"Hello" ==> text(@./greeting.txt)
"<svg>...</svg>" ==> SVG(@./icon.svg)
```

### Appending (`==>>`)
```parsley
newLine ==>> lines(@./log.txt)
message ==>> text(@./debug.log)
```

### Directory Operations
```parsley
let d = dir(@./images)
d.exists      // true
d.isDir       // true
d.count       // Number of entries
d.files       // Array of file handles

// Read directory
let files <== dir(@./images)

// File patterns
let images = files(@./images/*.jpg)
let sources = files(@./src/**/*.pars)
```

---

## Regex

### Literals
```parsley
/pattern/       // Basic regex
/pattern/i      // Case insensitive
/pattern/g      // Global
```

### Dynamic Creation
```parsley
regex("\\d+", "i")
```

### Methods
| Method | Description | Example |
|--------|-------------|---------|
| `.test(string)` | Test if matches | `/\d+/.test("abc123")` → `true` |
| `.format()` | Pattern only | `/\d+/i.format()` → `\d+` |
| `.format("literal")` | Literal form | `/\d+/i.format("literal")` → `/\d+/i` |
| `.format("verbose")` | Detailed form | `/\d+/i.format("verbose")` → `regex(\d+, i)` |
| `.toDict()` | Dictionary form | `/\d+/i.toDict()` → `{pattern: "\\d+", flags: "i", ...}` |

### String Conversion
Regex patterns convert to literal notation in templates:
```parsley
let r = /[a-z]+/i
"{r}"              // "/[a-z]+/i"
log(r)             // /[a-z]+/i
```

### Matching
```parsley
"test@example.com" ~ /\w+@\w+\.\w+/  // ["test@example.com"]
"hello" ~ /\d+/                       // null (no match)
"hello" !~ /\d+/                      // true
```

### Capture Groups
```parsley
let match = "Phone: (555) 123-4567" ~ /\((\d{3})\) (\d{3})-(\d{4})/
match[0]  // Full match
match[1]  // "555"
match[2]  // "123"
match[3]  // "4567"
```

### Replace and Split
```parsley
"hello world".replace(/world/, "Parsley")  // "hello Parsley"
"a1b2c3".split(/\d+/)                      // ["a", "b", "c"]
```

---

## Modules

### Creating a Module
```parsley
// math.pars
let PI = 3.14159
let add = fn(a, b) { a + b }

// Private (no 'let')
helper = fn(x) { x * 2 }
```

### Importing
```parsley
let math = import(@./math.pars)
math.add(2, 3)  // 5

// Destructure
let {add, PI} = import(@./math.pars)
```

---

## Tags

### HTML/XML Tags
```parsley
<div class="container">
    <h1>{title}</h1>
    <p>{content}</p>
</div>
```

### Self-Closing
```parsley
<br/>
<img src="photo.jpg" />
<meta charset="utf-8" />
```

### Components
```parsley
let Card = fn({title, body}) {
    <div class="card">
        <h2>{title}</h2>
        <p>{body}</p>
    </div>
}

<Card title="Hello" body="World" />
```

### Fragments
```parsley
<>
    <p>First</p>
    <p>Second</p>
</>
```

### Raw Mode (Style/Script)
Inside `<style>` and `<script>` tags, use `@{}` for interpolation:
```parsley
let color = "blue"
<style>.class { color: @{color}; }</style>
```

### Programmatic Tags
```parsley
tag("div", {class: "box"}, "content")
// Creates tag dictionary, use toString() to render
```

---

## Utility Functions

### Type Conversion
| Function | Description |
|----------|-------------|
| `toInt(str)` | String to integer |
| `toFloat(str)` | String to float |
| `toNumber(str)` | Auto-detect int/float |
| `toString(value)` | Convert to string |

### Debugging
| Function | Description |
|----------|-------------|
| `log(...)` | Output to stdout |
| `logLine(...)` | Output with file:line prefix |
| `toDebug(value)` | Debug representation |
| `repr(value)` | Dictionary representation of pseudo-types |

### The `repr()` Function
The `repr()` function returns a detailed dictionary representation of pseudo-types (datetime, duration, regex, path, url, file, dir, request). This is useful for debugging and introspection:

```parsley
let d = @1d2h30m
repr(d)    // {__type: "duration", days: 1, hours: 2, minutes: 30, ...}

let r = /\w+/i
repr(r)    // {__type: "regex", pattern: "\\w+", flags: "i"}

let p = @./src/main.go
repr(p)    // {__type: "path", path: "./src/main.go", basename: "main.go", ...}
```

For regular values, `repr()` returns them unchanged.

### The `toDict()` Method
All pseudo-types support a `.toDict()` method that returns their internal dictionary representation:

```parsley
@2024-12-25.toDict()    // {__type: "datetime", year: 2024, month: 12, day: 25, ...}
@1h30m.toDict()         // {__type: "duration", hours: 1, minutes: 30, ...}
/\d+/g.toDict()         // {__type: "regex", pattern: "\\d+", flags: "g"}
@./config.json.toDict() // {__type: "path", path: "./config.json", ...}
```

---

## Method Chaining

Methods return appropriate types, enabling fluent chains:

```parsley
"  hello world  ".trim().upper().split(" ")  // ["HELLO", "WORLD"]
[3, 1, 2].sort().reverse()                   // [3, 2, 1]
[1, 2, 3].map(fn(x) { x * 2 }).reverse()     // [6, 4, 2]
```

## Null Propagation

Methods called on null return null instead of erroring:

```parsley
let d = {a: 1}
d.b.upper()              // null (d.b is null)
d.b.split(",").reverse() // null (entire chain)
```
