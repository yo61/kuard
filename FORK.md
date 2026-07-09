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

Built for `linux/amd64` (the book's image was `kuard-amd64`).

### Using it

Wherever the book says `gcr.io/kuar-demo/kuard-amd64:blue`, use `ghcr.io/yo61/kuard:blue`:

```bash
kubectl run kuard --image=ghcr.io/yo61/kuard:blue
```

The GHCR package must be **public** for the cluster to pull it without credentials — set it in the package's
settings after the first push (Packages → kuard → Package settings → Change visibility → Public).

## Building

- The Action runs on push to `master` and on manual dispatch (Actions → build-kuard → Run workflow). It uses
  the repo's built-in `GITHUB_TOKEN`, so no PAT or secrets to configure.
- Build logic is in **`Dockerfile.ci`** (self-contained; the upstream `Makefile` targets the dead `gcr.io`
  and needs docker-in-docker, so it's ignored).
- The upstream 2019 build has bit-rotted; `Dockerfile.ci` works around it: the deleted `jteeuwen/go-bindata`
  dependency is swapped for the `go-bindata/go-bindata` community fork, and `golang.org/x/net` is bumped so
  the code links on current Go (the 2019 version used a `//go:linkname` into `syscall` that Go 1.26 rejects).
  The React client and Go binary are built in current `node`/`golang` stages, so nothing has to `apk add`
  from long-archived Alpine repos.

Everything else (the app, its source layout, licensing) is unchanged from upstream.
