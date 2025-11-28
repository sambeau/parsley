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

### Paths

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

### URLs

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
@2024-12-25                    // DateTime
@2h30m                         // Duration
@/usr/local/bin                // Path (starts with /)
@./config.json                 // Path (starts with ./)
@~/documents                   // Path (starts with ~/)
@https://example.com           // URL (has ://)
@http://localhost:8080/api     // URL (has ://)
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

## Phase 3: Operator Overloading ✅

### Path Operators

#### Joining with `+`

```parsley
let base = @/usr/local
let bin = base + "bin"              // "/usr/local/bin"
let tool = bin + "my-tool"          // "/usr/local/bin/my-tool"
```

#### Joining with `/` (Unix-style)

```parsley
let root = @/var
let logs = root / "log"             // "/var/log"
let syslog = logs / "system.log"    // "/var/log/system.log"
```

The `/` operator is equivalent to `+` but provides familiar Unix-style path syntax.

#### Chaining Operations

```parsley
let project = @~/Documents + "myproject"
let src = project / "src"
let main = src + "main.go"          // "~/Documents/myproject/src/main.go"
```

#### Comparison Operators

```parsley
let p1 = @/usr/local/bin
let p2 = @/usr/local/bin
let p3 = @/usr/bin

p1 == p2        // true (same path)
p1 != p3        // true (different paths)
p1 == p3        // false
```

Paths are compared by their string representation.

### URL Operators

#### Path Extension with `+`

```parsley
let api = @https://api.example.com
let users = api + "users"                // "https://api.example.com/users"
let user = users + "123"                 // "https://api.example.com/users/123"
```

#### Building API URLs

```parsley
let github = @https://api.github.com
let repos = github + "repos"
let project = repos + "sambeau" + "parsley"
// Result: "https://api.github.com/repos/sambeau/parsley"
```

#### Comparison Operators

```parsley
let u1 = @https://example.com/api
let u2 = @https://example.com/api
let u3 = @https://example.com/docs

u1 == u2        // true (same URL)
u1 != u3        // true (different URLs)
```

URLs are compared by their string representation.

### Operator Precedence

Path and URL operators take precedence over string concatenation:

```parsley
// This joins path segments (path operator)
let p = @/usr + "local"         // "/usr/local"

// Not string concatenation:
// (which would give "{__type: path...}local")
```

### Working with Properties

Operators work seamlessly with computed properties:

```parsley
let config = @./config + "app.json"
config.basename     // "app.json"
config.extension    // "json"

let api = @https://example.com + "v1" + "users"
api.pathname        // "v1/users"
api.origin          // "https://example.com"
```

### Implementation Notes

**Evaluator (pkg/evaluator/evaluator.go):**
- `evalPathInfixExpression()` - Handles path == path, path != path
- `evalPathStringInfixExpression()` - Handles path + string, path / string
- `evalUrlInfixExpression()` - Handles url == url, url != url
- `evalUrlStringInfixExpression()` - Handles url + string
- Operator precedence: dict+string cases come **before** general string concatenation
- Fixed `urlDictToString()` to prevent double slashes in URL paths

**Supported Operators:**

| Type | Operators | Description |
|------|-----------|-------------|
| Path + String | `+`, `/` | Append path segments |
| Path == Path | `==`, `!=` | Compare paths by string |
| URL + String | `+` | Extend URL path |
| URL == URL | `==`, `!=` | Compare URLs by string |

### Examples

See `examples/operator_demo.pars` for comprehensive examples.

### Testing

Comprehensive tests in `operator_test.go`:
- `TestPathPlusOperator` - Path joining with +
- `TestPathSlashOperator` - Path joining with /
- `TestPathComparisonOperators` - Path == and !=
- `TestPathOperatorWithProperties` - Operators + computed properties
- `TestUrlPlusOperator` - URL path extension
- `TestUrlComparisonOperators` - URL == and !=
- `TestUrlOperatorWithProperties` - URL operators + properties
- `TestMixedOperations` - Complex operator combinations
- `TestOperatorPrecedence` - Chaining and precedence

All tests passing ✅

## Phase 4: Computed Properties ✅

### New Path Properties

Phase 4 adds convenient aliases and additional computed properties:

**Aliases (backward compatible):**
- `name` → alias for `basename`
- `suffix` → alias for `extension`
- `parts` → alias for `components`

**New Properties:**
- `suffixes` → Array of all file extensions (e.g., `["tar", "gz"]` from `file.tar.gz`)
- `isAbsolute` → Boolean, true if path is absolute
- `isRelative` → Boolean, true if path is relative (opposite of `isAbsolute`)

```parsley
let archive = @/backups/database.tar.gz
archive.name         // "database.tar.gz" (same as basename)
archive.suffix       // "gz" (same as extension)
archive.suffixes     // ["tar", "gz"]
archive.suffixes[0]  // "tar"
len(archive.suffixes) // 2

let config = @./app.config
config.isAbsolute    // false
config.isRelative    // true
```

### New URL Properties

**Aliases (backward compatible):**
- `hostname` → alias for `host`

**New Properties:**
- `protocol` → Scheme with `:` suffix (e.g., `"https:"` - matches Web API)
- `search` → Query string with `?` prefix (e.g., `"?key=value&foo=bar"`)
- `href` → Full URL as string (alias for `toString(url)`)

```parsley
let api = @https://example.com:8080/search?q=hello&lang=en
api.hostname         // "example.com" (same as host)
api.protocol         // "https:"
api.search           // "?q=hello&lang=en"
api.href             // "https://example.com:8080/search?q=hello&lang=en"

// Check security
if (api.protocol == "https:") {
    logLine("✓ Secure connection")
}
```

### Implementation

**Evaluator (pkg/evaluator/evaluator.go):**

Added to `evalPathComputedProperty()`:
- `name` → calls `evalPathComputedProperty(dict, "basename", env)`
- `suffix` → calls `evalPathComputedProperty(dict, "extension", env)`
- `suffixes` → splits basename on `.` and returns array of extensions
- `parts` → evaluates `components` expression
- `isAbsolute` → evaluates `absolute` field
- `isRelative` → returns `!absolute`

Added to `evalUrlComputedProperty()`:
- `hostname` → evaluates `host` field
- `protocol` → evaluates `scheme` and appends `:`
- `search` → builds query string from `query` dict with `?` prefix
- `href` → calls `urlDictToString(dict)`

### Testing

Comprehensive tests in `computed_properties_test.go`:
- `TestPathNameProperty` - name alias
- `TestPathSuffixProperty` - suffix alias
- `TestPathSuffixesProperty` - multiple extensions
- `TestPathPartsProperty` - parts alias
- `TestPathIsAbsoluteProperty` - absolute flag
- `TestPathIsRelativeProperty` - relative flag
- `TestUrlHostnameProperty` - hostname alias
- `TestUrlProtocolProperty` - protocol with colon
- `TestUrlSearchProperty` - query string
- `TestUrlHrefProperty` - full URL string
- `TestCombinedComputedProperties` - complex usage
- `TestBackwardCompatibility` - old properties still work

Demo script: `examples/computed_properties_demo.pars`

All tests passing ✅

## Complete Property Reference

### Path Properties

| Property | Type | Description | Example |
|----------|------|-------------|---------|
| `basename` | String | Last path component | `"file.txt"` |
| `name` | String | Alias for basename | `"file.txt"` |
| `dirname` | Path | Parent directory | `@/usr/local` |
| `parent` | Path | Alias for dirname | `@/usr/local` |
| `extension` | String | File extension | `"txt"` |
| `ext` | String | Alias for extension | `"txt"` |
| `suffix` | String | Alias for extension | `"txt"` |
| `stem` | String | Filename without extension | `"file"` |
| `suffixes` | Array | All extensions | `["tar", "gz"]` |
| `components` | Array | Path segments | `["", "usr", "local", "bin"]` |
| `parts` | Array | Alias for components | `["", "usr", "local", "bin"]` |
| `absolute` | Boolean | Is absolute path? | `true` |
| `isAbsolute` | Boolean | Alias for absolute | `true` |
| `isRelative` | Boolean | Opposite of absolute | `false` |

### URL Properties

| Property | Type | Description | Example |
|----------|------|-------------|---------|
| `scheme` | String | Protocol | `"https"` |
| `protocol` | String | Scheme with colon | `"https:"` |
| `host` | String | Hostname | `"example.com"` |
| `hostname` | String | Alias for host | `"example.com"` |
| `port` | Integer | Port number | `8080` |
| `path` | Array | Path segments | `["api", "v1"]` |
| `pathname` | String | Path as string | `"/api/v1"` |
| `query` | Dictionary | Query parameters | `{q: "hello"}` |
| `search` | String | Query string with `?` | `"?q=hello"` |
| `fragment` | String | URL fragment | `"section"` |
| `username` | String | Auth username | `"user"` |
| `password` | String | Auth password | `"pass"` |
| `origin` | String | Scheme + host + port | `"https://example.com:8080"` |
| `href` | String | Full URL | `"https://example.com/api"` |

## Future Enhancements

Potential additions:
1. Path normalization: `p.normalize()` for `..` and `.` handling
2. URL query manipulation: `u + {query: {key: "value"}}`
3. Relative path resolution: `p.resolve(base)`
4. File extension change: `p.withExtension("md")`
5. Custom comparison operators: `<`, `>` for path depth
6. Path joining with arrays: `p + ["dir1", "dir2"]`

