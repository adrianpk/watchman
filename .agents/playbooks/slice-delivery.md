# Slice Delivery

Use this playbook only for formal delivery work. Exploration and design
discussion use the work-modes playbook and do not create delivery artifacts.

## Creation Gate

A formal slice is executable only when an approved, committed plan and tracker
contain the complete delivery-set map.

The map must include:

- the delivery-set name;
- the selected slicing strategy;
- ordered slices beginning at `Slice 1`;
- the active slice and task IDs;
- the exact branch and expected pull-request title;
- the expected report path;
- validation expectations and external prerequisites.

If the map is missing or the plan and tracker disagree, correct them before
creating a branch or changing runtime code.

## Naming and Traceability

- Use `Slice 1`, `Slice 2`, and so on within one delivery set.
- Restart at `Slice 1` for each delivery set.
- Use `T1.1`, `T1.2`, `T2.1`, and so on for task IDs.
- Do not derive slice numbers from pull-request history or another plan.
- Keep the plan, tracker, branch, expected pull-request title, report path, task
  IDs, commits, and validation results correlated.

## Workflow

1. Verify current `dev`, `origin/dev`, working tree, plan, tracker, and relevant
   pull-request state.
2. Create the recorded branch from current `dev`.
3. Implement only the active slice tasks in order.
4. Use one Conventional Commit per tracker task when practical.
5. Run focused checks, then the broadest relevant available gate.
6. Update tracker tasks, commits, and validation results.
7. Create the report required by `.agents/playbooks/slice-reporting.md`.
8. Set the report to `reviewing` and open one pull request targeting `dev`.
9. Use the exact title recorded in the plan and tracker.
10. Complete review and wait for the maintainer-controlled merge.
11. After merge confirmation, return to `dev`, verify the merge, and close the
    report and tracker.
12. If the same approved delivery set has a next slice, activate and execute it
    without requesting another confirmation.

## Pull Request Body

```md
## Summary

- <functional outcome>

## Changes

- <high-signal path and change>
- <slice report path>

## Validation

- <command> (pass|fail)
```

Order validation from focused to broad. Do not describe a focused subset as the
full suite.

## Blockers

- Capture the failing command and exact error.
- Attempt the least-invasive in-scope alternative.
- Record unresolved blockers in the tracker and report.
- Stop before risky or out-of-scope changes.

## Completion Checklist

- Scope is implemented.
- Relevant checks ran or their blockers are recorded.
- Tracker tasks and commits are current.
- The slice report is current.
- Pull-request state and final handoff are explicit.

## Main Alignment

- Slice pull requests target `dev`, not `main`.
- Never work directly on `main`.
- Align `dev -> main` only through the release-alignment playbook after an
  explicit request.

## Technical Pauses

Technical status reports are allowed during an authorized slice. State the
completed work, active delivery set, `Slice N of M`, validation, and blockers.
A status report does not request permission and does not expand scope.
