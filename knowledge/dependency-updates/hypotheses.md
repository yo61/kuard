# Dependency updates — hypotheses (need more data)

Promote to `rules.md` after three confirmations; record each one below.

## Merge the broadest Renovate security PR first

When Renovate opens several `[security]` PRs for modules in the same family (e.g.
`golang.org/x/*`), the one whose diff moves the most modules is a superset of the rest.
Merging it first leaves the others redundant, avoiding a conflict on every one of them.

What would disprove it: a narrower PR that pins some module *higher* than the broad one does,
so closing it would lose a fix.

Confirmations:

- 2026-10-12 — #45 (x/net) superseded #43, #44 and #46.
