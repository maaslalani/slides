# Repository Guidelines

## Project Structure & Module Organization

The CLI entry point lives in `main.go`, backed by Go modules declared in `go.mod`. The `internal/` tree houses production packages: `internal/code` renders and executes code blocks, `internal/navigation` handles slide traversal, `internal/process` wraps pre-processing, and `internal/file` plus `internal/meta` provide I/O helpers. Presentation themes, fonts, and style tests reside in `styles/`. Sample decks and walkthroughs live in `examples/` and `docs/`, while reusable assets (GIFs, screenshots) sit under `assets/`. Keep new features inside the appropriate `internal/` subpackage to preserve clear boundaries.

## Build, Test, and Development Commands

Use `make` to launch `go run main.go examples/slides.md` for a quick manual check. Compile binaries with `make build`, which writes the `slides` executable to the repository root. Run the focused test suite via `make test` (alias for `go test ./... -short`). For experiments, `go run ./cmd/...` is unnecessary—stick to the provided entry point to avoid drift.

## Coding Style & Naming Conventions

Format every Go file with `gofmt` (tabs, no alignment hacks) and organize imports with `goimports`. Follow standard Go casing: exported identifiers in `CamelCase`, locals in `mixedCase`, constants in `UPPER_SNAKE` when multiple words. Group files by responsibility inside their package; avoid duplicate helper functions across subpackages—share via `internal/...` instead. Markdown examples should mirror the decks in `examples/` and keep filenames kebab-case (e.g., `live-demo.md`).

## Testing Guidelines

All tests use Go’s built-in `testing` package, naming files `*_test.go` and functions `TestXxx`. Before submitting changes, run `go test ./...` and, when touching rendering logic, add assertions in the relevant package (e.g., `internal/navigation/navigation_test.go`). Use table-driven tests to cover edge cases yet stay concise. Optional coverage checks (`go test ./... -cover`) help catch regressions in parsing and execution flows.

## Commit & Pull Request Guidelines

Prefer short, descriptive commit subjects in the imperative mood, mirroring the existing history (`fix: adjust gif path`, `build(deps): bump lipgloss`). Squash noisy work-in-progress commits locally. Pull requests should summarize user-visible changes, link related issues, and note any UI or output differences (GIFs or short recordings help). Confirm `go test ./...` passes and mention any follow-up tasks explicitly.
