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

## Project Structure

- `main.go` - Application entry point
- `pkg/` - Public packages for external use
- `internal/` - Private packages for internal use
- `cmd/` - Command-line applications
- `api/` - API definitions and protocol files
- `configs/` - Configuration files

## Dependencies

- Use `go mod` for dependency management
- Keep dependencies minimal and well-maintained
- Pin dependency versions for reproducible builds

## Testing

- Write table-driven tests when appropriate
- Use the standard `testing` package
- Aim for high test coverage on critical paths
- Use `testify` for more complex assertions if needed
