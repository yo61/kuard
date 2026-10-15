# Quality criteria

The gate run before marking a task complete in this repo. Each category names where it came
from. A criterion triggers when it catches a real issue, not when it is applied and passes.
Three triggers promote a criterion to ALWAYS CHECK; ten evaluations without a trigger make it
a candidate for pruning.

---

## Category: Go changes

## Criteria:

    - Run the same steps as CI's `go tests` job, in its order: `go generate ./cmd/... ./pkg/...`,
      then gofmt over `git ls-files '*.go'`, then `go vet ./...`, then `go test ./...`.
    - A change that touches `pkg/dnsapi` or bumps `github.com/miekg/dns` is checked by running
      the binary and making a `/dns/api` lookup, because the package has no tests.
    - The `go` directive in `go.mod` and the `golang:` tag in `Dockerfile.ci` name the same
      exact version (`go 1.27.2` and `golang:1.27.2-alpine`).

## Severity: blocking

## Source: `.github/workflows/ci.yaml` (`go tests` job); 2026-10-12 Go 1.27 / miekg/dns
1.1.73 update (#52).

## Last triggered: never

---

## Category: Client changes

## Criteria:

    - `npm ci` succeeds in `client/`, so `package-lock.json` is in sync with `package.json`.
    - `npm test` passes in `client/`. It builds the production bundle and asserts against the
      output.

## Severity: blocking

## Source: `.github/workflows/ci.yaml` (`client tests` job).

## Last triggered: never

---

## Category: Pins and their owners

## Criteria:

    - A new pin (module, image, tool version, action SHA) is updated by exactly one of
      Renovate or Dependabot. Check it against `enabledManagers` and `customManagers` in
      `renovate.json5` and against `.github/dependabot.yaml`.
    - Base images in `Dockerfile.ci` are pinned by digest.
    - The `node:` tag in `Dockerfile.ci` names the exact version in `ci.yaml`'s `NODE_VERSION`
      (`node:24.21.0-alpine` and `24.21.0`). A floating tag like `node:24-alpine` fails this:
      its digest updates change the Node version without saying so.

## Severity: blocking

## Source: #38 (give every pin exactly one owner); `renovate.json5` `pinDigests` rule;
#57 (`node:24-alpine` moved to Node 24.22.0 while CI stayed on 24.21.0).

## Last triggered: 2026-10-15 (node tag criterion, #57)

---

## Category: Dependency queue

## Criteria:

    - Before closing a bot PR as superseded, main's `go.mod` (or lockfile) meets every version
      in that PR's diff, not just the module in its title.
    - After the queue is worked, `gh pr list -R yo61/kuard` shows no bot PR left open
      without a stated reason, and the `build-kuard` run for the newest commit on main that
      changed something other than `**.md` has succeeded. `build-kuard.yaml` ignores
      markdown-only pushes, so a docs commit at main's head has no run.

## Severity: blocking

## Source: 2026-10-12 dependency queue (#39–#52).

## Last triggered: never

---

## Category: Docs that name versions

## Criteria:

    - A change that moves a version named in `FORK.md` or in a `Dockerfile.ci` comment updates
      that text in the same PR, or the text is confirmed to still hold (for example, it
      describes history).

## Severity: warning

## Source: `FORK.md` names the node, webpack, Babel, React and Go versions the fork moved to.

## Last triggered: never
