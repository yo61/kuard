## Decision: prek hooks use only builtin and local hooks; tool versions come from mise.toml

## Context

The repo had no git hooks, so gofmt, actionlint and zizmor ran only when someone remembered
`task fmt:check` or `task lint:workflows`. CI runs gofmt but neither workflow linter.

## Alternatives considered

- **Remote hook repos** (`pre-commit-hooks`, `rhysd/actionlint`, `zizmor-pre-commit`) pinned by
  `rev`. Each `rev` is a new pin that needs an owner: Renovate's `pre-commit` manager would have
  to be enabled, and the versions would duplicate the ones in `mise.toml`.
- **`prek auto-update --cooldown-days 7`** to bump those revs. A second updater outside
  Renovate's queue and its `minimumReleaseAge`.
- **prek's `builtin` hooks plus `local` hooks that call the mise-pinned tools.**

## Reasoning

`prek.toml` pins nothing. actionlint and zizmor run at the version `mise.toml` pins, the same one
`task lint:workflows` uses, and Renovate's mise manager already owns those. prek itself is
pinned in `mise.toml` too.

## Trade-offs accepted

- The hooks need the tools on `PATH` (`mise install`) instead of prek installing them.
- The `builtin` hooks change only when prek is updated, not separately.
- CI does not run prek, so the hooks catch nothing for someone who hasn't run `prek install`.

## Supersedes

None.
