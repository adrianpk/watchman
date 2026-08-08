# Work Modes

Use this playbook whenever work may move from conversation into delivery.

## Explore

Discuss possibilities, tradeoffs, examples, and design direction. Do not create
files, plans, tickets, branches, scripts, commits, pull requests, or other
delivery state.

## Plan

Plan only when the maintainer requests planning. Define scope, exclusions,
validation, delivery mode, and whether the work needs slices. Create persistent
artifacts only when authorized.

Each delivery set begins at `Slice 1`. Exploration, a plan, a merge
confirmation, or a status acknowledgment does not expand the delivery set.

## Deliver

Start implementation only after explicit authorization of the agreed scope. An
approved delivery set continues through its planned slices after each confirmed
merge. A confirmation or `OK` continues only that set. It does not start another
set or authorize adjacent work.

## Close

Prepare the authorized handoff or pull request after validation. The maintainer
controls merges. After a reported merge, verify it, synchronize `dev`, and close
the completed work. Continue only with the next planned slice in the same
approved delivery set.

## Technical Pauses

Technical status reports are always allowed. State what changed, the active
scope, delivery-set status, `Slice N of M`, validation, and blockers. Reporting
status does not request permission and does not expand the authorized scope.
