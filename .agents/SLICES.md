# Slice-Based Development Strategy

Use formal slices for non-trivial coordinated work. Slices form a defined,
maintainer-visible delivery loop. Small independent maintenance can use a
ticket or direct `dev` delivery.

This document owns slice definitions, strategy, planning gates, and
traceability. Use `.agents/playbooks/slice-delivery.md` for the execution
procedure. Use `.agents/playbooks/layered-slicing.md` when substantial package
work benefits from review by layer.

## Definitions and Naming

- A delivery set is one approved plan with an ordered local slice map.
- A slice is the smallest meaningful, testable behavior increment or reviewable
  implementation layer.
- Slice labels use `Slice 1`, `Slice 2`, and so on.
- Every delivery set restarts at `Slice 1`.
- Task labels use `T1.1`, `T1.2`, `T2.1`, and so on.
- Completion reporting uses `Slice N of M completed`.
- Reports use
  `ops/<namespace>/report/slices/<delivery-slug>/slice-<n>-<slice-slug>.md`.
- Slice numbers are not repository-global, PR numbers, or continuations of a
  prior delivery set.

Do not use alphabetic or Roman-numeral slice and task labels.

## Planning Gate

Before implementing a formal slice, an approved, committed plan and tracker
must define:

- the delivery-set name and purpose;
- the selected slicing strategy and reason when it is not behavior-first;
- the full ordered slice list;
- the active slice number and short name;
- the exact branch name and expected PR title;
- the expected report path;
- task IDs and validation expectations;
- prerequisites that are already delivered or externally blocked.

If the plan and tracker disagree, correct them before creating a branch or
changing runtime code.

## Strategy Selection

- Prefer behavior-first vertical slices for user-visible or end-to-end
  increments that can remain independently valid.
- Prefer layered slices for substantial package work when contracts, tests,
  mappings, implementation, and integration are easier to review separately.
- State the selected strategy in the plan. Do not select it implicitly or
  switch strategies during implementation without updating and approving the
  plan.

## Slice Rules

Each slice must:

- deliver a meaningful behavior, contract, or reviewable layer;
- be testable and independently reviewable;
- keep the repository in a valid state;
- stay within its recorded task list;
- produce one durable implementation report;
- use one Conventional Commit per tracker task when practical;
- use one branch and one pull request targeting `dev`.

## Execution Sequence

1. Commit the approved plan and tracker to `dev`.
2. Create the recorded branch from current `dev`.
3. Implement tasks in order.
4. Run focused checks, then the broadest relevant available gate.
5. Update the tracker with tasks, commits, and validation.
6. Create or update the slice report and set it to `reviewing`.
7. Open the recorded pull request targeting `dev`.
8. Complete review and wait for the maintainer-controlled merge.
9. Return to `dev`, verify the merge, and close tracker and report state.
10. If the same approved delivery set has a next slice, start it without a new
    confirmation.

## Integration and Traceability

- Keep a one-to-one mapping between tracker tasks and commits when practical.
- Open one pull request per slice. Do not open one per task or commit.
- Keep the repository passing its required checks after each integrated slice.
- Do not integrate slices directly into `main`.
- Align `dev -> main` only through the release-alignment playbook.

## Scope Control

- Do not derive slice numbers from PR numbers or another delivery set.
- Do not add a tuning or integration slice automatically. Record the need and
  update the approved plan first.
