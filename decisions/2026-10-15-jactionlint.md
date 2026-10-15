## Decision: kuard lints its workflows with jdx/jactionlint 2.2.0 on the `correctness` profile

## Context

kuard pinned upstream `rhysd/actionlint` 1.7.12, which is unmaintained and needed
`.github/actionlint.yaml` to accept the `ubuntu-26.04` runner label. The yo61-wide decision of
2026-10-09 (`yo61/github-repos`) chose the `kjanat/actionlint` fork. mise's registry now marks
`actionlint` deprecated in favour of `jdx/jactionlint`.

## Alternatives considered

- **`kjanat/actionlint` 1.17.0.** Drop-in, compatible CLI and config, already used across yo61
  and horopter-dev.
- **`jdx/jactionlint` 2.2.0.** Covers actionlint's checks and adds zizmor-style security and
  policy rules, profiles, baselines and `--fix`. Created 2026-10-05; v2.0.0 changed its CLI.

Both fork rhysd's history, are maintained by one person and ship releases that pass
`gh attestation verify`.

## Reasoning

Robin chose jactionlint for kuard on 2026-10-15 and will move the rest of the org separately. Its
`default` profile found eight errors on kuard (four in `build-kuard.yaml`, including
workflow-wide `packages: write` and a cache restored in the job that pushes the image, and a
missing timeout on every job). `correctness` checks what actionlint checks, so this change swaps
the linter without changing the workflows; moving to `default` comes with those fixes.

## Trade-offs accepted

- 2.2.0 was one day old when pinned, inside the 7-day `minimumReleaseAge`; no v2 release was
  older. Robin accepted it.
- kuard differs from the rest of yo61 until the org moves.
- zizmor stays, though jactionlint's `default` profile overlaps it. Whether one replaces the
  other is for the org-wide decision.

## Supersedes

For kuard only: item 3 of `yo61/github-repos/decisions/2026-10-09-yo61-dependency-update-remediation.md`.
