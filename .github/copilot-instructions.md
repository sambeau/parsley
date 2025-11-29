<!-- Use this file to provide workspace-specific custom instructions to Copilot. For more details, visit https://code.visualstudio.com/docs/copilot/copilot-customization#_use-a-githubcopilotinstructionsmd-file -->

# Go Project Instructions

This is a Go project using standard Go practices and conventions.

## Code Style and Conventions

- Follow standard Go formatting with `gofmt`
- Write 'idiomatic' Go code that is easy to read and understand by other Go developers
- Follow Go naming conventions (camelCase for private, PascalCase for public)
- Use meaningful variable and function names
- Use meaningful package names that are lowercase and concise
- Keep functions small and focused on a single task
- Use comments to explain the "why" behind complex logic, not the "what"
- Use interfaces to define behavior contracts
- Handle errors explicitly and appropriately
- Write unit tests for all public functions
- Use Go's std library functions, interfaces and types in preference to creating your own

## Project Structure

- `main.go` - Application entry point
- `pkg/` - Public packages for external use
  - `lexer/` - Tokenizes input into lexical tokens
  - `parser/` - Converts tokens into an Abstract Syntax Tree
  - `ast/` - Defines the Abstract Syntax Tree nodes
  - `evaluator/` - Evaluates the AST and executes the program
  - `repl/` - Read-Eval-Print Loop for interactive usage
- `docs/` - Documentation and design documents
  - `examples/` - Example scripts and usage
  - `design/` - Language design documents
- `tests/` - Unit and integration tests

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

## CHANGELOG.md

- Maintain a `CHANGELOG.md` file at the root of the repository
- Document all notable changes in the `CHANGELOG.md`
- Follow the "Keep a Changelog" format for consistency
- Link changelog entries to corresponding version numbers

## Docs & README

- Maintain a `README.md` file at the root of the repository
- Keep the README with quick summary of language features, with short examples where needed
- Keep the README up to date with installation, usage, and contribution instructions
- Maintain a `docs/reference.md` file for language reference
- Keep the reference documentation up to date with comprehensive list of language features
- Maintain a `docs/` directory for design documents and specifications
- Keep design documents up to date with the current implementation status
- Document any deviations from the original design in the docs
- Strikethrough completed items in design/TODO.md with version numbers where applicable

## VS Code Extension

- Maintain a `.vscode-extension/README.md` file for the VS Code extension
- Keep the extension version aligned with the main project version
- Document installation and usage instructions in the extension README

## On feature/bugfix completion

- use `docs/Pre-flight for Git Commit.md` as a checklist to prepare for git commit

## Planning and design documents

- Store planning and design documents in `docs/design`
- Keep design documents up to date with implementation status
- Document any deviations from the original design in the docs
- Use relevant design documents to guide implementation and testing

## Design Philosophy

- Adhere to Parsley's core design philosophy when planning to add new features or make changes. (see docs/design/Design Philosophy.md)

