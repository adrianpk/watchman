# Quiz Execution Playbook

Use this playbook when an agent conducts a structured quiz with a user for a software engineering use case.

## Purpose

Define how an agent conducts structured, bounded, and high-signal quiz sessions.

The objective is to extract enough information to define or refine a use case without introducing noise, redundancy, or unbounded exploration.

This playbook governs interaction dynamics. It complements separate playbooks that define how results are curated and recorded.

## Core Principles

- Boundedness: every session operates within explicit limits for question count and nesting depth.
- Sequentiality: ask one question at a time. Do not batch questions.
- Relevance over completeness: stop when sufficient clarity is reached, not when limits are exhausted.
- Controlled exploration: follow-up questions are allowed, but strictly constrained.
- Coverage discipline: do not over-explore one dimension at the expense of others.
- Non-transcriptive outcome: the session produces insights, not logs.

## Session Initialization

Before asking the first question, the agent must initialize the session.

### Precompute a Question Set

- Define a finite set of candidate questions.
- Establish a maximum cap of 50 questions.
- Organize questions into logical groups or themes.
- Identify possible follow-up chains when applicable.

### Define Completion Criteria

- Internally determine what sufficient understanding means for the given use case.
- Treat sufficient understanding as the primary stopping condition, above the numeric cap.

### Initialize Tracking State

- Total questions asked.
- Remaining question budget.
- Current topic or group.
- Nesting depth when inside a follow-up chain.
- Coverage distribution across topics.

## Question Execution Model

### Sequential Delivery

- Ask exactly one question at a time.
- Wait for the user response before proceeding.
- Do not anticipate or pre-ask future questions.

### Adaptive Progression

After each answer, decide one of the following:

- Continue with the next planned question.
- Enter a follow-up question.
- Switch topic.
- Terminate the session.

The decision must be based on information gain, not on exhausting the precomputed list.

## Nested Questioning

Follow-up questions may refine or clarify an answer, but they must follow strict limits.

### Maximum Depth

- Level 1: primary question.
- Level 2: follow-up.
- Level 3: deep clarification.

No further depth is allowed beyond level 3.

After reaching maximum depth, return to the original question level or switch topic.

### Constraints

- Do not over-invest in a single branch.
- Resume broader coverage after closing a follow-up chain.
- Do not let nested questioning create local overfitting.

## Question Budgeting

- The hard cap is 50 questions per session.
- The cap is a limit, not a target.

Stop early when:

- The use case is sufficiently defined.
- Additional questions produce diminishing returns.
- Remaining questions are low-value or redundant.

## Topic Coverage Control

The agent must ensure that all relevant dimensions of the use case are explored at least minimally.

No single dimension should dominate the session disproportionately.

If excessive depth occurs in one area:

- Force a topic shift.
- Resume pending high-level questions.

## State Awareness

At all times, the agent must maintain:

- Current question index.
- Remaining budget.
- Active topic.
- Nesting level.
- Coverage distribution.

Loss of this state is a protocol violation.

## Termination Conditions

End the session when any of the following is true:

- Clarity achieved: the agent can confidently describe the use case and its constraints.
- Budget exhausted: 50 questions have been reached.
- Diminishing returns: new answers do not materially improve understanding.

After termination, no additional exploratory questions are allowed.

## Output Handoff

The quiz session does not produce a transcript.

It produces curated insights extracted from the interaction and structured according to the canonical format defined by the relevant curation playbook.

Quiz artifacts live under `ops/<namespace>/quiz/`.

Use these filename patterns:

- Curated quiz notes: `ops/<namespace>/quiz/YYYY-MM-DD-<topic>-quiz.md`
- Raw quiz capture: `ops/<namespace>/quiz/YYYY-MM-DD-<topic>-quiz-raw.md`

When the quiz has enough coverage and a curated `ops/<namespace>/quiz/...-quiz.md` exists, that quiz artifact may be sufficient pre-spec alignment. In that case, a separate spec draft is optional rather than mandatory, and the next artifact may be a formal spec under `ops/<namespace>/spec/stable/`.

## Raw Capture During Execution

For long quiz sessions, maintain a sibling raw capture file while the session progresses.

- Use the filename pattern `ops/<namespace>/quiz/YYYY-MM-DD-<topic>-quiz-raw.md`.
- Update this file incrementally after each answer or small group of related turns.
- Preserve enough question and answer context to prevent information loss during long sessions.
- Record user answers in English with polished style and fully corrected spelling and grammar while preserving the user's intent.
- Treat the raw capture as temporary working memory, not as the final curated artifact.
- After the quiz ends, curate the raw capture into the canonical `ops/<namespace>/quiz/YYYY-MM-DD-<topic>-quiz.md` quiz notes.
- Do not map the raw file one-to-one into the curated file; retain only distilled, decision-relevant content.
- Raw capture files may be versioned when they preserve useful alignment context, but they remain working-memory artifacts rather than curated final notes.

### Key Rule

- Do not map questions one-to-one to outputs.
- Do not include conversational artifacts.
- Retain only distilled, decision-relevant information.

## Failure Modes

The following indicate incorrect execution:

- Asking multiple questions at once.
- Creating infinite or uncontrolled follow-up chains.
- Ignoring the question cap.
- Treating the cap as a goal.
- Overfitting to a single topic.
- Producing raw transcripts instead of insights.
- Losing track of session state.

## Summary

The agent behaves as a constrained interviewer:

- Plan ahead, but adapt in real time.
- Explore, but within strict boundaries.
- Stop early when appropriate.
- Produce signal, not verbosity.

The result is a controlled, high-efficiency extraction process aligned with downstream structuring rules.
