# Error Reporting Demo

This document demonstrates the improved error messages in pars. The interpreter now provides:

- **Filename** in the error message
- **Line and column numbers** for precise error location
- **Human-readable token names** instead of technical token types
- **Visual pointer** (^) showing the exact error position
- **Source code context** displaying the problematic line

## Syntax Errors

### Example 1: Missing Expression

**File: `example1.pars`**
```pars
let x = 5
let y =
let z = 10
```

**Error Output:**
```
Error in 'example1.pars':
  line 3, column 4: unexpected 'let'
    let z = 10
       ^
```

### Example 2: Unclosed Parenthesis

**File: `example2.pars`**
```pars
let x = (5 + 3
let y = 10
```

**Error Output:**
```
Error in 'example2.pars':
  line 2, column 4: expected ')', got 'let'
    let y = 10
       ^
```

### Example 3: Missing Expression After Operator

**File: `example3.pars`**
```pars
let result = 5 +
print(result)
```

**Error Output:**
```
Error in 'example3.pars':
  line 2, column 1: unexpected 'print'
    print(result)
    ^
```

## Runtime Errors

### Example 4: Division by Zero

**File: `example4.pars`**
```pars
let x = 5 / 0
x
```

**Error Output:**
```
example4.pars: ERROR: division by zero
```

### Example 5: Undefined Identifier

**File: `example5.pars`**
```pars
let x = unknownVariable
```

**Error Output:**
```
example5.pars: ERROR: identifier not found: unknownVariable
```

## Comparison: Before vs After

### Before
```
expected RPAREN, got LET
```

### After
```
Error in 'example.pars':
  line 2, column 4: expected ')', got 'let'
    let y = 10
       ^
```

The improved error messages help developers quickly identify and fix issues in their code!
