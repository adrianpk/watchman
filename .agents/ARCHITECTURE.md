# Go Architecture

This document describes Go repository structure and ownership boundaries.

## Objective

Keep the Go codebase explicit, package-oriented, and easy to reason about without hiding workflows in oversized packages or command entrypoints.

## Baseline structure

```text
cmd/          # executable entrypoints
internal/     # application code private to this module
pkg/          # public/shared contracts only when external reuse is intentional
assets/       # embedded non-code assets, schemas, templates, and manifests
docs/         # user-facing and technical documentation
ops/          # specs, plans, tickets, quiz notes, dumps, and delivery artifacts
testdata/     # test fixtures owned by package tests
.agents/      # collaboration and execution rules for this repository
AGENTS.md     # top-level agent operating contract
```

## Boundary rules

- Command packages under `cmd/` should assemble dependencies, parse process-level configuration, start lifecycle components, and stay thin.
- Application behavior should live in `internal/` by default.
- Use `pkg/` only for contracts, envelope types, or libraries that are intentionally shared outside the module.
- Keep package boundaries aligned to behavior and ownership, not generic utility buckets.
- Isolate IO adapters, persistence, transport, and external service integration from domain decisions.
- Long-running components should expose explicit lifecycle methods when startup, readiness, or shutdown behavior matters.
- Embedded assets should be grouped under `assets/` or a clearly owned package-local `testdata/` directory.
- Tests should live beside the package behavior they validate.
- Delivery planning, specs, tickets, quiz notes, dumps, and cycle records should live in `ops/`.

## Runtime assembly

- Runtime startup should follow a deterministic assembly order.
- Constructors should store dependencies and lightweight configuration only.
- Failable initialization belongs in explicit setup/start paths, not constructors.
- Readiness and health checks should be visible at package boundaries when the component is operationally significant.
- Shutdown should be explicit and context-aware.

## Change discipline

- Structural changes must be explicit in commit scope.
- Do not introduce additional top-level directories or shared packages without clear ownership need.
- Prefer extending existing Go package boundaries before creating new framework layers.
