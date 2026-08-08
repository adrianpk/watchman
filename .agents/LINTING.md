# Go Linting and Quality Gates

This repository uses focused quality gates for Go application code.

## Baseline checks

- Formatting:
  - `go fmt ./...`
- Test suite:
  - `go test ./...`
- Race-sensitive validation when concurrency behavior changes:
  - `go test -race ./...`

## Optional stricter checks

If enabled by project setup:
- `golangci-lint run ./...` for static analysis.
- `golangci-lint run --timeout=5m ./...` for CI or stricter pre-merge validation.
- `go vet ./...` when it is not already covered by the configured linter.
- `go test -cover ./...` or the project's configured coverage command.

## Recommended golangci-lint shape

Prefer a focused baseline before adding noisy style-only rules:

- `govet`
- `staticcheck`
- `ineffassign`
- `errcheck`
- `unused`
- `misspell`
- `errorlint`
- `bodyclose` for HTTP-heavy projects
- `contextcheck` for context-sensitive code

Formatters should include:

- `gofmt`
- `goimports` or `gci`

## Hook policy

If git hooks are configured, pre-commit should run fast checks only:
- formatting checks
- changed-package tests where practical
- focused lint checks for changed code

Full lint and test gates should run before merge.

## Rules

- Do not skip tests for behavior changes.
- Fix failing checks before declaring work complete.
- For docs-only changes, validate the diff and whitespace instead of running unrelated application tests.
- Keep linter strictness high enough to catch defects, but avoid rules that encourage mechanical churn without improving maintainability.

## Root-level gates

- Repository-level quality tasks should run from the repository root unless the project documents a different working directory.
- Prefer project-owned commands such as `make test`, `make lint`, `make verify`, or CI-equivalent scripts when present.
- Report the exact command and outcome when closing work.
