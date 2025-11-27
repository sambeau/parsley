# Parsley Language Support for Visual Studio Code

Syntax highlighting and language support for the Parsley programming language.

## Features

- **Syntax Highlighting** for:
  - Keywords (`if`, `else`, `return`, `for`, `in`, `let`, `fn`, `as`)
  - Comments (`//`)
  - Strings (double-quoted) with escape sequences and template interpolation `{expr}`
  - Template literals (backticks) with interpolation
  - Regular expression literals (`/pattern/flags`)
  - Path literals (`@/usr/local/bin`, `@./file.pars`)
  - URL literals (`@https://example.com/api`)
  - Numbers (integers and floats)
  - Boolean constants (`true`, `false`, `null`)
  - Special `_` variable (write-only)
  - Built-in functions (`import`, `map`, `sort`, `toString`, `log`, `regex`, `path`, `url`, etc.)
  - HTML/XML tags (singleton and paired)
  - Operators (arithmetic, comparison, logical, assignment, regex match `~` and `!~`)
  - Destructuring syntax

- **Language Features**:
  - Auto-closing pairs for brackets, quotes, and tags
  - Bracket matching for all paired delimiters
  - Comment toggling with line comments (`//`)
  - Code folding with region markers

## Installation

### From Source

1. Copy the `.vscode-extension` directory to your VS Code extensions folder:

   **macOS/Linux:**
   ```bash
   cp -r .vscode-extension ~/.vscode/extensions/parsley-language-0.9.0
   ```

   **Windows:**
   ```powershell
   Copy-Item -Recurse .vscode-extension "$env:USERPROFILE\\.vscode\\extensions\\parsley-language-0.9.0"
   ```

2. Reload VS Code (`Cmd/Ctrl + Shift + P` → "Developer: Reload Window")

### From VSIX (if packaged)

```bash
# Install vsce (VS Code Extension Manager)
npm install -g @vscode/vsce

# Package the extension
cd .vscode-extension
vsce package

# Install
code --install-extension parsley-language-0.1.0.vsix
```

## File Extensions

Files with the `.pars` extension will automatically use Parsley syntax highlighting.

## Example

```parsley
// Function definitions
let greeting = fn(name) {
  `Hello, {name}!`
}

// Array operations
numbers = [1, 2, 3, 4, 5]
doubled = map(numbers, fn(x) { x * 2 })

// Destructuring
head, tail = numbers
x, y, z = 1, 2, 3

// HTML generation
Page = fn({title, contents}) {
  <html>
    <head>
      <title>{title}</title>
      <style>{"
        body { margin: 0; padding: 0; }
      "}</style>
    </head>
    <body>
      {contents}
    </body>
  </html>
}

// Sorting with custom comparator
reverseOrder = fn(a, b) {
  reverse(sort([a, b]))
}

sorted = sortBy([5, 2, 8, 1], reverseOrder)
log(sorted) // [8, 5, 2, 1]
```

## Language Syntax

### Comments
```parsley
// This is a line comment
```

### Variables
```parsley
let x = 42
y = "hello"
_ = "ignored value"  // Write-only variable
```

### Functions
```parsley
square = fn(x) { x * x }
add = fn(a, b) { a + b }
```

### Modules
```parsley
// Import entire module
let utils = import(@./lib/utils.pars)

// Destructure imports
let {add, multiply} = import(@./math.pars)

// Aliasing
let {square as sq} = import(@./math.pars)
```

### Literals
```parsley
// Regular expressions
let emailRegex = /^[\w.+-]+@[\w.-]+\.[a-zA-Z]{2,}$/
let match = "user@example.com" ~ emailRegex

// Paths
let configPath = @./config/settings.json
let binPath = @/usr/local/bin

// URLs
let apiUrl = @https://api.example.com/v1/users
let localUrl = @http://localhost:8080/api
```

### Destructuring
```parsley
// Arrays
x, y, z = 1, 2, 3
head, tail = [1, 2, 3, 4]

// Dictionaries
person = { name: "Sam", age: 57 }
{name, age} = person
```

### Control Flow
```parsley
if (x > 10) {
  "large"
} else {
  "small"
}

for (item in items) {
  log(item)
}
```

### Tags
```parsley
// Singleton tag
<img src="photo.jpg" />

// Paired tags
<div>
  <h1>Title</h1>
  <p>Content</p>
</div>

// With interpolation
<title>{pageTitle}</title>
```

### Strings
```parsley
// Regular strings
name = "Sam"
multiline = "
  This is
  a multiline
  string
"

// Template literals
message = `Hello, {name}!`
html = `<p>{content}</p>`
```

## Contributing

Report issues or contribute at: https://github.com/sambeau/parsley

## License

MIT
