# Dump Execution Playbook

Use this playbook when the user wants to freely dump knowledge about a project,
feature, design concern, implementation area, or unresolved thought cluster.

## Purpose

Capture user-provided context with less constraint than a quiz while still
turning it into reusable repository artifacts.

A dump is not an interview. The user leads with a broad knowledge transfer; the
agent preserves, lightly corrects, groups, and curates it.

## Core Principles

- User-led flow: let the user dump freely before steering.
- High preservation: raw capture keeps the user's intent and useful nuance.
- Light correction: fix spelling, punctuation, and phrasing without changing meaning.
- Center of gravity: every dump has a central topic even when it branches.
- Thread awareness: detect derived threads without forcing them into the main topic.
- Reversible curation: the curated dump should make it easy to decide whether to
  create a draft, spec, plan, ticket, quiz, or another dump later.

## Artifact Locations

Dump artifacts live under `ops/<namespace>/dump/`.

Use these filename patterns:

- Raw dump capture: `ops/<namespace>/dump/YYYY-MM-DD-<topic>-dump-raw.md`
- Curated dump notes: `ops/<namespace>/dump/YYYY-MM-DD-<topic>-dump.md`

## Raw Capture

Maintain raw capture when the dump is long, multi-message, or likely to feed a
future artifact.

Raw capture rules:

1. Preserve the user-provided sequence enough to recover the thinking path.
2. Correct spelling, punctuation, grammar, and obvious slips while preserving intent.
3. Keep first-person meaning when it carries useful ownership or uncertainty.
4. Mark obvious topic shifts with headings.
5. Record derived threads separately instead of flattening them into the central topic.
6. Do not invent decisions that the user did not state.
7. Do not turn the raw capture into a polished spec.

## Curated Dump

Create or update the curated dump after the raw material has enough signal.

The curated dump should include:

- central topic
- source/raw capture path
- distilled summary
- key concepts
- decisions or strong preferences
- tensions and tradeoffs
- assumptions
- derived threads
- possible downstream artifacts
- open questions

Do not map the raw dump one-to-one into the curated file. Preserve signal, not
transcript shape.

## Steering Rules

The agent may briefly steer when:

- the topic loses its center of gravity
- a derived thread becomes large enough to name
- the user appears to contradict an earlier point
- a concrete downstream artifact is becoming obvious
- the dump needs a pause for summarization

The agent should not over-interview. If many questions are needed, suggest
switching to the quiz playbook.

## Derived Threads

A derived thread is a branch that may deserve its own artifact later.

For each derived thread, record:

- name
- why it emerged
- relation to the central topic
- whether it looks like a draft, spec, plan, ticket, quiz, or later dump
- whether the user explicitly wants to pursue it now

Do not create derived artifacts unless the user asks.

## Handoff

When the dump is complete, report:

- raw dump path
- curated dump path
- central topic
- derived threads
- recommended next artifact options
- unresolved questions

Do not create specs, plans, tickets, or derived drafts unless explicitly requested.
