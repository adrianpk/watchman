# Agent Operating Contract

Read this document before starting work in this workspace.

This file routes agent work to the normative documents under `.agents/`. It is
not a product specification, implementation guide, or duplicate policy source.

## Required Loading

- Before changing workspace or external state, read `.agents/INVARIANTS.md`
  and `.agents/ROUTING.md`.
- Before performing an action, read the document that owns that concern.
- Use `.agents/PLAYBOOKS.md` to select an applicable procedure.

## Concern Routing

- User collaboration and responses: `.agents/COMMUNICATION.md`
- Repository language and locale boundaries: `.agents/LANGUAGE.md`
- Request classification and authorization: `.agents/ROUTING.md`
- Work-state transitions and completion: `.agents/LIFECYCLE.md`
- Safety and scope boundaries: `.agents/INVARIANTS.md`
- Branches, commits, remotes, pushes, and pull requests:
  `.agents/VERSIONING.md`
- Documentation and operational artifacts: `.agents/DOCUMENTING.md`
- Architecture and ownership: `.agents/ARCHITECTURE.md`
- Code and stack conventions: `.agents/CONVENTIONS.md`
- Validation and quality gates: `.agents/LINTING.md`
- Formal delivery sets and slices: `.agents/SLICES.md`
- Persistent implementation reports: `.agents/REPORTING.md`
- Operational procedures: `.agents/PLAYBOOKS.md`
- Repository activity recording: `.agents/ACTIVITY.md`

## Ownership

Each rule has one normative owner. Other documents may reference that owner but
must not restate or paraphrase the rule.

Keep product behavior and project-specific technical decisions outside this
file.
