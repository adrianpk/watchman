# Versioning Rules

This document defines commit, branch, remote, and merge rules.

Before changing version-control state, verify the current branch, working tree,
canonical remote, relevant upstream refs, and active operational artifacts.

## Commit Format

Use Conventional Commits:

```text
<type>(<optional-scope>): <description>
```

Common types are `feat`, `fix`, `docs`, `refactor`, `test`, `chore`, `ci`,
`build`, and `perf`.

## Remote Policy

- `origin` is the canonical remote and receives active branches.
- When `mirror` exists, it is a publication remote and receives only `main`
  unless local instructions explicitly define another policy.
- Open and manage pull requests against the canonical forge, not the mirror.

## Branch Policy

- `main` is protected. Do not commit directly to it.
- `dev` is the mandatory trunk for planned and day-to-day delivery work.
- Ad hoc fixes, maintenance, documentation, and unsliced implementation happen
  directly on `dev` unless the user requests branch-and-PR delivery.
- Create a short-lived branch only when the user explicitly requests one or an
  approved slicing plan and tracker record it.
- Do not infer branch or pull-request authorization from task size, code risk,
  validation scope, review value, or the fact that implementation changed.
- Formal slice branches use the exact name recorded in the plan and tracker.
- `main` receives changes only through an explicitly requested `dev -> main`
  release-alignment pull request.

## Pull Request Policy

- Pull requests target `dev` and require explicit authorization, directly or
  through an approved slicing plan and tracker.
- Do not open pull requests for ad hoc fixes, maintenance, documentation,
  tickets, or unsliced implementation unless the user requests one.
- Pull requests targeting `main` are not allowed during normal delivery.
- The formal `dev -> main` alignment also requires an explicit request.

## Changelog Policy

- Treat `CHANGELOG.md` as reader-facing release notes, not as a commit log.
- Describe what users can now do and meaningful user-visible problems that
  were fixed.
- Use plain English and group entries by user task or outcome.
- Omit implementation details, package names, tests, internal refactors,
  operational records, development topology, and minor maintenance.
- Before a release, compare `dev` with the latest tag and ensure every
  release-worthy user-facing change is represented.
- Update and validate the User Guide before finalizing the changelog.
- Propose the next SemVer with a rationale and wait for explicit maintainer
  approval before writing it into the changelog.
- When the next `main` alignment will be tagged immediately, replace
  `Unreleased` on `dev` with `## [x.y.z] - YYYY-MM-DD`.
- Commit and push the final changelog to `dev` before opening the
  `dev -> main` pull request.
- Do not modify release notes directly on `main`; tag the exact integrated
  `main` commit.

## Alignment Policy

- Keep regular delivery linear through `dev`.
- Use rebase and fast-forward for an approved `dev -> main` alignment when the
  forge and repository support that method.
- If an unintended pull request targets `main`, stop and realign through `dev`
  without destructive history edits.
- After `main` changes, push it to `mirror` when that remote exists and local
  policy requires publication.

## Workflow and Safety

- Prefer small, reviewable commits.
- Keep commit scope aligned with one active task or slice concern.
- Treat force pushes and published-history rewrites as destructive Git
  operations governed by `.agents/INVARIANTS.md`.
- Planning and tracker maintenance can be committed directly to `dev` when it
  does not include runtime implementation.
- Follow `.agents/LANGUAGE.md` for commit and pull-request language.
- Keep the working tree state explicit before declaring a delivery cycle
  complete.
