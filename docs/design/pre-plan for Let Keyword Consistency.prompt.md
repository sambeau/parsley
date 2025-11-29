# Plan: Let Keyword Consistency

## Current Behavior Analysis

### How `let` Currently Works

1. **Variable Declaration**: `let` is NOT required for variable declarations
   - `x = 5` works (bare assignment)
   - `let x = 5` also works
   - Both create variables in the current scope

2. **Module Exports**: `let` has a special secondary purpose
   - Only `let` bindings are exported from modules
   - Bare assignments are NOT exported
   - `IsLetBinding()` in environment tracks this

3. **No `const` keyword**: All variables are mutable
4. **No `export` keyword**: Export is implicit via `let`

### Example Behavior

```parsley
// In a module file (math.pars)
let add = fn(a, b) { a + b }  // EXPORTED
multiply = fn(a, b) { a * b }  // NOT exported (internal only)

// In main file
import "math.pars" as math
math.add(1, 2)       // Works - 3
math.multiply(2, 3)  // Error - not exported
```

## Issues with Current Design

1. **Confusing dual purpose**: `let` means both "declare variable" AND "export from module"
2. **Easy to forget**: Writing `x = 5` instead of `let x = 5` silently makes it non-exportable
3. **Inconsistent with other languages**: Most languages use explicit `export` keyword
4. **No way to have module-level constants**: Everything is mutable

## Recommendations

### Option 1: Minimal Changes (Document Current Behavior)
- Keep current behavior
- Add clear documentation about `let` = exportable
- Pros: No breaking changes
- Cons: Confusing semantics remain

### Option 2: Add Explicit `export` Keyword (Recommended)
```parsley
// Clearer intent
export let add = fn(a, b) { a + b }
export add = fn(a, b) { a + b }  // Also works

// Internal only
let multiply = fn(a, b) { a * b }
```
- Pros: Explicit, clear, familiar pattern
- Cons: Breaking change for existing modules

### Option 3: Full Consistency Overhaul
- Make `let` required for all declarations
- Add `const` for immutable bindings
- Add `export` for module exports
```parsley
let x = 5           // Mutable, local
const PI = 3.14159  // Immutable, local
export let add = fn(a, b) { a + b }      // Mutable, exported
export const VERSION = "1.0"             // Immutable, exported
```
- Pros: Most consistent, full-featured
- Cons: Significant implementation effort, breaking changes

## Implementation Tasks

### If implementing Option 2 (Recommended):

1. **Lexer Changes** (`pkg/lexer/lexer.go`)
   - Add `EXPORT` token type
   - Recognize `export` keyword

2. **Parser Changes** (`pkg/parser/parser.go`)
   - Parse `export` prefix on statements
   - Create `ExportStatement` AST node or flag

3. **AST Changes** (`pkg/ast/ast.go`)
   - Add export flag to `LetStatement`
   - Or create new `ExportStatement` wrapper

4. **Evaluator Changes** (`pkg/evaluator/evaluator.go`)
   - Track exports separately from let bindings
   - Update module import logic

5. **Migration Path**
   - Keep `let` = export as deprecated behavior
   - Warn when `let` used without `export` in modules
   - Remove deprecated behavior in next major version

## Testing Strategy

1. Test bare `let` still works (backward compat)
2. Test `export let` works
3. Test `export` without `let` works
4. Test non-exported `let` stays local
5. Test module import respects export status

## Decision Required

Which option should we implement?
- [ ] Option 1: Document only
- [ ] Option 2: Add `export` keyword (Recommended)
- [ ] Option 3: Full overhaul with `const` + `export`
