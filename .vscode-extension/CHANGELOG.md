# VSCode Extension Updates - v0.14.0

## Summary

Updated the Parsley VSCode extension to match the current grammar (v0.14.0).

## Changes in v0.14.0

### 1. Updated `package.json`
- Bumped version from `0.13.2` to `0.14.0`

### 2. Updated `syntaxes/parsley.tmLanguage.json`
- Added syntax highlighting for stdio path literals: `@-`, `@stdin`, `@stdout`, `@stderr`
- New pattern for stdio streams with `constant.other.path.stdio.parsley` scope

### Notes
- Parsley now supports reading from stdin and writing to stdout/stderr
- `@-` follows Unix convention: stdin for reads, stdout for writes
- `@stdin`, `@stdout`, `@stderr` are explicit aliases

---

# VSCode Extension Updates - v0.13.2

## Summary

Updated the Parsley VSCode extension to match the current grammar (v0.13.2).

## Changes in v0.13.2

### 1. Updated `package.json`
- Bumped version from `0.13.1` to `0.13.2`

### Notes
- `delete` is no longer a reserved keyword in Parsley
- Dictionary key deletion now uses `.delete(key)` method instead of `delete d.key` statement
- No syntax highlighting changes needed - `delete` was already listed as a builtin function

---

# VSCode Extension Updates - v0.9.13

## Summary

Updated the Parsley VSCode extension to match the current grammar (v0.9.13), including support for interpolated datetime templates.

## Changes in v0.9.13

### 1. Updated `package.json`
- Bumped version from `0.9.12` to `0.9.13`

### 2. Updated `syntaxes/parsley.tmLanguage.json`

#### Added Interpolated Datetime Template Highlighting
- **Date Templates**: `@(2024-{month}-{day})` - dates with embedded expressions
- **Time Templates**: `@({hour}:30)` - times with embedded expressions
- **Datetime Templates**: `@(2024-12-25T{hour}:30:00)` - full datetime with embedded expressions
- Interpolation expressions `{expr}` are highlighted within the template

---

# VSCode Extension Updates - v0.9.12

## Summary

Updated the Parsley VSCode extension to match the current grammar (v0.9.12), including support for interpolated path and URL templates.

## Changes in v0.9.12

### 1. Updated `package.json`
- Bumped version from `0.9.11` to `0.9.12`

### 2. Updated `syntaxes/parsley.tmLanguage.json`

#### Added Interpolated Path/URL Template Highlighting
- **Path Templates**: `@(./path/{name}/file)` - paths with embedded expressions
- **URL Templates**: `@(https://api.com/{version}/users)` - URLs with embedded expressions
- Interpolation expressions `{expr}` are highlighted within the template

---

# VSCode Extension Updates - v0.9.11

## Summary

Updated the Parsley VSCode extension to match the current grammar (v0.9.11), including support for datetime kind tracking and time-only literals.

## Changes in v0.9.11

### 1. Updated `package.json`
- Bumped version from `0.9.0` to `0.9.11`

### 2. Updated `syntaxes/parsley.tmLanguage.json`

#### Added Datetime/Duration Literal Highlighting
- **DateTime Literals**: `@2024-11-26T14:30:00Z` (full datetime with optional timezone)
- **Date Literals**: `@2024-11-26` (date only)
- **Time Literals**: `@12:30` or `@12:30:45` (time only, with or without seconds)
- **Duration Literals**: `@1d`, `@2h30m`, `@1y6mo`, `@-1d` (including negative durations)

#### Updated Built-in Functions
- **Added**: `repr` - Debug representation of pseudo-types

---

# VSCode Extension Updates - v0.9.0

## Summary

Updated the Parsley VSCode extension to match the current grammar (v0.9.0), including support for all new language features.

## Changes Made

### 1. Updated `package.json`
- Bumped version from `0.1.0` to `0.9.0`

### 2. Updated `syntaxes/parsley.tmLanguage.json`

#### Added New Literal Types
- **Regular Expression Literals**: `/pattern/flags` syntax with proper highlighting
- **Paths**: `@/path/to/file` and `@./relative/path` syntax
- **URLs**: `@https://example.com/api` syntax

#### Updated Built-in Functions
Removed deprecated functions and added new ones:
- **Added**: `import`, `now`, `time`, `path`, `url`, `regex`, `replace`, `split`
- **Removed**: Deprecated functions like `first`, `last`, `rest`, `push`, `pop`, `join`, `joinWith`, `filter`, `reduce`, `print`

Current complete list:
- `import` - Module imports
- `len`, `map`, `sort`, `sortBy`, `reverse` - Array operations
- `toString`, `toDebug`, `repr`, `toNumber`, `toInt`, `toFloat` - Type conversions
- `toUpper`, `toLower` - String operations
- `log`, `logLine` - Debugging
- `has`, `keys`, `values`, `toArray`, `toDict`, `delete` - Dictionary operations
- `sin`, `cos`, `tan`, `asin`, `acos`, `atan`, `sqrt`, `round`, `pow`, `pi` - Math
- `now`, `time` - Date/time
- `path`, `url` - Path and URL parsing
- `regex`, `replace`, `split` - Regular expressions

#### Added New Operators
- **Regex Match Operators**: `~` (match) and `!~` (not match)

#### Enhanced String Interpolation
- Added `{expr}` interpolation support in regular double-quoted strings (not just backtick templates)

### 3. Updated `README.md`
- Updated installation instructions with correct version number (0.9.0)
- Added examples for new features:
  - Module imports with `import()`
  - Regular expression literals
  - Paths
  - URLs
- Updated built-in functions list
- Added sections demonstrating new syntax

### 4. Created Test File
- Added `test-syntax.pars` to demonstrate all syntax features

## Features Now Supported

### ✅ Module System
```parsley
let utils = import(@./lib/utils.pars)
let {add, multiply} = import(@./math.pars)
let {square as sq} = import(@./helpers.pars)
```

### ✅ Regular Expressions
```parsley
let emailRegex = /^[\w.+-]+@[\w.-]+\.[a-zA-Z]{2,}$/
let match = "user@example.com" ~ emailRegex
let notMatch = "invalid" !~ emailRegex
```

### ✅ Paths
```parsley
let configPath = @./config/settings.json
let binPath = @/usr/local/bin
let homePath = @~/documents/file.txt
```

### ✅ URLs
```parsley
let apiUrl = @https://api.example.com/v1/users
let localUrl = @http://localhost:8080/api
```

### ✅ String Interpolation
```parsley
let name = "Alice"
let greeting = "Hello, {name}!"  // Works in double-quoted strings
let template = `Name: {name}`    // Also works in template literals
```

### ✅ All Built-in Functions
All current built-in functions are properly highlighted, including:
- Module system: `import()`
- Date/time: `now()`, `time()`
- Paths/URLs: `path()`, `url()`
- Regex: `regex()`, `replace()`, `split()`

## Testing

To test the extension:

1. Install the extension:
   ```bash
   cd .vscode-extension
   # macOS/Linux
   ln -s $(pwd) ~/.vscode/extensions/parsley-language
   ```

2. Reload VSCode

3. Open `test-syntax.pars` to see syntax highlighting for all features

## Version Alignment

| Component | Version |
|-----------|---------|
| Parsley Language | 0.9.0 |
| VSCode Extension | 0.9.0 |
| README.md | 0.9.0 |
| VERSION file | 0.9.0 |

All components are now synchronized at version 0.9.0.
