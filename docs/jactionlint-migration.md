# Migrating a repository from actionlint to jactionlint

How kuard moved from `rhysd/actionlint` to [`jdx/jactionlint`](https://github.com/jdx/jactionlint),
written as the playbook for the other yo61 and horopter-dev repositories. kuard did it in two pull
requests: #60 swapped the linter, and the follow-up turned on the `default` profile and fixed what
it found. Why jactionlint was chosen is in
[`decisions/2026-10-15-jactionlint.md`](../decisions/2026-10-15-jactionlint.md).

## Why two phases

jactionlint's `correctness` profile runs the checks actionlint runs. Its `default` profile adds
security and policy rules that fail the build. Swapping the tool and turning those rules on in one
change mixes a tooling change with workflow changes, and the workflow changes can alter how a
repository builds or publishes. Keep them apart:

1. **Swap**: replace actionlint with jactionlint on `profile: correctness`. No workflow changes.
2. **Tighten**: set `profile: default` and fix every finding, in its own pull request.

## Before you start

- **An existing `actionlint.yaml` changes behaviour.** jactionlint reads `.github/actionlint.yaml`
  as its own config. If that file sets no `profile`, the `default` profile runs and jactionlint
  prints a `note:` saying so. Phase 1 then fails on findings it was meant to postpone. Replace the
  file with `.github/jactionlint.yaml`.
- **Release age.** Repositories follow a 7-day `minimumReleaseAge`. jactionlint releases often, so
  the newest version may be younger than that. kuard took 2.2.0 at one day old by explicit
  exception, because no v2 release was old enough and v1.8.x predates v2's command-line changes.
  Record any such exception in the repository's decision.
- **Who owns the pin.** Every pin needs exactly one updater. Choose one:
  - `mise.toml`: `"github:jdx/jactionlint" = "<version>"`. Renovate's mise manager reads the
    `github:` backend from GitHub releases. Don't rely on the registry short name `jactionlint`;
    the explicit backend is what Renovate is known to resolve.
  - `.pre-commit-config.yaml`: `repo: https://github.com/jdx/jactionlint` with `rev: v<version>`,
    owned by Renovate's `pre-commit` manager. That manager is off by default, so set
    `"pre-commit": {enabled: true}` and add it to `enabledManagers`.

  kuard uses the first, because CI installs its tools from `mise.toml` with `jdx/mise-action` and
  the pre-commit hook runs the same binary as a `local` hook.
- **mise's deprecation warning names jactionlint.** mise warns that `actionlint` is deprecated in
  favour of `jactionlint`. That points at `jdx/jactionlint`, not the `kjanat/actionlint` fork the
  other repositories use. The warning goes away once `mise.toml` no longer names `actionlint`.

## Phase 1: swap the linter

1. In `mise.toml`, replace the actionlint entry with `"github:jdx/jactionlint" = "<version>"` and
   run `mise install`.
2. Create `.github/jactionlint.yaml`:

   ```yaml
   # yaml-language-server: $schema=https://raw.githubusercontent.com/jdx/jactionlint/main/schema/jactionlint-config.json
   profile: correctness
   ```

   Don't use `jactionlint --init-config` in a repository you're editing: it writes
   `.github/jactionlint.yaml` even when `-c` names another file, and the file it writes is mostly
   comments.
3. Delete `.github/actionlint.yaml` once nothing in it is still needed. jactionlint knows the
   `ubuntu-26.04` runner label, so kuard's only entry went away.
4. Point every caller at the new binary: the pre-commit hook, Taskfile or Makefile targets, and CI.
   With no arguments, jactionlint checks the workflows, composite actions and `dependabot.yaml`, so
   kuard's hook matches `^\.github/(workflows/.*|dependabot)\.ya?ml$`.
5. Check that `jactionlint` exits 0, then break a workflow on purpose and confirm the hook fails.

## Phase 2: turn on the default profile

1. See what you're facing without changing the config: `jactionlint --profile default`.
2. Fix each finding (see below), then set `profile: default`.
3. Optionally run `jactionlint --profile pedantic` to see the next level. kuard took one pedantic
   finding (below) and left the style rules (`require-shell`, `anonymous-definition`).

For scripts, `--format json` gives each finding's rule in `id`. The `kind` field can differ: the
timeout rule has `kind: timeout-check` but `id: missing-timeout`, and the config and
`--explain` accept only `missing-timeout`.

## Findings in kuard and their fixes

| Rule | Where | Fix |
|---|---|---|
| `missing-timeout` | all 5 jobs | `timeout-minutes` per job |
| `concurrency-limits` | build-kuard | `concurrency:`, `cancel-in-progress: false` |
| `excessive-permissions` | build-kuard | `permissions: {}`; job gets its own |
| `cache-poisoning` | `build-push-action` | drop `cache-from`, keep `cache-to` |
| `cache-poisoning` (pedantic) | `setup-buildx-action` | `cache-binary: false` |

### Timeouts

A job without `timeout-minutes` can run for GitHub's six-hour limit, and a matrix multiplies it.
A stalled job has used up a whole organisation's free minutes before. Set each timeout from the
job's observed run time, about 3-5×, not a round comfortable number: the timeout is the most a hung
job can cost. Get durations with
`gh run view <run-id> --json jobs -q '.jobs[] | "\(.name) \(.startedAt) \(.completedAt)"'`.
kuard's jobs take 15 s to 3 min, so they got 5, 10 and 15 minutes.

### Concurrency on a publishing workflow

Use `cancel-in-progress: false` on a workflow that publishes, so a newer run waits instead of
cancelling a push halfway. A waiting run is replaced by a newer waiting run. That's fine when
each publish overwrites the same tags, as kuard's do. Add `queue: max` where every publish must
happen.

### Permissions

Job-level `permissions:` replace the workflow's rather than adding to them, so the job lists
everything it needs, including `contents: read` for checkout.

### Caches in a publishing job

A publishing job shouldn't restore a cache that another run wrote. kuard's `build-push-action`
keeps `cache-to`, because `ci.yaml`'s PR image build reads that cache, and drops `cache-from`.
Check the cost first: kuard's first-ever build, with no cache, took about as long as its cached
builds (about 2-3 min per matrix job), so dropping the restore cost nothing noticeable.

## Overlap with zizmor

jactionlint's security rules follow zizmor's audits, and at its default settings it reported
things zizmor did not (`cache-poisoning`, `excessive-permissions`). kuard runs both for now.
Whether jactionlint replaces zizmor is an org-wide decision. If it does,
`jactionlint --migrate-ignores` turns `# zizmor: ignore[...]` comments into jactionlint's own.
