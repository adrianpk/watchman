# Dev-to-Main Release Alignment

Use this playbook only after an explicit request to publish `dev` to `main`.

1. Verify the working tree, `dev`, `origin/dev`, `main`, `origin/main`, and open
   release pull-request state.
2. Verify that included delivery artifacts are closed.
3. Run the broad relevant release gate from `.agents/LINTING.md`.
4. For a Markdown-only diff, validate paths, links, and `git diff --check`
   instead of unrelated runtime tests.
5. Open a pull request from `dev` to `main` titled
   `chore(release): align dev with main`.
6. Use this body shape:

```md
## Summary

- align `main` with the latest merged state from `dev`
- include all previously reviewed and merged delivery work

## Validation

- <command and result>
```

7. Report the gate result and stop for maintainer review.
8. The maintainer owns the final merge unless explicit merge authorization is
   given.
9. When authorized, merge through the pull request using the repository's
   approved rebase and fast-forward method. Verify the resulting refs.
10. When `mirror` exists and local policy requires publication, push updated
    `main` to it and verify its allowed refs.
