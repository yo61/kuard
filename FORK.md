# yo61/kuard — fork notes

A fork of [`kubernetes-up-and-running/kuard`](https://github.com/kubernetes-up-and-running/kuard), the demo
app used throughout *Kubernetes: Up and Running*.

## Why this fork exists

The book (and every tutorial) references images at `gcr.io/kuar-demo/…`, e.g.
`gcr.io/kuar-demo/kuard-amd64:blue`. **Those images are gone** — Google shut down Container Registry
(`gcr.io`) in favour of Artifact Registry, and the `kuar-demo` project was never migrated, so pulls now fail:

```
Failed to pull image "gcr.io/kuar-demo/kuard-amd64:blue": ... 403 Forbidden
# or: Container Registry is deprecated and shutting down ...
```

The upstream maintainers confirm there's no replacement image
([kuard#60](https://github.com/kubernetes-up-and-running/kuard/issues/60)). So this fork builds kuard itself
and publishes it to **GitHub Container Registry** under `ghcr.io/yo61/kuard`.

## Images

A GitHub Action (`.github/workflows/build-kuard.yml`) builds and pushes the three "fake versions" the book
uses for rollout demos:

| Image | equivalent to the book's |
|---|---|
| `ghcr.io/yo61/kuard:blue`   | `gcr.io/kuar-demo/kuard-amd64:blue`   |
| `ghcr.io/yo61/kuard:green`  | `gcr.io/kuar-demo/kuard-amd64:green`  |
| `ghcr.io/yo61/kuard:purple` | `gcr.io/kuar-demo/kuard-amd64:purple` |

Each tag is a multi-arch manifest covering `linux/amd64`, `linux/arm64` and `linux/arm/v7`, so the same
image runs on cloud x86, 64-bit ARM (Graviton, Pi 4/5) and 32-bit ARM (Pi 2/3). The book's image was
amd64-only (`kuard-amd64`); pulling `ghcr.io/yo61/kuard:blue` gets you the right architecture automatically.

### Using it

Wherever the book says `gcr.io/kuar-demo/kuard-amd64:blue`, use `ghcr.io/yo61/kuard:blue`:

```bash
kubectl run kuard --image=ghcr.io/yo61/kuard:blue
```

The GHCR package must be **public** for the cluster to pull it without credentials — set it in the package's
settings after the first push (Packages → kuard → Package settings → Change visibility → Public).

## Building

- **`build-kuard.yml`** publishes: it runs on push to `main` and on manual dispatch (Actions → build-kuard →
  Run workflow), using the repo's built-in `GITHUB_TOKEN`, so there is no PAT or secret to configure.
- **`ci.yml`** validates pull requests: it builds the same `Dockerfile.ci` image with pushing disabled, and
  runs `npm ci` to catch a `client/package-lock.json` that has drifted out of sync with `package.json`.
  It never authenticates to GHCR.
- Build logic is in **`Dockerfile.ci`** (self-contained; the upstream `Makefile` targets the dead `gcr.io`
  and needs docker-in-docker, so it's ignored).
- The upstream 2019 build has bit-rotted; `Dockerfile.ci` works around it: the deleted `jteeuwen/go-bindata`
  dependency is swapped for the `go-bindata/go-bindata` community fork, and `golang.org/x/net` is bumped so
  the code links on current Go (the 2019 version used a `//go:linkname` into `syscall` that Go 1.26 rejects).
  The React client and Go binary are built in current `node`/`golang` stages, so nothing has to `apk add`
  from long-archived Alpine repos.

## Divergence from upstream

The source layout, the app's behaviour and its licensing are unchanged from upstream. The dependencies are
not: upstream's 2019 stack no longer builds or installs cleanly, so it has been brought forward.

- **Client build**: webpack 4 → 5, webpack-dev-server 3 → 6, Babel 7 → 8, on node 24. This removed the
  `NODE_OPTIONS=--openssl-legacy-provider` workaround webpack 4 needed under OpenSSL 3.
- **Client app**: React 15 → 19. `react-router-component` (abandoned) → `react-router-dom`,
  `react-jsonschema-form` (deprecated) → `@rjsf/core`, and `marked` 0.3 → 18.
- **Go**: `viper`, `client_golang`, `miekg/dns` and friends brought to current. The metrics endpoint moved
  from `prometheus.Handler()` to `promhttp.Handler()`, which client_golang removed in 1.0.

Every page still renders and behaves as the book describes; the UI is visually unchanged.
