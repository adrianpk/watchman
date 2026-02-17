# Changelog

All notable changes to Watchman will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.2] - 2026-02-17

### Changed

- **Sentinel**: Structured violation messages for actionable feedback
  - Block messages now include: rule name, file, line numbers, and specific detail
  - Format: `[RULE] in [FILE:LINE]. [DETAIL]`
  - Multiple violations are listed with numbered entries
  - Ambiguous cases indicate uncertainty for human review

### Fixed

- Sentinel no longer returns vague messages like "Code style violation detected"

## [0.1.1] - 2026-02-16

### Added

- **Sentinel plugin**: AI-powered code standards enforcement
  - Evaluates code against project CONVENTIONS.md
  - Supports Anthropic, OpenAI, and Ollama providers
  - Fallback chain when multiple providers configured
- **jj (Jujutsu) support**: Sentinel now intercepts `jj commit` and `jj describe`
- **Commit evaluation mode**: Option to evaluate only on git add/commit instead of every write
- Deny logging to watchman log file

### Changed

- Use working directory from Claude Code input instead of os.Getwd()
- Separated behavior rules (AGENTS.md) from code standards (CONVENTIONS.md)

### Fixed

- Reminders implementation corrected

## [0.1.0] - 2026-01-13

### Added

- Initial release
- **Core rules engine** with declarative YAML configuration
- **Scope rule**: Restrict file access to specific directories
- **Versioning rule**: Enforce version bumps on protected files
- **Incremental rule**: Prevent large-scale changes in single operations
- **Invariants rule**: Block patterns that should never appear in code
- **External hooks**: Execute custom scripts on tool invocations
- **Periodic reminders**: Inject context reminders into agent workflow
- **Command blocklist**: Prevent dangerous shell commands
- `watchman init` command for configuration scaffolding
- `watchman setup` command for Claude Code integration
- Normalized path handling for cross-platform compatibility
