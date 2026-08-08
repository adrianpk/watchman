# Activity Recording

This document defines the repository activity-recording requirement.

- After creating a commit for completed work that is not hypertrivial, run
  `acta record [flags] [note]` immediately while that commit is still `HEAD`.
- Record features, fixes, refactors, meaningful documentation, tests,
  maintenance, decisions, specifications, plans, and operational milestones.
  During initial adoption, record the commit when uncertain.
- Omit only abandoned or work-in-progress commits, generated-only changes,
  punctuation-only changes, isolated formatting, or similarly trivial edits.
- Add a concise note when the commit subject lacks useful retrieval terms or
  does not explain why the change matters. Supply provenance flags only when
  their values are known reliably.
- Never record credentials, tokens, personal data, or hidden reasoning.
