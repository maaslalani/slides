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
