# Testing Strategy

Use this playbook when implementing or reviewing behavior changes.

## Baseline

- Use the repository-native test framework and commands defined in
  `.agents/LINTING.md`.
- Use table-driven or parameterized tests for deterministic rules, state
  transitions, policy decisions, mappings, parsers, and calculations.
- Cover expected failures, edge cases, boundary contracts, and input
  immutability when relevant.
- For bug fixes, reproduce the defect with a failing test first when practical.
- Do not skip tests for behavior changes.

## Risk-Based Techniques

These techniques are tools, not mandatory ceremony:

- property-based tests for broad invariants over generated inputs;
- fuzz tests for parsers, decoders, tokenizers, and untrusted input;
- mutation testing when branch confidence is weak and tooling exists;
- complexity or CRAP-style review for branching and high-impact decisions;
- concurrency or race tests for shared mutable state and transactions;
- benchmarks for code on a measured performance path;
- focused integration tests for persistence, transport, provider, browser, or
  external-service boundaries;
- test-mode or read-only smoke checks for external integrations.

Select techniques from the changed behavior, risk, and available tooling. Do
not add unrelated test ceremony.

## Validation Order

1. Run the smallest focused checks that give fast feedback.
2. Run the broadest relevant available repository gate.
3. Run external or deployment checks only with the required authority and
   safe credentials or test mode.

## Reporting

- Record exact commands and outcomes in the tracker and slice report.
- Distinguish focused checks from the full suite.
- Distinguish local, emulator, staging, and production verification.
- Record omitted high-value techniques only when the decision affects review
  risk.
