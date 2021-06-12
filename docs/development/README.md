# Development

Make changes, and test them by running:
```
make
```

This will run `go run main.go examples/slides.md`, you can then ensure
everything still works.

If you're adding a feature that requires a specific piece of markdown, you can
add a file with your test case into `examples/<test>.md` and iterate on that file.

### Breaking Changes
Most changes should be entirely backwards compatible.
Ensure that `slides examples/slides.md` still works.

### Codebase
Initialization (command-line interface, defaults) happens in [`cmd/root.go`](../../cmd/root.go).
Interaction (controls, input, output) happens in [`model.go`](../../internal/model/model.go)
Optional configuration (e.g. `theme: dark`) can be added to [`meta.go`](../../internal/meta/meta.go)
