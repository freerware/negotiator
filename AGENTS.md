# Go Coding Style Guide

This document outlines the Go coding conventions and style preferences for this project. Follow these guidelines when writing or modifying Go code.

---

## 1. Package Naming

Package naming is foundational to writing idiomatic Go. These guidelines follow the principles from Effective Go:

- **Use single-word, lowercase names** without underscores. Package names should be simple and clear. Example: `http`, `json`, `template`.
- **Avoid multi-word package names**. If you need multiple words, use the directory structure instead. Example: `encoding/json`, `encoding/base64` — the package name is still `json` and `base64`, not `encoding_json`.
- **Name constructors `New` when there's a single exported type**. If a package exports one primary type, the constructor should be `New()`, not `NewX()`. The package name provides the context. Example: `json.NewDecoder`, not `json.NewJSONDecoder`.
- **Avoid stuttering in exported names**. The package name is already a prefix, so don't repeat it. Example: `http.Server` and `http.Client`, not `http.HTTPServer` and `http.HTTPClient`.
- **Don't worry about name collisions prematurely**. If two packages have the same name, the importer can alias them. Example: `import json "myproject/internal/json"`. Choose the clearest name for your package's purpose.
- **Use nouns for packages**, not verbs. Packages are namespaces for types and functions, not actions. Example: `reader`, `parser`, `handler` — not `read`, `parse`, `handle`.

---

## 2. Project Layout

Follow the standard Go project layout conventions for clear, maintainable structure:

- **Use `internal/` for private packages** that should not be importable by external consumers. Place helper utilities, internal implementations, and application-specific logic here. For server or API applications, place nearly all functionality under this directory — only expose what's truly meant for external use.
- **Use `cmd/` for executable commands**. Each subdirectory represents a distinct command or application entry point. Name subdirectories after the command they produce (e.g., `cmd/server`, `cmd/cli`). Each subdirectory should contain a `main.go` file as the entry point. For server applications, this is where your webserver, CLI tools, and other executables live.
- **Use `pkg/` for public packages** that are safe and intended for external consumption. If you're building a server or API with an accompanying SDK, place the SDK code here. Only include code in `pkg/` that you're committed to supporting as a public API — once external consumers import it, changing or removing it becomes a breaking change.
- **Keep the root directory minimal**. Place module-level documentation (`README.md`, `doc.go`), configuration files, and build artifacts at the root. Avoid placing Go source files directly in the root unless they define the primary package for the module.
- **Organize by functionality, not file type**. Group related types, functions, and interfaces together in the same package. Let the purpose of the code dictate the structure, not a rigid template.
- **Use subpackages for distinct concerns**. When a package grows or has clearly separated responsibilities, extract to subpackages. Example: `pkg/client` and `pkg/server` instead of one bloated `pkg/`.

---

## 3. Documentation

Write clear, comprehensive documentation for all exported identifiers:

- **Document every exported type, function, method, and variable** with godoc-style comments.
- **Use complete sentences** starting with the identifier name. Example: `// New constructs a new negotiator...`
- **Reference relevant standards and specifications** when applicable (RFCs, protocols, etc.). Include links where helpful.
- **Document default behaviors** in package-level comments when providing constructors with defaults.
- **Explain the "why" not just the "what"** for non-obvious implementation decisions.
- **Use package-level `doc.go` files** to provide overview documentation for each package. Include examples.
- **Keep comments up-to-date** when modifying code. Outdated documentation is worse than none.

---

## 4. Testing

Write thorough, maintainable tests that give confidence in the codebase:

- **Use table-driven tests** for testing multiple scenarios. Define test cases with descriptive names.
- **Structure tests as suites** when tests share setup/teardown logic or common helpers. Use test suites.
- **Name test cases descriptively** to communicate intent. Example: `"WithFallback"`, `"EmptyInput"`, `"InvalidConfiguration"`.
- **Test both success and failure paths**. Verify error conditions return expected errors.
- **Use assertion helpers** for cleaner test code. Prefer `Require()` for failures that should stop the test, `Assert()` for failures that can continue.
- **Mirror test file structure to source**. Test files should be named `*_test.go` and live alongside the code they test.
- **Include edge cases** in test coverage: empty inputs, maximum values, nil values, boundary conditions.
- **Favor black-box testing** by using external test packages (`package_test`). Write tests in `mypackage_test"` instead of `mypackage` to ensure you're testing only the exported API, not internal implementation details.
- **Favor generated mocks** over hand-written ones. This reduces boilerplate and keeps mocks in sync with interface changes.
- **Include setup functions in table-driven test cases** when per-test setup is needed. Define a `setup` or `fn` field in your test case struct that applies the necessary setup tasks.

---

## 5. Code Conventions

Follow consistent code style and conventions:

### Variable Declarations

- **Use `var` blocks for grouped declarations** at package level.
- **Use short declarations (`:=`) inside functions** for brevity.
- **Choose descriptive variable names** that communicate intent. Avoid single-letter names except for loops or well-known conventions.

### Error Handling

- **Return errors explicitly**. Don't ignore errors without good reason.
- **Use named return values** when using `defer` for cleanup or recovery.
- **Wrap errors with context** when appropriate, but avoid over-wrapping. Use `fmt.Errorf("%w", err)` with the `%w` verb to wrap errors while preserving the original for `errors.Is()` and `errors.As()` checks.
- **Prefix error messages with the package or component name**. Example: `"domain: insufficient funds"`, `"http: client timeout"`. This makes error sources clear when debugging or logging.
- **Use sentinel errors for known, actionable error conditions**. Define package-level error variables (e.g., `var ErrNotFound = errors.New("domain: not found")`) that callers can check with `errors.Is()`.
- **Implement custom error types** when you need to convey structured information beyond a simple message. Custom error types should implement the `error` interface and can include fields for error codes, operation context, or recovery hints. Use these when callers need to inspect error details programmatically.

### Functions and Methods

- **Keep functions focused** on a single responsibility.
- **Use helper functions inline** when they're small and used only once.
- **Prefer value receivers** unless mutation or avoiding copies is necessary.
- **Favor object-oriented design over functional**. Prefer instance methods on types over bare functions.
- **Embrace encapsulation**. Use unexported fields and types to hide implementation details. Export only what external consumers need. Encapsulation is king.
- **Favor pure functions**. When writing functions or methods, prefer pure functions that have no side effects and return the same output for the same input.
- **Use `*Parameters` structs for functions with more than three arguments**. When a function or method requires more than three parameters, define a `*Parameters` struct to group them. This improves readability and makes adding optional parameters easier. Example: `func (s *Service) Create(ctx context.Context, params *CreateParameters) error`.
- **Never put `context.Context` as a field on Parameters structs**. The context should always be the first explicit parameter to the function or method, not embedded in the Parameters struct. This keeps context handling consistent and idiomatic.

### Interfaces

- **Define interfaces on the receiving side**. Don't create interfaces in the package that implements them. Instead, define interfaces where they're consumed. This follows Go's implicit interface satisfaction and keeps interfaces small and focused on the caller's needs.
- **Prefer single-method interfaces**. Keep interfaces minimal with one method when possible. Example: `io.Reader`, `io.Writer`, `fmt.Stringer`. Compose single-method interfaces into larger interfaces when needed using interface embedding. Example: `type ReadWriter interface { Reader; Writer }`.

### Returning Types

- **Leverage unexported types for controlled instantiation**. When you need control over how types are created or want to enforce invariants, return unexported concrete types that implement exported interfaces. This pattern lets you control instantiation while still allowing callers to use the interface.

### Comments and TODOs

- **Use `// TODO(name)` for pending work**. Include attribution for tracking. Example: `// TODO(john): support encoding extension.`
- **Comment non-obvious logic** but trust the code to speak for itself when clear. Only comment where the "why" is difficult to glean.
- **Reference specifications** in comments when implementing protocol behavior.

### Formatting

- **Run `gofmt`/`goimports`** before committing. No exceptions.
- **Use a consistent line length limit** (e.g., 140 characters) for readability.
- **Order imports**: standard library first, then external packages, then internal/project packages.

### Build and Tooling

- **Use Makefiles with terse target names**. Define simple, memorable target names like `test`, `build`, `lint`, `clean`, `tidy`. The Makefile should hide the actual commands executed — callers shouldn't need to know the underlying tool flags or complexity.
- **Exclude mocks and non-testable packages from test coverage**. Use `go list` to dynamically discover packages and filter out mocks, test helpers, and other non-testable code. Maintain an exclusion list in the Makefile for packages that should not be included in coverage reports.
- **Generate coverage output to `bin/coverage.out`**. Place build artifacts and coverage reports in a `bin/` directory. Example: `go test -race -covermode=atomic -coverprofile=bin/coverage.out $(go list ./... | grep -v -E "$(EXCLUDE)")`.

---

## 6. Preferred Libraries

When external dependencies are necessary, use the following libraries for common use cases:

- **Configuration**: [`viper`](https://github.com/spf13/viper) — for configuration management with support for multiple formats, environment variables, and remote config stores.
- **Command-line interfaces**: [`cobra`](https://github.com/spf13/cobra) — for building powerful CLI applications with subcommands, flags, and help text.
- **Database**: [`sqlx`](https://github.com/jmoiron/sqlx) — a light extension over `database/sql` that provides convenience methods for scanning rows, named parameters, and struct mapping.
- **HTTP**: Use the **standard library** (`net/http`) — it's robust, well-tested, and handles most use cases without additional dependencies.
- **JSON**: Use the **standard library** (`encoding/json`) — it's performant and sufficient for most needs. Avoid third-party JSON libraries unless you have specific requirements (e.g., JSON Schema validation, streaming).
- **Testing**: [`testify/suite`](https://github.com/stretchr/testify) — for organizing tests into suites with shared setup/teardown logic.
- **Mock generation**: [`mockery`](https://github.com/vektra/mockery) — for generating mock implementations of Go interfaces.

**Guiding principle**: Prefer the standard library when it meets your needs. Only reach for external libraries when they provide clear, significant benefits over standard library solutions.

---

## Quick Reference

| Convention        | Style                                      |
| ----------------- | ------------------------------------------ |
| Package naming    | Single word, lowercase, no stuttering      |
| Package structure | `internal/` for private                    |
| Documentation     | Godoc on all exported                      |
| Test structure    | Table-driven + suites                      |
| Error handling    | Explicit returns, named returns with defer |
| Variable names    | Descriptive, intent-revealing              |
| TODO format       | `// TODO(identifier): description`         |
| Formatting        | `gofmt`, `goimports` required              |
| Libraries         | Prefer standard library                    |

---

_Last updated: 2026-06-20_
