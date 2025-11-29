## Create Parsley 'Cheat Sheet'

**tl;dr**: Write yourself a small help file / primer to quickly bring yourself up-to-date about Parsley's grammar and quirks—especially compared to Javascript and Python (and to a lesser degree to Rust and Go).  I will also find it useful.

- This document will be saved and refered to when writing tests and examples script during development. It ewill also be a handy guide when working on Plans and design documents.

- The two main tripping points you run up against are that Parsley uses 'log()' instead of 'print'. When debugging a multi-line script it is recommended to use logLine() as it prints out the line number.

- There are also major dirrerences with how 'for' and 'if' work, as Parslef is concatentative - you can assign to for, as-if it's ?: from other languages, and 'for' returns an array of values, so it is more akin to 'map' in other languages.

The Most Used Parsley Features (Tests & Examples Analysis) are:-

## Top Functions & Keywords

### 1. Core Output/Logging
- `log()` - 523 uses in examples, 710 in tests - **Most used function**
- `logLine()` - 259 uses

### 2. Control Flow
- `let` - 382 uses - **Primary variable declaration**
- `if` - 72 uses
- `for` - 44 uses

### 3. File I/O Factories
- `file()` - 46 uses (tests) + 8 (examples)
- `JSON()` - 20 uses
- `dir()` - 23 uses (tests) + 15 (examples)
- `text()` - 27 uses (tests) + 7 (examples)
- `SFTP()` - 6 uses

### 4. String & Collection Functions
- `len()` - 189 uses (tests) + 14 (examples)
- `split()` - 23 uses (tests) + 6 (examples)
- `sort()` - 7 uses
- `map()` - 22 uses (tests) + 1 (example)

### 5. DateTime
- `time()` - 159 uses (tests) + 22 (examples)
- `now()` - 28 uses (tests) + 8 (examples)

### 6. Operators
- File operators (`==>`, `<==`, `==>>`) - 50 uses in examples
- Network operators (`=/=>`, `<=/=`, `=/=>>`) for SFTP/HTTP

### 7. HTML/XML Tags
- `<p>` - 10 uses
- `<div>` - 6 uses
- Custom elements supported

### 8. Path Literals
- URL literals (`@https://...`) - Heavy usage for API testing
- File path literals (`@/path/to/file`)

### 9. Other Common Functions
- `replace()` - 4 uses
- `filter()` - 2 uses
- `import()` - 5 uses (module system)

## Summary
The data shows Parsley is heavily used for **logging/output**, **file I/O**, **string/collection manipulation**, and **HTML templating** in the examples.