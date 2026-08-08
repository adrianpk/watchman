# Slice Reporting

Use this playbook for every formal implementation slice.

## Timing

1. Create or update the report while implementation context is current.
2. Set its status to `reviewing` before opening the pull request.
3. Include the report path in the pull-request `Changes` section.
4. After merge, set its status to `delivered`, record the pull-request number,
   and close the tracker state on `dev`.

## Location

```text
ops/<namespace>/report/slices/<delivery-slug>/slice-<n>-<slice-slug>.md
```

## Template

```md
# Slice N: <Slice Name>

Status: drafting
Delivery set: <name>
Plan: `<path>`
Tracker: `<path>`
Branch: `<branch>`
PR: pending

## Purpose

## Delivered Behavior

## Implementation Notes

## Contracts Added or Changed

## Files of Interest

## Validation

## Risks and Follow-ups
```

## Rules

- Follow `.agents/REPORTING.md`.
- Keep the report in English and concise.
- Record commands exactly as run and distinguish focused, full, local,
  emulator, staging, and production outcomes.
- Link the plan and tracker. Do not copy the tracker table.
- Mention high-signal files only. Do not paste raw diffs.
- Use `None recorded.` when no risk or follow-up exists.
- Markdown is canonical. HTML can only be a derived rendering.
