# Specification Workflow

Use this playbook when a topic is in active design, ready for stable
specification, or needs model and persistence alignment before implementation.

## Intent

- Preserve design context across sessions.
- Keep brainstorming material curated and editable.
- Preserve decisions, tensions, assumptions, and open questions.

## Artifact Inputs

A specification draft can exist without a dump or quiz when the available
context is sufficient. A real quiz or dump can be used as input and linked from
the specification. Do not create synthetic artifacts to satisfy ceremony.

Typical paths are:

- Draft: `ops/<namespace>/spec/drafts/<slug>-draft.md`
- Stable specification: `ops/<namespace>/spec/stable/<slug>.md`
- Quiz input: `ops/<namespace>/quiz/<slug>-quiz.md`
- Dump input: `ops/<namespace>/dump/<slug>-dump.md`

A stable specification can be created directly when no draft cycle is useful.

## Model Companions

When a specification introduces business models, persistence, or durable data
contracts, create a sibling file with the same basename and a `-model.md`
suffix.

Use a model companion for:

- business entities and value objects;
- lifecycle or status fields;
- roles and permissions;
- tables, documents, indexes, and uniqueness rules;
- queries, migrations, and event mappings;
- durable API or message payloads that mirror stored data.

Keep behavioral or workflow-only specifications model-free when no durable
model changes.

## Rules

1. Keep every artifact in English.
2. Preserve explicit decisions, assumptions, constraints, and open questions.
3. Mark provisional decisions clearly.
4. Create quiz and dump artifacts only when those workflows occurred.
5. Do not infer that ordinary discussion was a quiz session.
6. Add cross-references only between artifacts that exist and have a direct
   relationship.
7. Do not use title or heading suffixes in parentheses such as `(Draft)`.
8. Use semantic status metadata when it adds operational value.
9. Review business models together with their persistence shape, indexes,
   queries, migrations, rules, and external-event mappings.
10. Keep model fields and stored fields in an order that supports direct review
    when practical.
11. Run a full consistency pass before promotion or commit.
12. Prefer one cohesive specification update over corrective micro-commits.

## Promotion

- Iterate under `ops/<namespace>/spec/drafts/` while important questions remain.
- Promote to `ops/<namespace>/spec/stable/` when behavior and boundaries are
  sufficiently decided.
- Promote the model companion with the specification when it is part of the
  stable contract.
- Create a plan and tracker only after implementation scope is approved.

## Quiz-Driven Direct Specification

A completed quiz can provide:

- raw capture: `ops/<namespace>/quiz/<slug>-quiz-raw.md`;
- curated notes: `ops/<namespace>/quiz/<slug>-quiz.md`.

When these artifacts preserve decisions and open questions, the curated quiz
can be the direct source for a stable specification. A separate draft is
optional. Do not create a draft only to satisfy ceremony.
