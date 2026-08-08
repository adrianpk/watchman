# Repository Language Rules

This document defines language and locale boundaries for versioned artifacts.

## Repository Content

- Keep repository content in English unless the repository explicitly requires
  localized content.
- This includes code, names, comments, configuration, documentation, examples,
  quoted text, test data, commit metadata, pull-request text, and operational
  artifacts.
- Use neutral English examples unless a product requirement defines another
  language or locale.

## Conversation Boundary

- Conversation language does not determine repository language or locale.
- Localized product content does not change the language of surrounding code,
  documentation, metadata, or operational artifacts.

## Commit Check

Before committing, inspect staged natural-language content and remove unintended
non-English text.
