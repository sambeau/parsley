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
- [Process Execution](#process-execution)
- [Database](#database)
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
| Date | `@2024-11-26` | Date only |
| DateTime | `@2024-11-26T15:30:00` | Date and time |
| Time | `@12:30`, `@12:30:45` | Time only (uses current date internally) |
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
| `-` | Array subtraction | `[1,2,3] - [2]` → `[1, 3]` |
| `-` | Dictionary subtraction | `{a:1, b:2} - {b:0}` → `{a: 1}` |
| `*` | Multiplication | `4 * 3` → `12` |
| `*` | String repetition | `"ab" * 3` → `"ababab"` |
| `*` | Array repetition | `[1,2] * 3` → `[1, 2, 1, 2, 1, 2]` |
| `/` | Division | `10 / 4` → `2.5` |
| `/` | Array chunking | `[1,2,3,4] / 2` → `[[1, 2], [3, 4]]` |
| `%` | Modulo | `10 % 3` → `1` |
| `++` | Concatenation | `[1] ++ [2]` → `[1, 2]` |
| `++` | Scalar to array | `1 ++ [2,3]` → `[1, 2, 3]` |
| `++` | Array to scalar | `[1,2] ++ 3` → `[1, 2, 3]` |
| `..` | Range (inclusive) | `1..5` → `[1, 2, 3, 4, 5]` |

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
| Operator | Description | Example |
|----------|-------------|---------|
| `&&` | Boolean AND | `true && false` → `false` |
| `&&` | Array intersection | `[1,2,3] && [2,3,4]` → `[2, 3]` |
| `&&` | Dictionary intersection | `{a:1, b:2} && {b:3, c:4}` → `{b: 2}` |
| `||` | Boolean OR | `true || false` → `true` |
| `||` | Array union | `[1,2] || [2,3]` → `[1, 2, 3]` |
| `!` | NOT | `!true` → `false` |

### Set Operations

**Array Intersection** (`&&`): Returns elements present in both arrays (deduplicated).
```parsley
[1, 2, 3] && [2, 3, 4]           // [2, 3]
[1, 2, 2, 3] && [2, 3, 3, 4]     // [2, 3] (duplicates removed)
[1, 2] && [3, 4]                 // [] (no common elements)
```

**Array Union** (`||`): Merges arrays, removing duplicates.
```parsley
[1, 2] || [2, 3]                 // [1, 2, 3]
[1, 1, 2] || [2, 3, 3]           // [1, 2, 3] (duplicates removed)
[1, 2] || []                     // [1, 2]
```

**Array Subtraction** (`-`): Removes elements from left array that exist in right.
```parsley
[1, 2, 3, 4] - [2, 4]            // [1, 3]
[1, 2, 2, 3] - [2]               // [1, 3] (all instances removed)
[1, 2, 3] - [4, 5]               // [1, 2, 3] (no change)
```

**Dictionary Intersection** (`&&`): Returns dictionary with keys present in both (left values kept).
```parsley
{a: 1, b: 2, c: 3} && {b: 99, c: 99, d: 4}  // {b: 2, c: 3}
{a: 1} && {b: 2}                             // {}
```

**Dictionary Subtraction** (`-`): Removes keys from left that exist in right (values in right don't matter).
```parsley
{a: 1, b: 2, c: 3} - {b: 0, d: 0}  // {a: 1, c: 3}
{a: 1, b: 2} - {c: 3}              // {a: 1, b: 2} (no change)
```

**Array Chunking** (`/`): Splits array into chunks of specified size.
```parsley
[1, 2, 3, 4, 5, 6] / 2    // [[1, 2], [3, 4], [5, 6]]
[1, 2, 3, 4, 5] / 2       // [[1, 2], [3, 4], [5]]
[1, 2] / 5                // [[1, 2]]
[1, 2, 3] / 0             // ERROR: chunk size must be positive
```

**String Repetition** (`*`): Repeats string N times.
```parsley
"abc" * 3                 // "abcabcabc"
"x" * 5                   // "xxxxx"
"test" * 0                // ""
"hi" * -1                 // "" (negative treated as 0)
```

**Array Repetition** (`*`): Repeats array contents N times.
```parsley
[1, 2] * 3                // [1, 2, 1, 2, 1, 2]
["a"] * 4                 // ["a", "a", "a", "a"]
[1, 2, 3] * 0             // []
```

**Scalar Concatenation** (`++`): Wraps scalars in arrays for concatenation.
```parsley
1 ++ [2, 3, 4]            // [1, 2, 3, 4]
[1, 2, 3] ++ 4            // [1, 2, 3, 4]
1 ++ 2 ++ 3               // [1, 2, 3]
"a" ++ ["b", "c"]         // ["a", "b", "c"]
```

### Range Operator

**Range** (`..`): Creates inclusive integer ranges from start to end.
```parsley
1..5                      // [1, 2, 3, 4, 5]
0..10                     // [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
5..1                      // [5, 4, 3, 2, 1] (reverse)
-2..2                     // [-2, -1, 0, 1, 2]
10..10                    // [10] (single element)
```

**Common Use Cases:**
```parsley
// Loop over a range
for (i in 1..10) { log(i) }

// Generate sequences
let evens = (1..10).filter(fn(x) { x % 2 == 0 })  // [2, 4, 6, 8, 10]
let squares = (1..5).map(fn(x) { x * x })         // [1, 4, 9, 16, 25]

// Array indexing
let first10 = data[0..9]
let countdown = (10..1).join(", ")  // "10, 9, 8, 7, 6, 5, 4, 3, 2, 1"

// With variables
let start = 5
let end = 15
let range = start..end
```

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

### Process Execution
| Operator | Description | Example |
|----------|-------------|---------|
| `<=#=>` | Execute command with input | `let result = COMMAND("ls") <=#=> null` |

### Database
| Operator | Description | Example |
|----------|-------------|---------|
| `<=?=>` | Query single row | `let user = db <=?=> "SELECT * FROM users WHERE id = 1"` |
| `<=??=>` | Query multiple rows | `let users = db <=??=> "SELECT * FROM users"` |
| `<=!=>` | Execute mutation | `let result = db <=!=> "INSERT INTO users (name) VALUES ('Alice')"` |

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
1 ++ [2, 3]       // [1, 2, 3] (scalar concatenation)
[1, 2] ++ 3       // [1, 2, 3]
```

### Set Operations
```parsley
[1, 2, 3] && [2, 3, 4]  // [2, 3] (intersection)
[1, 2] || [2, 3]        // [1, 2, 3] (union)
[1, 2, 3] - [2]         // [1, 3] (subtraction)
```

### Other Operations
```parsley
[1, 2, 3, 4] / 2  // [[1, 2], [3, 4]] (chunking)
[1, 2] * 3        // [1, 2, 1, 2, 1, 2] (repetition)
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

### Set Operations
```parsley
{a: 1, b: 2} && {b: 3, c: 4}  // {b: 2} (intersection, left values kept)
{a: 1, b: 2} - {b: 0}         // {a: 1} (subtract keys)
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
now()                                    // Current datetime
time("2024-11-26")                       // Parse ISO date
time("2024-11-26T15:30:00")              // With time
time(1732579200)                         // Unix timestamp
time({year: 2024, month: 12, day: 25})   // From components
```

### Literals
Parsley supports three kinds of datetime literals, each with its own display format:

```parsley
@2024-11-26           // Date only
@2024-11-26T15:30:00  // Full datetime
@12:30                // Time only (HH:MM)
@12:30:45             // Time only with seconds (HH:MM:SS)
```

### Literal Kinds
Each datetime literal tracks its kind, which determines how it displays when converted to a string:

| Literal | Kind | String Output |
|---------|------|---------------|
| `@2024-11-26` | `"date"` | `"2024-11-26"` |
| `@2024-11-26T15:30:00` | `"datetime"` | `"2024-11-26T15:30:00Z"` |
| `@12:30` | `"time"` | `"12:30"` |
| `@12:30:45` | `"time_seconds"` | `"12:30:45"` |

```parsley
// Access the kind
@2024-11-26.kind           // "date"
@2024-11-26T15:30:00.kind  // "datetime"
@12:30.kind                // "time"
@12:30:45.kind             // "time_seconds"

// String conversion respects kind
toString(@2024-11-26)           // "2024-11-26"
toString(@2024-11-26T15:30:00)  // "2024-11-26T15:30:00Z"
toString(@12:30)                // "12:30"
toString(@12:30:45)             // "12:30:45"
```

### Time-Only Literals
Time-only literals (`@HH:MM` or `@HH:MM:SS`) use the current UTC date internally but display as time only:

```parsley
let meeting = @14:30
meeting.hour     // 14
meeting.minute   // 30
meeting.kind     // "time"

// Internal date is today (UTC)
meeting.year     // Current year
meeting.month    // Current month
meeting.day      // Current day

// But string output shows time only
toString(meeting)  // "14:30"
```

### Kind Preservation
The kind is preserved through arithmetic operations:

```parsley
// Date arithmetic stays date
(@2024-12-25 + 86400).kind        // "date"
(@2024-12-25 + @1d).kind          // "date"

// Datetime arithmetic stays datetime
(@2024-12-25T14:30:00 + 3600).kind  // "datetime"

// Time arithmetic stays time
(@12:30 + 3600).kind              // "time"
(@12:30:45 + 60).kind             // "time_seconds"
```

### Interpolated Datetime Templates
Use `@(...)` syntax for datetime literals with embedded expressions:

```parsley
// Date interpolation
month = "06"
day = "15"
dt = @(2024-{month}-{day})
dt.year    // 2024
dt.month   // 6
dt.day     // 15
dt.kind    // "date"

// Full datetime interpolation
year = "2025"
hour = "14"
dt2 = @({year}-12-25T{hour}:30:00)
dt2.year   // 2025
dt2.hour   // 14
dt2.kind   // "datetime"

// Time-only interpolation
h = "09"
m = "15"
meeting = @({h}:{m})
meeting.hour    // 9
meeting.minute  // 15
meeting.kind    // "time"

// Expressions in interpolations
baseDay = 10
dt3 = @(2024-12-{baseDay + 5})
dt3.day    // 15

// Dictionary-based construction
date = { year: "2024", month: "07", day: "04" }
dt4 = @({date.year}-{date.month}-{date.day})
dt4.month  // 7
```

The kind is automatically determined:
- Date templates (YYYY-MM-DD) → `"date"`
- Full datetime templates → `"datetime"`
- Time templates (HH:MM) → `"time"`

Static datetime literals (`@2024-12-25`) remain unchanged and don't require parentheses.

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
| `.kind` | Literal kind ("date", "datetime", "time", "time_seconds") |
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
| `.toDict()` | Dictionary form | `dt.toDict()` → `{year: 2024, month: 11, kind: "datetime", ...}` |

Style options: `"short"`, `"medium"`, `"long"`, `"full"`

### Comparisons
All datetime kinds can be compared:

```parsley
@12:30 < @14:00           // true
@2024-12-25 > @2024-12-24 // true
@12:30:45 == @12:30:45    // true
```

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

### Path Cleaning
Paths are automatically cleaned when created, following [Rob Pike's cleanname algorithm](https://9p.io/sys/doc/lexnames.html):
- `.` (current directory) elements are eliminated
- `..` elements eliminate the preceding component
- `..` at the start of absolute paths is eliminated (`/../foo` → `/foo`)
- `..` at the start of relative paths is preserved (`../foo` stays as is)

```parsley
let p = @/foo/../bar
p.string  // "/bar"

let p = @./a/b/../../c
p.string  // "./c"
```

### Interpolated Path Templates
Use `@(...)` syntax for paths with embedded expressions:
```parsley
name = "config"
p = @(./data/{name}.json)
p.string  // "./data/config.json"

dir = "src"
file = "main"
p = @(./{dir}/{file}.go)
p.string  // "./src/main.go"

// Expressions in interpolations
n = 1
p = @(./file{n + 1}.txt)
p.string  // "./file2.txt"
```

Static path literals (`@./path`) remain unchanged and don't require parentheses.

### Properties
| Property | Description | Example |
|----------|-------------|---------|
| `.basename` | Filename | `"config.json"` |
| `.ext` | Extension | `"json"` |
| `.stem` | Name without ext | `"config"` |
| `.dirname` | Parent directory | Path object |
| `.dir` | Parent directory as string | `"./data"` |
| `.string` | Full path as string | `"./data/config.json"` |

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

### Interpolated URL Templates
Use `@(...)` syntax for URLs with embedded expressions:
```parsley
version = "v2"
u = @(https://api.example.com/{version}/users)
u.string  // "https://api.example.com/v2/users"

host = "api.test.com"
u = @(https://{host}/data)
u.string  // "https://api.test.com/data"

// Port interpolation
port = 8080
u = @(http://localhost:{port}/api)
u.port    // "8080"

// Fragment interpolation
section = "intro"
u = @(https://docs.com/guide#{section})
u.fragment  // "intro"
```

Static URL literals (`@https://...`) remain unchanged and don't require parentheses.

### Properties
| Property | Description |
|----------|-------------|
| `.scheme` | Protocol ("https") |
| `.host` | Hostname |
| `.port` | Port number |
| `.path` | Path component |
| `.query` | Query parameters dict |
| `.string` | Full URL as string |

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

### File Handle Methods
| Method | Description |
|--------|-------------|
| `.remove()` | Removes/deletes the file from the filesystem. Returns `null` on success, error on failure. |

```parsley
// Remove a file
let f = file(@./temp.txt)
f.remove()  // Deletes the file

// With error handling
let result = f.remove()
if (result != null) {
    log("Error:", result)
}
```

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

## HTTP Requests

Fetch content from URLs using the `<=/=` operator with request handles.

### Fetch Operator

| Operator | Description | Example |
|----------|-------------|---------|
| `<=/=` | Fetch from URL | `let data <=/= JSON(@https://api.example.com)` |

### Request Handle Factories

| Factory | Format | Returns |
|---------|--------|---------|
| `JSON(url)` | JSON | Parsed JSON (dict/array) |
| `text(url)` | Plain text | String |
| `YAML(url)` | YAML | Parsed YAML |
| `lines(url)` | Lines | Array of strings |
| `bytes(url)` | Binary | Array of integers |

### Basic Usage

```parsley
// Fetch JSON data
let users <=/= JSON(@https://api.example.com/users)
log(users[0].name)

// Fetch text content
let html <=/= text(@https://example.com)

// Direct URL fetch (defaults to text)
let content <=/= @https://example.com
```

### Request Options

Pass a second argument to customize the request:

```parsley
// POST with JSON body
let response <=/= JSON(@https://api.example.com/users, {
    method: "POST",
    body: {name: "Alice", email: "alice@example.com"},
    headers: {"Authorization": "Bearer token123"}
})

// Custom timeout (milliseconds)
let data <=/= JSON(@https://slow-api.com/data, {
    timeout: 10000  // 10 seconds
})

// PUT request
let updated <=/= JSON(@https://api.example.com/users/1, {
    method: "PUT",
    body: {name: "Bob"},
    headers: {"Content-Type": "application/json"}
})
```

### Error Handling

Use destructuring to capture errors and response metadata:

```parsley
// Basic error capture
let {data, error} <=/= JSON(@https://api.example.com/data)
if (error != null) {
    log("Fetch failed:", error)
} else {
    log("Success:", data)
}

// Access HTTP status and headers
let {data, error, status, headers} <=/= JSON(@https://api.example.com/users)
log("Status code:", status)
log("Content-Type:", headers["Content-Type"])

// Handle errors gracefully
let {data, error} <=/= JSON(@https://unreliable-api.com/data)
let users = data ?? []  // Default to empty array on error
```

### HTTP Methods

Supported methods: GET (default), POST, PUT, PATCH, DELETE, HEAD, OPTIONS

```parsley
// GET (default)
let data <=/= JSON(@https://api.example.com/items)

// POST
let created <=/= JSON(@https://api.example.com/items, {
    method: "POST",
    body: {title: "New Item"}
})

// DELETE
let {data, status} <=/= JSON(@https://api.example.com/items/123, {
    method: "DELETE"
})

// PATCH
let updated <=/= JSON(@https://api.example.com/items/123, {
    method: "PATCH",
    body: {title: "Updated Title"}
})
```

### Request Headers

Customize headers for authentication, content negotiation, etc.

**Note**: Parsley dictionary syntax requires identifier keys (no hyphens or special characters). For HTTP headers with hyphens like "Content-Type" or "User-Agent", you may need to work around this limitation or use simple header names.

```parsley
// Simple headers without hyphens work fine
let data <=/= JSON(@https://api.example.com/data, {
    headers: {
        Authorization: "Bearer " + apiToken
    }
})

// For headers requiring hyphens, consider alternative approaches
// or wait for future Parsley enhancements
```

### Response Structure

When using error capture pattern `{data, error, status, headers}`:

| Field | Type | Description |
|-------|------|-------------|
| `data` | Varies | Parsed response body (based on format) |
| `error` | String/Null | Error message if request failed, `null` on success |
| `status` | Integer | HTTP status code (200, 404, 500, etc.) |
| `headers` | Dictionary | Response HTTP headers |

```parsley
let {data, error, status, headers} <=/= JSON(@https://api.example.com/data)

if (status == 200) {
    log("Success!")
} else if (status == 404) {
    log("Not found")
} else if (status >= 500) {
    log("Server error")
}
```

### Practical Examples

**API Integration:**
```parsley
// Fetch and process API data
let {data, error} <=/= JSON(@https://api.github.com/users/octocat)
if (error == null) {
    log("User: " + data.login)
    log("Repos: " + data.public_repos)
}
```

**Form Submission:**
```parsley
let formData = {
    username: "alice",
    password: "secret123"
}

let {data, error, status} <=/= JSON(@https://example.com/login, {
    method: "POST",
    body: formData,
    headers: {"Content-Type": "application/json"}
})

if (status == 200) {
    log("Login successful!")
} else {
    log("Login failed:", error)
}
```

**Download Text Content:**
```parsley
let {data, error} <=/= text(@https://raw.githubusercontent.com/user/repo/main/README.md)
if (error == null) {
    data ==> text(@./downloaded_readme.md)
}
```

**Multiple API Calls:**
```parsley
let users <=/= JSON(@https://api.example.com/users)
let posts <=/= JSON(@https://api.example.com/posts)

for (user in users) {
    let userPosts = posts.filter(fn(p) { p.userId == user.id })
    log(user.name + " has " + userPosts.length() + " posts")
}
```

### Best Practices

1. **Always handle errors** - Use `{data, error}` pattern for robust code
2. **Set reasonable timeouts** - Default is 30 seconds, adjust as needed
3. **Check status codes** - Don't assume 200 OK, verify response status
4. **Use appropriate formats** - JSON for APIs, text for HTML, bytes for binary
5. **Secure credentials** - Never hardcode API keys, use environment variables

```parsley
// Good: Error handling and timeout
let {data, error, status} <=/= JSON(@https://api.example.com/data, {
    timeout: 5000,
    headers: {"Authorization": "Bearer " + getToken()}
})

if (error != null) {
    log("Request failed:", error)
} else if (status >= 400) {
    log("HTTP error:", status)
} else {
    // Process data
    log("Success:", data)
}
```

---

## Database

Parsley provides first-class support for SQLite databases with clean, expressive operators.

### Database Operators

| Operator | Description | Returns |
|----------|-------------|---------|
| `<=?=>` | Query single row | Dictionary or `null` |
| `<=??=>` | Query multiple rows | Array of dictionaries |
| `<=!=>` | Execute mutation | `{affected, lastId}` |

### Connection Factory

```parsley
// SQLite (only supported driver currently)
let db = SQLITE(":memory:")           // In-memory database
let db = SQLITE(@./data.db)           // File-based database
let db = SQLITE("/path/to/data.db")   // String path also works
```

### Querying Data

#### Single Row Query (`<=?=>`)

Returns a dictionary if a row is found, or `null` if no match:

```parsley
let user = db <=?=> "SELECT * FROM users WHERE id = 1"
// Returns: {id: 1, name: "Alice", email: "alice@example.com"} or null

// Using with conditional
if (user) {
    log("Found user: {user.name}")
} else {
    log("User not found")
}

// With nullish coalescing
let user = db <=?=> "SELECT * FROM users WHERE id = 999" ?? {name: "Guest"}
```

#### Multiple Row Query (`<=??=>`)

Returns an array of dictionaries (empty array if no matches):

```parsley
let users = db <=??=> "SELECT * FROM users WHERE age > 25"
// Returns: [{id: 1, name: "Alice", age: 30}, {id: 2, name: "Bob", age: 28}]

// Iterate over results
for (user in users) {
    log("{user.name}: {user.email}")
}

// Get count
let count = len(users)
```

### Executing Mutations (`<=!=>`)

Execute INSERT, UPDATE, DELETE, or DDL statements:

```parsley
// CREATE TABLE
let _ = db <=!=> "CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE,
    age INTEGER
)"

// INSERT
let result = db <=!=> "INSERT INTO users (name, email, age) VALUES ('Alice', 'alice@example.com', 30)"
// Returns: {affected: 1, lastId: 1}

log("Inserted {result.affected} row(s), last ID: {result.lastId}")

// UPDATE
let result = db <=!=> "UPDATE users SET age = 31 WHERE id = 1"
// Returns: {affected: 1, lastId: 1}

// DELETE
let result = db <=!=> "DELETE FROM users WHERE id = 5"
// Returns: {affected: 1, lastId: 5}
```

### Transactions

```parsley
// Begin transaction
db.begin()

// Execute multiple statements
let _ = db <=!=> "INSERT INTO users (name) VALUES ('Alice')"
let _ = db <=!=> "INSERT INTO posts (user_id, title) VALUES (1, 'First Post')"

// Commit or rollback
if (someCondition) {
    db.commit()     // Returns true on success
} else {
    db.rollback()   // Returns true on success
}

// Check transaction status
if (db.inTransaction) {
    log("Still in transaction")
}
```

### Connection Methods

| Method | Returns | Description |
|--------|---------|-------------|
| `db.ping()` | Boolean | Test if connection is alive |
| `db.begin()` | Boolean | Start transaction |
| `db.commit()` | Boolean | Commit transaction |
| `db.rollback()` | Boolean | Rollback transaction |
| `db.close()` | Null | Close connection |

### Connection Properties

```parsley
db.type           // "sqlite"
db.connected      // true/false
db.inTransaction  // true/false
db.lastError      // Error message string or empty
```

### Data Type Mapping

SQLite types are automatically converted to Parsley types:

| SQLite Type | Parsley Type | Example |
|-------------|--------------|---------|
| INTEGER | Integer | `42` |
| REAL | Float | `3.14` |
| TEXT | String | `"hello"` |
| BLOB | String | (converted to string) |
| NULL | Null | `null` |

### Working with NULL Values

NULL database values are represented as `null` in Parsley:

```parsley
let user = db <=?=> "SELECT name, age FROM users WHERE id = 1"
// If age is NULL in database: {name: "Alice", age: null}

if (user.age == null) {
    log("Age not set")
}

// Use nullish coalescing for defaults
let age = user.age ?? 0
```

### Error Handling

Database errors are returned as Parsley errors:

```parsley
// Syntax error
let result = db <=!=> "INVALID SQL"
// Returns: ERROR: near "INVALID": syntax error

// Table doesn't exist
let users = db <=??=> "SELECT * FROM nonexistent"
// Returns: ERROR: no such table: nonexistent

// Check last error
if (db.lastError != "") {
    log("Database error: {db.lastError}")
}
```

### Complete Example

```parsley
// Create database
let db = SQLITE(@./app.db)

// Set up schema
let _ = db <=!=> "CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    created_at INTEGER DEFAULT (strftime('%s', 'now'))
)"

// Insert data
let result = db <=!=> "INSERT INTO users (name, email) VALUES ('Alice', 'alice@example.com')"
log("Created user with ID: {result.lastId}")

// Query single user
let user = db <=?=> "SELECT * FROM users WHERE email = 'alice@example.com'"

if (user) {
    log("Welcome back, {user.name}!")
    
    // Update
    let _ = db <=!=> "UPDATE users SET name = 'Alice Smith' WHERE id = {user.id}"
    
    // Query all users
    let allUsers = db <=??=> "SELECT name, email FROM users ORDER BY created_at DESC"
    
    for (u in allUsers) {
        log("{u.name} <{u.email}>")
    }
}

// Close when done
db.close()
```

### Best Practices

1. **Use single-operator syntax**: `let result = db <=!=> query` (not double-operator)
2. **Handle NULL values**: Always check for `null` in query results
3. **Use transactions for multiple operations**: Ensures data consistency
4. **Close connections**: Call `db.close()` when done (especially for file-based DBs)
5. **Check errors**: Use `db.lastError` or handle ERROR returns
6. **Avoid SQL injection**: Future versions will support parameterized queries

---

## Process Execution

Execute external commands and capture their output.

### Creating a Command

Use `COMMAND()` to create a command handle:

```parsley
// Simple command
let cmd = COMMAND("echo")

// Command with arguments
let cmd = COMMAND("ls", ["-la", "/tmp"])

// Command with options
let cmd = COMMAND("node", ["script.js"], {
    env: {NODE_ENV: "production"},
    dir: "/path/to/project",
    timeout: @30s
})
```

### Command Options

| Option | Type | Description |
|--------|------|-------------|
| `env` | Dictionary | Environment variables (merged with system env) |
| `dir` | String/Path | Working directory for command execution |
| `timeout` | Duration | Maximum execution time (process killed if exceeded) |

### Executing Commands

Use the `<=#=>` operator to execute a command:

```parsley
// Execute without input
let result = COMMAND("echo", ["hello"]) <=#=> null

// Command can also have input data (passed to stdin)
let result = COMMAND("cat") <=#=> "input data"
```

### Result Structure

Execution returns a dictionary with:

| Field | Type | Description |
|-------|------|-------------|
| `stdout` | String | Standard output from command |
| `stderr` | String | Standard error from command |
| `exitCode` | Integer | Exit code (0 for success) |
| `error` | String/Null | Error message if execution failed, `null` otherwise |

### Examples

```parsley
// Basic command
let result = COMMAND("date") <=#=> null
log("Current date:", result.stdout)
log("Exit code:", result.exitCode)

// Command with arguments
let result = COMMAND("ls", ["-la", "/tmp"]) <=#=> null
if (result.exitCode == 0) {
    log("Files:")
    log(result.stdout)
}

// Command with custom environment
let cmd = COMMAND("printenv", ["MY_VAR"], {
    env: {MY_VAR: "custom value"}
})
let result = cmd <=#=> null
log("Environment variable:", result.stdout)

// Command with working directory
let result = COMMAND("pwd", [], {dir: "/tmp"}) <=#=> null
log("Current directory:", result.stdout)

// Command with timeout
let result = COMMAND("sleep", ["60"], {timeout: @5s}) <=#=> null
if (result.error != null) {
    log("Command timed out or failed:", result.error)
}
```

### Security

Process execution requires explicit permission via command-line flags:

```bash
# Allow all process execution
./pars --allow-execute-all script.pars
./pars -x script.pars

# Allow execution from specific directories
./pars --allow-execute=/usr/bin,/bin script.pars
```

Without these flags, `COMMAND()` will return a security error.

### Error Handling

```parsley
// Command doesn't exist
let result = COMMAND("nonexistent_cmd") <=#=> null
if (result.error != null) {
    log("Error:", result.error)
}

// Non-zero exit code
let result = COMMAND("ls", ["/nonexistent"]) <=#=> null
if (result.exitCode != 0) {
    log("Command failed with code:", result.exitCode)
    log("Error output:", result.stderr)
}
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

### Format Conversion Functions

#### JSON Functions

**`parseJSON(string)`**
Parse a JSON string into Parsley objects:

```parsley
let jsonStr = "{\"name\":\"Alice\",\"age\":30}"
let obj = parseJSON(jsonStr)
log(obj.name)  // Alice
log(obj.age)   // 30

// Arrays
let arr = parseJSON("[1, 2, 3]")
log(arr[0])    // 1

// Nested structures
let data = parseJSON("{\"users\":[{\"id\":1,\"name\":\"Bob\"}]}")
log(data.users[0].name)  // Bob
```

**`stringifyJSON(object)`**
Convert Parsley objects to JSON string:

```parsley
let obj = {name: "Alice", age: 30, active: true}
let json = stringifyJSON(obj)
log(json)  // {"active":true,"age":30,"name":"Alice"}

// Arrays
let arr = [1, 2, 3]
log(stringifyJSON(arr))  // [1,2,3]

// Nested objects
let data = {user: {id: 1, name: "Bob"}, tags: ["a", "b"]}
log(stringifyJSON(data))
```

Supported types: dictionaries, arrays, strings, integers, floats, booleans, null.

#### CSV Functions

**`parseCSV(string, options?)`**
Parse CSV string into array of arrays or dictionaries:

```parsley
// Basic parsing (array of arrays)
let csv = "a,b,c\n1,2,3\n4,5,6"
let rows = parseCSV(csv)
log(rows)  // [["a","b","c"], ["1","2","3"], ["4","5","6"]]

// Parse with header (array of dictionaries)
let csv = "name,age,city\nAlice,30,NYC\nBob,25,LA"
let people = parseCSV(csv, {header: true})
log(people[0].name)   // Alice
log(people[1].city)   // LA

for (person in people) {
    log("{person.name} is {person.age} years old")
}
```

**`stringifyCSV(array)`**
Convert array of arrays to CSV string:

```parsley
let data = [
    ["Name", "Age", "City"],
    ["Alice", "30", "NYC"],
    ["Bob", "25", "LA"]
]
let csv = stringifyCSV(data)
log(csv)
// Output:
// Name,Age,City
// Alice,30,NYC
// Bob,25,LA
```

#### Practical Examples

**JSON API Response Processing:**
```parsley
// Simulate fetching JSON from API
let response = parseJSON("{\"users\":[{\"id\":1,\"name\":\"Alice\"}]}")
for (user in response.users) {
    log("User #{user.id}: {user.name}")
}

// Create JSON for API request
let request = {
    method: "POST",
    data: {username: "alice", email: "alice@example.com"}
}
let jsonRequest = stringifyJSON(request)
```

**CSV Data Processing:**
```parsley
// Read CSV with header
let csvData = "product,price,quantity\nApple,1.50,100\nBanana,0.75,200"
let inventory = parseCSV(csvData, {header: true})

// Calculate total value
let total = 0
for (item in inventory) {
    let value = parseFloat(item.price) * parseInt(item.quantity)
    total = total + value
}
log("Total inventory value: ${total}")

// Export to CSV
let report = [
    ["Product", "Value"],
    ["Apple", "150.00"],
    ["Banana", "150.00"]
]
let csvOutput = stringifyCSV(report)
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

---

## Security

Parsley provides file system access control through command-line flags. By default, write and execute operations are restricted for security.

### Security Model

| Operation | Default Behavior | Override Flags |
|-----------|-----------------|----------------|
| **Read** | ✅ Allowed | `--restrict-read=PATHS`, `--no-read` |
| **Write** | ❌ Denied | `--allow-write=PATHS`, `-w` |
| **Execute** | ❌ Denied | `--allow-execute=PATHS`, `-x` |

### Command-Line Flags

#### Read Control

```bash
--restrict-read=PATHS    # Blacklist: deny reading from paths
--no-read                # Deny all file reads
```

**Examples:**

```bash
# Prevent reading sensitive directories
./pars --restrict-read=/etc,/var script.pars

# stdin-only processing (no file reads)
./pars --no-read < data.json
```

#### Write Control

```bash
--allow-write=PATHS      # Whitelist: allow writes to specific paths
--allow-write-all        # Allow unrestricted writes (old behavior)
-w                       # Shorthand for --allow-write-all
```

**Examples:**

```bash
# Allow writes only to output directory
./pars --allow-write=./output build.pars

# Allow writes to multiple directories
./pars --allow-write=./data,./cache process.pars

# Development mode: unrestricted writes
./pars -w dev-script.pars
```

#### Execute Control

```bash
--allow-execute=PATHS    # Whitelist: allow imports from specific paths
--allow-execute-all      # Allow unrestricted module imports
-x                       # Shorthand for --allow-execute-all
```

**Examples:**

```bash
# Allow importing only from lib directory
./pars --allow-execute=./lib app.pars

# Allow imports from multiple directories
./pars --allow-execute=./lib,./modules app.pars

# Development mode: unrestricted imports
./pars -x dev-script.pars
```

### Path Resolution

All paths in security flags are:
- Resolved to absolute paths at startup
- Cleaned using filepath.Clean
- Applied to the directory and all subdirectories
- Support `~` for home directory expansion

```bash
# These are equivalent
./pars --allow-write=./output script.pars
./pars --allow-write=$(pwd)/output script.pars

# Home directory expansion
./pars --allow-write=~/Documents/output script.pars
```

### Combined Flags

Mix and match security flags for precise control:

```bash
# Static site generator: read freely, write to public
./pars --allow-write=./public build.pars

# API processor: restrict sensitive reads, write results, import libs
./pars --restrict-read=/etc --allow-write=./output --allow-execute=./lib process.pars

# Development: unrestricted writes and imports
./pars -w -x dev-script.pars

# Paranoid: specific write path, no reads, no imports
./pars --no-read --allow-write=./output template.pars
```

### Security Errors

When access is denied, clear error messages indicate the issue:

```
Error: security: file write not allowed: ./output/result.json (use --allow-write or -w)
Error: security: file read restricted: /etc/passwd
Error: security: script execution not allowed: ../tools/module.pars (use --allow-execute or -x)
```

### Migration from v0.9.x

**Breaking Changes in v0.10.0:**

- **Write operations** now denied by default
- **Module imports** (execute) now denied by default
- **Read operations** remain unrestricted (no change)

**Quick Fix:**

```bash
# Old (v0.9.x) - everything allowed
./pars build.pars

# New (v0.10.0) - add -w for old behavior
./pars -w build.pars

# Or specify allowed paths
./pars --allow-write=./output build.pars
```

### Protected Operations

The following operations are subject to security checks:

| Operation | Security Check | Example |
|-----------|----------------|---------|
| File read | `read` | `content <== text("file.txt")` |
| File write | `write` | `"data" ==> text("file.txt")` |
| File delete | `write` | `file("temp.txt").remove()` |
| Directory list | `read` | `dir("./folder").files` |
| Module import | `execute` | `import("./module.pars")` |

### Best Practices

1. **Production**: Use specific allow-lists
   ```bash
   ./pars --allow-write=./output --allow-execute=./lib app.pars
   ```

2. **Development**: Use shorthands for convenience
   ```bash
   ./pars -w -x dev-script.pars
   ```

3. **CI/CD**: Minimal permissions
   ```bash
   ./pars --allow-write=./dist build.pars
   ```

4. **Untrusted scripts**: Maximum restrictions
   ```bash
   ./pars --no-read --allow-write=./sandbox untrusted.pars
   ```
