<!-- Use this file to provide workspace-specific custom instructions to Copilot. For more details, visit https://code.visualstudio.com/docs/copilot/copilot-customization#_use-a-githubcopilotinstructionsmd-file -->

# Go Project Instructions

This is a Go project using standard Go practices and conventions.

## Code Style and Conventions

- Follow standard Go formatting with `gofmt`
- Use meaningful package names that are lowercase and concise
- Follow Go naming conventions (camelCase for private, PascalCase for public)
- Use interfaces to define behavior contracts
- Handle errors explicitly and appropriately
- Write unit tests for all public functions
- Use meaningful variable and function names
- Use Go's std library functions, interfaces and types where possible

## Project Structure

- `main.go` - Application entry point
- `pkg/` - Public packages for external use
  - `lexer/` - Tokenizes input into lexical tokens
  - `parser/` - Converts tokens into an Abstract Syntax Tree
  - `ast/` - Defines the Abstract Syntax Tree nodes
  - `evaluator/` - Evaluates the AST and executes the program
  - `repl/` - Read-Eval-Print Loop for interactive usage

## Dependencies

- Use `go mod` for dependency management
- Keep dependencies minimal and well-maintained
- Pin dependency versions for reproducible builds

## Testing

- Write table-driven tests when appropriate
- Use the standard `testing` package
- Aim for high test coverage on critical paths
- Use `testify` for more complex assertions if needed

## Version Numbering
- Use a `VERSION` file at the root of the repository to track the current version
- Keep the version number in in the README.md file up to date
- Increment minor version when major changes are made to the documentation
- Follow Semantic Versioning (SemVer) principles
- Increment major version for breaking changes
- Increment minor version for new features in a backward-compatible manner
- Increment patch version for backward-compatible bug fixes