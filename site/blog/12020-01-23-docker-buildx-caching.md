---
description: finally looking at buildx
title: docker buildx caching
---

### _buildkit_

docker [BuildKit](https://docs.docker.com/develop/develop-images/build_enhancements/)
is getting some interesting new features
setting it apart from the other container building tools.
Unfortunately, it means a new cli subcommand that's not entirely backwards compatible,
say hello to [`docker buildx`](https://github.com/docker/buildx)

#### _tldr_

##### _stateless_ build workers

You get a fresh build environment every time, use registry caching.

```Dockerfile
FROM golang:rc-alpine AS build
WORKDIR /workspace
COPY go.mod go.sum .
RUN go mod download
COPY . .
RUN go build -o app .

FROM scratch
COPY --from=build /workspace/app /app
ENTRYPOINT ["/app"]
```

Probably need to spend some time cleaning up the registry every once in a while.
You could make life easier by using a separate registry for caching
and nuking it every 2 weeks...

```sh
docker buildx \
  --cache-from type=registry,ref=your.registry/image \
  --cache-to   type=registry,ref=your.registry/image,mode=max \
  --tag your.registry/image:tag \
  --push .
```

###### _stateless_ with ci cache

like the above, but use an external system to save/restore the cache

```sh
# use some tool to restore a cache from previous builds
cache-restore /tmp/docker-cache
docker buildx \
  --cache-from type=local,ref=/tmp/docker-cache \
  --cache-to   type=local,ref=/tmp/docker-cache,mode=max \
  --tag your.registry/image:tag \
  --push .
# use some tool to save the cache for future use
cache-save /tmp/docker-cache
```

##### _statefull_ build workers

Your workers have a chance to reuse the local cache between builds,
and even across builds for different apps.

```Dockerfile
FROM golang:rc-alpine AS build
WORKDIR /workspace
COPY . .
RUN --mount=type=cache,id=gomod,dst=/go/pkg/mod \
    --mount=type=cache,id=gobuild,dst=/root/.cache/go-build \
    go build -o app .

FROM scratch
COPY --from=build /workspace/app /app
ENTRYPOINT ["/app"]
```

You could still use `--cache-to/from` if your image is more complex
and you'd like to reuse layers

```sh
docker buildx \
  --tag your.registry/image:tag \
  --push .
```

#### _caching_

The fun of trying to munge together the:

- language dependency cache
- language build cache
- docker layer cache/reuse
- ci system cache

##### _dependencies_

If you use any sane language package manager,
somewhere on disk will be a global cache for your dependencies.
This should be reproducible, ie given the same package list,
the resulting cache should always be the same.

If you use java/maven, despair as your developers don't properly declare all their deps
and dynamically add them at test time.

##### _build_

Some tools like [go](https://golang.org/cmd/go/#hdr-Build_and_test_caching)
have a global build cache with proper keys.
Others just use the local directory
and some time checks to guess if something is out of date.
The first could probably be shared across machines, the second probably can't.

##### _docker_

docker caches by layer, keyed by the previous layer and the context copied in so far.
Since dependencies change relatively infrequently compared to source code,
a good strategy is to copy the dependency manifest and download them,
allowing these layers to be cached.
However, since layers are cached as a whole,
any change in dependency invalidates the entire layer.
Also, since the dependency manifest is part of the cache key,
you're unlikely to be able to share layer caches across applications
(unless you purposefully construct a dedicated fat caching layer).

`--cache-from` and `--cache-to` support several outputs,
but the interesting ones are to a local directory
(so it can be cached by your ci system)
and to a remote registry.
The important flag is `mode=max` for `--cache-to` to include all layers in multistage builds.
Which is faster and whether you want to fill a registry with cache images
is up to you, though not all registries support cache manifests,
ex: [gcr](https://github.com/moby/buildkit/issues/1143)

The other fun feature is `RUN --mount=type=cache,dst=/some/path your command`.
This creates a directory that's excluded from the image
and can be reused across invocations, making it a good place to put the
dependency/build cache, assuming your machines build multiple images.
Being able to control where this is on the host
(for ci systems to save/restore)
is still [an open issue](https://github.com/moby/buildkit/issues/1512)

##### _ci_

There are 2 main flavours of CI: stateful and stateless.
Stateful systems persist between builds, giving you the ability to share state
(eg cache directories) between builds on the same machine.
Stateless systems give you a clean environment every time,
good for ensuring no external factors affect your build,
bad if you trust your build/dependency cache.
Either way, they will sometimes provide you with tools to save/restore
from persistent caches shared by all workers.
