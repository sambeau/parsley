# Module System Design for Parsley

## Executive Summary

This document outlines a minimalist module system for Parsley that leverages existing language features (path literals, dictionary destructuring, and functions) to achieve code organization, reusability, and encapsulation without introducing complex new syntax or multiple built-in functions.

## Design Philosophy

Parsley's module system should embody:
- **Simplicity**: Use existing language constructs where possible
- **Minimalism**: Single built-in function for importing
- **Completeness**: Support all essential module use cases
- **Composability**: Modules are just dictionaries that compose naturally

## Inspiration from Other Languages

### Lesser-known Elegant Approaches

1. **Lua's module system (pre-5.2)**: Modules are simply tables returned from files
   - Clean, minimal, no special syntax
   - Natural composition through table operations

2. **Elixir's explicit imports**: You explicitly choose what to import
   - Prevents namespace pollution
   - Makes dependencies clear

3. **OCaml's module system**: Modules are first-class values
   - Can be passed to functions, stored in data structures
   - Powerful but we'll keep it simpler

4. **Nim's include vs import**: Different mechanisms for different needs
   - `include` = textual inclusion (we skip this)
   - `import` = namespace isolation (we use this)

## Core Concepts

### What Modules Achieve

1. **Code Reusability**: Share functions, constants, and utilities across files
2. **Encapsulation**: Control what is exposed vs. internal
3. **Scope Isolation**: Module code runs in its own scope
4. **Code Organization**: Logical separation of concerns

### Module as Dictionary

In Parsley, a **module is simply a dictionary**. This aligns perfectly with the language's philosophy and existing features.

```parsley
// math_utils.pars
{
    PI: 3.141592653589793,
    
    square: fn(x) { x * x },
    
    distance: fn(x1, y1, x2, y2) {
        let dx = x2 - x1
        let dy = y2 - y1
        sqrt(square(dx) + square(dy))
    }
}
```

The last expression in a file is its export value. This is natural and requires no special syntax.

## Proposed Implementation

### Single Built-in Function: `import(path)`

```parsley
import(path) -> Dictionary
```

- **Input**: A path (string or path literal) to a `.pars` file
- **Output**: A dictionary containing all variables defined in the module's scope
- **Behavior**: 
  - Executes the file in an isolated environment
  - Captures all `let` bindings from the module's scope
  - Returns them as a dictionary
  - Caches results (files are imported once)
  - Paths are relative to the importing file

**Key Innovation**: Modules are just normal Parsley scripts. The module's scope is automatically converted to a dictionary.

### Basic Usage

```parsley
// mathlib.pars - just a normal script!
let add = fn(a, b) { a + b }
let multiply = fn(a, b) { a * b }
let PI = 3.141592653589793
```

```parsley
// main.pars
let math = import(@./mathlib.pars)

log("PI is:", math.PI)
log("2 + 3 =", math.add(2, 3))
log("4 * 5 =", math.multiply(4, 5))
```

**Behind the scenes**: When `mathlib.pars` is imported, all its `let` bindings (`add`, `multiply`, `PI`) are captured and returned as:
```parsley
{add: fn(a,b){...}, multiply: fn(a,b){...}, PI: 3.141592653589793}
```

### Importing Specific Items (Dictionary Destructuring)

```parsley
// Using dictionary destructuring to import only what you need
let {add, multiply} = load(@./mathlib.pars)

log("Sum:", add(5, 7))
log("Product:", multiply(5, 7))
```

### Aliasing Imports

```parsley
let {add as sum, PI as pi} = load(@./mathlib.pars)

log("Pi:", pi)
log("Sum:", sum(10, 20))
```

### Re-exporting (Module Composition)

```parsley
// geometry.pars
let math = load(@./mathlib.pars)

{
    // Re-export from math
    ...math,
    
    // Add geometry-specific functions
    circleArea: fn(radius) {
        math.PI * radius * radius
    },
    
    rectangleArea: fn(width, height) {
        math.multiply(width, height)
    }
}
```

### Nested Module Structure

```parsley
// utils/index.pars
{
    strings: load(@./strings.pars),
    arrays: load(@./arrays.pars),
    math: load(@./math.pars)
}
```

```parsley
// main.pars
let utils = load(@./utils/index.pars)

let {upper} = utils.strings
let {sum} = utils.arrays
```

## Advanced Patterns

### Private vs Public

Since modules are just dictionaries, you control what's exported:

```parsley
// api.pars
let privateHelper = fn(x) { x * 2 }  // Not exported
let publicFunction = fn(x) { 
    privateHelper(x) + 1 
}

{
    // Only export what you want public
    public: publicFunction
}
```

### Module Initialization

```parsley
// config.pars
let env = "production"
let port = 8080

log("Config loaded for:", env)

{
    environment: env,
    port: port,
    debug: env == "development"
}
```

The log statement runs when the module is loaded (once, due to caching).

### Circular Dependencies

**Design Decision**: Disallow circular dependencies. 
- Detect at load time
- Return error: "Circular dependency detected: A -> B -> A"
- Encourages better architecture

### Conditional Loading

```parsley
let utils = if (needsUtils) load(@./utils.pars) else null

if (utils) {
    utils.doSomething()
}
```

### Dynamic Module Paths

```parsley
let moduleName = "math"
let module = load(@./modules + moduleName + ".pars")
```

## Implementation Details

### File Loading (`load` builtin)

```go
// In getBuiltins()
"load": {
    Fn: func(args ...Object, env *Environment) Object {
        if len(args) != 1 {
            return newError("wrong number of arguments. got=%d, want=1", len(args))
        }
        
        // Accept path dictionary or string
        var pathStr string
        switch arg := args[0].(type) {
        case *Dictionary:
            // Handle path literal
            if typeVal, ok := arg.Pairs["__type"]; ok {
                if typeExpr, ok := typeVal.(*ast.StringLiteral); ok {
                    if typeExpr.Value == "path" {
                        pathStr = pathDictToString(arg)
                    }
                }
            }
        case *String:
            pathStr = arg.Value
        default:
            return newError("argument to `load` must be a path or string")
        }
        
        // Resolve path relative to current file
        // Check module cache
        // Load and execute file
        // Cache result
        // Return last expression value
    }
}
```

### Module Cache

```go
type ModuleCache struct {
    modules map[string]Object  // path -> result
    loading map[string]bool     // path -> currently loading (for cycle detection)
}

var moduleCache = &ModuleCache{
    modules: make(map[string]Object),
    loading: make(map[string]bool),
}
```

### Execution Context

Each module gets its own environment:
- Inherits global built-ins
- No access to importing file's variables
- Variables defined in module don't leak out

```go
func loadModule(path string) Object {
    // Create isolated environment
    moduleEnv := NewEnvironment()
    
    // Load file content
    content := readFile(path)
    
    // Parse
    program := parse(content)
    
    // Evaluate
    result := Eval(program, moduleEnv)
    
    return result
}
```

### Path Resolution

Paths are resolved relative to the **importing file**, not the current working directory:

```
project/
  main.pars              (loads ./lib/utils.pars)
  lib/
    utils.pars           (loads ./helpers.pars)
    helpers.pars
```

From `main.pars`: `load(@./lib/utils.pars)`
From `utils.pars`: `load(@./helpers.pars)` (relative to utils.pars location)

## What Data Types Are Importable?

**Answer**: Any valid Parsley value:
- Dictionaries (most common - collections of exports)
- Functions (module that exports a single function)
- Strings, Numbers, Booleans (module that exports a constant)
- Arrays (module that exports a list)
- null (valid but unusual)

Examples:

```parsley
// version.pars
"1.0.0"
```

```parsley
// main.pars
let version = load(@./version.pars)
log("Version:", version)  // "1.0.0"
```

## Symbol Renaming

**Q**: Is it important to import with a different name?

**A**: Yes, and we get it for free with dictionary destructuring's `as` syntax:

```parsley
let {longFunctionName as fn} = load(@./module.pars)
```

This is already implemented and requires no new syntax.

## Edge Cases and Pitfalls

### 1. Scope Binding Issues

**Problem**: What if a module function needs to reference other module variables?

**Solution**: Use closures (already works in Parsley):

```parsley
// module.pars
let secret = 42

let getSecret = fn() { secret }  // Closure captures 'secret'
let useSecret = fn(x) { x + secret }
```

The module's environment is preserved through closure.

### 2. Mutation of Imported Values

**Problem**: If multiple files load the same module and mutate it, what happens?

**Solution**: Module cache returns the **same object**:

```parsley
// counter.pars
let value = 0
```

```parsley
// a.pars
let counter = import(@./counter.pars)
counter.value = counter.value + 1
```

```parsley
// b.pars  
let counter = import(@./counter.pars)
log(counter.value)  // Will see mutations from a.pars
```

**Design decision**: This is a feature, not a bug. Allows shared state (like singletons). 
Document clearly that modules are cached and shared.

### 3. Module Not Found

```parsley
let missing = load(@./nonexistent.pars)
// ERROR: Module not found: ./nonexistent.pars
```

### 4. Parse/Evaluation Errors in Module

```parsley
// broken.pars
let x = unknownFunction()
```

```parsley
// main.pars
let broken = load(@./broken.pars)
// ERROR: In module ./broken.pars: identifier not found: unknownFunction
```

Error messages must indicate which module failed.

### 5. Path Types

Support both:
- Path literals: `load(@./module.pars)`
- Strings: `load("./module.pars")`

Path literals are preferred (type-safe, IDE-friendly).

## Should We Tackle File Functions First?

**Analysis**: The module system needs to read files, so we need basic file I/O. However, the initial implementation can be **load-only**:

### Minimal Viable Module System

1. **Phase 1** (This proposal):
   - `load(path)` built-in
   - File reading capability (internal to load)
   - Module caching
   - Circular dependency detection
   
2. **Phase 2** (Future):
   - `read(path)` - read file as string
   - `write(path, content)` - write string to file
   - `exists(path)` - check file existence
   - `listDir(path)` - list directory contents

The module system only needs read capability, which can be internal to `load()`. Full file I/O can come later.

## Examples: Complete Use Cases

### Example 1: Utility Library

```parsley
// lib/strings.pars
{
    isEmpty: fn(s) { len(s) == 0 },
    capitalize: fn(s) { 
        if (len(s) == 0) s else toUpper(s[0:1]) + s[1:len(s)]
    },
    reverse: fn(s) {
        let chars = for(i in range(len(s) - 1, -1, -1)) { s[i] }
        toString(chars)
    }
}
```

```parsley
// app.pars
let {capitalize, isEmpty} = load(@./lib/strings.pars)

let name = "alice"
log(capitalize(name))  // "Alice"
log(isEmpty(name))     // false
```

### Example 2: Configuration Management

```parsley
// config/database.pars
{
    host: "localhost",
    port: 5432,
    database: "myapp",
    pool: {
        min: 2,
        max: 10
    }
}
```

```parsley
// app.pars
let dbConfig = load(@./config/database.pars)

log("Connecting to:", dbConfig.host + ":" + dbConfig.port)
```

### Example 3: Validators

```parsley
// validators.pars
{
    isEmail: fn(str) {
        str ~ /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    },
    
    isStrongPassword: fn(pwd) {
        len(pwd) >= 8 & 
        (pwd ~ /[A-Z]/) & 
        (pwd ~ /[a-z]/) & 
        (pwd ~ /[0-9]/)
    },
    
    isURL: fn(str) {
        str ~ /^https?:\/\/.+/
    }
}
```

```parsley
// main.pars
let validate = load(@./validators.pars)

let email = "user@example.com"
if (validate.isEmail(email)) {
    log("Valid email")
}
```

### Example 4: Plugin Architecture

```parsley
// plugins/markdown.pars
{
    name: "Markdown Renderer",
    version: "1.0",
    render: fn(text) {
        // markdown rendering logic
        text
    }
}
```

```parsley
// plugins/textile.pars  
{
    name: "Textile Renderer",
    version: "1.0",
    render: fn(text) {
        // textile rendering logic
        text
    }
}
```

```parsley
// app.pars
let plugins = [
    load(@./plugins/markdown.pars),
    load(@./plugins/textile.pars)
]

for (plugin in plugins) {
    log("Loaded:", plugin.name, "v" + plugin.version)
}

let content = "# Hello World"
let rendered = plugins[0].render(content)
```

## Testing Strategy

### Unit Tests for Module System

```parsley
// tests/module_loading_test.pars
let testBasicLoad = fn() {
    let mod = load(@./fixtures/simple.pars)
    assert(mod.value == 42)
}

let testDestructuring = fn() {
    let {add} = load(@./fixtures/math.pars)
    assert(add(2, 3) == 5)
}

let testCaching = fn() {
    let mod1 = load(@./fixtures/counter.pars)
    let mod2 = load(@./fixtures/counter.pars)
    // Should be same object
    assert(mod1 == mod2)
}
```

## Migration Path

Existing Parsley code continues to work unchanged. Modules are purely opt-in.

### Before (single file):

```parsley
// app.pars
let square = fn(x) { x * x }
let distance = fn(x1, y1, x2, y2) {
    sqrt(square(x2-x1) + square(y2-y1))
}

log(distance(0, 0, 3, 4))
```

### After (with modules):

```parsley
// math.pars
{
    square: fn(x) { x * x },
    distance: fn(x1, y1, x2, y2) {
        sqrt(square(x2-x1) + square(y2-y1))
    }
}
```

```parsley
// app.pars
let {distance} = load(@./math.pars)
log(distance(0, 0, 3, 4))
```

## Summary: Why This Design is Parsleyish

1. **Minimal syntax**: Only one new built-in (`load`), no new keywords
2. **Leverages existing features**: Dictionaries, destructuring, closures
3. **Natural composition**: Modules compose like any other data
4. **Simple mental model**: "A module is a dictionary returned from a file"
5. **Explicit over implicit**: You see exactly what's being imported
6. **No magic**: Clear execution model, predictable behavior
7. **Flexible**: Supports multiple patterns without rigid structure

## Next Steps for Implementation

1. ✅ Design review (this document)
2. ⏳ Implement `load()` built-in
   - File reading
   - Path resolution  
   - Module caching
   - Circular dependency detection
3. ⏳ Add tests
4. ⏳ Update documentation
5. ⏳ Create example modules
6. ⏳ Consider file I/O built-ins for Phase 2

## Open Questions

1. **Error handling**: Should `load()` have a try-catch variant?
   - `tryLoad(path)` returns `null` on error instead of error object?
   
2. **Reloading**: Should there be a way to force reload?
   - `load(path, force: true)` or separate `reload(path)`?
   
3. **Module metadata**: Should modules have standard metadata?
   - `{__module: {name, version, exports}}`?

4. **Standard library**: Should Parsley ship with standard modules?
   - `std/math.pars`, `std/strings.pars`, etc.?

These can be decided during implementation and testing.
