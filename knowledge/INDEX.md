# Knowledge index

Routes to domain knowledge folders. Each domain holds `knowledge.md` (facts and patterns),
`hypotheses.md` (need more data) and, once a hypothesis is confirmed three or more times,
`rules.md` (apply by default).

- [dependency-updates](dependency-updates/knowledge.md) — working through the Renovate and
  Dependabot queue: which bot owns which pin, why Renovate's per-module security PRs overlap so
  that one supersedes the others, why a superseded or hand-combined bot PR needs closing by
  hand or by keyword in this squash-only repo, and the smoke test `pkg/dnsapi` needs because
  it has no tests. Open hypothesis in
  [dependency-updates/hypotheses.md](dependency-updates/hypotheses.md).

## Decisions

Cross-cutting decisions live in [`../decisions/`](../decisions/).

- `2026-10-12-combine-conflicting-renovate-prs.md` — when Renovate PRs conflict with main,
  redo them as one hand-made PR that closes them by keyword, instead of waiting for
  Renovate's cron to rebase.

## Quality gate

[`../quality/criteria.md`](../quality/criteria.md) — the checks run before marking a task
complete in this repo.
