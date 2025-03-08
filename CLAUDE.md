# Go Development Quick Reference

## Commands
- `task` - Run default workflow (clean, format, lint)
- `task all` - Run all tasks: clean, info, format, lint, build, test
- `task test` - Run all tests
- `go test ./path/to/package -run TestName` - Run a single test
- `task lint` - Run linters
- `task format` - Format code
- `task build` - Build project
- `task cover` - Run tests with coverage

## Code Style
- **Formatting**: 120 char line limit, using `gofumpt` with extra rules
- **Imports**: Strict ordering: standard lib → golang.org → 3rd party → github.com/idelchi → default
- **Naming**: PascalCase for exported, camelCase for unexported
- **Error handling**: Always check errors, propagate them up, avoid log.Fatal/os.Exit
- **Tests**: Files end with `_test.go`, use package `packagename_test`
- **Documentation**: All exported items must have proper Go doc comments
- **Linting**: Uses golangci-lint with most linters enabled, disable with `//nolint:lintername`

## Project Structure
- `internal/`: Private application code
- `pkg/`: Public reusable packages