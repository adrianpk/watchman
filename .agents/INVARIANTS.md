# Collaboration Invariants

This document defines non-negotiable scope, approval, safety, and truthfulness
rules.

## Scope Integrity

- Execute only approved scope.
- Record out-of-scope ideas instead of implementing them.
- Do not infer authorization from task size, technical value, or prior work.

## Approval Gates

- High-impact changes require explicit approval.
- Production deployment and external billing, credential, access, or account
  mutations require explicit approval.
- Start, stop, or restart long-running development services only when the user
  requests it or authorized validation requires it.

## Safe Delivery

- Prefer incremental, reversible changes.
- Keep the repository in a working state between delivery units.
- Do not use destructive Git operations without explicit approval.
- Keep secrets and local agent configuration out of version control.
- Preserve unrelated user changes and report any overlap with requested work.

## Quality and Truthfulness

- Do not mark work complete with failing required checks.
- Do not describe focused checks as the full validation suite.
- Distinguish local implementation from deployed or externally verified state.
- Record checks that could not run and the reason.
- Keep repository, tracker, commit, and pull-request state consistent.
