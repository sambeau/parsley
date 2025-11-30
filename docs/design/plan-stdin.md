# Plan: STDIN Support for Parsley

**STATUS**: ✅ Implemented in v0.14.0

**TL;DR**: Add support for reading from STDIN using the `@-` path literal (Unix convention), fitting seamlessly into Parsley's existing file I/O API. Optionally support `@stdin`, `@stdout`, `@stderr` as readable aliases.

---

## Design Principles

1. **Fits existing API** — STDIN is just another "path" that works with all file operations
2. **Unix-familiar** — `@-` follows the standard Unix convention for STDIN/STDOUT
3. **Format-aware** — Works with `JSON(@-)`, `CSV(@-)`, `lines(@-)`, `text(@-)`, `bytes(@-)`
4. **Error handling** — Supports `{data, error} <== JSON(@-)` and `?? fallback`
5. **Composable** — Can be stored in variables, passed to functions

---

## Syntax

### Primary: `@-` (Unix Convention)

The `-` path is universally understood in Unix to mean STDIN (for input) or STDOUT (for output):

```parsley
// Read from STDIN
input <== text(@-)
data <== JSON(@-)
lines <== lines(@-)
records <== CSV(@-)

// Write to STDOUT (explicit - normally output is implicit)
result ==> text(@-)
data ==> JSON(@-)
```

### Aliases (Optional, Self-Documenting)

For clarity, named aliases could also be supported:

```parsley
input <== text(@stdin)
output ==> text(@stdout)
errors ==> text(@stderr)
```

---

## Reading from STDIN

### Basic Reading

```parsley
// Read all STDIN as text
input <== text(@-)

// Read STDIN as lines array
for line in lines(@-) {
    <p>{line}</p>
}

// Parse STDIN as JSON
data <== JSON(@-)

// Parse STDIN as CSV
records <== CSV(@-, {header: true})
```

### With Error Handling

```parsley
// Destructure for error capture
{data, error} <== JSON(@-)

if error {
    <div class="error">Invalid JSON input: {error}</div>
} else {
    <h1>{data.title}</h1>
}

// With fallback
config <== JSON(@-) ?? {defaults: true}
```

### With File Handle

```parsley
// Create a file handle for STDIN
let input = JSON(@-)

// Check if there's input (non-blocking check?)
// Note: This may not be practical for STDIN
data <== input
```

---

## Writing to STDOUT/STDERR

By default, Parsley outputs to STDOUT. Explicit operators allow control:

### Explicit STDOUT

```parsley
// Format output as JSON
data ==> JSON(@-)

// Write raw text
"Processing complete\n" ==> text(@-)
```

### STDERR for Errors/Logs

```parsley
// Log warnings to STDERR
"Warning: file not found\n" ==> text(@stderr)

// Error messages
"Error: invalid input\n" ==> text(@stderr)

// Structured logging
{level: "error", msg: "failed"} ==> JSON(@stderr)
```

---

## Use Cases

### 1. JSON Pipeline Filter

```bash
$ echo '{"name":"Alice","age":30}' | pars filter.pars
```

```parsley
// filter.pars
let user <== JSON(@-)

<div class="user">
    <h1>{user.name}</h1>
    <p>Age: {user.age}</p>
</div>
```

### 2. Line-by-Line Processing

```bash
$ cat urls.txt | pars process.pars
```

```parsley
// process.pars
for url in lines(@-) {
    let trimmed = url.trim()
    if trimmed.length() > 0 {
        <a href="{trimmed}">{trimmed}</a>
    }
}
```

### 3. CSV to HTML Table

```bash
$ cat data.csv | pars table.pars > report.html
```

```parsley
// table.pars
let records <== CSV(@-, {header: true})

<table>
    <thead>
        <tr>
            for key in records[0].keys() {
                <th>{key}</th>
            }
        </tr>
    </thead>
    <tbody>
        for row in records {
            <tr>
                for key in row.keys() {
                    <td>{row[key]}</td>
                }
            </tr>
        }
    </tbody>
</table>
```

### 4. JSON-to-JSON Transform

```bash
$ curl -s https://api.example.com/users | pars transform.pars | jq .
```

```parsley
// transform.pars
let users <== JSON(@-)

let transformed = users.map(fn(u) {
    {
        fullName: u.firstName + " " + u.lastName,
        email: u.email.toLower(),
        active: u.status == "active"
    }
})

transformed ==> JSON(@-)
```

### 5. Multi-Source Processing

```bash
$ cat config.json | pars render.pars --data users.json
```

```parsley
// render.pars
// STDIN for config, file for data
let config <== JSON(@-)
let users <== JSON(@./users.json)

<html lang="{config.lang ?? "en"}">
    <head><title>{config.title}</title></head>
    <body>
        for user in users {
            <div>{user.name}</div>
        }
    </body>
</html>
```

### 6. Error Handling with STDERR

```parsley
{data, error} <== JSON(@-)

if error {
    // Log error to STDERR, output fallback to STDOUT
    "Error parsing input: " + error + "\n" ==> text(@stderr)
    {error: true, message: error} ==> JSON(@-)
} else {
    // Process and output
    data.map(fn(x) { x * 2 }) ==> JSON(@-)
}
```

---

## Implementation

### Phase 1: Lexer Changes

Recognize `@-` as a valid path literal:

```go
// In lexer, when parsing @ literals:
// If next char is '-' followed by non-path char, treat as stdin path
case '-':
    if !isPathChar(l.peekChar()) {
        return Token{Type: PATH, Literal: "-"}
    }
```

Alternatively, handle in the path parsing logic to accept `-` as a complete path.

### Phase 2: Path Dictionary Support

When creating a path dictionary for `-`:

```go
func evalPathLiteral(pathStr string) *Dictionary {
    if pathStr == "-" {
        return &Dictionary{
            Pairs: map[string]ast.Expression{
                "__type":    &ast.StringLiteral{Value: "path"},
                "__stdin":   &ast.BooleanLiteral{Value: true},
                "path":      &ast.StringLiteral{Value: "-"},
                "name":      &ast.StringLiteral{Value: "stdin"},
                // ... other path properties as appropriate
            },
        }
    }
    // ... normal path handling
}
```

### Phase 3: File Handle Support

Format factories should recognize STDIN paths:

```go
func evalJSONFactory(args []Object, env *Environment) Object {
    pathDict := args[0].(*Dictionary)
    
    // Check if this is a stdin path
    if isStdinPath(pathDict) {
        return &Dictionary{
            Pairs: map[string]ast.Expression{
                "__type":   &ast.StringLiteral{Value: "file"},
                "__stdin":  &ast.BooleanLiteral{Value: true},
                "format":   &ast.StringLiteral{Value: "json"},
                "path":     pathDict,
            },
        }
    }
    // ... normal file handle creation
}
```

### Phase 4: Read Operator Support

In `evalReadStatement`, handle STDIN file handles:

```go
func evalReadStatement(node *ast.ReadStatement, env *Environment) Object {
    source := Eval(node.Source, env)
    
    if fileDict, ok := source.(*Dictionary); ok && isFileDict(fileDict) {
        if isStdinFile(fileDict) {
            // Read from os.Stdin
            content, err := io.ReadAll(os.Stdin)
            if err != nil {
                return makeErrorResult(err.Error())
            }
            
            // Apply format parsing
            format := getFileFormat(fileDict)
            return parseContent(content, format)
        }
        // ... normal file reading
    }
}
```

### Phase 5: Write Operator Support (Optional)

For explicit STDOUT/STDERR output:

```go
func evalWriteStatement(node *ast.WriteStatement, env *Environment) Object {
    target := Eval(node.Target, env)
    
    if fileDict, ok := target.(*Dictionary); ok && isFileDict(fileDict) {
        if isStdoutFile(fileDict) {
            // Write to os.Stdout
            content := formatContent(data, getFileFormat(fileDict))
            os.Stdout.Write(content)
            return NULL
        }
        if isStderrFile(fileDict) {
            // Write to os.Stderr
            content := formatContent(data, getFileFormat(fileDict))
            os.Stderr.Write(content)
            return NULL
        }
        // ... normal file writing
    }
}
```

### Phase 6: Named Aliases (Optional)

Support `@stdin`, `@stdout`, `@stderr` as aliases:

```go
// In lexer or path evaluation
var stdioAliases = map[string]string{
    "stdin":  "-",
    "stdout": "-",
    "stderr": "/dev/stderr",  // or special marker
}

func normalizePathLiteral(path string) string {
    if alias, ok := stdioAliases[path]; ok {
        return alias
    }
    return path
}
```

---

## Security Considerations

### Flag Control

STDIN access should respect security flags:

```bash
# STDIN reading enabled by default (safe - just reads input)
pars script.pars

# Could be disabled if needed
pars --no-stdin script.pars
```

### No TTY Detection

Consider whether to detect if STDIN is a TTY (interactive terminal) vs pipe:

```go
// Optional: warn or behave differently if STDIN is a TTY
if isatty.IsTerminal(os.Stdin.Fd()) {
    // Interactive mode - maybe prompt or timeout?
}
```

---

## Edge Cases

### 1. Empty STDIN

```parsley
// Empty input should return empty/null, not hang
data <== JSON(@-) ?? []

// For lines, empty input = empty array
for line in lines(@-) {  // Zero iterations if empty
    ...
}
```

### 2. STDIN Already Consumed

STDIN can only be read once. Subsequent reads return empty:

```parsley
first <== text(@-)   // Gets all input
second <== text(@-)  // Empty string (already consumed)
```

**Potential Enhancement**: Buffer STDIN on first read for re-use within same script.

### 3. Binary Input

```parsley
// Read binary data
data <== bytes(@-)

// Useful for base64 encoding, etc.
encoded = base64(data)
```

### 4. Large Input

For very large inputs, consider streaming (future enhancement):

```parsley
// Future: streaming API
for line in stream(@-) {
    // Process line by line without loading all into memory
}
```

---

## Comparison with Other Languages

| Language | STDIN Syntax | Notes |
|----------|--------------|-------|
| Perl 5 | `<STDIN>`, `<>` | Magic diamond reads ARGV or STDIN |
| Raku | `$*IN.lines`, `$*IN.slurp` | Dynamic variable |
| Python | `sys.stdin.read()` | File object |
| Ruby | `STDIN.read`, `$stdin` | Global constant |
| Node.js | `process.stdin` | Stream object |
| Bash | `read`, `cat -` | `-` convention |
| **Parsley** | `text(@-)`, `JSON(@-)` | Path literal, format-aware |

Parsley's approach is most similar to the Unix `-` convention but with the added benefit of format binding.

---

## Summary

| Feature | Syntax | Description |
|---------|--------|-------------|
| Read text | `text <== text(@-)` | Read STDIN as string |
| Read lines | `lines <== lines(@-)` | Read STDIN as line array |
| Read JSON | `data <== JSON(@-)` | Parse STDIN as JSON |
| Read CSV | `records <== CSV(@-)` | Parse STDIN as CSV |
| Read bytes | `data <== bytes(@-)` | Read STDIN as byte array |
| Write STDOUT | `data ==> text(@-)` | Explicit STDOUT output |
| Write STDERR | `msg ==> text(@stderr)` | Write to STDERR |
| Error handling | `{data, error} <== JSON(@-)` | Capture parse errors |
| Fallback | `data <== JSON(@-) ?? {}` | Default on error/empty |

---

## Implementation Phases

| Phase | Description | Complexity |
|-------|-------------|------------|
| 1 | Lexer: recognize `@-` as path | Low |
| 2 | Path dict: mark stdin paths | Low |
| 3 | File handles: stdin-aware | Low |
| 4 | Read operator: read from os.Stdin | Medium |
| 5 | Write operator: stdout/stderr | Medium |
| 6 | Aliases: @stdin, @stdout, @stderr | Low |

---

## Open Questions

1. **Should STDIN be buffered for re-reading?**
   - Pro: More intuitive behavior
   - Con: Memory usage for large inputs
   - Recommendation: No buffering initially, document "read once" behavior

2. **Should we detect TTY and warn/timeout?**
   - Pro: Better UX when accidentally run interactively
   - Con: Complexity, may interfere with intentional interactive use
   - Recommendation: No detection initially

3. **Should `@-` for write default to STDOUT or be an error?**
   - `data ==> text(@-)` → STDOUT
   - `data ==> JSON(@-)` → STDOUT with JSON formatting
   - Recommendation: Allow, treat as explicit STDOUT

4. **Should we support `/dev/stdin` path on Unix?**
   - Already works if filesystem allows
   - Could normalize to `@-` internally
   - Recommendation: Let filesystem handle it naturally
