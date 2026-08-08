# Ticket Workflow

Use tickets for concrete development work that does not need a full spec, plan, or tracker.

This playbook defines a repository-local Markdown ticket model under `ops/<namespace>/ticket/`.

## Scope

- `ops/<namespace>/ticket/` contains repository-development tickets.
- It does not store application runtime tickets unless the project explicitly chooses to share the same representation.
- Ticket capture must preserve useful work without requiring full classification upfront.

## Directory Layout

`ops/<namespace>/ticket/` represents ticket columns with status directories:

```text
ops/<namespace>/ticket/
  TEMPLATE.md
  open/
  ready/
  in_progress/
  reviewing/
  blocked/
  solved/
  closed/
```

The directory is the operational column. The `status` frontmatter must match the directory name and acts as a consistency check.

## File Rules

1. Create ticket files under `ops/<namespace>/ticket/open/`.
2. Use sortable filename format: `YYYYMMDDHHMMSS-<slug>.md` with a UTC timestamp.
3. Use frontmatter fields from `ops/<namespace>/ticket/TEMPLATE.md` when that template exists.
4. Keep code and repository docs in English.
5. Move ticket files between status directories when status changes.
6. Keep the filename stable after creation unless correcting a timestamp/id mismatch.
7. If scope expands beyond a bounded ticket, promote it to the formal spec/plan/tracker flow.

## Delivery Unit

A ticket is expected to be one focused delivery unit.

- Default: one ticket -> one focused delivery unit -> `solved` or `closed`.
- Small ticket work may be committed directly to `dev`.
- Use one branch and one pull request when the user requests review delivery or
  when the ticket explicitly records that mode.
- If a ticket cannot be delivered in one focused branch/PR, split it into related tickets or promote it to the formal spec/plan/tracker workflow.
- Use slicing for long specs, plans, or coordinated series of tickets, not as the default shape for a single ticket.
- A ticket series may act as a lightweight roadmap, but each ticket should remain independently reviewable.

## Helper Scripts

If the repository provides ticket helper scripts, use them for repeatable operations.

Common helper responsibilities:

- create an `open` ticket with canonical id, filename, frontmatter, and body scaffold
- list or search tickets with filters such as status, tag, kind, and output format
- validate ticket frontmatter, status directories, ids, filenames, enum values, and tag limits

Run the repository's ticket validation command before committing ticket workflow changes.

## Required Properties

- `id`: stable ticket id, `TKT-YYYYMMDDHHMMSS`.
- `title`: short concrete title.
- `status`: lifecycle state.
- `kind`: work category.
- `severity`: impact if unresolved.
- `priority`: execution urgency.
- `scope`: affected area.
- `tags`: comma-separated labels for grouping, retrieval, or later triage.
- `source`: where the ticket came from.
- `reported_at`: UTC timestamp for initial capture.
- `commits`: list of related commit hashes.

## Status Values

- `open`: captured, but not yet sufficiently refined for execution.
- `ready`: clear enough to execute when someone picks it up.
- `in_progress`: actively being worked.
- `reviewing`: implementation or proposal is ready for review, usually with a PR.
- `blocked`: cannot advance because a dependency, decision, or external condition is missing.
- `solved`: resolved positively, normally by merged implementation.
- `closed`: closed without implemented resolution.

## Status Transitions

- `ops/<namespace>/ticket/open/ -> ops/<namespace>/ticket/ready/`
- `ops/<namespace>/ticket/ready/ -> ops/<namespace>/ticket/in_progress/`
- `ops/<namespace>/ticket/in_progress/ -> ops/<namespace>/ticket/reviewing/`
- `ops/<namespace>/ticket/reviewing/ -> ops/<namespace>/ticket/solved/`
- `ops/<namespace>/ticket/reviewing/ -> ops/<namespace>/ticket/closed/`
- `ops/<namespace>/ticket/open/ -> ops/<namespace>/ticket/closed/`
- `ops/<namespace>/ticket/ready/ -> ops/<namespace>/ticket/closed/`
- `any active directory -> ops/<namespace>/ticket/blocked/`
- `ops/<namespace>/ticket/blocked/ -> ops/<namespace>/ticket/open/`
- `ops/<namespace>/ticket/blocked/ -> ops/<namespace>/ticket/ready/`
- `ops/<namespace>/ticket/blocked/ -> ops/<namespace>/ticket/in_progress/`

When moving a ticket, update the `status` frontmatter in the same change.

Use `solved` when work was done and the issue is resolved. Use `closed` when the ticket is discarded, duplicated, obsolete, moved elsewhere, or answered without implementation.

## Kind Values

- `unclassified`: not classified yet.
- `bug`: incorrect behavior.
- `task`: concrete implementation or maintenance work.
- `cleanup`: simplification, pruning, or hygiene.
- `follow_up`: lightweight item captured from chat, review, or implementation.
- `research`: investigation before execution.
- `decision`: explicit choice needed before execution.

## Severity Values

- `unclassified`: impact cannot be inferred yet.
- `low`: minor inconvenience or cleanup.
- `medium`: noticeable workflow or correctness issue.
- `high`: significant functional, safety, or delivery impact.
- `critical`: blocks essential work or risks data/state integrity.

## Priority Values

- `unclassified`: scheduling importance cannot be inferred yet.
- `low`: can wait.
- `normal`: default priority.
- `high`: should be handled soon.
- `urgent`: needs immediate attention.

## Scope Values

- `unclassified`: affected area cannot be inferred yet.
- `domain`: core domain behavior and business rules.
- `api`: HTTP, RPC, CLI, or external interface behavior.
- `persistence`: storage, migrations, repositories, schemas, or data access.
- `ui`: user interface and browser interactions.
- `ops`: repository workflow, playbooks, tickets, specs, plans, or trackers.
- `docs`: durable documentation.
- `infra`: scripts, CI, environment, dependencies, or local tooling.

## Resolution Values

Set `resolution` only when `status` is `solved` or `closed`.

- `fixed`: implemented and resolved.
- `wont_do`: intentionally not doing it.
- `duplicate`: represented by another ticket or artifact.
- `obsolete`: no longer relevant.
- `moved`: moved to spec, plan, tracker, or another ticket.
- `answered`: resolved by clarification without implementation.

## Source Values

- `chat`: captured from user/agent chat.
- `manual_test`: found during manual testing.
- `review`: found during PR or code review.
- `trace`: found from logs, traces, or diagnostics.
- `user_report`: reported as observed behavior.
- `implementation`: discovered while implementing another change.

## Tags

Use `tags` for lightweight grouping without changing ticket workflow semantics.

- Use a comma-separated list.
- Prefer lowercase kebab-case tags.
- One tag is valid.
- Three tags is the ideal range.
- Five tags is the maximum.
- Keep tags descriptive and sparse.
- Do not block capture when tags are unclear; leave `tags` empty.

## Lifecycle Fields

- Set `ready_at` on `open -> ready`.
- Set `started_at` on `ready -> in_progress`.
- Set `reviewed_at` when entering or completing `reviewing`.
- Set `closed_at` only when status becomes `solved` or `closed`.
- Set `branch`, `pr`, and `commits` when implementation work exists.

## Branch, Commit, and PR Naming

Use ticket naming that mirrors the slice workflow style while keeping ticket identity visible.

- Branch format: `<type>/ticket-<timestamp>-<short-slug>`.
- PR title format: `<type>(ticket-<timestamp>): <short outcome>`.
- Commit format: regular Conventional Commits, scoped to the affected area.
- PR body format: use the slice delivery sections `Summary`, `Changes`, and `Validation`.

## Workflow Invocation

Canonical workflow ids use `snake_case`, but users do not need to say them literally. Map natural English phrases to the canonical workflow.

- `capture`: `ticket`, `capture ticket`, `record ticket`
- `capture_round`: `tickets`, `capture tickets`, `ticket round`, `record several tickets`
- `triage_manual`: `triage`, `manual triage`, `classify tickets`, `review tickets`
- `triage_auto`: `triage auto`, `auto triage`, `automatic triage`, `classify tickets automatically`

## Capture Workflow

Use `capture` for one quick ticket.

1. Accept a raw description as sufficient input.
2. Do not ask classification questions by default.
3. Ask only when the description is impossible to understand or appears to contain multiple unrelated tickets.
4. Create the ticket under `ops/<namespace>/ticket/open/`.
5. Default uncertain classification fields to `unclassified`.

Default capture frontmatter:

```yaml
status: open
kind: unclassified
severity: unclassified
priority: unclassified
scope: unclassified
tags:
source: chat
```

## Capture Round Workflow

Use `capture_round` when the user wants to dump multiple informal items.

1. Accept a rough list, paragraphs, or short fragments.
2. Split clearly separate items into separate tickets.
3. Merge obvious duplicates.
4. Ask questions only when an item is impossible to preserve as a ticket.
5. Do not spend user energy on classification during capture.
6. Create all preserved items under `ops/<namespace>/ticket/open/`.
7. Report created ticket ids and any deferred candidates.

## Assisted Capture

If the user asks for help clarifying before capture, ask at most 10 questions total for the ticket or round.

- Ask only questions that improve the ticket description, expected outcome, actual behavior, validation, or separability.
- Do not preserve the questions as quiz artifacts.
- The questions are scaffolding; the ticket is the artifact.

## Manual Triage Workflow

Use `triage_manual` to classify and refine existing tickets with user input.

1. Select open or unclassified tickets.
2. Ask the minimum questions needed to classify or clarify.
3. Let the user decide or confirm `kind`, `scope`, `severity`, and `priority`.
4. Move tickets from `open/` to `ready/` only when the expected outcome is actionable.
5. Leave unclear tickets in `open/` with `unclassified` fields.

## Automatic Triage Workflow

Use `triage_auto` when the user explicitly asks the agent to classify tickets by judgment.

1. Classify for execution, not drama.
2. Do not ask questions unless ambiguity is dangerous.
3. Do not invent facts.
4. Prefer conservative classification.
5. It is acceptable to leave fields as `unclassified`.
6. Do not move a ticket to `ready/` unless the expected outcome is actionable.
7. Add a short rationale note only when the classification is not obvious.

Automatic triage criteria:

- `severity` measures impact if unresolved.
- `priority` measures when to schedule it.
- Do not raise priority just because something is annoying.
- Do not raise severity because information is missing.
- Use `unclassified` when critical context is missing.

## Ticket-Driven Implementation Workflow

1. Capture the ticket in `ops/<namespace>/ticket/open/` with `status: open`.
2. Refine it to `status: ready` when the expected outcome is clear enough to execute.
3. Set `status: in_progress` and `started_at` when implementation starts.
4. Create a branch from `dev` and record it only when review delivery is
   requested.
5. Implement the change.
6. For direct `dev` delivery, record the commits and validation, then set
   `status: solved`, `resolution: fixed`, and `closed_at` in the delivery commit
   or a focused closure commit.
7. For branch-and-PR delivery, set `status: reviewing`, `pr`, and `commits` when
   the pull request is ready.
8. After the pull request is merged, set `status: solved`,
   `resolution: fixed`, and `closed_at` on `dev`.
9. Push post-merge ticket metadata directly to `dev`. Do not open a dedicated
   pull request only for closure metadata.

## Capture Rule

If the user mentions a small or informal task that should not be solved immediately, create an `open` ticket under `ops/<namespace>/ticket/open/` instead of leaving the work buried in chat.
