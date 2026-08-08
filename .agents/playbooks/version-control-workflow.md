# Version Control Workflow

1. Verify the branch, working tree, canonical remote, and active operational
   artifacts.
2. Use `dev` for approved ad hoc maintenance, unsliced implementation, and
   planning-only commits.
3. Create a branch and pull request only when the user requests that delivery
   mode or an approved plan and tracker record it.
4. For formal slices, create the exact recorded branch from current `dev`.
5. Never work directly on `main`.
6. Push active branches to `origin`. Do not push feature or `dev` branches to a
   publication mirror.
7. Before `dev -> main`, run the release-alignment gate.
8. Align `dev -> main` only through an explicitly requested pull request.
9. After `main` changes, publish it to `mirror` only when that remote exists and
   local policy requires it.
10. Avoid force pushes and history rewrites unless explicitly approved.
11. For docs-only changes, validate paths, links, and `git diff --check` instead
    of unrelated runtime tests.
