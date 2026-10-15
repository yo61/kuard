# Knowledge index

Routes to domain knowledge folders. Each domain holds `knowledge.md` (facts and patterns) and,
once a hypothesis is confirmed three or more times, `rules.md` (apply by default). Hypotheses
that span repositories live in
[yo61/global `knowledge/`](https://github.com/yo61/global/blob/main/knowledge/INDEX.md).

- [dependency-updates](dependency-updates/knowledge.md) — working through the Renovate and
  Dependabot queue: which bot owns which pin, why Renovate's per-module security PRs overlap so
  that one supersedes the others, why a superseded or hand-combined bot PR needs closing by
  hand or by keyword in this squash-only repo, and the smoke test `pkg/dnsapi` needs because
  it has no tests, why digest-only image updates skip the release-age gate, and which actionlint
  kuard runs. The jactionlint playbook, with kuard's findings and fixes, is
  [yo61/global `docs/jactionlint-migration.md`](https://github.com/yo61/global/blob/main/docs/jactionlint-migration.md).
- [ci](ci/knowledge.md) — why `ci.yaml` skips its required jobs on Markdown-only PRs instead of
  path-filtering the workflow, and why their `if:` starts with `!cancelled()`.

## Decisions

Cross-cutting decisions live in [`../decisions/`](../decisions/).

- `2026-10-12-combine-conflicting-renovate-prs.md` — when Renovate PRs conflict with main,
  redo them as one hand-made PR that closes them by keyword, instead of waiting for
  Renovate's cron to rebase.
- `2026-10-15-exact-base-image-tags.md` — `Dockerfile.ci`'s node and golang images use exact
  version tags plus digest, so a new release is a version update paired with CI's pin.
- `2026-10-15-prek-hooks-without-pins.md` — hooks live in `.pre-commit-config.yaml`, run by prek
  locally and in CI, with tool versions pinned once in `mise.toml`.
- `2026-10-15-jactionlint.md` — kuard lints workflows with `jdx/jactionlint`, the first repository
  to do so; yo61/global's `decisions/2026-10-15-jactionlint-org-wide.md` extends it to the
  other repositories.

## Quality gate

[`../quality/criteria.md`](../quality/criteria.md) — the checks run before marking a task
complete in this repo.
