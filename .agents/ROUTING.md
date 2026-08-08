# Work Request Routing

This document defines request classification and execution authorization.

## Execution Modes

- exploration
- planning
- implementation
- investigation
- maintenance or meta-work

## Classification

- Exploration covers questions, possibilities, design discussion, and tradeoff
  analysis. It does not create persistent artifacts or delivery state.
- Planning occurs only when requested. Persistent plans, trackers, tickets, and
  specifications require explicit authorization.
- Implementation starts only with explicit instruction and stays within the
  approved scope.
- Investigation verifies and reports state without changing it unless the user
  also authorizes a fix.
- Maintenance or meta-work changes repository workflow, tooling, documentation,
  or operating rules within the explicitly requested scope.
- Use tickets for bounded independent work. Use specifications, plans,
  trackers, slices, and reports for non-trivial coordinated work.

## Continuation

- A clear request to start, continue, resume, or proceed authorizes execution
  when the next workflow unit is already defined.
- A merge confirmation or status acknowledgment continues only the active
  approved delivery set. It does not authorize a new delivery set, plan,
  branch, or unrelated work.
- Do not turn exploration, planning, prior work, or status reporting into
  open-ended autonomous execution.

## Ambiguity

- Ask for clarification only when a reasonable assumption could materially
  change the result, scope, or external effect.
- During exploration, work inside the direction under evaluation. Name concrete
  tradeoffs and push back only for a specific technical, safety, scope, or
  product risk.

After classification, use `.agents/LIFECYCLE.md` for state transitions and the
applicable concern documents for execution.

## Non-Goal

This document does not define product behavior, Git policy, or delivery
procedures.
