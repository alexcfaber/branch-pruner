# Contributing

Thank you for considering contributing to branch-pruner! Small, focused contributions are welcome.

Getting started

1. Fork the repository and create a feature branch.
2. Run tests locally before opening a PR:

```bash
go test ./...
```

Linting

We use `golangci-lint` for linting. To run locally:

```bash
golangci-lint run
```

Configuration is in `.golangci.yml`.

Guidelines

- Keep changes small and well-scoped.
- Add unit tests for new behavior and integration tests for end-to-end changes.
- Update `README.md` where relevant.
- Follow Go formatting (`gofmt`) and idiomatic patterns.


CI

This repository has a GitHub Actions workflow that runs `go test ./...` on push and pull requests.

License

By contributing, you agree that your contributions will be licensed under the project's MIT License.
