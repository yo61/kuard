## Decision: Redo conflicting Renovate PRs as one hand-made PR that closes them

## Context

On 2026-10-12, merging #45 (`golang.org/x/*` security bumps) left Renovate's #49 (go modules
minor/patch) and #50 (Go 1.27) conflicting on `go.mod` and `go.sum`. Renovate only runs on the
cron in `yo61/infrastructure`, so it would not rebase them until that next fired.

## Alternatives considered

- **Wait for Renovate to rebase**, or tick its rebase checkbox, which also waits for the cron.
  Keeps bot authorship, but leaves the queue half-merged for hours.
- **Push a conflict resolution to Renovate's branch.** Renovate stops managing a branch someone
  else has pushed to, so later updates to those dependencies would stall behind it.
- **Make one branch from main that applies the same versions** with `go get`, `go mod edit` and
  `go mod tidy`, plus the PRs' exact image digests. Open it with `Closes #N` for each PR.

## Reasoning

The third option lands the update the same session, with the same versions and digests
Renovate chose and CI run on the combined result. Renovate's own branches stay untouched. When
main already has the versions, Renovate has nothing to reopen.

## Trade-offs accepted

- The commit is authored by a person, not the bot, and the Renovate PR's release notes are not
  in the squash commit. They remain readable on the closed PR.
- A hand-made PR has no `renovate/stability-days` check. Taking versions and digests exactly
  from Renovate's PRs keeps its release-age gate.
- Renovate's `go` group keeps the `go` directive and the `golang:` image together; by hand,
  that is up to whoever makes the branch.

## Supersedes

None.
