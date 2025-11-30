# Parsley Library

The `parsley` package provides a public API for embedding the Parsley language interpreter in Go applications.

## Installation

```bash
go get github.com/sambeau/parsley/pkg/parsley
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/sambeau/parsley/pkg/parsley"
)

func main() {
    // Simple evaluation
    result, err := parsley.Eval(`1 + 2`)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(result.String()) // "3"
}
```

## Usage Examples

### Evaluate with Variables

```go
result, err := parsley.Eval(`name ++ "!"`,
    parsley.WithVar("name", "Hello"),
)
// result.String() == "Hello!"
```

### Evaluate a File

```go
result, err := parsley.EvalFile("script.pars",
    parsley.WithSecurity(policy),
)
```

### Custom Logger

```go
logger := parsley.NewBufferedLogger()

result, err := parsley.Eval(`log("hello"); 42`,
    parsley.WithLogger(logger),
)

fmt.Println(logger.String()) // "hello\n"
```

### HTTP Server Example

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // Map HTTP request to Go types
    reqMap := map[string]interface{}{
        "method":  r.Method,
        "path":    r.URL.Path,
        "query":   queryToMap(r.URL.Query()),
        "headers": headersToMap(r.Header),
        "body":    readBody(r),
    }
    
    // Evaluate script with request data
    result, err := parsley.EvalFile("handler.pars",
        parsley.WithVar("request", reqMap),
        parsley.WithLogger(requestLogger),
        parsley.WithSecurity(serverPolicy),
    )
    
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    
    w.Write([]byte(result.String()))
}
```

## API Reference

### Core Functions

#### Eval

```go
func Eval(source string, opts ...Option) (*Result, error)
```

Evaluates Parsley source code and returns the result.

#### EvalFile

```go
func EvalFile(filename string, opts ...Option) (*Result, error)
```

Evaluates a Parsley file and returns the result.

### Options

- `WithVar(name string, value interface{})` - Pre-populate a variable
- `WithEnv(env *Environment)` - Use a pre-configured environment
- `WithSecurity(policy *SecurityPolicy)` - Set file system security policy
- `WithLogger(logger Logger)` - Set the logger for log()/logLine()
- `WithFilename(name string)` - Set the filename for error messages

### Result

```go
type Result struct {
    Value Object  // The Parsley object result
    Err   error   // Runtime error (if any)
}

func (r *Result) GoValue() interface{}  // Convert to Go value
func (r *Result) String() string        // String representation
func (r *Result) IsNull() bool          // Check if null
func (r *Result) IsError() bool         // Check if error
```

### Type Conversion

#### ToParsley

```go
func ToParsley(v interface{}) (Object, error)
```

Converts Go values to Parsley objects:
- `int`, `int64`, etc. → `Integer`
- `float64` → `Float`
- `string` → `String`
- `bool` → `Boolean`
- `[]interface{}` → `Array`
- `map[string]interface{}` → `Dictionary`
- `time.Time` → DateTime dictionary
- `time.Duration` → Duration dictionary
- `nil` → `Null`

#### FromParsley

```go
func FromParsley(obj Object) interface{}
```

Converts Parsley objects to Go values:
- `Integer` → `int64`
- `Float` → `float64`
- `String` → `string`
- `Boolean` → `bool`
- `Array` → `[]interface{}`
- `Dictionary` → `map[string]interface{}`
- `Null` → `nil`

### Loggers

- `StdoutLogger()` - Writes to stdout (default)
- `WriterLogger(w io.Writer)` - Writes to any io.Writer
- `NewBufferedLogger()` - Captures output for later retrieval
- `NullLogger()` - Discards all output

### Error Types

- `ParseError` - Syntax/parse errors
- `RuntimeError` - Runtime evaluation errors

## Testing

```go
func TestMyScript(t *testing.T) {
    logger := parsley.NewBufferedLogger()
    
    result, err := parsley.Eval(`log("test"); 42`,
        parsley.WithLogger(logger),
    )
    
    assert.NoError(t, err)
    assert.Equal(t, int64(42), result.GoValue())
    assert.Equal(t, "test\n", logger.String())
}
```
