# Path and URL Support in Parsley

## Overview

Parsley now has comprehensive support for file paths and URLs through both constructor functions and literal syntax.

## Phase 1: Constructor Functions ✅

### Path Constructor

```parsley
let p = path("/usr/local/bin/parsley")
```

Creates a path dictionary with:
- `__type`: "path"
- `components`: Array of path parts (["", "usr", "local", "bin", "parsley"])
- `absolute`: Boolean indicating if path is absolute

**Computed Properties:**
- `basename` - Last component ("parsley")
- `dirname` - Parent path ("/usr/local/bin")
- `extension` - File extension ("txt" from "file.txt")
- `stem` - Filename without extension ("file" from "file.txt")

### URL Constructor

```parsley
let u = url("https://example.com:8080/api?q=hello")
```

Creates a URL dictionary with:
- `__type`: "url"
- `scheme`: Protocol ("https")
- `host`: Hostname ("example.com")
- `port`: Port number (8080, or 0 if default)
- `path`: Array of path segments (["api"])
- `query`: Dictionary of query parameters ({"q": "hello"})
- `fragment`: URL fragment (optional)

**Computed Properties:**
- `origin` - Scheme + host + port ("https://example.com:8080")
- `pathname` - Path as string ("/api")

## Phase 2: Literal Syntax ✅

### Path Literals

Paths can be created using the `@` prefix:

```parsley
let p1 = @/usr/local/bin           // Absolute path
let p2 = @./config.json            // Relative path (current dir)
let p3 = @~/documents/notes.txt    // Home directory path
let p4 = @../parent/file.txt       // Parent directory path
```

**Property Access:**
```parsley
let p = @/usr/local/bin/file.txt
p.basename    // "file.txt"
p.dirname     // "/usr/local/bin"
p.extension   // "txt"
p.stem        // "file"
```

### URL Literals

URLs can be created using the `@` prefix:

```parsley
let u1 = @https://example.com/api
let u2 = @http://localhost:8080/test?q=hello&lang=en
let u3 = @ftp://ftp.example.org/files
```

**Property Access:**
```parsley
let u = @https://example.com:8080/api/v1?q=test
u.scheme      // "https"
u.host        // "example.com"
u.port        // 8080
u.pathname    // "/api/v1"
u.origin      // "https://example.com:8080"
u.query.q     // "test"
```

## @ Prefix Disambiguation

The lexer intelligently distinguishes between different `@` literal types:

```parsley
@2024-12-25                    // DateTime literal
@2h30m                         // Duration literal
@/usr/local/bin                // Path literal (starts with /)
@./config.json                 // Path literal (starts with ./)
@~/documents                   // Path literal (starts with ~/)
@https://example.com           // URL literal (has ://)
@http://localhost:8080/api     // URL literal (has ://)
```

Detection logic:
1. Check for `scheme://` → URL_LITERAL
2. Check for `/`, `./`, `../`, `~/` → PATH_LITERAL
3. Check for `YYYY-MM-DD` pattern → DATETIME_LITERAL
4. Default → DURATION_LITERAL

## Implementation Details

### Lexer (pkg/lexer/lexer.go)

- Added `PATH_LITERAL` and `URL_LITERAL` token types
- Added `detectAtLiteralType()` to distinguish between @ prefix types
- Added `readPathLiteral()` to read path strings (stops at property access `.`)
- Added `readUrlLiteral()` to read URL strings (handles `.com`, `.org`, etc.)
- Added `isWhitespace()` helper
- Added `isPathChar()` helper

### AST (pkg/ast/ast.go)

- Added `PathLiteral` node type
- Added `UrlLiteral` node type

### Parser (pkg/parser/parser.go)

- Registered `parsePathLiteral()` prefix function
- Registered `parseUrlLiteral()` prefix function

### Evaluator (pkg/evaluator/evaluator.go)

- Added `evalPathLiteral()` - converts path literal to dictionary
- Added `evalUrlLiteral()` - converts URL literal to dictionary
- Reuses Phase 1 functions: `parsePathString()`, `parseUrlString()`, `pathToDict()`, `urlDictToString()`

## Examples

### Path Manipulation

```parsley
let config = @./config.json
logLine("Config file:", config.basename)           // "config.json"
logLine("Directory:", toString(config.dirname))    // "."
logLine("Extension:", config.extension)            // "json"

let bin = @/usr/local/bin/parsley
logLine("Binary:", bin.basename)                   // "parsley"
logLine("Bin dir:", toString(bin.dirname))         // "/usr/local/bin"
```

### URL Parsing

```parsley
let api = @https://api.example.com:8080/v1/users?page=2&limit=10
logLine("API origin:", api.origin)                 // "https://api.example.com:8080"
logLine("Path:", api.pathname)                     // "/v1/users"
logLine("Page:", api.query.page)                   // "2"
logLine("Limit:", api.query.limit)                 // "10"
```

### Constructor vs Literal Equivalence

```parsley
// These are equivalent:
let p1 = path("/usr/local/bin")
let p2 = @/usr/local/bin
toString(p1) == toString(p2)  // true

// These are equivalent:
let u1 = url("https://example.com/api")
let u2 = @https://example.com/api
toString(u1) == toString(u2)  // true
```

## Testing

Comprehensive tests in `literal_syntax_test.go`:
- `TestPathLiterals` - Basic path literal creation
- `TestUrlLiterals` - Basic URL literal creation
- `TestUrlLiteralsWithQuery` - URL query parameters
- `TestLiteralConstructorEquivalence` - Literal vs constructor equivalence
- `TestPathLiteralsComputedProperties` - Path computed properties
- `TestUrlLiteralsComputedProperties` - URL computed properties
- `TestAtPrefixDisambiguation` - @ prefix type detection

Demo script: `examples/literal_syntax_test.pars`

All tests passing ✅

## Future Enhancements

Potential additions:
1. Path joining: `p.join("subdir", "file.txt")`
2. URL query manipulation: `u.addQuery("key", "value")`
3. Path normalization: `p.normalize()`
4. Relative path resolution: `p.resolve(base)`
5. File extension change: `p.withExtension("md")`
6. URL fragment handling improvements
