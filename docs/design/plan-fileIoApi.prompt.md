# Plan: File I/O API for Parsley

**TL;DR**: Add File Handle objects with format binding, directional operators (`<==`, `==>`, `==>>`) for read/write/append, and `??` for null coalescing. Files are first-class objects that know their format and can be inspected, passed around, and composed.

---

## Design Principles

1. **Files are objects** — Like `path`, `url`, `datetime`, files are special dictionaries
2. **Format is bound to the file** — `JSON(@./data.json)` creates a JSON-aware file handle
3. **I/O is explicit** — The `<==` operator is when reading actually happens
4. **Errors are values** — Destructure `{data, error}` or use `??` for fallbacks
5. **Composable** — File handles can be stored, passed, compared

---

## The `??` Nullish Coalescing Operator

Prerequisite operator for elegant error handling:

```parsley
value ?? default           // Returns default only if value is null

null ?? "fallback"         // "fallback"
"hello" ?? "fallback"      // "hello"
0 ?? 42                    // 0 (not null, so no fallback)
false ?? true              // false (not null)

// Chains naturally
a ?? b ?? c ?? "default"   // First non-null value
```

---

## File Handle Objects

### Creating File Handles

```parsley
// Generic file (format inferred from extension)
let f = file(@./readme.txt)

// Format-specific factories
let jf = JSON(@./config.json)
let cf = CSV(@./data.csv)
let lf = lines(@./log.txt)
let bf = bytes(@./image.png)
```

### File Handle Structure

```parsley
{
    __type: "file",
    path: {__type: "path", components: [...], absolute: true},
    format: "json",        // or "csv", "lines", "bytes", "text", null
    options: {},           // format-specific options like {header: true}
}
```

### File Handle Properties

All properties are lazy-evaluated (filesystem access on demand):

| Property | Type | Description |
|----------|------|-------------|
| `f.path` | Path | Underlying path object |
| `f.format` | String | Bound format ("json", "csv", etc.) |
| `f.options` | Dict | Format options |
| `f.exists` | Boolean | Does file exist? |
| `f.size` | Integer | Size in bytes |
| `f.modified` | Datetime | Last modification time |
| `f.created` | Datetime | Creation time |
| `f.isDir` | Boolean | Is a directory? |
| `f.isFile` | Boolean | Is a regular file? |
| `f.mode` | String | Permissions ("rwxr-xr-x") |
| `f.ext` | String | Extension (from path) |
| `f.basename` | String | Filename (from path) |

```parsley
let f = JSON(@./config.json)

f.exists      // true
f.size        // 2048
f.modified    // {__type: datetime, year: 2025, ...}
f.format      // "json"
f.path.ext    // "json"
```

---

## Format Factories

| Factory | Format | Read Returns | Write Accepts |
|---------|--------|--------------|---------------|
| `file(path)` | (inferred) | Depends on extension | String |
| `JSON(path)` | `"json"` | Dict or Array | Dict or Array |
| `CSV(path)` | `"csv"` | Array of Arrays | Array of Arrays |
| `CSV(path, {header: true})` | `"csv"` | Array of Dicts | Array of Dicts |
| `lines(path)` | `"lines"` | Array of Strings | Array of Strings |
| `text(path)` | `"text"` | String | String |
| `bytes(path)` | `"bytes"` | Byte Array | Byte Array |

### Auto-Detection from Extension

When using `file()`, format is inferred:

| Extension | Inferred Format |
|-----------|-----------------|
| `.json` | json |
| `.csv` | csv |
| `.txt`, `.md`, `.html`, `.xml` | text |
| `.log` | lines |
| (binary extensions) | bytes |
| (unknown) | text |

```parsley
// These are equivalent:
data <== JSON(@./config.json)
data <== file(@./config.json)    // Infers JSON from .json

// Explicit format overrides:
raw <== text(@./config.json)     // Read JSON file as raw text
```

---

## Read Operator: `<==`

Data flows FROM file TO variable.

### Basic Reading

```parsley
// Read with bound format
config <== JSON(@./config.json)       // Dict
records <== CSV(@./data.csv)          // Array of Arrays
logLines <== lines(@./app.log)        // Array of Strings

// Read with auto-detection
config <== file(@./config.json)       // Infers JSON
```

### With Fallback

```parsley
// Use ?? for default on error/missing
config <== JSON(@./config.json) ?? {}
users <== JSON(@./users.json) ?? []

// Cascade through multiple files
settings <== JSON(@./user.json) ?? JSON(@./default.json) ?? {}
```

### With Error Capture

```parsley
// Destructure to get both data and error
{data, error} <== JSON(@./config.json)

if (error) {
    <p class="error">Failed to load: {error}</p>
} else {
    <h1>{data.title}</h1>
}
```

### Reading from File Handle Variables

```parsley
let configFile = JSON(@./config.json)

// Check before reading
if (configFile.exists) {
    config <== configFile
} else {
    config = {defaults: true}
}
```

---

## Write Operator: `==>`

Data flows FROM variable TO file.

### Basic Writing

```parsley
// Write with format encoding
myDict ==> JSON(@./output.json)       // Encodes as JSON
records ==> CSV(@./export.csv)        // Encodes as CSV
"Hello" ==> text(@./greeting.txt)     // Writes string
```

### Error Handling

```parsley
// Returns error or null
error = data ==> JSON(@./output.json)

if (error) {
    log("Write failed:", error)
}

// Or ignore (fire and forget)
data ==> JSON(@./cache.json)
```

---

## Append Operator: `==>>`

Data flows and appends to file.

```parsley
// Append line to log
logEntry ==>> lines(@./app.log)

// Append text with newline
(message + "\n") ==>> text(@./debug.log)

// Append CSV row
newRecord ==>> CSV(@./data.csv)
```

---

## Directory Operations

### Listing Directories

```parsley
let d = dir(@./images/)

d.exists       // true
d.isDir        // true

// List contents (returns array of file handles)
{files, error} <== d

for (f in files ?? []) {
    if (f.isFile) {
        <p>{f.basename} - {f.size} bytes</p>
    }
}
```

### Glob Patterns

```parsley
// Match files by pattern
{images, error} <== glob(@./images/*.jpg)
{sources, error} <== glob(@./src/**/*.pars)

for (img in images ?? []) {
    <img src="{img.path}" />
}
```

---

## Complete Examples

### Example 1: Load Config with Defaults

```parsley
config <== JSON(@./config.json) ?? {}

<html 
    lang="{config.lang ?? "en"}"
    data-theme="{config.theme ?? "light"}"
>
    <head>
        <title>{config.title ?? "My App"}</title>
    </head>
</html>
```

### Example 2: Display Users or Error

```parsley
{users, error} <== JSON(@./users.json)

if (error) {
    <div class="error">
        <h2>Could not load users</h2>
        <pre>{error}</pre>
    </div>
} else {
    <ul class="user-list">
        for (user in users) {
            <li>
                <strong>{user.name}</strong>
                <span class="email">{user.email ?? "No email"}</span>
            </li>
        }
    </ul>
}
```

### Example 3: Process CSV Data

```parsley
{sales, error} <== CSV(@./sales.csv, {header: true})

if (error) {
    <p class="error">Failed to load sales data</p>
} else {
    let total = 0
    for (row in sales) {
        total = total + toFloat(row.amount ?? "0")
    }
    
    <table>
        <thead>
            <tr><th>Date</th><th>Product</th><th>Amount</th></tr>
        </thead>
        <tbody>
            for (row in sales) {
                <tr>
                    <td>{row.date}</td>
                    <td>{row.product}</td>
                    <td>{toFloat(row.amount).currency("USD")}</td>
                </tr>
            }
        </tbody>
        <tfoot>
            <tr><td colspan="2">Total</td><td>{total.currency("USD")}</td></tr>
        </tfoot>
    </table>
}
```

### Example 4: Image Gallery with Metadata

```parsley
{images, _} <== glob(@./gallery/*.jpg)
metadata <== JSON(@./gallery/meta.json) ?? {}

<div class="gallery">
    for (img in images ?? []) {
        let info = metadata[img.path.stem] ?? {}
        <figure>
            <img 
                src="{img.path}" 
                alt="{info.alt ?? img.path.stem}"
                loading="lazy"
            />
            <figcaption>
                <h3>{info.title ?? img.path.stem}</h3>
                <p>{info.description ?? ""}</p>
            </figcaption>
        </figure>
    }
</div>
```

### Example 5: Transform and Save Data

```parsley
{input, readErr} <== JSON(@./raw-data.json)

if (readErr) {
    log("Read failed:", readErr)
} else {
    // Transform data
    let processed = input
        .filter(fn(x) { x.active })
        .map(fn(x) { 
            {
                name: x.name.upper(),
                score: round(x.score * 1.1),
                processed: @now
            }
        })
    
    // Save result
    writeErr = processed ==> JSON(@./processed-data.json)
    
    if (writeErr) {
        log("Write failed:", writeErr)
    } else {
        <p>Processed {processed.length()} records</p>
    }
}
```

### Example 6: Append to Log File

```parsley
let entry = {
    timestamp: @now.format(),
    level: "INFO",
    message: "Page rendered successfully"
}

let line = entry.timestamp + " [" + entry.level + "] " + entry.message + "\n"
line ==>> text(@./app.log)
```

### Example 7: File Handle Composition

```parsley
// Store file handles for later use
let configs = [
    JSON(@./config/base.json),
    JSON(@./config/production.json),
    JSON(@./config/local.json)
]

// Merge configs (later files override earlier)
let merged = {}
for (configFile in configs) {
    if (configFile.exists) {
        let cfg <== configFile ?? {}
        merged = merge(merged, cfg)
    }
}
```

---

## Operator Summary

| Operator | Name | Usage | Returns |
|----------|------|-------|---------|
| `??` | Nullish Coalesce | `a ?? b` | `a` if not null, else `b` |
| `<==` | Read | `x <== file` | Data (or null on error) |
| `<==` | Read (destructure) | `{data, error} <== file` | Dict with data and error |
| `==>` | Write | `data ==> file` | Error or null |
| `==>>` | Append | `data ==>> file` | Error or null |

---

## Implementation Phases

### Phase 1: `??` Operator
- Add `NULLISH_COALESCE` token (`??`)
- Parse as infix operator (low precedence, right-associative)
- Evaluate: return left if not `*Null`, else evaluate right

### Phase 2: File Handle Objects
- Add `file()` builtin returning file dictionary
- Add `isFileDict()` helper
- Implement lazy property evaluation for metadata
- Add format factories: `JSON()`, `CSV()`, `lines()`, `text()`, `bytes()`

### Phase 3: Read Operator `<==`
- Add `READ_FROM` token (`<==`)
- Parse as special assignment: `let x <== expr` and `{pattern} <== expr`
- Implement read with format decoding

### Phase 4: Write Operators
- Add `WRITE_TO` (`==>`) and `APPEND_TO` (`==>>`) tokens
- Parse as expression statements
- Implement write/append with format encoding

### Phase 5: Directory Operations
- Add `dir()` factory
- Add `glob()` factory
- Return arrays of file handles

---

## Security Model

```bash
# Default: no file access (safe)
pars template.pars

# Allow reading from specific paths
pars --allow-read=./data,./config template.pars

# Allow writing to specific paths
pars --allow-write=./output template.pars

# Allow both
pars --allow-read=./data --allow-write=./output template.pars

# Full filesystem access (dangerous)
pars --allow-fs template.pars
```

---

## Open Questions

1. **Should `file()` accept URL objects for HTTP fetch?**
   ```parsley
   data <== JSON(url("https://api.example.com/data"))
   ```
   *Recommendation: Separate `fetch()` function for network requests*

2. **Should file handles cache their content?**
   ```parsley
   let f = JSON(@./data.json)
   a <== f
   b <== f  // Re-read or cached?
   ```
   *Recommendation: Always re-read for consistency*

3. **How to handle encoding?**
   ```parsley
   text <== text(@./latin1.txt, {encoding: "latin1"})
   ```
   *Recommendation: Default UTF-8, optional encoding parameter*

4. **Should glob return file handles or paths?**
   ```parsley
   let images <== glob(@./images/*.jpg)  // Array of file handles?
   ```
   *Recommendation: Return file handles for consistency*
