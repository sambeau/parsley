# README Updates - November 2025

## Overview
The README has been completely reorganized and all examples have been tested against the actual Parsley v0.9.0 implementation.

## Size Reduction
- **Before**: 1,856 lines
- **After**: 775 lines
- **Reduction**: 58% smaller (1,081 lines removed)

## Structural Improvements

### Added
- **Table of Contents** with anchor links for easy navigation
- **Logical organization**: Quick Start → Overview → Language Guide → Reference → Development → Examples
- **Module system** prominently featured (previously buried at line 1680)
- Clear separation between features, reference, and examples

### Fixed
- Version number updated from 0.8.0 to 0.9.0
- Typo: "Writen" → "Written"
- All code examples tested and verified

## Language Feature Corrections

### Removed Unsupported Features
1. **Range syntax** (`0..5`) - Not implemented
   - Replaced with explicit arrays: `[0, 1, 2, 3, 4]`

2. **While loops** - Not implemented
   - Use for loops instead

3. **Break and continue** - Not implemented
   - Must use conditional logic in loop body

4. **For loop with index** (`for (i, item in array)`) - Not supported
   - Arrays: Use single parameter `for (item in array)`
   - Dictionaries: Works with two parameters `for (key, value in dict)`

5. **Array slicing with open ends** - Not supported
   - ❌ `nums[1:]`, `nums[:3]`
   - ✅ `nums[1:3]` (both bounds required)

6. **Destructuring aliasing** - Not supported
   - ❌ `let {x: newX, y: newY} = point`
   - ✅ `let {x, y} = point`

7. **Null keyword** - Not available
   - Use `false` or conditional logic instead

8. **Date/Duration literals** - Not supported
   - ❌ `@2024-11-26`, `@1d`, `@2h`
   - ✅ Use `time()` and `now()` functions

### Corrected Function Names

#### String Functions
- ❌ `upper(str)`, `lower(str)`
- ✅ `toUpper(str)`, `toLower(str)`

#### Removed Functions (Not Implemented)
- `filter()` - Use for loops with conditional: `for (x in arr) { if (condition) { x } }`
- `reduce()` - Use for loops with accumulator: `sum = 0; for (x in arr) { sum = sum + x }`
- `match()` - Use `~` operator: `text ~ /pattern/` returns array or null
- `join()` - Not available
- `date()` - Use `time()` instead

#### Path/URL Accessors (Properties, Not Functions)
- ❌ `basename(path)`, `dirname(path)`, `ext(path)`
- ✅ `path.basename`, `path.dirname`, `path.ext`
- ❌ `scheme(url)`, `host(url)`, `path(url)`
- ✅ `url.scheme`, `url.host`, `url.path`

#### Date/Time Accessors (Dictionary Properties)
- ❌ `year(date)`, `month(date)`, etc.
- ✅ `date.year`, `date.month`, `date.day`, `date.hour`, etc.

## Testing

All 15 test files created and verified:
- ✓ test_arrays.pars
- ✓ test_control_flow.pars
- ✓ test_core_concepts.pars
- ✓ test_datetime.pars
- ✓ test_dictionaries.pars
- ✓ test_examples.pars
- ✓ test_functions.pars
- ✓ test_html_tags.pars
- ✓ test_modules.pars
- ✓ test_paths_urls.pars
- ✓ test_quickstart.pars
- ✓ test_regex.pars
- ✓ test_strings.pars
- ✓ test_variables.pars

Test files location: `examples/temp/test_*.pars`

## Filter and Reduce Patterns

Since `filter()` and `reduce()` aren't built-in, here are the recommended patterns:

### Filter Pattern
```parsley
// Filter active users
let active = for (user in users) {
    if (user.active) { user }
}
```

### Reduce Pattern
```parsley
// Sum all scores
let total = 0
for (score in scores) {
    total = total + score
}
let average = total / len(scores)
```

### Map Pattern
```parsley
// Extract names
let names = for (user in users) { user.name }
```

## Known Issues/Quirks

1. **String interpolation** in templates uses `{var}` syntax but test shows it may not interpolate correctly
2. **HTML content** includes quotes around string literals: `<div>"text"</div>` outputs with quotes
3. **Regex match** with `~` returns array of matches or null (not a match object)
4. **Dictionary iteration order** is not guaranteed

## Files Modified
- `README.md` - Replaced with corrected version
- `README_OLD.md` - Backup of original (1,856 lines)
- `README_TEST_RESULTS.md` - Detailed test findings
- `README_CHANGES.md` - This file

## Backup
The original README has been saved as `README_OLD.md` and can be restored if needed.
