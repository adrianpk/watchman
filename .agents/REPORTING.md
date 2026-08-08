# Slice Reporting Contract

Every completed formal slice produces one concise, human-readable Markdown
report.

## Purpose

The report provides:

- fast human review;
- long-term traceability;
- architecture and contract context;
- correlation between plans, trackers, commits, and pull requests;
- enough semantic detail to avoid reconstructing intent from a raw diff.

The report does not replace the pull-request diff or duplicate the tracker.

## Scope and Timing

A report is required once per formal implementation slice. Create or update it
while implementation context is current. Set it to `reviewing` before opening
the pull request. Set it to `delivered` after the merge is verified and the
slice is closed on `dev`.

Small ad hoc work committed directly to `dev` does not require a slice report.

## Location

```text
ops/<namespace>/report/slices/<delivery-slug>/slice-<n>-<slice-slug>.md
```

Each slice has exactly one canonical report.

## Required Metadata

```md
# Slice N: <Slice Name>

Status: drafting | reviewing | delivered
Delivery set: <name>
Plan: `<path>`
Tracker: `<path>`
Branch: `<branch>`
PR: pending | `#<number>`
```

## Required Sections

```md
## Purpose

## Delivered Behavior

## Implementation Notes

## Contracts Added or Changed

## Files of Interest

## Validation

## Risks and Follow-ups
```

## Content Rules

- Keep reports in English and concise.
- Record exact validation commands and pass or fail outcomes.
- Distinguish focused, full-suite, local, emulator, staging, and production
  results when applicable.
- Link the plan and tracker. Do not duplicate the tracker task table.
- Mention only high-signal files, decisions, and contracts.
- Do not paste raw diffs.
- Use `None recorded.` when the risk or follow-up section is empty.
- In the pull-request body, link the report under `Changes`; do not paste it.
