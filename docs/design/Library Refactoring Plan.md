# Parsley Library Refactoring Plan

**Author:** Sam Phillips  
**Date:** 30 November 2025  
**Version:** Draft 1.0

## TL;DR

Refactor Parsley into a clean library (`pkg/parsley/`) with a generic Go↔Parsley API. Move CLI to `cmd/pars/`. This enables multiple consumers (CLI, server, static site builder) without HTTP or domain-specific knowledge in the core library.

---

## Goals

1. **Clean separation** - Library knows nothing about HTTP, CLI args, or other consumers
2. **Generic conversion** - Go types ↔ Parsley types, bidirectionally
3. **Backward compatible** - CLI and REPL behavior unchanged
4. **Extensible** - Easy to add server, static site builder, etc.
5. **Testable** - Each layer independently testable

---

## Current Architecture

```
main.go                     # CLI + REPL entry point (mixed concerns)
pkg/
├── ast/                    # AST definitions
├── evaluator/              # Core evaluation engine
├── formatter/              # HTML formatting
├── lexer/                  # Tokenizer  
├── parser/                 # Parser
└── repl/                   # REPL implementation
```

**Current evaluation flow in main.go:**
```go
l := lexer.NewWithFilename(string(content), filename)
p := parser.New(l)
program := p.ParseProgram()
env := evaluator.NewEnvironment()
env.Filename = filename
env.Security = policy
evaluated := evaluator.Eval(program, env)
```

---

## Proposed Architecture

```
cmd/
├── pars/
│   └── main.go             # CLI entry point
└── pars-server/            # (future) HTTP server
    └── main.go
pkg/
├── parsley/                # NEW: Public facade API
│   ├── parsley.go          # Core Eval functions
│   ├── convert.go          # Go↔Parsley conversion
│   ├── options.go          # Functional options
│   ├── logger.go           # Logger interface + implementations
│   └── README.md           # Library documentation
├── ast/                    # (unchanged)
├── evaluator/              # (minor changes for logger)
├── formatter/              # (unchanged)
├── lexer/                  # (unchanged)
├── parser/                 # (unchanged)
└── repl/                   # (uses new library API)
```

---

## Public API Design

### Core Evaluation

```go
package parsley

// Eval evaluates Parsley source code and returns the result
func Eval(source string, opts ...Option) (*Result, error)

// EvalFile evaluates a Parsley file
func EvalFile(filename string, opts ...Option) (*Result, error)

// Result wraps evaluation output
type Result struct {
    Value  Object      // The Parsley object result
    Output string      // Captured output (if using BufferedLogger)
}

// Convert result to Go value
func (r *Result) GoValue() interface{}

// Get string representation
func (r *Result) String() string
```

### Options Pattern

```go
// Option configures evaluation
type Option func(*Config)

// WithVar pre-populates a variable in the environment
func WithVar(name string, value interface{}) Option

// WithEnv uses a pre-configured environment
func WithEnv(env *Environment) Option

// WithSecurity sets the security policy
func WithSecurity(policy *SecurityPolicy) Option

// WithLogger sets the logger for log()/logLine() output
func WithLogger(logger Logger) Option

// WithFilename sets the filename for error messages
func WithFilename(name string) Option
```

### Go↔Parsley Conversion

```go
// ToParsley converts a Go value to a Parsley Object
func ToParsley(v interface{}) (Object, error)

// Supported conversions:
//   int, int64, int32, etc.  → Integer
//   float64, float32         → Float
//   string                   → String
//   bool                     → Boolean
//   []interface{}            → Array
//   []T (any slice)          → Array
//   map[string]interface{}   → Dictionary
//   map[string]T             → Dictionary
//   []byte                   → Bytes
//   time.Time                → DateTime dictionary
//   time.Duration            → Duration dictionary
//   nil                      → Null

// FromParsley converts a Parsley Object to a Go value
func FromParsley(obj Object) interface{}

// Returns:
//   Integer    → int64
//   Float      → float64
//   String     → string
//   Boolean    → bool
//   Array      → []interface{}
//   Dictionary → map[string]interface{}
//   Bytes      → []byte
//   Null       → nil
```

### Environment

```go
// NewEnvironment creates a fresh evaluation environment
func NewEnvironment() *Environment

// SetVar sets a variable (converts Go value to Parsley)
func (e *Environment) SetVar(name string, value interface{}) error

// GetVar gets a variable (converts Parsley value to Go)
func (e *Environment) GetVar(name string) (interface{}, bool)

// SetParsley sets a variable with a Parsley Object directly
func (e *Environment) SetParsley(name string, obj Object)

// GetParsley gets a variable as a Parsley Object
func (e *Environment) GetParsley(name string) (Object, bool)
```

### Logging

```go
// Logger interface for log()/logLine() output
type Logger interface {
    Log(values ...interface{})
    LogLine(values ...interface{})
}

// StdoutLogger returns a logger that writes to stdout (default)
func StdoutLogger() Logger

// WriterLogger returns a logger that writes to an io.Writer
func WriterLogger(w io.Writer) Logger

// BufferedLogger captures output for later retrieval
type BufferedLogger struct { ... }
func NewBufferedLogger() *BufferedLogger
func (l *BufferedLogger) String() string
func (l *BufferedLogger) Lines() []string

// NullLogger discards all output
func NullLogger() Logger
```

### Security (existing, re-exported)

```go
// SecurityPolicy controls file system access
type SecurityPolicy struct {
    RestrictRead    []string
    NoRead          bool
    AllowWrite      []string
    AllowWriteAll   bool
    AllowExecute    []string
    AllowExecuteAll bool
}

func NewSecurityPolicy() *SecurityPolicy
```

---

## Usage Examples

### CLI (simple)

```go
// cmd/pars/main.go
func main() {
    result, err := parsley.EvalFile(os.Args[1])
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    fmt.Print(result.String())
}
```

### CLI with arguments

```go
result, err := parsley.EvalFile(filename,
    parsley.WithVar("args", os.Args[2:]),
    parsley.WithSecurity(policy),
)
```

### HTTP Server

```go
// cmd/pars-server/main.go
func handler(w http.ResponseWriter, r *http.Request) {
    // Map HTTP request to Go types (server's responsibility)
    reqMap := map[string]interface{}{
        "method":  r.Method,
        "path":    r.URL.Path,
        "query":   queryToMap(r.URL.Query()),
        "headers": headersToMap(r.Header),
        "body":    readBody(r),
    }
    
    // Create request-scoped logger
    logger := NewRequestLogger(r.Context())
    
    // Evaluate script
    result, err := parsley.EvalFile(scriptPath,
        parsley.WithVar("request", reqMap),
        parsley.WithLogger(logger),
        parsley.WithSecurity(serverPolicy),
    )
    
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    
    // Interpret result (server's responsibility)
    writeResponse(w, result.GoValue())
}
```

### Static Site Builder

```go
for _, page := range pages {
    result, err := parsley.EvalFile(page.Template,
        parsley.WithVar("page", map[string]interface{}{
            "title":   page.Title,
            "date":    page.Date,
            "content": page.Content,
            "tags":    page.Tags,
        }),
        parsley.WithVar("site", siteConfig),
    )
    
    os.WriteFile(page.OutputPath, []byte(result.String()), 0644)
}
```

### Testing

```go
func TestMyScript(t *testing.T) {
    logger := parsley.NewBufferedLogger()
    
    result, err := parsley.Eval(`log("hello"); 1 + 2`,
        parsley.WithLogger(logger),
    )
    
    assert.NoError(t, err)
    assert.Equal(t, int64(3), result.GoValue())
    assert.Equal(t, "hello", logger.String())
}
```

---

## Refactoring Steps

### Phase 1: Create Library Facade

1. **Create `pkg/parsley/` package**
   - `parsley.go` - `Eval()`, `EvalFile()`, `Result` type
   - `convert.go` - `ToParsley()`, `FromParsley()` 
   - `options.go` - Options pattern implementation
   - `logger.go` - Logger interface and implementations

2. **Wire up to existing evaluator**
   - Facade wraps `lexer` → `parser` → `evaluator` pipeline
   - Handle parse errors and runtime errors uniformly

3. **Add tests for new API**
   - Conversion tests (Go↔Parsley round-trips)
   - Option tests
   - Logger tests

### Phase 2: Integrate Logger

4. **Add Logger to Environment**
   - Add `Logger` field to `evaluator.Environment`
   - Default to `StdoutLogger()` if not set

5. **Update `log()`/`logLine()` builtins**
   - Change from `fmt.Print()` to `env.Logger.Log()`
   - Backward compatible (default logger = stdout)

6. **Test logger integration**
   - Verify CLI/REPL unchanged
   - Test buffered logger captures output

### Phase 3: Move CLI

7. **Create `cmd/pars/main.go`**
   - Move CLI logic from root `main.go`
   - Use new library API
   - Keep all existing flags

8. **Update `pkg/repl/`**
   - Use library API for evaluation
   - Keep existing liner integration

9. **Update root `main.go`**
   - Thin wrapper that calls cmd/pars or just remove

10. **Update build scripts**
    - Makefile targets for `cmd/pars`
    - VS Code tasks

### Phase 4: Documentation

11. **Write `pkg/parsley/README.md`**
    - Installation
    - Quick start
    - Full API reference
    - Examples for each consumer type

12. **Update root README**
    - Note library availability
    - Link to library docs

---

## Internal Changes Required

### evaluator.Environment

```go
// Add to Environment struct
type Environment struct {
    // ... existing fields ...
    Logger Logger  // NEW: for log()/logLine() output
}

// Update NewEnvironment
func NewEnvironment() *Environment {
    return &Environment{
        store:       make(map[string]Object),
        letBindings: make(map[string]bool),
        exports:     make(map[string]bool),
        Logger:      defaultStdoutLogger,  // NEW
    }
}
```

### log() and logLine() builtins

```go
// Before
fmt.Print(values...)

// After  
env.Logger.Log(values...)
```

### Module cache (verify thread safety)

The existing `moduleCache` uses a mutex - verify this is sufficient for concurrent server requests.

---

## Consumer Responsibility Matrix

| Concern | Library | CLI | Server | Static Builder |
|---------|---------|-----|--------|----------------|
| Parse/evaluate scripts | ✓ | | | |
| Go↔Parsley conversion | ✓ | | | |
| Logger interface | ✓ | | | |
| Security policy | ✓ | | | |
| Command-line args | | ✓ | | |
| HTTP request parsing | | | ✓ | |
| Multipart handling | | | ✓ | |
| Response writing | | | ✓ | |
| Request logging | | | ✓ | |
| Frontmatter parsing | | | | ✓ |
| File output | | | | ✓ |
| DB connection pools | | | ✓ | |

---

## Open Questions

1. **Connection pooling** - Should the library provide helpers for passing DB/SFTP connections, or leave entirely to consumers?

2. **Context/cancellation** - Should we thread `context.Context` through evaluation for request timeouts? (More invasive change)

3. **Concurrency** - Is the evaluator safe for concurrent use with separate environments? (Needs audit)

4. **Error types** - Should we define specific error types (`ParseError`, `RuntimeError`) for programmatic handling?

5. **Streaming output** - Current model buffers all output. Is streaming needed for large responses?

---

## Success Criteria

- [ ] CLI behaves identically to current version
- [ ] REPL behaves identically to current version  
- [ ] All existing tests pass
- [ ] New library API has >90% test coverage
- [ ] Library can be imported independently
- [ ] Documentation complete with examples
- [ ] No HTTP/domain knowledge in library code
