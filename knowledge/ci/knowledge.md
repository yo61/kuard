# CI — knowledge (facts & patterns)

Scope: `.github/workflows/ci.yaml` and the checks kuard's `main` ruleset requires.

## Required checks and path filters

The ruleset (managed in yo61/github-repos, `data/yo61/kuard.yaml`) requires `go tests`,
`image builds` and `client tests`, all from `ci.yaml`. Until 2026-10-15 `ci.yaml` had
`paths-ignore: '**.md'`, so a Markdown-only PR never started the workflow and the three checks
never reported: GitHub showed "Waiting for status to be reported" and only an admin bypass could
merge it (#62). GitHub's docs say not to path-filter a workflow whose checks are required.

The workflow now runs on every PR. A `detect changes` job lists the PR's files with
`gh api .../pulls/<n>/files` (renames count by their old name too) and sets `docs_only`. The
three required jobs skip when it is `true`, and a job skipped by its `if:` reports success.

- **Their `if:` starts with `!cancelled()`.** Without a status function, `if:` implies
  `success()`, so a failed `detect changes` would skip the jobs, and skipped counts as passed.
  With it, a failure leaves `docs_only` empty and the full CI runs.
- **Don't add a second workflow reporting the same job names** for docs PRs. A PR touching
  both Markdown and code would then report two runs under one check name.
- `prek hooks` isn't required and runs on every PR, Markdown included.
