# Plan: Method-Call Style API Redesign

**TL;DR**: Consolidate ~15 function-style builtins into method-style calls on their respective types (e.g., `toUpper(str)` → `str.upper()`), with full chaining support and null propagation. Keep factory functions, type converters, and math functions as standalone. Existing properties remain unchanged.

## Proposed API by Type

### 1. **String** - Add method support
| Current | Proposed |
|---------|----------|
| `toUpper(str)` | `str.upper()` |
| `toLower(str)` | `str.lower()` |
| `split(str, delim)` | `str.split(delim)` |
| `replace(str, pat, rep)` | `str.replace(pat, rep)` |
| `trim(str)` | `str.trim()` |
| `len(str)` | `str.length()` |

### 2. **Array** - Add method support
| Current | Proposed |
|---------|----------|
| `map(fn, arr)` | `arr.map(fn)` |
| `sort(arr)` | `arr.sort()` |
| `sortBy(arr, fn)` | `arr.sortBy(fn)` |
| `reverse(arr)` | `arr.reverse()` |
| `len(arr)` | `arr.length()` |
| `format(arr, style, loc)` | `arr.format(style?, loc?)` |

### 3. **Dictionary** - Add method support
| Current | Proposed |
|---------|----------|
| `keys(dict)` | `dict.keys()` |
| `values(dict)` | `dict.values()` |
| `has(dict, key)` | `dict.has(key)` |

### 4. **Number** (Integer/Float) - Add method support
| Current | Proposed |
|---------|----------|
| `formatNumber(n, loc)` | `n.format(loc?)` |
| `formatCurrency(n, code, loc)` | `n.currency(code, loc?)` |
| `formatPercent(n, loc)` | `n.percent(loc?)` |

### 5. **Datetime** - Add method
| Current | Proposed |
|---------|----------|
| `formatDate(dt, style, loc)` | `dt.format(style?, loc?)` |

### 6. **Duration** - Add method
| Current | Proposed |
|---------|----------|
| `format(dur, loc)` | `dur.format(loc?)` |

## Keep as Standalone Functions

| Category | Functions |
|----------|-----------|
| **Factories** | `now()`, `time()`, `path()`, `url()`, `regex()`, `tag()` |
| **Type converters** | `toInt()`, `toFloat()`, `toNumber()`, `toString()`, `toArray()`, `fromArray()` |
| **Math** | `sin`, `cos`, `tan`, `sqrt`, `pow`, `pi`, `round` |
| **I/O & Modules** | `log()`, `logLine()`, `import()` |

## Method Chaining Support

Each method returns a typed value enabling fluent chains:

```parsley
// String → String → Array → Array
"  Hello World  ".trim().upper().split(" ").reverse()
// → ["WORLD", "HELLO"]

// Array → Array → String
["banana", "apple", "cherry"].sort().format()
// → "apple, banana, and cherry"

// Datetime → String → String
now().format("long", "fr-FR").upper()
// → "27 NOVEMBRE 2025"
```

**Return type reference:**
| Method | Returns | Chainable to |
|--------|---------|--------------|
| `str.upper/lower/trim/replace` | String | String methods |
| `str.split(d)` | Array | Array methods |
| `str.length()` | Integer | Number methods |
| `arr.map/sort/sortBy/reverse` | Array | Array methods |
| `arr.format()` | String | String methods |
| `arr.length()` | Integer | Number methods |
| `n.format/currency/percent` | String | String methods |
| `dt.format()` | String | String methods |
| `dur.format()` | String | String methods |
| `dict.keys/values` | Array | Array methods |
| `dict.has(key)` | Boolean | — |

## Properties vs Methods

Properties and methods are distinguished by context:

| Expression | Type | Behavior |
|------------|------|----------|
| `x = d.unix` | Read | Property access or method call |
| `d.unix = x` | Write | Property assignment only |

**Rule**: Keep writable properties as properties. Only add methods for read-only operations.

| Type | Writable Properties (keep) | Read-only Methods (new) |
|------|---------------------------|------------------------|
| **Datetime** | `year`, `month`, `day`, `hour`, `minute`, `second`, `unix` | `format()`, `dayOfYear()`, `week()`, `timestamp()` |
| **Duration** | `months`, `seconds` | `format()` |
| **Path** | `components`, `absolute` | `isAbsolute()`, `isRelative()` |
| **URL** | `scheme`, `host`, `port`, `path`, `query`, `fragment`, `username`, `password` | `origin()`, `pathname()`, `search()`, `href()` |
| **Dictionary** | Any key | `keys()`, `values()`, `has()` |

## Null Propagation

Methods called on `null` propagate `null` through the chain instead of erroring:

```parsley
let name = null
name.upper()              // → null (no error)
name.upper().split(",")   // → null (chain continues)

let user = { name: "Alice", address: null }
user.address.city         // → null (not an error)
```

**Benefits for templates:**
- Missing/optional data doesn't break rendering
- Cleaner code without defensive null checks
- Combines well with `??` operator for defaults:
  ```parsley
  {user.nickname.upper() ?? "ANONYMOUS"}
  ```

## Future: Optional Chaining (`.?`)

For cases where strict null checking is preferred, add explicit optional chaining syntax:

| Syntax | Behavior |
|--------|----------|
| `x.method()` | Propagates null (lenient) |
| `x!.method()` | Errors if x is null (strict) |

```parsley
// Lenient — null propagates (default)
user.address.city.upper()  // → null if any part is null

// Strict — error if null (opt-in)
user!.address!.city.upper()  // ERROR if user or address is null
```

This allows explicit strictness when data is required:
```parsley
// Title is required, nickname is optional
<h1>{title!.upper()}</h1>
<p>Welcome, {user.nickname.upper() ?? user.name}</p>
```

## Implementation Steps

1. Extend `evalDotExpression` in [pkg/evaluator/evaluator.go](pkg/evaluator/evaluator.go) to handle String, Array, Integer, Float types
2. Add method dispatch for each primitive type
3. Implement null propagation for method calls
4. Convert read-only computed properties to methods
5. (Future) Add `!.` strict chaining operator
6. Update documentation

## Further Considerations

1. **Strict syntax choice?** `x!.method()` vs `x.method()!` vs `x?.method()` (inverted meaning)
2. **Null in arguments?** Should `str.split(null)` error or use default behavior?
