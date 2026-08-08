# Documentation Guidelines

This repository uses Diataxis for user-facing documentation under `docs/`:
https://diataxis.fr

## Ownership Boundaries

- Keep product and runtime documentation under `docs/` or the approved
  `ops/<namespace>/spec/` location.
- Keep agent operating rules in `AGENTS.md` and `.agents/`.
- Keep delivery state and operational artifacts under `ops/<namespace>/`.
- Resolve the active namespace from explicit user input, then `OPS_USER`, then
  `default`.
- Follow `.agents/LANGUAGE.md` for artifact language and locale rules.

## Categories

- `docs/tutorials/` for learning paths
- `docs/how-to/` for task procedures
- `docs/reference/` for contracts and exact behavior
- `docs/explanation/` for rationale and tradeoffs

Keep the top level of `docs/` limited to `index.md` and the four Diataxis
directories when the repository uses a documentation index. Put additional
structure inside the applicable quadrant.

## Examples and references

- Keep setup and execution examples aligned with the stack, package manager, runtime, and entrypoints actually present in the project.
- API and reference docs should use the naming, signatures, paths, and conventions used by the codebase.
- If API docs are generated, keep source comments compatible with the project's documentation tool.

## Contract documentation

- Document stable wire formats, configuration fields, persistence shapes,
  public APIs, message payloads, and error contracts under `docs/reference/`.
- Document architecture decisions and tradeoffs under `docs/explanation/`.
- Keep examples aligned with actual project commands, paths, and symbols.
- Document exported or public contracts with the repository-native source
  documentation convention when that information helps callers.

## Operational documents

- `ops/<namespace>/spec/drafts/` stores evolving specifications.
- `ops/<namespace>/spec/stable/` stores accepted specifications.
- `ops/<namespace>/quiz/` stores raw and curated quiz artifacts.
- `ops/<namespace>/dump/` stores raw and curated knowledge dumps.
- `ops/<namespace>/plan/` stores cycle plans.
- `ops/<namespace>/tracker/` stores execution trackers.
- `ops/<namespace>/ticket/` stores bounded repository work.
- `ops/<namespace>/report/` stores durable cycle outputs, generated maps, and reports.
- Do not mix operational delivery records into `docs/` unless they become durable product documentation.
- Keep temporary command output under an ignored temporary directory, not in
  versioned operational artifacts.

## General rules

- Each document should have one primary intent.
- Split documents when intent changes.
- Preserve technical meaning when reorganizing docs.
- Avoid title/heading suffixes in parentheses (for example `(Draft)`).
- Prefer semantic metadata fields (for example `Status`) over title decoration when state needs to be explicit.
- For structural documentation updates, complete one consistency pass before commit.
- Prefer one cohesive commit per documentation iteration; avoid corrective micro-commits caused by partial updates.

## Scope boundary

Diataxis applies to `docs/*` only.
Operational instructions in `.agents/*` and delivery records under `ops/*` follow their own templates.
