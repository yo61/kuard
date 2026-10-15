## Decision: Pin Dockerfile.ci's node and golang images to exact version tags plus digest

## Context

`Dockerfile.ci` pinned `node:24-alpine@sha256:…` and `golang:1.27-alpine@sha256:…`. A floating
tag's digest update has no version, so Renovate titled #57 "update node.js to 44fce69" and could
not tell a Node release from an Alpine rebuild. #57 moved the image to Node 24.22.0, released
the day before, while CI's `NODE_VERSION` stayed at 24.21.0: the `node` group never paired them,
and `minimumReleaseAge` did not hold a digest update back.

## Alternatives considered

- **Keep floating tags.** Digests still make builds reproducible, but versions drift from CI's
  pins without a PR saying so.
- **Exact tag, pinned to the newer 24.22.0, bumping CI to match.** Consistent, but takes a
  one-day-old release past the repo's 7-day release-age policy.
- **Exact tag at CI's version (24.21.0) and go.mod's (1.27.2).**

## Reasoning

With exact tags, a new Node or Go release is a version update that Renovate groups with the CI
pin it belongs with and holds for `minimumReleaseAge`. An Alpine rebuild of the same version is
still a digest update. Going back to 24.21.0 keeps CI and the image on one version and leaves
24.22.0 to arrive through Renovate once it is 7 days old.

## Trade-offs accepted

- The image builds with Node 24.21.0 again until Renovate proposes 24.22.0.
- `alpine:3.24` stays a minor-version tag: nothing else pins Alpine, so there is nothing for it
  to drift from.

## Supersedes

None.
