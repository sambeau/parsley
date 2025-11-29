# Parsley Cheat Sheet

Quick reference for Copilot AI agent developing Parsley. Focus on key differences from JavaScript, Python, Rust, and Go.

---

## 🚨 Major Gotchas (Common Mistakes)

### 1. Output Functions
```parsley
// ❌ WRONG (JavaScript/Python style)
print("hello")
println("hello")
console.log("hello")

// ✅ CORRECT
log("hello")           // Most common - concatenates args with spaces
logLine("hello")       // Includes line number - USE FOR DEBUGGING
```

### 2. Comments
```parsley
// ✅ CORRECT - C-style comments only
// This is a comment

/* Multi-line
   comments work too */

// ❌ WRONG - No Python/Shell style
# This will ERROR
```

### 3. For Loops Return Arrays (Like map)
```parsley
// ❌ WRONG (JavaScript thinking)
for (n in [1,2,3]) {
    console.log(n)  // Expecting side effects only
}

// ✅ CORRECT - For is expression-based, returns array
let doubled = for (n in [1,2,3]) { n * 2 }  // [2, 4, 6]

// Filter pattern - if returns null, omitted from result
let evens = for (n in [1,2,3,4]) {
    if (n % 2 == 0) { n }  // [2, 4]
}
```

### 4. If is an Expression (Like Ternary)
```parsley
// ✅ CORRECT - If returns value
let status = if (age >= 18) "adult" else "minor"

// Can use in concatenation
let msg = "You are " + if (premium) "premium" else "regular"

// Block style
let result = if (x > 0) {
    "positive"
} else if (x < 0) {
    "negative" 
} else {
    "zero"
}
```

### 5. Path Literals Use @
```parsley
// ✅ CORRECT
let path = @./config.json
let url = @https://example.com
let date = @2024-11-29
let time = @14:30
let duration = @1d

// ❌ WRONG
let path = "./config.json"  // This is just a string
```

---

## 📊 Most Used Features (from actual usage data)

### Core (use these constantly)
- `log()` - 710 uses in tests, 523 in examples - **#1 most used**
- `let` - 382 uses - variable declaration
- `if` - 72 uses - conditional expression
- `for` - 44 uses - iteration/mapping

### File I/O (very common)
- `file(@path)` - 46 test uses, 8 example uses
- `JSON(@path)` - 20 uses
- `dir(@path)` - 23 test uses, 15 example uses
- `text(@path)` - 27 test uses, 7 example uses

### String/Array (frequent)
- `len()` - 189 test uses, 14 example uses
- `split()` - 23 test uses, 6 example uses
- `sort()` - 7 uses

### DateTime (common)
- `time()` - 159 test uses, 22 example uses
- `now()` - 28 test uses, 8 example uses

---

## 🔤 Syntax Comparison

### Variables & Functions

| Feature | JavaScript | Python | Parsley |
|---------|-----------|--------|---------|
| Variable | `let x = 5` | `x = 5` | `let x = 5` |
| Destructure | `const {x, y} = obj` | `x, y = obj` | `let {x, y} = obj` |
| Multiple | `let a=1, b=2` | `a, b = 1, 2` | `let a, b = 1, 2` |
| Function | `(x) => x*2` | `lambda x: x*2` | `fn(x) { x*2 }` |
| Named func | `function f(x) {}` | `def f(x):` | `let f = fn(x) {}` |

### Control Flow

| Feature | JavaScript | Python | Parsley |
|---------|-----------|--------|---------|
| If expr | `x ? "yes" : "no"` | `"yes" if x else "no"` | `if (x) "yes" else "no"` |
| If block | `if (x) { } else { }` | `if x:\nelse:` | `if (x) {} else {}` |
| For loop | `for (let x of arr)` | `for x in arr:` | `for (x in arr) {}` |
| Map | `arr.map(x => x*2)` | `[x*2 for x in arr]` | `for (x in arr) { x*2 }` |
| Filter | `arr.filter(x => x>0)` | `[x for x in arr if x>0]` | `for (x in arr) { if (x>0) {x} }` |
| Index | `arr.forEach((x,i) => )` | `for i, x in enumerate(arr):` | `for (i, x in arr) {}` |

### Data Types

| Type | JavaScript | Python | Parsley |
|------|-----------|--------|---------|
| Array | `[1, 2, 3]` | `[1, 2, 3]` | `[1, 2, 3]` |
| Dict | `{x: 1, y: 2}` | `{"x": 1, "y": 2}` | `{x: 1, y: 2}` |
| String | `` `Hi ${x}` `` | `f"Hi {x}"` | `"Hi {x}"` |
| Regex | `/abc/i` | `re.compile(r"abc", re.I)` | `/abc/i` |
| Null | `null` | `None` | `null` |

---

## 🎯 Key Language Features

### 1. Concatenative/Expression-Based
Everything is an expression that returns a value:
```parsley
// If returns value
let x = if (true) 10 else 20  // x = 10

// For returns array
let squares = for (n in 1..5) { n * n }  // [1, 4, 9, 16, 25]

// Tags return strings
let html = <p>Hello</p>  // "<p>Hello</p>"
```

### 2. String Interpolation (like template literals)
```parsley
let name = "Alice"
let msg = "Hello, {name}!"      // "Hello, Alice!"
let calc = "2 + 2 = {2 + 2}"    // "2 + 2 = 4"

// In attributes
<div class="user-{id}">Content</div>
```

### 3. HTML/XML as First-Class
```parsley
// Tags return strings
<p>Hello</p>                    // "<p>Hello</p>"

// Components are just functions
let Card = fn({title, body}) {
    <div class="card">
        <h3>{title}</h3>
        <p>{body}</p>
    </div>
}

// Use like JSX
<Card title="Welcome" body="Hello world"/>
```

### 4. Literal Syntax with @
```parsley
// Paths
@./relative/path
@~/home/path
@/absolute/path

// URLs
@https://example.com
@https://api.github.com/users

// Dates/Times
@2024-11-29                  // Date
@2024-11-29T14:30:00        // DateTime
@14:30                       // Time
@14:30:45                    // Time with seconds

// Durations
@1d                          // 1 day
@2h30m                       // 2 hours 30 minutes
@-1w                         // Negative: 1 week ago

// Interpolated (dynamic)
let month = "11"
let day = "29"
let date = @(2024-{month}-{day})  // Builds from variables
```

### 5. Operators Are Overloaded
```parsley
// Arithmetic
5 + 3                        // 8
"Hello" + " World"           // "Hello World"
@/usr + "local"              // @/usr/local (path join)
@https://api.com + "/v1"     // @https://api.com/v1

// Multiplication
3 * 4                        // 12
"ab" * 3                     // "ababab"
[1, 2] * 3                   // [1, 2, 1, 2, 1, 2]

// Division
10 / 3                       // 3.333...
[1,2,3,4,5,6] / 2           // [[1,2], [3,4], [5,6]] (chunk)

// Logical become set operations on collections
[1,2,3] && [2,3,4]          // [2, 3] (intersection)
[1,2] || [2,3]              // [1, 2, 3] (union)
[1,2,3] - [2]               // [1, 3] (subtraction)

// Range
1..5                         // [1, 2, 3, 4, 5]
```

### 6. File I/O with Special Operators
```parsley
// Read operators
let data <== JSON(@./config.json)       // Read file
let {name, error} <== JSON(@./data.json) // With error capture

// Write operators  
data ==> JSON(@./output.json)           // Write/overwrite
data ==>> text(@./log.txt)              // Append

// Network operators (HTTP/SFTP)
let {response, error} <=/= Fetch(@https://api.example.com)
data =/=> conn(@/remote/file.json).json
```

### 7. Method Chaining
```parsley
// String methods
"hello".toUpper()              // "HELLO"
"  trim  ".trim()           // "trim"
"a,b,c".split(",")          // ["a", "b", "c"]

// Array methods
[3,1,2].sort()              // [1, 2, 3]
[1,2,3].reverse()           // [3, 2, 1]
[1,2,3].join(",")           // "1,2,3"

// Path methods
@./file.txt.exists          // true/false
@./file.txt.basename        // "file.txt"
@./file.txt.ext             // "txt"

// Chaining
"  HELLO  ".trim().toLower()  // "hello"
```

---

## 📁 File I/O Patterns

### Factory Functions
```parsley
file(@path)      // Auto-detect format from extension
JSON(@path)      // Parse as JSON
CSV(@path)       // Parse as CSV
MD(@path)        // Markdown with frontmatter
text(@path)      // Plain text (use for HTML files)
lines(@path)     // Array of lines
bytes(@path)     // Byte array
SVG(@path)       // SVG (strips prolog)
dir(@path)       // Directory listing
```

### Read Patterns
```parsley
// Simple read
let config <== JSON(@./config.json)

// With error handling
let {data, error} <== JSON(@./data.json)
if (error) {
    log("Error:", error)
}

// With fallback
let config <== JSON(@./config.json) ?? {default: true}
```

### Write Patterns
```parsley
// Overwrite
data ==> JSON(@./output.json)

// Append
log_entry ==>> text(@./log.txt)
```

### Directory Operations (NEW in v0.12.1)
```parsley
// Create directories
file(@./new-dir).mkdir()
file(@./parent/child).mkdir({parents: true})  // Recursive

// Remove directories
file(@./old-dir).rmdir()
file(@./tree).rmdir({recursive: true})        // With contents

// Works with dir() too
dir(@./test).mkdir()
dir(@./test).rmdir()
```

---

## 🌐 Network Operations

### HTTP (Fetch)
```parsley
// Simple GET
let {data, error} <=/= Fetch(@https://api.example.com)

// POST with body
let payload = {name: "Alice", age: 30}
let {response, error} =/=> Fetch(@https://api.example.com/users, {
    body: payload
})
```

### SFTP (NEW in v0.12.0)
```parsley
// Connect with SSH key
let conn = SFTP(@sftp://user@host, {
    keyFile: @~/.ssh/id_rsa,
    timeout: @10s
})

// Read remote file
let {config, error} <=/= conn(@/remote/config.json).json

// Write remote file
data =/=> conn(@/remote/output.json).json

// Directory operations
conn(@/remote/new-dir).mkdir()
conn(@/remote/old-dir).rmdir({recursive: true})

// Close connection
conn.close()
```

---

## 🔧 Common Patterns

### Error Handling
```parsley
// Capture pattern
let {data, error} <== JSON(@./file.json)
if (error) {
    log("Failed:", error)
    return
}
log("Success:", data)

// Fallback pattern
let config <== JSON(@./config.json) ?? {default: "settings"}
```

### Map/Filter
```parsley
// Map
let doubled = for (n in numbers) { n * 2 }

// Filter  
let evens = for (n in numbers) {
    if (n % 2 == 0) { n }
}

// Map + Filter
let processed = for (item in items) {
    if (item.active) {
        item.name.toUpper()
    }
}
```

### Components
```parsley
// Define component
let Button = fn({text, onClick}) {
    <button onclick="{onClick}">{text}</button>
}

// Use component
<Button text="Click Me" onClick="handleClick()"/>

// With children
let Card = fn({title}, ...children) {
    <div class="card">
        <h3>{title}</h3>
        {children}
    </div>
}

<Card title="Welcome">
    <p>Body content</p>
    <p>More content</p>
</Card>
```

### Modules
```parsley
// Export from module
export({
    greet: fn(name) { "Hello, {name}!" },
    PI: 3.14159
})

// Import in another file
let utils <== import(@./utils.pars)
log(utils.greet("Alice"))
```

---

## 🎨 String Formatting

### Numbers
```parsley
1234567.format()                  // "1,234,567"
99.99.currency("USD")             // "$99.99"
99.99.currency("EUR", "de-DE")    // "99,99 €"
0.15.percent()                    // "15%"
```

### Dates
```parsley
now().format("short")             // "11/29/24"
now().format("medium")            // "Nov 29, 2024"
now().format("long")              // "November 29, 2024"
now().format("long", "de-DE")     // "29. November 2024"

@2024-11-29.format("full")        // "Friday, November 29, 2024"
```

### Durations
```parsley
@1d.format()                      // "tomorrow"
@-1d.format()                     // "yesterday"
@2h30m.format()                   // "2 hours"
```

---

## 🔍 Type Checking

```parsley
type(42)                          // "INTEGER"
type(3.14)                        // "FLOAT"
type("hi")                        // "STRING"
type([1,2])                       // "ARRAY"
type({x: 1})                      // "DICTIONARY"
type(fn() {})                     // "FUNCTION"
type(null)                        // "NULL"
type(true)                        // "BOOLEAN"
type(@2024-11-29)                 // "DATE"
type(@14:30)                      // "TIME"
type(@1d)                         // "DURATION"
```

---

## 🚀 Quick Examples

### Simple Script
```parsley
// Read, transform, write
let data <== JSON(@./input.json)
let processed = for (item in data) {
    {
        name: item.name.toUpper(),
        score: item.score * 2
    }
}
processed ==> JSON(@./output.json)
log("Processed {len(processed)} items")
```

### HTML Generation
```parsley
let users = [
    {name: "Alice", role: "Admin"},
    {name: "Bob", role: "User"}
]

let UserTable = fn(users) {
    <table>
        <tr><th>Name</th><th>Role</th></tr>
        {for (user in users) {
            <tr>
                <td>{user.name}</td>
                <td>{user.role}</td>
            </tr>
        }}
    </table>
}

<UserTable users={users}/>
```

### API Integration
```parsley
let {posts, error} <=/= Fetch(@https://jsonplaceholder.typicode.com/posts)

if (error) {
    log("API Error:", error)
} else {
    for (post in posts) {
        log("Post {post.id}: {post.title}")
    }
}
```

---

## 📝 Testing/Debugging Tips

1. **Use `logLine()` in multi-line scripts** - shows line numbers
2. **Check types with `type()`** when confused
3. **Remember `for` returns an array** - don't expect side effects
4. **Use error capture** `{data, error}` for file/network ops
5. **Path literals need @** - `@./file` not `"./file"`
6. **Comments are //** not #
7. **Output is `log()`** not `print()` or `console.log()`

---

## 🎯 When Writing Tests/Examples

### Most Common Test Pattern
```parsley
logLine("=== Test Description ===")
let input = [1, 2, 3]
let result = for (n in input) { n * 2 }
log("Result:", result)
log("Length:", len(result))
logLine()  // Blank line
```

### File Operation Pattern
```parsley
// Write test data
testData ==> JSON(@./test-file.json)

// Read it back
let {data, error} <== JSON(@./test-file.json)

// Verify
if (error) {
    log("ERROR:", error)
} else {
    log("SUCCESS:", data)
}

// Cleanup
file(@./test-file.json).remove()
```

### Component Testing Pattern
```parsley
let TestComponent = fn({title, items}) {
    <div>
        <h1>{title}</h1>
        <ul>
            {for (item in items) {
                <li>{item}</li>
            }}
        </ul>
    </div>
}

let result = <TestComponent 
    title="Test" 
    items={["a", "b", "c"]}
/>

log(result)
```
