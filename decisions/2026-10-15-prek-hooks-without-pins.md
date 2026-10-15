## Decision: Git hooks live in `.pre-commit-config.yaml`, run by prek locally and in CI

## Context

The repo had no git hooks, so gofmt, actionlint and zizmor ran only when someone remembered
`task fmt:check` or `task lint:workflows`. CI ran gofmt but neither workflow linter.

## Alternatives considered

- **`prek.toml` with `repo: builtin`.** Pins nothing, but upstream pre-commit reads neither the
  file nor `builtin`, so contributors using pre-commit get no hooks.
- **Remote repos for every tool** (`rhysd/actionlint`, `zizmor-pre-commit`), each with a `rev`.
  Duplicates the versions `mise.toml` already pins, giving each tool two pins.
- **`.pre-commit-config.yaml`: `pre-commit/pre-commit-hooks` by `rev`, plus `local` hooks
  that call the mise-pinned tools.**

## Reasoning

The YAML file works under both prek and pre-commit, and prek runs `pre-commit-hooks` through its
Rust fast path. The only new pin is that `rev`, owned by Renovate's `pre-commit` manager (off by
default, so enabled in `renovate.json5`). actionlint, zizmor and prek stay pinned once, in
`mise.toml`. CI's `prek hooks` job installs them with `jdx/mise-action`, so a commit made
without `prek install` is still checked.

## Trade-offs accepted

- Local hooks need the tools on `PATH` (`mise install`), not installed by the hook manager.
- CI skips gofmt in the hooks job; the `go` job checks it with go.mod's Go version.
- `jdx/mise-action` v5.1.2 was 3 days old when pinned, inside the 7-day cooldown; Robin
  accepted it.

## Supersedes

None.
