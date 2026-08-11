# Demo application for "Kubernetes Up and Running"

![screenshot](docs/images/screenshot.png)

### Running

```
kubectl run --restart=Never --image=ghcr.io/yo61/kuard:blue kuard
kubectl port-forward kuard 8080:8080
```

Open your browser to [http://localhost:8080](http://localhost:8080).

The book and older tutorials say `gcr.io/kuar-demo/kuard-amd64:blue`. Those images are gone --
Google shut down Container Registry and the project was never migrated -- which is why this fork
publishes its own. See [FORK.md](FORK.md).

### Building

Build logic lives in **`Dockerfile.ci`**: a multi-stage build that compiles the React client and
the Go binary and emits a minimal Alpine image. `.github/workflows/build-kuard.yml` runs it to
publish the images described in [FORK.md](FORK.md).

Locally, [Task](https://taskfile.dev) wraps the workflow:

```
task              # list every task
task build        # bundle, embedded assets and binary into bin/
task run          # build and run on :8080
task test         # Go tests and the client bundle test
task image        # build the container image the way CI does
```

Note that `pkg/sitedata/bindata.go` is generated and not committed, so `go build` and `go test` do
not work on a fresh clone until `task generate` has run. Every Go task depends on it, so using
Task means you do not have to remember that.

#### Multi-stage Dockerfile

A feature of Docker added since the book was published is the "multi-stage" build: one image
compiles the software, and a second, smaller one runs it.

`Dockerfile.ci` is an example. It builds the React client in a `node` stage and the Go binary in a
`golang` stage, then copies only the resulting binary into a bare `alpine` image.

```
docker build -f Dockerfile.ci -t kuard:blue .
docker run --rm -ti --name kuard --publish 8080:8080 kuard:blue
```

`task image` and `task image:run` do the same thing.

#### Inserting a prebuilt binary

Closer to what the book describes: build the binary outside Docker and insert it. `task build`
puts one in `bin/kuard`, so a minimal image is just:

```
FROM alpine
COPY bin/kuard /kuard
ENTRYPOINT [ "/kuard" ]
```

To upload to a registry you will have to tag it and push. Refer to your registry's documentation
for details.

### KeyGen Workload

To help simulate batch workers, we have a synthetic workload of generating 4096 bit RSA keys.  This can be configured through the UI or the command line.

```
--keygen-enable               Enable KeyGen workload
--keygen-exit-code int        Exit code when workload complete
--keygen-exit-on-complete     Exit after workload is complete
--keygen-memq-queue string    The MemQ server queue to use. If MemQ is used, other limits are ignored.
--keygen-memq-server string   The MemQ server to draw work items from.  If MemQ is used, other limits are ignored.
--keygen-num-to-gen int       The number of keys to generate. Set to 0 for infinite
--keygen-time-to-run int      The target run time in seconds. Set to 0 for infinite
```

### MemQ server

We also have a simple in memory queue with REST API.  This is based heavily on https://github.com/kelseyhightower/memq.

The API is as follows with URLs being relative to `<server addr>/memq/server`.  See `pkg/memq/types.go` for the data structures returned.

| Method | Url | Desc
| --- | --- | ---
| `GET` | `/stats` | Get stats on all queues
| `PUT` | `/queues/:queue` | Create a queue
| `DELETE` | `/queues/:queue` | Delete a queue
| `POST` | `/queues/:queue/drain` | Discard all items in queue
| `POST` | `/queues/:queue/enqueue` | Add item to queue.  Body is plain text. Response is message object.
| `POST` | `/queues/:queue/dequeue` | Grab an item off the queue and return it. Returns a 204 "No Content" if queue is empty.

### Versions

Images built will automatically have the git version (based on tag) applied.  In addition, there is an idea of a "fake version".  This is used so that we can use the same basic server to demonstrate upgrade scenarios.

Originally (and in the Kubernetes Up & Running book) we had `1`, `2`, and `3`.  This confused people so going forward we will be using colors instead: `blue`, `green` and `purple`. This translates into the following container images:

```
gcr.io/kuar-demo/kuard-amd64:v0.9-blue
gcr.io/kuar-demo/kuard-amd64:blue
gcr.io/kuar-demo/kuard-amd64:v0.9-green
gcr.io/kuar-demo/kuard-amd64:green
gcr.io/kuar-demo/kuard-amd64:v0.9-purple
gcr.io/kuar-demo/kuard-amd64:purple
```

For documentation where you want to demonstrate using versions but use the latest version of this server, you can simply reference `gcr.io/kuar-demo/kuard-amd64:blue`.  You can then demonstrate an upgrade with `gcr.io/kuar-demo/kuard-amd64:green`.

(Another way to think about it is that `:blue` is essentially `:latest-blue`)

We also build versions for `arm`, `arm64`, and `ppc64le`.  Just substitute the appropriate architecture in the image name.  These aren't as well tested as the `amd64` version but seem to work okay.

### Development

If you just want to do Go server development, you can build the client as part of a build `make`.  It'll drop the result in to `sitedata/built/`.

If you want to do both Go server and React.js client dev, you need to do the following:

1. Have Node installed
2. In one terminal

  * `cd client`
  * `npm install`
  * `npm run start`
  * This will start a debug node server on `localhost:8081`.  It'll proxy all unhandled requests to `localhost:8080`

3. In another terminal
  * Ensure that $GOPATH is set to the directory with your go source code and binaries + ensure that $GOPATH is part of $PATH.
  * `go get -u github.com/jteeuwen/go-bindata/...`
  * `go generate ./pkg/...`
  * `GO111MODULE=on go run cmd/kuard/*.go --debug`
4. Open your browser to http://localhost:8081.

This should support live reload of any changes to the client.  The Go server will need to be exited and restarted to see changes.

### TODO
* [ ] Make file system browser better.  Show size, permissions, etc.  Might be able to do this by faking out an `index.html` as part of the http.FileSystem stuff.
* [ ] Clean up form for keygen workload.  It is too big and the form build doesn't have enough flexibility to really shrink it down.
* [ ] Get rid of go-bindata as it is abandoned.
