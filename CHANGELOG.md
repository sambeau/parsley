# Changelog

All notable changes to Parsley will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

---

## [0.12.2] - 2025-11-29

### Added
- **Enhanced REPL**: Professional interactive shell with modern editing features
  - Cursor movement: ← → to move within line, ↑ ↓ for command history
  - Tab completion: Auto-complete keywords and builtins (let, if, for, log, file, etc.)
  - Multi-line input: Automatic detection of unclosed braces/brackets/parentheses
  - Persistent history: Commands saved to `~/.parsley_history` across sessions
  - Improved controls: Ctrl+C to abort current line, Ctrl+D to exit
  - Visual feedback: ".." continuation prompt for multi-line expressions
  - Better UX: "^C (cleared)" message when aborting multi-line input

### Changed
- REPL now uses `github.com/peterh/liner` library for line editing
- Replaced basic `bufio.Scanner` with feature-rich interactive input

---

## [0.12.1] - 2025-11-29

### Added
- **Local Directory Methods**: Directory manipulation for local file paths
  - `file(@/path).mkdir(options?)` - Create directory
  - `file(@/path).rmdir(options?)` - Remove directory
  - `dir(@/path).mkdir(options?)` - Create directory
  - `dir(@/path).rmdir(options?)` - Remove directory
  - Options: `{parents: true}` for mkdir, `{recursive: true}` for rmdir
  - Achieves feature parity with SFTP directory operations
- Example script: `examples/directory_operations_demo.pars`
- Comprehensive test suite: `tests/local_directory_test.go` with 8 tests

---

## [0.12.0] - 2025-12-01

### Added
- **SFTP Support**: Full SFTP file operations with network operators
  - `SFTP(url, options?)` - Create SFTP connection with authentication
  - SSH key authentication (preferred) and password authentication
  - Connection caching by user@host:port for efficiency
  - `known_hosts` verification for security
  - Connection timeout support via options
- **SFTP Network Operators**: Use existing network syntax for SFTP
  - `<=/=` - Read from SFTP (e.g., `data <=/= conn(@/file.json).json`)
  - `=/=>` - Write to SFTP (e.g., `data =/=> conn(@/file.json).json`)
  - `=/=>>` - Append to SFTP file
- **SFTP Format Support**: All file formats work over SFTP
  - `.json` - JSON data
  - `.text` - Plain text
  - `.csv` - CSV data
  - `.lines` - Line-separated data
  - `.bytes` - Binary data (array of integers)
  - `.file` - Auto-detect format
  - `.dir` - Directory listing with metadata
- **SFTP Directory Operations**: Manage remote directories
  - `.mkdir(options?)` - Create directory
  - `.rmdir()` - Remove empty directory
  - `.remove()` - Delete file
- **SFTP Connection Management**:
  - `.close()` - Close connection and free resources
  - Callable syntax: `conn(@/path)` returns file handle
  - Error capture pattern: `{data, error} <=/= conn(@/file).json`
- **Dependencies**:
  - Added `github.com/pkg/sftp v1.13.10`
  - Upgraded `golang.org/x/crypto` to `v0.45.0`
- Comprehensive test suite: `sftp_test.go` with 11 test suites
- Example script: `examples/sftp_demo.pars`

### Implementation Notes
- SFTP follows same pattern as database connections (connection-based)
- Path-first syntax matches local file I/O: `conn(@/path).format`
- No special security flags (aligns with HTTP/Fetch pattern)
- SSH keys stored in `~/.ssh/` directory (standard location)
- Connection errors returned via error capture pattern

---

## [0.11.0] - 2025-11-30

### Added
- **Process Execution**: Execute external commands and capture output
  - `COMMAND(binary, args?, options?)` - Create command handle
  - `<=#=>` operator - Execute command with optional input
  - Command options support: `env` (environment variables), `dir` (working directory), `timeout` (duration)
  - Result dictionary with: `stdout`, `stderr`, `exitCode`, `error`
  - Security integration with `--allow-execute` flags
- **JSON Format Functions**: Parse and stringify JSON data
  - `parseJSON(string)` - Parse JSON string to Parsley objects
  - `stringifyJSON(object)` - Convert Parsley objects to JSON string
  - Supports objects, arrays, strings, numbers, booleans, and null
- **CSV Format Functions**: Parse and stringify CSV data
  - `parseCSV(string, options?)` - Parse CSV string to array of arrays or dictionaries
  - `stringifyCSV(array)` - Convert array of arrays to CSV string
  - Options: `header: true` parses first row as column names
- Comprehensive test suite for process execution and format conversion
- Example script: `examples/process_demo.pars`

### Changed
- Added new lexer token `EXECUTE_WITH` for `<=#=>` operator
- Extended AST with `ExecuteExpression` node type
- Parser supports process execution syntax at EQUALS precedence level

### Documentation
- Updated reference.md with process execution documentation
- Added format conversion function examples

---

## [0.10.0] - 2025-11-29

### Added
- **File System Security**: Command-line flags to restrict file system access
  - `--restrict-read=PATHS` - Deny reading from comma-separated paths (blacklist)
  - `--no-read` - Deny all file reads
  - `--allow-write=PATHS` - Allow writing to comma-separated paths (whitelist)
  - `--allow-write-all` / `-w` - Allow unrestricted writes
  - `--allow-execute=PATHS` - Allow executing scripts from paths (whitelist)
  - `--allow-execute-all` / `-x` - Allow unrestricted script execution
  - Security checks integrated into all file operations, directory listings, and module imports
  - Comprehensive test coverage for all security scenarios

### Changed
- **BREAKING**: Write operations now denied by default (use `--allow-write` or `-w` to enable)
- **BREAKING**: Script execution (module imports) now denied by default (use `--allow-execute` or `-x` to enable)
- Read operations remain unrestricted by default (use `--restrict-read` or `--no-read` to restrict)

### Security
- File writes now require explicit permission via command-line flags
- Module imports now require explicit execute permission
- Path validation ensures security restrictions are enforced consistently

---

## [0.9.18] - 2025-11-29

### Added
- **File Handle Methods**: File objects now have methods
  - `.remove()` - Delete files from the filesystem
  - Returns `null` on success, error on failure

---

## [0.9.17] - 2025-11-29

### Added
- **Range Operator (`..`)**: Create inclusive integer ranges
  - Forward ranges: `1..5` → `[1, 2, 3, 4, 5]`
  - Reverse ranges: `10..1` → `[10, 9, 8, 7, 6, 5, 4, 3, 2, 1]`
  - Works seamlessly with for loops and array methods
  - Supports negative numbers and single-element ranges

---

## [0.9.16] - 2025-11-29

### Added
- **Enhanced Array Operators**: Powerful set operations and transformations
  - **Scalar Concatenation (`++`)**: `1 ++ [2,3]` → `[1, 2, 3]`, `[1,2] ++ 3` → `[1, 2, 3]`
  - **Array Intersection (`&&`)**: `[1,2,3] && [2,3,4]` → `[2, 3]`
  - **Array Union (`||`)**: `[1,2] || [2,3]` → `[1, 2, 3]` (deduplicated)
  - **Array Subtraction (`-`)**: `[1,2,3,4] - [2,4]` → `[1, 3]`
  - **Array Chunking (`/`)**: `[1,2,3,4,5,6] / 2` → `[[1,2], [3,4], [5,6]]`
  - **String Repetition (`*`)**: `"abc" * 3` → `"abcabcabc"`
  - **Array Repetition (`*`)**: `[1,2] * 3` → `[1, 2, 1, 2, 1, 2]`
  
- **Dictionary Set Operations**:
  - **Dictionary Intersection (`&&`)**: `{a:1, b:2} && {b:3, c:4}` → `{b: 2}` (keeps left values)
  - **Dictionary Subtraction (`-`)**: `{a:1, b:2} - {b:0}` → `{a: 1}` (removes keys)

---

## [0.9.15] - 2025-11-29

### Added
- **SQLite Database Support**: First-class database integration
  - Database operators: `<=?=>` (single row), `<=??=>` (multiple rows), `<=!=>` (execute)
  - `SQLITE()` connection factory for in-memory and file-based databases
  - Transaction support: `db.begin()`, `db.commit()`, `db.rollback()`
  - Connection methods: `db.ping()`, `db.close()`
  - Automatic type mapping between SQLite and Parsley types
  - Error handling with `db.lastError` property

### Changed
- **Path Cleaning**: Paths are now automatically cleaned using Rob Pike's cleanname algorithm
  - Eliminates `.` and resolves `..` components
  - Normalizes path structure for consistency

---

## [0.9.13] - 2025-11-28

### Added
- **Interpolated Datetime Templates**: Dynamic datetime construction with `@(...)`
  - Date interpolation: `@(2024-{month}-{day})`
  - Datetime interpolation: `@({year}-12-25T{hour}:30:00)`
  - Time interpolation: `@({h}:{m})`
  - Automatic kind detection (date, datetime, time)

---

## [0.9.12] - 2025-11-28

### Added
- **Interpolated Path Templates**: Dynamic path construction with `@(...)`
  - Example: `@(./data/{name}.json)`
  - Supports expressions: `@(./file{n + 1}.txt)`
  
- **Interpolated URL Templates**: Dynamic URL construction with `@(...)`
  - Example: `@(https://api.example.com/{version}/users)`
  - Supports port, fragment, and query interpolation

---

## [0.9.11] - 2025-11-28

### Added
- **HTTP Requests (Fetch)**: Fetch content from URLs using the `<=/=` operator
  - Format factories: `JSON(url)`, `text(url)`, `YAML(url)`, `lines(url)`, `bytes(url)`
  - Request options: `method`, `headers`, `body`, `timeout`
  - Response destructuring: `{data, error, status, headers} <=/= request`
  - Support for GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS
  - Error capture pattern for robust error handling
  - Example: `let users <=/= JSON(@https://api.example.com/users)`
- **Datetime "Kinds"**: Datetime literals now track their display format
  - `@2024-11-26` → kind: `"date"` (displays as "2024-11-26")
  - `@2024-11-26T15:30:00` → kind: `"datetime"` (displays as ISO datetime)
  - `@12:30` → kind: `"time"` (displays as "12:30")
  - `@12:30:45` → kind: `"time_seconds"` (displays as "12:30:45")
  - Kind preserved through arithmetic operations

### Changed
- Improved consistency in `toString()` and `toDict()` methods across all types

### Documentation
- Added comprehensive HTTP Requests documentation in reference.md
- Added HTTP request examples to README.md  
- Created `examples/fetch_demo.pars` with working examples

**Note**: HTTP Requests feature was implemented in v0.9.11 but documentation was added retroactively in v0.11.0.

---

## [0.9.10] - 2025-11-28

### Added
- **Markdown Importer**: Load markdown files with YAML frontmatter
  - `let post <== MD(@./blog.md)` returns `{html, title, date, ...}` from frontmatter
  - Automatic datetime parsing for ISO-formatted dates in frontmatter
  - Separates metadata from rendered HTML content

- **YAML Support**: Import YAML files as dictionaries
  - Integrated with markdown frontmatter parsing

---

## [0.9.9] - 2025-11-28

### Added
- **Nullish Coalescing Operator (`??`)**: Return default value if null
  - `value ?? "default"` returns default only when value is null
  - Short-circuit evaluation
  - Chainable: `a ?? b ?? c ?? "fallback"`
  
- **File Handle Objects**: First-class file I/O with format binding
  - Factories: `JSON()`, `CSV()`, `text()`, `lines()`, `bytes()`, `file()`
  - Properties: `.exists`, `.size`, `.modified`, `.basename`, `.stem`
  - Type-aware reading and writing based on file format

- **File I/O Operators**:
  - **Read (`<==`)**: `let data <== JSON(@./config.json)`
  - **Write (`==>`)**: `myDict ==> JSON(@./output.json)`
  - **Append (`==>>`)**: `newLine ==>> lines(@./log.txt)`
  - Error capture pattern: `let {data, error} <== JSON(@./file.json)`

- **Directory Operations**:
  - `dir()` for directory handles with `.exists`, `.count`, `.files` properties
  - `files()` for glob pattern matching: `files(@./src/**/*.pars)`

- **SVG Importer**: Load SVG files as reusable components
  - `let Arrow <== SVG(@./icons/arrow.svg)` strips XML prolog
  - Use as component: `<button><Arrow/> Next</button>`

- **Module Export Keyword**: Explicit `export` for clearer module APIs
  - `export let version = "1.0"` marks exported bindings
  - Backward compatible: `let` bindings still exported by default

### Changed
- `this` binding now works correctly in dictionary methods
- Merged development documentation sections for clarity

---

## [0.9.7] - 2025-11-27

### Added
- **Internationalization (i18n)**: Locale-aware formatting
  - Number formatting with locales: `1234.format("de-DE")` → `"1.234"`
  - Currency formatting: `99.currency("EUR", "de-DE")` → `"99,00 €"`
  - Date formatting: `date.format("long", "fr-FR")` → `"27 novembre 2024"`
  - Duration formatting: `@1d.format("de-DE")` → `"morgen"`

- **Array Formatting**: Convert arrays to natural language
  - `["a","b","c"].format()` → `"a, b, and c"`
  - `["a","b"].format("or")` → `"a or b"`

### Changed
- Number formatting handles locales natively (eliminates need for separate Decimal type)

---

## [0.9.2] - 2025-11-27

### Added
- **For Loop Indexing**: Access element index in for loops
  - `for (i, item in items) { <li>{i}. {item}</li> }`
  - Works with arrays and dictionaries

---

## [0.9.1] - 2025-11-27

### Added
- **Open-Ended Slicing**: Slice arrays and strings without end index
  - `arr[2:]` - from index 2 to end
  - `str[:3]` - from start to index 3
  - Complements existing range slicing: `arr[1:3]`

---

## [0.9.0] - 2025-11-27

### Added
- **Module System**: Import and organize code across files
  - `import(@./module.pars)` to load modules
  - Dictionary destructuring: `let {add, PI} = import(@./math.pars)`
  - Module caching (files loaded once)
  - Circular dependency detection
  - Only `let` bindings are exported (private scope for other variables)

- **Path Type**: File system paths as first-class values
  - Path literals: `@./config.json`, `@/usr/local/bin`
  - Properties: `.basename`, `.ext`, `.stem`, `.dirname`, `.string`
  - Methods: `.isAbsolute()`, `.isRelative()`

- **URL Type**: Web addresses as first-class values
  - URL literals: `@https://api.example.com/users`
  - Properties: `.scheme`, `.host`, `.port`, `.path`, `.query`
  - Methods: `.origin()`, `.pathname()`, `.href()`

- **Computed Properties**: Dynamic property access for pseudo-types
  - Datetime: `.iso`, `.unix`, `.dayOfYear`, `.week`
  - Path and URL computed properties

### Changed
- Updated VSCode extension grammar to v0.9.0
- Major README reorganization with tested examples

---

## [0.8.0] - 2025-11-27

### Added
- **Better Error Messages**: Runtime errors now show context
  - File and line information in error messages
  - Improved error formatting for debugging

### Fixed
- Assignment in if conditions now produces helpful error message

---

## [0.7.0] - 2025-11-26

### Added
- **Datetime Literals**: `@` syntax for dates and times
  - Date: `@2024-11-26`
  - DateTime: `@2024-11-26T15:30:00`
  - Time: `@12:30` or `@12:30:45`

- **Duration Literals**: Time spans with `@` syntax
  - `@1d` (1 day), `@2h` (2 hours), `@30m` (30 minutes)
  - Combined: `@1d2h30m`
  - Negative durations: `@-1d` (yesterday)
  - Relative formatting: `@1d.format()` → `"tomorrow"`

---

## [0.6.0] - 2025-11-26

### Added
- **Regular Expressions**: Pattern matching with regex literals
  - Regex literals: `/pattern/`, `/pattern/i`, `/pattern/g`
  - Match operator: `"test" ~ /\w+/` → `["test"]`
  - Not-match operator: `"abc" !~ /\d/` → `true`
  - Methods: `.test()`, `.format()`
  - Dynamic creation: `regex("\\d+", "i")`

- **Datetime Type**: Date and time manipulation
  - Factory: `now()`, `time("2024-11-26")`, `time(timestamp)`
  - Properties: `.year`, `.month`, `.day`, `.hour`, `.minute`, `.second`
  - Methods: `.format()`, `.format(style)`, `.format(style, locale)`
  - Arithmetic: datetime + seconds, datetime - datetime

- **Datetime Operators**: Work with dates and times
  - Addition: `@2024-12-25 + 86400` (add seconds)
  - Subtraction: `dt1 - dt2` (duration between dates)
  - Comparisons: `dt1 < dt2`

---

## [0.5.0] - 2025-11-26

### Added
- **HTML/XML Components**: Reusable tag components
  - Component definition: `let Card = fn({title, body}) { <div>...</div> }`
  - Component usage: `<Card title="Hello" body="World" />`
  - Uppercase tags are treated as components

- **Tag Pairs**: Proper opening and closing tags
  - `<div>content</div>` for paired tags
  - Nested tag support
  - Self-closing tags: `<img src="..." />`

- **Function Parameter Destructuring**: Extract values from dictionary arguments
  - `fn({title, contents}) { ... }` destructures props
  - Works with components and regular functions

### Changed
- Tag interpolation uses `{...}` instead of `${...}` for consistency

---

## [0.4.0] - 2025-11-26

### Added
- **Dictionary Destructuring**: Extract values from dictionaries
  - Variable destructuring: `let {x, y} = {x: 10, y: 20}`
  - Works in assignments and function parameters

- **Pretty Printer**: Format HTML output
  - Command line flag: `-pp` or `--pretty`
  - Indented, readable HTML output

### Changed
- Improved error messages with better context
- Better runtime error formatting

---

## [0.3.0] - 2025-11-25

### Added
- **Singleton Tags**: Self-closing HTML tags
  - `<img src="photo.jpg" />`, `<br/>`, `<meta charset="utf-8" />`

- **Debug & Logging Functions**:
  - `log(...)` for output to stdout
  - `logLine(...)` with file:line prefix for debugging

- **Array Methods**: Enhanced array manipulation
  - `.sort()` - ascending sort
  - `.reverse()` - reverse array order
  - `.sortBy(fn)` - custom comparison function

### Changed
- Renamed `print()` to `toString()` for clarity
- Improved REPL prompt display

---

## [0.2.0] - 2025-11-25

### Added
- **Dictionaries**: Key-value objects
  - Literal syntax: `{name: "Alice", age: 30}`
  - Dot access: `user.name`
  - Bracket access: `user["age"]`
  - Methods: `.keys()`, `.values()`, `.has(key)`

- **Array Destructuring**: Extract array elements to variables
  - `let a, b, c = [1, 2, 3]` assigns each element

- **Modulo Operator**: `%` for remainder
  - `10 % 3` → `1`

- **Type Converters**: String to number conversion
  - `toInt(str)`, `toFloat(str)`, `toNumber(str)`

### Changed
- Better error messages showing line numbers
- Fixed errors appearing on line after actual error

---

## [0.1.0] - 2025-11-24

### Added
- **Core Language Features**:
  - Variables with `let` keyword
  - Integers and floats
  - String literals with interpolation: `"Hello, {name}!"`
  - Boolean values and operators: `&&`, `||`, `!`
  - Arithmetic operators: `+`, `-`, `*`, `/`
  - Comparison operators: `==`, `!=`, `<`, `<=`, `>`, `>=`

- **Control Flow**:
  - Ternary if-then-else: `if (condition) value1 else value2`
  - For loops: `for (item in array) { ... }`
  - Map and filter with for loops

- **Functions**:
  - Function definitions: `fn(x, y) { x + y }`
  - First-class functions (can be stored, passed, returned)
  - Closures with lexical scoping

- **Arrays**:
  - Array literals: `[1, 2, 3]`
  - Indexing: `arr[0]`, negative indexing: `arr[-1]`
  - Slicing: `arr[1:3]`
  - Methods: `.length()`, `.map()`, `.filter()`
  - Multidimensional array support

- **Strings**:
  - String methods: `.upper()`, `.lower()`, `.trim()`, `.split()`, `.replace()`
  - String interpolation
  - Escape sequences: `\{`, `\}`

- **REPL**: Interactive read-eval-print loop

- **File Execution**: Run `.pars` files from command line

- **Comments**: `//` line comments

### Infrastructure
- Built with Go
- Lexer, Parser, AST, and Evaluator architecture
- Command line argument support

---

## Legend

- **Added**: New features
- **Changed**: Changes to existing functionality
- **Deprecated**: Soon-to-be removed features
- **Removed**: Removed features
- **Fixed**: Bug fixes
- **Security**: Security improvements

---

*For version numbers before v0.1.0, see git commit history*
