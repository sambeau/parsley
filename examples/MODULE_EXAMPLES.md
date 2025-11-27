# Parsley Module System Examples

This document provides comprehensive examples of using Parsley's module system.

## Table of Contents
- [Basic Usage](#basic-usage)
- [Dictionary Destructuring](#dictionary-destructuring)
- [Aliasing Imports](#aliasing-imports)
- [Module Caching](#module-caching)
- [Practical Examples](#practical-examples)

## Basic Usage

### Creating a Module

Modules are just normal Parsley scripts. All `let` bindings are automatically exported:

```parsley
// mathlib.pars
let add = fn(a, b) { a + b }
let multiply = fn(a, b) { a * b }
let PI = 3.141592653589793
```

### Importing a Module

Use the `import()` builtin with a path (string or path literal):

```parsley
// main.pars
let math = import(@./mathlib.pars)

log("2 + 3 =", math.add(2, 3))
log("4 * 5 =", math.multiply(4, 5))
log("PI =", math.PI)
```

## Dictionary Destructuring

Import only what you need using dictionary destructuring:

```parsley
let {add, PI} = import(@./mathlib.pars)

log("Sum:", add(10, 20))
log("Pi:", PI)
```

## Aliasing Imports

Rename imports using the `as` keyword:

```parsley
let {add as sum, multiply as product} = import(@./mathlib.pars)

log(sum(5, 7))      // 12
log(product(3, 4))  // 12
```

## Module Caching

Modules are loaded and evaluated once, then cached:

```parsley
let mod1 = import(@./mathlib.pars)
let mod2 = import(@./mathlib.pars)

log(mod1 == mod2)  // true - same cached object
```

This means:
- Module initialization code runs only once
- Module state is shared across imports
- Efficient - no redundant file reads or evaluations

## Practical Examples

### Example 1: String Utilities Module

```parsley
// modules/strings.pars
let isEmpty = fn(s) { len(s) == 0 }

let capitalize = fn(s) {
    if (len(s) == 0) {
        s
    } else {
        toUpper(s[0:1]) + s[1:len(s)]
    }
}

let repeat = fn(s, n) {
    if (n <= 0) {
        ""
    } else {
        s + repeat(s, n - 1)
    }
}
```

Usage:

```parsley
let {capitalize, repeat} = import(@./modules/strings.pars)

log(capitalize("hello"))  // "Hello"
log(repeat("*", 5))       // "*****"
```

### Example 2: Validation Module

```parsley
// modules/validators.pars
let isEmail = fn(str) {
    str ~ /^[^\s@]+@[^\s@]+\.[^\s@]+$/
}

let isStrongPassword = fn(pwd) {
    len(pwd) >= 8 & 
    (pwd ~ /[A-Z]/) & 
    (pwd ~ /[a-z]/) & 
    (pwd ~ /[0-9]/)
}

let isURL = fn(str) {
    str ~ /^https?:\/\/.+/
}
```

Usage:

```parsley
let validate = import(@./modules/validators.pars)

if (validate.isEmail("user@example.com")) {
    log("Valid email!")
}

if (validate.isStrongPassword("MyPass123")) {
    log("Strong password!")
}
```

### Example 3: Array Utilities Module

```parsley
// modules/arrays.pars
let sum = fn(arr) {
    let total = 0
    for (item in arr) {
        total = total + item
    }
    total
}

let average = fn(arr) {
    if (len(arr) == 0) {
        0
    } else {
        sum(arr) / len(arr)
    }
}

let max = fn(arr) {
    if (len(arr) == 0) {
        0
    } else {
        let maximum = arr[0]
        for (item in arr) {
            if (item > maximum) {
                maximum = item
            }
        }
        maximum
    }
}
```

Usage:

```parsley
let {sum, average, max} = import(@./modules/arrays.pars)

let numbers = [10, 25, 5, 30, 15]
log("Sum:", sum(numbers))      // 85
log("Average:", average(numbers)) // 17
log("Max:", max(numbers))      // 30
```

### Example 4: Configuration Module

```parsley
// config/database.pars
let host = "localhost"
let port = 5432
let database = "myapp"
let connectionString = host + ":" + port + "/" + database
```

Usage:

```parsley
let db = import(@./config/database.pars)

log("Connecting to:", db.connectionString)
log("Host:", db.host)
log("Port:", db.port)
```

### Example 5: Module Composition

You can re-export from other modules:

```parsley
// utils/index.pars
let strings = import(@./strings.pars)
let arrays = import(@./arrays.pars)
let validators = import(@./validators.pars)
```

```parsley
// main.pars
let utils = import(@./utils/index.pars)

log(utils.strings.capitalize("hello"))
log(utils.validators.isEmail("test@example.com"))
```

## Path Resolution

Paths are resolved relative to the importing file:

```
project/
  main.pars
  lib/
    utils.pars
    helpers.pars
```

From `main.pars`:
```parsley
let utils = import(@./lib/utils.pars)
```

From `utils.pars`:
```parsley
// Relative to utils.pars location (inside lib/)
let helpers = import(@./helpers.pars)
```

## Best Practices

1. **Keep modules focused**: Each module should have a single, clear purpose
2. **Use descriptive names**: Module filenames should clearly indicate their contents
3. **Export consistently**: Use clear, consistent naming for exported bindings
4. **Avoid circular dependencies**: Design module hierarchy to prevent circular imports
5. **Document exports**: Add comments describing each exported binding
6. **Use path literals**: Prefer `@./path` over string paths for better tooling support

## Error Handling

### Module Not Found
```parsley
let mod = import(@./nonexistent.pars)
// ERROR: failed to read module file ...
```

### Parse Errors in Module
```parsley
let broken = import(@./broken-syntax.pars)
// ERROR: parse errors in module ./broken-syntax.pars:
//   line X, column Y: ...
```

### Circular Dependencies
```parsley
// a.pars
let b = import(@./b.pars)

// b.pars
let a = import(@./a.pars)

// ERROR: circular dependency detected when importing: ...
```

## Advanced Patterns

### Private Module Variables

By convention, prefix private bindings with underscore:

```parsley
// api.pars
let _privateHelper = fn(x) { x * 2 }
let publicFunction = fn(x) { _privateHelper(x) + 1 }

// Both are exported, but _privateHelper signals "internal use"
```

Users can choose to ignore private exports:

```parsley
let {publicFunction} = import(@./api.pars)
// Only import what's meant for public use
```

### Conditional Exports

```parsley
// feature.pars
let isProduction = false  // Could read from environment

let debug = if (isProduction) {
    fn(msg) { }  // No-op in production
} else {
    fn(msg) { log("DEBUG:", msg) }
}
```

## Complete Working Example

See `examples/module_demo.pars` for a complete demonstration of the module system.

Run it with:
```bash
./pars examples/module_demo.pars
```
