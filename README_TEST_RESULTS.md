# README Testing Results - Issues Found

## Functions That Don't Exist
- `upper()` → should be `toUpper()`
- `lower()` → should be `toLower()`
- `filter()` - NOT FOUND (not in builtins)
- `reduce()` - NOT FOUND (not in builtins)  
- `match()` - NOT FOUND (only regex test with ~ operator)
- `date()` - NOT FOUND (only `now()` and `time()` exist)
- `basename()` - exists but only as PATH property (path.basename)
- `dirname()` - exists but only as PATH property (path.dirname)
- `ext()` - exists but only as PATH property (path.ext)
- `scheme()` - exists but only as URL property (url.scheme)
- `host()` - exists but only as URL property (url.host)
- `path()` - exists but only as URL property (url.path)

## Syntax Issues

### 1. Range Syntax
❌ `for (i in 0..5)` - DOES NOT EXIST
✅ `for (i in [0,1,2,3,4,5])`

The language doesn't support range literals. You must use actual arrays.

### 2. Dictionary Destructuring Aliasing
❌ `let {x: px, y: py} = {x: 10, y: 20}` - NOT SUPPORTED
✅ `let {x, y} = {x: 10, y: 20}` - extract with same names

The parser doesn't support aliasing in destructuring.

### 3. Array Slicing
✅ Array slicing DOES work: `nums[1:]`, `nums[:3]`, `nums[1:3]`
(This was a false positive in my testing)

## Correct Built-in Functions

### String Functions
- `toUpper(str)` - Convert to uppercase  
- `toLower(str)` - Convert to lowercase
- `split(str, delimiter)` - Split string
- `replace(str, pattern, replacement)` - Replace text
- `trim(str)` - Remove whitespace (needs verification)
- `contains(str, substr)` - Check if contains
- `starts_with(str, prefix)` - Check prefix
- `ends_with(str, suffix)` - Check suffix
- `join(array, separator)` - Join array elements
- `len(str)` - Get length

### Array Functions
- `len(array)` - Get length
- `sort(array)` - Sort array
- `reverse(array)` - Reverse array
- `map(fn, array)` - Map function over array
- `sortBy(array, compareFn)` - Custom sort

### Dictionary Functions
- `keys(dict)` - Get keys as array
- `values(dict)` - Get values as array
- `has(dict, key)` - Check if key exists

### Date/Time Functions
- `now()` - Get current time (returns dictionary)
- `time(str|int|dict, delta?)` - Parse/create time
- `year(datetime)`, `month(datetime)`, `day(datetime)` - Extract parts (properties, not functions)
- `hour(datetime)`, `minute(datetime)`, `second(datetime)` - Extract time parts (properties)
- `weekday(datetime)` - Get weekday (property)

### Path/URL Functions
- `path(str)` - Parse path string to path dictionary
- `url(str)` - Parse URL string to URL dictionary

Path properties (not functions):
- `path.basename`
- `path.dirname`
- `path.ext`

URL properties (not functions):
- `url.scheme`
- `url.host`
- `url.path`

### Regular Expression
- `regex(pattern, flags?)` - Create regex dictionary
- `~` operator - Test if string matches regex
- `!~` operator - Test if string doesn't match regex

Note: Regex match with ~ returns array of matches or null, not a separate match() function.

## Features That Need Removal or Correction

1. **Range syntax** - Remove all `0..5` examples
2. **filter() and reduce()** - Remove or note as "not yet implemented"
3. **match() function** - Replace with `~` operator examples
4. **date() function** - Replace with `time()` or `now()`
5. **Aliasing in destructuring** - Remove `{x: newName}` syntax
6. **Path/URL functions** - Show as properties, not function calls:
   - `path.basename` not `basename(path)`
   - `url.scheme` not `scheme(url)`
7. **DateTime extractors** - These might be properties too:
   - Need to verify if `year(date)` or `date.year`

## Template String Interpolation Issue
Test showed: `"Hello, ${name}!"` printed literally as `"Hello, ${name}!"` not interpolated.
Need to verify if this is a bug or wrong syntax.

## HTML Tag Content Issue
Test showed: `<div>"Hello, Alice!"</div>` includes the quotes in output.
This might be intentional (preserves string literals) or might need fixing.
