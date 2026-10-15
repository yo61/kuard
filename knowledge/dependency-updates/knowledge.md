# Dependency updates — knowledge (facts & patterns)

Scope: reviewing and merging the bot PRs that keep kuard's pins current.

## Who opens what

`renovate.json5` is the source of truth for the split; read it rather than a copy here. In
short, Dependabot owns npm and action SHAs, and Renovate owns everything else. That shapes the
conflicts: npm PRs can only collide on `client/package-lock.json`, and Renovate PRs collide on
`go.mod`/`go.sum` and `Dockerfile.ci`.

Renovate runs only on the cron in `yo61/infrastructure`'s workflow, not on this repo's events.
A Renovate PR left conflicting by a merge stays conflicting until that cron next fires. See
`../../decisions/2026-10-12-combine-conflicting-renovate-prs.md` for what to do instead of
waiting.

## Security PRs overlap

`vulnerabilityAlerts` makes Renovate open one PR **per vulnerable module**. Each also pulls
whatever its target version requires, so the PRs overlap. On 2026-10-12:

| PR | Named module | Also moved |
|---|---|---|
| #43 | x/crypto 0.56 | text 0.41, go 1.26 |
| #44 | x/mod 0.40 | crypto 0.55, net 0.58, text 0.41, tools 0.49 |
| #45 | x/net 0.60 | crypto 0.57, mod 0.41, sync 0.23, sys 0.48, text 0.42, tools 0.49, go 1.26 |
| #46 | x/text 0.41 | — |

#45 was a superset of the others. Merging it first left nothing for #43, #44 or #46 to do, so
they were closed. Before closing a PR as superseded, check main's `go.mod` meets every
version in that PR's diff, not just the module named in its title.

## Closing bot PRs in a squash-only repo

The repo allows squash merges only. A bot PR's commits never reach main, so GitHub never marks
it merged when its change lands some other way. Close superseded PRs explicitly with a comment
naming what superseded them. A PR that redoes bot PRs puts `Closes #N` in its body so they
close on merge.

## `gh` targets upstream unless told otherwise

This clone has both `origin` (yo61/kuard) and `upstream` (kubernetes-up-and-running/kuard).
`gh` resolved its default to upstream, so `gh pr list` silently listed upstream's PRs. This is
set per clone with `gh repo set-default yo61/kuard` (stored in `.git/config`), so a fresh clone
needs it again. Until then, pass `-R yo61/kuard`.

## Verifying a Go bump

- `pkg/sitedata/bindata.go` is generated and gitignored, so run `go generate ./cmd/...
  ./pkg/...` (after `go install` of the go-bindata version pinned in `ci.yaml`) before
  `go vet` or `go test` compiles.
- `pkg/dnsapi` is the only consumer of `github.com/miekg/dns`, and it has **no tests**, so a
  green `go test ./...` says nothing about a dns bump. Run the binary instead:

  ```bash
  go build -o "$TMPDIR/kuard" ./cmd/kuard
  "$TMPDIR/kuard" --address=127.0.0.1:18080 &
  curl -s -X POST -H 'Content-Type: application/json' \
    -d '{"type":"A","name":"example.com"}' http://127.0.0.1:18080/dns/api
  ```

  A good response has `status: NOERROR` and an ANSWER section.
- CI's `go tests` job reads the Go version from `go.mod` (`go-version-file`), so the `go`
  directive and `Dockerfile.ci`'s `golang:` tag must move together. Renovate's `go` group
  enforces this; a hand-made PR has to do it itself.

## Digest-only updates skip the release-age gate

A Docker digest update evidently carries no release timestamp, and `minimumReleaseAgeBehaviour:
"timestamp-optional"` lets it through. On a floating tag like `node:24-alpine` that let #57 ship
Node 24.22.0 one day after release, unpaired with CI's `NODE_VERSION`. Base images that have a
matching pin elsewhere use exact tags (`node:24.21.0-alpine@…`, `golang:1.27.2-alpine@…`) so a
new release is a version update. See `../../decisions/2026-10-15-exact-base-image-tags.md`.
To see which version a floating tag resolves to:

```bash
docker buildx imagetools inspect node:24-alpine \
  --format '{{json (index .Image "linux/amd64").Config.Env}}'
```

## Workflows are linted by jactionlint, not actionlint

kuard uses `jdx/jactionlint`, pinned in `mise.toml` as `"github:jdx/jactionlint"` (Renovate's mise
manager updates it from GitHub releases). This departs from the yo61-wide choice of
`kjanat/actionlint`; see `../../decisions/2026-10-15-jactionlint.md`. Its config is
`.github/jactionlint.yaml`; jactionlint also reads `actionlint.yaml`, and one without a `profile`
runs the stricter `default` profile. `jactionlint --init-config` writes the config file into
`.github/` even when `-c` names another file. kuard runs the `default` profile; how to migrate
another repository, and the fixes its findings needed, are in `../../docs/jactionlint-migration.md`.
