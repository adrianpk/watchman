# Planned Delivery Workflow

Use this workflow for a maintainer-approved large feature, architecture change,
or multi-slice delivery set. It is a defined loop of slices, pull requests,
maintainer-controlled merges, and checkpoints. It is not open-ended autonomous
work.

1. Start from an approved long-form idea, ticket, or specification.
2. Create `ops/<namespace>/plan/` and `ops/<namespace>/tracker/` artifacts.
3. Define the slicing strategy, ordered slices, tasks, validation, branches,
   pull-request titles, and report paths.
4. Commit the approved plan and tracker to `dev` before implementation.
5. Implement tasks in order. Use one Conventional Commit per task when
   practical.
6. Produce one implementation report per slice.
7. Open one pull request per slice, targeting `dev`.
8. Complete review and wait for the maintainer-controlled merge.
9. Report slice completion as `Slice N of M completed`.
10. After merge confirmation, verify the merge, close the slice artifacts, and
    start the next slice in the same approved delivery set.
11. After the final slice, publish one end-to-end cycle report when the plan
    requires it.
12. Align `dev -> main` only through the release-alignment playbook after an
    explicit request.
