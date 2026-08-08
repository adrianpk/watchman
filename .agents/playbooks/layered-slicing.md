# Layered Slicing

Use this playbook for substantial package work when review by implementation
layer is clearer than behavior-first vertical slicing.

## Purpose

Layered slicing keeps the standard slice workflow but separates contracts,
tests, mappings, behavior, and integration into reviewable units. The approved
plan must state the selected strategy and explain why it fits the work.

## Typical Layer Order

Use only the layers that the work needs:

1. Types and package contracts.
2. Function, method, and interface skeletons.
3. Behavior tests that describe the target logic.
4. Mapping and adapter-neutral glue.
5. Logic implementation.
6. Edge-case hardening and broader tests.
7. Package documentation and report closure.
8. Integration adjustments.

## Rules

- Keep the standard plan, tracker, branch, pull request, report, review, merge,
  and `dev` closure workflow.
- Every slice must compile and pass its required checks.
- A preparatory or test-first slice is permitted when it improves review, but
  do not commit intentionally failing tests.
- Use skipped or staged assertions only when the reason is explicit and useful.
- Do not add an integration-adjustment slice automatically. Include it in the
  approved plan or update the plan before starting it.
- Keep shared services and workflows adapter-neutral. Put adapter mappings in
  adapter-owned files or layers.

## Planning Record

Record the strategy in the plan:

```text
Slice strategy: layered
Reason: <why layered review fits this work>
```

For behavior-first slicing of substantial package work, record:

```text
Slice strategy: behavior-first
Reason: <why vertical behavior increments fit this work better>
```
