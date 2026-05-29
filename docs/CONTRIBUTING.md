# Contributing

Short contribution guide and plugin conventions.

Checklist

- Fork + branch: `git checkout -b feature/your-change`
- Format Go code: `gofmt -w ./`
- Run tests: `go test ./...`
- Keep PRs small and documented.

Plugin conventions

- Plugin layout: `plugins/<name>/` with Go wrapper and optional `android/` for Java/Kotlin sources.
- Native Android plugin must implement the `WailsPlugin` interface and return a stable domain string via `getDomain()`.
- Register plugins in `_examples/helloworld/` for manual testing.

Docs

- Keep README concise; place detailed install/cli/contrib docs in `docs/`.
